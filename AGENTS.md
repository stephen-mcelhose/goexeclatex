# Agent Instructions — goexeclatex

## Starting a new session

**Open `docs/roadmap.md` first.** It shows current focus, open issues, and the
next milestone. Then read `docs/plan.md` for architecture context if needed.

## What this project is

A Go CLI that evaluates a LaTeX math expression and prints its numeric result.
Reads from stdin by default or via `-e`; supports variable bindings (`-v`) and
output precision (`-p`).

## Working style

- **Plan aloud, then execute** — state the next step before multi-step work;
  do not silently chain large sequences.
- **Update tracking issues** when decisions lock (GitHub comments), not only
  local notes.
- End-of-session: update `docs/roadmap.md` (and `docs/log.md` when material).

## Issue planning loop (hosted skills)

For non-trivial issues, prefer the hosted skills under `.agents/skills/`:

| Skill | Role |
| ----- | ---- |
| `ac-plan-loop` | Orchestrate clarify → RFC2119 AC → breakdown with Verify → gates |
| `clarify-issue` | Explore + clarification doc → `.agents/issues/issue-<N>/` |
| `breakdown-issue` | Ordered tasks + `### Verify` per task |
| `llm-wiki` | Ingest / query / lint the `docs/` wiki (**canonical copy is in-repo**) |

Artifacts live at `.agents/issues/issue-<IDENTIFIER>/`. Do not use personal
`~/.config/csgdaa-code/...` paths for this repo.

Personal `~/.cursor/skills/` copies may exist; **prefer the in-repo skill** when
both are available. Tier B skills (cognitive locality, define-contracts, …)
are not hosted here yet — file a follow-up if needed ([#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15) recorded C+B).

## Implementation strategy

Every new package or subsystem follows this workflow — **no exceptions**:

### Phase 1 — Write the spec first

Before writing any implementation code, produce a normative spec document in
`docs/specs/<name>.md` (OKF frontmatter, RFC 2119 language). The spec must be
grounded in two sources:

1. **`docs/plan.md`** — the authoritative architecture for this project.
2. **The evaluatex reference implementation** — ingested at
   `docs/raw/evaluatex-reference-implementation.md` and summarised in
   `docs/evaluatex-reference-implementation.md`. Where our behaviour
   intentionally diverges, the spec must document the divergence explicitly.

The lexer spec at `docs/specs/lexer.md` is the worked example to follow.

### Phase 2 — TDD the implementation

1. Write failing tests that pin every normative behaviour in the spec.
2. Confirm the tests fail (red).
3. Implement until all tests pass (green).
4. Commit with a conventional commit message (`feat:` / `fix:` / `test:` etc.).

Do not write implementation code before the spec exists. Do not skip the
red-phase confirmation.

### Phase 3 — Closing review (process)

> Not the same as `docs/plan.md` “Phase 3 — Evaluator”. That is the architecture
> build order for `internal/eval/`. **This** Phase 3 is the end-of-milestone
> discriminator pass before calling work done.

Before calling a milestone done:

1. **GAN-style docs pass** — specs/ADRs/roadmap match shipped behaviour;
   divergences and deferred items are explicit.
2. **Wiki lint** — run `llm-wiki` lint (or equivalent) on touched `docs/` pages;
   clear punch-list items.
3. **Goldens / package tests** — `go test` for affected packages and CLI
   goldens green (avoid raw-shell replay that mangles `\t`/`\f`/`\b`).
4. **Tracking** — comment on the GitHub issue with what shipped / deferred;
   update `docs/roadmap.md`.

## Package layout

```
internal/lexer/    — done (Tier 1 + v0.3 UNDERSCORE/NORM + v0.4 FLOOR/CEIL/BMOD)
internal/parser/   — done (Tier 1 + v0.3 + v0.4 floor/ceil, sqrt[n], variadic, log_b, bmod)
internal/eval/     — done (Tier 1 + v0.3 + v0.4 builtins; ADR-015 odd roots)
cmd/goexeclatex/   — done (Tier 1 complete; goldens through v0.4)
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
| `.agents/skills/`                                | Hosted agent skills (planning + llm-wiki)            |

## Conventions

- Conventional Commits (`feat:`, `fix:`, `chore:`, `test:`, `docs:`)
- Table-driven tests in `_test.go` files per package
- No implementation code without a spec; no spec without a plan anchor
- Update `docs/roadmap.md` at the end of each session
