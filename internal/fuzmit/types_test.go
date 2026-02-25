package fuzmit

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestFindCommitType(t *testing.T) {
	ct, ok := FindCommitType("fix")
	if !ok {
		t.Fatalf("expected fix to resolve")
	}
	if ct.Name != "fix" {
		t.Fatalf("expected fix, got %s", ct.Name)
	}

	if _, ok := FindCommitType("unknown"); ok {
		t.Fatalf("unknown type should not resolve")
	}
}

func TestFormatTypeLabelAlignsDashNoEmoji(t *testing.T) {
	build, _ := FindCommitType("build")
	refactor, _ := FindCommitType("refactor")

	buildLabel := FormatTypeLabel(build, true)
	refactorLabel := FormatTypeLabel(refactor, true)

	buildDash := prefixDisplayWidth(buildLabel)
	refactorDash := prefixDisplayWidth(refactorLabel)
	if buildDash == -1 || refactorDash == -1 {
		t.Fatalf("expected labels to contain separator: %q | %q", buildLabel, refactorLabel)
	}
	if buildDash != refactorDash {
		t.Fatalf("expected aligned separator, got %d and %d", buildDash, refactorDash)
	}
}

func TestFormatTypeLabelAlignsDashWithEmoji(t *testing.T) {
	chore, _ := FindCommitType("chore")
	perf, _ := FindCommitType("perf")

	choreLabel := FormatTypeLabel(chore, false)
	perfLabel := FormatTypeLabel(perf, false)

	choreDash := prefixDisplayWidth(choreLabel)
	perfDash := prefixDisplayWidth(perfLabel)
	if choreDash == -1 || perfDash == -1 {
		t.Fatalf("expected labels to contain separator: %q | %q", choreLabel, perfLabel)
	}
	if choreDash != perfDash {
		t.Fatalf("expected aligned separator, got %d and %d", choreDash, perfDash)
	}
}

func TestFormatTypeLabelAlignsDashBuildPerfRefactor(t *testing.T) {
	build, _ := FindCommitType("build")
	perf, _ := FindCommitType("perf")
	refactor, _ := FindCommitType("refactor")

	buildLabel := FormatTypeLabel(build, false)
	perfLabel := FormatTypeLabel(perf, false)
	refactorLabel := FormatTypeLabel(refactor, false)

	buildDash := prefixDisplayWidth(buildLabel)
	perfDash := prefixDisplayWidth(perfLabel)
	refactorDash := prefixDisplayWidth(refactorLabel)

	if buildDash == -1 || perfDash == -1 || refactorDash == -1 {
		t.Fatalf("expected labels to contain separator: %q | %q | %q", buildLabel, perfLabel, refactorLabel)
	}
	if buildDash != perfDash || buildDash != refactorDash {
		t.Fatalf("expected aligned separator, got build=%d perf=%d refactor=%d", buildDash, perfDash, refactorDash)
	}
}

func prefixDisplayWidth(label string) int {
	prefix, _, ok := strings.Cut(label, " - ")
	if !ok {
		return -1
	}
	return runewidth.StringWidth(prefix)
}
