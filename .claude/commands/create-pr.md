---
name: create-pr
description: Commit changes and create or update a pull request
argument-hint: PR details or branch/commit context
agent: agent
tools:
  - "execute/getTerminalOutput"
  - "execute/runInTerminal"
  - "read/terminalSelection"
  - "read/terminalLastCommand"
  - "read/readFile"
  - "search"
  - "web/githubRepo"
---

Commit staged changes and manage pull requests (PR).

Task:

- Use specific skills to create a branch, commit the changes, and open/change a PR with a meaningful title and description
- Incorporate provided details or context about the changes
- Detect whether you need to create a new PR or update an existing one based on context
- When updating, verify the PR description still accurately reflects the changes
- Use conventional commit messages when appropriate
