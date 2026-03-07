package fuzmit

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateFlagDependencies(t *testing.T) {
	tests := []struct {
		name    string
		opts    runOptions
		wantErr string
	}{
		{
			name: "scope with explicit type",
			opts: runOptions{
				Type:     "fix",
				ScopeSet: true,
				Scope:    "parser",
			},
		},
		{
			name: "scope prompt with explicit type",
			opts: runOptions{
				Type:     "feat",
				ScopeSet: true,
				AskScope: true,
			},
		},
		{
			name: "scope without type fails",
			opts: runOptions{
				ScopeSet: true,
				Scope:    "auth",
			},
			wantErr: "--type is required when --scope is provided",
		},
		{
			name: "no scope no type allowed for interactive flow",
			opts: runOptions{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFlagDependencies(tc.opts)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("validateFlagDependencies() unexpected error: %v", err)
			}
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("validateFlagDependencies() expected error containing %q", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("validateFlagDependencies()=%q want substring %q", err.Error(), tc.wantErr)
				}
			}
		})
	}
}

func TestEnvCommandOutput(t *testing.T) {
	cmd := newEnvCommandWithGetenv(func(key string) string {
		switch key {
		case EnvScope:
			return "true"
		case EnvGeoScope:
			return "maybe"
		case EnvNoEmojis:
			return ""
		default:
			return ""
		}
	})

	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("env command failed: %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"VARIABLE",
		"VALUE",
		"SOURCE",
		EnvScope,
		"env (\"true\")",
		EnvGeoScope,
		"default (invalid: \"maybe\")",
		EnvNoEmojis,
		"default",
		"FUZMIT_JIRA_SCOPE=true ignores --scope and FUZMIT_SCOPE",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected output to contain %q, got %q", want, got)
		}
	}
}

func TestShouldUseUnifiedInteractiveFlow(t *testing.T) {
	tests := []struct {
		name string
		args []string
		opts runOptions
		want bool
	}{
		{
			name: "fully interactive with no flags",
			args: nil,
			opts: runOptions{},
			want: true,
		},
		{
			name: "explicit type disables unified flow",
			args: nil,
			opts: runOptions{Type: "fix"},
			want: false,
		},
		{
			name: "message flag disables unified flow",
			args: nil,
			opts: runOptions{Message: "fix parser panic"},
			want: false,
		},
		{
			name: "positional description disables unified flow",
			args: []string{"fix parser panic"},
			opts: runOptions{},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldUseUnifiedInteractiveFlow(tc.args, tc.opts)
			if got != tc.want {
				t.Fatalf("shouldUseUnifiedInteractiveFlow()=%t want %t", got, tc.want)
			}
		})
	}
}

func TestTypeFlagCompletionIncludesSupportedTypes(t *testing.T) {
	cmd := NewRootCommand()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"__complete", "--type", ""})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("completion command failed: %v", err)
	}

	got := out.String()
	for _, name := range typeNames() {
		if !strings.Contains(got, name) {
			t.Fatalf("expected completion output to contain type %q, got %q", name, got)
		}
	}
}
