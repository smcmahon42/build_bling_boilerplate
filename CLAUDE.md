# CLAUDE.md — Claude Code adapter for {{PROJECT_NAME}}

Claude Code should treat `AGENTS.md` as the canonical project router and `agent-instructions/` as the shared doctrine. This file exists only for Claude-specific startup and slash-command discovery.

## Startup

1. Read `AGENTS.md` first.
2. Read the relevant topic files in `agent-instructions/`.
3. Scan `.claude/skills/README.md` and `skills-index.yaml` before non-trivial work.
4. Follow Claude Code permission settings in `.claude/settings.json`.

## Claude-specific workflows

- Use `/bootstrap-project` after cloning the template.
- Use `/start-session` and `/end-session` for non-trivial session handoff.
- Use `/security-review`, `/dep-audit`, `/test-gaps`, `/new-adr`, `/new-contract`, and `/new-skill` when their triggers apply.
- Store optional Claude Code memory under `~/.claude/projects/<encoded-path>/memory/`; do not commit it.

## Doctrine boundary

Do not duplicate house rules here. If project doctrine changes, update `AGENTS.md` or `agent-instructions/` and treat this file as an adapter shim.
