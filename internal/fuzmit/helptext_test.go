package fuzmit

import (
	"strings"
	"testing"
)

func TestRootLongDescriptionCommitFormat(t *testing.T) {
	got := rootLongDescription(false)

	if !strings.Contains(got, "<type>(optional scope): <description>") {
		t.Fatalf("missing commit format line: %q", got)
	}
	if strings.Contains(got, "<type>(optional scope)!: <description>") {
		t.Fatalf("commit format should not force breaking marker: %q", got)
	}
}

func TestHelpCommitTypeRowsNoEmojiAlignment(t *testing.T) {
	rows := helpCommitTypeRows(true)

	want := []string{
		"build     Project build or dependencies changes",
		"refactor  Code refactoring without behavior change",
	}
	for _, w := range want {
		found := false
		for _, row := range rows {
			if row == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing row %q in %q", w, rows)
		}
	}
}
