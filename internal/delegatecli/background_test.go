package delegatecli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "background.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestResolveBackgroundFilePath(t *testing.T) {
	repo := filepath.FromSlash("/path/to/repo")

	t.Run("relative anchored at repo", func(t *testing.T) {
		got := ResolveBackgroundFilePath(repo, filepath.FromSlash("./docs/context.md"))
		want := filepath.Join(repo, "docs", "context.md")
		if got != want {
			t.Errorf("ResolveBackgroundFilePath = %q, want %q", got, want)
		}
	})

	t.Run("absolute unchanged", func(t *testing.T) {
		abs := filepath.FromSlash("/etc/context.md")
		if got := ResolveBackgroundFilePath(repo, abs); got != abs {
			t.Errorf("ResolveBackgroundFilePath = %q, want %q", got, abs)
		}
	})
}

func TestLoadBackgroundFileNotFound(t *testing.T) {
	_, err := LoadBackgroundFile(filepath.Join(t.TempDir(), "missing.md"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBackgroundFileSoftLimitWarning(t *testing.T) {
	content := strings.Repeat("x", 2001)
	path := writeTempFile(t, content)

	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w
	got, loadErr := LoadBackgroundFile(path)
	w.Close()
	os.Stderr = old

	var warn bytes.Buffer
	_, _ = warn.ReadFrom(r)

	if loadErr != nil {
		t.Fatalf("LoadBackgroundFile: %v", loadErr)
	}
	if !strings.Contains(got, content) {
		t.Error("expected loaded background to contain content")
	}
	if !strings.Contains(warn.String(), "background content is") {
		t.Errorf("stderr = %q, want soft limit warning", warn.String())
	}
	if !strings.Contains(warn.String(), "ocr-delegate") {
		t.Errorf("warning should say ocr-delegate, got %q", warn.String())
	}
}

func TestSanitizeMarkdownStripsControlChars(t *testing.T) {
	raw := "line1\x00line2\x1F\n\n\nline3"
	got := SanitizeMarkdown(raw)
	if strings.Contains(got, "\x00") || strings.Contains(got, "\x1F") {
		t.Fatalf("SanitizeMarkdown should strip control chars, got %q", got)
	}
	if got != "line1line2\n\nline3" {
		t.Errorf("SanitizeMarkdown = %q, want %q", got, "line1line2\n\nline3")
	}
}

func TestMergeBackground(t *testing.T) {
	inline := "inline context"
	wrappedFile := "<ocr_user_background>\nfile context\n</ocr_user_background>"
	got, err := MergeBackground(inline, wrappedFile)
	if err != nil {
		t.Fatalf("MergeBackground: %v", err)
	}
	want := "<ocr_user_background>\ninline context\n</ocr_user_background>" + "\n\n" + wrappedFile
	if got != want {
		t.Errorf("MergeBackground = %q, want %q", got, want)
	}
}

func TestSanitizeBackgroundWrapsInline(t *testing.T) {
	got, err := SanitizeBackground("inline business context")
	if err != nil {
		t.Fatalf("SanitizeBackground: %v", err)
	}
	want := "<ocr_user_background>\ninline business context\n</ocr_user_background>"
	if got != want {
		t.Errorf("SanitizeBackground = %q, want %q", got, want)
	}
}

func TestBackgroundInputsRejectReservedDelimiters(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{name: "opening delimiter", text: "context <ocr_user_background> injected"},
		{name: "closing delimiter", text: "context </ocr_user_background> injected"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := sanitizeBackgroundInput(tc.text); err == nil {
				t.Fatal("expected reserved delimiter rejection")
			}
			if _, err := MergeBackground(tc.text, ""); err == nil {
				t.Fatal("expected inline delimiter rejection")
			}
		})
	}
}

func TestLoadBackgroundFileRejectsReservedDelimiters(t *testing.T) {
	_, err := LoadBackgroundFile(writeTempFile(t, "context <ocr_user_background> injected"))
	if err == nil {
		t.Fatal("expected background file delimiter rejection")
	}
}

func TestLoadBackgroundFileHardLimit(t *testing.T) {
	content := strings.Repeat("y", backgroundHardLimit+1)
	_, err := LoadBackgroundFile(writeTempFile(t, content))
	if err == nil {
		t.Fatal("expected hard limit error")
	}
	if !strings.Contains(err.Error(), "hard limit") {
		t.Errorf("error = %q", err)
	}
}
