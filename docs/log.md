# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-14] init | wiki initialized at ~/wiki
## [2026-08-14] ingest | LaTeX Math — Evaluable Subset Spec (updated: AMS short-math-guide.tex read via pandoc)
## [2026-08-14] ingest | evaluatex — Reference Implementation
## [2026-08-14] ingest | goexeclatex — Gap Analysis
## [2026-08-15] lint | 22 pages checked, 18 issues found, 18 fixed

Pages checked: 6 indexed + 16 discovered via glob. Issues found and fixed:
- Index gaps (16): how-to, specs/cli, specs/eval, specs/lexer, specs/parser, adrs/adr-001 through adr-011 — all added to index.md.
- Missing OKF frontmatter (2): how-to.md had none (added full OKF block); examples.md missing timestamp (added).
Advisory (not fixed — human judgment needed):
- specs/*.md and adrs/*.md use domain-specific frontmatter (id/status/version or id/status/date) that conflicts with OKF. They have title but lack type, description, timestamp. Upgrading would interfere with proto-spec-writer conventions. Recommend deciding on canonical frontmatter format for specs/ADRs.
- No wikilinks from plan.md body into specs, and specs don't wikilink back. Cross-link graph is sparse in this direction.

## [2026-08-15] ingest | runbooks/ingest-reference-implementation — how to ingest a reference impl
