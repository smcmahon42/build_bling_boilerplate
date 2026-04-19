# Observability practices

Tracing, metrics, and their correlation with logs and agent primitives. The project uses OpenTelemetry as the tracing standard. See ADR [`0007`](../docs/decisions/0007-opentelemetry-and-correlation.md) for the decision.

Logging conventions live in [`logging-practices.md`](logging-practices.md) — this file is specifically about traces and metrics, and how they thread through the `trace_id` / `step_id` fields carried by `Provenance`.

## The correlation contract

Three identifiers thread through every request and every produced artifact:

| Identifier | Defined in | Populated in | Also appears on |
| --- | --- | --- | --- |
| `trace_id` | OTel (32 hex chars) | Root span at ring 0 | Every `Provenance.trace_id`, every log line, every `Result.provenance.trace_id` |
| `span_id` | OTel (16 hex chars) | Each span | Log lines emitted inside that span (optional but preferred) |
| `step_id` | `Provenance.step_id` | Producer at its ring | The corresponding OTel span's attributes; log lines carry it too |

The effect: a single `trace_id` filter in the observability backend surfaces every log line, every span, and every Result contributing to one logical workflow — across processes, languages, and rings.

## What to span

Minimum spans per process:

1. **Outer dispatch.** The handler's top-level `Dispatch` (or the equivalent MCP tool-call entrypoint) is always a span. Attributes:
   - `agent.id` — `Provenance.produced_by.agent_id`
   - `task.kind`
   - `task.id`
   - `task.version`
   - `idempotency.hit` (boolean — whether the cache served the result)
   - `ring` — depth in the call tree
2. **Capability verification.** A child span named `capability.check`. Attribute `capability.action`, and on denial `error.code`.
3. **Producer invocation.** A child span named `producer.<kind>`. Attributes:
   - `result.status` (ok / partial / error)
   - On `status = error`: `error.code`, `error.retryable`
   - `metrics.tokens_in`, `metrics.tokens_out`, `metrics.wall_seconds`, `metrics.cost_usd` (when populated)
3. **Outbound calls.** Any HTTP / gRPC / MCP call to another ring is a span. Use the OTel instrumentation for the client library where available.

Spans that would be purely decorative are not worth the cost. A span is warranted when a question like "how long did X take?" or "how often does X fail?" gets asked during incident review.

## Propagating across rings

Outbound calls propagate the trace context using W3C TraceContext headers (`traceparent`, `tracestate`). The agent primitives reinforce this at the Task level: `Task.provenance.trace_id` is the same value as the outgoing `traceparent`'s trace-id, and `Task.provenance.parent_step_id` corresponds to the caller's current `span_id`.

Receivers:

- Extract OTel context from transport headers.
- Also extract from `Task.provenance` — whichever arrives first wins; both SHOULD agree.
- Reject the Task if the two disagree — that is a signal of a misbehaving or compromised caller.

## Metrics

OTel metrics complement traces. The minimum useful set at every producer:

| Metric | Type | Labels |
| --- | --- | --- |
| `tasks_dispatched_total` | counter | `kind`, `status` |
| `task_duration_seconds` | histogram | `kind`, `status` |
| `capability_denials_total` | counter | `action`, `code` |
| `idempotency_cache_hits_total` | counter | `kind` |
| `producer_tokens_in_total` / `producer_tokens_out_total` | counter | `kind`, `model` |
| `producer_cost_usd_total` | counter | `kind`, `model` |

Labels are low-cardinality. `task.id`, `result.id`, and `agent.id` go on spans (and log lines), not on metrics.

## Exporters, sampling, and backends

The **exporter** choice is per-project — OTLP over gRPC to a collector is the pragmatic default; vendors (Honeycomb, Grafana Tempo, Datadog) expose OTLP endpoints. Keep the exporter configuration in env vars so it swaps without code change.

**Sampling** is per-project too. Defaults:

- **Parent-based sampling** — a span inherits its parent's sampling decision. Prevents a trace with some spans sampled and some not.
- **Head-based at the entry ring** — the ring-0 dispatcher decides. 100% during early development; reduce (e.g., 10%) once volume warrants.
- **Always sample errors** — a tail sampler (or a router rule at the collector) keeps every trace with `error = true`.

## What NOT to trace

- **Tight inner loops.** Per-iteration spans drown the pipeline; reach for a counter + histogram instead.
- **Serialization / deserialization of small payloads.** Single-digit microsecond work does not warrant a span.
- **Pure pure-function helpers.** Span what does I/O or crosses a process/ring boundary; leave the rest to the language's profiler.

## Relationship to logs

- Every span-carrying handler writes log lines that include `trace_id` and the current `span_id`. This gives backends like Honeycomb / Tempo the ability to jump from a trace to the logs produced inside any of its spans.
- Logs are narrative ("cache populated", "retry triggered"); traces are timing. You need both; neither replaces the other.

## Adopting this in a service

Checklist when wiring OTel into a new or existing service:

1. Pick the OTel SDK for the language. Follow `logging-practices.md`'s structured-logger pattern for shared context.
2. Configure the OTLP exporter via env vars.
3. Initialize the tracer at `main.go` startup; shut down cleanly on exit so pending spans flush.
4. Wrap `handler.Dispatch` (or equivalent) in an outer span.
5. Instrument outbound HTTP / gRPC / MCP calls with the OTel client middleware.
6. Bind the logger's required fields (`trace_id`, `step_id`) to the current span on every entry point.
7. Add the minimum metric set from the table above.
8. Verify one full trace end-to-end: ring-0 call → ring-1 producer → ring-2 downstream tool. All three spans should appear in the backend with a shared `trace_id`.

## Related

- [ADR 0007](../docs/decisions/0007-opentelemetry-and-correlation.md) — decision record.
- [`logging-practices.md`](logging-practices.md) — the sister file; logs carry the same `trace_id`.
- [`agent-primitives.md`](agent-primitives.md) — `Provenance.trace_id` / `step_id` / `parent_step_id` are the anchors.
- [`testing-practices.md`](testing-practices.md) — integration tests SHOULD verify trace propagation end-to-end.
