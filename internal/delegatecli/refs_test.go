package delegatecli

import (
	"strings"
	"testing"
)

func TestValidateReviewRefs_RejectsDashPrefix(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")
	cases := []struct {
		name string
		opts Options
		flag string
	}{
		{"from", Options{From: "-evil", To: hash}, "--from"},
		{"to", Options{From: hash, To: "-evil"}, "--to"},
		{"commit", Options{Commit: "-evil"}, "--commit"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateReviewRefs(repo, tc.opts)
			if err == nil {
				t.Fatal("expected error for ref starting with '-'")
			}
			if !strings.Contains(err.Error(), tc.flag) {
				t.Errorf("error = %q, want mention of %s", err, tc.flag)
			}
			if !strings.Contains(err.Error(), "must not start with '-'") {
				t.Errorf("error = %q", err)
			}
		})
	}
}

func TestValidateReviewRefs_InvalidRef(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")
	err := ValidateReviewRefs(repo, Options{Commit: hash + "deadbeef"})
	if err == nil {
		t.Fatal("expected error for invalid commit ref")
	}
	if !strings.Contains(err.Error(), "not a valid commit ref") {
		t.Errorf("error = %q", err)
	}
}

func TestValidateReviewRefs_ValidCommit(t *testing.T) {
	repo := initRepo(t)
	writeRepoFile(t, repo, "x.go", "package main\n")
	hash := commitAll(t, repo, "initial")
	if err := ValidateReviewRefs(repo, Options{Commit: hash}); err != nil {
		t.Fatalf("ValidateReviewRefs: %v", err)
	}
}
