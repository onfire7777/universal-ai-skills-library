# Scenario Tests

Six structured tests that go beyond page-by-page evaluation. Each simulates a real-world situation that exposes problems traditional QA misses. Run all six during thorough mode.

---

## 1. New Hire Onboarding Test

**Premise**: You just started this job. Nobody has shown you the app. Figure it out.

**How to run**:
1. Navigate to the app's entry point (login page or dashboard)
2. Do NOT read any CLAUDE.md, README, or source code. No insider knowledge.
3. Try to answer these questions purely from using the app:
   - What is this app for?
   - What are the main things I can do?
   - Where do I start?
   - How do I do the most common task?
4. Time yourself (count clicks/pages visited). How many screens before you can do something useful?
5. Note every moment of confusion — that's a real onboarding failure

**What to report**:
- **Time to first value**: How many clicks/pages before you accomplished something real?
- **Self-explanatory score**: Could you figure out the app without help? (1-5 scale)
- **Missing onboarding**: Is there a welcome screen, guided tour, or empty-state guidance? Or just a blank dashboard?
- **Terminology barriers**: Any labels or section names you didn't understand?
- **Dead ends**: Places where you got stuck with no indication of what to do next?

**Severity guide**: If a new user can't figure out the primary task within 5 minutes of clicking around, that's a **High** finding.

---

## 2. Interrupted Workflow Test

**Premise**: Real users don't work in uninterrupted sessions. They start something, get a phone call, close the laptop, come back hours later.

**How to run**:
1. Start a multi-step task (creating a record, filling a long form, configuring settings)
2. Get halfway through — don't finish
3. Simulate interruption:
   - Navigate away to a completely different section
   - Close the tab and reopen the app from the URL bar
   - Refresh the page mid-form
   - Open a different page in the same tab, then hit back
4. Try to resume the task

**What to report**:

| Interruption | Data preserved? | Could resume? | Notes |
|-------------|----------------|---------------|-------|
| Navigate away mid-form | Yes/No | Yes/No | |
| Close tab + reopen | Yes/No | Yes/No | |
| Page refresh | Yes/No | Yes/No | |
| Back button | Yes/No | Yes/No | |

Also check:
- **Draft/autosave**: Does the app save partial work? Is there a "drafts" or "recent" section?
- **Return navigation**: Is there a "continue where you left off" or "recently viewed" feature?
- **Warning on abandon**: Does it warn you before losing unsaved work?
- **Session persistence**: After reopening, are you still logged in? Does the app remember your last location?

**Severity guide**: Losing form data silently is **Critical**. No way to find where you left off is **High**.

---

## 3. Explain It to a Colleague Test

**Premise**: After using the app, write a 2-minute guide for a colleague. If you can't explain it simply, the app is too complicated.

**How to run**:
1. Complete the main workflow threads (from the workflow comprehension step)
2. Now write a brief, plain-English explanation for each thread:
   ```
   "To [accomplish task], go to [section], click [button], fill in [fields],
   then click [action]. You'll see [confirmation]. The [thing] will appear
   in [location]."
   ```
3. Evaluate what you just wrote

**What to report**:
- **Threads that were easy to explain**: These are well-designed workflows
- **Threads that required caveats**: "But first you need to..." or "Make sure you don't..." signals friction
- **Threads you couldn't explain clearly**: Where you had to say "then you kind of need to figure out..." — that's a comprehension failure
- **Jargon you had to translate**: Words in the app that you replaced with plain language in your explanation
- **Steps you had to invent**: Things the app doesn't guide you through that you had to figure out

**Output**: Include the actual written explanations in the report. They double as user documentation — and the gaps in them map directly to UX improvements needed.

**Severity guide**: Any thread you can't explain in under 5 steps without caveats is a **Medium** finding. Any thread you can't explain at all is **High**.

---

## 4. Wrong Turn Recovery Test

**Premise**: Users click the wrong thing constantly. How forgiving is the app?

**How to run**:
1. For each main section of the app, deliberately go to the wrong place:
   - Click Policies when you meant Clients
   - Open the wrong record from a list
   - Start creating the wrong type of item
   - Apply wrong filters to a list
2. Now try to recover — get back to what you actually wanted to do

**What to report**:

| Wrong turn | Recovery method | Steps to recover | Context lost? |
|-----------|----------------|-----------------|---------------|
| Wrong nav section | [back button / nav click / breadcrumb] | [count] | [Yes/No — what was lost] |
| Wrong record opened | [back to list / browser back] | [count] | [filters lost? scroll position?] |
| Wrong form started | [cancel / back / close] | [count] | [any data lost from previous work?] |
| Wrong filter applied | [clear / reset / undo] | [count] | [other filters preserved?] |

Also evaluate:
- **Cancel safety**: Does cancelling a form always work without side effects?
- **Undo availability**: Can you undo any action? Which ones are irreversible?
- **Navigation memory**: When you go back to a list, are your filters/sort/scroll position preserved?
- **Confirmation dialogs**: Do destructive actions have confirmations? Are they clear about what will happen?

**Severity guide**: Losing work because of a wrong turn is **Critical**. More than 3 clicks to recover from a wrong nav click is **High**. Lost filter state is **Medium**.

---

## 5. Day Two Test

**Premise**: Day one is about discovery. Day two is about efficiency. Does the app reward returning users?

**How to run**:
1. After completing the full walkthrough (you've already used the app extensively)
2. Repeat the main workflow threads from the beginning
3. This time, measure the difference:
   - Is it faster the second time?
   - Are there shortcuts you can now use?
   - Does the app remember your preferences?
   - Can you skip steps that were necessary the first time?

**What to report**:

| Thread | First time (clicks) | Second time (clicks) | Improvement | Shortcuts found |
|--------|-------------------|---------------------|-------------|-----------------|
| [Thread 1] | [count] | [count] | [faster/same/slower] | [list] |
| [Thread 2] | [count] | [count] | [faster/same/slower] | [list] |

Also check:
- **Recent items**: Does the app surface recently viewed/edited items?
- **Defaults**: Does it remember your last choices (filters, sort order, selected view)?
- **Keyboard shortcuts**: Are there any? Are they discoverable?
- **Bulk operations**: For repeat tasks, can you do them in batch?
- **Search**: Can you jump directly to what you need instead of navigating?
- **Bookmarkable deep links**: Can you save links to specific records/views?

**Power user friction**: Things that are fine on day one but maddening on day 100:
- Confirmation dialogs you can't skip ("Are you sure?" for routine operations)
- Mandatory fields that are always the same value
- No way to duplicate a record (have to re-enter everything)
- Pagination that resets when you navigate back

**Severity guide**: No improvement between first and second use is **High** — it means the app doesn't learn from the user or offer any efficiency gains. Missing search on a data-heavy app is **High**.

---

## 6. What Changed? Test

**Premise**: The user logs in after being away for a day/week. Can they tell what happened while they were gone?

**How to run**:
1. Look at the dashboard or landing page after the walkthrough has created/modified data
2. Try to answer: what happened since the last time I was here?
3. Check for these signals:

| Signal | Present? | Useful? |
|--------|----------|---------|
| Activity feed / recent activity | Yes/No | Helpful / Noise / Missing |
| Notification badges/counts | Yes/No | Accurate / Stale / Missing |
| "New since last visit" indicators | Yes/No | |
| Modified/updated timestamps | Yes/No | Visible or buried? |
| Assigned-to-me / needs-attention view | Yes/No | |

**What to report**:
- **Awareness score**: On a 1-5 scale, how quickly could a returning user understand what needs attention?
- **Notification quality**: Do notifications tell you what to do, or just that something happened?
- **Staleness risk**: Could a user miss something important because there's no alert or indicator?
- **Team awareness**: In a multi-user app, can you see what colleagues did?

**Severity guide**: No way to tell what changed since last visit on a team app is **High**. Dashboard that shows static content instead of actionable status is **Medium**.

---

## Running These in Thorough Mode

Run all six tests after the layout sweep and QA sweep (you need the app populated with test data). Recommended order:

1. **New hire onboarding** — do this first while you still have fresh eyes from the layout sweep
2. **Workflow thread testing** — follow the main threads (already in the workflow step)
3. **Interrupted workflow** — test mid-workflow during the thread testing
4. **Wrong turn recovery** — test after you know the right paths
5. **Day two** — repeat the threads, measure improvement
6. **Explain it** — write the explanations last, as a synthesis of everything observed
7. **What changed?** — check the dashboard/notifications view after all the test data creation

The "explain it" test output goes directly into the report and doubles as user documentation draft.
