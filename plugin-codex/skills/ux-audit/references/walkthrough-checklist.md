# UX Walkthrough Checklist

Evaluate each screen against these categories during a walkthrough. The goal is to **dogfood** the app — use it as a real person with a real job to do, not as a tester running through checkboxes.

## User Persona (Set Before Starting)

Before evaluating any screen, know who you are:

| Question | Why It Matters |
|----------|---------------|
| Who is the user? | A developer and a receptionist notice different things |
| How much time do they have? | Time-poor users won't explore — they need obvious paths |
| How tech-comfortable are they? | Determines tolerance for complexity |
| What device are they on? | Desktop at a desk vs phone between tasks |
| How often do they use this? | Daily power user vs monthly visitor |
| What's their emotional state? | Calm and focused vs stressed and rushing |

Default persona if not specified: "First-time user, moderate tech comfort, slightly distracted, wants to get this done quickly and move on."

## Per-Screen Evaluation

| Category | What to Check |
|----------|---------------|
| **First Impression** | Does the page orient me? Do I know what I can do here? |
| **Navigation** | Can I find what I need in 3 clicks or fewer? Is the current location clear? |
| **Labels & Icons** | Do they describe what they do? Would a first-time user understand them? |
| **Visual Hierarchy** | Is the most important information prominent? Is secondary info de-emphasised? |
| **Call to Action** | Is the primary action obvious? Are secondary actions visually distinct? |
| **Forms** | Are required fields marked? Is validation immediate and clear? Are error messages helpful? |
| **Feedback** | Does the app confirm my actions? Are there loading states? Success/error toasts? |
| **Error Recovery** | Can I undo mistakes? Is there a back button? Are destructive actions guarded? |
| **Consistency** | Same patterns used for similar features? Same terminology throughout? |
| **Data Display** | Tables sorted sensibly? Pagination? Empty states helpful? Long text truncated? |
| **Layout & Visual** | All text visible? Nothing clipped by sidebar or container? Spacing consistent? No overlapping elements? |

## Emotional Friction (Dogfooding Focus)

These go beyond "does it work?" to "how does it feel?"

| Category | What to Check |
|----------|---------------|
| **Trust** | Do I feel confident about what will happen when I click this? Or am I nervous? |
| **Certainty** | After an action, am I sure it worked? Or do I need to go check somewhere else? |
| **Momentum** | Does the app keep me moving forward? Or do I hit walls and have to backtrack? |
| **Cognitive load** | Am I thinking about my task, or thinking about the interface? |
| **Anxiety points** | Are there moments where I'm afraid I'll lose work or break something? |
| **Satisfaction** | After completing the task, do I feel accomplished or exhausted? |

### Trust Signals to Look For

- Confirmation before destructive actions (delete, overwrite, send)
- Clear indication of saved vs unsaved state
- Undo available for reversible actions
- Preview before commit (e.g. "This will email 50 people")
- No surprising side effects from simple actions

### Anxiety Red Flags

- No autosave on long forms
- Ambiguous button labels ("Process", "Submit", "Continue" — continue to WHERE?)
- Destructive actions styled the same as safe actions
- No way to see what just happened (no activity log, no confirmation)
- Modal dialogs with no escape route

## Flow Efficiency

Track the mechanical cost of completing the task:

| Metric | What to Measure |
|--------|----------------|
| **Click count** | Total clicks from start to task completion |
| **Decision points** | Moments where the user has to stop and think about what to do next |
| **Dead ends** | Times the user goes down the wrong path and has to backtrack |
| **Redundant steps** | Actions that feel unnecessary (re-entering info, confirming the obvious) |
| **Shortcuts available?** | For repeat users — keyboard shortcuts, bulk actions, defaults, recent items |

A good flow feels like sliding downhill. A bad flow feels like climbing stairs. Count the stairs.

## Resilience Testing

What happens when things go wrong or the user behaves unexpectedly?

| Scenario | What Should Happen |
|----------|-------------------|
| Navigate away mid-form | Data preserved, or clear warning before losing it |
| Submit with bad/missing data | Specific error messages on the right fields, form state preserved |
| Hit the back button | Sensible navigation, not a broken state |
| Refresh the page | State survives (especially filters, scroll position, form data) |
| Double-click a submit button | No duplicate submission |
| Slow/no network | Graceful degradation, not a white screen |
| Very long text input | Handled gracefully (truncated, scrolled, not overflowing) |

## Visual & Layout Inspection

Catch layout bugs that break the visual presentation. These are easy to miss because the app still "works" — but it looks broken.

**On every page, actively look for:**

| Issue | What to look for | Common cause |
|-------|-----------------|--------------|
| **Clipped/truncated text** | Headings, labels, or descriptions cut off by a container edge. First characters missing. | Content area missing left margin/padding, or sidebar overlapping main content |
| **Overlapping elements** | Sidebar covering main content. Modals under nav bars. Floating buttons over text. | Z-index conflicts, missing margin on main content, fixed positioning errors |
| **Overflow / scrollbars** | Unexpected horizontal scroll. Content wider than viewport. Rogue scrollbar on a container. | Element with fixed width exceeding parent, uncontained images or tables |
| **Broken grid/alignment** | Cards or columns not aligned. Uneven spacing. Items jumping when data loads. | CSS grid/flex issues, missing min-width constraints, layout shift from async data |
| **Text contrast** | Text hard to read against background. Light grey on white. Dark text on dark background. | Missing dark mode styles, low-contrast colour choices, text over images without overlay |
| **Misaligned sidebar/content** | Main content starts behind or underneath the sidebar. Content pushed too far right with gap. | Sidebar width not accounted for in main content margin/padding |
| **Broken responsive transitions** | Layout looks fine at desktop and mobile but breaks at tablet widths. Nav items wrapping oddly. | Missing breakpoint styles for mid-range widths |
| **Image issues** | Broken image icons. Images stretched or squished. Oversized images causing slow load. | Missing src, no aspect-ratio/object-fit, unoptimised originals |
| **Invisible elements** | Buttons or links that exist but can't be seen (same colour as background). | Dark mode missing styles, transparent text, hidden by z-index |
| **Spacing inconsistency** | Some sections have generous padding, others are cramped. Cards with different internal spacing. | Inconsistent use of spacing utilities, missing design tokens |
| **Form input/button mismatch** | Input fields and their submit buttons are different heights, or misaligned horizontally. Search bar taller than its search button. Inline form elements not vertically centred. | Missing matching height utilities, inconsistent padding on inputs vs buttons |
| **Proportional imbalance** | A tiny action button next to a huge empty area. A narrow sidebar with a vast unused content zone. Header taking up half the viewport. | Missing max-width constraints, fixed sizes not scaling, poor use of available space |
| **Collapsible section issues** | Accordion/collapsible panels that don't indicate they're expandable. Sections that collapse and lose the user's scroll position. Nested collapsibles that fight each other. | Missing expand/collapse indicators (chevrons), broken state management |
| **Alignment across components** | Labels don't align with their inputs. Left edges of cards in a grid don't line up. Table columns have inconsistent text alignment (left vs centre vs right). | Mixed alignment rules, inconsistent text-align, grid items not aligned to baseline |
| **Touch target sizing** | Buttons, links, or checkboxes too small to tap on mobile. Icon-only buttons without adequate padding. Close/dismiss buttons that are tiny. | Missing min-height/min-width (44px minimum), padding too tight on interactive elements |

**How to check**: On each page, before evaluating usability, do a quick visual scan:
1. Does all text render fully? No clipped first/last characters?
2. Does the sidebar and main content sit side by side cleanly?
3. Are all elements inside their containers?
4. Is spacing consistent between similar elements?
5. In dark mode: can you read everything? Any elements disappear?

**Severity guide for visual issues**:
- **Critical**: Content is unreadable or unreachable (text fully hidden behind sidebar)
- **High**: Content is partially clipped or overlapping (first letters cut off, as in the screenshot issue)
- **Medium**: Visual inconsistency that doesn't block usage (uneven spacing, slight misalignment)
- **Low**: Minor polish (1-2px alignment, subtle colour inconsistency)

## Cross-Cutting Checks

| Category | What to Check |
|----------|---------------|
| **Mobile (375px)** | Touch targets at least 44px. No horizontal scroll. Text readable. Forms usable. |
| **Dark Mode** | All text readable. No invisible elements. Borders/separators visible. Images appropriate. |
| **Keyboard** | Tab order logical. Focus indicators visible. Modals trap focus. Escape closes dialogs. |
| **Loading States** | Skeleton screens or spinners. No layout shift when data loads. Buttons disabled during submit. |
| **Empty States** | Helpful message when no data. Clear call to action to add first item. |

## Friction Scoring

When you find an issue, classify it:

| Severity | Definition | Example |
|----------|-----------|---------|
| **Critical** | User cannot complete their task | Submit button does nothing, form data lost |
| **High** | User is confused or takes wrong path | Unclear labels cause wrong selection, no undo on delete |
| **Medium** | User succeeds but with unnecessary effort | Required field not marked, have to scroll to find action |
| **Low** | Minor polish issue | Inconsistent capitalisation, alignment off by a few pixels |

## Performance Perception

Not benchmarks — how fast does the app *feel*?

| What to check | Good | Bad |
|--------------|------|-----|
| **Button response** | Instant visual feedback on click (colour change, loading state) | Button does nothing for 1-2 seconds, user clicks again |
| **Page transitions** | Content appears within 300ms, skeleton if longer | Blank white screen while loading, layout jumps when data arrives |
| **List/table loading** | Skeleton rows, then data fills in smoothly | Empty table → sudden appearance of all rows (layout shift) |
| **Form submission** | Button shows spinner, disables, then confirms | "Submitting..." for 5 seconds with no progress indication |
| **Search** | Results appear as you type or within 500ms of pressing enter | Full page reload, or results take 3+ seconds with no feedback |
| **Saving** | Auto-save indicator or instant "Saved" confirmation | No indication whether changes are saved or lost |
| **Navigation** | Feels instant between pages (client-side routing) | Full page reload on every nav click, flash of white |
| **Optimistic UI** | Action appears to succeed immediately, rolls back if fails | Everything waits for server confirmation before showing the change |

**Key question**: At any point, did you wonder "did that work?" or "is it still loading?" — that's a perception problem.

## Copy and Microcopy

The words in the UI matter as much as the layout.

| What to check | Example of good | Example of bad |
|--------------|----------------|---------------|
| **Button labels** | "Save Client" — specific, tells you what happens | "Submit" — submit to where? "Process" — process what? |
| **Error messages** | "Email address is required" on the email field | "Validation error" at the top of the page with no field highlighted |
| **Empty states** | "No policies yet. Create your first policy." with a button | Blank white space, or just a table with column headers and zero rows |
| **Confirmation messages** | "Jenny O'Brien's policy has been saved" — names the thing | "Success" — success at what? "Done" — done with what? |
| **Tooltips and help text** | Short, specific, answers "what is this?" or "why is this here?" | Missing entirely, or repeating the label ("Email: Enter your email") |
| **Placeholder text** | "e.g. smith@company.com" — shows the expected format | "Enter value here" — tells you nothing useful |
| **Destructive action warnings** | "Delete Jenny O'Brien? This cannot be undone. 3 policies will also be removed." | "Are you sure?" — sure about what? What happens if I click yes? |
| **Loading/progress text** | "Sending email to 12 recipients..." — specific, shows progress | "Please wait..." — wait for what? How long? |
| **Navigation labels** | "Clients", "Policies", "Email Outbox" — matches the user's vocabulary | "Entities", "Records", "Queue" — developer vocabulary |
| **Tone consistency** | All messages use the same voice (friendly, professional, casual) | Mix of "Oops! Something went wrong!" and "Error: HTTP 500" |

**Spelling and grammar**: Check every page for typos, especially in error messages, empty states, and tooltips. These are often written once and never reviewed.

## Feedback Completeness

After every user action, the app should confirm what happened. Check the full action lifecycle:

| Action | Expected feedback | Common failure |
|--------|------------------|----------------|
| **Create** | "Created [thing]" + navigate to the new item or show it in the list | Silent success — thing is created but user doesn't know where it went |
| **Save/Update** | "Saved" indicator near the save button, or auto-save badge | No indication — user saves, nothing changes visually, they're unsure |
| **Delete** | Confirmation dialog with specifics, then "Deleted" toast, item removed from view | Item disappears with no confirmation or toast |
| **Send** (email, message) | "Sent to [recipient]" with details | "Message sent" — sent to whom? Did it really send? |
| **Bulk action** | "Updated 12 of 12 items" with count | "Done" — how many? Did they all succeed? |
| **Error** | Specific error on the specific field, form state preserved | Page-level error banner, form cleared, user has to re-enter everything |
| **Long-running action** | Progress bar or step indicator, estimated time | Spinner for 30 seconds with no indication of progress |
| **Background action** | Toast "Processing in background, we'll notify you when done" | Nothing — user doesn't know if it's still running or failed silently |

**The test**: After every action you take during the walkthrough, ask: "Do I know what just happened? Am I confident it worked? Do I know what to do next?"

## Data Edge Cases

Test how the app handles unusual but realistic data:

| Scenario | What to check |
|----------|--------------|
| **Very long text** | Paste a 500-char name, title, or description. Does it truncate, wrap, or overflow? |
| **Special characters** | O'Brien, café, Müller, emoji 🎉, `<script>alert(1)</script>` — does it render or break? |
| **Empty strings** | Submit forms with spaces-only input. Does validation catch it? |
| **Large numbers** | Enter $999,999.99 or 100,000 items. Does formatting hold? |
| **Many items** | If a list/table has 100+ rows, does pagination work? Does it lag? |
| **No data** | What does the page look like with zero records? Helpful empty state or blank void? |
| **Concurrent edits** | Open the same record in two tabs, edit both, save both. What happens? |
| **Rapid clicks** | Click a submit/action button 5 times fast. Does it fire once or create duplicates? |

## Navigation & Routing

| Scenario | What to check |
|----------|--------------|
| **Direct URL access** | Paste a deep link (e.g. `/app/clients/123`). Does it load or redirect to login then back? |
| **Browser back/forward** | After navigating 3 pages deep, does back button work correctly? |
| **Bookmark-able** | Does the URL change when you navigate? Can you bookmark and return? |
| **404 handling** | Navigate to a non-existent route (e.g. `/app/doesnotexist`). Graceful error or white screen? |
| **Stale URLs** | If a record is deleted, what happens when you visit its URL? |

## The Big Questions

Ask these after completing (or failing) the task. These are the dogfooding questions — they go beyond "does it work?" to "is it good?"

1. **Would I come back?** Or would I look for an alternative tool next time?
2. **Could I teach someone else?** In under 2 minutes, with no documentation?
3. **What's the one thing?** If I could change one thing to make this twice as easy, what would it be?
4. **Where did I hesitate?** Every hesitation is a design failure — the interface made me think instead of act.
5. **Did I trust it?** At any point, was I unsure whether my data was saved, my action was processed, or my input was correct?
6. **What would a busy person skip?** Which steps feel optional, tedious, or "I'll do that later"?
7. **Would I recommend this?** Not just "it works" but "you should use this, it's good."
