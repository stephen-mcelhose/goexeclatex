---
id: spec-lexer
title: Lexer Specification
status: active
version: 0.1
sources:
  - docs/plan.md                              # token table, pre-pass, post-lex passes
  - docs/evaluatex-reference-implementation.md  # char-mode algorithm (elevated to normative)
---

# Lexer Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHALL NOT**, **SHOULD**,
**MAY**, and **OPTIONAL** in this document are to be interpreted as described
in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the behaviour of `Lex(input string) ([]Token, error)`
in `internal/lexer`.

The lexer transforms a UTF-8 LaTeX math string into a flat `[]Token` slice
terminated by an `EOF` token.  It operates exclusively in **LaTeX mode**; there
is no ASCII mode.

## 3. Token Types

The lexer **MUST** produce tokens of the following types and **MUST NOT**
produce any other type.

| Type         | Source pattern / origin                              | Tier |
| ------------ | ---------------------------------------------------- | ---- |
| `NUMBER`     | `\d+(\.\d+)?([eE][+-]?\d+)?`                        | 0.1  |
| `SYMBOL`     | `[A-Za-z_][A-Za-z_0-9]*`                            | 0.1  |
| `COMMAND`    | `\\[A-Za-z]+`, normalised (see §7)                   | 0.1  |
| `PLUS`       | `+`                                                  | 0.1  |
| `MINUS`      | `-`                                                  | 0.1  |
| `TIMES`      | `*`; also remapped from `\times`, `\cdot`            | 0.1  |
| `DIVIDE`     | `/`; also remapped from `\div`                       | 0.1  |
| `POWER`      | `^`                                                  | 0.1  |
| `LPAREN`     | `(`, `[`, `{`; also from `\left(`, `\left[`, `\left{` | 0.1 |
| `RPAREN`     | `)`, `]`, `}`; also from `\right)`, `\right]`, `\right}` | 0.1 |
| `PIPE`       | `\|`; also remapped from `\lvert`, `\rvert`          | 0.1  |
| `BANG`       | `!`                                                  | 0.1  |
| `COMMA`      | `,`                                                  | 0.1  |
| `UNDERSCORE` | `_`                                                  | 0.3  |
| `EOF`        | Synthetic end-of-input sentinel                      | 0.1  |

The `Value` field of a token **MUST** contain the raw matched text, except:

- `COMMAND` tokens **MUST** have their value normalised per §7.3 before being
  returned to the caller.
- `TIMES` and `DIVIDE` tokens that arose from a remap **MUST** carry the value
  `"*"` and `"/"` respectively.
- `PIPE` tokens that arose from a remap **MUST** carry the value `"|"`.

## 4. Pre-Pass: Escaped Underscore

Before tokenisation begins, the lexer **MUST** replace every occurrence of the
two-character sequence `\_` (backslash followed by underscore) in the input
string with a single `_` character.

This ensures that `\_` within a symbol name (e.g. `MY\_CODE`) is treated as a
literal underscore and does **not** trigger the `UNDERSCORE` subscript token
that a bare `_` would produce in Tier 0.3.

> **Source:** `docs/plan.md` §Phase 1 escape rule.

## 5. Pattern Matching

### 5.1 Ordering

The lexer **MUST** try patterns in the order listed in §5.2.  The **first**
pattern that matches at the current position **MUST** be used; subsequent
patterns **MUST NOT** be tried for that position.

### 5.2 Pattern Table

Patterns are anchored to the current position (equivalent to a `^` anchor).
Whitespace **MUST** be skipped before each pattern attempt.

| Priority | Type         | Anchored pattern                                     |
| -------- | ------------ | ---------------------------------------------------- |
| 1        | `LPAREN`     | `(?:\\left[(\[{]\|[(\[{])`                           |
| 2        | `RPAREN`     | `(?:\\right[)\]}]\|[)\]}])`                          |
| 3        | `PLUS`       | `\+`                                                 |
| 4        | `MINUS`      | `-`                                                  |
| 5        | `TIMES`      | `\*`                                                 |
| 6        | `DIVIDE`     | `/`                                                  |
| 7        | `POWER`      | `\^`                                                 |
| 8        | `PIPE`       | `\|`                                                 |
| 9        | `BANG`       | `!`                                                  |
| 10       | `COMMA`      | `,`                                                  |
| 11       | `COMMAND`    | `\\[A-Za-z]+`                                        |
| 12       | `SYMBOL`     | `[A-Za-z_][A-Za-z_0-9]*`                            |
| 13       | `NUMBER`     | `\d+(?:\.\d+)?(?:[eE][+\-]?\d+)?`                   |

> **Rationale for ordering:**
> `LPAREN`/`RPAREN` precede `COMMAND` so that `\left(` is consumed as a single
> `LPAREN` token rather than as `COMMAND(\left)` + `LPAREN(()`.
> `COMMAND` precedes `SYMBOL` so that `\sin` is not tokenised as `SYMBOL(sin)`.

### 5.3 Whitespace

The lexer **MUST** skip any sequence of Unicode whitespace characters
(`U+0009`, `U+000A`, `U+000D`, `U+0020`) before each token read.
Whitespace **MUST NOT** produce a token.

### 5.4 Unknown Characters

If no pattern matches at the current position, the lexer **MUST** return an
error of the form:

```
lexer: unexpected character <q> at position <n>
```

where `<q>` is the quoted character and `<n>` is its zero-based byte offset.

## 6. Char-Mode Rule (LaTeX Single-Character Arguments)

In LaTeX, the argument to `^` or a `\command` is either a single character or
a `{...}` group.  The lexer **MUST** enforce this by switching to **char mode**
after emitting a `POWER` or `COMMAND` token.

> **Source:** `docs/evaluatex-reference-implementation.md` §Lexer, char mode;
> algorithm derived from `evaluatex/src/lexer.js` `lexExpression(charMode)`.

### 6.1 Char-Mode Token Read

When reading in char mode, the lexer **MUST** behave as follows:

1. Skip whitespace.
2. If the next character is `\`, read a complete `COMMAND` token (greedy,
   not single-character).
3. Otherwise, read exactly **one byte** and match it against the pattern table.

### 6.2 Argument Counts

After emitting a `POWER` token, the lexer **MUST** read exactly **1** argument
in char mode.

After emitting a `COMMAND` token, the lexer **MUST** read exactly **N**
arguments in char mode, where N is the command's arity from the table in §6.4.
If the command is not in the arity table, N is **0** (no char-mode arguments
are consumed).

### 6.3 Group Expansion Within Char Mode

When a char-mode read produces an `LPAREN` token whose value is `"{"`, the
lexer **MUST** continue reading tokens in **normal mode** until it emits the
matching `RPAREN` token whose value is `"}"`.  Only after the closing `}` is
the char-mode argument considered complete.

This allows `2^{12}` to produce `NUMBER(2) POWER LPAREN({) NUMBER(12) RPAREN(})`
rather than treating `{` alone as the exponent.

### 6.4 Command Arity Table

| Command  | Arity |
| -------- | ----- |
| `frac`   | 2     |
| `dfrac`  | 2     |
| `tfrac`  | 2     |
| `cfrac`  | 2     |
| `sqrt`   | 1     |
| `sin`    | 1     |
| `cos`    | 1     |
| `tan`    | 1     |
| `asin`   | 1     |
| `acos`   | 1     |
| `atan`   | 1     |
| `sec`    | 1     |
| `csc`    | 1     |
| `cot`    | 1     |
| `asec`   | 1     |
| `acsc`   | 1     |
| `acot`   | 1     |

Commands not listed here have arity 0.

## 7. Post-Lex Remapping

After the tokenisation loop completes, the lexer **MUST** apply the following
three passes **in order** before appending the `EOF` token.

### 7.1 Pass 1 — Drop `\left` / `\right` Decorators

Any `COMMAND` token whose normalised name is `"left"` or `"right"` **MUST** be
removed from the token stream.

> These tokens arise only when `\left` or `\right` appear without a directly
> adjacent bracket (e.g. `\left\lvert`).  When adjacent (e.g. `\left(`), they
> are already consumed as a single `LPAREN` by pattern priority 1.

### 7.2 Pass 2 — Operator Remaps

| COMMAND name | Replacement type | Replacement value |
| ------------ | ---------------- | ----------------- |
| `times`      | `TIMES`          | `"*"`             |
| `cdot`       | `TIMES`          | `"*"`             |
| `div`        | `DIVIDE`         | `"/"`             |
| `lvert`      | `PIPE`           | `"\|"`            |
| `rvert`      | `PIPE`           | `"\|"`            |

### 7.3 Pass 3 — COMMAND Normalisation

Every remaining `COMMAND` token **MUST** have its `Value` set to the
backslash-stripped, lowercased form of the original matched text.

> Example: `\Sin` → `"sin"`, `\FRAC` → `"frac"`.

Passes 7.1 and 7.2 **MUST** use the normalised name for matching (i.e. they
**MUST** be applied after the backslash is stripped and the name is lowercased,
or equivalently, normalisation **MAY** be performed as a preliminary step within
the same pass).

## 8. EOF Token

The lexer **MUST** append a single `EOF` token as the last element of the
returned slice.  The `EOF` token's `Pos` field **MUST** be set to `len(input)`
(after the pre-pass substitution).

## 9. Error Conditions

| Condition                                    | Behaviour                         |
| -------------------------------------------- | --------------------------------- |
| No pattern matches at current position       | Return error (§5.4)               |
| Unexpected end of input inside a `{` group   | Return error                      |
| Unexpected end of input in char mode         | Return error                      |

On error the lexer **MUST** return a `nil` token slice and a non-nil `error`.
