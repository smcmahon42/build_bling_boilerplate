# Architectural Decision Records

This directory holds the project's ADRs — short, numbered records of decisions that shape the system. ADRs are **append-only**: once accepted, an ADR is not edited. Decisions that overturn an earlier ADR create a new ADR that marks the old one "Superseded by NNNN."

See [0001](0001-record-architecture-decisions.md) for why this project uses ADRs and how they fit into the broader workflow.

## Index

| # | Title | Status | Date |
| --- | --- | --- | --- |
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions | Accepted | <!-- YYYY-MM-DD filled on merge --> |
| [0002](0002-boundary-annotation-enforcement.md) | Enforce boundary conventions via in-code annotations + allowlist | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |
| [0003](0003-agent-primitives-as-foundation.md) | Typed agent primitives are the canonical handoff between rings | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |
| [0004](0004-contracts-as-source-of-truth.md) | Machine-readable contracts are the source of truth; code is generated | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |
| [0005](0005-structured-logging-and-redaction.md) | Structured JSON logs with mandatory correlation fields and a redaction layer | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |
| [0006](0006-human-surfaces-are-views-over-contracts.md) | Human surfaces are thin views over contracts, never the other way around | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |
| [0007](0007-opentelemetry-and-correlation.md) | OpenTelemetry for tracing and metrics, with Provenance as the correlation anchor | Accepted (pattern) | <!-- YYYY-MM-DD filled on merge --> |

<!--
Add new rows here in number order. When a new ADR supersedes an old one, keep both rows and update the old row's Status to "Superseded by NNNN".
-->

## Unformalized decisions

Decisions that have been made but not yet written up as ADRs. When a decision here gains enough weight (or anyone asks "why did we decide X?"), promote it to a full ADR.

<!--
Template row:
- <decision in one sentence> — rationale lives in <link>. To be ADR'd when: <trigger>.
-->

## How to add an ADR

Run `/new-adr` (the skill at `.claude/skills/new-adr.md`) — it handles numbering, template copy, status line, and index update.

Manual process:

1. Find the next number: `ls -1 *.md | grep -E '^[0-9]{4}-' | tail -1`.
2. Copy the template: `cp TEMPLATE.md NNNN-<slug>.md`.
3. Fill Context / Decision / Alternatives / Consequences.
4. Set `Status: Proposed`.
5. Open a PR; on merge, flip to `Status: Accepted` and update this index.

## Template

See [TEMPLATE.md](TEMPLATE.md).
