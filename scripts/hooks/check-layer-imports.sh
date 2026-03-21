#!/usr/bin/env bash
# check-layer-imports.sh — Detect layer boundary violations after edit
#
# Hook type: postToolUse
# Advisory only (exit 0) — warns via stderr, never blocks.
# Checks Go import statements against the dependency graph:
#
#   domain       <- (no internal imports)
#   logging      <- (no internal imports)
#   config       <- domain, logging
#   prompt       <- domain, logging
#   registry     <- domain, logging
#   persistence  <- domain, config, logging
#   workspace    <- domain, config, persistence, logging
#   workflow     <- config, prompt, logging
#   tracker/*    <- domain, registry, logging
#   agent/*      <- domain, registry, logging
#   orchestrator <- everything above
#
# Why no false positives: import rules are deterministic. A domain package
# importing orchestrator is always wrong regardless of context.
# Input: JSON with toolName, toolArgs, toolResult
# Output: Warning to stderr if violation detected

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName')
RESULT_TYPE=$(echo "$INPUT" | jq -r '.toolResult.resultType')

if [ "$TOOL_NAME" != "edit" ] && [ "$TOOL_NAME" != "create" ]; then
  exit 0
fi

if [ "$RESULT_TYPE" != "success" ]; then
  exit 0
fi

FILE_PATH=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.path // .file // empty')

case "$FILE_PATH" in
  *.go) ;;
  *) exit 0 ;;
esac

[ -f "$FILE_PATH" ] || exit 0

PKG=$(dirname "$FILE_PATH" | sed 's|.*/internal/||')
VIOLATION=""

case "$PKG" in
  domain*)
    if grep -qE '".*sortie/internal/' "$FILE_PATH"; then
      VIOLATION="domain must not import other internal packages"
    fi
    ;;
  logging*)
    if grep -E '".*sortie/internal/' "$FILE_PATH" | grep -qvE '".*sortie/internal/logging'; then
      VIOLATION="logging must not import other internal packages"
    fi
    ;;
  config*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|persistence|tracker|registry|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="config may only import domain, logging"
    fi
    ;;
  prompt*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|persistence|tracker|registry|config|workflow|server)' "$FILE_PATH"; then
      VIOLATION="prompt may only import domain, logging"
    fi
    ;;
  registry*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|persistence|tracker|config|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="registry may only import domain, logging"
    fi
    ;;
  persistence*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|tracker|registry|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="persistence may only import domain, config, logging"
    fi
    ;;
  workspace*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|tracker|registry|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="workspace may only import domain, config, persistence, logging"
    fi
    ;;
  workflow*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|persistence|tracker|registry|domain|server)' "$FILE_PATH"; then
      VIOLATION="workflow may only import config, prompt, logging"
    fi
    ;;
  tracker/*)
    if grep -qE '".*sortie/internal/(orchestrator|agent|workspace|persistence|config|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="tracker adapters may only import domain, registry, logging"
    fi
    ;;
  agent/*)
    if grep -qE '".*sortie/internal/(orchestrator|tracker|workspace|persistence|config|prompt|workflow|server)' "$FILE_PATH"; then
      VIOLATION="agent adapters may only import domain, registry, logging"
    fi
    ;;
esac

if [ -n "$VIOLATION" ]; then
  echo "Layer violation in $FILE_PATH: $VIOLATION" >&2
fi

exit 0
