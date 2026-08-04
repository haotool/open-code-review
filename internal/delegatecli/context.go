package delegatecli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/open-code-review/internal/config/rules"
	"github.com/alibaba/open-code-review/internal/config/template"
	"github.com/alibaba/open-code-review/internal/gitcmd"
)

// Context holds shared delegate state after LoadContext.
type Context struct {
	RepoDir    string
	Resolver   rules.Resolver
	FileFilter *rules.FileFilter
	GitRunner  *gitcmd.Runner
	Options    Options
}

// LoadContext validates flags, loads rules, and prepares git/background state.
func LoadContext(opts Options) (*Context, error) {
	cc, err := loadCommonContext(opts.RepoDir, opts.RulePath, opts.MaxGitProcs)
	if err != nil {
		return nil, err
	}
	applyCLIExcludesOnContext(cc, SplitPaths(opts.Excludes))

	if err := ValidateReviewRefs(cc.RepoDir, opts); err != nil {
		return nil, err
	}

	if opts.BackgroundFile != "" {
		bgPath := ResolveBackgroundFilePath(cc.RepoDir, opts.BackgroundFile)
		fileBackground, err := LoadBackgroundFile(bgPath)
		if err != nil {
			return nil, err
		}
		opts.Background, err = MergeBackground(opts.Background, fileBackground)
		if err != nil {
			return nil, err
		}
	} else if opts.Background != "" {
		opts.Background, err = SanitizeBackground(opts.Background)
		if err != nil {
			return nil, fmt.Errorf("inline background: %w", err)
		}
	}

	if opts.Commit != "" && opts.Background == "" {
		if msg, err := getCommitMessage(cc.RepoDir, opts.Commit); err == nil && msg != "" {
			opts.Background, err = SanitizeBackground(msg)
			if err != nil {
				return nil, fmt.Errorf("commit message background: %w", err)
			}
		}
	}

	return &Context{
		RepoDir:    cc.RepoDir,
		Resolver:   cc.Resolver,
		FileFilter: cc.FileFilter,
		GitRunner:  cc.GitRunner,
		Options:    opts,
	}, nil
}

type commonContext struct {
	RepoDir    string
	Resolver   rules.Resolver
	FileFilter *rules.FileFilter
	GitRunner  *gitcmd.Runner
}

func loadCommonContext(repoDirInput, rulePath string, maxGitProcs int) (*commonContext, error) {
	tpl, err := template.LoadDefault()
	if err != nil {
		return nil, fmt.Errorf("load default template: %w", err)
	}
	if err := tpl.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	repoDir, err := resolveWorkingDir(repoDirInput)
	if err != nil {
		return nil, err
	}

	resolver, fileFilter, err := rules.NewResolver(repoDir, rulePath)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}

	return &commonContext{
		RepoDir:    repoDir,
		Resolver:   resolver,
		FileFilter: fileFilter,
		GitRunner:  gitcmd.New(maxGitProcs),
	}, nil
}

func resolveWorkingDir(input string) (string, error) {
	if input == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		input = wd
	}
	absPath, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}
	if _, statErr := os.Stat(absPath); statErr != nil {
		return "", fmt.Errorf("stat %s: %w", absPath, statErr)
	}
	out, err := runGitCmd(absPath, "rev-parse", "--git-dir")
	isGit := err == nil && len(out) > 0
	if !isGit {
		return "", fmt.Errorf("%s is not a git repository", absPath)
	}
	top, topErr := runGitCmdStdout(absPath, "rev-parse", "--show-toplevel")
	t := strings.TrimSpace(string(top))
	if topErr != nil || t == "" {
		return "", fmt.Errorf("%s is a git repository without a work tree (bare repo?); cannot resolve its top level for review", absPath)
	}
	return t, nil
}

// ReviewMode returns the string mode identifier for the current options.
func (dc *Context) ReviewMode() string {
	switch {
	case dc.Options.Commit != "":
		return "commit"
	case dc.Options.From != "" && dc.Options.To != "":
		return "range"
	default:
		return "workspace"
	}
}
