---
name: test-gaps
description: Analyze test coverage across the repo and identify uncovered critical paths. Runs framework-native coverage (go test -cover, vitest --coverage, pytest --cov, etc.), ranks gaps by impact, and recommends the next ten tests to write.
---

# Test gaps

Produce a coverage snapshot plus a prioritized list of "tests worth writing next," weighted by impact rather than raw coverage percentage.

## Stack detection

| File | Stack | Coverage command |
| --- | --- | --- |
| `go.mod` | Go | `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out` |
| `package.json` with `vitest` | Vitest (JS/TS) | `npx vitest run --coverage` |
| `package.json` with `jest` | Jest (JS/TS) | `npx jest --coverage` |
| `pyproject.toml` | Python | `pytest --cov` |
| `Cargo.toml` | Rust | `cargo llvm-cov` (or `tarpaulin`) |
| `playwright.config.*` | Playwright | `npx playwright test` (for pass/fail; coverage via `--coverage` if configured) |

## Output format

### Coverage summary table

| Module | Runner | Line % | Branch % | Below floor? |
| --- | --- | --- | --- | --- |

Floor comes from the project's `COVERAGE.md` if it exists; otherwise note "no floor set."

### Top ten uncovered paths

Rank by impact:

1. **Critical-path uncovered** — auth, payments, data mutations, external calls.
2. **High-blast-radius uncovered** — utilities called from many places.
3. **Recently-changed uncovered** — code changed in the last 30 days without tests.
4. **Error paths uncovered** — the happy path is tested but the failure branch isn't.
5. **Integration seams uncovered** — boundaries between modules/services.

Produce a table:

| Rank | File:function | Why it matters | Suggested test |
| --- | --- | --- | --- |

### Quick wins

Tests that would be cheap to add (under ~20 lines) and close a gap on an important path. List up to five.

### Tests to consider deleting

Coverage noise: tests that assert tautologies, snapshot tests on structures nobody reviews, `skip`/`xit` leftovers. List them — deletion is the operator's call.

## Invariants

- **Coverage percentage is an input, not the goal.** A 95% covered module where the critical path is the uncovered 5% is worse than an 80% covered module with full critical-path coverage.
- Don't write the tests yourself in this skill — recommend them. Writing follows in the normal TDD loop.
- If no test runner is detected, report that and recommend setup instead of hallucinating results.

## Related

- `claude-instructions/testing-practices.md` — pyramid, framework conventions, what tests to write.
- `claude-instructions/development-workflow.md` — TDD loop.
