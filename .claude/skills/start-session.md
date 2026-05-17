---
name: start-session
description: Orient at the top of a work session by reading STATE.md, surfacing blocked or in-progress items, and either continuing prior work or asking the operator for direction.
---

# Start session

Run this at the top of every non-trivial work session. It loads the
project-local handoff ledger (`STATE.md`) and grounds the agent on what's
open, what's blocked, and what came up since last session.

See [`../../agent-instructions/session-handoff.md`](../../agent-instructions/session-handoff.md)
for the doctrine and the STATE.md schema.

## Pre-flight

1. `STATE.md` exists at the project root. If not, copy from
   `templates/state/STATE.md`, replace `{{PROJECT_NAME}}`, and commit
   before continuing.
2. Working tree should be clean. Uncommitted changes are a signal the
   previous session didn't run `/end-session` cleanly. If dirty, ask the
   operator whether to discard, stash, or commit before proceeding.

## Steps

1. **Read `STATE.md`** in full. It's deliberately small — read all of it.
2. **Surface blocked items first.** Anything in *Open work items* with
   `State: blocked` is reported immediately — these are the items most at
   risk of being silently dropped.
3. **Surface unreviewed agent-authored items.** Anything in *Open work items*
   with `Review: unreviewed` is content the operator hasn't seen yet. Surface
   these together with their `Provenance` so the operator can confirm,
   reject, or mark stale before any work continues. This is the agent's
   handoff of "things I noticed while you weren't looking."
4. **Surface in-progress items.** Anything with `State: in-progress`
   suggests work was interrupted. Ask: continue this, or switch?
5. **Surface open questions.** Anything in the *Open questions* section
   waiting on the operator gets surfaced now — they may answer right away.
6. **Report a short orientation summary** (see output format below).
7. **Ask for direction:** continue an existing item, address an unreviewed
   one, start a new one, or answer an open question.

## Output format

A short summary block, e.g.:

```
Session state for {{PROJECT_NAME}}:
  Open:   1 blocked, 1 in-progress, 3 pending (2 unreviewed)
  Tabled: 2
  Open questions: 1 (waiting on operator)

Most recent completed: add-rate-limit-middleware — landed 2026-05-14.

Blocked items:
  - mcp-server-tls — waiting on cert provisioning (opened 2026-05-10)

Unreviewed agent-authored items (need your confirm/reject/stale):
  - error-redaction-gap — agent noticed missing PII filter in audit logs
    (skill: /security-review, session: 2026-05-15-002)
  - api-rate-limit-tightening — agent suggested 100→60 rpm after load test
    (skill: manual, session: 2026-05-15-003)

In-progress:
  - migrate-auth-to-jwt — last touched 2026-05-15

Open question:
  - Should the audit log include redacted PII or omit fully?
    (asked by agent, 2026-05-15)

How would you like to proceed?
```

## Invariants

- **Read STATE.md before any other repo exploration.** The whole point is
  to avoid re-reading the codebase to figure out what's open.
- **Never silently modify STATE.md from this skill.** Mutation is
  `/end-session`'s job.
- **Don't skip surfacing blocked or unreviewed items**, even if the operator
  names a different task. Mention them once, then proceed.
- **Don't flip `Review: unreviewed → confirmed`** on agent-authored entries
  without the operator's explicit confirmation. That promotion belongs to
  the operator; this skill surfaces, it does not promote.
- **Don't expand the ledger.** If STATE.md needs schema changes, that's an
  ADR, not a silent edit.

## Related

- [`../../agent-instructions/session-handoff.md`](../../agent-instructions/session-handoff.md) — doctrine and schema.
- [`end-session.md`](end-session.md) — the paired closing skill.
- [`../../templates/state/STATE.md`](../../templates/state/STATE.md) — the seed.
