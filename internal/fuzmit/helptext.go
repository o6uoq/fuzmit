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
