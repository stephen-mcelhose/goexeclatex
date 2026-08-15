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
- The last-char heuristic is load-bearing but implicit: it holds only because all
  lexer bracket patterns end with the actual bracket glyph. This assumption is not
  enforced by the lexer spec.

## Revisit criteria

Supersede this ADR if **any** of the following occur:

1. **A feature needs bracket-type information downstream.** If the eval layer,
   a pretty-printer, or any AST consumer needs to distinguish `\left(` from `(`
   (e.g. for large-delimiter rendering or LaTeX round-trip output), the right fix
   is to normalise in the lexer — strip the `\left`/`\right` prefix and emit bare
   bracket values — so the distinction never reaches the parser.

2. **A new bracket form breaks the last-char assumption.** If the lexer gains a
   bracket pattern whose token value does not end with the bracket glyph (e.g.
   `\Bigl\{` where the last char is `{` — actually fine — or some pathological
   form), audit `closerFor` and `parseGroup` immediately.

3. **The lexer spec is updated to guarantee normalised bracket values.** At that
   point the last-char logic becomes dead code and should be replaced with direct
   value comparison.

This decision is **stable** as long as the project remains evaluation-only with no
AST-to-LaTeX rendering requirement.

