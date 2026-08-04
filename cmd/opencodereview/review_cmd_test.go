package main

import (
	"strings"
	"testing"

	"github.com/alibaba/open-code-review/internal/delegatecli"
)

func TestValidateReviewRefsRejectsOptionLikeCommit(t *testing.T) {
	err := delegatecli.ValidateReviewRefs(t.TempDir(), delegatecli.Options{Commit: "-O./pwn.sh"})
	if err == nil {
		t.Fatal("expected option-like --commit ref to be rejected")
	}
	if !strings.Contains(err.Error(), "--commit") || !strings.Contains(err.Error(), "must not start with '-'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateReviewRefsRejectsOptionLikeRangeRef(t *testing.T) {
	err := delegatecli.ValidateReviewRefs(t.TempDir(), delegatecli.Options{To: "-O./pwn.sh"})
	if err == nil {
		t.Fatal("expected option-like --to ref to be rejected")
	}
	if !strings.Contains(err.Error(), "--to") || !strings.Contains(err.Error(), "must not start with '-'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseReviewFlagsRejectsToWithoutFrom(t *testing.T) {
	_, err := parseReviewFlags([]string{"--to", "HEAD"})
	if err == nil {
		t.Fatal("expected --to without --from to fail")
	}
	if !strings.Contains(err.Error(), "--from is required when --to is specified") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseReviewFlagsRejectsFromWithoutTo(t *testing.T) {
	_, err := parseReviewFlags([]string{"--from", "main"})
	if err == nil {
		t.Fatal("expected --from without --to to fail")
	}
	if !strings.Contains(err.Error(), "--to is required when --from is specified") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseReviewFlagsAllowsFromAndTo(t *testing.T) {
	opts, err := parseReviewFlags([]string{"--from", "main", "--to", "HEAD"})
	if err != nil {
		t.Fatalf("expected --from/--to to pass, got: %v", err)
	}
	if opts.from != "main" || opts.to != "HEAD" {
		t.Fatalf("unexpected opts: from=%q to=%q", opts.from, opts.to)
	}
}

func TestRunReviewPreviewRejectsReservedBackgroundDelimiter(t *testing.T) {
	repo := initTestGitRepo(t)
	err := runReview([]string{
		"--repo", repo,
		"--preview",
		"--background", "context <ocr_user_background> injected",
	})
	if err == nil {
		t.Fatal("expected full review path to reject inline background delimiter")
	}
	if !strings.Contains(err.Error(), "inline background") {
		t.Errorf("error = %q, want inline background context", err)
	}
}
