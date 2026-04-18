---
name: security-review
description: Walk the current branch's changes against OWASP Web Top 10, OWASP LLM Top 10, and OWASP Agentic Top 10. Produce a prioritized findings table covering all 30 categories, with file references and suggested fixes. Cross-check against CI coverage (is security-scan wired up? is codeql on?).
---

# Security review

Run this skill before shipping any change that touches authentication, authorization, data flow, tool use, memory, logs, or external integrations. It covers the three OWASP lists used across this project — see `claude-instructions/security-practices.md` for the full reference.

## What to review

**Default scope:** the diff between the current branch and the project's default branch (`main` or `master`). Use:

```bash
git diff origin/main...HEAD --name-only
git diff origin/main...HEAD
```

If the operator names a different base, use that.

## Categories to check (all 30)

### OWASP Web Top 10 (2021)

- **A01 Broken access control** — is every new/changed endpoint gated by the right permission check? Look for missing authorization middleware, object-level access checks, default-allow policies.
- **A02 Cryptographic failures** — any hand-rolled crypto? Any plaintext secret storage? TLS required on every new endpoint?
- **A03 Injection** — string concatenation into SQL, shell, HTML, LDAP, XPath? Template rendering of untrusted input without escaping?
- **A04 Insecure design** — is there a threat model? Are abuse cases listed?
- **A05 Security misconfiguration** — sample accounts in seed data? Overly permissive CORS? Missing security headers?
- **A06 Vulnerable components** — any new deps without a `/dep-audit` run?
- **A07 Authentication failures** — rate limits on login? Session rotation on privilege change? MFA available where it should be?
- **A08 Software & data integrity** — any unsigned artifact pulls? Supply-chain risks in new deps?
- **A09 Logging & monitoring failures** — auth events logged? Admin actions audited?
- **A10 SSRF** — any new outbound URL construction from user input? Is private IP range blocked?

### OWASP LLM Top 10

- **LLM01 Prompt injection** — does the prompt concatenate untrusted input? Are tool calls gated by authorization?
- **LLM02 Insecure output handling** — is the LLM output rendered as HTML, passed to `eval`, used in shell commands, or inserted into SQL?
- **LLM03 Training data poisoning** — for fine-tunes or RAG corpora: provenance verified?
- **LLM04 Model DoS** — input size caps, max tokens, timeouts, rate limits?
- **LLM05 Supply chain** — model provider and version pinned?
- **LLM06 Sensitive info disclosure** — is PII/secret-redaction applied before the model sees the payload?
- **LLM07 Insecure plugin/tool design** — are tool schemas narrow? Do tools validate arguments?
- **LLM08 Excessive agency** — does the agent confirm before destructive writes, external messages, or large financial actions?
- **LLM09 Overreliance** — are confidence and provenance surfaced? Safe fallback when uncertain?
- **LLM10 Model theft** — rate-limit, watermark, anti-scraping?

### OWASP Agentic AI Top 10

- **ASI01 Agent identity** — distinct identity per agent? Auditable?
- **ASI02 Memory poisoning** — memory writes sanitized and signed? Short-term scratch separated from long-term memory?
- **ASI03 Tool misuse** — narrow schemas, argument validation, sandboxed execution?
- **ASI04 Privilege escalation** — least privilege on agent credentials?
- **ASI05 Insecure orchestration** — plan/act loops detect cycles, budget caps, runaway cost?
- **ASI06 Reasoning manipulation** — plan monitoring? Confirmation on deviation?
- **ASI07 Cross-agent collusion / cascade** — blast radius isolated? Circuit breakers?
- **ASI08 Unauthorized external actions** — outbound allowlist? Confirmation on external writes?
- **ASI09 Observability gaps** — structured traces per tool call/prompt/decision? Redaction at storage?
- **ASI10 Agent supply chain** — prompt/tool/sub-agent provenance verified? Versions pinned?

## Output format

Produce a table:

| ID | Category | Severity | File:line | Finding | Suggested fix |
| --- | --- | --- | --- | --- | --- |

Severity scale: **critical** (fix before merge) / **high** (fix before release) / **medium** (open an issue) / **low** (consider) / **n/a** (not applicable to this diff).

After the table, include:

- **Summary counts** by severity.
- **CI coverage audit:** confirm `security-scan.yml`, `codeql.yml`, `license-check.yml`, and `pre-commit.yml` exist and run on PRs. Call out any missing workflow that should have caught a finding.
- **Recommended follow-ups:** any items that should be tracked as issues (use `security` + `priority:p*` labels).

## Invariants

- Don't invent vulnerabilities that don't exist — "n/a" is an acceptable row.
- If a finding is theoretical but the diff doesn't actually introduce it, still document it under "n/a" or omit.
- Reference OWASP docs by ID, not by rephrased summary. The ID is stable; summaries drift.
- If a category doesn't apply (e.g., LLM01 in a project with no LLM), mark the whole row "n/a" and move on.

## Related

- `claude-instructions/security-practices.md` — full reference for the three lists.
- `/dep-audit` — companion skill for component-level vulnerability scanning.
