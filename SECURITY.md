# Security policy

Thank you for helping keep this project and its users safe. This document covers how to report a vulnerability and what to expect from us in return.

## Reporting a vulnerability

**Do not file a public GitHub issue for security reports.** Public disclosure before a fix is available puts users at risk.

Instead, use **GitHub's private vulnerability reporting**:

- Go to the repository's **Security** tab → **Advisories** → **Report a vulnerability**.
- Fill out the form with as much detail as you can provide.

If that channel isn't available to you, reach out to the maintainer listed in `CODEOWNERS` by GitHub mention in a draft PR marked private, or by whatever contact method is published in the repo README.

When you report, please include:

- The affected version or commit.
- A clear description of the vulnerability.
- Steps to reproduce (a minimal test case is ideal).
- The impact you believe it has (data exposure, privilege escalation, DoS, etc.).
- Any suggested mitigation.

## What happens next

- **Acknowledgement** within 3 business days.
- **Initial assessment** within 7 days — severity, affected versions, rough timeline.
- **Fix development** in a private branch.
- **Coordinated disclosure** — we'll agree with you on a disclosure date. Default is after a fix ships and users have had a reasonable window to update.
- **Credit** in the advisory and release notes, unless you prefer to remain anonymous.

## Scope

In scope:

- Source code in this repository.
- CI/CD workflows under `.github/workflows/`.
- Documented deployment patterns.

Out of scope:

- Third-party dependencies (report to their maintainers; we'll pick up the fix via dependency updates).
- Hosting infrastructure not maintained by this project.
- Social engineering of maintainers or users.
- Denial-of-service via resource exhaustion that requires authenticated abuse at implausible scale.

## Security posture

This project follows OWASP guidance across three lists: Web Top 10, LLM Top 10, and Agentic AI Top 10. Baseline controls include:

- Pre-commit secret scanning (`gitleaks`).
- CI vulnerability scanning per detected stack (`govulncheck`, `npm audit`, `pip-audit`, `cargo audit`).
- CodeQL SAST on every PR.
- SBOM generation on every push to main.
- License compliance checks (copyleft denied by default).
- Least-privilege workflow permissions.

See [`claude-instructions/security-practices.md`](claude-instructions/security-practices.md) for the full checklist.

## Supported versions

Which versions receive security fixes depends on the project's release practice. Document it here during `/bootstrap-project`. Typical policy:

| Version | Security fixes |
| --- | --- |
| Latest minor | Yes |
| Previous minor | Yes for 90 days after latest minor ships |
| Older | No |

## Related

- [`.github/workflows/security-scan.yml`](.github/workflows/security-scan.yml) — CI security baseline.
- [`.github/workflows/dast-scan.yml`](.github/workflows/dast-scan.yml) — dynamic application security testing.
- [`docs/security/`](docs/security/) — project-specific threat model, remediation plan, testing plan.
