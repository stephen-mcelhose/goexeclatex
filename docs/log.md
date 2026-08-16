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

## [2026-08-15] feat | v0.4 parser extensions implemented (#8) — GAN pending

Branch `feat/8-parser-extensions-v0.4`: floor/ceil, `\sqrt[n]`, min/max/gcd, `\log_{b}(x)`, `\bmod`. Spec [[specs/parser-extensions]]; ADR-014 mixed round out of scope. Deferred: `\pmod`/`\mod`/`\pod`. Next: GAN-style closing review then merge/close #8.

## [2026-08-15] decision | ADR-015 — odd-integer roots of negatives

Accepted [[adrs/adr-015-odd-integer-roots-of-negatives]]: supersedes parser-extensions §4.3 blanket ban; `\sqrt[3]{-8}` → −2 via `−Pow(|x|,1/n)`.

## [2026-08-15] lint | GAN punch-list docs drift cleared (v0.4)

Pages updated: specs/eval §9, specs/parser §10, specs/subscripts-largeops §8, examples.md (deduped + v0.4 examples), latex-math-evaluable-spec modular table, gap-analysis intros, roadmap GAN checkbox, index (ADR-015), AGENTS.md package layout. Spec parser-extensions §4.3 amended for ADR-015.

## [2026-08-15] lint | 30 content pages checked, 12 issues found, 12 fixed

Inventory: 33 md under docs/ (excl. counting raw separately). Content pages ~30 + index/log/AGENTS.

Issues found and fixed:
- Stale (2): specs/parser.md §4.2 still deferred `\lVert` to Tier 2 → updated to NORM/v0.3 + ADR links; specs/eval.md §6.4 still deferred `\log_{b}` to Tier 0.3 → points at [[specs/parser-extensions]].
- Contradiction (1): specs/lexer.md §3 SYMBOL table still had `_` in pattern → aligned with ADR-013; added EQUALS/NORM/FLOOR/CEIL/BMOD rows.
- Orphans (8): ADRs 001–006, 008 and runbooks/ingest had only index inbound links → added body `[[wikilink]]`s from lexer/parser/eval/cli/subscripts/evaluatex Sources or related sections.
- Frontmatter (1): adr-013 missing `timestamp` → added.
- Index summary (1): roadmap blurb still said v0.2/v0.3 → v0.4/#8.
- Wrong ADR cite (1): parser.md implicit-multiply row cited ADR-003 → [[adrs/adr-006-implicit-multiply-at-product]].

Advisory (not fixed — human judgment / intentional):
- Specs/ADRs use hybrid frontmatter (`id`/`status`/`date` plus OKF `type`); full OKF-only upgrade left as prior advisory.
- `evaluatex-reference-implementation.md` SYMBOL row still shows evaluatex’s `[A-Za-z_…]` pattern — correct for that reference, not goexeclatex (divergence documented in ADR-013).
- `raw/` is read-only; not indexed (by design).

## [2026-08-15] ingest | #15 C+B — host Tier A skills + encode Phase 3

Decision C+B on [#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15): encode planning/Phase 3 in root AGENTS.md + roadmap; vendor genericized `clarify-issue`, `breakdown-issue`, `ac-plan-loop` under `.agents/skills/`; declare in-repo `llm-wiki` canonical. Artifacts: `.agents/issues/issue-<N>/`. Tier B / community pack deferred. Issue closed.

## [2026-08-15] lint | PR #21 GAN punch-list (process skills)

Fixed: roadmap v0.4 “pending GAN” stale header; AGENTS/roadmap Phase 3 process vs plan.md Evaluator disambiguation; docs/AGENTS.md points at in-repo llm-wiki; skills.json registers local clarify/breakdown/ac-plan-loop. Deferred: Cursor slash `/skill` discovery vs description-only; personal ~/.cursor/skills shadowing; #15 closed before merge.

## [2026-08-15] docs | ac-plan-loop v1.2 — restore #8 fidelity

Gate A defaults-OK loop; require Completion strategy + Dangerous status quo; Stage C GAN discriminator table; gate-exit issue comments; worked example `.agents/issues/issue-8/`. clarify-issue gains matching sections.

## [2026-08-15] decision | ADR-016 self-host Tier A agent skills (C+B)

Accepted [[adrs/adr-016-self-host-agent-skills]]: encode loop + host Tier A; Tier B/community deferred; revisit on degraded experience (loop ignored, shadowing, net slowdown, proven better spare parts).

## [2026-08-15] spec | #20 Task 0 — public library API ADR + spec

Added [[adrs/adr-017-public-library-api]] and [[specs/library]]; roadmap current focus → [#20](https://github.com/stephen-mcelhose/goexeclatex/issues/20) on `feat/20-public-library-api`. Planning artifacts: `.agents/issues/issue-20/` (process twin of issue-8).

## [2026-08-15] feat | #20 public library API (Tasks 1–4)

Shipped module-root `goexeclatex.Eval` + `SyntaxError`/`EvalError`; CLI rewires through public API (goldens green); `min`/`max`/`gcd` arity guards; README/AGENTS/plan/how-to/cli spec updated. Float precision helper omitted (ADR-017). Branch: `feat/20-public-library-api`.

## [2026-08-15] lint | #20 GAN punch-list

Fixed: library tests for empty expr (`*SyntaxError`), builtin override, `+\infty` success; `specs/library` §4.4 documents chosen empty behaviour; `examples.md` points at public API. Deferred: public `Format`; broader arity table; AST export; optional module split so library consumers need not pull Cobra via `go.mod`.

## [2026-08-15] docs | zsh echo `\frac` footgun

zsh `echo` interprets `\f`/`\b` escapes, so `echo '\frac{…}' | goexeclatex` fails as `undefined symbol: rac`. how-to/examples/cli/plan/README/CLI help now prefer `goexeclatex -e` or `printf '%s\n'`.

## [2026-08-15] docs | ADR-016 — #20 packaging worked example

Amended [[adrs/adr-016-self-host-agent-skills]]: capability twin (#8) vs packaging/facade twin (#20); same spine, different density. `ac-plan-loop` v1.2.1 points at both artifact dirs.

## [2026-08-15] docs | examples — named AUPC/WI/TH/subsample cases

Expanded [[examples]] with named calculated-process cases: AUPC step + `\sum` forms (PestProgress/EpidemicCurve/SpikeAndFade), WI/TH severity-class sets, PlotStripFive subsample aggregates. Links [[specs/subscripts-largeops]] for the sum engine.

## [2026-08-15] spike | matrix math charter (#24)

Opened [#24](https://github.com/stephen-mcelhose/goexeclatex/issues/24) and [[spike-matrix-math]]: explore evaluable subset (`\det` vs matrix×matrix), gonum dependency, usability vs scalar `Eval`. Roadmap vFuture + [[latex-math-evaluable-spec]] cross-linked; no implementation.
