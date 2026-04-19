# cli-go

Minimal Go CLI that exercises the OpenAPI contract as a thin client. It builds a `Task`, POSTs it to the service defined in [`templates/contracts/openapi/service.yaml`](../../contracts/openapi/service.yaml), and renders the `Result`.

## What this demonstrates

1. **Input ergonomics** — CLI flags map to `Task.inputs` / `Task.params` fields. The flag layer is the *only* thing that exists because humans are the consumer.
2. **Capability-token handling** — the token is read from an env var (`CAPABILITY_TOKEN`), base64-decoded, and sent as a bearer header. Never captured in shell history.
3. **Idempotency** — the CLI sets `Idempotency-Key` so retries dedupe.
4. **Error rendering** — a `Result` with `status: error` maps to a non-zero exit code and a structured stderr message carrying `error.code`. Scripts grep on the code; humans read the message.
5. **No domain logic.** The CLI contains argument parsing, HTTP plumbing, and rendering — nothing else. The server is authoritative for everything else.

## Layout

```
cmd/cli/main.go   # the whole CLI in one file (~150 LOC)
```

Keeping the CLI in one file is intentional: it makes the thin-skin rule visible. If this file starts growing past ~200 lines, the CLI is probably absorbing logic that belongs in the contract or the service.

## Adoption path

1. Copy `templates/human-surface/cli-go/` into the project (e.g., as `cmd/<project>-cli/`).
2. Replace the placeholder module path `github.com/example/project` in `go.mod` and the single import in `main.go` with the project's real module path.
3. Replace the hardcoded server URL (`http://localhost:8080/v1`) with the project's base URL (flag or env).
4. Adjust the flag set to match the task kinds the CLI exposes. One subcommand per kind, if you move to multiple.
5. If the project generates Go types from OpenAPI / JSON Schema, swap the hand-rolled `Task`/`Result` struct tags for the generated types — do not duplicate.

## Related

- [`../../contracts/openapi/service.yaml`](../../contracts/openapi/service.yaml) — the contract this CLI binds to.
- [`../../agent-primitives/`](../../agent-primitives/) — the Task/Result shape.
- [ADR 0006](../../../docs/decisions/0006-human-surfaces-are-views-over-contracts.md) — why surfaces are thin.
