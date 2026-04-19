# contracts

Machine-readable contracts that are **the source of truth** for this project's interfaces. Code is generated from these; hand-written types that duplicate them are drift waiting to happen.

Three complementary forms live here. Pick by use case, not preference.

## Decision matrix

| Use case | Form | Subdir |
| --- | --- | --- |
| Public HTTP API (consumed by humans, browsers, external services, other teams) | **OpenAPI 3.1** | [`openapi/`](openapi/) |
| Task-kind input/output schemas (the catalog entries referenced from `Task.inputs` and `Result.output`) | **JSON Schema** | [`json-schema/`](json-schema/) |
| Internal agent-to-agent RPC (hot path, cross-language codegen, streaming, bidirectional) | **Protocol Buffers** | [`proto/`](proto/) |

A project may use one, two, or all three. The shapes compose — OpenAPI schemas `$ref` JSON Schemas; `.proto` messages mirror the JSON Schema shape by convention.

## The five rules

1. **Contracts are versioned.** Semver. Breaking changes bump major and require an ADR.
2. **Contracts reference [agent primitives](../agent-primitives/).** `Task`, `Result`, `Evidence`, `Provenance`, `CapabilityToken` are never redefined per service — they are referenced.
3. **Code is generated from contracts.** The project ships a `make contracts` (or equivalent) target that regenerates types, clients, and servers. Hand-written types that overlap with generated ones are a smell.
4. **CI verifies regeneration is stable.** A workflow regenerates and fails if the working tree diffs — catches the case where a contract changed but committed code did not.
5. **Contracts are the first thing to change.** A new endpoint, a new RPC, a new task kind — the contract lands first (as a failing test expecting it), code follows.

See ADR [`0004`](../../docs/decisions/0004-contracts-as-source-of-truth.md) for the full decision and alternatives.

## Adding a new contract

Run `/new-contract` (the skill at [`.claude/skills/new-contract.md`](../../.claude/skills/new-contract.md)). It asks which form, scaffolds from the matching template, registers the contract, and reminds you about codegen wiring.

Manual path: copy the relevant subdirectory's starter file, rename, edit, add a codegen target, file an ADR if the contract is load-bearing.

## Related

- [`agent-primitives/`](../agent-primitives/) — the typed primitives these contracts reference.
- [`claude-instructions/contract-discipline.md`](../../claude-instructions/contract-discipline.md) — operating manual.
- [ADR 0004](../../docs/decisions/0004-contracts-as-source-of-truth.md) — the decision record.
