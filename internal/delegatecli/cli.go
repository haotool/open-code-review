package delegatecli

import (
	"fmt"
	"io"
	"os"
)

// Run is the top-level CLI entry point for ocr-delegate.
func Run(args []string) error {
	return RunWithIO(os.Stdout, os.Stderr, args)
}

// RunWithIO runs the CLI writing stdout to out and errors to errOut (for tests).
func RunWithIO(out, errOut io.Writer, args []string) error {
	if len(args) == 0 {
		PrintUsage(out)
		return nil
	}

	sub := args[0]
	switch sub {
	case "-h", "--help":
		PrintUsage(out)
		return nil
	case "preview":
		return RunPreview(out, args[1:])
	case "rule":
		return RunRule(out, args[1:])
	default:
		fmt.Fprintf(errOut, "unknown sub-command: %s\nRun 'ocr-delegate -h' for usage\n", sub)
		return fmt.Errorf("unknown sub-command: %s", sub)
	}
}

// PrintUsage prints top-level ocr-delegate usage to w.
func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, `OpenCodeReview - Delegation Mode

Usage:
  ocr-delegate <sub-command> [flags]

Sub-commands:
  preview       Preview reviewable files with mode/ref metadata
  rule          Output resolved review rules grouped by content

Shared Flags:
  --from string           source ref to start diff from (e.g., 'main')
  --to string             target ref to end diff at (e.g., 'feature-branch')
  -c, --commit string     single commit hash or tag to review
  --repo string           root directory of the git repository (default: current dir)
  --rule string           path to JSON file with system review rules
  --exclude string        comma-separated gitignore-style patterns to exclude
  -b, --background string optional requirement/business context
  -B, --background-file   path to a Markdown file used as background
  --max-git-procs int     max concurrent git subprocesses (default 16)

Examples:
  # Preview which files will be reviewed
  ocr-delegate preview --from main --to feature

  # Preview workspace changes
  ocr-delegate preview

  # Get rules for multiple files (grouped by content)
  ocr-delegate rule internal/agent/agent.go internal/llm/client.go`)
}
