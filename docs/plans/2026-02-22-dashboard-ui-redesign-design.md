# Dashboard UI/UX Redesign Design Document

**Date:** 2026-02-22
**Status:** Approved
**Approach:** Unified Component Architecture

---

## Overview

Redesign the Dashboard page and all its components to match the landing page's iOS 26 liquid glass design system, ensuring consistent dark/light mode theming across the application.

---

## Design Goals

1. **Consistency:** Match landing page FeatureCard styling for KPI cards
2. **Theme Support:** Proper `isDark` prop propagation for dark/light mode
3. **Maintainability:** Create shared `MetricCard` component for reuse
4. **Design System Compliance:** Use `liquid-glass-content-card` and related classes

---

## Architecture

### Theme Propagation Flow

```
App.tsx (derives isDark from theme context)
    └── DashboardPage.tsx (receives isDark)
            ├── LiquidGlassHeader (already has theme support)
            ├── KPI Section with MetricCard components (receives isDark)
            ├── DashboardGrid (receives isDark)
            │   └── PinnedChartCard (receives isDark)
            └── Quick Actions Section (receives isDark)
```

### Files to Modify/Create

| File | Action | Purpose |
|------|--------|---------|
| `components/dashboard/MetricCard.tsx` | **Create** | New shared metric card component |
| `pages/DashboardPage.tsx` | **Modify** | Add isDark prop, update KPI section |
| `components/dashboard/DashboardGrid.tsx` | **Modify** | Add isDark prop support |
| `components/dashboard/PinnedChartCard.tsx` | **Modify** | Add isDark prop for theme-aware styling |

---

## Component Specifications

### 1. MetricCard Component (New)

**Location:** `frontend/src/components/dashboard/MetricCard.tsx`

**Props Interface:**
```typescript
interface MetricCardProps {
  isDark: boolean
  icon: React.ReactNode
  label: string
  value: string | number
  trend?: { value: string; isPositive: boolean }
  // FeatureCard-style gradient props
  gradientLight: string      // e.g., "from-blue-100 to-cyan-100"
  gradientDark: string       // e.g., "from-blue-500/20 to-cyan-400/20"
  iconColorLight: string     // e.g., "text-blue-600"
  iconColorDark: string      // e.g., "text-cyan-400"
  shadowLight?: string       // e.g., "shadow-md shadow-blue-500/15"
  borderLight?: string       // e.g., "border-2 border-blue-200"
  borderDark?: string        // e.g., "border border-cyan-500/20"
}
```

**Styling:**
- Container: `liquid-glass-content-card rounded-2xl overflow-hidden group animate-fade-in-up`
- Icon Container: 12x12 rounded-xl with gradient background
- Hover: Scale icon to 110% on hover

**Light Mode Colors:**
- Label: `text-slate-900`
- Value: `text-slate-900`
- Secondary text: `text-slate-600`
- Positive trend: `text-emerald-600`
- Negative trend: `text-red-600`

**Dark Mode Colors:**
- Label: `text-white`
- Value: `text-white`
- Secondary text: `text-slate-400`
- Positive trend: `text-emerald-400`
- Negative trend: `text-red-400`

---

### 2. DashboardPage Updates

**KPI Cards Configuration:**

| Card | Gradient Light | Gradient Dark | Icon Color Light | Icon Color Dark |
|------|---------------|---------------|------------------|-----------------|
| Revenue | `from-blue-100 to-cyan-100` | `from-blue-500/20 to-cyan-400/20` | `text-blue-600` | `text-cyan-400` |
| Patients | `from-emerald-100 to-teal-100` | `from-emerald-500/20 to-teal-400/20` | `text-emerald-600` | `text-teal-400` |
| Appointments | `from-purple-100 to-pink-100` | `from-purple-500/20 to-pink-400/20` | `text-purple-600` | `text-purple-400` |
| Inventory | `from-amber-100 to-orange-100` | `from-amber-500/20 to-orange-400/20` | `text-amber-600` | `text-amber-400` |

**Quick Actions:**
- Use `liquid-glass-button-prominent` class
- Add theme-aware hover states
- Icon size: 2xl
- Transition: all duration-200

**Footer:**
- Border: `isDark ? 'border-white/10' : 'border-slate-200'`
- Text: `liquid-text-secondary`

---

### 3. DashboardGrid Updates

**Props Update:**
```typescript
interface DashboardGridProps {
  onChartClick?: (chart: PinnedChart) => void
  className?: string
  isDark: boolean  // Add this
}
```

**Section Header:**
```tsx
<h2 className={`text-xl font-semibold ${isDark ? 'text-white' : 'text-slate-900'}`}>
  {t('title')}
</h2>
<p className={`text-sm mt-1 ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
  Your pinned business insights
</p>
```

**Empty State:**
- Container: `liquid-glass-content-card text-center py-16`
- Icon container: Gradient background matching theme
- Use `isDark` conditional styling

---

### 4. PinnedChartCard Updates

**Props Update:**
```typescript
interface PinnedChartCardProps {
  chart: PinnedChart
  locale: string
  isDark: boolean  // Add this
  onDelete: () => void
  onRefresh: () => void
  onToggle: (active: boolean) => void
  onClick?: () => void
}
```

**Border Updates:**
- Header border: `isDark ? 'border-white/10' : 'border-slate-200'`
- Footer border: `isDark ? 'border-white/10' : 'border-slate-200'`
- Footer background: `isDark ? 'bg-white/5' : 'bg-slate-50'`

---

## Implementation Order

1. Create `MetricCard.tsx` component
2. Update `DashboardPage.tsx`:
   - Add `isDark` prop
   - Replace KPI cards with `MetricCard`
   - Update Quick Actions styling
   - Update footer styling
3. Update `DashboardGrid.tsx`:
   - Add `isDark` prop
   - Update section headers
   - Update empty state
   - Pass `isDark` to PinnedChartCard
4. Update `PinnedChartCard.tsx`:
   - Add `isDark` prop
   - Update borders and backgrounds

---

## Testing Checklist

- [ ] KPI cards match FeatureCard styling in light mode
- [ ] KPI cards match FeatureCard styling in dark mode
- [ ] Hover effects work correctly
- [ ] Quick Actions buttons have correct theme
- [ ] Dashboard section headers are theme-aware
- [ ] PinnedChartCard borders adapt to theme
- [ ] Empty state displays correctly in both themes
- [ ] Footer adapts to theme
- [ ] RTL support maintained

---

## Design System References

- `liquid-glass-content-card` - iOS 26 style content cards
- `liquid-glass-button-prominent` - iOS 26 style buttons
- `liquid-text-primary` / `liquid-text-secondary` - Theme-aware text colors
- FeatureCard component pattern for gradient icon containers
