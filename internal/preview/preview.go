// Package preview implements the LLM-free file-selection preview shared by
// `ocr review --preview` and delegation mode: it loads diffs (workspace,
// range, or commit mode) and applies the exclusion filters without
// dispatching any LLM calls. Extracted from internal/agent so delegate-only
// builds can link it without pulling in llm/llmloop/session/telemetry.
package preview

import (
	"context"
	"fmt"
	"strings"

	allowedext "github.com/alibaba/open-code-review/internal/config/allowlist"
	"github.com/alibaba/open-code-review/internal/config/rules"
	"github.com/alibaba/open-code-review/internal/diff"
	"github.com/alibaba/open-code-review/internal/gitcmd"
	"github.com/alibaba/open-code-review/internal/model"
)

// Args holds the inputs needed to compute a preview. It mirrors the
// diff-selection subset of agent.Args.
type Args struct {
	// RepoDir is the root of the git repository.
	RepoDir string

	// From and To define the diff range (e.g., "main..feature-branch").
	From string
	To   string

	// Commit is a single commit hash to review (vs its parent).
	Commit string

	// FileFilter holds user-configured include/exclude path patterns from rule.json.
	// When nil, only the default extension and path filters apply.
	FileFilter *rules.FileFilter

	// GitRunner limits the total number of concurrent git subprocesses.
	// When nil, subprocesses are spawned without a global limit.
	GitRunner *gitcmd.Runner
}

// Run loads diffs and applies the filter algorithm, returning structured
// preview data without dispatching any LLM calls.
func Run(ctx context.Context, args Args) (*model.Preview, error) {
	var provider *diff.Provider

	switch {
	case args.Commit != "":
		provider = diff.NewCommitProvider(args.RepoDir, args.Commit, args.GitRunner)
	case args.From != "" && args.To != "":
		provider = diff.NewProvider(args.RepoDir, args.From, args.To, args.GitRunner)
	default:
		provider = diff.NewWorkspaceProvider(args.RepoDir, args.GitRunner)
	}

	parsed, err := provider.GetDiff(ctx)
	if err != nil {
		// 保留 agent.Preview 既有的「load diffs: get diffs:」前綴，
		// 使 CLI 錯誤輸出在抽取前後逐字一致。
		return nil, fmt.Errorf("load diffs: get diffs: %w", err)
	}

	result := &model.Preview{
		TotalFiles: len(parsed),
	}
	for i := range parsed {
		result.TotalInsertions += parsed[i].Insertions
		result.TotalDeletions += parsed[i].Deletions
	}

	for _, d := range parsed {
		entry := model.PreviewEntry{
			Path:       EffectivePath(d),
			Insertions: d.Insertions,
			Deletions:  d.Deletions,
			Status:     DiffStatus(d),
		}

		reason := WhyExcluded(args.FileFilter, d)
		if reason == model.ExcludeNone && d.IsDeleted {
			reason = model.ExcludeDeleted
		}

		entry.WillReview = reason == model.ExcludeNone
		entry.ExcludeReason = reason

		if entry.WillReview {
			result.ReviewableCount++
		} else {
			result.ExcludedCount++
		}

		result.Entries = append(result.Entries, entry)
	}

	return result, nil
}

// WhyExcluded applies the file-selection filter algorithm and returns the
// specific reason a file is excluded, or model.ExcludeNone when it will be
// reviewed. f may be nil, in which case only the default extension and path
// filters apply.
func WhyExcluded(f *rules.FileFilter, d model.Diff) model.ExcludeReason {
	if d.IsBinary {
		return model.ExcludeBinary
	}

	path := EffectivePath(d)

	if f != nil && f.IsUserExcluded(path) {
		return model.ExcludeUserRule
	}

	if f != nil && f.HasInclude() && f.IsUserIncluded(path) {
		return model.ExcludeNone
	}

	ext := ExtFromPath(path)
	if ext != "" && !allowedext.IsAllowedExt(ext) {
		return model.ExcludeExtension
	}

	if allowedext.IsExcludedPath(path) {
		return model.ExcludeDefaultPath
	}

	return model.ExcludeNone
}

// EffectivePath returns the user-facing path of a diff: NewPath, unless the
// file was deleted (NewPath == /dev/null), in which case OldPath.
func EffectivePath(d model.Diff) string {
	if d.NewPath == "/dev/null" {
		return d.OldPath
	}
	return d.NewPath
}

// DiffStatus returns the display status string for a diff entry.
func DiffStatus(d model.Diff) string {
	switch {
	case d.IsBinary:
		return "binary"
	case d.IsNew:
		return "added"
	case d.IsDeleted:
		return "deleted"
	case d.IsRenamed:
		return "renamed"
	case d.OldPath != d.NewPath && d.OldPath != "" && d.OldPath != "/dev/null":
		return "renamed"
	default:
		return "modified"
	}
}

// ExtFromPath returns the file extension with leading dot, lowercased.
// Files without an extension (including dotfiles) return "".
func ExtFromPath(path string) string {
	basename := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		basename = path[idx+1:]
	}
	dot := strings.LastIndex(basename, ".")
	if dot <= 0 {
		return ""
	}
	return strings.ToLower(basename[dot:])
}
