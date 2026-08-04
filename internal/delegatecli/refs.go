package delegatecli

import (
	"fmt"
	"strings"
)

// ValidateReviewRefs rejects ref-option injection and verifies commit refs exist.
func ValidateReviewRefs(repoDir string, opts Options) error {
	refs := []struct {
		flag string
		ref  string
	}{
		{"--from", opts.From},
		{"--to", opts.To},
		{"--commit", opts.Commit},
	}
	for _, item := range refs {
		if item.ref == "" {
			continue
		}
		if strings.HasPrefix(item.ref, "-") {
			return fmt.Errorf("%s value %q is not a valid git ref: refs must not start with '-'", item.flag, item.ref)
		}
		if out, err := runGitCmd(repoDir, "rev-parse", "--verify", "--end-of-options", item.ref+"^{commit}"); err != nil {
			msg := strings.TrimSpace(string(out))
			if msg != "" {
				return fmt.Errorf("%s value %q is not a valid commit ref: %s", item.flag, item.ref, msg)
			}
			return fmt.Errorf("%s value %q is not a valid commit ref", item.flag, item.ref)
		}
	}
	return nil
}
