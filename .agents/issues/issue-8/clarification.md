---
type: spike
title: "Issue #8 clarification — parser extensions (floor/ceil, nth-root, variadic, mod, log_b)"
description: Clarification, RFC2119 AC, and completion strategy for GitHub issue #8.
tags: [goexeclatex, issue-8, clarification]
timestamp: 2026-08-15T15:55:00Z
---

# Issue #8 — Clarification

**Repo:** stephen-mcelhose/goexeclatex  
**Title:** v0.2.1 Tier 2.1: parser extensions for floor/ceil, nth-root, variadic paren functions, mod  
**Artifacts:** `.agents/issues/issue-8/`  
**Explore:** [[explore.md]]

---

## What We Know

Issue #8 is the deferred **parser-capability** milestone that unblocks several Tier 2/3 surface features. The builtins themselves (`floor`, `min`, `mod`, …) are small; the hard part is teaching the lexer/parser four new syntactic shapes without inventing non-LaTeX syntax:

1. **Paired named delimiters** — `\lfloor…\rfloor`, `\lceil…\rceil` (today false-parse as implicit multiply).
2. **Optional bracket args** — `\sqrt[n]{x}` (today mis-parses `[n]` as the sqrt argument).
3. **Variadic paren args** — `\min(a,b)`, `\max(a,b)`, `\gcd(a,b)`, and `\log_{b}(x)` (COMMA is lexed but never consumed; `\log_{b}` fails at parse).
4. **Infix modulo commands** — `\mod` / `\pmod` / `\bmod` / `\pod` (today false-parse as multiply).

Project process for this issue (from AGENTS.md / roadmap, plus closing review):

1. **Spec first** — normative RFC 2119 spec(s) grounded in plan + evaluatex; divergences explicit.
2. **TDD** — red → green → conventional commits per capability slice.
3. **GAN-style review** — end-of-milestone discriminator pass on specs + behaviour (same class of work as closed #2): hunt gaps, contradictions, undocumented divergences, and “false success” parses.

---

## Codebase Touchpoints

| File / Package | Role | Why Affected |
| -------------- | ---- | ------------ |
| `internal/lexer/token.go` | Token types | Floor/ceil (and possibly mod) tokens; COMMA already exists |
| `internal/lexer/lexer.go` | Lex + arity + remap | Dual `commandArities`; patterns; `postLexRemap`; char-mode |
| `internal/parser/parser.go` | Grammar | `parseAtom` / `parseCommand` / abs+norm depth; COMMA unused |
| `internal/parser/node.go` | AST | `FunctionNode.Args` already `[]Node`; maybe dedicated nodes |
| `internal/eval/{eval,functions}.go` | Semantics | New builtins + binary `mod` |
| `cmd/goexeclatex/testdata/stdin.txt` | E2E golden | End-to-end cases |
| `docs/specs/{lexer,parser,eval}.md` + new Tier 2.1 spec | Normative | Phase 1 before code |
| `docs/adrs/*` | Design decisions | Optional `[n]`, variadic parens, floor tokens; ADR-013 vs `\_` |
| `docs/roadmap.md` | Tracker | Checklist + stale `\_` bullet |

**Dangerous status quo:** several #8 inputs already **parse successfully with the wrong AST** (`\lfloor…\rfloor`, `\sqrt[n]{x}`, `a \bmod b`). Closing #8 means converting silent wrong answers / wrong products into correct eval or hard errors.

---

## Assumptions Being Made

1. **No invented syntax** — brace-only shims (e.g. `\min{a}{b}`, `\log_{b}{x}`) are rejected; LaTeX-shaped forms only.
2. **Floor/ceil follow NORM, not PIPE remap** — dedicated open/close tokens + depth counters (absDepth/normDepth pattern).
3. **`FunctionNode` stays the AST for min/max/gcd/floor/ceil/sqrt/log`** unless an ADR says otherwise; mod is `BinaryNode`.
4. **Dual arity tables** (lexer + parser) must stay in sync for every command touched.
5. **`\_` is NOT in scope** — ADR-013 already removed the pre-pass; roadmap L71 is stale (confirm Q1).
6. **Existing brace forms remain valid** — `\sqrt{x}`, `\log{10}`, `\sin{x}` must not regress; `\sin(x)` accidental paren form should keep working or be deliberately specified.
7. **evaluatex is a reference, not a ceiling** — evaluatex lacks floor/nth-root/mod; we still implement them per evaluable-subset intent, documenting divergence.

---

## Clarifying Questions

**Q1 — `\_` scope** ✅ answered
> Roadmap still defers `\_` → literal `_` to #8, but ADR-013 rejected that pre-pass and the issue body omits it. Confirm `\_` is **out of scope** for #8 (roadmap cleanup only)?
*Why it matters: otherwise we reopen a closed ADR and add a lexer pre-pass workstream.*

**Answer:** Out of scope. Decision already made in ADR-013. Roadmap updated to “won't do”.

**Q2 — Milestone naming**
> Issue title says “v0.2.1 Tier 2.1” while it also carries v0.3’s `\log_{b}`. Prefer **v0.4** (next after v0.3), keep **v0.2.1**, or rename to “parser extensions” without a version bump label?
*Why it matters: roadmap / commit / release framing.*

**Q3 — Variadic arity**
> For `\min` / `\max` / `\gcd`: require **exactly 2** args, or allow **n ≥ 2** (evaluatex-style `sum {',' sum}`)?
*Why it matters: grammar + eval loops + error messages.*

**Q4 — Modulo family semantics**
> In evaluable math, should `\mod`, `\pmod`, `\bmod`, `\pod` all evaluate to the same binary remainder (e.g. Go `%` / `math.Mod` policy), with only syntax/shape differing — or does `\pmod`/`\pod` require a different parse shape (second operand in braces/parens)?
*Why it matters: whether mod is pure infix remap or command+arg hybrids.*

**Q5 — Spec packaging**
> One new `docs/specs/parser-extensions.md` (Tier 2.1) that owns all four grammar shapes, with surgical updates to lexer/parser/eval specs — or fold everything into existing `parser.md` / `eval.md` only?
*Why it matters: Phase 1 structure and GAN review surface.*

**Q6 — `\lfloor x \rceil` (round)** ✅ answered
> `docs/latex-math-evaluable-spec.md` lists mixed floor/ceil as “round”. In or out of #8?
*Why it matters: extra delimiter pairing matrix vs floor/ceil only.*

**Answer:** Out of scope. [[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]] — not AMS semantics; revisit only under that ADR’s criteria.

**Q7 — Single-arg paren policy**
> Today `\sin(x)` works via char-mode fallback. When we add real paren-arg lists, should the spec **MUST** keep unary `(arg)` for arity-1 commands, or only guarantee braces + the new variadic commands?
*Why it matters: regression risk vs intentional narrowing.*

---

## Out of Scope (Proposed)

- Restoring `\_` pre-pass (ADR-013) — pending Q1
- `\begin{cases}`, matrices/`\det`, `-json` (vFuture / #7)
- Vector norms / multi-arg `\lVert`
- Invented brace syntax for min/max/log-base
- Continuous integrals / non-integer sum bounds (already out of v0.3)
- Mixed `\lfloor…\rceil` round — out of scope (ADR-014)

---

## Answers

| Q | Answer |
| - | ------ |
| Q1 | `\_` out of scope (ADR-013); roadmap cleaned |
| Q2 | Milestone **v0.4** |
| Q3 | `\min`/`\max`/`\gcd`: **n ≥ 2** |
| Q4 | Ship **`\bmod` only**; defer `\pmod`/`\mod`/`\pod` |
| Q5 | New `docs/specs/parser-extensions.md` |
| Q6 | Mixed round out of scope (ADR-014) |
| Q7 | Keep unary paren form (`\sin(x)`) — SHOULD |

## Open Items

- ~~Resolve Q1–Q7~~ — locked 2026-08-15
- GAN-style review checklist as closing task (Task 7)
- False-success regressions required (AC-MUSTNOT-2)

---

## Completion strategy (process)

Aligned with AGENTS.md; closing review made explicit.

```
Phase 0  Clarify + AC + plan          ← this document / issue updates
Phase 1  Spec first (+ ADRs as needed)
Phase 2  TDD by capability slice      (red → green → commit)
Phase 3  GAN-style discriminator review → fix gaps → close #8
```

### Phase 1 — Spec first

- Write normative Tier 2.1 behaviour (RFC 2119) covering the four grammar shapes + eval semantics.
- Ground in `docs/plan.md` + evaluatex reference; **document divergences** (evaluatex has no floor / nth-root / mod).
- ADRs for any non-obvious design (floor token strategy, optional `[n]` vs grouping `[`, variadic arity, mod family).
- Update stale deferred tables in `parser.md` / `eval.md` / roadmap (`\_` bullet).

### Phase 2 — TDD (suggested capability order)

Order minimizes “wrong success” hazard and builds shared paren infrastructure once:

1. **Paired floor/ceil** (NORM-like tokens + depth) — kills false multiply parses.
2. **Optional `\sqrt[n]{x}`** — isolates `[` optional-arg from grouping.
3. **Variadic paren grammar + COMMA** — shared substrate.
4. **`\min` / `\max` / `\gcd`** on that substrate.
5. **`\log_{b}(x)`** — subscript base + paren arg (special-case like sum/prod).
6. **Infix mod family** — remap or product-level dispatch.
7. **CLI golden + docs/examples** refresh.

Each slice: failing table-driven tests → confirm red → implement → green → conventional commit.

### Phase 3 — GAN-style review (end)

Modeled on closed #2 (“spec gaps identified in GAN review”):

| Pass | Discriminator asks |
| ---- | ------------------ |
| Spec vs plan/evaluatex | Undocumented divergence? Missing MUST/MUST NOT? |
| Spec vs tests | Every normative rule pinned? |
| False-success hunt | Inputs that used to wrong-parse now correct or hard-error? |
| Dual tables / depth | lexer↔parser arity sync; abs/norm/floor depth interactions |
| Docs drift | roadmap, gap-analysis, examples, how-to |

Output: punch-list issues or inline fixes; #8 closes only when punch-list is empty (or explicitly deferred with issue links).

---

## Acceptance Criteria (RFC2119)

### MUST

- **AC-MUST-1** — A normative Tier 2.1 / parser-extensions spec MUST exist before implementation code for #8 features lands.  
  **Verify:** `test -f docs/specs/<agreed-name>.md` and greppable RFC 2119 language; CI/docs review in GAN pass.

- **AC-MUST-2** — `\lfloor expr \rfloor` MUST evaluate to `math.Floor` of `expr`; `\lceil expr \rceil` MUST evaluate to `math.Ceil` of `expr`. Unmatched / empty pairs MUST error at parse.  
  **Verify:** `go test ./internal/{parser,eval}/...` table cases; CLI golden lines.

- **AC-MUST-3** — `\sqrt[n]{x}` MUST evaluate to the principal n-th root of `x` for supported numeric `n`; `\sqrt{x}` MUST remain square root.  
  **Verify:** parser/eval table tests incl. `\sqrt[3]{27}` → `3`; regression `\sqrt{9}` → `3`.

- **AC-MUST-4** — `\min(a,b)`, `\max(a,b)`, `\gcd(a,b)` MUST parse comma-separated paren args and evaluate correctly (arity per Q3).  
  **Verify:** table-driven parser+eval tests; probe that COMMA inside these forms no longer errors incorrectly.

- **AC-MUST-5** — `\log_{b}(x)` MUST parse subscript base + paren argument and evaluate as `\ln(x)/\ln(b)` (or equivalent); plain `\log{x}` MUST remain base-10.  
  **Verify:** eval tables for `\log_{2}(8)` → `3`; `\log{100}` → `2`.

- **AC-MUST-6** — Supported mod forms (set per Q4) MUST evaluate as binary modulo / remainder with documented float policy — not as implicit multiply.  
  **Verify:** `a \bmod b` (and agreed siblings) eval tests; AST is binary mod, not `*`.

- **AC-MUST-7** — Each capability slice MUST follow TDD: failing tests committed or demonstrated red before implementation green.  
  **Verify:** session/PR narrative + `go test` history; conventional commits (`test:` then `feat:` or combined with note).

- **AC-MUST-8** — Before closing #8, a GAN-style review MUST run against the new/updated specs and the false-success probe set, and resulting gaps MUST be fixed or explicitly deferred with links.  
  **Verify:** issue comment checklist completed; no open punch-list items without deferral links.

### MUST NOT

- **AC-MUSTNOT-1** — The implementation MUST NOT invent non-LaTeX brace shims for these features (`\min{a}{b}`, `\log_{b}{x}` as the supported form, etc.).  
  **Verify:** negative tests / spec Non-Scope section; GAN review.

- **AC-MUSTNOT-2** — `\lfloor…\rfloor`, `\sqrt[n]{x}`, and infix mod inputs MUST NOT silently succeed as implicit-multiply products once the feature is implemented.  
  **Verify:** regression tests asserting correct node type / value (or hard error if intentionally unsupported shape).

- **AC-MUSTNOT-3** — #8 MUST NOT reintroduce the `\_` → `_` pre-pass unless a new ADR supersedes ADR-013 *(provisional on Q1)*.  
  **Verify:** no pre-pass code; roadmap bullet removed or retargeted.

### SHOULD

- **AC-SHOULD-1** — Floor/ceil SHOULD use dedicated tokens + depth guards modeled on NORM/absDepth, not PIPE remapping. *(provisional — confirm)*  
  **Verify:** lexer token types present; depth tests analogous to abs/norm.

- **AC-SHOULD-2** — Dual `commandArities` tables SHOULD be updated together in the same commit for any arity change.  
  **Verify:** code review / grep both maps in the feat commit.

- **AC-SHOULD-3** — Unary paren form for existing arity-1 commands (e.g. `\sin(x)`) SHOULD remain accepted once paren lists exist. *(provisional — Q7)*  
  **Verify:** regression `\sin(\pi/6)` / `\sin(x)` cases.

- **AC-SHOULD-4** — CLI `testdata/stdin.txt`, `examples.md`, and roadmap checkboxes SHOULD be updated in the same milestone as the features.  
  **Verify:** docs diff + golden run.

### MAY

- **AC-MAY-1** — ~~Mixed `\lfloor x \rceil` “round”~~ — **withdrawn**; out of scope per ADR-014.
- **AC-MAY-2** — n-ary min/max/gcd beyond two args MAY be supported if Q3 chooses n ≥ 2.

---

## Sources

- GitHub issue #8 (+ comment on `\log_{b}` deferral from #6)
- `.agents/issues/issue-8/explore.md`
- `AGENTS.md`, `docs/roadmap.md`, `docs/plan.md`
- `docs/specs/{lexer,parser,eval,subscripts-largeops}.md`
- `docs/adrs/adr-013-drop-underscore-from-symbol.md`
- Closed issue #2 (GAN review precedent)
