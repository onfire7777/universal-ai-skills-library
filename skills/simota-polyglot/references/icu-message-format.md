# ICU Message Format & Translation Key Conventions

Patterns for ICU MessageFormat syntax and translation file organization.

---

## ICU Message Format

### Basic Plural

```json
// en.json
{
  "items_count": "{count, plural, =0 {No items} one {# item} other {# items}}"
}

// ja.json
{
  "items_count": "{count, plural, other {#個のアイテム}}"
}
```

```typescript
t('items_count', { count: 0 });  // → "No items"
t('items_count', { count: 1 });  // → "1 item"
t('items_count', { count: 5 });  // → "5 items"
```

### Select (Gender/Type)

```json
{
  "greeting": "{gender, select, male {He} female {She} other {They}} liked your post.",
  "notification_type": "{type, select, comment {commented on} like {liked} share {shared} other {interacted with}} your post"
}
```

```typescript
t('greeting', { gender: 'female' });           // → "She liked your post."
t('notification_type', { type: 'comment' });    // → "commented on your post"
```

### SelectOrdinal

```json
{
  "ranking": "You came in {place, selectordinal, one {#st} two {#nd} few {#rd} other {#th}} place!"
}
```

```typescript
t('ranking', { place: 1 });  // → "You came in 1st place!"
t('ranking', { place: 3 });  // → "You came in 3rd place!"
```

### Nested Messages

```json
{
  "notification": "{count, plural, =0 {No new notifications} one {{name} sent you a message} other {{name} and # others sent you messages}}"
}
```

### Date and Number in Messages

```json
{
  "last_login": "Last login: {date, date, medium}",
  "account_balance": "Your balance is {amount, number, currency}"
}
```

### Complex Example

```json
{
  "order_summary": "{itemCount, plural, =0 {Your cart is empty.} one {You have # item ({price, number, currency}) ready for checkout.} other {You have # items (total: {price, number, currency}) ready for checkout.}}"
}
```

---

## Translation Key Naming Conventions

### Flat vs Nested Structure

```json
// BAD: Flat (hard to maintain)
{
  "homeHeroTitle": "Welcome",
  "homeHeroDescription": "Description",
  "authLoginTitle": "Login"
}

// GOOD: Nested (organized by feature/page)
{
  "home": {
    "hero": {
      "title": "Welcome",
      "description": "Description"
    }
  },
  "auth": {
    "login": {
      "title": "Login",
      "button": "Sign In"
    }
  }
}
```

### Namespace Design

```
locales/
├── en/
│   ├── common.json      # Shared across app (buttons, labels)
│   ├── auth.json         # Login, signup, password reset
│   ├── dashboard.json    # Dashboard-specific
│   ├── settings.json     # Settings page
│   ├── errors.json       # Error messages
│   └── validation.json   # Form validation messages
├── ja/
│   ├── common.json
│   ├── auth.json
│   └── ...
```

### Common Namespace Examples

**common.json** — Shared UI elements:
```json
{
  "actions": {
    "save": "Save",
    "cancel": "Cancel",
    "delete": "Delete",
    "edit": "Edit",
    "submit": "Submit",
    "back": "Back",
    "next": "Next",
    "close": "Close"
  },
  "status": {
    "loading": "Loading...",
    "saving": "Saving...",
    "success": "Success!",
    "error": "An error occurred"
  },
  "pagination": {
    "previous": "Previous",
    "next": "Next",
    "page": "Page {{current}} of {{total}}"
  }
}
```

**errors.json** — Error messages with context:
```json
{
  "network": {
    "offline": "You appear to be offline. Please check your connection.",
    "timeout": "Request timed out. Please try again.",
    "server": "Server error. Please try again later."
  },
  "auth": {
    "invalid_credentials": "Invalid email or password.",
    "session_expired": "Your session has expired. Please log in again."
  },
  "validation": {
    "required": "This field is required.",
    "email_invalid": "Please enter a valid email address.",
    "password_weak": "Password must be at least 8 characters."
  }
}
```

### Key Naming Patterns

| Pattern | Example | Use Case |
|---------|---------|----------|
| `feature.element.action` | `auth.login.submit` | Button actions |
| `feature.element.state` | `order.status.pending` | Status text |
| `feature.message.type` | `cart.error.empty` | Error/success messages |
| `feature.label.field` | `profile.label.email` | Form labels |
| `feature.placeholder.field` | `search.placeholder.query` | Input placeholders |
| `feature.title.page` | `settings.title.page` | Page titles |

### Context-Aware Keys

```json
{
  "user_profile": {
    "page_title": "User Profile",
    "form": {
      "submit": "Update Profile"
    }
  }
}
```

### Translator Comments

```json
{
  "greeting": "Hello, {{name}}!",
  "_greeting_comment": "Appears at the top of the dashboard. 'name' is the user's first name.",

  "items_count": "{count, plural, one {# item} other {# items}}",
  "_items_count_comment": "Shopping cart item count. Keep it short for mobile.",

  "delete_confirm": "Are you sure you want to delete \"{{itemName}}\"?",
  "_delete_confirm_comment": "Confirmation dialog. itemName can be long (up to 50 chars)."
}
```
