# Agent Instructions — goexeclatex

## Starting a new session

**Open `docs/roadmap.md` first.** It shows current focus, open issues, and the
next milestone. Then read `docs/plan.md` for architecture context if needed.

## What this project is

A Go CLI that evaluates a LaTeX math expression and prints its numeric result.
Reads from stdin by default or via `-e`; supports variable bindings (`-v`) and
output precision (`-p`).

## Implementation strategy

Every new package or subsystem follows this two-phase workflow — **no exceptions**:

### Phase 1 — Write the spec first

Before writing any implementation code, produce a normative spec document in
`docs/specs/<name>.md` (OKF frontmatter, RFC 2119 language).  The spec must be
grounded in two sources:

1. **`docs/plan.md`** — the authoritative architecture for this project.
2. **The evaluatex reference implementation** — ingested at
   `docs/raw/evaluatex-reference-implementation.md` and summarised in
   `docs/evaluatex-reference-implementation.md`.  Where our behaviour
   intentionally diverges, the spec must document the divergence explicitly.

The lexer spec at `docs/specs/lexer.md` is the worked example to follow.

### Phase 2 — TDD the implementation

1. Write failing tests that pin every normative behaviour in the spec.
2. Confirm the tests fail (red).
3. Implement until all tests pass (green).
4. Commit with a conventional commit message (`feat:` / `fix:` / `test:` etc.).

Do not write implementation code before the spec exists.  Do not skip the
red-phase confirmation.

## Package layout

```
internal/lexer/    — done (Tier 1 + v0.3: UNDERSCORE, NORM tokens)
internal/parser/   — done (Tier 1 + v0.3: subscript grammar, parseLargeOp, parseNorm)
internal/eval/     — done (Tier 1 + v0.3: SubscriptNode, LargeOpNode, NormNode, ScopeLookup)
cmd/goexeclatex/   — done (Tier 1 complete)
```

## Key references

| Path                                             | Purpose                                              |
| ------------------------------------------------ | ---------------------------------------------------- |
| `docs/roadmap.md`                                | **Start here** — current focus, open issues, milestones |
| `docs/plan.md`                                   | Stable architecture: pipeline, grammar, error policy |
| `docs/specs/lexer.md`                            | Normative lexer spec (worked example)                |
| `docs/evaluatex-reference-implementation.md`     | Reference impl summary                               |
| `docs/raw/evaluatex-reference-implementation.md` | Raw reference source (read-only)                     |
| `docs/log.md`                                    | Append-only session log                              |

## Conventions

- Conventional Commits (`feat:`, `fix:`, `chore:`, `test:`, `docs:`)
- Table-driven tests in `_test.go` files per package
- No implementation code without a spec; no spec without a plan anchor
- Update `docs/roadmap.md` at the end of each session
