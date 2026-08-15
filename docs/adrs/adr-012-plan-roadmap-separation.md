---
id: adr-012
type: decision
title: "ADR-012: Separate architecture plan from forward-looking roadmap"
description: Splits docs/plan.md into a stable architecture document and a separate docs/roadmap.md that tracks upcoming milestones and current priorities.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T03:04:00Z
tags: [adr, docs, project-management]
---

# ADR-012: Separate architecture plan from forward-looking roadmap

## Status

Accepted

## Context

`docs/plan.md` was written as an implementation proposal and grew to contain
two distinct kinds of content:

1. **Stable architecture** — pipeline design, grammar, package layout, CLI
   interface, error handling policy, testing strategy. These are settled
   decisions that change only when the design changes.

2. **Forward-looking milestones** — v0.2, v0.3, vFuture checklists with
   GitHub issue links and open `- [ ]` items. These change every session as
   work is completed, priorities shift, and new issues are filed.

Mixing them in one document causes two problems:

- **Update frequency mismatch.** Every session that makes progress bumps the
  plan for what is effectively project-management bookkeeping, creating noise
  in the git history of a document that should be largely settled.
- **Audience mismatch.** An agent or developer reading the architecture needs
  the stable design. A PM or agent picking up a new session needs the current
  priorities. A single document serves neither audience well.

## Decision

Split `docs/plan.md` into two documents:

| Document | Type | Purpose | Update frequency |
| -------- | ---- | ------- | ---------------- |
| `docs/plan.md` | `concept` | Stable architecture record — pipeline design, grammar, package layout, CLI interface, error policy, testing strategy, completed milestones for posterity | Rarely; only when design decisions change |
| `docs/roadmap.md` | `proposal` | Forward-looking milestone tracker — upcoming versions with issue links, current focus, process reminder | Every session |

`plan.md` type changes from `proposal` to `concept`, reflecting that the
architecture is implemented and settled, not proposed.

Completed milestones (v0.1 ✅) remain in `plan.md` as a historical record.
Open milestones (v0.2, v0.3, vFuture) move to `roadmap.md`.

`plan.md` links to `roadmap.md` at the end of its milestones section.
`roadmap.md` links back to `plan.md` for architecture context.

## Consequences

- Agents starting a new session open `roadmap.md` to orient on current
  priorities, then follow links to the relevant spec or ADR.
- `plan.md` history in git reflects genuine design changes, not sprint
  bookkeeping.
- `roadmap.md` must be kept up-to-date; stale roadmaps are worse than none.

## Sources

- `docs/plan.md` — document being split
- `docs/index.md` — wiki index updated to register `roadmap`
