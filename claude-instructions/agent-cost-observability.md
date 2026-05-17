# Agent cost observability

`claude-instructions/observability-practices.md` covers *runtime* observability —
tracing requests through services with OpenTelemetry. This file covers the
parallel concern: *agent operation* observability — measuring how much
agent work each open work item consumes over its lifetime.

The motivation is the "rely on skills more than recall" intuition. Skills
are pre-compressed knowledge; recalling and re-reading the repo is the
expensive alternative. Without measurement, that intuition stays a vibe.
With measurement, cost drift becomes visible — and projects can spot
when a workflow has degraded into "re-read everything every session."

## What's measurable, what isn't

The agent itself does not have authoritative per-session token counts.
Those numbers live in the harness (Claude Code's session log, the
provider's billing export, an OpenTelemetry trace from a custom proxy).
This topic file describes the *proxies* an agent can record in STATE.md;
projects compose those with their preferred external measurement layer
for absolute numbers.

**Reliably available to the agent:**

- Count of skills invoked in a session.
- Count of files read and approximate line totals.
- Count of operator turns (messages from the human).
- Count of sessions an entry has been touched in (cumulative).
- Whether the work was small / medium / large per the task-scale matrix.

**Not reliably available to the agent — defer to external tooling:**

- Exact input and output token counts.
- Exact wall-clock duration of the conversation.
- Dollar cost per session.
- Cache hit rate.

## Cost signals on STATE.md entries

Every entry in *Open work items*, *Recently completed*, and *Tabled* may
carry an optional `Cost signals` block, captured by `/end-session`:

```
- **Cost signals:** (optional, best-effort)
  - **Sessions to date:** <integer>
  - **Skills used (cumulative):** <list of skill names>
  - **Operator turns:** <integer>
  - **Context read (approx):** <files-touched count>, <total lines if known>
  - **Notes:** <free-form — e.g., "rewrote twice", "blocked then recovered">
```

All fields are optional. The two most useful in practice are *Sessions to
date* and *Operator turns* — both are cheap to record and reveal cost
drift quickly. An entry that took 5 sessions and 18 operator turns to
land is not free; that signal should be visible to the next planner.

For absolute token / dollar measurements, projects layer in external
tooling (see below) and copy the figures into `Cost signals.Notes` at
session end.

## Composing with external tooling

The boilerplate is deliberately neutral on which external tool provides
measurement. Three common shapes:

- **Claude Code session log** — the harness writes structured logs of
  each session under `~/.claude/projects/<encoded-path>/`. A `jq` query
  at session end can extract approximate token counts and append them
  to `Cost signals.Notes`.
- **Provider billing export** — Anthropic / OpenAI billing exports
  carry exact per-request token counts. A nightly aggregation can
  attribute spend by session id (recorded in Provenance, see ADR-0009)
  back to the STATE.md entry.
- **Custom OTel proxy** — if requests are routed through a proxy that
  emits OpenTelemetry spans (see `observability-practices.md`), the
  trace itself carries token counts, latency, and cache-hit attributes.
  Pair `session id` in Provenance with `trace_id` in the proxy logs
  for end-to-end attribution.

All three are project-level wiring; the boilerplate ships the *schema*
that those wirings target.

## When to record, when to skip

Cost signals pay off when:

- An entry is multi-session — sessions-to-date alone is then a useful
  signal.
- A skill is suspected of degrading (taking more turns than it used to)
  and you want a baseline.
- A project is approaching a budget ceiling and per-feature attribution
  is needed.

Skip when:

- The entry is one-session and trivially scoped — the recording
  overhead exceeds the value.
- External tooling already attributes cost at the right granularity
  (e.g. a dashboard keyed by session id is sufficient on its own).

Default for the boilerplate's `/end-session` skill: record *Sessions to
date* and *Operator turns* on every multi-session entry; record the
rest only when the operator asks for full cost capture.

## Why this composes with prior slices

- **Slice 7 (STATE.md)** gives this slice an entry to attribute cost to.
- **Slice 8 (Provenance — session id)** gives cost signals a join key
  back to external billing / OTel data.
- **Slice 9 (Agent-autonomy)** identifies which operator turns happened
  at L3 (propose-and-confirm) — those are the most expensive turns in
  attention terms, and recording them per entry surfaces L3-heavy
  workflows that might benefit from being demoted to L2 (or promoted to
  L1) with documented confidence.

Without the prior slices, this one wouldn't have anywhere to write the
data or a way to correlate it to external measurements.

## Related

- `claude-instructions/observability-practices.md` — the runtime
  analog, covering OpenTelemetry traces and metrics through services.
- `claude-instructions/session-handoff.md` — STATE.md schema and the
  `Cost signals` block sits alongside `Provenance`.
- `claude-instructions/agent-autonomy.md` — L3 operator turns are the
  costliest in attention and are worth tracking.
- `.claude/skills/end-session.md` — the skill that records cost
  signals.
- `docs/decisions/0011-agent-cost-observability.md` — the architectural
  decision proposing this.
