# 0003. Typed agent primitives are the canonical handoff between rings

- **Status:** Accepted as a recommended pattern. Projects adopt this ADR by copying `templates/agent-primitives/schemas/` into their contract tree and filing a follow-up ADR naming their task catalog and any project-specific extensions.
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0002](0002-boundary-annotation-enforcement.md); `claude-instructions/agent-primitives.md`; `templates/agent-primitives/`

## Context

Agent systems compose. A human-facing agent (ring 0) invokes a worker agent (ring 1) which invokes a sub-worker (ring 2) which invokes a tool (ring 3). At each hop the caller hands the callee a unit of work and receives a result. The shape of that handoff decides whether the system can scale past two rings without collapsing into unverifiable prose.

The default — handing off free-text prompts and parsing free-text replies — fails predictably:

- **Fidelity decays at every hop.** Ring 2 re-parses ring 1's prose, loses structure, introduces interpretation error. By ring 3 the original intent is a rumor.
- **No caching.** Identical work produced in two branches of the call tree cannot be deduped because there is no stable identifier for "this exact unit of work."
- **No trust calibration.** Ring N+1 cannot tell whether ring N's output is high-confidence or a guess.
- **No auditability.** When an action is taken six rings deep, there is no record of which agent authorized what, on whose behalf.
- **No training signal.** Feedback ("this answer was wrong") cannot be attributed to the specific claim, source, or step that produced it.
- **Delegation is unbounded.** A free-text "you have access to these tools" grant cannot be narrowed when ring 1 delegates to ring 2.

Each failure above has the same root cause: **the handoff medium was designed for a single human reader, not for a composed system of agents.** Fixing it at each project independently produces five incompatible shapes; fixing it at the boilerplate gives every project the same Lego brick.

Agents will increasingly produce and consume the outputs of *other agents* — the human is sometimes at ring 0, sometimes not in the loop at all. The primitives that the boilerplate ships therefore need to be agent-first: typed, content-addressed, capability-scoped, provenance-carrying, machine-parseable by default.

## Decision

**Every project derived from this boilerplate treats five typed objects as the canonical language of agent-to-agent handoff, defined once in `templates/agent-primitives/schemas/` and copied into the project's contract tree:**

1. **`Task`** — a typed unit of work. A prompt is a *rendering* of a Task for an LLM; the Task itself is structured. `task_id` is content-addressed: `hash(kind, version, inputs, params)`. Identical Tasks produce identical ids.
2. **`Result`** — a typed output. The consumer — often another agent — acts on structured fields. Free text lives inside `output` under a schema-defined field, never as the Result itself. `result_id = hash(task_id, output, step_id)`.
3. **`Evidence`** — a claim paired with a citation. Anchors agent assertions to verifiable sources.
4. **`Provenance`** — who produced an object, with what inputs, at which step. Carried by every produced object. Enables walking the call tree, reasoning about trust, and replaying failures.
5. **`CapabilityToken`** — scoped, attenuable authority issued by ring N to ring N+1. Inner rings MAY narrow; they MUST NOT widen. Modeled on macaroons: authority is a bundle of caveats, not an identity claim.

Five non-negotiable wiring rules follow:

- **Every produced object carries `Provenance`.** No exceptions.
- **`task_id` and `result_id` are content-addressed.** A stale cache cannot mis-serve; retries converge; dedup across the call tree is automatic.
- **Capability tokens only attenuate on delegation.** A sub-agent constructs a new token with *narrower* capabilities (e.g., same action, narrower scope, shorter TTL) before handing it down.
- **Errors are structured.** `Result.error.code` is a stable, machine-readable code; agents branch on it. Humans read `message`.
- **Primitives are validated at ring boundaries.** This is the boundary validation the house rules permit; internal code trusts the resulting types.

Projects adopt this pattern by (a) copying the schemas, (b) defining a **task catalog** — a registry of `kind` values with input and output schemas — and (c) filing a follow-up ADR naming their catalog location and any project-specific extensions (extra Provenance fields, extra capability actions, extra Evidence source kinds).

## Alternatives considered

- **Free-text prompts and free-text replies (status quo for most agent systems).** Rejected: fails at every dimension above. Retained only as an *escape hatch* for the ring-0 human-facing turn, where the human's input genuinely is prose.
- **Project-specific typed shapes per service.** Rejected: produces a different handoff shape per repo, preventing composition across projects and making agent-to-agent calls across service boundaries costly. The primitives are worth defining once, in the boilerplate, even though most of the surface value is local.
- **Adopt an existing agent protocol (e.g., LangChain's Runnable, OpenAI Assistant tool-call shape, AutoGen messages).** Rejected as the default: those shapes are tightly coupled to specific runtimes and do not cleanly express content-addressed ids, attenuable capability tokens, or multi-ring provenance. They can be adapters *into* these primitives, not replacements.
- **Skip `CapabilityToken`; reuse the project's auth tokens.** Rejected: normal auth tokens are identity-bound and not attenuable. A sub-agent that inherits its parent's full auth token is as dangerous as giving it the parent's credentials. Macaroon-shaped tokens keep delegation bounded.
- **Skip `Evidence`; treat citations as a free-form field.** Rejected: the whole point of the primitives is to make agent outputs verifiable and trainable. Free-form citations defeat both; the cost of a schema is small.
- **Use protobuf instead of JSON Schema as the canonical form.** Deferred, not rejected. JSON Schema wins for readability, for integrating with MCP and OpenAPI, and for the zero-dependency path. A future ADR may add a canonical protobuf form for hot paths; until then, project-specific protobuf is a *generated* form, not the source of truth.

## Consequences

- **Easier.** A new service that consumes another service's Results reads a schema and generates types — no bespoke parsing. Agent calls deduplicate and cache by construction. Feedback signals attach to specific Evidence with specific sources, giving `iris-slm-lab` and similar learning loops a training-ready substrate. Audit walks the `parent_step_id` chain. Ring-2 and ring-3 agents compose without losing context. The same typed surface serves both human labelers (via a thin UI) and agent critics (directly) — the human-in-loop → agent-in-loop transition becomes a config change, not a rewrite.
- **Harder.** Projects must maintain a task catalog of kind → (input schema, output schema). Adding a new task kind is a schema change that follows the same discipline as an API change (semver, ADR if breaking). Producers must hash correctly — canonical serialization matters; an ADR may follow to lock the canonicalization. Capability-token verification requires per-issuer key management. Engineers used to prompt-only agent code will initially feel the overhead; the payoff appears when the call tree grows past two rings or when incidents force a post-mortem.
- **Follow-ups.**
  - Contracts slice: `templates/contracts/` will reference these primitives from OpenAPI / protobuf definitions.
  - MCP server skeleton: demonstrates accepting a CapabilityToken, emitting structured Results, and carrying Provenance.
  - Logging slice: structured logs carry `trace_id` and `step_id` from Provenance, unifying log/trace/result.
  - A canonicalization ADR may be needed once two independent services need to agree on `task_id` hashing.
  - `iris` and `iris-slm-lab` migrate their internal objects to these primitives as the first consuming projects.

## See also

- `templates/agent-primitives/` — the schemas themselves, with a short orientation README.
- `claude-instructions/agent-primitives.md` — doctrine, the ring model, when to use typed primitives vs. free text, how to add a new task kind.
- `claude-instructions/security-practices.md` — OWASP Agentic Top 10 items that the capability-token and provenance patterns mitigate (over-permission, unbounded delegation, unverifiable outputs).
- [0002](0002-boundary-annotation-enforcement.md) — the adjacent pattern for enforcing conventions at every instance of a boundary; primitive wiring rules are a candidate for enforcement once the pattern is in use.
