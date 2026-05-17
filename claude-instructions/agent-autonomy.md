# Agent autonomy scope

Not every action an agent takes carries the same risk. Reading a file is
cheap and reversible. Pushing to main is irreversible. Without an explicit
doctrine, agents drift toward "ask about everything" (which is annoying and
slow) or "do everything" (which is dangerous). This file draws the lines.

The doctrine has four levels. Each operation an agent can perform falls
into exactly one level. The level determines whether the agent acts
autonomously, proposes-and-confirms, or refuses without explicit operator
direction.

## The four levels

| Level | Operator role | Examples |
| --- | --- | --- |
| **L1 — Autonomous** | None. Agent acts. | Read files, run linters, `git status`/`diff`/`log`, `gh pr view` |
| **L2 — Autonomous on feature branch; propose-first on medium/large tasks** | Plan-review before edits start; otherwise none. | Edit code on a feature branch, commit, run `/end-session` |
| **L3 — Propose-and-confirm** | Reviews the proposed action and explicitly approves before agent proceeds. | Merge to main, push, modify CI / settings.json, modify CLAUDE.md house rules, accept ADRs |
| **L4 — Human-only** | Agent never performs; refers the operation back. | Push to main without PR, force-push to main, repo settings, external messages, financial / production / infra actions |

The split between L2 and L3 maps onto the existing task-scale matrix
(see `development-workflow.md#task-scale-matrix`) and house rule 3 of
`CLAUDE.md`: medium and large tasks share the plan before editing. Small
tasks are L2-autonomous; medium/large are L2-propose-then-act on the plan.

## L1 — Autonomous (read & inspect)

The agent performs these without asking. They are reversible and cheap.

- Read any file in the repo.
- `git status`, `git diff`, `git log`, `git show`, `git branch` (no
  checkout), `git remote -v`.
- `gh pr view`, `gh issue view`, `gh repo view`, `gh pr list`,
  `gh issue list`.
- `ls`, `find`, `grep`, `pwd`, `which` — read-only inspection.
- Run linters, type checkers, test runners (deterministic, no side
  effects).
- Read from user-local Claude memory at
  `~/.claude/projects/<path>/memory/`.
- Run `pre-commit run` (no-op on clean trees; only shows what would
  change).

These are the operations on the `.claude/settings.json` allowlist today.
The allowlist *is* the enforcement layer for this level.

## L2 — Autonomous on feature branch; plan-first on medium/large

The agent edits, scaffolds, and commits on a feature branch. House rule 3
applies: medium and large tasks share the plan (what files, what shape,
what tests) before editing starts. Small tasks proceed directly.

- Edit code files on a feature branch.
- Add new files in `templates/`, `claude-instructions/`, `.claude/skills/`,
  `docs/`, etc.
- Run skills that include their own "ask first" gates (`/new-adr`,
  `/new-skill`, `/new-contract`).
- Make Conventional Commits on the feature branch.
- Update STATE.md via `/end-session`. Agent-authored entries default to
  `Review: unreviewed`; operator-directed ones to `Review: confirmed`
  (per `session-handoff.md`).
- Run pre-commit hooks that may modify formatting / fix lint.

The branch is the containment boundary. Nothing the agent does at L2
becomes visible to main until L3 (merge) succeeds.

## L3 — Propose and confirm

The agent **prepares** the action and **shows it** to the operator. The
operator explicitly approves before the action is performed. "Implicit
approval by silence" is not allowed.

- **Merge a feature branch into main.** Show the merge plan; confirm.
- **Push to origin** (any branch).
- **Modify `.claude/settings.json`.** Permissions affect the agent's own
  capabilities — never silent edits.
- **Modify `CLAUDE.md` house rules or process checkpoints.** Doctrine
  changes require explicit confirmation.
- **Flip an ADR from Proposed to Accepted.** Acceptance is operator's act.
- **Modify or supersede an existing Accepted ADR.** Write a superseding
  ADR; do not edit the original.
- **Modify CI workflows** (`.github/workflows/`).
- **Add or change a project dependency.** Run `/dep-audit` as part of the
  proposal.
- **Modify `.gitignore`.** Adding entries can mask real problems.
- **Modify the project license or `CODE_OF_CONDUCT.md`** — these are legal
  and community artifacts; never silent edits.
- **Modify the skill review criteria** in `.claude/skills/references/`.

The proposal pattern: agent presents (a) the change as a unified diff or
file write preview, (b) the rationale, (c) the alternatives considered.
Operator approves, modifies, or rejects.

## L4 — Human-only

The agent never performs these. If asked, the agent declines, explains
why this is L4, and asks the operator to perform the action themselves.

- **Push to main directly without a PR / merge commit.** Even if the local
  history is clean.
- **Force-push to main or any protected branch.** Loss of history.
- **`git rebase --interactive` on shared history** that has already been
  pushed.
- **Delete remote branches** (`git push --delete`, `gh repo delete-branch`).
- **Modify repo settings** (visibility, access, protected branches,
  webhooks, secrets).
- **Send external messages** — Slack, email, `gh issue create` or
  `gh pr comment` to other users, anything fan-out.
- **Financial / commerce actions** — any monetary write, payment, refund,
  subscription change.
- **Production data writes** — direct production database mutations,
  destructive migrations, schema drops, `TRUNCATE`, unqualified `DELETE`.
- **Cloud infrastructure changes** — `terraform apply`, `gcloud` mutations,
  `aws` writes, deploys to production environments.
- **Modify the project license.**

These are the actions where the blast radius justifies a hard stop. The
agent's job at L4 is to *refer*: "this is L4 per `agent-autonomy.md`;
please perform it yourself or grant a one-shot exception with explicit
direction."

## One-shot exceptions

The operator can grant a one-shot exception to any level — "go ahead and
merge this one" or "yes, push directly this time." The agent honors the
exception **only for the specific action stated**, not as a permanent
upgrade. A subsequent action of the same type re-enters the doctrine.

One-shot exceptions are recorded in the relevant STATE.md entry's `Notes`
field for traceability.

## Why these particular lines

- **L1/L2 boundary** is the *visibility* line — can other contributors see
  what the agent did before the operator has? Feature branches contain;
  main does not.
- **L2/L3 boundary** is the *doctrine* line — does this action change how
  the project itself decides what to do? Settings, house rules, ADR
  acceptance, dependencies all shape future decisions; they get a
  proposal step.
- **L3/L4 boundary** is the *blast-radius* line — can the action damage
  other people, other systems, or unrecoverable state? Force-pushes,
  production writes, external messages all qualify.

If a new operation type doesn't fit cleanly, default to the higher level.
Promoting an action from L4 to L3 (or L3 to L2) requires an ADR; the
reverse only requires updating this file.

## Related

- `CLAUDE.md` house rule 3 — process checkpoints, including the
  share-plan-before-editing requirement that L2 inherits.
- `claude-instructions/development-workflow.md` — the task-scale matrix
  that distinguishes small from medium/large.
- `claude-instructions/session-handoff.md` — STATE.md authorship and
  Review rules that interact with L2 mutations.
- `.claude/settings.json` — the L1 enforcement allowlist.
- `docs/decisions/0010-agent-autonomy-scope.md` — the architectural
  decision proposing this doctrine.
