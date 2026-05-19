# Agent client strategy

This repo is designed for multiple agent clients working against one shared project contract. Do not fork the repo by LLM provider unless the project itself diverges.

## Strategy

| Layer | Shared or adapter-specific | Examples |
| --- | --- | --- |
| Project doctrine | Shared | `AGENTS.md`, `agent-instructions/`, ADRs |
| Session state | Shared | `STATE.md`, Review state, Provenance, Cost signals |
| Runtime contracts | Shared | Task / Result / Evidence / CapabilityToken schemas |
| Client invocation | Adapter-specific | Claude slash commands, Codex skills/plugins, homegrown commands |
| Client memory | Adapter-specific | Claude Code `~/.claude/.../memory/`, other local stores |
| Permissions / sandbox | Adapter-specific | `.claude/settings.json`, Codex sandbox config, internal ACLs |

The rule is simple: shared behavior goes in the repo contract; client mechanics stay in adapters.

## Shared workflows, adapter invocations

Shared docs refer to workflows by purpose. Each agent client can expose the
same workflow through slash commands, skills, plugins, prompts, scripts, or a
manual checklist.

| Workflow purpose | Shared trigger | Claude Code adapter | Other clients |
| --- | --- | --- | --- |
| Bootstrap a cloned template | First setup after using the template | `/bootstrap-project` | Follow the bootstrap steps in `README.md` and `AGENTS.md` manually or through a client-specific adapter |
| Start a work session | Non-trivial session start | `/start-session` | Read `STATE.md`, surface blocked / unreviewed / in-progress items, then ask for direction |
| End a work session | Non-trivial session close | `/end-session` | Reconcile `STATE.md`, move completed items, record open questions and cost signals |
| Security review | Security-relevant diff or pre-release review | `/security-review` | Walk `agent-instructions/security-practices.md` against the current diff |
| Dependency audit | Adding dependencies or scheduled dependency review | `/dep-audit` | Run the language-specific scanners and the five-point check in `dependency-vetting.md` |
| Test gap review | Release prep or suspected coverage gap | `/test-gaps` | Inspect changed behavior against `testing-practices.md` and list missing critical-path tests |
| New ADR | Architectural decision | `/new-adr` | Copy `docs/decisions/TEMPLATE.md`, assign the next number, update the ADR index |
| New contract | New OpenAPI / JSON Schema / protobuf contract | `/new-contract` | Copy the matching starter from `templates/contracts/`, wire codegen, update docs |
| New adapter skill | Recurring client-specific workflow | `/new-skill` | Add a thin client adapter that links back to shared doctrine |

Adapter documentation may name the concrete command. Shared doctrine should
prefer the workflow purpose unless the client-specific invocation is the topic.

## Supported operating modes

### Single-agent mode

One agent client works the repo. It reads `AGENTS.md`, follows `agent-instructions/`, and uses its own adapter workflows where available.

### Alternating-agent mode

Claude Code and Codex both work the repo over time. `STATE.md` is the handoff point. Agent-authored entries stay `Review: unreviewed` until the operator confirms them, regardless of which client wrote them.

### Hybrid mode

A human or orchestration layer uses multiple agents for different strengths: one for implementation, one for review, one for research. Agents exchange durable facts only through shared artifacts: PRs, ADRs, `STATE.md`, contracts, and tests. Do not rely on private conversation context crossing clients.

### Homegrown-agent mode

A custom agent should implement the same contract:

1. Read `AGENTS.md` first.
2. Load only relevant files from `agent-instructions/`.
3. Treat `STATE.md` as project-local handoff state.
4. Mark agent-authored observations as `Review: unreviewed`.
5. Follow the L1–L4 autonomy doctrine.
6. Emit typed Task/Result/Evidence/Provenance when delegating work across rings.

## Adapter requirements

Every adapter should answer these questions without changing shared doctrine:

- How does the operator invoke common workflows?
- Where does the client store local/private memory?
- What commands can the client run without confirmation?
- How does the client run review-only mode?
- How does the client record session handoff in `STATE.md`?

If the answer affects all agents, update `AGENTS.md` or `agent-instructions/`. If it affects only one client, update that adapter.

## Codex baseline

Codex can operate this template directly from `AGENTS.md`; a separate
repo-local `CODEX.md` is not required for the shared doctrine. Add Codex-specific
skills, plugin metadata, or local memory wiring only when they provide concrete
workflow mechanics, and keep those adapter files thin links back to
`AGENTS.md` and `agent-instructions/`.

## Conflict resolution

If two agents disagree:

1. Prefer tests, contracts, and cited evidence over model confidence.
2. Ask the operator to decide when the trade-off is product or doctrine-level.
3. Record durable decisions as ADRs.
4. Record temporary operational decisions in `STATE.md` with explicit `Review` state.

## Why not separate repos

Separate repos create drift in doctrine, templates, security posture, and ADR history. They make sense only if the generated projects themselves diverge by product, not because different agents operate them. The adapter pattern keeps one source of truth while allowing each client to expose its best workflow ergonomics.
