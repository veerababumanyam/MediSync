# AnySync Design System v3.0.0

**iOS 26 Liquid Glass · WCAG 3.0 Bronze · Light + Dark Modes**

Version 3.0.0 · February 23, 2026

---

## Table of Contents

1. [Logo Analysis & Color Extraction](#1-logo-analysis--color-extraction)
2. [Brand Color System](#2-brand-color-system)
3. [Extended Palette & Neutral Scale](#3-extended-palette--neutral-scale)
4. [Gradient Tokens](#4-gradient-tokens)
5. [Semantic Colors](#5-semantic-colors)
6. [Surface & Text Colors](#6-surface--text-colors)
7. [WCAG 3.0 Contrast Compliance](#7-wcag-30-contrast-compliance)
8. [iOS 26 Liquid Glass System](#8-ios-26-liquid-glass-system)
9. [Shadow & Elevation System](#9-shadow--elevation-system)
10. [Spacing & Layout](#10-spacing--layout)
11. [Border Radius & Corners](#11-border-radius--corners)
12. [Typography System](#12-typography-system)
13. [Animation & Motion](#13-animation--motion)
14. [Z-Index Architecture](#14-z-index-architecture)
15. [Component Tokens](#15-component-tokens)
16. [Accessibility Standards](#16-accessibility-standards)
17. [File Architecture & Usage](#17-file-architecture--usage)

---

## 1. Logo Analysis & Color Extraction

### 1.1 Visual Structure

The AnySync logo is a composite symbol merging two visual metaphors into a single cohesive mark:

```
 ▌▌▌▌∿∿→
 ^^^^  ^^^^
 BARS  PULSE
```

**Left side — Bar chart (Data/Analytics):** Four vertical bars in descending height. Colors transition from deep royal blue through teal to cyan. Represents structure, measurement, and intelligence.

**Right side — ECG heartbeat line with arrow (Vitality):** An ECG-style waveform with sharp peak and forward-pointing arrow. Colors shift from cyan through emerald to bright green. Represents life, health, real-time monitoring, and forward momentum.

**Background — Pure black (#000000):** Places this brand firmly in a premium, tech-forward, medical-grade context.

### 1.2 Color Extraction Points

Colors sampled at seven key points along the logo's gradient path:

| Element | Hex | Color Name | Role |
|---------|-----|-----------|------|
| Bar 1 (tallest, left) | `#1B3FA0` | Royal Blue | Primary brand |
| Bar 2 | `#1A5DAE` | Mid Blue | Transition |
| Bar 3 | `#1A7A9E` | Blue-Teal | Bridge |
| Bar 4 (shortest) | `#17B5A6` | Teal-Cyan | Secondary brand |
| Heartbeat start | `#19CDA0` | Cyan-Green | Transition |
| Heartbeat peak | `#1CD760` | Emerald | Accent / Action |
| Arrow tip | `#2AE668` | Bright Green | Vitality CTA |

### 1.3 Simplified 3-Color System

The seven extraction points consolidate into a practical three-color brand system:

- **Primary `#1B3FA0` (Royal Blue)** — Authority, structure, data intelligence. Used for headings, primary buttons, links, and brand anchoring.
- **Secondary `#17B5A6` (Teal-Cyan)** — Flow, transition, connectivity. Used for secondary actions, accents, gradients, and interactive highlights.
- **Accent `#1CD760` (Emerald)** — Vitality, pulse, action. Used for success states, CTAs, status indicators, and the heartbeat-pulse animation motif.

### 1.4 Color Narrative

The gradient flows left-to-right telling the brand story: **structured data → flowing transition → living pulse → forward momentum.** This narrative informs every gradient, animation, and visual hierarchy decision in the system.

---

## 2. Brand Color System

Each brand color has a full 11-step scale (50–950) for maximum flexibility. The **500** step is always the logo-extracted exact color.

### 2.1 Primary — Royal Blue

| Step | Hex | RGB | Usage |
|------|-----|-----|-------|
| 50 | `#EBF0FF` | 235, 240, 255 | Subtle backgrounds, hover states |
| 100 | `#D6E0FF` | 214, 224, 255 | Light backgrounds, borders |
| 200 | `#ADC1FF` | 173, 193, 255 | Disabled states |
| 300 | `#7596F5` | 117, 150, 245 | Dark mode `--text-brand` |
| 400 | `#4A6FDB` | 74, 111, 219 | Dark mode buttons |
| **500** | **`#1B3FA0`** | **27, 63, 160** | **Logo exact — Primary brand** |
| 600 | `#163489` | 22, 52, 137 | Hover states |
| 700 | `#112972` | 17, 41, 114 | Active / pressed states |
| 800 | `#0D1F5B` | 13, 31, 91 | Deep emphasis |
| 900 | `#091544` | 9, 21, 68 | Ultra-deep |
| 950 | `#050C2D` | 5, 12, 45 | Near-black brand |

**CSS:** `--brand-primary` through `--brand-primary-950`

### 2.2 Secondary — Teal-Cyan

| Step | Hex | RGB | Usage |
|------|-----|-----|-------|
| 50 | `#ECFDF9` | 236, 253, 249 | Success-tinted backgrounds |
| 100 | `#D1FAF0` | 209, 250, 240 | Light teal fills |
| 200 | `#A5F3E2` | 165, 243, 226 | Light teal accents |
| 300 | `#6DE8D0` | 109, 232, 208 | Highlights |
| 400 | `#33D4BB` | 51, 212, 187 | Dark mode secondary |
| **500** | **`#17B5A6`** | **23, 181, 166** | **Logo exact — Secondary brand** |
| 600 | `#0F9185` | 15, 145, 133 | Hover |
| 700 | `#0B7068` | 11, 112, 104 | Deep teal |
| 800 | `#08544E` | 8, 84, 78 | Very deep |
| 900 | `#053B37` | 5, 59, 55 | Ultra-deep |
| 950 | `#022421` | 2, 36, 33 | Near-black teal |

**CSS:** `--brand-secondary` through `--brand-secondary-950`

### 2.3 Accent — Emerald Green (Heartbeat Pulse)

| Step | Hex | RGB | Usage |
|------|-----|-----|-------|
| 50 | `#EDFFF3` | 237, 255, 243 | Success backgrounds |
| 100 | `#D5FFE4` | 213, 255, 228 | Light green fills |
| 200 | `#ADFFC9` | 173, 255, 201 | Light green accents |
| 300 | `#72F59E` | 114, 245, 158 | Highlights |
| 400 | `#3DE878` | 61, 232, 120 | Dark mode `--text-accent` |
| **500** | **`#1CD760`** | **28, 215, 96** | **Logo heartbeat — Accent** |
| 600 | `#14B04E` | 20, 176, 78 | Hover |
| 700 | `#0F8A3D` | 15, 138, 61 | Light mode `--text-accent` (WCAG safe) |
| 800 | `#0B6A2F` | 11, 106, 47 | Deep green |
| 900 | `#074D22` | 7, 77, 34 | Ultra-deep |
| 950 | `#042F15` | 4, 47, 21 | Near-black green |

**CSS:** `--brand-accent` through `--brand-accent-950`

### 2.4 Deep Surface (Logo Background)

| Token | Hex | Usage |
|-------|-----|-------|
| `--brand-deep` | `#0A0F1A` | Dark mode primary surface |
| `--brand-navy` | `#0F172A` | Dark mode secondary surface |

---

## 3. Extended Palette & Neutral Scale

Neutrals are **tinted toward the blue brand family** for visual cohesion. Pure grays look disconnected; blue-tinted slates feel integrated and intentional.

| Step | Hex | RGB | Usage |
|------|-----|-----|-------|
| 0 | `#FFFFFF` | 255, 255, 255 | Pure white — cards, inputs |
| 25 | `#FAFBFD` | 250, 251, 253 | Near-white backgrounds |
| 50 | `#F4F6FA` | 244, 246, 250 | Page background (light mode) |
| 100 | `#E8ECF3` | 232, 236, 243 | Dividers, subtle borders |
| 200 | `#D1D8E5` | 209, 216, 229 | Default borders |
| 300 | `#B0BBCE` | 176, 187, 206 | Strong borders, disabled text |
| 400 | `#8694AE` | 134, 148, 174 | Placeholder text |
| 500 | `#64748B` | 100, 116, 139 | Tertiary text |
| 600 | `#475569` | 71, 85, 105 | Secondary text |
| 700 | `#334155` | 51, 65, 85 | Strong secondary text |
| 800 | `#1E293B` | 30, 41, 59 | Dark surfaces |
| 900 | `#0F172A` | 15, 23, 42 | Primary text (light mode) |
| 950 | `#080D19` | 8, 13, 25 | Ultra-dark background |

**CSS:** `--neutral-0` through `--neutral-950`

---

## 4. Gradient Tokens

Gradients mirror the logo's left-to-right color narrative. The full spectrum gradient recreates the exact logo flow across all seven extraction points.

| Token | Value | Usage |
|-------|-------|-------|
| `--gradient-brand` | `135deg, Primary → Secondary` | Default brand gradient for CTAs, headers |
| `--gradient-brand-extended` | `135deg, Primary → Secondary → Accent` | Full 3-color brand expression |
| `--gradient-brand-reverse` | `135deg, Secondary → Primary` | Reverse for variety |
| `--gradient-logo-spectrum` | `90deg, #1B3FA0 → #1A7A9E → #17B5A6 → #1CD760 → #2AE668` | Hero backgrounds, marketing materials |
| `--gradient-pulse` | `90deg, Secondary → Accent` | Heartbeat-line inspired elements |
| `--gradient-pulse-glow` | `90deg, Secondary 40% → Accent 40%` | Subtle pulse glow behind elements |
| `--gradient-mesh-light` | 3 radial-gradients overlaid on white | Hero section ambient lighting |
| `--gradient-mesh-dark` | 3 radial-gradients on `#0A0F1A` | Dark mode hero ambient |
| `--gradient-glass-blue` | Brand-tinted glass overlay | iOS 26 colored glass panels |
| `--gradient-glass-teal` | Teal-tinted glass overlay | Secondary glass panels |
| `--gradient-glass-accent` | Green-tinted glass overlay | Accent glass panels |
| `--gradient-shine` | `transparent → white 40% → transparent` | Shimmer / loading animation |
| `--gradient-surface-subtle` | `180deg, surface-primary → surface-secondary` | Subtle page gradient |

### Mesh Background Example

```css
/* Hero section with brand-colored ambient lighting */
.hero {
  background: var(--gradient-mesh-light);
  /* Resolves to:
     radial-gradient(ellipse at 20% 50%, rgba(27,63,160, 0.08), transparent 60%),
     radial-gradient(ellipse at 80% 20%, rgba(23,181,166, 0.06), transparent 50%),
     radial-gradient(ellipse at 60% 80%, rgba(28,215,96, 0.05), transparent 50%),
     #FFFFFF;
  */
}
```

---

## 5. Semantic Colors

Status colors are visually distinct from brand colors while harmonizing with the palette. Success aligns with the accent green family for coherence.

| Status | Default | Light | Dark | Background | Border |
|--------|---------|-------|------|-----------|--------|
| Success | `#10B981` | `#34D399` | `#059669` | `rgba(16,185,129, 0.08)` | `rgba(16,185,129, 0.25)` |
| Warning | `#F59E0B` | `#FBBF24` | `#D97706` | `rgba(245,158,11, 0.08)` | `rgba(245,158,11, 0.25)` |
| Error | `#EF4444` | `#F87171` | `#DC2626` | `rgba(239,68,68, 0.08)` | `rgba(239,68,68, 0.25)` |
| Info | `#3B82F6` | `#60A5FA` | `#2563EB` | `rgba(59,130,246, 0.08)` | `rgba(59,130,246, 0.25)` |

### UI State Colors

| Token | Value | Usage |
|-------|-------|-------|
| `--color-hover` | `rgba(27,63,160, 0.06)` | Hover background tint |
| `--color-active` | `rgba(27,63,160, 0.12)` | Active / pressed |
| `--color-selected` | `rgba(27,63,160, 0.08)` | Selected row / item |
| `--color-disabled` | `rgba(0,0,0, 0.26)` | Disabled foreground |
| `--color-disabled-bg` | `rgba(0,0,0, 0.06)` | Disabled background |

---

## 6. Surface & Text Colors

### 6.1 Light Mode

| Token | Value | Purpose |
|-------|-------|---------|
| `--surface-primary` | `#FFFFFF` | Cards, inputs, main content |
| `--surface-secondary` | `#F4F6FA` | Page background |
| `--surface-tertiary` | `#E8ECF3` | Inset panels, wells |
| `--surface-elevated` | `#FFFFFF` | Elevated cards, modals |
| `--surface-sunken` | `#EDF0F7` | Recessed areas |
| `--surface-overlay` | `rgba(15,23,42, 0.5)` | Modal scrim |
| `--text-primary` | `#0F172A` | Headings, body text |
| `--text-secondary` | `#475569` | Descriptions, labels |
| `--text-tertiary` | `#64748B` | Captions, metadata |
| `--text-disabled` | `#B0BBCE` | Disabled text |
| `--text-brand` | `#1B3FA0` | Brand-colored text |
| `--text-accent` | `#0F8A3D` | Accent text (WCAG-safe dark green) |
| `--text-link` | `#1B3FA0` | Hyperlinks |
| `--text-link-hover` | `#112972` | Hovered links |

### 6.2 Dark Mode

| Token | Value | Purpose |
|-------|-------|---------|
| `--surface-primary` | `#0A0F1A` | Main background (logo deep) |
| `--surface-secondary` | `#111827` | Card backgrounds |
| `--surface-tertiary` | `#1E293B` | Inset panels |
| `--surface-elevated` | `#1A2236` | Elevated surfaces |
| `--surface-sunken` | `#070B14` | Recessed areas |
| `--surface-overlay` | `rgba(0,0,0, 0.6)` | Modal scrim |
| `--text-primary` | `#F1F5F9` | Headings, body text |
| `--text-secondary` | `#94A3B8` | Descriptions, labels |
| `--text-tertiary` | `#64748B` | Captions, metadata |
| `--text-disabled` | `#475569` | Disabled text |
| `--text-brand` | `#7596F5` | Brand text (Primary-300) |
| `--text-accent` | `#3DE878` | Accent text (vibrant green) |
| `--text-link` | `#7596F5` | Hyperlinks |
| `--text-link-hover` | `#ADC1FF` | Hovered links |

### 6.3 Border Tokens

| Token | Light | Dark |
|-------|-------|------|
| `--border-default` | `#D1D8E5` | `rgba(255,255,255, 0.10)` |
| `--border-subtle` | `#E8ECF3` | `rgba(255,255,255, 0.06)` |
| `--border-strong` | `#B0BBCE` | `rgba(255,255,255, 0.18)` |
| `--border-brand` | `var(--brand-primary)` | `var(--brand-primary-400)` |

---

## 7. WCAG 3.0 Contrast Compliance

All text/surface combinations verified against WCAG 3.0 APCA (Accessible Perceptual Contrast Algorithm) Bronze requirements:

- **Body text (<24px):** ≥ 60 Lc (lightness contrast)
- **Large text (24px+):** ≥ 45 Lc
- **Non-text (icons, borders):** ≥ 45 Lc
- **Focus indicators:** ≥ 3:1 against adjacent colors

### 7.1 Light Mode Contrast Matrix

| Color Pair | APCA Lc | Status |
|-----------|---------|--------|
| `--text-primary` (#0F172A) on white | ~106 Lc | ✅ Exceeds all requirements |
| `--text-secondary` (#475569) on white | ~68 Lc | ✅ Body text compliant |
| `--text-tertiary` (#64748B) on white | ~55 Lc | ✅ Large text / non-text |
| `--brand-primary` (#1B3FA0) on white | ~82 Lc | ✅ All text sizes |
| `--brand-accent-dark` (#0F8A3D) on white | ~60 Lc | ✅ Body text minimum |
| White on `--brand-primary` (#1B3FA0) | ~82 Lc | ✅ Button text |

### 7.2 Dark Mode Contrast Matrix

| Color Pair | APCA Lc | Status |
|-----------|---------|--------|
| `--text-primary` (#F1F5F9) on #0A0F1A | ~107 Lc | ✅ Exceeds all requirements |
| `--text-secondary` (#94A3B8) on #0A0F1A | ~62 Lc | ✅ Body text compliant |
| `--text-brand` (#7596F5) on #0A0F1A | ~68 Lc | ✅ All text sizes |
| `--text-accent` (#3DE878) on #0A0F1A | ~72 Lc | ✅ All text sizes |
| White on `--brand-primary-400` (#4A6FDB) | ~68 Lc | ✅ Dark mode buttons |

---

## 8. iOS 26 Liquid Glass System

The glassmorphism system follows Apple's iOS 26 "Liquid Glass" design philosophy. Every glass panel simulates physical glass with translucency, frosted blur, light refraction, and specular highlights.

### 8.1 Design Principles

1. **Translucency** — Background content bleeds through at calibrated opacity levels
2. **Backdrop Blur** — Frosted glass effect via `backdrop-filter`, ranging from 8px to 100px
3. **Light Refraction** — Borders simulate light bending through glass with varying opacity
4. **Specular Highlights** — Top-edge shimmer simulates overhead light source (the "glass edge" effect)
5. **Brand Tinting** — Glass absorbs underlying brand color for visual cohesion (unique to iOS 26)
6. **Saturation Boost** — Blurred content is saturated 1.4–1.8× for vibrancy
7. **Continuous Curvature** — Rounded corners use iOS-style "squircle" (superellipse) approximation

### 8.2 Glass Layer Hierarchy

Five opacity tiers from most transparent to most opaque:

| Token | Light Mode | Dark Mode | Usage |
|-------|-----------|-----------|-------|
| `--glass-bg-subtle` | `rgba(255,255,255, 0.35)` | `rgba(15,23,42, 0.55)` | Background hints |
| `--glass-bg-light` | `rgba(255,255,255, 0.60)` | `rgba(15,23,42, 0.70)` | Subtle panels |
| `--glass-bg-medium` | `rgba(255,255,255, 0.72)` | `rgba(15,23,42, 0.80)` | Standard panels |
| `--glass-bg-heavy` | `rgba(255,255,255, 0.82)` | `rgba(15,23,42, 0.88)` | Cards, nav bars |
| `--glass-bg-pronounced` | `rgba(255,255,255, 0.92)` | `rgba(15,23,42, 0.94)` | Modals, sheets |

### 8.3 Brand-Tinted Glass (iOS 26 Exclusive)

Glass panels absorb the color of underlying content. These tokens add a subtle brand tint:

| Token | Light Mode | Dark Mode |
|-------|-----------|-----------|
| `--glass-bg-brand-blue` | `rgba(27,63,160, 0.06)` | `rgba(27,63,160, 0.12)` |
| `--glass-bg-brand-teal` | `rgba(23,181,166, 0.06)` | `rgba(23,181,166, 0.10)` |
| `--glass-bg-brand-green` | `rgba(28,215,96, 0.06)` | `rgba(28,215,96, 0.08)` |

### 8.4 Glass Borders

| Token | Light Mode | Dark Mode | Purpose |
|-------|-----------|-----------|---------|
| `--glass-border-bright` | `rgba(255,255,255, 0.65)` | `rgba(255,255,255, 0.14)` | Primary glass edge |
| `--glass-border-soft` | `rgba(255,255,255, 0.40)` | `rgba(255,255,255, 0.08)` | Secondary edge |
| `--glass-border-dim` | `rgba(255,255,255, 0.18)` | `rgba(255,255,255, 0.04)` | Subtle dividers |
| `--glass-border-tinted` | `rgba(27,63,160, 0.12)` | `rgba(113,150,245, 0.15)` | Brand-tinted |

### 8.5 Specular Highlights

| Token | Light | Dark |
|-------|-------|------|
| `--glass-highlight` | `rgba(255,255,255, 0.85)` | `rgba(255,255,255, 0.12)` |
| `--glass-highlight-soft` | `rgba(255,255,255, 0.45)` | `rgba(255,255,255, 0.06)` |
| `--glass-highlight-tint` | `rgba(23,181,166, 0.15)` | `rgba(23,181,166, 0.10)` |

### 8.6 Backdrop Blur Scale

| Token | Value | Usage |
|-------|-------|-------|
| `--blur-xs` | `blur(8px)` | Subtle depth hint |
| `--blur-sm` | `blur(16px)` | Cards, subtle panels |
| `--blur-md` | `blur(32px)` | Standard glass panels |
| `--blur-lg` | `blur(48px)` | Navigation bars, heavy glass |
| `--blur-xl` | `blur(64px)` | Modals, overlays |
| `--blur-xxl` | `blur(80px)` | Full-screen frosted |
| `--blur-max` | `blur(100px)` | Maximum frost effect |

### 8.7 Glass Utility Classes

```css
/* Ready-to-use glass panels */
.glass-panel         /* Standard: medium bg, md blur, soft border, 24px radius */
.glass-panel-subtle  /* Lightweight: subtle bg, sm blur, dim border */
.glass-panel-heavy   /* Prominent: heavy bg, lg blur, bright border */
.glass-card          /* Card with specular highlight ::before pseudo-element */
.glass-card-gradient /* Card with gradient border (brand-tinted) */
.glass-nav           /* Fixed navigation bar with heavy glass */
```

### 8.8 Composing a Glass Panel

```css
.my-component {
  background: var(--glass-bg-medium);
  backdrop-filter: var(--blur-md) var(--glass-saturation);
  -webkit-backdrop-filter: var(--blur-md) var(--glass-saturation);
  border: 1px solid var(--glass-border-soft);
  border-radius: var(--radius-2xl);
  box-shadow: var(--shadow-raised), var(--glass-inner-glow);
}
```

---

## 9. Shadow & Elevation System

A five-tier elevation model where each level adds more shadow depth. Every shadow uses three layers: a ring (border simulation), a key shadow (direct light), and an ambient shadow (diffused fill).

### 9.1 Elevation Levels

| Token | Level | Usage | Light Mode |
|-------|-------|-------|-----------|
| `--shadow-ambient` | 0 | Resting cards, default state | Subtle 3-layer shadow |
| `--shadow-raised` | 1 | Hovered cards, buttons | Medium 3-layer shadow |
| `--shadow-floating` | 2 | Dropdowns, popovers | Strong 3-layer shadow |
| `--shadow-elevated` | 3 | Modals, dialogs | White ring + deep shadow |
| `--shadow-glass-floating` | 4 | Hero glass panels | White ring + maximum depth |

### 9.2 Inner Shadows

| Token | Usage |
|-------|-------|
| `--shadow-inset` | Pressed / inset button state |
| `--glass-inner-glow` | `inset 0 1px 0 0 rgba(255,255,255, 0.5)` — Top-edge luminosity |
| `--glass-inner-glow-soft` | `inset 0 1px 0 0 rgba(255,255,255, 0.25)` — Subtle version |

### 9.3 Brand Glow Shadows

| Token | Color | Usage |
|-------|-------|-------|
| `--shadow-glow-primary` | Royal Blue | Primary action emphasis |
| `--shadow-glow-secondary` | Teal | Secondary highlights |
| `--shadow-glow-accent` | Emerald | Success / pulse indicators |
| `--shadow-glow-success` | Green | Success state glow |
| `--shadow-glow-warning` | Amber | Warning state glow |
| `--shadow-glow-error` | Red | Error state glow |
| `--shadow-pulse` | Animated emerald ring | Heartbeat-line indicator |

Dark mode versions have **increased glow intensity** (30–40px spread vs 24px) for visibility against dark surfaces.

---

## 10. Spacing & Layout

A **4px base unit system** provides consistent spatial rhythm. All spacing values are multiples of 4px (0.25rem).

### 10.1 Spacing Scale

| Token | rem | Pixels | Usage |
|-------|-----|--------|-------|
| `--space-0` | 0 | 0 | Reset |
| `--space-px` | — | 1px | Hairline |
| `--space-0_5` | 0.125 | 2px | Micro gaps |
| `--space-1` | 0.25 | 4px | Tight inline gaps |
| `--space-1_5` | 0.375 | 6px | Small inline |
| `--space-2` | 0.5 | 8px | Icon gaps, compact padding |
| `--space-3` | 0.75 | 12px | Input padding, small gaps |
| `--space-4` | 1.0 | 16px | Standard padding, card gaps |
| `--space-5` | 1.25 | 20px | Medium gaps |
| `--space-6` | 1.5 | 24px | Section gaps, card padding |
| `--space-8` | 2.0 | 32px | Large gaps |
| `--space-10` | 2.5 | 40px | Extra large gaps |
| `--space-12` | 3.0 | 48px | Section padding (small) |
| `--space-16` | 4.0 | 64px | Section padding (medium) |
| `--space-20` | 5.0 | 80px | Section padding (large) |
| `--space-24` | 6.0 | 96px | Hero spacing |
| `--space-32` | 8.0 | 128px | Maximum section spacing |
| `--space-40` | 10.0 | 160px | Full-bleed hero |

### 10.2 Section Spacing

| Token | Value | Usage |
|-------|-------|-------|
| `--section-padding-sm` | 48px | Compact sections |
| `--section-padding-md` | 80px | Standard sections |
| `--section-padding-lg` | 112px | Prominent sections |
| `--section-padding-xl` | 144px | Hero sections |

### 10.3 Gap System

| Token | Value | Usage |
|-------|-------|-------|
| `--gap-xs` | 8px | Tight grid gaps |
| `--gap-sm` | 16px | Card grid gaps |
| `--gap-md` | 24px | Standard grid gaps |
| `--gap-lg` | 32px | Large grid gaps |
| `--gap-xl` | 48px | Section grid gaps |

### 10.4 Layout Tokens

| Token | Value | Usage |
|-------|-------|-------|
| `--main-content-top` | 24px | Space between header and main |
| `--banner-to-content` | 24px | Space below announcement banner |

---

## 11. Border Radius & Corners

iOS uses "squircle" continuous curvature corners (superellipse), which appear smoother than standard CSS `border-radius`. We approximate with generous radius values and assign semantic roles.

### 11.1 Radius Scale

| Token | Value | Component |
|-------|-------|-----------|
| `--radius-none` | 0 | Sharp corners |
| `--radius-xs` | 4px | Inline badges, tags |
| `--radius-sm` | 8px | Tooltips, small elements |
| `--radius-md` | 12px | Inputs, compact cards |
| `--radius-lg` | 16px | Buttons |
| `--radius-xl` | 20px | Standard cards |
| `--radius-2xl` | 24px | Large cards |
| `--radius-3xl` | 28px | Modals, sheets |
| `--radius-4xl` | 32px | Hero panels |
| `--radius-dynamic` | 40px | iOS Dynamic Island inspired |
| `--radius-pill` | 9999px | Pill buttons, capsule badges |
| `--radius-full` | 9999px | Circles (avatars) |

### 11.2 Component Assignments

| Component | Token | Value |
|-----------|-------|-------|
| Button | `--radius-button` | 16px |
| Card | `--radius-card` | 24px |
| Modal | `--radius-modal` | 28px |
| Input | `--radius-input` | 12px |
| Badge | `--radius-badge` | 9999px (pill) |
| Avatar | `--radius-avatar` | 9999px (circle) |
| Tooltip | `--radius-tooltip` | 8px |

### 11.3 Border Widths

| Token | Value | Usage |
|-------|-------|-------|
| `--border-thin` | 1px | Default borders |
| `--border-medium` | 1.5px | Emphasized borders |
| `--border-thick` | 2px | Focus rings, strong borders |
| `--border-heavy` | 3px | Maximum emphasis |

---

## 12. Typography System

### 12.1 Font Stacks

```css
--font-sans:    'SF Pro Display', 'SF Pro Text', -apple-system, BlinkMacSystemFont,
                'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;

--font-mono:    'SF Mono', 'Fira Code', 'JetBrains Mono', 'Cascadia Code',
                'Consolas', 'Liberation Mono', monospace;

--font-display: 'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
```

RTL Arabic: `'Cairo', 'Noto Sans Arabic', sans-serif` with `line-height: 1.8`.

### 12.2 Type Scale

| Token | Size | Pixels | Usage |
|-------|------|--------|-------|
| `--text-2xs` | 0.625rem | 10px | Micro labels |
| `--text-xs` | 0.75rem | 12px | Captions, metadata |
| `--text-sm` | 0.875rem | 14px | Secondary text, labels |
| `--text-base` | 1rem | 16px | Body text |
| `--text-lg` | 1.125rem | 18px | Large body, subtitles |
| `--text-xl` | 1.25rem | 20px | Card titles |
| `--text-2xl` | 1.5rem | 24px | Section headings |
| `--text-3xl` | 1.875rem | 30px | Page headings |
| `--text-4xl` | 2.25rem | 36px | Hero subheadings |
| `--text-5xl` | 3rem | 48px | Hero headings |
| `--text-6xl` | 3.75rem | 60px | Display headings |
| `--text-7xl` | 4.5rem | 72px | Maximum display |

### 12.3 Fluid Display Sizes

For hero sections, headings use `clamp()` for responsive scaling:

| Token | Value | Range |
|-------|-------|-------|
| `--text-display-sm` | `clamp(1.875rem, 4vw, 2.25rem)` | 30–36px |
| `--text-display-md` | `clamp(2.25rem, 5vw, 3rem)` | 36–48px |
| `--text-display-lg` | `clamp(3rem, 6vw, 4.5rem)` | 48–72px |
| `--text-display-xl` | `clamp(3.75rem, 8vw, 6rem)` | 60–96px |

### 12.4 Font Weights

| Token | Value | Usage |
|-------|-------|-------|
| `--font-thin` | 100 | Decorative display |
| `--font-light` | 300 | Light emphasis |
| `--font-regular` | 400 | Body text |
| `--font-medium` | 500 | Labels, navigation |
| `--font-semibold` | 600 | Buttons, card titles |
| `--font-bold` | 700 | Headings |
| `--font-extrabold` | 800 | Hero headings |
| `--font-black` | 900 | Maximum weight display |

### 12.5 Line Heights & Letter Spacing

| Line Height Token | Value | Letter Spacing Token | Value |
|-------------------|-------|---------------------|-------|
| `--leading-none` | 1 | `--tracking-tighter` | -0.05em |
| `--leading-tight` | 1.2 | `--tracking-tight` | -0.025em |
| `--leading-snug` | 1.35 | `--tracking-normal` | 0 |
| `--leading-normal` | 1.5 | `--tracking-wide` | 0.025em |
| `--leading-relaxed` | 1.65 | `--tracking-wider` | 0.05em |
| `--leading-loose` | 1.8 | `--tracking-widest` | 0.1em |

---

## 13. Animation & Motion

The motion system follows **iOS 26 Liquid Motion** principles: fluid organic movement, slight overshoot for playfulness, and crisp settle for precision.

### 13.1 Duration Scale

| Token | Value | Usage |
|-------|-------|-------|
| `--duration-instant` | 75ms | Micro-feedback (ripples, checkmarks) |
| `--duration-fast` | 150ms | Hover states, toggles |
| `--duration-base` | 250ms | Standard transitions |
| `--duration-moderate` | 350ms | Panel reveals, slides |
| `--duration-slow` | 500ms | Modal entry/exit |
| `--duration-slower` | 700ms | Page transitions |
| `--duration-slowest` | 1000ms | Complex orchestrated reveals |

### 13.2 Easing Curves

| Token | Cubic Bezier | Character |
|-------|-------------|-----------|
| `--ease-default` | `0.25, 0.1, 0.25, 1` | Standard smooth |
| `--ease-liquid` | `0.4, 0, 0.2, 1` | Material standard |
| `--ease-liquid-out` | `0, 0, 0.2, 1` | Decelerate (entry) |
| `--ease-liquid-in` | `0.4, 0, 1, 1` | Accelerate (exit) |
| `--ease-spring` | `0.34, 1.56, 0.64, 1` | iOS spring bounce |
| `--ease-spring-gentle` | `0.22, 1.2, 0.36, 1` | Subtle spring overshoot |
| `--ease-spring-snappy` | `0.68, -0.2, 0.32, 1.2` | Crisp spring |
| `--ease-elastic` | `0.68, -0.6, 0.32, 1.6` | Playful elastic stretch |
| `--ease-bounce` | `0.34, 1.56, 0.64, 1` | Bouncy settle |
| `--ease-dramatic-in` | `0.6, 0, 0.7, 0.2` | Modal entrance |
| `--ease-dramatic-out` | `0.2, 0.8, 0.2, 1` | Modal / sheet exit |

### 13.3 Logo-Inspired Animations

**Pulse Ring** — Heartbeat-line inspired concentric ring animation using `--brand-accent` green. For status indicators and live data points.
```css
@keyframes pulse-ring {
  0%   { box-shadow: 0 0 0 0 rgba(28, 215, 96, 0.40); }
  50%  { box-shadow: 0 0 0 8px rgba(28, 215, 96, 0.15); }
  100% { box-shadow: 0 0 0 16px rgba(28, 215, 96, 0.0); }
}
```

**Float** — Gentle 20s figure-eight drift for background orbs. Creates ambient depth without distraction.

**Skeleton Shimmer** — 1.5s gradient sweep for loading states. Uses neutral scale colors for subtlety.

**Gradient Border Shift** — Animated gradient position for hero section borders, creating a slowly breathing color effect.

### 13.4 Animation Utility Classes

| Class | Description |
|-------|-------------|
| `.animate-float` | 20s figure-eight background orb drift |
| `.animate-pulse` | 2s opacity pulse for loading |
| `.animate-pulse-glow` | 3s scale + opacity pulse |
| `.animate-pulse-ring` | 2s heartbeat ring expansion |
| `.skeleton` | Shimmer loading placeholder |

---

## 14. Z-Index Architecture

A structured layering system prevents z-index conflicts. Every UI element has a designated layer range.

| Token | Value | Layer |
|-------|-------|-------|
| `--z-deep` | -10 | Decorative backgrounds |
| `--z-background` | -1 | Animated mesh, orbs |
| `--z-base` | 0 | Normal content flow |
| `--z-raised` | 10 | Cards, interactive elements |
| `--z-dropdown` | 100 | Dropdown menus |
| `--z-sticky` | 200 | Sticky headers, navigation bars |
| `--z-fixed` | 300 | Fixed-position elements |
| `--z-modal-backdrop` | 400 | Modal overlay / scrim |
| `--z-modal` | 500 | Modal dialogs |
| `--z-popover` | 600 | Popovers, command palette |
| `--z-tooltip` | 700 | Tooltips |
| `--z-toast` | 800 | Toast notifications |
| `--z-spotlight` | 900 | Spotlight / onboarding |
| `--z-max` | 9999 | Skip-to-main, emergencies |

---

## 15. Component Tokens

Pre-composed composite tokens that combine primitives into ready-to-use component values.

### 15.1 Card

| Token | Value |
|-------|-------|
| `--card-bg` | `var(--glass-bg-heavy)` |
| `--card-border` | `1px solid var(--glass-border-soft)` |
| `--card-radius` | `var(--radius-2xl)` = 24px |
| `--card-shadow` | `var(--shadow-ambient)` |
| `--card-shadow-hover` | `var(--shadow-raised)` |
| `--card-backdrop` | `var(--blur-sm)` = blur(16px) |
| `--card-padding` | `var(--space-6)` = 24px |

### 15.2 Button

| Token | Value |
|-------|-------|
| `--btn-radius` | `var(--radius-lg)` = 16px |
| `--btn-padding-sm` | `8px 16px` |
| `--btn-padding-md` | `12px 24px` |
| `--btn-padding-lg` | `16px 32px` |
| `--btn-font-weight` | `600` (semibold) |
| `--btn-transition` | `all 150ms ease-liquid` |

### 15.3 Input

| Token | Value |
|-------|-------|
| `--input-bg` | `var(--surface-primary)` |
| `--input-border` | `1px solid var(--border-default)` |
| `--input-border-focus` | `2px solid var(--brand-primary)` |
| `--input-radius` | `var(--radius-md)` = 12px |
| `--input-padding` | `12px 16px` |

### 15.4 Navigation Bar (iOS-style glass)

| Token | Value |
|-------|-------|
| `--nav-height` | 64px (56px mobile) |
| `--nav-bg` | `var(--glass-bg-heavy)` |
| `--nav-backdrop` | `var(--blur-lg)` = blur(48px) |
| `--nav-border` | `1px solid var(--glass-border-dim)` |
| `--nav-shadow` | `var(--shadow-ambient)` |

### 15.5 Modal / Sheet

| Token | Value |
|-------|-------|
| `--modal-bg` | `var(--glass-bg-pronounced)` |
| `--modal-backdrop` | `var(--blur-md)` |
| `--modal-radius` | `var(--radius-modal)` = 28px |
| `--modal-shadow` | `var(--shadow-glass-floating)` |
| `--modal-overlay` | Light: `rgba(10,15,26, 0.45)` / Dark: `rgba(0,0,0, 0.65)` |

### 15.6 Badge / Tag

| Token | Value |
|-------|-------|
| `--badge-radius` | `var(--radius-pill)` |
| `--badge-padding` | `4px 12px` |
| `--badge-font-size` | `var(--text-xs)` = 12px |
| `--badge-font-weight` | `600` (semibold) |

### 15.7 Tooltip

| Token | Value |
|-------|-------|
| `--tooltip-bg` | Light: `--neutral-900` / Dark: `--neutral-200` |
| `--tooltip-text` | Light: `--neutral-0` / Dark: `--neutral-900` |
| `--tooltip-radius` | `var(--radius-sm)` = 8px |
| `--tooltip-padding` | `8px 12px` |

### 15.8 Data Visualization (Logo-Inspired)

Chart colors follow the exact logo gradient spectrum:

| Token | Hex | Role |
|-------|-----|------|
| `--chart-color-1` | `#1B3FA0` | Series 1 — Primary data |
| `--chart-color-2` | `#1A5DAE` | Series 2 |
| `--chart-color-3` | `#1A7A9E` | Series 3 |
| `--chart-color-4` | `#17B5A6` | Series 4 |
| `--chart-color-5` | `#1CD760` | Series 5 — Accent |
| `--chart-color-6` | `#2AE668` | Series 6 — Highlight |

Supporting tokens:
- `--chart-grid`: `rgba(100,116,139, 0.12)` (light) / `rgba(148,163,184, 0.08)` (dark)
- `--chart-axis`: `--neutral-400` (light) / `--neutral-500` (dark)
- `--chart-label`: `--text-secondary`

### 15.9 Pulse Indicator (Heartbeat-Inspired)

| Token | Value |
|-------|-------|
| `--pulse-color` | `var(--brand-accent)` |
| `--pulse-ring-1` | `rgba(28, 215, 96, 0.3)` |
| `--pulse-ring-2` | `rgba(28, 215, 96, 0.15)` |
| `--pulse-ring-3` | `rgba(28, 215, 96, 0.05)` |

---

## 16. Accessibility Standards

The design system is built to **WCAG 3.0 Bronze** compliance with additional support for platform-specific accessibility modes.

### 16.1 Supported Accessibility Modes

| Mode | Detection | Behavior |
|------|----------|----------|
| Dark mode | `prefers-color-scheme: dark` | Full dark mode with inverted surfaces, adjusted glass, enhanced glows |
| High contrast | `prefers-contrast: high` | Opaque glass fallbacks, stronger borders, 3px focus rings, boosted text |
| Reduced motion | `prefers-reduced-motion: reduce` | All animations disabled (0.01ms), scroll-behavior auto |
| Forced colors | `forced-colors: active` | Windows HC. Glass → Canvas, borders → CanvasText, no box-shadows |
| Legacy HC | `-ms-high-contrast` | IE/Edge legacy with WindowText focus indicators |
| Reduced transparency | `prefers-reduced-transparency` | `--surface-opaque` fallback for glass |

### 16.2 Focus System

| Context | Light Mode | Dark Mode |
|---------|-----------|-----------|
| Default | 3px solid `#1B3FA0`, offset 3px | 3px solid `#33D4BB`, offset 3px |
| High contrast | 4px solid `currentColor`, offset 3px | 4px solid `currentColor`, offset 3px |
| Forced colors | 3px solid `Highlight`, offset 3px | 3px solid `Highlight`, offset 3px |
| Gold (fallback) | `#FFD700` for maximum visibility | `#FFD700` for maximum visibility |

Focus ring implementation:
```css
--focus-ring: 0 0 0 var(--focus-ring-offset) var(--surface-primary),
              0 0 0 calc(var(--focus-ring-offset) + var(--focus-ring-width)) var(--focus-color);
```

### 16.3 Touch Targets

- **Minimum:** 44×44px (WCAG 2.5.5) enforced on all interactive elements at mobile breakpoints
- **Comfortable:** 48×48px recommended for primary actions
- Tokens: `--touch-target-min: 44px` and `--touch-target-comfortable: 48px`

### 16.4 Skip Navigation

A skip-to-main-content link (WCAG 2.4.1 Bypass Blocks) positioned off-screen, slides into view on focus. Uses `--brand-primary` background with white text.

### 16.5 Screen Reader Utilities

`.sr-only` class: visually hidden but accessible to screen readers. Becomes visible on `:focus`, `:active`, and `:focus-within`.

---

## 17. File Architecture & Usage

### 17.1 File Structure

```
├── design-tokens.css    ← Single source of truth for ALL token values
├── globals.css           ← Global styles, reset, utilities, animations, a11y
└── tailwind.config.ts    ← Tailwind v4 integration via @theme bridge
```

**Load order:** `design-tokens.css` must be loaded **before** `globals.css`.

### 17.2 Consumption Patterns

**CSS Variables:**
```css
color: var(--brand-primary);
background: var(--glass-bg-medium);
box-shadow: var(--shadow-floating);
```

**Tailwind Classes:**
```html
<div class="bg-brand-primary text-white shadow-floating rounded-2xl">
```

**Inline Styles (React):**
```jsx
<div style={{ background: 'var(--gradient-brand)' }}>
```

**Glass Utility Classes:**
```html
<div class="glass-panel hover-glass-lift">
<div class="glass-card">
<nav class="glass-nav">
```

**Gradient Text:**
```html
<h1 class="text-gradient-brand">Brand gradient heading</h1>
<h1 class="text-gradient-spectrum">Full logo spectrum heading</h1>
```

### 17.3 Theme Switching

The system supports three theme mechanisms simultaneously:

| Mechanism | Trigger | Usage |
|-----------|---------|-------|
| `@media (prefers-color-scheme: dark)` | OS-level | Automatic detection |
| `[data-theme='dark']` | JavaScript | next-themes or manual toggle |
| `.dark` class | JavaScript | Tailwind dark mode class |

All three override the same CSS custom properties, ensuring consistent rendering regardless of toggle method.

### 17.4 Legacy Aliases

For backward compatibility with existing code:

| Legacy Token | Maps To |
|-------------|---------|
| `--brand-blue` | `var(--brand-primary)` |
| `--brand-teal` | `var(--brand-secondary)` |
| `--brand-green` | `var(--brand-accent)` |
| `--ms-deep-blue` | `var(--brand-primary)` |
| `--ms-blue` | `var(--brand-primary)` |
| `--ms-teal` | `var(--brand-secondary)` |
| `--color-trust-blue` | `var(--brand-primary)` |
| `--color-growth-teal` | `var(--brand-secondary)` |

---

*AnySync Design System v3.0.0 · February 2026*