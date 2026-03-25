# Vue 3 & Svelte 5 Patterns

Production-quality patterns for Vue 3.5 Composition API and Svelte 5 Runes.

---

## Vue 3 Composition API

### Component with Props, Emits, and Reactivity

```vue
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';

interface Props {
  userId: string;
  initialName?: string;
}

const props = withDefaults(defineProps<Props>(), {
  initialName: '',
});

const emit = defineEmits<{
  (e: 'update', value: string): void;
  (e: 'submit'): void;
}>();

const name = ref(props.initialName);
const isLoading = ref(false);
const isValid = computed(() => name.value.length >= 3);

watch(name, (newValue) => { emit('update', newValue); });

onMounted(async () => {
  isLoading.value = true;
  // fetch data...
  isLoading.value = false;
});

const handleSubmit = () => { if (isValid.value) emit('submit'); };
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <input v-model="name" :disabled="isLoading" />
    <button type="submit" :disabled="!isValid">Submit</button>
  </form>
</template>
```

---

## Vue 3.5 "Tengu" (September 2024)

### Reactive Props Destructure (Stable)

```vue
<script setup lang="ts">
// Vue 3.5+: destructured props are reactive — withDefaults deprecated for this pattern
const { userId, count = 0, title = 'Default' } = defineProps<{
  userId: string;
  count?: number;
  title?: string;
}>();

// `count` and `title` are reactive with defaults
// No need for withDefaults() when using destructure
const doubled = computed(() => count * 2);
</script>
```

### useTemplateRef

```vue
<script setup lang="ts">
import { useTemplateRef, onMounted } from 'vue';

// Type-safe template ref (replaces string-based ref())
const inputRef = useTemplateRef<HTMLInputElement>('input');

onMounted(() => {
  inputRef.value?.focus();
});
</script>

<template>
  <input ref="input" />
</template>
```

### useId

```vue
<script setup lang="ts">
import { useId } from 'vue';

// SSR-stable unique IDs (like React useId)
const id = useId();
</script>

<template>
  <label :for="id">Email</label>
  <input :id="id" type="email" />
</template>
```

### Lazy Hydration (SSR)

```vue
<script setup>
import { defineAsyncComponent, hydrateOnVisible, hydrateOnIdle, hydrateOnInteraction } from 'vue';

// Hydrate when element becomes visible
const HeavyChart = defineAsyncComponent({
  loader: () => import('./HeavyChart.vue'),
  hydrate: hydrateOnVisible(),
});

// Hydrate during idle time
const Analytics = defineAsyncComponent({
  loader: () => import('./Analytics.vue'),
  hydrate: hydrateOnIdle(2000), // timeout: 2s
});

// Hydrate on user interaction
const CommentSection = defineAsyncComponent({
  loader: () => import('./CommentSection.vue'),
  hydrate: hydrateOnInteraction(['click', 'focus']),
});
</script>
```

### Other 3.5 Improvements

| Feature | Details |
|---------|---------|
| `onWatcherCleanup()` | Cleanup function for `watch`/`watchEffect` — replaces `onCleanup` parameter |
| Reactive `v-bind` memory | ~56% reduction in memory usage |
| Deferred Teleport | `<Teleport defer>` — renders target after current update cycle |

---

## Vue Composables (Custom Hooks)

```typescript
// composables/useAsync.ts
import { ref, type Ref } from 'vue';

export function useAsync<T>(asyncFn: () => Promise<T>): {
  data: Ref<T | null>; error: Ref<Error | null>;
  isLoading: Ref<boolean>; execute: () => Promise<void>;
} {
  const data = ref<T | null>(null) as Ref<T | null>;
  const error = ref<Error | null>(null);
  const isLoading = ref(false);

  const execute = async () => {
    isLoading.value = true;
    error.value = null;
    try { data.value = await asyncFn(); }
    catch (e) { error.value = e instanceof Error ? e : new Error(String(e)); }
    finally { isLoading.value = false; }
  };

  return { data, error, isLoading, execute };
}
```

---

## Pinia Store

```typescript
import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    isAuthenticated: false,
  }),
  getters: {
    userName: (state) => state.user?.name ?? 'Guest',
  },
  actions: {
    async login(email: string, password: string) {
      const user = await authApi.login(email, password);
      this.user = user;
      this.isAuthenticated = true;
    },
    logout() {
      this.user = null;
      this.isAuthenticated = false;
    },
  },
});
```

---

## Svelte 5 Runes (Stable — October 2024)

### Component with $state, $derived, $effect

```svelte
<script lang="ts">
  interface Props {
    userId: string;
    initialCount?: number;
  }

  let { userId, initialCount = 0 }: Props = $props();

  let count = $state(initialCount);
  let name = $state('');
  let doubled = $derived(count * 2);
  let isValid = $derived(name.length >= 3);

  $effect(() => {
    console.log(`Count changed to ${count}`);
    return () => { console.log('Cleanup'); };
  });

  function increment() { count++; }
  function handleSubmit() { if (isValid) { /* submit */ } }
</script>

<div>
  <p>Count: {count} (doubled: {doubled})</p>
  <button onclick={increment}>Increment</button>

  <form onsubmit={handleSubmit}>
    <input bind:value={name} />
    <button type="submit" disabled={!isValid}>Submit</button>
  </form>
</div>
```

### Svelte 5 Migration Guide

| Svelte 4 | Svelte 5 | Notes |
|-----------|----------|-------|
| `export let prop` | `let { prop } = $props()` | Props via Runes |
| `$:` reactive | `$derived()` / `$effect()` | Explicit reactivity |
| `<slot>` | `{@render children()}` with Snippets | Type-safe composition |
| `on:click={handler}` | `onclick={handler}` | Standard DOM event attributes |
| `createEventDispatcher()` | Callback props | Pass functions as props |
| Stores (`$store`) | `$state` + context | Runes replace most store patterns |

### $bindable (Two-Way Binding)

```svelte
<script lang="ts">
  // Child component with bindable prop
  let { value = $bindable('') }: { value: string } = $props();
</script>

<input bind:value={value} />

<!-- Parent usage -->
<script lang="ts">
  let searchQuery = $state('');
</script>
<SearchInput bind:value={searchQuery} />
```

### $state.raw (Large Data)

```svelte
<script lang="ts">
  // Skip deep reactivity for large, immutable datasets
  let items = $state.raw<Item[]>([]);

  async function loadItems() {
    const data = await fetchItems();
    items = data; // Replace entire array (no deep proxy overhead)
  }
</script>
```

### Snippets (Replacing Slots)

```svelte
<!-- Card.svelte -->
<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title: string;
    children: Snippet;
    footer?: Snippet;
  }

  let { title, children, footer }: Props = $props();
</script>

<div class="card">
  <header><h2>{title}</h2></header>
  <div class="card-body">{@render children()}</div>
  {#if footer}
    <footer>{@render footer()}</footer>
  {/if}
</div>

<!-- Usage -->
<Card title="My Card">
  <p>Card content</p>
  {#snippet footer()}
    <button>Action</button>
  {/snippet}
</Card>
```

### SvelteKit $app/state (2.12+)

```svelte
<script>
  import { page } from '$app/state';

  // Reactive access to page data (replaces $app/stores)
  // No need for $page store subscription
  const currentPath = $derived(page.url.pathname);
</script>

<nav>
  <a href="/" class:active={currentPath === '/'}>Home</a>
</nav>
```

---

## Vue Performance Hints

### v-memo

```vue
<template>
  <!-- Only re-render rows where selection changed -->
  <div v-for="item in items" :key="item.id" v-memo="[item.id === selectedId]">
    <ItemComponent :item="item" :selected="item.id === selectedId" />
  </div>
</template>
```

Use for: large lists (100+ rows), selection toggling, table row highlighting. `v-memo="[]"` equals `v-once`.

### markRaw

```ts
import { markRaw, ref } from 'vue';

// Skip reactivity for large third-party objects
const map = ref(markRaw(new MapLibreGL.Map({ /* ... */ })));
const chart = markRaw(new Chart(canvas, config));
const editor = markRaw(new Monaco.Editor(container));
```

Use for: map libraries, chart libraries, editor instances, immutable config objects.

**Source:** [Vue 3.5 Blog](https://blog.vuejs.org/posts/vue-3-5) · [Vue Composition API](https://vuejs.org/guide/extras/composition-api-faq) · [Svelte 5 Docs](https://svelte.dev/docs/svelte) · [Svelte 5 Migration Guide](https://svelte.dev/docs/svelte/v5-migration-guide) · [SvelteKit $app/state](https://svelte.dev/docs/kit/$app-state) · [Pinia Docs](https://pinia.vuejs.org/)
