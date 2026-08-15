# Task verification format

Every task in `breakdown.md` MUST end with:

```markdown
### Verify
- **Command:** `go test ./internal/parser/... -count=1`
- **Expect:** PASS (or specific failing stubs until implement: undefined Symbol)
- **Covers:** AC-MUST-1, AC-SHOULD-2
- **Tests:** Prefer table-driven cases for multi-input behavior; name the test function if known
```

## Preferences

1. **Automated over manual** — `go test`, package-scoped first.
2. **Table-driven** when ≥2 cases share the same assertion shape.
3. **Smallest package scope** that proves the task — avoid repo-wide `./...`
   until a final gate.
4. **CLI goldens** when the AC is end-to-end user-visible behaviour.
5. If Verify is manual-only, state the exact observable (stderr substring,
   wiki page link, issue comment URL).

## Bad vs good

Bad: `### Verify` — "make sure it works"
Good: `### Verify` — `go test ./internal/eval/... -run TestOddRoot -count=1` covers AC-MUST-3
