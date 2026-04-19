# 0005. Structured JSON logs with mandatory correlation fields and a redaction layer

- **Status:** Accepted as a recommended pattern. Projects adopt this ADR by wiring a structured logger that enforces the required fields and by installing the redaction layer as a boundary check.
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0003](0003-agent-primitives-as-foundation.md); `claude-instructions/logging-practices.md`; `claude-instructions/security-practices.md`

## Context

Multi-ring agent systems produce a lot of log lines — each agent on each ring emits its own — and the lines mean almost nothing in isolation. A `warn: retry triggered` without `trace_id`, `step_id`, and `agent_id` cannot be placed inside the call tree that produced it. Prose logs (`"processing user request, got 47 items back in 2.3s"`) defeat any structured query: finding "every capability.denied error in the last hour across all services" requires grep-and-parse rather than a log-store filter.

Three specific failure modes motivate this decision:

1. **Prose logs are uncorrelatable.** Nothing in the line ties it to the Task, Result, trace, or Provenance record that explains *why* the line exists. Post-mortems become archaeology.
2. **Log pipelines leak secrets.** Without a redaction discipline enforced at the logger, every new feature is a new chance to ship an auth token to the log store. Detection is post-hoc; containment is expensive.
3. **Different log shapes per environment.** Teams that render human-friendly prose in dev and JSON in prod hit the first production incident with code paths that were never exercised under the structured format.

The house rule "trust types and tests" applies here: if the log schema is enforced at the logger (required fields, redaction), call sites can trust the shape. If it's enforced by convention, it will drift.

## Decision

**Every project derived from this boilerplate emits structured JSON logs to stdout with a mandatory schema, via a logger that enforces the required fields and redacts known sensitive values before serialization.**

Required fields on every line: `ts`, `level`, `msg`, `agent_id`, `trace_id`, `step_id`. Level is one of `debug`, `info`, `warn`, `error` — no other values. `msg` is a short stable phrase; values are carried in their own fields, not interpolated into `msg`.

Logs mirror Provenance. `agent_id`, `trace_id`, and `step_id` on a log line are the same values as on the Provenance of the Result the line contributed to. This is the single correlation key across logs, traces, and Results — no separate join table, no bespoke glue.

Every `error` line additionally carries `err_code` (from `Result.error.code`), `err_msg`, `retryable`, and `task_id` when available. This makes error-code-based queries a primary tool in diagnosis.

A **redaction layer** runs before serialization and strips or masks known-sensitive fields: `authorization`, `api_key`, `password`, `token`, `signature`, any `*_secret` / `*_key`, and project-configured PII patterns. CapabilityTokens log only their `token_id`. `Task.inputs` and `Result.output` are not logged in full at `info`; a hash and size are logged instead, with full-body logging available at `debug` behind a capability-gated flag.

Log shape is identical across environments. Local development gets a pretty-printer over the JSON lines; the underlying stream stays structured.

The project chooses a mainstream structured logger for its language (`slog`, `pino`, `structlog`, `tracing`), wraps it with the required-field binding and redaction layer, and exposes it as the single project-wide logger. Call sites never build log lines by string concatenation.

## Alternatives considered

- **Prose logs with post-hoc parsing.** Rejected: grep-and-regex against unstable strings is brittle, expensive, and defeats multi-ring correlation. Works fine for single-process CLIs; fails for agent systems.
- **Traces only; skip logs.** Rejected: traces capture spans and structured attributes well but are a poor fit for operational narrative ("cache populated," "retry triggered") and for values that don't cleanly attach to a span. Logs and traces are complements, both carrying the same `trace_id`.
- **Structured logs without a redaction layer (convention-only).** Rejected: every new feature is a new leakage risk. Redaction at the logger is the only place it reliably holds. Convention-only redaction fails the same way convention-only security checks fail (ADR 0002).
- **A custom home-grown logger.** Rejected: mainstream loggers (per-language) get correctness, performance, and ecosystem integration right. Wrap them; do not reinvent.
- **Different log schemas per service.** Rejected: cross-service queries and cross-ring tracing are the reason for the schema. One schema, one shape.
- **Log full Task/Result payloads at `info` level.** Rejected: produces large log volumes, maximizes PII/secret exposure surface, and buries the operational signal in payload noise. Hashes + sizes at `info`; full payloads at `debug` behind a flag.

## Consequences

- **Easier.** Cross-ring diagnosis is a single `trace_id` filter. Error-code queries surface patterns across services ("every `capability.denied` in the last hour"). Log storage costs drop because payload bodies aren't duplicated into logs. Agent consumers of logs (incident-triage bots, anomaly detectors) read a stable schema. Compliance audits find that logging is structurally bounded against secret leakage.
- **Harder.** Call sites must log with structured fields rather than prose ("value was 47" → `value=47`); small retraining cost. The redaction layer must be updated when new sensitive fields appear — this should be covered by boundary checks (ADR 0002 pattern) to avoid drift. Pretty-printer for local dev is a small extra component to maintain.
- **Follow-ups.**
  - A boundary-check enforcer (ADR 0002 pattern) can verify every logger construction uses the wrapped project logger, not the raw library.
  - Redaction rules are covered by unit tests feeding known-sensitive payloads; treat these tests as regression tests — never delete, only add.
  - Slice 6 (OpenTelemetry / correlation) binds `trace_id` across logs, traces, and Provenance automatically; once shipped, `trace_id` generation moves from ad-hoc to OTel-managed.

## See also

- `claude-instructions/logging-practices.md` — day-to-day operating manual.
- `claude-instructions/security-practices.md` — OWASP items around prompt/output logging and PII handling.
- [0003](0003-agent-primitives-as-foundation.md) — the Provenance fields that log fields mirror.
- [0002](0002-boundary-annotation-enforcement.md) — the pattern for enforcing "every logger goes through the wrapper" at every call site.
