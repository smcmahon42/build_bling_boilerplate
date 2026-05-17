# Contract discipline

Contracts under `contracts/` (copied from `templates/contracts/`) are **the source of truth** for this project's interfaces. Code is generated from them. Hand-written types that duplicate generated ones are drift in progress. See ADR [`0004`](../docs/decisions/0004-contracts-as-source-of-truth.md) for the decision; this file is the operating manual.

## Picking a form

| Use case | Form | Subdir |
| --- | --- | --- |
| External HTTP API (browsers, CLIs, partners, other teams) | OpenAPI 3.1 | `contracts/openapi/` |
| Task-kind input/output schemas (referenced from `Task.inputs` / `Result.output`) | JSON Schema | `contracts/json-schema/` |
| Internal agent-to-agent RPC (hot path, streaming, cross-language) | Protocol Buffers | `contracts/proto/` |

One contract, one form. Don't hand-translate a proto into an OpenAPI; reference or regenerate.

## Non-negotiable rules

1. **Contracts reference [agent primitives](agent-primitives.md).** `Task`, `Result`, `Evidence`, `Provenance`, `CapabilityToken` live in one place (`contracts/agent-primitives/` by convention in a consuming project). Never redefine them.
2. **Contracts are versioned.** Semver. OpenAPI: in the path (`/v1/tasks`). Proto: in the package (`example.v1`). JSON Schema: in the `$id`.
3. **Breaking changes bump major + require an ADR.** Additive changes (new optional field, new response status) are minor. The ADR names the migration plan for consumers.
4. **Codegen is one command.** Usually `make contracts`. It regenerates server stubs, client SDKs, and shared types.
5. **CI diffs the working tree after regeneration.** A workflow runs `make contracts` and fails the PR on a non-empty diff.
6. **Contract changes land first.** New endpoint / RPC / task kind = contract edit first (with a failing test expecting the new shape), code follows.

## Adding a new contract

Run `/new-contract` (`.claude/skills/new-contract.md`). It picks the form, scaffolds from the template, and registers the contract.

Manual path:

1. **Pick the form** using the table above.
2. **Copy the matching starter** from `templates/contracts/<form>/` into `contracts/<form>/`.
3. **Rename and edit** — replace placeholders (`{{SERVICE_NAME}}`, `{{OWNER}}`, `{{PROJECT_NAME}}`), define your operations/messages/schemas.
4. **Reference the primitives** via `$ref` (JSON Schema / OpenAPI) or by matching message shape (proto).
5. **Wire codegen.** Add a rule to `Makefile` (or equivalent) so `make contracts` picks up the new file. Include the output path in `.gitignore` if generated types are committed separately, or commit them to the tree if the project keeps generated code in-repo.
6. **Add a failing test** that instantiates the generated type and verifies the shape round-trips. TDD (house rule 1).
7. **Register the contract** in the project's catalog if it's a task kind (see [agent-primitives.md](agent-primitives.md#adding-a-new-task-kind)).
8. **File an ADR** if the contract is load-bearing — new external API surface, new cross-cutting task kind, new RPC surface.
9. **Document** — a short README in the contract subdirectory explaining what it is and who consumes it.

## Versioning playbook

| Change | Rule |
| --- | --- |
| Add an optional field | Minor bump. Same path / package / `$id` version. |
| Add a new endpoint, RPC, or task kind | Minor bump. |
| Add a new response status or error code | Minor. Errors are strings, not enums — additive by design. |
| Remove or rename a field | **Major.** Do not remove from the old version; ship a new major with the field absent. |
| Change a field's type or meaning | **Major.** |
| Change required / optional status | **Major** if going from optional to required; minor the other way. |
| Change URL path, proto package, or `$id` version | Only on major bump. |

Sunsetting an old major: keep serving it alongside the new one for a deprecation window (named in the ADR), then remove in a tracked PR.

## Task-kind contracts specifically

Task kinds are JSON Schema pairs (input + output) registered in the project's task catalog (e.g., `contracts/catalog.yaml`). The catalog is itself contract-adjacent — treat it like a contract for the dispatch layer.

- **One pair of schemas per kind per version.** `summarize.thread` at v1 has `summarize-thread.input.schema.json` and `summarize-thread.output.schema.json`.
- **`Result.output` includes evidence-anchored claims by convention** when the producer is an LLM — even if the first consumer doesn't require `Evidence`, the shape is there for the day another ring does.
- **`Result.error.code` values are part of the contract.** A new error code is an additive change; renaming one is a major bump.

## Escape hatches

- **Experimental endpoints / kinds.** Prefix with `x-` (OpenAPI `/v1/x-experimental/...`, kind `x.summarize.thread_experimental`) and mark them explicitly unstable in the schema `description`. Not covered by the breaking-change rules until they graduate.
- **Prototype / throwaway code.** Code in `prototypes/` or `examples/` may skip codegen. If it matures to a real service, it starts a contract then.
- **Third-party contracts.** When consuming someone else's OpenAPI / proto, copy their artifact into `contracts/external/<source>/` and regenerate clients. Never edit in place.

## What NOT to treat as a contract

- **Internal function signatures.** The type system handles this.
- **Private package-to-package calls within one service.** Language types suffice.
- **Ad-hoc scripts / one-offs.** If it's not reused across a boundary, it's not a contract.

## Related

- [ADR 0004](../docs/decisions/0004-contracts-as-source-of-truth.md) — the decision record.
- [`agent-primitives.md`](agent-primitives.md) — the primitives contracts reference.
- [`templates/contracts/README.md`](../templates/contracts/README.md) — decision matrix and form details.
- [`documentation-discipline.md`](documentation-discipline.md) — how contracts fit into the doc-routing table.
- [`.claude/skills/new-contract.md`](../.claude/skills/new-contract.md) — scaffolding skill.
