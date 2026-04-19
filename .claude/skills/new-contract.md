---
name: new-contract
description: Scaffold a new contract (OpenAPI, JSON Schema, or protobuf). Picks the form using the decision matrix, copies the starter, registers the contract, and reminds about codegen wiring.
---

# New contract

Create a new machine-readable contract under `contracts/`. Contracts are the source of truth; code is generated from them. See ADR [`0004`](../../docs/decisions/0004-contracts-as-source-of-truth.md).

## Steps

1. **Confirm `contracts/` exists.** If not, copy `templates/contracts/` into `contracts/` and offer to stage the copy as a dedicated commit (`chore(contracts): bootstrap contracts tree`). Stop and ask before committing.

2. **Pick the form using the decision matrix.** Ask the operator which applies:
   - **OpenAPI 3.1** — external HTTP API (browsers, CLIs, partners, other teams).
   - **JSON Schema** — task-kind input/output pair, or generic data contract (event payload, stored record).
   - **Protocol Buffers (proto3)** — internal agent-to-agent RPC; hot path; streaming; cross-language codegen.
   If unclear, read `templates/contracts/README.md` back to the operator and let them pick. Do not guess.

3. **Ask for a name / slug** appropriate to the form:
   - OpenAPI: service name (`orders-service.yaml`, lowercase kebab).
   - JSON Schema: for a task kind, the dotted kind (`summarize.thread`), which becomes `summarize-thread.input.schema.json` + `summarize-thread.output.schema.json`. For a generic data contract, a single kebab-case name.
   - Proto: service name (`orders/v1/service.proto`). Ask for a version segment (default `v1`).

4. **Copy the starter** from `templates/contracts/<form>/` into `contracts/<form>/<target>`. Replace placeholders (`{{SERVICE_NAME}}`, `{{SERVICE_DESCRIPTION}}`, `{{OWNER}}`, `{{PROJECT_NAME}}`) using values from the project's `CLAUDE.md` header and `CODEOWNERS`. Ask the operator for anything missing.

5. **Validate `$ref` / import paths.** Adjust the references to the agent primitives so they resolve from the new contract's location:
   - JSON Schema / OpenAPI: `$ref` to `../agent-primitives/schemas/<primitive>.schema.json` (or whatever path the project uses).
   - Proto: if the project has agent-primitive `.proto` files in `contracts/proto/agent_primitives/`, `import` them; otherwise, keep the messages inline from the starter and flag that a future refactor could extract them.

6. **Edit the body** to reflect the actual contract the operator wants:
   - OpenAPI: the real paths and operations.
   - JSON Schema (task kind): the real input fields and output shape. Inputs are the semantic arguments; params are non-semantic knobs; output includes a named free-text field (if any) alongside structured siblings.
   - Proto: the real service RPCs and domain messages.

7. **Register the contract.**
   - OpenAPI: add it to `contracts/openapi/README.md` table if one exists; otherwise just leave the file.
   - JSON Schema task kind: add an entry to `contracts/catalog.yaml` (or the project's catalog file). If the catalog doesn't exist, scaffold it with one entry using the shape in `templates/contracts/json-schema/README.md`.
   - Proto: add the package and service to `contracts/proto/README.md` table if one exists.

8. **Wire codegen.** Add a rule to the project's `Makefile` (or equivalent) so `make contracts` picks up the new file. If `make contracts` doesn't exist yet, scaffold it with a minimal target that calls the appropriate generator for this form. Do not guess generator paths — ask the operator which toolchain (oapi-codegen, openapi-generator, buf, protoc, etc.) the project uses.

9. **Add a failing test.** Create or extend a contract-shape test that instantiates the generated type, round-trips it through JSON/proto, and validates against the schema. TDD — write the test before running codegen.

10. **Decide on an ADR.** Ask: is this contract load-bearing?
    - **Yes** (new external API surface; new cross-cutting task kind; new RPC surface): run `/new-adr` next with a slug like `adopt-<contract>`. Link from the ADR back to this contract file.
    - **No** (small, local, additive): skip the ADR.

11. **Run the quality gate locally.**
    - Lint the contract (`openapi-validate`, `buf lint`, or a JSON Schema validator — use what the project has).
    - Run `make contracts` and verify the tree diff is clean.
    - Run the failing test to confirm it fails for the *expected* reason.

12. **Stage and commit (ask first).** Include:
    - the new contract file(s)
    - the catalog / README registration edit
    - the Makefile / codegen wiring edit
    - the failing test
    - the ADR, if filed
    Commit message: `feat(contracts): add <form> contract for <name>`.

## Invariants

- **One form per contract.** A kind's input/output is JSON Schema; a service is OpenAPI *or* protobuf, not both. Choose using the matrix.
- **References, not duplication.** Reuse the agent primitives via `$ref` or import — never inline their shape.
- **Versioned paths.** OpenAPI major version in the path; proto major version in the package; JSON Schema version in `$id`.
- **Breaking changes need an ADR.** If this skill run produces a breaking change to an existing contract, stop and require `/new-adr` before proceeding.

## When NOT to run this skill

- **The change is additive to an existing contract** (new optional field, new response status) — edit the existing file directly; no scaffolding needed.
- **You're documenting an internal function.** Contracts are for cross-boundary interfaces.
- **You're writing a prototype under `prototypes/`.** Prototypes skip the contract tree until they graduate.

## Related

- `templates/contracts/` — the form-specific starters and decision matrix.
- `claude-instructions/contract-discipline.md` — day-to-day operating manual.
- `claude-instructions/agent-primitives.md` — how task kinds plug into the primitives.
- `.claude/skills/new-adr.md` — companion skill for the load-bearing-contract case.
