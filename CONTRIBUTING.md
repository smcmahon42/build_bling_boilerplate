# Contributing

Thanks for your interest in contributing. Whether you're fixing a typo, proposing a practice, or generalizing something from your own project, the flow is the same.

## Ground rules

1. **Keep it generic.** This is a public template for *anyone's* project. Don't land changes that only fit one project, one stack, or one organization.
2. **Practices over tools.** When the choice is "explain the practice" vs. "ship a bespoke tool," explain the practice. Tools rot; practices travel.
3. **Small, focused PRs.** One change per PR. Split if it sprawls.
4. **Docs change with code.** If the shape of the boilerplate changes, the docs explaining it change in the same PR.

## Quick start

```bash
# Clone
git clone https://github.com/<your-org>/build_bling_boilerplate
cd build_bling_boilerplate

# Install hooks
pip install pre-commit
pre-commit install

# Work on a branch
git checkout -b feat/<slug>
```

Then follow the flow in [`claude-instructions/development-workflow.md`](claude-instructions/development-workflow.md).

## Commit messages

Conventional Commits. See [`claude-instructions/commit-conventions.md`](claude-instructions/commit-conventions.md).

## What we will and won't accept

**Will:**

- Sharper or clearer wording in the instruction files.
- New topic files that capture a recurring practice.
- New skills in `.claude/skills/` that are broadly useful.
- New CI workflows that are language-detecting and cost-free.
- ADRs that record the boilerplate's own design decisions.

**Won't (without discussion first):**

- Stack-specific scaffolding (Go project layout, Next.js project layout, etc.). These belong in a downstream fork.
- Project-management tools that assume a specific platform beyond GitHub.
- Dependencies. The boilerplate aims to be runnable with nothing but git, `gh`, and pre-commit.
- Long-form tutorials. Keep docs action-oriented.

## Proposing a significant change

Open an ADR proposal issue using `.github/ISSUE_TEMPLATE/adr_proposal.md`. Discuss there; on acceptance, the content moves into `docs/decisions/NNNN-<slug>.md` via a PR. See `.claude/skills/new-adr.md`.

## Code of conduct

Be decent. Disagreement is fine; contempt is not. Report issues to the maintainer via GitHub.

## License

By contributing, you agree that your contributions are licensed under the project's MIT license.
