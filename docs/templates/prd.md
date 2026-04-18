# PRD — <Feature name>

- **Author:** <name>
- **Status:** Draft  <!-- Draft | In review | Approved | Superseded | Shipped -->
- **Last updated:** YYYY-MM-DD
- **Tracking issue / epic:** #
- **Related ADRs:** <NNNN, or "none yet">
- **Design Doc:** <link, if one exists>

## Problem

Who has this problem, what happens to them today, and why the current state is unacceptable. One or two paragraphs. Cite data if available (user counts, support ticket volume, latency numbers) — otherwise call out that the data is anecdotal.

**Not**: a solution. Stay with the problem until the reader feels it.

## Goals

What "done" looks like from the user's or system's perspective. Each goal is measurable.

- [ ] <Goal, e.g. "A new user can sign up in under 30 seconds on mobile.">
- [ ] <Goal>
- [ ] <Goal>

## Non-goals

What is explicitly out of scope. Non-goals prevent scope creep more than goals do.

- <Non-goal>
- <Non-goal>

## Users and use cases

Who are the users? What are the main use cases this solves?

| User type | Use case | Frequency |
| --- | --- | --- |
| <role> | <they do X to accomplish Y> | <daily / weekly / rare> |

## Requirements

Functional requirements — what the system must do.

1. <Requirement>
2. <Requirement>
3. <Requirement>

Non-functional requirements — performance, reliability, security, compliance, UX.

- **Performance:** <e.g. "P95 response under 300ms at 1k rps.">
- **Reliability:** <e.g. "No data loss on single-node failure.">
- **Security:** <e.g. "All three OWASP lists apply; auth covered by existing SSO.">
- **Accessibility:** <e.g. "WCAG 2.1 AA for all new UI.">

## Success metrics

How we'll know this shipped the intended outcome.

| Metric | Baseline | Target | Measurement window |
| --- | --- | --- | --- |
| <metric> | <current> | <desired> | <e.g. "30 days post-launch"> |

## Risks and open questions

What could derail this, and what we don't yet know. Name owners for the unknowns.

- **Risk:** <what could go wrong> — *mitigation:* <how we plan to handle it>
- **Open question:** <unresolved issue> — *owner:* <name>, *by:* <date>

## Constraints

Fixed inputs that shape the design — timelines, budgets, dependencies, compliance.

- <Constraint>
- <Constraint>

## Rollout plan

How this ships to users. Feature flag? Staged rollout? Migration needed?

- **Launch gate:** <what must be true before turning on>
- **Flag:** <name, default>
- **Staged rollout:** <e.g. "5% → 25% → 100% over 2 weeks, backed by metric X">
- **Rollback:** <how to turn it off if something breaks>

## Related

- Tracking issue: #
- Design Doc: <link>
- ADRs: <list>
- Prior art / competitive reference: <links>
