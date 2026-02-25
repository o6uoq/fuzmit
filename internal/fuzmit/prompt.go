package fuzmit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
)

var errSelectionAborted = errors.New("selection aborted")

// SelectCommitType runs the interactive type picker.
func SelectCommitType(noEmojis bool) (CommitType, error) {
	types := pickerTypes()
	var selectedType string
	options := make([]huh.Option[string], 0, len(types))
	for _, ct := range types {
		options = append(options, huh.NewOption(FormatTypeLabel(ct, noEmojis), ct.Name))
	}

	err := huh.NewSelect[string]().
		Title("Pick commit type").
		Description("Press / to filter, Enter to select").
		Options(options...).
		Value(&selectedType).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return CommitType{}, errSelectionAborted
		}
		return CommitType{}, err
	}

	ct, ok := FindCommitType(selectedType)
	if !ok {
		return CommitType{}, fmt.Errorf("fuzmit: selected unknown commit type %q", selectedType)
	}
	return ct, nil
}

// pickerTypes returns commit types in deterministic alphabetical order.
func pickerTypes() []CommitType {
	types := make([]CommitType, len(SupportedTypes))
	copy(types, SupportedTypes)
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})
	return types
}

// PromptLine asks for a single-line text input.
func PromptLine(in io.Reader, out io.Writer, prompt string) (string, error) {
	if _, err := fmt.Fprint(out, prompt); err != nil {
		return "", err
	}
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	line = strings.TrimSpace(line)
	if errors.Is(err, io.EOF) && line == "" {
		return "", io.EOF
	}
	return line, nil
}
