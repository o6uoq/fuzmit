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

func helpEnvLines() []string {
	return []string{
		"If FUZMIT_JIRA_SCOPE=true, both --scope and FUZMIT_SCOPE are ignored since Jira scope is auto-detected.",
		"If FUZMIT_SCOPE=true and --scope is not provided, fuzmit prompts for optional scope.",
		"If FUZMIT_NO_EMOJIS=true, commit-type picker/help output is shown without emojis.",
	}
}
