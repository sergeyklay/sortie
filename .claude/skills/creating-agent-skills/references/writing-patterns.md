# Writing Patterns for Skill Instructions

Concrete patterns for structuring skill instructions. Each pattern includes when to use it, a complete example, and common pitfalls.

## Contents

- Degrees of freedom
- Template pattern
- Examples (multishot) pattern
- Workflow with checklists
- Conditional workflow pattern
- Feedback loop pattern
- Progressive disclosure patterns
- Script integration patterns
- Cross-references without force-loading
- MCP tool references

## Degrees of Freedom

Match specificity to how fragile and variable the task is.

**High freedom** - multiple valid approaches, context-dependent:

```markdown
## Code review process

1. Analyze code structure and organization
2. Check for potential bugs or edge cases
3. Suggest improvements for readability
4. Verify adherence to project conventions
```

Why high freedom works here: code reviews depend heavily on context. The agent needs latitude to adapt its approach based on what it finds.

**Medium freedom** - preferred pattern exists, some variation acceptable:

```markdown
## Generate report

Use this template and customize based on the data:

# [Analysis Title]

## Executive summary

[One-paragraph overview]

## Key findings

[Adapt sections based on what you discover]

## Recommendations

[Tailor to the specific context]
```

**Low freedom** - fragile operations, consistency critical:

```markdown
## Database migration

Run exactly this script:
python scripts/migrate.py --verify --backup
Do not modify the command or add additional flags.
The --verify flag prevents destructive changes without confirmation.
The --backup flag creates a rollback point before any schema changes.
```

Why low freedom works here: database migrations can cause data loss. The specific flags exist for safety, and skipping them creates real risk.

**Choosing the right level**: If getting it wrong causes data loss or broken output, use low freedom. If the agent's creativity would improve the result, use high freedom. When in doubt, start with medium freedom and adjust based on testing.

## Template Pattern

Provide output format templates. Match strictness to requirements.

**Strict** (API responses, data formats, compliance documents):

```markdown
## Report structure

ALWAYS use this exact template structure:

# [Analysis Title]

## Executive summary

[One-paragraph overview of key findings]

## Key findings

- Finding 1 with supporting data
- Finding 2 with supporting data

## Recommendations

1. Specific actionable recommendation
2. Specific actionable recommendation
```

**Flexible** (when adaptation improves quality):

```markdown
## Report structure

Sensible default format - adapt based on what you discover:

# [Analysis Title]

## Executive summary

## Key findings (adapt sections based on analysis)

## Recommendations (tailor to context)

Adjust sections as needed. Add "Risks" if significant concerns emerge.
Drop "Recommendations" if the task is purely analytical.
```

## Examples (Multishot) Pattern

Input/output pairs communicate style and expectations more effectively than descriptions alone. Include 3-5 diverse examples covering edge cases.

```markdown
## Commit message format

**Example 1** (new feature):
Input: Added user authentication with JWT tokens
Output:
feat(auth): implement JWT-based authentication
Add login endpoint and token validation middleware

**Example 2** (bug fix):
Input: Fixed bug where dates displayed incorrectly
Output:
fix(reports): correct date formatting in timezone conversion
Use UTC timestamps consistently across report generation

**Example 3** (multiple changes):
Input: Updated dependencies and refactored error handling
Output:
chore: update dependencies and refactor error handling

- Upgrade lodash to 4.17.21
- Standardize error response format across endpoints

Follow this pattern: type(scope): brief description, then detailed explanation.
```

**Tip**: Include at least one edge case example. If your skill handles "normal" cases well but fails on edge cases, add examples specifically for those.

## Workflow Pattern with Checklists

For complex multi-step processes, provide a checklist the agent can copy into its response and track. This prevents skipped steps and gives the user visibility into progress.

```markdown
## PDF form filling workflow

Copy this checklist and check off items as you complete them:
Task Progress:

- [ ] Step 1: Analyze the form (run analyze_form.py)
- [ ] Step 2: Create field mapping (edit fields.json)
- [ ] Step 3: Validate mapping (run validate_fields.py)
- [ ] Step 4: Fill the form (run fill_form.py)
- [ ] Step 5: Verify output (run verify_output.py)

**Step 1: Analyze the form**
Run: python scripts/analyze_form.py input.pdf
This extracts form fields, their types, and coordinates.
Expected output: fields.json with field definitions.

**Step 2: Create field mapping**
Edit fields.json to add values for each field.
Match field names to the data source provided by the user.

**Step 3: Validate mapping**
Run: python scripts/validate_fields.py fields.json
Fix any validation errors before continuing.
Common issues: missing required fields, type mismatches.

**Step 4: Fill the form**
Run: python scripts/fill_form.py input.pdf fields.json output.pdf

**Step 5: Verify output**
Run: python scripts/verify_output.py output.pdf
If verification fails, return to Step 2.
```

## Conditional Workflow Pattern

Guide the agent through decision points. Keep branches concise in SKILL.md; push large branches into reference files.

**Inline branches** (when branches are short):

```markdown
## Document modification workflow

1. Determine the modification type:
   **Creating new content?** -> Follow "Creation workflow" below
   **Editing existing content?** -> Follow "Editing workflow" below

2. Creation workflow:
   - Use docx-js library
   - Build document from scratch
   - Export to .docx format

3. Editing workflow:
   - Unpack existing document
   - Modify XML directly
   - Validate after each change
   - Repack when complete
```

**Reference branches** (when branches are large):

```markdown
## Cloud deployment

Detect target platform, then load the appropriate guide:
**AWS?** -> See [references/aws.md](references/aws.md)
**GCP?** -> See [references/gcp.md](references/gcp.md)
**Azure?** -> See [references/azure.md](references/azure.md)

Each guide covers authentication, resource provisioning,
deployment commands, and rollback procedures.
```

The agent reads only the relevant reference file, keeping context usage minimal.

## Feedback Loop Pattern

For quality-critical operations, implement a validate-fix-repeat cycle. This pattern catches errors early and prevents cascading failures.

```markdown
## Document editing process

1. Make edits to word/document.xml
2. Validate immediately: python scripts/validate.py unpacked_dir/
3. If validation fails:
   - Review the error message carefully
   - Fix the issue in the XML
   - Run validation again
4. Only proceed when validation passes
5. Rebuild: python scripts/pack.py unpacked_dir/ output.docx
6. Open and test the output document
```

**When to use feedback loops**: Any time the output format is strict (XML, JSON schemas, binary formats), when errors compound (each step builds on the previous), or when the cost of failure is high (data loss, broken documents).

## Progressive Disclosure Patterns

### Pattern 1: High-level guide with on-demand references

```markdown
# PDF Processing

## Quick start

Extract text with pdfplumber:
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
text = pdf.pages[0].extract_text()

## Advanced features

**Form filling**: See [FORMS.md](FORMS.md) for complete guide
**API reference**: See [REFERENCE.md](REFERENCE.md) for all methods
```

### Pattern 2: Domain-specific organization

```plaintext
bigquery-skill/
├── SKILL.md (overview + navigation)
└── references/
    ├── finance.md (revenue, billing)
    ├── sales.md (pipeline, accounts)
    └── product.md (API usage, features)
```

SKILL.md points to domain files. Agent reads only the relevant one based on the user's query. Include a quick-search command for discoverability:

```markdown
## Quick search

grep -i "revenue" references/finance.md
grep -i "pipeline" references/sales.md
```

### Pattern 3: Conditional details

```markdown
# DOCX Processing

## Creating documents

Use docx-js for new documents. See [DOCX-JS.md](DOCX-JS.md).

## Editing documents

For simple edits, modify XML directly.
**For tracked changes**: See [REDLINING.md](REDLINING.md)
```

**Key rule**: Keep references one level deep from SKILL.md. Avoid chains like SKILL.md -> advanced.md -> details.md. Agents may only partially read deeply nested files (using commands like `head -100`), resulting in incomplete information.

## Script Integration Patterns

Scripts run via bash without loading into context. Only their output enters the context window. This makes them ideal for deterministic operations.

**Execute pattern** (most common):

```markdown
Run analyze_form.py to extract fields:
python scripts/analyze_form.py input.pdf > fields.json
```

**Execute with clear I/O documentation**:

```markdown
## scripts/validate_boxes.py

Check for overlapping bounding boxes in form fields.
Input: python scripts/validate_boxes.py fields.json
Output: "OK" if no overlaps, or lists specific conflicts with coordinates
Exit code: 0 = pass, 1 = conflicts found
```

**Read as reference** (rare - only when the agent needs to understand the algorithm):

```markdown
See scripts/scoring.py for the ranking algorithm.
The scoring weights are: recency (0.4), relevance (0.3), authority (0.3).
```

## Cross-References Without Force-Loading

When skills relate to each other, avoid `@` syntax that force-loads referenced files and burns tokens before they are needed. Use plain text references with explicit markers:

```markdown
## Prerequisites

Before deploying, ensure all tests pass. If tests fail,
follow the procedures in the `test-runner` skill (REQUIRED).

For performance-critical deployments, also consult the
`load-testing` skill (OPTIONAL) for pre-deployment benchmarks.
```

The REQUIRED/OPTIONAL markers tell the agent how to prioritize without forcing immediate loading. The agent loads referenced skills only when it reaches the step that needs them.

**When to use markers:**
- REQUIRED: Agent must load this skill before proceeding
- OPTIONAL: Agent loads only if the specific condition applies
- RECOMMENDED: Agent should consider loading but can proceed without

## MCP Tool References

When referencing MCP tools in skill instructions, use fully qualified names to prevent "tool not found" errors:

```markdown
Use the BigQuery:bigquery_schema tool to retrieve table schemas.
Use the GitHub:create_issue tool to create issues.
Use the Jira:search_issues tool to find related tickets.
```

Format: `ServerName:tool_name`. The server name must match the MCP server configuration name exactly.
