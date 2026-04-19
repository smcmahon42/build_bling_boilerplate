# json-schema

JSON Schemas for task-kind input/output shapes and generic data contracts. These are the schemas that `Task.inputs` and `Result.output` validate against.

## What lives here

- **One schema per task-kind per side.** For a kind `summarize.thread` at version `1.0.0`, ship `summarize-thread.input.schema.json` and `summarize-thread.output.schema.json`.
- **Shared data types.** If two kinds use the same value object (`Citation`, `Speaker`, `ModelConfig`), factor it into its own schema and `$ref` from both.
- **Persistent-record schemas.** Event payloads and stored records that aren't part of an HTTP contract per se.

## Starters

- [`summarize-thread.input.schema.json`](summarize-thread.input.schema.json) — example input schema for a `summarize.thread` task kind.
- [`summarize-thread.output.schema.json`](summarize-thread.output.schema.json) — example output schema for the same kind.

Copy these as the starting shape for a new task kind; they demonstrate versioning, evidence-anchored claims in `output`, and how to reference agent primitives.

## Conventions

- **File name mirrors kind.** `classify.intent` → `classify-intent.input.schema.json` / `classify-intent.output.schema.json`. Dots become hyphens; other kind segments are preserved.
- **`$id`** follows `https://<project-schema-base>/<kind-path>/<version>/<side>.schema.json` — pick a base and keep it stable; treat it as opaque if the project doesn't host schemas publicly.
- **`version`** in `$id` is the kind's semver, not the schema draft version. Bump major when a field is removed or changes meaning.
- **Additive changes** (new optional field) are minor-version bumps, same `$id` path. Clients on the older version keep working.
- **Output schemas include `evidence[]` by convention** when the producer is an LLM — even if the consumer doesn't require it today, the shape is there for the day it does.

## The task catalog

A project's catalog is the registry mapping kinds to their schemas. Minimal form:

```yaml
# contracts/catalog.yaml (or similar)
kinds:
  summarize.thread:
    version: 1.0.0
    input: json-schema/summarize-thread.input.schema.json
    output: json-schema/summarize-thread.output.schema.json
    owner: "{{OWNER}}"
```

Producers read the catalog at startup to know which kinds they serve; routers read it to dispatch Tasks. A catalog-lint CI job verifies every catalog entry points to an existing schema file.

## Related

- [`../README.md`](../README.md) — contract-form decision matrix.
- [`../../agent-primitives/`](../../agent-primitives/) — the `Task`/`Result` primitives these kinds plug into.
- [`../../../claude-instructions/agent-primitives.md`](../../../claude-instructions/agent-primitives.md) — how to add a new task kind end-to-end.
