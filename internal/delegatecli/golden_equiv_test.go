// Golden equivalence regression suite (task 007 / issue #8).
//
// # Updating golden files
//
// When intentional CLI output changes land, refresh fixtures from this package:
//
//	go test ./internal/delegatecli/... -run TestGoldenEquivalence -golden-update -count=1
//
// Review diffs under testdata/golden/ before committing. Each scenario stores
// stdout.txt, stderr.txt, and exit.txt. Dynamic repo paths and 40-char git
// hashes are normalized to <REPO> and <HASH> before compare/update.
//
// # Runners
//
//   - ocr-delegate: delegatecli.RunWithIO (canonical SSOT path)
//   - ocr-delegate: runOCRDelegateSim mirrors cmd/opencodereview/delegate_cmd.go
//
// Scenarios with equiv=true require identical stdout, stderr, and exit code
// between both runners. Top-level help/usage (branded strings) use per-runner
// golden subdirs instead.
package delegatecli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var updateGoldens = flag.Bool("golden-update", false, "rewrite testdata/golden fixtures")

type goldenRunner string

const (
	runnerOCRDelegate    goldenRunner = "ocr_delegate"
	runnerOCRDelegateBin goldenRunner = "ocr_delegate_bin"
)

type goldenEnv struct {
	repo      string
	rulePath  string
	bgFile    string
	extraHash string
}

type goldenScenario struct {
	name  string
	equiv bool
	setup func(t *testing.T) goldenEnv
	args  []string
	// redirectGlobalStderr captures os.Stderr writes (e.g. background soft limit).
	redirectGlobalStderr bool
}

func goldenDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "testdata", "golden")
}

var hashPattern = regexp.MustCompile(`\b[0-9a-f]{40}\b`)

func normalizeGoldenOutput(s, repo string) string {
	if repo != "" {
		s = strings.ReplaceAll(s, repo, "<REPO>")
	}
	return hashPattern.ReplaceAllString(s, "<HASH>")
}

type cliOutcome struct {
	stdout string
	stderr string
	exit   int
	err    error
}

func (o cliOutcome) normalized(repo string) cliOutcome {
	return cliOutcome{
		stdout: normalizeGoldenOutput(o.stdout, repo),
		stderr: normalizeGoldenOutput(o.stderr, repo),
		exit:   o.exit,
		err:    o.err,
	}
}

func runOCRDelegateBin(out, errOut io.Writer, args []string) error {
	return RunWithIO(out, errOut, args)
}

// runOCRDelegateSim mirrors cmd/opencodereview/delegate_cmd.go routing.
func runOCRDelegateSim(out, errOut io.Writer, args []string) error {
	if len(args) == 0 {
		printOCRDelegateUsage(out)
		return nil
	}
	switch args[0] {
	case "-h", "--help":
		printOCRDelegateUsage(out)
		return nil
	case "preview":
		return RunPreview(out, args[1:])
	case "rule":
		return RunRule(out, args[1:])
	default:
		fmt.Fprintf(errOut, "unknown delegate sub-command: %s\nRun 'ocr delegate -h' for usage\n", args[0])
		return fmt.Errorf("unknown delegate sub-command: %s", args[0])
	}
}

func printOCRDelegateUsage(w io.Writer) {
	fmt.Fprintln(w, `OpenCodeReview - Delegation Mode

Usage:
  ocr delegate <sub-command> [flags]
  ocr d <sub-command> [flags]       (alias)

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
  ocr delegate preview --from main --to feature

  # Preview workspace changes
  ocr delegate preview

  # Get rules for multiple files (grouped by content)
  ocr delegate rule internal/agent/agent.go internal/llm/client.go`)
}

func execGoldenRunner(runner goldenRunner, stdout, stderr io.Writer, args []string) cliOutcome {
	var err error
	switch runner {
	case runnerOCRDelegate:
		err = runOCRDelegateSim(stdout, stderr, args)
	case runnerOCRDelegateBin:
		err = runOCRDelegateBin(stdout, stderr, args)
	default:
		return cliOutcome{
			exit: 1,
			err:  fmt.Errorf("unknown runner: %q", runner),
		}
	}
	exit := 0
	if err != nil {
		exit = 1
	}
	return cliOutcome{
		stdout: stdout.(*bytes.Buffer).String(),
		stderr: stderr.(*bytes.Buffer).String(),
		exit:   exit,
		err:    err,
	}
}

func runGoldenScenario(runner goldenRunner, args []string, redirectGlobalStderr bool) cliOutcome {
	var stdout, stderr bytes.Buffer
	if redirectGlobalStderr {
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		outcome := execGoldenRunner(runner, &stdout, &stderr, args)
		w.Close()
		os.Stderr = old
		var extra bytes.Buffer
		_, _ = extra.ReadFrom(r)
		outcome.stderr += extra.String()
		return outcome
	}
	return execGoldenRunner(runner, &stdout, &stderr, args)
}

func readGoldenFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeGoldenFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func assertGoldenFiles(t *testing.T, dir string, repo string, got cliOutcome, update bool) {
	t.Helper()
	n := got.normalized(repo)
	files := map[string]string{
		"stdout.txt": n.stdout,
		"stderr.txt": n.stderr,
		"exit.txt":   fmt.Sprintf("%d\n", n.exit),
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if update {
			if err := writeGoldenFile(path, content); err != nil {
				t.Fatalf("write golden %s: %v", path, err)
			}
			continue
		}
		want, err := readGoldenFile(path)
		if err != nil {
			t.Fatalf("read golden %s: %v (run with -update to create)", path, err)
		}
		if name == "exit.txt" {
			want = strings.TrimSpace(want)
			gotExit := strings.TrimSpace(content)
			if want != gotExit {
				t.Fatalf("exit mismatch: got %q want %q", gotExit, want)
			}
			continue
		}
		if want != content {
			t.Fatalf("%s mismatch\n--- got ---\n%s\n--- want ---\n%s", name, content, want)
		}
	}
}

func assertEquivalentOutcomes(t *testing.T, name string, a, b cliOutcome) {
	t.Helper()
	if a.exit != b.exit {
		t.Fatalf("%s: exit mismatch ocr=%d bin=%d\nocr err=%v\nbin err=%v", name, a.exit, b.exit, a.err, b.err)
	}
	if a.stdout != b.stdout {
		t.Fatalf("%s: stdout mismatch\n--- ocr delegate ---\n%s\n--- ocr-delegate ---\n%s", name, a.stdout, b.stdout)
	}
	if a.stderr != b.stderr {
		t.Fatalf("%s: stderr mismatch\n--- ocr delegate ---\n%s\n--- ocr-delegate ---\n%s", name, a.stderr, b.stderr)
	}
}

func setupWorkspaceRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() { println(\"hi\") }\n")
	return goldenEnv{repo: repo}
}

func setupRangeRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	gitRun(t, repo, "branch", "-M", "main")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	commitAll(t, repo, "initial")
	gitRun(t, repo, "checkout", "-q", "-b", "feature")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() { println(\"hi\") }\n")
	commitAll(t, repo, "feature change")
	return goldenEnv{repo: repo}
}

func setupCommitRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "app.go", "package main\n")
	hash1 := commitAll(t, repo, "first")
	writeRepoFile(t, repo, "app.go", "package main\n\nfunc f() {}\n")
	commitAll(t, repo, "second")
	return goldenEnv{repo: repo, extraHash: hash1}
}

func setupExcludeRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "src/keep.go", "package main\n")
	writeRepoFile(t, repo, "src/extra.go", "package main\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "src/keep.go", "package main\n\nfunc k() {}\n")
	writeRepoFile(t, repo, "src/extra.go", "package main\n\nfunc e() {}\n")
	return goldenEnv{repo: repo}
}

func setupRuleRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "handler.go", "package main\n")
	commitAll(t, repo, "init")
	return goldenEnv{repo: repo}
}

func setupMultiRuleRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "a.go", "package main\n")
	writeRepoFile(t, repo, "b.go", "package main\n")
	commitAll(t, repo, "init")
	return goldenEnv{repo: repo}
}

func setupCustomRule(t *testing.T) goldenEnv {
	repo := setupRuleRepo(t).repo
	customRule := `{"rules":[{"path":"*.go","rule":"Custom Go review rule for golden equivalence test."}]}`
	rulePath := filepath.Join(t.TempDir(), "custom-rule.json")
	if err := os.WriteFile(rulePath, []byte(customRule), 0o644); err != nil {
		t.Fatalf("write rule: %v", err)
	}
	return goldenEnv{repo: repo, rulePath: rulePath}
}

func setupBackgroundRepos(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "main.go", "package main\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	writeRepoFile(t, repo, "docs/context.md", "business requirement context\n")
	return goldenEnv{repo: repo}
}

func setupSoftLimitBackground(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "main.go", "package main\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	content := strings.Repeat("x", 2001)
	writeRepoFile(t, repo, "docs/large.md", content+"\n")
	return goldenEnv{repo: repo}
}

func setupErrorRepo(t *testing.T) goldenEnv {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")
	return goldenEnv{repo: repo, extraHash: hash}
}

func setupBareRepo(t *testing.T) goldenEnv {
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "--bare", "-q", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %v\n%s", err, out)
	}
	return goldenEnv{repo: dir}
}

func goldenScenarios() []goldenScenario {
	return []goldenScenario{
		{name: "01_top_level_help", equiv: false, setup: func(t *testing.T) goldenEnv { return goldenEnv{} }, args: []string{"-h"}},
		{name: "02_empty_usage", equiv: false, setup: func(t *testing.T) goldenEnv { return goldenEnv{} }, args: []string{}},
		{name: "03_preview_subcommand_help", equiv: true, setup: setupWorkspaceRepo, args: []string{"preview", "-h"}},
		{name: "04_rule_subcommand_help", equiv: true, setup: setupRuleRepo, args: []string{"rule", "-h"}},
		{name: "05_preview_workspace", equiv: true, setup: setupWorkspaceRepo, args: []string{"preview", "--repo", "<repo>"}},
		{name: "06_preview_range", equiv: true, setup: setupRangeRepo, args: []string{"preview", "--repo", "<repo>", "--from", "main", "--to", "feature"}},
		{name: "07_preview_commit", equiv: true, setup: setupCommitRepo, args: []string{"preview", "--repo", "<repo>", "--commit", "<hash1>"}},
		{name: "08_preview_exclude", equiv: true, setup: setupExcludeRepo, args: []string{"preview", "--repo", "<repo>", "--exclude", "src/extra.go"}},
		{name: "09_preview_custom_rule", equiv: true, setup: setupCustomRule, args: []string{"preview", "--repo", "<repo>", "--rule", "<rule>"}},
		{name: "10_preview_background_inline", equiv: true, setup: setupBackgroundRepos, args: []string{"preview", "--repo", "<repo>", "-b", "inline business context"}},
		{name: "11_preview_background_file", equiv: true, setup: setupBackgroundRepos, args: []string{"preview", "--repo", "<repo>", "-B", "docs/context.md"}},
		{name: "12_preview_background_both", equiv: true, setup: setupBackgroundRepos, args: []string{"preview", "--repo", "<repo>", "-b", "inline", "-B", "docs/context.md"}},
		{name: "13_preview_explicit_repo", equiv: true, setup: setupWorkspaceRepo, args: []string{"preview", "--repo", "<repo>"}},
		{name: "14_rule_single_path", equiv: true, setup: setupRuleRepo, args: []string{"rule", "--repo", "<repo>", "handler.go"}},
		{name: "15_rule_multiple_paths", equiv: true, setup: setupMultiRuleRepo, args: []string{"rule", "--repo", "<repo>", "a.go", "b.go"}},
		{name: "16_rule_custom_rule", equiv: true, setup: setupCustomRule, args: []string{"rule", "--repo", "<repo>", "--rule", "<rule>", "handler.go"}},
		{name: "17_error_unknown_subcommand", equiv: false, setup: func(t *testing.T) goldenEnv { return goldenEnv{} }, args: []string{"nope"}},
		{name: "18_error_unknown_flag", equiv: true, setup: setupWorkspaceRepo, args: []string{"preview", "--repo", "<repo>", "--not-a-flag"}},
		{name: "19_error_mode_exclusivity", equiv: true, setup: setupErrorRepo, args: []string{"preview", "--repo", "<repo>", "--from", "main", "--to", "HEAD", "--commit", "<hash1>"}},
		{name: "20_error_from_without_to", equiv: true, setup: setupErrorRepo, args: []string{"preview", "--repo", "<repo>", "--from", "main"}},
		{name: "21_error_ref_injection", equiv: true, setup: setupErrorRepo, args: []string{"preview", "--repo", "<repo>", "--commit", "-evil"}},
		{name: "22_error_invalid_commit_ref", equiv: true, setup: setupErrorRepo, args: []string{"preview", "--repo", "<repo>", "--commit", "<badhash>"}},
		{name: "23_error_rule_no_paths", equiv: true, setup: setupRuleRepo, args: []string{"rule", "--repo", "<repo>"}},
		{name: "24_edge_bare_repo", equiv: true, setup: setupBareRepo, args: []string{"preview", "--repo", "<repo>"}},
		{name: "25_edge_non_git_dir", equiv: true, setup: func(t *testing.T) goldenEnv { return goldenEnv{repo: t.TempDir()} }, args: []string{"preview", "--repo", "<repo>"}},
		{name: "26_edge_background_soft_limit", equiv: true, setup: setupSoftLimitBackground, args: []string{"preview", "--repo", "<repo>", "-B", "docs/large.md"}, redirectGlobalStderr: true},
	}
}

func substituteArgs(args []string, env goldenEnv) []string {
	out := make([]string, len(args))
	copy(out, args)
	for i, a := range out {
		switch a {
		case "<repo>":
			out[i] = env.repo
		case "<rule>":
			out[i] = env.rulePath
		case "<hash1>":
			out[i] = env.extraHash
		case "<badhash>":
			out[i] = env.extraHash + "deadbeef"
		}
	}
	return out
}

func TestGoldenEquivalence(t *testing.T) {
	base := goldenDir(t)
	update := *updateGoldens

	for _, sc := range goldenScenarios() {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			env := sc.setup(t)
			args := substituteArgs(sc.args, env)

			ocr := runGoldenScenario(runnerOCRDelegate, args, sc.redirectGlobalStderr)
			bin := runGoldenScenario(runnerOCRDelegateBin, args, sc.redirectGlobalStderr)

			if sc.equiv {
				assertEquivalentOutcomes(t, sc.name, ocr, bin)
				dir := filepath.Join(base, sc.name, string(runnerOCRDelegateBin))
				assertGoldenFiles(t, dir, env.repo, bin, update)
				return
			}

			for _, runner := range []goldenRunner{runnerOCRDelegate, runnerOCRDelegateBin} {
				got := ocr
				if runner == runnerOCRDelegateBin {
					got = bin
				}
				dir := filepath.Join(base, sc.name, string(runner))
				assertGoldenFiles(t, dir, env.repo, got, update)
			}
		})
	}
}
