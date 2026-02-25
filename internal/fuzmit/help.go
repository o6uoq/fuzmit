package fuzmit

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	helpColorTitle       = "#51A4FF"
	helpColorSubtitle    = "#A7ADB8"
	helpColorSection     = "#7D56F4"
	helpColorType        = "#0CB37F"
	helpColorScope       = "#5EA0FF"
	helpColorDescription = "#FF8C5A"
	helpColorBorder      = "#7D56F4"
)

type helpStyleSet struct {
	headerBox   lipgloss.Style
	title       lipgloss.Style
	subtitle    lipgloss.Style
	section     lipgloss.Style
	typeToken   lipgloss.Style
	scopeToken  lipgloss.Style
	descToken   lipgloss.Style
	link        lipgloss.Style
	bullet      lipgloss.Style
	codeExample lipgloss.Style
}

func helpStyles() helpStyleSet {
	return helpStyleSet{
		headerBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(helpColorBorder)).
			Padding(0, 1),
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpColorTitle)),
		subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorSubtitle)),
		section: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpColorSection)),
		typeToken: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorType)),
		scopeToken: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorScope)),
		descToken: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorDescription)),
		link: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorTitle)),
		bullet: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorSection)),
		codeExample: lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpColorTitle)),
	}
}

func PrintHelp(w io.Writer, noEmojis bool) error {
	s := helpStyles()

	header := lipgloss.JoinVertical(
		lipgloss.Left,
		s.title.Render("fuzmit"),
		s.subtitle.Render("Conventional Commits, but Fuzzy."),
	)
	if _, err := fmt.Fprintln(w, s.headerBox.Render(header)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "COMMIT FORMAT", []string{
		s.typeToken.Render("<type>") + s.scopeToken.Render("(optional scope)") + ": " + s.descToken.Render("<description>"),
	}); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "COMMIT TYPES", commitTypeRows(noEmojis)); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "EXAMPLE", []string{
		s.codeExample.Render("feat(authentication): add support for password reset"),
	}); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "USAGE", []string{
		"  fuzmit [--type <type>] [--scope <scope>|--scope] [--jira-scope] [--override] [-m <message>]",
		"  fuzmit defaults",
		"  fuzmit scope [on|off]",
		"  fuzmit jira-scope [on|off]",
	}); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "OPTIONS", optionRows()); err != nil {
		return err
	}

	if err := writeHelpSection(w, s, "NOTES", []string{
		s.bullet.Render("•") + " Commit subjects are emoji-free conventional commits",
		s.bullet.Render("•") + " If FUZMIT_JIRA_SCOPE=true, scope prompt/defaults are ignored and Jira scope is auto-detected",
		s.bullet.Render("•") + " No staged changes exits cleanly without committing",
		s.bullet.Render("•") + " Conventional Commits: " + s.link.Render("https://www.conventionalcommits.org/en/v1.0.0/#specification"),
	}); err != nil {
		return err
	}

	_, err := fmt.Fprintln(w)
	return err
}

func writeHelpSection(w io.Writer, s helpStyleSet, title string, lines []string) error {
	if _, err := fmt.Fprintln(w, s.section.Render(title)); err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w)
	return err
}

func commitTypeRows(noEmojis bool) []string {
	maxName := 0
	for _, ct := range SupportedTypes {
		if len(ct.Name) > maxName {
			maxName = len(ct.Name)
		}
	}

	rows := make([]string, 0, len(SupportedTypes))
	for _, ct := range SupportedTypes {
		if noEmojis {
			rows = append(rows, fmt.Sprintf("  %-*s  %s", maxName, ct.Name, ct.Description))
			continue
		}
		rows = append(rows, fmt.Sprintf("  %s  %-*s  %s", ct.Emoji, maxName, ct.Name, ct.Description))
	}
	return rows
}

func optionRows() []string {
	rows := []struct {
		flag string
		desc string
	}{
		{flag: "--jira-scope", desc: "Auto-detect Jira scope from branch name (e.g. ABC-123)"},
		{flag: "-m, --message <description>", desc: "Set commit description directly"},
		{flag: "--no-emojis", desc: "Disable emojis in commit-type menus/help"},
		{flag: "--override", desc: "Bypass main branch restriction and allow committing"},
		{flag: "--scope", desc: "Prompt for optional scope"},
		{flag: "--scope <scope>", desc: "Set optional scope directly"},
		{flag: "--type <type>", desc: "Set commit type directly (build|chore|ci|docs|feat|fix|perf|refactor|style|test)"},
	}

	maxFlagLen := 0
	for _, row := range rows {
		if len(row.flag) > maxFlagLen {
			maxFlagLen = len(row.flag)
		}
	}

	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, fmt.Sprintf("  %-*s  %s", maxFlagLen, strings.TrimSpace(row.flag), row.desc))
	}
	return out
}
