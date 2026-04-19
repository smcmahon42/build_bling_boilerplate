# openapi

OpenAPI 3.1 contracts for HTTP APIs consumed by humans, browsers, external services, or other teams.

## When to use OpenAPI (and when not to)

**Use OpenAPI when:**

- The caller is a browser, a CLI run by a human, an external partner, or a service in a different org.
- You want SwaggerUI / Redoc / auto-generated clients in many languages.
- The interaction is request/response over HTTP — not streaming, not bidirectional, not sub-millisecond-sensitive.

**Use protobuf instead when:**

- The caller is another agent or service inside this system at a hot path.
- You want strongly typed streaming.
- Binary efficiency matters.

**Use plain JSON Schema instead when:**

- You're defining a task-kind's `inputs` or `output` shape (these are referenced from OpenAPI or protobuf; they do not need an HTTP wrapper of their own).
- You're defining a persisted data contract (event payload, stored record).

## What's here

- [`service.yaml`](service.yaml) — a minimal OpenAPI 3.1 starter showing a service that accepts a `Task` (POST /tasks) and returns a `Result`. Demonstrates: referencing agent primitives via `$ref`, CapabilityToken as a security scheme, structured errors, and server variables for versioning.

## Conventions

- **Version in the path.** `/v1/tasks`, not a header. Major-version bumps get a new path; minor/patch ship in-place.
- **Every request that triggers work takes a `Task`.** Every response that produces work returns a `Result`. HTTP status codes are transport-level (`200`, `400`, `5xx`); semantic status lives in `Result.status`.
- **CapabilityToken is bearer-token shaped.** The token's JSON serialization is `Base64URL`-encoded; the server deserializes and validates (signature, expiry, scope).
- **Idempotency via header.** Clients pass `Idempotency-Key: <Task.idempotency_key>`. Servers dedupe.
- **Errors are `Result` with `status: error`.** 4xx/5xx bodies use the same schema — one parser.

## Codegen

- **Server stubs:** use a framework appropriate to the project's stack (oapi-codegen for Go, openapi-generator for many languages, FastAPI for Python).
- **Client SDKs:** generate per-language from the same `service.yaml`.
- **Wire it into `make contracts`** so regeneration is one command and CI can diff the working tree.

## Related

- [`../README.md`](../README.md) — decision matrix across contract forms.
- [`../../agent-primitives/schemas/`](../../agent-primitives/schemas/) — the primitive schemas `service.yaml` references.
