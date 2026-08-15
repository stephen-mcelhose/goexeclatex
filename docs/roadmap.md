---
type: proposal
title: goexeclatex — Roadmap
description: Forward-looking milestone tracker for goexeclatex. Current focus, upcoming versions, issue links, and process reminder. Update every session.
tags: [goexeclatex, roadmap, milestones]
timestamp: 2026-08-15T03:04:00Z
---

# goexeclatex — Roadmap

For architecture context see [[plan]]. For the decision to split these documents see [[adrs/adr-012-plan-roadmap-separation]].

---

## Current focus

**Process / agent workflow:** [#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15) C+B on branch `chore/15-host-agent-skills` (encode Phase 3 + host Tier A skills). Product next: vFuture (#7) or remaining `\mod`/`\pmod`/`\pod` if needed. Spec for shipped v0.4: [[specs/parser-extensions]].

| Issue | Title | Priority |
| ----- | ----- | -------- |
| ~~[#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)~~ | ~~v0.4 parser extensions~~ | ✅ closed — merged #14 |
| ~~[#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15)~~ | ~~Host agent skills (C+B)~~ | ✅ closed — [[adrs/adr-016-self-host-agent-skills]]; landing via `chore/15-host-agent-skills` |
| ~~[#1](https://github.com/stephen-mcelhose/goexeclatex/issues/1)~~ | ~~Pre-pass byte offset drift corrupts downstream error positions~~ | ✅ closed — ADR-013 |
| ~~[#2](https://github.com/stephen-mcelhose/goexeclatex/issues/2)~~ | ~~Spec gaps identified in GAN review (§6.1.3, whitespace, group context)~~ | ✅ closed |
| [#3](https://github.com/stephen-mcelhose/goexeclatex/issues/3) | UTF-8 char-mode error message shows raw byte instead of full rune | Low — deferred (ADR-002) |
| [#4](https://github.com/stephen-mcelhose/goexeclatex/issues/4) | Revisit: allow positional arg as expression shorthand | Design decision |

---

## ⚠️ Process — mandatory for all milestone work

See root [`AGENTS.md`](../AGENTS.md) for the full loop. Summary:

**Planning (non-trivial issues).** Prefer `/ac-plan-loop` → clarify + RFC2119 AC
→ breakdown with `### Verify` per task. Artifacts:
`.agents/issues/issue-<N>/`. Hosted skills live under `.agents/skills/`.

**Phase 1 — Spec first.**
Write or update the normative spec in `docs/specs/<name>.md` (OKF frontmatter,
RFC 2119 language). Ground every normative statement in [[plan]] and
[[evaluatex-reference-implementation]]. Document any divergence from the
reference explicitly.

**Phase 2 — TDD.**
1. Write failing tests that pin every normative behaviour in the spec.
2. Confirm the tests fail (red).
3. Implement until all tests pass (green).
4. Commit with a conventional commit message.

**Phase 3 — Closing review (process).**
GAN-style docs pass, wiki lint on touched pages, goldens/tests green, tracking
issue + roadmap updated. (Distinct from [[plan]] architecture “Phase 3 —
Evaluator”.)

Do not write implementation code before the spec exists.
Do not skip the red-phase confirmation.
Plan aloud, then execute.

---

## ~~v0.2 — Tier 2: function completeness ([#5](https://github.com/stephen-mcelhose/goexeclatex/issues/5))~~ ✅ closed

- [x] `\arcsin` / `\arccos` / `\arctan` (canonical names alongside `\asin` etc.)
- [x] `\ln`, `\exp`, `\log` (base 10)
- [x] `\dfrac`, `\tfrac`, `\cfrac` (aliases for `\frac`)
- [x] `\sinh` / `\cosh` / `\tanh` / `\coth` / `\sech` / `\csch`
- [x] `\infty` → `math.Inf(1)`
- [x] `\binom{n}{k}`, `\dbinom`, `\tbinom`
- [x] Greek letter variables (`\alpha`, `\beta`, … — user-supplied via `-v`)
- [x] `\min(a,b)`, `\max(a,b)` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4 (n≥2)
- [x] `\lfloor x \rfloor`, `\lceil x \rceil` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4
- [x] `\sqrt[n]{x}` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4
- [x] `\gcd(a,b)` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4 (n≥2)
- [x] `\bmod` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4
- [ ] `\mod`, `\pmod`, `\pod` — still deferred (AMS congruence annotations; see [[specs/parser-extensions]] §7.4)

---

## ~~v0.3 — Tier 3: subscripts + big operators ([#6](https://github.com/stephen-mcelhose/goexeclatex/issues/6))~~ ✅ closed

- [x] `_` token; subscript grammar rule in parser
- [x] `x_{i}` — subscript variable lookup
- [x] `\\sum_{i=a}^{b} f(i)` — discrete summation engine
- [x] `\\prod_{i=a}^{b} f(i)` — discrete product engine
- [x] `\\lVert v \\rVert` — norm (double-pipe)
- [x] `\\_` → literal `_` pre-pass — **won't do** ([[adrs/adr-013-drop-underscore-from-symbol]]; not evaluable-math)
- [x] `\\log_{b}(x)` — shipped in [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) / v0.4

---

## ~~v0.4 — Parser extensions ([#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8))~~ ✅ closed

- [x] Floor / ceil paired delimiters
- [x] `\sqrt[n]{x}` (incl. odd roots of negatives — [[adrs/adr-015-odd-integer-roots-of-negatives]])
- [x] Variadic `\min` / `\max` / `\gcd`
- [x] `\log_{b}(x)`
- [x] `\bmod`
- [x] GAN-style closing review (docs punch-list cleared)

---

## vFuture ([#7](https://github.com/stephen-mcelhose/goexeclatex/issues/7))

**Depends on:** v0.3 (#6) closed. Each item needs its own issue, spec, and ADR before implementation.

- `\begin{cases}…\end{cases}` — conditional piecewise evaluation
- Matrix input + `\det` — determinant of a matrix literal
- JSON output mode (`-json`) — structured output for tooling integration

---

## Sources

- [[plan]] — architecture, pipeline design, error policy
- [[adrs/adr-012-plan-roadmap-separation]] — why this document exists
