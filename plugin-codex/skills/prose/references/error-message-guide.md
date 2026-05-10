# Error Message Guide

エラーメッセージ構造（何/なぜ/次に何を）、重要度別テンプレートのリファレンス。

---

## Error Message Structure

### The Three-Part Framework

```
1. WHAT happened (state the problem clearly)
2. WHY it happened (brief explanation, if helpful)
3. WHAT TO DO next (actionable resolution)

Example:
  WHAT: "Your payment couldn't be processed."
  WHY:  "The card on file has expired."
  NEXT: "Update your payment method to continue."
```

### Writing Rules

| Rule | Good | Bad |
|------|------|-----|
| **Be specific** | "File must be under 10 MB" | "File too large" |
| **Be human** | "We couldn't find that page" | "404 Error" |
| **Blame the system** | "We couldn't save your changes" | "You entered invalid data" |
| **Offer a path forward** | "Try again or contact support" | "An error occurred" |
| **Avoid jargon** | "Something went wrong on our end" | "Internal server error" |
| **Use active voice** | "We couldn't send your message" | "Message could not be sent" |

---

## Severity-Based Templates

### Critical Errors (Data loss, account issues)

```markdown
## Template
[Specific problem statement.]
[Brief cause if it helps the user.]
[Specific recovery action + alternative.]

## Examples
"Your account has been locked after too many login attempts.
Reset your password to regain access, or contact support
if you need help."

"We couldn't save your document due to a connection issue.
Your latest changes are stored locally and will sync
when you're back online."
```

### Validation Errors (Form input)

```markdown
## Template
[What's wrong with the input — be specific.]

## Rules
- Show inline, next to the field
- Use red/error color indicator
- Be specific about the requirement
- Show the requirement, not just the violation

## Examples
✅ "Password must be at least 8 characters"
❌ "Invalid password"

✅ "Enter a valid email (e.g., name@company.com)"
❌ "Invalid email format"

✅ "Username can only contain letters, numbers, and hyphens"
❌ "Special characters not allowed"
```

### System Errors (Server issues, timeouts)

```markdown
## Template
"Something went wrong [optional: specific area].
[Recovery action]. If this keeps happening, [escalation path]."

## Examples
"Something went wrong loading your dashboard.
Refresh the page to try again. If this keeps happening,
check our status page or contact support."

"We're having trouble connecting to the server.
Check your internet connection and try again."
```

### Permission Errors

```markdown
## Template
"You don't have access to [specific resource].
[How to get access]."

## Examples
"You don't have permission to edit this project.
Ask the project owner to grant you editor access."

"This feature is available on the Pro plan.
Upgrade your plan to unlock it."
```

### Not Found Errors

```markdown
## Template
"We couldn't find [what they were looking for].
[Suggestions or alternatives]."

## Examples
"We couldn't find a page at this address.
Check the URL or go back to the homepage."

"No results for 'acme widget'.
Try different keywords or browse categories."
```

---

## Error Message Patterns

### Inline Validation

```
Timing:
  - Validate on blur (when leaving field), not on every keystroke
  - Show success state only after previously showing error
  - Clear error as soon as input becomes valid

Position:
  - Below the field, left-aligned
  - Use aria-describedby for accessibility
  - Red color + icon for error, green + icon for success
```

### Toast/Snackbar Errors

```
Use for:
  - Background operation failures (save, sync, send)
  - Non-critical errors that don't block workflow

Duration:
  - Auto-dismiss after 5-8 seconds
  - Allow manual dismiss
  - Persist if action is needed

Content:
  - One line: "[What failed]. [Action link]"
  - Example: "Couldn't save changes. Try again"
```

### Full-Page Errors

```
Use for:
  - Page load failures
  - Authentication/authorization errors
  - Critical system outages

Include:
  - Illustration (friendly, not dramatic)
  - Clear message (no error codes unless technical audience)
  - Primary action (retry, go home, contact support)
  - Secondary action (status page, help docs)
```

---

## Error Message Checklist

```markdown
- [ ] States the problem clearly (no jargon)
- [ ] Uses active voice, human language
- [ ] Doesn't blame the user
- [ ] Includes a specific recovery action
- [ ] Provides alternative if primary action might fail
- [ ] Appropriate severity and display pattern
- [ ] Accessible (aria roles, color not sole indicator)
- [ ] Tested with real users for clarity
```
