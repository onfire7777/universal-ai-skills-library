# Bolt React Performance Patterns

React-specific optimization patterns and code examples.

---

## Detecting Re-renders

```tsx
// Development-only re-render tracker
function useWhyDidYouUpdate(name: string, props: Record<string, unknown>) {
  const previousProps = useRef<Record<string, unknown>>();

  useEffect(() => {
    if (previousProps.current) {
      const allKeys = Object.keys({ ...previousProps.current, ...props });
      const changesObj: Record<string, { from: unknown; to: unknown }> = {};

      allKeys.forEach(key => {
        if (previousProps.current![key] !== props[key]) {
          changesObj[key] = {
            from: previousProps.current![key],
            to: props[key]
          };
        }
      });

      if (Object.keys(changesObj).length) {
        console.log('[why-did-you-update]', name, changesObj);
      }
    }
    previousProps.current = props;
  });
}
```

---

## React.memo Patterns

```tsx
// ❌ Bad: Inline object causes re-render every time
<UserCard user={{ name, email }} />

// ✅ Good: Memoized object
const user = useMemo(() => ({ name, email }), [name, email]);
<UserCard user={user} />

// ✅ Good: Custom comparison for complex props
const UserCard = memo(
  ({ user, onSelect }: Props) => { /* ... */ },
  (prevProps, nextProps) => {
    return prevProps.user.id === nextProps.user.id &&
           prevProps.user.updatedAt === nextProps.user.updatedAt;
  }
);
```

---

## useMemo vs useCallback

```tsx
// useMemo: Cache computed values
const sortedItems = useMemo(() => {
  return items.slice().sort((a, b) => a.name.localeCompare(b.name));
}, [items]);

// useMemo: Cache expensive calculations
const statistics = useMemo(() => {
  return calculateStatistics(data); // O(n) operation
}, [data]);

// useCallback: Cache functions passed to children
const handleSubmit = useCallback((values: FormValues) => {
  submitForm(values);
}, [submitForm]);

// useCallback: Cache event handlers for memoized children
const handleItemClick = useCallback((id: string) => {
  setSelectedId(id);
}, []); // Empty deps if setSelectedId is stable
```

---

## Context Optimization

```tsx
// ❌ Bad: Single context causes all consumers to re-render
const AppContext = createContext<{ user: User; theme: Theme; settings: Settings }>();

// ✅ Good: Split contexts by update frequency
const UserContext = createContext<User | null>(null);
const ThemeContext = createContext<Theme>('light');
const SettingsContext = createContext<Settings>(defaultSettings);

// ✅ Good: Memoize context value
function UserProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const value = useMemo(() => ({ user, setUser }), [user]);

  return (
    <UserContext.Provider value={value}>
      {children}
    </UserContext.Provider>
  );
}

// ✅ Good: Separate state and dispatch contexts
const StateContext = createContext<State>(initialState);
const DispatchContext = createContext<Dispatch<Action>>(() => {});

function Provider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(reducer, initialState);

  return (
    <StateContext.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </StateContext.Provider>
  );
}
```

---

## Lazy Loading Components

```tsx
// Route-based code splitting
const Dashboard = lazy(() => import('./pages/Dashboard'));
const Settings = lazy(() => import('./pages/Settings'));

function App() {
  return (
    <Suspense fallback={<PageSkeleton />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
}

// Named export lazy loading
const UserProfile = lazy(() =>
  import('./components/UserProfile').then(module => ({
    default: module.UserProfile
  }))
);

// Preload on hover/focus
const SettingsPage = lazy(() => import('./pages/Settings'));

function NavLink() {
  const preload = () => import('./pages/Settings');

  return (
    <Link
      to="/settings"
      onMouseEnter={preload}
      onFocus={preload}
    >
      Settings
    </Link>
  );
}
```

---

## List Virtualization

```tsx
// Using @tanstack/react-virtual
import { useVirtualizer } from '@tanstack/react-virtual';

function VirtualList({ items }: { items: Item[] }) {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50, // Estimated row height
    overscan: 5, // Render 5 extra items above/below viewport
  });

  return (
    <div ref={parentRef} style={{ height: '400px', overflow: 'auto' }}>
      <div style={{ height: `${virtualizer.getTotalSize()}px`, position: 'relative' }}>
        {virtualizer.getVirtualItems().map(virtualRow => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${virtualRow.size}px`,
              transform: `translateY(${virtualRow.start}px)`,
            }}
          >
            <ItemRow item={items[virtualRow.index]} />
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## Debounce and Throttle

```tsx
// Debounced search input
function SearchInput({ onSearch }: { onSearch: (query: string) => void }) {
  const [value, setValue] = useState('');

  const debouncedSearch = useMemo(
    () => debounce((query: string) => onSearch(query), 300),
    [onSearch]
  );

  useEffect(() => {
    return () => debouncedSearch.cancel();
  }, [debouncedSearch]);

  return (
    <input
      value={value}
      onChange={e => {
        setValue(e.target.value);
        debouncedSearch(e.target.value);
      }}
    />
  );
}

// Throttled scroll handler
function useThrottledScroll(callback: () => void, delay: number) {
  useEffect(() => {
    const throttled = throttle(callback, delay);
    window.addEventListener('scroll', throttled);
    return () => {
      window.removeEventListener('scroll', throttled);
      throttled.cancel();
    };
  }, [callback, delay]);
}
```
