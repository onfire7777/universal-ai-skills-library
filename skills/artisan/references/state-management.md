# State Management Decision Guide

State classification, tool selection framework, RSC-era patterns, and anti-patterns.

---

## 1. State Classification

| Type | Description | Recommended Tool |
|------|-------------|-----------------|
| **Remote State** | API/DB data; needs fetch, cache, sync | TanStack Query v5 / SWR |
| **URL State** | UI state reflected in URL params | nuqs v2 / searchParams |
| **Local State** | Scoped to a single component | `useState` / `useReducer` |
| **Shared State** | Shared across distant components | Context (few) → Zustand (many) |

---

## 2. Decision Flowchart

```
Need state?
├── Server data? → TanStack Query / SWR
├── Should reflect in URL? → nuqs / useSearchParams
├── Single component? → useState / useReducer
└── Shared across components?
    ├── 2-3 levels of parent-child? → Props drilling (OK)
    ├── 1-2 concerns? → React Context
    └── 3+ concerns or high-frequency updates? → Zustand
```

---

## 3. Library Selection

| Library | Good fit | Poor fit |
|---------|----------|----------|
| **Zustand** | Default for shared client state. Lightweight, selective re-renders | Server data management |
| **TanStack Query v5** | Standard for server state. Auto cache/dedup/invalidation | Pure client state |
| **Jotai** | Fine-grained reactivity, many derived/computed values | Simple global state |
| **Redux Toolkit** | Large enterprise, strict architecture, DevTools required | Small-to-mid apps |
| **XState** | Complex state machines (Figma-level) | Simple on/off state |
| **React Context** | Theme, locale, low-frequency updates | High-frequency state |

---

## 4. TanStack Query v5

### Migration from v4

| v4 | v5 | Notes |
|----|-----|-------|
| `onSuccess`/`onError`/`onSettled` in `useQuery` | **Removed** | Use `useEffect` or handle in `queryFn` |
| `cacheTime` | `gcTime` | Renamed for clarity |
| `isLoading` (first load only) | `isPending` | `isLoading` = `isPending && isFetching` |
| `useQuery({ queryKey, queryFn })` | Same | Object syntax only (no positional args) |

### Patterns

```tsx
// Custom hook encapsulating query
function useProducts(category?: string) {
  return useQuery({
    queryKey: ['products', { category }],
    queryFn: () => fetchProducts(category),
    staleTime: 5 * 60 * 1000,
  });
}

// Suspense-first query
function ProductList({ category }: { category: string }) {
  const { data } = useSuspenseQuery({
    queryKey: ['products', category],
    queryFn: () => fetchProducts(category),
  });
  return <ul>{data.map(p => <li key={p.id}>{p.name}</li>)}</ul>;
}

// Optimistic update mutation
function useUpdateProduct() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateProduct,
    onMutate: async (updated) => {
      await queryClient.cancelQueries({ queryKey: ['products'] });
      const prev = queryClient.getQueryData(['products']);
      queryClient.setQueryData(['products'], (old) =>
        old.map(p => p.id === updated.id ? { ...p, ...updated } : p)
      );
      return { prev };
    },
    onError: (_err, _vars, ctx) => {
      queryClient.setQueryData(['products'], ctx?.prev);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });
}

// RSC hydration support
import { HydrationBoundary, dehydrate } from '@tanstack/react-query';

// Server Component (page.tsx)
export default async function Page() {
  const queryClient = new QueryClient();
  await queryClient.prefetchQuery({
    queryKey: ['products'],
    queryFn: fetchProducts,
  });
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <ProductList />
    </HydrationBoundary>
  );
}
```

---

## 5. Zustand Patterns

```tsx
import { create } from 'zustand';

// Split stores by concern
const useUIStore = create<UIState>()((set) => ({
  isSidebarOpen: false,
  theme: 'light' as const,
  toggleSidebar: () => set((s) => ({ isSidebarOpen: !s.isSidebarOpen })),
  setTheme: (theme) => set({ theme }),
}));

// Subscribe only to needed values (prevents unnecessary re-renders)
function Sidebar() {
  const isOpen = useUIStore((s) => s.isSidebarOpen);
  // Does not re-render when theme changes
}
```

---

## 6. URL State with nuqs v2

```tsx
import { useQueryState, parseAsInteger } from 'nuqs';

// Client component
function ProductFilter() {
  const [category, setCategory] = useQueryState('category');
  const [page, setPage] = useQueryState('page', parseAsInteger.withDefault(1));
  // URL synced automatically: ?category=shoes&page=2
}

// Server-side parsing (Next.js App Router)
import { createSearchParamsCache, parseAsString } from 'nuqs/server';

const searchParamsCache = createSearchParamsCache({
  category: parseAsString.withDefault('all'),
  page: parseAsInteger.withDefault(1),
});

export default async function Page({ searchParams }: { searchParams: Promise<Record<string, string>> }) {
  const { category, page } = await searchParamsCache.parse(searchParams);
  const products = await fetchProducts({ category, page });
  return <ProductList products={products} />;
}
```

---

## 7. Anti-Patterns

| # | Anti-pattern | Problem | Fix |
|---|-------------|---------|-----|
| 1 | Server data in Zustand/Redux | Manual cache/sync/dedup reimplementation | Delegate to TanStack Query |
| 2 | Monolithic Context | Unrelated updates re-render all consumers | Split Context by concern |
| 3 | Providers Hell (5+ nested) | Poor readability, hard to debug | Consolidate with Zustand or compose utility |
| 4 | Over-globalizing local state | Modal open/close in global store | Component-level `useState` |
| 5 | `useState` + `useEffect` for URL sync | Sync bugs | nuqs / useSearchParams |
| 6 | Client-only data in RSC era | Ignores server-side benefits | Read in RSC, client stores for interactive parts only |
| 7 | `onSuccess` callbacks in v5 queries | Removed API — silent failures | Move logic to `queryFn` or use `useEffect` on `data` |

**Source:** [TanStack Query v5 Migration](https://tanstack.com/query/latest/docs/framework/react/guides/migrating-to-v5) · [nuqs v2 Docs](https://nuqs.47ng.com/) · [React State Management 2025](https://www.developerway.com/posts/react-state-management-2025) · [Zustand Best Practices](https://zustand.docs.pmnd.rs/)
