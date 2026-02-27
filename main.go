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
	args = normalizeArgs(args)
	helpMode := helpRequested(args)
	cmd := fuzmit.NewRootCommand()
	cmd.SetArgs(args)
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

func normalizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			out = append(out, args[i:]...)
			break
		}
		if arg == "--scope" || arg == "-s" {
			if i+1 < len(args) {
				next := args[i+1]
				if next != "--" && !strings.HasPrefix(next, "-") {
					out = append(out, arg+"="+next)
					i++
					continue
				}
			}
		}
		out = append(out, arg)
	}
	return out
}
