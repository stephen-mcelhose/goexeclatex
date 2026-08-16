---
type: spike
title: Spike — Matrix Math
description: Exploration charter for matrix LaTeX in goexeclatex — evaluable subset, gonum dependency, and usability against the scalar Eval contract. Findings TBD.
tags: [goexeclatex, spike, matrix, linear-algebra, gonum, evaluable]
timestamp: 2026-08-15T20:00:00Z
---

# Spike — Matrix Math

**Tracking:** [#24](https://github.com/stephen-mcelhose/goexeclatex/issues/24) — `SPIKE: matrix math — evaluable subset, gonum, usability`
**Parent:** [#7](https://github.com/stephen-mcelhose/goexeclatex/issues/7) (vFuture — matrix + `\det`)

Status: **open** — charter only; no recommendation yet.

## High-level (see #24)

| | |
| --- | --- |
| **Value prop** | Matrix exprs without distorting scalar path — via **separate library API + CLI subcommand** when needed |
| **Feasible?** | Scalar slice on `Eval`: yes. Matrix×matrix / bindings: yes as **API + subcommand**, not via `float64` |
| **Size** | Spike **S**; scalar-on-`Eval` **M**; API + subcommand **L–XL** |

**Design note:** if non-scalar, the intended separate surface is a **different public API** (beside `Eval`) **and** a **CLI subcommand** (root command stays scalar).

---

## Problem

goexeclatex evaluates LaTeX math to a **scalar** (`float64`). Matrix notation is common in LaTeX, and [[latex-math-evaluable-spec]] flags `\det` as evaluable *if* matrices exist, while listing matrix environments as out-of-scope. [[plan]] still lists matrix environments under “Out of scope (all versions)”. [[goexeclatex-gap-analysis]] parked both under vFuture.

This spike asks whether that boundary should move, and whether matrix support belongs on **`Eval`**, on a **separate API + CLI subcommand**, or nowhere.

---

## Hypotheses to test

1. **Scalar-only slice on `Eval` is enough** — parse matrix *literals* only as arguments to scalar ops (`\det`, maybe `vmatrix`), never as first-class values. Keeps `Eval` and the root CLI unchanged.
2. **Non-scalar matrix ops use a different surface: library API + CLI subcommand** — e.g. matrix eval beside ADR-017 `Eval`, and `goexeclatex matrix …` (name TBD); root command stays scalar.
3. **Matrix multiplication is not Eval-evaluable** — belongs on the matrix API/subcommand (or only inside a scalar wrapper like `\det(AB)` if we allow nested matrix exprs on `Eval`).
4. **gonum may be unnecessary for `\det` of small literals** — closed-form 2×2/3×3 (or a tiny local det) might beat a mat dependency; gonum becomes attractive if the matrix surface takes general n×n, multiply, or norms.
5. **Usability hinge** — a correct `\det` on `Eval` may be enough; a matrix API + subcommand only pays off if bindings / print story are clear.

---

## Exploration checklist

Copy answers into this page (and comment on #24) when done.

### Evaluable subset

| Construct | In / out / defer | Result type | Rationale |
| --------- | ---------------- | ----------- | --------- |
| `\det` of matrix literal | | scalar? | |
| `\begin{vmatrix}…\end{vmatrix}` | | scalar? | |
| Matrix × matrix (incl. juxtaposition) | | matrix? | |
| Trace / matrix norms | | | |
| Inverse / linear solve | | | |
| Matrix-valued variables (`-v` / `Eval` vars) | | | |
| Element access `A_{ij}` | | | |

### LaTeX surface

- [ ] Environments to support (or reject): `matrix`, `pmatrix`, `bmatrix`, `Bmatrix`, `vmatrix`, `Vmatrix`, `smallmatrix`, `array`
- [ ] Row/column rules: `&`, `\\`, empty cells, non-rectangular input → errors
- [ ] evaluatex coverage (expect none) — confirm

### gonum

- [ ] Need general dense mat ops? → lean gonum
- [ ] Only `\det` of small literals? → consider no new dependency
- [ ] Record module impact (size, API wrap vs leak, Go version)

### Usability

- [ ] Compare surfaces: scalar-on-`Eval` vs **matrix library API + CLI subcommand** (preserve [[adrs/adr-017-public-library-api]] `Eval` and root CLI)
- [ ] Subcommand sketch: name, flags, how matrix results print (plain text vs JSON), bindings
- [ ] Error mapping (shape, singular, domain) vs existing `EvalError` / new types

---

## Decision gate (spike exit)

Spike closes when #24 has:

1. Recommended scope **and surface** (scalar-on-`Eval` / **API + CLI subcommand** / won’t-do).
2. Dependency recommendation (gonum / none / defer).
3. Explicit updates needed for [[plan]], [[latex-math-evaluable-spec]], [[roadmap]].
4. Follow-up milestone issue **or** won’t-do — no implementation under #24.

---

## Sources

- [#24](https://github.com/stephen-mcelhose/goexeclatex/issues/24) — spike tracking issue
- [#7](https://github.com/stephen-mcelhose/goexeclatex/issues/7) — vFuture umbrella
- [[latex-math-evaluable-spec]] — `\det` / matrix rows
- [[plan]] — out-of-scope list
- [[goexeclatex-gap-analysis]] — matrix / `\det` deferred
- [[specs/library]] / [[adrs/adr-017-public-library-api]] — public value model
- [[raw/short-math-guide]] — AMS matrices section
- gonum mat (external): https://pkg.go.dev/gonum.org/v1/gonum/mat
