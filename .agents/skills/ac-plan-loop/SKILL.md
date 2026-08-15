---
name: ac-plan-loop
description: >-
  Orchestrate clarify → RFC2119 acceptance criteria → plan/task breakdown with
  per-step verification → human confirm gates for goexeclatex. Use when the
  user wants an issue planning loop, AC clarification alongside the codebase,
  plan with verification steps, confirm-then-loop, or says "ac-plan-loop",
  "clarify and plan", or "plan this issue with MUST/SHOULD". Thin orchestrator:
  runs clarify-issue and breakdown-issue; does not reimplement them.
license: MIT
metadata:
  version: "1.1.0"
  hosted-in: goexeclatex
---

# ac-plan-loop

Thin orchestrator for:

1. Read issue + codebase → fill missing/ambiguous AC (RFC2119)
2. Plan + task breakdown with a verification step after every task
3. Confirm with the human → loop (revise or advance)
4. Optionally continue into implementation per `AGENTS.md` (spec → TDD → Phase 3)

This skill **does not** replace `clarify-issue` or `breakdown-issue`. It
sequences them and enforces language/verify rules at the gates.

Optional later stages (`cognitive-locality-grouper`, `define-contracts`) are
**Tier B — not hosted in this repo yet**. Skip them unless the user provides
those skills or asks to improvise lightly.

```
issue → clarify (+ RFC2119 AC) → [GATE A]
      → breakdown (+ Verify per task) → [GATE B]
      → implement per AGENTS.md (spec → TDD) → [GATE C / Phase 3]
```

---

## Hard rules

1. **No code changes** until Gate B is approved (unless the user explicitly
   skips planning for a tiny fix — record that skip in the tracking issue).
2. **Every MUST/SHOULD AC** MUST map to at least one verification (automated
   preferred).
3. **Every task** MUST end with a `### Verify` subsection (command + expected
   result). See `references/task-verify.md`.
4. Prefer **table-driven Go tests** when the package already uses them or the
   behavior is multi-case.
5. Use **RFC2119** keywords in AC only (see `references/rfc2119-ac.md`).
6. Prefer **invoking hosted skills** over rewriting their workflows inline.
7. **Plan aloud, then execute** — state the next stage before running it.
8. **Update the tracking GitHub issue** when decisions lock (comments), not
   only local artifacts.

---

## Artifact directory

```
ARTIFACT_DIR = <workspace>/.agents/issues/issue-<IDENTIFIER>/
```

Announce once. Create if missing. All stage outputs go here.

| Stage | Artifact |
| ----- | -------- |
| Clarify | `clarification.md`, `explore.md` |
| Breakdown | `breakdown.md` |
| Optional plan | `plan.md` (large / multi-package issues) |

---

## Stage A — Clarify + RFC2119 AC

1. **Obtain the issue** (in order):
   1. Explicit number, issue URL, pasted body, or board deep-link that embeds
      an issue → fetch with `gh issue view <N> --repo <owner>/<repo>` (default
      repo = this workspace remote).
   2. Else ask once for a number/URL/paste.
   3. Do **not** invent an issue from `gh issue list` / search.
2. **Follow `clarify-issue`** end-to-end.
3. **Additionally** append `## Acceptance Criteria (RFC2119)` to
   `clarification.md` using `references/rfc2119-ac.md`.
4. Present Gate A — do not advance without approval.

### Gate A

> Clarification + AC ready at `<ARTIFACT_DIR>/clarification.md`.
> Review MUST/SHOULD/MAY. How proceed?

Options: answer questions / AC look good — advance / pause.

---

## Stage B — Task breakdown (with verify)

1. For large/ambiguous issues, optionally write `plan.md` (steps + verification
   + success criteria) before the breakdown.
2. Always run **`breakdown-issue`** into `breakdown.md` with:
   - Each task mapped to one or more AC ids when AC exist
   - Each task ending with `### Verify` per `references/task-verify.md`
   - Ordering that respects `AGENTS.md` (spec before code; lexer → parser →
     eval → CLI when spanning the pipeline)
3. Present Gate B.

### Gate B

> Plan/breakdown ready. Task list summary:
> 1. …
> Confirm plan + breakdown?

Options: looks good — implement Task 1 / adjust / pause.

---

## Stage C — Implement (AGENTS.md)

After Gate B approval:

1. For each task (or coherent slice): Phase 1 spec if required → Phase 2 TDD
   (red confirmed → green) → conventional commit when the user asks to commit.
2. Run each task’s `### Verify` before marking the task done.
3. When the milestone/slice is complete: **Phase 3** closing review (see
   root `AGENTS.md`) — GAN-style docs/wiki punch-list, goldens, tracking-issue
   update, roadmap update.

### Gate C (per task or end of unit)

Next task / stop / revisit AC (back to A or B).

---

## Resume

If `ARTIFACT_DIR` already has files:

1. Read existing artifacts
2. Resume at the earliest incomplete gate
3. Ask before overwriting

---

## Anti-patterns

- Skipping Gate A because "the issue was clear"
- Tasks without `### Verify`
- AC written as user stories without MUST/SHOULD
- Implementation before spec when `AGENTS.md` requires a spec
- Reimplementing clarify/breakdown inside this skill instead of invoking them
- Guessing the "next" issue from search when none was provided
