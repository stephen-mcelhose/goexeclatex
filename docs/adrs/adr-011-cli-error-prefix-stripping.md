---
id: adr-011
title: "ADR-011: CLI strips internal package prefixes from error messages"
status: accepted
date: 2026-08-15
---

# ADR-011: CLI strips internal package prefixes from error messages

## Status

Accepted

## Context

The internal packages follow Go convention and prefix their error strings with
the package name:

| Package          | Example error string                            |
| ---------------- | ----------------------------------------------- |
| `internal/lexer` | `lexer: unclosed brace group`                   |
| `internal/parser`| `parser: unexpected token ...`                  |
| `internal/eval`  | `eval: division by zero`                        |
| `internal/eval`  | `eval: domain error: sqrt of negative`          |
| `internal/eval`  | `eval: undefined symbol: x`                     |

These prefixes are useful for Go callers (tests, library users) to identify
the origin of an error programmatically.

However, they are noise at the CLI surface. A user who types
`echo '1/0' | goexeclatex` does not benefit from seeing:

```
error: eval: division by zero
```

The double prefix (`error:` + `eval:`) is redundant and looks unpolished.

Two options were considered:

**Option A — Strip the package prefix, re-wrap with `error:`**
```
error: division by zero
```
Clean, user-facing. Requires a small string-stripping step at the CLI boundary.

**Option B — Pass the raw error string through**
```
error: eval: division by zero
```
Zero transformation. Slightly exposes internal structure to the user.

## Decision

**Option A.** The CLI **MUST** strip the leading `<package>: ` prefix from
internal error strings before printing. The output format is:

```
error: <message>
```

Implementation:

```go
func userMessage(err error) string {
    s := err.Error()
    for _, prefix := range []string{"eval: ", "lexer: ", "parser: "} {
        s = strings.TrimPrefix(s, prefix)
    }
    return s
}
```

This is applied at the single point where errors are written to stderr.

## Consequences

- User-facing error messages are clean: `error: division by zero`, not
  `error: eval: division by zero`.
- The internal package error strings are **not changed** — they remain
  prefixed for testability and library use.
- The CLI spec §8 documents the stripped format as the contract; internal
  strings are an implementation detail.
- If a new internal package is added (e.g. `internal/formatter`), its prefix
  **MUST** be added to the strip list in `cmd/goexeclatex/main.go`.

## Revisit criteria

If the internal packages switch to structured errors (e.g. a typed `EvalError`
with a `UserMessage()` method), the strip list can be removed in favour of a
type assertion.
