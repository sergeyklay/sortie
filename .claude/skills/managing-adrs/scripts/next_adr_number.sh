#!/usr/bin/env bash
# Outputs the next available ADR number(s) in docs/decisions/.
#
# Usage:
#   next_adr_number.sh              # prints one number, e.g. 0004
#   next_adr_number.sh --count 3    # prints three: 0004, 0005, 0006 (one per line)
#
# Exit codes:
#   0  success
#   1  docs/decisions/ not found

set -euo pipefail

DECISIONS_DIR="docs/decisions"
COUNT=1

while [[ $# -gt 0 ]]; do
  case "$1" in
    --count) COUNT="$2"; shift 2 ;;
    *) echo "Unknown argument: $1" >&2; exit 1 ;;
  esac
done

if [[ ! -d "$DECISIONS_DIR" ]]; then
  echo "Error: $DECISIONS_DIR not found. Run from repository root." >&2
  exit 1
fi

# Find the highest existing ADR number
max=0
for f in "$DECISIONS_DIR"/[0-9][0-9][0-9][0-9]-*.md; do
  [[ -e "$f" ]] || continue
  num=$(basename "$f" | grep -oE '^[0-9]+' | sed 's/^0*//' )
  [[ -z "$num" ]] && num=0
  (( num > max )) && max=$num
done

# Output sequential numbers
for (( i = 1; i <= COUNT; i++ )); do
  printf "%04d\n" $(( max + i ))
done
