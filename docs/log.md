# Wiki Log

<!-- Append-only. Never edit existing entries. -->

## [2026-08-14] init | wiki initialized at ~/wiki
## [2026-08-14] ingest | LaTeX Math — Evaluable Subset Spec (updated: AMS short-math-guide.tex read via pandoc)
## [2026-08-14] ingest | evaluatex — Reference Implementation
## [2026-08-14] ingest | goexeclatex — Gap Analysis
## [2026-08-15] restructure | plan/roadmap split (ADR-012); roadmap.md created; index updated
## [2026-08-15] lint | 22 pages checked, 18 issues found, 18 fixed

Pages checked: 6 indexed + 16 discovered via glob. Issues found and fixed:
- Index gaps (16): how-to, specs/cli, specs/eval, specs/lexer, specs/parser, adrs/adr-001 through adr-011 — all added to index.md.
- Missing OKF frontmatter (2): how-to.md had none (added full OKF block); examples.md missing timestamp (added).
Advisory (not fixed — human judgment needed):
- specs/*.md and adrs/*.md use domain-specific frontmatter (id/status/version or id/status/date) that conflicts with OKF. They have title but lack type, description, timestamp. Upgrading would interfere with proto-spec-writer conventions. Recommend deciding on canonical frontmatter format for specs/ADRs.
- No wikilinks from plan.md body into specs, and specs don't wikilink back. Cross-link graph is sparse in this direction.

## [2026-08-15] ingest | runbooks/ingest-reference-implementation — how to ingest a reference impl
## [2026-08-15] lint | 29 pages checked, 9 issues found, 7 fixed

Pages checked: 29 (all docs/**/*.md). Issues found:
- Stale (fixed): adrs/adr-003 — amended to reflect v0.2 arcsin/arccos/arctan aliases; status → amended.
- Stale (fixed): goexeclatex-gap-analysis.md — Tier 2 table updated with ✅/⏳ status column; v0.2 items marked done.
- Stale (fixed): how-to.md — Supported syntax table expanded with v0.2 functions (hyperbolic, ln/log/exp, binom, dfrac, arcsin); stale \acos example replaced with \arccos.
- Stale (fixed): examples.md — paren-syntax shell examples corrected to brace syntax (\sin{\pi/6} etc.); \sqrt[3]{27} (v0.2) note corrected to "deferred to future milestone".
- Index gap (fixed): index.md ADR-003 summary updated to reflect amended status.
- Orphan (fixed): examples.md — added [[examples]] wikilink from how-to.md.
Advisory (not fixed — human judgment needed):
- Orphan: ADRs 001–003, 005–008 referenced only from index.md. Adding inbound links from specs would improve graph density but requires editing spec bodies non-trivially.
- Orphan: runbooks/ingest-reference-implementation.md — standalone runbook, acceptable.

## [2026-08-15] lint | 30 pages checked, 4 issues found, 4 fixed

Fixed: index gap (specs/subscripts-largeops added to index.md); ADR-004 status updated from accepted → superseded (lVert implemented v0.3); gap-analysis Tier 3 table updated with status column (4 items ✅, 1 ⏳ #8); roadmap "Current focus" updated to reflect v0.3 complete.

## [2026-08-15] feat | v0.3 complete — subscripts, big operators, norm (closes #6)

Implemented: UNDERSCORE + NORM lexer tokens; subscript grammar (`x_{i}`, `x_i`); `\sum`/`\prod` with integer bounds and iteration scope; `\lVert…\rVert` norm. Key design decisions: `ScopeLookup` interface for inner-scope shadowing; `parseSuperArg` to enforce `x_{i}^2` ordering; `normDepth` guard mirrors `absDepth`. `\\_` and `\log_b` deferred to #8. Specs: `docs/specs/subscripts-largeops.md`. testdata: `cmd/goexeclatex/testdata/stdin.txt` extended with v0.3 cases.
