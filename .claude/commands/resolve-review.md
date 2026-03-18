---
name: resolve-review
description: Process and resolve code review comments
argument-hint: Instructions for processing PR review comments
agent: Coder
---

## Task

You will retrieve, analyze, and act on reviewer comments from the current Pull Request using the `gh` CLI tool. Your objective is to produce code of exceptional quality by critically evaluating each piece of feedback, applying only the changes that genuinely improve the codebase, and respectfully declining those that do not.

---

## Workflow

<step_1_retrieve_comments>
Use the `gh` CLI to fetch **all review comments** on the current Pull Request:

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
</step_1_retrieve_comments>

<step_2_classify_each_comment>
For **every** comment, determine its category by reasoning through the following taxonomy:

| Category                              | Description                                                                                                                                                                             | Action                     |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| **Valid & Actionable**                | The comment identifies a real bug, security flaw, performance issue, readability improvement, or idiomatic best practice that aligns with the project philosophy.                       | **Apply the fix.**         |
| **Valid but Already Addressed**       | The concern was correct at the time of writing but has since been resolved in a subsequent commit or is no longer applicable to the current diff.                                       | **Skip with explanation.** |
| **Subjective / Stylistic Preference** | The comment suggests a stylistic change that is neither better nor worse — merely different. It does not align with the existing project conventions or adds no measurable improvement. | **Skip with explanation.** |
| **Incorrect or Counterproductive**    | The suggestion would introduce a bug, degrade performance, violate project architecture, break existing conventions, reduce type safety, or otherwise lower code quality.               | **Reject with rationale.** |
| **Outdated / Stale**                  | The comment references code that no longer exists in the current diff, or the surrounding context has changed enough that the comment is no longer relevant.                            | **Skip with explanation.** |

Think step by step. Before accepting or rejecting any comment, first quote the relevant reviewer comment and the code it references. Then reason through: (a) what is the reviewer asking for, (b) is it technically correct, (c) does it align with the project's conventions and philosophy, (d) would it improve or degrade the code quality.
</step_2_classify_each_comment>

<step_3_apply_changes>
For each comment classified as **Valid & Actionable**:

1. Locate the exact file and line range referenced by the comment.
2. Implement the change with surgical precision — modify only what is necessary.
3. Ensure the fix does not introduce regressions (run tests if available).
4. If the reviewer's suggestion is directionally correct but the proposed implementation is suboptimal, implement a **better version** that addresses the underlying concern while maintaining higher code quality.
   </step_3_apply_changes>

<step_4_produce_summary>
After processing all comments, output a structured summary:

```
## PR Review Processing Summary

### Applied (N comments)
- [Comment by @reviewer, file:line] — Brief description of what was changed and why.

### Skipped — Already Addressed (N comments)
- [Comment by @reviewer, file:line] — Brief explanation of why it's already resolved.

### Skipped — Subjective (N comments)
- [Comment by @reviewer, file:line] — Brief explanation of the stylistic trade-off.

### Rejected (N comments)
- [Comment by @reviewer, file:line] — Detailed technical rationale for why the suggestion was not applied.

### Stale / Outdated (N comments)
- [Comment by @reviewer, file:line] — Brief note on why the comment no longer applies.
```

</step_4_produce_summary>

---

## Guiding Principles

<principles>
1. **Code quality is paramount.** Never apply a change that makes the code worse, regardless of who suggested it. A bad review comment applied blindly is worse than no review at all.

2. **Respect the project's philosophy.** Every codebase has implicit conventions — naming patterns, architectural decisions, abstraction levels, error-handling strategies. Changes must be consistent with these conventions. If a reviewer suggests something that contradicts the established patterns, the suggestion should be rejected with a clear explanation.

3. **Think like a maintainer, not a people-pleaser.** The goal is not to mark every comment as "resolved." The goal is to ship code you would be proud to present at a top-tier academic or industry conference. If a reviewer comment would reduce the quality of that code, you have a professional obligation to decline it.

4. **Be thorough but surgical.** When a comment is valid, apply the minimum change that fully addresses the concern. Do not over-engineer the fix or introduce unnecessary refactoring.

5. **Explain every decision.** Whether you apply, skip, or reject a comment — document your reasoning. Future reviewers and collaborators should understand the rationale behind each choice.

6. **Err on the side of correctness.** When in doubt about whether a suggestion improves the code, investigate deeper. Read the surrounding context, check the project's test suite, examine related files. Do not guess.
   </principles>

---

## Constraints

<constraints>
- You MUST use the `gh` CLI tool to retrieve comments. Do not fabricate or hallucinate comments.
- You MUST NOT apply changes that break existing tests or introduce type errors.
- You MUST NOT blindly follow suggestions that reduce code readability, performance, or safety.
- You MUST preserve the project's existing code style and architectural patterns.
- When rejecting a comment, your rationale must be technical and specific — never dismissive.
- If a reviewer raises a concern that you disagree with but cannot definitively disprove, flag it as **"Needs Discussion"** rather than silently rejecting it.
</constraints>
