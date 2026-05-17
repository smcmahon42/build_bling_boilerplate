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
- **`STATE.md` + `/start-session` + `/end-session`** — project-local session-handoff ledger so the next session grounds in one file read instead of re-grepping the repo. See [`claude-instructions/session-handoff.md`](claude-instructions/session-handoff.md) and [ADR-0008](docs/decisions/0008-session-handoff-state.md).
- **Review state + Provenance on every STATE.md row** — agent-authored entries start as evidence (`Review: unreviewed`) until the operator confirms; an optional `Provenance` block records skill, session, prompt summary, and context. Keeps the line between agent observation and operator direction explicit. See [ADR-0009](docs/decisions/0009-provenance-and-review-state-on-state-rows.md).
- **Four-level agent-autonomy doctrine** — L1 read, L2 autonomous-on-branch (plan-first on medium/large), L3 propose-and-confirm, L4 human-only. Resolves the in-between of "this isn't a destructive op but it's not trivial either." See [`claude-instructions/agent-autonomy.md`](claude-instructions/agent-autonomy.md) and [ADR-0010](docs/decisions/0010-agent-autonomy-scope.md).
- **Per-feature cost signals on STATE.md entries** — sessions to date, skills used, operator turns, files read. Captured by `/end-session` on multi-session entries; external billing / OpenTelemetry data composes via the session id already in Provenance. Makes the "skills over recall" intuition testable. See [`claude-instructions/agent-cost-observability.md`](claude-instructions/agent-cost-observability.md) and [ADR-0011](docs/decisions/0011-agent-cost-observability.md).

## Who this is for

You're starting a new project — especially one involving AI agents, LLM-driven features, or just a codebase where you want Claude Code to be a useful collaborator from day one — and you'd rather not reinvent the same dev/doc/security practices every time.

The session-handoff, review-state, autonomy, and cost-signal primitives assume work will span many sessions and possibly many contributors (or many agents) over time. If your project genuinely lives in a single session by a single person, those four pieces are overhead you can ignore — the rest of the template still applies.

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

If you prefer to do it manually: grep for `{{` and replace the placeholders, then copy `templates/memory/` into `~/.claude/projects/<encoded-path>/memory/` using the instructions in [`templates/memory/README.md`](templates/memory/README.md). Also copy `templates/state/STATE.md` to the project root and replace `{{PROJECT_NAME}}` — that's the session-handoff ledger described below.

## A typical session

The boilerplate is built for projects worked across many sessions and contributors. Each non-trivial session follows the same loop:

1. **`/start-session`** — the agent reads `STATE.md` before any other repo exploration. It surfaces blocked items, *unreviewed* agent-authored items (entries the agent added since you last looked), in-progress work, and open questions. You direct what's next instead of asking the agent to re-derive ground state.
2. **Do the work** — per `CLAUDE.md`'s house rules. House rule 3 (share-plan-before-editing) governs medium and large tasks. The agent operates at **L2 autonomy** on a feature branch, escalating to **L3 (propose-and-confirm)** for anything touching `main`, settings, doctrine, dependencies, or CI; **L4** operations (force-push, production writes, external messages, infra) are human-only.
3. **`/end-session`** — the agent reconciles `STATE.md` against what actually happened: completed items move to recent with a one-sentence summary, last-touched dates update, new items get added with `Review` set per author defaults (`confirmed` when operator-directed, `unreviewed` when agent-noticed), tabled items are marked with reason and un-table condition, open questions are recorded. `Cost signals` (sessions to date, operator turns, skills used) are captured on multi-session entries. `STATE.md` is committed separately from code.

Architectural decisions that crystallize get an ADR (`/new-adr`); decisions still in flux stay in `STATE.md` and graduate later.

### How the four primitives compose

| Primitive | Carries | Read at | Updated at |
| --- | --- | --- | --- |
| `STATE.md` (slice 7) | Open / completed / tabled / open-questions | `/start-session` | `/end-session` |
| Review + Provenance (slice 8) | Evidence-vs-instruction line, agent-side origin | every read | when new entries are added or operator reviews |
| Autonomy levels (slice 9) | What the agent does without asking | every action | rarely (doctrine change → ADR) |
| Cost signals (slice 10) | Per-entry cost over its lifetime | `/start-session` (visible on items) | `/end-session` |

Every ledger entry carries its origin (Provenance), authority (Review), and the cost it has consumed (Cost signals). The autonomy doctrine governs what the agent can do *to* the ledger and the rest of the repo. The session-handoff skills are the routine that keeps the whole thing accurate.

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
│   ├── session-handoff.md              # STATE.md doctrine + Review/Provenance
│   ├── agent-autonomy.md               # L1-L4 autonomy doctrine
│   ├── agent-cost-observability.md     # Cost signals on STATE.md entries
│   ├── epics-and-projects.md           # doc-embedded vs Projects v2
│   └── scaling-claude-instructions.md  # growing the instruction system
├── docs/
│   └── decisions/                      # ADRs + template + index
├── .claude/
│   ├── settings.json                   # shared Claude permissions
│   ├── skills/                         # bootstrap-project, security-review, dep-audit, test-gaps,
│   │                                   # new-adr, new-skill, new-contract, start-session, end-session
│   └── agents/                         # (empty — add project-specific agents here)
├── .github/
│   ├── PULL_REQUEST_TEMPLATE.md
│   ├── ISSUE_TEMPLATE/                 # bug, feature, epic, adr_proposal
│   ├── labels.yml                      # baseline taxonomy
│   ├── dependabot.yml
│   └── workflows/                      # pre-commit, adr-lint, security-scan, codeql, dast-scan, sbom, license-check
├── templates/
│   ├── memory/                         # seed files for user-local Claude memory
│   └── state/                          # STATE.md seed (project-local session-handoff ledger)
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
