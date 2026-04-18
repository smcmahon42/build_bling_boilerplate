# Security testing plan

> Copy this template to `testing-plan.md` and fill in the project specifics. Delete this intro paragraph.

A phased plan for verifying the project's security posture. Phases are cumulative — you don't skip ahead; Phase 3 assumes Phase 1 and 2 are in place.

Each phase states:

- **What** — the tests or reviews performed.
- **Trigger** — when it runs.
- **Owner** — who is responsible.
- **Artifact** — where the result is recorded.

## Phase 1 — Static / automated baseline

Continuous verification from code. Fast, cheap, runs on every change.

| Check | Tool / workflow | Trigger | Owner | Artifact |
| --- | --- | --- | --- | --- |
| Secret scanning | `gitleaks` (pre-commit + `.github/workflows/pre-commit.yml`) | Every commit / PR | — | GH Actions log |
| Dependency vulnerabilities | `.github/workflows/security-scan.yml` (govulncheck, npm audit, pip-audit, cargo audit) | Every PR, weekly | — | GH Actions log + `/dep-audit` report |
| SAST | `.github/workflows/codeql.yml` | Every PR, weekly | — | CodeQL alerts in Security tab |
| License compliance | `.github/workflows/license-check.yml` | Every PR | — | GH Actions log |
| Agent/diff review | `/security-review` skill | Security-relevant PRs | Author + reviewer | PR comment |
| OWASP walk (30 items) | `/security-review` | Before each release | — | `remediation-plan.md` |

Phase 1 is **non-negotiable.** These checks run before anything in later phases is planned.

## Phase 2 — Dynamic application testing (DAST)

Runs against a deployed instance. Catches what static analysis can't.

| Check | Tool / workflow | Trigger | Owner | Artifact |
| --- | --- | --- | --- | --- |
| OWASP ZAP API scan | `.github/workflows/dast-scan.yml` | After deploy to preview/staging | — | GH issue auto-filed; `remediation-plan.md` |
| Auth flow fuzzing | Manual or scripted; Burp / OWASP Amass / custom | Quarterly or on auth change | — | Report in `docs/security/audit/` |
| API contract testing | Schemathesis / Dredd / similar | CI (if OpenAPI spec exists) | — | GH Actions log |

Enable Phase 2 once a preview or staging URL exists. Update `.github/workflows/dast-scan.yml` triggers accordingly.

## Phase 3 — Manual penetration testing

Human adversarial review. Expensive, high-signal.

| Check | Cadence | Owner | Artifact |
| --- | --- | --- | --- |
| Internal pentest | Before major releases | Security engineer / rotating reviewer | Report in `docs/security/audit/YYYY-MM-internal.md` |
| External pentest | Annually (or before compliance milestones) | Contracted firm | Report in `docs/security/audit/YYYY-MM-external.md` |
| Threat model refresh | When architecture changes materially | Owner of `threat-model.md` | Updated threat model |

Pentest scope should include:

- Abuse cases from the threat model.
- Recently-added endpoints / tools / agent surfaces.
- Known-weak areas (new auth flows, new data classes handled).
- All three OWASP lists — Web Top 10, LLM Top 10, Agentic Top 10.

## Phase 4 — Infrastructure hardening

Configuration and runtime posture, not code.

| Check | Trigger | Owner | Artifact |
| --- | --- | --- | --- |
| TLS configuration | On deploy change | Platform | Qualys SSL Labs or internal scan log |
| Secret rotation | Documented cadence (e.g. 90-day) | Secrets owner | Audit log |
| IAM review | Quarterly | Platform | Review notes in `docs/security/audit/` |
| Runtime sandboxing (agent tools) | On tool change | Agent owner | ADR reference |
| Logging / monitoring coverage | On new surface | Owner of affected module | Runbook |

## Out-of-scope

<!-- List what this plan explicitly does not cover.
     Example: physical security, employee offboarding, vendor security reviews. -->

## Review cadence

- Phase 1 is continuous — monitor results, don't re-plan it.
- Phase 2 plan reviewed quarterly.
- Phase 3 scoped before each pentest.
- Phase 4 reviewed quarterly or on infrastructure change.

## Related

- `remediation-plan.md` (sibling file; copy from `remediation-plan-template.md`) — where findings go.
- [`claude-instructions/security-practices.md`](../../claude-instructions/security-practices.md) — the OWASP reference.
- [`SECURITY.md`](../../SECURITY.md) — disclosure policy.
