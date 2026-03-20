---
name: changelog-maintenance
description: >
  Use when asked to update the changelog, document version changes, prepare
  a release, or add entries for recent work. Handles CHANGELOG.md updates
  following Keep a Changelog format and Semantic Versioning. Do NOT use for
  committing (use git-commit) or creating release notes outside CHANGELOG.md.
---

# Changelog Maintenance

Sortie's changelog speaks to operators and integrators who deploy and configure
the service. Every entry must answer: "Does this change affect someone who
installs, upgrades, configures, or integrates with Sortie?" If not, omit it.

Authoritative references:
- [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
- [Common Changelog](https://common-changelog.org/)
- [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)

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

### Step 3: Filter — decide what belongs

The changelog records **notable changes to the distributed software**. A change
is notable when it alters what a consumer of Sortie can observe: new
capabilities, changed behavior, fixed bugs, security patches, removed features,
or deprecation notices.

Apply the following filter to every commit or change before writing an entry.

**ALWAYS include:**

| Signal | Why it matters to consumers |
| --- | --- |
| New user-facing feature (CLI flag, adapter, config option) | Operators discover new capabilities |
| Changed behavior of existing feature | Operators must adjust usage |
| Bug fix for incorrect behavior | Operators know issues are resolved |
| Security or vulnerability fix | Operators must act on upgrades |
| Deprecation of public interface | Operators prepare for removal |
| Removal of feature or public interface | Operators must adapt before upgrading |
| Performance improvement with measurable impact | Operators benefit from upgrading |
| New or changed persistence schema (migration) | Operators plan upgrade procedures |

**NEVER include — these are noise, not signal:**

| Noise | Why it does not belong |
| --- | --- |
| Internal variable/function/type renames | No observable effect on consumers |
| Code formatting, whitespace, linting fixes | No observable effect on consumers |
| Test-only changes (new tests, test refactors) | Not shipped to consumers |
| CI/CD pipeline changes (workflows, actions) | Not shipped to consumers |
| Dotfile changes (`.gitignore`, `.github/*`, `CODEOWNERS`) | Not shipped to consumers |
| Documentation-only changes (README, AGENTS.md, comments) | Not shipped to consumers |
| Merge commits | Infrastructure artifact, not a change |
| Internal refactoring with no behavior change | No observable effect on consumers |
| Dev-only dependency bumps | Not shipped to consumers |
| Project scaffolding and repo housekeeping | Not shipped to consumers |

**Edge cases — include only when the threshold is met:**

| Change | Include when... | Omit when... |
| --- | --- | --- |
| Dependency bump | Major version, security fix, or changed behavior | Routine patch/minor with no user impact |
| Refactoring | It changes observable performance, error messages, or log output | Purely internal restructuring |
| New internal module/package | It introduces a new adapter or public API surface | It reorganizes existing code |
| ADR or architecture doc update | It records a decision that changes system behavior | It clarifies existing behavior |

When in doubt, ask: "If I were an operator reading this before upgrading, would
I need to know this?" If the answer is no, leave it out.

### Step 4: Classify each change

Place every surviving entry under exactly one category:

| Category | When to use |
| --- | --- |
| **Added** | New user-facing capability: CLI command, adapter, config option, API surface |
| **Changed** | Existing behavior altered in a way consumers can observe |
| **Deprecated** | Still works but scheduled for removal in a future version |
| **Removed** | Previously available feature or interface deleted |
| **Fixed** | Bug fix — incorrect behavior corrected |
| **Security** | Vulnerability patch, dependency CVE fix |

Writing rules:
- One bullet per logical change. Combine related sub-changes into one bullet.
- Start each bullet with what changed, not with "Fixed" or "Added" (the heading
  already says that).
- Be specific: "`coroutine 'main' was never awaited` bug after async migration"
  not "Fixed async bug".
- Identify the subsystem when it helps locate the change: "Jira adapter:",
  "SQLite store:", "CLI:".
- Reference types or functions in backticks when they help the reader.
- Do not copy git commit messages verbatim — rewrite for a human reader.

### Step 5: Write the entry

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

### Step 6: Determine the version bump

When cutting a release, choose the version number:

| Bump | Trigger |
| --- | --- |
| **Major** | Breaking API/CLI change, removed public functionality |
| **Minor** | New feature, backward-compatible behavior change |
| **Patch** | Bug fix, security patch |

To cut a release:
1. Replace `## [Unreleased]` with `## [X.Y.Z] - YYYY-MM-DD`.
2. Add a fresh empty `## [Unreleased]` section above it.
3. Update the comparison links at the bottom.

### Step 7: Verify

- [ ] Every entry passes the filter from Step 3 (no noise).
- [ ] Newest version is at the top.
- [ ] Every version has a date (except Unreleased).
- [ ] Bottom links are correct and complete.
- [ ] No empty category headings.
- [ ] No git-log copy-paste — entries are human-readable.
- [ ] Entries identify the subsystem where helpful.

## Error Recovery

| Problem | Fix |
| --- | --- |
| Missing comparison links | Reconstruct from `git tag --sort=-version:refname` |
| Duplicate entries | Deduplicate, keep the more descriptive version |
| Entry under wrong category | Move it; if ambiguous, prefer Changed over Added |
| No tags in repository | Use commit SHAs in comparison links as a temporary measure |
| Noise entry slipped in | Remove it — a leaner changelog is more trustworthy |
