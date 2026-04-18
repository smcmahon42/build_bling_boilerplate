---
name: new-adr
description: Scaffold a new Architectural Decision Record. Finds the next ADR number, copies the template, fills the Status line, opens it for editing, and updates the index.
---

# New ADR

Create the next ADR in `docs/decisions/`.

## Steps

1. **List existing ADRs:** `ls docs/decisions/*.md`. Ignore `README.md` and `TEMPLATE.md`.
2. **Compute the next number:** highest existing `NNNN` + 1, zero-padded to four digits.
3. **Ask the operator for a slug:** short, lowercase, hyphen-separated (e.g., `use-postgres-for-primary-store`). Validate it matches `^[a-z0-9][a-z0-9-]*$`.
4. **Copy `docs/decisions/TEMPLATE.md` to `docs/decisions/NNNN-<slug>.md`.**
5. **Fill the header fields:**
   - `# <NNNN>. <Title>` — title-case version of the slug.
   - `Status: Proposed`
   - `Date: <today's date in YYYY-MM-DD>`
6. **Open the new file for editing.** Tell the operator which sections they need to fill:
   - **Context** — the problem, forces at play.
   - **Decision** — what we're doing.
   - **Alternatives considered** — what else we looked at and why we rejected it.
   - **Consequences** — what becomes easier, what becomes harder.
7. **Update `docs/decisions/README.md`:** add a row to the index table with the ADR number, title, status, and date. If an "Unformalized decisions" tracker entry matches, offer to remove it.
8. **Stage and commit (ask first):** `git add docs/decisions/NNNN-<slug>.md docs/decisions/README.md` with message `docs(adr): add NNNN-<slug>`.

## Status transitions

- `Proposed` — draft, under discussion.
- `Accepted` — merged, in force.
- `Superseded by NNNN` — retired by a later ADR. The old ADR is **not deleted**; it stays in the tree for history. The new ADR references it.
- `Deprecated` — no longer in force but not replaced. Rare.

## When NOT to write an ADR

- Coding style or naming conventions → `claude-instructions/` or linter config.
- Minor refactors → commit message.
- Choices that don't have meaningful alternatives → not a decision, just a step.

## Invariants

- Never skip a number. Append-only.
- Never edit an accepted ADR's decision. Write a superseding ADR instead.
- Keep each ADR ~1 page. Link out to longer context if needed.

## Related

- `docs/decisions/TEMPLATE.md` — the blank template.
- `docs/decisions/0001-record-architecture-decisions.md` — the meta-ADR explaining the practice.
