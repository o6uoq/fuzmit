package fuzmit

import "charm.land/lipgloss/v2"

const conventionalCommitsSpecURL = "https://www.conventionalcommits.org/en/v1.0.0/#specification"

func rootLongDescription(colorize bool) string {
	return helpTagline(colorize)
}

func helpTagline(colorize bool) string {
	tagline := "fuzmit: Conventional Commits, but Fuzzy."
	if !colorize {
		return tagline
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")).Render(tagline)
}

func helpNotesLine(colorize bool) string {
	line := "Conventional Commits: " + conventionalCommitsSpecURL
	if !colorize {
		return line
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")).Render(line)
}

func helpEnvLines(colorize bool) []string {
	env := func(s string) string {
		if !colorize {
			return s
		}
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")).Render(s)
	}
	return []string{
		"If " + env("FUZMIT_JIRA_SCOPE=true") + ", both --scope and " + env("FUZMIT_SCOPE") + " are ignored and Jira scope is detected from branch name",
		"If " + env("FUZMIT_SCOPE=true") + " and --scope is not provided, fuzmit prompts for optional scope",
		"If " + env("FUZMIT_NO_EMOJIS=true") + ", emojis are omitted from commit subjects, picker menus, and help output",
	}
}
