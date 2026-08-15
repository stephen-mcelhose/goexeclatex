---
name: ac-plan-loop
description: >-
  Orchestrate clarify → RFC2119 acceptance criteria → plan/task breakdown with
  per-step verification → human confirm gates for goexeclatex. Use when the
  user wants an issue planning loop, AC clarification alongside the codebase,
  plan with verification steps, confirm-then-loop, or says "ac-plan-loop",
  "clarify and plan", or "plan this issue with MUST/SHOULD". Thin orchestrator:
  runs clarify-issue and breakdown-issue; does not reimplement them. Worked
  example: .agents/issues/issue-8/.
license: MIT
metadata:
  version: "1.2.0"
  hosted-in: goexeclatex
---

# ac-plan-loop

Thin orchestrator for:

1. Read issue + codebase → fill missing/ambiguous AC (RFC2119)
2. Plan + task breakdown with a verification step after every task
3. Confirm with the human → loop (revise or advance)
4. Continue into implementation per `AGENTS.md` (spec → TDD → Phase 3 GAN)

This skill **does not** replace `clarify-issue` or `breakdown-issue`. It
sequences them and enforces language/verify rules at the gates.

Optional later stages (`cognitive-locality-grouper`, `define-contracts`) are
**Tier B — not hosted in this repo yet**. Skip them unless the user provides
those skills or asks to improvise lightly.

**Worked example (organic #8 run):** `.agents/issues/issue-8/` — clarification
with completion strategy + dangerous status quo, breakdown with Verify, gate
comments on the tracking issue, then slice TDD and GAN close. Decision record:
`docs/adrs/adr-016-self-host-agent-skills.md`.

```
issue → clarify (+ AC + strategy + false-success) → [GATE A]
      → breakdown (+ Verify per task) → [GATE B]
      → implement per AGENTS.md (spec → TDD) → [GATE C / Phase 3 GAN]
```

---

## Hard rules

1. **No code changes and no feature branch** until Gate B is approved (unless
   the user explicitly skips planning for a tiny fix — record that skip on the
   tracking issue).
2. **Every MUST/SHOULD AC** MUST map to at least one verification (automated
   preferred).
3. **Every task** MUST end with a `### Verify` subsection (command + expected
   result). See `references/task-verify.md`.
4. Prefer **table-driven Go tests** when the package already uses them or the
   behavior is multi-case.
5. Use **RFC2119** keywords in AC only (see `references/rfc2119-ac.md`). Do not
   invent soft synonyms ("ideally", "try to") for requirements.
6. Prefer **invoking hosted skills** over rewriting their workflows inline.
7. **Plan aloud, then execute** — state the next stage before running it.
8. **Comment on the tracking GitHub issue at every gate exit** (and when
   individual Answers/ADRs lock), not only at the end. Local artifacts are not
   enough.
9. Design forks that are not “just AC wording” (semantics, out-of-scope
   product choices) MUST get an **ADR** (or explicit deferral) before Gate B
   advances.

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
| Optional ADR | `docs/adrs/adr-NNN-*.md` (when Gate A locks a design fork) |

---

## Stage A — Clarify + RFC2119 AC

1. **Obtain the issue** (in order):
   1. Explicit number, issue URL, pasted body, or board deep-link that embeds
      an issue → fetch with `gh issue view <N> --repo <owner>/<repo>` (default
      repo = this workspace remote).
   2. Else ask once for a number/URL/paste.
   3. Do **not** invent an issue from `gh issue list` / search.
2. **Follow `clarify-issue`** end-to-end (explore + clarification doc).
3. **Additionally** ensure `clarification.md` contains all of the following
   before presenting Gate A (append if `clarify-issue` omitted them):

   ### Acceptance Criteria (RFC2119)

   Follow `references/rfc2119-ac.md`:
   - Derive MUST / SHOULD / MAY / MUST NOT from the issue + exploration.
   - Flag gaps: implied-but-unstated behaviour → clarifying question **or**
     provisional SHOULD with `*(provisional — confirm)*`.
   - Each MUST/SHOULD line MUST include a **Verify** hint (test name, command,
     or observable).

   ### Dangerous status quo / false-success probes

   List inputs or behaviours that today **succeed wrongly** (wrong AST, wrong
   value, silent accept) or that must become hard errors. This seeds Phase 3
   GAN. If none exist, write that explicitly.

   ### Completion strategy

   Short Phase 0→3 (or equivalent) plan: clarify/AC → spec(+ADRs) → TDD slice
   order → GAN close. Name the intended capability/slice order when the work
   spans the pipeline.

4. Present Gate A — do **not** advance without approval.

### Gate A

> Clarification + AC ready at `<ARTIFACT_DIR>/clarification.md`.
> Review MUST/SHOULD/MAY, dangerous status quo, and completion strategy.
> How proceed?

Options:

- Answer clarifying questions now
- **Propose defaults** for unanswered questions (agent shows a defaults table) →
  user says “defaults OK” / edits → lock under `## Answers` → revise AC →
  **re-present Gate A**
- AC + strategy look good — advance to plan/breakdown
- Pause (I'll review the doc)

**On answers (partial or full):**

1. Record under `## Answers` in `clarification.md` (attribute by Q number).
2. Revise AC / out-of-scope / open items as needed.
3. If a locked answer is a design fork → draft or update an ADR before Gate B.
4. Comment on the tracking issue with what locked.
5. Re-present Gate A until the user advances.

Do **not** treat “advance” as valid while blocking Open Items remain unanswered
unless the user explicitly accepts provisional defaults for them.

**Gate A exit comment** on the tracking issue: locked Answers (or defaults),
AC summary, ADR links, pointer to `clarification.md`.

---

## Stage B — Task breakdown (with verify)

1. For large/ambiguous issues, optionally write `plan.md` (steps + verification
   + success criteria) before the breakdown.
2. Always run **`breakdown-issue`** into `breakdown.md` with:
   - Each task mapped to one or more AC ids when AC exist
   - Each task ending with `### Verify` per `references/task-verify.md`
   - Ordering that respects `AGENTS.md` (spec before code; lexer → parser →
     eval → CLI when spanning the pipeline) and the clarification **Completion
     strategy** slice order
   - False-success probes from clarification reflected in Verify or a dedicated
     closing task (GAN / goldens)
3. Present Gate B.

### Gate B

> Plan/breakdown ready. Task list summary:
> 1. …
> Confirm plan + breakdown?

Options:

- Looks good — cut feature branch (if needed) and implement Task 1
- Adjust tasks / AC (describe changes) → update artifacts → re-present Gate B
- Pause

**Gate B exit comment** on the tracking issue: task list summary, locked
process (spec→TDD per slice→GAN), branch name if created, pointer to
`breakdown.md`.

Only after Gate B approval: create `feat/<issue>-…` (or equivalent) and start
implementation.

---

## Stage C — Implement (AGENTS.md) + Phase 3 GAN

After Gate B approval:

1. For each task (or coherent slice):
   - Phase 1 spec if required (red-phase tests may land in the same slice)
   - Phase 2 TDD: **confirm red** → implement → green
   - Run that task’s `### Verify` before marking it done
   - Conventional commit when the user asks (prefer one commit per green slice
     when working a milestone)
2. Gate C after each slice (or on request): next task / stop / revisit AC
   (back to A or B).
3. When the milestone/slice is complete: **Phase 3 — GAN-style closing review**
   (process Phase 3 in root `AGENTS.md`, not `docs/plan.md` Evaluator phase).

### Phase 3 — GAN discriminator table

Walk every row. Fix punch-list items or defer with an issue/ADR link. Do not
call the milestone done while a row is open without an explicit deferral.

| Pass | Discriminator asks |
| ---- | ------------------ |
| Spec vs plan / evaluatex | Undocumented divergence? Missing MUST/MUST NOT? |
| Spec vs tests | Every normative rule pinned? |
| False-success hunt | Inputs from clarification’s dangerous-status-quo list now correct or hard-error? |
| Dual tables / depth | lexer↔parser arity sync; delimiter depth interactions; related invariants |
| Docs drift | roadmap, gap-analysis, examples, how-to, index, log |

Also: wiki lint on touched `docs/` pages; package + CLI goldens green.

**Gate C / close comment** on the tracking issue: GAN checklist with each row
checked or deferred; what shipped vs deferred; PR/branch link.

---

## Resume

If `ARTIFACT_DIR` already has files:

1. Read `clarification.md` / `breakdown.md` / `plan.md`
2. Resume at the earliest incomplete gate
3. Ask before overwriting
4. Prefer matching the tone/structure of `.agents/issues/issue-8/` when unsure

---

## Anti-patterns

- Skipping Gate A because "the issue was clear"
- Advancing Gate A with unanswered blockers and no defaults-OK lock
- Clarification without **Completion strategy** or **Dangerous status quo**
- Tasks without `### Verify`
- AC written as user stories without MUST/SHOULD
- Implementation or feature-branch cut before Gate B
- Implementation before spec when `AGENTS.md` requires a spec
- Phase 3 that only skims docs and skips the false-success probe set
- Reimplementing clarify/breakdown inside this skill instead of invoking them
- Guessing the "next" issue from search when none was provided
- Local artifacts only — no gate-exit issue comments
