---
type: spike
title: goexeclatex — Gap Analysis
description: What evaluatex covers vs the full evaluable LaTeX math spec, partitioned into tractable issues for goexeclatex v0.
resource: https://github.com/stephen-mcelhose/goexeclatex
tags: [goexeclatex, gap-analysis, reference-implementation, latex]
timestamp: 2026-08-14T22:43:35Z
---

# goexeclatex — Gap Analysis

Derived from [[evaluatex-reference-implementation]] (what the JS reference does)
vs [[latex-math-evaluable-spec]] (what evaluable LaTeX math requires).

The gaps are categorised by difficulty to guide milestone planning.

---

## Tier 1 — Already in evaluatex (port these first)

These should be the baseline of goexeclatex v0.

| Feature               | LaTeX                     | Notes                              |
| --------------------- | ------------------------- | ---------------------------------- |
| Arithmetic            | `+`, `-`, `*`, `/`, `^`   | Core grammar                       |
| Implicit multiply     | `2x`, `(2)(3)`            | Parser lookahead rule              |
| Grouping              | `()`, `[]`, `{}`          | Same precedence                    |
| `\left\right`         | `\left( \right)`          | Token-level remap to parens        |
| `\frac{a}{b}`         | Fraction                  | 2-arity command                    |
| `\sqrt{x}`            | Square root               | 1-arity command                    |
| `\times`, `\cdot`     | Multiplication            | Token remap                        |
| `\sin/\cos/\tan`      | Trig                      | Forward to `math.Sin` etc.         |
| `\sec/\csc/\cot`      | Recip trig                | `1/cos(x)` etc.                    |
| `\asin/\acos/\atan`   | Inverse trig (short form) | Forward to `math.Asin` etc.        |
| `\asec/\acsc/\acot`   | Inverse recip trig        |                                    |
| `\|x\|` abs value     | `\|expr\|`                | Pipe-delimited                     |
| `x!` factorial        | Postfix                   | Iterative implementation           |
| Constants             | `\pi`, `e`, `\tau`, `\phi`| Pre-seeded symbol table            |
| Custom vars/fns       | User-supplied scope       | Map[string]float64 / func          |

---

## Tier 2 — Small gaps; straightforward to add

Well-understood semantics; no new grammar machinery needed beyond Tier 1.

✅ = implemented in v0.2 unless noted. Parser-extension items shipped in **v0.4** (#8).
`\mod`/`\pmod`/`\pod` remain deferred (congruence annotations, not `\bmod`).

| Feature                    | LaTeX                          | Effort   | Status | Notes                                      |
| -------------------------- | ------------------------------ | -------- | ------ | ------------------------------------------ |
| `\arcsin/\arccos/\arctan`  | Arc trig (canonical names)     | XS       | ✅ v0.2 | Aliases to asin/acos/atan                 |
| `\ln`                      | Natural log                    | XS       | ✅ v0.2 | `math.Log`                                |
| `\exp`                     | Exponential                    | XS       | ✅ v0.2 | `math.Exp`                                |
| `\log` (base 10)           | Log base 10                    | XS       | ✅ v0.2 | `math.Log10`                              |
| `\dfrac`, `\tfrac`, `\cfrac` | Display/text fraction        | XS       | ✅ v0.2 | Same value as `\frac`                     |
| `\sinh/\cosh/\tanh`        | Hyperbolic trig                | XS       | ✅ v0.2 | `math.Sinh` etc.                          |
| `\coth/\sech/\csch`        | Hyperbolic recip               | XS       | ✅ v0.2 | `1/tanh(x)` etc.                          |
| `\infty`                   | Infinity constant              | XS       | ✅ v0.2 | `math.Inf(1)` seeded in NewScope          |
| `\binom{n}{k}`             | Binomial coefficient           | S        | ✅ v0.2 | `n! / (k! * (n-k)!)`                      |
| `\min(a,b)`, `\max(a,b)`   | Min / max                      | XS       | ✅ v0.4 | Variadic n≥2 paren args (#8)              |
| `\lfloor x \rfloor`        | Floor                          | S        | ✅ v0.4 | FLOOR tokens + depth (#8)                 |
| `\lceil x \rceil`          | Ceiling                        | S        | ✅ v0.4 | CEIL tokens + depth (#8)                  |
| `\sqrt[n]{x}`              | n-th root                      | S        | ✅ v0.4 | Optional `[n]` (#8)                       |
| `\gcd(a,b)`                | GCD                            | S        | ✅ v0.4 | Variadic n≥2 (#8)                         |
| `\bmod`                    | Binary modulo                  | S        | ✅ v0.4 | BMOD token (#8)                           |
| `\mod`, `\pmod`, `\pod`    | Congruence annotations         | S        | ⏳ deferred | Not binary remainder; see ADR/spec §7.4 |
| Greek letter variables      | `\alpha`, `\beta`, `\mu`…      | S        | ✅ v0.2 | Work via `-v alpha=1.5` with no changes   |

---

## Tier 3 — Significant new grammar

Requires extending the parser with new production rules (subscripts, large operators).

✅ = implemented in v0.3. `\log_{b}(x)` shipped in v0.4 (#8).

| Feature                     | LaTeX                        | Effort | Status   | Notes                                           |
| --------------------------- | ---------------------------- | ------ | -------- | ----------------------------------------------- |
| `\\log_{b}(x)`               | Log base b (subscript)       | M      | ✅ v0.4  | `#8` / parser-extensions §6                |
| `x_{i}` subscript variables | Indexed variables            | M      | ✅ v0.3  | `_` token + symbol table with index             |
| `\\sum_{i=a}^{b} f(i)`       | Discrete summation           | L      | ✅ v0.3  | Big-op token, bounds parsing, iteration engine  |
| `\\prod_{i=a}^{b} f(i)`      | Discrete product             | L      | ✅ v0.3  | Same engine as `\\sum`                           |
| `\\lVert v \\rVert` norm      | Norm / Euclidean length      | M      | ✅ v0.3  | NORM token + `normDepth` guard (see ADR-004)    |

---

## Tier 4 — Out of scope for numeric evaluation

These are symbolic or require a CAS. Intentionally excluded from goexeclatex.

| Feature              | Reason                                      |
| -------------------- | ------------------------------------------- |
| `\int_{a}^{b}`       | Symbolic integration                        |
| `\lim_{x\to a}`      | Symbolic limits                             |
| `\sum_{i=1}^{\infty}` | Requires convergence analysis              |
| `\partial`           | Symbolic partial derivative                 |
| Matrix environments  | Linear algebra; separate domain            |
| `\det` of a matrix   | Needs matrix type                           |
| `\begin{cases}`      | Conditional — possible future extension     |

---

## Recommended Milestone Sequence

```
v0.1 — Tier 1 parity with evaluatex
  Lexer, parser (recursive descent), evaluator, symbol table, implicit multiply

v0.2 — Tier 2 functions & constants
  Aliases, hyperbolic, floor/ceil, nth root, gcd, binom, modulo, Greek vars

v0.3 — Tier 3: subscripts + large operators
  _ token, \log_b, \sum, \prod with discrete bounds

vFuture — \begin{cases} conditional eval, matrix support
```

## Sources

- [[evaluatex-reference-implementation]]
- [[latex-math-evaluable-spec]]
- `~/repos/evaluatex/src/`
