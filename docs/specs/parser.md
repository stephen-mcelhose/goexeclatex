---
id: spec-parser
title: Parser Specification
status: active
version: 0.1
sources:
  - docs/plan.md                              # grammar, error policy, milestones
  - docs/evaluatex-reference-implementation.md  # reference grammar and node types
  - docs/specs/lexer.md                        # token types and post-lex invariants
  - docs/latex-math-evaluable-spec.md          # evaluable subset catalogue
---

# Parser Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHALL NOT**, **SHOULD**,
**MAY**, and **OPTIONAL** in this document are to be interpreted as described
in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the behaviour of `Parse(tokens []Token) (Node, error)`
in `internal/parser`.

The parser transforms a flat `[]Token` slice (produced by the lexer) into an
Abstract Syntax Tree (AST) using a recursive-descent strategy that mirrors the
evaluatex reference implementation.

## 3. Input Requirements

### 3.1 Token Stream

The input **MUST** be a `[]Token` as produced by `internal/lexer.Lex`. The
slice **MUST** be terminated by an `EOF` token. If the slice is empty or does
not end with `EOF`, the parser **MUST** return an error.

### 3.2 Post-Lex Invariants

The parser **MAY** rely on the following invariants guaranteed by the lexer:

- `\left` / `\right` sizing decorators have been dropped; their brackets appear
  as plain `LPAREN` / `RPAREN` tokens.
- `\times`, `\cdot` → `TIMES`; `\div` → `DIVIDE`; `\lvert`, `\rvert` → `PIPE`.
- Every `COMMAND` value is backslash-stripped and lowercased (`"frac"`, `"sin"`).
- Every `{` group opened during lexing is matched by a closing `}` — unclosed
  groups are a lex-time error and will not reach the parser.

## 4. Grammar

The parser **MUST** implement the following grammar (low → high precedence):

```
expression → sum
sum        → product { (PLUS | MINUS) product }
product    → power { (TIMES | DIVIDE) power }
           | power power*          // implicit multiply — see §5
power      → unary [POWER power]   // right-associative
unary      → MINUS power
           | postfix
postfix    → atom [BANG]
atom       → NUMBER
           | SYMBOL
           | COMMAND args          // §6
           | LPAREN('(') sum RPAREN(')')
           | LPAREN('[') sum RPAREN(']')
           | LPAREN('{') sum RPAREN('}')
           | PIPE sum PIPE         // absolute value — §4.2
```

### 4.1 Associativity

- `+`, `-`, `*`, `/` **MUST** be left-associative.
- `^` **MUST** be right-associative, implemented via right-recursive `power`.

### 4.2 Absolute Value

`PIPE sum PIPE` **MUST** handle both literal `|` and the `\lvert` / `\rvert`
forms. Both are remapped to `PIPE` by the lexer (§3.2), so the parser sees only
`PIPE` tokens.

The result **MUST** be represented as `FunctionNode{Name: "abs", Args: [inner]}`.

**Divergence from evaluable subset spec:** `docs/latex-math-evaluable-spec.md`
also lists `\lVert…\rVert` (capital V, double-pipe norm). The lexer does **not**
remap `\lVert`/`\rVert` to `PIPE`; they arrive as COMMAND tokens. Support for
the norm form is **deferred to Tier 2**.

### 4.3 Bracket Matching

For each `LPAREN`, the parser **MUST** verify the closing `RPAREN` carries the
matching value:

| Opening value | Required closing value |
| ------------- | ---------------------- |
| `(`           | `)`                    |
| `[`           | `]`                    |
| `{`           | `}`                    |

A mismatch **MUST** produce a parse error (§8).

Brackets are transparent grouping — the parser **MUST** return the inner
expression node directly, with no wrapper node.

## 5. Implicit Multiplication

Implicit multiplication is resolved at the **`product`** rule. After parsing a
`power`, the parser **MUST** check whether the next token can begin a new
`power`. If so, it inserts an implicit `TIMES` and continues.

A token can begin a `power` if it is one of:

- `NUMBER`
- `SYMBOL`
- `COMMAND`
- `LPAREN` (any opening bracket)
- `PIPE`

**Examples:**

| Input        | Parsed as        | Note                                       |
| ------------ | ---------------- | ------------------------------------------ |
| `2x`         | `2 * x`          |                                            |
| `2(3+4)`     | `2 * (3+4)`      |                                            |
| `x y`        | `x * y`          | `xy` (no space) is one identifier (ADR-003 evaluatex convention) |
| `2\sin x`    | `2 * sin(x)`     |                                            |
| `2\|x\|`     | `2 * abs(x)`     | PIPE only starts implicit-multiply when absDepth == 0 (ADR-007) |

## 6. Command Arity Dispatch

When the parser encounters a `COMMAND` token it **MUST** look up the command
name in the arity table (§6.1) and consume that many arguments (§6.2).

### 6.1 Arity Table

| Command                                                                          | Arity | Node type    |
| -------------------------------------------------------------------------------- | ----- | ------------ |
| `frac`                                                                           | 2     | FunctionNode |
| `sqrt`                                                                           | 1     | FunctionNode |
| `sin`, `cos`, `tan`, `sec`, `csc`, `cot`                                        | 1     | FunctionNode |
| `asin`, `acos`, `atan`, `asec`, `acsc`, `acot`                                  | 1     | FunctionNode |
| `pi`, `tau`, `phi`                                                               | 0     | SymbolNode   |
| Any command not in the table                                                     | 0     | SymbolNode   |

**Divergence from evaluable subset spec:** Standard LaTeX spells inverse-trig as
`\arcsin`, `\arccos`, `\arctan`. goexeclatex (following evaluatex) uses the
abbreviated forms `\asin`, `\acos`, `\atan`. Inputs using `\arcsin` etc. will
be treated as arity-0 unknown commands and produce a SymbolNode; the evaluator
will then fail with "undefined symbol". This is intentional for Tier 1 and will
be addressed in Tier 2.

### 6.2 Argument Consumption

The lexer emits each `{…}` command argument as:

```
LPAREN('{')  <tokens for the argument>  RPAREN('}')
```

For each of the N required arguments the parser **MUST**:

1. Expect `LPAREN('{')` — return an error if not found.
2. Call `sum` to parse the argument expression.
3. Expect `RPAREN('}')` — return an error if not found.

**Char-mode single-token arguments:** When the lexer consumed a single-token
char-mode argument (e.g. `\sin x`), the argument arrives without braces — the
SYMBOL `x` appears directly after the COMMAND. If no `LPAREN('{')` follows a
command that expects an argument, the parser **MUST** fall back to calling `atom`
to consume exactly one token as the argument.

## 7. AST Node Types

All node types **MUST** be defined in `internal/parser/node.go` and implement
the `Node` interface:

```go
type Node interface {
    node() // unexported marker method
}
```

### 7.1 NumberNode

```go
type NumberNode struct {
    Value float64
    Pos   int
}
```

### 7.2 SymbolNode

Represents a variable or arity-0 constant command.

```go
type SymbolNode struct {
    Name string // normalised: backslash-stripped, lowercased
    Pos  int
}
```

### 7.3 BinaryNode

```go
type BinaryNode struct {
    Op    string // "+", "-", "*", "/", "^"
    Left  Node
    Right Node
}
```

### 7.4 UnaryNode

Negation (`-`) is prefix; factorial (`!`) is postfix.

```go
type UnaryNode struct {
    Op      string // "-" or "!"
    Operand Node
}
```

### 7.5 FunctionNode

Represents any arity ≥ 1 command, plus absolute value.

```go
type FunctionNode struct {
    Name string // normalised command name, e.g. "frac", "sqrt", "abs"
    Args []Node
    Pos  int
}
```

## 8. Error Conditions

The parser **MUST** return a non-nil error for each of the following:

| Condition                    | Error message pattern                                               |
| ---------------------------- | ------------------------------------------------------------------- |
| Empty slice / missing EOF    | `parser: invalid token stream`                                      |
| Unexpected token             | `parser: unexpected <TYPE> "<value>" at position <pos>`             |
| Unexpected EOF               | `parser: unexpected EOF`                                            |
| Mismatched bracket           | `parser: expected "<close>" at position <pos>, got "<actual>"`      |
| Missing command argument     | `parser: <name> expects {…} argument at position <pos>`             |

## 9. API

```go
// Parse transforms a lexer token stream into an AST.
// tokens must be the slice returned by lexer.Lex — terminated by EOF.
func Parse(tokens []lexer.Token) (Node, error)
```

## 10. Deferred (out of scope for v0.1)

| Feature                                     | Target |
| ------------------------------------------- | ------ |
| `\arcsin`, `\arccos`, `\arctan` canonical names | Tier 2 |
| `\lVert…\rVert` norm (double-pipe)          | Tier 2 |
| `\sqrt[n]{x}` optional argument            | Tier 2 |
| `\min(a,b)`, `\max(a,b)` variadic          | Tier 2 |
| Subscript `_` operator                      | Tier 3 |
| `\sum`, `\prod` big operators               | Tier 3 |
