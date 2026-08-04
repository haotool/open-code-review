package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Security regression suite (T6/S7): rules file read boundaries per
// system_rules.go tryReadRuleFile / readRuleFileSafe. Findings outside this
// closure are reported to PM per ADR-001 — do not patch production here.

func TestSecurity_ReadRuleFileSafe_RejectsUnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "evil.json")
	if err := os.WriteFile(f, []byte(`{"weakened":true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := readRuleFileSafe(f)
	if err == nil {
		t.Fatal("expected rejection for non-whitelisted extension")
	}
	if !strings.Contains(err.Error(), "unsupported extension") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_ReadRuleFileSafe_RejectsOversize512KB(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "big.md")
	big := make([]byte, 513*1024)
	if err := os.WriteFile(f, big, 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := readRuleFileSafe(f)
	if err == nil {
		t.Fatal("expected rejection for file exceeding 512KB cap")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("error = %q", err)
	}
}

func TestSecurity_ReadRuleFileSafe_RejectsSymlinkToNonWhitelistedExt(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "secret.pem")
	if err := os.WriteFile(target, []byte("PRIVATE KEY"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "looks.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	_, err := readRuleFileSafe(link)
	if err == nil {
		t.Fatal("expected rejection when symlink resolves to non-whitelisted extension")
	}
}

func TestSecurity_TryReadRuleFile_BlocksRelativeTraversal(t *testing.T) {
	repoDir := t.TempDir()
	outsideDir := t.TempDir()
	outside := filepath.Join(outsideDir, "outside.md")
	if err := os.WriteFile(outside, []byte("traversal secret"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Relative path that escapes repoDir must not load outside content.
	got := tryReadRuleFile(filepath.Join("..", filepath.Base(outsideDir), "outside.md"), repoDir)
	if got != nil {
		t.Fatalf("relative traversal should be blocked, got content %q", *got)
	}
}

func TestSecurity_TryReadRuleFile_AbsPathOutsideRepoAllowed(t *testing.T) {
	// Characterization (S7): absolute rule paths bypass repo boundary by design.
	outsideDir := t.TempDir()
	absRule := filepath.Join(outsideDir, "external.md")
	if err := os.WriteFile(absRule, []byte("external rule text"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := tryReadRuleFile(absRule, t.TempDir())
	if got == nil || *got != "external rule text" {
		t.Fatalf("absolute path outside repo is allowed by design, got %v", got)
	}
}

func TestSecurity_TryReadRuleFile_SymlinkOutsideRepo(t *testing.T) {
	// Characterization: readRuleFileSafe resolves symlinks before reading but does
	// not re-check repo boundary after resolution. Document current behaviour;
	// if outside content loads, report to PM (upstream rules path, not delegate).
	repoDir := t.TempDir()
	outsideDir := t.TempDir()
	outside := filepath.Join(outsideDir, "outside.md")
	if err := os.WriteFile(outside, []byte("OUTSIDE RULE CONTENT"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(repoDir, "link.md")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	got := tryReadRuleFile("link.md", repoDir)
	if got != nil && *got == "OUTSIDE RULE CONTENT" {
		t.Log("PM-DEFECT: symlink inside repo can load rule file outside repo boundary")
	}
}

func TestSecurity_ResolveRuleEntries_EmptyRepoDirRejectsRelative(t *testing.T) {
	entries := []ProjectRuleEntry{
		{Path: "**/*.go", Rule: "rules.md"},
	}
	resolveRuleEntries(entries, "")

	if entries[0].Rule != "" {
		t.Errorf("relative rule with empty repoDir should be cleared, got %q", entries[0].Rule)
	}
}
