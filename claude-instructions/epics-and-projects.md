# Epics and projects

"Epic" means different things to different teams. This file covers two patterns for tracking multi-issue work, and when to pick each.

## Pattern A — Doc-embedded epics (lightweight)

Epics live in a markdown document under `docs/` (typically `docs/roadmap/` or `docs/implementation-plan.md`). The document lists phases or workstreams; each phase links to GitHub issues by number.

**Shape:**

```markdown
# Implementation plan

## Phase 3 — Ingestion pipeline (epic #123)

Goal: end-to-end ingestion with back-pressure and retries.

### Milestones

- [ ] Batch reader (#124)
- [ ] Back-pressure queue (#125)
- [ ] Retry with exponential backoff (#126)
- [ ] Observability: metrics + traces (#127)

### Exit criteria

- Sustained 10k events/sec in staging for 24h.
- Retry loop exercised by fault injection tests.
- Runbook published in `docs/runbooks/ingestion.md`.
```

**Pros:**
- Version-controlled with the rest of the docs.
- Reviewed in PRs like any other design change.
- No external tooling dependency.

**Cons:**
- GitHub UI doesn't render progress bars or roll up.
- Filtering/views are bounded by markdown.

**Pick this when:**
- The team is small (< 10 contributors).
- Epics rarely span more than a handful of issues.
- You prefer docs over project boards.

## Pattern B — GitHub Projects v2 (board-driven)

Epics are GitHub Projects v2 boards. Each epic has a board with columns (To do / In progress / In review / Done), and issues are the cards. Parent-child relationships are tracked via the `tracks`/`tracked by` field on issues.

**Shape:**
- One board per epic (or one board per quarter, with an "epic" field on each issue).
- Every issue on the board has: `type` (epic | feature | bug | chore), `priority`, `status`.
- Epic issue holds a checklist of child issue references, auto-synced by the `actions/add-to-project` action.

**Pros:**
- Rich filtering, views, and roll-ups.
- Works well for stakeholders who live in GitHub UI.
- Progress visible at a glance.

**Cons:**
- Requires up-front setup.
- Second source of truth next to the code.
- Config changes aren't version-controlled as neatly.

**Pick this when:**
- The team is larger and coordinates work across modules.
- Non-engineers (PM, design, support) are active issue participants.
- You want cross-epic dashboards.

## Hybrid (both)

Many teams land here: docs hold the **narrative** (why, goals, exit criteria, links to ADRs), and a Project board holds the **tracking** (status, ownership, burn-down). The narrative doc links to the board; each board issue links back to the doc section.

This works but costs discipline to keep the two in sync.

## Issue templates that support either pattern

The boilerplate ships four issue templates:

| Template | When |
| --- | --- |
| `bug_report.md` | A defect or regression. Reproducer, expected vs actual, environment. |
| `feature_request.md` | New capability. Problem, proposal, acceptance criteria. |
| `epic.md` | Parent issue for a multi-issue effort. Context, milestones, exit criteria, child issues. |
| `adr_proposal.md` | Proposing an architectural decision. Mirrors the ADR template; moves into `docs/decisions/` once accepted. |

## Labels (baseline taxonomy)

Shipped in `.github/labels.yml`. Sync with `gh label sync .github/labels.yml` (or a label-sync action on PR).

```
type:feature          new capability
type:bug              defect
type:chore            infra, tooling, deps
type:docs             doc-only
type:epic             parent issue tracking a larger effort
type:adr              ADR proposal

priority:p0           drop everything
priority:p1           next up
priority:p2           scheduled
priority:p3           nice to have

status:blocked        waiting on something
status:in-review      PR open
status:needs-info     waiting on reporter

area:api              adjust per project
area:web
area:docs
area:ci

security              security-relevant
breaking              breaking change
dependencies          dependency bump
```

Adjust `area:*` to match the project's module list.

## Epic lifecycle

1. **Propose.** File an epic issue (template), describe context + goals + exit criteria.
2. **Break down.** Create child issues; link them to the epic (checklist, `tracks` field, or doc reference).
3. **Track.** Pick Pattern A or B (or hybrid). Update progress as child issues close.
4. **Close.** When exit criteria are met, close the epic with a summary comment linking the ADR(s) it produced and the runbook (if applicable). Archive the board (Pattern B) or move the doc section to `docs/roadmap/completed/` (Pattern A).

## Related

- [documentation-discipline.md](documentation-discipline.md) — where epic docs live.
- [commit-conventions.md](commit-conventions.md) — how to reference epic/child issues in commits.
