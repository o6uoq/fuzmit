package main

import (
	"context"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/o6uoq/fuzmit/internal/fuzmit"
)

var version = "dev"

func main() {
	args := os.Args[1:]
	helpMode := helpRequested(args)
	cmd := fuzmit.NewRootCommand()
	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
	); err != nil {
		os.Exit(1)
	}
	if helpMode {
		fuzmit.PrintHelpNotes(os.Stdout)
	}
}

func helpRequested(args []string) bool {
	if len(args) > 0 && args[0] == "help" {
		return true
	}

	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-h" || arg == "--help" || strings.HasPrefix(arg, "--help=") {
			return true
		}
	}
	return false
}
