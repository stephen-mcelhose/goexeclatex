---
id: spec-library
type: rfc2119-spec
title: Public Library API Specification
description: >
  Normative RFC 2119 specification for the importable goexeclatex package —
  Eval, variable bindings, built-in scope, typed stage errors, and non-scope
  (AST, float precision helpers).
status: active
version: 0.1
tags: [library, spec, goexeclatex, public-api, issue-20]
timestamp: 2026-08-15T18:20:00Z
sources:
  - docs/plan.md
  - docs/specs/cli.md
  - docs/specs/eval.md
  - docs/adrs/adr-017-public-library-api.md
  - docs/evaluatex-reference-implementation.md
  - docs/adrs/adr-009-explicit-domain-errors-over-nan.md
  - docs/adrs/adr-010-inf-result-exits-zero.md
  - docs/adrs/adr-011-cli-error-prefix-stripping.md
---

# Public Library API Specification

## 1. Key Words

The key words **MUST**, **MUST NOT**, **SHALL**, **SHOULD**, and **MAY** are to
be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. Scope

This specification defines the behaviour of the **public** Go package
`github.com/stephen-mcelhose/goexeclatex` (`package goexeclatex`).

The package is a thin facade over `internal/lexer`, `internal/parser`, and
`internal/eval`. Normative math/grammar behaviour remains in [[specs/lexer]],
[[specs/parser]], [[specs/eval]], [[specs/subscripts-largeops]], and
[[specs/parser-extensions]].

Packaging decisions are locked in [[adrs/adr-017-public-library-api]].

## 3. Relationship to the CLI

| Concern | Library | CLI (`cmd/goexeclatex`) |
| ------- | ------- | ----------------------- |
| Evaluate expression | `Eval` | **MUST** call `Eval` |
| `-e` / stdin / `-v` parsing | — | CLI-local |
| `-p` / stdout formatting | — | CLI-local ([[specs/cli]] §7) |
| Exit codes | — | CLI maps stage errors → 1 vs 2 ([[specs/cli]] §8) |
| ADR-011 prefix stripping | — | CLI-local user-facing stderr |

The library **MUST NOT** depend on Cobra or other CLI-only packages.

## 4. API

```go
package goexeclatex

// Eval parses and evaluates a LaTeX math expression.
// Variables in vars are merged onto the built-in scope before evaluation.
// A nil vars map MUST be treated as empty (no user bindings).
// Returns the full-precision float64 result or a stage-discriminable error.
func Eval(expr string, vars map[string]float64) (float64, error)
```

### 4.1 Pipeline

`Eval` **MUST** execute the same pipeline as the historical CLI engine path:

```
expr → lexer.Lex → parser.Parse → eval.NewScope() + merge vars → eval.Eval
```

`Eval` **MUST NOT** duplicate lexer/parser/eval logic; it **MUST** delegate to
the `internal/` packages.

### 4.2 Variable bindings

- `vars == nil` **MUST** be accepted and treated as an empty map.
- User bindings **MUST** be merged after built-in constants are seeded
  (same policy as [[specs/cli]] §5 / [[specs/eval]] §5): user keys **MAY**
  override built-ins.
- Key normalisation for library callers is the map key as provided; the CLI’s
  `{`/`}` stripping for `-v` names remains CLI-local.

### 4.3 Return value

- On success, `Eval` **MUST** return the `float64` produced by `internal/eval`
  without display rounding or truncation.
- ±Inf results **MUST** be returned as success with the Inf value (same policy
  as [[adrs/adr-010-inf-result-exits-zero]] for the CLI).
- Domain violations and other eval failures **MUST** return a non-nil error
  ([[adrs/adr-009-explicit-domain-errors-over-nan]]).

### 4.4 Empty expression

An expression that is empty after the same trimming rules the CLI applies to
`-e`/stdin **MAY** be rejected by the library with a non-nil error, or left to
the lexer/parser. Callers that mirror the CLI **SHOULD** reject empty input
before calling `Eval` (the CLI does).

**Chosen behaviour (v0.1):** `Eval("")` and whitespace-only expressions are
passed through unchanged; the parser returns unexpected EOF, wrapped as
`*SyntaxError`. This **MUST** be pinned by tests (`TestEvalEmptyExpression`).

## 5. Errors

### 5.1 Stage discrimination

Public errors **MUST** be distinguishable by stage so callers can map failures
without brittle new string protocols:

| Stage | Typical causes | CLI exit (when used by CLI) |
| ----- | -------------- | --------------------------- |
| Syntax (lex and/or parse) | Tokenisation or grammar failure | 1 |
| Eval | Domain error, undefined symbol, arity, … | 2 |

The implementation **MUST** provide typed errors (or an equivalent
`errors.As`-friendly discriminator) covering at least **syntax** vs **eval**.
Separate lex vs parse types **MAY** be provided; sharing one syntax type is
allowed for v0.1.

### 5.2 Message text

`error.Error()` strings **SHOULD** remain compatible with CLI presentation:
either preserve `lexer:` / `parser:` / `eval:` prefixes that
[[adrs/adr-011-cli-error-prefix-stripping]] strips, or ensure the CLI maps
typed errors to the same user-facing text as existing goldens.

### 5.3 Wrapping

Typed public errors **SHOULD** wrap the underlying internal error so
`errors.Is` / `errors.Unwrap` remain useful where applicable.

## 6. Defense-in-depth (eval builtins)

When `min`, `max`, or `gcd` are invoked with fewer than two evaluated numeric
arguments, `internal/eval` **MUST** return an error and **MUST NOT** panic.
(The parser already enforces arity for string `Eval`; this guard closes the
review note on [#20](https://github.com/stephen-mcelhose/goexeclatex/issues/20).)

Broader arity hardening for all builtins **MAY** be deferred.

## 7. Non-Scope

The public package **MUST NOT** export:

- AST types (`Node` and concrete node structs)
- Lexer or parser entry points (`Lex`, `Parse`) or token types
- A float-returning “precision” / rounding helper that alters the numeric
  result for display purposes

A public string formatter matching [[specs/cli]] §7 **MAY** be added in a
later revision without changing `Eval`.

Compile-then-evaluate (separate parse and run steps), as in evaluatex, is
**out of scope**; this API is one-shot string evaluation. That divergence
**MUST** remain documented here and in [[adrs/adr-017-public-library-api]].

## 8. Divergence from evaluatex

| evaluatex | goexeclatex public API |
| --------- | ---------------------- |
| Compile expression → reusable evaluator function | One-shot `Eval(expr, vars)` |
| JS package import | Go module-root import |

Semantic coverage of expressions remains driven by the internal specs, not by
matching evaluatex’s two-step API shape.

## 9. Conformance

Implementations of this package **MUST**:

1. Satisfy §4–§7.
2. Pass table-driven tests that cover happy-path arithmetic, `vars == nil`,
   user var override, syntax errors, and eval errors with stage checks.
3. Remain importable as `github.com/stephen-mcelhose/goexeclatex`.

## Sources

- [[plan]]
- [[specs/cli]], [[specs/eval]]
- [[adrs/adr-017-public-library-api]]
- [[evaluatex-reference-implementation]]
- GitHub issue [#20](https://github.com/stephen-mcelhose/goexeclatex/issues/20)
