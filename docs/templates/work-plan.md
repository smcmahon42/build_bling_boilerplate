# Work plan — <Feature or task>

- **Author:** <name>
- **Status:** Draft  <!-- Draft | In progress | Done -->
- **Last updated:** YYYY-MM-DD
- **Tracking issue / epic:** #
- **Design Doc:** <link, if any>

## Scope

One paragraph restating what this plan covers. If it's narrower than the linked Design Doc or PRD, say so explicitly.

## Milestones

Each milestone is a shipped, visible thing. Not a check-in or a status update — an artifact a reviewer can open.

- [ ] **M1 — <milestone>** — *acceptance:* <what's true when this is done>
- [ ] **M2 — <milestone>** — *acceptance:* ...
- [ ] **M3 — <milestone>** — *acceptance:* ...

## Tasks

Sequenced. One PR each where possible. Reference the milestone they advance.

| # | Task | Milestone | Owner | Est. effort | Status |
| --- | --- | --- | --- | --- | --- |
| 1 | <task> | M1 | <@handle> | <S / M / L> | Todo |
| 2 | <task> | M1 | | | Todo |
| 3 | <task> | M2 | | | Todo |

Effort guide (matches `claude-instructions/development-workflow.md#task-scale-matrix`): S = 1–2 files, M = 3–5 files, L = 6+ files.

## Dependencies

Anything external this plan depends on — other teams, infra, vendor work.

- <Dependency> — *owner:* <name>, *needed by:* <date>

## Risks

What could derail the plan. Triage weekly.

- <Risk> — *mitigation:* <how we'll handle it>

## Exit criteria

All of these must be true for the plan to be "done":

- [ ] All milestones checked.
- [ ] Tests for each milestone shipped and green.
- [ ] Docs updated (ADR / runbook / module README as applicable).
- [ ] Security review completed (`/security-review` run on the cumulative diff).
- [ ] Tracking issue / epic closed.

## Related

- Design Doc: <link>
- PRD: <link>
- ADRs: <list>
