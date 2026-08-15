---
type: runbook
title: Ingest a Reference Implementation
description: Steps to read a reference implementation into the wiki — clone, read source files, synthesize a spike page, and propagate cross-references.
tags: [llm-wiki, ingest, reference-implementation, runbook]
timestamp: 2026-08-15T00:46:01Z
---

# Ingest a Reference Implementation

Use this runbook whenever a new external library or codebase is designated as
a reference implementation for goexeclatex.

## Prerequisites

- The repo is cloned locally (or is accessible via GitHub).
- The `llm-wiki` skill is available in the session.
- `WIKI_ROOT` is `docs/` (confirmed by `docs/AGENTS.md`).

## Steps

### 1. Clone the repo (if not already local)

```bash
cd ~/repos
git clone <repo-url> <local-name>
# e.g. git clone https://github.com/arthanzel/evaluatex evaluatex
```

### 2. Identify the files to read

Focus on the files that reveal architecture, not tests or build config:

| Priority | What to read                                         |
| -------- | ---------------------------------------------------- |
| High     | `src/` or `lib/` — the core implementation           |
| High     | `README.md` — overview, API surface, known limits    |
| Medium   | Key algorithm files (lexer, parser, evaluator)       |
| Low      | Tests — useful to confirm edge-case behaviour        |
| Skip     | `node_modules/`, build artifacts, lock files         |

### 3. Run the llm-wiki ingest operation

In the agent session, trigger ingest:

```
ingest ~/repos/<local-name>
```

The skill will:
1. Read the source files you identified above.
2. Surface 3–5 key takeaways in chat for your review.
3. Write a `spike` page to `docs/<slug>.md` (e.g. `docs/evaluatex-reference-implementation.md`).
4. Propagate cross-references to related wiki pages.
5. Update `docs/index.md` and append to `docs/log.md`.

### 4. Review the spike page

Check that the generated page covers:

- [ ] Architecture overview (phases: lex → parse → evaluate)
- [ ] Feature coverage (what LaTeX constructs are handled)
- [ ] Known gaps or limitations
- [ ] Cross-references to `[[plan]]`, `[[goexeclatex-gap-analysis]]`, or other relevant pages

Edit the page directly if anything is missing or wrong.

### 5. Run the gap analysis

After ingesting a new reference, update or create a gap analysis page:

```
ingest — gap analysis comparing <reference> vs our spec
```

Target slug: `goexeclatex-gap-analysis` (update in place if it already exists).

### 6. Verify wiki hygiene

```
lint the wiki
```

Confirm no orphans, missing index entries, or broken wikilinks were introduced.

## Notes

- The evaluatex ingest (2026-08-14) read `~/repos/evaluatex/src/` directly —
  specifically `lexer.js`, `Token.js`, `parser.js`, and `math-functions.js`.
- The AMS short-math-guide was read from `docs/raw/short-math-guide.tex` after
  conversion with pandoc (`pandoc short-math-guide.tex -t markdown`).
- Raw source files go in `docs/raw/` and are never modified by the LLM.

## Sources

- `docs/evaluatex-reference-implementation.md`
- `docs/log.md`
- `docs/AGENTS.md`
