# Logging practices

Structured logs by default. Every log line is a JSON object with a stable schema, emitted to stdout, consumed by whatever the project's log pipeline is. Prose logs defeat correlation across rings and defeat agentic consumers. See ADR [`0005`](../docs/decisions/0005-structured-logging-and-redaction.md) for the decision.

## The log line shape

Every log line, minimum:

```json
{
  "ts": "2026-04-19T14:22:08.123Z",
  "level": "info",
  "msg": "task dispatched",
  "agent_id": "iris.summarizer.v3",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "step_id": "step:abc123..."
}
```

Required fields:

| Field | Meaning |
| --- | --- |
| `ts` | RFC 3339 UTC. Millisecond precision minimum; nanosecond if the platform supports it. |
| `level` | One of `debug`, `info`, `warn`, `error`. Never `trace`, never `fatal` — traces belong in tracing, fatals belong in the exit code. |
| `msg` | A short, stable, kebab-or-phrase message. NOT a templated string with interpolated values — put values in their own fields. |
| `agent_id` | The `Provenance.produced_by.agent_id` that owns the line. |
| `trace_id` | The `Provenance.trace_id` for the current workflow. Empty string is not acceptable; generate one at ring 0. |
| `step_id` | The `Provenance.step_id` for the current production step. |

Additional context fields are per-line; prefer flat keys over nested objects where reasonable.

## Tying logs to Tasks, Results, and tracing

Log fields mirror Provenance fields. This is the single correlation key across **logs, traces, and Results**:

- A `Result` with `result_id = X` → logs carrying `result_id=X` are the logs produced while computing it.
- A trace span with `trace_id = T` → logs carrying `trace_id=T` are the logs emitted inside that span.
- A `Task` with `parent_step_id = S` → logs carrying `step_id=S` are the upstream caller that invoked this step.

Callers who wire OTel (slice 6) get matching `trace_id` and `span_id` values on both logs and traces; no separate glue.

## Levels — what goes where

- **`debug`** — noisy, developer-local detail (raw payloads, branch taken, cache hit/miss). Disabled in production by default; enabled per-request via a capability-gated flag when diagnosing.
- **`info`** — normal operational events. Task received, task dispatched, Result emitted, cache populated, tool started. The steady rhythm of healthy operation.
- **`warn`** — recoverable anomalies. Budget approaching limit, retry triggered, fallback path taken, capability token about to expire.
- **`error`** — unrecoverable-for-this-step failures. Accompany every `Result.status = error` (and map `Result.error.code` to a field on the log line).

No `fatal`. A process that cannot continue writes one final `error` log line and exits with a non-zero code.

## Error logging

Every `error` level log carries these fields:

| Field | Source |
| --- | --- |
| `err_code` | `Result.error.code` (dotted, stable, machine-readable). |
| `err_msg` | `Result.error.message` (prose, human-read). |
| `retryable` | `Result.error.retryable`. |
| `task_id` | Optional but preferred — the Task that triggered the error. |

This makes "show me every `capability.denied` in the last hour" a single query against the log store.

## Redaction

Secrets and PII do not reach the log pipeline. The logger has a **redaction layer** that runs before serialization:

- **Known-secret fields** — `authorization`, `api_key`, `password`, `token`, `signature`, any field matching `*_secret` or `*_key` — are replaced with a fixed marker (`"***REDACTED***"`).
- **Configurable PII fields** — per project, values matching known PII patterns (emails, phone numbers, national IDs) are either hashed or marked per project policy.
- **CapabilityTokens** — log only `token_id`, never the token body or `signature`.
- **Arbitrary payloads (`Task.inputs`, `Result.output`)** — do not log in full by default. Log a hash of the payload + the size. If diagnosing, elevate to `debug` with a capability-gated flag, and still redact at the field level.

Adding a redaction rule is a one-line change to the redaction layer; that layer is covered by unit tests that feed known-sensitive input and assert the marker appears.

## Context propagation

Logger instances carry context so the required fields don't have to be passed manually at every call site. Pattern:

1. At ring entry, the handler creates a logger bound to `trace_id`, `step_id`, `agent_id`.
2. The bound logger is passed (or placed in `context.Context` for Go) through the call chain.
3. Call sites log to the bound logger; the required fields are attached automatically.
4. On delegation to ring N+1, the bound logger re-binds with the new `step_id` and incremented `ring`.

This mirrors the Provenance rebinding rule in [`agent-primitives.md`](agent-primitives.md) — logs and primitives evolve together.

## What NOT to log

- **Entire request/response bodies at `info`.** Too noisy; too easy to leak. Log hashes + sizes; raise to `debug` only behind a flag.
- **Free-text status messages that parse out structured data.** `"processed 47 items in 2.3s"` → instead, log with `count=47 duration_seconds=2.3`.
- **Secrets, ever.** If a new field carries a secret, it is added to the redaction layer before the feature merges.
- **Stack traces at `info` or `warn`.** Reserve for `error`. Stack traces are expensive to store and search.
- **Anything that duplicates Result fields.** The Result is authoritative for its `output`/`error`; logs describe the *operational event*, not the domain result.

## Local development

In local runs, a pretty-printer may render the JSON lines as human-readable output — the *storage* format stays JSON, the *display* is a development convenience. Never generate different log structures for dev vs prod; that's the path to surprises in production.

## Library choice

This boilerplate is language-agnostic. Pick a mainstream structured logger for the project's language (`slog` in Go stdlib, `pino` in Node, `structlog` in Python, `tracing` in Rust) and wrap it with the redaction layer and required-field binding. Do not invent a logger.

## Related

- [ADR 0005](../docs/decisions/0005-structured-logging-and-redaction.md) — decision record.
- [`agent-primitives.md`](agent-primitives.md) — Provenance fields that log fields mirror.
- [`security-practices.md`](security-practices.md) — redaction ties into OWASP LLM items around prompt/output logging and PII handling.
- [`observability-practices.md`](observability-practices.md) — tracing and metrics; shares `trace_id` with logs.
