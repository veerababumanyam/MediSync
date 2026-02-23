# Hero Carousel Redesign Design Document

**Date:** 2026-02-23
**Status:** Approved
**Author:** Claude AI + User Collaboration

## Overview

Complete visual overhaul of the landing page Hero Carousel to increase demo bookings through a conversion-optimized narrative, eye-catching animations, and redesigned SVG illustrations.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Business Goal** | Increase demo bookings | Focus on conversion optimization with clear CTA progression |
| **Narrative Flow** | Progressive build | Hook → Capabilities → Proof → Demo creates trust before ask |
| **Visual Treatment** | Eye-catching animations | Showcase innovation through dynamic, modern design |
| **Dot Navigation** | Hidden (sr-only) | Cleaner look, progress bar indicates timing |
| **Approach** | Full Visual Overhaul | Complete SVG redesign + capability strip for maximum impact |

## Section 1: Copy and Narrative Structure

### Slide 1 - Hook + Problem

**Role:** Create "aha" moment, stop the scroll

| Element | Value |
|---------|-------|
| **Title** | "Don't Replace Your Software. Make It Speak." |
| **Subtitle** | "Any System. Any Database. Zero Migration." |
| **Description** | "Ask questions in plain language. Get instant answers from HIMS, Tally, SQL, or any legacy system. No rip-and-replace required." |
| **CTA** | "Start Free Trial" |
| **Stat 1** | "50+ Integrations" |
| **Stat 2** | "< 2 Min Setup" |
| **Stat 3** | "Zero Code Changes" |

### Slide 2 - Capabilities + How It Works

**Role:** Show full platform breadth in one glance

| Element | Value |
|---------|-------|
| **Title** | "One Platform. Every Capability." |
| **Subtitle** | "Conversational BI · AI Accountant · Smart Reports · Deep Analytics" |
| **Description** | "From natural language queries to automated accounting sync — all powered by 58 specialized AI agents working together." |
| **CTA** | "See It In Action" |
| **Stat 1** | "58 AI Agents" |
| **Stat 2** | "NL → SQL → Charts" |
| **Stat 3** | "Auto Ledger Mapping" |

**New Element:** Capability strip with 4 icons below description

### Slide 3 - Proof + Demo CTA

**Role:** Build trust, drive demo booking

| Element | Value |
|---------|-------|
| **Title** | "Built for Healthcare & Finance." |
| **Subtitle** | "HIPAA Compliant · SOC 2 Certified · 99.9% Uptime" |
| **Description** | "Trusted by clinics, labs, and hospitals across 12 countries. Average savings: ₹2Cr+ annually with full compliance." |
| **CTA** | "Book a Demo" |
| **Stat 1** | "₹2Cr+ Avg Savings" |
| **Stat 2** | "12 Countries" |
| **Stat 3** | "500+ Clinics" |

## Section 2: Visual Layout and Hierarchy

### Inner Card Treatment

Each slide becomes a "hero card" within the liquid-glass-pronounced container:

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo] AnySync                                             │  ← Logo moves to top
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────────────────────┐  ┌────────────────────┐   │
│   │  HEADLINE                   │  │                    │   │
│   │  (gradient text)            │  │    ILLUSTRATION    │   │
│   │                             │  │    (animated SVG)  │   │
│   │  Subtitle                   │  │                    │   │
│   │                             │  │    [subtle glow]   │   │
│   │  Description paragraph      │  │                    │   │
│   │                             │  └────────────────────┘   │
│   │  [CAPABILITY STRIP*]        │                           │
│   │                             │                           │
│   │  ┌────────────────────┐     │                           │
│   │  │  PRIMARY CTA       │     │                           │
│   │  └────────────────────┘     │                           │
│   │                             │                           │
│   │  ★ 4.9/5  •  Stat1  •  Stat2  •  Stat3                │
│   └─────────────────────────────┘                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Layout Changes

| Element | Current | New |
|---------|---------|-----|
| **Logo position** | Inline with badge row | Top-left of card, more prominent |
| **Inner card** | None | Optional inset shadow/border using existing tokens |
| **Illustration backdrop** | None | Subtle radial gradient (var(--brand-secondary) at 10% opacity) |
| **Capability strip** | None (slide 2 only) | 4 icons with labels: Chat, Document, Chart, Sparkle |
| **CTA prominence** | Standard | Slightly larger with enhanced glow on hover |
| **Typography scale** | text-3xl → xl:text-5xl | Bump to xl:text-5xl/6xl for more impact |

### Responsive Behavior

| Breakpoint | Layout |
|------------|--------|
| **Mobile (sm)** | Stacked - text above, illustration below, CTA centered |
| **Tablet (md-lg)** | Side-by-side, illustration at 280px max |
| **Desktop (xl+)** | Side-by-side, illustration at 340px max, larger typography |

## Section 3: Illustration Redesign (3 New SVGs)

### Slide 1 - "Conversation with Legacy" (Hook)

**Concept:** AI brain/nodes interacting with legacy database icons, showing "speech" bridging old and new

**Visual Elements:**
- Glow ring (animated pulse)
- Central AI node (accent color, float animation)
- Legacy system nodes: HIMS, Tally, SQL (neutral, floatReverse)
- Chat bubbles flowing between AI and legacy
- Sparkle particles on connection lines

**Animations:**
- Central AI node: `floatY` 4s
- Legacy nodes: `floatYReverse` 4.5s with staggered delays
- Connection lines: Dash animation (data flowing)
- Sparkle particles: `pulseGlow` 2-3s
- Chat bubbles: `floatY` alternating

### Slide 2 - "Platform Capabilities" (Capabilities)

**Concept:** Central hub with 4 capability nodes radiating out

**Visual Elements:**
- Central hub with "AnySync" text
- 4 capability nodes with icons:
  - Chat/BI: Speech bubble with spark
  - Accountant: Document with check
  - Reports: Chart bars
  - Analytics: Brain/network
- Connection lines radiating from hub
- Data particles flowing from hub to nodes

**Animations:**
- Central hub: `pulseGlow` 5s
- Capability nodes: `floatY` 3-4s staggered
- Connection lines: Radial dash animation
- Data particles: Flowing from hub to nodes

### Slide 3 - "Trust Dashboard" (Proof)

**Concept:** Premium analytics dashboard with trust badges and outcome metrics

**Visual Elements:**
- Window chrome with traffic lights
- Floating metric cards: +127%, ₹2Cr, 500+
- Growing bar chart
- Trend line overlay
- Trust badges: HIPAA, 99.9%, SOC2

**Animations:**
- Dashboard frame: `floatY` 5s
- Metric cards: `floatYReverse` 3-4s staggered
- Bar chart: Sequential grow animation (CSS)
- Trust badges: Subtle `pulseGlow` 4s

### SVG Technical Specifications

| Property | Value |
|----------|-------|
| viewBox | `"0 0 320 280"` |
| Max width | `max-w-[340px]` |
| Colors | Design tokens (`var(--brand-primary)`, etc.) |
| Animations | CSS keyframes (`floatY`, `floatYReverse`, `pulseGlow`) |
| Reduced motion | Respects `prefers-reduced-motion` |
| Accessibility | `aria-hidden="true"` |

## Section 4: Animations and Interactions

### Animation Palette

| Animation | Duration | Use Case |
|-----------|----------|----------|
| `floatY` | 3-5s | Upward float, use on primary elements |
| `floatYReverse` | 3-5s | Downward float, use on secondary elements |
| `pulseGlow` | 3-5s | Glow pulse, use on accent nodes |
| `heroSlideIn` | 0.7s | Slide entrance (existing) |
| `fade-in-up` | 0.5s | Text entrance with delay |
| `fade-in-down` | 0.5s | Logo entrance |

### Entrance Sequence (Per Slide)

| Time | Element | Animation |
|------|---------|-----------|
| 0ms | Logo row | fade-in-down |
| 100ms | Headline | fade-in-up |
| 200ms | Subtitle | fade-in-up |
| 300ms | Description | fade-in-up |
| 400ms | Capability strip (slide 2) | fade-in-up |
| 500ms | CTA button | fade-in-up |
| 600ms | Stats + rating badge | fade-in-up |
| 700ms | Illustration | heroSlideIn + floatY |

### Hover Interactions

| Element | Hover Effect |
|---------|-------------|
| **CTA button** | Lift (-4px), enhanced glow, shimmer |
| **Capability icons** | Subtle scale (1.05), glow accent |
| **Trust badges** | Subtle glow pulse |
| **Entire carousel** | Pause auto-rotate |

### Reduced Motion Compliance

```css
@media (prefers-reduced-motion: reduce) {
  .hero-carousel-container *,
  .hero-carousel-container *::before,
  .hero-carousel-container *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

### Performance Considerations

- **will-change**: Applied to animated elements only (`transform`, `opacity`)
- **CSS-only**: No JavaScript-driven animations
- **GPU acceleration**: Transforms use `translate3d()` where supported
- **Animation budget**: Max 5 concurrent animations per slide

## Section 5: Implementation Summary

### Files to Modify

| File | Changes |
|------|---------|
| `frontend/src/i18n/index.ts` | Add new copy keys (EN), add AR translations |
| `frontend/src/components/landing/HeroCarousel.tsx` | New layout, capability strip, inner card |
| `frontend/src/components/landing/icons.tsx` | 3 new HeroIllustration SVGs |
| `frontend/src/styles/globals.css` | Add capability strip styles, reduced motion |

### Implementation Order

1. **Copy & i18n** - Add new keys to i18n/index.ts (EN + AR)
2. **Illustrations** - Create 3 new SVGs in icons.tsx
3. **HeroCarousel** - Update layout, add capability strip, inner card
4. **Styles** - Add any new CSS classes needed
5. **Testing** - Verify animations, reduced motion, RTL, accessibility

### i18n Keys Structure

```typescript
heroCarousel: {
  slide1: {
    title: string
    subtitle: string  // NEW
    description: string
    cta: string
    stat1: string
    stat2: string
    stat3: string
  }
  slide2: {
    title: string
    subtitle: string  // NEW
    description: string
    cta: string
    stat1: string
    stat2: string
    stat3: string
  }
  slide3: {
    title: string
    subtitle: string  // NEW
    description: string
    cta: string
    stat1: string
    stat2: string
    stat3: string
  }
  capabilities: {     // NEW
    bi: string
    accountant: string
    reports: string
    analytics: string
  }
  trustStrip: { ... }  // existing
}
```

## Accessibility & Compliance

| Requirement | Implementation |
|-------------|----------------|
| **Touch targets** | `min-height: 44px` (WCAG 2.5.5) |
| **Focus indicators** | `liquid-focus-gold` with 2px gold outline |
| **Color contrast** | APCA-compliant text colors |
| **Reduced motion** | Media query disables all animations |
| **High contrast mode** | Inherits from existing system |
| **RTL support** | Logical properties, i18n keys |
| **Keyboard navigation** | sr-only dots, tab to CTA |
| **Screen readers** | aria-live, aria-roledescription |

## Design Tokens Used

| Token | Purpose |
|-------|---------|
| `--brand-primary` | Logo blue (#2750a8) |
| `--brand-secondary` | Logo teal (#18929d) |
| `--gradient-brand` | Pre-defined brand gradient |
| `--shadow-glow-secondary` | Teal glow effect |
| `--focus-gold` | WCAG 3.0 Gold focus color |
| `--radius-xl` | 20px border radius |
| `--color-text-primary` | Primary text |
| `--color-text-secondary` | Secondary text |

## Success Metrics

| Metric | Target |
|--------|--------|
| Demo click-through rate | +20% from current baseline |
| Time on carousel | 12+ seconds (2 full rotations) |
| Slide 3 CTA clicks | 40%+ of total carousel CTA clicks |
| Reduced motion users | No animation-related complaints |

## Approval

- [x] Design approved by user
- [ ] Implementation completed
- [ ] Accessibility verified
- [ ] Cross-browser testing completed
- [ ] RTL testing completed

---

*Generated by Claude AI with user collaboration on 2026-02-23*
