#!/usr/bin/env bash
# makefile-only.sh — Enforce Makefile targets instead of raw go commands
#
# Hook type: preToolUse
# Blocks: go test, go build, go vet, go run
# Allows: make *, go mod *, go get * (dependency management)
# Input: JSON with toolName, toolArgs
# Output: JSON with permissionDecision if blocked
#
# Why always wrong: Makefile targets include -race, correct linter config,
# and asdf-managed Go version. Direct go commands bypass all of these.

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName')

if [ "$TOOL_NAME" != "bash" ] && [ "$TOOL_NAME" != "powershell" ]; then
  exit 0
fi

COMMAND=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.command // empty')

if [ -z "$COMMAND" ]; then
  exit 0
fi

if echo "$COMMAND" | grep -qE '(^|\s|&&|\|)\s*go\s+(test|build|vet|run)(\s|$)'; then
  jq -n '{
    permissionDecision: "deny",
    permissionDecisionReason: "Use Makefile targets: make test, make build, make lint, make fmt. Consult with .github/instructions/go-environment.instructions.md for setup."
  }'
fi
