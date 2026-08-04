package preview

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alibaba/open-code-review/internal/config/rules"
	"github.com/alibaba/open-code-review/internal/model"
)

// ---- WhyExcluded：副檔名允許清單 ----

func TestWhyExcluded_ExtensionAllowlist(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected model.ExcludeReason
	}{
		{"unsupported extension txt", "README.txt", model.ExcludeExtension},
		{"unsupported extension md", "docs/guide.md", model.ExcludeExtension},
		{"supported extension go", "main.go", model.ExcludeNone},
		{"supported extension java", "src/Main.java", model.ExcludeNone},
		{"supported extension ts", "app.ts", model.ExcludeNone},
		{"extension check is case-insensitive", "server.GO", model.ExcludeNone},
		{"file without extension", "Makefile", model.ExcludeNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhyExcluded(nil, model.Diff{NewPath: tt.path})
			if got != tt.expected {
				t.Errorf("WhyExcluded(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

// ---- WhyExcluded：預設排除 glob（含 {a,b} 花括號展開）----

func TestWhyExcluded_DefaultExcludeGlobs(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected model.ExcludeReason
	}{
		{"go test file at root", "foo_test.go", model.ExcludeDefaultPath},
		{"go test file in nested dir", "internal/agent/foo_test.go", model.ExcludeDefaultPath},
		{"java standard test dir", "service/src/test/java/com/example/FooTest.java", model.ExcludeDefaultPath},
		{"brace expansion matches .test.ts", "src/app.test.ts", model.ExcludeDefaultPath},
		{"brace expansion matches .test.js", "lib/util.test.js", model.ExcludeDefaultPath},
		{"brace expansion matches .spec.jsx", "components/Button.spec.jsx", model.ExcludeDefaultPath},
		{"brace expansion does not match .test.go", "src/app.test.go", model.ExcludeNone},
		{"regular source file", "handler.go", model.ExcludeNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhyExcluded(nil, model.Diff{NewPath: tt.path})
			if got != tt.expected {
				t.Errorf("WhyExcluded(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

// ---- WhyExcluded：使用者 include/exclude ----

func TestWhyExcluded_UserExclude(t *testing.T) {
	f := &rules.FileFilter{Exclude: []string{"vendor/**", "*.gen.go"}}

	tests := []struct {
		name     string
		path     string
		expected model.ExcludeReason
	}{
		{"path matching exclude glob", "vendor/foo/bar.go", model.ExcludeUserRule},
		{"generated file excluded", "api.gen.go", model.ExcludeUserRule},
		{"unrelated file not excluded", "main.go", model.ExcludeNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhyExcluded(f, model.Diff{NewPath: tt.path})
			if got != tt.expected {
				t.Errorf("WhyExcluded(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestWhyExcluded_UserIncludeBypassesDefaults(t *testing.T) {
	f := &rules.FileFilter{Include: []string{"**/*_test.go"}}

	tests := []struct {
		name     string
		path     string
		expected model.ExcludeReason
	}{
		{"included test file bypasses default-path exclusion", "foo_test.go", model.ExcludeNone},
		{"non-included file falls through to defaults", "main.go", model.ExcludeNone},
		{"non-included file still hit by extension filter", "README.md", model.ExcludeExtension},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhyExcluded(f, model.Diff{NewPath: tt.path})
			if got != tt.expected {
				t.Errorf("WhyExcluded(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestWhyExcluded_UserExcludeWinsOverInclude(t *testing.T) {
	f := &rules.FileFilter{
		Include: []string{"src/**/*.go"},
		Exclude: []string{"src/generated/**"},
	}

	tests := []struct {
		name     string
		path     string
		expected model.ExcludeReason
	}{
		{"file matching both include and exclude is excluded", "src/generated/api.go", model.ExcludeUserRule},
		{"file matching include only is reviewed", "src/handler.go", model.ExcludeNone},
		{"file outside include with valid ext still reviewed (additive)", "lib/utils.go", model.ExcludeNone},
		{"file outside include still hit by default path", "lib/utils_test.go", model.ExcludeDefaultPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhyExcluded(f, model.Diff{NewPath: tt.path})
			if got != tt.expected {
				t.Errorf("WhyExcluded(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

// ---- WhyExcluded：binary 優先於其他規則 ----

func TestWhyExcluded_BinaryTakesPriority(t *testing.T) {
	f := &rules.FileFilter{Exclude: []string{"vendor/**"}}

	if got := WhyExcluded(nil, model.Diff{NewPath: "image.png", IsBinary: true}); got != model.ExcludeBinary {
		t.Errorf("WhyExcluded(binary) = %q, want %q", got, model.ExcludeBinary)
	}
	if got := WhyExcluded(f, model.Diff{NewPath: "vendor/image.png", IsBinary: true}); got != model.ExcludeBinary {
		t.Errorf("WhyExcluded(binary in excluded dir) = %q, want %q (binary checked first)", got, model.ExcludeBinary)
	}
}

// ---- 匯出 helpers ----

func TestEffectivePath(t *testing.T) {
	tests := []struct {
		name     string
		diff     model.Diff
		expected string
	}{
		{"normal new path", model.Diff{OldPath: "old.go", NewPath: "new.go"}, "new.go"},
		{"deleted file falls back to old path", model.Diff{OldPath: "deleted.go", NewPath: "/dev/null"}, "deleted.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EffectivePath(tt.diff); got != tt.expected {
				t.Errorf("EffectivePath() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDiffStatus(t *testing.T) {
	tests := []struct {
		name     string
		diff     model.Diff
		expected string
	}{
		{"binary file", model.Diff{IsBinary: true}, "binary"},
		{"new file", model.Diff{IsNew: true}, "added"},
		{"deleted file", model.Diff{IsDeleted: true}, "deleted"},
		{"renamed flag", model.Diff{IsRenamed: true}, "renamed"},
		{"renamed by path mismatch", model.Diff{OldPath: "old.go", NewPath: "new.go"}, "renamed"},
		{"modified file", model.Diff{OldPath: "main.go", NewPath: "main.go"}, "modified"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffStatus(tt.diff); got != tt.expected {
				t.Errorf("DiffStatus() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExtFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"nested go file", "a/b/main.go", ".go"},
		{"uppercase extension lowered", "SERVER.GO", ".go"},
		{"no extension", "Makefile", ""},
		{"hidden file without extension", ".gitignore", ""},
		{"multi-dot keeps last segment", "archive.tar.gz", ".gz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtFromPath(tt.path); got != tt.expected {
				t.Errorf("ExtFromPath(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

// ---- Run：workspace / range / commit 三模式 diff 載入 ----

func runGitTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGitTest(t, repo, "init", "-q", "-b", "main")
	runGitTest(t, repo, "config", "user.email", "test@example.com")
	runGitTest(t, repo, "config", "user.name", "Test User")
	runGitTest(t, repo, "config", "commit.gpgsign", "false")
	return repo
}

func writeFile(t *testing.T, repo, rel, content string) {
	t.Helper()
	path := filepath.Join(repo, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

func entryByPath(t *testing.T, p *model.Preview, path string) model.PreviewEntry {
	t.Helper()
	for _, e := range p.Entries {
		if e.Path == path {
			return e
		}
	}
	t.Fatalf("entry %q not found in preview entries %+v", path, p.Entries)
	return model.PreviewEntry{}
}

func TestRun_WorkspaceMode(t *testing.T) {
	repo := initRepo(t)
	writeFile(t, repo, "a.go", "package a\n\nvar x = 1\n")
	runGitTest(t, repo, "add", "-A")
	runGitTest(t, repo, "commit", "-q", "-m", "initial")

	// tracked 修改 + untracked 新檔，兩者都須進入 workspace preview
	writeFile(t, repo, "a.go", "package a\n\nvar x = 2\n")
	writeFile(t, repo, "b.go", "package a\n\nvar y = 1\n")

	p, err := Run(context.Background(), Args{RepoDir: repo})
	if err != nil {
		t.Fatalf("Run(workspace) error: %v", err)
	}
	if p.TotalFiles != 2 {
		t.Fatalf("TotalFiles = %d, want 2", p.TotalFiles)
	}
	a := entryByPath(t, p, "a.go")
	if a.Status != "modified" || !a.WillReview {
		t.Errorf("a.go = {Status:%q WillReview:%v}, want {modified true}", a.Status, a.WillReview)
	}
	b := entryByPath(t, p, "b.go")
	if b.Status != "added" || !b.WillReview {
		t.Errorf("b.go = {Status:%q WillReview:%v}, want {added true}", b.Status, b.WillReview)
	}
	if p.ReviewableCount != 2 || p.ExcludedCount != 0 {
		t.Errorf("counts = {reviewable:%d excluded:%d}, want {2 0}", p.ReviewableCount, p.ExcludedCount)
	}
}

func TestRun_RangeMode(t *testing.T) {
	repo := initRepo(t)
	writeFile(t, repo, "a.go", "package a\n\nvar x = 1\n")
	runGitTest(t, repo, "add", "-A")
	runGitTest(t, repo, "commit", "-q", "-m", "base")
	runGitTest(t, repo, "checkout", "-q", "-b", "feature")
	writeFile(t, repo, "a.go", "package a\n\nvar x = 2\n")
	runGitTest(t, repo, "commit", "-q", "-am", "change on feature")

	p, err := Run(context.Background(), Args{RepoDir: repo, From: "main", To: "feature"})
	if err != nil {
		t.Fatalf("Run(range) error: %v", err)
	}
	if p.TotalFiles != 1 {
		t.Fatalf("TotalFiles = %d, want 1", p.TotalFiles)
	}
	a := entryByPath(t, p, "a.go")
	if a.Status != "modified" || !a.WillReview {
		t.Errorf("a.go = {Status:%q WillReview:%v}, want {modified true}", a.Status, a.WillReview)
	}
	if p.TotalInsertions != 1 || p.TotalDeletions != 1 {
		t.Errorf("totals = {+%d -%d}, want {+1 -1}", p.TotalInsertions, p.TotalDeletions)
	}
}

func TestRun_CommitMode(t *testing.T) {
	repo := initRepo(t)
	writeFile(t, repo, "a.go", "package a\n\nvar x = 1\n")
	runGitTest(t, repo, "add", "-A")
	runGitTest(t, repo, "commit", "-q", "-m", "c1")
	writeFile(t, repo, "b.go", "package a\n\nvar y = 1\n")
	writeFile(t, repo, "README.md", "hello\n")
	runGitTest(t, repo, "add", "-A")
	runGitTest(t, repo, "commit", "-q", "-m", "c2")

	p, err := Run(context.Background(), Args{RepoDir: repo, Commit: "HEAD"})
	if err != nil {
		t.Fatalf("Run(commit) error: %v", err)
	}
	// 只能包含 c2 引入的變更，不得回溯 c1
	for _, e := range p.Entries {
		if e.Path == "a.go" {
			t.Fatal("commit mode must only include the reviewed commit's changes")
		}
	}
	if p.TotalFiles != 2 {
		t.Fatalf("TotalFiles = %d, want 2", p.TotalFiles)
	}
	b := entryByPath(t, p, "b.go")
	if b.Status != "added" || !b.WillReview {
		t.Errorf("b.go = {Status:%q WillReview:%v}, want {added true}", b.Status, b.WillReview)
	}
	md := entryByPath(t, p, "README.md")
	if md.WillReview || md.ExcludeReason != model.ExcludeExtension {
		t.Errorf("README.md = {WillReview:%v Reason:%q}, want {false %q}", md.WillReview, md.ExcludeReason, model.ExcludeExtension)
	}
	if p.ReviewableCount != 1 || p.ExcludedCount != 1 {
		t.Errorf("counts = {reviewable:%d excluded:%d}, want {1 1}", p.ReviewableCount, p.ExcludedCount)
	}
	// 總計須涵蓋被排除檔案的行數（與 agent.Preview 舊行為一致）
	if p.TotalInsertions != 4 || p.TotalDeletions != 0 {
		t.Errorf("totals = {+%d -%d}, want {+4 -0}", p.TotalInsertions, p.TotalDeletions)
	}
}

func TestRun_DeletedFileExcluded(t *testing.T) {
	repo := initRepo(t)
	writeFile(t, repo, "gone.go", "package a\n\nvar x = 1\n")
	runGitTest(t, repo, "add", "-A")
	runGitTest(t, repo, "commit", "-q", "-m", "initial")
	runGitTest(t, repo, "rm", "-q", "gone.go")

	p, err := Run(context.Background(), Args{RepoDir: repo})
	if err != nil {
		t.Fatalf("Run(workspace with deletion) error: %v", err)
	}
	e := entryByPath(t, p, "gone.go")
	if e.Status != "deleted" || e.WillReview || e.ExcludeReason != model.ExcludeDeleted {
		t.Errorf("gone.go = {Status:%q WillReview:%v Reason:%q}, want {deleted false %q}",
			e.Status, e.WillReview, e.ExcludeReason, model.ExcludeDeleted)
	}
	if p.ReviewableCount != 0 || p.ExcludedCount != 1 {
		t.Errorf("counts = {reviewable:%d excluded:%d}, want {0 1}", p.ReviewableCount, p.ExcludedCount)
	}
}

func TestRun_ErrorKeepsAgentPreviewFormat(t *testing.T) {
	// 錯誤字串須保留 agent.Preview 既有的 "load diffs: get diffs: " 前綴，
	// 讓 CLI 輸出在抽取前後逐字一致。
	_, err := Run(context.Background(), Args{RepoDir: t.TempDir()})
	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
	if !strings.Contains(err.Error(), "load diffs: get diffs: ") {
		t.Fatalf("error = %q, want it to contain %q", err.Error(), "load diffs: get diffs: ")
	}
}
