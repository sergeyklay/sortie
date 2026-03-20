#!/usr/bin/env bash
# block-secrets.sh — Block creation of files that may contain secrets
#
# Hook type: preToolUse
# Blocks: .env, .pem, .key files and secrets/ directory
# Why always wrong: Secrets come via environment variables and WORKFLOW.md
# $VAR expansion. Secret files in a Git repo are always a mistake.
# Input: JSON with toolName, toolArgs
# Output: JSON with permissionDecision if blocked

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName')

if [ "$TOOL_NAME" != "edit" ] && [ "$TOOL_NAME" != "create" ]; then
  exit 0
fi

FILE_PATH=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.path // .file // empty')

if [ -z "$FILE_PATH" ]; then
  exit 0
fi

REASON=""

case "$FILE_PATH" in
  *.env|*/.env.*)        REASON="Environment files (.env) contain secrets" ;;
  *.pem)                 REASON="PEM files may contain private keys" ;;
  *.key)                 REASON="Key files may contain private keys" ;;
  */secrets/*|secrets/*)  REASON="Files in secrets/ directory are protected" ;;
esac

if [ -n "$REASON" ]; then
  jq -n \
    --arg reason "$REASON. Use CI/CD variables or WORKFLOW.md \$VAR expansion." \
    '{permissionDecision: "deny", permissionDecisionReason: $reason}'
fi
