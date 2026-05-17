# 0011. Agent-cost observability via STATE.md Cost signals

- **Status:** Accepted (pattern)
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0007](0007-opentelemetry-and-correlation.md); [0008](0008-session-handoff-state.md); [0009](0009-provenance-and-review-state-on-state-rows.md); [0010](0010-agent-autonomy-scope.md); `claude-instructions/agent-cost-observability.md`

## Context

A premise of this boilerplate is that agents should rely on skills (pre-
compressed knowledge) more than on recalling and re-reading the repo. That
intuition is unfalsifiable without measurement. Today the project has
runtime observability — `observability-practices.md` covers OpenTelemetry
spans, metrics, and correlation through services. It does not cover the
parallel concern: how much agent work each open work item consumes over
its lifetime.

The cost the boilerplate wants to make visible is "this feature took 5
sessions and 18 operator turns" — not in aggregate but *per feature*. That
attribution is what reveals workflow drift: a skill that used to land in
2 sessions now consistently takes 5, or an entry that quietly accrued 14
review turns when a similar one took 3.

Two constraints shape the decision:

1. **The agent does not have authoritative per-session token counts.**
   Those live in the harness (Claude Code's session log, the provider's
   billing export, a custom OpenTelemetry proxy). The boilerplate cannot
   require the agent to know them.
2. **The boilerplate is language- and tooling-agnostic.** It cannot
   prescribe which external measurement system a project must use.

The right level of intervention is therefore: define the *schema* that
projects target, document the *proxies* the agent can record on its own,
and leave the absolute-cost wiring as a project-level concern.

## Decision

**STATE.md entries gain an optional `Cost signals` block, recorded by
`/end-session` on every multi-session entry.**

Block schema:

```
- **Cost signals:** (optional, best-effort)
  - **Sessions to date:** <integer>
  - **Skills used (cumulative):** <list>
  - **Operator turns:** <integer>
  - **Context read (approx):** <files-touched count>, <lines if known>
  - **Notes:** <free-form — e.g., absolute token counts pasted in from
    external tooling, anomalies observed during the work>
```

Wiring rules:

1. **All fields are optional.** The two strongly encouraged ones are
   *Sessions to date* and *Operator turns* — both are cheap and high-
   signal.
2. **`/end-session` records on multi-session entries by default;** for
   single-session trivial work it skips. The operator can ask for full
   capture on any entry.
3. **External tooling supplies absolute numbers.** When a project wires
   in a billing export or OpenTelemetry proxy, the per-session token
   counts get appended to `Cost signals.Notes` keyed by the session id
   recorded in `Provenance` (per ADR-0009).
4. **The boilerplate does not prescribe a measurement tool.** The
   composition layer is the join between `Provenance.Session` and the
   external system's session-id field.

## Alternatives considered

- **Cryptographic / signed cost attribution.** Rejected: overkill at the
  current trust scope. STATE.md is committed to a single repo; signing
  cost reports would imply distrust between contributors that the
  boilerplate isn't built around.
- **Boilerplate-shipped measurement (require Claude Code session log
  parsing in `/end-session`).** Rejected: ties the boilerplate to a
  specific harness. Projects using other clients (Claude API directly,
  third-party orchestrators) would inherit a dependency they can't fulfill.
- **No agent-side recording — rely entirely on external dashboards.**
  Rejected: external dashboards don't know about STATE.md entries.
  Attribution requires the join key, and the join key has to be in
  both systems. Recording *something* on the entry side is necessary.
- **Token counts as required fields.** Rejected on the constraint above:
  the agent does not have them reliably. Required-but-unenforced fields
  decay into noise.
- **Separate cost-tracking file (e.g., `COST.md` or `.cost-log.jsonl`).**
  Rejected: adds a second source of truth for "what did this entry
  cost" that has to stay in sync with STATE.md. One file with optional
  fields is simpler than two files with a sync contract.
- **OpenTelemetry-only attribution (no STATE.md fields).** Rejected: OTel
  is excellent for trace-level joins but operators reading STATE.md
  shouldn't have to switch to a dashboard to know "is this entry
  expensive?" The summary signal lives on the entry; OTel carries the
  detail.

## Consequences

- **Easier.** Cost drift becomes visible per feature. The "skills over
  recall" intuition becomes testable: an entry that consistently
  consumes more sessions or operator turns than its peers signals a
  skill gap or a documentation hole. Future planners reading STATE.md
  see the actual cost history, not just the outcome.
- **Harder.** `/end-session` does more work per session. Multi-session
  entries must carry running totals correctly — incrementing each
  session, not double-counting. The risk is that `Cost signals` stays
  blank for too long because the discipline lapses; that failure mode
  is the same as STATE.md drift and gets the same fix (a future
  `/state-audit` skill).
- **Follow-ups.**
  - **`/state-audit` skill** to flag entries that appear multi-session
    but have empty `Cost signals`.
  - **Boilerplate example wiring** showing how to parse a Claude Code
    session log into `Cost signals.Notes` via `jq`. Optional add-on,
    not part of the core decision.
  - **Dashboard composability** — once projects accumulate enough
    `Cost signals` data, a simple `awk`/`jq` pass over STATE.md across
    sister projects can build a portfolio-wide cost picture without
    any vendor tool.
  - **Cost feedback into the autonomy doctrine** (ADR-0010) — L3-heavy
    workflows that consistently consume high operator-turn counts are
    candidates for promotion to L2 (or demotion if the cost is
    revealing a confidence gap). The cost signals are the input data.

## See also

- `claude-instructions/agent-cost-observability.md` — operating
  doctrine and the cost-signal schema.
- `claude-instructions/observability-practices.md` — the runtime
  analog this slice deliberately mirrors.
- [0008](0008-session-handoff-state.md), [0009](0009-provenance-and-review-state-on-state-rows.md),
  [0010](0010-agent-autonomy-scope.md) — the prior slices this
  observability layer composes onto.
