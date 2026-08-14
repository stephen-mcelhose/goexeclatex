---
type: spike
title: evaluatex — Reference Implementation
description: Deep read of arthanzel/evaluatex, a JS LaTeX/ASCIIMath parser-evaluator. Documents architecture, covered features, and known limitations as the reference implementation for goexeclatex.
resource: https://github.com/arthanzel/evaluatex
tags: [reference-implementation, javascript, latex, math-parser, evaluatex]
timestamp: 2026-08-14T22:43:35Z
---

# evaluatex — Reference Implementation

[evaluatex](https://github.com/arthanzel/evaluatex) is a JavaScript/TypeScript
library that parses and numerically evaluates LaTeX and ASCIIMath expressions
without using `eval()`. It is the primary reference implementation for
[[goexeclatex-gap-analysis]].

Source: `~/repos/evaluatex` (cloned locally).

---

## Architecture

Three-phase pipeline: **Lexer → Parser → Evaluator**.

### Lexer (`src/lexer.js`)

Reads the input string into a `Token[]`. Two modes:

- **Normal mode** — greedy regex match against `Token.patterns` (ordered Map).
- **Char mode** (`lexExpression(charMode=true)`) — reads exactly one character at
  a time; activated after single-arg tokens like `^` and `\command` in LaTeX
  mode. This is how `2^24` correctly tokenizes as `2 ^ 2 4` rather than `2 ^ 24`
  when `{...}` are absent.

After tokenization, two post-passes:
1. `replaceConstants()` — promotes SYMBOL tokens that match known constants/functions into NUMBER or FUNCTION tokens at lex time.
2. `replaceCommands()` — strips leading `\` from COMMAND tokens and resolves them against the constants map (which includes `localFunctions`).

**Token patterns (in precedence order):**

| Type       | Pattern              |
| ---------- | -------------------- |
| LPAREN     | `[\(\[\{]` + `\left` |
| RPAREN     | `[\)\]\}]` + `\right`|
| PLUS       | `+`                  |
| MINUS      | `-`                  |
| TIMES      | `*`                  |
| DIVIDE     | `/`                  |
| COMMAND    | `\\[A-Za-z]+`        |
| SYMBOL     | `[A-Za-z_][A-Za-z_0-9]*` |
| WHITESPACE | `\s+`                |
| ABS        | `\|`                 |
| BANG       | `!`                  |
| COMMA      | `,`                  |
| POWER      | `^`                  |
| NUMBER     | `\d+(\.\d+)?`        |

### Parser (`src/parser.js`)

Recursive-descent parser. Grammar:

```
expression : sum
sum        : product { ('+'|'-') product }
product    : power { ('*'|'/') power }
           | power ('(' | SYMBOL | NUMBER | FUNCTION)   ← implicit multiplication
power      : val ['^' power]                             ← right-associative
val        : SYMBOL
           | NUMBER
           | COMMAND val{arity}                          ← arity from arities.js
           | FUNCTION ['(' sum {',' sum} ')' | power]   ← parens optional for 1-arg
           | '-' power
           | '(' sum ')'
           | '|' sum '|'
           | val '!'
```

### Evaluator (`src/Node.js` — `evaluate()`)

Depth-first AST traversal. Node types:

| Type     | Behaviour                          |
| -------- | ---------------------------------- |
| NUMBER   | returns `this.value`               |
| SYMBOL   | looks up `variables[this.value]`   |
| FUNCTION | calls `this.value(...children)`    |
| POWER    | `children[0] ** children[1]`       |
| PRODUCT  | `children.reduce(×)` where INVERSE wraps `/` |
| SUM      | `children.reduce(+)`               |
| NEGATE   | `-children[0]`                     |
| INVERSE  | `1 / children[0]`                  |

---

## Supported Features

### LaTeX Commands (hard-coded in `arities.js`)

| Command    | Arity | Implementation          |
| ---------- | ----- | ----------------------- |
| `\frac`    | 2     | `frac(a,b) = a/b`       |
| `\sqrt`    | 1     | `Math.sqrt`             |
| `\sin`     | 1     | `Math.sin`              |
| `\cos`     | 1     | `Math.cos`              |
| `\tan`     | 1     | `Math.tan`              |
| `\asin`    | 1     | `Math.asin`             |
| `\acos`    | 1     | `Math.acos`             |
| `\atan`    | 1     | `Math.atan`             |
| `\sec`     | 1     | `sec(x) = 1/cos(x)`     |
| `\csc`     | 1     | `csc(x) = 1/sin(x)`     |
| `\cot`     | 1     | `cot(x) = 1/tan(x)`     |
| `\asec`    | 1     | `Math.acos(1/x)`        |
| `\acsc`    | 1     | `Math.asin(1/x)`        |
| `\acot`    | 1     | `Math.atan(1/x)`        |
| `\times`   | —     | token → TIMES           |
| `\cdot`    | —     | token → TIMES           |
| `\left`    | —     | token → LPAREN (ignored) |
| `\right`   | —     | token → RPAREN (ignored) |

All of `Math.*` plus `fact`, `frac`, `logn`, `rootn`, `sec`, `csc`, `cot`
are merged into the default constants map.

### Built-in Constants

`pi`, `e`, `tau`, `phi`, `theta` (and all of `Math.*` that are numbers —
`Math.PI`, `Math.E`, `Math.LN2`, etc.)

### Other

- Implicit multiplication: `2x`, `(2)(3)`, `2(a+b)`
- Absolute value: `|expr|`
- Factorial: `expr!`
- Paren-less function application: `sin PI`
- Mixed brackets: `{1 + [2 - {3}]}` all treated as grouping
- Two-step compile/evaluate API

---

## Known Gaps (vs [[latex-math-evaluable-spec]])

| Feature                   | evaluatex | Notes                                      |
| ------------------------- | --------- | ------------------------------------------ |
| `\sqrt[n]{x}` n-th root   | ❌         | `rootn` function exists but no `[n]` syntax |
| `\arcsin/\arccos/\arctan` | ❌         | Has `\asin/\acos/\atan` only               |
| `\ln`, `\exp`             | ❌         | `Math.log` (= ln) exists but not as `\ln`  |
| `\log_{b}(x)` subscript   | ❌         | No subscript parsing at all                |
| `\lfloor \rfloor`         | ❌         | No floor/ceiling                           |
| `\lceil \rceil`           | ❌         | No ceiling                                 |
| `\min`, `\max`            | ❌         | `Math.min/max` in scope but not as commands |
| `\gcd`                    | ❌         |                                            |
| `\binom{n}{k}`            | ❌         |                                            |
| `\sinh/\cosh/\tanh`       | ❌         | `Math.sinh` etc. exist but not as commands  |
| `\sum_{i=a}^{b}`          | ❌         | No big operator support                    |
| `\prod_{i=a}^{b}`         | ❌         | No big operator support                    |
| `\infty`                  | ❌         | Not mapped                                 |
| `\pm`                     | ❌         |                                            |
| `\mod / \pmod`            | ❌         |                                            |
| `\dfrac / \tfrac`         | ❌         | Only `\frac`                               |
| Greek letters as vars      | ⚠️ partial | Only pi, e, tau, phi, theta built-in       |
| `\|...\|` norm syntax     | ❌         | Only single `|...|`                        |
| `x_{i}` subscript vars    | ❌         |                                            |

---

## Architecture Observations for Go Port

1. **Char mode** is an elegant hack for LaTeX's `^` single-char-arg rule. In Go this could be a parser state flag.
2. **`replaceToken`** (util) handles the `\left`/`\right` → paren mapping before the parser sees them. Good pattern to keep.
3. **Arity table** (`arities.js`) is a thin lookup — easy to replicate as a Go map.
4. **Implicit multiplication** is handled in the `product()` parser rule by checking the lookahead token type. No separate pass needed.
5. **The grammar has no subscript rule**. Adding subscripts (`_`) requires a new token type and grammar extension.

## Sources

- Source code: `~/repos/evaluatex/src/`
- Tests: `~/repos/evaluatex/test/evaluatex.spec.js`
- TODO: `~/repos/evaluatex/TODO.txt`
