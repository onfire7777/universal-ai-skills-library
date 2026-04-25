# Microinteraction Patterns

Purpose: Choose the correct feedback state, loading pattern, notification, and destructive-action safeguard for a small UI interaction.

## Contents

- Button feedback
- Validation feedback
- Loading states
- Notifications
- Destructive actions
- Coding standards

## Button Feedback

Standard state progression:

`idle -> hover -> pressed -> loading -> success/error`

Use it for any async action that changes visible state.

## Validation Feedback

| Pattern | Use it for |
|---------|------------|
| Real-time | format validation |
| On-blur | most text entry |
| Submit-time | cross-field validation |

Always reflect invalid state with semantic attributes.

## Loading States

| Pattern | Use it when |
|---------|-------------|
| Skeleton | content structure is known |
| Spinner | short action-level waits |
| Progressive loading | large lists or staged content |
| Optimistic update | the failure rate is very low and rollback is natural |

## Notifications

| Type | Suggested behavior |
|------|--------------------|
| Success toast | auto-dismiss around `3s` |
| Error toast | longer visibility or manual dismiss |
| Undo toast | show for around `5s` with an action |

Use `aria-live="polite"` for non-urgent success or undo feedback.

## Destructive Actions

| Pattern | Use it when |
|---------|-------------|
| Confirmation dialog | the action is permanent or high-risk |
| Soft delete + undo | recovery is possible and safer for flow continuity |

If the action cannot be undone, state that explicitly.

## Coding Standards

- keep feedback immediate and stateful
- do not leave async actions visually silent
- avoid optimistic UI on high-risk operations
- expose recovery whenever the action is recoverable
