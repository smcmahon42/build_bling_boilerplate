# mcp-server

Language-indexed MCP server skeletons. Each subdirectory is a working reference for one language; projects pick the one matching their stack, copy it, and bind to an MCP SDK.

## Subdirectories

| Language | Subdir | Status |
| --- | --- | --- |
| Go | [`go/`](go/) | Shipped. Builds standalone (has its own `go.mod`). SDK binding is a stub — pick a Go MCP SDK and wire it in `cmd/server/main.go`. |

More language skeletons (TypeScript, Python, Rust) can be added as projects need them. Each follows the same shape: a transport-agnostic core that translates tool calls to `Task` → `Result` via the shared agent primitives, plus an SDK-specific entrypoint.

## Why one per language

The shape of the pipeline (capability verification → idempotency check → producer dispatch → Result construction with Provenance) is the same across languages. The SDK integration and the idiomatic type system are what change. A per-language skeleton keeps each idiom natural without fragmenting the shape.

## Related

- [`../agent-primitives/`](../agent-primitives/) — the primitives every skeleton serializes.
- [`../contracts/`](../contracts/) — the contracts that define the task kinds a skeleton serves.
- [ADR 0003](../../docs/decisions/0003-agent-primitives-as-foundation.md) — why the primitives.
- [ADR 0004](../../docs/decisions/0004-contracts-as-source-of-truth.md) — why the contracts are the source of truth.
