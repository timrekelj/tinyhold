---
name: git-pr
description: Creates a well-structured pull request from recent commits
license: MIT
compatibility: Git 2.x+
user-invocable: true
allowed-tools:
  - bash
  - ask_user_question
---

# Git PR Skill

Creates a well-structured pull request from recent commits. This skill helps you generate a clear PR title, description, and labels based on commit history.

---

## Workflow

### 1. Check current branch and remote

```bash
git branch --show-current
git remote -v
```

If not on a feature/fix branch or no remote configured, notify the user and suggest setup.

### 2. Gather recent commits

```bash
git log --oneline -n 20
```

Show the user the recent commits and ask which range to include in the PR.

### 3. Analyze commit types

Categorize commits by type (feat, fix, refactor, etc.) to suggest appropriate PR labels and title.

### 4. Generate PR title

Follow conventional commit style:
```
<type>(<scope>): <short description>
```

Examples:
- `feat(auth): Add OAuth2 provider integration`
- `fix(ui): Resolve mobile menu overflow issue`
- `refactor(api): Simplify endpoint response structure`

### 5. Generate PR description

Structure:
```
## Summary
<1-2 sentence summary of changes>

## Changes
- Bullet point list of key changes
- Group by component/area if applicable

## Related Issues
<list any related issue numbers>

## Testing
<brief testing notes if available from commits>
```

### 6. Suggest labels

Based on commit types:
- `enhancement` for feat commits
- `bug` for fix commits
- `refactor` for refactor commits
- `documentation` for docs commits
- `performance` for perf commits

### 7. Show PR preview

Display the complete PR template:
```
Title: <generated title>

Description:
<generated description>

Labels: <suggested labels>
```

Ask user: **"Does this PR look good, or would you like any changes?"**

### 8. Create PR (optional)

If user approves and wants to create it:
```bash
gh pr create --title "<title>" --body "<description>" --label "<label1>,<label2>"
```

Or provide the command for user to run manually.

---

## Notes

- Requires GitHub CLI (`gh`) for automatic PR creation
- Without `gh`, provide the PR template for manual creation
- Never push automatically — user controls when to create PR
- If no commits found, suggest user make commits first
