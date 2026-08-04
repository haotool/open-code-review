package main

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestPrintVersion_Dev(t *testing.T) {
	origVersion := Version
	origCommit := GitCommit
	origDate := BuildDate
	t.Cleanup(func() {
		Version = origVersion
		GitCommit = origCommit
		BuildDate = origDate
	})

	Version = "dev"
	GitCommit = ""
	BuildDate = ""

	got := captureStdout(t, printVersion)
	if !strings.Contains(got, "ocr-delegate dev") {
		t.Errorf("expected 'ocr-delegate dev', got %q", got)
	}
	if !strings.Contains(got, runtime.GOOS+"/"+runtime.GOARCH) {
		t.Errorf("expected OS/ARCH, got %q", got)
	}
}

func TestPrintVersion_WithCommitAndDate(t *testing.T) {
	origVersion := Version
	origCommit := GitCommit
	origDate := BuildDate
	t.Cleanup(func() {
		Version = origVersion
		GitCommit = origCommit
		BuildDate = origDate
	})

	Version = "1.2.3"
	GitCommit = "abc1234"
	BuildDate = "2026-01-01"

	got := captureStdout(t, printVersion)
	if !strings.Contains(got, "1.2.3") {
		t.Errorf("expected version, got %q", got)
	}
	if !strings.Contains(got, "abc1234") {
		t.Errorf("expected commit, got %q", got)
	}
	if !strings.Contains(got, "2026-01-01") {
		t.Errorf("expected build date, got %q", got)
	}
}
