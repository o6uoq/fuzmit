package fuzmit

import (
	"strings"
	"testing"
)

func TestRootLongDescription(t *testing.T) {
	got := rootLongDescription(false)

	if got != "fuzmit: Conventional Commits, but Fuzzy." {
		t.Fatalf("unexpected long description: %q", got)
	}
}

func TestHelpNotesLine(t *testing.T) {
	got := helpNotesLine(false)

	if !strings.Contains(got, conventionalCommitsSpecURL) {
		t.Fatalf("missing notes URL: %q", got)
	}
}

func TestHelpTaglineColorized(t *testing.T) {
	got := helpTagline(true)
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("expected ANSI styling, got %q", got)
	}
}
