---
name: dep-audit
description: Run the vulnerability scanners appropriate to the project's stacks (Go, npm, Python, Rust, etc.), then audit newly-introduced dependencies against the five-point check. Produce a prioritized upgrade list.
---

# Dep audit

Detect which stacks are present, run the matching scanners, and produce a prioritized report of vulnerabilities and pending upgrades.

## Stack detection

Check for these files in the repo:

| File | Stack | Primary scanner |
| --- | --- | --- |
| `go.mod` | Go | `govulncheck ./...` |
| `package.json` | Node.js / JS / TS | `npm audit --audit-level=high` (or `pnpm audit`, `yarn audit`) |
| `pyproject.toml` or `requirements.txt` | Python | `pip-audit` |
| `Cargo.toml` | Rust | `cargo audit` |
| `Gemfile` | Ruby | `bundle audit check --update` |
| `composer.json` | PHP | `composer audit` |
| `pom.xml` or `build.gradle` | Java | `mvn dependency-check:check` or `./gradlew dependencyCheckAnalyze` |

Run the scanner for every detected stack. If a stack is detected but the scanner isn't installed, report the install command and skip.

## What to check for each finding

1. **Severity** (critical / high / medium / low).
2. **Affected component** and version.
3. **Fixed version** (if available).
4. **Exploitability in this project** — is the vulnerable code path actually reached? Often a scanner flags a transitive dep that the project doesn't exercise.
5. **Upgrade cost** — is the fix a patch bump, minor, or major? Major bumps may require code changes.

## Also check: newly-introduced deps

For every dep added on the current branch (compare to `origin/main`):

- **Existence:** confirm the package exists on the official registry. (AI-hallucinated package names are a real supply-chain risk.)
- **License:** flag anything outside `MIT / Apache-2.0 / BSD / ISC / MPL-2.0`.
- **Maintenance:** last commit older than 12 months? Flag.
- **Necessity:** is this replaceable with ~50 lines of code? Ask the author in the PR.

See `claude-instructions/dependency-vetting.md` for the full five-point check.

## Output format

### Vulnerability report

| Severity | Package | Current | Fixed in | Scanner | Exploitable? | Action |
| --- | --- | --- | --- | --- | --- | --- |

### New dependency report (current branch vs. main)

| Package | Version | License | Last commit | 5-point check | Action |
| --- | --- | --- | --- | --- | --- |

### Prioritized upgrade list

1. Critical and high that are actually exploitable.
2. Critical and high not currently reached (document why).
3. Medium where fixes are cheap.
4. Low — defer.

### CI coverage check

Confirm `security-scan.yml` exists and includes a job for each detected stack. Call out gaps.

## Invariants

- Don't upgrade anything — produce recommendations only.
- Don't delete a dep the project uses; if necessity is suspect, open it as a discussion.
- Be explicit about what was and wasn't scanned. If Python is present but `pip-audit` isn't installed, say so.

## Related

- `claude-instructions/dependency-vetting.md` — the five-point check in full.
- `.github/workflows/security-scan.yml` — CI version of the same scans.
