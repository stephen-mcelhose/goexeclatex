# Wiki Index

| Page                                        | Summary                                                                     | Tags                                          |
| ------------------------------------------- | --------------------------------------------------------------------------- | --------------------------------------------- |
| [[plan]]                                          | Stable architecture record: pipeline design, grammar, package layout, error policy   | goexeclatex, architecture, cli                      |
| [[roadmap]]                                       | Forward-looking milestone tracker: current focus, v0.2/v0.3/vFuture, issue links    | goexeclatex, roadmap, milestones                    |
| [[examples]]                                      | Concrete CLI invocations: arithmetic, trig, AUPC, WI, TH                            | goexeclatex, examples, cli                          |
| [[how-to]]                                        | User guide: stdin/-e/-p/-v flags, supported syntax, error handling, shell compose   | goexeclatex, cli, how-to                            |
| [[latex-math-evaluable-spec]]                     | Catalogue of LaTeX math constructs with numeric evaluation semantics                 | latex, math, spec, goexeclatex                      |
| [[evaluatex-reference-implementation]]            | Deep read of arthanzel/evaluatex JS parser-evaluator; architecture + gaps            | reference-implementation, javascript, latex         |
| [[goexeclatex-gap-analysis]]                      | Gap analysis: evaluatex coverage vs full spec; tiered roadmap                        | goexeclatex, gap-analysis, latex                    |
| [[specs/lexer]]                                   | Normative RFC 2119 lexer spec: token types, char-mode algorithm, pre/post-lex passes | lexer, spec, goexeclatex                           |
| [[specs/parser]]                                  | Normative RFC 2119 parser spec: grammar, AST node types, error policy               | parser, spec, goexeclatex                           |
| [[specs/eval]]                                    | Normative RFC 2119 evaluator spec: scope, function dispatch, domain errors           | eval, spec, goexeclatex                             |
| [[specs/cli]]                                     | Normative RFC 2119 CLI spec: flags, input resolution, exit codes, output format      | cli, spec, goexeclatex                              |
| [[adrs/adr-001-lexer-ingroup-param]]              | ADR-001: Add inGroup param to lexExpression for unclosed-brace detection             | adr, lexer                                          |
| [[adrs/adr-002-utf8-error-deferred]]              | ADR-002: Defer UTF-8 char-mode error quality to post-Tier-1                         | adr, lexer, deferred                                |
| [[adrs/adr-003-asin-not-arcsin]]                  | ADR-003: Use asin/acos/atan abbreviations (evaluatex convention) over arcsin etc.   | adr, eval, trig                                     |
| [[adrs/adr-004-lVert-deferred]]                   | ADR-004: Defer \\lVert…\\rVert norm support to Tier 2                                | adr, lexer, deferred                                |
| [[adrs/adr-005-no-group-node]]                    | ADR-005: Brackets are transparent — no GroupNode in the AST                          | adr, parser, ast                                    |
| [[adrs/adr-006-implicit-multiply-at-product]]     | ADR-006: Implicit multiplication resolved at the product rule                        | adr, parser, implicit-multiply                      |
| [[adrs/adr-007-absdepth-pipe-ambiguity]]          | ADR-007: Track absDepth to resolve PIPE implicit-multiply ambiguity                  | adr, parser, abs                                    |
| [[adrs/adr-008-lparen-last-char-matching]]        | ADR-008: Match brackets by last character of token value, not full value             | adr, parser, brackets                               |
| [[adrs/adr-009-explicit-domain-errors-over-nan]]  | ADR-009: Return explicit errors for domain violations instead of NaN                 | adr, eval, errors                                   |
| [[adrs/adr-010-inf-result-exits-zero]]            | ADR-010: ±Inf results print to stdout and exit 0                                     | adr, cli, inf                                       |
| [[adrs/adr-011-cli-error-prefix-stripping]]       | ADR-011: CLI strips internal package prefixes from user-facing error messages        | adr, cli, errors                                    |
| [[adrs/adr-012-plan-roadmap-separation]]          | ADR-012: Separate stable architecture plan from forward-looking roadmap              | adr, docs, project-management                       |
| [[adrs/adr-013-drop-underscore-from-symbol]]      | ADR-013: Drop underscore from SYMBOL pattern; remove `\_` pre-pass                  | adr, lexer, symbol, underscore                      |
| [[runbooks/ingest-reference-implementation]]      | How to clone, read, and ingest a reference implementation into the wiki              | runbook, llm-wiki, ingest, reference-implementation |
