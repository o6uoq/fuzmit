package fuzmit

import (
	"strings"
	"testing"
)

func TestCLIContracts_TypeScopeAndEnvMatrix(t *testing.T) {
	tests := []struct {
		name        string
		stdin       string
		env         map[string]string
		args        []string
		setup       func(t *testing.T, repoDir string)
		verify      func(t *testing.T, repoDir, output string)
		wantOutput  string
		wantSubject string
		wantErr     string
	}{
		{
			name:        "type and message with emoji",
			args:        []string{"--override", "--type", "fix", "-m", "patch parser"},
			wantSubject: "🔧 fix: patch parser",
		},
		{
			name:    "invalid type errors",
			args:    []string{"--override", "--type", "invalid", "-m", "patch parser"},
			wantErr: "invalid --type \"invalid\"",
		},
		{
			name:        "type scope and message with emoji",
			args:        []string{"--override", "--type", "fix", "--scope=auth", "-m", "patch parser"},
			wantSubject: "🔧 fix(auth): patch parser",
		},
		{
			name:        "type scope prompt with value uses scope",
			stdin:       "auth\n",
			args:        []string{"--override", "--type", "feat", "--scope", "-m", "add login"},
			wantSubject: "🚀 feat(auth): add login",
		},
		{
			name:        "type scope prompt with blank keeps no scope",
			stdin:       "\n",
			args:        []string{"--override", "--type", "feat", "--scope", "-m", "add login"},
			wantSubject: "🚀 feat: add login",
		},
		{
			name:        "no-emojis omits emoji from subject",
			args:        []string{"--override", "--type", "fix", "--no-emojis", "-m", "patch parser"},
			wantSubject: "fix: patch parser",
		},
		{
			name: "no-emojis via env omits emoji",
			env: map[string]string{
				EnvNoEmojis: "true",
			},
			args:        []string{"--override", "--type", "fix", "-m", "patch parser"},
			wantSubject: "fix: patch parser",
		},
		{
			name:    "scope without type errors",
			args:    []string{"--override", "--scope=auth", "-m", "patch parser"},
			wantErr: "--type is required when --scope is provided",
		},
		{
			name: "jira env auto detects scope",
			env: map[string]string{
				EnvGeoScope: "true",
			},
			args:        []string{"--override", "--type", "fix", "-m", "patch parser"},
			wantSubject: "🔧 fix(ABC-123): patch parser",
		},
		{
			name:        "jira scope flag auto detects scope",
			args:        []string{"--override", "--jira-scope", "--type", "fix", "-m", "patch parser"},
			wantSubject: "🔧 fix(ABC-123): patch parser",
		},
		{
			name:        "jira scope flag overrides explicit scope",
			args:        []string{"--override", "--jira-scope", "--type", "fix", "--scope=auth", "-m", "patch parser"},
			wantSubject: "🔧 fix(ABC-123): patch parser",
		},
		{
			name:    "invalid scope value errors",
			args:    []string{"--override", "--type", "fix", "--scope=bad?", "-m", "patch parser"},
			wantErr: "invalid --scope value \"bad?\"",
		},
		{
			name: "main branch requires override",
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				runGitInDir(t, repoDir, "checkout", "-B", "main")
			},
			args:    []string{"--type", "fix", "-m", "patch parser"},
			wantErr: "you are on the main branch; use --override to bypass this check",
		},
		{
			name: "no staged changes exits cleanly",
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				runGitInDir(t, repoDir, "commit", "-m", "test: initial")
			},
			args: []string{"--override", "--type", "fix", "-m", "patch parser"},
			verify: func(t *testing.T, repoDir, _ string) {
				t.Helper()
				if got := runGitInDir(t, repoDir, "rev-list", "--count", "HEAD"); got != "1" {
					t.Fatalf("expected no new commit when nothing is staged, got commit count %q", got)
				}
			},
		},
		{
			name: "scope env prompt blank keeps no scope",
			env: map[string]string{
				EnvScope: "true",
			},
			stdin:       "\n",
			args:        []string{"--override", "--type", "docs", "-m", "update readme"},
			wantSubject: "📚 docs: update readme",
		},
		{
			name: "scope env prompt provided value uses scope",
			env: map[string]string{
				EnvScope: "true",
			},
			stdin:       "cli\n",
			args:        []string{"--override", "--type", "docs", "-m", "update readme"},
			wantSubject: "📚 docs(cli): update readme",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoDir := newTempRepoWithStagedFile(t)
			if tc.setup != nil {
				tc.setup(t, repoDir)
			}
			output, err := runRootCommandInDir(t, repoDir, tc.stdin, tc.env, tc.args...)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error=%q want substring %q", err.Error(), tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantOutput != "" && !strings.Contains(output, tc.wantOutput) {
				t.Fatalf("output=%q want substring %q", output, tc.wantOutput)
			}
			if tc.verify != nil {
				tc.verify(t, repoDir, output)
			}
			if tc.wantSubject != "" {
				got := lastCommitSubjectInDir(t, repoDir)
				if got != tc.wantSubject {
					t.Fatalf("subject=%q want %q", got, tc.wantSubject)
				}
			}
		})
	}
}
