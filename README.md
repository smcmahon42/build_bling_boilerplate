# build_bling_boilerplate

A language-agnostic template for new AI / engineering projects. Ships the practices — not the scaffolding — that make a codebase legible to both humans and Claude from day one.

What you get:

- A **slim-router `CLAUDE.md`** that points at on-demand topic files in `claude-instructions/`.
- **Architectural Decision Records** in `docs/decisions/` with a template, index, and meta-ADR explaining the practice.
- **Three OWASP lists wired in** (Web Top 10, LLM Top 10, Agentic AI Top 10) with a security-review skill and language-detecting CI.
- **TDD as the starting posture**, with Playwright as the standard for frontend e2e.
- **GitHub templates** for PRs, bugs, features, epics, and ADR proposals, plus a baseline label taxonomy.
- **Pre-commit hooks** (gitleaks, hygiene, markdownlint) and CI workflows (security-scan, CodeQL, SBOM, license-check, ADR-lint).
- **Claude memory templates** you can seed into `~/.claude/projects/<path>/memory/` for this project.
- A **`/bootstrap-project` skill** that walks you through setup: placeholders, optional trim, memory seeding, GitHub labels, first ADR.

## Who this is for

You're starting a new project — especially one involving AI agents, LLM-driven features, or just a codebase where you want Claude Code to be a useful collaborator from day one — and you'd rather not reinvent the same dev/doc/security practices every time.

You're **not** looking for a runtime scaffold (no Go / Next.js / Python project structure here). Pair this with whatever language scaffold you normally use.

## Use as a template

1. Click **"Use this template"** on GitHub (or `gh repo create <name> --template=...`), then clone your new repo.
2. Open the clone in Claude Code.
3. Run the bootstrap skill:

   ```
   /bootstrap-project
   ```

   It will:
   - Ask for project name, description, stack, frontend yes/no.
   - Replace `{{PLACEHOLDERS}}` across the tree.
   - Offer to seed your user-local Claude memory for this project.
   - Offer to `gh label sync` the label taxonomy.
   - Offer to scaffold a Playwright `e2e/` directory if you have a frontend.
   - Offer to kick off your project's first ADR.
   - Remove itself when done.
4. Commit and push. You're set up.

If you prefer to do it manually: grep for `{{` and replace the placeholders, then copy `templates/memory/` into `~/.claude/projects/<encoded-path>/memory/` using the instructions in [`templates/memory/README.md`](templates/memory/README.md).

## File map

```
.
├── CLAUDE.md                           # slim router + house rules
├── claude-instructions/                # topic files loaded on-demand
│   ├── development-workflow.md         # TDD loop, branching, PR flow
│   ├── testing-practices.md            # test pyramid, Playwright e2e
│   ├── security-practices.md           # OWASP Web / LLM / Agentic
│   ├── documentation-discipline.md     # doc routing rules
│   ├── commit-conventions.md           # Conventional Commits
│   ├── dependency-vetting.md           # 5-point dep check
│   ├── component-explainability.md     # README at every boundary
│   ├── claude-memory.md                # user-local memory conventions
│   ├── epics-and-projects.md           # doc-embedded vs Projects v2
│   └── scaling-claude-instructions.md  # growing the instruction system
├── docs/
│   └── decisions/                      # ADRs + template + index
├── .claude/
│   ├── settings.json                   # shared Claude permissions
│   ├── skills/                         # bootstrap-project, security-review, dep-audit, test-gaps, new-adr
│   └── agents/                         # (empty — add project-specific agents here)
├── .github/
│   ├── PULL_REQUEST_TEMPLATE.md
│   ├── ISSUE_TEMPLATE/                 # bug, feature, epic, adr_proposal
│   ├── labels.yml                      # baseline taxonomy
│   ├── dependabot.yml
│   └── workflows/                      # pre-commit, adr-lint, security-scan, codeql, dast-scan, sbom, license-check
├── templates/
│   └── memory/                         # seed files for user-local Claude memory
├── .pre-commit-config.yaml
├── .editorconfig
├── .gitignore
├── .markdownlint.yaml
├── CONTRIBUTING.md
└── LICENSE                              # MIT
```

## Growing the template

`CLAUDE.md` starts as a slim router to `claude-instructions/`. As the project grows into multiple modules, each module gets its own `CLAUDE.md` + `claude-instructions/` subdirectory. See [`claude-instructions/scaling-claude-instructions.md`](claude-instructions/scaling-claude-instructions.md) for the Stage 1 → 2 → 3 pattern.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Short version: keep it generic, ship docs with code, Conventional Commits, TDD where code changes.

## License

MIT — see [LICENSE](LICENSE).
