package fuzmit

import (
	"strings"
	"testing"
)

func TestInteractiveFlow_HappyPathCommitsSubject(t *testing.T) {
	stubPromptInteractiveFlow(t, func(_ bool, askScope bool) (interactiveCommitAnswers, error) {
		if askScope {
			t.Fatalf("askScope should be false for default env settings")
		}
		return interactiveCommitAnswers{
			Type:        "test",
			Description: "add coverage",
		}, nil
	})

	repoDir := newTempRepoWithStagedFile(t)
	_, err := runRootCommandInDir(t, repoDir, "", nil, "--override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := lastCommitSubjectInDir(t, repoDir)
	if got != "test: add coverage" {
		t.Fatalf("subject=%q want %q", got, "test: add coverage")
	}
}

func TestInteractiveFlow_AbortReturnsFriendlyError(t *testing.T) {
	stubPromptInteractiveFlow(t, func(_ bool, _ bool) (interactiveCommitAnswers, error) {
		return interactiveCommitAnswers{}, errSelectionAborted
	})

	repoDir := newTempRepoWithStagedFile(t)
	_, err := runRootCommandInDir(t, repoDir, "", nil, "--override")
	if err == nil {
		t.Fatal("expected error from interactive abort")
	}
	if !strings.Contains(err.Error(), "no commit type selected, aborting") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInteractiveFlow_ScopeEnabledByEnvPromptsForScope(t *testing.T) {
	askedForScope := false
	stubPromptInteractiveFlow(t, func(_ bool, askScope bool) (interactiveCommitAnswers, error) {
		askedForScope = askScope
		return interactiveCommitAnswers{
			Type:        "feat",
			Scope:       "auth",
			Description: "add login",
		}, nil
	})

	repoDir := newTempRepoWithStagedFile(t)
	_, err := runRootCommandInDir(t, repoDir, "", map[string]string{
		EnvScope: "true",
	}, "--override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !askedForScope {
		t.Fatal("expected interactive prompt to ask for scope when FUZMIT_SCOPE=true")
	}

	got := lastCommitSubjectInDir(t, repoDir)
	if got != "feat(auth): add login" {
		t.Fatalf("subject=%q want %q", got, "feat(auth): add login")
	}
}

func TestInteractiveFlow_JiraScopeEnvDisablesScopePrompt(t *testing.T) {
	askedForScope := true
	stubPromptInteractiveFlow(t, func(_ bool, askScope bool) (interactiveCommitAnswers, error) {
		askedForScope = askScope
		return interactiveCommitAnswers{
			Type:        "fix",
			Scope:       "ignored",
			Description: "patch parser",
		}, nil
	})

	repoDir := newTempRepoWithStagedFile(t)
	_, err := runRootCommandInDir(t, repoDir, "", map[string]string{
		EnvGeoScope: "true",
	}, "--override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if askedForScope {
		t.Fatal("scope prompt should be disabled when FUZMIT_JIRA_SCOPE=true")
	}

	got := lastCommitSubjectInDir(t, repoDir)
	if got != "fix(ABC-123): patch parser" {
		t.Fatalf("subject=%q want %q", got, "fix(ABC-123): patch parser")
	}
}
