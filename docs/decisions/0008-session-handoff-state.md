# 0008. Project-local session-handoff state via STATE.md

- **Status:** Accepted (pattern)
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0001](0001-record-architecture-decisions.md); [0003](0003-agent-primitives-as-foundation.md); `agent-instructions/session-handoff.md`; `agent-instructions/operator-memory.md`

## Context

Projects built with this boilerplate assume agent collaboration across many
sessions and contributors. Today the repo already carries three forms of
durable context:

- **CLAUDE.md + `agent-instructions/`** — instruction surface, mostly static,
  read at the top of every session.
- **`docs/decisions/`** — append-only Architectural Decision Records.
- **User-local memory** at `~/.claude/projects/<encoded-path>/memory/` —
  biographical and time-bound project context for an individual operator.

Missing from that set: the *currently-open* operational state of the
project — what work is in flight, what's blocked, what got tabled and why,
what questions are waiting on the operator. Without this, every new session
pays the cost of re-reading the repo to reconstruct ground truth. That cost
shows up in two ways: tokens consumed before any useful work begins, and
silent drops where work-in-progress is forgotten between sessions.

ADRs are too heavy for sub-decision-level state. User-local memory is the
wrong scope — `operator-memory.md` explicitly says "transient task state —
use Claude's task tracker, not memory." `git log` records what *was* done
but not what's *currently open*. External issue trackers (GitHub Issues,
Linear) are valid mirrors but introduce a network round-trip and a tooling
dependency that a fresh clone shouldn't require.

The pattern OB1 (Open Brain) uses — a project-local `.planning/STATE.md`
referenced from their `CLAUDE.md` — solves this with a file-system primitive
that any agent or contributor can read on first clone. This ADR adopts that
pattern with adjustments for the boilerplate's language-agnostic posture and
its "platform not built in one session by one person" framing.

## Decision

**Every project derived from this boilerplate carries a committed `STATE.md`
at the project root**, holding four sections — Open work items, Recently
completed, Tabled, Open questions — under the schema documented in
`agent-instructions/session-handoff.md`.

Paired skills `/start-session` and `/end-session` make reading and updating
STATE.md a routine session boundary, not a manual chore.

Wiring rules:

1. **STATE.md is committed**, not gitignored. The whole point of the handoff
   is legibility across humans and agents on first clone.
2. **STATE.md is project-local**, not user-local. It belongs to the project,
   not the operator. User-local memory continues to hold biographical and
   time-bound facts; STATE.md holds operational state.
3. **Sub-decision-level state lives here**, not in ADRs. An open work item
   that crystallizes into an architectural decision *graduates* into an ADR
   and is removed from STATE.md.
4. **The schema is fixed at boilerplate level.** Open work-item fields:
   `Summary`, `State` (pending | in-progress | blocked), `Opened by`,
   `Opened at`, `Last touched`, `Next step`. Project-specific extensions are
   allowed; the base fields are required.
5. **`/start-session` reads STATE.md before any other repo exploration.** That
   is how this primitive earns its keep — by replacing "re-read the repo"
   with "read one file."
6. **`/end-session` updates STATE.md.** Completed items move to recent;
   tabled items move to tabled with a reason; new items get added; open
   questions get recorded. STATE.md commits are separate from code commits.
7. **Recently completed is a rolling window** — ~10 entries or 14 days,
   whichever is shorter. Older entries age out. Git retains everything;
   STATE.md stays scannable.

## Alternatives considered

- **User-local memory (`~/.claude/projects/.../memory/`).** Rejected: the
  existing `operator-memory.md` explicitly draws the boundary — memory is for
  biographical and time-bound facts, not transient operational state.
  Putting STATE.md content there crosses that line and breaks
  cross-contributor legibility on first clone.
- **Embedding state inside CLAUDE.md.** Rejected: CLAUDE.md is instruction
  (read every session, mostly stable). Mixing mutable state into it either
  risks the agent rewriting instructions when it meant to update state, or
  causes CLAUDE.md to grow unbounded as history accumulates — exactly the
  token-cost problem this ADR is trying to reduce.
- **ADRs for in-progress work.** Rejected: ADRs are append-only records of
  *accepted decisions*. Most agent work doesn't reach that threshold; the
  friction of writing an ADR for "I started on the rate-limit middleware"
  suppresses useful state-keeping.
- **External issue tracker (GitHub Issues, Linear) as the primary store.**
  Rejected: introduces a network round-trip on every read and a tooling
  dependency that a fresh clone shouldn't require. Issue trackers remain
  useful mirrors for cross-team visibility; STATE.md stays the local source
  of truth.
- **`.planning/STATE.md` gitignored, OB1-style.** Rejected: OB1 chose
  gitignored because their `.planning/` directory carries maintainer-local
  plans they don't intend to share with contributors. This boilerplate's
  framing is opposite — making operational state legible across humans and
  agents on first clone is the whole point.

## Consequences

- **Easier.** Session N grounds itself in seconds by reading one file instead
  of grepping the codebase. Humans see "what's open" without asking an
  agent. Tabled work doesn't get silently lost. Cross-contributor handoff
  is a file read, not a conversation.
- **Harder.** STATE.md needs discipline to stay accurate — a stale ledger is
  worse than no ledger. The `/end-session` skill is the enforcement
  mechanism, but only if it actually runs. Drift is the main failure mode
  and warrants a future audit skill (see follow-ups).
- **Follow-ups.**
  - **Slice 8** adds provenance fields and a review state machine on STATE.md
    rows, leveraging the ledger as a carrier.
  - **Slice 9** defines agent-autonomy scope policy — what kinds of STATE.md
    updates an agent can make without human confirmation.
  - **Slice 10** wires token-budget observability into `/end-session`, using
    STATE.md entries as the per-feature accounting unit.
  - `/bootstrap-project` should copy `templates/state/STATE.md` to the
    project root on first run. Tracked separately to keep this slice
    scope-bounded.
  - A future `/state-audit` skill can detect drift between STATE.md and
    actual repo state.

## See also

- `agent-instructions/session-handoff.md` — operating doctrine and schema.
- `templates/state/STATE.md` — the seed.
- `.claude/skills/start-session.md` and `.claude/skills/end-session.md` —
  the paired session-boundary skills.
- `agent-instructions/operator-memory.md` — the user-local memory system this
  ADR carefully does *not* extend.
- [0007](0007-opentelemetry-and-correlation.md) — `Provenance` and
  correlation primitives that slice 8 will layer onto STATE.md rows.
