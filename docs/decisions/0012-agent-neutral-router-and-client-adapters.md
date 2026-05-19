# 0012. Agent-neutral router with client-specific adapters

- **Status:** Proposed
- **Date:** 2026-05-17
- **Deciders:** Template maintainers
- **Related:** [0003](0003-agent-primitives-as-foundation.md); [0008](0008-session-handoff-state.md); [0009](0009-provenance-and-review-state-on-state-rows.md); [0010](0010-agent-autonomy-scope.md); [0011](0011-agent-cost-observability.md); `AGENTS.md`; `agent-instructions/`; `docs/agent-clients/README.md`

## Context

The template started with Claude Code as the primary agent client. That made the first workflow concrete: `CLAUDE.md`, `.claude/skills/`, Claude memory templates, and a Claude review workflow.

The expected use is broader. Operators may use Claude Code, Codex, a third-party agent, a homegrown orchestrator, or a hybrid of several clients. If each client gets its own repository or duplicated instruction tree, doctrine will drift: security rules, session handoff, autonomy levels, and ADR discipline will slowly diverge.

The repo needs one shared project contract with thin client adapters.

## Decision

Use `AGENTS.md` as the canonical agent-neutral router. Move shared topic guidance into `agent-instructions/`. Keep `CLAUDE.md` and `.claude/` as Claude Code adapters that point back to the shared doctrine. Codex and homegrown integrations must add adapter mechanics without duplicating shared rules.

## Alternatives considered

- **Claude-only repository.** Keep `CLAUDE.md` and `.claude/skills/` as the source of truth. *Why rejected:* it makes Codex and homegrown agents second-class and encourages duplicated prompt instructions.
- **One repo per agent client.** Maintain separate Claude, Codex, and custom-agent boilerplates. *Why rejected:* this guarantees drift in security posture, ADRs, templates, and session-handoff schema.
- **Lowest-common-denominator docs only.** Remove all client-specific adapters and keep only generic markdown. *Why rejected:* client ergonomics matter; Claude slash commands, Codex review flows, and custom orchestrators each need adapter surfaces.

## Consequences

- **Easier:** Claude, Codex, and homegrown agents can operate the same repo through one doctrine.
- **Easier:** Agent handoff stays durable through `STATE.md`, ADRs, tests, and contracts rather than private conversation context.
- **Easier:** New clients can be added as adapters without re-litigating project rules.
- **Harder:** Documentation has to maintain a clear line between shared doctrine and adapter-specific mechanics.
- **Harder:** Existing Claude-branded references must be reviewed and either neutralized or intentionally marked as Claude adapter details.
- **Follow-ups:** Add Codex adapter examples when the first concrete Codex workflow is worth standardizing; keep `.claude/skills/` as an adapter, not a shared doctrine layer.

## See also

- `AGENTS.md` — canonical router.
- `CLAUDE.md` — Claude Code adapter shim.
- `agent-instructions/README.md` — shared topic index.
- `docs/agent-clients/README.md` — operating strategy for Claude, Codex, hybrid, and homegrown agents.
