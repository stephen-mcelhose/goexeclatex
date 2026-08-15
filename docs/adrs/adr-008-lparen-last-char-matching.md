---
id: adr-008
title: "ADR-008: Match brackets by last character of token value, not full value"
status: accepted
date: 2026-08-15
---

# ADR-008: Match brackets by last character of token value, not full value

## Status

Accepted

## Context

The lexer's LPAREN/RPAREN patterns capture full forms including the `\left`/`\right`
prefix:

- `\left(` → `LPAREN{Value: "\\left("}`
- `\right)` → `RPAREN{Value: "\\right)"}`
- `(` → `LPAREN{Value: "("}`
- `)` → `RPAREN{Value: ")"}`

The parser's `closerFor` function and bracket-matching check in `parseGroup`
originally compared token values directly (e.g. `tok.Value != closer`). This
caused `\left( expr \right)` to fail: `closerFor("\\left(")` returned `""` (no
match), and `"\\right)" != ")"` triggered a false mismatch error.

## Decision

Key on the **last character** of the token value in both directions:

- `closerFor(opener)` switches on `opener[len(opener)-1]`.
- The bracket-match check in `parseGroup` compares
  `string(tok.Value[len(tok.Value)-1])` against `closer`.

The last character is always the actual bracket character (`(`, `)`, `[`, `]`,
`{`, `}`) regardless of whether a `\left`/`\right` prefix is present.

## Consequences

- `\left( expr \right)` parses correctly and identically to `( expr )`.
- Mismatched brackets (`\left( expr \right]`) are still caught correctly because
  `)` ≠ `]`.
- The fix is contained to `closerFor` and one comparison in `parseGroup`.
- No changes to the lexer or token types are required.
