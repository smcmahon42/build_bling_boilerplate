# claude-instructions

Topic files loaded on-demand by Claude during a session. The root `CLAUDE.md` is a slim **router** — it points here by topic; content lives in these files so no single file has to be read in full every time.

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
| [claude-memory.md](claude-memory.md) | User-local Claude memory system, namespacing, when to save. |
| [epics-and-projects.md](epics-and-projects.md) | Two patterns for epic tracking: doc-embedded vs GitHub Projects v2. |
| [scaling-claude-instructions.md](scaling-claude-instructions.md) | How to grow this directory as the project grows (Stage 1 → 2 → 3). |

## When to add a new topic file

- A recurring task needs guidance that doesn't fit in an existing file.
- A section of an existing file grows past ~300 lines.
- A module-specific concern emerges that would distract from the cross-cutting file.

If the new topic is **module-specific** (e.g. "how we do auth in the API service"), prefer putting it in the module's own `claude-instructions/` directory once you've moved to Stage 3. See [scaling-claude-instructions.md](scaling-claude-instructions.md).

## When to edit an existing topic file

Default to editing. Add a new section if the concept fits; split into a new file only when the existing one crosses the size threshold or the topic is genuinely separate.

Keep these files:

- **Scannable.** Lead with tables, checklists, numbered steps.
- **Action-oriented.** Describe what to do, not background theory.
- **Honest about escape hatches.** If a rule has an exception, name it.
