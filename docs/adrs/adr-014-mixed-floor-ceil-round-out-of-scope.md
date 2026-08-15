---
id: adr-014
type: decision
title: "ADR-014: Mixed floor/ceil delimiters are out of scope (no round convention)"
description: >
  Rejects treating \lfloor x \rceil as nearest-integer round. That pairing is
  not standard AMS/LaTeX semantics; only matched floor and ceiling pairs are in
  scope. Includes revisit criteria if a real need appears later.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T16:05:00Z
tags: [adr, parser, floor, ceil, round, out-of-scope]
---

# ADR-014: Mixed floor/ceil delimiters are out of scope (no round convention)

## Status

Accepted — out of scope for issue [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)
and subsequent floor/ceil work. Supersedes the “Round (nearest)” row in
[[latex-math-evaluable-spec]] for product behaviour.

## Context

`docs/latex-math-evaluable-spec.md` listed:

| LaTeX | Meaning | Evaluable? |
| ----- | ------- | ---------- |
| `\lfloor x \rfloor` | Floor | ✅ |
| `\lceil x \rceil` | Ceiling | ✅ |
| `\lfloor x \rceil` | Round (nearest) | ✅ (convention) |

Matched `\lfloor…\rfloor` and `\lceil…\rceil` are normal LaTeX/AMS paired
delimiters. The **mixed** form `\lfloor x \rceil` (left floor + right ceiling)
is still typesettable — the macros exist — but AMS does **not** define it as
“round to nearest integer.” The “(convention)” tag was a project invention for
the evaluable subset, not core LaTeX.

Issue #8 needs paired-delimiter support for floor and ceiling. Adopting mixed
pairs as `math.Round` would:

1. Encode non-standard semantics that authors cannot look up in AMS docs.
2. Complicate open/close matching (four pairings instead of two matched pairs).
3. Invite further invented delimiter matrices without a LaTeX anchor.

Project policy remains: do not invent syntax or evaluation meanings that LaTeX
does not already carry.

## Decision

1. **Out of scope:** `\lfloor … \rceil` and `\lceil … \rfloor` MUST NOT be
   given round (or other) evaluation semantics.
2. **In scope (when #8 / floor-ceil lands):** only matched pairs
   `\lfloor … \rfloor` → floor and `\lceil … \rceil` → ceiling.
3. **Mismatch behaviour:** if floor/ceil delimiter tokens are implemented,
   a mixed open/close pair MUST be a parse error (unmatched / mismatched
   delimiter), not silent multiply and not round.
4. **Catalogue update:** the evaluable-spec “Round (nearest)” row is marked
   out of scope per this ADR (not a promised feature).

## Revisit criteria

Reopen or supersede this ADR only if **all** of the following hold:

1. **External demand** — a concrete user, corpus, or downstream tool needs
   nearest-integer round expressed in LaTeX-like input (not merely a nice-to-have).
2. **No better surface** — an explicit form is unavailable or rejected
   (e.g. a documented `\operatorname{round}` / `\round{x}` style), **or** a
   cited external standard defines mixed `\lfloor…\rceil` as round.
3. **Spec + tests first** — a normative RFC 2119 rule, divergence note from
   AMS, and table-driven tests are written before implementation (AGENTS.md
   Phase 1 / Phase 2).

Until then, do not implement mixed-pair round “for completeness.”

## Consequences

- Issue #8 planning treats Q6 (mixed round) as **answered: out of scope**.
- Floor/ceil implementation stays two matched pairings; no round builtin via
  delimiter mixing.
- Authors wanting nearest-integer round must use a future explicit command
  (if ever added under the revisit criteria) or compute it outside goexeclatex.
- GAN / closing review for #8 should confirm mixed pairs error rather than
  evaluate as round.

## Sources

- [[latex-math-evaluable-spec]] — Floor, Ceiling, Rounding table
- AMS Short Math Guide — delimiter inventory (`\lfloor`/`\rfloor`, `\lceil`/`\rceil`); no mixed-pair round rule
- GitHub issue [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) — parser extensions planning
- [[adrs/adr-013-drop-underscore-from-symbol]] — precedent for rejecting non-evaluable invented rules
