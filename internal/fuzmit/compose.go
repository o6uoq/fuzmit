package fuzmit

import (
	"fmt"
	"regexp"
	"strings"
)

var jiraScopePattern = regexp.MustCompile(`(?i)[A-Z][A-Z0-9_]+-[0-9]+`)
var scopePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
var conventionalSubjectPattern = regexp.MustCompile(`^(?:[^\x00-\x7F]+ )?[a-z][a-z0-9-]*(?:\([A-Za-z0-9][A-Za-z0-9._/-]*\))?(?:!)?: [^\r\n]+$`)

// ExtractJiraScope extracts a Jira issue key from a branch name.
func ExtractJiraScope(branch string) string {
	branch = strings.TrimSpace(branch)
	locs := jiraScopePattern.FindAllStringIndex(branch, -1)
	for _, loc := range locs {
		start, end := loc[0], loc[1]
		if start > 0 && isJiraKeyChar(branch[start-1]) {
			continue
		}
		if end < len(branch) && isJiraKeyChar(branch[end]) {
			continue
		}
		return strings.ToUpper(branch[start:end])
	}
	return ""
}

func isJiraKeyChar(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	case b == '_':
		return true
	default:
		return false
	}
}

// BuildCommitMessage assembles a conventional commit message.
// If emoji is non-empty it is prepended to the subject.
func BuildCommitMessage(emoji, commitType, scope, description string) string {
	emoji = strings.TrimSpace(emoji)
	commitType = strings.TrimSpace(commitType)
	scope = strings.TrimSpace(scope)
	description = strings.TrimSpace(description)

	var subject string
	if scope != "" {
		scope = strings.Trim(scope, "()")
		subject = commitType + "(" + scope + "): " + description
	} else {
		subject = commitType + ": " + description
	}

	if emoji != "" {
		return emoji + " " + subject
	}
	return subject
}

// ValidateScope validates scope characters for conventional commits.
func ValidateScope(scope string) error {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return nil
	}
	if !scopePattern.MatchString(scope) {
		return fmt.Errorf("scope must match %s", scopePattern.String())
	}
	return nil
}

// ValidateConventionalSubject validates the conventional commit subject line format.
func ValidateConventionalSubject(subject string) error {
	subject = strings.TrimSpace(subject)
	if !conventionalSubjectPattern.MatchString(subject) {
		return fmt.Errorf("commit message must match %s", conventionalSubjectPattern.String())
	}
	return nil
}
