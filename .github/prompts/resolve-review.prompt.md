---
name: resolveReview
description: Process and resolve code review comments from a pull request
argument-hint: Path to PR or review context
agent: Coder
---

## Task

Retrieve, analyze, and act on reviewer comments from the current Pull Request. Your objective is to produce code of exceptional quality by critically evaluating each piece of feedback, applying only the changes that genuinely improve the codebase, and respectfully declining those that do not.

## Workflow

### Step 1: Retrieve Comments

Use the terminal to fetch **all review comments** on the current Pull Request:

```bash
# Get the current PR number
PR_NUMBER=$(gh pr view --json number --jq '.number')

# Fetch all review comments (inline code comments)
gh api "repos/{owner}/{repo}/pulls/${PR_NUMBER}/comments" --paginate

# Fetch all top-level PR review bodies
gh api "repos/{owner}/{repo}/pulls/${PR_NUMBER}/reviews" --paginate

# Fetch general issue-level comments (if any)
gh pr view "$PR_NUMBER" --json comments --jq '.comments'
```

Collect every comment regardless of its resolution status.

### Step 2: Classify Each Comment

For **every** comment, determine its category:

| Category | Description | Action |
|---|---|---|
| **Valid & Actionable** | Real bug, security flaw, performance issue, readability improvement, or idiomatic best practice that aligns with the project philosophy. | **Apply the fix.** |
| **Valid but Already Addressed** | Concern was correct at the time but has since been resolved in a subsequent commit. | **Skip with explanation.** |
| **Subjective / Stylistic Preference** | Stylistic change that is neither better nor worse — merely different. Does not align with existing project conventions. | **Skip with explanation.** |
| **Incorrect or Counterproductive** | Suggestion would introduce a bug, degrade performance, violate project architecture, break conventions, or reduce code quality. | **Reject with rationale.** |
| **Outdated / Stale** | Comment references code that no longer exists in the current diff. | **Skip with explanation.** |

Before accepting or rejecting any comment: quote the relevant reviewer comment and the code it references. Then reason through: (a) what is the reviewer asking for, (b) is it technically correct, (c) does it align with the project's conventions, (d) would it improve or degrade the code quality.

### Step 3: Apply Changes

For each comment classified as **Valid & Actionable**:

1. Locate the exact file and line range referenced.
2. Implement the change with surgical precision — modify only what is necessary.
3. Ensure the fix does not introduce regressions (run `make test`).
4. If the reviewer's suggestion is directionally correct but the proposed implementation is suboptimal, implement a **better version** that addresses the underlying concern.

### Step 4: Produce Summary

```markdown
## PR Review Processing Summary

### Applied (N comments)
- [Comment by @reviewer, file:line] — Brief description of what was changed and why.

### Skipped — Already Addressed (N comments)
- [Comment by @reviewer, file:line] — Brief explanation.

### Skipped — Subjective (N comments)
- [Comment by @reviewer, file:line] — Brief explanation of the stylistic trade-off.

### Rejected (N comments)
- [Comment by @reviewer, file:line] — Detailed technical rationale.

### Stale / Outdated (N comments)
- [Comment by @reviewer, file:line] — Brief note on why it no longer applies.
```

## Guiding Principles

1. **Code quality is paramount.** Never apply a change that makes the code worse, regardless of who suggested it.
2. **Respect the project's philosophy.** Changes must be consistent with established conventions. Reject suggestions that contradict architectural patterns.
3. **Think like a maintainer, not a people-pleaser.** The goal is not to mark every comment as "resolved." The goal is to ship correct, maintainable code.
4. **Be thorough but surgical.** Apply the minimum change that fully addresses the concern.
5. **Explain every decision.** Document your reasoning for every apply, skip, or reject.
6. **Err on the side of correctness.** When in doubt, investigate deeper before deciding.

## Constraints

- You MUST use the `gh` CLI tool to retrieve comments. Do not fabricate comments.
- You MUST NOT apply changes that break existing tests or introduce type errors.
- You MUST NOT blindly follow suggestions that reduce code readability, performance, or safety.
- You MUST preserve the project's existing code style and architectural patterns.
- When rejecting a comment, your rationale must be technical and specific — never dismissive.
- If a reviewer raises a concern that you disagree with but cannot definitively disprove, flag it as **"Needs Discussion"** rather than silently rejecting it.

Apply coding standards from: [Go documentation guidelines](../instructions/go-documentation.instructions.md) and [Go environment guidelines](../instructions/go-environment.instructions.md)
