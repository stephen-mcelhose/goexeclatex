---
id: adr-017
type: decision
title: "ADR-017: Public library API at module root; AST stays internal"
description: >
  Locks #20 packaging decisions: module-root package goexeclatex with string
  Eval, typed stage errors, no AST/Lex/Parse export, no float precision API;
  CLI remains a thin client that formats with -p and maps errors to exit codes.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T18:20:00Z
tags: [adr, library, public-api, packaging, errors, issue-20]
---

# ADR-017: Public library API at module root; AST stays internal

## Status

Accepted — records Gate A answers on
[#20](https://github.com/stephen-mcelhose/goexeclatex/issues/20).

## Context

All evaluation lived under `internal/{lexer,parser,eval}`, so external Go
modules could not import the engine. The README already claimed a “library.”
Issue #20 asks for a stable, importable facade without rewriting internals.

Design forks considered:

| Fork | Options |
| ---- | ------- |
| Layout | Module-root `package goexeclatex` vs `pkg/goexeclatex` |
| Precision | Float-rounding helper vs omit vs string `Format` matching CLI `-p` |
| Errors | String prefixes only vs typed/stage-discriminable errors |
| Surface | String `Eval` only vs also export AST / Lex / Parse |

Community guidance ([Organizing a Go module](https://go.dev/doc/modules/layout)):
prefer a root library package; use `internal/` for private engine code. `pkg/`
is a public-folder convention with no toolchain force — not a privacy boundary
(privacy is `internal/` only).

CLI `-p` formats with `strconv.FormatFloat` (`'g'` / `'f'` decimal places).
Rounding a returned `float64` to “N decimals” is not the same and can hurt
usable numeric accuracy.

## Decision

1. **Public package** — module-root `package goexeclatex`, import path
   `github.com/stephen-mcelhose/goexeclatex`.
2. **API (v1)** — `Eval(expr string, vars map[string]float64) (float64, error)`
   runs Lex → Parse → Eval via `internal/*`. `vars == nil` means no user
   bindings; builtins still seeded (`eval.NewScope` policy, same as CLI).
3. **Errors** — export stage-discriminable typed errors (at least
   syntax/lex-parse vs eval) so library callers and the CLI can distinguish
   stages. CLI exit 1 vs 2 wiring uses `errors.As` in the same milestone.
4. **Non-export** — AST/`Node`, lexer tokens, and parser/lexer entry points
   stay in `internal/`. The public surface is a facade, not a compiler API.
5. **Precision** — **MUST NOT** ship a float-returning precision/rounding
   helper. Display/`-p` stays CLI-local (`formatResult`). A public string
   `Format` matching CLI §7 **MAY** be added later without changing `Eval`.
6. **Engine location** — no duplication; shim only. Optional defense-in-depth
   arity guards for `min`/`max`/`gcd` may land with #20.

Normative detail: [[specs/library]].

## Consequences

- Semver applies to the root package API once published.
- Expanding to AST/compile-then-eval requires a new issue and an ADR that
  supersedes or amends this one.
- README / AGENTS / plan package trees must describe the library + CLI split.
- evaluatex’s two-step compile/evaluate API remains a documented divergence
  (we ship one-shot string `Eval`).

## Revisit criteria

Reopen or amend when any of:

1. Callers need a stable AST or partial pipeline (Lex/Parse only).
2. A public `Format` (or other display helper) is required for embedders to
   match CLI output without copying `formatResult`.
3. Finer error taxonomy (separate lex vs parse types, positions, codes) is
   needed beyond stage discrimination.

## Sources

- GitHub issue [#20](https://github.com/stephen-mcelhose/goexeclatex/issues/20)
- `.agents/issues/issue-20/{clarification,breakdown,explore}.md`
- [Organizing a Go module](https://go.dev/doc/modules/layout)
- [[specs/cli]] — `-p` / exit codes
- [[adrs/adr-011-cli-error-prefix-stripping]]
- [[evaluatex-reference-implementation]] — compile/evaluate API
