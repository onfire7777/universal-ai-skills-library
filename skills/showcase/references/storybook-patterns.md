# Storybook Patterns & Templates

Comprehensive CSF 3.0 templates, Storybook 8.5+ features, and best practices.

---

## CSF 3.0 Templates

### Basic Component Story

```typescript
import type { Meta, StoryObj } from '@storybook/react';
import { within, userEvent, expect } from '@storybook/test';
import { ComponentName } from './ComponentName';

const meta = {
  title: 'Category/ComponentName',
  component: ComponentName,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component: 'Component description here',
      },
    },
  },
// ...
```

### Form Component Story

```typescript
import type { Meta, StoryObj } from '@storybook/react';
import { within, userEvent, expect } from '@storybook/test';
import { Input } from './Input';

const meta = {
  title: 'Forms/Input',
  component: Input,
  tags: ['autodocs'],
  argTypes: {
    type: { control: 'select', options: ['text', 'email', 'password', 'number'] },
    size: { control: 'select', options: ['sm', 'md', 'lg'] },
  },
} satisfies Meta<typeof Input>;

export default meta;
// ...
```

---

## Storybook 8.5+ Features

### Vitest Browser Mode Integration

```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';
import { storybookTest } from '@storybook/experimental-addon-test/vitest-plugin';

export default defineConfig({
  plugins: [storybookTest()],
  test: {
    browser: {
      enabled: true,
      provider: 'playwright',
      name: 'chromium',
    },
    setupFiles: ['.storybook/vitest.setup.ts'],
  },
});
// ...
```

### React Server Components (RSC) Stories

```typescript
// ServerComponent.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { ServerComponent } from './ServerComponent';

const meta = {
  title: 'RSC/ServerComponent',
  component: ServerComponent,
  parameters: {
    // Mark as server component for proper rendering
    nextjs: { appDirectory: true },
  },
  tags: ['autodocs'],
} satisfies Meta<typeof ServerComponent>;

export default meta;
// ...
```

### @storybook/test (Unified Testing)

```typescript
import { fn, expect, within, userEvent, waitFor } from '@storybook/test';
import type { Meta, StoryObj } from '@storybook/react';
import { Form } from './Form';

const meta = {
  component: Form,
  args: {
    // Type-safe mock functions
    onSubmit: fn(),
    onCancel: fn(),
  },
} satisfies Meta<typeof Form>;

export default meta;
type Story = StoryObj<typeof meta>;
// ...
```

### Portable Stories (Test Reuse)

```typescript
// Button.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';

const meta = {
  component: Button,
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: { variant: 'primary', children: 'Click me' },
};

// ...
```

### beforeEach / afterEach Lifecycle

```typescript
import type { Meta, StoryObj } from '@storybook/react';
import { within, userEvent } from '@storybook/test';
import { Modal } from './Modal';

const meta = {
  component: Modal,
  beforeEach: async () => {
    localStorage.clear();
  },
  afterEach: async () => {
    // Cleanup
  },
} satisfies Meta<typeof Modal>;

export default meta;
// ...
```

### Tags for Organization & Filtering

```typescript
const meta = {
  component: Button,
  tags: [
    'autodocs',      // Auto-generate docs
    'component',     // Category tag
    'visual-test',   // Include in visual regression
    '!dev',          // Exclude from dev sidebar
    '!test',         // Exclude from test runs
  ],
} satisfies Meta<typeof Button>;

// Filter in test runner:
// test-storybook --tags="component"
// test-storybook --tags="visual-test"
```

### Theme Testing (Dark Mode)

```typescript
import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';

const meta = {
  component: Button,
  parameters: {
    backgrounds: {
      default: 'light',
      values: [
        { name: 'light', value: '#ffffff' },
        { name: 'dark', value: '#1a1a1a' },
      ],
    },
  },
  decorators: [
// ...
```

---

## MDX 3 Documentation

### Component Documentation

```mdx
{/* Button.mdx */}
import { Meta, Stories, Primary, Controls, Story } from '@storybook/blocks';
import * as ButtonStories from './Button.stories';

<Meta of={ButtonStories} />

# Button

A button component used as the trigger for user actions.

## Usage

\`\`\`tsx
import { Button } from '@/components/Button';

// ...
```

### MDX with Custom Blocks

```mdx
{/* ComponentDoc.mdx */}
import { Meta, Canvas, Source } from '@storybook/blocks';
import * as Stories from './Component.stories';

<Meta of={Stories} />

# Component Name

## Interactive Demo

<Canvas of={Stories.Default} />

## Code Example

<Source of={Stories.Default} />
// ...
```

---

## Figma Integration

### Figma → Storybook Sync

```typescript
// .storybook/preview.ts
import { withDesign } from 'storybook-addon-designs';

export default {
  decorators: [withDesign],
  parameters: {
    design: {
      type: 'figma',
      url: 'https://www.figma.com/file/xxx',
    },
  },
};

// Button.stories.tsx
export const Primary: Story = {
// ...
```

### Design Tokens from Figma

```typescript
// tokens.ts - Generated from Figma Variables
export const tokens = {
  colors: {
    primary: { 50: '#eff6ff', 500: '#3b82f6', 900: '#1e3a8a' },
  },
  spacing: { 1: '4px', 2: '8px', 4: '16px' },
  radius: { sm: '4px', md: '8px', lg: '12px' },
};

// .storybook/preview.ts
import { tokens } from '../src/tokens';

export default {
  parameters: {
    backgrounds: {
// ...
```

---

## Audit Report Format

```markdown
## Showcase Audit Report: [Project Name]

### Coverage Summary

| Category | Total Components | With Stories | Coverage |
|----------|------------------|--------------|----------|
| Atoms | X | Y | Z% |
| Molecules | X | Y | Z% |
| Organisms | X | Y | Z% |
| Templates | X | Y | Z% |
| **Total** | **X** | **Y** | **Z%** |

### Story Quality Scores

| Component | Variants | A11y | Interactions | Docs | Grade |
...
```

---

## Forge Enhancement Workflow

### Forge Preview → Showcase Full Coverage

```
Forge (Preview Story)              Showcase (Full Story)
├─ Default state only              ├─ All variants
├─ Prototypes/ hierarchy           ├─ Components/ hierarchy
├─ tags: ['prototype']             ├─ tags: ['autodocs', 'component']
├─ No interactions                 ├─ Play functions
├─ No a11y config                  ├─ A11y rules configured
└─ TODO comments                   └─ MDX documentation
```

### Enhancement Checklist

```markdown
## Showcase Enhancement Checklist

### Story Location
- [ ] Move from `Prototypes/` to appropriate category
- [ ] Update title path in meta

### Variant Coverage
- [ ] Add size variants (sm, md, lg)
- [ ] Add color/theme variants
- [ ] Add state variants (default, hover, focus, active, disabled)
- [ ] Add content variants (empty, minimal, maximal)

### Interaction Tests
- [ ] Add play function for primary interaction
- [ ] Add keyboard navigation test
...
```

### Enhancement Template

```typescript
// BEFORE (Forge generated)
const meta = {
  component: ComponentName,
  title: 'Prototypes/ComponentName',
  tags: ['prototype'],
} satisfies Meta<typeof ComponentName>;

export const Preview: Story = {
  args: { /* default only */ },
};

// AFTER (Showcase enhanced)
const meta = {
  component: ComponentName,
  title: 'Components/ComponentName',
// ...
```
