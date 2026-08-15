---
type: concept
title: LaTeX Math — Evaluable Subset Spec
description: Catalogue of LaTeX math-mode constructs that have a numeric evaluation semantics, partitioned by category and complexity — the target spec for goexeclatex.
resource: https://ctan.org/pkg/short-math-guide
tags: [latex, math, spec, goexeclatex]
timestamp: 2026-08-14T22:43:35Z
---

# LaTeX Math — Evaluable Subset Spec

LaTeX math mode is far richer than what can be numerically evaluated. This page
catalogues the subset that is *evaluable* — i.e. can be reduced to a number —
and is therefore in scope for [[goexeclatex-gap-analysis]]. Purely typographic
commands (spacing, font sizing, alignment) are excluded.

Source: AMS Short Math Guide v2.0 (2017), Michael Downes / AMS. Raw file:
`~/wiki/raw/short-math-guide.tex` (LPPL 1.3c). Converted to Markdown via pandoc
3.10.2 and read in full. Overleaf symbol reference used only for initial framing.

---

## Arithmetic Operators

| LaTeX          | Meaning                  | Evaluable? |
| -------------- | ------------------------ | ---------- |
| `+`            | Addition                 | ✅          |
| `-`            | Subtraction / negation   | ✅          |
| `*`            | Multiplication (ASCII)   | ✅          |
| `/`            | Division (ASCII)         | ✅          |
| `^`            | Exponentiation           | ✅          |
| `\times`       | Multiplication           | ✅          |
| `\cdot`        | Multiplication (dot)     | ✅          |
| `\div`         | Division symbol          | ✅          |
| `\pm`          | Plus-or-minus            | ⚠️ ambiguous — yields ± pair |
| `!`            | Factorial (postfix)      | ✅          |

---

## Fractions & Roots

| LaTeX                  | Meaning             | Evaluable? |
| ---------------------- | ------------------- | ---------- |
| `\frac{a}{b}`          | Fraction a/b        | ✅          |
| `\dfrac{a}{b}`         | Display fraction    | ✅ (same value as `\frac`) |
| `\tfrac{a}{b}`         | Text fraction       | ✅ (same value as `\frac`) |
| `\cfrac{a}{b}`         | Continued fraction  | ✅ (same value as `\frac`) |
| `\sqrt{x}`             | Square root         | ✅          |
| `\sqrt[n]{x}`          | n-th root           | ✅          |
| `\binom{n}{k}`         | Binomial coefficient | ✅         |
| `\dbinom{n}{k}`        | Display binomial    | ✅ (same value) |
| `\tbinom{n}{k}`        | Text binomial       | ✅ (same value) |

---

## Standard Functions (LaTeX operator names)

| LaTeX         | Function          | Evaluable? |
| ------------- | ----------------- | ---------- |
| `\sin`        | Sine              | ✅          |
| `\cos`        | Cosine            | ✅          |
| `\tan`        | Tangent           | ✅          |
| `\csc`        | Cosecant          | ✅          |
| `\sec`        | Secant            | ✅          |
| `\cot`        | Cotangent         | ✅          |
| `\arcsin`     | Arc sine          | ✅          |
| `\arccos`     | Arc cosine        | ✅          |
| `\arctan`     | Arc tangent       | ✅          |
| `\sinh`       | Hyperbolic sine   | ✅          |
| `\cosh`       | Hyperbolic cosine | ✅          |
| `\tanh`       | Hyperbolic tangent| ✅          |
| `\coth`       | Hyperbolic cotangent | ✅       |
| `\ln`         | Natural logarithm | ✅          |
| `\log`        | Logarithm (base 10 or base e) | ✅ |
| `\log_{b}`    | Logarithm base b (subscript) | ✅ |
| `\exp`        | Exponential       | ✅          |
| `\abs`        | Absolute value    | ✅ (via `|...|`) |
| `\min`        | Minimum           | ✅          |
| `\max`        | Maximum           | ✅          |
| `\gcd`        | Greatest common divisor | ✅   |
| `\lcm`        | Least common multiple | ✅     |
| `\det`        | Determinant       | ⚠️ needs matrix |
| `\operatorname{name}` | Custom named op | ✅ (lookup) |

---

## Floor, Ceiling, Rounding

| LaTeX                        | Meaning       | Evaluable? |
| ---------------------------- | ------------- | ---------- |
| `\lfloor x \rfloor`          | Floor         | ✅          |
| `\lceil x \rceil`            | Ceiling       | ✅          |
| `\lfloor x \rceil`           | Round (nearest) | ❌ out of scope — [[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]] |

> **Note (from AMS guide §Vertical bar notations):** `|...|` is ambiguous for
> spacing. The guide recommends `\lvert...\rvert` for absolute value and
> `\lVert...\rVert` for norms. Both forms must be recognised by goexeclatex.

---

## Grouping & Delimiters

| LaTeX                     | Meaning              | Evaluable? |
| ------------------------- | -------------------- | ---------- |
| `(...)`, `[...]`, `{...}` | Grouping             | ✅          |
| `\left( \right)`          | Scaled parens        | ✅ (same as parens) |
| `\left[ \right]`          | Scaled brackets      | ✅          |
| `\left\{ \right\}`        | Scaled braces        | ✅          |
| `\|...\|`                 | Norm / abs value     | ✅          |
| `\langle \rangle`         | Angle brackets       | ⚠️ inner product needs 2 args |

---

## Subscripts & Superscripts

| LaTeX           | Meaning                   | Evaluable? |
| --------------- | ------------------------- | ---------- |
| `x^{n}`         | Exponentiation            | ✅          |
| `x_{i}`         | Subscript (indexing)      | ⚠️ requires symbol table |
| `\log_{b}(x)`   | Log base b                | ✅ (special case) |

---

## Big Operators (Discrete Bounds)

| LaTeX                          | Meaning          | Evaluable?                     |
| ------------------------------ | ---------------- | ------------------------------ |
| `\sum_{i=a}^{b} f(i)`         | Summation        | ✅ with discrete numeric bounds |
| `\prod_{i=a}^{b} f(i)`        | Product          | ✅ with discrete numeric bounds |
| `\int_{a}^{b} f(x) dx`        | Integral         | ❌ symbolic; out of scope       |
| `\lim_{x \to a} f(x)`         | Limit            | ❌ symbolic; out of scope       |
| `\sum_{i=a}^{\infty}`         | Infinite series  | ❌ symbolic; out of scope       |

---

## Constants

| LaTeX       | Value            | Evaluable? |
| ----------- | ---------------- | ---------- |
| `\pi`       | π ≈ 3.14159…     | ✅          |
| `e`         | e ≈ 2.71828…     | ✅          |
| `\tau`      | 2π ≈ 6.28318…    | ✅          |
| `\phi`      | φ ≈ 1.61803…     | ✅          |
| `\infty`    | +∞               | ✅ (as `math.Inf(1)` in Go) |
| `\theta`    | Variable         | ⚠️ context-dependent |
| `\alpha`, `\beta`, `\gamma`, `\delta`, `\mu`, `\sigma`, `\omega` | Variables | ⚠️ user-supplied |

---

## Modular Arithmetic

| LaTeX               | Meaning       | Evaluable? |
| ------------------- | ------------- | ---------- |
| `a \mod b`          | Modulo (spacing variant)   | ✅ |
| `a \pmod{b}`        | Modulo with parens         | ✅ |
| `a \bmod b`         | Modulo binary op form      | ✅ |
| `a \pod{b}`         | Modulo parens, no "mod"    | ✅ |

---

## Out-of-Scope (Symbolic / Typographic Only)

- `\int`, `\oint`, `\iint`, `\iiint` — integration
- `\lim` with symbolic variable (`x \to \infty`)
- `\begin{matrix}...\end{matrix}` and all matrix environments
- `\begin{cases}...\end{cases}` — piecewise (could be conditional eval)
- `\partial` — partial derivative (symbolic)
- `\nabla` — gradient (symbolic)
- All spacing commands: `\,`, `\;`, `\!`, `\quad`, `\qquad`
- Font/style commands: `\mathbf`, `\mathbb`, `\mathcal`, `\text`
- `\forall`, `\exists` — logic, not numeric

## Sources

- AMS Short Math Guide v2.0 (2017), Michael Downes / AMS: https://ctan.org/pkg/short-math-guide
- Raw source (read in full): `~/wiki/raw/short-math-guide.tex` — LPPL 1.3c, © AMS
- Overleaf symbol reference (initial framing only): https://www.overleaf.com/learn/latex/List_of_Greek_letters_and_math_symbols
