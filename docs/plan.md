---
type: proposal
title: goexeclatex — Implementation Plan
description: Formal implementation plan for a Go CLI that numerically evaluates LaTeX math expressions read from stdin.
tags: [goexeclatex, plan, cli, go]
timestamp: 2026-08-14T23:00:00Z
---

# goexeclatex — Implementation Plan

## What it is

A Go CLI that reads a LaTeX math expression and prints its numeric result.

```
echo "\frac{1}{2} + \sqrt{9}" | goexeclatex
# => 3.5
```

Pipe-first design. No interactive mode, no streaming protocol — just stdin in, float64 out. Composable with shell pipelines.

---

## Architecture

Three-phase pipeline, mirroring evaluatex:

```
stdin → Lexer → []Token → Parser → AST → Evaluator → float64 → stdout
```

### Phase 1 — Lexer (`internal/lexer/`)

Reads a `string` (entire stdin, trimmed), emits `[]Token`.

**Escape rule:** `\\_` is a literal underscore within a symbol name, not a subscript trigger. The lexer must resolve `\\_` → `_` before the `_` token rule fires, so that `MY\\_CODE_{0}` lexes as `SYMBOL("MY_CODE")` `SUBSCRIPT` `NUMBER(0)`.

Token types:

| Type      | Pattern / source                              |
| --------- | --------------------------------------------- |
| NUMBER    | `\d+(\.\d+)?([eE][+-]?\d+)?`                 |
| SYMBOL    | `[A-Za-z_][A-Za-z_0-9]*`                     |
| COMMAND   | `\\[A-Za-z]+`                                 |
| PLUS      | `+`                                           |
| MINUS     | `-`                                           |
| TIMES     | `*`                                           |
| DIVIDE    | `/`                                           |
| POWER     | `^`                                           |
| LPAREN    | `(`, `[`, `{`, `\left` prefix                 |
| RPAREN    | `)`, `]`, `}`, `\right` prefix                |
| PIPE      | `\|`                                          |
| BANG      | `!`                                           |
| COMMA     | `,`                                           |
| UNDERSCORE| `_`  (Tier 3)                                 |
| EOF       | end of input                                  |

Post-lex passes:
1. Remap `\left` / `\right` → `LPAREN` / `RPAREN` (drop sizing info)
2. Remap `\times`, `\cdot`, `\div` → `TIMES` / `DIVIDE`
3. Resolve COMMAND tokens against the built-in function/constant table

**Char-mode rule (from evaluatex):** After `^` or a single-arg command, consume exactly one character or one `{...}` group. This ensures `2^24` parses as `2^2 * 4` (LaTeX single-char superscript), not `2^24`.

### Phase 2 — Parser (`internal/parser/`)

Recursive-descent. Grammar (precedence low → high):

```
expression → sum
sum        → product { ('+' | '-') product }
product    → power { ('*' | '/') power }
           | power  <implicit multiply: next token is LPAREN | SYMBOL | NUMBER | COMMAND>
power      → unary ['^' power]           // right-associative
unary      → '-' power
           | postfix
postfix    → atom ['!']
atom       → NUMBER
           | SYMBOL
           | COMMAND args                // arity from table
           | '(' sum ')'
           | '[' sum ']'
           | '{' sum '}'
           | '|' sum '|'                 // absolute value
           | '\lvert' sum '\rvert'       // preferred abs value form
```

Arity table (`internal/parser/arities.go`) drives how many `{...}` arguments each `\command` consumes.

### Phase 3 — Evaluator (`internal/eval/`)

Depth-first AST walk. Returns `(float64, error)`.

Symbol table (`Scope`) is a `map[string]float64` pre-seeded with built-in constants and extended with user-supplied vars via CLI flags.

---

## CLI interface (`cmd/goexeclatex/main.go`)

Built with [Cobra](https://github.com/spf13/cobra).

```
goexeclatex [flags]

Reads a LaTeX math expression from stdin, prints the result to stdout.

Flags:
  -e, --expr string       expression to evaluate (skips stdin read)
  -v, --var stringArray   bind a variable: -v x=3.14 (repeatable)
  -p, --prec int          decimal places in output (-1 = full precision)
  -h, --help              help for goexeclatex

Exit codes:
  0   success
  1   parse error
  2   evaluation error (e.g. division by zero, domain error)
```

Usage examples:

```sh
echo '\frac{22}{7}'                  | goexeclatex          # 3.142857...
echo '\sin(\pi / 6)'                 | goexeclatex          # 0.5
echo '2^{10}'                        | goexeclatex          # 1024
goexeclatex -e '\sqrt{2}'                                   # 1.4142...
goexeclatex -v x=3 -e 'x^2 + 2x + 1'                       # 16
```

---

## Package layout

```
goexeclatex/
├── cmd/
│   └── goexeclatex/
│       └── main.go          CLI entry point; reads stdin, calls Parse+Eval
├── internal/
│   ├── lexer/
│   │   ├── lexer.go         Lex(input string) ([]Token, error)
│   │   ├── token.go         Token type + TokenType enum
│   │   └── lexer_test.go
│   ├── parser/
│   │   ├── parser.go        Parse(tokens []Token, scope Scope) (Node, error)
│   │   ├── arities.go       map[string]int — command → arg count
│   │   ├── node.go          AST node types
│   │   └── parser_test.go
│   └── eval/
│       ├── eval.go          Eval(node Node, scope Scope) (float64, error)
│       ├── scope.go         Scope type, built-in seeds
│       ├── functions.go     built-in function implementations
│       └── eval_test.go
├── docs/
│   ├── plan.md              (this file)
│   ├── latex-math-evaluable-spec.md
│   ├── evaluatex-reference-implementation.md
│   └── goexeclatex-gap-analysis.md
├── go.mod
└── README.md
```

---

## Milestones

### v0.1 — Tier 1: evaluatex parity

**Goal:** all features the JS reference covers, in Go.

- [ ] `internal/lexer` — tokeniser + post-lex remaps
- [ ] `internal/parser` — recursive descent, arity table, implicit multiply
- [ ] `internal/eval` — AST evaluator, Scope, built-ins
- [ ] Built-in constants: `\pi`, `e`, `\tau`, `\phi`
- [ ] Trig: `\sin/\cos/\tan/\sec/\csc/\cot/\asin/\acos/\atan/\asec/\acsc/\acot`
- [ ] `\frac{a}{b}`, `\sqrt{x}`, `\times`, `\cdot`
- [ ] `|expr|` and `\lvert expr \rvert` absolute value
- [ ] `expr!` factorial (integer only, error otherwise)
- [ ] Implicit multiplication
- [ ] `\left` / `\right` grouping
- [ ] `cmd/goexeclatex` — stdin pipe, `-var`, `-prec`, `-e` flags
- [ ] Table-driven unit tests for all of the above

### v0.2 — Tier 2: function completeness

- [ ] `\arcsin/\arccos/\arctan` (canonical names)
- [ ] `\ln`, `\exp`, `\log` (base 10)
- [ ] `\dfrac`, `\tfrac`, `\cfrac` (aliases for `\frac`)
- [ ] `\sinh/\cosh/\tanh/\coth/\sech/\csch`
- [ ] `\min(a,b)`, `\max(a,b)` (variadic)
- [ ] `\infty` → `math.Inf(1)`
- [ ] `\lfloor x \rfloor`, `\lceil x \rceil`
- [ ] `\sqrt[n]{x}` — optional arg before brace
- [ ] `\gcd(a,b)` — Euclidean
- [ ] `\binom{n}{k}`, `\dbinom`, `\tbinom`
- [ ] `\mod`, `\pmod`, `\bmod`, `\pod`
- [ ] Greek letter variables (`\alpha`, `\beta`, … — user-supplied via `-var`)

### v0.3 — Tier 3: subscripts + big operators

- [ ] `\\_` → literal `_` pre-pass in lexer (before subscript tokenisation)
- [ ] `_` token; subscript grammar rule
- [ ] `\log_{b}(x)` — log base b
- [ ] `x_{i}` — subscript variable lookup
- [ ] `\sum_{i=a}^{b} f(i)` — discrete summation engine
- [ ] `\prod_{i=a}^{b} f(i)` — discrete product engine
- [ ] `\lVert v \rVert` — norm (double-pipe)

### vFuture

- `\begin{cases}` conditional evaluation
- Matrix input + `\det`
- JSON output mode (`-json`)

---

## Error handling policy

- **Parse errors** → print `error: <message> at position <n>` to stderr, exit 1
- **Eval errors** → print `error: <message>` to stderr, exit 2
  - Division by zero → `error: division by zero`
  - Domain error (e.g. `\sqrt{-1}`) → `error: domain error: sqrt of negative`
  - Undefined symbol → `error: undefined symbol: \foo`
  - Non-integer factorial → `error: factorial requires a non-negative integer`
- **Success** → result on stdout, newline-terminated, exit 0

---

## Testing strategy

- Table-driven `_test.go` files per package
- `eval_test.go` tests end-to-end (raw LaTeX string → float64) as integration
- Golden files for CLI output (`testdata/`)
- Reference: port the relevant cases from `evaluatex/test/evaluatex.spec.js`

---

## Out of scope (all versions)

- `\int`, `\lim`, `\partial` — symbolic ops
- `\sum_{i=1}^{\infty}` — infinite series
- Matrix environments
- Interactive / REPL mode
- Streaming protocol (use shell pipes instead)
