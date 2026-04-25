# Form Interaction Patterns

Purpose: Choose the right validation strategy, error display, step flow, defaults, and submission feedback for forms.

## Contents

- Validation strategy
- Error display
- Multi-step forms
- Affordances and help
- Defaults
- Submission and unsaved changes
- Accessibility checklist

## Validation Strategy

| Strategy | Use it for | Trade-off |
|----------|------------|-----------|
| Real-time | format-specific fields such as email, URL, phone, password strength | immediate help, but can interrupt typing |
| On-blur | most text fields | low interruption, still timely |
| On-submit | cross-field or complex validation | all issues at once |
| Debounced | async checks such as username availability | slower, but avoids request spam |

Rules:

- use `aria-invalid` on invalid fields
- connect errors with `aria-describedby`
- announce blocking errors via `role="alert"` or `aria-live`

## Error Display

- keep field-level errors specific and actionable
- use an error summary when submit-time validation returns multiple errors
- provide recovery actions for non-field failures

Preferred pattern:

1. identify the field or action
2. explain what is wrong
3. tell the user how to fix it

## Multi-Step Forms

- show clear step progress
- allow safe back navigation
- persist step data
- avoid more than `7+` steps before reconsidering the flow design

## Affordances And Help

- use real labels; placeholder is not a label
- add helper text when the format is non-obvious
- use character counters only when limits matter
- use tooltips for optional help, not critical instructions

## Defaults

- prefer smart defaults when confidence is high
- make defaults easy to inspect and change
- never hide risky auto-filled values

## Submission And Unsaved Changes

- submit buttons need idle, loading, success, and error states
- disable duplicate submission while the request is active
- warn before discarding unsaved changes
- for destructive form actions, prefer confirm or undo patterns

## Accessibility Checklist

- labels and grouping are explicit
- errors are announced
- tab order is logical
- instructions remain visible
- the form works with keyboard only
