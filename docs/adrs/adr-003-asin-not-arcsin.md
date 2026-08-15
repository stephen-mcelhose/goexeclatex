---
id: adr-003
type: decision
title: "ADR-003: asin/acos/atan abbreviations for Tier 1; arcsin/arccos/arctan added as aliases in v0.2"
description: Tier 1 followed evaluatex's abbreviated asin/acos/atan convention. v0.2 added arcsin/arccos/arctan as canonical LaTeX aliases — both forms now work.
status: amended
date: 2026-08-15
timestamp: 2026-08-15T00:00:00Z
tags: [adr, eval, trig, evaluatex]
---

# ADR-003: `asin`/`acos`/`atan` for Tier 1; `arcsin`/`arccos`/`arctan` aliases added in v0.2

## Status

Amended (v0.2 — see Amendment below)

## Context

Standard LaTeX spells inverse trigonometric functions as `\arcsin`, `\arccos`,
`\arctan`. The evaluatex reference implementation uses abbreviated forms:
`\asin`, `\acos`, `\atan`, `\asec`, `\acsc`, `\acot`.

goexeclatex uses the lexer's `commandArities` table (and the parser's arity
table) as the single source of truth for recognised commands.

## Decision

Follow evaluatex for Tier 1: recognise `asin`, `acos`, `atan`, `asec`, `acsc`,
`acot` as arity-1 trig commands. The canonical LaTeX names (`arcsin` etc.) are
**not** in the arity table for Tier 1 and will be treated as arity-0 unknown
commands, causing the evaluator to fail with "undefined symbol".

This divergence is explicitly documented in `docs/specs/parser.md §6.1`.

## Consequences

- Users who write standard LaTeX `\\arcsin` will get a clear evaluator error
  rather than a wrong result (Tier 1 only — see Amendment).
- No special handling needed at the parser level — unknown commands produce
  `SymbolNode` which the evaluator resolves against the scope.

## Amendment — v0.2 (2026-08-15)

`\\arcsin`, `\\arccos`, `\\arctan` were added as arity-1 aliases in v0.2 alongside
the existing `\\asin`/`\\acos`/`\\atan` abbreviations. Both spellings are now valid.
`\\asec`, `\\acsc`, `\\acot` remain evaluatex-convention only (no `arc` prefix) since
those are non-standard in LaTeX anyway.

See [[specs/eval]] §6.3 for the normative function table.
