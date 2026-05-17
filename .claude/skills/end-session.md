---
name: end-session
description: Close a work session cleanly by updating STATE.md — move completed items, mark tabled or blocked entries with reasons, record open questions, and write a short session summary.
---

# End session

Run this before wrapping up any non-trivial work session. It reconciles
`STATE.md` against what actually happened, so the next session starts from
accurate ground state.

See [`../../claude-instructions/session-handoff.md`](../../claude-instructions/session-handoff.md)
for the doctrine and the STATE.md schema.

## Pre-flight

1. `STATE.md` exists at the project root.
2. The operator has indicated they're wrapping the session — don't run this
   mid-task.

## Steps

1. **Compare STATE.md to actual work.** Use `git diff --stat
   origin/main...HEAD` and recent commit summaries as the signal for what
   was touched this session. Note any drift between the ledger and reality.
2. **Move completed items** from *Open work items* to *Recently completed*.
   Each completed entry gets:
   - `Completed at: YYYY-MM-DD`
   - `Session summary:` one sentence — what changed.
   - `Commit / PR:` link if useful.
3. **Update last-touched dates** on items that were worked on but not
   completed. State transitions (`pending → in-progress`, `→ blocked`)
   happen here too.
4. **Capture new items** that came up during the session. Add them to
   *Open work items* with `State: pending`.
5. **Mark tabled items.** Anything the operator explicitly parked moves to
   *Tabled* with `Why tabled` and `Un-table when`.
6. **Record open questions** raised during the session — typically asks to
   the operator that didn't get answered.
7. **Age out** entries in *Recently completed* that exceed the rolling
   window (~10 entries or 14 days). Don't delete history — git keeps it;
   STATE.md stays scannable.
8. **Commit.** Conventional Commits:
   `docs(state): update STATE.md — <one-line session summary>`. Stage only
   `STATE.md` — this skill doesn't touch code.

## Output format

After updating, report:

```
STATE.md updated:
  Moved to completed: <slug-id> (<one-line summary>)
  Last-touched updates: 2
  New items added: 1 (<slug-id>)
  Tabled: 0
  Open questions: 1 added

Commit staged: docs(state): update STATE.md — finished rate-limit middleware
```

## Invariants

- **Don't lose information.** Tabling is preferred over deletion when in
  doubt. Aging out is acceptable — git history retains everything.
- **One sentence for session summary.** Long summaries belong in commit
  bodies, not STATE.md.
- **Commit STATE.md changes separately from code changes** unless they
  belong to the same logical unit (rare). Mixing them muddies `git log`.
- **Ask the operator before deleting an item outright.** Tabling is the
  default for "not doing this right now."

## Related

- [`../../claude-instructions/session-handoff.md`](../../claude-instructions/session-handoff.md) — doctrine and schema.
- [`start-session.md`](start-session.md) — the paired opening skill.
- [`../../claude-instructions/commit-conventions.md`](../../claude-instructions/commit-conventions.md) — commit message format.
