package delegatecli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

const (
	backgroundSoftLimit    = 2000
	backgroundHardLimit    = 8000
	backgroundOpenTag      = "<ocr_user_background>"
	backgroundCloseTag     = "</ocr_user_background>"
	maxBackgroundFileBytes = 1 << 20 // 1 MB
)

var multiNewline = regexp.MustCompile(`\n{3,}`)

// ResolveBackgroundFilePath resolves a --background-file path relative to repoDir.
func ResolveBackgroundFilePath(repoDir, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(repoDir, path)
}

// MergeBackground combines inline --background with --background-file content.
// Both sides are validated and wrapped with reserved delimiters before merge.
func MergeBackground(inline, fromFile string) (string, error) {
	var err error
	inline, err = prepareBackground(inline, "inline background")
	if err != nil {
		return "", err
	}
	switch {
	case inline == "":
		return fromFile, nil
	case fromFile == "":
		return inline, nil
	default:
		return inline + "\n\n" + fromFile, nil
	}
}

// LoadBackgroundFile reads and sanitises a Markdown background file.
func LoadBackgroundFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("read background file %q: %w", path, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("background file %q is a directory, not a file", path)
	}
	if info.Size() > maxBackgroundFileBytes {
		return "", fmt.Errorf(
			"background file %q is %d bytes, exceeding the maximum of %d bytes; please provide a smaller file",
			path, info.Size(), maxBackgroundFileBytes,
		)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read background file %q: %w", path, err)
	}

	cleaned, err := sanitizeBackgroundInput(string(raw))
	if err != nil {
		return "", fmt.Errorf("background file %q: %w", path, err)
	}
	if cleaned == "" {
		return "", fmt.Errorf("background file %q is empty after sanitisation", path)
	}

	wrapped, err := wrapBackgroundContent(cleaned)
	if err != nil {
		return "", fmt.Errorf("background file %q: %w", path, err)
	}
	return wrapped, nil
}

func sanitizeBackgroundInput(s string) (string, error) {
	if err := rejectBackgroundDelimiters(s); err != nil {
		return "", err
	}
	cleaned := SanitizeMarkdown(s)
	if err := rejectBackgroundDelimiters(cleaned); err != nil {
		return "", err
	}
	return cleaned, nil
}

func rejectBackgroundDelimiters(s string) error {
	if strings.Contains(s, backgroundOpenTag) || strings.Contains(s, backgroundCloseTag) {
		return fmt.Errorf(
			"background must not contain the reserved delimiters %q or %q",
			backgroundOpenTag, backgroundCloseTag,
		)
	}
	return nil
}

// SanitizeBackground validates, sanitises, and wraps caller-supplied background text.
// Inline --background and commit-message fallback use the same delimiter wrapping as
// --background-file so host agents can treat all background channels uniformly.
func SanitizeBackground(s string) (string, error) {
	return prepareBackground(s, "inline background")
}

func prepareBackground(s, errPrefix string) (string, error) {
	cleaned, err := sanitizeBackgroundInput(s)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errPrefix, err)
	}
	if cleaned == "" {
		return "", nil
	}
	wrapped, err := wrapBackgroundContent(cleaned)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errPrefix, err)
	}
	return wrapped, nil
}

func wrapBackgroundContent(cleaned string) (string, error) {
	if n := len([]rune(cleaned)); n > backgroundHardLimit {
		return "", fmt.Errorf(
			"background content is %d characters, exceeding the hard limit of %d (aborting)",
			n, backgroundHardLimit,
		)
	} else if n > backgroundSoftLimit {
		fmt.Fprintf(os.Stderr,
			"[ocr-delegate] background content is %d characters, exceeding the recommended %d (continuing but review quality might be impacted)\n",
			n, backgroundSoftLimit,
		)
	}
	return backgroundOpenTag + "\n" + cleaned + "\n" + backgroundCloseTag, nil
}

func SanitizeMarkdown(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		switch r {
		case '\n', '\t':
			b.WriteRune(r)
			continue
		case '\r':
			continue
		}
		if isForbiddenChar(r) {
			continue
		}
		b.WriteRune(r)
	}

	collapsed := multiNewline.ReplaceAllString(b.String(), "\n\n")
	return strings.TrimSpace(collapsed)
}

func isForbiddenChar(r rune) bool {
	switch {
	case r <= 0x1F:
		return true
	case r >= 0x7F && r <= 0x9F:
		return true
	}

	switch r {
	case '\u200B', '\u200C', '\u200D', '\u200E', '\u200F', '\u2060', '\u00AD', '\uFEFF':
		return true
	}

	return unicode.Is(unicode.Cf, r)
}
