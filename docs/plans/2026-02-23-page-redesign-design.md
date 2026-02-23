# Page Redesign with Design System v3.0.0

**iOS 26 Liquid Glass · WCAG 3.0 · Light + Dark Mode**

Date: 2026-02-23
Status: Approved

---

## Overview

Comprehensive overhaul of all 4 application pages (`Home.tsx`, `ChatPage.tsx`, `DashboardPage.tsx`, `CouncilPage.tsx`) to fully align with Design System v3.0.0.

### Goals

1. **No Hardcoded Colors**: All colors via CSS variables from `design-tokens.css`
2. **Lucide-Style Icons**: Replace emojis with consistent outline SVG icons
3. **Clear Dark/Light Separation**: Proper contrast and visibility in both modes
4. **Moderate Animations**: Professional hover effects and transitions
5. **WCAG 3.0 Bronze Compliance**: Accessible contrast ratios

---

## 1. Design Token Strategy

All pages use ONLY CSS variables. No hardcoded Tailwind colors.

### Token Mapping

| Hardcoded (Old) | Token (New) | CSS Variable |
|-----------------|-------------|--------------|
| `text-slate-900` | `text-primary` | `--text-primary` |
| `text-slate-400` | `text-secondary` | `--text-secondary` |
| `bg-blue-500/10` | `bg-surface-glass` or brand glass | `--glass-bg-brand-blue` |
| `text-emerald-400` | `text-accent` | `--text-accent` |
| `border-blue-200` | `border-subtle` | `--border-subtle` |
| `bg-gray-400` | `bg-surface-tertiary` | `--surface-tertiary` |
| `text-green-500` | `text-success` | `--color-success` |

### Utility Classes

Use predefined glass classes instead of ad-hoc styles:

```tsx
// Instead of:
className="bg-white/80 backdrop-blur-xl border border-white/20"

// Use:
className="glass-panel"
```

---

## 2. Icon System

### Lucide-Style SVG Icons

- **Stroke Width**: 1.5px
- **Size**: 24px default (w-6 h-6)
- **Color**: `currentColor` for theme adaptation
- **Accessibility**: `aria-hidden="true"` with visible labels

### Icons Needed

| Icon Name | Component | Usage |
|-----------|-----------|-------|
| ChatBubble | `ChatBubbleIcon` | Conversational BI feature |
| ArrowsExchange | `ArrowsExchangeIcon` | Tally Sync feature |
| DocumentAI | `DocumentAIIcon` | AI Accountant feature |
| ShieldCheck | `ShieldCheckIcon` | PII Protection feature |
| ChartLine | `ChartLineIcon` | Prescriptive Analytics |
| Link | `LinkIcon` | HIMS Connectivity |
| QuestionCircle | `QuestionCircleIcon` | FAQ section |
| Revenue | `RevenueIcon` | Dashboard revenue KPI |
| Users | `UsersIcon` | Dashboard patients |
| Calendar | `CalendarIcon` | Dashboard appointments |
| Package | `PackageIcon` | Dashboard inventory |
| Lightning | `LightningIcon` | Quick actions |
| Error | `ErrorIcon` | Error states |
| Check | `CheckIcon` | Success states |
| Heart | `HeartIcon` | Health indicators |

---

## 3. Page-Specific Updates

### 3.1 Home.tsx

**Changes:**
- Remove all conditional `isDark ? '...' : '...'` color classes
- Use semantic tokens that automatically switch
- Replace 6 feature card emojis with Lucide SVG icons
- Add `.glass-card` hover effects
- Animate feature cards with staggered `animate-fade-in-up`
- Hero section with `backdrop-filter` blur

**Token Usage:**
```tsx
// Badge
className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full
  text-xs font-semibold tracking-wider uppercase
  bg-glass-bg-brand-blue border border-glass-border-tinted text-brand-primary"

// Feature Cards
className="glass-card hover:shadow-floating transition-all duration-150"
```

### 3.2 ChatPage.tsx

**Changes:**
- Chat interface container: `.liquid-glass-pronounced`
- Error states: Use `--color-error-bg` semantic token
- Focus rings: Gold in dark mode, brand blue in light
- Add `liquid-focus-gold` for keyboard navigation

**Token Usage:**
```tsx
// Error container
className="liquid-glass-pronounced mt-4 p-4
  border-s-4 border-s-[var(--color-error)]
  bg-[var(--color-error-bg)]"
```

### 3.3 DashboardPage.tsx

**Changes:**
- KPI cards: Use `GlassBrandCard`, `GlassTealCard`, `GlassBlueCard`, `GlassGreenCard`
- Replace emoji icons (💰, 👥, 📅, 📦) with Lucide SVGs
- Add `animate-pulse-slow` to live indicators
- Quick actions: `.liquid-glass-hover-lift` on hover

**Token Usage:**
```tsx
// KPI Icon container
className="w-12 h-12 rounded-xl bg-gradient-to-br
  from-brand-primary to-brand-secondary
  flex items-center justify-center
  shadow-glow-primary animate-pulse-slow"
```

### 3.4 CouncilPage.tsx

**Changes:**
- Health status indicators: Use semantic success/warning/error tokens
- Header: `.glass-nav` pattern
- Query input section: `.glass-panel`
- Response display: `.glass-card`
- Add proper dark mode contrast for all text

**Token Usage:**
```tsx
// Health indicator
const healthColors = {
  healthy: 'bg-[var(--color-success)]',
  degraded: 'bg-[var(--color-warning)]',
  failed: 'bg-[var(--color-error)]',
}
```

---

## 4. Animation System

### Moderate/Professional Animations

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Cards hover | `translateY(-2px) + shadow-raised` | 150ms | `ease-liquid` |
| Page entry | Staggered fade-in | 250ms each | `ease-liquid-out` |
| Buttons hover | `scale(1.02)` | 150ms | `ease-spring-gentle` |
| Icons hover | `scale(1.1)` | 150ms | `ease-liquid` |
| Focus rings | Ring expansion | 150ms | `ease-liquid` |
| Panel reveal | `opacity + translateY` | 350ms | `ease-liquid-out` |

### CSS Implementation
```css
/* Already defined in design-tokens.css */
--duration-fast: 150ms;
--duration-base: 250ms;
--ease-liquid: cubic-bezier(0.4, 0, 0.2, 1);
--ease-spring-gentle: cubic-bezier(0.22, 1.2, 0.36, 1);
```

---

## 5. Dark/Light Mode Separation

### Light Mode Palette

| Token | Value | Usage |
|-------|-------|-------|
| `--surface-primary` | `#FFFFFF` | Cards, inputs |
| `--surface-secondary` | `#F4F6FA` | Page background |
| `--text-primary` | `#0F172A` | Headings, body |
| `--text-secondary` | `#475569` | Descriptions |
| `--text-brand` | `#1B3FA0` | Brand text |
| `--glass-bg-medium` | `rgba(255,255,255,0.72)` | Glass panels |
| `--focus-color` | `#1B3FA0` | Focus rings |

### Dark Mode Palette

| Token | Value | Usage |
|-------|-------|-------|
| `--surface-primary` | `#0A0F1A` | Main background |
| `--surface-secondary` | `#111827` | Card backgrounds |
| `--text-primary` | `#F1F5F9` | Headings, body |
| `--text-secondary` | `#94A3B8` | Descriptions |
| `--text-brand` | `#7596F5` | Brand text |
| `--glass-bg-medium` | `rgba(15,23,42,0.80)` | Glass panels |
| `--focus-color` | `#FFD700` | Gold focus rings |

---

## 6. Accessibility (WCAG 3.0 Bronze)

### Contrast Requirements

| Element | Light Mode | Dark Mode |
|---------|------------|-----------|
| Body text | ≥ 60 Lc | ≥ 60 Lc |
| Large text | ≥ 45 Lc | ≥ 45 Lc |
| Focus indicators | ≥ 3:1 | ≥ 3:1 (gold) |
| Non-text | ≥ 45 Lc | ≥ 45 Lc |

### Implementation

- Focus rings: 3px solid, 3px offset
- Touch targets: 44×44px minimum
- Screen reader: `.sr-only` for hidden labels
- Reduced motion: `prefers-reduced-motion` respected

---

## 7. Implementation Approach

### Parallel Agent Team

4 agents update pages simultaneously:

1. **Agent-Home**: Update `Home.tsx` with icons and tokens
2. **Agent-Chat**: Update `ChatPage.tsx` with glass effects
3. **Agent-Dashboard**: Update `DashboardPage.tsx` with KPI icons
4. **Agent-Council**: Update `CouncilPage.tsx` with status tokens

### Shared Dependencies

All agents use:
- `design-tokens.css` (already exists)
- New `Icons.tsx` component (create first)
- Glass utility classes from `globals.css`

---

## 8. Success Criteria

- [ ] No hardcoded Tailwind color classes (slate-*, blue-*, etc.)
- [ ] All icons are Lucide-style SVGs
- [ ] Dark mode has clear visibility and contrast
- [ ] All pages pass WCAG 3.0 Bronze
- [ ] Hover effects on all interactive elements
- [ ] Consistent glass morphism styling
- [ ] Build passes without errors

---

*Design approved: 2026-02-23*
