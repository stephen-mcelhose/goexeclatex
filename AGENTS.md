# Agent Instructions — goexeclatex

## What this project is

A Go CLI that reads a LaTeX math expression from stdin and prints its numeric
result.  See `docs/plan.md` for the full architecture, grammar, milestones, and
error-handling policy.

## Implementation strategy

Every new package or subsystem follows this two-phase workflow — **no exceptions**:

### Phase 1 — Write the spec first

Before writing any implementation code, produce a normative spec document in
`docs/specs/<name>.md` (OKF frontmatter, RFC 2119 language).  The spec must be
grounded in two sources:

1. **`docs/plan.md`** — the authoritative design for this project.
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
internal/lexer/    — done (Tier 1 complete)
internal/parser/   — next
internal/eval/     — after parser
cmd/goexeclatex/   — after eval
```

## Milestones

See `docs/plan.md §Milestones` for the full breakdown.  Current target: **v0.1
(Tier 1 — evaluatex parity)**.

## Key references

| Path                                          | Purpose                              |
| --------------------------------------------- | ------------------------------------ |
| `docs/plan.md`                                | Architecture, grammar, milestones    |
| `docs/specs/lexer.md`                         | Normative lexer spec (worked example)|
| `docs/evaluatex-reference-implementation.md`  | Reference impl summary               |
| `docs/raw/evaluatex-reference-implementation.md` | Raw reference source (read-only)  |
| `docs/log.md`                                 | Append-only session log              |

## Conventions

- Conventional Commits (`feat:`, `fix:`, `chore:`, `test:`, `docs:`)
- Table-driven tests in `_test.go` files per package
- No implementation code without a spec; no spec without a plan anchor
