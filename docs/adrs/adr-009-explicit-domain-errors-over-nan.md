---
id: adr-009
title: "ADR-009: Return explicit errors for domain violations instead of propagating NaN"
status: accepted
date: 2026-08-15
---

# ADR-009: Return explicit errors for domain violations instead of propagating NaN

## Status

Accepted

## Context

evaluatex silently propagates `NaN` when a function receives an out-of-domain
argument — e.g. `asin(2)` returns `NaN`, and any expression containing it
evaluates to `NaN` with no error raised. The caller receives a meaningless
result with no indication of what went wrong.

goexeclatex is a CLI evaluator: the user supplies an expression and expects
either a number or a clear error message. Silent `NaN` propagation means a
wrong answer exits 0, which violates the contract implied by the CLI's exit
codes (0 = success, 2 = evaluation error).

Affected functions for Tier 1:
- `sqrt(x)` where `x < 0`
- `asin(x)`, `acos(x)` where `|x| > 1`
- `asec(x)`, `acsc(x)` where `|x| < 1`

## Decision

The evaluator **MUST** check domain conditions explicitly and return a non-nil
`error` rather than allowing Go's `math` package to return `NaN`. The error
strings follow the plan §Error handling policy:

| Condition               | Error                                             |
| ----------------------- | ------------------------------------------------- |
| `sqrt(x < 0)`           | `eval: domain error: sqrt of negative`            |
| `asin`/`acos` `|x| > 1` | `eval: domain error: <name> argument out of range` |
| `asec`/`acsc` `|x| < 1` | `eval: domain error: <name> argument out of range` |

Reciprocal trig at poles (`sec`, `csc`, `cot`, `tan`) is **not** checked —
the result is `±Inf`, which is a valid IEEE 754 value and will be printed as
`+Inf`/`-Inf` by the CLI. This is consistent with evaluatex behaviour and
deferred to a future policy decision (see Revisit criteria below).

## Consequences

- Out-of-domain inputs produce `exit 2` and a clear error message rather than
  `exit 0` with `NaN` printed to stdout.
- The evaluator diverges from evaluatex for these cases; the divergence is
  documented in `docs/specs/eval.md §6.3`.
- `NaN` cannot appear in the output of a successful evaluation for Tier 1
  (all NaN-producing paths are either caught or unreachable).

## Revisit criteria

1. **`±Inf` at trig poles policy.** When the CLI spec is written, decide
   whether `+Inf` should exit 0 (valid result) or exit 2 (evaluation error).
   If exit 2, extend the domain-check approach to reciprocal trig poles.

2. **NaN propagation mode.** A future flag (e.g. `--nan-ok`) could suppress
   domain errors and allow `NaN` to propagate, for users who prefer IEEE 754
   semantics. If added, this ADR should be superseded.
