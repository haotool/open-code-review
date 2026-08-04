package agent

import (
	"context"

	"github.com/alibaba/open-code-review/internal/model"
	"github.com/alibaba/open-code-review/internal/preview"
)

// ExcludeReason / DiffPreview / DiffPreviewEntry are now type aliases of
// the mode-agnostic preview types in internal/model. Kept for backwards
// compatibility with existing call sites; internal/scan returns the same
// model.Preview shape directly.
type ExcludeReason = model.ExcludeReason
type DiffPreview = model.Preview
type DiffPreviewEntry = model.PreviewEntry

// Re-export the constants so callers can keep writing agent.ExcludeBinary.
const (
	ExcludeNone        = model.ExcludeNone
	ExcludeUserRule    = model.ExcludeUserRule
	ExcludeExtension   = model.ExcludeExtension
	ExcludeDefaultPath = model.ExcludeDefaultPath
	ExcludeDeleted     = model.ExcludeDeleted
	ExcludeBinary      = model.ExcludeBinary
)

// whyExcluded applies the filter algorithm as shouldReview but returns the
// specific reason a file is excluded. The logic lives in internal/preview
// so delegate-only builds can share it without linking this package.
func (a *Agent) whyExcluded(d model.Diff) ExcludeReason {
	return preview.WhyExcluded(a.args.FileFilter, d)
}

// Preview loads diffs and applies the filter algorithm, returning structured
// preview data without dispatching any LLM calls.
func (a *Agent) Preview(ctx context.Context) (*DiffPreview, error) {
	return preview.Run(ctx, preview.Args{
		RepoDir:    a.args.RepoDir,
		From:       a.args.From,
		To:         a.args.To,
		Commit:     a.args.Commit,
		FileFilter: a.args.FileFilter,
		GitRunner:  a.args.GitRunner,
	})
}

func effectivePath(d model.Diff) string {
	return preview.EffectivePath(d)
}

func diffStatus(d model.Diff) string {
	return preview.DiffStatus(d)
}
