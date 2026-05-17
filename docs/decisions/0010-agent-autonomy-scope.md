# 0010. Agent-autonomy scope — a four-level doctrine

- **Status:** Accepted (pattern)
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0008](0008-session-handoff-state.md); [0009](0009-provenance-and-review-state-on-state-rows.md); `agent-instructions/agent-autonomy.md`; `CLAUDE.md`

## Context

This boilerplate ships agents into projects where they will perform many
classes of action — read files, edit code, run tools, modify settings,
merge branches, push to remotes. The risk profile of those actions varies
enormously: reading a file is cheap and reversible, force-pushing to main
is irreversible.

Today the doctrine that governs these distinctions is scattered:

- **`CLAUDE.md` house rule "Process checkpoints"** says "For destructive
  actions (deletions, force-pushes, production writes, external messages,
  schema drops), stop and confirm before acting." That covers the
  highest-risk operations but leaves the middle ground unaddressed.
- **`.claude/settings.json`** carries an allowlist of safe read-only
  Bash commands. That is enforcement for the lowest-risk level, but
  there's no corresponding doctrine that explains *why* those particular
  commands are on the list or what should be on the list as the project
  grows.
- **`agent-instructions/development-workflow.md` task-scale matrix**
  governs "share the plan before editing" for medium/large tasks. That
  is autonomy guidance, but framed as a scale concern rather than a
  capability concern.
- **OWASP Agentic ASI04 (Privilege escalation)** and **ASI08 (Unauthorized
  external actions)** flag the gaps explicitly: an agent operating
  without a documented autonomy doctrine is a privilege-escalation
  vector by default.

The OB1/Open Brain pattern of scope discipline (`personal → channel →
project → workspace`) is the inspiration, but it's framed for memory
access, not capability. The boilerplate-relevant version is simpler:
agent-autonomous / agent-proposes-human-confirms / human-only — plus a
fourth level that captures "autonomous on a feature branch with plan
review for non-trivial tasks," which is where the existing house rule 3
already lives.

## Decision

**Every action an agent can perform falls into exactly one of four
autonomy levels, documented in
`agent-instructions/agent-autonomy.md`.**

| Level | Operator role | Containment boundary |
| --- | --- | --- |
| **L1 — Autonomous** | None | Read-only, no side effects |
| **L2 — Autonomous on feature branch; plan-first on medium/large** | Plan review for medium/large tasks per house rule 3; none for small | Feature branch |
| **L3 — Propose and confirm** | Explicit per-action confirmation; silence is not consent | Action-by-action |
| **L4 — Human-only** | Agent refers the operation; never performs | Hard refusal |

Rules:

1. **Default upward.** When a new operation type doesn't fit cleanly,
   default to the higher level. Promoting an action *down* (L4→L3,
   L3→L2) requires an ADR; promoting *up* only requires updating
   `agent-autonomy.md`.
2. **One-shot exceptions are scoped.** The operator can grant
   "go ahead just this once" on a specific action. The exception does
   not generalize; the next action of the same type re-enters the
   doctrine. The exception is logged in the relevant STATE.md entry's
   `Notes` for traceability.
3. **`.claude/settings.json` enforces L1.** The allowlist is the
   machine-readable expression of L1; any command not on the allowlist
   triggers a confirmation prompt and falls under L2+ doctrine.
4. **L4 includes a specific list of recurring high-blast-radius
   operations** — force-pushes to main, production data writes, cloud
   infra changes, external messages, financial actions, license
   modification. The list lives in `agent-autonomy.md` and is the
   normative reference.
5. **STATE.md authorship rules (ADR-0009) compose with autonomy levels.**
   An L2 agent edit that adds an unreviewed STATE.md entry is consistent
   doctrine: autonomous at the action level, evidence-not-instruction at
   the data level. The operator's L3-style confirmation on the STATE.md
   entry happens at the next `/start-session`.

## Alternatives considered

- **Two levels only (autonomous / require-confirmation).** Rejected:
  collapses meaningfully different operations into one bucket. Editing
  code on a feature branch is not in the same risk class as merging to
  main; force-pushing to main is not in the same risk class as merging.
  Two levels would either over-restrict editing or under-restrict
  pushing.
- **Per-tool ACLs (allow list / deny list in `.claude/settings.json` only,
  no doctrine).** Rejected: settings.json captures *what's allowed by
  default*, but it doesn't explain *why* or guide what to add next. New
  projects derived from this boilerplate would re-derive the doctrine
  from scratch every time. The settings.json allowlist remains the
  enforcement for L1; the doctrine in `agent-autonomy.md` is what tells
  contributors how to expand the list.
- **Role-based access (read-only agent, contributor agent, maintainer
  agent).** Rejected: roles introduce identity management complexity
  that is overkill for solo and small-team projects. The four-level
  doctrine maps onto a single agent identity and lets the operator
  selectively grant one-shot exceptions.
- **Trust by intent ("the agent should infer the right level from
  context").** Rejected: that is exactly the OWASP Agentic ASI04
  failure mode. Inference is not a doctrine; explicit rules are.
- **OB1's scope hierarchy (personal / channel / project / workspace)
  applied directly.** Rejected: it's a memory-access model, not a
  capability model. The boilerplate-relevant question is what an agent
  can *do*, not where its writes propagate.

## Consequences

- **Easier.** Future agents (this one or others) reading the repo learn
  what they can do autonomously vs what needs operator turn-taking from
  a single page. New operation types get classified by the contributor
  who adds them, with the doctrine to ground the call. The "agent did
  the wrong thing autonomously" failure mode gets a documented escape
  hatch (one-shot exceptions) instead of a permanent override.
- **Harder.** Every new operation type costs a doctrine decision —
  which level does this belong to? That cost is paid once per
  operation, not once per task; net favorable, but real.
- **Follow-ups.**
  - **`.claude/settings.json` expansion** to enforce L1 more
    comprehensively and add denials for L4 — tracked as a separate
    slice once the doctrine settles.
  - **Hook integration** — pre-commit / pre-push hooks could enforce
    L3/L4 boundaries automatically (e.g., refuse a push to main
    without a corresponding merge commit reference).
  - **`/agent-autonomy-check` skill** — a future skill that classifies
    a proposed action against the doctrine before the agent performs
    it. Useful for ambiguous edge cases.
  - **Slice 10** adds token-budget observability that composes with
    this — operator turns at L3 are themselves observable events worth
    measuring.

## See also

- `agent-instructions/agent-autonomy.md` — the operating doctrine.
- `CLAUDE.md` — house rule 3 (process checkpoints), updated in this
  slice to reference the new topic file.
- `agent-instructions/development-workflow.md` — task-scale matrix
  the L2 boundary inherits from.
- `.claude/settings.json` — the L1 enforcement allowlist.
- [0008](0008-session-handoff-state.md), [0009](0009-provenance-and-review-state-on-state-rows.md)
  — STATE.md authorship and Review rules compose with this doctrine.
