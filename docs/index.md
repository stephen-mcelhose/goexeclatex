# Wiki Index

| Page                                        | Summary                                                                     | Tags                                          |
| ------------------------------------------- | --------------------------------------------------------------------------- | --------------------------------------------- |
| [[plan]]                                          | Stable architecture record: pipeline design, grammar, package layout, error policy   | goexeclatex, architecture, cli                      |
| [[roadmap]]                                       | Forward-looking milestone tracker: current focus, v0.4/#8, vFuture                  | goexeclatex, roadmap, milestones                    |
| [[examples]]                                      | Concrete CLI invocations: arithmetic, trig, AUPC, WI, TH                            | goexeclatex, examples, cli                          |
| [[how-to]]                                        | User guide: stdin/-e/-p/-v flags, supported syntax, error handling, shell compose   | goexeclatex, cli, how-to                            |
| [[latex-math-evaluable-spec]]                     | Catalogue of LaTeX math constructs with numeric evaluation semantics                 | latex, math, spec, goexeclatex                      |
| [[evaluatex-reference-implementation]]            | Deep read of arthanzel/evaluatex JS parser-evaluator; architecture + gaps            | reference-implementation, javascript, latex         |
| [[goexeclatex-gap-analysis]]                      | Gap analysis: evaluatex coverage vs full spec; tiered roadmap                        | goexeclatex, gap-analysis, latex                    |
| [[specs/lexer]]                                   | Normative RFC 2119 lexer spec: token types, char-mode algorithm, pre/post-lex passes | lexer, spec, goexeclatex                           |
| [[specs/parser]]                                  | Normative RFC 2119 parser spec: grammar, AST node types, error policy               | parser, spec, goexeclatex                           |
| [[specs/eval]]                                    | Normative RFC 2119 evaluator spec: scope, function dispatch, domain errors           | eval, spec, goexeclatex                             |
| [[specs/cli]]                                     | Normative RFC 2119 CLI spec: flags, input resolution, exit codes, output format      | cli, spec, goexeclatex                              |
| [[specs/subscripts-largeops]]                     | Normative RFC 2119 spec: subscript grammar, \sum/\prod bounds, \lVert norm (v0.3)   | parser, eval, spec, goexeclatex, subscripts         |
| [[specs/parser-extensions]]                       | Normative RFC 2119 spec: floor/ceil, \sqrt[n], variadic parens, log_b, \bmod (v0.4) | parser, eval, spec, goexeclatex, issue-8            |
| [[adrs/adr-001-lexer-ingroup-param]]              | ADR-001: Add inGroup param to lexExpression for unclosed-brace detection             | adr, lexer                                          |
| [[adrs/adr-002-utf8-error-deferred]]              | ADR-002: Defer UTF-8 char-mode error quality to post-Tier-1                         | adr, lexer, deferred                                |
| [[adrs/adr-003-asin-not-arcsin]]                  | ADR-003: asin/acos/atan for Tier 1; arcsin/arccos/arctan added as aliases in v0.2   | adr, eval, trig                                     |
| [[adrs/adr-004-lVert-deferred]]                   | ADR-004: Defer \lVert…\rVert norm support to Tier 2 — superseded: implemented v0.3  | adr, lexer, deferred                                |
| [[adrs/adr-005-no-group-node]]                    | ADR-005: Brackets are transparent — no GroupNode in the AST                          | adr, parser, ast                                    |
| [[adrs/adr-006-implicit-multiply-at-product]]     | ADR-006: Implicit multiplication resolved at the product rule                        | adr, parser, implicit-multiply                      |
| [[adrs/adr-007-absdepth-pipe-ambiguity]]          | ADR-007: Track absDepth to resolve PIPE implicit-multiply ambiguity                  | adr, parser, abs                                    |
| [[adrs/adr-008-lparen-last-char-matching]]        | ADR-008: Match brackets by last character of token value, not full value             | adr, parser, brackets                               |
| [[adrs/adr-009-explicit-domain-errors-over-nan]]  | ADR-009: Return explicit errors for domain violations instead of NaN                 | adr, eval, errors                                   |
| [[adrs/adr-010-inf-result-exits-zero]]            | ADR-010: ±Inf results print to stdout and exit 0                                     | adr, cli, inf                                       |
| [[adrs/adr-011-cli-error-prefix-stripping]]       | ADR-011: CLI strips internal package prefixes from user-facing error messages        | adr, cli, errors                                    |
| [[adrs/adr-012-plan-roadmap-separation]]          | ADR-012: Separate stable architecture plan from forward-looking roadmap              | adr, docs, project-management                       |
| [[adrs/adr-013-drop-underscore-from-symbol]]      | ADR-013: Drop underscore from SYMBOL pattern; remove `\_` pre-pass                  | adr, lexer, symbol, underscore                      |
| [[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]] | ADR-014: Mixed `\lfloor…\rceil` round convention is out of scope                 | adr, parser, floor, ceil, round, out-of-scope       |
| [[adrs/adr-015-odd-integer-roots-of-negatives]]   | ADR-015: Allow real odd-integer roots of negative radicands                     | adr, eval, sqrt, domain, v0.4                       |
| [[runbooks/ingest-reference-implementation]]      | How to clone, read, and ingest a reference implementation into the wiki              | runbook, llm-wiki, ingest, reference-implementation |
