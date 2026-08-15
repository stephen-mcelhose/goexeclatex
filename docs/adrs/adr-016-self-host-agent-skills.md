---
id: adr-016
type: decision
title: "ADR-016: Self-host Tier A agent skills (C+B); revisit on degraded experience"
description: >
  Locks the #15 C+B decision: encode the milestone planning loop in AGENTS.md
  and host genericized Tier A skills under .agents/skills/. Wholesale community
  skill packs and Tier B remain out of scope until revisit criteria fire.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T17:43:37Z
tags: [adr, agents, process, skills, self-host]
---

# ADR-016: Self-host Tier A agent skills (C+B); revisit on degraded experience

## Status

Accepted — records the C+B decision on
[#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15), landing via
PR [#21](https://github.com/stephen-mcelhose/goexeclatex/pull/21).

## Context

Milestone [#8](https://github.com/stephen-mcelhose/goexeclatex/issues/8) was
planned and closed with an organic loop: clarify → RFC2119 AC → Gate A →
breakdown with Verify → Gate B → spec → TDD slices → GAN-style Phase 3. That
loop used personal Cursor skills (`ac-plan-loop`, `clarify-issue`,
`breakdown-issue`) plus in-repo `llm-wiki`.

Issue #15 asked whether to self-host those skills. Options considered:

| Path | Meaning |
| ---- | ------- |
| **A** | Document recommended personal skills only |
| **B** | Vendor genericized Tier A into the repo |
| **C** | Encode the loop in `AGENTS.md` without vendoring |
| **Buy (community)** | Adopt packs such as mattpocock/skills wholesale or as a replacement spine |

**Tension:** community packs are more battle-proven *across* repos; the personal
Tier A spine is battle-proven *on this repo’s #8 run*. Organic fit won for the
default spine; community skills stay spare parts, not an automatic replace.

Closing #15 before the hosting PR merged risked “done on paper”; this ADR
pins the durable decision and when to reopen it.

## Decision

1. **C+B (accepted):**
   - **C** — Encode working style, planning loop pointers, and process Phase 3
     (closing / GAN review) in root `AGENTS.md` and `docs/roadmap.md`.
   - **B** — Host genericized Tier A under `.agents/skills/`:
     `clarify-issue`, `breakdown-issue`, `ac-plan-loop`, and canonical
     in-repo `llm-wiki`.
2. **Artifact root:** `.agents/issues/issue-<IDENTIFIER>/` (workspace-local).
   Do not use personal `~/.config/csgdaa-code/...` paths for this repo.
3. **Prefer in-repo skills** when both personal (`~/.cursor/skills/`) and
   hosted copies exist.
4. **Out of scope for now:**
   - Tier B vendor (`cognitive-locality-grouper`, `define-contracts`,
     `issue-workflow`, …)
   - Wholesale adoption of a community skill pack as the planning spine
   - Org-specific defaults (Ready boards, foreign default repos) in hosted skills
5. **Fidelity:** `ac-plan-loop` MUST preserve the #8 discriminators that made
   the loop work (Gate A defaults-OK, completion strategy, false-success
   probes, gate-exit issue comments, Stage C GAN table). See worked example
   `.agents/issues/issue-8/`.

## Revisit criteria

Reopen or supersede this ADR when **any** of the following holds (degraded
experience — not “nice to have” churn):

1. **Loop ignored** — across ≥2 milestone-sized issues, agents systematically
   skip Gate A defaults negotiation, dangerous-status-quo / false-success
   probes, gate-exit issue comments, or Stage C GAN rows despite the hosted
   skills being present.
2. **Shadowing** — personal `~/.cursor/skills/` copies repeatedly win over
   in-repo skills and produce conflicting artifacts or process (e.g. foreign
   `ARTIFACT_DIR`), and prefer-in-repo guidance is insufficient.
3. **Net slowdown** — for ≥2 consecutive milestones, the hosted loop costs
   more wall-clock / reconciliation than clarifying and breaking down by hand
   (or a thinner `AGENTS.md`-only checklist), with concrete session evidence.
4. **Proven better spare parts** — a community skill (or Tier B) is spiked on a
   real issue and measurably reduces planning friction **without** rewriting
   this repo’s planning vocabulary (Gate A/B, RFC2119 AC, Verify-per-task,
   spec→TDD→GAN).

Until then:

- Do **not** unwind C+B “for completeness.”
- Do **not** bulk-replace hosted Tier A with a community pack.
- Do **not** vendor Tier B solely because it exists in personal skill dirs.

## Consequences

- New agents SHOULD start from `AGENTS.md` + `.agents/skills/ac-plan-loop`.
- Process regressions are fixed in hosted skills (and this ADR’s fidelity
  clause), not by silently falling back to personal copies.
- Tier B / community adoption requires superseding or amending this ADR with
  spike evidence against the revisit criteria.
- `#15` remains the historical decision record; this ADR is the durable one.

## Sources

- GitHub issue [#15](https://github.com/stephen-mcelhose/goexeclatex/issues/15) — self-host skills (C+B)
- Pull request [#21](https://github.com/stephen-mcelhose/goexeclatex/pull/21) — host Tier A + encode Phase 3
- Root `AGENTS.md` — planning loop + process Phase 3
- `.agents/issues/issue-8/` — organic worked example
- [[adrs/adr-012-plan-roadmap-separation]] — process/docs ADR precedent
- [[adrs/adr-014-mixed-floor-ceil-round-out-of-scope]] — revisit-criteria shape
