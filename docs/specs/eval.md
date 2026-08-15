---
id: spec-eval
type: rfc2119-spec
title: Evaluator Specification
description: Normative RFC 2119 specification for the goexeclatex evaluator — scope resolution, arithmetic and function dispatch, domain error handling, and IEEE 754 special-value policy.
status: active
version: 0.2
tags: [eval, spec, goexeclatex]
timestamp: 2026-08-15T00:00:00Z
sources:
  - docs/plan.md                              # error policy, milestone scope
  - docs/evaluatex-reference-implementation.md  # reference function implementations
  - docs/latex-math-evaluable-spec.md          # evaluable subset catalogue
  - docs/specs/parser.md                       # AST node types
---

# Evaluator Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHALL NOT**, **SHOULD**,
**MAY**, and **OPTIONAL** in this document are to be interpreted as described
in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the behaviour of `Eval(node Node, scope Scope) (float64, error)`
in `internal/eval`.

The evaluator performs a depth-first walk of the AST produced by `internal/parser`
and reduces it to a single `float64`.

## 3. API

```go
// Scope maps variable names to their float64 values.
// Keys are normalised identifiers (backslash-stripped, lowercased) matching
// the SymbolNode.Name convention from the parser.
type Scope map[string]float64

// NewScope returns a Scope pre-seeded with all built-in constants (§5).
// Callers MAY add user-defined variables before passing to Eval.
func NewScope() Scope

// Eval reduces node to a float64 using scope for symbol lookup.
// Returns a non-nil error for all conditions in §7.
func Eval(node parser.Node, scope Scope) (float64, error)
```

## 4. Node Evaluation Rules

The evaluator **MUST** handle every node type produced by the parser.

### 4.1 NumberNode

**MUST** return `node.Value` directly. No error possible.

### 4.2 SymbolNode

**MUST** look up `node.Name` in `scope`.

- If found, **MUST** return the associated `float64`.
- If not found, **MUST** return `eval: undefined symbol: <name>` (§7).

### 4.3 BinaryNode

**MUST** evaluate both `Left` and `Right` sub-trees recursively, then apply `Op`.

| Op  | Semantics                         | Error condition                  |
| --- | --------------------------------- | -------------------------------- |
| `+` | `left + right`                    | none                             |
| `-` | `left - right`                    | none                             |
| `*` | `left * right`                    | none                             |
| `/` | `left / right`                    | `right == 0` → division by zero  |
| `^` | `math.Pow(left, right)`           | none (Go propagates NaN/Inf)     |

Division by zero **MUST** return the error `eval: division by zero`. The check
**MUST** be `right == 0.0` (exact float comparison); the result is never `+Inf`.

### 4.4 UnaryNode

| Op  | Semantics                    | Error condition              |
| --- | ---------------------------- | ---------------------------- |
| `-` | `-operand`                   | none                         |
| `!` | factorial of operand (§4.4.1)| non-integer or negative (§7) |

#### 4.4.1 Factorial

The operand **MUST** be a non-negative integer. The evaluator **MUST** check:

1. `operand >= 0` — if not, return `eval: factorial requires a non-negative integer`.
2. `operand == math.Trunc(operand)` — if not, return the same error.
3. Compute iteratively: `result = 1; for i = 2; i <= operand; i++ { result *= i }`.
4. `0!` **MUST** return `1`.

Overflow (large n) **MAY** produce `+Inf`; no error is required for Tier 1.

### 4.5 FunctionNode

**MUST** evaluate each element of `Args` recursively, then dispatch on `Name`
to the built-in function table (§6).

If `Name` is not in the table, **MUST** return `eval: unknown function: <name>`.

## 5. Built-in Constants (Scope seeds)

`NewScope()` **MUST** return a `Scope` containing at least:

| Key     | Value                               | LaTeX source |
| ------- | ----------------------------------- | ------------ |
| `pi`    | `math.Pi` (≈ 3.14159…)              | `\pi`        |
| `e`     | `math.E`  (≈ 2.71828…)              | `e`          |
| `tau`   | `2 * math.Pi` (≈ 6.28318…)          | `\tau`       |
| `phi`   | `(1 + math.Sqrt(5)) / 2` (≈ 1.618…) | `\phi`      |
| `infty` | `math.Inf(1)` (+∞)                  | `\infty`     |

`infty` **MUST** be `math.Inf(1)`. Per ADR-010, a `+Inf` result prints to
stdout and exits 0. Arithmetic on `infty` follows IEEE 754 (e.g.
`\infty + 1 = +Inf`, `1 / \infty = 0`).

## 6. Built-in Function Table

All function names are normalised (backslash-stripped, lowercased) as set by
the parser's `FunctionNode.Name`.

### 6.1 Arithmetic

| Name    | Arity | Implementation       | Domain error condition            |
| ------- | ----- | -------------------- | --------------------------------- |
| `frac`  | 2     | `args[0] / args[1]`  | `args[1] == 0` → division by zero |
| `dfrac` | 2     | alias for `frac`     | same                              |
| `tfrac` | 2     | alias for `frac`     | same                              |
| `cfrac` | 2     | alias for `frac`     | same                              |
| `sqrt`  | 1     | `math.Sqrt(args[0])` | `args[0] < 0` → domain error      |
| `abs`   | 1     | `math.Abs(args[0])`  | none                              |

`frac` division-by-zero **MUST** produce `eval: division by zero`.

`sqrt` domain error **MUST** produce `eval: domain error: sqrt of negative`.

`dfrac`, `tfrac`, `cfrac` are display-mode variants with identical numerical
semantics to `frac`.

### 6.2 Trigonometric

All angles are in **radians**.

| Name  | Implementation          | Notes                               |
| ----- | ----------------------- | ----------------------------------- |
| `sin` | `math.Sin(args[0])`     |                                     |
| `cos` | `math.Cos(args[0])`     |                                     |
| `tan` | `math.Tan(args[0])`     | ±Inf near poles — no error required |
| `sec` | `1 / math.Cos(args[0])` | ±Inf near poles — no error required |
| `csc` | `1 / math.Sin(args[0])` | ±Inf near poles — no error required |
| `cot` | `1 / math.Tan(args[0])` | ±Inf near poles — no error required |

### 6.3 Inverse Trigonometric

| Name     | Implementation         | Domain: error if…           |
| -------- | ---------------------- | --------------------------- |
| `asin`   | `math.Asin(args[0])`   | `\|args[0]\| > 1`           |
| `acos`   | `math.Acos(args[0])`   | `\|args[0]\| > 1`           |
| `atan`   | `math.Atan(args[0])`   | none                        |
| `asec`   | `math.Acos(1/args[0])` | `\|args[0]\| < 1`           |
| `acsc`   | `math.Asin(1/args[0])` | `\|args[0]\| < 1`           |
| `acot`   | `math.Atan(1/args[0])` | none                        |
| `arcsin` | alias for `asin`       | same domain checks          |
| `arccos` | alias for `acos`       | same domain checks          |
| `arctan` | alias for `atan`       | none                        |

`arcsin`, `arccos`, `arctan` are canonical LaTeX names. Their numerical result
and domain checks **MUST** be identical to `asin`, `acos`, `atan` respectively.

Domain error message **MUST** be `eval: domain error: <name> argument out of range`
where `<name>` is the name as it appeared in the expression (e.g. `arcsin`).

**Divergence from evaluatex:** evaluatex supports only `\asin`/`\acos`/`\atan`
and propagates `NaN` silently. goexeclatex supports both spellings and returns
explicit domain errors.

### 6.4 Logarithmic and Exponential

| Name  | Implementation        | Domain: error if… |
| ----- | --------------------- | ----------------- |
| `ln`  | `math.Log(args[0])`   | `args[0] <= 0`    |
| `log` | `math.Log10(args[0])` | `args[0] <= 0`    |
| `exp` | `math.Exp(args[0])`   | none              |

`\ln` is the natural logarithm (base e). `\log` is the base-10 logarithm.

Domain error **MUST** be `eval: domain error: <name> argument out of range`.

> **Divergence from evaluatex:** evaluatex does not expose `\ln`, `\log`, or
> `\exp` as LaTeX commands. goexeclatex follows standard LaTeX convention:
> `\ln` = natural log, `\log` = base-10 log. `\log_{b}(x)` subscript notation
> is deferred to Tier 0.3 (subscript support).

### 6.5 Hyperbolic

| Name   | Implementation           | Notes                              |
| ------ | ------------------------ | ---------------------------------- |
| `sinh` | `math.Sinh(args[0])`     | none                               |
| `cosh` | `math.Cosh(args[0])`     | none                               |
| `tanh` | `math.Tanh(args[0])`     | none                               |
| `coth` | `1 / math.Tanh(args[0])` | ±Inf at x=0 — no error required    |
| `sech` | `1 / math.Cosh(args[0])` | none (`cosh(x) ≥ 1` always)        |
| `csch` | `1 / math.Sinh(args[0])` | ±Inf at x=0 — no error required    |

`coth(0)` and `csch(0)` yield ±Inf per IEEE 754. Per ADR-010, `±Inf` is a
success result (exit 0).

> **Divergence from evaluatex:** evaluatex does not expose hyperbolic functions
> as LaTeX commands.

### 6.6 Combinatorics

| Name     | Arity | Implementation            | Domain error condition                   |
| -------- | ----- | ------------------------- | ---------------------------------------- |
| `binom`  | 2     | `n! / (k! * (n-k)!)`      | non-integer, negative, or `k > n`        |
| `dbinom` | 2     | alias for `binom`         | same                                     |
| `tbinom` | 2     | alias for `binom`         | same                                     |

`binom(n, k)` **MUST**:

1. Verify both `n` and `k` satisfy `v >= 0 && v == math.Trunc(v)`. If not:
   `eval: domain error: binom requires non-negative integer arguments`.
2. Verify `k <= n`. If not: `eval: domain error: binom requires k ≤ n`.
3. Compute via the multiplicative formula:
   `result = 1; for i in [0, k): result = result * (n-i) / (i+1)`.

`dbinom` and `tbinom` are display-mode variants; their result **MUST** be
identical to `binom`.

> **Divergence from evaluatex:** evaluatex does not support `\binom`.

### 6.7 Greek Letter Variables

Greek letters used as variables (`\alpha`, `\beta`, `\gamma`, …) lex as
arity-0 COMMAND tokens that normalise to `SymbolNode` with the backslash-
stripped, lowercased name (e.g. `\alpha` → `SymbolNode{Name: "alpha"}`).

These names are **not** seeded in `NewScope()`. The caller **MUST** provide
values via the scope. From the CLI, users supply values with `-v alpha=1.5`.

No implementation change is required; this behaviour is inherited from v0.1.

## 7. Error Conditions

The evaluator **MUST** return a non-nil error for each of the following:

| Condition                              | Error string                                                        |
| -------------------------------------- | ------------------------------------------------------------------- |
| Undefined symbol                       | `eval: undefined symbol: <name>`                                    |
| Division by zero (`/` or `frac`)       | `eval: division by zero`                                            |
| `sqrt` of negative                     | `eval: domain error: sqrt of negative`                              |
| Inverse trig out of range              | `eval: domain error: <name> argument out of range`                  |
| `ln` or `log` of non-positive          | `eval: domain error: <name> argument out of range`                  |
| Non-integer or negative factorial      | `eval: factorial requires a non-negative integer`                   |
| `binom` non-integer/negative args      | `eval: domain error: binom requires non-negative integer arguments` |
| `binom` k > n                          | `eval: domain error: binom requires k ≤ n`                          |
| Unknown function name                  | `eval: unknown function: <name>`                                    |

All error strings **MUST** begin with the prefix `eval: `.

## 8. Package Layout

```
internal/eval/
├── eval.go        Eval() — AST walker and dispatch
├── scope.go       Scope type and NewScope()
├── functions.go   built-in function implementations
└── eval_test.go   table-driven tests
```

## 9. Deferred (out of scope for v0.2)

Items below require parser changes or new syntax not yet supported. They are
deferred rather than implemented with non-standard syntax.

| Feature                                 | Target   | Blocker                                 |
| --------------------------------------- | -------- | --------------------------------------- |
| `\lfloor x \rfloor`, `\lceil x \rceil` | Tier 2.1 | New token types + paired-bracket parser |
| `\sqrt[n]{x}` n-th root                | Tier 2.1 | Optional `[n]` argument in parser       |
| `\min(a,b)`, `\max(a,b)`               | Tier 2.1 | Variadic paren-arg parser support       |
| `\gcd(a,b)`                            | Tier 2.1 | Variadic paren-arg parser support       |
| `\mod`, `\pmod`, `\bmod`, `\pod`       | Tier 2.1 | Binary infix operator parser support    |
| `\log_{b}(x)` subscript base           | Tier 0.3 | Subscript support                       |
| NaN propagation mode                   | vFuture  |                                         |
