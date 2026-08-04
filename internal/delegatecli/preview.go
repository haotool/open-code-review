package delegatecli

import (
	"context"
	"fmt"
	"io"

	"github.com/alibaba/open-code-review/internal/diff"
	"github.com/alibaba/open-code-review/internal/model"
	"github.com/alibaba/open-code-review/internal/preview"
)

// RunPreview executes the preview subcommand and writes formatted output to w.
func RunPreview(w io.Writer, args []string) error {
	opts, _, err := ParseFlags(args)
	if err != nil {
		return err
	}
	if opts.ShowHelp {
		fmt.Fprintln(w, "Usage: ocr-delegate preview [flags]")
		fmt.Fprintln(w, "\nOutputs reviewable file list with mode/ref metadata for the host agent to construct git commands.")
		return nil
	}

	dc, err := LoadContext(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()
	pv, err := preview.Run(ctx, preview.Args{
		RepoDir:    dc.RepoDir,
		From:       dc.Options.From,
		To:         dc.Options.To,
		Commit:     dc.Options.Commit,
		FileFilter: dc.FileFilter,
		GitRunner:  dc.GitRunner,
	})
	if err != nil {
		return fmt.Errorf("preview failed: %w", err)
	}

	return FormatPreview(w, dc, pv)
}

// FormatPreview renders preview output in delegate markdown format.
func FormatPreview(w io.Writer, dc *Context, pv *model.Preview) error {
	ctx := context.Background()

	fmt.Fprintf(w, "# Files (%d reviewable / %d total)\n\n", pv.ReviewableCount, pv.TotalFiles)
	fmt.Fprintf(w, "- mode: %s\n", dc.ReviewMode())
	if dc.Options.From != "" {
		fmt.Fprintf(w, "- from: %s\n", dc.Options.From)
	}
	if dc.Options.To != "" {
		fmt.Fprintf(w, "- to: %s\n", dc.Options.To)
	}
	if dc.Options.Commit != "" {
		fmt.Fprintf(w, "- commit: %s\n", dc.Options.Commit)
	}
	if mergeBase := dc.mergeBase(ctx); mergeBase != "" {
		fmt.Fprintf(w, "- merge_base: %s\n", mergeBase)
	}
	if dc.Options.Background != "" {
		fmt.Fprintf(w, "- background: %s\n", dc.Options.Background)
	}
	fmt.Fprintf(w, "- total_insertions: %d\n", pv.TotalInsertions)
	fmt.Fprintf(w, "- total_deletions: %d\n\n", pv.TotalDeletions)

	for _, entry := range pv.Entries {
		marker := "  "
		if !entry.WillReview {
			marker = "~~"
		}
		fmt.Fprintf(w, "%s- `%s` [%s] +%d/-%d", marker, entry.Path, entry.Status, entry.Insertions, entry.Deletions)
		if !entry.WillReview {
			fmt.Fprintf(w, " (excluded: %s)", entry.ExcludeReason)
			fmt.Fprint(w, "~~")
		}
		fmt.Fprintln(w)
	}

	return nil
}

func (dc *Context) mergeBase(ctx context.Context) string {
	if dc.Options.From == "" || dc.Options.To == "" {
		return ""
	}
	provider := diff.NewProvider(dc.RepoDir, dc.Options.From, dc.Options.To, dc.GitRunner)
	return provider.MergeBase(ctx)
}
