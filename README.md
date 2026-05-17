# build_bling_boilerplate

A language-agnostic template for new AI / engineering projects — built for codebases where agents are first-class collaborators and work spans many sessions and contributors, not just one.

It ships the **practices** (not the runtime scaffolding) that make a repo legible to humans and multiple agent clients from day one. The session-handoff, review-state, autonomy, and cost-observability primitives mean a new agent or contributor opening the project tomorrow already knows: what's open right now, what was just done, what the agent before them *proposed* (and whether you've reviewed it), what each feature has cost so far, *and* who you are and how you work — all from a single read of two files (`STATE.md` at the repo root, and optional client-local operator memory such as Claude Code's `~/.claude/projects/<path>/memory/`).

What you get:

- A **neutral `AGENTS.md` router** that points all agents at shared topic files in `agent-instructions/`.
- **Architectural Decision Records** in `docs/decisions/` with a template, index, and meta-ADR explaining the practice.
- **Three OWASP lists wired in** (Web Top 10, LLM Top 10, Agentic AI Top 10) with a security-review skill and language-detecting CI.
- **TDD as the starting posture**, with Playwright as the standard for frontend e2e.
- **GitHub templates** for PRs, bugs, features, epics, and ADR proposals, plus a baseline label taxonomy.
- **Pre-commit hooks** (gitleaks, hygiene, markdownlint) and CI workflows (security-scan, CodeQL, SBOM, license-check, ADR-lint).
- **Client-local operator memory templates** — so the agent starts every session knowing your role, stack, preferences, and current project context, not asking. Four memory types (`user`, `feedback`, `project`, `reference`) seeded into `~/.claude/projects/<path>/memory/`. See [How operator memory works in your projects](#how-operator-memory-works-in-your-projects) below.
- A **Claude `/bootstrap-project` adapter skill** that walks you through setup: placeholders, optional trim, memory seeding, GitHub labels, first ADR.
- **`STATE.md` + start/end-session workflows** — project-local session-handoff ledger so the next session grounds in one file read instead of re-grepping the repo. See [`agent-instructions/session-handoff.md`](agent-instructions/session-handoff.md) and [ADR-0008](docs/decisions/0008-session-handoff-state.md).
- **Review state + Provenance on every STATE.md row** — agent-authored entries start as evidence (`Review: unreviewed`) until the operator confirms; an optional `Provenance` block records skill, session, prompt summary, and context. Keeps the line between agent observation and operator direction explicit. See [ADR-0009](docs/decisions/0009-provenance-and-review-state-on-state-rows.md).
- **Four-level agent-autonomy doctrine** — L1 read, L2 autonomous-on-branch (plan-first on medium/large), L3 propose-and-confirm, L4 human-only. Resolves the in-between of "this isn't a destructive op but it's not trivial either." See [`agent-instructions/agent-autonomy.md`](agent-instructions/agent-autonomy.md) and [ADR-0010](docs/decisions/0010-agent-autonomy-scope.md).
- **Per-feature cost signals on STATE.md entries** — sessions to date, skills used, operator turns, files read. Captured by the end-session workflow on multi-session entries; external billing / OpenTelemetry data composes via the session id already in Provenance. Makes the "skills over recall" intuition testable. See [`agent-instructions/agent-cost-observability.md`](agent-instructions/agent-cost-observability.md) and [ADR-0011](docs/decisions/0011-agent-cost-observability.md).

## Why this exists

Without primitives like these, every new agent session pays the same costs over and over. Four failure modes recur:

- **Re-deriving ground state.** The next session re-reads the codebase to figure out "what's open?" — exactly the friction that makes "skills over recall" feel impossible. `STATE.md` replaces that with one file read.
- **Agent observation silently becoming instruction.** Without a review-state schema, an entry an agent added because it noticed something looks identical to one you confirmed. Future agents treat both as direction. Review + Provenance draw the line: agent-authored entries start as `unreviewed`, and only the operator promotes to `confirmed` or `rejected`.
- **Over-asking or over-acting.** Without a doctrine for what's autonomous vs what needs your turn, agents either ping you constantly (annoying) or quietly take actions you'd have wanted to see first (dangerous). The four-level autonomy doctrine pins it down: L1 read, L2 autonomous-on-branch, L3 propose-and-confirm, L4 human-only.
- **Cost drift going unnoticed.** Without per-feature cost attribution, you can't tell whether a workflow has degraded from 2 sessions to 5 over time — or whether one of your "skills over recall" intuitions is actually paying off. Cost signals make it testable.

The original boilerplate slices (1–6) gave projects ADRs, contracts, agent primitives, structured logs, and OpenTelemetry. Slices 7–10 made the *session loop itself* something agents and humans share durably. The neutral-router pattern (see [ADR-0012](docs/decisions/0012-agent-neutral-router-and-client-adapters.md)) extended that loop **across agent clients**, so the same doctrine drives Claude Code, Codex, a homegrown agent, or several of them working together — without forking the repo per provider.

The whole template is opinionated on one principle: **agent-written content starts as evidence, not instruction.** It's a position client-local operator memory should take by default ("transient task state — use the active agent client's task tracker, not memory"). This template extends the same discipline to project-local state, governance, and cost — so an agent can compound work across sessions without quietly authoring the rules it operates under.

## How this works with one or many agents

This template is designed for **any agent client** — Claude Code, Codex, a homegrown agent, or several working alongside each other — operating against one shared project contract. The framework draws a hard line between *shared doctrine* (in the repo, every agent reads it) and *client adapters* (thin shims that translate invocation mechanics for a specific client).

### The neutral router pattern

| Layer | Lives in | Audience |
| --- | --- | --- |
| **Shared doctrine** | `AGENTS.md`, `agent-instructions/`, `STATE.md`, ADRs, templates | Every agent and contributor |
| **Runtime contracts** | `templates/agent-primitives/`, `templates/contracts/` | Every agent acting across rings |
| **Client adapters** | `CLAUDE.md`, `.claude/skills/`, `.claude/settings.json` for Claude Code; analogous files for Codex or homegrown clients | One specific client |

`AGENTS.md` is the canonical entry point — every agent client reads it on session start, then fetches the relevant topic files from `agent-instructions/`. Adapters add invocation mechanics (slash commands, skill bundles, permission config, local-memory paths) but **must link to shared doctrine, not duplicate it**. Doctrine drift between clients is the failure mode the pattern exists to prevent. Full strategy in [`docs/agent-clients/README.md`](docs/agent-clients/README.md).

### Four operating modes

The template supports four ways agents work the repo:

- **Single-agent.** One client opens the repo, reads `AGENTS.md`, follows `agent-instructions/`, uses its own adapter workflows where available. The default case.
- **Alternating-agent.** Multiple clients (e.g., Claude Code one day, Codex the next) work the same repo over time. `STATE.md` is the handoff point; agent-authored entries default to `Review: unreviewed` regardless of which client wrote them, so the operator stays in control of cross-client promotion.
- **Hybrid.** Multiple agents in parallel for different strengths — one for implementation, one for review, one for research. Agents exchange durable facts only through shared artifacts: PRs, ADRs, `STATE.md`, contracts, tests. Private conversation context does *not* cross clients.
- **Homegrown.** A custom agent implements the same contract: read `AGENTS.md` first, treat `STATE.md` as project-local handoff, mark agent-authored observations as `Review: unreviewed`, follow the L1–L4 autonomy doctrine, emit typed Task / Result / Evidence / Provenance when delegating across rings.

### How agents coordinate when more than one is in play

Three primitives carry the coordination load:

- **`STATE.md`** is the shared ledger. Any agent reads it on session start; any agent appends to it; only the operator transitions `Review: unreviewed → confirmed | rejected`. The next agent, in this session or another, reads the same ground truth.
- **`Provenance`** records *which* agent and which session authored each entry — so when something later needs to be traced back, the origin is durable even after the conversation that produced it is long gone.
- **Conflict resolution is operator-mediated.** Tests, contracts, and cited evidence outrank model confidence. When agents disagree on a doctrine-level call, the operator decides; the decision is recorded as an ADR or a confirmed `STATE.md` entry. The full flow lives in [`docs/agent-clients/README.md`](docs/agent-clients/README.md#conflict-resolution).

The promise of the pattern is simple: **one repo, one doctrine, many agents**. No fork-per-provider. No silent doctrine drift. The operator stays the resolver of last resort.

## Who this is for

You're starting a new project — especially one involving AI agents, LLM-driven features, or just a codebase where you want Claude Code, Codex, or another agent to be a useful collaborator from day one — and you'd rather not reinvent the same dev/doc/security practices every time.

The session-handoff, review-state, autonomy, and cost-signal primitives assume work will span many sessions and possibly many contributors (or many agents) over time. If your project genuinely lives in a single session by a single person, those four pieces are overhead you can ignore — the rest of the template still applies.

You're **not** looking for a runtime scaffold (no Go / Next.js / Python project structure here). Pair this with whatever language scaffold you normally use.

## Use as a template

1. Click **"Use this template"** on GitHub (or `gh repo create <name> --template=...`), then clone your new repo.
2. Open the clone in your agent client. Claude Code users can run the bundled bootstrap skill; other clients can follow the same steps manually.
3. For Claude Code, run the bootstrap skill:

   ```
   /bootstrap-project
   ```

   Other clients can follow the same steps manually. The bootstrap workflow will:
   - Ask for project name, description, stack, frontend yes/no.
   - Replace `{{PLACEHOLDERS}}` across the tree.
   - Offer to seed optional Claude Code memory for this project.
   - Offer to `gh label sync` the label taxonomy.
   - Offer to scaffold a Playwright `e2e/` directory if you have a frontend.
   - Offer to kick off your project's first ADR.
   - Remove itself when done.
4. Commit and push. You're set up.

If you prefer to do it manually: grep for `{{` and replace the placeholders, then optionally copy `templates/memory/` into Claude Code's `~/.claude/projects/<encoded-path>/memory/` using the instructions in [`templates/memory/README.md`](templates/memory/README.md). Also copy `templates/state/STATE.md` to the project root and replace `{{PROJECT_NAME}}` — that's the session-handoff ledger described below.

## A typical session

The boilerplate is built for projects worked across many sessions and contributors. Each non-trivial session follows the same loop:

1. **Start-session workflow** — the agent reads `STATE.md` before any other repo exploration. It surfaces blocked items, *unreviewed* agent-authored items (entries the agent added since you last looked), in-progress work, and open questions. You direct what's next instead of asking the agent to re-derive ground state.
2. **Do the work** — per `AGENTS.md`'s house rules. The share-plan-before-editing rule governs medium and large tasks. The agent operates at **L2 autonomy** on a feature branch, escalating to **L3 (propose-and-confirm)** for anything touching `main`, settings, doctrine, dependencies, or CI; **L4** operations (force-push, production writes, external messages, infra) are human-only.
3. **End-session workflow** — the agent reconciles `STATE.md` against what actually happened: completed items move to recent with a one-sentence summary, last-touched dates update, new items get added with `Review` set per author defaults (`confirmed` when operator-directed, `unreviewed` when agent-noticed), tabled items are marked with reason and un-table condition, open questions are recorded. `Cost signals` (sessions to date, operator turns, skills used) are captured on multi-session entries. `STATE.md` is committed separately from code.

Architectural decisions that crystallize get an ADR (`/new-adr`); decisions still in flux stay in `STATE.md` and graduate later.

### How the four primitives compose

| Primitive | Carries | Read at | Updated at |
| --- | --- | --- | --- |
| `STATE.md` (slice 7) | Open / completed / tabled / open-questions | start-session workflow | end-session workflow |
| Review + Provenance (slice 8) | Evidence-vs-instruction line, agent-side origin | every read | when new entries are added or operator reviews |
| Autonomy levels (slice 9) | What the agent does without asking | every action | rarely (doctrine change → ADR) |
| Cost signals (slice 10) | Per-entry cost over its lifetime | start-session workflow (visible on items) | end-session workflow |

Every ledger entry carries its origin (Provenance), authority (Review), and the cost it has consumed (Cost signals). The autonomy doctrine governs what the agent can do *to* the ledger and the rest of the repo. The session-handoff skills are the routine that keeps the whole thing accurate.

## How operator memory works in your projects

This template assumes two distinct stores of cross-session context, and ships templates for both. They answer different questions:

| Store | Answers | Lives in | Audience |
| --- | --- | --- | --- |
| **Client-local operator memory** | "Who am I working with, and how do they like to work?" | Client-specific local store, such as Claude Code's `~/.claude/projects/<encoded-path>/memory/` | This operator only |
| **Project-local `STATE.md`** | "Where is the project right now?" | Repo root (committed) | Every contributor and agent |

### Client-local operator memory

Some agent clients maintain a per-project memory store on your machine, separate from the repo. It is where the agent remembers things about **you** across sessions — your role, your stack fluency, your communication preferences, the corrections and confirmations from past sessions, time-bound project facts (deadlines, stakeholders).

Without it, every conversation opens with the same orientation tax: "I'm a frontend engineer doing solutions-architect work, prefer concise responses, no acronyms…" With it, the agent has continuity of relationship — and the time you spend re-orienting the model becomes time spent on the work.

The boilerplate ships templates for four memory types:

| Type | Holds | Example |
| --- | --- | --- |
| **`user`** | Your role, what you do, what you already know. Slow-moving. | "Self-taught frontend engineer, 8 years as a Solutions Architect for Customer Data Platforms; new to deployment ops." |
| **`feedback`** | Corrections and confirmations from past sessions. | "Don't mock the database in integration tests — got burned by divergence before." |
| **`project`** | Time-bound project facts. Decays as the project moves. | "Release cut scheduled 2026-05-12; only bug fixes after that date." |
| **`reference`** | Pointers to external systems. | "Pipeline bugs tracked in Linear project INGEST." |

Each memory file is a markdown document with YAML frontmatter naming the `type` and a one-line description. A `MEMORY.md` index keeps the directory scannable on every load. The schema and naming conventions are **client-agnostic** — Claude Code reads them from `~/.claude/projects/<path>/memory/`, but a Codex adapter or a homegrown client can point at the same files from whatever local-memory path its client uses. Full conventions live in [`agent-instructions/operator-memory.md`](agent-instructions/operator-memory.md).

The Claude `/bootstrap-project` skill offers to seed the templates into the right Claude Code path on first run. Other clients can copy them manually using their own local-memory conventions — see [`templates/memory/README.md`](templates/memory/README.md).

### How memory composes with STATE.md

The two stores carry orthogonal information and have different lifecycles:

- **Client-local memory** is *biographical*. It travels with you across all projects you work on, slowly accumulating preferences and reusable context. It is not shared with teammates and not part of the repo.
- **`STATE.md`** is *operational*. It lives in the repo, is the same for every contributor and agent, and changes every session.

Together, when a new session opens, the agent knows both **who you are** (memory) and **where the project is** (STATE.md). That single change — replacing the first 5 minutes of every session with a single orientation read — is the entire payoff loop this template is designed around.

**The boundary matters.** Putting operational state in client-local memory breaks cross-contributor legibility (only you see it). Putting biographical context in `STATE.md` leaks operator preferences into the committed repo. The Memory-vs-project-state comparison table in [`agent-instructions/operator-memory.md`](agent-instructions/operator-memory.md) is the normative reference.

### What this looks like in practice

After bootstrap, a client-local memory directory may hold files like:

```
~/.claude/projects/-Users-you-projects-your-app/memory/
  MEMORY.md                  # index — always loaded into context
  user_role.md               # who you are
  feedback_pr_style.md       # "PRs stay small, single-purpose"
  project_release_q2.md      # "Q2 release cut 2026-06-30; bug-only after"
  reference_grafana.md       # "Latency dashboard: grafana.internal/d/api-latency"
```

The repo root holds `STATE.md` with sections for open work, recently completed, tabled, and open questions. Both stores get richer over time. Project-local state is shared across every agent and contributor; client-local memory remains private to that operator and that client (a Claude Code memory file is not visible to Codex, and vice versa). Neither requires you to type "remember that…" — the agent and the operating doctrine do that work for you.

## File map

```
.
├── AGENTS.md                           # neutral router + house rules
├── CLAUDE.md                           # Claude Code adapter shim
├── agent-instructions/                 # shared topic files loaded on-demand
│   ├── development-workflow.md         # TDD loop, branching, PR flow
│   ├── testing-practices.md            # test pyramid, Playwright e2e
│   ├── security-practices.md           # OWASP Web / LLM / Agentic
│   ├── documentation-discipline.md     # doc routing rules
│   ├── commit-conventions.md           # Conventional Commits
│   ├── dependency-vetting.md           # 5-point dep check
│   ├── component-explainability.md     # README at every boundary
│   ├── operator-memory.md              # client-local memory conventions
│   ├── session-handoff.md              # STATE.md doctrine + Review/Provenance
│   ├── agent-autonomy.md               # L1-L4 autonomy doctrine
│   ├── agent-cost-observability.md     # Cost signals on STATE.md entries
│   ├── epics-and-projects.md           # doc-embedded vs Projects v2
│   └── scaling-agent-instructions.md   # growing the instruction system
├── docs/
│   ├── agent-clients/                  # multi-agent adapter strategy
│   └── decisions/                      # ADRs + template + index
├── .claude/
│   ├── settings.json                   # Claude Code permissions
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
│   ├── memory/                         # seed files for Claude Code memory
│   └── state/                          # STATE.md seed (project-local session-handoff ledger)
├── .pre-commit-config.yaml
├── .editorconfig
├── .gitignore
├── .markdownlint.yaml
├── CONTRIBUTING.md
└── LICENSE                              # MIT
```

## Growing the template

`AGENTS.md` starts as a slim router to `agent-instructions/`. As the project grows into multiple modules, each module gets its own `AGENTS.md` + `agent-instructions/` subdirectory. Client-specific adapters should stay thin and point back to the neutral router. See [`agent-instructions/scaling-agent-instructions.md`](agent-instructions/scaling-agent-instructions.md) for the Stage 1 → 2 → 3 pattern.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Short version: keep it generic, ship docs with code, Conventional Commits, TDD where code changes.

## License

MIT — see [LICENSE](LICENSE).
