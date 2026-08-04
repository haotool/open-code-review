// Package delegatecli is the single source of truth for the ocr-delegate
// binary: flag parsing, context loading, ref validation, background files,
// and preview/rule subcommands without LLM, telemetry, or session I/O.
//
// Exported helper API (task 003 SSOT):
//
//	| Symbol                  | Kind   | Purpose                                      |
//	|-------------------------|--------|----------------------------------------------|
//	| Options                 | type   | Parsed delegate flags                        |
//	| Context                 | type   | Loaded repo/rules/git state                  |
//	| ParseFlags              | func   | Flag parsing + mode validation               |
//	| LoadContext             | func   | Rules/repo/git/background load (no LLM)      |
//	| ValidateReviewRefs      | func   | Ref-injection rejection + git verify         |
//	| MergeBackground         | func   | Inline + file background merge               |
//	| LoadBackgroundFile      | func   | Markdown background load/sanitize            |
//	| ResolveBackgroundFilePath | func | Repo-relative background path                |
//	| ApplyCLIExcludes        | func   | Append --exclude patterns to delegate Context|
//	| ApplyExcludesToFilter   | func   | Append --exclude patterns to FileFilter      |
//	| SplitPaths              | func   | Comma-separated exclude/path split           |
//	| Run                     | func   | Top-level CLI entry                          |
//	| RunWithIO               | func   | CLI entry with injectable stdout/stderr      |
//	| RunPreview              | func   | Preview subcommand                           |
//	| RunRule                 | func   | Rule subcommand                              |
//	| FormatPreview           | func   | Preview markdown output                      |
//	| PrintUsage              | func   | Top-level usage text                         |
//	| Context.ReviewMode      | method | workspace/range/commit mode string           |
package delegatecli
