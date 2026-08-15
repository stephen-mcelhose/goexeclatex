---
id: adr-004
type: decision
title: "ADR-004: Defer \\lVert…\\rVert norm support to Tier 2"
description: Defers double-bar norm syntax (\\lVert…\\rVert) to Tier 2, keeping Tier-1 scope to single-bar absolute value only.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T00:00:00Z
tags: [adr, lexer, deferred]
---

# ADR-004: Defer `\lVert…\rVert` norm support to Tier 2

## Status

Accepted

## Context

`docs/latex-math-evaluable-spec.md` lists `\‖…\‖` (norm / double-pipe) as an
evaluable construct. The AMS Short Math Guide recommends `\lVert…\rVert` for
norms and `\lvert…\rvert` for absolute value.

The lexer's post-lex remap pass handles lowercase `\lvert`/`\rvert` → `PIPE`.
It does **not** remap `\lVert`/`\rVert` (capital V). Those arrive at the parser
as raw `COMMAND` tokens.

## Decision

Do not support `\lVert`/`\rVert` in Tier 1. The absolute value form (`|…|` and
`\lvert…\rvert`) is sufficient for evaluatex parity.

## Consequences

- `\lVert…\rVert` input will produce a SymbolNode for `lVert` (arity 0) and the
  evaluator will fail with "undefined symbol".
- Tier 2 fix: add `lVert`/`rVert` to the lexer's remap table alongside
  `lvert`/`rvert`, so the parser sees them as `PIPE` tokens automatically.
