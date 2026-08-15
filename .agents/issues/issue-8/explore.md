---
type: spike
title: "Issue #8 explore — parser extensions for floor/ceil, nth-root, variadic parens, mod"
description: >
  Codebase relevance analysis for GitHub issue #8 (v0.2.1 Tier 2.1 deferred
  parser features). Maps touch points, precedents, data flows, contradictions,
  and spec anchors.
tags: [goexeclatex, issue-8, explore, parser, lexer]
timestamp: 2026-08-15T15:52:00Z
---

# Issue #8 — Codebase explore

**Issue:** v0.2.1 Tier 2.1: parser extensions for floor/ceil, nth-root, variadic
paren functions, mod.

**Roadmap status** (`docs/roadmap.md` L17): v0.3 complete; next is vFuture or #8.

**Note on `\_`:** Roadmap L71 defers `\_` → literal `_` pre-pass to #8; that item
is **not** in the issue body. See §5.6 and ADR-013 — the codebase has already
**rejected** that pre-pass.

---

## 1. Files / packages / interfaces most likely touched

### Must change (implementation)

| Layer | Path | Why |
| ----- | ---- | --- |
| Lexer tokens | `internal/lexer/token.go` | New types likely: floor/ceil paired delimiters (cf. `PIPE` L20, `NORM` L25). `COMMA` already exists (L22). |
| Lexer scan + arity | `internal/lexer/lexer.go` | `commandArities` (L13–31), `patterns` (L40–64), `postLexRemap` (L259+), char-mode after COMMAND (`lexCharArgs`). |
| Lexer tests | `internal/lexer/lexer_test.go`, new `*_test.go` | Remap / token-stream cases (existing `TestRemapLvert`, `TestMisc` COMMA). |
| Parser AST | `internal/parser/node.go` | Possibly extend `FunctionNode` (already variadic `Args []Node` L47–51) or add nodes; `BinaryNode.Op` (L27–31) for mod. |
| Parser grammar | `internal/parser/parser.go` | `commandArities` (L12–30), `parseAtom` (L293–328), `parseCommand` / `parseCommandArg` (L330–378), `parseAbsValue` / `parseNorm` (L405–443), `parseProduct` / `parseSum` for infix mod, `canBeginPower` (L108–118). |
| Parser tests | `internal/parser/parser_test.go` | Abs/PIPE, sqrt, command-arity tables. |
| Eval dispatch | `internal/eval/functions.go` | `callBuiltin` (L11–112) — add `floor`, `ceil`, `min`, `max`, `gcd`, nth-root / `rootn`, log-base. |
| Eval walker | `internal/eval/eval.go` | `applyBinary` (L59–77) if mod is `BinaryNode`; `evalFunction` (L96+). |
| Eval tests | `internal/eval/eval_test.go` | Parallel tables for new builtins. |
| CLI golden | `cmd/goexeclatex/testdata/stdin.txt` | End-to-end cases (pipeline already wires Lex→Parse→Eval in `cmd/goexeclatex/main.go` L52–67). |

### Spec / docs (Phase 1 before code — AGENTS.md / roadmap process)

| Path | Role |
| ---- | ---- |
| New or extended `docs/specs/*.md` | Normative RFC 2119 for Tier 2.1 (no dedicated #8 spec yet). |
| `docs/specs/parser.md` | Grammar, arity, deferred table L260–269 (stale vs v0.3). |
| `docs/specs/lexer.md` | Token types, arity table §6.4, remap §7. |
| `docs/specs/eval.md` | Deferred table §9 L275–288. |
| `docs/specs/subscripts-largeops.md` | Explicitly out-of-scopes `\log_{b}(x)` → #8 (L377–378). |
| `docs/roadmap.md` | Checklist L56–60, L71–72. |
| `docs/goexeclatex-gap-analysis.md` | Effort/status table L60–65, L78. |
| ADRs | Likely new ADRs for optional `[n]`, variadic parens, floor tokens; revisit ADR-013 if `\_` stays deferred. |

### Key interfaces / signatures

```go
// lexer
func Lex(input string) ([]Token, error)                    // lexer.go:76

// parser
func Parse(tokens []lexer.Token) (Node, error)             // parser.go:33
type Node interface { node() }                             // node.go:3
type FunctionNode struct { Name string; Args []Node; Pos int }  // node.go:47
type BinaryNode struct { Op string; Left, Right Node }     // node.go:27
type NormNode struct { Arg Node }                          // node.go:78

// eval
func Eval(node parser.Node, scope ScopeLookup) (float64, error)  // eval.go:12
func callBuiltin(name string, args []float64) (float64, error)   // functions.go:11
```

Duplicate arity tables that **must stay in sync**:

- `lexer.commandArities` — `internal/lexer/lexer.go:13–31`
- `parser.commandArities` — `internal/parser/parser.go:12–30`

---

## 2. Existing patterns to follow

### 2.1 Paired delimiters (abs / norm) — primary precedent for floor/ceil

**Abs (`|…|` / `\lvert…\rvert`):**

1. Lexer: bare `|` → `PIPE` (`lexer.go:50`); post-lex remap `\lvert`/`\rvert` → `PIPE` (`lexer.go:282–284`).
2. Parser: `parseAtom` → `parseAbsValue` (`parser.go:316–317`, `405–419`).
3. Depth guard: `absDepth` + `canBeginPower` only allows `PIPE` when `absDepth == 0` (ADR-007; `parser.go:50–57`, `108–116`).
4. AST: `FunctionNode{Name: "abs", Args: [inner]}` (`parser.go:419`).
5. Eval: `callBuiltin` case `"abs"` → `math.Abs` (`functions.go:27–28`).

**Norm (`\lVert…\rVert`) — closer model for named open/close commands:**

1. Lexer: dedicated pattern **before** generic `COMMAND` so capital-V is not lowercased into `lvert` (`lexer.go:55–57`; spec `docs/specs/subscripts-largeops.md` §3.3).
2. Token: `NORM` with value `"||"` (`token.go:25`).
3. Parser: `parseNorm` + `normDepth` (`parser.go:319–320`, `422–443`).
4. AST: dedicated `NormNode` (not `FunctionNode`) — `node.go:77–83`.
5. Eval: `evalNorm` → `math.Abs` (`eval.go:194–201`).

**Recommendation for `\lfloor`/`\rfloor`/`\lceil`/`\rceil`:** follow **NORM** (dedicated token types before COMMAND + depth counters), not remap-to-PIPE. Open/close are distinct command names; case-fold collision is not the issue, but dedicated tokens avoid arity-0 SymbolNode + implicit-multiply false parses (see §5.1).

### 2.2 Command arity + brace / char-mode args

- Lexer char-mode after COMMAND: `lexer.go:156–162` uses arity table.
- Parser: `parseCommand` (`parser.go:330–354`) → N× `parseCommandArg` (`360–378`).
- Brace path: only `LPAREN` with value `"{"` is a true command argument group (`parser.go:361–370`).
- Char-mode fallback: `parseAtom()` if `canBeginPower` (`372–377`) — this is why `\sin(x)` accidentally works (see §8).

### 2.3 Subscripts

- Token: `UNDERSCORE` (`token.go:23`); pattern `lexer.go:53`.
- Grammar: `parsePostfix` attaches `_` after atom (`parser.go:250–267`); `parseSubArg` (`275–291`).
- Large-op special case: `parseCommand` branches `sum`/`prod` before arity table (`parser.go:335–338`) → `parseLargeOp`.
- **Pattern for `\log_{b}(x)`:** special-case `log` in `parseCommand` (or post-atom) like `sum`/`prod`, then require paren-arg list — subscript alone is not enough (see §5.4).

### 2.4 Command / remap tables

| Mechanism | Location | Examples |
| --------- | -------- | -------- |
| Arity | `commandArities` in lexer + parser | `sqrt:1`, `log:1`, `frac:2` |
| Remap | `postLexRemap` | `times`/`cdot`→TIMES, `lvert`/`rvert`→PIPE |
| Dedicated pre-COMMAND patterns | `patterns` | `NORM` for `\lVert`/`\rVert` |
| Parser special-case | `parseCommand` | `sum`/`prod` → `parseLargeOp` |

### 2.5 evaluatex grammar hint for variadic parens

`docs/evaluatex-reference-implementation.md` L71:

```
FUNCTION ['(' sum {',' sum} ')' | power]   ← parens optional for 1-arg
```

goexeclatex has **no** `FUNCTION` vs `COMMAND` split and **no** COMMA consumption in the parser today.

---

## 3. Adjacent code that could break

| Risk | Detail |
| ---- | ------ |
| Dual arity tables | Changing `sqrt` arity or adding `min`/`max`/`gcd` in only one package breaks lex↔parse contract. |
| Char-mode + `[` | `[` is already `LPAREN` (`lexer.go:42`). `\sqrt[n]{x}` currently lexes optional-arg brackets as grouping tokens; changing sqrt char-mode / optional-arg handling can break `\sqrt{x}` and any `expr[…]` grouping. |
| Implicit multiply | Unknown COMMAND → `SymbolNode` then `canBeginPower` treats following `LPAREN`/`COMMAND`/`NUMBER` as `*` (`parser.go:161–168`). Floor/mod inputs **already parse wrongly as products** (probe §5). |
| `parseCommandArg` + `(` | 1-arg functions accept bare `(` via char-mode fallback (`\sin(x)` works). Variadic `(a,b)` must not break single-arg `\sin(x)` / `\log{10}`. |
| Subscript + log | `parsePostfix` can subscript **any** atom, including `FunctionNode` if args were consumed first — order of “args vs subscript” for `\log` is the design pivot. |
| `absDepth` / `normDepth` | New paired tokens need analogous depth guards or `canBeginPower` will steal closers (ADR-007 pattern). |
| Spec staleness | `docs/specs/parser.md` §4.2 still says `\lVert` deferred to Tier 2 (L88–91); §10 deferred list outdated (arcsins, norm, subscripts already done). Tests/docs may disagree during #8 work. |
| `BinaryNode.Op` set | Only `+ - * / ^` today (`eval.go:61–73`). Adding `"mod"` etc. requires eval + any AST printers/tests assuming that set. |
| CLI exit codes | Lex/parse → exit 1; eval → exit 2 (`main.go:52–69`). Wrong success parses (floor as multiply) currently may reach eval and fail with “undefined symbol” — changing to real floor must not regress exit policy. |

---

## 4. Data flow: lexer → parser → eval

```
stdin / -e string
        │
        ▼
 cmd/goexeclatex/main.go  RunE
        │  lexer.Lex(input)           // main.go:52
        ▼
 []lexer.Token  (… EOF)
        │  patterns → char-mode args → postLexRemap → EOF
        │  lexer.go:76–87
        ▼
 parser.Parse(tokens)                 // main.go:57
        │  parseSum → product → power → unary → postfix → atom
        │  COMMAND → parseCommand / parseLargeOp
        │  PIPE → parseAbsValue → FunctionNode("abs")
        │  NORM → parseNorm → NormNode
        ▼
 parser.Node (AST)
        │  eval.Eval(node, scope)     // main.go:67
        ▼
 float64 → formatResult → stdout
```

**Scope:** `eval.NewScope()` + `-v` bindings (`main.go:62–65`). Greek letters remain arity-0 symbols resolved from scope.

**Error policy:** plan.md / CLI — lex+parse exit 1, eval exit 2.

---

## 5. Partial addresses / contradictions

### 5.1 Floor / ceil — **dangerous false parse today**

Probe: `\lfloor 3.2 \rfloor` → tokens `COMMAND(lfloor) NUMBER COMMAND(rfloor)`; **PARSE OK** as `*parser.BinaryNode` (implicit multiply of three Symbol/Number nodes). No FLOOR tokens; both commands arity 0.

### 5.2 `\sqrt[n]{x}` — **lexes, wrong AST**

Probe: `\sqrt[3]{27}` →

```
COMMAND(sqrt) LPAREN([) NUMBER(3) RPAREN(]) LPAREN({) NUMBER(27) RPAREN(}) EOF
PARSE OK: *parser.BinaryNode
```

Interpretation: sqrt’s single char-mode/parser arg becomes the `[3]` group → `FunctionNode("sqrt", [3])`, then `{27}` implicit-multiplies → product, **not** cube root. Square-root-only path `\sqrt{4}` remains correct `FunctionNode("sqrt", 1 arg)`.

evaluatex also lacks `[n]` syntax but has a `rootn` builtin (`docs/evaluatex-reference-implementation.md` L143).

### 5.3 Variadic `\min(a,b)` / `\max` / `\gcd` — **COMMA rejected**

Probe: `\min(1,2)` → `COMMAND(min) LPAREN NUMBER COMMA NUMBER RPAREN`; parse error `expected ")" … got ","`. `min` is arity-0 SymbolNode; `(1,2)` is a group that cannot contain COMMA.

### 5.4 `\log_{b}(x)` — **blocked after COMMAND**

Probe: `\log_{2}(8)` → parse error `log expects {…} argument at position 4` (UNDERSCORE). Lexer emits `COMMAND(log) UNDERSCORE LPAREN({) …` because arity-1 char-mode does not consume `_`. Parser insists on a command arg immediately.

Plain `\log{10}` works as base-10 (`functions.go:82–86`; eval spec §6.4).

`docs/specs/subscripts-largeops.md` L377–378: explicitly deferred to #8 (needs paren-arg support).
`docs/specs/eval.md` L195–196 / L287: still says deferred to “Tier 0.3” — **stale** relative to roadmap (#8).

### 5.5 `\mod` / `\pmod` / `\bmod` / `\pod` — **implicit multiply**

Probe: `a \bmod b` → `SYMBOL COMMAND(bmod) SYMBOL`; **PARSE OK** as nested `BinaryNode("*")`. Not infix mod. No remap to an operator token.

### 5.6 `\_` literal underscore — **roadmap vs ADR contradiction**

| Source | Claim |
| ------ | ----- |
| `docs/roadmap.md` L71 | `\_` → literal `_` pre-pass deferred to #8 |
| `docs/log.md` (v0.3 session) | `\_` and `\log_b` deferred to #8 |
| **ADR-013** (accepted) | **Remove** the pre-pass entirely; `\_` is invalid evaluable math; SYMBOL has no `_` |
| Probe | `\_` → `lexer: unexpected character "\" at position 0` |

**Finding:** Issue body does not mention `\_`. Roadmap still lists it under #8, but ADR-013 already closed the design against restoring the pre-pass. #8 work should **drop or re-scope** that roadmap bullet, not reintroduce the pre-pass without a new ADR superseding ADR-013.

### 5.7 Stale “deferred” docs

- `parser.md` §4.2 / §10: norm & arcsin still “Tier 2” though implemented.
- `eval.md` §9: `\log_{b}` target “Tier 0.3”.
- ADR-004 superseded text suggested remapping `\lVert`→PIPE; v0.3 used NORM instead — prefer that lesson for floor/ceil.

---

## 6. Relevant spec / plan / reference sections

| Document | Sections relevant to #8 |
| -------- | ------------------------ |
| `docs/roadmap.md` | L17 current focus; L56–60 Tier 2 deferred checklist; L71–72 `\_` and `\log_{b}`; L28–43 Phase 1 spec / Phase 2 TDD process |
| `docs/plan.md` | Token table incl. COMMA/PIPE (≈L48–58); grammar abs `|` / `\lvert` (≈L88–89); pipeline; error policy; Tier 1 done list |
| `docs/specs/lexer.md` | §3 tokens (PIPE, COMMA); §5.2 pattern order; §6 char-mode + §6.4 arity; §7 remap (`lvert`→PIPE) |
| `docs/specs/parser.md` | §4.2 abs; §4.3 bracket matching (`[`/`]` already grouping!); §5 implicit multiply; §6 arity/args; §10 deferred (stale) |
| `docs/specs/eval.md` | §6.1 `sqrt`/`abs`; §6.4 `ln`/`log`/`exp` + log-base deferral note; **§9 deferred Tier 2.1 table L280–288** |
| `docs/specs/subscripts-largeops.md` | §3.3 NORM; §4.3 norm grammar; §5.3 NormNode; **§8 out-of-scope `\log_{b}(x)` → #8**; §9.3 norm vs pipe |
| `docs/evaluatex-reference-implementation.md` | Pipeline; grammar L71 paren/`','` sums; gaps table L143–157 (nth-root, floor, min/max, gcd, mod, log base) |
| `docs/goexeclatex-gap-analysis.md` | L60–65, L78 status ⏳ #8 |
| `docs/latex-math-evaluable-spec.md` | Floor/ceil L91–93; min/max/gcd L78–80; mod forms L155–156 |
| `docs/examples.md` | Notes that `\sqrt[3]{27}` deferred |
| `docs/adrs/adr-007-absdepth-pipe-ambiguity.md` | Depth-counter pattern |
| `docs/adrs/adr-004-lVert-deferred.md` | Superseded; NORM approach |
| `docs/adrs/adr-013-drop-underscore-from-symbol.md` | Contradicts roadmap `\_` deferral |

**No** `docs/specs/` file yet dedicated to Tier 2.1 / issue #8 — Phase 1 must create one (or extend parser+lexer+eval specs) before implementation.

---

## 7. How abs-value (`|`) and `\lVert` are implemented today

### Absolute value

| Step | Implementation |
| ---- | -------------- |
| Lex `|` | `PIPE` pattern `^\|` — `lexer.go:50` |
| Lex `\lvert`/`\rvert` | Fall through as COMMAND, then `postLexRemap` → `PIPE` value `"\|"` — `lexer.go:282–284` |
| Parse | `parseAtom` case `PIPE` → `parseAbsValue` — `parser.go:316–317`, `405–419` |
| Depth | `absDepth++` around inner `parseSum`; `canBeginPower(PIPE)` only if `absDepth==0` — ADR-007 |
| AST | `FunctionNode{Name: "abs", Args: []Node{inner}}` |
| Eval | `callBuiltin("abs")` → `math.Abs` — `functions.go:27–28` |

Tests: `parser_test.go` abs/`\lvert`/`2|x|` (~L349–380); `eval_test.go` abs cases; lexer `TestRemapLvert`.

### Norm (`\lVert` / `\rVert`)

| Step | Implementation |
| ---- | -------------- |
| Lex | Pattern `^\\(?:lVert|rVert)` as `NORM` **before** COMMAND — `lexer.go:55–57` (avoids lowercasing collision with `\lvert`) |
| Parse | `parseNorm` — `parser.go:422–443`; empty body and unmatched errors |
| Depth | `normDepth` mirrors abs — `parser.go:57`, `112–113`, `427–434` |
| AST | `NormNode{Arg}` — `node.go:77–83` |
| Eval | `evalNorm` → `math.Abs` — `eval.go:194–201` (scalar only; vector norm out of scope) |

Tests: `internal/lexer/subscript_test.go` `TestLVert`, `TestNormDistinctFromPipe`, `TestNormSequence`; eval `TestEvalNorm*` (~L966–1018).

---

## 8. How `\sqrt` and `\log` are implemented today (brace-oriented)

### `\sqrt`

- Arity **1** in both arity tables (`lexer.go:16`, `parser.go:15`).
- Lexer: after `\sqrt`, one char-mode arg; `{…}` expands to `LPAREN({) … RPAREN(})`.
- Parser: `parseCommandArg` prefers `{…}`; otherwise single atom (including `[…]` groups!).
- Eval: `math.Sqrt`; domain error on negative (`functions.go:21–25`).
- **Brace-only for correct square root in practice**; `\sqrt 4` char-mode works for a single token; `\sqrt[n]{x}` is mis-parsed (§5.2).
- No optional-argument grammar; `[` is ordinary grouping (`parser.md` §4.3).

### `\log`

- Arity **1**; eval = **base-10** `math.Log10` (`functions.go:82–86`).
- `\ln` = natural log; `\exp` = `math.Exp`.
- Intended usage: `\log{10}`, `\log 2` (char-mode) — **not** `\log_{b}(x)`.
- Subscript form fails at parse before eval (§5.4).

### Accidental paren support for 1-arg commands

Probe: `\sin(x)` → `FunctionNode("sin", 1 arg)` **succeeds**. Mechanism: char-mode / `parseCommandArg` treats `LPAREN("(")` as the single argument atom via `parseGroup`. That gives **paren form for unary functions**, but **not** comma-separated lists.

---

## 9. COMMA tokens in the parser

| Layer | Status |
| ----- | ------ |
| Lexer | **Yes** — type `COMMA` (`token.go:22`); pattern `^,` (`lexer.go:52`); tested in `TestMisc` (`lexer_test.go` ~L224–232). |
| Parser | **No consumption** — `parseAtom` `default` errors on COMMA. Probe `1,2` → `parser: unexpected COMMA "," at position 1`. Inside `(1,2)`, error is `expected ")" … got ","`. |
| Eval | N/A |

**Implication:** Variadic `\min(a,b)` requires new parser productions (e.g. after COMMAND, optional `LPAREN("(") sum { COMMA sum } RPAREN(")")`) and likely treating `min`/`max`/`gcd` as arity-0 at **lexer** char-mode (so `(` is not eaten as a single char-mode arg) while the **parser** reads the paren list — or a dedicated “paren-function” arity convention.

evaluatex reference grammar already has `{',' sum}` (`evaluatex-reference-implementation.md` L71).

---

## 10. GAN review process in docs

**Only mention in the repo:**

```
docs/roadmap.md:22
| ~~[#2]~~ | ~~Spec gaps identified in GAN review (§6.1.3, whitespace, group context)~~ | ✅ closed |
```

- No definition of “GAN”, no runbook, no ADR, no `docs/log.md` entry explaining the process.
- No files matching `*GAN*`.
- Closest **process** documentation is the mandatory Phase 1 (spec) / Phase 2 (TDD) block in `docs/roadmap.md` L28–43 and `AGENTS.md` — not labeled GAN.

**Conclusion:** GAN appears only as a historical label on closed issue #2 (lexer spec gaps). It is **not** a documented ongoing review process for #8; follow the roadmap/AGENTS spec-first + TDD workflow instead.

---

## Probe appendix (live Lex/Parse, 2026-08-15)

| Input | Lex (abbrev) | Parse |
| ----- | ------------ | ----- |
| `\sqrt[3]{27}` | sqrt, `[`, 3, `]`, `{`, 27, `}` | OK BinaryNode (wrong) |
| `\sqrt{4}` | sqrt, `{`, 4, `}` | OK FunctionNode sqrt |
| `\log{10}` | log, `{`, 10, `}` | OK FunctionNode log |
| `\log_{2}(8)` | log, `_`, `{`, 2, `}`, `(`, 8, `)` | ERR expects `{…}` arg |
| `\min(1,2)` | min, `(`, 1, COMMA, 2, `)` | ERR COMMA vs `)` |
| `\lfloor 3.2 \rfloor` | lfloor, 3.2, rfloor | OK BinaryNode (wrong) |
| `|x|` / `\lvert x \rvert` | PIPE…PIPE | OK FunctionNode abs |
| `\lVert x \rVert` | NORM…NORM | OK NormNode |
| `\sin(x)` | sin, `(`, x, `)` | OK FunctionNode sin |
| `a \bmod b` | a, bmod, b | OK BinaryNode `*` (wrong) |
| `\_` | — | LEX ERR |
| `1,2` | 1, COMMA, 2 | ERR unexpected COMMA |

---

## Suggested implementation order (for later planning)

1. **Spec** covering all #8 bullets + explicit decision on roadmap `\_` (likely “won’t do — ADR-013”).
2. Paired floor/ceil tokens + parse/eval (copy NORM/absDepth pattern).
3. Optional `[n]` for `\sqrt` (lexer/parser coordination so `[` is not a false char-mode arg).
4. Variadic paren-arg grammar + COMMA + `min`/`max`/`gcd` builtins.
5. `\log` special-case: subscript base + paren argument; eval `log(x)/log(b)`.
6. Infix `\bmod` / `\mod` / `\pmod` / `\pod` (token remap or COMMAND special-case at product/sum level).

---

## Sources

- GitHub issue #8 body (user-supplied summary; `gh issue view` blocked in this environment)
- `docs/roadmap.md`, `docs/plan.md`, `docs/specs/{lexer,parser,eval,subscripts-largeops}.md`
- `docs/evaluatex-reference-implementation.md`, `docs/goexeclatex-gap-analysis.md`
- `docs/adrs/adr-004-lVert-deferred.md`, `adr-007-absdepth-pipe-ambiguity.md`, `adr-013-drop-underscore-from-symbol.md`
- `internal/{lexer,parser,eval}/**`, `cmd/goexeclatex/main.go`
- Live Lex/Parse probe of representative inputs
