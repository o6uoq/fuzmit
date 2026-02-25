package fuzmit

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestPrintHelpCommitFormatOmitsForcedBreakingBang(t *testing.T) {
	var out bytes.Buffer
	if err := PrintHelp(&out, false); err != nil {
		t.Fatalf("print help failed: %v", err)
	}

	plain := ansiRE.ReplaceAllString(out.String(), "")
	if !strings.Contains(plain, "<type>(optional scope): <description>") {
		t.Fatalf("commit format line missing: %q", plain)
	}
	if strings.Contains(plain, "<type>(optional scope)!: <description>") {
		t.Fatalf("help should not force a breaking-change marker in the format line: %q", plain)
	}
}

func TestCommitTypeRowsNoEmojiAlignment(t *testing.T) {
	rows := commitTypeRows(true)

	contains := func(want string) {
		for _, row := range rows {
			if row == want {
				return
			}
		}
		t.Fatalf("missing row %q in %q", want, rows)
	}

	contains("  build     Project build or dependencies changes")
	contains("  refactor  Code refactoring without behavior change")
}
