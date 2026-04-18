# Security remediation plan

> Copy this template to `remediation-plan.md` when you have your first finding to track. Delete this intro paragraph.

Living document. Active security findings with status, owner, and verification. Organized by severity, not by source — a finding from `/security-review`, CI, DAST, or an external audit is tracked the same way.

## How this document is used

- **Every finding** from `/security-review`, CI scans, DAST, external audits, or bug bounty lands here.
- **Rows are not deleted.** When a finding is fixed or explicitly accepted, the status changes but the history stays.
- **Review cadence:** weekly triage; monthly review of accepted risks.

## Status values

| Status | Meaning |
| --- | --- |
| Open | Confirmed, not yet being worked |
| In progress | Actively being fixed |
| Fixed — unverified | Fix merged, verification pending |
| Fixed — verified | Fix merged and verified via test / scan / retest |
| Accepted risk | Explicitly accepted with ADR; see `docs/decisions/` |
| Duplicate | Merged into another finding |
| Invalid | Not reproducible or not applicable |

## Severity ladder

| Severity | Response SLO |
| --- | --- |
| Critical | Fix before next release; active mitigation if exploitable |
| High | Fix within 30 days |
| Medium | Fix within 90 days |
| Low | Fix when convenient or next time the area is touched |

---

## Critical

| ID | OWASP | Finding | Discovered | Owner | Status | Verification |
| --- | --- | --- | --- | --- | --- | --- |
| <!-- e.g. SEC-001 --> | <!-- e.g. A01 / LLM08 / ASI03 --> | <!-- one line --> | <!-- YYYY-MM-DD --> | <!-- @handle --> | Open | <!-- how we'll prove it's fixed --> |

## High

| ID | OWASP | Finding | Discovered | Owner | Status | Verification |
| --- | --- | --- | --- | --- | --- | --- |

## Medium

| ID | OWASP | Finding | Discovered | Owner | Status | Verification |
| --- | --- | --- | --- | --- | --- | --- |

## Low

| ID | OWASP | Finding | Discovered | Owner | Status | Verification |
| --- | --- | --- | --- | --- | --- | --- |

## Accepted risks

For each, link to the ADR that documents the decision.

| ID | Finding | Rationale | ADR | Review date |
| --- | --- | --- | --- | --- |

## Fixed (historical)

Move rows here once status is `Fixed — verified`. Keep them for audit trail.

| ID | OWASP | Finding | Fixed | Verified by |
| --- | --- | --- | --- | --- |
