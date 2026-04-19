# Agent primitives

Typed objects passed between agents in this project. Use these instead of free-text prompts whenever one agent calls another. See ADR [`0003`](../docs/decisions/0003-agent-primitives-as-foundation.md) for the decision; this file is the operating manual.

## The five primitives

| Primitive | Use it when | Schema |
| --- | --- | --- |
| `Task` | Any agent-to-agent call. The handoff object. | [`task.schema.json`](../templates/agent-primitives/schemas/task.schema.json) |
| `Result` | Any agent-to-agent response. | [`result.schema.json`](../templates/agent-primitives/schemas/result.schema.json) |
| `Evidence` | Inside a Result, to back a claim with a citation. | [`evidence.schema.json`](../templates/agent-primitives/schemas/evidence.schema.json) |
| `Provenance` | Carried by every produced object. Never omit. | [`provenance.schema.json`](../templates/agent-primitives/schemas/provenance.schema.json) |
| `CapabilityToken` | Delegating scoped authority from ring N to ring N+1. | [`capability-token.schema.json`](../templates/agent-primitives/schemas/capability-token.schema.json) |

## The ring model

Think of an agent call graph as concentric rings:

- **Ring 0** — the human-facing turn. Free text is fine *here*.
- **Ring 1+** — every agent-to-agent hop. Typed primitives only.

When you write code that crosses a ring boundary, you are serializing a `Task` and deserializing a `Result`. The prompt you send to the LLM inside ring 1 is a *rendering* of the Task; the Task is the canonical object.

Walking the tree works because every Result carries `provenance.parent_step_id` back to the Task that caused it. Walking the fan-out works because every Task has a content-addressed `task_id` — identical work at different points in the tree collapses.

## Non-negotiable wiring rules

1. **Every produced object carries `Provenance`.** If you're tempted to skip it "because this is internal," you are building ring-N+1's audit gap.
2. **`task_id` is `hash(kind, version, inputs, params)`.** Producers compute it; consumers verify. Never accept a `task_id` a caller supplied without recomputing.
3. **`result_id` is `hash(task_id, output, provenance.step_id)`.** Same rule: recompute on receipt.
4. **CapabilityTokens only attenuate.** When ring N delegates to ring N+1, construct a *narrower* token (same or fewer actions, narrower scope, same-or-shorter TTL). Set `attenuated_from` to the parent token_id.
5. **Errors are structured.** `Result.error.code` is machine-readable, stable, dotted (e.g. `budget.exceeded`, `input.invalid`, `capability.denied`). Prose goes in `message`.
6. **Free text lives inside `output`.** Under a schema-defined field. Not as the Result itself.
7. **Validate at ring boundaries.** Every Task received by a producer and every Result received by a consumer is schema-validated. Internal code trusts the types (see house rule 6 in root `CLAUDE.md`).

## Adding a new task kind

A *kind* is a registered pair of input and output schemas. Follow this when adding one:

1. **Name it.** Dotted, lowercase, stable: `summarize.thread`, `classify.intent`, `extract.entities.v1`. The prefix groups related kinds.
2. **Write its input and output schemas** under the project's `contracts/` tree (or wherever the project's schemas live — the contracts slice defines this). The schemas reference the primitives (`$ref` to `provenance.schema.json` etc.).
3. **Register it in the task catalog** — one file per project listing kind → (input schema path, output schema path, version, owner). Callers and producers look up kinds here.
4. **Add a failing test** that constructs a Task of this kind, hands it to the producer, and asserts the Result shape. TDD first (house rule).
5. **File an ADR if the kind introduces a new capability action, a new Evidence source kind, or a breaking change to an existing kind.** Adding a new kind alongside existing ones doesn't require one.
6. **Document the kind** in the owning module's README — one paragraph, linking to the schemas.

## When free text is the right answer

Be honest about escape hatches:

- **Ring 0 (human turn).** Humans write prose. Ring 0 converts prose to a Task.
- **Rendering for an LLM.** Inside a producer, you turn the Task into a prompt. That prompt is free text by nature — but it is derived from the Task, not a substitute for it.
- **The `output.summary` or `output.explanation` field of a Result.** Natural-language outputs live here, under a named field, alongside structured siblings (`output.claims[]`, `output.confidence`).

If a handoff *between* rings is free text, you have not adopted the pattern. Stop and reshape.

## Capability tokens in practice

- **Issue narrowly.** When ring 0 issues a token to ring 1, include only the actions ring 1 needs — not every action the human is authorized for.
- **Attenuate before delegating.** Ring 1 calls ring 2 with a token whose `capabilities` are a strict subset, `attenuated_from` set, `expires_at` at-or-before the parent's.
- **Verify on receipt.** Check the signature against the issuer's known key, check `expires_at`, check `subject == self.agent_id`, check the action against what's being attempted, reject unseen nonces after TTL or track them.
- **Reject empty capability arrays.** A token granting nothing is a misconfiguration, not a no-op.

Token management (key rotation, issuer trust roots, revocation) is a per-project concern — file a project-specific ADR when you adopt this.

## Evidence discipline

Every non-trivial claim a producer makes in `Result.output` should appear in `evidence[]` with a pointer to its source:

- **`kind`** — `document`, `url`, `tool_output`, `memory`, `human_input`. Consumers may weight these differently; e.g., `memory` citations are less trustworthy than `tool_output` for factual claims.
- **`id`** — prefer a content-addressed id (`sha256:...`). A locator URL is acceptable when content addressing isn't available; add a `range` for spans.
- **`confidence`** — the producer's confidence that the source *supports* the claim (not the confidence in the claim itself; that lives at the Result level).

Evidence is not decoration. It is the training substrate: when a downstream system labels an output as wrong, the label attaches to the specific Evidence record, which carries the specific source — and that's what feedback loops need.

## How this integrates with the rest of the repo

- **Validation at boundaries** (house rule 6): primitive validation *is* the boundary check. Don't also hand-write input validators for the same fields.
- **Structured logging** (slice 4, forthcoming): logs carry `trace_id` and `step_id` from Provenance. A log line and a Result can be correlated without bespoke plumbing.
- **MCP servers**: MCP tool calls arrive as Tasks, return as Results. See [`templates/mcp-server/go/`](../templates/mcp-server/go/) for the reference Go skeleton.
- **Contracts**: [`templates/contracts/`](../templates/contracts/) references these schemas from OpenAPI and protobuf definitions. See [`contract-discipline.md`](contract-discipline.md).
- **Security** (OWASP Agentic Top 10): content-addressed ids + capability attenuation + Provenance materially reduce over-permission, unbounded delegation, and unverifiable-output risks. See `security-practices.md`.

## See also

- [ADR 0003](../docs/decisions/0003-agent-primitives-as-foundation.md) — the decision record.
- [`templates/agent-primitives/README.md`](../templates/agent-primitives/README.md) — short orientation for the schemas themselves.
- [`development-workflow.md`](development-workflow.md) — TDD loop a new task kind follows.
- [`security-practices.md`](security-practices.md) — OWASP Agentic Top 10 context.
