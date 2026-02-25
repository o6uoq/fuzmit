package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/o6uoq/fuzmit/internal/fuzmit"
)

var version = "dev"

func main() {
	args := os.Args[1:]
	noEmojis := noEmojiRequested(args, os.Getenv("FUZMIT_NO_EMOJIS"))
	if rootHelpRequested(args) {
		_ = fuzmit.PrintHelp(os.Stdout, noEmojis)
		return
	}

	cmd := fuzmit.NewRootCommand()
	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
		fang.WithErrorHandler(errorHandler()),
	); err != nil {
		os.Exit(1)
	}
}

func rootHelpRequested(args []string) bool {
	if len(args) == 0 {
		return false
	}
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return arg == "help"
	}
	return false
}

func errorHandler() fang.ErrorHandler {
	return func(w io.Writer, _ fang.Styles, err error) {
		msg := strings.TrimSpace(err.Error())
		if msg == "" {
			return
		}
		if strings.HasPrefix(strings.ToLower(msg), "fuzmit:") {
			_, _ = fmt.Fprintf(w, "❌ %s\n\n", msg)
			return
		}
		_, _ = fmt.Fprintf(w, "❌ fuzmit: %s\n\n", msg)
	}
}

func noEmojiRequested(args []string, env string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--no-emojis" || strings.HasPrefix(arg, "--no-emojis=") {
			return true
		}
	}
	env = strings.TrimSpace(strings.ToLower(env))
	switch env {
	case "1", "true", "yes", "on", "y":
		return true
	default:
		return false
	}
}
