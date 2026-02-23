# Animation Timing Guidelines

Reference guide for consistent, accessible animations in AnySync following iOS 26 Liquid Motion principles.

## Duration Scale

| Duration | Value | Use Case |
|----------|-------|----------|
| Instant | 75ms | Micro-interactions, hover states |
| Fast | 150ms | Button press, toggle switch |
| Base | 250ms | Default transitions, focus states |
| Moderate | 350ms | Card hover, panel expansion |
| Slow | 500ms | Page transitions, modal open |
| Slower | 700ms | Complex animations |
| Slowest | 1000ms | Hero animations, onboarding |

## Easing Functions

### Standard Easings

```css
/* Default - natural feel */
--ease-default: cubic-bezier(0.25, 0.1, 0.25, 1);

/* Liquid - smooth, organic */
--ease-liquid: cubic-bezier(0.4, 0, 0.2, 1);

/* Liquid out - decelerate */
--ease-liquid-out: cubic-bezier(0, 0, 0.2, 1);

/* Liquid in - accelerate */
--ease-liquid-in: cubic-bezier(0.4, 0, 1, 1);
```

### Spring Easings (iOS-inspired)

```css
/* Spring - bouncy feel */
--ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);

/* Spring gentle - subtle bounce */
--ease-spring-gentle: cubic-bezier(0.22, 1.2, 0.36, 1);

/* Spring snappy - quick bounce */
--ease-spring-snappy: cubic-bezier(0.68, -0.2, 0.32, 1.2);
```

### Dramatic Easings

```css
/* For modals, sheets */
--ease-dramatic-in: cubic-bezier(0.6, 0, 0.7, 0.2);
--ease-dramatic-out: cubic-bezier(0.2, 0.8, 0.2, 1);
```

## Common Animation Patterns

### Button Press

```css
.button {
  transition: transform var(--duration-fast) var(--ease-liquid);
}

.button:active {
  transform: scale(0.97);
}
```

### Card Hover Lift

```css
.card {
  transition:
    transform var(--duration-base) var(--ease-liquid),
    box-shadow var(--duration-base) var(--ease-liquid);
}

.card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-floating);
}
```

### Modal Open

```css
@keyframes modal-enter {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal {
  animation: modal-enter var(--duration-moderate) var(--ease-spring-gentle);
}
```

### Shimmer Effect

```css
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

.skeleton {
  background: linear-gradient(90deg, #E8ECF3 25%, #F4F6FA 50%, #E8ECF3 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s var(--ease-liquid) infinite;
}
```

### Pulse Effect

```css
@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.8;
  }
}

.pulse-indicator {
  animation: pulse 2s var(--ease-liquid) infinite;
}
```

### Float Animation

```css
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.floating-element {
  animation: float 4s var(--ease-liquid) infinite;
}
```

## Staggered Animations

For sequential entrance animations:

```css
.animate-stagger > *:nth-child(1) { animation-delay: 0ms; }
.animate-stagger > *:nth-child(2) { animation-delay: 100ms; }
.animate-stagger > *:nth-child(3) { animation-delay: 200ms; }
.animate-stagger > *:nth-child(4) { animation-delay: 300ms; }
.animate-stagger > *:nth-child(5) { animation-delay: 400ms; }
```

Or using CSS custom properties:

```tsx
<div className="animate-stagger">
  {items.map((item, index) => (
    <div
      key={item.id}
      style={{ '--delay': `${index * 100}ms` }}
      className="animate-item"
    >
      {item.content}
    </div>
  ))}
</div>
```

```css
.animate-item {
  animation: liquid-fade-in var(--duration-slow) var(--ease-liquid-out) forwards;
  animation-delay: var(--delay);
  opacity: 0;
}
```

## Reduced Motion

Always respect user preferences:

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

### Component-Specific Reduced Motion

```css
.card {
  transition: transform var(--duration-base) var(--ease-liquid);
}

@media (prefers-reduced-motion: reduce) {
  .card {
    transition: none;
  }
}

/* Still show focus even with reduced motion */
.card:focus-visible {
  outline: 2px solid var(--focus-gold);
  outline-offset: 2px;
}
```

## Performance Tips

1. **Use transform and opacity** - These are GPU-accelerated
2. **Avoid animating layout properties** - width, height, margin, padding
3. **Use will-change sparingly** - Only for known animations
4. **Batch animations** - Use animation-delay for stagger effects
5. **Test on low-end devices** - Animations can cause frame drops

## Animation Checklist

- [ ] All animations respect `prefers-reduced-motion`
- [ ] Animations don't cause layout thrashing
- [ ] Focus states work even with animations disabled
- [ ] Loading states provide feedback without being distracting
- [ ] Animations enhance understanding, not just decoration
