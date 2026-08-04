package delegatecli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRule_OutputFormat(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "handler.go", "package main\n")
	commitAll(t, repo, "init")

	var out bytes.Buffer
	err := RunRule(&out, []string{"--repo", repo, "handler.go"})
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}
	got := out.String()
	for _, want := range []string{
		"### Rule Group 1:",
		"Applies to:",
		"- handler.go",
		"#### Content",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, got)
		}
	}
}

func TestRunRule_CustomRule(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "svc.go", "package main\n")
	commitAll(t, repo, "init")

	customRule := `{"rules":[{"path":"*.go","rule":"Custom Go review rule for delegate test."}]}`
	rulePath := filepath.Join(t.TempDir(), "custom-rule.json")
	if err := os.WriteFile(rulePath, []byte(customRule), 0o644); err != nil {
		t.Fatalf("write rule: %v", err)
	}

	var out bytes.Buffer
	if err := RunRule(&out, []string{"--repo", repo, "--rule", rulePath, "svc.go"}); err != nil {
		t.Fatalf("RunRule: %v", err)
	}
	if !strings.Contains(out.String(), "Custom Go review rule for delegate test.") {
		t.Errorf("output missing custom rule:\n%s", out.String())
	}
}

func TestRunRule_RequiresPath(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	commitAll(t, repo, "init")

	var out bytes.Buffer
	err := RunRule(&out, []string{"--repo", repo})
	if err == nil {
		t.Fatal("expected error when no paths given")
	}
	if !strings.Contains(err.Error(), "at least one file path") {
		t.Errorf("error = %q", err)
	}
}

func TestRunRule_SubcommandHelp(t *testing.T) {
	var out bytes.Buffer
	if err := RunRule(&out, []string{"-h"}); err != nil {
		t.Fatalf("RunRule: %v", err)
	}
	if !strings.Contains(out.String(), "ocr-delegate rule") {
		t.Errorf("help = %q", out.String())
	}
}

func TestRunRule_MultiplePaths(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "a.go", "package main\n")
	writeRepoFile(t, repo, "b.go", "package main\n")
	commitAll(t, repo, "init")

	var out bytes.Buffer
	err := RunRule(&out, []string{"--repo", repo, "a.go", "b.go"})
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "- a.go") || !strings.Contains(got, "- b.go") {
		t.Errorf("output missing file paths:\n%s", got)
	}
}
