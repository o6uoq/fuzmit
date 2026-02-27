package fuzmit

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runGitInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func newTempRepoWithStagedFile(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGitInDir(t, dir, "init")
	runGitInDir(t, dir, "config", "user.name", "fuzmit-test")
	runGitInDir(t, dir, "config", "user.email", "fuzmit-test@example.com")
	runGitInDir(t, dir, "checkout", "-b", "feat/ABC-123-test")

	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("staged content\n"), 0o644); err != nil {
		t.Fatalf("write staged file: %v", err)
	}
	runGitInDir(t, dir, "add", "test.txt")
	return dir
}

func runRootCommandInDir(t *testing.T, dir string, stdin string, env map[string]string, args ...string) (string, error) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir to temp repo: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousDir)
	})

	t.Setenv(EnvScope, "")
	t.Setenv(EnvGeoScope, "")
	t.Setenv(EnvNoEmojis, "")
	for key, value := range env {
		t.Setenv(key, value)
	}

	output := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(args)

	return output.String(), cmd.Execute()
}

func lastCommitSubjectInDir(t *testing.T, dir string) string {
	t.Helper()
	return runGitInDir(t, dir, "log", "-1", "--pretty=%s")
}
