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

**v0.3 complete.** Next: vFuture items or #8 (paired brackets, variadic args, `\log_b`). No blocking issues.

| Issue | Title | Priority |
| ----- | ----- | -------- |
| ~~[#1](https://github.com/stephen-mcelhose/goexeclatex/issues/1)~~ | ~~Pre-pass byte offset drift corrupts downstream error positions~~ | ✅ closed — ADR-013 |
| ~~[#2](https://github.com/stephen-mcelhose/goexeclatex/issues/2)~~ | ~~Spec gaps identified in GAN review (§6.1.3, whitespace, group context)~~ | ✅ closed |
| [#3](https://github.com/stephen-mcelhose/goexeclatex/issues/3) | UTF-8 char-mode error message shows raw byte instead of full rune | Low — deferred (ADR-002) |
| [#4](https://github.com/stephen-mcelhose/goexeclatex/issues/4) | Revisit: allow positional arg as expression shorthand | Design decision |

---

## ⚠️ Process — mandatory for all milestone work

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

Do not write implementation code before the spec exists.
Do not skip the red-phase confirmation.

---

## ~~v0.2 — Tier 2: function completeness ([#5](https://github.com/stephen-mcelhose/goexeclatex/issues/5))~~ ✅ closed

- [x] `\arcsin` / `\arccos` / `\arctan` (canonical names alongside `\asin` etc.)
- [x] `\ln`, `\exp`, `\log` (base 10)
- [x] `\dfrac`, `\tfrac`, `\cfrac` (aliases for `\frac`)
- [x] `\sinh` / `\cosh` / `\tanh` / `\coth` / `\sech` / `\csch`
- [x] `\infty` → `math.Inf(1)`
- [x] `\binom{n}{k}`, `\dbinom`, `\tbinom`
- [x] Greek letter variables (`\alpha`, `\beta`, … — user-supplied via `-v`)
- [ ] `\min(a,b)`, `\max(a,b)` — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) (parser: variadic paren args)
- [ ] `\lfloor x \rfloor`, `\lceil x \rceil` — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) (parser: paired brackets)
- [ ] `\sqrt[n]{x}` — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) (parser: optional arg)
- [ ] `\gcd(a,b)` — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) (parser: variadic paren args)
- [ ] `\mod`, `\pmod`, `\bmod`, `\pod` — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) (parser: binary infix)

---

## ~~v0.3 — Tier 3: subscripts + big operators ([#6](https://github.com/stephen-mcelhose/goexeclatex/issues/6))~~ ✅ closed

- [x] `_` token; subscript grammar rule in parser
- [x] `x_{i}` — subscript variable lookup
- [x] `\\sum_{i=a}^{b} f(i)` — discrete summation engine
- [x] `\\prod_{i=a}^{b} f(i)` — discrete product engine
- [x] `\\lVert v \\rVert` — norm (double-pipe)
- [ ] `\\_` → literal `_` pre-pass — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)
- [ ] `\\log_{b}(x)` — log base b — deferred to [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8)

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
