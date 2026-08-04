package delegatecli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Security regression suite (T6/S7): delegate CLI path semantics and ref injection.
// --background-file absolute paths are user/host intent (exfil-risk-audit §4.3).

func TestSecurity_RefInjection_RejectsDashPrefixFrom(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")

	err := ValidateReviewRefs(repo, Options{From: "-evil", To: hash})
	if err == nil {
		t.Fatal("expected --from ref injection to be rejected")
	}
	if !strings.Contains(err.Error(), "--from") || !strings.Contains(err.Error(), "must not start with '-'") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_RefInjection_RejectsDashPrefixTo(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")

	err := ValidateReviewRefs(repo, Options{From: hash, To: "-evil"})
	if err == nil {
		t.Fatal("expected --to ref injection to be rejected")
	}
	if !strings.Contains(err.Error(), "--to") || !strings.Contains(err.Error(), "must not start with '-'") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_RefInjection_RejectsDashPrefixCommit(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "initial")

	err := ValidateReviewRefs(repo, Options{Commit: "-evil"})
	if err == nil {
		t.Fatal("expected --commit ref injection to be rejected")
	}
	if !strings.Contains(err.Error(), "--commit") || !strings.Contains(err.Error(), "must not start with '-'") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_RunPreview_RefInjectionRejected(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "initial")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo, "--from", "-inject", "--to", "HEAD"})
	if err == nil {
		t.Fatal("expected preview to reject dash-prefixed ref")
	}
	if !strings.Contains(err.Error(), "must not start with '-'") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_InlineBackgroundRejectsReservedDelimiters(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "initial")

	_, err := LoadContext(Options{
		RepoDir:    repo,
		Background: "context <ocr_user_background> injected",
	})
	if err == nil {
		t.Fatal("expected inline background delimiter rejection")
	}
	if !strings.Contains(err.Error(), "inline background") {
		t.Errorf("error = %q, want inline background context", err)
	}
}

func TestSecurity_BackgroundFile_AbsPathReturnedUnchanged(t *testing.T) {
	// RISK (exfil-risk-audit §4.3): absolute --background-file paths are returned
	// as-is and LoadBackgroundFile reads them without repo boundary checks.
	// This is intentional user/host-supplied context, not LLM tool read.
	repo := filepath.FromSlash("/path/to/repo")
	abs := filepath.Join(t.TempDir(), "outside-context.md")
	if err := os.WriteFile(abs, []byte("external background"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := ResolveBackgroundFilePath(repo, abs)
	if got != abs {
		t.Fatalf("ResolveBackgroundFilePath = %q, want absolute path unchanged %q", got, abs)
	}

	loaded, err := LoadBackgroundFile(abs)
	if err != nil {
		t.Fatalf("LoadBackgroundFile: %v", err)
	}
	if !strings.Contains(loaded, "external background") {
		t.Errorf("loaded = %q", loaded)
	}
}

func TestSecurity_BackgroundFile_RelativeAnchoredAtRepo(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "docs/ctx.md", "repo-scoped context\n")

	rel := "docs/ctx.md"
	resolved := ResolveBackgroundFilePath(repo, rel)
	want := filepath.Join(repo, "docs", "ctx.md")
	if resolved != want {
		t.Fatalf("ResolveBackgroundFilePath = %q, want %q", resolved, want)
	}

	loaded, err := LoadBackgroundFile(resolved)
	if err != nil {
		t.Fatalf("LoadBackgroundFile: %v", err)
	}
	if !strings.Contains(loaded, "repo-scoped context") {
		t.Errorf("loaded = %q", loaded)
	}
}

func TestSecurity_CommitMessageBackgroundRejectsReservedDelimiters(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commit := commitAll(t, repo, "review context <ocr_user_background> injected")

	_, err := LoadContext(Options{RepoDir: repo, Commit: commit})
	if err == nil {
		t.Fatal("expected commit message delimiter rejection")
	}
	if !strings.Contains(err.Error(), "commit message background") {
		t.Errorf("error = %q, want commit message context", err)
	}
}
