<!--
Thanks for the PR. Keep the description tight — the code tells the "what," this section tells the "why."
See claude-instructions/development-workflow.md for the full PR flow.
-->

## Summary

<!-- 1–3 bullets: what and why. -->

-
-

## Related issues

<!-- Use `Fixes #N` if this PR closes the issue; `Part of #N` if it's a contribution. -->

Fixes #

## Test plan

<!-- Exactly how you verified it. What should a reviewer re-run to confirm? -->

- [ ] Unit tests added/updated and passing
- [ ] Integration tests added/updated (if boundaries changed)
- [ ] e2e tests (Playwright or equivalent) added/updated (if UI or user flow changed)
- [ ] Regression test added (if this is a bug fix — required)
- [ ] Manual verification steps:

## Security considerations

<!-- Does this change touch auth, data flow, tool use, memory, logs, external integrations, or
     any OWASP Web/LLM/Agentic category? If yes, list the items and what you did. If no, write "n/a". -->

## Docs updated

- [ ] Module README updated (if the module's contract changed)
- [ ] ADR added (if this implements an architectural decision)
- [ ] claude-instructions updated (if Claude's behavior on this project changed)
- [ ] Other docs: <bucket/file>

## Breaking changes

<!-- "None" or describe what breaks and the migration path. -->

## Checklist

- [ ] Branch follows `feat/N-slug | fix/N-slug | chore/slug | docs/slug` naming.
- [ ] Commits follow Conventional Commits (see `claude-instructions/commit-conventions.md`).
- [ ] Pre-commit hooks pass locally.
- [ ] Security workflows green (or explained if skipped).
