---
id: adr-001
title: "ADR-001: Add inGroup parameter to lexExpression for unclosed-brace detection"
status: accepted
date: 2026-08-15
---

# ADR-001: Add `inGroup` parameter to `lexExpression` for unclosed-brace detection

## Status

Accepted

## Context

The original `lexExpression(charMode bool)` signature had no way to distinguish
between reaching EOF at the top level (normal end of input) and reaching EOF
inside a `{…}` group (which means the input is malformed). As a result, inputs
like `2^{12`, `\sqrt{4`, and `\frac{1}{2` silently produced a truncated token
stream instead of returning an error.

## Decision

Add a second parameter: `lexExpression(charMode, inGroup bool)`.

- Top-level calls: `lexExpression(false, false)` — EOF is a normal termination.
- Group-opener calls: `lexExpression(false, true)` — EOF is an error
  (`"lexer: unexpected end of input: unclosed '{' group at position N"`).
- Char-mode calls: `lexExpression(true, false)` — EOF is already an error via
  the charMode branch.

## Consequences

- The parser can now rely on the invariant that every `{` token in the stream has
  a matching `}` — unclosed groups are caught at lex time.
- All recursive call sites in `lexer.go` must pass both arguments explicitly.
- `Lex()` (the public entry point) passes `(false, false)`.
