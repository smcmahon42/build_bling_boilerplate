# mcp-server (Go)

A minimal, SDK-agnostic Go skeleton for an MCP server that speaks the agent-primitives language: every tool call becomes a `Task`, every response becomes a `Result`, every call validates a `CapabilityToken`, every Result carries `Provenance`.

## What this is (and isn't)

**Is.** A working reference shape: capability verification, idempotency cache, producer interface, Result construction with Provenance. The core (~300 LOC across `internal/`) is transport-agnostic — it takes a parsed tool call and returns a typed response.

**Isn't.** A production server. The MCP transport binding in `cmd/server/main.go` is a stub: pick an MCP Go SDK (the official `modelcontextprotocol/go-sdk`, `mark3labs/mcp-go`, or an equivalent) and wire its tool-call handler to the `handler.Dispatch` function. The capability verifier uses a placeholder signature check — swap in real crypto (Ed25519, a KMS, whatever) before production.

## Layout

```
cmd/server/main.go               # entrypoint — bind to your chosen MCP SDK here
internal/primitives/             # Go types mirroring the JSON Schemas (hand-written; projects may codegen instead)
internal/capability/             # capability token verification
internal/idempotency/            # in-memory idempotency cache (swap for durable store in prod)
internal/handler/                # tool-call → Task → Result pipeline
internal/producer/               # producer interface + a mock implementation
```

## Adoption path

1. Copy `templates/mcp-server/go/` into the project (e.g. as `cmd/mcp-server/`).
2. Replace the placeholder module path **`github.com/example/project`** in `go.mod` and every `import` site with the project's real module path (one find-and-replace across the tree). Replace the `{{AGENT_ID}}` literal in `cmd/server/main.go` with the agent's stable identifier.
3. Pick an MCP Go SDK. Bind its tool-call handler to `handler.Dispatch`.
4. Replace `capability.Verify` with real signature verification against your issuer's public key(s).
5. Swap `idempotency.MemoryCache` for a durable store (Redis, Postgres) if calls span processes.
6. Register real producers for the task kinds this server serves. One tool = one kind.
7. Decide if Go types should be codegenerated from the JSON Schemas (`go-jsonschema`, `quicktype`) — if yes, move `internal/primitives/` to the generated tree and keep the hand-written file as a fallback.

## Why SDK-agnostic

MCP Go SDKs are still stabilizing. The template separates the **shape** (Task/Result pipeline, capability, idempotency, provenance) from the **transport** so the shape survives an SDK swap. The SDK owns JSON-RPC framing, tool discovery, and connection management; this skeleton owns what happens *after* a tool call is parsed.

## Related

- [ADR 0003](../../../docs/decisions/0003-agent-primitives-as-foundation.md) — agent primitives.
- [ADR 0004](../../../docs/decisions/0004-contracts-as-source-of-truth.md) — contracts are the source of truth; Go types here mirror (or regenerate from) the JSON Schemas.
- [`../../agent-primitives/`](../../agent-primitives/) — the canonical schemas this skeleton mirrors.
- [`../../contracts/`](../../contracts/) — contracts referenced by the tool catalog.
