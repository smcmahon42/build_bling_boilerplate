# Claude memory

Claude Code stores per-project memory on the **user's local machine**, not in the repo. This is deliberate: memories hold preferences, feedback, and context that's specific to the individual user collaborating with Claude. They aren't shared with teammates, and they aren't committed to git.

This file explains:
- Where memory lives.
- How to structure it (namespacing, types).
- How the boilerplate seeds it.

## Where memory lives

Claude Code's memory system writes to a per-project directory under your home:

```
~/.claude/projects/<encoded-project-path>/memory/
```

The `<encoded-project-path>` is your repo's absolute path with `/` replaced by `-`. For example, a repo at `/Users/alex/code/myapp` becomes `-Users-alex-code-myapp`.

Inside `memory/` you'll find:

```
memory/
  MEMORY.md              # index, always loaded into Claude's context
  user_role.md           # who you are, how you like to work
  feedback_testing.md    # feedback about testing approach
  project_deadline.md    # time-boxed project context
  reference_dashboards.md
  ...
```

## The four memory types

| Type | Purpose | Example |
| --- | --- | --- |
| **user** | Who you are, your role, what you already know | "Ten years of Go experience, new to React." |
| **feedback** | What to keep doing, what to stop doing — based on past interactions | "Don't mock the DB in integration tests — got burned by divergence before." |
| **project** | Time-bound facts about the current project: deadlines, stakeholders, constraints | "Release cut scheduled 2026-05-12; only bug fixes after that date." |
| **reference** | Pointers to external systems (dashboards, Linear projects, Slack channels) | "Pipeline bugs tracked in Linear project INGEST." |

Each memory file carries frontmatter:

```markdown
---
name: feedback_testing
description: Preferred testing approach in this project
type: feedback
---

<content>
```

## Namespacing convention

Use a prefix that matches the memory type, followed by a short topic slug. This makes the index scan-friendly and grep-friendly.

```
user_role.md
user_background_go.md
feedback_testing.md
feedback_pr_style.md
project_release_cut.md
project_team_ownership.md
reference_linear.md
reference_grafana.md
```

When a project grows, prefix with the module:

```
feedback_api_testing.md        # testing advice specific to the api module
feedback_web_e2e.md
project_web_redesign.md
```

## The MEMORY.md index

`MEMORY.md` is always loaded into Claude's context — it's how Claude knows which detailed memory files exist. Keep it **concise**: one line per memory, under 200 lines total (lines beyond ~200 get truncated).

Format:

```markdown
- [Title](file.md) — one-line hook
```

Example:

```markdown
# Memory index

## User
- [Role](user_role.md) — Staff engineer, backend focus, recent React work
- [Background](user_background_go.md) — 10y Go; prefers minimalism

## Feedback
- [Testing approach](feedback_testing.md) — No mocks for owned code; real DB in integration
- [PR style](feedback_pr_style.md) — Small, single-purpose PRs; no squash of meaningful history

## Project
- [Release cut](project_release_cut.md) — 2026-05-12; bug-only after

## Reference
- [Linear](reference_linear.md) — Project INGEST for pipeline bugs
```

No frontmatter on `MEMORY.md` itself — it's an index, not a memory.

## What NOT to save

- **Code patterns or conventions** — these live in `CLAUDE.md` and `claude-instructions/` and are readable from the repo.
- **Git history or recent changes** — `git log` is authoritative.
- **Debugging recipes** — the fix is in the code; the commit message holds the context.
- **Transient task state** — use Claude's task tracker, not memory.

## Seeding memory for a new project

The boilerplate ships examples in `templates/memory/`. The `/bootstrap-project` skill offers to copy them into the right user-local path, with project-specific placeholders filled in.

You can also seed manually:

```bash
# Replace <your-project-path> with what Claude Code uses for this repo
cp templates/memory/*.example ~/.claude/projects/<your-project-path>/memory/
# Rename .example files to .md and edit to taste
```

## Maintenance

- **Update memories that turn out to be wrong.** Don't leave stale entries.
- **Delete memories that don't matter anymore.** Fewer, better memories beat many stale ones.
- **Don't cross-pollinate.** Memory for *this* project should be about this project. Put cross-project preferences in a separate profile or in `~/.claude/CLAUDE.md`.
