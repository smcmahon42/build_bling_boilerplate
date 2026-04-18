# 0002. Enforce boundary conventions via in-code annotations + allowlist

- **Status:** Accepted as a recommended pattern (not yet applied to a specific boundary — projects adopt this ADR by filing a follow-up ADR naming the boundary, annotation, and enforcer).
- **Date:** <!-- filled on merge -->
- **Deciders:** project maintainers
- **Related:** [0001](0001-record-architecture-decisions.md); `claude-instructions/component-explainability.md`; `claude-instructions/security-practices.md`

## Context

Any non-trivial project accumulates conventions that must hold at *every* instance of a boundary: every HTTP handler needs an auth check, every agent tool needs a permission scope, every data access needs a tenancy filter, every logged event uses the structured logger. Conventions documented in prose alone rot — new code lands without the convention, old code is refactored and loses it, no one notices until an incident.

The common workarounds each fail in a recognizable way:

- **Code review alone** — relies on humans remembering a checklist that grows with the project. Silently fails when reviewers rotate.
- **Runtime checks** — move the problem to production. By the time it fires, the bad code has shipped.
- **Whole-codebase linter rules** — heavy to author and noisy; often degrade to ignored warnings.

What does work, consistently, is a **small, purpose-built enforcer** paired with an **in-code annotation** and an **explicit opt-out allowlist**. Each boundary that should satisfy the convention carries an annotation (`// foo:bar=value` or equivalent). A CI check walks the codebase, finds every boundary, and fails if any boundary lacks the annotation *and* isn't on the allowlist. The allowlist is plain text, diffable, and reviewable — removing the social cost of documenting why this particular instance is exempt.

Worked example: "every public HTTP handler must register an agent tool or explicitly opt out." Each handler carries `//agent:tool=<name>` or `//agent:skip`. A small program under `tools/agent-check/` walks the handler registration, verifies each is annotated or listed in `tools/agent-check/allow.txt`, and fails CI otherwise. The same shape works for tenancy filters, logging gates, permission scopes, or any rule that must hold at every instance of a boundary.

## Decision

**When a convention must hold at every instance of a boundary, enforce it with:**

1. **An annotation** (comment, attribute, decorator, or equivalent) that carries the convention's payload. Format is stable and greppable. Convention names live in a well-known prefix (e.g., `//auth:…`, `//mcp:…`, `//tenant:…`).
2. **A small enforcer** — a script or program under `tools/<convention>-check/` that:
   - Discovers every instance of the boundary (AST walk, regex, or reflection appropriate to the language).
   - For each instance, verifies the annotation is present *or* the instance is on the allowlist.
   - Fails CI if any instance is neither annotated nor allowlisted.
3. **An allowlist** — a plain text file (`enforced.txt` or `<convention>-allow.txt`) of paths or symbols that are deliberately exempt. Each entry has a comment explaining *why* it's exempt. The allowlist is under `CODEOWNERS` so additions are reviewed.
4. **A CI workflow** that runs the enforcer on every PR and blocks merge on failure. See `.github/workflows/` for the pattern.
5. **A matching rule in `claude-instructions/`** so Claude knows the convention and writes annotated code by default.

The convention is *not* introduced without all four pieces in the same PR. Partial adoption is worse than none — it trains contributors to think the rule is optional.

## Alternatives considered

- **Whole-codebase linter rules (golangci-lint custom linter, ESLint plugin, etc.).** Rejected as the default: the tooling overhead is high, the error messages are poor, and adoption drifts across language ecosystems. Use them when a linter ecosystem already encodes the convention (e.g., gosec for security); fall back to a bespoke enforcer otherwise.
- **Runtime middleware / framework hooks.** Rejected: moves discovery to production. Keep runtime hooks as defense-in-depth alongside the enforcer, not as a replacement.
- **Pure documentation + code review.** Rejected: does not survive reviewer turnover or AI-assisted code generation, which often skips conventions that aren't machine-enforced.
- **A `//nolint` style suppression inline** (instead of an allowlist file). Rejected as the default: scatters exemptions through the tree where they're invisible. The allowlist file makes exemptions auditable as a set.

## Consequences

- **Easier.** New boundaries get the convention by default (Claude reads `claude-instructions/`, human reviewers trust the CI gate). Exemptions are visible. Conventions can be added incrementally without rewriting the tree.
- **Harder.** Each convention costs a small tool to build (~a day for a simple AST walker). Projects adopting this pattern should reserve capacity for enforcer authorship. Conventions accumulate — review them annually and retire ones that no longer carry weight.
- **Follow-ups.** Each specific convention gets its own ADR naming: the boundary (e.g., "every public HTTP handler"), the annotation format, the enforcer location, the allowlist path, and the initial exemption rationale. File those ADRs as `0NNN-enforce-<convention>.md`.

## See also

- `claude-instructions/component-explainability.md` — adjacent pattern for the README-at-every-boundary rule.
- `claude-instructions/security-practices.md` — uses this pattern implicitly for security gates (e.g., the `/security-review` skill's CI coverage audit).
- `.github/workflows/adr-lint.yml` — a concrete tiny example of the enforcer shape (annotation = filename format; allowlist = `README.md` / `TEMPLATE.md`; CI fails on missing format).
