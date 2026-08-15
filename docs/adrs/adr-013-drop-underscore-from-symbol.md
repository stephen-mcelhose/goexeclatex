---
id: adr-013
type: decision
title: "ADR-013: Drop underscore from SYMBOL pattern; remove \\_  pre-pass"
description: >
  Removes underscore from the SYMBOL token pattern and eliminates the \_ → _
  pre-pass. In evaluable LaTeX math, _ is always the subscript operator; it has
  no place in a symbol name. The pre-pass was solving a problem that does not
  exist in this domain.
status: accepted
date: 2026-08-15
supersedes: plan.md §Phase 1 "Escape rule" and SYMBOL pattern row
tags: [adr, lexer, symbol, underscore]
---

# ADR-013: Drop underscore from SYMBOL pattern; remove `\_` pre-pass

## Status

Accepted — supersedes the "Escape rule" paragraph and the SYMBOL pattern row
(`[A-Za-z_][A-Za-z_0-9]*`) in `docs/plan.md` §Phase 1 (Lexer).

## Context

`plan.md` originally specified:

> **Escape rule:** `\_` is a literal underscore within a symbol name, not a
> subscript trigger. The lexer must resolve `\_` → `_` before the `_` token
> rule fires, so that `MY\_CODE_{0}` lexes as `SYMBOL("MY_CODE") SUBSCRIPT …`.

This implied that:

1. The SYMBOL pattern should include underscore: `[A-Za-z_][A-Za-z_0-9]*`.
2. A `\_` → `_` string substitution must run before scanning (the "pre-pass").

Issue [#1](https://github.com/stephen-mcelhose/goexeclatex/issues/1) exposed
that the pre-pass shrinks the input string, causing every token position
*after* a `\_` occurrence to be off by one byte in error messages.

When investigating the fix we consulted `docs/latex-math-evaluable-spec.md`
(derived from the AMS Short Math Guide). The spec is unambiguous:

- `_` is exclusively the **subscript operator** (`x_{i}`, `\sum_{i=a}^{b}`).
- Variable names are **single letters** or **Greek-letter commands**
  (`x`, `n`, `\alpha`, `\beta`, …).
- Multi-word identifiers such as `MY_CODE` do not appear anywhere in evaluable
  LaTeX math.

The escape rule was solving a problem that does not exist in this domain.

## Decision

1. **Remove `_` from the SYMBOL pattern.**
   New pattern: `[A-Za-z][A-Za-z0-9]*`

2. **Remove the `\_` pre-pass entirely.**
   `\_` in input produces an "unexpected character" error — it is not valid
   evaluable LaTeX math.

3. **Remove the `\_` scan pattern** (priority 12) and §7.0 rewrite added
   during the aborted token-level-rewrite approach.

4. **Update `plan.md`** to reflect the new SYMBOL pattern and remove the
   escape-rule paragraph.

5. **Update `docs/specs/lexer.md`** §4 and §5.2 accordingly.

## Consequences

- Issue [#1](https://github.com/stephen-mcelhose/goexeclatex/issues/1)
  (pre-pass byte-offset drift) is closed: there is no pre-pass.
- The existing tests `TestEscapedUnderscore` and
  `TestEscapedUnderscoreSubscriptContext` encoded the old (now-superseded)
  behaviour and must be deleted.
- The new tests `TestEscapedUnderscoreToken`, `TestEscapedUnderscoreInSymbol`,
  `TestEscapedUnderscorePositions`, and `TestEscapedUnderscoreMultiple` added
  during the aborted approach must also be deleted — they tested a feature that
  no longer exists.
- In Tier 0.3, `_` will be added as the `UNDERSCORE` token type (already
  reserved in `token.go`). No change is needed there now.
- Callers passing `MY\_CODE`-style expressions will receive a lexer error.
  This is correct: such expressions are not valid evaluable LaTeX math.
