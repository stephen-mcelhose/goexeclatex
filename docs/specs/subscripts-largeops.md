---
id: spec-subscripts-largeops
type: rfc2119-spec
title: Subscripts and large operators Specification
description: Normative RFC 2119 specification for v0.3 — SUBSCRIPT token, subscript variable resolution, discrete summation (\\sum) and product (\\prod) engines, and \\lVert norm support.
status: active
version: 0.3
tags: [subscript, bigop, sum, prod, norm, spec, goexeclatex]
timestamp: 2026-08-15T00:00:00Z
sources:
  - docs/plan.md                              # pipeline architecture, error policy
  - docs/evaluatex-reference-implementation.md  # grammar patterns; no subscript/bigop support in reference
  - docs/latex-math-evaluable-spec.md          # evaluable subset: §Subscripts, §large operators
  - docs/specs/lexer.md                        # token types, char-mode, post-lex passes
  - docs/specs/parser.md                       # grammar baseline being extended
  - docs/specs/eval.md                         # evaluator scope, error policy
  - docs/adrs/adr-004-lVert-deferred.md        # prior deferral of norm; now lifted for v0.3
  - docs/adrs/adr-013-drop-underscore-from-symbol.md  # why _ is now a free token
---

# Subscripts and large operators Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHALL NOT**, **SHOULD**,
**MAY**, and **OPTIONAL** in this document are to be interpreted as described
in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the v0.3 extensions to `internal/lexer`,
`internal/parser`, and `internal/eval` that introduce:

1. The `SUBSCRIPT` token (`_`)
2. Subscript variable resolution (`x_{i}`)
3. The discrete summation engine (`\sum_{i=a}^{b} expr`)
4. The discrete product engine (`\prod_{i=a}^{b} expr`)
5. The `\lVert ... \rVert` norm form (scalar absolute value)

> **Divergence from evaluatex:** The evaluatex reference implementation has
> no subscript token, no large-operator support, and no `\lVert` support. All
> of v0.3 is genuinely new ground — reference behaviours are not applicable.
> The canonical authority for semantics is the LaTeX math evaluable subset
> (see `docs/latex-math-evaluable-spec.md §Subscripts` and `§large operators`).

## 3. Lexer Extensions

### 3.1 SUBSCRIPT Token

The lexer **MUST** recognise a bare `_` character (U+005F) as a new token
type `SUBSCRIPT` with value `"_"`. This token **MUST** be added to the
pattern table at the **same priority position as `POWER`** (`^`), i.e. after
`BANG` and before `NUMBER`.

> **Rationale:** ADR-013 removed `_` from the SYMBOL pattern. The character
> is now free for use as a structural operator. `\_` is not handled specially;
> the `\` character requires letters to form a valid `COMMAND` token
> (`\\[A-Za-z]+`), so `\_` produces a lexer "unexpected character" error
> naturally without any additional logic.

### 3.2 EQUALS Token

The lexer **MUST** recognise `=` (U+003D) as a new token type `EQUALS` with
value `"="`. This token is used exclusively inside large-operator bound
expressions (`\sum_{i=a}^{b}`) and **MUST NOT** be valid in a general
expression context.

If an `EQUALS` token is encountered outside a large-operator bound, the parser
**MUST** return an error.

### 3.3 NORM Token (lVert / rVert)

The lexer **MUST** recognise `\lVert` and `\rVert` (capital `V`,
case-sensitive) as `NORM` tokens with value `"\|\|"`.

> **Token type:** `NORM` is distinct from `PIPE`, so the parser can
> distinguish `\lVert ... \rVert` from `\lvert ... \rvert` / `| ... |`.

`\lVert` and `\rVert` **MUST** be matched as **dedicated patterns** in the
initial tokenisation phase, placed **before** the `COMMAND` pattern in the
pattern table. They **MUST NOT** enter the COMMAND→normalise→remap pipeline.

> **Rationale:** The existing normalisation pass (lexer spec §7.3) lowercases
> all COMMAND values. `\lVert` → `lvert` after lowercasing, which collides
> with the existing `\lvert` → `PIPE` remap. Matching as dedicated patterns
> bypasses the collision entirely and keeps the 3-pass pipeline unmodified.

### 3.4 Char-Mode Behaviour for SUBSCRIPT

After emitting a `SUBSCRIPT` token, the lexer **MUST** enter char mode and
read exactly **1** argument, using the same char-mode rules defined in lexer
spec §6. This mirrors the existing behaviour for `POWER`.

The arity of `_` for char-mode purposes is **1**.

## 4. Parser Grammar Extensions

The parser **MUST** extend the grammar defined in `docs/specs/parser.md §4`
as follows.

### 4.1 Subscript Rule

A new `subscript` production **MUST** be inserted between `power` and
`atom`:

```
power      → subscript [POWER power]   // unchanged; subscript replaces atom as base
subscript  → atom [SUBSCRIPT sub_arg]
sub_arg    → NUMBER
           | SYMBOL
           | LPAREN('{') sum RPAREN('}')
```

> **Position in precedence stack:** Subscript binds tighter than power. This
> reflects standard LaTeX behaviour where `x_i^2` means `(x_i)^2`, not
> `x_(i^2)`. The subscript is part of the base, not the exponent.

The subscript production **MUST** be non-recursive (no chained subscripts):
`x_{i}_{j}` **MUST** produce a parse error.

**Order constraint:** When both `_` and `^` modify the same atom, the
subscript **MUST** appear first. `x_i^2` is valid and parses as
`(x_i)^2`. `x^2_i` **MUST** produce a parse error:

```
parser: superscript before subscript — write x_{i}^{2} not x^{2}_{i}
```

> **Divergence from standard LaTeX:** LaTeX allows `_` and `^` in either
> order after a base (`x^2_i` and `x_i^2` are equivalent). This
> implementation requires subscript-before-superscript because the recursive-
> descent grammar cannot unambiguously resolve both orders without
> backtracking. The constraint is documented rather than silently misparssed.

### 4.2 large operator Rule

When the parser encounters `COMMAND("sum")` or `COMMAND("prod")` in the
`atom` production, it **MUST** parse a **large-operator expression** rather than
a normal command-with-args.

The large-operator production is:

```
bigop → COMMAND("sum"|"prod")
        SUBSCRIPT LPAREN('{') SYMBOL EQUALS sum RPAREN('}')
        POWER super_arg
        expression
```

where `super_arg` is identical to `sub_arg` (§4.1).

Both bounds are **required** for numeric evaluation. A `\sum` or `\prod`
without bounds **MUST** produce a parse error. Bounds **MUST** appear in
`_{lower}^{upper}` order; `^{upper}_{lower}` order **MUST** produce a parse
error:

```
parser: expected '_{' after \sum/\prod — both bounds required; write \sum_{i=a}^{b}
```

> **Divergence from standard LaTeX:** Standard LaTeX permits `\sum` with only
> one bound or no bounds (typographic use). For numeric evaluation, both
> bounds are required because an open-ended sum has no finite value. The
> fixed ordering `_{...}^{...}` is a deliberate implementation simplification
> consistent with how learners write summations.

The components are:

| Component | Description |
| --------- | ----------- |
| `SYMBOL` | The iteration variable (bound within the body). |
| `sum` (lower bound) | Expression evaluated once to yield the start value. |
| `super_arg` (upper bound) | Expression evaluated once to yield the end value. |
| `expression` (body) | Evaluated once per integer step of the iteration variable. |

The iteration variable **MUST** be a single `SYMBOL` token. An expression in
the variable position **MUST** produce a parse error.

The lower and upper bound expressions **MUST NOT** reference the iteration
variable (they are evaluated in the outer scope, before iteration begins).

The body **MUST** be parsed as a full `expression` at the current precedence
level. The body **MAY** reference the iteration variable.

> **Precedence of body:** The body expression consumes the remainder of the
> current expression at `sum` level. This means `\sum_{i=1}^{3} i + 1` parses
> as `(\sum_{i=1}^{3} i) + 1 = 7`, not `\sum_{i=1}^{3} (i+1) = 9`. Users
> **MUST** use braces to group the intended body: `\sum_{i=1}^{3} {i+1}`.
>
> **Rationale:** This is consistent with standard LaTeX semantic parsing
> where the large operator's scope is determined by explicit grouping, not
> lookahead.

### 4.3 Norm Rule

The `atom` production **MUST** include:

```
atom → ... | NORM sum NORM
```

The two `NORM` tokens **MUST** be matched as an opening/closing pair,
analogous to the `PIPE` rule for absolute value (parser spec §4.2).

A `NORM` token encountered without a matching closing `NORM` **MUST**
produce a parse error.

## 5. AST Node Types

### 5.1 SubscriptNode

```go
type SubscriptNode struct {
    Base Node  // the variable or expression being subscripted
    Sub  Node  // the subscript argument
}
```

`SubscriptNode` **MUST** be produced by the `subscript` rule whenever a
`SUBSCRIPT` token follows an `atom`.

### 5.2 LargeOpNode

```go
type LargeOpNode struct {
    Op   string // "sum" or "prod"
    Var  string // iteration variable name
    From Node   // lower bound expression
    To   Node   // upper bound expression
    Body Node   // body expression
}
```

### 5.3 NormNode

```go
type NormNode struct {
    Arg Node
}
```

`NormNode` **MUST** be produced by the `NORM sum NORM` rule.

## 6. Evaluator Behaviour

### 6.1 SubscriptNode — Variable Resolution

To evaluate `SubscriptNode{Base, Sub}`:

1. `Base` **MUST** be a `SymbolNode`. If `Base` is any other node type, the
   evaluator **MUST** return an error:
   ```
   eval: subscript base must be a symbol, got <type>
   ```

2. Evaluate `Sub` to yield a float64 `v`. If evaluation fails, propagate the
   error.

3. If `v` is not a non-negative integer, the evaluator **MUST** return an
   error:
   ```
   eval: subscript index must be a non-negative integer, got <v>
   ```
   The integer check **MUST** use an epsilon tolerance to account for
   floating-point rounding in upstream arithmetic:
   ```go
   rounded := math.Round(v)
   if math.Abs(v-rounded) > 1e-9 || rounded < 0 {
       // error
   }
   v = rounded
   ```

4. Form the composite key: `key = Base.Name + "_" + strconv.Itoa(int(v))`.
   For example, `x_{0}` → key `"x_0"`.

5. Look up `key` in the scope. If not found, return:
   ```
   eval: undefined symbol: <key>
   ```

> **Variable binding:** Users supply subscript variables via `-v x_0=3.14`.
> The key format is `<name>_<index>` with no spaces.

### 6.2 LargeOpNode — Iteration Engine

To evaluate `LargeOpNode{Op, Var, From, To, Body}`:

1. Evaluate `From` in the current scope to yield `from` (float64). If
   `from` is not an integer, return:
   ```
   eval: sum/prod lower bound must be an integer, got <from>
   ```

2. Evaluate `To` in the current scope to yield `to` (float64). If `to` is
   not an integer, return:
   ```
   eval: sum/prod upper bound must be an integer, got <to>
   ```

3. If `to < from`, the result **MUST** be the identity value for the
   operation: `0` for `sum`, `1` for `prod`. No iterations are performed.

4. Create an inner scope that inherits all bindings from the outer scope
   and adds the iteration variable. The outer scope **MUST NOT** be mutated.

5. For each integer `i` from `int(from)` to `int(to)` **inclusive**, stepping
   by 1:
   a. Bind `Var = float64(i)` in the inner scope (shadowing any existing
      binding).
   b. Evaluate `Body` in the inner scope to yield `val`. If evaluation fails,
      propagate the error immediately.
   c. Accumulate: `\sum` adds `val` to the accumulator; `\prod` multiplies.

6. Return the accumulator.

> **Bounds constraint:** The specification intentionally limits to discrete
> integer bounds. `\sum_{i=0}^{\infty}` (infinite upper bound) **MUST**
> produce an error:
> ```
> eval: sum/prod upper bound must be a finite integer
> ```
> `\infty` evaluates to `math.Inf(1)`, which is not an integer, so §6.2.2
> already handles this case.

> **Iteration limit:** Implementations **SHOULD** impose a maximum iteration
> count (e.g. 10,000 steps) to prevent runaway evaluation. When the limit is
> exceeded the evaluator **MUST** return:
> ```
> eval: sum/prod iteration limit exceeded (<n> steps)
> ```

### 6.3 NormNode — Scalar Absolute Value

For a scalar expression, the norm equals the absolute value.

To evaluate `NormNode{Arg}`:

1. Evaluate `Arg` to yield `v`.
2. Return `math.Abs(v)`.

> **Rationale:** `\lVert x \rVert` for a scalar `x` is defined as `|x|`. This
> is the only case goexeclatex handles; vector/matrix norms are out of scope.

> **Divergence from ADR-004:** ADR-004 deferred `\lVert` support because the
> token-level distinction between single `|` and double `\|` was unresolved.
> v0.3 resolves this by introducing a distinct `NORM` token type (§3.3),
> making disambiguation unambiguous at the lexer level.

## 7. Error Conditions

| Condition | Phase | Error message |
| --------- | ----- | ------------- |
| `EQUALS` outside large-op bound | Parser | `parser: unexpected '=' outside iteration bound` |
| Chained subscripts `x_{i}_{j}` | Parser | `parser: chained subscripts are not supported` |
| Superscript before subscript `x^2_i` | Parser | `parser: superscript before subscript — write x_{i}^{2} not x^{2}_{i}` |
| large-op missing subscript bound | Parser | `parser: expected '_{' after \sum/\prod — both bounds required; write \sum_{i=a}^{b}` |
| large-op missing superscript bound | Parser | `parser: expected '^' after \sum/\prod bound` |
| large-op variable not a symbol | Parser | `parser: iteration variable must be a symbol` |
| Unmatched `\lVert` | Parser | `parser: unmatched \lVert` |
| Empty norm `\lVert\rVert` | Parser | `parser: empty \lVert...\rVert expression` |
| SubscriptNode base not a symbol | Eval | `eval: subscript base must be a symbol, got <type>` |
| Subscript index not a non-negative integer | Eval | `eval: subscript index must be a non-negative integer, got <v>` |
| Composite key not in scope | Eval | `eval: undefined symbol: <key>` |
| Lower bound not an integer | Eval | `eval: sum/prod lower bound must be an integer, got <v>` |
| Upper bound not an integer | Eval | `eval: sum/prod upper bound must be an integer, got <v>` |
| Upper bound is `\infty` | Eval | (covered by non-integer check above) |
| Iteration limit exceeded | Eval | `eval: sum/prod iteration limit exceeded (<n> steps)` |

All errors propagate as per the existing error policy in `docs/plan.md §Error
handling policy`. Parse errors exit 1; eval errors exit 2.

## 8. Constraints and Explicit Non-Scope

The following items are explicitly **out of scope** for this specification:

- **`\log_{b}(x)`** — shipped in v0.4; see [[specs/parser-extensions]] §6
  (no longer deferred here).
- **`x_{i}` where the subscript is a non-integer expression** — undefined;
  produces an eval error (§6.1).
- **Nested large operators** — `\sum_{i=0}^{n} \sum_{j=0}^{m} f(i,j)` is
  not explicitly prohibited but is expected to work naturally because the
  inner `\sum` is just the body expression of the outer one. No special
  handling is required.
- **Symbolic subscripts** — `x_{n-1}` where the subscript is an expression
  that evaluates to an integer IS in scope (§6.1 handles it via float→int
  conversion).
- **`\sum_{i=a}^{b}` with non-integer step** — step is always 1. No
  half-integer or floating-point stepping.
- **`\lVert v \rVert` for vector `v`** — out of scope; only scalar norm
  (= absolute value) is supported. (Scalar form shipped in v0.3; prior deferral
  [[adrs/adr-004-lVert-deferred]] is superseded.)

## 9. Implementation Notes

### 9.1 Inner Scope for large operators

The iteration variable **MUST** shadow outer bindings without mutating the
outer scope. The recommended implementation is a thin wrapper scope that
overrides a single key:

```go
type innerScope struct {
    outer Scope
    key   string
    val   float64
}
```

### 9.2 EQUALS Token Validity Guard

The `EQUALS` token is valid only inside the lower-bound position of a
large-op. The parser **MUST** check for stray `EQUALS` tokens in the general
`sum` rule and return an error, ensuring users get a clear message rather
than a confusing "unexpected token" if they accidentally write `=` in an
expression.

### 9.3 Norm vs Pipe Ambiguity

The `NORM` token is distinct from `PIPE`, so no depth-tracking (cf.
ADR-007) is needed for `\lVert ... \rVert`. The `absDepth` counter
introduced in ADR-007 applies only to `PIPE` tokens and is unaffected by
this change.
