---
id: adr-006
title: "ADR-006: Implicit multiplication is resolved at the product rule"
status: accepted
date: 2026-08-15
---

# ADR-006: Implicit multiplication is resolved at the `product` rule

## Status

Accepted

## Context

Implicit multiplication (`2x`, `2(3+4)`, `xy`) must be inserted somewhere in
the recursive-descent grammar. Two candidate locations are `product` and `power`.

Inserting at `power` would mean `2x^2` parses as `(2x)^2 = 4x^2` rather than
the correct `2(x^2)`. This is wrong: exponentiation binds tighter than implicit
multiply in standard mathematical convention.

## Decision

Implicit multiplication is resolved at the `product` rule — after a complete
`power` sub-expression has been parsed. If the next token can begin another
`power` (NUMBER, SYMBOL, COMMAND, LPAREN, PIPE), a `BinaryNode{Op: "*"}` is
inserted between the two `power` nodes.

This gives implicit multiply the same precedence as explicit `*` and `/`, and
lower precedence than `^`.

```
product → power { ('*' | '/') power }
        | power power*    // implicit multiply
```

## Consequences

- `2x^2` → `2 * (x^2)` ✓ (correct)
- `2(3+4)` → `2 * (3+4)` ✓
- `xy` → `x * y` ✓
- The `product` parse function loops and checks `canBeginPower(peek())` after
  each `power` call, inserting an implicit `*` node when true.
