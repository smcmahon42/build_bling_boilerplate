# Component explainability

Every meaningful boundary in the codebase ships with a `README.md` that answers four questions. A **boundary** is any unit a new developer would treat as a separate thing to understand: a package, a top-level directory, a React feature group, a service, a subsystem.

## The four questions

Each README answers, in this order:

1. **What is this?** One or two sentences. No jargon. A person new to the team should understand the role of this component.
2. **Why is it separate?** What makes it a distinct unit rather than part of its neighbor? What's the boundary buying us?
3. **How does it connect?** Who calls in (callers/consumers)? Who does it call out to (dependencies/downstream)? What contract does it expose?
4. **Who owns it?** A team, an individual, or a rotation. "The whole team" is a non-answer — name someone.

## Target shape

- 20–40 lines total.
- Plain English. No ASCII diagrams unless genuinely useful.
- Link to relevant ADRs, runbooks, and deeper docs instead of duplicating them.
- **No API reference.** The code (and generated docs) are the source of truth. The README is the *context*.

## Example skeleton

```markdown
# <component name>

## What

A short description of the component's role in the system.

## Why separate

What problem justifies keeping this isolated from <neighbor>. Often: a different
stability contract, a different stack, a different team, or a different rate of
change.

## Connections

- **Called by:** `<upstream consumers>`
- **Depends on:** `<downstream services, libs>`
- **Contract:** `<public interface summary>` — defined in `<file>`

## Owner

<name or team> — see `CODEOWNERS` for the authoritative mapping.

## See also

- ADR 0042: <decision that shaped this component>
- Runbook: <link>
```

## When to write one

- **On creation.** A new package/module/service doesn't merge without its README.
- **After a non-trivial change to its contract.** Update the README in the same PR.
- **When you find yourself answering the same question twice.** That answer belongs in the README.

## When to delete or merge

- If two components share the same answer to "Why is it separate?", they probably shouldn't be separate. Collapse them.
- If a component is removed, delete its README. Don't leave orphans.

## How this connects to the rest of the system

- The root `README.md` is the top-level component README for the repo itself.
- A module README is the Stage 2/3 equivalent for each module (see [scaling-claude-instructions.md](scaling-claude-instructions.md)).
- ADRs (`docs/decisions/`) capture the *decisions* that shaped the component; the README captures the *current state* and the *why*.

Together they form the answer to: "I'm new here, where do I start?" → read the repo README, then the relevant module READMEs, then the ADRs that look interesting.
