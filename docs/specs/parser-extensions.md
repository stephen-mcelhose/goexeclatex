---
id: spec-parser-extensions
type: rfc2119-spec
title: Parser extensions Specification (v0.4)
description: >
  Normative RFC 2119 specification for issue #8 / v0.4 — paired floor/ceil,
  optional nth-root argument, variadic paren arguments, log base b, and \bmod.
status: active
version: 0.4
tags: [parser, lexer, eval, floor, ceil, sqrt, min, max, gcd, log, bmod, spec, goexeclatex]
timestamp: 2026-08-15T16:12:00Z
sources:
  - docs/plan.md
  - docs/evaluatex-reference-implementation.md
  - docs/latex-math-evaluable-spec.md
  - docs/specs/lexer.md
  - docs/specs/parser.md
  - docs/specs/eval.md
  - docs/adrs/adr-007-absdepth-pipe-ambiguity.md
  - docs/adrs/adr-014-mixed-floor-ceil-round-out-of-scope.md
  - .agents/issues/issue-8/clarification.md
---

# Parser extensions Specification (v0.4)

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHALL NOT**, **SHOULD**,
**MAY**, and **OPTIONAL** in this document are to be interpreted as described
in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the v0.4 / issue #8 extensions:

1. Matched floor / ceiling delimiters (`\lfloor…\rfloor`, `\lceil…\rceil`)
2. Optional nth-root index (`\sqrt[n]{x}`)
3. Variadic comma-separated paren arguments (`\min`, `\max`, `\gcd`)
4. Logarithm with subscript base (`\log_{b}(x)`)
5. Binary modulo (`\bmod`)

> **Divergence from evaluatex:** evaluatex lacks floor/ceil, `\sqrt[n]`, named
> `\min`/`\max`/`\gcd` commands, `\log_{b}`, and `\bmod`. Behaviours here are
> grounded in [[latex-math-evaluable-spec]] and AMS short-math-guide shapes.

Locked planning decisions: milestone **v0.4**; min/max/gcd arity **n ≥ 2**;
unary paren form for arity-1 commands **SHOULD** remain valid; `\pmod` /
`\mod` / `\pod` deferred; mixed floor/ceil round out of scope ([[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]]).

## 3. Floor and ceiling

### 3.1 Lexer

The lexer **MUST** recognise the following as dedicated tokens **before** the
generic `COMMAND` pattern (same rationale as `NORM`):

| Input | Token type | Value (normative) |
| ----- | ---------- | ----------------- |
| `\lfloor` | `FLOOR` | `lfloor` |
| `\rfloor` | `FLOOR` | `rfloor` |
| `\lceil` | `CEIL` | `lceil` |
| `\rceil` | `CEIL` | `rceil` |

Values **MUST** be backslash-stripped (no leading `\`). Case **MUST** match
the LaTeX command names above (these commands are lowercase in AMS).

### 3.2 Parser

Opening `\lfloor` **MUST** parse as:

```
floor → FLOOR("lfloor") sum FLOOR("rfloor")
```

Opening `\lceil` **MUST** parse as:

```
ceil → CEIL("lceil") sum CEIL("rceil")
```

The AST **MUST** be a `FunctionNode` with `Name` `"floor"` or `"ceil"` and a
single argument (the inner sum), matching the abs → `FunctionNode("abs")`
pattern.

The parser **MUST** track `floorDepth` and `ceilDepth` analogous to
`absDepth` / `normDepth` (ADR-007) so a closing delimiter is not consumed as
an implicit-multiply trigger.

Empty body (`\lfloor\rfloor`) **MUST** be a parse error.
Unmatched open **MUST** be a parse error.
Mixed pairs (`\lfloor…\rceil`, `\lceil…\rfloor`) **MUST** be a parse error
([[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]]); they **MUST NOT**
evaluate as round.

### 3.3 Eval

| Name | Semantics |
| ---- | --------- |
| `floor` | `math.Floor(arg)` |
| `ceil` | `math.Ceil(arg)` |

### 3.4 False-success regression

After this section is implemented, `\lfloor 3.2 \rfloor` **MUST NOT** parse as
an implicit-multiply `BinaryNode` product of symbols.

## 4. Optional nth-root (`\sqrt[n]{x}`)

### 4.1 Lexer

`\sqrt` **MUST NOT** consume char-mode arguments in the lexer (arity **0** in
the lexer arity table). The parser owns optional `[n]` and the radicand so that
`[` is not mistaken for the sole command argument.

> **Divergence from Tier 1:** Tier 1 listed `\sqrt` as arity 1 in both tables.
> v0.4 moves argument collection to the parser for `\sqrt` only.

### 4.2 Parser

```
sqrt → COMMAND("sqrt") [ LPAREN('[') sum RPAREN(']') ] command_arg
```

- If an optional `[…]` index is present, the AST **MUST** be
  `FunctionNode{Name: "sqrt", Args: [radicand, index]}` (radicand first).
- If absent, the AST **MUST** be `FunctionNode{Name: "sqrt", Args: [radicand]}`.
- `\sqrt{x}` and `\sqrt x` (single-token arg) **MUST** remain valid.

`\sqrt[n]{x}` **MUST NOT** parse as an implicit-multiply product of a one-arg
sqrt and a brace group.

### 4.3 Eval

| Args | Semantics |
| ---- | --------- |
| 1 (`x`) | `math.Sqrt(x)`; domain error if `x < 0` |
| 2 (`x`, `n`) | `math.Pow(x, 1/n)`; domain error if `n == 0`; domain error if `x < 0` |

### 4.4 False-success regression

`\sqrt[3]{27}` **MUST NOT** parse as `BinaryNode("*")`.

## 5. Variadic paren arguments (`\min`, `\max`, `\gcd`)

### 5.1 Lexer

`\min`, `\max`, and `\gcd` **MUST** have lexer arity **0** (no char-mode
expansion). The parser reads the paren list.

### 5.2 Parser

```
paren_args → LPAREN('(') sum { COMMA sum } RPAREN(')')
variadic_cmd → COMMAND("min"|"max"|"gcd") paren_args
```

- Arity **MUST** be **n ≥ 2**. Fewer than two arguments **MUST** be a parse
  error.
- The AST **MUST** be `FunctionNode` with `Name` equal to the command and
  `Args` the parsed list in order.
- Bare `COMMA` outside these lists **MUST** remain a parse error.

### 5.3 Unary paren form (arity-1 commands)

Existing arity-1 commands that accept a single parenthesized expression
(e.g. `\sin(x)`) **SHOULD** continue to parse successfully. Variadic lists
**MUST NOT** break that path.

### 5.4 Eval

| Name | Semantics |
| ---- | --------- |
| `min` | Minimum of all arguments |
| `max` | Maximum of all arguments |
| `gcd` | Greatest common divisor of all arguments; each argument **MUST** be a non-negative integer (same integer check spirit as `\binom`); pairwise Euclidean algorithm left-associative |

## 6. `\log_{b}(x)` — logarithm with subscript base

### 6.1 Lexer

`\log` **MUST** have lexer arity **0** so a following `_` is available to the
parser (not consumed as a char-mode argument).

### 6.2 Parser

```
log → COMMAND("log") [ UNDERSCORE sub_arg ] log_value
log_value → LPAREN('(') sum RPAREN(')')   // when subscript base present: required
          | command_arg                   // when no base: brace or single-token arg
```

- With base: AST **MUST** be `FunctionNode{Name: "log", Args: [value, base]}`.
- Without base: AST **MUST** be `FunctionNode{Name: "log", Args: [value]}`
  (base-10), preserving `\log{10}` / `\log 10`.
- When a subscript base is present, the value argument **MUST** be a
  parenthesized `(…)`. Brace-only `\log_{b}{x}` is **not** the supported form
  (Non-Scope).

### 6.3 Eval

| Args | Semantics |
| ---- | --------- |
| 1 (`x`) | `math.Log10(x)`; domain error if `x ≤ 0` |
| 2 (`x`, `b`) | `math.Log(x)/math.Log(b)`; domain error if `x ≤ 0`, `b ≤ 0`, or `b == 1` |

## 7. `\bmod` — binary modulo

### 7.1 Lexer

The lexer **MUST** remap COMMAND `\bmod` to token type `BMOD` with value
`"bmod"` in the post-lex remap pass (same pattern as `\times` → `TIMES`).

### 7.2 Parser

`BMOD` **MUST** be parsed at the **product** precedence level as a binary
operator (same band as `*` / `/`):

```
product → power { ('*' | '/' | BMOD) power | <implicit multiply> }
```

AST: `BinaryNode{Op: "bmod", Left, Right}`.

`a \bmod b` **MUST NOT** parse as implicit multiply of `a`, symbol `bmod`, and
`b`.

### 7.3 Eval

`bmod` **MUST** evaluate as `math.Mod(left, right)`. Division by zero (right
`== 0`) **MUST** return an eval domain / division error.

### 7.4 Deferred mod spellings

`\pmod`, `\mod`, and `\pod` are **out of scope** for v0.4 (AMS congruence
annotations, not binary remainder). They **MUST NOT** be remapped to `BMOD`.

## 8. Non-Scope

The following are **out of scope** for this specification:

- `\_` literal underscore pre-pass ([[adrs/adr-013-drop-underscore-from-symbol]])
- Mixed `\lfloor…\rceil` / `\lceil…\rfloor` as round ([[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]])
- `\pmod`, `\mod`, `\pod` (deferred; not binary remainder)
- Invented brace shims (`\min{a}{b}`, `\log_{b}{x}` as the supported form)
- vFuture items (#7): cases, matrices, `-json`

## 9. Sources

- [[plan]], [[roadmap]], [[latex-math-evaluable-spec]]
- [[evaluatex-reference-implementation]]
- [[adrs/adr-007-absdepth-pipe-ambiguity]], [[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]]
- GitHub issue [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)
