---
name: git-commit
description: Creates clean, conventional git commits from staged files
license: MIT
compatibility: Git 2.x+
user-invocable: true
allowed-tools:
  - bash
  - ask_user_question
---

# Git Commit Skill

Creates clean, conventional git commits from staged files. No pushing — that's manual.

---

## Workflow

### 1. Check what's staged

```bash
git diff --cached --name-only
git diff --cached --stat
```

If nothing is staged, tell the user and stop. Do not auto-stage anything.

### 2. Scan staged content for potential issues

Run:
```bash
git diff --cached
```

Look for these red flags in the diff:

| Category | Examples |
|---|---|
| **Credentials / secrets** | API keys, tokens, passwords, private keys, `.env` values hardcoded in source |
| **Test/debug SQL** | `SELECT *` queries, hardcoded IDs, queries with `LIMIT 1` that look exploratory |
| **Debug artifacts** | `console.log`, `print(`, `debugger`, `TODO`, `FIXME`, `HACK`, `XXX`, `// temp`, `// remove` |
| **Hardcoded debug variables** | Variables with names like `DEBUG_*`, `TEST_*`, `FAKE_*`, `TEMP_*` or containing hardcoded fake/test values (e.g. `'abc-123-fake'`, `'test@example.com'`) — **only flag if the variable is actually used or exported**, not if it's just declared and unused |
| **Test data** | Hardcoded emails, phone numbers, fake names that look like test fixtures left in prod code |
| **Commented-out code** | Large blocks of commented code that look accidental |

If issues are found — regardless of severity (even minor things like `console.log` or a `TODO`):
- List them clearly with file + line context
- Ask the user: **"Found potential issues — do you want to fix them before committing, or proceed anyway?"**
- Wait for confirmation before continuing. Do not proceed until the user explicitly says so.

If no issues found, proceed.

### 3. Understand the changes

From the diff, determine:
- What changed (files, functionality)
- Why it likely changed (feature, fix, refactor, etc.)
- Whether a scope makes sense (e.g., `ui`, `db`, `api`, `auth`)

### 4. Write the commit message

Follow this structure:

```
<type>[optional scope]: <Description starting with capital letter>

[optional body — what and why, not how]

[optional footer — only if user mentioned an issue number]
```

**Types:**
- `feat` — new feature
- `fix` — bug fix
- `refactor` — restructure without changing behavior
- `chore` — dependencies, config, non-src changes
- `perf` — performance improvement
- `ci` — CI/CD changes
- `ops` — infra/deployment
- `build` — build system changes
- `docs` — documentation
- `style` — formatting, whitespace
- `revert` — reverting a previous commit
- `test` — adding or fixing tests

**Rules:**
- Description: imperative mood, capitalized, brief but informative ("Add", "Fix", "Remove" — not "Added")
- Add `!` after type for breaking changes (e.g. `feat!:`)
- Use scope in parentheses if the change is clearly scoped to one area: `feat(ui):`, `fix(db):`
- Body: only include for complex changes (non-obvious reasoning, multiple concerns, breaking changes, or significant architectural decisions). For simple, self-explanatory changes, omit the body entirely. Separate from header with a blank line when included.
- Footer: only include `Closes: #<issue>` if the user explicitly mentioned an issue number. Otherwise omit entirely.
- Do not add co-author lines or mention AI involvement.

### 5. Show the commit message and confirm

Show the user the proposed commit message in a code block. Ask:
**"Does this look good, or would you like any changes?"**

Wait for approval. If they suggest changes, update and show again.

### 6. Create the commit

Once approved:

```bash
git commit -m "<header>" -m "<body if any>" -m "<footer if any>"
```

Or for simple single-line commits:
```bash
git commit -m "<header>"
```

Confirm success by showing the output of the commit command.

---

## Examples

**Simple:**
```
chore: bump version from 1.0.0+8 to 1.0.0+9
```

**With scope:**
```
feat(ui): Add like/dislike buttons on AI chat

On pogovori/[id]/page.tsx added two buttons Like and Dislike which change
colour based on state. On click a supabase row chat_messages also gets updated.

Closes: #CLDA-6138
```

**Breaking change:**
```
fix!: Upgrade Next.js and React versions due to security vulnerabilities

Bumps Next.js and React to latest stable releases to address identified
security issues. The upgrade includes breaking changes in framework APIs,
so existing integrations may require adjustments.

Closes: #CLDA-4556
```

---

## Notes

- Never stage additional files — only commit what was already staged.
- Never push — the user does that manually.
- If `git` is not available or there's no repo, say so clearly.
