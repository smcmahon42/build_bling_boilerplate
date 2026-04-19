# agent-primitives

The canonical typed objects passed between agents in this project. Copy these schemas into any service that participates in an agent call graph; do not reinvent them per service.

## What

Five JSON Schemas that define the **language of handoff between rings** of an agent system:

| Schema | Role |
| --- | --- |
| [`task.schema.json`](schemas/task.schema.json) | A typed unit of work. Replaces free-text prompts as the canonical call object. |
| [`result.schema.json`](schemas/result.schema.json) | Typed output of a Task. Consumers act on structured fields, not prose. |
| [`evidence.schema.json`](schemas/evidence.schema.json) | A claim with a citation. Turns assertions into verifiable, trainable records. |
| [`provenance.schema.json`](schemas/provenance.schema.json) | Who produced an object, with what inputs, at which step. Carried by every produced object. |
| [`capability-token.schema.json`](schemas/capability-token.schema.json) | Scoped, attenuable authority passed from outer rings to inner rings. |

## Why separate

Agents composed across rings (agent → sub-agent → sub-sub-agent) lose fidelity at every hop when the handoff medium is free text. Typed primitives preserve structure, provenance, confidence, and capability across arbitrary depth. Content-addressed ids (`task_id`, `result_id`, source ids in Evidence) make caching and dedup work automatically across the tree.

## How to use these in a new project

1. **Copy** `schemas/` into the project's contract directory (e.g., `contracts/agent-primitives/` or wherever the project's schemas live).
2. **Replace** the `$id` base URL (`https://schemas.example.com/...`) with a URL the project controls, or keep the placeholder if `$id`s are treated as opaque within the project. If you keep the placeholder, do so across all five files so `$ref`s resolve.
3. **Define task kinds** as pairs of input/output schemas registered in a catalog. Each kind references the primitives above — for example, a `summarize.thread` kind's Result has `output` matching its own schema and `evidence[]` following `evidence.schema.json`.
4. **Generate code** from the schemas using the project's codegen path (see the contracts slice — not yet scaffolded). Prefer generated types over hand-written ones; hand-written types drift.
5. **Validate at ring boundaries.** Every Task received by a producer and every Result received by a consumer is validated against the schema. This is the boundary validation the house rules allow — internal code trusts the types.

## Wiring rules (non-negotiable for projects that adopt these)

- **Every produced object carries Provenance.** No exceptions.
- **`task_id` and `result_id` are content-addressed.** A stale cache cannot mis-serve because identical inputs produce identical ids.
- **Capability tokens only attenuate.** A sub-agent MUST NOT construct a token granting more than it received.
- **Errors are structured.** `Result.error.code` is machine-readable and stable. Prose goes in `message`.
- **Free text lives inside `output`**, under a schema-defined field. It does not live as the Result.

## Related

- [`claude-instructions/agent-primitives.md`](../../claude-instructions/agent-primitives.md) — doctrine and design guidance.
- [`docs/decisions/0003-agent-primitives-as-foundation.md`](../../docs/decisions/0003-agent-primitives-as-foundation.md) — the decision record.

## Owner

Project maintainers — see `CODEOWNERS`.
