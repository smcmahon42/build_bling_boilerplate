# 0001. Record architecture decisions

- **Status:** Accepted
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [TEMPLATE.md](TEMPLATE.md)

## Context

Projects accumulate decisions that are hard to recover later. "Why do we use Postgres here?" "Why is this service separated from that one?" "Why does the agent confirm before writes?" Answers often live in the heads of whoever made the call, in Slack threads, or in outdated design docs. New contributors re-litigate settled debates. Old decisions are forgotten and then contradicted.

Architectural Decision Records (ADRs) — short, numbered, append-only markdown files — are a well-trodden solution. Popularized by Michael Nygard; adopted widely. They capture *context*, *decision*, *alternatives*, and *consequences*, and stay in the repo alongside the code they shape.

## Decision

This project records non-trivial architectural decisions as ADRs in `docs/decisions/`. Each ADR follows the template in [TEMPLATE.md](TEMPLATE.md). ADRs are numbered sequentially (`NNNN-slug.md`), append-only. A decision that overturns an earlier one creates a new ADR that marks the old one "Superseded by NNNN."

An ADR is warranted when the decision:

- Introduces a new boundary (module, service, layer, bounded context).
- Commits to a technology, vendor, or pattern the project hasn't used before.
- Changes an auth boundary, data-flow boundary, or observability boundary.
- Has a reasonable alternative that was considered and rejected.

Decisions that don't warrant an ADR: coding style, naming conventions, minor refactors, file locations, linter rules. These live in `claude-instructions/`, `CONTRIBUTING.md`, or the lint configs.

## Alternatives considered

- **No formal record; rely on commit messages and Slack.** Rejected: doesn't survive team turnover; search is poor; no stable structure.
- **Design docs in a wiki.** Rejected: lives outside version control; drifts from the code; no PR review; hard to link to specific commits.
- **Long-form RFCs only.** Rejected: too heavyweight for routine decisions; the ceremony discourages recording lower-stakes choices that still matter.
- **ADR tool (e.g., `adr-tools` CLI).** Rejected as a mandatory dependency but compatible — contributors can use it if they prefer. The `/new-adr` skill in this repo does the same job without extra tooling.

## Consequences

- **Easier:** onboarding (new contributors read the ADR index); defending decisions in code review (link to the ADR); retiring old decisions (write a superseding ADR); auditing (each decision is dated and reviewed).
- **Harder:** writing ADRs is extra work at decision time. Mitigated by the template and the `/new-adr` skill.
- **Follow-ups:**
  - `docs/decisions/README.md` maintains the index.
  - Reviewers flag PRs that implement architectural changes without a matching ADR.
  - `claude-instructions/documentation-discipline.md` codifies when an ADR is required.

## See also

- [TEMPLATE.md](TEMPLATE.md) — ADR template.
- [Michael Nygard — Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) — original proposal.
- [ADR GitHub org](https://adr.github.io/) — reference patterns.
- `../../claude-instructions/documentation-discipline.md` — the broader doc discipline rule.
