# STATE — {{PROJECT_NAME}}

> Project-local session-handoff ledger. **Not** user-local memory (which lives
> at `~/.claude/projects/<encoded-path>/memory/`). **Not** an Architectural
> Decision Record (which lives in `docs/decisions/`). This file is the running
> record of what's open, in-progress, blocked, tabled, or recently completed —
> so the next session can ground itself in seconds without re-reading the repo.
>
> Schema and lifecycle live in
> [`claude-instructions/session-handoff.md`](../../claude-instructions/session-handoff.md).
> Run `/start-session` at the top of each work session and `/end-session`
> before wrapping up.

---

## Open work items

<!--
One entry per open work item. Schema:

### <slug-id>
- **Summary:** <one line>
- **State:** pending | in-progress | blocked
- **Opened by:** <human or agent name>
- **Opened at:** YYYY-MM-DD
- **Last touched:** YYYY-MM-DD
- **Next step:** <one line — what unblocks progress>
- **Notes:** <optional; keep short>

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
-->

<!-- (none yet) -->
