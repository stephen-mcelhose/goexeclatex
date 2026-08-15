---
id: adr-007
type: decision
title: "ADR-007: Track absDepth to resolve PIPE implicit-multiply ambiguity"
description: An absDepth counter tracks nesting of absolute-value bars to correctly distinguish an opening PIPE from a closing PIPE and from implicit-multiply insertion sites.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T00:00:00Z
tags: [adr, parser, abs, implicit-multiply]
---

# ADR-007: Track `absDepth` to resolve PIPE implicit-multiply ambiguity

## Status

Accepted

## Context

The grammar allows PIPE to begin an absolute-value expression (`|expr|`), which
means it can also trigger implicit multiplication (`2|x|` → `2 * |x|`). However,
inside an absolute-value group the closing `|` must not be consumed as the start
of a new implicit-multiply abs-value. Without disambiguation, `|3+4|` fails:
after parsing `4`, the product loop sees `|` and tries to parse it as an
implicit-multiply abs-value, consuming the closing pipe and leaving the outer
abs-value without a closer.

## Decision

Add an `absDepth int` field to the `parser` struct. `parseAbsValue` increments
it before parsing the inner expression and decrements it after. The `canBeginPower`
method only returns `true` for `PIPE` when `absDepth == 0`.

This ensures that a `|` seen while inside an open abs-value group is treated as
the closing delimiter, not as the start of a new implicit-multiply expression.

## Consequences

- `|3+4|` parses correctly as `abs(3+4)`.
- `2|x|` parses correctly as `2 * abs(x)` (PIPE fires implicit-multiply only at
  depth 0, before the outer abs group is opened).
- Nested absolute values `||x||` are not a supported LaTeX construct; if
  encountered they will produce a parse error (no special handling needed).
- The `absDepth` approach is O(1) state — no stack required for Tier 1.
