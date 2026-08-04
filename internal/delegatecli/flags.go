package delegatecli

import (
	"flag"
	"fmt"
)

// Options holds parsed delegate subcommand flags shared by preview and rule.
type Options struct {
	RepoDir        string
	From           string
	To             string
	Commit         string
	Excludes       string
	RulePath       string
	Background     string
	BackgroundFile string
	MaxGitProcs    int
	ShowHelp       bool
}

type ocrFlagSet struct {
	fs       *flag.FlagSet
	shortMap map[string]string
	showHelp bool
}

func newOcrFlagSet(name string) *ocrFlagSet {
	return &ocrFlagSet{
		fs:       flag.NewFlagSet(name, flag.ContinueOnError),
		shortMap: make(map[string]string),
	}
}

func (a *ocrFlagSet) StringVarP(p *string, name, shorthand, value, usage string) {
	suffix := ""
	if shorthand != "" {
		a.shortMap[shorthand] = name
		suffix = fmt.Sprintf(" (shorthand: -%s)", shorthand)
	}
	a.fs.StringVar(p, name, value, usage+suffix)
}

func (a *ocrFlagSet) StringVar(p *string, name, value, usage string) {
	a.fs.StringVar(p, name, value, usage)
}

func (a *ocrFlagSet) IntVar(p *int, name string, value int, usage string) {
	a.fs.IntVar(p, name, value, usage)
}

func (a *ocrFlagSet) Parse(arguments []string) error {
	expanded := expandShortFlags(arguments, a.shortMap)

	for _, arg := range expanded {
		if arg == "-h" || arg == "--help" {
			a.showHelp = true
			return nil
		}
	}

	return a.fs.Parse(expanded)
}

func expandShortFlags(args []string, shortMap map[string]string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if len(arg) == 2 && arg[0] == '-' && arg[1] != '-' {
			key := string(arg[1])
			if full, ok := shortMap[key]; ok {
				out = append(out, "--"+full)
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}

// ParseFlags parses shared delegate flags and returns remaining positional args.
func ParseFlags(args []string) (Options, []string, error) {
	a := newOcrFlagSet("ocr-delegate")

	opts := Options{}
	a.StringVar(&opts.RepoDir, "repo", "", "root directory of the git repository (default: current dir)")
	a.StringVar(&opts.From, "from", "", "source ref to start diff from (e.g., 'main')")
	a.StringVar(&opts.To, "to", "", "target ref to end diff at (e.g., 'feature-branch')")
	a.StringVarP(&opts.Commit, "commit", "c", "", "single commit hash or tag to review (vs its parent)")
	a.StringVar(&opts.Excludes, "exclude", "", "comma-separated gitignore-style patterns to exclude")
	a.StringVar(&opts.RulePath, "rule", "", "path to JSON file with system review rules")
	a.StringVarP(&opts.Background, "background", "b", "", "optional requirement/business context")
	a.StringVarP(&opts.BackgroundFile, "background-file", "B", "", "path to a Markdown file used as background")
	a.IntVar(&opts.MaxGitProcs, "max-git-procs", 16, "max concurrent git subprocesses")

	if err := a.Parse(args); err != nil {
		return opts, nil, fmt.Errorf("parse flags: %w", err)
	}

	opts.ShowHelp = a.showHelp
	if opts.ShowHelp {
		return opts, nil, nil
	}

	modeCount := 0
	if opts.From != "" || opts.To != "" {
		modeCount++
	}
	if opts.Commit != "" {
		modeCount++
	}
	if modeCount > 1 {
		return opts, nil, fmt.Errorf("only one review mode allowed (--from/--to or --commit)")
	}
	if opts.From != "" && opts.To == "" {
		return opts, nil, fmt.Errorf("--to is required when --from is specified")
	}
	if opts.To != "" && opts.From == "" {
		return opts, nil, fmt.Errorf("--from is required when --to is specified")
	}

	remaining := a.fs.Args()
	return opts, remaining, nil
}
