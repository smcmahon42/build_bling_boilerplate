# Docs

Project documentation lives under this directory. The structure is intentional — each bucket has a clear purpose, so new docs land in the right place and old docs are easy to find.

## Buckets

| Bucket | Purpose |
| --- | --- |
| [decisions/](decisions/) | Architectural Decision Records (ADRs). Immutable once accepted; superseded by later ADRs. |
| architecture/ | (create on first use) Narrative architecture docs: system overview, module descriptions, data flow. |
| runbooks/ | (create on first use) How to operate the system — deployment, incident response, recurring ops tasks. |
| security/ | (create on first use) Threat models, security posture, audit logs. |
| qa/ | (create on first use) Test strategy, test plans, fixtures. |
| guides/ | (create on first use) Long-form internal docs and how-tos. |
| api/ | (create on first use) Narrative supplementing generated API references. |
| roadmap/ | (create on first use) Product/engineering roadmap and historical phases. |
| archive/ | (create on first use) Retired docs kept for history. Move rather than delete. |

Directories marked "create on first use" don't exist yet — add them when you have something to put there. Don't scaffold empty directories.

## Rules

- **Every new doc goes in a bucket.** If it doesn't fit, propose a new bucket (ADR-worthy decision).
- **Each bucket has a `README.md`** that indexes its contents.
- **Every change ships docs.** See [claude-instructions/documentation-discipline.md](../claude-instructions/documentation-discipline.md).
- **Archive, don't delete.** When a doc is retired, move to `archive/<year>/` with a note at the top.

## ADR quick links

- [ADR index](decisions/README.md)
- [ADR template](decisions/TEMPLATE.md)
- [0001 — Record architecture decisions](decisions/0001-record-architecture-decisions.md)

## See also

- `claude-instructions/documentation-discipline.md` — doc routing rules.
- `claude-instructions/component-explainability.md` — the module-README rule.
