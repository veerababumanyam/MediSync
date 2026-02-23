# Design Tokens Reference

Complete reference for AnySync design tokens. These are the single source of truth for all design values.

## Brand Colors

Extracted from the AnySync logo's visual story: bar chart (data) → heartbeat line (vitality).

### Primary: Royal Blue

Authority, data, structure.

```css
--brand-primary:      #1B3FA0;  /* Logo exact */
--brand-primary-50:   #EBF0FF;
--brand-primary-100:  #D6E0FF;
--brand-primary-200:  #ADC1FF;
--brand-primary-300:  #7596F5;
--brand-primary-400:  #4A6FDB;
--brand-primary-500:  #1B3FA0;
--brand-primary-600:  #163489;
--brand-primary-700:  #112972;
--brand-primary-800:  #0D1F5B;
--brand-primary-900:  #091544;
--brand-primary-950:  #050C2D;
--brand-primary-rgb:  27, 63, 160;
```

### Secondary: Teal-Cyan

Bridge, transition, flow.

```css
--brand-secondary:      #17B5A6;  /* Logo exact */
--brand-secondary-50:   #ECFDF9;
--brand-secondary-100:  #D1FAF0;
--brand-secondary-200:  #A5F3E2;
--brand-secondary-300:  #6DE8D0;
--brand-secondary-400:  #33D4BB;
--brand-secondary-500:  #17B5A6;
--brand-secondary-600:  #0F9185;
--brand-secondary-700:  #0B7068;
--brand-secondary-800:  #08544E;
--brand-secondary-900:  #053B37;
--brand-secondary-950:  #022421;
--brand-secondary-rgb:  23, 181, 166;
```

### Accent: Emerald Green

Vitality, pulse, action.

```css
--brand-accent:      #1CD760;  /* Logo heartbeat */
--brand-accent-50:   #EDFFF3;
--brand-accent-100:  #D5FFE4;
--brand-accent-200:  #ADFFC9;
--brand-accent-300:  #72F59E;
--brand-accent-400:  #3DE878;
--brand-accent-500:  #1CD760;
--brand-accent-600:  #14B04E;
--brand-accent-700:  #0F8A3D;
--brand-accent-800:  #0B6A2F;
--brand-accent-900:  #074D22;
--brand-accent-950:  #042F15;
--brand-accent-rgb:  28, 215, 96;
```

## Semantic Colors

### Success
```css
--color-success:       #10B981;
--color-success-light: #34D399;
--color-success-dark:  #059669;
--color-success-bg:    rgba(16, 185, 129, 0.08);
```

### Warning
```css
--color-warning:       #F59E0B;
--color-warning-light: #FBBF24;
--color-warning-dark:  #D97706;
--color-warning-bg:    rgba(245, 158, 11, 0.08);
```

### Error
```css
--color-error:         #EF4444;
--color-error-light:   #F87171;
--color-error-dark:    #DC2626;
--color-error-bg:      rgba(239, 68, 68, 0.08);
```

### Info
```css
--color-info:          #3B82F6;
--color-info-light:    #60A5FA;
--color-info-dark:     #2563EB;
--color-info-bg:       rgba(59, 130, 246, 0.08);
```

## Glassmorphism Tokens

### Light Mode Glass

```css
--glass-bg-pronounced:  rgba(255, 255, 255, 0.92);
--glass-bg-heavy:       rgba(255, 255, 255, 0.82);
--glass-bg-medium:      rgba(255, 255, 255, 0.72);
--glass-bg-light:       rgba(255, 255, 255, 0.60);
--glass-bg-subtle:      rgba(255, 255, 255, 0.35);

--glass-border-bright:  rgba(255, 255, 255, 0.65);
--glass-border-soft:    rgba(255, 255, 255, 0.40);
--glass-border-dim:     rgba(255, 255, 255, 0.18);

--glass-highlight:      rgba(255, 255, 255, 0.85);
--glass-highlight-soft: rgba(255, 255, 255, 0.45);
```

### Dark Mode Glass

```css
--glass-bg-pronounced:  rgba(15, 23, 42, 0.94);
--glass-bg-heavy:       rgba(15, 23, 42, 0.88);
--glass-bg-medium:      rgba(15, 23, 42, 0.80);
--glass-bg-light:       rgba(15, 23, 42, 0.70);
--glass-bg-subtle:      rgba(15, 23, 42, 0.55);

--glass-border-bright:  rgba(255, 255, 255, 0.14);
--glass-border-soft:    rgba(255, 255, 255, 0.08);
--glass-border-dim:     rgba(255, 255, 255, 0.04);
```

### Blur Scale

```css
--blur-xs:   blur(8px);
--blur-sm:   blur(16px);
--blur-md:   blur(32px);
--blur-lg:   blur(48px);
--blur-xl:   blur(64px);
--blur-xxl:  blur(80px);
--blur-max:  blur(100px);
```

## Shadow System

```css
/* Elevation Level 0 — Ambient */
--shadow-ambient:
  0 0 0 1px rgba(15, 23, 42, 0.04),
  0 1px 3px rgba(15, 23, 42, 0.03),
  0 6px 12px rgba(15, 23, 42, 0.04);

/* Elevation Level 1 — Raised */
--shadow-raised:
  0 0 0 1px rgba(15, 23, 42, 0.05),
  0 4px 8px rgba(15, 23, 42, 0.04),
  0 12px 24px rgba(15, 23, 42, 0.08);

/* Elevation Level 2 — Floating */
--shadow-floating:
  0 0 0 1px rgba(15, 23, 42, 0.06),
  0 8px 16px rgba(15, 23, 42, 0.06),
  0 24px 48px rgba(15, 23, 42, 0.10);

/* Elevation Level 3 — Elevated */
--shadow-elevated:
  0 0 0 1px rgba(255, 255, 255, 0.25),
  0 16px 48px rgba(15, 23, 42, 0.12),
  0 32px 64px rgba(15, 23, 42, 0.08);

/* Brand-colored glows */
--shadow-glow-primary:   0 0 24px rgba(27, 63, 160, 0.25);
--shadow-glow-secondary: 0 0 24px rgba(23, 181, 166, 0.25);
--shadow-glow-accent:    0 0 24px rgba(28, 215, 96, 0.30);
```

## Spacing Scale

4px base unit system:

```css
--space-0:   0;
--space-px:  1px;
--space-1:   4px;
--space-2:   8px;
--space-3:   12px;
--space-4:   16px;
--space-5:   20px;
--space-6:   24px;
--space-8:   32px;
--space-10:  40px;
--space-12:  48px;
--space-16:  64px;
--space-20:  80px;
--space-24:  96px;
```

## Border Radius

iOS continuous curvature approximation:

```css
--radius-xs:   4px;
--radius-sm:   8px;
--radius-md:   12px;
--radius-lg:   16px;
--radius-xl:   20px;
--radius-2xl:  24px;
--radius-3xl:  28px;
--radius-4xl:  32px;
--radius-full: 9999px;

/* Component-specific */
--radius-button: 16px;
--radius-card:   24px;
--radius-modal:  28px;
--radius-input:  12px;
```

## Typography

### Font Families

```css
--font-sans: 'SF Pro Display', 'SF Pro Text', -apple-system, BlinkMacSystemFont,
  'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
--font-mono: 'SF Mono', 'Fira Code', 'JetBrains Mono', monospace;
```

### Font Sizes

```css
--text-xs:   0.75rem;   /* 12px */
--text-sm:   0.875rem;  /* 14px */
--text-base: 1rem;      /* 16px */
--text-lg:   1.125rem;  /* 18px */
--text-xl:   1.25rem;   /* 20px */
--text-2xl:  1.5rem;    /* 24px */
--text-3xl:  1.875rem;  /* 30px */
--text-4xl:  2.25rem;   /* 36px */
--text-5xl:  3rem;      /* 48px */
```

### Fluid Display Sizes

```css
--text-display-sm: clamp(1.875rem, 4vw, 2.25rem);
--text-display-md: clamp(2.25rem, 5vw, 3rem);
--text-display-lg: clamp(3rem, 6vw, 4.5rem);
--text-display-xl: clamp(3.75rem, 8vw, 6rem);
```

## Animation Timing

```css
/* Durations */
--duration-instant:  75ms;
--duration-fast:     150ms;
--duration-base:     250ms;
--duration-moderate: 350ms;
--duration-slow:     500ms;

/* Easing */
--ease-liquid:        cubic-bezier(0.4, 0, 0.2, 1);
--ease-liquid-out:    cubic-bezier(0, 0, 0.2, 1);
--ease-spring:        cubic-bezier(0.34, 1.56, 0.64, 1);
--ease-spring-gentle: cubic-bezier(0.22, 1.2, 0.36, 1);
```

## Focus Tokens

```css
--focus-color:       #1B3FA0;  /* Light mode */
--focus-gold:        #FFD700;  /* Dark mode / high contrast */
--focus-gold-rgb:    255, 215, 0;
--focus-ring-width:  2px;
--focus-ring-offset: 2px;

/* Pre-composed focus ring */
--focus-ring: 0 0 0 var(--focus-ring-offset) var(--surface-primary),
              0 0 0 calc(var(--focus-ring-offset) + var(--focus-ring-width)) var(--focus-color);
```

## Touch Targets

```css
--touch-target-min:        44px;
--touch-target-comfortable: 48px;
```
