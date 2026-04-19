# proto

Protocol Buffers (proto3) contracts for internal agent-to-agent RPC. Use this form when the caller is another service inside this system and you want generated types in multiple languages, streaming, or efficient binary wire format.

## When to use protobuf (and when not to)

**Use protobuf when:**

- The caller is an internal service or agent — not an external consumer.
- Calls are frequent enough that binary encoding and schema compactness matter.
- You want streaming or bidirectional RPC (gRPC).
- Multiple languages need matching types (Go server, Python client, etc.).

**Use OpenAPI instead when:**

- External consumers (humans, browsers, partners) need documentation and clients.
- You want HTTP semantics visible to intermediaries (caches, proxies, gateways).

**Use JSON Schema instead when:**

- You're describing task-kind input/output shapes (referenced from both OpenAPI and proto).

## What's here

- [`service.proto`](service.proto) — minimal proto3 starter showing a service that accepts Task messages and returns Result messages. Messages mirror the JSON Schema shape of the agent primitives by convention.

## Conventions

- **Package names are reverse-DNS.** `package com.example.service.v1;`. Versioning happens in the package name — `v1` → `v2` is a new package, not a flag.
- **No renumbering.** Field numbers are permanent. To remove a field, mark it `reserved`. To repurpose meaning, add a new field.
- **Mirror the JSON Schema shape.** `Task`, `Result`, `Provenance`, `Evidence`, `CapabilityToken` messages use field names matching the JSON Schema properties. This keeps the two forms in lockstep even though each is hand-maintained (until `make contracts` generates the proto from JSON Schema, or vice versa).
- **Error codes are strings.** `Result.error.code` is a string matching the dotted convention, not an enum — new codes are additive and don't break wire compatibility.
- **Streaming outputs use `stream Result`.** Long-running productions stream partial Results with `status = PARTIAL` and a final one with `status = OK` or `ERROR`.

## Codegen

- **Go:** `buf generate` with `protoc-gen-go` + `protoc-gen-go-grpc`, or the equivalent Bazel/Make path.
- **Other languages:** buf or protoc with the respective plugin.
- **Wire into `make contracts`** alongside OpenAPI codegen.

## Related

- [`../README.md`](../README.md) — contract-form decision matrix.
- [`../../agent-primitives/schemas/`](../../agent-primitives/schemas/) — the canonical (JSON Schema) shape of the primitives this proto mirrors.
