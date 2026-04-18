# Commit conventions

[Conventional Commits](https://www.conventionalcommits.org/) with optional scope. Consistent history makes changelogs, release notes, and `git log` scanning actually useful.

## Format

```
<type>(<scope>): <short summary>

<body: why, not what>

<footer: issue refs, co-authors, breaking changes>
```

## Types

| Type | Use |
| --- | --- |
| `feat` | New user-facing capability |
| `fix` | Bug fix |
| `docs` | Doc-only change |
| `refactor` | Internal reshape, no behavior change |
| `test` | Tests added or changed, no code change |
| `chore` | Tooling, deps, config |
| `perf` | Performance improvement |
| `style` | Formatting only (pre-commit fixups, etc.) |
| `build` | Build system or external deps |
| `ci` | CI config |
| `revert` | Reverts a previous commit |

## Scope (optional)

A short, stable identifier for the area touched: `(api)`, `(web)`, `(auth)`, `(ci)`, `(claude-instructions)`. Pick from the project's module list where possible.

## Short summary

- Imperative mood: "add", "fix", "remove" — not "added" or "adds".
- Lowercase (except proper nouns and identifiers).
- No trailing period.
- Aim for ≤ 72 characters.

## Body

- Explain **why** the change is needed, not what the diff already shows.
- Reference the problem, the trade-offs considered, and why this approach.
- One paragraph is usually enough.

## Footer

- **Issue refs:** `Fixes #123`, `Part of #456`. Use `Fixes` only when the commit actually closes the issue.
- **Co-authors:** `Co-authored-by: Name <email>` — one per line. Include the real author when pair-programming or when Claude contributed substantively.
- **Breaking changes:** `BREAKING CHANGE: <description>`. Also bump major version.

## Examples

```
feat(api): add rate-limit middleware with per-tenant bucket

Previously, a single noisy tenant could exhaust the API's request
budget for everyone. This adds a token-bucket limiter scoped by
tenant ID, with per-plan limits loaded at startup.

Fixes #412
```

```
fix(web): stop leaking auth token into query string on redirect

Callback URLs were built with the token appended for debugging;
that debug path reached production. The fix drops the token from
URL construction and adds a regression test that fails if any
code path ever reintroduces it.

Fixes #488
Co-authored-by: Priya Patel <priya@example.com>
```

```
chore(deps): bump go modules (weekly)
```

## Squash vs merge

Default to squash on PR merge. The PR title becomes the squashed commit summary — treat it with the same care as a commit summary.

Exception: if the PR contains a sequence of commits that each tell a distinct story and you want them preserved (e.g. a carefully staged refactor), merge with `--no-ff` to keep history. Rare.

## Co-authoring with Claude

When Claude contributed substantively (wrote code, made design choices, not just boilerplate), include:

```
Co-authored-by: Claude <noreply@anthropic.com>
```

Don't credit Claude for trivial edits (a one-line fix you dictated). Do credit when the reasoning or implementation was meaningfully Claude's.

## Signing

If the project requires signed commits (GPG or SSH), document the setup in `CONTRIBUTING.md`. The boilerplate doesn't require signing by default.
