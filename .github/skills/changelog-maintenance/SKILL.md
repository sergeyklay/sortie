---
name: changelog-maintenance
description: >
  Use when asked to update the changelog, document version changes, prepare
  a release, or add entries for recent work. Handles CHANGELOG.md updates
  following Keep a Changelog format and Semantic Versioning. Do NOT use for
  committing (use git-commit) or creating release notes outside CHANGELOG.md.
---

# Changelog Maintenance

## When to use

- Adding entries for new features, fixes, or breaking changes.
- Preparing a release: moving Unreleased entries under a versioned heading.
- Creating CHANGELOG.md from scratch when it does not exist.

## Workflow

### Step 1: Read the current changelog

```bash
cat CHANGELOG.md
```

If the file does not exist, create it with the preamble from Step 4.

### Step 2: Gather changes

Determine what changed since the last release. Use the sources that fit the
situation — not all are needed every time.

```bash
# Changes since last tag
git log --oneline "$(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~20)"..HEAD

# Or changes in the Unreleased section (already documented)
head -60 CHANGELOG.md
```

If the user describes changes verbally, use that as the primary source.

### Step 3: Classify each change

Place every entry under exactly one category. Use this decision order:

| Category       | When to use                                            |
| -------------- | ------------------------------------------------------ |
| **Added**      | New feature, new file, new dependency, new CLI command  |
| **Changed**    | Existing behavior altered, refactored, or restructured  |
| **Deprecated** | Still works but scheduled for removal                   |
| **Removed**    | Deleted feature, file, or dependency                    |
| **Fixed**      | Bug fix, corrected behavior                             |
| **Security**   | Vulnerability patch, dependency CVE fix                 |

Rules:
- One bullet per logical change. Combine related sub-changes into one bullet.
- Start each bullet with what changed, not with "Fixed" or "Added" (the heading
  already says that).
- Be specific: "`coroutine 'main' was never awaited` bug after async migration"
  not "Fixed async bug".
- Reference function/class names in backticks when they help the reader locate
  the change.
- Do not copy git commit messages verbatim — rewrite for a human reader.
- Omit routine dependency bumps unless they are significant (major version,
  security fix, or the user explicitly asks).

### Step 4: Write the entry

Format (Keep a Changelog 1.1.0):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Description of new capability.

## [X.Y.Z] - YYYY-MM-DD

### Fixed

- Description of what was broken and how it was fixed.

[Unreleased]: https://github.com/OWNER/REPO/compare/vX.Y.Z...HEAD
[X.Y.Z]: https://github.com/OWNER/REPO/compare/vA.B.C...vX.Y.Z
```

Structural rules:
- Reverse chronological order (newest first).
- `[Unreleased]` section always present at the top.
- Dates in ISO 8601 (`YYYY-MM-DD`).
- Comparison links at the bottom for every version.
- Empty categories are omitted (no `### Removed` if nothing was removed).

### Step 5: Determine the version bump

When cutting a release, choose the version number:

| Bump      | Trigger                                               |
| --------- | ----------------------------------------------------- |
| **Major** | Breaking API/CLI change, removed public functionality |
| **Minor** | New feature, backward-compatible behavior change      |
| **Patch** | Bug fix, security patch, docs-only (if versioned)     |

To cut a release:
1. Replace `## [Unreleased]` with `## [X.Y.Z] - YYYY-MM-DD`.
2. Add a fresh empty `## [Unreleased]` section above it.
3. Update the comparison links at the bottom.

### Step 6: Verify

- [ ] Newest version is at the top.
- [ ] Every version has a date (except Unreleased).
- [ ] Bottom links are correct and complete.
- [ ] No empty category headings.
- [ ] No git-log copy-paste — entries are human-readable.

## Error Recovery

| Problem                    | Fix                                                        |
| -------------------------- | ---------------------------------------------------------- |
| Missing comparison links   | Reconstruct from `git tag --sort=-version:refname`         |
| Duplicate entries           | Deduplicate, keep the more descriptive version             |
| Entry under wrong category | Move it; if ambiguous, prefer Changed over Added           |
| No tags in repository      | Use commit SHAs in comparison links as a temporary measure |
