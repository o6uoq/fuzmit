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
	nameCol := "VARIABLE"
	valueCol := "VALUE"
	sourceCol := "SOURCE"

	nameWidth := len(nameCol)
	valueWidth := len(valueCol)
	for _, setting := range settings {
		if len(setting.Name) > nameWidth {
			nameWidth = len(setting.Name)
		}
	}

	note := "FUZMIT_JIRA_SCOPE=true ignores --scope and FUZMIT_SCOPE."
	if !supportsStyling(w) {
		_, _ = fmt.Fprintf(w, "%-*s  %-*s  %s\n", nameWidth, nameCol, valueWidth, valueCol, sourceCol)
		for _, setting := range settings {
			_, _ = fmt.Fprintf(w, "%-*s  %-*t  %s\n", nameWidth, setting.Name, valueWidth, setting.Value, envSettingSource(setting))
		}
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintf(w, "NOTE  %s\n", note)
		return
	}

	cs := defaultFangColorScheme()
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(cs.Codeblock).
		Background(cs.Title)
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(cs.Command)
	valueTrueStyle := lipgloss.NewStyle().Bold(true).Foreground(cs.Flag)
	valueFalseStyle := lipgloss.NewStyle().Foreground(cs.Base)
	sourceStyle := lipgloss.NewStyle().Foreground(cs.Base)
	noteLabelStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(cs.Codeblock).
		Background(cs.Title)
	noteTextStyle := lipgloss.NewStyle().Foreground(cs.Base)

	_, _ = fmt.Fprintf(
		w,
		"%s %s %s\n",
		headerStyle.Render(fmt.Sprintf("%-*s", nameWidth, nameCol)),
		headerStyle.Render(fmt.Sprintf("%-*s", valueWidth, valueCol)),
		headerStyle.Render(sourceCol),
	)

	for _, setting := range settings {
		nameCell := nameStyle.Render(fmt.Sprintf("%-*s", nameWidth, setting.Name))
		valueRaw := fmt.Sprintf("%-*t", valueWidth, setting.Value)
		valueCell := valueFalseStyle.Render(valueRaw)
		if setting.Value {
			valueCell = valueTrueStyle.Render(valueRaw)
		}
		sourceCell := sourceStyle.Render(envSettingSource(setting))
		_, _ = fmt.Fprintf(w, "%s  %s  %s\n", nameCell, valueCell, sourceCell)
	}

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "%s  %s\n", noteLabelStyle.Render("NOTE"), noteTextStyle.Render(note))
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
