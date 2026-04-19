# 0007. OpenTelemetry for tracing and metrics, with Provenance as the correlation anchor

- **Status:** Accepted as a recommended pattern. Projects adopt this ADR by wiring an OTel SDK at startup, instrumenting Dispatch and outbound calls, and verifying `trace_id` flows end-to-end through Provenance and logs.
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0003](0003-agent-primitives-as-foundation.md); [0005](0005-structured-logging-and-redaction.md); `claude-instructions/observability-practices.md`

## Context

An agent system composed across multiple rings and services produces a call graph that is invisible without correlated traces. When ring 0 calls ring 1 calls ring 2 calls a tool, and the tool fails, the post-mortem needs to answer: which ring-0 request triggered this? What did ring 1 decide? Was ring 2 also slow? Did the same trace fan out to other leaves that failed in ways nobody noticed?

Three failure modes follow from not having a standard tracing story:

1. **Per-service divergence.** Each service invents its own tracing shape (or skips tracing entirely). Cross-service joins require bespoke glue per pair.
2. **Identifier sprawl.** Without a shared convention, each log store and each trace backend uses different keys for the same thing. Correlating logs to traces becomes a heroic act.
3. **Vendor lock-in by accretion.** Teams that reach for a vendor-specific SDK find later that switching is a rewrite, not a config change.

OpenTelemetry (OTel) solves the protocol side: a vendor-neutral API, a stable wire format (OTLP), instrumented libraries for most languages, and a collector pattern that routes to many backends. What OTel does not prescribe is *what* to span, *which* attributes to attach, or *how* the tracing identifiers relate to the project's domain primitives.

ADR 0003 already defined the domain's identifiers — `Provenance.trace_id`, `Provenance.step_id`, `Provenance.parent_step_id`. This ADR commits to making those identifiers the same ones OTel propagates, so the two systems never drift.

## Decision

**Every project derived from this boilerplate uses OpenTelemetry for distributed tracing and metrics, with these wiring rules:**

1. **`Provenance.trace_id` is an OTel trace-id.** 32 hex chars. Generated at the ring-0 entry point; propagated unchanged to every downstream ring.
2. **`Provenance.step_id` corresponds to an OTel span.** A producer emits one span per Dispatch; that span's id is reflected (in a stable form) in the Provenance the producer attaches to its Result.
3. **`Provenance.parent_step_id` corresponds to the parent span's id.** The delegation chain recovered from Provenance matches the span tree in the tracing backend.
4. **Outbound calls use W3C TraceContext.** `traceparent` and `tracestate` headers carry the OTel context; the receiving ring extracts, verifies against `Task.provenance.trace_id`, and attaches to its own span.
5. **The minimum spans are fixed.** Outer Dispatch, capability check, producer invocation, outbound calls. A project may add more; it may not skip these.
6. **The minimum span attributes are fixed.** `agent.id`, `task.kind`, `task.id`, `ring`, `idempotency.hit`. On error: `error.code`, `error.retryable`. `result.status` on completion.
7. **Minimum metrics are fixed.** Counters for tasks dispatched, capability denials, idempotency-cache hits. Histograms for task duration. Labels are low-cardinality (`kind`, `status`, `action`, `code`); high-cardinality identifiers stay on spans and log lines.
8. **Logs carry `trace_id` and `span_id`.** Per ADR 0005. This is the join key between the two systems.
9. **Disagreement between `traceparent` and `Task.provenance.trace_id` is a rejection.** Receivers treat the mismatch as a signal of a misbehaving or compromised caller.
10. **Exporter, sampling, and backend are project-level config.** OTLP to a collector is the pragmatic default; env vars drive the rest.

Projects adopting this ADR wire OTel in their `main.go` (or equivalent), instrument Dispatch and outbound calls, and verify one end-to-end trace before closing the adoption PR.

## Alternatives considered

- **Vendor-specific tracing (Datadog, NewRelic, Honeycomb SDKs directly).** Rejected as the default: switching vendors requires rewrites. OTel's collector pattern lets a project route to any of these backends via configuration, keeping the code portable.
- **Custom lightweight tracing (home-grown span ids in logs, no dedicated tracing backend).** Rejected: works for single-service systems; collapses at the multi-ring scale that the agent primitives are designed for. Loses span timing, loses fan-out visualization.
- **Tracing only; no metrics.** Rejected: per-span-rollup is expensive at scale and not the right aggregation for counts / rates / percentiles. Metrics and traces are complements.
- **Metrics only; no tracing.** Rejected: flat aggregates hide the causal chain that incident review needs. Metrics say *what* broke; traces say *why*.
- **Let each service pick its own tracing story.** Rejected: defeats cross-service correlation, which is the whole point.
- **Make `trace_id` independent of OTel trace-ids.** Rejected: two identifiers for the same concept drift within weeks. Make them the same value.

## Consequences

- **Easier.** End-to-end traces across rings and services work with a single trace-id filter. Logs, traces, and Results are joinable by `trace_id` with no bespoke glue. Vendor swaps are a collector-config change, not a code rewrite. New services inherit instrumentation from library-level OTel instrumentations rather than re-instrumenting by hand. Error-rate and latency dashboards fall out of the fixed metric set.
- **Harder.** Every service carries an OTel SDK dependency and a startup/shutdown lifecycle — cold-start cost is real, mostly small. Sampling decisions need deliberate tuning at scale. Disagreement-detection between `traceparent` and `Task.provenance.trace_id` is extra code at every boundary — an opportunity for an ADR-0002-style enforcer.
- **Follow-ups.**
  - A boundary-check enforcer verifies every ring entrypoint calls the trace-context extraction helper — drift-free via CI.
  - Tail-sampling rules (keep traces with `error = true`) land at the OTel collector layer, not in application code.
  - Cost-aware sampling (keep expensive LLM traces more often than cheap deterministic ones) is a project-level tuning decision once volume warrants.
  - Agent-facing dashboards (for automated incident-triage agents) read the same OTel backend that humans read — one set of instrumentation serves both consumers.

## See also

- `claude-instructions/observability-practices.md` — the operating manual; what to span, what attributes, metrics list.
- [0003](0003-agent-primitives-as-foundation.md) — `Provenance` is the domain-side anchor.
- [0005](0005-structured-logging-and-redaction.md) — logs carry the same `trace_id` / `span_id`.
- [0002](0002-boundary-annotation-enforcement.md) — pattern for enforcing "every ring entrypoint extracts trace context" at every call site.
