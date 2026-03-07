package fuzmit

import "testing"

func TestExtractJiraScope(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		want   string
	}{
		{name: "ticket branch prefix", branch: "ABC-123/feature", want: "ABC-123"},
		{name: "ticket branch nested", branch: "feat/ABC-123-login-flow", want: "ABC-123"},
		{name: "ticket lowercase branch", branch: "bugfix/abcd-12345-redirect", want: "ABCD-12345"},
		{name: "ticket with underscore in project key", branch: "feature/APP_CORE-12-add", want: "APP_CORE-12"},
		{name: "ticket only", branch: "PAY-9", want: "PAY-9"},
		{name: "missing hyphen", branch: "feat/ABCD12345-login", want: ""},
		{name: "no ticket", branch: "feat/new-ui", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractJiraScope(tt.branch)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildCommitMessage(t *testing.T) {
	if got := BuildCommitMessage("🔧", "fix", "auth", "patch nil check"); got != "🔧 fix(auth): patch nil check" {
		t.Fatalf("unexpected message: %s", got)
	}

	if got := BuildCommitMessage("📚", "docs", "", "update readme"); got != "📚 docs: update readme" {
		t.Fatalf("unexpected message: %s", got)
	}

	if got := BuildCommitMessage("", "fix", "auth", "patch nil check"); got != "fix(auth): patch nil check" {
		t.Fatalf("unexpected message (no emoji): %s", got)
	}
}

func TestValidateScope(t *testing.T) {
	ok := []string{"auth", "API-123", "pkg/http", "build.v2", "core_utils"}
	for _, scope := range ok {
		if err := ValidateScope(scope); err != nil {
			t.Fatalf("scope %q should be valid, got %v", scope, err)
		}
	}

	bad := []string{"", "has space", "(auth)", "auth:", "auth?"}
	for _, scope := range bad {
		err := ValidateScope(scope)
		if scope == "" {
			if err != nil {
				t.Fatalf("empty scope should be allowed, got %v", err)
			}
			continue
		}
		if err == nil {
			t.Fatalf("scope %q should be invalid", scope)
		}
	}
}

func TestValidateConventionalSubject(t *testing.T) {
	valid := []string{
		"fix: patch nil check",
		"feat(parser): add array support",
		"perf(api-v2)!: remove deprecated endpoint",
		"🔧 fix: patch nil check",
		"🚀 feat(parser): add array support",
	}
	for _, s := range valid {
		if err := ValidateConventionalSubject(s); err != nil {
			t.Fatalf("subject %q should be valid, got %v", s, err)
		}
	}

	invalid := []string{
		"Fix: uppercase type",
		"feat parser: missing scope parens",
		"feat: ",
	}
	for _, s := range invalid {
		if err := ValidateConventionalSubject(s); err == nil {
			t.Fatalf("subject %q should be invalid", s)
		}
	}
}
