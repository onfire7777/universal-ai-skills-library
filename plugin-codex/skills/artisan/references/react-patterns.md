# React Patterns & Server Components

Production-quality React patterns: components, hooks, error boundaries, React 19 hooks, React Compiler, RSC, Server Actions, and form handling.

---

## 1. Component Structure

### Compound Component Pattern

```tsx
interface CardProps {
  children: React.ReactNode;
  className?: string;
}

interface CardComponent extends React.FC<CardProps> {
  Header: typeof CardHeader;
  Body: typeof CardBody;
  Footer: typeof CardFooter;
}

const Card: CardComponent = ({ children, className }) => (
  <div className={cn("rounded-lg border", className)}>{children}</div>
);

const CardHeader: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="border-b p-4 font-semibold">{children}</div>
);

const CardBody: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="p-4">{children}</div>
);

const CardFooter: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="border-t p-4">{children}</div>
);

Card.Header = CardHeader;
Card.Body = CardBody;
Card.Footer = CardFooter;
```

---

## 2. Custom Hooks

```tsx
function useAsync<T>(asyncFn: () => Promise<T>, deps: unknown[] = []) {
  const [state, setState] = useState<{
    data: T | null;
    error: Error | null;
    isLoading: boolean;
  }>({ data: null, error: null, isLoading: true });

  useEffect(() => {
    let cancelled = false;
    setState(prev => ({ ...prev, isLoading: true }));
    asyncFn()
      .then(data => { if (!cancelled) setState({ data, error: null, isLoading: false }); })
      .catch(error => { if (!cancelled) setState({ data: null, error, isLoading: false }); });
    return () => { cancelled = true; };
  }, deps);

  return state;
}
```

---

## 3. Error Boundaries

```tsx
class ErrorBoundary extends React.Component<
  { fallback: React.ReactNode; children: React.ReactNode },
  { hasError: boolean }
> {
  state = { hasError: false };
  static getDerivedStateFromError() { return { hasError: true }; }
  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }
  render() {
    return this.state.hasError ? this.props.fallback : this.props.children;
  }
}
```

---

## 4. React 19 Hooks

### useActionState

```tsx
import { useActionState } from 'react';

function ContactForm() {
  const [state, formAction] = useActionState(
    async (prevState, formData: FormData) => {
      const name = formData.get('name') as string;
      const res = await submitContact({ name });
      if (!res.ok) return { error: 'Failed to submit' };
      return { success: true, message: `Thanks ${name}!` };
    },
    {}
  );

  return (
    <form action={formAction}>
      <input name="name" required />
      <SubmitButton />
      {state.error && <p role="alert">{state.error}</p>}
      {state.success && <p>{state.message}</p>}
    </form>
  );
}
```

### useFormStatus

```tsx
import { useFormStatus } from 'react-dom';

function SubmitButton() {
  const { pending } = useFormStatus();
  return (
    <button type="submit" disabled={pending}>
      {pending ? 'Submitting...' : 'Submit'}
    </button>
  );
}
```

### useOptimistic

```tsx
import { useOptimistic } from 'react';

function TodoList({ todos, addTodoAction }) {
  const [optimisticTodos, addOptimisticTodo] = useOptimistic(
    todos,
    (state, newTodo) => [...state, { ...newTodo, pending: true }]
  );

  return (
    <form action={async (formData) => {
      addOptimisticTodo({ text: formData.get('text'), id: crypto.randomUUID() });
      await addTodoAction(formData);
    }}>
      {optimisticTodos.map(todo => (
        <li key={todo.id} style={{ opacity: todo.pending ? 0.5 : 1 }}>{todo.text}</li>
      ))}
      <input name="text" />
      <button type="submit">Add</button>
    </form>
  );
}
```

### use()

Promise or Context in render. Can be called conditionally (exception to Rules of Hooks).

```tsx
import { use, Suspense } from 'react';

function UserProfile({ userPromise }: { userPromise: Promise<User> }) {
  const user = use(userPromise);
  return <h1>{user.name}</h1>;
}

function ThemeText({ showTheme }: { showTheme: boolean }) {
  if (showTheme) {
    const theme = use(ThemeContext); // OK inside conditional
    return <p>Current theme: {theme}</p>;
  }
  return <p>Theme hidden</p>;
}
```

---

## 5. React 19 Breaking Changes

| Change | Migration |
|--------|-----------|
| `forwardRef` deprecated | Use `ref` as a regular prop: `function Input({ ref, ...props })` |
| `Context.Provider` deprecated | Use `<MyContext>` directly as the provider |
| `ref` callbacks get cleanup | Return a cleanup function from ref callbacks |
| Document metadata hoisting | `<title>`, `<meta>`, `<link>` in components hoist to `<head>` automatically |

### React 19.2 (October 2025)

| Feature | Purpose |
|---------|---------|
| `<Activity>` | Offscreen rendering — hide/show without unmounting (replaces old `<Offscreen>`) |
| `<ViewTransition>` | Built-in View Transitions API integration for route/state transitions |
| `useEffectEvent` | Stable reference for event handlers used inside effects without adding to deps |

---

## 6. React Compiler (v1.0 Stable — October 2025)

Auto-memoization for components and hooks. No manual `useMemo`/`useCallback`/`React.memo` needed.

| Before | With React Compiler |
|--------|-------------------|
| Manual `useMemo` / `useCallback` | **Auto-memoized** — not needed |
| `React.memo` wrapper | Compiler decides |
| Dependency array management | Compiler tracks |

**Rules:**
- New code: skip `useMemo`/`useCallback` entirely
- Existing manual memoization: safe to keep (no harm)
- Requires Rules of React compliance (pure render, no side effects during render)
- Opt-out directive: `"use no memo"` at file/function level
- Next.js 16: `reactCompiler: true` in `next.config.ts`

---

## 7. Server Components (RSC)

### Composition: Push `'use client'` to Leaves

```tsx
// Bad: top-level 'use client' — all children become Client Components
'use client';
export function Layout() {
  const [isOpen, setIsOpen] = useState(false);
  return <ServerHeavyComponent />; // loses server execution
}

// Good: Composition pattern — pass Server Component as children
'use client';
export function InteractiveShell({ children }: { children: React.ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);
  return <div>{children}</div>;
}

// page.tsx (Server Component)
export default function Page() {
  return (
    <InteractiveShell>
      <ServerHeavyComponent /> {/* runs on server */}
    </InteractiveShell>
  );
}
```

### When to Use Client Components

| Condition | Server | Client |
|-----------|:---:|:---:|
| DB/API direct access | ✅ | — |
| Display-only (props) | ✅ | — |
| `onClick`/`onChange` handlers | — | ✅ |
| `useState`/`useEffect` | — | ✅ |
| Browser APIs (localStorage) | — | ✅ |

### Suspense Streaming

```tsx
// Bad: top-level await blocks entire shell
export default async function Home() {
  const analytics = await getAnalytics(); // slow
  return <div><Header /><AnalyticsWidget data={analytics} /></div>;
}

// Good: Suspense sends shell immediately, streams slow parts
export default function Home() {
  return (
    <div>
      <Header />
      <Suspense fallback={<AnalyticsSkeleton />}>
        <AnalyticsWidget />
      </Suspense>
    </div>
  );
}

async function AnalyticsWidget() {
  const data = await getAnalytics(); // blocks only within Suspense
  return <Dashboard data={data} />;
}
```

---

## 8. Server Actions

```tsx
// app/actions.ts
'use server';
import { z } from 'zod';
import { revalidatePath } from 'next/cache';

const TodoSchema = z.object({
  title: z.string().min(1, 'Title is required'),
});

export async function createTodo(prevState: any, formData: FormData) {
  const result = TodoSchema.safeParse({ title: formData.get('title') });
  if (!result.success) return { errors: result.error.flatten().fieldErrors };
  await db.todo.create({ data: result.data });
  revalidatePath('/todos');
  return { success: true };
}
```

### Server Actions vs Route Handlers

| Use case | Server Actions | Route Handlers |
|----------|:---:|:---:|
| Form mutations | ✅ | — |
| UI data changes | ✅ | — |
| Public API / Webhooks | — | ✅ |
| Large file uploads | — | ✅ |
| External service calls | — | ✅ |

### Cache & Revalidation

```tsx
const data = await fetch(url, { cache: 'force-cache' });         // static (build-time)
const data = await fetch(url, { cache: 'no-store' });            // dynamic (per-request)
const data = await fetch(url, { next: { revalidate: 3600 } });   // time-based (ISR)
const data = await fetch(url, { next: { tags: ['products'] } }); // tag-based
// After mutation:
revalidateTag('products');
```

---

## 9. Form Handling Selection Guide

| Criteria | React Hook Form | React 19 Native | Conform | TanStack Form |
|----------|:---:|:---:|:---:|:---:|
| Large/complex forms | ✅ | — | — | ✅ |
| Server Actions first | △ | ✅ | ✅ | — |
| Works without JS | — | ✅ | ✅ | — |
| Dynamic field arrays | ✅ | manual | ○ | ✅ |
| UI library integration | ✅ | — | — | ✅ |
| Bundle size | +8.6KB | 0KB | light | +12KB |

**Recommendations:** Complex forms → RHF + Zod / Simple + Server Actions → React 19 native / Remix/Next.js PE → Conform / Dynamic/nested → TanStack Form

---

## 10. Anti-Patterns

### Hooks

| # | Pattern | Problem | Fix |
|---|---------|---------|-----|
| 1 | Derived state via `useEffect` | Unnecessary re-renders | Compute directly in render |
| 2 | Raw `fetch` in `useEffect` | Memory leak, race conditions | TanStack Query |
| 3 | Conditional hook calls | Rules of Hooks violation | Always call + early return |
| 4 | Excessive `useMemo`/`useCallback` | Complexity without benefit | React Compiler / measure first |
| 5 | `use` prefix on non-hook functions | Misleading naming | Use regular function name |

### RSC

| # | Pattern | Problem | Fix |
|---|---------|---------|-----|
| 1 | Top-level `'use client'` | All children become Client Components | Push to leaves + Composition |
| 2 | Page-level `await` | Blocks entire shell | Suspense isolation |
| 3 | Serial data fetching | Server-side waterfall | `Promise.all` for parallel |
| 4 | Passing large objects across boundary | HTML payload bloat | Pass only needed fields |
| 5 | Missing `revalidate` after mutation | Stale UI | `revalidatePath`/`revalidateTag` |
| 6 | `await` in layouts | Delays entire app | Layouts sync, data in pages |

**Source:** [React 19 Blog](https://react.dev/blog/2024/12/05/react-19) · [React 19.2 Blog](https://react.dev/blog/2025/10/07/react-19-2) · [React Compiler v1.0](https://react.dev/blog/2025/10/07/react-compiler-rc) · [RSC Mental Models 2025](https://dev.to/eva_clari_289d85ecc68da48/the-complete-guide-to-react-server-components-mental-models-for-2025-390d) · [Next.js 16](https://nextjs.org/blog/next-16) · [React Hook Form vs React 19](https://blog.logrocket.com/react-hook-form-vs-react-19/)
