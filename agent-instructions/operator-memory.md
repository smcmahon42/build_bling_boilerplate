# Operator memory

Some agent clients support client-local memory: facts about the operator, feedback from past sessions, and short-lived project context stored outside the repo. This file defines the boundary and shape for that memory without making the repo depend on one LLM provider.

Claude Code currently has first-class memory templates in `templates/memory/`. Codex or homegrown agents can adopt the same content model if they provide their own local memory store.

## What belongs in client-local memory

| Type | Purpose | Example |
| --- | --- | --- |
| **user** | Who the operator is, their role, what they already know | "Ten years of Go experience, new to React." |
| **feedback** | What to keep doing or stop doing based on past interactions | "Don't mock the DB in integration tests." |
| **project** | Time-bound facts about the current project | "Release cut scheduled 2026-05-12; bug-only after." |
| **reference** | Pointers to external systems | "Pipeline bugs tracked in Linear project INGEST." |

Each memory file should carry frontmatter:

```markdown
---
name: feedback_testing
description: Preferred testing approach in this project
type: feedback
---

<content>
```

## What does not belong in memory

- **Code patterns or conventions** — these live in `AGENTS.md` and `agent-instructions/`.
- **Git history or recent changes** — `git log` is authoritative.
- **Debugging recipes** — the fix belongs in code; durable context belongs in commit messages or docs.
- **Transient task state** — use the active agent client's task tracker.
- **Project-local session handoff** — open, blocked, tabled, and recently completed work lives in committed `STATE.md`.

## Memory vs project state

| Concern | Lives in | Lifecycle | Audience |
| --- | --- | --- | --- |
| Who the operator is and how they work | Client-local memory | Slow-moving | This operator only |
| Time-bound project facts visible only to one operator | Client-local memory | Decays as project changes | This operator only |
| What's open / in-progress / blocked / tabled | Repo root `STATE.md` | Updated by session workflows | Every contributor and agent |
| Accepted architectural decisions | `docs/decisions/` | Append-only | Every contributor and agent |
| Agent doctrine and process rules | `AGENTS.md` + `agent-instructions/` | Updated by PR/ADR as needed | Every contributor and agent |

When in doubt: if the next contributor on first clone needs to see it, it belongs in the repo. If only one operator's agent client needs it, it belongs in client-local memory.

## Claude Code adapter

Claude Code stores per-project memory under:

```text
~/.claude/projects/<encoded-project-path>/memory/
```

The boilerplate ships Claude-compatible examples in `templates/memory/`. The `/bootstrap-project` skill can seed them, or the operator can copy them manually.

## Portable memory shape

Use namespaced markdown files so memory can migrate between clients:

```text
user_role.md
feedback_testing.md
feedback_pr_style.md
project_release_cut.md
reference_linear.md
```

Keep the index concise, one line per memory. Delete stale entries aggressively. Fewer accurate memories beat many stale ones.
