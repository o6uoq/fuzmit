package fuzmit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

var errSelectionAborted = errors.New("selection aborted")

// SelectCommitType runs the interactive fuzzy type picker.
func SelectCommitType(noEmojis bool) (CommitType, error) {
	idx, err := fuzzyfinder.Find(
		SupportedTypes,
		func(i int) string {
			return FormatTypeLabel(SupportedTypes[i], noEmojis)
		},
		fuzzyfinder.WithPromptString("Pick commit type > "),
	)
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return CommitType{}, errSelectionAborted
		}
		return CommitType{}, err
	}
	return SupportedTypes[idx], nil
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
