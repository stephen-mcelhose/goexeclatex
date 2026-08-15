---
type: spike
title: "Issue #20 clarification — public Go library API"
description: Clarification, RFC2119 AC, and completion strategy for GitHub issue #20.
tags: [goexeclatex, issue-20, clarification, public-api, library]
timestamp: 2026-08-15T18:00:00Z
---

# Issue #20 — Clarification

**Repo:** stephen-mcelhose/goexeclatex  
**Title:** feat: expose a public Go library API (separate from CLI)  
**Artifacts:** `.agents/issues/issue-20/`  
**Explore:** [[explore.md]]  
**Tracking:** https://github.com/stephen-mcelhose/goexeclatex/issues/20  
**Process twin:** `.agents/issues/issue-8/` (organic ac-plan-loop worked example)

---

## What We Know

Issue #20 closes a packaging gap: all evaluation lives under `internal/{lexer,parser,eval}`, so external Go modules cannot import the engine. The CLI (`cmd/goexeclatex`) is currently the only consumer and wires Lex → Parse → Eval itself. The proposed fix is a **thin public facade** (`Eval(expr, vars)` and possibly a precision helper) that calls the existing internals without duplicating logic. The CLI then becomes a client of that facade; flags, stdin, stdout formatting, ADR-011 message stripping, and exit codes stay in `cmd/`.

This is **not** a grammar/eval capability milestone like #8. It is a **cross-cutting packaging + docs** milestone: new public package, normative library spec, ADR(s) for layout/errors/AST/precision, CLI rewire with behaviour freeze, and README/AGENTS/plan/roadmap honesty (README already claims “library”).

A comment on #20 (from the #8 review) flags `min`/`max`/`gcd` panics if a future API accepts hand-built ASTs. For a **string-only** `Eval`, the parser still enforces arity — defensive checks are optional defense-in-depth for this issue’s MVP.

---

## Process compare vs issue-8

| Dimension | Issue #8 (parser extensions) | Issue #20 (public library API) |
| --------- | ---------------------------- | ------------------------------ |
| Kind | Capability / grammar+eval slices | Packaging / facade + docs |
| False-success hazard | High — inputs already wrong-parsed | Low for math — packaging gap, not silent wrong math |
| Spec surface | New `parser-extensions.md` + lexer/parser/eval updates | New `docs/specs/library.md` (+ plan/README/AGENTS/roadmap) |
| ADR density | Several (floor, odd roots, mod deferrals, …) | Fewer but sharp: layout, typed errors, AST, precision |
| TDD slices | Many capability slices | Few: API tests → CLI rewire → docs |
| Phase 3 GAN | Spec vs evaluatex + false-success hunt | Spec vs CLI contract + importability + docs drift |
| Artifact layout | Same: `.agents/issues/issue-N/{explore,clarification,breakdown}.md` | Same |
| Gate comments | On #8 | On #20 |

**What we deliberately keep identical:** ac-plan-loop gates, RFC2119 AC with Verify hints, `### Verify` per breakdown task, no feature branch before Gate B, Gate A/B/C comments on the tracking issue, Phase 3 discriminator table before close.

**What differs:** #20’s “dangerous status quo” is mostly **docs/marketing false success** (README says library) and **CLI exit-code collapse risk** if a single `error` erases lex/parse vs eval staging — not silent wrong numeric results.

---

## Proposed CLI vs library export (for Gate A)

```
External Go caller ──► package goexeclatex (public) ──► internal/{lexer,parser,eval}
CLI (cmd/goexeclatex) ──► same public package      ──► same internals
```

| Layer | Owns | Does not own |
| ----- | ---- | ------------ |
| Public package | String→float pipeline; optional typed errors; semver surface | Flags, stdin, exit codes, Cobra |
| `internal/*` | Lexer/parser/eval semantics (unchanged) | Stability contract for outsiders |
| `cmd/goexeclatex` | `-e`/`-v`/`-p`, formatting, ADR-011 strip, exit 1 vs 2 | Grammar/eval logic (after rewire) |

**Preferred default (provisional until Answers lock):** module-root `package goexeclatex` so import is `github.com/stephen-mcelhose/goexeclatex`; CLI stays `cmd/goexeclatex`. Alternative `pkg/goexeclatex` is longer and name-stutters.

---

## Codebase Touchpoints

| File / Package | Role | Why Affected |
| -------------- | ---- | ------------ |
| `pkg/goexeclatex/` **or** root `goexeclatex.go` | New public API | Issue MVP surface |
| `*_test.go` beside public package | Table-driven API tests | Phase 2 TDD |
| `cmd/goexeclatex/main.go` | CLI | Rewire Lex/Parse/Eval → public `Eval`; keep format/exit |
| `cmd/goexeclatex/testdata/*`, `main_test.go` | Goldens | Behaviour freeze regression |
| `internal/eval/functions.go` | Builtins | Optional arity guards (issue comment) |
| `docs/specs/library.md` | New normative spec | Phase 1 |
| `docs/adrs/adr-NNN-*.md` | Design forks | Layout / errors / AST / precision |
| `docs/specs/cli.md`, `docs/plan.md` | Architecture | Package tree + “CLI calls library” |
| `README.md`, `AGENTS.md`, `docs/roadmap.md`, `docs/index.md`, `docs/log.md` | Docs | Close README/plan contradiction; package layout |
| `docs/adrs/adr-011-*.md` | Error prefixes | May interact with typed errors / strip list |

---

## Assumptions Being Made

1. **Thin shim only** — no move of lexer/parser/eval out of `internal/`; no logic duplication.
2. **CLI behaviour freeze** — same exit codes, same stdout/stderr messages (ADR-011), same `-p` formatting; only the call path changes.
3. **AST stays internal** for this issue unless Answers say otherwise (keeps semver surface small; `Node.node()` is unexported today).
4. **Built-in scope** — public `Eval` seeds the same builtins as CLI (`NewScope`) and merges `vars` with the same override policy as `cli.md` §5.
5. **Nil `vars`** — treated as empty map (provisional).
6. **`EvalWithPrecision` wording in the issue is wrong relative to CLI** — CLI `-p` is decimal places via `'f'` string format, not significant-digit float rounding. Must resolve in Answers before implementing that helper.
7. **evaluatex divergence** — reference uses compile-then-evaluate; we ship one-shot string `Eval` and document it.
8. **Docs updates are in-scope** — README/AGENTS/plan/roadmap are part of “done,” not a follow-up.

---

## Clarifying Questions

**Q1 — Package layout**
> Prefer **module-root** `package goexeclatex` (`import "github.com/stephen-mcelhose/goexeclatex"`) or **`pkg/goexeclatex`**?
*Why it matters: import path, ADR, README examples, AGENTS package layout.*

**Q2 — `EvalWithPrecision`**
> For v1 public API: **(A)** omit it (callers / CLI format floats), **(B)** add a **string** helper matching CLI §7 decimal places, or **(C)** implement true significant-digit **float** rounding as the issue wording suggests (diverges from `-p`)?
*Why it matters: float rounding ≠ string formatting; wrong choice confuses library vs CLI.*

**Q3 — Typed errors**
> MUST v1 export stage-discriminable errors (e.g. `*LexError`/`*ParseError`/`*EvalError` or a `Stage` enum) so the CLI can keep exit **1 vs 2** without string-prefix sniffing? Or keep prefix strings and sniff for now?
*Why it matters: “no CLI behaviour change” + semver-friendly caller story.*

**Q4 — AST exposure**
> Confirm AST / `Node` / Lex / Parse are **out of scope** for #20 (string `Eval` only)?
*Why it matters: panic surface (issue comment), unexported `node()`, semver width.*

**Q5 — Defensive arity in `callBuiltin`**
> For string-only MVP: **(A)** skip (parser already enforces), **(B)** add ≥2 guards for `min`/`max`/`gcd` in this PR as defense-in-depth matching the #20 comment, or **(C)** require a broader arity table for all builtins?
*Why it matters: scope creep vs closing the review note from #8.*

**Q6 — Nil / empty `vars`**
> Confirm `vars == nil` MUST be accepted and treated as no user bindings (builtins still present)?
*Why it matters: API ergonomics and tests.*

**Q7 — Spec / ADR packaging**
> One new `docs/specs/library.md` plus **one** ADR covering layout + typed errors + AST-out + precision decision — or separate ADRs per fork?
*Why it matters: Phase 1 structure and GAN surface.*

---

## Out of Scope (Proposed)

- Changing CLI flag semantics, exit-code policy, or ADR-011 user-facing messages (beyond rewiring through the public API)
- Rewriting `internal/{lexer,parser,eval}` (except optional small arity guards)
- Exposing AST / compile-then-eval two-step API (unless Q4 reverses)
- JSON output (`-json`), cases/matrices (vFuture / #7)
- Moving off Cobra or changing module path
- Guaranteeing bit-identical float formatting across platforms beyond existing CLI tests

---

## Dangerous status quo / false-success probes

| Probe | Today | After #20 success looks like |
| ----- | ----- | ---------------------------- |
| External `import` of engine | **Fails** (`internal/` rule) | Public import compiles and `Eval` works |
| README “Go library / embeddable” | **Marketing false success** | README shows real import path + example |
| CLI exit 1 vs 2 after rewire | Risk of **collapse to one code** if errors undifferentiated | Goldens / `TestFlagErrors` still distinguish lex/parse vs eval |
| ADR-011 stderr text | Risk of new `goexeclatex:` prefix breaking goldens | Same stripped messages as today |
| `\min()` empty / under-arity via string API | Hard parse error (OK) | Remains hard error; no panic |
| Hand-built AST → `callBuiltin` | N/A externally | Still N/A if AST out of scope; optional internal guards |

**Math false-success:** none expected from the shim itself (same Lex→Parse→Eval path).

---

## Completion strategy (process)

Aligned with AGENTS.md / ac-plan-loop; mirrors #8’s Phase 0→3 shape with packaging-shaped slices.

```
Phase 0  Clarify + AC + Answers     ← this document / Gate A
Phase 1  Spec first + ADR(s)        ← library.md + layout/errors/precision ADR
Phase 2  TDD slices                 ← public API → CLI rewire → docs
Phase 3  GAN-style discriminator    ← docs honesty + goldens + issue close
```

### Phase 1 — Spec first

- Write `docs/specs/library.md` (RFC 2119): API signatures, vars/nil, builtins, errors, non-scope (AST), divergences from evaluatex (one-shot vs compile/eval) and from CLI (`-p` formatting stays CLI-local unless Q2 says otherwise).
- ADR locking Q1–Q4 (and Q2 precision).
- Touch `docs/plan.md` package tree; do not implement until spec+ADR exist.

### Phase 2 — TDD (suggested slice order)

1. **Public `Eval` + error contract** — table tests (happy path, undefined var, parse fail, div0); confirm red.
2. **Implement shim** — green; optional arity guards per Q5.
3. **CLI rewire** — call public API; prove goldens / exit codes unchanged (red only if broken).
4. **Docs** — README, AGENTS package layout, roadmap #20, plan, index/log; library how-to snippet.

Each slice: Verify commands from breakdown; conventional commits when asked.

### Phase 3 — GAN-style review (end)

| Pass | Discriminator asks |
| ---- | ------------------ |
| Spec vs plan / evaluatex | Undocumented divergence (one-shot API)? Missing MUST/MUST NOT? |
| Spec vs tests | Every normative library rule pinned? |
| False-success hunt | External import works; README not aspirational; exit 1 vs 2 intact |
| Dual tables / depth | N/A for grammar; check error-stage mapping CLI↔library |
| Docs drift | AGENTS, README, plan tree, roadmap, cli.md cross-links |

---

## Acceptance Criteria (RFC2119)

### MUST

- **AC-MUST-1** — A normative library spec (`docs/specs/library.md` or agreed path) MUST exist before public-API implementation code lands.  
  **Verify:** file exists with RFC 2119 language; Gate B / PR review.

- **AC-MUST-2** — An importable **module-root** `package goexeclatex` MUST expose `Eval(expr string, vars map[string]float64) (float64, error)` that runs the same Lex→Parse→Eval pipeline as today’s CLI and returns full-precision `float64` (no display rounding).  
  **Verify:** `go test` on the public package; external-style import in tests (`package goexeclatex_test` recommended).

- **AC-MUST-3** — `cmd/goexeclatex` MUST call the public API for evaluation (not `internal/{lexer,parser,eval}` directly) while preserving existing CLI behaviour: exit codes (via typed/stage errors + `errors.As`), ADR-011-stripped stderr, and stdout formatting including `-p` (CLI-local `formatResult`).  
  **Verify:** `go test ./cmd/goexeclatex/...` goldens / flag tests green; grep that `main.go` does not import `internal/lexer|parser|eval`.

- **AC-MUST-4** — Design forks (root layout, typed errors, AST-out, omit float precision API) MUST be recorded in an ADR before Gate B advances to implementation.  
  **Verify:** `docs/adrs/adr-NNN-*.md` linked from clarification Answers / issue comment.

- **AC-MUST-5** — `README.md` (with basic math `Eval` examples), `AGENTS.md` package layout, and `docs/roadmap.md` / `docs/plan.md` package tree MUST describe the public library and CLI split accurately before #20 closes.  
  **Verify:** docs review in Phase 3 GAN; README import path + example; no aspirational-only “library” claim.

- **AC-MUST-6** — Public errors MUST be stage-discriminable (typed errors or equivalent) covering at least lex/parse vs eval so callers and the CLI can distinguish failure stages.  
  **Verify:** library tests with `errors.As`; CLI exit 1 vs 2 goldens.

- **AC-MUST-7** — `callBuiltin` (or equivalent) MUST reject `min`/`max`/`gcd` with fewer than 2 evaluated args via error, not panic.  
  **Verify:** `go test ./internal/eval/...` table case (direct or via under-arity path).

- **AC-MUST-8** — Each implementation slice MUST follow TDD: failing tests demonstrated red before green (or CLI-only slice proven unchanged via existing goldens).  
  **Verify:** session/PR narrative; `go test` history.

- **AC-MUST-9** — Before closing #20, a GAN-style review MUST run against the library spec + false-success probe set above; gaps MUST be fixed or explicitly deferred with links.  
  **Verify:** Gate C / close comment checklist on #20.

### MUST NOT

- **AC-MUSTNOT-1** — #20 MUST NOT export AST/`Node`, lexer tokens, or parser entry points.  
  **Verify:** public package API surface review / `go doc`; GAN.

- **AC-MUSTNOT-2** — Rewiring the CLI MUST NOT change exit codes or user-facing error text for existing golden cases.  
  **Verify:** `go test ./cmd/goexeclatex/...` unchanged expectations.

- **AC-MUSTNOT-3** — The public package MUST NOT duplicate lexer/parser/eval logic; it MUST delegate to `internal/` packages.  
  **Verify:** code review of shim (thin file(s)).

- **AC-MUSTNOT-4** — The public API MUST NOT ship a float-returning “precision” helper that rounds/truncates the numeric result for display purposes.  
  **Verify:** no such symbol in public API; spec Non-Scope / ADR.

### SHOULD

- **AC-SHOULD-1** — Public API SHOULD accept `vars == nil` as empty bindings (locked Answer Q6 — elevate to MUST in implementation).  
  **Verify:** table test `Eval(expr, nil)`.

### MAY

- **AC-MAY-1** — A public `Format(v float64, prec int) string` matching CLI §7 MAY ship in #20 or a follow-up (display parity, no math cost).
- **AC-MAY-2** — Broader `callBuiltin` arity hardening beyond min/max/gcd MAY be deferred.
- **AC-MAY-3** — Lex vs parse may share one “syntax” error type if finer split is deferred; eval MUST remain distinct for exit 2.

---

## Open Items

- ~~Q1–Q7~~ — locked 2026-08-15 (Gate A advanced)
- Draft ADR-017 + `docs/specs/library.md` in Stage B / Task 0–1
- Gate B → feature branch after breakdown approval
- Add #20 to `docs/roadmap.md` current focus when work starts (post Gate B)

---

## Answers

Locked from Gate A discussion 2026-08-15 (user replies + agent community guidance).

| Q | Answer |
| - | ------ |
| Q1 | **Module-root** `package goexeclatex`. `pkg/` is *not* `internal/` — it is a *public* convention with no toolchain force; Go docs prefer root library + `internal/` for private engine ([Organizing a Go module](https://go.dev/doc/modules/layout)). Russ Cox: vast majority of ecosystem packages are **not** under `pkg/`. |
| Q2 | **Omit** float-returning `EvalWithPrecision`. `Eval` MUST return full-precision `float64` (no accuracy sacrifice). Display/`-p` behaviour stays at the edge: CLI keeps `formatResult`. A public **string** `Format` helper matching CLI §7 is **MAY** (same display behaviour, zero math cost) — not required for MVP. |
| Q3 | **Typed / stage-discriminable errors in #20** (library contract). **CLI exit 1 vs 2 mapping stays in #20** — not a separate issue: rewire + behaviour freeze already requires *some* stage discrimination; `errors.As` is a few lines and avoids shipping a throwaway prefix-sniff then churning the public error API later. Punting typed errors while rewiring would either break exit codes or add throwaway sniff code. |
| Q4 | **Out of scope.** AST / Lex / Parse remain `internal/`. Public surface is the string `Eval` facade. ADR records this (gist: internals stay internals). |
| Q5 | **(B)** Add ≥2 guards for `min`/`max`/`gcd` in this PR. Note: panic under parser invariant is defensible; guards are cheap defense-in-depth + close the #20 comment. |
| Q6 | **Yes** — `vars == nil` accepted as empty bindings; builtins still seeded. **Also:** README / library docs MUST include basic math examples (import + `Eval`). |
| Q7 | **One** `docs/specs/library.md` + **one** ADR covering layout, typed errors, AST-out, precision-omit. Locked 2026-08-15. |

### Q1 — community note (why not `pkg/`)

| Path | Meaning |
| ---- | ------- |
| `internal/` | **Compiler-enforced** private — other modules cannot import |
| Root package | **Public** — idiomatic for a single-purpose module ([go.dev layout](https://go.dev/doc/modules/layout)) |
| `pkg/` | **Also public** — folder name only; signals “for import” in large repos; **not** standard and not recommended by the Go team |

So `pkg/` ≠ private. Private is `internal/`. For goexeclatex (library + one CLI), root `goexeclatex` + existing `internal/*` + `cmd/goexeclatex` matches official guidance.

### Q2 — trade-offs (accuracy vs “same behaviour”)

| Option | Same as CLI `-p`? | Math accuracy | Verdict |
| ------ | ----------------- | -------------- | ------- |
| (A) Omit precision API | Display stays CLI-only | Full `float64` from `Eval` | **Chosen for MVP** |
| (B) Public `Format(v, prec) string` | Yes (same `'g'`/`'f'` rules) | Unchanged (format only) | Good follow-on / MAY |
| (C) Round `float64` to N digits | Illusory — binary float ≠ decimal places | **Degrades** usable precision | Rejected |

“Same behaviour” for `-p` is **string formatting**, not a different numeric value. Library callers who need CLI-identical output can use a Format helper later or format themselves; they always get the true eval result from `Eval`.

### Q3 — scope: typed errors vs CLI wiring

```
Library needs typed errors  →  public contract (do in #20)
CLI needs exit 1 vs 2       →  already in #20 via “behaviour freeze”
Connecting them            →  errors.As in main.go (too small to be its own issue)
```

Punting **only** the CLI `errors.As` wiring would force either broken goldens or temporary string sniffing — worse than including the wiring. Punting **typed errors entirely** is a coherent alternate (sniff prefixes forever / follow-up issue) but conflicts with “typed errors” preference and creates API churn. **Decision: both in #20.**

---

## Updated 2026-08-15 — AC / strategy deltas after Answers

- Layout default → **locked** root package (AC-SHOULD-4 → MUST via ADR).
- `EvalWithPrecision` float form → **MUST NOT** ship; Format string helper remains MAY.
- Typed errors + CLI stage mapping → **MUST** in #20 (not a follow-up micro-issue).
- AST → **MUST NOT** export; ADR.
- Arity guards min/max/gcd → **SHOULD** → treat as in-scope MUST for this PR.
- Nil vars + basic math examples in README → **MUST**.

---

## Sources

- GitHub issue #20 (+ comment on minArgs/maxArgs panic risk), 2026-08-15
- `.agents/issues/issue-20/explore.md`
- `.agents/issues/issue-8/{clarification,breakdown,explore}.md` (process twin)
- `AGENTS.md`, `docs/roadmap.md`, `docs/plan.md`
- `docs/specs/cli.md`, `docs/adrs/adr-011-cli-error-prefix-stripping.md`
- `cmd/goexeclatex/main.go`, `internal/eval/{eval,functions,scope}.go`
