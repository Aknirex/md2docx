---
description: Multi-agent dev loop — code implements, debug reviews. Max 3 rounds.
agent: code
---

Loop until approved (max 3 rounds):

1. Round 1: task agent `code` → $ARGUMENTS
2. Every round after code: task agent `debug` → "Review changes. Reply with one line: `PASS: <commit-msg>` or `FAIL: <issues>`. Add ` |END` if no more work needed."
3. If FAIL → next round: task agent `code` → "Fix: <debug feedback>"
4. If PASS |END → stop. If PASS → commit with the suggested message.

All task prompts ≤4 sentences. Reference files by path.
