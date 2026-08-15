---
id: adr-015
type: decision
title: "ADR-015: Allow real odd-integer roots of negative radicands"
description: >
  Supersedes parser-extensions §4.3’s blanket ban on negative radicands for
  \sqrt[n]{x}. When n is an odd integer, evaluate the real nth root (sign-preserving).
  Even / non-integer indices still domain-error on x < 0.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T16:22:00Z
supersedes: specs/parser-extensions.md §4.3 (negative radicand row for 2-arg sqrt)
tags: [adr, eval, sqrt, domain, v0.4]
---

# ADR-015: Allow real odd-integer roots of negative radicands

## Status

Accepted — amends [[specs/parser-extensions]] §4.3. Does not change `\sqrt{x}`
(one-arg square root still rejects `x < 0`).

## Context

v0.4 §4.3 specified:

> 2 (`x`, `n`) → `math.Pow(x, 1/n)`; domain error if `n == 0`; domain error if `x < 0`

That refused `\sqrt[3]{-8}`, which has a well-defined **real** cube root (−2).
GAN review flagged the gap between the written rule and real-analysis /
AMS-facing expectations for odd roots.

Go’s `math.Pow` returns **NaN** for a negative base with a fractional exponent
(`Pow(-8, 1/3) == NaN`), so a naïve `Pow(x, 1/n)` cannot implement the real
odd root even if the domain check is lifted.

## Decision

1. For two-arg `\sqrt[n]{x}` (`FunctionNode("sqrt", [x, n])`):
   - If `n == 0` → domain error (unchanged).
   - If `x >= 0` → `math.Pow(x, 1/n)` (unchanged).
   - If `x < 0` and `n` is an **odd integer** (finite, `n == trunc(n)`,
     `|n| >= 1`, `int64(n) % 2 != 0`) → evaluate
     `−math.Pow(−x, 1/n)` (real principal odd root; avoids NaN).
   - If `x < 0` otherwise (even integer index, non-integer index, etc.) →
     domain error (unchanged spirit of “no complex results”).
2. One-arg `\sqrt{x}` remains `math.Sqrt` with domain error on `x < 0`.
3. Update [[specs/parser-extensions]] §4.3 and tests accordingly.

## Consequences

- `\sqrt[3]{-8}` → approximately `−2` (float noise allowed in tests/goldens).
- `\sqrt[2]{-4}` still domain-errors.
- `\sqrt[-3]{-8}` is allowed if −3 is treated as an odd integer index
  (`1/n` negative → reciprocal of an odd root); implementers MUST apply the
  same odd-integer test to `n` before the sign-preserving formula.
- Spec §4.3’s blanket “domain error if `x < 0`” for two-arg sqrt is
  **superseded** by this ADR.

## Sources

- GAN review on issue [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)
- [[specs/parser-extensions]] §4.3 (prior rule)
- Go `math.Pow` / `math.Cbrt` behaviour on negative inputs
