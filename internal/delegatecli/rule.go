package delegatecli

import (
	"fmt"
	"io"

	"github.com/alibaba/open-code-review/internal/delegate"
)

// RunRule executes the rule subcommand and writes formatted output to w.
func RunRule(w io.Writer, args []string) error {
	opts, remaining, err := ParseFlags(args)
	if err != nil {
		return err
	}
	if opts.ShowHelp {
		fmt.Fprintln(w, "Usage: ocr-delegate rule [flags] <path...>")
		fmt.Fprintln(w, "\nOutputs resolved review rules grouped by content. Accepts multiple paths.")
		return nil
	}
	if len(remaining) == 0 {
		return fmt.Errorf("at least one file path is required\nUsage: ocr-delegate rule [flags] <path...>")
	}

	dc, err := LoadContext(opts)
	if err != nil {
		return err
	}

	groups := delegate.GroupRules(dc.Resolver, remaining)
	fmt.Fprint(w, delegate.RuleGroupsMarkdown(groups))
	return nil
}
