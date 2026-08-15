---
id: adr-003
type: decision
title: "ADR-003: Use asin/acos/atan abbreviations (evaluatex convention) over arcsin/arccos/arctan (standard LaTeX)"
description: Adopts evaluatex's abbreviated asin/acos/atan/asec/acsc/acot names over standard LaTeX arcsin etc., for consistency with the reference implementation.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T00:00:00Z
tags: [adr, eval, trig, evaluatex]
---

# ADR-003: Use `asin`/`acos`/`atan` abbreviations over `arcsin`/`arccos`/`arctan`

## Status

Accepted

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

- Users who write standard LaTeX `\arcsin` will get a clear evaluator error
  rather than a wrong result.
- Tier 2 will add `arcsin`/`arccos`/`arctan` as aliases (see plan.md §v0.2).
- No special handling needed at the parser level — unknown commands produce
  `SymbolNode` which the evaluator resolves against the scope.
