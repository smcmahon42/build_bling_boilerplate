---
name: end-session
description: Close a work session cleanly by updating STATE.md — move completed items, mark tabled or blocked entries with reasons, record open questions, and write a short session summary.
---

# End session

Run this before wrapping up any non-trivial work session. It reconciles
`STATE.md` against what actually happened, so the next session starts from
accurate ground state.

See [`../../agent-instructions/session-handoff.md`](../../agent-instructions/session-handoff.md)
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
   - `Review at completion:` `confirmed (by operator)` if the operator
     directly confirmed the work, or `self-completed (agent, under operator
     direction)` if the agent landed it under earlier direction without an
     explicit final review.
3. **Update last-touched dates** on items that were worked on but not
   completed. State transitions (`pending → in-progress`, `→ blocked`)
   happen here too.
4. **Capture new items** that came up during the session. Add them to
   *Open work items* with `State: pending`. Set `Review` per the defaults:
   - **`confirmed`** if the operator directed the entry's creation.
   - **`unreviewed`** if the agent added it from its own observation
     (e.g., spotted a gap while doing other work). Include a `Provenance`
     block so the operator can ground their later review.
5. **Update Review transitions** for any unreviewed entries the operator
   acted on during the session (`unreviewed → confirmed | rejected | stale`).
   Don't infer the operator's review from silence — only update Review when
   they actually said.
6. **Mark tabled items.** Anything the operator explicitly parked moves to
   *Tabled* with `Why tabled`, `Un-table when`, and a carried-forward
   `Review` field.
7. **Record open questions** raised during the session — typically asks to
   the operator that didn't get answered. Include a `Provenance` block for
   agent-asked questions.
8. **Record cost signals** on multi-session entries (entries that have
   been touched in more than one session, including this one). Update
   the `Cost signals` block per
   [`../../agent-instructions/agent-cost-observability.md`](../../agent-instructions/agent-cost-observability.md):
   - Increment `Sessions to date`.
   - Append any new skills invoked to `Skills used (cumulative)`,
     de-duplicated.
   - Increment `Operator turns` by the count of operator messages this
     session.
   - Update `Context read (approx)` with this session's touched-file
     count (cumulative or per-session as the operator prefers).
   - Append any anomalies or absolute token counts to `Notes`.
   For single-session trivial entries, skip this step unless the
   operator asks for full cost capture.
9. **Age out** entries in *Recently completed* that exceed the rolling
   window (~10 entries or 14 days). Don't delete history — git keeps it;
   STATE.md stays scannable.
10. **Commit.** Conventional Commits:
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
- **Never set `Review: confirmed` on an agent-authored entry without an
  explicit operator turn confirming it.** Silence is not consent. Default
  to `unreviewed` and let the next `/start-session` surface it.
- **Always include a `Provenance` block on agent-authored entries** added
  this session. Skill name and session id alone usually suffice.

## Related

- [`../../agent-instructions/session-handoff.md`](../../agent-instructions/session-handoff.md) — doctrine and schema.
- [`../../agent-instructions/agent-cost-observability.md`](../../agent-instructions/agent-cost-observability.md) — `Cost signals` schema and external-tooling composition.
- [`start-session.md`](start-session.md) — the paired opening skill.
- [`../../agent-instructions/commit-conventions.md`](../../agent-instructions/commit-conventions.md) — commit message format.
