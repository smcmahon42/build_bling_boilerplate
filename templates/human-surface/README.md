# human-surface

Thin skins over the project's contracts. Humans interact through CLIs, dashboards, and UIs; none of those own domain logic. A human surface's job is input ergonomics (arguments, forms, prompts) and output ergonomics (rendering, colors, layouts) — nothing more. Business logic lives behind the contract; the surface just calls it.

See [ADR 0006](../../docs/decisions/0006-human-surfaces-are-views-over-contracts.md) for the decision.

## Subdirectories

| Surface | Subdir | Status |
| --- | --- | --- |
| CLI (Go) | [`cli-go/`](cli-go/) | Shipped. Minimal single-command CLI that POSTs a Task and renders the Result. |

More surfaces (TUI, web dashboard, mobile app) can be added as projects need them. Each follows the same rule: **the surface is a view; the contract is the truth.**

## The cardinal rule

- **A new button, tab, or flag MUST NOT introduce new server logic.** If it does, the logic belongs in the contract first and the UI-facing endpoint is already there.
- **If the UI needs data the contract doesn't return, add the field to the contract.** Don't inline-compute it in the client. Don't proxy-and-enrich from a "UI backend."
- **If two surfaces render the same data differently, that's fine.** It's what surfaces are for. Nothing in the contract changes.

## Related

- [`../contracts/`](../contracts/) — the source of truth surfaces render.
- [`../agent-primitives/`](../agent-primitives/) — same primitives, different consumer.
- [ADR 0006](../../docs/decisions/0006-human-surfaces-are-views-over-contracts.md) — the decision record.
