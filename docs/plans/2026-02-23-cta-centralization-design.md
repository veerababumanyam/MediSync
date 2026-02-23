# CTA Centralization Design Document

**Date:** 2026-02-23
**Status:** Approved
**Author:** Claude AI + User Collaboration

## Overview

Create a single source of truth for all CTA (Call-to-Action) button styles across the AnySync frontend. This design uses a **hybrid approach**: CSS semantic classes for 90% of cases (theme-driven via `.dark` selector), with TypeScript helper functions for special cases needing JS theme control.

## Problem Statement

Currently, CTA styles are hardcoded in multiple places:
- **HeroCarousel.tsx** (line ~106): Long inline class string `liquid-glass-pronounced liquid-glass-brand liquid-glass-hover-shimmer px-8 py-4 rounded-xl font-semibold text-white shadow-glow-secondary hover:shadow-glow-secondary transition-all duration-300 hover:-translate-y-1 liquid-focus-gold`
- **FinalCTA.tsx** (lines 28-31, 40-42): Two buttons with duplicated base classes and `isDark`-based Tailwind for primary vs secondary
- **ChatHeader.tsx** (line 51): Primary-style CTA with its own class string
- **ChatInput.tsx** (line 95): Submit button uses `liquid-glass-button-primary liquid-glass-hover-glow`

This leads to:
1. **Inconsistency**: Each button may have slightly different spacing, sizing, or effects
2. **Maintenance burden**: Changing a CTA style requires updating multiple files
3. **No visual language**: No semantic naming like "hero" or "landing-primary"

## Design Solution

### Hybrid Approach

| Use Case | Implementation | Rationale |
|----------|---------------|-----------|
| Hero CTA, Landing CTAs, Header CTAs | **CSS semantic classes** | Theme driven by root `.dark` class; fast path; no JS overhead |
| FinalCTA (special gradients) | **TS helper functions** | Needs per-component theme control; JS flexibility for conditional styling |

### CSS Semantic Classes

Add a new "CTA Variants" section to `liquid-glass.css` with three semantic classes:

#### `.cta-hero` - Hero Carousel CTA

The most prominent CTA, used in the hero carousel.

**Composition:**
- Layout: `inline-flex`, items-center, justify-center
- Size: `padding: 1rem 2rem` (px-8 py-4), `min-height: 44px`
- Typography: `font-weight: 600`, `font-size: 1rem`, `color: white`
- Shape: `border-radius: 0.75rem` (rounded-xl)
- Glass: `liquid-glass-pronounced`, `liquid-glass-brand`, `liquid-glass-hover-shimmer`
- Shadow: `--shadow-glow-secondary`
- Focus: `liquid-focus-gold`
- Transition: `all 0.3s cubic-bezier(0.4, 0, 0.2, 1)`
- Hover: `translateY(-4px)`, enhanced shadow

**Accessibility:**
- Touch target: 44px min-height (WCAG 2.5.5)
- Gold focus indicator (WCAG 3.0 Bronze)
- High contrast mode support

#### `.cta-landing-primary` - Landing Page Primary

Primary action button on landing pages and FinalCTA section.

**Light Mode:**
- Background: `linear-gradient(135deg, #ffffff 0%, #f8fafc 100%)`
- Text: `#2563eb` (blue-600) - APCA compliant contrast
- Shadow: White glow effect

**Dark Mode:**
- Background: `linear-gradient(135deg, #3b82f6 0%, #22d3ee 100%)`
- Text: `#0f172a` (slate-900) - APCA compliant contrast
- Shadow: Cyan glow effect

**Shared:**
- Layout: `inline-flex`, items-center, justify-center
- Size: `padding: 1rem 2rem`, `min-height: 44px`
- Typography: `font-weight: 700`, `font-size: 1.125rem` (text-lg)
- Shape: `border-radius: 0.75rem`
- Base: `liquid-glass-button-primary` (for structure)
- Focus: `liquid-focus-gold`
- Transition: `all 0.3s cubic-bezier(0.4, 0, 0.2, 1)`
- Hover: `translateY(-4px)`, enhanced shadow

#### `.cta-landing-secondary` - Landing Page Secondary

Secondary action button with bordered glass style.

**Light Mode:**
- Border: `1px solid rgba(255, 255, 255, 0.3)`
- Text: `white`

**Dark Mode:**
- Border: `1px solid rgba(255, 255, 255, 0.2)`
- Text: `white`

**Shared:**
- Layout: `inline-flex`, items-center, justify-center
- Size: `padding: 1rem 2rem`, `min-height: 44px`
- Typography: `font-weight: 700`, `font-size: 1.125rem`
- Shape: `border-radius: 0.75rem`
- Glass: `liquid-glass-pronounced`
- Focus: `liquid-focus-gold`
- Transition: `all 0.2s ease`
- Hover: Slight background tint

### TypeScript Module

Create `frontend/src/lib/cta-classes.ts`:

**Constants:**
```typescript
export const CTA_CLASSES = {
  HERO: 'cta-hero',
  LANDING_PRIMARY: 'cta-landing-primary',
  LANDING_SECONDARY: 'cta-landing-secondary',
  HEADER_PRIMARY: 'liquid-glass-button-primary liquid-focus-gold',
} as const
```

**Helper Functions** (for FinalCTA special case):
```typescript
/**
 * Get FinalCTA primary button classes with theme-specific styling
 * Special case: JS theme control for custom gradients
 */
export function getFinalCTAPrimaryClasses(isDark: boolean): string

/**
 * Get FinalCTA secondary button classes with theme-specific styling
 * Special case: JS theme control for border/text colors
 */
export function getFinalCTASecondaryClasses(isDark: boolean): string
```

### Component Refactoring

| Component | Current | New |
|-----------|---------|-----|
| **HeroCarousel.tsx** | Inline string | `className={CTA_CLASSES.HERO}` |
| **FinalCTA.tsx** | `isDark` branches | `getFinalCTAPrimaryClasses(isDark)` / `getFinalCTASecondaryClasses(isDark)` |
| **ChatHeader.tsx** | `liquid-glass-button-primary liquid-focus-gold liquid-glass-hover-shimmer` | `className={CTA_CLASSES.HEADER_PRIMARY}` |
| **ChatInput.tsx** | `liquid-glass-button-primary liquid-glass-hover-glow` | `className={CTA_CLASSES.HEADER_PRIMARY}` |

## Accessibility & Compliance

All CTA classes inherit WCAG 3.0 compliance from existing liquid-glass infrastructure:

| Requirement | Implementation |
|-------------|----------------|
| **Touch targets** | `min-height: 44px` (WCAG 2.5.5) |
| **Focus indicators** | `liquid-focus-gold` with 2px gold outline + glow |
| **Color contrast** | APCA-compliant text colors (primary: ~106 Lc, secondary: ~68 Lc) |
| **Reduced motion** | Inherits from `@media (prefers-reduced-motion: reduce)` |
| **High contrast mode** | Inherits from `@media (prefers-contrast: high)` |
| **Forced colors mode** | Inherits from `@media (forced-colors: active)` |
| **RTL support** | Text alignment uses logical properties |
| **Keyboard navigation** | Native button element with visible focus |

## Visual Differentiation: Light vs Dark Mode

### Light Mode CTAs
- **Hero**: Brand gradient on white glass backdrop
- **Primary**: White background with blue text, clean shadow
- **Secondary**: Transparent glass with white border

### Dark Mode CTAs
- **Hero**: Brand gradient on dark glass backdrop, enhanced glow
- **Primary**: Blue-to-cyan gradient with dark text, cyan shadow
- **Secondary**: Transparent glass with subtle white border

## Design Tokens Used

| Token | Purpose |
|-------|---------|
| `--brand-primary` | Logo blue (#2750a8) |
| `--brand-secondary` | Logo teal (#18929d) |
| `--gradient-brand` | Pre-defined brand gradient |
| `--shadow-glow-secondary` | Teal glow effect |
| `--shadow-raised` | Elevated shadow |
| `--focus-gold` | WCAG 3.0 Gold focus color (#FFD700) |
| `--radius-xl` | 20px border radius |
| `--duration-base` | 250ms transition duration |
| `--ease-liquid` | Cubic-bezier(0.4, 0, 0.2, 1) |

## Files to Create/Modify

**Create:**
- `frontend/src/lib/cta-classes.ts` — CTA class name constants and helper functions

**Modify:**
- `frontend/src/styles/liquid-glass.css` — Add CTA Variants section (~150 lines)
- `frontend/src/components/landing/HeroCarousel.tsx` — Use `CTA_CLASSES.HERO`
- `frontend/src/components/landing/FinalCTA.tsx` — Use helper functions
- `frontend/src/components/chat/ChatHeader.tsx` — Use `CTA_CLASSES.HEADER_PRIMARY`
- `frontend/src/components/chat/ChatInput.tsx` — Use `CTA_CLASSES.HEADER_PRIMARY`

## Implementation Notes

1. **Pronounced Glass Effect**: All CTAs use `liquid-glass-pronounced` or `liquid-glass-button-primary` as base for iOS 26 style liquid glass effect
2. **Hover Effects**: Consistent `translateY(-4px)` lift on hover for tactile feedback
3. **Focus States**: All CTAs use `liquid-focus-gold` for WCAG 3.0 Bronze compliance
4. **Spacing**: Consistent `px-8 py-4` (1rem 2rem) padding and `min-h-[44px]` for touch targets
5. **Theme Switching**: Automatic via `.dark` selector; no JS needed for most CTAs

## Future Enhancements

1. **Additional Variants**: Add `.cta-sm`, `.cta-lg` for size variants if needed
2. **CVA Integration**: Add variants to `LiquidGlassButton` component for React props API
3. **Animation Variants**: Add pulse, bounce variants for special CTAs
4. **Loading States**: Add `.cta-loading` with spinner for async actions

## Approval

- [x] Design approved by user
- [ ] Implementation completed
- [ ] Accessibility verified
- [ ] Cross-browser testing completed
