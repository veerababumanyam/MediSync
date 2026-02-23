# RTL (Right-to-Left) Patterns

Comprehensive guide for building RTL-compatible components that work seamlessly for Arabic and other RTL languages.

## Core Principles

1. **Use logical properties** instead of physical (left/right)
2. **Test with actual RTL content** - not just CSS direction
3. **Consider icon directionality** - some icons need mirroring
4. **Handle numbers and dates correctly** - LTR within RTL

## CSS Logical Properties

### Inset Properties

```css
/* Physical (AVOID) */
left: 0;
right: 0;

/* Logical (PREFER) */
inset-inline-start: 0;
inset-inline-end: 0;
```

### Margin Properties

```css
/* Physical (AVOID) */
margin-left: 1rem;
margin-right: 1rem;

/* Logical (PREFER) */
margin-inline-start: 1rem;
margin-inline-end: 1rem;
```

### Padding Properties

```css
/* Physical (AVOID) */
padding-left: 1rem;
padding-right: 1rem;

/* Logical (PREFER) */
padding-inline-start: 1rem;
padding-inline-end: 1rem;
```

### Border Properties

```css
/* Physical (AVOID) */
border-left: 1px solid;
border-right: 1px solid;

/* Logical (PREFER) */
border-inline-start: 1px solid;
border-inline-end: 1px solid;
```

### Text Alignment

```css
/* Physical (AVOID) */
text-align: left;
text-align: right;

/* Logical (PREFER) */
text-align: start;
text-align: end;
```

### Float

```css
/* Physical (AVOID) */
float: left;
float: right;

/* Logical (PREFER) */
float: inline-start;
float: inline-end;
```

## Tailwind RTL Utilities

Tailwind CSS provides logical property utilities:

```tsx
// Start/End instead of Left/Right
<div className="ms-4">     {/* margin-inline-start: 1rem */}
<div className="me-4">     {/* margin-inline-end: 1rem */}
<div className="ps-4">     {/* padding-inline-start: 1rem */}
<div className="pe-4">     {/* padding-inline-end: 1rem */}
<div className="text-start">  {/* text-align: start */}
<div className="text-end">    {/* text-align: end */}
```

## Component Patterns

### RTL-Safe Flex Layout

```tsx
function IconText({ icon, text }: { icon: React.ReactNode; text: string }) {
  return (
    <div className="flex items-center gap-2">
      <span className="shrink-0">{icon}</span>
      <span className="text-start">{text}</span>
    </div>
  );
}
```

### RTL-Safe Card with Arrow

```tsx
function CardWithArrow({ title }: { title: string }) {
  return (
    <div className="liquid-container-card flex items-center justify-between">
      <h3 className="text-start">{title}</h3>
      <ChevronRightIcon className="shrink-0 rtl:rotate-180" />
    </div>
  );
}
```

### RTL-Safe Navigation Breadcrumb

```tsx
function Breadcrumb({ items }: { items: string[] }) {
  return (
    <nav aria-label="Breadcrumb">
      <ol className="flex items-center gap-2">
        {items.map((item, index) => (
          <li key={index} className="flex items-center gap-2">
            {index > 0 && (
              <ChevronRightIcon className="h-4 w-4 rtl:rotate-180 text-muted" />
            )}
            <span className="text-start">{item}</span>
          </li>
        ))}
      </ol>
    </nav>
  );
}
```

## Icon Mirroring

Some icons need to be mirrored for RTL:

### Icons That Need Mirroring

- Arrows (back, forward, chevrons)
- Bulleted lists
- Progress indicators
- Navigation icons
- Speech bubbles

### Icons That Stay the Same

- Checkmarks
- X marks
- Warning triangles
- Plus/minus signs
- Play/pause

### Tailwind RTL Icon Rotation

```tsx
// Using Tailwind's rtl: variant
<ChevronRightIcon className="rtl:rotate-180" />

// Or using CSS
.icon-directional {
  transform: rotate(0deg);
}

[dir="rtl"] .icon-directional {
  transform: rotate(180deg);
}
```

## Form Patterns

### RTL-Safe Form Labels

```tsx
function FormField({ label, id, ...inputProps }: FormFieldProps) {
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="text-start font-medium">
        {label}
      </label>
      <input id={id} {...inputProps} className="liquid-glass-input" />
    </div>
  );
}
```

### RTL-Safe Input with Icon

```tsx
function SearchInput() {
  return (
    <div className="relative">
      <input
        type="search"
        className="liquid-glass-input ps-10"
        placeholder="Search..."
      />
      <SearchIcon className="absolute start-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted" />
    </div>
  );
}
```

## Table Patterns

### RTL-Safe Table

```tsx
function DataTable({ columns, data }: TableProps) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className="text-start p-4 border-b border-default"
              >
                {col.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, index) => (
            <tr key={index}>
              {columns.map((col) => (
                <td key={col.key} className="text-start p-4">
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

## Numbers and Dates in RTL

Numbers and dates should remain LTR within RTL content:

```tsx
// Wrap numbers in LTR span
function RTLTextWithNumber({ text, number }: { text: string; number: number }) {
  return (
    <span className="text-start">
      {text} <span dir="ltr" className="inline-block">{number}</span>
    </span>
  );
}

// Format dates properly
function FormattedDate({ date }: { date: Date }) {
  const formatted = new Intl.DateTimeFormat('ar-SA', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);

  return <span dir="ltr" className="inline-block">{formatted}</span>;
}
```

## Testing RTL

### Enable RTL in i18next

```tsx
// In your i18n config
i18n.on('languageChanged', (lng) => {
  const dir = i18n.dir(lng);
  document.documentElement.dir = dir;
  document.documentElement.lang = lng;
});
```

### CSS Testing

```css
/* Force RTL for testing */
html {
  direction: rtl;
}

/* Or use class */
.rtl-test {
  direction: rtl;
}
```

### Browser Extensions

- RTL Switch (Chrome/Firefox)
- Language Force

## Common Pitfalls

1. **Fixed positioning** - May not flip correctly
2. **Absolute positioning with left/right** - Use inset-inline-*
3. **Transform translateX** - Use logical values or CSS custom properties
4. **Box shadows with spread** - Direction-independent
5. **Gradients** - Use logical keywords (to inline-start, etc.)

## RTL Checklist

- [ ] All layouts use logical properties (start/end)
- [ ] Text alignment uses text-start/text-end
- [ ] Directional icons are mirrored
- [ ] Forms have proper label alignment
- [ ] Tables work in RTL
- [ ] Numbers and dates display correctly
- [ ] No overflow issues with long RTL text
- [ ] Tested with Arabic content
