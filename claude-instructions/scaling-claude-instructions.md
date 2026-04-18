# Scaling Claude instructions as the project grows

This directory structure — slim `CLAUDE.md` router + topic files under `claude-instructions/` — is **Stage 2** of a three-stage pattern. Most projects start here and stay here for a long time. This file explains when and how to move between stages so the instruction system grows *with* the project instead of collapsing under its own weight.

## The three stages

### Stage 1 — Single file

**Shape:**

```
<repo>/
  CLAUDE.md     # house rules, conventions, links to key files
```

**Fits when:**
- The project is new (under a few thousand lines).
- A single developer (or Claude) can hold the whole thing in context.
- There aren't yet enough distinct concerns to warrant separation.

**Signals it's time to move to Stage 2:**
- `CLAUDE.md` crosses ~200 lines.
- Different task types (testing, security, deploys) start competing for the same space.
- You find yourself saying "skip the security section" when Claude reads the file for unrelated work.

### Stage 2 — Slim router + topic files

**Shape:** where this boilerplate starts.

```
<repo>/
  CLAUDE.md                          # slim router: ~100 lines, links to topic files
  claude-instructions/
    README.md                        # topic index
    development-workflow.md
    testing-practices.md
    security-practices.md
    ...                              # one file per recurring concern
```

**Fits when:**
- The project has several distinct concerns each worth its own file.
- Claude benefits from loading only the relevant topic.
- The team wants the index-style ergonomics of "if you're doing X, read Y."

**Signals it's time to move to Stage 3:**
- The repo now hosts **multiple modules with distinct stacks** (e.g. Go API + React admin + JS SDK).
- A single topic file has to hedge ("in the API service, do this; in the frontend, do that").
- Onboarding a developer to one module requires reading a page of rules that don't apply to their module.

### Stage 3 — Per-module routers

**Shape:**

```
<repo>/
  CLAUDE.md                                # top-level map of modules + links to cross-cutting topics
  claude-instructions/                     # CROSS-CUTTING topics only
    README.md
    commit-conventions.md
    security-practices.md                  # principles; each module adds stack-specific notes
    ...
  apps/
    api/
      CLAUDE.md                            # slim router for the API service
      claude-instructions/
        architecture.md
        engineering-standards.md
        module-specific-security.md
    web/
      CLAUDE.md
      claude-instructions/
        architecture.md
        engineering-standards.md
        e2e-tests.md                       # Playwright config / fixtures particular to this app
    sdk/
      CLAUDE.md
      claude-instructions/
        architecture.md
        bundle-size-budget.md              # JS-specific concern
```

**What stays cross-cutting (root `claude-instructions/`):**
- Commit conventions.
- PR flow and review gates.
- TDD and testing *principles* (the pyramid stays; specific framework notes move to modules).
- Doc discipline and doc routing.
- Claude memory conventions.
- OWASP security principles (module-specific checks move to modules).

**What becomes module-local:**
- Architecture narrative.
- Engineering standards specific to the stack.
- Test framework particulars (vitest vs go test vs pytest).
- Stack-specific security concerns (e.g. CSP for web, SQL injection for DB-heavy services).
- Build and run commands.

## How to split from Stage 2 to Stage 3

1. **Create `<module>/CLAUDE.md`** as a slim router — same pattern as the root. Lead with the module's purpose in one sentence.
2. **Create `<module>/claude-instructions/`** with 2–5 topic files. Start with `architecture.md` and `engineering-standards.md`; add more only when earned.
3. **Move module-specific sections** out of root topic files into the module's topic files. Leave the *principle* in the root file and put the *application* in the module.
4. **Update root `CLAUDE.md`** with a table of modules:
   ```
   | Module | Stack | CLAUDE.md |
   | --- | --- | --- |
   | apps/api | Go 1.22 + Postgres | apps/api/CLAUDE.md |
   | apps/web | React 18 + Vite | apps/web/CLAUDE.md |
   ```
5. **Update root `claude-instructions/README.md`** to explain that module-specific guidance now lives under each module.
6. **Verify** by running `/security-review` or `/test-gaps` against one module and confirming Claude reads the right files.

## How much to split

- **Keep module `CLAUDE.md` under ~100 lines.** If it's longer, its own topic files aren't carrying their weight.
- **Don't duplicate.** If a rule applies to all modules, it lives at the root. If you find the same rule in two modules, it's cross-cutting.
- **Don't pre-split.** Three similar modules are easier to reason about than three modules hiding behind abstraction. Split when the divergence is real, not anticipated.

## Related

- The root `CLAUDE.md` of this boilerplate is the Stage 2 reference implementation. Mirror its shape for modules.
- See [documentation-discipline.md](documentation-discipline.md) for the parallel scaling pattern in `docs/`.
