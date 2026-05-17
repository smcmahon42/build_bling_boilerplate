# Scaling agent instructions as the project grows

This directory structure — neutral `AGENTS.md` router + topic files under `agent-instructions/` + optional client adapters — is **Stage 2** of a three-stage pattern. Most projects start here and stay here for a long time.

## The three stages

### Stage 1 — Single router

**Shape:**

```
<repo>/
  AGENTS.md     # house rules, conventions, links to key files
```

**Fits when:**

- The project is new.
- One agent or developer can hold the whole thing in context.
- There are not enough recurring concerns to warrant separate topic files.

**Signals it's time to move to Stage 2:**

- `AGENTS.md` crosses ~200 lines.
- Different task types compete for the same space.
- Agent clients repeatedly load guidance unrelated to the task.

### Stage 2 — Neutral router + topic files + adapters

**Shape:** where this boilerplate starts.

```
<repo>/
  AGENTS.md                         # neutral router
  CLAUDE.md                         # Claude Code adapter shim
  agent-instructions/
    README.md                       # topic index
    development-workflow.md
    testing-practices.md
    security-practices.md
    ...                             # one file per recurring concern
  .claude/                          # Claude-specific skills/settings
```

**Fits when:**

- The project has several distinct concerns each worth its own file.
- Multiple agent clients may work in the same repo.
- The team wants one doctrine and multiple client adapters.

**Signals it's time to move to Stage 3:**

- The repo hosts multiple modules with distinct stacks.
- A topic file has to hedge by module.
- Onboarding to one module requires reading rules that do not apply to it.

### Stage 3 — Per-module routers

**Shape:**

```
<repo>/
  AGENTS.md                               # top-level map of modules + cross-cutting topics
  CLAUDE.md                               # optional Claude root adapter
  agent-instructions/                     # cross-cutting topics only
    README.md
    commit-conventions.md
    security-practices.md
    ...
  apps/
    api/
      AGENTS.md                           # slim router for the API service
      agent-instructions/
        architecture.md
        engineering-standards.md
        module-specific-security.md
    web/
      AGENTS.md
      agent-instructions/
        architecture.md
        engineering-standards.md
        e2e-tests.md
```

**What stays cross-cutting:**

- Commit conventions.
- PR flow and review gates.
- TDD and testing principles.
- Doc discipline and doc routing.
- Operator-memory boundaries.
- OWASP security principles.
- Adapter policy.

**What becomes module-local:**

- Architecture narrative.
- Engineering standards specific to the stack.
- Test framework particulars.
- Stack-specific security concerns.
- Build and run commands.

## How to split from Stage 2 to Stage 3

1. **Create `<module>/AGENTS.md`** as a slim router. Lead with the module's purpose in one sentence.
2. **Create `<module>/agent-instructions/`** with 2–5 topic files. Start with `architecture.md` and `engineering-standards.md`.
3. **Move module-specific sections** out of root topic files into module topic files. Leave the principle in the root file and put the application in the module.
4. **Update root `AGENTS.md`** with a module table.
5. **Update root `agent-instructions/README.md`** to explain that module-specific guidance now lives under each module.
6. **Update client adapters** only if they need module-specific invocation mechanics.

## How much to split

- **Keep routers under ~120 lines.** If longer, topic files are not carrying their weight.
- **Do not duplicate.** If a rule applies to all modules, it lives at the root.
- **Do not pre-split.** Split when divergence is real, not anticipated.

## Related

- The root `AGENTS.md` of this boilerplate is the Stage 2 reference implementation.
- See [documentation-discipline.md](documentation-discipline.md) for the parallel scaling pattern in `docs/`.
