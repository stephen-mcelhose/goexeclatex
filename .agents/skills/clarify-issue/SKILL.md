---
name: clarify-issue
description: >-
  Fetch a GitHub issue, explore how it touches the codebase, surface assumptions
  and ambiguities, and produce a clarification doc. Use when the user says
  clarify issue, understand issue, start issue, begin work, or ac-plan-loop
  Stage A.
license: MIT
metadata:
  version: "1.2.1"
  hosted-in: goexeclatex
---

# Clarify Issue

Takes a GitHub issue from raw text to a structured clarification document. The
goal is to surface **what the issue doesn't say** — hidden assumptions,
ambiguous scope, missing context — before any planning or coding begins.

## Workflow

### 0. Establish Artifact Directory

The artifact directory is named `issue-<IDENTIFIER>/` where `<IDENTIFIER>` is:

- A GitHub issue number (e.g. `issue-8`) — most common
- A kebab-case slug when no number exists (e.g. `issue-add-json-output`)

`<IDENTIFIER>` is not known until Step 1. Finalise paths as soon as the number
or name is determined, then announce before proceeding.

```
ARTIFACT_DIR = <workspace>/.agents/issues/issue-<IDENTIFIER>/
EXPLORE_OUT  = <ARTIFACT_DIR>/explore.md
```

Announce: "Artifacts will be written to: `<ARTIFACT_DIR>`"

Create the directory if missing. Do **not** write under personal
`~/.config/...` paths for this repo.

### 1. Get the Issue

- Check if the user provided an issue number, URL, or text description.
- Resolve `OWNER/REPO` and `NUMBER` before fetching:
  - Full issue URL `https://github.com/<owner>/<repo>/issues/<N>` →
    `--repo <owner>/<repo>`, number `<N>`
  - Board deep-link with `issue=<owner>%7C<repo>%7C<N>` (or `|`) → same
  - Bare number only → default to **this repo** (`gh repo view --json
    nameWithOwner -q .nameWithOwner` from the workspace). Ask only if the
    remote is ambiguous or the user named another repo.
- **Always** pass `--repo`. Never run bare `gh issue view <N>`.
- If a number or URL was provided, fetch it:

  ```bash
  gh issue view <NUMBER> --repo <owner>/<repo> \
    --json number,title,body,labels,assignees,comments
  ```

- If not provided, ask for a number, URL, or paste.
- Store the issue title, body, and `ISSUE_REPO`.
- If comments exist, note only those that materially change scope — skip
  resolved discussions and +1s.
- **`gh` failure modes**:
  - Agent-shell proxy `Forbidden` while the user's terminal works → ask them
    to write JSON into the artifact dir, then read it:

    ```bash
    mkdir -p <ARTIFACT_DIR>
    gh issue view <NUMBER> --repo <owner>/<repo> \
      --json number,title,body,labels,assignees,comments \
      > <ARTIFACT_DIR>/issue.json
    ```

    Do **not** tell them to `gh auth login` when their own `gh` already succeeds.
  - Auth error on both agent and user terminal → `gh auth login -h github.com`.
  - 404 → confirm number and `--repo`.
  - Other failure → ask the user to paste the issue text.

### 2. Explore the Codebase

Check first: does `EXPLORE_OUT` already exist?

- **Yes** → read it and skip spawning an explore agent. Announce: "Using
  existing exploration from `<EXPLORE_OUT>`."
- **No** → announce "Exploring codebase for relevance to: `<title>`..." and
  run an explore agent (Task/`explore` or equivalent) with this brief:

  Analyze relevance to this issue (title + body). Report: (1) files/packages/
  interfaces likely touched, (2) patterns/conventions to follow, (3) adjacent
  breakage risk, (4) data flow for the affected area, (5) partial/contradictory
  existing code. Be specific — files, signatures, line numbers. Write findings
  to `<EXPLORE_OUT>`.

**WAIT** for completion, then read `<EXPLORE_OUT>`.

### 3. Synthesize the Clarification Document

Do **NOT** restate the issue verbatim. Each section must be earned by thinking.

#### 3a. What We Know

The issue in your own words. One paragraph.

#### 3b. Codebase Touchpoints

| File / Package | Role | Why Affected |
| -------------- | ---- | ------------ |
| ...            | ...  | ...          |

#### 3c. Assumptions Being Made

Every implicit assumption that, if wrong, breaks the implementation — data,
state, and scope.

#### 3d. Clarifying Questions

Minimum 3. Cut any question whose answer would not change the implementation.
Questions answerable by reading the code are not clarifying questions.

**Q1 — [Short label]**
> [The question]
*Why it matters: [one sentence]*

Prioritise by impact.

#### 3e. Out of Scope (Proposed)

Related work explicitly not being done based on the issue as written.

#### 3f. Open Items

Blockers to starting breakdown.

### 4. Write the Artifact

Write to `<ARTIFACT_DIR>/clarification.md`.

If the file already exists, **append** an `## Updated [date]` section rather
than overwriting.

### 5. Present and Discuss

- Tell the user where the artifact was written.
- Summarize the most impactful clarifying questions in chat (all if ≤5, else
  top 5).
- Ask how to proceed: answer now / review doc / looks good — next steps.

**If the user answers in chat:** record under `## Answers` in
`clarification.md`, then re-present the gate.

**If "Looks good":** say run `/breakdown-issue` or continue via `/ac-plan-loop`
Stage B. Do **not** call `breakdown-issue` yourself unless the user asked for
the full loop.

## Quality Bar

A good clarification document makes it possible to write a task breakdown
**without going back to the issue or re-reading the code**.

## Safety

- Do not modify any source code.
- Do not create branches or commits.
