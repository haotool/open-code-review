package delegatecli

import (
	"strings"

	"github.com/alibaba/open-code-review/internal/config/rules"
)

// SplitPaths splits a comma-separated path list, trimming whitespace.
func SplitPaths(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ApplyCLIExcludes appends user-supplied --exclude patterns onto ctx.FileFilter.
func ApplyCLIExcludes(ctx *Context, patterns []string) {
	if ctx == nil {
		return
	}
	ApplyExcludesToFilter(&ctx.FileFilter, patterns)
}

// ApplyExcludesToFilter appends exclude patterns to *filter, creating a filter if nil.
func ApplyExcludesToFilter(filter **rules.FileFilter, patterns []string) {
	if len(patterns) == 0 {
		return
	}
	if *filter == nil {
		*filter = &rules.FileFilter{}
	}
	(*filter).Exclude = append((*filter).Exclude, patterns...)
}

func applyCLIExcludesOnContext(cc *commonContext, patterns []string) {
	ApplyExcludesToFilter(&cc.FileFilter, patterns)
}
