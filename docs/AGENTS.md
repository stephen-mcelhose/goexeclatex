# Wiki Schema

This wiki is maintained by an LLM using the llm-wiki skill
(https://gist.github.com/karpathy/442a6bf555914893e9541c11519de94f).

## Wiki Root

`docs/` (the directory containing this file) is `WIKI_ROOT`.  All paths below
are relative to it.

## Domain

Engineering knowledge base: Go, LaTeX, compilers, parsers, math evaluation,
language design, and software architecture decisions encountered during
development of goexeclatex and adjacent projects.

## Directory Structure

Pages are organised into subdirectories by kind.  The subdirectory **is part
of the slug** (e.g. `specs/lexer`, `adrs/adr-001-char-mode`).

| Directory  | Purpose                                              |
| ---------- | ---------------------------------------------------- |
| `specs/`     | RFC 2119-style normative specifications              |
| `adrs/`      | Architecture Decision Records                        |
| `runbooks/`  | Operational procedures with steps and verification   |
| `raw/`       | Immutable raw sources (LLM reads; never writes)      |
| *(root)*     | Concept pages, proposals, spikes, how-tos            |

New subdirectories MAY be added as the domain grows; update this table when
one is added.

## Conventions

- **Page slugs**: `kebab-case`, relative to `WIKI_ROOT` including the
  subdirectory — e.g. `specs/lexer`, `adrs/adr-001-char-mode`, `plan`.
- **Filenames**: `<slug>.md` — the slug maps 1-to-1 to the file path under
  `WIKI_ROOT`.
- **Frontmatter**: OKF — required `type`, `title`, `description`, `timestamp`
  (ISO-8601 UTC); optional `resource`, `tags`.
- **Cross-references**: `[[slug]]` wikilinks where slug includes the
  subdirectory prefix when the target is not in the root —
  e.g. `[[specs/lexer]]`, `[[adrs/adr-001-char-mode]]`.
- **Sources section**: every page ends with `## Sources` listing its raw inputs.

## Operations

Run these via the `llm-wiki` skill:

- `ingest <source>` — read a new source, write a summary page, propagate to related pages
- `query <question>` — synthesize an answer from wiki pages, optionally write back
- `lint` — audit for orphans, contradictions, stale claims, missing links

## Raw Sources

Raw source files live in `raw/`. They are immutable — the LLM reads them but never writes to them.

## index.md

Structured catalog of all wiki pages. Updated on every write operation.

## log.md

Append-only chronological log. Format: `## [YYYY-MM-DD] operation | detail`
