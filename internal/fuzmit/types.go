package fuzmit

import "strings"

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
	{Name: "refactor", Emoji: "🛠️", Description: "Code refactoring without behavior change"},
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
	if noEmojis {
		return ct.Name + " - " + ct.Description
	}
	return ct.Emoji + " " + ct.Name + " - " + ct.Description
}
