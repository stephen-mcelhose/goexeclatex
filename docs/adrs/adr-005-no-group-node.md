---
id: adr-005
title: "ADR-005: Brackets are transparent — no GroupNode in the AST"
status: accepted
date: 2026-08-15
---

# ADR-005: Brackets are transparent — no `GroupNode` in the AST

## Status

Accepted

## Context

Grouping brackets `(…)`, `[…]`, and `{…}` in LaTeX math carry no numeric
semantics of their own — they exist only to control precedence. A candidate
design would wrap the inner expression in a `GroupNode` to preserve bracket
type in the AST.

## Decision

Do not introduce a `GroupNode`. The parser **MUST** return the inner expression
node directly when it encounters any of the three grouping forms. The bracket
tokens are consumed and discarded.

This mirrors the evaluatex reference implementation, which similarly unwraps
grouping forms during parsing.

The only bracket-like construct that produces a distinct AST node is absolute
value (`|…|`), which becomes `FunctionNode{Name: "abs"}` because it has
evaluation semantics.

## Consequences

- AST is simpler and the evaluator does not need to handle a GroupNode case.
- Bracket type (`(` vs `[` vs `{`) is not preserved in the AST. This is
  acceptable because the evaluator only needs the value, not the presentation.
- If a future feature needs bracket-type information (e.g. interval notation),
  this decision will need to be revisited and superseded.
