# 0006. Human surfaces are thin views over contracts, never the other way around

- **Status:** Accepted as a recommended pattern. Projects adopt this ADR by building CLIs, dashboards, and UIs as thin clients that bind to the project's contracts, and by refusing to let UI-shaped concerns leak into the contracts themselves.
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0003](0003-agent-primitives-as-foundation.md); [0004](0004-contracts-as-source-of-truth.md); `templates/human-surface/`

## Context

Every project that serves both humans and agents eventually faces a fork: is the user experience designed first (with the API retrofitted to match what the UI renders), or is the contract designed first (with UIs as thin skins that render whatever the contract returns)?

The UI-first path is intuitive — design a screen, then back-fill the API to serve it. It fails in predictable ways:

- **One-consumer lock-in.** The API shape reflects a specific screen; adding a second UI (mobile, CLI, embedded widget) requires either a parallel endpoint or a translation layer. Agents — the fastest-growing class of consumer — are the worst-served.
- **UI-shaped contracts.** Endpoints carry names, field groupings, and aggregations that exist because they matched a mockup. Six months later the mockup has changed and the API still bears its fossils.
- **Logic migrates into the UI.** When the API doesn't return "the number the dashboard wants," the UI computes it from other fields. That computation is now policy, and it lives in JavaScript rather than behind the contract.
- **Two backends.** Teams introduce a "UI backend" (BFF) to paper over the shape mismatch. That backend then accretes domain logic that duplicates or diverges from the real service.

The contract-first path inverts this. The contract is the truth; UIs render it. When a new UI needs a field the contract doesn't expose, the field is added to the contract — never inline-computed in the client, never proxied-and-enriched by a bespoke UI backend. Two UIs may render the same Result differently, and that is the whole point of having separable surfaces.

This pattern compounds with ADR 0004 (contracts are source of truth) and ADR 0003 (agent primitives). Because the contract emits typed Results, any surface — human UI, agent consumer, dashboard, event pipeline — reads the same shape.

## Decision

**Every human-facing surface in this project (CLI, web UI, TUI, dashboard, mobile app) is a thin view over the contracts. The surface's responsibilities are strictly:**

1. **Input ergonomics** — arguments, forms, prompts, validation of shape (not policy).
2. **Output ergonomics** — rendering, colors, tables, streams, layouts.
3. **Session concerns local to the surface** — cookies, bearer-token plumbing, keyboard shortcuts.

**The surface MUST NOT:**

- Compute domain values the contract doesn't return. If the UI wants "average handle time per agent," that field is added to the contract.
- Proxy-and-enrich via a dedicated UI backend. Generated clients from the contract, or direct calls, are the default.
- Bake policy decisions into client-side logic that a future UI would have to re-implement.

**Practical consequences:**

- **Generated clients are preferred.** The OpenAPI / protobuf contract generates a client; the CLI or UI wraps it. Hand-rolled HTTP is acceptable only for toy surfaces or until codegen is wired (see the `cli-go` template for the hand-rolled shape; treat it as a transitional example).
- **When a UI needs a new field, the contract changes first.** Follow ADR 0004's contract-change flow (versioning, ADR if breaking).
- **Two UIs rendering the same data differently is a feature, not a smell.** Separable surfaces is the whole point.
- **Thickness is a smell.** If a human-surface component grows past a small ratio of the service it fronts (rough heuristic: 10% of the service's LOC), inspect for migrated logic.

The same rule applies to dashboards that happen to be Grafana/Metabase/Retool — they are surfaces, the contract or the underlying data model is the truth. Pre-baked dashboard queries that compute domain aggregations belong in the contract (a metrics endpoint, a materialized view with a documented schema), not as SQL in the dashboard.

## Alternatives considered

- **UI-first: design the screens, retrofit APIs.** Rejected: produces the failure modes in Context. May be pragmatic for throwaway prototypes; do not carry into production.
- **BFF (Backend-for-Frontend) pattern.** Rejected as default: introduces a second backend that tends to accrete logic. Acceptable as a narrow *aggregation* layer (stitching two existing contracts) but never as a translation layer or a place for new policy.
- **Let each UI compute what it needs from raw data.** Rejected: replicates policy across clients; one UI disagrees with another; agent consumers re-derive. Centralize in the contract.
- **A single canonical UI owned with the contract.** Rejected: forecloses the right of separate surfaces to render the same data differently. Tying the canonical view to the contract is fine (docs can show a rendering); blessing it as *the* view is not.
- **Skip the rule for dashboards because "they're just reports."** Rejected: dashboards are surfaces. SQL-in-the-dashboard rots the same way UI-side computation rots. The escape hatch is deliberate: a dashboard exploring an un-promoted metric during investigation is fine; when the metric becomes a standing signal, it gets a contract.

## Consequences

- **Easier.** Multiple human surfaces (CLI + web + mobile) coexist cheaply. Agent consumers get the same shape for free. UI swaps don't require server changes. Post-mortems on "the dashboard showed X but the API said Y" disappear because there is one source. Codegen keeps clients in lockstep with the contract.
- **Harder.** Every "just display this calculated value" request becomes a contract change. Small UI improvements sometimes block on a contract PR. Teams must internalize the reflex: "where should this live?" → "behind the contract." Cultural work; not just technical.
- **Follow-ups.**
  - Adopting a codegen path (part of ADR 0004's build target) makes the thin-client rule cheap to follow; without it, teams reach for expedient hand-rolled clients.
  - A lint / architectural check (ADR 0002 pattern) can flag human-surface packages that import domain-logic packages directly, enforcing the thinness rule at the module boundary.
  - Dashboard queries that compute cross-cutting domain aggregations get promoted to contract endpoints in a tracked PR.

## See also

- `templates/human-surface/` — the thin-client templates.
- `templates/human-surface/cli-go/` — reference CLI implementation.
- [0003](0003-agent-primitives-as-foundation.md) — typed primitives the surface renders.
- [0004](0004-contracts-as-source-of-truth.md) — contracts as source of truth; this ADR applies that rule to human-facing surfaces specifically.
