# fuzmit

Conventional commits with fuzzy type selection and direct CLI flags.

`fuzmit` is a standalone binary. If you also want Git subcommand style, install the same binary as `git-fuzmit`.

## Run locally

```bash
go run . --help
go run . --type fix --scope auth -m "prevent nil panic"
go run . --scope -m "prompted scope example"
```

## Build locally

```bash
go build -o fuzmit .
./fuzmit --help
```

Optional Git-subcommand install:

```bash
cp ./fuzmit ~/.bin/git-fuzmit
git fuzmit --help
```

## Flags (A-Z)

- `--jira-scope` auto-detect Jira scope from current branch (for example `ABC-123`)
- `-m, --message <description>` set commit description directly
- `--no-emojis` disable emojis in commit-type picker/help output
- `--override` bypass main-branch protection
- `--scope <scope>` set optional scope directly
- `--scope` (no value) prompt for optional scope interactively
- `--type <type>` set commit type directly (`build|chore|ci|docs|feat|fix|perf|refactor|style|test`)

Commit subjects are always emoji-free conventional commits.

## Defaults and config

Environment overrides:

- `FUZMIT_SCOPE=true|false` default prompt for scope on each run
- `FUZMIT_JIRA_SCOPE=true|false` default auto Jira scope extraction from branch name
- `FUZMIT_NO_EMOJIS=true|false` default no emoji in picker/help list items

Persist defaults:

```bash
fuzmit scope on
fuzmit jira-scope on
fuzmit defaults
```

Config path:

- macOS: `~/Library/Application Support/fuzmit/config.json`
- Linux: `${XDG_CONFIG_HOME:-~/.config}/fuzmit/config.json`

## Jira scope detection

When Jira scope extraction is enabled, `fuzmit` scans the current branch name for an issue key and uses the first match as commit scope.

- Case-insensitive match
- Normalized to uppercase
- Expected Jira key shape: `<PROJECT_KEY>-<NUMBER>` (for example `ABCD-12345`)

## Pre-commit hooks

Install and enable both `pre-commit` and `commit-msg` hooks:

```bash
pre-commit install --hook-type pre-commit --hook-type commit-msg
pre-commit run -a
```

Configured checks:

- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
- `golangci-lint` (official upstream pre-commit hook)
- Conventional Commit message validation on `commit-msg`

## Homebrew release model

Recommended path for Homebrew users:

1. Build release binaries in GitHub Actions (not local Docker)
2. Attach tarballs + checksums to GitHub Releases
3. Update a tap formula from release artifacts (typically automated via GoReleaser)

This keeps local development lightweight (`go run`, `go build`) and pushes packaging to CI.
