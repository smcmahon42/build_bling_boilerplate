# Session handoff

Projects don't get built in one session by one person. `STATE.md` is the
project-local ledger that lets the next session ground itself in seconds
instead of re-reading the repo to figure out what's open, what's done, and
what got parked.

Pair this with [`templates/state/STATE.md`](../templates/state/STATE.md) (the
seed) and the [`/start-session`](../.claude/skills/start-session.md) and
[`/end-session`](../.claude/skills/end-session.md) skills.

## What STATE.md is

A single committed file at the project root with four sections:

| Section | Holds | Lifecycle |
| --- | --- | --- |
| **Open work items** | Active or queued work | Moves to "Recently completed" when done; moves to "Tabled" when parked |
| **Recently completed** | Rolling window of finished items | Ages out (~10 entries or 14 days, whichever is shorter) |
| **Tabled** | Deliberately parked work | Stays until the un-table condition is met |
| **Open questions** | Asks waiting on a human or external answer | Resolves into an open work item or gets removed |

The full row schema lives in `templates/state/STATE.md` — copy it to the
project root on bootstrap and treat it as living history.

## What STATE.md is NOT

- **Not user-local Claude memory.** That lives at
  `~/.claude/projects/<encoded-path>/memory/` and holds *biographical* facts
  about the operator plus time-bound *project* context. STATE.md holds
  *operational handoff* — different lifecycle, different scope. See
  [`claude-memory.md`](claude-memory.md).
- **Not an Architectural Decision Record.** ADRs in `docs/decisions/` are
  append-only records of accepted decisions. STATE.md rows are mutable
  working state. An open-work item that crystallizes into an architectural
  decision *graduates* into an ADR (and is removed from STATE.md).
- **Not the commit log.** `git log` is authoritative for what *was* done.
  STATE.md is authoritative for what's *currently open*.
- **Not the external issue tracker.** STATE.md is the file-system primitive
  that works inside the repo without external tooling. Projects that adopt
  GitHub Issues or Linear can mirror, but STATE.md stays the local source
  of truth that any agent or contributor can read on first clone.

## Lifecycle

1. **A new item appears** in *Open work items* when the operator (human or
   agent) starts on something or queues it for later. Use a short slug as
   the heading (e.g. `### add-rate-limit-middleware`).
2. **State transitions** happen in place — `pending → in-progress` when
   work starts; `→ blocked` when something stops progress and the reason is
   non-trivial.
3. **On completion** the item moves to *Recently completed* with a
   one-sentence session summary and (optionally) a commit/PR link. The
   original entry is removed from *Open work items* in the same commit.
4. **On tabling** the item moves to *Tabled* with `Why tabled` and
   `Un-table when`.
5. ***Recently completed* ages out** at the next `/end-session` after an
   entry exceeds the rolling window. Git history retains everything; STATE.md
   stays scannable.
6. **Open questions** resolve into open work items when the answer changes
   what should happen next, or get removed if the answer is simply "no,
   don't do that."

## When to skip session-handoff skills

For very small tasks (one-line fixes, typo edits) the start/end-session
ceremony is overhead. Update STATE.md only if the work is non-trivial or
affects future sessions. Rule of thumb: if the change merits a commit, it
merits a STATE.md update; if it doesn't, it doesn't.

## Why this exists

The "platform isn't built in one session by one person" framing means every
agent at session N inherits the work of every prior session. The existing
primitives — ADRs for decisions, `git log` for history, user-local memory for
biographical context — leave a gap: the *currently-open* state of the
project. STATE.md fills that gap. It makes agent collaboration cheap by
replacing "re-read the repo" with "read one file."

## Related

- [`templates/state/STATE.md`](../templates/state/STATE.md) — the seed.
- [`.claude/skills/start-session.md`](../.claude/skills/start-session.md) —
  opens a session against STATE.md.
- [`.claude/skills/end-session.md`](../.claude/skills/end-session.md) —
  closes a session and updates STATE.md.
- [`claude-memory.md`](claude-memory.md) — the user-local memory system, a
  different concern.
- [`../docs/decisions/0008-session-handoff-state.md`](../docs/decisions/0008-session-handoff-state.md)
  — the architectural decision proposing this pattern.
