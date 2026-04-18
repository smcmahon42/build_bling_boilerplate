# Security docs

Project-specific security artifacts live here. CI gates (in `.github/workflows/`) are the automated baseline; the docs in this directory are the *narrative* that accompanies them — threat model, remediation plan, testing plan.

## Contents

| File | Purpose |
| --- | --- |
| `threat-model.md` *(create when ready)* | Assets, adversaries, entry points, trust boundaries, abuse cases. Living doc; revisit when the architecture shifts. |
| `remediation-plan.md` *(copy from the template below)* | Active list of security findings with status, owner, and verification steps. |
| `testing-plan.md` *(copy from the template below)* | Phased security testing: DAST, pentesting, infrastructure hardening. |
| `audit/` | Dated audit reports (internal or third-party). One file per audit. |

Templates to copy when you start using each:

- [`remediation-plan-template.md`](remediation-plan-template.md)
- [`testing-plan-template.md`](testing-plan-template.md)

## When to update

- **Remediation plan:** whenever a security finding is accepted from `/security-review`, CI, DAST, or an external report. Update status when findings are fixed or accepted as risk.
- **Testing plan:** revisit quarterly or when the system changes materially (new attack surface, new data class, new integration).
- **Threat model:** revisit when architecture changes (new service, new auth boundary, new external integration, new data class handled).

## Related

- [`claude-instructions/security-practices.md`](../../claude-instructions/security-practices.md) — the three OWASP lists.
- [`/security-review`](../../.claude/skills/security-review.md) — run against the current diff.
- [`/dep-audit`](../../.claude/skills/dep-audit.md) — language-detecting vuln scans.
- [`SECURITY.md`](../../SECURITY.md) — disclosure policy.
