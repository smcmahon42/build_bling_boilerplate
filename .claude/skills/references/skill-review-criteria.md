# Skill review criteria

Quality bar for skills in `.claude/skills/`. Use this checklist both when authoring (`/new-skill`) and when reviewing a PR that adds or modifies a skill.

## Frontmatter

- [ ] `name` is a lowercase hyphen-separated slug (e.g. `security-review`).
- [ ] `name` matches the filename stem and the `skills-index.yaml` entry.
- [ ] `description` is one sentence, under 200 characters, action-oriented. Starts with a verb. Explains *when* to invoke, not just what it does.

## Purpose

- [ ] The skill solves a **recurring** task. One-off tasks don't need a skill.
- [ ] The skill has a **single responsibility**. If the description needs "and" to explain it, split into two skills.
- [ ] The task couldn't be handled adequately with a prompt alone — the skill adds structure, invariants, or mandated steps.
- [ ] There isn't already a skill that covers this task. If there's overlap, update the existing skill instead.

## Body

- [ ] **Starts with context** — one paragraph stating what the skill does and when to invoke.
- [ ] **Scannable** — uses headings, tables, and numbered lists rather than prose paragraphs. A reader should find the right section in under 10 seconds.
- [ ] **Action-oriented** — "do X", "ask Y", "check Z" — not narrative or theory.
- [ ] **Steps are numbered** where order matters; bullets where it doesn't.
- [ ] **Invariants are called out.** What must the skill never do? What must it always do? (e.g., "never push without confirmation", "always cite file:line for findings").
- [ ] **Output format is explicit** when the skill produces an artifact. Table columns, section structure, or example output are shown.
- [ ] **Confirmation gates** are explicit for any destructive action (file deletion, external-facing writes, home-dir edits, `gh` write operations).
- [ ] **Self-contained.** Doesn't rely on information the user would have to provide separately. If it does need input, the skill asks for it via `AskUserQuestion`.

## Scope

- [ ] **Language-agnostic** — if the skill detects a stack, it degrades gracefully when the stack isn't present. No silent failures.
- [ ] **No hard-coded secrets, credentials, or domain-specific identifiers.**
- [ ] **No opinions that conflict with `CLAUDE.md` or the cross-cutting `claude-instructions/` files.** If it needs to, it should update those files instead.

## Index integration

- [ ] Added to `skills-index.yaml` with accurate fields.
- [ ] `triggers` use vocabulary from the index header; don't invent new tags casually.
- [ ] `scale` reflects which task sizes this skill applies to (small / medium / large).
- [ ] `tools` lists exactly the tools the skill uses (don't over-request).
- [ ] `lifecycle` is `once` (run once per project) or `recurring`.
- [ ] The human-readable table in `.claude/skills/README.md` is updated.

## Related topic files

- [ ] If the skill references cross-cutting rules (testing, security, docs), it links to the topic file rather than duplicating content.
- [ ] If the skill introduces a new recurring practice, the corresponding topic file in `claude-instructions/` is updated in the same PR.

## Testability

- [ ] A dry-run of the skill (reading the markdown only, without editing) makes the intended behavior obvious.
- [ ] The skill degrades sensibly on the common failure modes (missing dependency, no detected stack, user declines a prompt).

## Size

- [ ] Under ~250 lines. Skills longer than this should split, extract to a `references/` file, or link out to a topic file.

## Final check

- [ ] Read the skill top-to-bottom out loud (mentally). Does the shape match what you'd want Claude to do? Are there implicit steps you'd need to remember?
- [ ] Imagine a reviewer who's never seen this skill invoking it for the first time. Is anything unclear?
