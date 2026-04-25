# Feature Flag Patterns

## Flag Types

| Type | Lifecycle | Example |
|------|-----------|---------|
| Release | Short (days-weeks) | New checkout flow |
| Experiment | Medium (weeks) | A/B test variant |
| Ops | Permanent | Kill switch, rate limit |
| Permission | Permanent | Premium features |

## LaunchDarkly Integration

```typescript
import { init } from 'launchdarkly-js-client-sdk';

const client = init('client-key', { key: userId });

const showNewCheckout = client.variation('new-checkout', false);
```

## Custom Feature Flag

```typescript
interface FeatureFlag {
  key: string;
  enabled: boolean;
  variants?: Record<string, unknown>;
  targeting?: { percentage: number; segments?: string[] };
}
```

## Cleanup Checklist
- [ ] Remove flag evaluation code
- [ ] Remove losing variant code
- [ ] Archive flag in management system
- [ ] Update documentation

---

## Basic Feature Flag Setup

```typescript
// lib/featureFlags.ts
interface FeatureFlag {
  name: string;
  variants: string[];
  allocation: number[];  // Percentage for each variant
  enabled: boolean;
}

const flags: Record<string, FeatureFlag> = {
  'new_checkout_flow': {
    name: 'new_checkout_flow',
    variants: ['control', 'treatment'],
    allocation: [50, 50],
    enabled: true
  }
};

export function getVariant(
  flagName: string,
  userId: string
): string {
  const flag = flags[flagName];
  if (!flag || !flag.enabled) {
    return 'control';
  }

  // Deterministic assignment based on user ID
  const hash = hashUserId(userId, flagName);
  const bucket = hash % 100;

  let cumulative = 0;
  for (let i = 0; i < flag.variants.length; i++) {
    cumulative += flag.allocation[i];
    if (bucket < cumulative) {
      return flag.variants[i];
    }
  }

  return 'control';
}

function hashUserId(userId: string, salt: string): number {
  const str = `${userId}:${salt}`;
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash & hash;
  }
  return Math.abs(hash);
}
```

## React Integration

```tsx
// components/ExperimentProvider.tsx
import { createContext, useContext, ReactNode } from 'react';
import { getVariant } from '@/lib/featureFlags';
import { useUser } from '@/hooks/useUser';
import { trackEvent } from '@/lib/analytics';

interface ExperimentContextValue {
  getExperimentVariant: (experimentName: string) => string;
  trackExposure: (experimentName: string) => void;
}

const ExperimentContext = createContext<ExperimentContextValue | null>(null);

export function ExperimentProvider({ children }: { children: ReactNode }) {
  const { user } = useUser();

  const getExperimentVariant = (experimentName: string): string => {
    return getVariant(experimentName, user?.id || 'anonymous');
  };

  const trackExposure = (experimentName: string): void => {
    const variant = getExperimentVariant(experimentName);
    trackEvent('experiment_exposure', {
      experiment_name: experimentName,
      variant: variant,
      user_id: user?.id
    });
  };

  return (
    <ExperimentContext.Provider value={{ getExperimentVariant, trackExposure }}>
      {children}
    </ExperimentContext.Provider>
  );
}

export function useExperiment(experimentName: string) {
  const context = useContext(ExperimentContext);
  if (!context) {
    throw new Error('useExperiment must be used within ExperimentProvider');
  }

  const variant = context.getExperimentVariant(experimentName);

  // Track exposure on first render
  useEffect(() => {
    context.trackExposure(experimentName);
  }, [experimentName]);

  return {
    variant,
    isControl: variant === 'control',
    isTreatment: variant === 'treatment'
  };
}
```

## Usage Example

```tsx
function CheckoutPage() {
  const { variant, isTreatment } = useExperiment('new_checkout_flow');

  return (
    <div>
      {isTreatment ? (
        <NewCheckoutFlow />
      ) : (
        <CurrentCheckoutFlow />
      )}
    </div>
  );
}
```
