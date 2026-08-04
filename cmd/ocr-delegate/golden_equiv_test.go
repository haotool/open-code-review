// Golden regression for the ocr-delegate main entry (task 007 / issue #8).
//
// Reuses fixtures from internal/delegatecli/testdata/golden/ written by
// TestGoldenEquivalence. Update goldens via:
//
//	go test ./internal/delegatecli/... -run TestGoldenEquivalence -golden-update -count=1
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/alibaba/open-code-review/internal/delegatecli"
)

var hashPattern = regexp.MustCompile(`\b[0-9a-f]{40}\b`)

func normalizeGoldenOutput(s, repo string) string {
	if repo != "" {
		s = strings.ReplaceAll(s, repo, "<REPO>")
	}
	return hashPattern.ReplaceAllString(s, "<HASH>")
}

func delegatecliGoldenDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "internal", "delegatecli", "testdata", "golden")
}

type entryScenario struct {
	name string
	args []string
}

func entryScenarios() []entryScenario {
	return []entryScenario{
		{name: "01_top_level_help", args: []string{"-h"}},
		{name: "03_preview_subcommand_help", args: []string{"preview", "-h"}},
		{name: "17_error_unknown_subcommand", args: []string{"nope"}},
	}
}

func TestGoldenEquivalence_OcrDelegateEntry(t *testing.T) {
	base := delegatecliGoldenDir(t)
	for _, sc := range entryScenarios() {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := delegatecli.RunWithIO(&stdout, &stderr, sc.args)
			exit := 0
			if err != nil {
				exit = 1
			}
			dir := filepath.Join(base, sc.name, "ocr_delegate_bin")
			for _, tc := range []struct {
				file string
				got  string
			}{
				{"stdout.txt", normalizeGoldenOutput(stdout.String(), "")},
				{"stderr.txt", normalizeGoldenOutput(stderr.String(), "")},
				{"exit.txt", fmt.Sprintf("%d\n", exit)},
			} {
				wantBytes, readErr := os.ReadFile(filepath.Join(dir, tc.file))
				if readErr != nil {
					t.Fatalf("read golden %s: %v", tc.file, readErr)
				}
				want := string(wantBytes)
				if tc.file == "exit.txt" {
					want = strings.TrimSpace(want)
					got := strings.TrimSpace(tc.got)
					if want != got {
						t.Fatalf("exit: got %q want %q", got, want)
					}
					continue
				}
				if want != tc.got {
					t.Fatalf("%s mismatch\n--- got ---\n%s\n--- want ---\n%s", tc.file, tc.got, want)
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
