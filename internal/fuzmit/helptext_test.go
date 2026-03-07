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

func TestHelpEnvLines(t *testing.T) {
	lines := helpEnvLines(false)
	if len(lines) != 3 {
		t.Fatalf("expected 3 env notes lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], EnvGeoScope) {
		t.Fatalf("missing %s guidance: %q", EnvGeoScope, lines[0])
	}
	if !strings.Contains(lines[1], EnvScope) {
		t.Fatalf("missing %s guidance: %q", EnvScope, lines[1])
	}
	if !strings.Contains(lines[2], EnvNoEmojis) {
		t.Fatalf("missing %s guidance: %q", EnvNoEmojis, lines[2])
	}
}
