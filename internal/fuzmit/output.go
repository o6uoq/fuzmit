package fuzmit

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

type outputLevel string

const (
	outputInfo   outputLevel = "info"
	outputCommit outputLevel = "commit"
)

func printInfo(cmd *cobra.Command, message string) {
	printStatus(cmd.OutOrStdout(), outputInfo, message)
}

func printInfof(cmd *cobra.Command, format string, args ...any) {
	printInfo(cmd, fmt.Sprintf(format, args...))
}

func printCommit(cmd *cobra.Command, subject string) {
	printStatus(cmd.OutOrStdout(), outputCommit, fmt.Sprintf("Commit message %q", subject))
}

func printKeyValue(cmd *cobra.Command, key string, value any) {
	w := cmd.OutOrStdout()
	if !supportsStyling(w) {
		_, _ = fmt.Fprintf(w, "%s: %v\n", key, value)
		return
	}

	cs := defaultFangColorScheme()
	keyStyle := lipgloss.NewStyle().Foreground(cs.Command).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(cs.Base)
	_, _ = fmt.Fprintf(w, "%s: %s\n", keyStyle.Render(key), valueStyle.Render(fmt.Sprint(value)))
}

func printEnvSettings(cmd *cobra.Command, settings []EnvSetting) {
	w := cmd.OutOrStdout()
	nameCol, valueCol, sourceCol := "VARIABLE", "VALUE", "SOURCE"
	rows := make([][3]string, 0, len(settings))
	nameWidth, valueWidth, sourceWidth := len(nameCol), len(valueCol), len(sourceCol)
	for _, setting := range settings {
		row := [3]string{
			setting.Name,
			fmt.Sprintf("%t", setting.Value),
			envSettingSource(setting),
		}
		rows = append(rows, row)
		if len(row[0]) > nameWidth {
			nameWidth = len(row[0])
		}
		if len(row[1]) > valueWidth {
			valueWidth = len(row[1])
		}
		if len(row[2]) > sourceWidth {
			sourceWidth = len(row[2])
		}
	}
	note := "FUZMIT_JIRA_SCOPE=true ignores --scope and FUZMIT_SCOPE."
	if !supportsStyling(w) {
		_, _ = fmt.Fprintf(w, "%-*s  %-*s  %-*s\n", nameWidth, nameCol, valueWidth, valueCol, sourceWidth, sourceCol)
		for _, row := range rows {
			_, _ = fmt.Fprintf(w, "%-*s  %-*s  %-*s\n", nameWidth, row[0], valueWidth, row[1], sourceWidth, row[2])
		}
		_, _ = fmt.Fprintln(w)
		printStatus(w, outputInfo, note)
		return
	}

	cs := defaultFangColorScheme()
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(cs.Command)
	nameStyle := lipgloss.NewStyle().Bold(true)
	valueTrueStyle := lipgloss.NewStyle().Bold(true).Foreground(cs.Flag)
	valueFalseStyle := lipgloss.NewStyle()
	sourceStyle := lipgloss.NewStyle()

	_, _ = fmt.Fprintf(
		w,
		"%s  %s  %s\n",
		headerStyle.Render(fmt.Sprintf("%-*s", nameWidth, nameCol)),
		headerStyle.Render(fmt.Sprintf("%-*s", valueWidth, valueCol)),
		headerStyle.Render(fmt.Sprintf("%-*s", sourceWidth, sourceCol)),
	)

	for _, row := range rows {
		nameCell := nameStyle.Render(fmt.Sprintf("%-*s", nameWidth, row[0]))
		valueRaw := fmt.Sprintf("%-*s", valueWidth, row[1])
		valueCell := valueFalseStyle.Render(valueRaw)
		if row[1] == "true" {
			valueCell = valueTrueStyle.Render(valueRaw)
		}
		sourceCell := sourceStyle.Render(fmt.Sprintf("%-*s", sourceWidth, row[2]))
		_, _ = fmt.Fprintf(w, "%s  %s  %s\n", nameCell, valueCell, sourceCell)
	}

	_, _ = fmt.Fprintln(w)
	printStatus(w, outputInfo, note)
}

func printStatus(w io.Writer, level outputLevel, message string) {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return
	}

	if !supportsStyling(w) {
		_, _ = fmt.Fprintf(w, "fuzmit: %s\n", msg)
		return
	}

	cs := defaultFangColorScheme()
	badge := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(cs.Codeblock).
		Background(cs.Title)
	if level == outputCommit {
		badge = badge.Background(cs.Flag)
	}
	text := lipgloss.NewStyle().Foreground(cs.Base)
	_, _ = fmt.Fprintf(w, "%s %s\n", badge.Render(strings.ToUpper(string(level))), text.Render(msg))
}

func PrintHelpNotes(w io.Writer) {
	colorize := supportsStyling(w)
	title := "NOTES"
	if colorize {
		title = lipgloss.NewStyle().Bold(true).Render(title)
	}

	_, _ = fmt.Fprintf(w, "  %s\n\n", title)
	_, _ = fmt.Fprintf(w, "    %s\n\n", helpNotesLine(colorize))
	for _, line := range helpEnvLines() {
		_, _ = fmt.Fprintf(w, "    %s\n", line)
	}
	_, _ = fmt.Fprintln(w)
}

func supportsStyling(w io.Writer) bool {
	f, ok := w.(term.File)
	return ok && term.IsTerminal(f.Fd())
}

func defaultFangColorScheme() fang.ColorScheme {
	isDark := false
	if term.IsTerminal(os.Stdout.Fd()) {
		isDark = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	}
	return fang.DefaultColorScheme(lipgloss.LightDark(isDark))
}
