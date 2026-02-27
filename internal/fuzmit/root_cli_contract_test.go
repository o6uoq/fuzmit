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
		wantSubject string
		wantErr     string
	}{
		{
			name:        "type and message",
			args:        []string{"--override", "--type", "fix", "-m", "patch parser"},
			wantSubject: "fix: patch parser",
		},
		{
			name:        "type scope and message",
			args:        []string{"--override", "--type", "fix", "--scope=auth", "-m", "patch parser"},
			wantSubject: "fix(auth): patch parser",
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
			wantSubject: "fix(ABC-123): patch parser",
		},
		{
			name: "scope env prompt blank keeps no scope",
			env: map[string]string{
				EnvScope: "true",
			},
			stdin:       "\n",
			args:        []string{"--override", "--type", "docs", "-m", "update readme"},
			wantSubject: "docs: update readme",
		},
		{
			name: "scope env prompt provided value uses scope",
			env: map[string]string{
				EnvScope: "true",
			},
			stdin:       "cli\n",
			args:        []string{"--override", "--type", "docs", "-m", "update readme"},
			wantSubject: "docs(cli): update readme",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoDir := newTempRepoWithStagedFile(t)
			_, err := runRootCommandInDir(t, repoDir, tc.stdin, tc.env, tc.args...)

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

			got := lastCommitSubjectInDir(t, repoDir)
			if got != tc.wantSubject {
				t.Fatalf("subject=%q want %q", got, tc.wantSubject)
			}
		})
	}
}
