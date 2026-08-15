---
id: spec-cli
type: rfc2119-spec
title: CLI Specification
description: Normative RFC 2119 specification for the goexeclatex CLI — flag definitions, input resolution priority, variable binding, output format, exit codes, and error presentation.
status: active
version: 0.1
tags: [cli, spec, goexeclatex]
timestamp: 2026-08-15T03:04:00Z
sources:
  - docs/plan.md
  - docs/specs/eval.md
  - docs/adrs/adr-009-explicit-domain-errors-over-nan.md
  - docs/adrs/adr-010-inf-result-exits-zero.md
  - docs/adrs/adr-011-cli-error-prefix-stripping.md
---

# CLI Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHOULD**, and **MAY** are to
be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Invocation

```
goexeclatex [flags]
```

The tool **MUST** read a LaTeX math expression, evaluate it numerically, and
print the result to stdout.

## 3. Flags

| Flag               | Type         | Default | Description                                    |
| ------------------ | ------------ | ------- | ---------------------------------------------- |
| `-e`, `--expr`     | string       | `""`    | Expression to evaluate (skips stdin)           |
| `-v`, `--var`      | string array | `nil`   | Bind a variable: `-v name=value` (repeatable)  |
| `-p`, `--prec`     | int          | `-1`    | Decimal places in output (`-1` = full precision) |
| `-h`, `--help`     | —            | —       | Print help and exit 0                          |

## 4. Input Resolution

The tool **MUST** obtain the expression by the following priority:

1. If `-e` / `--expr` is non-empty, use its value (trimmed of leading/trailing whitespace).
2. Otherwise, read all of stdin, trim whitespace, and use that as the expression.
3. If the result is an empty string after trimming, **MUST** print
   `error: no expression provided (use stdin or -e)` to stderr and exit 1.

## 5. Variable Binding (`-v`)

Each `-v` value **MUST** have the form `name=value` where:

- `name` is trimmed of whitespace before insertion into the scope.
- `name` **MUST** have all `{` and `}` characters removed after trimming, so
  that `x_{0}` and `x_0` are equivalent keys. This allows the `-v` key to
  mirror LaTeX subscript notation used in the expression.
- `value` **MUST** be parseable by `strconv.ParseFloat(..., 64)`.

On format violation the tool **MUST** print
`error: invalid -v value "<arg>": expected name=value` to stderr and exit 1.

On unparseable value the tool **MUST** print
`error: invalid -v value "<arg>": <strconv error>` to stderr and exit 1.

User-supplied variables are merged into `eval.NewScope()` after the built-in
constants; they **MAY** override built-ins (e.g. `-v e=5` replaces Euler's
number).

## 6. Pipeline

The tool **MUST** obtain a numeric result by calling the public library API
([[specs/library]]):

```
expression → goexeclatex.Eval(expr, vars) → format → stdout
```

`Eval` runs `lexer.Lex` → `parser.Parse` → `eval.Eval` internally. The CLI
**MUST NOT** import `internal/{lexer,parser,eval}` directly.

Lex/parse failures (**MUST** be detectable as [[specs/library]] syntax-stage
errors) exit 1. Eval failures exit 2 (§8).
## 7. Output Format

On success the tool **MUST** print the result followed by a single newline to
stdout.

### 7.1 Full precision (default, `-p -1`)

**MUST** format with `strconv.FormatFloat(v, 'g', -1, 64)` — the shortest
decimal representation that round-trips, in fixed or scientific notation as
appropriate.

| Value            | Output                   |
| ---------------- | ------------------------ |
| `0.5`            | `0.5`                    |
| `1024.0`         | `1024`                   |
| `math.Sqrt(2)`   | `1.4142135623730951`     |
| `math.Pi`        | `3.141592653589793`      |
| `math.Inf(1)`    | `+Inf`                   |
| `math.Inf(-1)`   | `-Inf`                   |

### 7.2 Fixed precision (`-p N`, N ≥ 0)

**MUST** format with `strconv.FormatFloat(v, 'f', N, 64)` — exactly N decimal
places.

| Value          | `-p 2` | `-p 4`    |
| -------------- | ------ | --------- |
| `0.5`          | `0.50` | `0.5000`  |
| `math.Pi`      | `3.14` | `3.1416`  |
| `math.Sqrt(2)` | `1.41` | `1.4142`  |

### 7.3 `±Inf`

`±Inf` results **MUST** be printed and exit 0 ([[adrs/adr-010-inf-result-exits-zero]]). They are valid IEEE
754 values — the evaluator succeeded.

### 7.4 `NaN`

`NaN` **SHALL NOT** appear in v0.1 output. All NaN-producing paths are caught
as domain errors (ADR-009) or are unreachable.

## 8. Exit Codes and Error Reporting

All error messages **MUST** be written to stderr. The format is:

```
error: <message>
```

Where `<message>` is the internal error string with the package prefix
(`eval: `, `lexer: `, `parser: `) stripped ([[adrs/adr-011-cli-error-prefix-stripping]]).

| Condition              | Exit | Example stderr                                      |
| ---------------------- | ---- | --------------------------------------------------- |
| Success                | 0    | —                                                   |
| Empty input            | 1    | `error: no expression provided (use stdin or -e)`   |
| Bad `-v` format        | 1    | `error: invalid -v value "x": expected name=value`  |
| Lex error              | 1    | `error: unclosed brace group`                       |
| Parse error            | 1    | `error: unexpected token ...`                       |
| Eval error             | 2    | `error: division by zero`                           |
| Undefined symbol       | 2    | `error: undefined symbol: x`                        |
| Domain error           | 2    | `error: domain error: sqrt of negative`             |

**Note:** Position information (`at position N`) is a v0.2 enhancement.

## 9. Examples

```sh
printf '%s\n' '\frac{1}{2} + \sqrt{9}' | goexeclatex          # 3.5
printf '%s\n' '\sin{\pi/2}'            | goexeclatex          # 1
printf '%s\n' '2^{10}'                 | goexeclatex          # 1024
goexeclatex -e '\sqrt{2}'                                 # 1.4142135623730951
goexeclatex -e '\frac{1}{2} + \sqrt{9}'                   # 3.5
goexeclatex -e '\sqrt{2}' -p 4                            # 1.4142
goexeclatex -v x=3 -e 'x^2 + 2*x + 1'                    # 16
goexeclatex -v x=3 -v y=4 -e '\sqrt{x^2+y^2}'            # 5
```

> **Shell note:** On zsh, `echo '\frac{…}'` interprets `\f` as a form feed.
> Prefer `printf '%s\n' '…'` or `goexeclatex -e '…'` (see [[how-to]]).

## Sources

- [[plan]], [[specs/eval]]
- [[adrs/adr-009-explicit-domain-errors-over-nan]], [[adrs/adr-010-inf-result-exits-zero]], [[adrs/adr-011-cli-error-prefix-stripping]]
