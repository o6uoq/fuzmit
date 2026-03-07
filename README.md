# fuzmit

fuzmit: Conventional Commits, but Fuzzy.

## 🎬 Demo

![fuzmit demo](./fuzmit.gif)

## 🍺 Install (Homebrew)

From [o6uoq/homebrew-tap](https://github.com/o6uoq/homebrew-tap):

```bash
brew tap o6uoq/homebrew-tap
brew install o6uoq/homebrew-tap/fuzmit
```

## ⚡ Quick Start

```bash
fuzmit --help
fuzmit --type fix --scope auth -m "handle nil panic"
fuzmit --type feat -m "add jira scope detection"
fuzmit --type feat --scope -m "prompt for optional scope"
```

## 🎌 Flags

- `-h, --help` show help for fuzmit
- `-j, --jira-scope` detect Jira scope from branch name (example: `ABC-123`)
- `-m, --message <description>` set commit message directly
- `--no-emojis` disable emojis in picker/help output
- `--override` bypass main branch protection
- `-s, --scope <scope>` set optional scope; requires `--type` (pass `--scope` without a value to prompt)
- `-t, --type <type>` commit type: `build|chore|ci|docs|feat|fix|perf|refactor|style|test`
- `-v, --version` show fuzmit version

## 🌱 Environment Variables

- `FUZMIT_SCOPE=true|false` prompt for scope if `--scope` is not provided
- `FUZMIT_JIRA_SCOPE=true|false` detect Jira scope from branch and ignore `--scope` / `FUZMIT_SCOPE`
- `FUZMIT_NO_EMOJIS=true|false` disable emojis in commit subjects, menus, and help output

Inspect resolved env settings:

```bash
fuzmit env
```
