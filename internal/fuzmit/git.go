package fuzmit

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// EnsureGitRepo verifies the current directory is inside a git work tree.
func EnsureGitRepo() error {
	_, err := runGit("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return errors.New("fuzmit: not inside a git repository")
	}
	return nil
}

// CurrentBranch resolves the current branch name.
func CurrentBranch() (string, error) {
	branch, err := runGit("symbolic-ref", "--short", "HEAD")
	if err == nil {
		return strings.TrimSpace(branch), nil
	}

	branch, err = runGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("fuzmit: unable to resolve current branch: %w", err)
	}
	return strings.TrimSpace(branch), nil
}

// HasStagedChanges returns true when there are staged changes.
func HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err == nil {
		return false, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return true, nil
	}
	return false, fmt.Errorf("fuzmit: unable to inspect staged changes: %w", err)
}

// Commit executes git commit with the provided message.
func Commit(message string) (string, error) {
	return runGit("commit", "-m", message)
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	text := strings.TrimSpace(out.String())
	if err != nil {
		if text == "" {
			return "", err
		}
		return text, fmt.Errorf("%w: %s", err, text)
	}
	return text, nil
}
