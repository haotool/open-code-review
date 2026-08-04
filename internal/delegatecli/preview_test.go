package delegatecli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPreview_WorkspaceFormat(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() { println(\"hi\") }\n")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo})
	if err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	got := out.String()
	for _, want := range []string{
		"# Files (",
		"- mode: workspace",
		"- total_insertions:",
		"- total_deletions:",
		"`main.go`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, got)
		}
	}
}

func TestRunPreview_CommitMode(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "app.go", "package main\n")
	hash1 := commitAll(t, repo, "first")
	writeRepoFile(t, repo, "app.go", "package main\n\nfunc f() {}\n")
	commitAll(t, repo, "second")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo, "--commit", hash1})
	if err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "- mode: commit") {
		t.Errorf("output missing commit mode:\n%s", got)
	}
	if !strings.Contains(got, "- commit: "+hash1) {
		t.Errorf("output missing commit ref:\n%s", got)
	}
}

func TestRunPreview_Exclude(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "src/keep.go", "package main\n")
	writeRepoFile(t, repo, "src/extra.go", "package main\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "src/keep.go", "package main\n\nfunc k() {}\n")
	writeRepoFile(t, repo, "src/extra.go", "package main\n\nfunc e() {}\n")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo, "--exclude", "src/extra.go"})
	if err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "src/extra.go") {
		t.Fatalf("expected extra file in listing:\n%s", got)
	}
	if !strings.Contains(got, "excluded: user_exclude") {
		t.Errorf("expected user_exclude for extra file:\n%s", got)
	}
}

func TestRunPreview_RefInjectionRejected(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "initial")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo, "--commit", "-evil"})
	if err == nil {
		t.Fatal("expected ref injection error")
	}
	if !strings.Contains(err.Error(), "must not start with '-'") {
		t.Errorf("error = %q", err)
	}
}

func TestRunPreview_NoSessionFiles(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	sessionsDir := filepath.Join(home, ".opencodereview", "sessions")
	before, _ := os.ReadDir(sessionsDir)

	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "init")
	writeRepoFile(t, repo, "x.go", "package main\n\nfunc x() {}\n")

	var out bytes.Buffer
	if err := RunPreview(&out, []string{"--repo", repo}); err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	if err := RunRule(&out, []string{"--repo", repo, "x.go"}); err != nil {
		t.Fatalf("RunRule: %v", err)
	}

	after, _ := os.ReadDir(sessionsDir)
	if len(after) > len(before) {
		t.Errorf("sessions dir grew: before=%d after=%d", len(before), len(after))
	}
}

func TestRunPreview_BackgroundFile(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "main.go", "package main\n")
	commitAll(t, repo, "initial")
	writeRepoFile(t, repo, "main.go", "package main\n\nfunc main() {}\n")
	writeRepoFile(t, repo, "docs/context.md", "business requirement context\n")

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", repo, "--background-file", "docs/context.md"})
	if err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "- background:") {
		t.Fatalf("output missing background line:\n%s", got)
	}
	if !strings.Contains(got, "<ocr_user_background>") || !strings.Contains(got, "business requirement context") {
		t.Errorf("output missing background-file content:\n%s", got)
	}
}

func TestRunPreview_NonGitRepo(t *testing.T) {
	nonGit := t.TempDir()

	var out bytes.Buffer
	err := RunPreview(&out, []string{"--repo", nonGit})
	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("error = %q", err)
	}
}

func TestRunPreview_SubcommandHelp(t *testing.T) {
	var out bytes.Buffer
	if err := RunPreview(&out, []string{"-h"}); err != nil {
		t.Fatalf("RunPreview: %v", err)
	}
	if !strings.Contains(out.String(), "ocr-delegate preview") {
		t.Errorf("help = %q", out.String())
	}
}
