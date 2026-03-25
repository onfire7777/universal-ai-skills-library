# Workflow Comprehension

Evaluate whether users can understand the system, follow their work through it, and build a mental model that makes them productive. This goes beyond individual screens to assess the **connective tissue** between pages, sections, and actions.

## The Core Question

> If someone sat down with this app on their first day at work, could they figure out how to do their job? Not just "click the right buttons" but actually understand what the system is for, how their work flows through it, and what they should do next at every point.

## Workflow Threads

A "thread" is a sequence of actions that makes up a real task in the user's day. Before evaluating, identify the main threads:

**Examples for an insurance app:**
- Thread 1: Client calls about a renewal → find their policy → review details → draft renewal email → send
- Thread 2: New enquiry comes in → create client record → add policy details → set follow-up
- Thread 3: Morning routine → check what needs attention today → work through the list

**Examples for a project management app:**
- Thread 1: Start the day → see what's assigned to me → pick a task → update progress → mark done
- Thread 2: Client asks for a status update → find the project → see overall progress → send summary
- Thread 3: New project comes in → create it → break into tasks → assign team → set deadlines

### How to Identify Threads

1. Look at the app's navigation structure — each top-level section usually maps to a thread
2. Think about the user's actual day, not the developer's route structure
3. Ask: "What would someone do first thing in the morning? What's the most common task? What's the most important task?"

### How to Test a Thread

Follow each thread end to end. At every step, evaluate:

| Question | What you're checking |
|----------|---------------------|
| **Do I know where to start?** | Is the entry point for this thread obvious from the dashboard or nav? |
| **Does the next step suggest itself?** | After each action, is the logical next step visible or do I have to hunt? |
| **Can I see my progress?** | Is there any indication of how far through the workflow I am? |
| **Can I leave and come back?** | If I get interrupted, can I find my place again tomorrow? |
| **Does it end clearly?** | When the thread is complete, do I know it's done? Is there a confirmation or just... nothing? |

## Wayfinding

Can the user orient themselves at any moment?

| Signal | What to check |
|--------|--------------|
| **Current location** | Does the nav/sidebar highlight where I am? Is there a breadcrumb? |
| **Context from URL** | Does the URL tell me where I am? (e.g. `/clients/jenny-obrien/policies` vs `/app/page/3`) |
| **Page title/heading** | Does the heading tell me what I can do here, not just what the page is called? |
| **Back navigation** | Can I always get back to where I came from? Is there a clear "up" direction? |
| **Section grouping** | Are related pages grouped together in the nav? Does the grouping match how I think about my work? |

### Wayfinding Red Flags

- Same page title on multiple different pages
- Nav doesn't highlight current page
- No way to get back to a parent page except the browser back button
- Breadcrumbs that show route structure instead of logical hierarchy
- Sidebar sections that don't match the user's mental model of their work

## Visual Cues and Nudges

The app should guide users through visual design, not just labels.

| Cue type | What to look for |
|----------|-----------------|
| **Primary action prominence** | Is the main thing I should do on this page visually dominant? (Size, colour, position) |
| **Status indicators** | Can I tell at a glance which items need attention? (Colour coding, badges, icons) |
| **Empty state guidance** | When a section has no data, does it tell me what to do? ("No policies yet — create one") |
| **Progressive disclosure** | Does the app show me what I need now and hide complexity until I need it? |
| **Affordances** | Do clickable things look clickable? Do input fields look like input fields? |
| **Grouping and proximity** | Are related actions near each other? Is there visual separation between unrelated things? |
| **Notification/attention magnets** | Are there badges, counts, or indicators that pull me toward things that need action? |

### Nudge Quality Check

For each page, ask:
1. If I squint at this page, what draws my eye first? Is that the right thing?
2. Are there visual "dead zones" where important actions are hidden?
3. Does the colour/weight hierarchy match the importance hierarchy?
4. Could someone who doesn't read English still roughly understand the page structure from visual cues alone?

## Page-to-Page Continuity

The transitions between pages are where comprehension breaks down most often.

| Transition | What to check |
|-----------|--------------|
| **After creating something** | Am I taken to the thing I just created? Or dumped back at a list? |
| **After saving/updating** | Is there confirmation? Can I see what changed? |
| **After deleting** | Where do I end up? Does the list update? Is the deleted item gone? |
| **Navigating from a list to a detail** | Does the detail page have context about where I came from? Can I get back to the filtered list? |
| **Moving between sections** | If I'm working on a client and switch to emails, does the email section know which client I was looking at? |
| **Deep linking** | If someone sends me a link to a specific record, do I have enough context when I land there? |

### Continuity Red Flags

- Creating a record and being sent to a different section entirely
- Losing filter/sort state when navigating back to a list
- Detail pages with no "back to list" or parent navigation
- Sections that feel disconnected — no cross-references between related data
- Having to re-enter context (re-select a client, re-apply a filter) when switching between related pages

## Mental Model Alignment

Does the app's structure match how the user thinks about their work?

### How to Test

1. **Name the sections**: Read the nav labels out loud. Would the user use these words to describe their work? Or are they developer terms?
2. **Draw the relationships**: If a client has policies, and policies have documents, and documents have versions — does the app's navigation reflect these relationships?
3. **Find the mismatch**: Where does the app organise by technical entity (Users, Records, Settings) instead of by task (Manage Clients, Handle Claims, Run Reports)?
4. **Check the verbs**: Are the action buttons labelled with what the user wants to achieve ("Send Renewal") or with what the system does ("Create Email")?

### Common Model Mismatches

| App says | User thinks |
|----------|------------|
| "Create Record" | "Add a new client" |
| "Entities" | "People" or "Companies" |
| "Workflow" | "What I need to do today" |
| "Submit" | "Save" or "Send" or "Done" — submit to where? |
| "Dashboard" | "My overview" or "What needs attention" |

## Productivity Assessment

After following all threads, ask:

1. **Time to value**: How long before a new user can do something useful? Minutes? Hours? Days?
2. **Repeat efficiency**: Does the second time doing a task feel faster than the first? Are there shortcuts for power users?
3. **Error recovery cost**: When something goes wrong, how many steps to get back on track?
4. **Context switching cost**: When the user moves between threads, how much re-orientation is needed?
5. **Completeness**: Can the user finish their real-world task entirely in this app, or do they need to jump to email/spreadsheets/another tool for some steps?

## Reporting Workflow Issues

Workflow comprehension issues are often **High** severity even though they're not bugs. A user who can't figure out how to use the app won't file a bug report — they'll just stop using it.

When reporting, frame the issue as what the user experiences:

```
BAD:  "Nav labels don't match domain terminology"
GOOD: "A new insurance broker wouldn't recognise 'Entities' as the place to manage their clients.
       They'd look for 'Clients' or 'Contacts'. The current label requires insider knowledge."
```

```
BAD:  "No post-create redirect"
GOOD: "After creating a new policy, the user is sent back to the policy list instead of the
       policy they just created. They have to hunt through the list to find it and continue
       adding details. This breaks the natural workflow of 'create → configure → review'."
```
