# Security practices

Security is a first-class concern. AI projects face the standard web attack surface **plus** a distinct set of LLM and agent-specific risks. The boilerplate tracks three OWASP lists and wires CI to enforce the baseline.

When in doubt, run `/security-review` — the skill walks all three lists against the current diff.

## The three OWASP lists

### OWASP Web Top 10 (2021)

| ID | Risk | What to check |
| --- | --- | --- |
| A01 | Broken access control | Authorization enforced on every endpoint; default-deny; object-level checks. |
| A02 | Cryptographic failures | TLS everywhere; no hand-rolled crypto; keys in a secret manager, not code. |
| A03 | Injection | Parameterized queries; no string concatenation into SQL/shell/HTML/LDAP. |
| A04 | Insecure design | Threat model the feature before building; list the abuse cases. |
| A05 | Security misconfiguration | Secure defaults; remove sample accounts; disable directory listing; strict CSP. |
| A06 | Vulnerable & outdated components | `security-scan` workflow + `/dep-audit` skill. |
| A07 | Identification & authentication failures | MFA option; rate-limit login; rotate sessions on privilege change. |
| A08 | Software & data integrity | Sign artifacts; verify supply chain; SBOM on every release. |
| A09 | Logging & monitoring failures | Log auth events, auth failures, admin actions; alert on suspicious patterns. |
| A10 | Server-side request forgery | Validate outbound URLs; block private IP ranges; use allowlists. |

### OWASP Top 10 for LLM Applications

| ID | Risk | What to check |
| --- | --- | --- |
| LLM01 | Prompt injection | Treat all model inputs as untrusted; segregate system/user/tool contexts; filter before tool calls. |
| LLM02 | Insecure output handling | Never eval/execute/render LLM output without the same sanitization you'd apply to user input. |
| LLM03 | Training data poisoning | Vet fine-tuning and RAG corpora; sign and version them. |
| LLM04 | Model denial of service | Rate-limit, cap input size, enforce max tokens, timeout long-running calls. |
| LLM05 | Supply chain | Pin model providers and versions; verify plugin/tool provenance. |
| LLM06 | Sensitive information disclosure | Redact PII/secrets before sending to the model; test with fuzzing prompts. |
| LLM07 | Insecure plugin/tool design | Least-privilege tool schemas; validate outputs from tools before using them. |
| LLM08 | Excessive agency | Require human confirmation for destructive or high-blast-radius actions. |
| LLM09 | Overreliance | Show confidence/provenance; fall back safely when the model is uncertain. |
| LLM10 | Model theft | Rate-limit, watermark outputs where it matters, monitor for scraping patterns. |

### OWASP Agentic AI Top 10

| ID | Risk | What to check |
| --- | --- | --- |
| ASI01 | Agent identity & impersonation | Every agent has a distinct, auditable identity; credentials scoped per agent. |
| ASI02 | Memory poisoning | Sanitize and sign memory writes; separate short-term scratchpad from long-term memory. |
| ASI03 | Tool misuse | Narrow tool schemas; validate arguments before execution; sandbox where possible. |
| ASI04 | Privilege escalation | Principle of least privilege; deny agent access to credentials of other agents or humans. |
| ASI05 | Insecure orchestration | Plan/act loops must detect cycles, caps, and runaway budgets. |
| ASI06 | Reasoning/plan manipulation | Monitor plans for anomalous steps; require confirmation for deviation from allowed workflows. |
| ASI07 | Cross-agent collusion / cascading failures | Isolate agent blast radius; circuit-break on cascade. |
| ASI08 | Unauthorized external actions | Allowlist outbound domains; confirm before external-facing writes (Slack, email, git push). |
| ASI09 | Observability gaps | Structured traces for every tool call, prompt, and decision; redact before storage. |
| ASI10 | Agent supply chain | Verify provenance of prompts, tools, and sub-agents; pin versions. |

> OWASP updates these lists; check the [OWASP site](https://owasp.org/) periodically and revise this file when the next revision lands.

## Secret handling

- **Never commit secrets.** `gitleaks` runs pre-commit and in CI. If a secret hits a remote branch, rotate it — don't just rewrite history.
- **Use a secret manager** (Doppler, AWS Secrets Manager, 1Password, Vault, `.env` files *outside* the repo). `.env.example` in the repo documents the keys required, with fake values.
- **Never log secrets or PII.** Structured logging with a redaction layer at the logger, not per-call-site.
- **Rotate regularly.** Credentials tied to an individual are rotated on departure.

## Review gates (what blocks merge)

- `pre-commit.yml` passes (gitleaks, hygiene).
- `security-scan.yml` passes (vuln scan for detected language).
- `codeql.yml` passes (SAST).
- `license-check.yml` passes (no copyleft deps unless explicitly waived in an ADR).
- Any OWASP-relevant change has a Security Considerations note in the PR body.

## When to file a security ADR

If the change:
- Moves an auth boundary,
- Introduces a new tool the agent can call,
- Changes what's in memory or logs,
- Adds a new outbound destination,

…it gets an ADR. See `docs/decisions/`.

## Related skills

- `/security-review` — walk all 30 OWASP items against current diff; produce prioritized findings.
- `/dep-audit` — run language-specific vulnerability scanners.
