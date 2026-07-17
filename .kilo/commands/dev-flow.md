---
description: Multi-agent dev workflow — code designs/implements, debug reviews/tests. Loops until approved.
agent: code
---

Orchestrate a dev loop between `code` (design+implement) and `debug` (review+test+decide). Max 3 rounds.

## Round 1
**Step 1** — task agent `code`: `$ARGUMENTS`
**Step 2** — task agent `debug`: "Review the changes. Reply with one line: `PASS: <commit-msg>` or `FAIL: <specific issues>`. Append `|END` if no further work is needed."

## Round 2+ (if FAIL)
**Step 1** — task agent `code`: "Fix these issues: <debug feedback from prior round>"
**Step 2** — task agent `debug`: "Re-review fixes. `PASS: <commit-msg>` or `FAIL: <issues>`. Append `|END` if done."

## On PASS
- If `|END` is present → stop the workflow.
- Otherwise → commit with the debug agent's commit message, then report "Ready for next task" to the user.

All task prompts must be ≤4 sentences. Reference file paths, not paste content.
