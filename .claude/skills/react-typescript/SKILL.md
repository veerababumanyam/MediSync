---
name: react-typescript
description: This skill should be used when the user asks to "create React components", "write TypeScript React code", "implement React hooks", "React patterns", "TypeScript best practices", "React state management", "React performance optimization", "custom hooks", or mentions React-specific concepts like useEffect, useState, context, memo, or Suspense.
---

# React + TypeScript Best Practices

React 19 with TypeScript provides type-safe, performant UI development. This skill covers idiomatic patterns, hooks, state management, and performance optimization for MediSync's frontend.

★ Insight ─────────────────────────────────────
MediSync's frontend stack:
1. **React 19** - Latest with Server Components support
2. **TypeScript 5.9** - Strict mode enabled
3. **CopilotKit** - Generative UI for AI interactions
4. **i18next** - English/Arabic with RTL support
5. **Apache ECharts** - Data visualization

Always use TypeScript strict mode and functional components.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Convention |
|--------|------------|
| **Components** | Functional components with arrow functions |
| **Typing** | Interface for props, type for unions/primitives |
| **State** | `useState` for local, Context/zustand for global |
| **Effects** | Prefer event handlers over effects |
| **Styling** | Tailwind CSS with logical properties for RTL |

## Component Patterns

### Basic Component Structure

```typescript
// components/UserCard.tsx
import { FC, memo } from 'react';

interface UserCardProps {
  user: User;
  onEdit?: (id: string) => void;
  className?: string;
}

export const UserCard: FC<UserCardProps> = memo(({ user, onEdit, className }) => {
  return (
    <div className={cn('p-4 rounded-lg', className)}>
      <h3>{user.name}</h3>
      {onEdit && (
        <button onClick={() => onEdit(user.id)}>
          Edit
        </button>
      )}
    </div>
  );
});

UserCard.displayName = 'UserCard';
```

### Component with Children

```typescript
interface CardProps {
  title: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
}

export const Card: FC<CardProps> = ({ title, children, footer }) => (
  <div className="card">
    <header>{title}</header>
    <main>{children}</main>
    {footer && <footer>{footer}</footer>}
  </div>
);
```

### Generic Components

```typescript
interface ListProps<T> {
  items: T[];
  renderItem: (item: T, index: number) => React.ReactNode;
  keyExtractor: (item: T) => string;
}

export function List<T>({ items, renderItem, keyExtractor }: ListProps<T>) {
  return (
    <ul>
      {items.map((item, index) => (
        <li key={keyExtractor(item)}>
          {renderItem(item, index)}
        </li>
      ))}
    </ul>
  );
}

// Usage
<List<User>
  items={users}
  renderItem={(user) => <span>{user.name}</span>}
  keyExtractor={(user) => user.id}
/>
```

## Hooks Patterns

### Custom Hook with Cleanup

```typescript
function useEventListener<K extends keyof WindowEventMap>(
  event: K,
  handler: (e: WindowEventMap[K]) => void,
  options?: AddEventListenerOptions
) {
  useEffect(() => {
    window.addEventListener(event, handler, options);
    return () => window.removeEventListener(event, handler, options);
  }, [event, handler, options]);
}
```

### Data Fetching Hook

```typescript
interface UseQueryResult<T> {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
  refetch: () => void;
}

function useQuery<T>(url: string): UseQueryResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await fetch(url);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const json = await response.json();
      setData(json);
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Unknown error'));
    } finally {
      setIsLoading(false);
    }
  }, [url]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, isLoading, error, refetch: fetchData };
}
```

### Debounced Value Hook

```typescript
function useDebouncedValue<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}
```

## State Management

### Context Pattern

```typescript
// context/UserContext.tsx
import { createContext, useContext, useReducer, FC, ReactNode } from 'react';

interface User {
  id: string;
  name: string;
  email: string;
}

interface UserState {
  user: User | null;
  isLoading: boolean;
}

type UserAction =
  | { type: 'SET_USER'; payload: User }
  | { type: 'CLEAR_USER' }
  | { type: 'SET_LOADING'; payload: boolean };

const UserContext = createContext<{
  state: UserState;
  dispatch: React.Dispatch<UserAction>;
} | null>(null);

function userReducer(state: UserState, action: UserAction): UserState {
  switch (action.type) {
    case 'SET_USER':
      return { ...state, user: action.payload, isLoading: false };
    case 'CLEAR_USER':
      return { ...state, user: null };
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    default:
      return state;
  }
}

export const UserProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(userReducer, {
    user: null,
    isLoading: true,
  });

  return (
    <UserContext.Provider value={{ state, dispatch }}>
      {children}
    </UserContext.Provider>
  );
};

export function useUser() {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within UserProvider');
  }
  return context;
}
```

## Performance Patterns

### Memoization

```typescript
// Memoize expensive computations
const sortedItems = useMemo(
  () => items.sort((a, b) => a.name.localeCompare(b.name)),
  [items]
);

// Memoize callbacks passed to children
const handleSubmit = useCallback(
  (data: FormData) => submitForm(data),
  [submitForm]
);
```

### List Virtualization (for large lists)

```typescript
import { useVirtualizer } from '@tanstack/react-virtual';

function VirtualList({ items }: { items: Item[] }) {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50,
  });

  return (
    <div ref={parentRef} className="h-screen overflow-auto">
      <div style={{ height: virtualizer.getTotalSize() }}>
        {virtualizer.getVirtualItems().map((virtualItem) => (
          <div
            key={virtualItem.key}
            style={{
              position: 'absolute',
              transform: `translateY(${virtualItem.start}px)`,
              height: virtualItem.size,
            }}
          >
            {items[virtualItem.index].name}
          </div>
        ))}
      </div>
    </div>
  );
}
```

### Lazy Loading

```typescript
const HeavyChart = lazy(() => import('./HeavyChart'));

function Dashboard() {
  return (
    <Suspense fallback={<ChartSkeleton />}>
      <HeavyChart data={chartData} />
    </Suspense>
  );
}
```

## i18n Integration

```typescript
// Use i18next with RTL support
import { useTranslation } from 'react-i18next';

function UserGreeting({ name }: { name: string }) {
  const { t, i18n } = useTranslation();
  const isRTL = i18n.dir() === 'rtl';

  return (
    <div className={isRTL ? 'text-right' : 'text-left'}>
      {t('greeting', { name })}
    </div>
  );
}
```

## Error Boundaries

```typescript
interface ErrorBoundaryProps {
  fallback: ReactNode;
  children: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.props.onError?.(error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children;
  }
}
```

## Form Handling

```typescript
import { useForm } from 'react-hook-form';

interface UserFormData {
  name: string;
  email: string;
  role: 'admin' | 'user';
}

function UserForm({ onSubmit }: { onSubmit: (data: UserFormData) => void }) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<UserFormData>();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input
        {...register('name', { required: 'Name is required' })}
        aria-invalid={!!errors.name}
      />
      {errors.name && <span role="alert">{errors.name.message}</span>}

      <input
        {...register('email', {
          required: 'Email is required',
          pattern: {
            value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
            message: 'Invalid email',
          },
        })}
      />

      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Saving...' : 'Save'}
      </button>
    </form>
  );
}
```

## Testing Patterns

```typescript
// Component test with Testing Library
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('UserCard', () => {
  it('calls onEdit when edit button is clicked', async () => {
    const user = userEvent.setup();
    const onEdit = vi.fn();
    const mockUser = { id: '1', name: 'Test User' };

    render(<UserCard user={mockUser} onEdit={onEdit} />);

    await user.click(screen.getByRole('button', { name: /edit/i }));

    expect(onEdit).toHaveBeenCalledWith('1');
  });
});
```

## Additional Resources

### Reference Files
- **`references/hooks.md`** - Comprehensive hook patterns
- **`references/state-management.md`** - State management strategies

### Example Files
- **`examples/DataGrid.tsx`** - Full data grid implementation
- **`examples/Modal.tsx`** - Accessible modal component
- **`examples/DataTable.tsx`** - Table with sorting/filtering
