package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alibaba/open-code-review/internal/delegatecli"
)

func TestMain_HelpSaysOcrDelegate(t *testing.T) {
	var out bytes.Buffer
	if err := delegatecli.RunWithIO(&out, &out, []string{"-h"}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	s := out.String()
	if strings.Contains(s, "ocr delegate") {
		t.Errorf("help must not say 'ocr delegate':\n%s", s)
	}
}

func TestRunMain_VersionFlags(t *testing.T) {
	for _, flag := range []string{"--version", "-V", "version"} {
		flag := flag
		t.Run(flag, func(t *testing.T) {
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = w
			code := runMain([]string{flag})
			_ = w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			if code != 0 {
				t.Fatalf("runMain(%q) = %d, want 0", flag, code)
			}
			if !strings.Contains(buf.String(), "ocr-delegate") {
				t.Errorf("stdout missing ocr-delegate: %q", buf.String())
			}
		})
	}
}

func TestRunMain_UnknownSubcommand(t *testing.T) {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	code := runMain([]string{"not-a-real-subcommand"})
	_ = w.Close()
	os.Stderr = old
	var stderr bytes.Buffer
	_, _ = io.Copy(&stderr, r)
	if code != 1 {
		t.Fatalf("runMain unknown subcommand = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "unknown sub-command") {
		t.Errorf("stderr = %q", stderr.String())
	}
}

func TestRunMain_Help(t *testing.T) {
	code := runMain([]string{"-h"})
	if code != 0 {
		t.Fatalf("runMain(-h) = %d, want 0", code)
	}
}

func TestRunMain_NoArgs(t *testing.T) {
	code := runMain(nil)
	if code != 0 {
		t.Fatalf("runMain(nil) = %d, want 0", code)
	}
}

func TestRunMain_PreviewHelp(t *testing.T) {
	code := runMain([]string{"preview", "-h"})
	if code != 0 {
		t.Fatalf("runMain(preview -h) = %d, want 0", code)
	}
}

func TestRunMain_RuleHelp(t *testing.T) {
	code := runMain([]string{"rule", "-h"})
	if code != 0 {
		t.Fatalf("runMain(rule -h) = %d, want 0", code)
	}
}
