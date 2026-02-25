#!/usr/bin/env sh
set -eu

msg_file="${1:-}"
if [ -z "$msg_file" ] || [ ! -f "$msg_file" ]; then
  echo "missing commit message file" >&2
  exit 1
fi

first_line=$(sed -n '1p' "$msg_file")

case "$first_line" in
  Merge\ *|Revert\ *|fixup!\ *|squash!\ *)
    exit 0
    ;;
esac

pattern='^[a-z][a-z0-9-]*(\([A-Za-z0-9][A-Za-z0-9._/-]*\))?(!)?: .+'

if printf '%s' "$first_line" | grep -Eq "$pattern"; then
  exit 0
fi

echo "commit message must follow Conventional Commits: <type>[optional scope][!]: <description>" >&2
echo "example: feat(parser): add array parsing" >&2
echo "spec: https://www.conventionalcommits.org/en/v1.0.0/#specification" >&2
echo "got: $first_line" >&2
exit 1
