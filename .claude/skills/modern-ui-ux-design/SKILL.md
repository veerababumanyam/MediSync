---
name: modern-ui-ux-design
description: This skill should be used when the user asks to "create UI components", "design interfaces", "implement WCAG 3.0 accessibility", "add liquid glass effects", "iOS 26 design", "responsive layouts", "modern UI patterns", "APCA contrast", "touch targets", "focus indicators", "design tokens", "RTL support", "animation patterns", or mentions UI/UX design, accessibility compliance, or visual design principles for AnySync.
---

# Modern UI/UX Design for AnySync

Comprehensive guide for building accessible, elegant, and responsive interfaces following WCAG 3.0 accessibility standards, Apple iOS 26 Interface Human Guidelines (IHG) liquid glass morphism, and modern design principles.

## Quick Reference

| Aspect | Standard |
|--------|----------|
| **Accessibility** | WCAG 3.0 Bronze (minimum) |
| **Contrast** | APCA 60 Lc (body), 45 Lc (large text) |
| **Touch Targets** | 44px minimum (WCAG 2.5.8) |
| **Focus Indicators** | Gold 2px ring with 2px offset |
| **Design Language** | iOS 26 Liquid Glass |
| **Motion** | Reduced motion support required |
| **RTL** | Full Arabic support via logical properties |

## Design Philosophy

AnySync's design embodies "Data Meets Vitality" through:
- **Liquid Glass** — Translucent, depth-layered surfaces with organic motion
- **Accessibility First** — WCAG 3.0 Bronze minimum, Silver target
- **Responsive** — Mobile-first, fluid scaling across devices
- **Elegant Simplicity** — Clear visual hierarchy, minimal cognitive load
- **Inclusive** — High contrast, reduced motion, forced-colors support

## 1. WCAG 3.0 Accessibility Compliance

### APCA Contrast Requirements

WCAG 3.0 uses APCA (Accessible Perceptual Contrast Algorithm) for perceptually accurate contrast:

```css
/* Light Mode Text Colors (on white #FFFFFF) */
--text-primary:     #0F172A;  /* ~106 Lc - excellent */
--text-secondary:   #475569;  /* ~68 Lc - body text compliant */
--text-tertiary:    #576574;  /* ~61 Lc - body text minimum */
--text-muted:       #6B7280;  /* ~57 Lc - large text only (24px+) */

/* Dark Mode Text Colors (on #0F172A) */
--text-primary:     #F8FAFC;  /* ~107 Lc - excellent */
--text-secondary:   #CBD5E1;  /* ~75 Lc - body text compliant */
--text-tertiary:    #94A3B8;  /* ~62 Lc - body text minimum */
```

### Touch Target Requirements

```css
/* Minimum touch target - WCAG 2.5.8 */
.touch-target {
  min-width: 44px;
  min-height: 44px;
}

/* Comfortable touch target */
.touch-target-comfortable {
  min-width: 48px;
  min-height: 48px;
}

/* Apply with padding to maintain visual size */
.button-compact {
  padding: 8px 16px;  /* Visual size */
  min-height: 44px;   /* Touch target */
}
```

### Focus Indicators (WCAG 3.0 Gold)

```css
/* Gold focus ring for maximum visibility */
.focus-gold:focus-visible {
  outline: 2px solid #FFD700;
  outline-offset: 2px;
  box-shadow:
    0 0 0 4px rgba(255, 215, 0, 0.15),
    var(--shadow-elevation);
}

/* Dark mode - lighter gold */
.dark .focus-gold:focus-visible {
  outline-color: #FFED4E;
  box-shadow:
    0 0 0 4px rgba(255, 237, 78, 0.2),
    0 0 20px rgba(255, 237, 78, 0.4);
}
```

### Accessibility Media Queries

```css
/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  .animate-element {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}

/* High contrast mode */
@media (prefers-contrast: high) {
  .focus-gold:focus-visible {
    outline-width: 3px;
    outline-offset: 3px;
  }
}

/* Windows High Contrast */
@media (forced-colors: active) {
  .liquid-glass-button {
    background: ButtonFace;
    border: 2px solid ButtonText;
    color: ButtonText;
  }
}

/* Reduced transparency */
@media (prefers-reduced-transparency: reduce) {
  .liquid-glass {
    backdrop-filter: none !important;
    background-color: #F8FAFC !important;
  }
}
```

## 2. iOS 26 Liquid Glass Morphism

### Core Principles

1. **Translucency** — Background shows through with layered opacity
2. **Blur** — Frosted glass backdrop-filter creates depth
3. **Refraction** — Borders simulate light bending
4. **Specular** — Top-edge highlights simulate light source above
5. **Tint** — Glass absorbs underlying brand color
6. **Depth** — Layered panels at different opacities
7. **Continuity** — Rounded corners with continuous curvature

### Glass Layer Hierarchy

```css
/* Light Mode Glass Layers (most to least opaque) */
--glass-bg-pronounced:  rgba(255, 255, 255, 0.92);  /* Modals, key UI */
--glass-bg-heavy:       rgba(255, 255, 255, 0.82);  /* Cards, panels */
--glass-bg-medium:      rgba(255, 255, 255, 0.72);  /* Default surfaces */
--glass-bg-light:       rgba(255, 255, 255, 0.60);  /* Headers, nav */
--glass-bg-subtle:      rgba(255, 255, 255, 0.35);  /* Backgrounds */

/* Dark Mode Glass Layers */
--glass-bg-pronounced:  rgba(15, 23, 42, 0.94);
--glass-bg-heavy:       rgba(15, 23, 42, 0.88);
--glass-bg-medium:      rgba(15, 23, 42, 0.80);
--glass-bg-light:       rgba(15, 23, 42, 0.70);
--glass-bg-subtle:      rgba(15, 23, 42, 0.55);
```

### Base Glass Component

```tsx
// components/ui/LiquidGlass.tsx
interface LiquidGlassProps {
  variant?: 'subtle' | 'light' | 'medium' | 'heavy' | 'pronounced';
  elevation?: 'ambient' | 'raised' | 'floating' | 'elevated';
  interactive?: boolean;
  className?: string;
  children: React.ReactNode;
}

export function LiquidGlass({
  variant = 'medium',
  elevation = 'raised',
  interactive = false,
  className,
  children
}: LiquidGlassProps) {
  const baseClasses = [
    'liquid-glass',
    `liquid-glass-${variant}`,
    `liquid-shadow-${elevation}`,
    interactive && 'liquid-glass-hover-lift liquid-focus-gold',
    className
  ].filter(Boolean).join(' ');

  return (
    <div
      className={baseClasses}
      tabIndex={interactive ? 0 : undefined}
    >
      {children}
    </div>
  );
}
```

### Blur Intensity Scale

```css
/* Backdrop blur values */
--blur-xs:   blur(8px);   /* Subtle glass */
--blur-sm:   blur(16px);  /* Light glass */
--blur-md:   blur(32px);  /* Medium glass */
--blur-lg:   blur(48px);  /* Heavy glass */
--blur-xl:   blur(64px);  /* Pronounced glass */
--blur-xxl:  blur(80px);  /* Hero sections */
```

## 3. Design Tokens Reference

### Brand Colors (Logo-Extracted)

```css
/* Primary: Royal Blue (Authority, Data, Structure) */
--brand-primary: #1B3FA0;

/* Secondary: Teal-Cyan (Bridge, Transition, Flow) */
--brand-secondary: #17B5A6;

/* Accent: Emerald Green (Vitality, Pulse, Action) */
--brand-accent: #1CD760;
```

### Spacing Scale (4px Base Unit)

```css
--space-1:  4px;
--space-2:  8px;
--space-3:  12px;
--space-4:  16px;
--space-5:  20px;
--space-6:  24px;
--space-8:  32px;
--space-10: 40px;
--space-12: 48px;
```

### Border Radius (iOS Continuous Curvature)

```css
--radius-sm:   8px;   /* Small elements */
--radius-md:   12px;  /* Buttons, inputs */
--radius-lg:   16px;  /* Cards */
--radius-xl:   20px;  /* Large cards */
--radius-2xl:  24px;  /* Modals, sheets */
--radius-3xl:  28px;  /* Hero elements */
--radius-full: 9999px; /* Pills, avatars */
```

### Animation Easing (iOS Liquid Motion)

```css
/* Standard easing */
--ease-liquid:     cubic-bezier(0.4, 0, 0.2, 1);
--ease-liquid-out: cubic-bezier(0, 0, 0.2, 1);

/* Spring-like curves */
--ease-spring:        cubic-bezier(0.34, 1.56, 0.64, 1);
--ease-spring-gentle: cubic-bezier(0.22, 1.2, 0.36, 1);

/* Duration */
--duration-fast:    150ms;
--duration-base:    250ms;
--duration-moderate: 350ms;
--duration-slow:    500ms;
```

## 4. Component Patterns

### Button Variants

```tsx
// Primary CTA - Hero/Landing pages
<button className="cta-hero">
  Get Started
</button>

// Header button - Compact
<button className="cta-header-primary">
  Subscribe
</button>

// Secondary button
<button className="cta-landing-secondary">
  Learn More
</button>

// Glass button - Transparent background
<button className="liquid-glass-button liquid-glass-button-prominent">
  <Icon />
  Action
</button>
```

### Card Patterns

```tsx
// Standard glass card
<div className="liquid-container-card">
  <h3>Card Title</h3>
  <p>Card content...</p>
</div>

// Interactive card with hover
<div className="liquid-container-card liquid-container-card--interactive" tabIndex={0}>
  <h3>Clickable Card</h3>
</div>

// Pronounced card (featured content)
<div className="liquid-container-card liquid-container-card--pronounced">
  <h2>Featured Content</h2>
</div>
```

### Form Patterns

```tsx
// WCAG 3.0 compliant form input
<div className="liquid-clear-label">
  <label htmlFor="email">Email Address</label>
  <input
    id="email"
    type="email"
    className="liquid-glass-input"
    aria-describedby="email-hint"
  />
  <span id="email-hint" className="text-tertiary">
    We'll never share your email
  </span>
</div>

// Error state
<input
  className="liquid-glass-input state-error"
  aria-invalid="true"
  aria-errormessage="email-error"
/>
<span id="email-error" role="alert" className="text-error">
  Please enter a valid email
</span>
```

## 5. Responsive Design Patterns

### Breakpoint System

```css
/* Mobile-first breakpoints */
--breakpoint-sm:  640px;   /* Small tablets */
--breakpoint-md:  768px;   /* Tablets */
--breakpoint-lg:  1024px;  /* Laptops */
--breakpoint-xl:  1280px;  /* Desktops */
--breakpoint-2xl: 1536px;  /* Large screens */
```

### Fluid Typography

```css
/* Display sizes with clamp() */
--text-display-sm: clamp(1.875rem, 4vw, 2.25rem);   /* 30-36px */
--text-display-md: clamp(2.25rem, 5vw, 3rem);       /* 36-48px */
--text-display-lg: clamp(3rem, 6vw, 4.5rem);        /* 48-72px */
```

### Responsive Glass Cards

```tsx
// Responsive grid with glass cards
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
  {items.map(item => (
    <div key={item.id} className="liquid-container-card">
      {/* Card content */}
    </div>
  ))}
</div>
```

## 6. RTL Support

### Logical Properties

```css
/* Use logical properties instead of physical */
.inset-start { inset-inline-start: 0; }
.inset-end { inset-inline-end: 0; }
.margin-start { margin-inline-start: 1rem; }
.margin-end { margin-inline-end: 1rem; }
.padding-start { padding-inline-start: 1rem; }
.padding-end { padding-inline-end: 1rem; }

/* RTL-safe text alignment */
.text-start { text-align: start; }
.text-end { text-align: end; }
```

### RTL Component Pattern

```tsx
function RTLCard({ icon, title, description }: RTLCardProps) {
  return (
    <div className="flex gap-4 p-4 liquid-container-card">
      <div className="shrink-0">{icon}</div>
      <div className="flex-1">
        <h3 className="text-start font-semibold">{title}</h3>
        <p className="text-start text-secondary">{description}</p>
      </div>
    </div>
  );
}
```

## 7. Animation Patterns

### Entrance Animations

```css
/* Liquid fade-in */
@keyframes liquid-fade-in {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.96);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.liquid-animate-in {
  animation: liquid-fade-in var(--duration-slow) var(--ease-liquid-out) forwards;
}

/* Staggered delays */
.liquid-delay-100 { animation-delay: 100ms; }
.liquid-delay-200 { animation-delay: 200ms; }
.liquid-delay-300 { animation-delay: 300ms; }
```

### Hover Effects

```css
/* Lift effect */
.liquid-glass-hover-lift:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-float);
}

/* Glow effect */
.liquid-glass-hover-glow:hover {
  box-shadow:
    var(--shadow-elevation),
    0 0 40px rgba(23, 181, 166, 0.3);
}

/* Shimmer effect */
.liquid-glass-hover-shimmer:hover::after {
  /* Moving light reflection */
}
```

### Loading States

```tsx
// Skeleton loader
<div className="liquid-skeleton h-4 w-3/4 rounded" />
<div className="liquid-skeleton h-4 w-1/2 rounded mt-2" />

// Spinner
<div className="liquid-spinner" role="status" aria-label="Loading">
  <span className="sr-only">Loading...</span>
</div>
```

## 8. Accessibility Checklist

Before committing UI components, verify:

- [ ] All interactive elements have visible focus indicators (2px gold ring)
- [ ] Touch targets are minimum 44x44px
- [ ] Color contrast meets APCA 60 Lc (body) / 45 Lc (large text)
- [ ] Animations respect `prefers-reduced-motion`
- [ ] Images have meaningful alt text
- [ ] Form inputs have associated labels
- [ ] Error messages are announced to screen readers
- [ ] Focus order follows visual layout
- [ ] Color is not the only means of conveying information
- [ ] High contrast mode is supported
- [ ] Windows forced-colors mode is supported

## Additional Resources

### Reference Files
- **`references/wcag-3-0-checklist.md`** — Complete WCAG 3.0 compliance checklist
- **`references/design-tokens.md`** — Full token reference
- **`references/animation-timings.md`** — Animation timing guidelines
- **`references/rtl-patterns.md`** — Comprehensive RTL patterns

### Example Files
- **`examples/LiquidGlass.tsx`** — Complete glass component implementation
- **`examples/CTAButtons.tsx`** — CTA button variants
- **`examples/FormFields.tsx`** — Accessible form components
- **`examples/Cards.tsx`** — Card component patterns

### Related Skills
- **react-typescript** — React component patterns
- **flutter-dart** — Mobile UI patterns
- **echarts-visualization** — Data visualization
- **i18next-i18n** — Internationalization and RTL

### External Resources
- [WCAG 3.0 Working Draft](https://www.w3.org/TR/wcag-3.0/)
- [Apple Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- [APCA Contrast Calculator](https://www.myndex.com/APCA/)
