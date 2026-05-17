# AGENTS.md — {{PROJECT_NAME}}

> **For the setup agent:** if you are reading this right after cloning the template, run the bootstrap workflow from `.claude/skills/bootstrap-project.md` if you are in Claude Code, or follow its steps manually from any other agent client. It replaces placeholders, trims sections, and seeds optional client-local memory.

This file is the **agent-neutral router**. Every agent client — Claude Code, Codex, a homegrown agent, or a hybrid orchestrator — should read it at session start, then fetch the relevant topic files in `agent-instructions/`.

Client-specific files are adapters, not doctrine:

- `CLAUDE.md` points Claude Code here and advertises Claude-only slash-command skills.
- Future Codex skills/plugins should point here instead of duplicating doctrine.
- Homegrown agents should implement this file plus the topic files as their repo contract.

---

## Project

- **Name:** {{PROJECT_NAME}}
- **One-liner:** {{PROJECT_DESCRIPTION}}
- **Primary stack:** {{PRIMARY_STACK}}
- **Status:** {{PROJECT_STATUS}}  <!-- e.g. exploring, alpha, production -->

## House rules (non-negotiable)

1. **TDD is the starting posture.** Every behavior change begins with a failing test unless the task is doc-only, generated output, or explicitly exempted. See `agent-instructions/testing-practices.md`.
2. **Security is a first-class gate.** OWASP Web Top 10, OWASP LLM Top 10, and OWASP Agentic Top 10 all apply. See `agent-instructions/security-practices.md` and run the applicable security-review workflow before shipping sensitive changes.
3. **Every change ships documentation.** If the change affects behavior, architecture, or operations, docs go in the matching bucket. See `agent-instructions/documentation-discipline.md`.
4. **Architectural decisions get an ADR.** See `docs/decisions/README.md` and use the ADR workflow to scaffold one.
5. **No business-specific code in shared layers.** Domain logic stays in its owned module.
6. **Trust the tests and the types; don't add validation that can't fail.** Validate at system boundaries only.
7. **Agent-written content starts as evidence, not instruction.** New agent-authored `STATE.md` entries default to `Review: unreviewed` until the operator confirms them.

## Process checkpoints

1. **Orient.** Read this file, then the matching topic files from the table below. If `STATE.md` exists, read it before broad repo exploration.
2. **Scope.** Classify the task as small / medium / large using `agent-instructions/development-workflow.md#task-scale-matrix`.
3. **Propose.** For medium or large tasks, share the plan before editing. For small tasks, proceed directly unless the action is L3 or L4.
4. **Test-first.** Write the failing test and confirm the expected failure when behavior changes.
5. **Edit.** Make the smallest coherent change. No unrelated cleanup.
6. **Verify.** Run targeted tests first, then the relevant local quality gate.
7. **Document.** Update the matching doc bucket in the same change.
8. **Commit policy.** Do not commit unless the operator explicitly asks, or the active client-specific adapter has an approved commit workflow.

For destructive actions, production writes, external messages, dependency changes, CI/settings/doctrine changes, pushes, or merges, follow `agent-instructions/agent-autonomy.md` and stop for explicit confirmation when required.

## When you're about to work on…

| Task | Read |
| --- | --- |
| A bug fix, feature, or refactor | `agent-instructions/development-workflow.md` |
| A commit or PR | `agent-instructions/commit-conventions.md` |
| Tests (unit, integration, e2e, Playwright) | `agent-instructions/testing-practices.md` |
| Security review or threat modeling | `agent-instructions/security-practices.md` |
| Adding a dependency | `agent-instructions/dependency-vetting.md` |
| Creating a module or boundary | `agent-instructions/component-explainability.md` |
| An agent-to-agent call, a new task kind, or a capability delegation | `agent-instructions/agent-primitives.md` |
| Adding or changing a contract (OpenAPI, JSON Schema, protobuf) | `agent-instructions/contract-discipline.md` |
| Logging, redaction, error-line shape, correlation fields | `agent-instructions/logging-practices.md` |
| Tracing, metrics, OpenTelemetry wiring across rings | `agent-instructions/observability-practices.md` |
| Measuring agent work per feature (sessions, turns, cost signals) | `agent-instructions/agent-cost-observability.md` |
| An architectural decision | `docs/decisions/README.md` |
| Planning an epic or filing issues | `agent-instructions/epics-and-projects.md` |
| Writing or editing client-local operator memory | `agent-instructions/operator-memory.md` |
| Starting or ending a work session (read or update `STATE.md`) | `agent-instructions/session-handoff.md` |
| Deciding whether an action is autonomous or needs operator confirmation | `agent-instructions/agent-autonomy.md` |
| Splitting agent instructions as the project grows | `agent-instructions/scaling-agent-instructions.md` |
| Writing or organizing docs | `agent-instructions/documentation-discipline.md` |

## Adapter policy

- Shared doctrine lives in `AGENTS.md`, `agent-instructions/`, `STATE.md`, ADRs, and templates.
- Client adapters may add invocation mechanics, permission config, and local-memory setup.
- Adapters must link to shared doctrine rather than copy it.
- If two agents disagree, the operator resolves the conflict and records durable decisions in an ADR or confirmed `STATE.md` entry.
