# 0004. Machine-readable contracts are the source of truth; code is generated

- **Status:** Accepted as a recommended pattern. Projects adopt this ADR by copying the relevant subtrees of `templates/contracts/` and wiring `make contracts` into their build.
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0003](0003-agent-primitives-as-foundation.md); `templates/contracts/`; `claude-instructions/contract-discipline.md`

## Context

Every project that exposes an interface — an HTTP API for users, an RPC for other services, a task-kind input/output pair for agents — faces the same fork in the road: is the contract defined *in code* (handlers, structs, decorators) with documentation rendered from the code, or is the contract defined *as a machine-readable artifact* (OpenAPI, JSON Schema, protobuf) with the code generated from it?

The code-first path is faster to start and lets each language feel idiomatic. It fails predictably at three points:

- **Cross-language consumers.** A Go-defined contract needs a hand-translated Python/TypeScript client, and those translations drift.
- **Cross-service agreement.** Two services that share a message shape each own a local copy; when one changes, the other discovers it at runtime.
- **Pre-code review.** The contract cannot be reviewed before the code exists, so breaking changes are found at merge time rather than at design time.

Agent-first systems feel all three problems more sharply than traditional services, because (a) agents generate code and need a stable target schema to hit, (b) contracts are often consumed by multiple runtimes simultaneously (MCP servers, OpenAPI clients, internal RPC), and (c) the task catalog — the registry of `kind → (input schema, output schema)` — is itself a cross-cutting contract that no single service owns.

The contract-first path inverts this: contracts are the artifact reviewers read, CI diffs, and codegen consumes. Generated types are disposable; the contract is the source of truth.

Three contract forms cover the common cases. Each is strongest in a different regime:

- **OpenAPI 3.1** for external HTTP APIs — browser clients, partner integrations, documentation rendered by SwaggerUI/Redoc.
- **JSON Schema** for task-kind input/output shapes and persistent-record shapes — the forms that plug into `Task.inputs` and `Result.output` from ADR 0003.
- **Protocol Buffers (proto3)** for internal agent-to-agent RPC — hot paths, streaming, cross-language codegen.

These are complementary, not competitive. OpenAPI operations `$ref` JSON Schemas; protobuf messages mirror JSON Schema shapes. A project may use all three without contradiction.

## Decision

**Every project derived from this boilerplate treats machine-readable contracts as the source of truth for its interfaces, co-located under `contracts/` (copied from `templates/contracts/`).**

1. **Form selection follows the decision matrix** in `templates/contracts/README.md`: OpenAPI for external HTTP, JSON Schema for task kinds and generic data, protobuf for internal RPC. A contract uses exactly one form; that form is the source, any other representation is generated.
2. **Contracts reference the agent primitives** (ADR 0003) via `$ref` (JSON Schema / OpenAPI) or by mirroring the message shape (protobuf). The primitives are never redefined per service.
3. **Codegen is a single project-level command** — `make contracts` or equivalent. It regenerates server stubs, client SDKs, and shared types. Hand-written types that overlap with generated ones are a drift smell.
4. **CI diffs the working tree after regeneration.** A workflow runs `make contracts` and fails the PR if the working tree changes — this catches "I edited the contract but forgot to regenerate" at PR time rather than runtime.
5. **Contracts are versioned.** Semver. OpenAPI versions appear in the URL path (`/v1/…`); proto versions appear in the package name (`example.v1`); JSON Schema versions appear in the `$id`. Breaking changes bump major *and* require an ADR describing the migration.
6. **Contracts change first.** A new endpoint, RPC, or task kind lands as a contract edit (with a failing test that expects the new shape), then code follows.

Scaffolding a new contract is a `/new-contract` skill invocation. Projects may extend it with project-specific registration (e.g., updating a task catalog, adding a CODEOWNERS entry, wiring a dedicated codegen target).

## Alternatives considered

- **Code-first with generated OpenAPI/JSON Schema (e.g., decorators → schema).** Rejected as default: the schema becomes an afterthought reviewed at merge time rather than before code is written, and cross-language / cross-service consumers are second-class. Adopt locally only when the project is genuinely a single-language monolith and unlikely to grow agent consumers.
- **Pick one contract form (OpenAPI only, or protobuf only).** Rejected: OpenAPI underserves internal RPC hot paths; protobuf underserves browser/human consumers. JSON Schema serves neither as an HTTP or RPC protocol. All three carry weight.
- **GraphQL as a replacement for OpenAPI.** Deferred: GraphQL is a valid choice when query flexibility matters more than request/response simplicity. Projects that need it can add `templates/contracts/graphql/` via a follow-up ADR without disturbing this decision.
- **gRPC + protobuf for everything (including external).** Rejected: browser and partner support is materially worse; documentation tooling is thinner; costs outweigh benefits for external-facing surfaces.
- **Skip codegen; hand-write clients and servers from the contract.** Rejected: hand-written translations drift. Even when the initial overhead is small, the ongoing maintenance tax is real and compounds.

## Consequences

- **Easier.** Cross-language consumers generate types from the contract. Cross-service changes fail at contract-diff time, not runtime. Agents (human and LLM) read a schema and write code against a known target. Task-kind contracts plug directly into the agent primitives from ADR 0003 — the task catalog *is* a directory of JSON Schemas. Review shifts left: reviewers evaluate the shape before implementations exist. MCP servers (slice 3) bind to the same contracts, so an MCP tool and an HTTP endpoint that serve the same Task kind share one source of truth.
- **Harder.** Each project pays a codegen setup cost (one-time, ~a day). `make contracts` must stay fast and reproducible; if it drifts, CI alerts. Contract churn is visible — cosmetic edits that don't change semantics show up in diffs, so teams learn to land contract changes deliberately. Breaking changes are more expensive (bump major, file ADR, migrate consumers) — which is the point, but should be acknowledged.
- **Follow-ups.**
  - MCP server skeleton (slice 3) demonstrates binding to the same contracts via the Task/Result primitives.
  - Structured-logging slice ties log schemas to contract schemas via `trace_id` and `step_id`.
  - A catalog-lint CI workflow may be added later — similar in shape to ADR 0002's annotation enforcer — to verify every registered `kind` in `catalog.yaml` points to existing schemas.
  - Canonicalization ADR (how to compute stable content-addressed ids for Tasks and Results across languages) becomes relevant once two independent services need to agree on `task_id`.

## See also

- `templates/contracts/` — the three-form template tree and decision matrix.
- `templates/agent-primitives/` — the primitives these contracts reference.
- `claude-instructions/contract-discipline.md` — operating manual for day-to-day contract work.
- `.claude/skills/new-contract.md` — scaffold a new contract.
- [0003](0003-agent-primitives-as-foundation.md) — the primitives that anchor this ADR's references.
