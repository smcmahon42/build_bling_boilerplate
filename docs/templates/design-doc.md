# Design Doc — <Feature or system name>

- **Author:** <name>
- **Status:** Draft  <!-- Draft | In review | Approved | Superseded | Implemented -->
- **Last updated:** YYYY-MM-DD
- **PRD:** <link, if one exists>
- **Tracking issue / epic:** #
- **ADRs produced by this design:** <NNNN, or "to be written">

## Summary

One paragraph. What we're building, the shape of the approach, and why this approach over others. A reader should leave this paragraph knowing whether to read on.

## Context

What's true today that motivates this change. Link to the PRD (if one exists), the relevant code, any prior ADRs, and any data.

## Goals

From the PRD, restated as *engineering* goals:

- <Goal>
- <Goal>

## Non-goals

- <Non-goal>
- <Non-goal>

## Design

The proposed solution. Lead with the big picture; descend into detail.

### High-level architecture

Text or a diagram. Name the components and how they connect. If a picture is worth a thousand words, use one; otherwise a short list beats a pasted ASCII diagram.

### Key abstractions

What new boundaries exist after this ships. What's their contract? What do they *not* promise?

| Component | Responsibility | Depends on | Exposes |
| --- | --- | --- | --- |
| <name> | <one line> | <components> | <interface / API> |

### Data model

Schemas, persistence, migration. If this adds or changes a data class, call out:

- What's stored, where, for how long.
- Who can read / write it.
- Any PII / secret / regulated data implications.

### Interfaces

Public contracts introduced or changed — HTTP endpoints, CLI commands, library functions, agent tools. Link to the generated contract (OpenAPI, gRPC, JSON schema) rather than pasting it.

### Flows

The 2–4 most important happy paths, step by step. Include the important failure paths too — what happens under retry, timeout, partial write, concurrent request.

## Alternatives considered

At least two. For each:

- **<Alternative name>** — brief description of the approach. *Why rejected:* one or two sentences.

If there are no real alternatives, the design is probably smaller than this template warrants — use a commit message or ADR instead.

## Consequences

- **Easier:** what this design enables.
- **Harder:** what it complicates.
- **Follow-ups:** work implied by this design that ships later.

## Testing strategy

How we'll know it works.

- **Unit:** what behavior is tested at the function level.
- **Integration:** what boundaries are exercised end-to-end in CI.
- **End-to-end:** Playwright (frontend) or contract tests (services) — which user flows.
- **Load / failure:** if relevant, how we'll prove the non-functional requirements hold.
- **Regression:** how this design prevents reintroducing the bugs it's replacing.

## Security review

Walk the relevant OWASP categories. Reuse the structure from `.claude/skills/security-review.md`.

- **OWASP Web Top 10:** <which apply; one line each>
- **OWASP LLM Top 10:** <which apply>
- **OWASP Agentic Top 10:** <which apply>
- **Mitigations:** <what this design does to address each applicable category>

## Observability

What logs / metrics / traces will exist? What alerts? Who gets paged when what breaks?

## Rollout

How this ships safely. Mirror the PRD's rollout section or expand on it.

- **Feature flag:** <name, default>
- **Migration:** <backfill strategy, downtime expectations>
- **Rollback:** <how to undo if something goes wrong>

## Open questions

Unresolved items. Owner and target-resolution date for each.

- <Question> — *owner:* <name>, *by:* <date>

## Related

- PRD: <link>
- ADRs: <list>
- Issue / epic: #
- Related designs: <links>
