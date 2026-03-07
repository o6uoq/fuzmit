package fuzmit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPrintInfoFallback(t *testing.T) {
	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	printInfo(cmd, "No staged changes to commit.")

	if got, want := out.String(), "fuzmit: No staged changes to commit.\n"; got != want {
		t.Fatalf("printInfo()=%q want %q", got, want)
	}
}

func TestPrintCommitFallback(t *testing.T) {
	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	printCommit(cmd, "feat(core): add parser")

	if got, want := out.String(), "fuzmit: Commit message \"feat(core): add parser\"\n"; got != want {
		t.Fatalf("printCommit()=%q want %q", got, want)
	}
}

func TestPrintKeyValueFallback(t *testing.T) {
	cmd := &cobra.Command{}
	out := &bytes.Buffer{}
	cmd.SetOut(out)

	printKeyValue(cmd, "scope", true)

	if got, want := out.String(), "scope: true\n"; got != want {
		t.Fatalf("printKeyValue()=%q want %q", got, want)
	}
}

func TestPrintHelpNotesFallback(t *testing.T) {
	out := &bytes.Buffer{}

	PrintHelpNotes(out)

	got := out.String()
	if !strings.Contains(got, "NOTES") {
		t.Fatalf("expected NOTES section, got %q", got)
	}
	if !strings.Contains(got, conventionalCommitsSpecURL) {
		t.Fatalf("expected Conventional Commits URL, got %q", got)
	}
	if !strings.Contains(got, EnvGeoScope) {
		t.Fatalf("expected %s note, got %q", EnvGeoScope, got)
	}
	if !strings.Contains(got, EnvScope) {
		t.Fatalf("expected %s note, got %q", EnvScope, got)
	}
	if !strings.Contains(got, EnvNoEmojis) {
		t.Fatalf("expected %s note, got %q", EnvNoEmojis, got)
	}
}
