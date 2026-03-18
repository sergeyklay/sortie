# Skills Ecosystem: When to Use What

Decision guide for choosing between skills and other agent customization approaches.

## Contents

- Comparison table
- Skills vs Custom Instructions
- Skills vs AGENTS.md
- Skills vs MCP Servers
- Skills vs Prompt Files / Commands
- Skills vs Cursor Rules
- The "skills replace MCP" myth
- Decision framework

## Comparison Table

| Artifact | Loading | Purpose | Use When |
|----------|---------|---------|----------|
| **Skills** | On-demand (model-invoked) | Procedural knowledge | Teaching HOW to do specific tasks |
| **Custom Instructions** | Always-on | Coding standards, preferences | Rules that apply to every session |
| **AGENTS.md** | Session start | Project context | Describing WHAT your project is |
| **MCP Servers** | Always connected | External system access | Live data, write operations, auth |
| **Prompt Files** | User-invoked | Templates | Slash commands, boilerplate |
| **Cursor Rules** | Pattern-matched | IDE-specific rules | Cursor-only environments |

## Skills vs Custom Instructions

Custom instructions (`.instructions.md`, `copilot-instructions.md`, User Rules) are always-on rules that apply to every session. They define coding standards, naming conventions, and general preferences.

Skills are on-demand capabilities that load only when relevant.

**Use custom instructions for:** "Always use TypeScript", "Follow our naming conventions", "Prefer async/await over callbacks"

**Use skills for:** "Here is the step-by-step procedure for creating a new API endpoint in our specific framework"

## Skills vs AGENTS.md

AGENTS.md files provide persistent context about your project that the agent reads at session start. They describe architecture, conventions, and project-specific knowledge.

Skills are procedural - they describe HOW to do something, not WHAT something is.

**AGENTS.md explains:** Project structure, tech stack, team conventions, architectural decisions

**Skills teach:** How to execute a specific workflow within that structure

## Skills vs MCP Servers

This is the most misunderstood distinction. MCP and skills operate at completely different layers.

**MCP provides:** Live connections to external systems (databases, APIs, issue trackers). Without MCP, the agent cannot query your database or create a Jira ticket, no matter how detailed its instructions are.

**Skills provide:** Procedural knowledge about HOW to use those connections effectively. Your MCP server connects to Jira. Your skill teaches the agent your team's conventions for writing tickets.

**Neither replaces the other.** MCP is a power cord. Skills are a user manual.

### Where MCP is essential (skills cannot help)

- **Live data access:** Querying databases, reading Slack messages, checking Kubernetes state
- **Real-time streaming:** Log monitoring, deployment pipelines, analytics
- **Write operations:** Creating issues, sending messages, deploying code
- **Authentication:** OAuth flows, API keys, service account credentials
- **Cross-agent communication:** Capability discovery in multi-agent systems

### Where skills excel (MCP is overkill)

- **Procedural knowledge:** Code review checklists, migration steps, commit conventions
- **Context efficiency:** Skills load ~100 tokens at startup; MCP tool definitions can consume 8,000+ tokens
- **Portability:** Skills are markdown files, no infrastructure required
- **Domain expertise transfer:** Teaching how your team handles incident response

### Where they work together

The most effective setups combine both:

```
MCP Server: Google Drive connection (read/write files)
Skill: Company document procedures (folder structure, naming conventions, templates)
```

The MCP server provides access. The skill ensures that access follows your team's conventions.

## Skills vs Prompt Files / Commands

Prompt files (`.prompt.md`) and Cursor Commands are user-invoked templates triggered by slash commands. You type `/create-component` and the prompt runs.

Skills are model-invoked - the agent decides when to activate them based on the task at hand.

**Use prompt files for:** Boilerplate generation, repetitive tasks you trigger manually

**Use skills for:** Procedures the agent should recognize and apply autonomously

## Skills vs Cursor Rules

Cursor Rules (`.cursor/rules/`) use a structured format with frontmatter controlling when they apply (always, intelligently, by file pattern, or manually).

Skills follow the universal SKILL.md format and work across platforms.

**Use Cursor Rules if:** You only use Cursor and want tight IDE integration

**Use skills if:** Your team uses multiple tools, or you want portable procedures

## The "Skills Replace MCP" Myth

When Agent Skills became an open standard in December 2025, a narrative emerged: "MCP is obsolete, skills do everything." This is a category error.

**What skills actually replace:** The informal knowledge developers carry in their heads about how to use tools effectively. Skills formalize that knowledge.

**What MCP provides that skills cannot:** Live connections, authentication, streaming, write operations to external systems.

**The real alternative to MCP:** Not skills - it's bash and curl. When developers argue against MCP, they're arguing for direct CLI access documented through skills.

### When MCP becomes unnecessary

- CLI tools already handle it (git, npm, kubectl via bash)
- Static data that rarely changes (embed in `references/` instead)
- Local development workflows where you're the only consumer

## Decision Framework

Ask these questions:

1. **Does this require live access to external state?**
   Yes → MCP (or equivalent tool access)

2. **Does this require teaching a procedure or convention?**
   Yes → Skill

3. **Does this require authentication across a team?**
   Yes → MCP (centralized credential management)

4. **Should this apply to every session automatically?**
   Yes → Custom Instructions

5. **Is this project context the agent needs at startup?**
   Yes → AGENTS.md

6. **Is this a user-triggered template?**
   Yes → Prompt file / Command

7. **Can existing CLI tools accomplish this?**
   Yes → Consider skill teaching CLI usage instead of MCP server

Most production environments use a hybrid: MCP for essential integrations, skills for team-specific procedures.
