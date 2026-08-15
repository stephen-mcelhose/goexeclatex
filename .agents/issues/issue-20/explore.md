---
type: spike
title: "Issue #20 explore — public Go library API"
description: >
  Codebase relevance analysis for GitHub issue #20 (expose a public Go library
  API separate from the CLI). Maps touch points, conventions, data flow,
  contradictions (README vs internal-only), precision wording, errors/exit
  codes, minArgs/maxArgs panic risk, and module import paths.
tags: [goexeclatex, issue-20, explore, public-api, library, cli]
timestamp: 2026-08-15T17:56:57Z
---

# Issue #20 — Codebase explore

**Issue:** [feat: expose a public Go library API (separate from CLI)](https://github.com/stephen-mcelhose/goexeclatex/issues/20)

**Roadmap status:** `#20` is **not** listed in `docs/roadmap.md` (current focus
is process/#15 residual and product next = vFuture [#7] or remaining `\mod` /
`\pmod` / `\pod`). This issue sits outside the named milestones — treat as a
cross-cutting packaging concern.

**Issue proposes (signatures):**

```go
package goexeclatex

func Eval(expr string, vars map[string]float64) (float64, error)
func EvalWithPrecision(expr string, vars map[string]float64, prec int) (float64, error)
```

Thin shim over `internal/{lexer,parser,eval}`; CLI should call the public API.
Out of scope: changing CLI behaviour (beyond wiring), rewriting internals.

**Design forks called out in the issue:** stability/semver surface; typed
errors; whether to expose AST. **Comment on issue** (2026-08-15): `minArgs` /
`maxArgs` panic if empty args ever reach `callBuiltin` once a public path can
bypass the parser (hand-built AST or future AST-accepting API).

---

## 1. Files / packages / interfaces likely touched

### Must add (new public surface)

| Layer | Path (candidate) | Why |
| ----- | ---------------- | --- |
| Public package | `pkg/goexeclatex/*.go` **or** root `goexeclatex.go` | Issue forks both layouts; none exist today (`pkg/` dir absent). |
| Public tests | `pkg/goexeclatex/*_test.go` (or root `*_test.go`) | Table-driven API tests (expr + vars → float / error). |
| Normative spec | **new** `docs/specs/library.md` (or similar) | AGENTS.md Phase 1 — no public-API spec exists. |
| ADR | **new** under `docs/adrs/` | Layout (`pkg/` vs root), typed errors, AST exposure, precision semantics. |
| Roadmap / plan | `docs/roadmap.md`, `docs/plan.md` §Package layout | Plan tree (L139–179) lists only `cmd/` + `internal/` — no public pkg. |
| README | `README.md` L5–9 | Already claims “Go library” / “embeddable” — must align with real import path. |

### Must change (CLI becomes consumer)

| Path | Why |
| ---- | --- |
| `cmd/goexeclatex/main.go` L10–14, L52–72 | Today imports `internal/{lexer,parser,eval}` and runs Lex→Parse→Eval inline. Should call public `Eval` (and keep local formatting / flags). |
| `cmd/goexeclatex/main_test.go` + `testdata/{stdin,errors}.txt` | Goldens **must stay green** without CLI behaviour change; exit codes 1 vs 2 are the hard constraint (see §7). |

### Likely touch only if design chooses defensive arity / typed errors

| Path | Why |
| ---- | --- |
| `internal/eval/functions.go` L11–138, `minArgs`/`maxArgs` L140–158 | Issue comment: add `len(args) < 2` before `min`/`max` (and preferably `gcd`). Same class of unchecked `args[0]`/`args[1]` for `frac`, `abs`, trig, etc. |
| `internal/eval/eval.go` | Only if wrapping errors or exporting stage metadata for CLI exit codes. |
| `internal/lexer`, `internal/parser` | **Out of scope** to rewrite; only if typed errors are constructed at origin. |

### Key signatures today (internal — not importable outside module)

```go
// internal/lexer/lexer.go:80
func Lex(input string) ([]Token, error)

// internal/parser/parser.go:34
func Parse(tokens []lexer.Token) (Node, error)

// internal/parser/node.go:3–6 — Node marker is unexported
type Node interface { node() }

// internal/eval/eval.go:12
func Eval(node parser.Node, scope ScopeLookup) (float64, error)

// internal/eval/scope.go:7–29
type ScopeLookup interface { Lookup(name string) (float64, bool) }
type Scope map[string]float64
func NewScope() Scope  // seeds pi, e, tau, phi, infty
```

**Note:** `docs/specs/eval.md` §3 documents `Eval(..., scope Scope)` but the
implementation takes `ScopeLookup` (`eval.go:12`). Public shim should hide this.

### CLI helpers that stay CLI-local (not library)

| Helper | Location | Role |
| ------ | -------- | ---- |
| `formatResult` | `main.go:137–148` | `-p` → string for stdout |
| `parseVars` | `main.go:116–135` | `-v name=value` → `map[string]float64` (+ brace strip) |
| `userMessage` / `die` | `main.go:84–98` | ADR-011 prefix strip + exit |
| `readInput` | `main.go:100–114` | `-e` vs stdin |

---

## 2. Patterns / conventions to follow

| Convention | Where it lives | Implication for #20 |
| ---------- | -------------- | ------------------- |
| Spec before code | `AGENTS.md` Phase 1; worked example `docs/specs/lexer.md` | Write `docs/specs/library.md` (RFC 2119) before implementing the shim. |
| Ground in plan + evaluatex | `docs/plan.md`, `docs/evaluatex-reference-implementation.md` | evaluatex has a two-step compile/evaluate API (`docs/evaluatex-reference-implementation.md` ~L135); our issue proposes a one-shot `Eval(string)` — document divergence. |
| ADR for design forks | `docs/adrs/adr-00N-*.md` | Stability, typed errors, AST yes/no, `pkg/` vs root, and **precision meaning** need ADR(s) before locking API. |
| TDD | Phase 2 | Failing public-API tests first; CLI goldens as regression gate after rewiring. |
| Table-driven tests | `internal/*/*_test.go`, `cmd/.../main_test.go` | Mirror style for `Eval` / error cases. |
| CLI goldens | `cmd/goexeclatex/testdata/stdin.txt`, `errors.txt`; `TestFlagCases` in `main_test.go:123+` | Must remain bit-identical for stdout/stderr substrings and exit codes after CLI→public-API switch. |
| Closing review | Phase 3 process | GAN docs pass, wiki lint, `go test`, issue comment, roadmap update. |
| Conventional commits | repo convention | e.g. `feat: add public goexeclatex Eval API`. |

**Process artifact dir:** `.agents/issues/issue-20/` (this file).

---

## 3. Adjacent breakage risk

| Risk | Severity | Detail |
| ---- | -------- | ------ |
| **Exit code 1 vs 2 collapse** | **High** | CLI today: lex/parse → `die(1, …)` (`main.go:52–59`); eval → `die(2, …)` (`main.go:67–69`). A single `Eval(expr) error` without stage discrimination forces either string-prefix sniffing (`lexer:`/`parser:` vs `eval:`) or typed errors — otherwise “no CLI behaviour change” is violated. |
| ADR-011 strip list | Medium | `userMessage` strips `eval: `, `lexer: `, `parser: ` (`main.go:94–96`). If the public package wraps as `goexeclatex: …` or typed errors with different `Error()` text, goldens in `testdata/errors.txt` / `TestFlagErrors` break. |
| Built-in scope merge | Medium | CLI merges user vars onto `eval.NewScope()` (`main.go:62–65`). Public `Eval(expr, vars)` must document: nil map OK? override builtins? Same as CLI (`cli.md` §5: MAY override). |
| `EvalWithPrecision` vs formatting | **High (semantic)** | Issue says “rounds … to prec **significant digits**”; CLI `-p` is **decimal places** via `'f'` (`main.go:140–147`, `cli.md` §7). Shipping both without ADR will confuse library vs CLI. |
| Exposing AST | High | `Node` uses unexported `node()` (`node.go:3–6`) — external packages cannot implement it. Exposing AST means moving/re-exporting types out of `internal/parser`, expanding semver surface, and making arity panics reachable (issue comment). |
| Semver / module | Medium | Root or `pkg/` export becomes the stability contract; accidental export of internal types (via return values) must be avoided. |
| Cobra / CLI only deps | Low | `go.mod` requires only `cobra`; library package must not pull CLI deps if placed at module root carefully (root file OK if it only imports `internal/*`). |

---

## 4. Data flow — CLI pipeline today

From `docs/specs/cli.md` §6 and `cmd/goexeclatex/main.go` `RunE` (L41–73):

```
stdin | -e
   → readInput (trim; empty → exit 1)
   → parseVars (-v → map; bad → exit 1)
   → lexer.Lex(input)           // error → die(1)
   → parser.Parse(tokens)       // error → die(1)
   → eval.NewScope() + merge vars
   → eval.Eval(node, scope)     // error → die(2)
   → formatResult(result, prec) // -p; print + newline
   → exit 0
```

`die` → stderr `error: ` + ADR-011-stripped message + `os.Exit(code)`.

**Target shape after #20 (behaviour-preserving):**

```
stdin | -e → parseVars → goexeclatex.Eval(expr, vars)
                         → formatResult → stdout
Errors: still map to exit 1 (lex/parse/input) vs 2 (eval) somehow.
```

Precision remains a **CLI formatting** concern unless `EvalWithPrecision` is
defined to mutate the float (see §6).

---

## 5. Partial / contradictory existing code

| Claim | Reality |
| ----- | ------- |
| `README.md` L5: “A Go library for parsing and evaluating…” | **No** importable library package. All logic under `internal/`. No `pkg/` directory. Outside consumers cannot `import` lexer/parser/eval (Go `internal/` rule). |
| README L9: “embeddable … dependency-free” | Embeddable only if you live inside this module or copy code; CLI depends on Cobra. |
| `docs/plan.md` L11–14, L139–179 | Describes a **CLI** product and package tree with only `cmd/` + `internal/`. No public API section. |
| ADR-011 L31, L85 | Mentions “library users” benefiting from `eval:`/`lexer:`/`parser:` prefixes — aspirational; those users do not exist yet outside the module. |
| `docs/specs/eval.md` §3 API | Documents internal `Eval(node, Scope)` — not a string API. Spec version 0.2; would remain the internal contract under a shim. |
| Roadmap | Does not mention #20. |
| evaluatex reference | Two-step compile/evaluate; issue proposes one-shot string `Eval` — intentional simplification vs reference. |

**Conclusion:** marketing/docs already speak as if a library exists; implementation is CLI-over-internal only. #20 closes that gap; README/plan/roadmap need an honest update in the same milestone.

---

## 6. Precision (`-p`) today vs issue “significant digits”

| Aspect | CLI today | Issue draft `EvalWithPrecision` |
| ------ | --------- | -------------------------------- |
| Flag | `-p` / `--prec`, default `-1` (`main.go:78`) | `prec int` parameter |
| Meaning | **Decimal places** after the point | Wording: “**significant digits**” |
| Default / full | `-1` → `FormatFloat(v, 'g', -1, 64)` (`main.go:141–142`; `cli.md` §7.1) | Unspecified if `prec < 0` |
| Fixed | `N ≥ 0` → `FormatFloat(v, 'f', N, 64)`; clamp `N > 100` → 100 (`main.go:144–147`; `cli.md` §7.2) | Unspecified algorithm |
| Return type | **string** on stdout | **`(float64, error)`** — rounding a float64 to N digits is not the same as `'f'` string formatting (binary float cannot hold arbitrary decimal places exactly) |
| Tests | `TestFlagCases` e.g. `prec 4 sqrt2` → `"1.4142"` (`main_test.go:136`) | N/A |

**Recommendation for clarification:** either (a) drop `EvalWithPrecision` from v1 public API and leave formatting to callers (CLI keeps `formatResult`), (b) redefine it as “format helper returning string” matching CLI §7, or (c) implement true significant-digit rounding and document divergence from `-p`. Do not ship under the issue’s current wording without an ADR.

---

## 7. Error types / exit codes / ADR-011

### Errors today

- Plain `fmt.Errorf("lexer: …")` / `"parser: …"` / `"eval: …"` — **no** typed error types, **no** `errors.Is`/`As` hierarchy in code.
- Domain policy: ADR-009 (explicit errors, not NaN); Inf: ADR-010 (success / exit 0).

### Exit codes (`cli.md` §8; `main.go`)

| Condition | Exit |
| --------- | ---- |
| Success / ±Inf | 0 |
| Empty input, bad `-v`, lex, parse | 1 |
| Eval (div0, domain, undefined symbol, …) | 2 |

### ADR-011

CLI strips leading `eval: `, `lexer: `, `parser: ` then prints `error: <message>`
(`main.go:91–97`; `docs/adrs/adr-011-cli-error-prefix-stripping.md`).

ADR-011 **revisit criteria** (L93–95): structured/`UserMessage()` errors could
replace the strip list — natural fit if #20 adds `*ParseError` / `*EvalError`.

### Implication for public API + CLI rewire

To keep CLI behaviour unchanged when calling a single public `Eval`:

1. **Typed errors** (issue design fork) with something CLI can switch on for exit 1 vs 2; **or**
2. Preserve internal prefixes in `error.Error()` and keep string strip + prefix-based exit mapping; **or**
3. Public API returns `(float64, Stage, error)` / multi-value — wider surface.

String-only matching on `"eval: "` is fragile but matches today’s internals.
Typed errors are the cleaner semver story and match the issue’s own design note.

---

## 8. `callBuiltin` / `minArgs` / `maxArgs` — defensive checks for string-only `Eval`?

### Current safety chain

1. Parser `\min`/`\max`/`\gcd` → `parseVariadicParen` (`parser.go:377–402`) **requires `len(args) ≥ 2`** (L399–401).
2. `evalFunction` evaluates args then `callBuiltin` (`eval.go:100–110`).
3. `minArgs` / `maxArgs` index `args[0]` with no length check (`functions.go:140–157`).
4. `gcdArgs` similarly assumes non-empty (`functions.go:160–172`).

### String-only public `Eval(expr, vars)`

If the shim is only:

`Lex → Parse → NewScope+merge → eval.Eval(node, scope)`

then **parser still enforces arity**. Empty-slice panics are **not reachable** via
well-formed LaTeX through that path. Same for other builtins that index
`args[0]`/`args[1]` without checks (`frac`, `abs`, trig, …) — parser
`commandArities` / special parsers supply the expected count.

### When checks become necessary

| Path | Need defensive arity? |
| ---- | --------------------- |
| String-only `Eval` (issue MVP) | **Not required for correctness**; optional defense-in-depth. |
| Public AST / `Eval(Node, …)` / bypassing parser | **Yes** — issue comment; panic becomes a library footgun. |
| Future hand-built `FunctionNode` inside module tests | Already a concern for internal callers. |

**Broader note:** if adding checks, do not stop at `min`/`max` — `callBuiltin`
trusts arity for nearly every case (e.g. `frac` uses `args[1]` at L16). A
single shared arity table or `len` guards would match the issue comment’s
intent.

**Practical recommendation:** string-only MVP can ship without changing
`functions.go`; still **SHOULD** add ≥2 guards for `min`/`max`/`gcd` in the
same PR or a tiny follow-up if AST exposure is deferred but anticipated — low
cost, matches the issue comment’s tracking note from #8 review.

---

## 9. Module path and import conventions

| Item | Value |
| ---- | ----- |
| Module path (`go.mod`) | `github.com/stephen-mcelhose/goexeclatex` |
| Go version | `1.26.5` |
| Current external import used by CLI | `github.com/stephen-mcelhose/goexeclatex/internal/{lexer,parser,eval}` |
| Candidate public imports | **A.** `github.com/stephen-mcelhose/goexeclatex` (root `package goexeclatex`) **B.** `github.com/stephen-mcelhose/goexeclatex/pkg/goexeclatex` |

**Go conventions:**

- Root package named `goexeclatex` is idiomatic for a single-purpose module
  (like many libraries). CLI stays in `cmd/goexeclatex` (`package main`).
- `pkg/` is optional/explicit; some style guides discourage `pkg/` as redundant
  with `internal/` contrast. Issue allows either — **ADR should pick one**.
- Do **not** put the library under `internal/` (defeats the issue).
- Avoid exporting symbols from `package main`.

**Import example (root layout):**

```go
import "github.com/stephen-mcelhose/goexeclatex"

result, err := goexeclatex.Eval(`\frac{1}{2}`, nil)
```

**Import example (`pkg/` layout):**

```go
import "github.com/stephen-mcelhose/goexeclatex/pkg/goexeclatex"
```

(Name stutter; sometimes mitigated by `import gel "…/pkg/goexeclatex"`.)

---

## 10. Summary — highest-leverage clarification questions

1. **Layout:** root `goexeclatex` vs `pkg/goexeclatex`?
2. **`EvalWithPrecision`:** drop, match CLI decimal-places formatting (string?), or true significant-digit float rounding?
3. **Typed errors:** required in v1 so CLI can preserve exit 1 vs 2 without string hacks?
4. **AST:** explicitly out of v1 (recommended given `node()` unexported + panic surface)?
5. **Defensive arity in `callBuiltin`:** same PR as public API or defer until AST is public?
6. **Nil / empty `vars`:** treat as empty map; always seed builtins via `NewScope`?

---

## Sources

- GitHub issue #20 body + comment (minArgs/maxArgs panic risk), 2026-08-15
- `cmd/goexeclatex/main.go`, `main_test.go`, `testdata/`
- `internal/eval/{eval.go,functions.go,scope.go}`
- `internal/parser/{parser.go,node.go}`
- `internal/lexer/lexer.go`
- `docs/specs/cli.md`, `docs/specs/eval.md`
- `docs/adrs/adr-009-*.md`, `adr-010-*.md`, `adr-011-*.md`
- `docs/plan.md`, `docs/roadmap.md`, `README.md`, `go.mod`
- `AGENTS.md` (spec → TDD → closing review)
