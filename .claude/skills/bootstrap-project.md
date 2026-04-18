---
name: bootstrap-project
description: One-time setup skill for a project cloned from the build_bling_boilerplate template. Replaces placeholders, trims optional pieces, seeds Claude memory, and offers to initialize GitHub labels/hooks. Delete this file after running.
---

# Bootstrap project

Use this skill **once**, right after cloning the boilerplate for a new project. Walk through each step with the operator; confirm before taking actions that touch anything outside the repo (home directory, `gh` commands).

## Step 0 — Sanity check

- Confirm the repo was cloned from `build_bling_boilerplate` (look for this skill file; the template README's "use as template" banner).
- Confirm the current directory is the intended new-project root, not the boilerplate itself.
- Confirm the operator is ready to commit changes — this skill makes many edits.

## Step 1 — Collect project identity

Ask the operator (use AskUserQuestion for each):

- **Project name** (e.g., `acme-recs`). Used in `{{PROJECT_NAME}}`.
- **One-line description**. Used in `{{PROJECT_DESCRIPTION}}`.
- **Primary stack** — free-form (e.g., "Go 1.22 + Postgres + React 18"). Used in `{{PRIMARY_STACK}}`.
- **Project status** — one of `exploring`, `alpha`, `beta`, `production`. Used in `{{PROJECT_STATUS}}`.
- **Frontend?** — yes/no. Determines whether to scaffold Playwright later.

Replace placeholders across every file in the repo:

```bash
# Operator: run or let Claude run
grep -rl "{{PROJECT_NAME}}" . --exclude-dir=.git
# Then edit each match; do not use blind sed — confirm each file
```

## Step 2 — Trim optional practice areas

Ask which practice areas to keep. Default-on for all; offer to drop any. For each dropped area, remove the relevant files and references:

- **ADRs** (`docs/decisions/`) — removing these is rare but possible for small scripts.
- **GitHub issue templates** — drop `.github/ISSUE_TEMPLATE/` if not using GitHub.
- **Pre-commit hooks** — drop `.pre-commit-config.yaml` if the project won't use them.
- **CI workflows** — let the operator cherry-pick which workflows in `.github/workflows/` to keep. Defaults: keep `pre-commit.yml`, `adr-lint.yml`, `security-scan.yml`, `codeql.yml`, `license-check.yml`. Optional: `dast-scan.yml` (requires a deployed preview URL), `sbom.yml`.
- **Memory seed** — offer to skip if the operator doesn't use Claude Code memory.

For every dropped section, also update `CLAUDE.md` and `claude-instructions/README.md` to remove dangling references.

## Step 3 — Seed Claude memory (confirm before touching home dir)

If the operator wants memory seeding:

1. Compute the project's Claude memory path: `~/.claude/projects/<encoded-absolute-path>/memory/` where the path has `/` replaced with `-`.
2. Show the operator the computed path and ask for confirmation.
3. Create the directory if it doesn't exist.
4. Copy `templates/memory/*.example` → `<memory-path>/*.md` (strip the `.example` suffix).
5. Edit the copied files: replace placeholders (`{{PROJECT_NAME}}`, `{{PRIMARY_STACK}}`) with project-specific values.
6. Update `MEMORY.md` index entries to match the seeded files.

**Do not** commit memory files to the repo — they are user-local by design.

## Step 4 — Scaffold Playwright (if frontend = yes)

If the operator indicated a frontend:

1. Ask where the frontend lives (repo root, `apps/web/`, or custom).
2. Create an `e2e/` directory there with:
   - `playwright.config.ts` using `trace: 'retain-on-failure'` and `video: 'retain-on-failure'`.
   - `tests/smoke.spec.ts` as a starter (a single spec that loads the homepage).
   - `fixtures/` (empty).
3. **Do not** run `npm install` — let the operator install Playwright in their package manager.
4. Add a note to the project README with the install commands.

## Step 5 — GitHub setup (confirm before any `gh` calls)

Ask the operator whether to run these:

- **Install pre-commit hooks locally:** `pre-commit install` (runs locally, fine to auto-run).
- **Sync labels:** `gh label sync .github/labels.yml`. Requires `gh` authenticated. Ask before running.
- **Create initial epic issue:** offer to use `.github/ISSUE_TEMPLATE/epic.md` to file issue #1 as the inaugural epic. Confirm title and body with operator before filing.

Never run `gh repo create`, `gh pr create`, or any write operation without explicit operator confirmation.

## Step 6 — First ADR

Offer to run the `/new-adr` skill to scaffold `docs/decisions/0002-<slug>.md` for the project's foundational architectural decision (e.g., "choice of runtime", "monolith vs services"). Explain that `0001-record-architecture-decisions.md` is the meta-ADR and should remain.

## Step 7 — Rewrite README

The template `README.md` currently describes how to use the template. Replace it with a project-appropriate README. Ask the operator if they want:

- A minimal stub (one-liner + setup instructions), or
- A full template with Install / Usage / Development / Contributing / License sections.

Preserve the LICENSE file as-is (MIT).

## Step 8 — Self-cleanup

Once the operator confirms the bootstrap is complete:

1. Offer to delete this skill file (`.claude/skills/bootstrap-project.md`) — it's no longer needed.
2. Offer to remove any `{{PLACEHOLDER}}` strings that remain (grep for `{{[A-Z_]*}}` and list remaining matches).
3. Offer to commit the result with: `chore: bootstrap project from build_bling_boilerplate template`.
4. Suggest the operator push to the remote and create the initial epic issue.

## Invariants

Throughout the skill:

- **Never** edit files outside the repo without explicit confirmation for the specific path.
- **Never** run `gh` write operations without confirmation.
- **Never** commit memory files to the repo.
- **Never** delete the operator's uncommitted work — check `git status` before destructive operations.
- If any step fails, stop and report. Don't retry the same command; diagnose first.
