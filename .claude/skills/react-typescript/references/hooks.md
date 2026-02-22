# React Hooks Comprehensive Guide

## Built-in Hooks Reference

### State Management

```typescript
// useState - local component state
const [count, setCount] = useState(0);
const [user, setUser] = useState<User | null>(null);

// useReducer - complex state logic
type State = { count: number };
type Action = { type: 'increment' } | { type: 'decrement' };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'increment': return { count: state.count + 1 };
    case 'decrement': return { count: state.count - 1 };
  }
}

const [state, dispatch] = useReducer(reducer, { count: 0 });
```

### Side Effects

```typescript
// useEffect - side effects with cleanup
useEffect(() => {
  const controller = new AbortController();

  fetchData(controller.signal);

  return () => controller.abort();
}, [dependency]);

// useLayoutEffect - DOM measurements before paint
useLayoutEffect(() => {
  const rect = ref.current.getBoundingClientRect();
  // Sync DOM updates here
}, []);

// useInsertionEffect - CSS-in-JS injection (React 18+)
useInsertionEffect(() => {
  insertStyles(css);
}, [css]);
```

### Performance

```typescript
// useMemo - memoize expensive computations
const sortedItems = useMemo(
  () => items.filter(filterFn).sort(sortFn),
  [items, filterFn, sortFn]
);

// useCallback - memoize callbacks
const handleSubmit = useCallback(
  (data: FormData) => {
    submit(data);
  },
  [submit]
);

// memo - prevent unnecessary re-renders
const MemoizedComponent = memo(ExpensiveComponent, (prev, next) => {
  return prev.id === next.id;
});
```

### Refs

```typescript
// useRef - mutable value that persists
const inputRef = useRef<HTMLInputElement>(null);
const previousValue = useRef<T>(value);

// useImperativeHandle - expose methods to parent
useImperativeHandle(ref, () => ({
  focus: () => inputRef.current?.focus(),
  reset: () => inputRef.current?.setValue(''),
}));
```

### Context

```typescript
// useContext - consume context
const ThemeContext = createContext<Theme>('light');

function Component() {
  const theme = useContext(ThemeContext);
  return <div className={theme}>...</div>;
}
```

## Custom Hook Patterns

### useAsync - Async State Management

```typescript
interface AsyncState<T> {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
}

function useAsync<T>(
  asyncFn: () => Promise<T>,
  deps: React.DependencyList = []
): AsyncState<T> & { execute: () => Promise<void> } {
  const [state, setState] = useState<AsyncState<T>>({
    data: null,
    isLoading: false,
    error: null,
  });

  const execute = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }));
    try {
      const data = await asyncFn();
      setState({ data, isLoading: false, error: null });
    } catch (error) {
      setState({ data: null, isLoading: false, error: error as Error });
    }
  }, deps);

  useEffect(() => {
    execute();
  }, [execute]);

  return { ...state, execute };
}
```

### useLocalStorage - Persisted State

```typescript
function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T | ((prev: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch {
      return initialValue;
    }
  });

  const setValue = useCallback((value: T | ((prev: T) => T)) => {
    setStoredValue(prev => {
      const newValue = value instanceof Function ? value(prev) : value;
      localStorage.setItem(key, JSON.stringify(newValue));
      return newValue;
    });
  }, [key]);

  return [storedValue, setValue];
}
```

### useMediaQuery - Responsive Design

```typescript
function useMediaQuery(query: string): boolean {
  const [matches, setMatches] = useState(
    () => window.matchMedia(query).matches
  );

  useEffect(() => {
    const mediaQuery = window.matchMedia(query);
    const handler = (e: MediaQueryListEvent) => setMatches(e.matches);

    mediaQuery.addEventListener('change', handler);
    return () => mediaQuery.removeEventListener('change', handler);
  }, [query]);

  return matches;
}

// Usage
const isMobile = useMediaQuery('(max-width: 768px)');
const prefersDark = useMediaQuery('(prefers-color-scheme: dark)');
```

### useIntersectionObserver - Element Visibility

```typescript
function useIntersectionObserver(
  ref: RefObject<Element>,
  options?: IntersectionObserverInit
): IntersectionObserverEntry | null {
  const [entry, setEntry] = useState<IntersectionObserverEntry | null>(null);

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      ([entry]) => setEntry(entry),
      { threshold: 0.1, ...options }
    );

    observer.observe(element);
    return () => observer.disconnect();
  }, [ref, options]);

  return entry;
}

// Usage for lazy loading
const ref = useRef<HTMLDivElement>(null);
const entry = useIntersectionObserver(ref);
const isVisible = entry?.isIntersecting ?? false;
```

### useDebounce and useThrottle

```typescript
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}

function useThrottle<T>(value: T, interval: number): T {
  const [throttledValue, setThrottledValue] = useState(value);
  const lastUpdated = useRef<number>(Date.now());

  useEffect(() => {
    const now = Date.now();
    if (now - lastUpdated.current >= interval) {
      lastUpdated.current = now;
      setThrottledValue(value);
    } else {
      const timer = setTimeout(() => {
        lastUpdated.current = Date.now();
        setThrottledValue(value);
      }, interval - (now - lastUpdated.current));

      return () => clearTimeout(timer);
    }
  }, [value, interval]);

  return throttledValue;
}
```

### useClickOutside - Detect Outside Clicks

```typescript
function useClickOutside<T extends HTMLElement>(
  handler: () => void
): RefObject<T> {
  const ref = useRef<T>(null);

  useEffect(() => {
    const listener = (event: MouseEvent | TouchEvent) => {
      if (!ref.current || ref.current.contains(event.target as Node)) {
        return;
      }
      handler();
    };

    document.addEventListener('mousedown', listener);
    document.addEventListener('touchstart', listener);

    return () => {
      document.removeEventListener('mousedown', listener);
      document.removeEventListener('touchstart', listener);
    };
  }, [handler]);

  return ref;
}

// Usage
const ref = useClickOutside<HTMLDivElement>(() => setIsOpen(false));
return <div ref={ref}>...</div>;
```

### useKeyboardShortcut

```typescript
type KeyHandler = (event: KeyboardEvent) => void;

function useKeyboardShortcut(
  key: string,
  callback: KeyHandler,
  modifiers: { ctrl?: boolean; shift?: boolean; alt?: boolean } = {}
): void {
  useEffect(() => {
    const handler = (event: KeyboardEvent) => {
      if (
        event.key.toLowerCase() === key.toLowerCase() &&
        (!modifiers.ctrl || event.ctrlKey) &&
        (!modifiers.shift || event.shiftKey) &&
        (!modifiers.alt || event.altKey)
      ) {
        callback(event);
      }
    };

    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [key, callback, modifiers]);
}

// Usage
useKeyboardShortcut('s', handleSave, { ctrl: true });
useKeyboardShortcut('Escape', closeModal);
```

### usePrevious - Track Previous Value

```typescript
function usePrevious<T>(value: T): T | undefined {
  const ref = useRef<T>();

  useEffect(() => {
    ref.current = value;
  }, [value]);

  return ref.current;
}

// Usage
const [count, setCount] = useState(0);
const previousCount = usePrevious(count);
```

### useToggle

```typescript
function useToggle(initialValue = false): [boolean, () => void] {
  const [value, setValue] = useState(initialValue);
  const toggle = useCallback(() => setValue(v => !v), []);
  return [value, toggle];
}

// Usage
const [isOpen, toggleOpen] = useToggle();
```

## Rules of Hooks

1. **Only call hooks at the top level** - Not inside loops, conditions, or nested functions
2. **Only call hooks from React functions** - Components or custom hooks
3. **Use ESLint plugin** - `eslint-plugin-react-hooks` enforces rules

```typescript
// WRONG
function Component({ items }) {
  items.forEach(item => {
    useEffect(() => {}, [item]); // Hook in loop
  });
}

// CORRECT
function Component({ items }) {
  useEffect(() => {
    items.forEach(item => {
      // Process items inside effect
    });
  }, [items]);
}
```
