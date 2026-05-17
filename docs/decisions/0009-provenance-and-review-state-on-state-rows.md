# 0009. Provenance and review-state fields on STATE.md rows

- **Status:** Accepted (pattern)
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0008](0008-session-handoff-state.md); [0003](0003-agent-primitives-as-foundation.md); [0007](0007-opentelemetry-and-correlation.md); `agent-instructions/session-handoff.md`

## Context

Slice 7 (ADR-0008) introduced STATE.md as the project-local handoff ledger.
Once both humans and agents start writing to it, two distinct problems emerge
that the original schema does not address.

**Problem one — authority ambiguity.** A row in STATE.md authored by an
agent without operator review is *evidence*: the agent noticed something
and wrote it down. The same row authored by the human operator is
*instruction*: this is what we're doing. Future agents and contributors
reading the ledger cannot tell these apart from the unaltered schema, which
means an agent observation can quietly graduate to operator-confirmed
direction by virtue of sitting in the ledger for a while. That is the
"agent-written memory becomes a hidden prompt" failure mode that
OB1/OpenBrain explicitly designs against — and that the user-local memory
discipline already encodes (agent-written memories require human
confirmation or trusted import before they can shape action).

**Problem two — opaque origins.** A row that reads "add rate-limit
middleware" tells the next session *what* but not *why* or *from where*.
When the next agent or operator wonders "why is this here?" they have to
re-derive the context — exactly the cost STATE.md was built to eliminate.
Skill name, session id, and the prompt that produced the entry are cheap to
record and high-value to read later.

## Decision

**Every STATE.md row carries a `Review` field, and agent-authored rows
carry a `Provenance` block.**

Review state machine (one of five values):

| Value | Meaning | Set by |
| --- | --- | --- |
| `unreviewed` | Agent-authored, not yet seen or accepted by the operator | Agent (default for new agent entries) |
| `confirmed` | Operator has reviewed and accepts the entry as direction | Operator only |
| `rejected` | Operator has reviewed and declined the entry | Operator only |
| `stale` | No longer relevant but kept for record | Operator (or agent, with operator confirmation) |
| `superseded` | Replaced by a newer entry | Agent or operator |

Defaults:

- **Human-authored entries** start `Review: confirmed`. The human authoring
  *is* the confirmation.
- **Agent-authored entries** start `Review: unreviewed` *unless* the agent
  is acting on explicit operator direction for *this entry*, in which case
  the direction is the confirmation.

Authority rules:

- **Only the operator transitions to/from `confirmed` or `rejected`** on
  entries the agent didn't author. This is the line between evidence and
  instruction.
- **The agent may mark its own entries `superseded`** when writing a
  replacement; the replacement starts `unreviewed` unless operator-directed.
- **The agent may not silently promote `unreviewed → confirmed`** on its
  own entries. That promotion is the operator's act.

Provenance block (optional but strongly encouraged for agent-authored
entries):

```
- **Provenance:**
  - **Skill:** <skill name, e.g. `/security-review`, or `manual`>
  - **Session:** <session id, e.g. `2026-05-17-001`>
  - **Prompt summary:** <one line — what the operator asked when this came up>
  - **Context:** <file:line refs or related entry slugs>
```

The point is not audit-grade cryptographic signing — it is helping future
agents and operators answer "why is this here?" without re-reading the
whole repo. Skill name and session id alone usually carry most of the
value.

Skill changes (delivered in this slice):

- `/start-session` surfaces unreviewed agent-authored items alongside
  blocked items and includes their `Provenance` in the orientation report.
- `/end-session` sets `Review` per the defaults when capturing new entries,
  records `Provenance` for agent-authored entries, and updates `Review`
  transitions for any unreviewed entries the operator acted on during the
  session.

## Alternatives considered

- **One unified state field combining work-state and review-state**
  (e.g. `pending-unreviewed`, `in-progress-confirmed`). Rejected: the two
  concerns are orthogonal — an entry can be `in-progress` AND `unreviewed`
  (agent started exploring without explicit operator confirmation) or
  `pending` AND `confirmed` (queued by the operator). Combining them
  multiplies the value space and obscures the distinction the review state
  machine exists to enforce.
- **No review state — rely on `Opened by` alone.** Rejected: who authored
  the entry is not the same question as whether the operator has accepted
  it. An entry authored by `agent` and explicitly confirmed by the operator
  is different from one the operator hasn't seen, and the schema needs to
  carry that distinction.
- **No provenance — let `git blame` carry origin.** Rejected: `git blame`
  on STATE.md tells you *who edited* a row, not *which skill or prompt*
  produced the underlying observation. Skill + session + prompt summary
  capture the agent-side context that git can't see.
- **Cryptographically signed provenance (per OWASP Agentic ASI02
  "memory poisoning" guidance).** Rejected for now: overkill for a markdown
  ledger inside a single repo. Worth revisiting when STATE.md crosses trust
  boundaries (e.g., mirrored across orgs, fed to autonomous agents acting
  on external systems). The current threat model — a small team or solo
  builder collaborating with their own agents — does not require it.
- **Separate `agent_memories`-style table or sidecar file.** Rejected:
  adds tooling complexity (parsing, schema migration, sync with STATE.md)
  for a benefit that fields-on-rows already deliver. The file-system
  primitive stays primary; if cross-project agent memory becomes a need,
  that is a different ADR.

## Consequences

- **Easier.** Agent suggestions no longer drift silently into operator-
  confirmed direction. The operator gets a clear "what did the agent
  propose that I haven't reviewed" view at every `/start-session`.
  Future agents and contributors can answer "why is this here?" by reading
  the entry's Provenance instead of re-deriving from the repo.
- **Harder.** Two more fields per entry, plus an optional block. The
  schema is bigger; agents must remember to set defaults correctly and
  record Provenance. The `/end-session` skill carries enforcement, but
  drift is possible — a future `/state-audit` skill (post-slice-10) can
  detect entries with missing Provenance or invalid Review transitions.
- **Follow-ups.**
  - **Slice 9** layers an agent-autonomy scope policy on top of this —
    *which* STATE.md mutations the agent can make autonomously vs which
    require operator confirmation.
  - **Slice 10** wires token-budget observability into Provenance — each
    entry can carry the rough token cost of the work it represents.
  - A future `/state-audit` skill can flag agent-authored entries with
    missing or invalid Provenance, and detect Review-state transitions
    that violate authority rules.

## See also

- [0008](0008-session-handoff-state.md) — the foundational STATE.md ADR
  this slice extends.
- `agent-instructions/session-handoff.md` — operating doctrine, updated
  in this slice.
- `templates/state/STATE.md` — schema, updated in this slice.
- `.claude/skills/start-session.md`, `.claude/skills/end-session.md` —
  skills updated in this slice.
