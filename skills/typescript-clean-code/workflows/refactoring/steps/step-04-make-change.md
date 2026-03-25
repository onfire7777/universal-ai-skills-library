---
name: 'step-04-make-change'
description: 'Make ONE structural change — no behavior change'
nextStepFile: './step-05-run-tests.md'
referenceFiles:
  - 'references/functions/rules.md'
  - 'references/classes/rules.md'
---

# Step 4: Make ONE Change

## STEP GOAL

Execute one small refactoring step. The change must be purely structural with no behavior change.

## REFERENCE LOADING

Before making changes, load the relevant reference files:
- `references/functions/rules.md` — if refactoring functions
- `references/classes/rules.md` — if refactoring classes

Apply the rules to guide the structural change.

## EXECUTION

### Checklist

- Make exactly ONE change
- The change is purely structural (no behavior change)
- You can describe the change in one sentence

### Common Refactorings

**Extract Method**:
```typescript
// Before
function processOrder(order: Order) {
  // validate
  if (!order.items.length) throw new Error('Empty order');
  if (!order.customer) throw new Error('No customer');

  // calculate
  const subtotal = order.items.reduce((sum, i) => sum + i.price, 0);
  const tax = subtotal * 0.1;
  const total = subtotal + tax;

  // save
  db.save({ ...order, total });
}

// After (one extraction)
function processOrder(order: Order) {
  validateOrder(order);  // Extracted

  const subtotal = order.items.reduce((sum, i) => sum + i.price, 0);
  const tax = subtotal * 0.1;
  const total = subtotal + tax;

  db.save({ ...order, total });
}

function validateOrder(order: Order) {
  if (!order.items.length) throw new Error('Empty order');
  if (!order.customer) throw new Error('No customer');
}
```

**Rename**: Change a name to better reveal intent
**Introduce Parameter Object**: Group related parameters
**Move Method**: Move to the class that uses the data
**Extract Class**: Split responsibilities

## PRESENT CHANGE

```
Step 4: Change Made
===================

Iteration: {{N}}
Change: {{one-sentence description}}
Type: {{Extract Method / Rename / etc.}}
File: {{file:line}}
Before: {{brief summary}}
After: {{brief summary}}
```

Then ask: **[C] Continue to Step 5: Run Tests**

## FRONTMATTER UPDATE

Update the output document:
- Add `4` to `stepsCompleted` (or update if looping)
- Append the change description to the log

## NEXT STEP

After user confirms `[C]`, load `step-05-run-tests.md`.
