---
id: adr-010
title: "ADR-010: ±Inf results print to stdout and exit 0"
status: accepted
date: 2026-08-15
supersedes: ~
related: docs/adrs/adr-009-explicit-domain-errors-over-nan.md
---

# ADR-010: ±Inf results print to stdout and exit 0

## Status

Accepted — resolves the first revisit criterion in ADR-009.

## Context

ADR-009 explicitly deferred the following question:

> **±Inf at trig poles policy.** When the CLI spec is written, decide whether
> `+Inf` should exit 0 (valid result) or exit 2 (evaluation error). If exit 2,
> extend the domain-check approach to reciprocal trig poles.

`±Inf` can appear from:
- Reciprocal trig at poles: `\sec{\pi/2}`, `\csc{0}`, `\cot{0}`
- Division of a non-zero by zero if we ever relax the `right == 0` guard
  (not applicable in v0.1)

Two options were considered:

**Option A — exit 0, print `+Inf`/`-Inf`**
The evaluator succeeded; it produced a mathematically defined IEEE 754 value.
The expression `\sec{\pi/2}` is ill-defined in the reals but produces `+Inf`
in floating-point arithmetic — the same result evaluatex would give. Treating
it as an error requires a guard on every reciprocal trig call, adds complexity,
and breaks parity with the reference implementation.

**Option B — exit 2**
Treats `±Inf` as an evaluation failure, consistent with the spirit of ADR-009
(prefer explicit errors over silent wrong answers). Requires checking
`math.IsInf(result, 0)` after each call and returning an error.

## Decision

**Option A.** `±Inf` results are printed to stdout and the tool exits 0.

Rationale:
- `±Inf` is a valid IEEE 754 value, not a malformed result. The evaluator
  behaved correctly.
- evaluatex parity: the reference implementation also propagates `±Inf`
  silently at trig poles (documented in ADR-009 §Consequences).
- Adding post-call `IsInf` guards across all reciprocal trig is scope creep
  for v0.1 and can be added in v0.2 behind a flag if needed.
- A downstream shell pipeline can detect `+Inf` in stdout if desired.

## Consequences

- `echo '\sec{\pi/2}' | goexeclatex` prints `+Inf` and exits 0.
- The spec §7.3 documents this explicitly so callers are not surprised.
- `NaN` is still treated as a bug (ADR-009); all NaN-producing paths are
  caught as domain errors or are unreachable in v0.1.

## Revisit criteria

If a future user survey or downstream tooling shows that `+Inf` silently
corrupting pipelines is a real problem, add `--strict` flag (exit 2 on `±Inf`)
and supersede this ADR.
