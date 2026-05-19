# STATE — {{PROJECT_NAME}}

> Project-local session-handoff ledger. **Not** user-local memory (which lives
> in client-local memory, such as a client-specific per-project memory directory). **Not** an Architectural
> Decision Record (which lives in `docs/decisions/`). This file is the running
> record of what's open, in-progress, blocked, tabled, or recently completed —
> so the next session can ground itself in seconds without re-reading the repo.
>
> Schema and lifecycle live in `agent-instructions/session-handoff.md` after
> this template is copied to the project root.
> Run the active client's start-session workflow at the top of each work
> session and the end-session workflow before wrapping up.

---

## Open work items

<!--
One entry per open work item. Schema:

### <slug-id>
- **Summary:** <one line>
- **State:** pending | in-progress | blocked
- **Review:** unreviewed | confirmed | rejected | stale | superseded
- **Opened by:** <human name or `agent`>
- **Opened at:** YYYY-MM-DD
- **Last touched:** YYYY-MM-DD
- **Next step:** <one line — what unblocks progress>
- **Provenance:** (optional but encouraged for agent-authored entries)
  - **Workflow:** <workflow that produced this entry, if any>
  - **Session:** <session id or YYYY-MM-DD-<seq>>
  - **Prompt summary:** <one line — what the operator asked when this came up>
  - **Context:** <file:line refs or related entries>
- **Cost signals:** (optional, recorded by the end-session workflow on multi-session entries)
  - **Sessions to date:** <integer>
  - **Skills used (cumulative):** <list>
  - **Operator turns:** <integer>
  - **Context read (approx):** <files-touched count>, <lines if known>
  - **Notes:** <free-form — e.g., absolute token counts from external tooling>
- **Notes:** <optional; keep short>

Review defaults: `confirmed` for human-authored entries, `unreviewed` for
agent-authored entries (entries the agent added without explicit operator
direction). The operator alone transitions Review to/from confirmed and
rejected. Agents may mark their own entries `superseded` (pointing at the
replacement) or update Provenance.

Remove the entry once it lands in "Recently completed".
-->

<!-- (none yet — add items as work begins) -->

## Recently completed

<!--
Rolling window — keep ~10 most recent or last 14 days, whichever is shorter.
Older entries age out (git history retains everything). Schema:

### <slug-id>
- **Summary:** <one line>
- **Completed at:** YYYY-MM-DD
- **Session summary:** <one sentence — what changed>
- **Commit / PR:** <link if useful>
- **Review at completion:** confirmed (by operator) | self-completed (agent, under operator direction)
- **Cost signals:** (optional, final totals at completion — see schema above)
-->

<!-- (none yet) -->

## Tabled

<!--
Work that's been deliberately parked. Don't just delete — the tabled entry is
the record that the trade-off was considered. Schema:

### <slug-id>
- **Summary:** <one line>
- **Tabled at:** YYYY-MM-DD
- **Why tabled:** <one sentence — the reason>
- **Un-table when:** <one sentence — what changes to bring it back>
- **Review:** stale | rejected | tabled-only
-->

<!-- (none yet) -->

## Open questions

<!--
Questions waiting on an answer (typically from the human operator). Resolves
into an open work item or gets removed once answered. Schema:

### <slug-id>
- **Question:** <one or two sentences>
- **Asked at:** YYYY-MM-DD
- **Asked by:** <agent or human>
- **Waiting on:** <who or what>
- **Provenance:** (optional)
  - **Workflow:** <workflow that produced the question>
  - **Session:** <session id>
  - **Context:** <file:line refs>
-->

<!-- (none yet) -->
