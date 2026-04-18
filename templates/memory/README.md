# Memory templates

Seed files for Claude Code's user-local memory system. These are **not** copied into the repo — they're copied into your home directory so Claude can load them across sessions on this specific project.

See [`claude-instructions/claude-memory.md`](../../claude-instructions/claude-memory.md) for the underlying concepts.

## Where this goes

Claude Code stores per-project memory at:

```
~/.claude/projects/<encoded-project-path>/memory/
```

where `<encoded-project-path>` is your repo's absolute path with `/` replaced by `-`. For a repo at `/Users/alex/code/myapp`, that's:

```
~/.claude/projects/-Users-alex-code-myapp/memory/
```

## What's in this directory

| File | Becomes |
| --- | --- |
| `MEMORY.md` | `~/.claude/projects/<path>/memory/MEMORY.md` (index, always loaded) |
| `user_role.md.example` | `user_role.md` — who the operator is |
| `feedback_example.md.example` | `feedback_<topic>.md` — preferences and corrections |
| `project_context.md.example` | `project_<topic>.md` — project-specific facts |
| `reference_example.md.example` | `reference_<topic>.md` — pointers to external systems |

## Automated seeding

Run the `/bootstrap-project` skill. It will:

1. Compute your project's memory path.
2. Confirm with you before creating anything under `~/.claude/`.
3. Copy these templates, strip the `.example` suffix, and replace placeholders.

## Manual seeding

If you prefer to do it by hand:

```bash
# 1. Compute the encoded path
PROJECT_ABS="$(pwd)"
ENCODED="${PROJECT_ABS//\//-}"
MEM_DIR="$HOME/.claude/projects/$ENCODED/memory"

# 2. Create it
mkdir -p "$MEM_DIR"

# 3. Copy the seed files (stripping the .example suffix)
for f in templates/memory/*.example; do
  name="$(basename "$f" .example)"
  cp "$f" "$MEM_DIR/$name"
done
cp templates/memory/MEMORY.md "$MEM_DIR/MEMORY.md"

# 4. Edit the copied files — remove placeholders, tailor to you.
```

## Do not commit memory files

Memory is user-local. The operator's role, feedback history, and project context are their own. Committing them would:

- Leak personal preferences into the shared repo.
- Create merge conflicts across teammates with different memories.
- Fight Claude Code's assumption that `~/.claude/` is authoritative.

The repo's `.gitignore` should exclude `.claude/projects/` at the home-dir level automatically — that path lives outside the repo by design.

## Maintenance

- Memories decay. Update or delete them as the project evolves.
- Keep `MEMORY.md` under ~200 lines (anything beyond is truncated).
- When a rule applies to *every* project, move it from memory into `~/.claude/CLAUDE.md` (your user-level config) instead.
