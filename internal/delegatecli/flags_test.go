package delegatecli

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFlags_ModeExclusivity(t *testing.T) {
	_, _, err := ParseFlags([]string{"--from", "main", "--to", "dev", "--commit", "abc123"})
	if err == nil {
		t.Fatal("expected error when both range and commit modes are set")
	}
	if !strings.Contains(err.Error(), "only one review mode") {
		t.Errorf("error = %q, want mode exclusivity message", err)
	}
}

func TestParseFlags_FromRequiresTo(t *testing.T) {
	_, _, err := ParseFlags([]string{"--from", "main"})
	if err == nil {
		t.Fatal("expected error when --from without --to")
	}
	if !strings.Contains(err.Error(), "--to is required") {
		t.Errorf("error = %q", err)
	}
}

func TestParseFlags_ToRequiresFrom(t *testing.T) {
	_, _, err := ParseFlags([]string{"--to", "dev"})
	if err == nil {
		t.Fatal("expected error when --to without --from")
	}
	if !strings.Contains(err.Error(), "--from is required") {
		t.Errorf("error = %q", err)
	}
}

func TestParseFlags_ShortCommit(t *testing.T) {
	opts, remaining, err := ParseFlags([]string{"-c", "abc123", "extra"})
	if err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}
	if opts.Commit != "abc123" {
		t.Errorf("Commit = %q, want abc123", opts.Commit)
	}
	if len(remaining) != 1 || remaining[0] != "extra" {
		t.Errorf("remaining = %v, want [extra]", remaining)
	}
}

func TestParseFlags_Help(t *testing.T) {
	opts, remaining, err := ParseFlags([]string{"-h"})
	if err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}
	if !opts.ShowHelp {
		t.Error("ShowHelp should be true")
	}
	if remaining != nil {
		t.Errorf("remaining = %v, want nil", remaining)
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, _, err := ParseFlags([]string{"--not-a-flag"})
	if err == nil {
		t.Fatal("expected parse error for unknown flag")
	}
	if !strings.Contains(err.Error(), "parse flags") {
		t.Errorf("error = %q", err)
	}
}

func TestSplitPaths(t *testing.T) {
	got := SplitPaths(" a , b , , c ")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("SplitPaths = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("SplitPaths[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestPrintUsage_SaysOcrDelegate(t *testing.T) {
	var buf bytes.Buffer
	PrintUsage(&buf)
	out := buf.String()
	if strings.Contains(out, "ocr delegate") {
		t.Errorf("usage must not contain 'ocr delegate':\n%s", out)
	}
	if !strings.Contains(out, "ocr-delegate") {
		t.Errorf("usage must contain 'ocr-delegate':\n%s", out)
	}
}

func TestRunWithIO_UnknownSubcommand(t *testing.T) {
	var out, errOut bytes.Buffer
	err := RunWithIO(&out, &errOut, []string{"nope"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(errOut.String(), "ocr-delegate -h") {
		t.Errorf("stderr = %q", errOut.String())
	}
}

func TestRunWithIO_Help(t *testing.T) {
	var out bytes.Buffer
	if err := RunWithIO(&out, &out, []string{"-h"}); err != nil {
		t.Fatalf("RunWithIO: %v", err)
	}
	if !strings.Contains(out.String(), "ocr-delegate preview") {
		t.Errorf("help output missing ocr-delegate preview:\n%s", out.String())
	}
}
