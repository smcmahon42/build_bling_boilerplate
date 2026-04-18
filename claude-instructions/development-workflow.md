# Development workflow

The loop below is the default for every change — bug fix, feature, refactor, or doc-only change. Deviating requires an explicit reason in the PR description.

## Task scale matrix

Classify the task **before** you start editing. The scale determines which docs are required, how much planning is required up front, and how many sub-tasks to break the work into.

| Scale | Size signals | Planning | Docs required | Example |
| --- | --- | --- | --- | --- |
| **Small** | 1–2 files changed; fits in one short session; no new boundary | Skip — proceed directly | Update any affected README if contracts change; no ADR | Typo fix, log-level tweak, one-line config change, single-function bug fix |
| **Medium** | 3–5 files; new function(s) or endpoint(s); crosses a single module | Share a short plan (intent, files, tests) before editing | Module README updated; commit-level rationale | Add an endpoint, new service layer, new test suite, refactor within a module |
| **Large** | 6+ files; new module/service; cross-module changes; new tech | Share a design sketch (see `docs/templates/design-doc.md`); prefer an ADR first | Design Doc + ADR + module README(s) + runbook if operational | New service, new data store, auth refactor, multi-module migration |

**Use the biggest scale that applies** — a task is Large if *any* dimension (file count, new tech, cross-module) qualifies. Do not downshift to avoid writing docs.

**Signals you're mid-task at the wrong scale:** the diff keeps growing past the tier you picked; you're touching modules you didn't plan to; you're discovering new decisions. Pause, re-scope, and promote the task.

## The TDD loop

1. **Write a failing test.** If you can't articulate the behavior as a test, the change isn't ready to write yet.
2. **Run the test.** Confirm it fails for the *right reason* (not a syntax error or missing import).
3. **Write the minimum code to pass.** No scope creep. No anticipatory abstraction.
4. **Run the test again.** Confirm it passes.
5. **Refactor.** With green tests as the safety net.
6. **Repeat.**

Non-negotiable: bug fixes add a **regression test** that reproduces the bug before the fix lands. The test prevents the class of bug from recurring. No regression test = the PR is not done.

See [testing-practices.md](testing-practices.md) for pyramid, framework choices, and Playwright e2e conventions.

## Branching

| Prefix | Use |
| --- | --- |
| `feat/<issue>-<slug>` | New capability tied to an issue |
| `fix/<issue>-<slug>` | Bug fix tied to an issue |
| `chore/<slug>` | Infra, tooling, dependency bumps |
| `docs/<slug>` | Doc-only changes |
| `refactor/<slug>` | Internal reshape, no behavior change |
| `dependabot/*` | Reserved for Dependabot |

One change per branch. Delete the branch after merge.

## PR flow

1. Create the branch. Make the change. Tests first.
2. Run the local quality gate (pre-commit hooks + test runner + any language-specific linters).
3. Open the PR using `.github/PULL_REQUEST_TEMPLATE.md`. Fill in:
   - **Summary** — 1–3 bullets of what and why.
   - **Related issues** — `Fixes #N` or `Part of #N`.
   - **Test plan** — exactly how you verified it; what a reviewer should re-run.
   - **Security considerations** — any OWASP category touched? Call it out.
4. Address review. One PR should tell one story — split if it sprawls.
5. Merge (squash preferred). Delete the branch. Close the issue.

See [commit-conventions.md](commit-conventions.md) for commit message format.

## Local-CI-first (optional)

If you find yourself pushing to GitHub just to run CI, consider wiring a `Makefile` + [`act`](https://github.com/nektos/act) to run the same workflows locally. Iterate on test failures in seconds instead of minutes, and reserve GitHub Actions minutes for integration on shared branches.

Shape of this pattern:
- `make test` — run tests across modules.
- `make lint` — run language-specific linters.
- `make ci-local` — run `.github/workflows/*.yml` via `act` end-to-end.
- PR template asks contributors to paste `make ci-local` output as an attestation.

This is a pattern, not a requirement. Adopt it when CI minutes or latency become a real constraint.

## Review gates (what blocks merge)

- All tests pass (unit + integration + e2e where applicable).
- Pre-commit hooks pass (secrets, formatting, yaml/json validity).
- Security workflows pass (`security-scan`, `codeql`, `license-check`).
- New ADR merged first if the change implements an architectural decision.
- Docs updated in the matching bucket ([documentation-discipline.md](documentation-discipline.md)).
- At least one human reviewer approved (even if a Claude review also ran).

## Escape hatches

- **Hotfix to production:** branch from the production tag, minimal change + regression test, ship, then backport to main. Still write the test.
- **Doc-only PR:** skip the TDD loop, but still fill the PR template.
- **Generated code (SDKs, migrations from schema, etc.):** the generator's tests cover the output; generated files don't need their own tests.
