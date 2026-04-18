# Testing practices

TDD is the starting posture for every change — see [development-workflow.md](development-workflow.md) for the loop. This file covers the *shape* of the test suite.

## The pyramid

```
        /\
       /e2e\          Playwright (any UI) or HTTP contract (services-only)
      /------\        Few, slow, high-confidence. Run on PRs + nightly.
     /integr. \       Real boundaries: DB, HTTP, queue. No mocks for the thing under test.
    /----------\      Moderate count, moderate speed. Run on PRs.
   /   unit     \     Framework-native runner. Fast, pure, isolated.
  /--------------\    Many. Run on every save (watch mode) and every CI build.
```

## Unit tests

- Use the framework's native runner (`go test`, `vitest`, `pytest`, `cargo test`, `jest`).
- **One behavior per test.** A test with three assertions about three different behaviors is three tests.
- **Test through the public API.** Don't reach into private fields to assert state.
- **No mocks of your own code.** Mock at the boundary (HTTP client, DB driver) — never mock a function you also own. If it's hard to test without mocking internal code, the seam is wrong.
- Keep them under ~10ms each. If a unit test needs more, it's probably an integration test.

## Integration tests

- Hit **real boundaries**: a real database (Docker or testcontainers), a real HTTP server, a real queue.
- Preferred over mocks when the boundary is easy to stand up locally. Mocks silently diverge from reality over time.
- Isolate state: each test gets a fresh schema/namespace, or cleans up before running.
- Run in CI on every PR. If a boundary is expensive (LLM API, paid third-party), use a recorded-cassette library or a test double *at the network layer*, not inside your code.

## End-to-end tests (frontend)

**Playwright is the standard** for any project with a frontend. It drives a real browser, captures traces/videos on failure, and supports component + page-level tests.

Directory layout:

```
<repo-root>/
  e2e/                     # or <app>/e2e/ in a monorepo
    playwright.config.ts
    tests/
      auth.spec.ts
      checkout.spec.ts
    fixtures/
      auth-state.json      # logged-in storage state, seeded in global setup
```

Required conventions:

- **Trace on retry, video on failure.** Debug loops matter.
  ```ts
  use: { trace: 'retain-on-failure', video: 'retain-on-failure' }
  ```
- **Auth via storage state**, not by filling a login form in every test. Seed it once in `globalSetup`.
- **Network mocks at the test layer**, not by rewriting the app.
- **One user journey per file.** Don't pile unrelated flows into `smoke.spec.ts`.
- **Run against a real backend** in CI (docker-compose) or a deployed preview URL. Don't stub the backend in e2e — that's an integration test's job.

## End-to-end tests (services-only)

If there's no UI, e2e = exercising the public API end-to-end against a running stack. Use the service's HTTP client (or gRPC client) against a docker-composed instance.

## Coverage

The boilerplate does not impose a specific threshold — that's per-project. Pick a floor, write it in a `COVERAGE.md` at repo root, and ratchet up over time. Treat coverage as a **floor**, not a goal: 100% coverage of trivial code is useless; 70% of the critical paths beats 90% of the total.

## Tests that don't count

- Tests that assert the code "compiles" or "doesn't throw" without checking behavior.
- Tests that mock the thing they're testing.
- Snapshot tests for structures you don't review on every change.
- Tests skipped with `skip` or `xit` that nobody is going to un-skip.

If a test falls into these buckets, delete it.

## Running locally

Each project documents its own commands (in the root `README.md` or `Makefile`). As a rule:

- Unit: `<runner> --watch` during development.
- Integration: `<runner>` without watch, requires local services up.
- e2e: `npx playwright test` after `npm run dev` (or equivalent) is running.

## Related skills

- `/test-gaps` — analyze coverage across the repo and list uncovered critical paths.
