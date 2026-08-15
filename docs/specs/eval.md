---
id: spec-eval
type: rfc2119-spec
title: Evaluator Specification
description: Normative RFC 2119 specification for the goexeclatex evaluator — scope resolution, arithmetic and function dispatch, domain error handling, and IEEE 754 special-value policy.
status: active
version: 0.1
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

| Key    | Value                         | LaTeX source |
| ------ | ----------------------------- | ------------ |
| `pi`   | `math.Pi` (≈ 3.14159…)        | `\pi`        |
| `e`    | `math.E`  (≈ 2.71828…)        | `e`          |
| `tau`  | `2 * math.Pi` (≈ 6.28318…)    | `\tau`       |
| `phi`  | `(1 + math.Sqrt(5)) / 2` (≈ 1.61803…) | `\phi` |

## 6. Built-in Function Table

All function names are normalised (backslash-stripped, lowercased) as set by
the parser's `FunctionNode.Name`.

### 6.1 Arithmetic

| Name   | Arity | Implementation                | Domain error condition            |
| ------ | ----- | ----------------------------- | --------------------------------- |
| `frac` | 2     | `args[0] / args[1]`           | `args[1] == 0` → division by zero |
| `sqrt` | 1     | `math.Sqrt(args[0])`          | `args[0] < 0` → domain error      |
| `abs`  | 1     | `math.Abs(args[0])`           | none                              |

`frac` division-by-zero **MUST** produce `eval: division by zero`.

`sqrt` domain error **MUST** produce `eval: domain error: sqrt of negative`.

### 6.2 Trigonometric

All angles are in **radians**.

| Name   | Implementation          | Notes                   |
| ------ | ----------------------- | ----------------------- |
| `sin`  | `math.Sin(args[0])`     |                         |
| `cos`  | `math.Cos(args[0])`     |                         |
| `tan`  | `math.Tan(args[0])`     | ±Inf near poles — no error required |
| `sec`  | `1 / math.Cos(args[0])` | ±Inf near poles — no error required |
| `csc`  | `1 / math.Sin(args[0])` | ±Inf near poles — no error required |
| `cot`  | `1 / math.Tan(args[0])` | ±Inf near poles — no error required |

### 6.3 Inverse Trigonometric

| Name   | Implementation           | Domain: error if… |
| ------ | ------------------------ | ----------------- |
| `asin` | `math.Asin(args[0])`     | `\|args[0]\| > 1` → domain error |
| `acos` | `math.Acos(args[0])`     | `\|args[0]\| > 1` → domain error |
| `atan` | `math.Atan(args[0])`     | none              |
| `asec` | `math.Acos(1/args[0])`   | `\|args[0]\| < 1` → domain error |
| `acsc` | `math.Asin(1/args[0])`   | `\|args[0]\| < 1` → domain error |
| `acot` | `math.Atan(1/args[0])`   | none              |

Domain error message **MUST** be `eval: domain error: <funcname> argument out of range`.

**Divergence from evaluatex:** evaluatex propagates `NaN` silently for out-of-range
inverse trig inputs. goexeclatex returns an explicit error instead, following the
error policy in `docs/plan.md §Error handling policy`.

## 7. Error Conditions

The evaluator **MUST** return a non-nil error for each of the following:

| Condition                         | Error string                                              |
| --------------------------------- | --------------------------------------------------------- |
| Undefined symbol                  | `eval: undefined symbol: <name>`                         |
| Division by zero (`/` or `frac`)  | `eval: division by zero`                                 |
| `sqrt` of negative                | `eval: domain error: sqrt of negative`                   |
| Inverse trig argument out of range| `eval: domain error: <name> argument out of range`       |
| Non-integer or negative factorial | `eval: factorial requires a non-negative integer`        |
| Unknown function name             | `eval: unknown function: <name>`                         |

All error strings **MUST** begin with the prefix `eval: `.

## 8. Package Layout

```
internal/eval/
├── eval.go        Eval() — AST walker and dispatch
├── scope.go       Scope type and NewScope()
├── functions.go   built-in function implementations
└── eval_test.go   table-driven tests
```

## 9. Deferred (out of scope for v0.1)

| Feature                                  | Target |
| ---------------------------------------- | ------ |
| `\arcsin`, `\arccos`, `\arctan`          | Tier 2 |
| `\ln`, `\log`, `\exp`                    | Tier 2 |
| `\sinh`, `\cosh`, `\tanh`, `\coth`       | Tier 2 |
| `\lfloor`, `\lceil`                      | Tier 2 |
| `\min`, `\max`, `\gcd`                   | Tier 2 |
| `\infty`                                 | Tier 2 |
| NaN propagation mode (no domain errors)  | vFuture |
