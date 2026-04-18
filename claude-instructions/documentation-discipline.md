# Documentation discipline

**Rule:** every change ships documentation. If the change affects behavior, architecture, operations, or the developer experience, the matching doc bucket is updated in the same PR.

## Doc routing table

| Change type | Goes in |
| --- | --- |
| Architectural decision (new boundary, tech, pattern) | `docs/decisions/NNNN-slug.md` (ADR) |
| New module or service | `docs/architecture/<module>.md` + module's own `README.md` |
| Runbook, operations, deployment | `docs/runbooks/<topic>.md` |
| Onboarding, contribution flow | `CONTRIBUTING.md` + `docs/guides/` |
| Security finding, threat model | `docs/security/` |
| QA plan, test strategy | `docs/qa/` |
| API surface change | auto-generated from source (OpenAPI, GraphQL schema, etc.) + `docs/api/` narrative |
| Claude-facing instruction update | `claude-instructions/*.md` or `<module>/claude-instructions/*.md` |
| User-visible behavior change | `CHANGELOG.md` (if the project ships versioned releases) |

If a change doesn't fit any bucket, it probably doesn't need doc — but ask: would a new teammate wonder why this exists in six months? If yes, doc it somewhere.

## Minimum bar per bucket

**Architectural decisions (`docs/decisions/`):**
- Template: Context / Decision / Alternatives Considered / Consequences / Status.
- Numbered append-only. Never delete; supersede with a new ADR and mark the old one "Superseded by NNNN."

**Runbooks:**
- The opening paragraph answers "what problem does this solve?"
- Steps are copy-pasteable commands, not prose.
- Include "how to tell it worked" and "what to do if it didn't."

**Module READMEs:**
- 20–40 lines. Plain English.
- Four questions: What is this? Why is it separate? How does it connect? Who maintains it?

**Claude instructions:**
- Action-oriented, not background theory.
- Scannable — tables, checklists, short sections.
- See [scaling-claude-instructions.md](scaling-claude-instructions.md) for structure.

## Index discipline

Each doc bucket has a `README.md` that indexes its contents. New docs require an index update in the same PR. This prevents the "where did I put that ADR?" problem.

## What NOT to document

- **Implementation detail a reader can learn from the code in 30 seconds.** The code is the source of truth.
- **Transient state** (current sprint, this week's priorities). Use issues/projects instead.
- **Personal opinions about why the old code was bad.** Describe the current state and what it enables.
- **Generated content duplicated by hand.** If the SDK docs are generated, don't paste them into `docs/`.

## Deletion

Docs rot. When a practice changes or a system is retired, either update the doc or move it to `docs/archive/<year>/` with a note at the top: "Archived YYYY-MM-DD because …". Don't leave stale docs in the active tree.

## Related

- [scaling-claude-instructions.md](scaling-claude-instructions.md) — the parallel scaling pattern for `CLAUDE.md`.
- [component-explainability.md](component-explainability.md) — the module-README rule in detail.
