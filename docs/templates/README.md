# Document templates

Starter templates for the structured documents referenced in `claude-instructions/development-workflow.md`'s task scale matrix.

## When to use each

| Template | Required for scale | Purpose |
| --- | --- | --- |
| [prd.md](prd.md) | Large | Product Requirements Document — what we're building and why, from the user's perspective. |
| [design-doc.md](design-doc.md) | Medium (optional) / Large (required) | Technical design — how we're building it, what we looked at, what we rejected. |
| [work-plan.md](work-plan.md) | Medium / Large | Break the work into sequenced tasks with acceptance criteria. Per-PR plan. |

ADRs are separate — see `docs/decisions/`. An ADR captures a *single decision* that shaped the design; a Design Doc captures the design *as a whole* and cites its ADRs.

## Storage

- **Active PRD / Design Doc** for a feature in progress: `docs/design/<feature-slug>/`.
- **Completed or shipped**: move to `docs/design/archive/<year>/<feature-slug>/`.
- **Work plans** are per-PR; they live in the PR body or as a checklist in the epic issue. Only commit a work plan to `docs/` if the project publishes its roadmap.

## Quality bar

All three templates share the same invariants:

- **Action-oriented.** A reader should know what to do, not just what the author thought about.
- **Decisions are cited, not hidden.** If the PRD excludes a use case, say so. If the Design Doc rejects an alternative, name it and give a one-line reason.
- **Linked.** PRDs link to Design Docs; Design Docs link to ADRs; everything links to the issue or epic that tracks it.
- **Dated and owned.** Every template has a header with author, date, and status.

## Related

- [`../../claude-instructions/development-workflow.md`](../../claude-instructions/development-workflow.md) — task scale matrix.
- [`../decisions/`](../decisions/) — ADRs.
- [`../../claude-instructions/documentation-discipline.md`](../../claude-instructions/documentation-discipline.md) — doc routing rules.
