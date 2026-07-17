---
description: Multi-agent dev loop — code-impl designs/implements, debug-review reviews/tests/commits. Max 3 rounds.
---

Execute a dev loop using only the `code-impl` and `debug-review` subagents. Max 3 rounds.

**Round 1**
1. Task `code-impl`: $ARGUMENTS
2. Task `debug-review`: Review all uncommitted changes and run tests.

**Round 2+ (only if debug-review returned FAIL)**
1. Task `code-impl`: Fix these issues: <exact debug-review FAIL text>
2. Task `debug-review`: Re-review all uncommitted changes and run tests.

**Outcome**
- `debug-review` commits on PASS — do nothing, it already committed.
- If result line contains `|END` → stop the workflow.
- If PASS without `|END` → report "Ready for next task" to the user.

Keep each task prompt ≤4 sentences. Reference files by path, never paste content.
