# fuzmit

Conventional commits with fuzzy type selection and direct CLI flags.

`fuzmit` is a standalone binary. If you also want Git subcommand style, install the same binary as `git-fuzmit`.

## Run locally

```bash
go run . --help
go run . --type fix --scope auth -m "fix nil panic"
go run . --type feat --scope -m "prompted scope example"
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
- `--scope <scope>` set optional scope directly (requires `--type`)
- `--scope` (no value) prompt for optional scope interactively (requires `--type`)
- `--type <type>` set commit type directly (`build|chore|ci|docs|feat|fix|perf|refactor|style|test`)

Commit subjects are always emoji-free conventional commits.

## Environment defaults

`fuzmit` is env-driven. If a variable is unset or invalid, it defaults to `false`.

- `FUZMIT_SCOPE=true|false` prompt for scope by default when `--scope` is not provided
- `FUZMIT_JIRA_SCOPE=true|false` auto-detect Jira scope from branch name; when true, both `--scope` and `FUZMIT_SCOPE` are ignored
- `FUZMIT_NO_EMOJIS=true|false` disable emojis in picker/help output

Inspect current resolved env settings:

```bash
fuzmit env
```

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

This repo uses tag-driven releases in GitHub Actions:

1. Push a tag matching `v*`
2. GitHub Actions runs tests and GoReleaser
3. GoReleaser publishes release artifacts + checksums to GitHub Releases
4. GoReleaser opens/updates a PR in `o6uoq/homebrew-tap` for `Formula/fuzmit.rb`

Required repository secret in `o6uoq/fuzmit`:

- `HOMEBREW_TAP_GITHUB_TOKEN` (PAT with write access to `o6uoq/homebrew-tap`)

Release workflow file:

- `.github/workflows/release.yml`

GoReleaser config:

- `.goreleaser.yaml`

Example prerelease (before `v0.1.0`):

```bash
git tag v0.0.1-beta.1
git push origin v0.0.1-beta.1
```

Example stable release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Install from tap:

```bash
brew tap o6uoq/homebrew-tap
brew install o6uoq/homebrew-tap/fuzmit
```
