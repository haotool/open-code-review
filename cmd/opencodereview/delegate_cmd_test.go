package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alibaba/open-code-review/internal/delegatecli"
)

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func initEquivRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q")
	gitRun(t, dir, "config", "user.email", "test@example.com")
	gitRun(t, dir, "config", "user.name", "test")
	gitRun(t, dir, "config", "commit.gpgsign", "false")
	return dir
}

func writeEquivFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func commitAllEquiv(t *testing.T, dir, msg string) string {
	t.Helper()
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-q", "-m", msg)
	return gitRun(t, dir, "rev-parse", "HEAD")
}

type equivResult struct {
	stdout string
	stderr string
	err    error
}

func runOCRDelegate(t *testing.T, args []string) equivResult {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := runDelegateWithIO(&stdout, &stderr, args)
	return equivResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

func runOCRDelegateBinary(t *testing.T, args []string) equivResult {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := delegatecli.RunWithIO(&stdout, &stderr, args)
	return equivResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

func assertEquivalent(t *testing.T, name string, ocr, bin equivResult) {
	t.Helper()
	ocrExit, binExit := 0, 0
	if ocr.err != nil {
		ocrExit = 1
	}
	if bin.err != nil {
		binExit = 1
	}
	if ocrExit != binExit {
		t.Fatalf("%s: exit mismatch ocr=%d bin=%d\nocr err=%v\nbin err=%v", name, ocrExit, binExit, ocr.err, bin.err)
	}
	if ocr.stdout != bin.stdout {
		t.Fatalf("%s: stdout mismatch\n--- ocr delegate ---\n%s\n--- ocr-delegate ---\n%s", name, ocr.stdout, bin.stdout)
	}
	if ocr.stderr != bin.stderr {
		t.Fatalf("%s: stderr mismatch\n--- ocr delegate ---\n%s\n--- ocr-delegate ---\n%s", name, ocr.stderr, bin.stderr)
	}
	if ocrExit == 1 && !errors.Is(ocr.err, bin.err) && ocr.err.Error() != bin.err.Error() {
		// Error types may differ; message must match for user-visible equivalence.
		if ocr.err.Error() != bin.err.Error() {
			t.Fatalf("%s: error message mismatch\nocr=%v\nbin=%v", name, ocr.err, bin.err)
		}
	}
}

func TestDelegateEquivalence_PreviewWorkspace(t *testing.T) {
	repo := initEquivRepo(t)
	writeEquivFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	commitAllEquiv(t, repo, "initial")
	writeEquivFile(t, repo, "main.go", "package main\n\nfunc main() { println(\"hi\") }\n")

	args := []string{"preview", "--repo", repo}
	assertEquivalent(t, "preview workspace", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateEquivalence_PreviewCommit(t *testing.T) {
	repo := initEquivRepo(t)
	writeEquivFile(t, repo, "app.go", "package main\n")
	hash1 := commitAllEquiv(t, repo, "first")
	writeEquivFile(t, repo, "app.go", "package main\n\nfunc f() {}\n")
	commitAllEquiv(t, repo, "second")

	args := []string{"preview", "--repo", repo, "--commit", hash1}
	assertEquivalent(t, "preview commit", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateEquivalence_PreviewExclude(t *testing.T) {
	repo := initEquivRepo(t)
	writeEquivFile(t, repo, "src/keep.go", "package main\n")
	writeEquivFile(t, repo, "src/extra.go", "package main\n")
	commitAllEquiv(t, repo, "initial")
	writeEquivFile(t, repo, "src/keep.go", "package main\n\nfunc k() {}\n")
	writeEquivFile(t, repo, "src/extra.go", "package main\n\nfunc e() {}\n")

	args := []string{"preview", "--repo", repo, "--exclude", "src/extra.go"}
	assertEquivalent(t, "preview exclude", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateEquivalence_PreviewRange(t *testing.T) {
	repo := initEquivRepo(t)
	gitRun(t, repo, "branch", "-M", "main")
	writeEquivFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	commitAllEquiv(t, repo, "initial")
	gitRun(t, repo, "checkout", "-q", "-b", "feature")
	writeEquivFile(t, repo, "main.go", "package main\n\nfunc main() { println(\"hi\") }\n")
	commitAllEquiv(t, repo, "feature change")

	args := []string{"preview", "--repo", repo, "--from", "main", "--to", "feature"}
	assertEquivalent(t, "preview range", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateEquivalence_RuleSinglePath(t *testing.T) {
	repo := initEquivRepo(t)
	writeEquivFile(t, repo, "handler.go", "package main\n")
	commitAllEquiv(t, repo, "init")

	args := []string{"rule", "--repo", repo, "handler.go"}
	assertEquivalent(t, "rule single", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateEquivalence_RuleMultiplePaths(t *testing.T) {
	repo := initEquivRepo(t)
	writeEquivFile(t, repo, "a.go", "package main\n")
	writeEquivFile(t, repo, "b.go", "package main\n")
	commitAllEquiv(t, repo, "init")

	args := []string{"rule", "--repo", repo, "a.go", "b.go"}
	assertEquivalent(t, "rule multiple", runOCRDelegate(t, args), runOCRDelegateBinary(t, args))
}

func TestDelegateUsageSaysOcrDelegate(t *testing.T) {
	var out bytes.Buffer
	if err := runDelegateWithIO(&out, io.Discard, []string{"-h"}); err != nil {
		t.Fatalf("runDelegateWithIO: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "ocr delegate") {
		t.Errorf("usage must contain 'ocr delegate':\n%s", s)
	}
	if strings.Contains(s, "ocr-delegate") {
		t.Errorf("usage must not contain 'ocr-delegate':\n%s", s)
	}
}
