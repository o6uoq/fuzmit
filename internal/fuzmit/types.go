package fuzmit

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// CommitType describes a conventional commit type.
type CommitType struct {
	Name        string
	Emoji       string
	Description string
}

// SupportedTypes is the canonical supported commit type list.
var SupportedTypes = []CommitType{
	{Name: "build", Emoji: "🏗", Description: "Project build or dependencies changes"},
	{Name: "chore", Emoji: "🧹", Description: "Routine tasks and maintenance"},
	{Name: "ci", Emoji: "🔄", Description: "Continuous integration changes"},
	{Name: "docs", Emoji: "📚", Description: "Documentation updates"},
	{Name: "feat", Emoji: "🚀", Description: "New feature addition"},
	{Name: "fix", Emoji: "🔧", Description: "Bug fixes"},
	{Name: "perf", Emoji: "🏎", Description: "Performance improvements"},
	{Name: "refactor", Emoji: "🛠", Description: "Code refactoring without behavior change"},
	{Name: "style", Emoji: "🎨", Description: "Code style or formatting changes"},
	{Name: "test", Emoji: "🧪", Description: "Adding or modifying tests"},
}

// FindCommitType resolves a type name to a supported commit type.
func FindCommitType(name string) (CommitType, bool) {
	n := strings.ToLower(strings.TrimSpace(name))
	for _, ct := range SupportedTypes {
		if ct.Name == n {
			return ct, true
		}
	}
	return CommitType{}, false
}

// FormatTypeLabel renders a commit type label for interactive selection.
func FormatTypeLabel(ct CommitType, noEmojis bool) string {
	maxNameWidth := maxCommitTypeNameWidth()
	if noEmojis {
		return fmt.Sprintf("%-*s - %s", maxNameWidth, ct.Name, ct.Description)
	}

	emojiWidth := maxCommitTypeEmojiWidth()
	emojiPad := emojiWidth - runewidth.StringWidth(ct.Emoji)
	if emojiPad < 0 {
		emojiPad = 0
	}
	emojiCell := ct.Emoji + strings.Repeat(" ", emojiPad)
	return fmt.Sprintf("%s %-*s - %s", emojiCell, maxNameWidth, ct.Name, ct.Description)
}

func maxCommitTypeNameWidth() int {
	maxWidth := 0
	for _, supported := range SupportedTypes {
		if l := len(supported.Name); l > maxWidth {
			maxWidth = l
		}
	}
	return maxWidth
}

func maxCommitTypeEmojiWidth() int {
	maxWidth := 0
	for _, supported := range SupportedTypes {
		if w := runewidth.StringWidth(supported.Emoji); w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}
