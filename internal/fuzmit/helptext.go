package fuzmit

import (
	"fmt"
	"strings"
)

func rootLongDescription(noEmojis bool) string {
	var b strings.Builder
	title := "fuzmit: Conventional Commits, but fuzzy."
	border := strings.Repeat("─", len(title)+2)
	b.WriteString("┌" + border + "┐\n")
	b.WriteString("│ " + title + " │\n")
	b.WriteString("└" + border + "┘\n\n")
	b.WriteString("Conventional Commit Message Format:\n")
	b.WriteString("  <type>(optional scope): <description>\n")
	b.WriteString("  Use \"!\" after type/scope for breaking changes (example: feat!: remove deprecated endpoint)\n\n")
	b.WriteString("Commit Types:\n")
	for _, row := range helpCommitTypeRows(noEmojis) {
		b.WriteString("  " + row + "\n")
	}
	b.WriteString("\nExample Commit Message:\n")
	b.WriteString("  feat(authentication): add support for password reset\n\n")
	b.WriteString("Notes:\n")
	b.WriteString("  - Commit subjects are emoji-free conventional commits\n")
	b.WriteString("  - If FUZMIT_JIRA_SCOPE=true, scope prompt/defaults are ignored and Jira scope is auto-detected\n")
	b.WriteString("  - No staged changes exits cleanly without committing\n")
	b.WriteString("  - Conventional Commits: https://www.conventionalcommits.org/en/v1.0.0/#specification")
	return b.String()
}

func helpCommitTypeRows(noEmojis bool) []string {
	maxName := 0
	for _, ct := range SupportedTypes {
		if len(ct.Name) > maxName {
			maxName = len(ct.Name)
		}
	}

	rows := make([]string, 0, len(SupportedTypes))
	for _, ct := range SupportedTypes {
		if noEmojis {
			rows = append(rows, fmt.Sprintf("%-*s  %s", maxName, ct.Name, ct.Description))
			continue
		}
		rows = append(rows, fmt.Sprintf("%s  %-*s  %s", ct.Emoji, maxName, ct.Name, ct.Description))
	}
	return rows
}
