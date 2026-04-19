# Skills

User-invocable skills live here. Each skill is a markdown file with frontmatter (`name`, `description`) plus a body that defines the workflow.

## Index

The machine-readable index lives in [`skills-index.yaml`](skills-index.yaml) — tag-driven mapping from task types to skills. Use it as the source of truth when deciding which skill to load. The table below is a human-readable summary.

| Skill | Invoke with | Use when |
| --- | --- | --- |
| [bootstrap-project](bootstrap-project.md) | `/bootstrap-project` | First-time setup after cloning the template. Runs once per project. |
| [security-review](security-review.md) | `/security-review` | Any change touching auth, data flow, tool use, memory, logs, or external integrations. Walks all 30 OWASP items (Web + LLM + Agentic). |
| [dep-audit](dep-audit.md) | `/dep-audit` | Before adding dependencies, before releases, or weekly. Language-detecting vulnerability scanners + 5-point check on new deps. |
| [test-gaps](test-gaps.md) | `/test-gaps` | Before a release, when coverage drops, or when you suspect critical-path tests are missing. |
| [new-adr](new-adr.md) | `/new-adr` | Recording an architectural decision (new boundary, new tech, changed auth/data/observability surface). |
| [new-skill](new-skill.md) | `/new-skill` | Authoring a new skill to add to this directory. |
| [new-contract](new-contract.md) | `/new-contract` | Adding a new contract (OpenAPI / JSON Schema / protobuf). Picks the form, copies the starter, registers and wires codegen. |

## How to pick a skill

1. Classify the task scale (see `claude-instructions/development-workflow.md#task-scale-matrix`).
2. Consult `skills-index.yaml` — match the task's tags against skill `triggers`.
3. Load the matching skill's markdown file and follow its workflow.
4. If no skill matches but the task recurs, consider `/new-skill`.

## How to author a new skill

Use `/new-skill`. It walks you through the template, review criteria, and index update. Manual path: copy an existing skill file as a starting point, fill in the frontmatter and body, update `skills-index.yaml`, and run through [`references/skill-review-criteria.md`](references/skill-review-criteria.md) before merging.

## Related

- `references/skill-review-criteria.md` — quality bar for skills.
- `.claude/agents/README.md` — when to prefer an agent over a skill.
