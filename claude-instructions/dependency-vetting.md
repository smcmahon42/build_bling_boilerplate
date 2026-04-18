# Dependency vetting

Before adding any dependency, run the five-point check. "Run" is literal — if you haven't answered all five, the dependency isn't approved.

## The five-point check

1. **Existence & identity.** Does the package actually exist on the official registry? Is the name spelled exactly right? AI-assisted code is prone to hallucinated package names — a plausible-looking `npm install fast-json-lite` may be malware sitting on a typo'd name. Verify on the registry before installing.

2. **Version.** What is the latest stable version? Pin to it or a compatible range. Beta/RC versions need a justification.

3. **License.** What license is it under? Acceptable by default: **MIT, Apache-2.0, BSD (2/3-clause), ISC, MPL-2.0**. Copyleft licenses (**GPL, AGPL, LGPL**) require an ADR documenting the decision and any redistribution consequences. If the license is unclear, treat that as a red flag — don't use it.

4. **Maintenance.** When was the last commit? When was the last release? A project whose last activity is over 12 months old is a risk. Look at the issue tracker: are security issues triaged or ignored?

5. **Necessity.** Can you implement the equivalent in under ~50 lines of your own code? If yes, strongly consider doing that. Every dependency is:
   - A surface for supply-chain attacks,
   - A transitive dependency burden,
   - A piece of behavior you didn't author and can't easily change.

## When to add a dependency despite small size

- It implements a correctness-critical algorithm (crypto, parsing, etc.) where a hand-rolled version would be worse.
- It's a well-known standard (e.g. `zod`, `uuid`, `requests`) with a large user base and active maintenance.
- It's required for interop (protocol libraries, official SDKs).

## When to reject

- Single-maintainer package with low download counts and no recent commits.
- Transitive dependency trees that pull in dozens of indirect packages for a small utility.
- Package that bundles its own network calls or telemetry (review what it phones home).
- Project whose maintainers have been caught shipping malicious updates in the past.

## Process

1. Run the five-point check.
2. If accepted: add it, commit with `chore(deps): add <pkg> for <reason>`, document *why* in the commit body.
3. If it's a significant architectural choice (framework, runtime, database client), write an ADR.
4. Run `/dep-audit` periodically to catch newly-disclosed vulnerabilities.

## Automation

- **Dependabot** (`.github/dependabot.yml`) opens weekly PRs for updates. Review them; don't auto-merge without tests.
- **`security-scan.yml`** runs language-specific vuln scanners per PR.
- **`license-check.yml`** blocks PRs that introduce copyleft deps.

## AI-specific caution

Claude and other coding assistants will occasionally invent plausible package names. Always:
- Verify on the official registry.
- Check download counts and star counts.
- Skim the package's actual `README` — if it's empty or generic, stop.

## Related skills

- `/dep-audit` — run the vulnerability scanners your project needs.
