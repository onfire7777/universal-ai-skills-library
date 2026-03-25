# Deprecated Library Catalog

## Date/Time Libraries

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `moment.js` | `date-fns`, `dayjs`, `Temporal API` | Moment is in maintenance mode. date-fns is tree-shakeable. Temporal API is the future standard. |
| `moment-timezone` | `Intl.DateTimeFormat`, `date-fns-tz` | Native Intl API handles most timezone needs. |

```typescript
// Before: moment
import moment from 'moment';
const formatted = moment().format('YYYY-MM-DD');

// After: date-fns (tree-shakeable)
import { format } from 'date-fns';
const formatted = format(new Date(), 'yyyy-MM-dd');

// After: Native Intl (no dependency)
const formatted = new Intl.DateTimeFormat('sv-SE').format(new Date());
```

## HTTP Libraries

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `request` | `node-fetch`, `undici`, native `fetch` | request is deprecated. Node 18+ has native fetch. |
| `axios` (consider) | native `fetch` | For simple cases, fetch is sufficient. axios still valid for interceptors/advanced features. |
| `superagent` | native `fetch` | fetch with AbortController covers most cases. |

```typescript
// Before: axios
import axios from 'axios';
const { data } = await axios.get('/api/users');

// After: Native fetch
const response = await fetch('/api/users');
const data = await response.json();
```

## Testing Libraries

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `enzyme` | `@testing-library/react` | Enzyme doesn't support React 18+. RTL encourages better testing patterns. |
| `sinon` (consider) | `jest.fn()`, `vitest.fn()` | Built-in mocking is often sufficient. |
| `karma` | `vitest`, `jest` | Modern test runners are faster and simpler. |

```typescript
// Before: Enzyme
import { shallow } from 'enzyme';
const wrapper = shallow(<MyComponent />);
expect(wrapper.find('.button').text()).toBe('Click');

// After: React Testing Library
import { render, screen } from '@testing-library/react';
render(<MyComponent />);
expect(screen.getByRole('button')).toHaveTextContent('Click');
```

## CSS/Styling Libraries

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `node-sass` | `sass` (dart-sass) | node-sass is deprecated. dart-sass is the primary implementation. |
| CSS-in-JS (runtime) | CSS Modules, Tailwind, vanilla-extract | Runtime CSS-in-JS has performance overhead. |
| `@emotion/core` | `@emotion/react` | Package renamed. |

## Utility Libraries

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `lodash` (full) | `lodash-es`, native methods | Import specific functions only. Many methods now native. |
| `underscore` | native ES6+ methods | Most utilities now built into JavaScript. |
| `uuid` (consider) | `crypto.randomUUID()` | Native in Node 19+, modern browsers. |
| `classnames` | `clsx` | clsx is smaller and faster. |

```typescript
// Before: lodash
import _ from 'lodash';
const result = _.uniq(array);

// After: Native Set
const result = [...new Set(array)];

// Before: uuid
import { v4 as uuidv4 } from 'uuid';
const id = uuidv4();

// After: Native crypto
const id = crypto.randomUUID();
```

## Build Tools

| Deprecated | Replacement | Migration Notes |
|------------|-------------|-----------------|
| `webpack` (consider) | `vite`, `esbuild`, `turbopack` | Vite offers faster DX. webpack still valid for complex setups. |
| `create-react-app` | `vite`, `next.js` | CRA is effectively deprecated. |
| `babel` (consider) | `swc`, `esbuild` | SWC/esbuild are faster. Babel still needed for some transforms. |
| `tslint` | `eslint` + `@typescript-eslint` | TSLint is officially deprecated. |
