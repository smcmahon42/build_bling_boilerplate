# Custom agents

This directory is empty by default — the boilerplate ships no custom agent definitions. Add project-specific agents here as the project grows.

## When to add an agent vs. a skill

**Skill (`.claude/skills/*.md`):** a reusable workflow the operator invokes by name (`/security-review`, `/new-adr`). Skills live inside the main conversation — they share context with everything the operator is doing. Use skills for recurring, well-scoped tasks.

**Agent (`.claude/agents/*.md`):** a specialized sub-agent the main Claude instance can delegate to. Agents have their own isolated context (they don't see the main conversation history). Use agents for long-running research, parallelizable investigation, or tasks where isolating context is valuable (keeps the main conversation's context window clean).

If in doubt, start with a skill. Promote to an agent when the workflow:
- Takes more than a few minutes,
- Produces large intermediate output you don't want in the main context,
- Can run in parallel with other work,
- Has a clearly-defined self-contained input and output.

## Adding an agent

Create `<name>.md` with frontmatter:

```markdown
---
name: <agent-name>
description: <one-line purpose — surfaced to Claude when choosing when to invoke>
tools: [<tools this agent needs>]  # optional; defaults to all
---

# <Agent name>

<Body: instructions for the agent>
```

Keep the body focused — it's the agent's system prompt. State the task, the expected output format, and any invariants.

## Examples that might fit here over time

- `pr-reviewer` — deep review of a diff: correctness, security, style, tests.
- `doc-writer` — generate first-draft docs from a code change.
- `branch-audit` — summarize what's on a branch vs. main.
- `release-notes` — compile user-visible changes since last tag.

These are illustrative, not shipped. Add them when the project has a clear use case.

## Related

- `.claude/skills/` — user-invokable skills.
- [Claude Code docs on sub-agents](https://docs.claude.com/en/docs/claude-code/sub-agents) for the schema reference.
