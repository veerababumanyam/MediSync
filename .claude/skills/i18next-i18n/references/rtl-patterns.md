# RTL Layout Patterns for React/i18next

## CSS Logical Properties

### Tailwind RTL Utilities

```tsx
// Instead of directional properties, use logical properties
// These automatically flip in RTL mode

// Margins
<div className="ms-4 me-8">  {/* ms = margin-start, me = margin-end */}

// Padding
<div className="ps-4 pe-8">  {/* ps = padding-start, pe = padding-end */}

// Text alignment
<div className="text-start">  {/* Use instead of text-left */}
<div className="text-end">    {/* Use instead of text-right */}

// Borders
<div className="border-s border-e border-s-2 border-e-4">
<div className="rounded-s-lg rounded-e-lg rounded-ss-lg rounded-se-lg">

// Positioning
<div className="start-0 end-0">  {/* Instead of left-0 right-0 */}
<div className="inset-start-4 inset-end-4">

// Floats
<div className="float-start float-end">  {/* Instead of float-left float-right */}
```

### Standard CSS Logical Properties

```css
/* Instead of margin-left/right */
.element {
  margin-inline-start: 16px;  /* Left in LTR, Right in RTL */
  margin-inline-end: 16px;    /* Right in LTR, Left in RTL */
  padding-inline-start: 8px;
  padding-inline-end: 8px;
}

/* Border radius */
.element {
  border-start-start-radius: 8px;  /* top-left in LTR, top-right in RTL */
  border-start-end-radius: 8px;    /* top-right in LTR, top-left in RTL */
  border-end-start-radius: 8px;
  border-end-end-radius: 8px;
}

/* Position */
.element {
  inset-inline-start: 0;
  inset-inline-end: 0;
}

/* Text alignment */
.element {
  text-align: start;  /* Left in LTR, Right in RTL */
  text-align: end;    /* Right in LTR, Left in RTL */
}
```

## Component Patterns

### RTL-Aware Layout Component

```tsx
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';

interface DirectionalProps {
  children: React.ReactNode;
  className?: string;
}

export function DirectionalLayout({ children, className }: DirectionalProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <div
      dir={isRTL ? 'rtl' : 'ltr'}
      className={cn(isRTL && 'rtl', className)}
    >
      {children}
    </div>
  );
}

// Usage
<DirectionalLayout>
  <YourComponent />
</DirectionalLayout>
```

### Directional Icon Component

```tsx
import { useTranslation } from 'react-i18next';
import { cn } from '@/lib/utils';

// Icons that should flip in RTL
const MIRROR_ICONS = new Set([
  'arrow-left',
  'arrow-right',
  'chevron-left',
  'chevron-right',
  'arrow-back',
  'arrow-forward',
  'reply',
  'forward',
  'undo',
  'redo',
  'indent',
  'outdent',
]);

// Icons that should NOT flip
const NO_MIRROR_ICONS = new Set([
  'language',
  'settings',
  'help',
  'search',
  'refresh',
  'clock',
  'calendar',
  'phone',
  'email',
  'camera',
]);

interface DirectionalIconProps {
  name: string;
  className?: string;
}

export function DirectionalIcon({ name, className }: DirectionalIconProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  const shouldMirror = isRTL && MIRROR_ICONS.has(name) && !NO_MIRROR_ICONS.has(name);

  return (
    <Icon
      name={name}
      className={cn(
        className,
        shouldMirror && 'scale-x-[-1]'
      )}
    />
  );
}
```

### RTL-Aware Flex Components

```tsx
import { cn } from '@/lib/utils';
import { useRTL } from '@/hooks/useRTL';

interface RowProps {
  children: React.ReactNode;
  className?: string;
  reverse?: boolean;
}

export function Row({ children, className, reverse }: RowProps) {
  const isRTL = useRTL();

  // Automatically reverse for RTL if needed
  const shouldReverse = reverse !== isRTL;

  return (
    <div
      className={cn(
        'flex',
        shouldReverse ? 'flex-row-reverse' : 'flex-row',
        className
      )}
    >
      {children}
    </div>
  );
}

// Usage
<Row>
  <Icon name="arrow-left" />
  <Text>Back</Text>
</Row>
```

## Form Patterns

### RTL Form Fields

```tsx
import { useTranslation } from 'react-i18next';

export function RTLInput({ label, ...props }: InputProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <div className="flex flex-col gap-2">
      <label
        className={cn(
          'text-sm font-medium',
          isRTL ? 'text-end' : 'text-start'
        )}
      >
        {label}
      </label>
      <input
        {...props}
        dir={isRTL ? 'rtl' : 'ltr'}
        className={cn(
          'w-full px-4 py-2 border rounded-lg',
          isRTL ? 'text-end' : 'text-start'
        )}
      />
    </div>
  );
}
```

### Number Input with RTL

```tsx
export function RTLNumberInput({ value, onChange, ...props }: NumberInputProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  // Numbers should always be LTR even in RTL context
  return (
    <input
      type="number"
      value={value}
      onChange={onChange}
      dir="ltr"
      className={cn(
        'w-full px-4 py-2 border rounded-lg',
        isRTL ? 'text-end' : 'text-start'
      )}
      {...props}
    />
  );
}
```

## Navigation Patterns

### RTL Breadcrumbs

```tsx
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { DirectionalIcon } from './DirectionalIcon';

export function Breadcrumbs({ items }: BreadcrumbsProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <nav className="flex items-center gap-2">
      {items.map((item, index) => (
        <Fragment key={item.href}>
          {index > 0 && (
            <DirectionalIcon
              name={isRTL ? 'chevron-left' : 'chevron-right'}
              className="w-4 h-4 text-gray-400"
            />
          )}
          <Link href={item.href}>{item.label}</Link>
        </Fragment>
      ))}
    </nav>
  );
}
```

### RTL Sidebar Navigation

```tsx
export function SidebarNavigation({ items }: SidebarProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <aside
      className={cn(
        'fixed top-0 bottom-0 w-64 bg-gray-100',
        isRTL ? 'right-0' : 'left-0',
        isRTL ? 'border-l' : 'border-r'
      )}
    >
      <nav className="p-4">
        {items.map((item) => (
          <NavItem
            key={item.href}
            {...item}
            className={cn(
              'flex items-center gap-3 px-4 py-2',
              isRTL ? 'flex-row-reverse' : 'flex-row'
            )}
          />
        ))}
      </nav>
    </aside>
  );
}
```

## Table Patterns

### RTL Table

```tsx
export function RTLTable({ columns, data }: TableProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <div className="overflow-x-auto">
      <table className="w-full" dir={isRTL ? 'rtl' : 'ltr'}>
        <thead>
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className={cn(
                  'px-4 py-2 font-medium',
                  col.align === 'start' && 'text-start',
                  col.align === 'end' && 'text-end',
                  col.align === 'center' && 'text-center'
                )}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row) => (
            <tr key={row.id}>
              {columns.map((col) => (
                <td
                  key={col.key}
                  className={cn(
                    'px-4 py-2',
                    col.align === 'end' && 'text-end'
                  )}
                >
                  {row[col.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

## Modal and Dialog Patterns

### RTL Dialog

```tsx
import * as Dialog from '@radix-ui/react-dialog';

export function RTLDialog({ children, ...props }: DialogProps) {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return (
    <Dialog.Root {...props}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50" />
        <Dialog.Content
          className={cn(
            'fixed top-1/2 -translate-y-1/2 bg-white rounded-lg p-6',
            isRTL ? 'right-4' : 'left-4'
          )}
          dir={isRTL ? 'rtl' : 'ltr'}
        >
          {children}
          <Dialog.Close className={cn(
            'absolute top-4',
            isRTL ? 'left-4' : 'right-4'
          )}>
            <X className="w-4 h-4" />
          </Dialog.Close>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
```

## Chart Patterns

### RTL Chart Container

```tsx
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

export function useRTLOption(baseOption: EChartsOption): EChartsOption {
  const { i18n } = useTranslation();
  const isRTL = i18n.language === 'ar';

  return useMemo(() => {
    if (!isRTL) return baseOption;

    return {
      ...baseOption,
      grid: {
        left: baseOption.grid?.right ?? '3%',
        right: baseOption.grid?.left ?? '4%',
        top: baseOption.grid?.top ?? '10%',
        bottom: baseOption.grid?.bottom ?? '3%',
        containLabel: baseOption.grid?.containLabel ?? true,
      },
      legend: {
        ...baseOption.legend,
        align: 'right',
      },
      xAxis: baseOption.xAxis ? {
        ...baseOption.xAxis,
        inverse: true,
      } : undefined,
    };
  }, [baseOption, isRTL]);
}
```

## Testing RTL

### RTL Test Utilities

```tsx
import { render, RenderOptions } from '@testing-library/react';
import { I18nextProvider } from 'react-i18next';
import i18n from 'i18next';

export function renderRTL(
  ui: React.ReactElement,
  { locale = 'ar', ...options }: { locale?: string } & RenderOptions = {}
) {
  i18n.changeLanguage(locale);

  return render(
    <I18nextProvider i18n={i18n}>
      <div dir={locale === 'ar' ? 'rtl' : 'ltr'}>
        {ui}
      </div>
    </I18nextProvider>,
    options
  );
}

// Usage in tests
test('renders correctly in RTL', () => {
  const { container } = renderRTL(<YourComponent />);

  expect(container.querySelector('[dir="rtl"]')).toBeInTheDocument();
});
```
