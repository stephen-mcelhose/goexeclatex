# RFC2119 Acceptance Criteria

Use these keywords **exactly** (see [RFC 2119](https://datatracker.ietf.org/doc/html/rfc2119)):

| Keyword | Meaning in this workflow |
| ------- | ------------------------ |
| **MUST** | Required for done; blocking. Every MUST has automated or explicit manual verify. |
| **MUST NOT** | Forbidden behavior; verify with a negative test when feasible. |
| **SHOULD** | Strong default; omit only with documented reason in Answers. |
| **SHOULD NOT** | Discouraged; call out if violated. |
| **MAY** | Optional; no blocking verify required. |

## Section template (append to `clarification.md`)

```markdown
## Acceptance Criteria (RFC2119)

### MUST
- **AC-MUST-1** — …
  **Verify:** `go test ./path/... -run TestName` (table-driven) / or manual: …
- **AC-MUST-2** — …
  **Verify:** …

### MUST NOT
- **AC-MUSTNOT-1** — …
  **Verify:** …

### SHOULD
- **AC-SHOULD-1** — … *(provisional — confirm)*
  **Verify:** …

### MAY
- **AC-MAY-1** — …
```

## Writing tips

- One behavior per bullet. Prefer observable outcomes over implementation detail.
- Tie each MUST to codebase touchpoints when known (`internal/parser`, CLI golden, …).
- If the issue is silent on error handling or reference divergence — ask, or mark provisional SHOULD.
- Do not use "will", "should probably", or "ideally" as soft requirements.
