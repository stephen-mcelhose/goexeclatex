# Implementation Layers Reference (goexeclatex)

Per-layer guidance for task breakdowns: coverage, detection signals, and
acceptance-criteria shape.

---

## Spec / docs

**Covers**: Normative specs under `docs/specs/`, ADRs under `docs/adrs/`,
roadmap / plan updates, OKF frontmatter.

**Detected when** touchpoints list `docs/specs/`, `docs/adrs/`, or assumptions
require a behaviour decision before code.

**AC format**:

```
- [ ] Spec exists at docs/specs/<name>.md with OKF frontmatter
- [ ] Normative statements use RFC 2119 (MUST/SHOULD/MAY)
- [ ] Divergences from evaluatex are documented explicitly
- [ ] Spec grounded in docs/plan.md (and reference summary as needed)
```

**Notes**: Spec tasks MUST precede implementation for new subsystems
(`AGENTS.md` Phase 1).

---

## Lexer

**Covers**: Token kinds, scanning, pre-pass / char-mode behaviour in
`internal/lexer/`.

**Detected when** touchpoints list `internal/lexer/` or assumptions mention
new LaTeX commands that need tokens.

**AC format**:

```
- [ ] go test ./internal/lexer/... -count=1 passes
- [ ] Table-driven cases cover happy path + rejected input
```

---

## Parser

**Covers**: Grammar, AST nodes, command parsing in `internal/parser/`.

**Detected when** touchpoints list `internal/parser/` or `node.go`.

**AC format**:

```
- [ ] go test ./internal/parser/... -count=1 passes
- [ ] Cases pin AST shape for the new construct
```

---

## Eval

**Covers**: Numeric evaluation, builtins, domain errors in `internal/eval/`.

**Detected when** touchpoints list `internal/eval/` or numeric semantics ADRs.

**AC format**:

```
- [ ] go test ./internal/eval/... -count=1 passes
- [ ] Error / edge cases covered (NaN, domain, empty args, … as applicable)
```

---

## CLI / goldens

**Covers**: `cmd/goexeclatex/`, flags, stdin/error golden fixtures.

**Detected when** touchpoints list `cmd/` or `testdata` goldens.

**AC format**:

```
- [ ] go test ./cmd/goexeclatex/... -count=1 passes
- [ ] New goldens match intended stdout/stderr (no shell-escape false fails)
```

**Notes**: Prefer Go/Python harnesses over raw bash when LaTeX contains
`\t`, `\f`, `\b`, etc.

---

## Wiki / process

**Covers**: `docs/` wiki lint, `AGENTS.md`, `docs/log.md`, skill docs under
`.agents/skills/`.

**Detected when** the issue is process/docs-only or closing review (Phase 3).

**AC format**:

```
- [ ] llm-wiki lint (or targeted docs check) is clean for touched pages
- [ ] docs/roadmap.md updated if milestone status changed
- [ ] Tracking issue comment records decisions locked in-session
```
