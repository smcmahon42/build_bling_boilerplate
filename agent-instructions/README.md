# agent-instructions

Topic files loaded on demand by any agent client during a session. The root `AGENTS.md` is the neutral router; client-specific files such as `CLAUDE.md`, Codex skills, or homegrown adapters should point here instead of duplicating doctrine.

## Files

| File | Purpose |
| --- | --- |
| [development-workflow.md](development-workflow.md) | TDD loop, branching, PR flow, local-CI-first option, review gates. |
| [testing-practices.md](testing-practices.md) | Test pyramid, framework-native unit tests, Playwright for frontend e2e, regression rule. |
| [security-practices.md](security-practices.md) | OWASP Web Top 10, LLM Top 10, Agentic Top 10; review gates; secret handling. |
| [documentation-discipline.md](documentation-discipline.md) | Doc buckets, routing table, every-change-ships-docs rule. |
| [commit-conventions.md](commit-conventions.md) | Conventional Commits, scope, PR references, co-author trailers. |
| [dependency-vetting.md](dependency-vetting.md) | 5-point checklist before adding any dependency. |
| [component-explainability.md](component-explainability.md) | README at every meaningful boundary. |
| [agent-primitives.md](agent-primitives.md) | Typed Task/Result/Evidence/Provenance/CapabilityToken — the language of agent-to-agent handoff. |
| [contract-discipline.md](contract-discipline.md) | OpenAPI / JSON Schema / protobuf as source of truth; versioning; codegen rules. |
| [logging-practices.md](logging-practices.md) | Structured JSON logs with mandatory correlation fields and redaction. |
| [observability-practices.md](observability-practices.md) | OpenTelemetry tracing and metrics; `Provenance` as the correlation anchor. |
| [operator-memory.md](operator-memory.md) | Client-local operator memory conventions, namespacing, and when to save. |
| [epics-and-projects.md](epics-and-projects.md) | Two patterns for epic tracking: doc-embedded vs GitHub Projects v2. |
| [scaling-agent-instructions.md](scaling-agent-instructions.md) | How to grow this directory as the project grows (Stage 1 → 2 → 3). |

## Adapter rule

Shared instructions live here. Agent clients may provide their own invocation mechanics:

- Claude Code uses `.claude/skills/` and `CLAUDE.md`.
- Codex can use repo `AGENTS.md`, local skills, plugins, or `codex review` prompts.
- Homegrown agents should read `AGENTS.md` and these topic files directly.

Adapters should link to these files rather than copy their content. If a rule applies to every agent, update this directory. If a rule applies only to one client, keep it in that client's adapter.

## When to add a new topic file

- A recurring task needs guidance that doesn't fit in an existing file.
- A section of an existing file grows past ~300 lines.
- A module-specific concern emerges that would distract from the cross-cutting file.

If the new topic is module-specific, prefer putting it in the module's own `agent-instructions/` directory once you've moved to Stage 3. See [scaling-agent-instructions.md](scaling-agent-instructions.md).

## When to edit an existing topic file

Default to editing. Add a new section if the concept fits; split into a new file only when the existing one crosses the size threshold or the topic is genuinely separate.

Keep these files:

- **Scannable.** Lead with tables, checklists, numbered steps.
- **Action-oriented.** Describe what to do, not background theory.
- **Honest about escape hatches.** If a rule has an exception, name it.
