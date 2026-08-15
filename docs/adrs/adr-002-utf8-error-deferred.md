---
id: adr-002
type: decision
title: "ADR-002: Defer UTF-8 char-mode error message quality to post-Tier-1"
description: Defers improvement of UTF-8 char-mode error messages to post-Tier-1, keeping the Tier-1 scope minimal while preserving a revisit marker.
status: accepted
date: 2026-08-15
timestamp: 2026-08-15T00:00:00Z
tags: [adr, lexer, deferred]
---

# ADR-002: Defer UTF-8 char-mode error message quality to post-Tier-1

## Status

Accepted

## Context

In `nextCharToken`, when an unexpected character is encountered, the error
message formats the offending byte as `string(l.input[l.pos])`. For a multi-byte
UTF-8 rune (e.g. α = 0xCE 0xB1), this produces `\xce` instead of `α`.

## Decision

Do not fix this in Tier 1. LaTeX math expressions are ASCII-dominant — non-ASCII
input in char-mode is already an error condition, so the exact character
representation in the error string is low priority.

A known-limitation comment has been added to `lexer_test.go`. Tracked as GitHub
issue #2.

## Consequences

- Error messages for non-ASCII unexpected characters will show a raw byte escape
  rather than the Unicode character until this is addressed.
- Fix direction (for later): use `utf8.DecodeRuneInString(l.input[l.pos:])` and
  format the rune directly.
- No impact on correctness for Tier 1 use cases.
