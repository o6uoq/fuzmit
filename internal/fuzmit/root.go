package fuzmit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type runOptions struct {
	Type     string
	Scope    string
	ScopeSet bool
	AskScope bool
	GeoScope bool
	NoEmojis bool
	Override bool
	Message  string
}

const scopePromptSentinel = "__PROMPT_SCOPE__"

// NewRootCommand builds the fuzmit command tree.
func NewRootCommand() *cobra.Command {
	opts := runOptions{}

	cmd := &cobra.Command{
		Use:   "fuzmit [flags] [description]",
		Short: "Conventional Commits, but fuzzy",
		Long: `fuzmit builds conventional commits with fuzzy selection or direct flags.

Examples:
  fuzmit
  fuzmit --type fix --scope auth -m "prevent nil panic"
  fuzmit --type docs --no-emojis -m "update README usage"
  fuzmit --scope
`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			localOpts := opts
			localOpts.ScopeSet = cmd.Flags().Changed("scope")
			if localOpts.ScopeSet && strings.TrimSpace(localOpts.Scope) == scopePromptSentinel {
				localOpts.Scope = ""
				localOpts.AskScope = true
			}
			if !localOpts.GeoScope && cmd.Flags().Changed("geoscope") {
				v, _ := cmd.Flags().GetBool("geoscope")
				localOpts.GeoScope = v
			}
			return runCommit(cmd, args, localOpts)
		},
	}

	typeList := strings.Join(typeNames(), "|")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Commit type: "+typeList)
	cmd.Flags().StringVar(&opts.Scope, "scope", "", "Set optional scope (e.g. auth or ABC-123); pass --scope without a value to prompt")
	if scopeFlag := cmd.Flags().Lookup("scope"); scopeFlag != nil {
		scopeFlag.NoOptDefVal = scopePromptSentinel
	}
	cmd.Flags().BoolVar(&opts.AskScope, "prompt-scope", false, "Deprecated: use --scope without a value")
	_ = cmd.Flags().MarkHidden("prompt-scope")
	_ = cmd.Flags().MarkDeprecated("prompt-scope", "use --scope without a value instead")
	cmd.Flags().BoolVar(&opts.GeoScope, "jira-scope", false, "Auto-detect Jira scope from branch name (e.g. ABC-123)")
	cmd.Flags().Bool("geoscope", false, "Deprecated alias for --jira-scope")
	_ = cmd.Flags().MarkHidden("geoscope")
	_ = cmd.Flags().MarkDeprecated("geoscope", "use --jira-scope instead")
	cmd.Flags().BoolVar(&opts.NoEmojis, "no-emojis", false, "Disable emojis in commit-type menus/help (commit subjects are emoji-free)")
	cmd.Flags().BoolVar(&opts.Override, "override", false, "Bypass main branch protection")
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Commit description")

	cmd.AddCommand(newScopeCommand())
	cmd.AddCommand(newJiraScopeCommand())
	cmd.AddCommand(newDefaultsCommand())
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		_ = PrintHelp(c.OutOrStdout(), resolveNoEmojiHelp(c))
	})
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return c.Help()
	})

	return cmd
}

func runCommit(cmd *cobra.Command, args []string, opts runOptions) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	defaults := ResolveDefaults(cfg, os.Getenv)

	noEmojis := defaults.NoEmojis || opts.NoEmojis
	geoScope := defaults.GeoScope || opts.GeoScope
	scopeEnabled := defaults.Scope || opts.AskScope

	if err := EnsureGitRepo(); err != nil {
		return err
	}

	branch, err := CurrentBranch()
	if err != nil {
		return err
	}
	if branch == "main" && !opts.Override {
		return errors.New("fuzmit: you are on the main branch; use --override to bypass this check")
	}

	hasStaged, err := HasStagedChanges()
	if err != nil {
		return err
	}
	if !hasStaged {
		cmd.Println("ℹ️  fuzmit: No staged changes to commit.")
		return nil
	}

	commitType, err := resolveCommitType(opts.Type, noEmojis)
	if err != nil {
		return err
	}

	scope, err := resolveScope(cmd, opts, branch, scopeEnabled, geoScope)
	if err != nil {
		return err
	}

	description, err := resolveDescription(cmd, args, opts.Message)
	if err != nil {
		return err
	}

	full := BuildCommitMessage(commitType.Name, scope, description)
	if err := ValidateConventionalSubject(full); err != nil {
		return fmt.Errorf("fuzmit: %w", err)
	}

	cmd.Printf("💾  fuzmit: Commit message - '%s'\n", full)

	output, err := Commit(full)
	if output != "" {
		cmd.Println(output)
	}
	if err != nil {
		return fmt.Errorf("fuzmit: git commit failed: %w", err)
	}
	return nil
}

func resolveCommitType(typeFlag string, noEmojis bool) (CommitType, error) {
	if typeFlag != "" {
		ct, ok := FindCommitType(typeFlag)
		if !ok {
			return CommitType{}, fmt.Errorf("fuzmit: invalid --type %q", typeFlag)
		}
		return ct, nil
	}

	ct, err := SelectCommitType(noEmojis)
	if err != nil {
		if errors.Is(err, errSelectionAborted) {
			return CommitType{}, errors.New("fuzmit: no commit type selected, aborting")
		}
		return CommitType{}, fmt.Errorf("fuzmit: unable to select commit type: %w", err)
	}
	return ct, nil
}

func resolveScope(cmd *cobra.Command, opts runOptions, branch string, scopeEnabled bool, geoScope bool) (string, error) {
	if geoScope {
		scope := ExtractJiraScope(branch)
		if scope != "" {
			if err := ValidateScope(scope); err != nil {
				return "", fmt.Errorf("fuzmit: invalid extracted Jira scope %q: %w", scope, err)
			}
			cmd.Printf("ℹ️  fuzmit: Auto-detected Jira scope '%s'\n", scope)
		}
		return scope, nil
	}

	if opts.ScopeSet {
		if opts.AskScope {
			return promptScope(cmd)
		}
		return normalizeScope(opts.Scope)
	}

	if !scopeEnabled {
		return "", nil
	}

	return promptScope(cmd)
}

func promptScope(cmd *cobra.Command) (string, error) {
	scope, err := PromptLine(cmd.InOrStdin(), cmd.OutOrStdout(), "Enter optional scope (leave empty if not needed): ")
	if err != nil {
		return "", fmt.Errorf("fuzmit: failed to read scope: %w", err)
	}
	return normalizeScope(scope)
}

func normalizeScope(scope string) (string, error) {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return "", nil
	}
	scope = strings.Trim(scope, "()")
	if err := ValidateScope(scope); err != nil {
		return "", fmt.Errorf("invalid --scope value %q: %w", scope, err)
	}
	return scope, nil
}

func resolveDescription(cmd *cobra.Command, args []string, messageFlag string) (string, error) {
	description := strings.TrimSpace(messageFlag)
	if description == "" && len(args) > 0 {
		description = strings.TrimSpace(strings.Join(args, " "))
	}
	if description == "" {
		var err error
		description, err = PromptLine(cmd.InOrStdin(), cmd.OutOrStdout(), "Enter commit description: ")
		if err != nil {
			return "", fmt.Errorf("fuzmit: failed to read commit description: %w", err)
		}
	}
	if description == "" {
		return "", errors.New("fuzmit: commit description cannot be empty, aborting")
	}
	return description, nil
}

func typeNames() []string {
	out := make([]string, 0, len(SupportedTypes))
	for _, ct := range SupportedTypes {
		out = append(out, ct.Name)
	}
	return out
}

func newScopeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "scope [on|off]",
		Short: "Get or set default scope prompting",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := LoadConfig()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				cmd.Printf("scope default: %t\n", cfg.Scope)
				return nil
			}

			v, err := ParseToggleArg(args[0])
			if err != nil {
				return err
			}
			cfg.Scope = v
			if err := SaveConfig(cfg); err != nil {
				return err
			}
			cmd.Printf("scope default set to %t\n", v)
			return nil
		},
	}
}

func newJiraScopeCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "jira-scope [on|off]",
		Aliases: []string{"geoscope"},
		Short:   "Get or set default Jira branch scope extraction",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := LoadConfig()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				cmd.Printf("jira-scope default: %t\n", cfg.GeoScope)
				return nil
			}

			v, err := ParseToggleArg(args[0])
			if err != nil {
				return err
			}
			cfg.GeoScope = v
			if err := SaveConfig(cfg); err != nil {
				return err
			}
			cmd.Printf("jira-scope default set to %t\n", v)
			return nil
		},
	}
}

func newDefaultsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "defaults",
		Short: "Show resolved defaults (config + env)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := LoadConfig()
			if err != nil {
				return err
			}
			d := ResolveDefaults(cfg, os.Getenv)
			path, err := ConfigPath()
			if err != nil {
				return err
			}
			cmd.Printf("config: %s\n", path)
			cmd.Printf("scope: %t\n", d.Scope)
			cmd.Printf("jira-scope: %t\n", d.GeoScope)
			cmd.Printf("no-emojis: %t\n", d.NoEmojis)
			return nil
		},
	}
}

func resolveNoEmojiHelp(cmd *cobra.Command) bool {
	cfg, err := LoadConfig()
	if err != nil {
		return false
	}
	noEmojis := ResolveDefaults(cfg, os.Getenv).NoEmojis
	if v, err := cmd.Flags().GetBool("no-emojis"); err == nil && v {
		noEmojis = true
	}
	return noEmojis
}

func PrintHelp(w io.Writer, noEmojis bool) error {
	c := func(color, text string) string { return color + text + "\033[0m" }
	writeLine := func(s string) {
		_, _ = fmt.Fprintln(w, s)
	}
	writeBlank := func() {
		_, _ = fmt.Fprintln(w)
	}
	writef := func(format string, a ...any) {
		_, _ = fmt.Fprintf(w, format, a...)
	}

	brightCyan := "\033[1;38;5;81m"
	white := "\033[1;37m"
	green := "\033[1;32m"
	purple := "\033[1;35m"
	blue := "\033[1;34m"
	cyan := "\033[1;36m"

	writeLine(c(brightCyan, "fuzmit: Conventional Commits, but Fuzzy."))
	writeBlank()
	writeLine(c(white, "# Conventional Commit Message Format:"))
	writef("%s%s%s: %s\n", c(green, "<type>"), c(blue, "(optional scope)"), c(purple, "!"), c(cyan, "<description>"))
	writeBlank()
	writeLine(c(white, "Commit Types:"))
	for _, ct := range SupportedTypes {
		if noEmojis {
			writef("    %s - %s\n", ct.Name, ct.Description)
		} else {
			writef("    %s  %s - %s\n", ct.Emoji, ct.Name, ct.Description)
		}
	}
	writeBlank()
	writeLine(c(white, "# Example Commit Message:"))
	writeLine(c(cyan, "feat(authentication): add support for password reset"))
	writeBlank()
	writeLine(c(white, "# Usage:"))
	writeLine("  fuzmit [--type <type>] [--scope <scope>|--scope] [--jira-scope] [--override] [-m <message>]")
	writeLine("  fuzmit defaults")
	writeLine("  fuzmit scope [on|off]")
	writeLine("  fuzmit jira-scope [on|off]")
	writeBlank()
	writeLine(c(white, "# Options:"))
	writeLine("  --jira-scope     Auto-detect Jira scope from branch name (e.g. ABC-123)")
	writeLine("  -m, --message    Set commit description directly")
	writeLine("  --no-emojis      Disable emojis in commit-type menus/help")
	writeLine("  --override       Bypass main branch restriction and allow committing")
	writeLine("  --scope <name>   Set optional scope directly")
	writeLine("  --scope          Prompt for optional scope")
	writeLine("  --type <type>    Set commit type directly (build|chore|ci|docs|feat|fix|perf|refactor|style|test)")
	writeBlank()
	writeLine(c(white, "# Notes:"))
	writeLine("  - Commit subjects are emoji-free conventional commits")
	writeLine("  - If FUZMIT_JIRA_SCOPE=true, scope prompt/defaults are ignored and Jira scope is auto-detected")
	writeLine("  - No staged changes exits cleanly without committing")
	writeBlank()
	writeLine(c(cyan, "Conventional Commits: https://www.conventionalcommits.org/en/v1.0.0/#specification"))
	writeBlank()
	return nil
}
