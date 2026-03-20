#!/usr/bin/env bash
# no-cgo.sh — Block CGo-dependent libraries and CGO_ENABLED=1
#
# Hook type: preToolUse
# Blocks: go get mattn/go-sqlite3, CGO_ENABLED=1, mattn/go-sqlite3 in go.mod
# Why always wrong: Sortie is a single static binary. CGo breaks this.
# Input: JSON with toolName, toolArgs
# Output: JSON with permissionDecision if blocked

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName')

if [ "$TOOL_NAME" = "bash" ] || [ "$TOOL_NAME" = "powershell" ]; then
  COMMAND=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.command // empty')

  if echo "$COMMAND" | grep -qE 'go\s+get\s+.*mattn/go-sqlite3'; then
    jq -n '{
      permissionDecision: "deny",
      permissionDecisionReason: "mattn/go-sqlite3 requires CGo. Use modernc.org/sqlite (pure Go)."
    }'
    exit 0
  fi

  if echo "$COMMAND" | grep -qE 'CGO_ENABLED=1'; then
    jq -n '{
      permissionDecision: "deny",
      permissionDecisionReason: "CGo is not allowed. Sortie must be a single static binary."
    }'
    exit 0
  fi
fi

if [ "$TOOL_NAME" = "edit" ] || [ "$TOOL_NAME" = "create" ]; then
  FILE_PATH=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.path // .file // empty')

  if [ "$(basename "$FILE_PATH" 2>/dev/null)" = "go.mod" ]; then
    CONTENT=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.new_str // .content // empty')
    if echo "$CONTENT" | grep -qE 'mattn/go-sqlite3'; then
      jq -n '{
        permissionDecision: "deny",
        permissionDecisionReason: "mattn/go-sqlite3 cannot be added to go.mod. Use modernc.org/sqlite."
      }'
      exit 0
    fi
  fi
fi
