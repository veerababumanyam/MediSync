# MediSync — Design System

**Version:** 3.0.0  
**Status:** Production Ready  
**Last Updated:** February 22, 2026  
**Maintainer:** MediSync Design Team  
**Standards:** WCAG 3.0,  Apple iOS 26 Liquid Glass HIG

> **Single source of truth.** This file supersedes `DESIGN-GUIDELINES.md` and `LIQUID-GLASS-DESIGN-SYSTEM.md`, both of which have been removed. Updated to align with **WCAG 3.0** (all Guideline 3.3 success criteria) and Apple's **iOS 26 Liquid Glass** design language announced at WWDC 2025.

---

## Table of Contents

1. [Brand Identity](#1-brand-identity)
2. [Color System](#2-color-system)
3. [Typography](#3-typography)
4. [Spacing & Layout](#4-spacing--layout)
5. [Glassmorphism System](#5-glassmorphism-system)
6. [Elevation & Shadows](#6-elevation--shadows)
7. [Border Radius](#7-border-radius)
8. [UI Components](#8-ui-components)
9. [Navigation](#9-navigation)
10. [Animation & Motion](#10-animation--motion)
11. [Accessibility Standards (WCAG 3.0 Bronze + Guideline 3.3)](#11-accessibility-standards-wcag-30-bronze--guideline-33)
12. [Component Library Reference](#12-component-library-reference)
13. [AI Accountant Module](#13-ai-accountant-module-dashboard--real-time-tally-integration)
14. [Design System Maintenance & Evolution](#14-design-system-maintenance--evolution)
15. [Testing & Quality Assurance](#15-testing--quality-assurance)
16. [Implementation Handoff to Development](#16-implementation-handoff-to-development)
17. [Key Takeaways](#17-key-takeaways)
18. [CSS Custom Properties Reference](#18-css-custom-properties-reference)
19. [Implementation Checklist](#19-implementation-checklist)
20. [Internationalisation & RTL Design System](#20-internationalisation--rtl-design-system)
21. [Future Spatial Design (v2.0)](#21-future-spatial-design-v20)

---

## 1. Brand Identity

### 1.1 Core Concept
**"The Interconnected Medical Brain"**
MedMentor AI transforms the overwhelming flood of medical information into a structured, interconnected network of knowledge. Powered by **Deep Research AI Agents**, the system synthesizes global medical insights in real-time, allowing users to query any medical topic or clinical scenario. Document ingestion (EHRs, notes, PDFs) acts as an additive feature to further personalize and ground the AI's research in specific curricula.

**Personalized Greeting Pattern:**
"Good Morning, [Name]. Your brain is ready for new connections." — Emphasizes the partnership between the user and the AI.

### 1.2 Design Philosophy
The visual language draws from three pillars:
*   **iOS 26 Liquid Glass (WWDC 2025):** Apple’s next-generation material system built on **translucency, specular highlights, physics-driven refraction, lensing (light bending), and adaptive tinting**. UI elements dynamically react to touch, ambient light, device tilt, content, and context — creating a fluid, volumetric feel across platforms. MediSync implements this via `backdrop-filter`, radial inner gradients, dynamic highlights, and the `prefers-reduced-transparency` fallback. The iOS 26.1 **Tinted Mode** (increased opacity + neutral color overlay) is supported as an accessibility option for improved contrast and legibility.
*   **iOS-Grade Precision:** Apple-level attention to spacing, typography, micro-interactions, and touch targets. Every pixel is intentional, with sharp text layers always positioned above blurred materials.
*   **Medical Trust & PII Protection**: A palette grounded in deep blues and clean neutrals that conveys clinical authority. This is reinforced by a robust PII protection layer (Microsoft Presidio) that ensures patient data remains anonymized across all AI interactions.

### 1.3 Generative UI
The MedMentor dashboard is a **Living UI**. It utilizes **Generative UI** patterns orchestrated via **CopilotKit** and the **WebMCP** standard to dynamically construct interfaces based on agent reasoning and user data context.
*   **Adaptability:** Interfaces change based on the complexity of the medical query.
*   **Predictive Layouts:** Components are rendered JUST-IN-TIME based on agentic confidence.
*   **Reference:** All generative components follow the [CopilotKit Generative UI Guide](https://github.com/CopilotKit/generative-ui/blob/main/assets/generative-ui-guide.pdf).

### 1.4 WebMCP (Agent-Responsive UI)
Interactive elements are enhanced with **WebMCP** attributes for autonomous discovery by browser-based AI agents.
*   **Declarative discovery**: Standardized `tool-name` and `tool-description` attributes on glass surfaces.
*   **Reference:** [WebMCP Explainer](https://github.com/web-mcp/explainer) | [MediSync WEBMCP.md](file:///Users/v13478/Desktop/MediSync/docs/WEBMCP.md)

### 1.5 Logo Concept
*   **Symbol:** A stylized neural network where nodes and connections form a medical cross in negative space.
*   **Metaphor:** The "spark" of connecting two concepts (synapse firing).
*   **App Icon:** Symbol on a *Trust Blue* or *Midnight Navy* background with a subtle glass-material overlay.
*   **Wordmark:** "MedMentor" in *Inter Bold*, "AI" in *Inter Light* to suggest precision.

---

## 2. Color System

The palette uses a trustworthy medical blue base with energetic teal accents to signify growth and active recall. Extended with a dark-mode glassmorphism palette for the immersive calendar and dashboard experiences.

### 2.1 Primary Brand Colors

| Color Name | Hex | RGB | Contrast (white) | Usage |
| :--- | :--- | :--- | :--- | :--- |
| **Logo Blue** | `#2750a8` | `39, 80, 168` | 5.37:1 ✅ | Core brand identity. Primary buttons, active states, key highlights. Replaces old Trust Blue. |
| **Logo Teal** | `#18929d` | `24, 146, 157` | 3.51:1 ⚠️ | Secondary brand context. Gradient endpoints, secondary glow, badges. Used carefully with large text or dark backgrounds. |
| **Midnight Navy** | `#0f172a` | `15, 23, 42` | 17.85:1 ✅ | Primary text, deep backgrounds, dark mode base. |

### 2.2 Apple HIG Liquid Glass Palette
Used for immersive experiences like the calendar, dashboards, and focus mode. Built on layered transparency over animated mesh backgrounds.

| Color Name | Hex / Value | Opacity | Usage |
| :--- | :--- | :--- | :--- |
| **Glass Background** | `rgba(255,255,255,0.12)` | 12% | Primary card/container surfaces |
| **Glass Hover** | `rgba(255,255,255,0.18)` | 18% | Hover state for glass containers |
| **Glass Border** | `rgba(255,255,255,0.20)` | 20% | Default border on glass surfaces |
| **Glass Border Strong** | `rgba(255,255,255,0.35)` | 35% | Active/focus state borders |
| **Glass Subtle** | `rgba(255,255,255,0.05)` | 5% | Muted backgrounds, toggle groups |
| **Deep Background** | `#0A0A1A` | 100% | Base canvas behind mesh gradients |
| **Elevated Surface** | `rgba(30,30,50,0.85)` | 85% | Dropdown panels, modal sheets |

### 2.3 iOS Accent Palette
Extended accent colors used for event categories, tags, and data visualization. Directly mapped from Apple’s Human Interface Guidelines.

| Color Name | Hex | CSS Variable | Usage |
| :--- | :--- | :--- | :--- |
| **System Blue** | `#007AFF` | `--accent` | Primary interactive, links, today indicator |
| **System Purple** | `#5856D6` | `--purple` | Secondary accent, gradient endpoints |
| **System Pink** | `#FF2D55` | `--pink` | Social events, alerts, notifications |
| **System Orange** | `#FF9500` | `--warning` | Warning states, personal events |
| **System Green** | `#34C759` | `--success` | Success, health events, completion |
| **System Red** | `#FF3B30` | `--danger` | Error states, destructive actions |
| **System Teal** | `#5AC8FA` | `--teal` | Travel, info badges, light accents |
| **System Indigo** | `#AF52DE` | `--purple` | Creative events, design reviews |

### 2.4 Neutral Scale (Slate)

| Token | Hex | Usage |
| :--- | :--- | :--- |
| **Slate 900** | `#0F172A` | Headings, Primary Text (Light Mode) |
| **Slate 700** | `#334155` | Secondary Text, Icons |
| **Slate 500** | `#64748B` | Captions, Placeholder Text, Disabled States |
| **Slate 300** | `#CBD5E1` | Subtle borders, dividers |
| **Slate 200** | `#E2E8F0` | Dividers, Borders |
| **Slate 100** | `#F1F5F9` | Page Backgrounds (Light Mode) |
| **Slate 50** | `#F8FAFC` | Card Backgrounds (Light Mode) |
| **White** | `#FFFFFF` | Surface Backgrounds, Cards, Modals |

### 2.5 Semantic Colors

| Context | Light Mode | Dark Mode | Usage |
| :--- | :--- | :--- | :--- |
| **Success** | `#10B981` | `#34C759` | Synced, correct answers, completion |
| **Warning** | `#F59E0B` | `#FF9500` | Pending, low confidence, uncertainty |
| **Error** | `#EF4444` | `#FF3B30` | Invalid, destructive actions, failed |
| **Info** | `#0EA5E9` | `#007AFF` | Tooltips, guidance, help |

### 2.6 Color Usage Rules

**DO:**
- Use semantic colors for their intended purpose only
- Maintain 4.5:1 contrast ratio for all body text (WCAG 3.0 Bronze)
- Test color combinations in both light and dark modes
- Reserve brand colors for primary actions

**DON'T:**
- Use color as the *only* indicator of state — always pair with an icon or text label
- Use red/green as the only visual differentiation (color-blind users)
- Use more than 3 brand colors in a single view
- Set brand color opacity below 20% for text

### 2.7 Gradients

**Light Mode Surface Gradient**
```css
background: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.4));
```

**Dark Mode Accent Gradient (Text-Bearing — WCAG-Safe)**
```css
background: linear-gradient(135deg, #0056D2, #0F766E);
```

**Dark Mode Accent Glow (Decorative Only — No Text)**
```css
background: linear-gradient(135deg, #2750a8, #18929d);
```

**Study Card Gradient (Cardiology)**
```css
background: linear-gradient(135deg, #1e3a8a 0%, #0891b2 100%); /* Deep Blue to Cyan-Teal */
```

**Animated Background Orbs (Premium Canvas)**
```css
/* Deep blurred orbs acting as light sources */
background:
  radial-gradient(circle at 10% 20%, rgba(39, 80, 168, 0.45) 0%, transparent 60%), /* Logo Blue */
  radial-gradient(circle at 80% 80%, rgba(24, 146, 157, 0.35) 0%, transparent 60%); /* Logo Teal */
```

**Glass Shine (Liquid Effect)**
```css
background: linear-gradient(135deg,
  transparent 40%,
  rgba(255,255,255,0.04) 45%,
  rgba(255,255,255,0.08) 50%,
  rgba(255,255,255,0.04) 55%,
  transparent 60%);
animation: shineSlide 6s ease-in-out infinite;
```

---

## 3. Typography

### 3.1 Font Family
*   **Primary (LTR / Latin):** [Inter](https://fonts.google.com/specimen/Inter) — Variable weight support, highly legible, neutral but modern.
*   **Display Alternative (LTR):** [Plus Jakarta Sans](https://fonts.google.com/specimen/Plus+Jakarta+Sans) — Used in immersive/glassmorphism contexts for display headings only. Adds warmth and character.
*   **Monospace:** JetBrains Mono or SF Mono — Code snippets, medical codes, technical references.
*   **Arabic Primary:** [Cairo](https://fonts.google.com/specimen/Cairo) — Geometric, professional Arabic sans-serif. Matches Inter's clean proportions. Used for all Arabic UI text, report body, and chat.
*   **Arabic Fallback:** [Noto Sans Arabic](https://fonts.google.com/noto/specimen/Noto+Sans+Arabic) — Universal coverage, multi-weight, ideal as fallback and for technical labels.
*   **Arabic Display:** [Tajawal](https://fonts.google.com/specimen/Tajawal) — Optionally used for Arabic hero/headline text (dashboard titles, empty-state headings).

### 3.1.1 Liquid Glass Typography Rules
Typography in the Liquid Glass system is optimized for extreme legibility against translucent backgrounds:
*   **Weight & Alignment**: Prefer **bolder weights** and **left-aligned** layouts to ground the eye.
*   **Contrast Layering**: Sharp text (Slate 900 or White) must always be layered *above* the glass material. **Never blur typography.**
*   **WCAG 3.0 Bronze Fixes**: Ensure ≥4.5:1 text contrast. Use media queries to disable gradients on text for users with vision impairments. See §11.9 for full Guideline 3.3 Input Assistance requirements.
*   **Spacing**: Adhere to WCAG 1.4.12 for line height (1.5) and paragraph spacing (2.0) to ensure readability.

### 3.1.2 Iconography
*   **Solid & Filled**: Use solid or filled icons for primary actions to ensure ≥3:1 contrast against glass backdrops.
*   **Layering**: Place icons in a distinct functional layer above the glass surface.

> All Arabic fonts are available on Google Fonts under the SIL Open Font License (OFL) — free for commercial use.

**CSS Font Stack:**
```css
/* Latin screens */
--font-primary: 'Inter', system-ui, -apple-system, sans-serif;

/* Arabic screens — applied when html[lang='ar'] */
--font-arabic: 'Cairo', 'Noto Sans Arabic', 'Tajawal', sans-serif;

/* Applied via :lang() selector */
:lang(ar) {
  font-family: var(--font-arabic);
  font-feature-settings: 'kern' 1;
  line-height: 1.8; /* Arabic text needs more line height than Latin */
}
```

### 3.2 Type Scale

| Style | Weight | Size (rem/px) | Line Height | Tracking | Usage |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Display H1** | ExtraBold (800) | `2.25rem` / `36px` | `1.1` | `-0.02em` | Marketing, Hero text, Calendar header |
| **H1** | Bold (700) | `1.875rem` / `30px` | `1.2` | `-0.01em` | Page Titles |
| **H2** | SemiBold (600) | `1.5rem` / `24px` | `1.3` | `-0.01em` | Section Headers, Card Titles |
| **H3** | Medium (500) | `1.25rem` / `20px` | `1.4` | `0` | Subsection Headers |
| **Body Large** | Regular (400) | `1.125rem` / `18px` | `1.6` | `0` | Introduction text, Focal content |
| **Body** | Regular (400) | `1rem` / `16px` | `1.5` | `0` | Standard paragraph text |
| **Small** | Regular (400) | `0.875rem` / `14px` | `1.5` | `0` | Metadata, Secondary info |
| **Caption** | SemiBold (600) | `0.75rem` / `12px` | `1.5` | `0.05em` | Labels, Uppercase tags, Status bar |
| **Micro** | SemiBold (600) | `0.625rem` / `10px` | `1.5` | `0.05em` | Nav labels, badge text |

### 3.3 Glass Mode Typography
In dark/glass mode, text uses opacity-based color hierarchy instead of hex values for seamless blending with translucent surfaces:

| Token | Value | Usage |
| :--- | :--- | :--- |
| `--text-primary` | `rgba(255, 255, 255, 0.95)` | Headings, active labels, key data |
| `--text-secondary` | `rgba(255, 255, 255, 0.60)` | Body text, descriptions, subtitles |
| `--text-tertiary` | `rgba(255, 255, 255, 0.35)` | Disabled, placeholder, weekday labels |

---

## 4. Spacing & Layout

### 4.1 Base Unit
**4px Grid System.** All spacing, sizing, and typography line-heights are multiples of 4.

| Token | Size | Value |
| :--- | :--- | :--- |
| `space-1` | 4px | `0.25rem` |
| `space-2` | 8px | `0.5rem` |
| `space-3` | 12px | `0.75rem` |
| `space-4` | 16px | `1rem` |
| `space-5` | 20px | `1.25rem` |
| `space-6` | 24px | `1.5rem` |
| `space-8` | 32px | `2rem` |
| `space-10` | 40px | `2.5rem` |
| `space-12` | 48px | `3rem` |
| `space-16` | 64px | `4rem` |

### 4.2 Containers
*   **Mobile:** 100% width, `16px` horizontal padding. Max calendar width: `480px`.
*   **Tablet:** Max-width `768px`, centered, `24px` padding.
*   **Desktop:** Max-width `1200px`, centered, `32px` padding.

### 4.3 Component Gap System
| Context | Gap | Notes |
| :--- | :--- | :--- |
| **Between sections** | `24px` (`space-6`) | Major content blocks |
| **Between cards** | `10–16px` | Event cards, stat cards |
| **Inside card padding** | `16–20px` | Standard glass containers |
| **Grid cell gap** | `2–3px` | Calendar day grid |
| **Icon + text** | `6–12px` | Inline label pairs |
| **Nav items** | `space-around` | Bottom navigation |

### 4.4 Alignment & RTL
*   **Text alignment:** Use `text-start` / `text-end` (not `text-left` / `text-right`) so content flips correctly in RTL. Use `text-end` for numeric columns and end-aligned UI (e.g. helper text under inputs).
*   **Logical borders:** Use `border-s-*` (inline-start) for accent bars (e.g. error cards) so the bar stays on the correct side in RTL.
*   **Positioning:** Use `start-0` / `end-0` (or Tailwind logical insets) for dropdowns and anchored UI so they open on the correct side in RTL.
*   **Implementation:** See `frontend/src/styles/liquid-glass.css` for `.liquid-text-start`, `.liquid-text-end`, and spacing utility comments.

---

## 5. Glassmorphism System (iOS 26 Liquid Glass)
The iOS 26 Liquid Glass system is the signature visual layer of MediSync’s immersive interfaces. Announced at WWDC 2025, it evolves traditional glassmorphism into a physics-driven material with **specular highlights, refraction, lensing, and adaptive tinting** that dynamically responds to content, touch, ambient light, and device orientation — creating a sense of depth, material, and space using layered translucency, blur, and light effects.

### 5.1 Glass Material Classes (Dark / Glass Mode)
| Class | Background | Blur | Border | Shadow |
| :--- | :--- | :--- | :--- | :--- |
| `.glass` | `rgba(255,255,255,0.12)` | `blur(40px)` | `rgba(255,255,255,0.20)` | `0 8px 32px rgba(0,0,0,0.3)` |
| `.glass-elevated` | `rgba(255,255,255,0.08)` | `blur(60px)` | `rgba(255,255,255,0.20)` | `0 20px 60px rgba(0,0,0,0.4)` |
| `.glass-subtle` | `rgba(255,255,255,0.05)` | `blur(20px)` | `rgba(255,255,255,0.08)` | None |
| `.glass-panel` (light) | `rgba(255,255,255,0.70)` | `blur(12px)` | `rgba(255,255,255,0.30)` | Level 1 |

### 5.2 Glass Materials (Apple HIG Equivalent)
| Material | Background (Light) | Background (Dark) | Blur | Stroke/Border | Usage |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Ultra Thin** | `rgba(255,255,255,0.4)` | `rgba(30,30,30,0.3)` | `blur(20px)` | `rgba(255,255,255,0.4)` top edge | Subtle backdrop regions |
| **Thin** | `rgba(255,255,255,0.55)` | `rgba(35,35,35,0.4)` | `blur(30px)` | `rgba(255,255,255,0.5)` top edge | Default cards, panels |
| **Regular** | `rgba(255,255,255,0.7)` | `rgba(45,45,45,0.6)` | `blur(45px)` | `rgba(255,255,255,0.6)` top edge | Modals, Prominent surfaces |
| **Thick** | `rgba(255,255,255,0.85)` | `rgba(60,60,60,0.8)` | `blur(60px)` | `rgba(255,255,255,0.7)` top edge | Navbars, highly readable zones |

### 5.3 Hover Effects
| Effect | Description | Duration |
| :--- | :--- | :--- |
| **Lift** | `translateY(-4px)` + stronger shadow | 300ms |
| **Glow** | Colored ambient light (teal/blue/green) | 300ms |
| **Shimmer** | Moving light reflection across surface | 1.5s |
| **Lift-Glow** | Combination of lift + colored glow | 300ms |

### 5.4 Liquid Glass Technical Implementation (iOS 26)
MediSync adopts Apple’s iOS 26 Liquid Glass foundation, emphasizing translucent, adaptive materials that create depth and hierarchy through **specular highlights, physics-driven refraction, and dynamic adaptive tinting**.

**Core Web Equivalent:**
```css
/* iOS 26 Liquid Glass Base Material */
.liquid-glass-regular {
  backdrop-filter: blur(45px) saturate(1.8);
  -webkit-backdrop-filter: blur(45px) saturate(1.8);
  background-color: rgba(255, 255, 255, 0.7); /* Adjust per theme */
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.05); /* Diffuse soft shadow */
  border: 1px solid rgba(255,255,255,0.2);
  border-top: 1px solid rgba(255, 255, 255, 0.5); /* iOS 26 specular top edge highlight */
}

/* Dark Mode */
.dark .liquid-glass-regular {
  background-color: rgba(45, 45, 45, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-top: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}

/* WCAG 3.0 Bronze — Reduced Transparency Fallback */
@media (prefers-reduced-transparency: reduce) {
  .liquid-glass-regular {
    backdrop-filter: none;
    background-color: var(--surface-opaque);
    opacity: 1;
  }
}

/* iOS 26.1 Tinted Mode — Accessibility Fallback */
/* Increases opacity + applies neutral tint for improved contrast & legibility */
.liquid-glass-tinted {
  backdrop-filter: blur(30px) saturate(1.2);
  -webkit-backdrop-filter: blur(30px) saturate(1.2);
  background-color: rgba(128, 128, 128, 0.65); /* Neutral tint */
  border: 1px solid rgba(255, 255, 255, 0.25);
}

.dark .liquid-glass-tinted {
  background-color: rgba(40, 40, 50, 0.85); /* Higher opacity for readability */
  border: 1px solid rgba(255, 255, 255, 0.12);
}
```

*   **Lensing & Refraction**: Simulated via radial inner gradients and dynamic highlights that shift based on background or scroll position.
*   **Adaptive Tinting**: Materials brighten at the center and vignetted at the edges to maintain volumetric depth.

### 5.5 Glass Shine (Liquid Reflection)
A signature animated gradient overlay that simulates light passing across a glass surface. Applied via the `::after` pseudo-element with `pointer-events: none`.
*   **Gradient angle:** 135deg diagonal sweep.
*   **Animation:** `shineSlide`, 6s ease-in-out infinite.
*   **Peak opacity:** 8% (`rgba(255,255,255,0.08)`) to maintain subtlety.

### 5.6 Immersive Study Workspace (Agent Studio)
- **Agent-First Layout**:
-    - **Center**: Dynamic AI Chat Integrated with Deep Research Specialists (Primary focus).
-    - **Left**: Collapsible Source Library (Document ingestion as an additive feature).
-    - **Right**: Collapsible "Research Drawer" (Citations, Tools, and Synthesis).
- **Research Breadcrumbs**: A dynamic progress UI element showing the agent's research steps (e.g., "Scanning NEJM..." → "Extracting RCT data..." → "Synthesizing...").
- **Evidence Badge System**: Color-coded badges (Gold/Silver/Bronze) indicating the evidence tier of a response.
- **Visual Pinpointing**: AI citations from both internet sources and uploaded documents must trigger visual highlights or link-outs.
- **Responsive Behavior**: Sidebars collapse to allow for a focused chat experience on small screens.

### 5.7 Animated Background Orbs
Floating radial gradient orbs create a living, breathing canvas behind glass surfaces. Three orbs with different colors, sizes, and animation timings prevent visual repetition.

| Orb | Size | Color Base | Blur | Animation Duration |
| :--- | :--- | :--- | :--- | :--- |
| **Orb 1 (Logo Blue)** | 500px | `rgba(39, 80, 168, 0.45)` | 100px | 20s |
| **Orb 2 (Logo Teal)** | 400px | `rgba(24, 146, 157, 0.35)` | 100px | 25s |

### 5.8 CSS Implementation
```css
.glass {
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.glass-elevated {
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(60px);
  -webkit-backdrop-filter: blur(60px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.08);
}
```

---

## 6. Elevation & Shadows

### 6.1 Light Mode Shadows
| Level | Token | Value | Usage |
| :--- | :--- | :--- | :--- |
| **1** | `shadow-sm` | `0 1px 3px rgba(0,0,0,0.1), 0 1px 2px rgba(0,0,0,0.06)` | Cards at rest |
| **2** | `shadow-md` | `0 4px 6px rgba(0,0,0,0.1), 0 2px 4px rgba(0,0,0,0.06)` | Hover/active cards |
| **3** | `shadow-lg` | `0 20px 25px rgba(0,0,0,0.1), 0 10px 10px rgba(0,0,0,0.04)` | Modals, dropdowns |

### 6.2 Dark / Glass Mode Shadows
| Level | Token | Value | Usage |
| :--- | :--- | :--- | :--- |
| **1** | `glass-shadow` | `0 8px 32px rgba(0,0,0,0.3)` | Standard glass cards |
| **2** | `glass-shadow-elevated` | `0 20px 60px rgba(0,0,0,0.4), inset glow` | Dropdowns, modals |
| **3** | `glass-shadow-fab` | `0 8px 30px rgba(0,122,255,0.5), ring 4px` | FAB, primary CTA |
| **4** | `today-pulse` | Animated 0–30px, color-cycling | Today’s date indicator |

### 6.3 Focus Ring (Accent Glow)
All interactive elements receive a visible focus ring for keyboard navigation:
*   **Light mode:** `box-shadow: 0 0 0 3px rgba(0, 86, 210, 0.2); border-color: #0056D2;`
*   **Glass mode:** `box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.2); border-color: #007AFF;`

---

## 7. Border Radius

| Token | Value | Usage |
| :--- | :--- | :--- |
| `radius-sm` | `8px` | Buttons, tags, input fields, small badges |
| `radius-md` | `12–16px` | Calendar day cells, cards (light mode), dropdown items |
| `radius-lg` | `22px` | Glass containers, event cards, month selector |
| `radius-xl` | `28px` | Calendar card wrapper, modal sheets |
| `radius-full` | `50%` / `9999px` | Icon buttons, avatar circles, FAB, color pickers |

---

## 8. UI Components

### 8.1 Buttons (Liquid Glass Context)
*   **Trait**: Glass with a tint shift on hover/press; specular highlights create a volumetric feel.
*   **Web Equivalent**: use `radial-gradient` + `scale(1.02)` on hover.
*   **WCAG 3.0 Bronze**: Minimum 44x44px touch targets; 3:1 contrast ratio between states; distinct ARIA labels for screen readers. Supports password manager paste (SC 3.3.8).

### 8.2 Cards (Volumetric Float)
*   **Apple Trait**: Volumetric floating effect that refracts background content.
*   **Web Equivalent**: 
    ```css
    .card-glass {
      backdrop-filter: blur(16px);
      border-radius: 20px;
      box-shadow: 0 12px 40px rgba(0,0,0,0.15);
    }
    ```
*   **WCAG**: Ensure text is ≥14px bold or 18px regular; use high-contrast color overrides for legibility.

### 8.3 Panels & Sheets (Modal Dynamics)
*   **Trait**: Backdrop dimming combined with glass material; growth effect on drag.
*   **Web Equivalent**: Overlay with `position: fixed`, `blur(24px)`.
*   **WCAG**: Support `Escape` key to close; maintain ≥3:1 contrast for hover/focus; disable elastic animations for `prefers-reduced-motion`. Backdrop dimming must respect `prefers-reduced-transparency` (iOS 26).

### 8.4 Scroll Views (Edge Effects)
*   **Trait**: Edge effects with soft blur under pinned controls.
*   **WCAG**: Ensure no blur is applied directly to readable text. Support WCAG 1.4.10 reflow.

### 8.5 Controls & Inputs
*   **Sliders/Toggles**: Glass material directly on the control thumb; flex effect on interaction. WCAG: Visible thumb ≥24px; live region announcements for values.
*   **Text Fields**: Glass container with sharp text; auto-adjusting tint for maximum legibility. WCAG: ≥4.5:1 placeholder/focus contrast; native autocomplete support. Must support password manager paste and autofill per SC 3.3.8. Labels required per SC 3.3.2.
*   **Switches/Checkboxes**: Subtle glass "lift" on toggle. WCAG: 3:1 checked/unchecked state contrast.

### 8.6 Alerts & Overlays
*   **Trait**: Glass with a dimming scrim; opacity ramps ensure focus on the message.
*   **WCAG**: Must be interruptible (Escape/Dismiss); ARIA `role="alert"` or `aria-live="polite"`.

### 8.7 Light Mode Card Variants
**Standard Light Cards**
*   Background: White (`#FFFFFF`)
*   Border: 1px solid Slate 200
*   Radius: 16px (`radius-md`)
*   Shadow: `shadow-sm`
*   Hover: `translateY(-2px)`, `shadow-md`, border Blue-200

**Event Cards**
*   Layout: Horizontal flex with 4px colored accent bar on the left edge.
*   Content: Title (15px/600), Meta row with time badge and date label.
*   Avatar: 36px circle with category emoji, background matches event color at 13% opacity.
*   Active: `scale(0.98)` on press for haptic-style feedback.

### 8.8 Calendar Grid
**Day Cells**
*   Size: Square (aspect-ratio: 1), in a 7-column CSS Grid with 3px gap.
*   Default: 15px/500 weight, `--text-primary` color.
*   Other Month: `--text-tertiary` (35% opacity white).
*   Hover: `rgba(255,255,255,0.08)` background, `scale(1.08)`.
*   Selected: `rgba(0,122,255,0.2)` background with 1.5px solid accent border.
*   Today: Accent gradient background, bold white text, animated pulse shadow (0–30px, 3s infinite).
*   Event Dots: Up to 3 colored dots (5px diameter) below the date number.

**Weekday Row**
*   12px uppercase, font-weight 600, `--text-tertiary`, letter-spacing 0.5px. Consistent 8px vertical padding.

### 8.9 Dropdown Menus
Custom dropdown panels with iOS-grade presentation and clear visibility.

**Month Selector Trigger**
*   Container: Glass material with glass-shine effect, `radius-lg` (22px).
*   Layout: Month label (20px/700) + Year label (20px/400, accent-light color) + Chevron icon.
*   Chevron: Rotates 180° on open (0.4s cubic-bezier).
*   Hover: Background transitions to `glass-bg-hover` (18% white).

**Dropdown Panel**
*   Background: `rgba(30, 30, 50, 0.85)` with `blur(60px)` — opaque enough for readability.
*   Border: 1px solid `rgba(255,255,255,0.15)`.
*   Shadow: `0 25px 70px rgba(0,0,0,0.5)`, `0 0 0 1px rgba(255,255,255,0.05)`.
*   Entry animation: opacity 0→1, `translateY(-12px)`→0, `scale(0.97)`→1 over 0.35s.
*   Max height: 320px with styled scrollbar (4px track, 15% white thumb).

**Dropdown Items**
*   Padding: 12px 16px, `radius-sm` (12px).
*   Font: 15px/500, `--text-secondary`.
*   Hover: `rgba(255,255,255,0.1)` background, promote to `--text-primary`.
*   Active/Selected: Accent gradient background, white text, font-weight 600, checkmark suffix.
*   Active shadow: `0 4px 15px rgba(0,122,255,0.4)`.

**Year Navigator**
*   Centered row above the month list with ‹ / › navigation buttons. Year displayed in 17px/700 weight. Separated from month list with a 1px divider at 8% white opacity.

**Select Dropdowns (Form Context)**
*   Styling: Custom `appearance:none` with matching glass-style input treatment.
*   Background: `rgba(255,255,255,0.06)`, border `rgba(255,255,255,0.12)`.
*   Arrow: Custom SVG chevron absolutely positioned (right: 14px).
*   Focus: Accent border + 3px accent glow ring.
*   Options: Background `#1A1A30`, white text for native dropdown readability.

### 8.10 Form Inputs
**Light Mode**
*   Background: Slate 50 (`#F8FAFC`).
*   Border: 1px solid Slate 300.
*   Radius: `radius-md` (12px).
*   Focus: Trust Blue ring (3px solid, 20% opacity) + blue border.

**Glass Mode**
*   Background: `rgba(255,255,255,0.06)`.
*   Border: 1px solid `rgba(255,255,255,0.12)`.
*   Radius: `radius-sm` (12px).
*   Focus: `border-color: #007AFF`, `box-shadow: 0 0 0 3px rgba(0,122,255,0.2)`, bg promotes to 0.08.
*   Placeholder: `--text-tertiary` (35% white).
*   Text: 15px/500, `--text-primary`.

### 8.11 Modal / Bottom Sheet
*   Overlay: `rgba(0,0,0,0.5)` with `blur(8px)` backdrop.
*   Sheet: Slides up from bottom, `rgba(25,25,45,0.92)` with `blur(60px)`.
*   Handle: Centered 40px × 4px bar, `rgba(255,255,255,0.2)`, radius 4px.
*   Border radius: `radius-xl` top corners only (28px 28px 0 0).
*   Animation: `translateY(100%)`→0 over 0.4s cubic-bezier(0.4, 0, 0.2, 1).
*   Shadow: `0 -20px 60px rgba(0,0,0,0.4)`.

### 8.12 Color Picker
*   Layout: Flex row with 10px gap, wrapping.
*   Options: 32px circles with 2px transparent border.
*   Hover: `scale(1.15)`.
*   Selected: 2px white border + 3px white ring (20% opacity) + checkmark overlay.

### 8.13 View Toggle (Segmented Control)
*   Container: `glass-subtle` background, `radius-md` (16px), 4px padding.
*   Buttons: Equal flex, 10px 16px padding, `radius-sm` (12px).
*   Active: Accent gradient background, white text, 0 4px 15px blue shadow.
*   Inactive hover: `rgba(255,255,255,0.08)`, promote text to `--text-primary`.

### 8.14 Dashboard Widgets (Analytics)
**Retention Card**
*   **Visuals:** Dark glass card with a vibrant teal/green smooth curve graph.
*   **Data:** "RETENTION [XX]%" header.
*   **Context:** "Forgetting curve optimization active." (Subtle text).
*   **Accent:** Growth Teal `#00E8C6` for graph line and connection points.

**Predictor Card**
*   **Visuals:** Dark glass card containing a circular progress ring.
*   **Data:** Large central number (e.g., "245"), Label "Step 1" or "Step 2".
*   **Badge:** "ON TRACK" pill badge (Teal background, low opacity).
*   **Progress Ring:** Gradient blue to teal glow.

**Technical Note:** All charts and analytics widgets are standardized on **Apache ECharts** for cross-platform consistency. The `go-echarts` library is used for server-side configuration in the Go backend.

### 8.15 Daily Queue Card
A primary action area for the Space Repetition System (SRS).
*   **Header:** "Daily Queue" + "Space Repetition System" subtext.
*   **Stat:** Large counter "120 CARDS DUE" (White display type).
*   **Content:**
    *   **Inner Glass Card:** "Mixed Review" container with icon and 3-dot complexity indicator (Orange, Yellow, Green).
    *   **Action Button:** Full-width "Start Review" button (Primary Blue/Teal gradient).

---

## 9. Navigation

### 9.1 Bottom Tab Bar (Liquid Lift)
*   **Behavior**: Liquid Glass lifts navigation above content. As content scrolls, the material adapts its opacity and blur dynamically.
*   **WCAG Fixes**: 
    *   Minimum 48px touch targets.
    *   Icons and labels must maintain ≥4.5:1 contrast.
    *   Active icons use a distinct glow (`#00E8C6`) for clear state communication.

### 9.2 Sidebars (Content Inset)
*   **Visuals**: Inset panels that allow underlying content to flow behind.
*   **Accessibility**: High-contrast text only; must be fully keyboard-navigable with clear focus rings.

### 9.3 Top Nav Bars (Adaptive Dimming)
*   **Scroll Behavior**: Translucent over hero sections, dimming dynamically on scroll to prioritize foreground content.
*   **WCAG**: Bold labels only; ensure no overlap with interactive search inputs.

### 9.4 Floating Action Button (FAB)
*   **Size:** 56px circle (prominent overlap on tab bar).
*   **Background:** **Cyan/Teal Gradient** (distinct from the purple/blue brand gradient).
    *   Matches "Growth Teal" aesthetic.
*   **Shadow:** Strong Cyan Glow `0 0 20px rgba(0, 232, 198, 0.4)`.
*   **Icon:** Large White "+" (28-32px).
*   **Position:** Center docked in the Bottom Tab Bar.

### 9.3 Status Bar (iOS-Style)
Simulated iOS status bar with time (15px/700), and system icons (WiFi, battery) at 16px/white fill. Provides platform-native immersion in mobile views.

---

## 10. Animation & Motion

### 10.1 Easing Functions
| Token | Value | Usage |
| :--- | :--- | :--- |
| `ease-default` | `cubic-bezier(0.4, 0, 0.2, 1)` | General transitions (Google Material standard) |
| `ease-spring` | `cubic-bezier(0.4, 0, 0.2, 1)` | Buttons, cards, scale interactions |
| `ease-bounce` | `cubic-bezier(0.34, 1.56, 0.64, 1)` | Playful feedback, checkmarks |

### 10.2 Liquid Dynamics
*   **Spring Flex**: Use spring-based timing for transforms to create an organic, physical feel.
*   **Lensing Effect**: Dynamic shifts in `backdrop-filter` values during motion to simulate light bending.
*   **Reduced Motion**: Automatically disable elastic animations and lensing when `@prefers-reduced-motion` is active.

### 10.3 Transition Defaults
*   **Duration:** 0.25–0.35s for micro-interactions, 0.4s for layout shifts.
*   **Properties:** `all` (for simple elements) or explicit (`transform`, `opacity`, `background`).

### 10.3 Key Animations
| Animation | Duration | Easing | Description |
| :--- | :--- | :--- | :--- |
| `fadeInUp` | 0.6s | `cubic-bezier(0.4,0,0.2,1)` | Page load entry. Staggered 50ms per element. |
| `orbFloat1/2/3` | 18–25s | `ease-in-out` | Background orb floating. Different per orb. |
| `shineSlide` | 6s | `ease-in-out` | Glass shine sweep across containers. |
| `todayPulse` | 3s | `ease-in-out` | Glow pulse on today’s calendar cell. |
| `dropdown-enter` | 0.35s | `cubic-bezier(0.4,0,0.2,1)` | Scale + fade + translateY for menus. |
| `modal-slide` | 0.4s | `cubic-bezier(0.4,0,0.2,1)` | Bottom sheet slide up from 100%. |

### 10.4 Stagger Pattern
On page load, elements animate in sequence with 50ms delay increments (`delay-1` through `delay-6`). This creates a cascading reveal that guides the eye from top to bottom.

---

## 11. Accessibility Standards (WCAG 3.0 Bronze + Guideline 3.3)

MediSync unifies hierarchy while staying premium and inclusive. Compliance target: **WCAG 3.0 Bronze** across all outcomes, with explicit coverage of **Guideline 3.3 (Input Assistance)** to help users avoid and correct mistakes, and **iOS 26 Liquid Glass accessibility adaptations** including Tinted Mode.

> **WCAG 3.0 Conformance Model:** WCAG 3.0 replaces the A/AA/AAA levels with **Bronze / Silver / Gold** tiers and uses outcome-based scoring (0–4) instead of binary pass/fail. Bronze is the baseline conformance tier, roughly equivalent to WCAG 2.2 AA. MediSync targets Bronze with aspirational Silver outcomes for critical healthcare and financial flows.

> **Legal context (2026):** The US DOJ ADA Title II rule (April 2024) mandates WCAG 2.1 AA for state/local government digital services by April 24, 2026. MediSync targets the forward-looking **WCAG 3.0 Bronze** standard, which subsumes WCAG 2.2 AA and adds outcome-based scoring and functional categories for visual, auditory, cognitive, and motor accessibility.

### 11.1 Color Contrast
| Element | Minimum Ratio | Notes |
| :--- | :--- | :--- |
| **Normal text** | 4.5:1 | Slate-900+ on white/light |
| **Large text (≥18px regular / ≥14px bold)** | 3:1 | Slate-700+ on white |
| **Interactive components / icons** | 3:1 | Brand color on adjacent background |
| **Glass text (dark mode)** | 7:1+ | `--text-primary` on `#0A0A1A` |
| **Glass text (Tinted Mode)** | 4.5:1+ | Neutral overlay ensures legibility on all glass surfaces |
| **Non-text contrast (UI components)** | 3:1 | Borders, focus indicators, form controls on adjacent color |
| **Placeholder text** | 4.5:1 | `--text-tertiary` must pass against glass backgrounds |

### 11.2 Focus Management
*   All interactive elements: visible focus ring (3px Trust Blue / System Blue, 20% opacity glow).
*   Tab order follows visual layout: header → selector → content → CTA → nav.
*   Modal traps focus; returns to trigger on close.
*   Focus indicators must be visible on **all glass materials** — test against both light and dark mode glass surfaces.
*   Focus rings must meet **3:1 contrast** against adjacent colors (WCAG 2.4.7 / 2.4.11).

### 11.3 Touch Targets
*   Minimum **44×44px** for all interactive elements (WCAG 2.5.8 Target Size).
*   Calendar day cells: `aspect-ratio: 1` ensures adequate tap area.
*   Icon buttons: 40px diameter + 4px invisible hit area extension (total ≥44px).
*   Inline links in dense text: minimum **24×24px** active area with adequate spacing (WCAG 2.5.8).

### 11.4 Motion & Reduced Motion
```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
  /* Disable Liquid Glass lensing/refraction effects */
  .liquid-glass-regular::after,
  .glass::after {
    animation: none !important;
    background: none !important;
  }
}
```

### 11.5 Reduced Transparency (iOS 26 Liquid Glass)
```css
/* WCAG adaptation: users who struggle with translucent backgrounds */
@media (prefers-reduced-transparency: reduce) {
  .liquid-glass-regular,
  .glass,
  .glass-elevated,
  .glass-subtle {
    backdrop-filter: none !important;
    -webkit-backdrop-filter: none !important;
    background-color: var(--surface-opaque) !important;
    opacity: 1 !important;
  }
}
```

### 11.6 High Contrast Mode
Supports `prefers-contrast: high`: increased border visibility (2px solid borders), stronger shadows, maintained brand colors. Glass surfaces gain opaque backgrounds and full-strength borders.

### 11.7 iOS 26.1 Tinted Mode (Accessibility Fallback)
*   When Tinted Mode is active, all Liquid Glass surfaces gain increased opacity and a neutral color overlay.
*   Use `.liquid-glass-tinted` class variant (see §5.4) for programmatic fallback.
*   Text contrast must remain ≥4.5:1 against the tinted surface at all times.
*   Specular highlights and refraction effects are muted (reduced `saturate()`) to avoid visual distraction.

### 11.8 Screen Reader Considerations
*   Calendar cells: `aria-label` with full date (e.g., *Thursday, February 19, 2026*) + `aria-current="date"` on today.
*   Dropdown: `aria-expanded` on trigger. Modal: `role="dialog"` + `aria-labelledby`.
*   Form errors: Announced via `aria-live="assertive"` or `role="alert"` (see §11.9).
*   Glass surfaces: Decorative glass effects (shine, orbs, gradients) must have `aria-hidden="true"` and `pointer-events: none`.
*   Authentication steps: Each step must have `aria-describedby` linking to help text (see §11.9.6).

### 11.9 Guideline 3.3 — Input Assistance (WCAG 3.0 Bronze)

MediSync implements **all** WCAG 3.0 Guideline 3.3 outcomes at Bronze conformance level. This is critical for healthcare forms, financial data entry (AI Accountant), and authentication flows where input errors can have clinical or financial consequences.

#### 11.9.1 Error Identification — SC 3.3.1 (Level A)
> *If an input error is automatically detected, the item that is in error is identified and the error is described to the user in text.*

*   **Implementation:**
    *   All form fields with validation errors display a text error message **adjacent to the field** (below, in `--danger` color, minimum 14px).
    *   The erroneous field receives a visible red border (`--danger` + 2px) **and** an error icon — color is never the sole indicator.
    *   Error messages use specific language: *"Invoice amount must be a positive number"* — not *"Invalid input."*
    *   On form submission with errors, programmatic focus moves to the **first** erroneous field.
    *   Errors are announced to assistive technology via `aria-invalid="true"` and `aria-describedby` linking to the error message.

```html
<!-- Correct MediSync error identification pattern -->
<div class="liquid-glass-input-wrapper">
  <label for="invoice-amount">Invoice Amount *</label>
  <input id="invoice-amount"
         type="number"
         aria-invalid="true"
         aria-describedby="invoice-amount-error"
         class="liquid-glass-input state-error" />
  <p id="invoice-amount-error" class="field-error" role="alert">
    ⚠ Invoice amount must be a positive number (e.g., 1500.00)
  </p>
</div>
```

#### 11.9.2 Labels or Instructions — SC 3.3.2 (Level A)
> *Labels or instructions are provided when content requires user input.*

*   **Implementation:**
    *   Every form field has a visible `<label>` element with `for` attribute, or uses `aria-label` / `aria-labelledby`.
    *   Required fields are marked with an asterisk (`*`) **and** the text "required" in the label or via `aria-required="true"`.
    *   Complex inputs (date ranges, currency amounts, invoicing codes) include helper text below the field explaining expected format.
    *   Upload zones display accepted file types and size limits prominently.
    *   Multi-step forms include a step indicator with clear "Step X of Y" labelling.

```html
<!-- Label + instruction pattern for financial forms -->
<label for="gst-number">GST Number *</label>
<input id="gst-number"
       type="text"
       aria-required="true"
       aria-describedby="gst-hint"
       placeholder="22AAAAA0000A1Z5" />
<p id="gst-hint" class="field-hint">
  15-character alphanumeric GST identification number
</p>
```

#### 11.9.3 Error Suggestion — SC 3.3.3 (Level AA)
> *If an input error is automatically detected and suggestions for correction are known, then the suggestions are provided to the user, unless it would jeopardize the security or purpose of the content.*

*   **Implementation:**
    *   When a field value is rejected, the error message includes a **specific correction suggestion** when feasible:
        *   Date field: *"Date must be in DD/MM/YYYY format. Did you mean 15/02/2026?"*
        *   Email field: *"Invalid email address. Did you mean user@hospital.com?"*
        *   Ledger mapping: *"No matching ledger found for 'Ofice Supplies'. Suggestion: 'Office Supplies' (Ledger 4301)."*
    *   AI-powered suggestions (OCR corrections, ledger auto-mapping) show confidence scores alongside suggestions.
    *   Password fields: Show criteria checklist with ✓/✗ indicators for each rule (length, uppercase, number, special char).
    *   For healthcare-specific inputs (ICD codes, drug names), suggest closest matches from the validated database.

#### 11.9.4 Error Prevention: Legal, Financial, Data — SC 3.3.4 (Level AA)
> *For pages with legal commitments, financial transactions, or user-controllable data modification, at least one of: submissions are reversible, data is checked and corrected, or a review/confirm step is provided.*

*   **Implementation (Critical for AI Accountant):**
    *   **Financial transactions** (Tally sync, invoice posting, bill payment):
        *   **Review step:** Confirmation modal summarizing all data before submission.
        *   **Reversibility:** "Undo" option available for 30 seconds after submission for non-destructive operations.
        *   **Data validation:** Server-side validation with inline error feedback before final commit.
    *   **Bulk operations** (Confirm All, Delete Selected):
        *   Explicit confirmation dialog: *"You are about to confirm 24 transactions totaling ₹15,40,000. This action will sync to Tally. Continue?"*
        *   Destructive actions require typing "DELETE" or similar explicit confirmation.
    *   **Patient data** (EHR uploads, clinical notes):
        *   PII detection warning before submission if sensitive data is detected.
        *   Preview/review screen showing anonymization results before AI processing.
    *   **Account settings changes:** Email change requires verification; password change shows preview of affected sessions.

#### 11.9.5 Redundant Entry — SC 3.3.7 (Bronze, carried from WCAG 2.2)
> *Information previously entered by or provided to the user that is required to be entered again in the same process is either auto-populated or available for the user to select.*

*   **Implementation:**
    *   Multi-step form flows (bill upload → mapping → confirmation) carry forward all previously entered data.
    *   Address / vendor / company fields auto-populate from previous entries or user profile.
    *   If a user navigates back in a multi-step flow, all fields retain their previous values.
    *   Search filters persist across pagination and navigation within the same session.
    *   Exception: Re-entering a password for security confirmation is permitted (not redundant entry).

#### 11.9.6 Accessible Authentication (Minimum) — SC 3.3.8 (Bronze, carried from WCAG 2.2)
> *A cognitive function test (such as remembering a password) is not required for any step in an authentication process unless an alternative is provided, a mechanism assists the user, or the test involves object recognition.*

*   **Implementation:**
    *   **Password manager support:** No `autocomplete="off"` on login fields. Native browser autofill and password manager paste must work.
    *   **Copy/paste allowed:** Users must be able to paste credentials from clipboard into all authentication fields.
    *   **Alternative authentication methods:** At least one non-memory method is available:
        *   Magic link (email-based OTP)
        *   Biometric authentication (Face ID / Touch ID on iOS 26)
        *   SSO / OAuth (Google Sign-In, Microsoft Entra ID)
    *   **CAPTCHA:** If used, must provide an audio alternative **and** be solvable without solving a puzzle (e.g., reCAPTCHA v3 silent scoring or Turnstile).
    *   **2FA:** TOTP codes support paste from authenticator apps. WebAuthn (passkeys) supported as zero-knowledge alternative.
    *   **Session management:** Reasonable session durations (≥30 minutes active) to avoid frequent re-authentication.

### 11.10 Keyboard Navigation
| Key | Action |
| :--- | :--- |
| `Tab` / `Shift+Tab` | Forward/backward navigation |
| `Enter` / `Space` | Activate button, toggle, or link |
| `Escape` | Close modal/dropdown/sheet |
| Arrow keys | Navigate date pickers, menus, comboboxes |
| `Home` / `End` | Jump to first/last item in lists |
| `Ctrl+Z` | Undo last action (where supported) |

### 11.11 Glassmorphism-Specific Accessibility Patterns
| Concern | Mitigation |
| :--- | :--- |
| **Blur reduces text legibility** | Text is always in a layer *above* the blur; never blurred. Minimum 4.5:1 contrast verified against blurred background. |
| **Translucent backgrounds vary** | Color contrast tested against worst-case (lightest) background content behind glass. |
| **Animated glass shine distracts** | Respects `prefers-reduced-motion`; decorative elements carry `aria-hidden="true"`. |
| **Low-end device performance** | Progressive enhancement: glass effects degrade to solid backgrounds on devices that don't support `backdrop-filter`. |
| **iOS 26 Liquid Glass refraction** | Refraction/lensing effects disabled under `prefers-reduced-motion` and `prefers-reduced-transparency`. |
| **Color-only state indicators** | All status indicators pair color with text label + icon (e.g., red ✗ "Failed" not just red dot). |

---

## 12. Component Library Reference

### 12.1 File Structure
```
frontend/src/
├── components/
│   ├── ui/
│   │   ├── LiquidGlassCard.tsx
│   │   ├── LiquidGlassButton.tsx
│   │   ├── LiquidGlassInput.tsx
│   │   └── index.ts
│   ├── chat/
│   ├── dashboard/
│   └── ...
├── styles/
│   ├── liquid-glass.css
│   └── globals.css
└── lib/
    └── cn.ts
```

### 12.2 Naming Conventions
| Type | Convention | Example |
| :--- | :--- | :--- |
| Components | PascalCase | `LiquidGlassCard` |
| Props interface | PascalCase + `Props` | `LiquidGlassCardProps` |
| CSS classes | kebab-case prefixed | `.liquid-glass-card` |
| CSS variables | `--ms-` prefix | `--ms-teal-light` |
| Files | PascalCase | `LiquidGlassCard.tsx` |

### 12.3 Component Template
```tsx
import React from 'react'
import { cn } from '@/lib/cn'

export interface ComponentNameProps {
  className?: string
}

export const ComponentName: React.FC<ComponentNameProps> = ({ className, ...props }) => {
  return (
    <div className={cn('liquid-glass', className)}>
      {/* Content */}
    </div>
  )
}

ComponentName.displayName = 'ComponentName'
```

### 12.4 Component Props Reference

**LiquidGlassCard**
```tsx
import { LiquidGlassCard, GlassBrandCard, GlassTealCard } from '@/components/ui'

<LiquidGlassCard intensity="medium" elevation="raised" hover="lift-glow" brand="teal" interactive>
  Content
</LiquidGlassCard>
```
*Props:* `intensity` (`subtle|light|medium|heavy`) · `elevation` (`none|base|raised|floating`) · `hover` (`none|lift|glow|glow-blue|glow-green|shimmer|lift-glow`) · `brand` (`none|blue|teal|green|brand`) · `radius` (`sm|md|lg|xl|2xl|full`) · `interactive` · `pulseGlow` · `float` · `gradientOverlay`

*Presets:* `GlassCard` · `GlassHeader` · `GlassModal` · `GlassInteractiveCard` · `GlassBrandCard` · `GlassBlueCard` · `GlassTealCard` · `GlassGreenCard`

**LiquidGlassButton**
```tsx
import { LiquidGlassButton, ButtonPrimary, IconButton } from '@/components/ui'

<ButtonPrimary icon={<PlusIcon />} isLoading={isLoading}>Create New</ButtonPrimary>
<IconButton icon={<SearchIcon />} onClick={handleSearch} />
```
*Props:* `variant` (`glass|primary|secondary|ghost|danger`) · `size` (`xs|sm|md|lg|xl`) · `radius` · `hover` (`none|lift|glow|scale`) · `icon` · `iconPosition` (`left|right|only`) · `isLoading` · `disabled`

*Presets:* `ButtonPrimary` · `ButtonSecondary` · `ButtonGhost` · `ButtonDanger` · `IconButton`

**LiquidGlassInput**
```tsx
import { LiquidGlassInput, LiquidGlassSearch, LiquidGlassTextarea } from '@/components/ui'

<LiquidGlassInput label="Email" placeholder="you@example.com" error={err} prefixIcon={<MailIcon />} />
<LiquidGlassSearch placeholder="Search..." value={q} onChange={setQ} onClear={() => setQ('')} />
<LiquidGlassTextarea label="Notes" maxLength={500} showCount rows={4} />
```
*Props:* `size` (`sm|md|lg`) · `state` (`default|error|success|warning`) · `label` · `error` · `helperText` · `prefixIcon` · `suffixIcon` · `showCount` · `maxLength` · `isLoading` · `disabled`

### 12.5 CSS Utility Classes
```css
.liquid-glass                     /* Base glass */
.liquid-glass-subtle|light|heavy  /* Intensity */
.liquid-shadow-ambient|elevation|float  /* Shadows */
.liquid-glass-hover-lift|glow|glow-blue|glow-green|shimmer /* Hover */
.liquid-glass-blue|teal|green|brand     /* Brand variants */
.liquid-radius-sm|md|lg|xl|2xl|full     /* Radius */
.liquid-blur-xs|sm|md|lg|xl             /* Blur */
.liquid-text-primary|secondary|tertiary|inverted|blue|teal|green
.liquid-animate-in  .liquid-pulse-glow  .liquid-float
.liquid-delay-100|200|300|400|500
.liquid-glass-scroll  .liquid-glass-divider  .liquid-skeleton  .liquid-spinner
```

### 12.6 Browser Support
| Browser | Min Version |
| :--- | :--- |
| Chrome / Edge | 90+ |
| Safari | 14+ |
| Firefox | 88+ |
| Mobile Safari (iOS) | 14+ |
| Chrome Android | 90+ |

> Graceful degradation: falls back to solid backgrounds, maintaining accessibility.

---

## 13. AI Accountant Module: Dashboard & Real-Time Tally Integration

This section defines the design patterns and UI components specific to the AI Accountant cloud dashboard, which powers real-time financial data synchronization, automated transaction mapping, and intelligent reconciliation workflows.

### 13.1 Core Design Philosophy for Accounting Interface
The AI Accountant module adopts the same glassmorphism and precision-first design language while introducing financial-specific patterns:

* **Real-Time Status Indicators:** Pulsing sync badges, connection status monitors, and refresh timestamps.
* **Dual-Pane Workflows:** Left sidebar for document management/queues, center for detail views, right for actions/summaries.
* **Color-Coded Transactions:** Success (Green), Pending (Amber), Error (Red), Under Review (Blue) state indicators.
* **Hierarchical Data Views:** Drill-down from summary dashboards → ledgers → transaction details → supporting documents.

### 13.2 Real-Time Synchronization UI

**Sync Status Badge**
*   **Position:** Top-right of dashboard header.
*   **States:**
    *   **Syncing:** Animated pulse, `cyan` glow (`#00E8C6`), text "Syncing..." (12px/600).
    *   **Synced:** Static checkmark, green background (`#34C759`), text "Last sync: 2 mins ago" (12px/400).
    *   **Error:** Red icon (`#FF3B30`), text "Sync failed. Click to retry." (12px/600).
*   **Animation:** Pulse scale 1→1.15, opacity 0.6→1, 2s infinite (syncing state only).
*   **Interactive:** Clicking opens history drawer showing last 20 sync events with timestamps and status.

**Connection Indicator (Tally Connector)**
*   **Visual:** Small dot (8px) beside company/instance name.
*   **Colors:**
    *   **Connected (Green):** `#34C759` steady.
    *   **Disconnected (Red):** `#FF3B30` blinking.
    *   **Connecting (Blue):** `#007AFF` animated rotation of loading spinner.
*   **Hover Tooltip:** Shows API endpoint, last activity timestamp, connection uptime, and manual sync button.

### 13.3 Dashboard Overview (Main Analytics View)

**Header Structure**
*   **Company/Period Selector:** Dropdown to switch between company instances (for multi-entity firms) and date ranges (Monthly, YTD, Custom).
*   **KPI Metric Cards** (2×2 or 3×3 grid, glassmorphic):
    *   **Total Receivable:** Large number with trend arrow (↑/↓), contextual color (green if <30 days, red if >60 days).
    *   **Total Payable:** Similar treatment.
    *   **Cash Position:** Live ticker showing current bank balance with minute-level updates.
    *   **Reconciliation Status:** "X out of Y invoices reconciled" with progress bar and "Review Pending" count badge.

**Chart Widgets**
*   **Income Trend Chart:** Multi-line graph showing revenue vs. expenses over time (last 12 months). Interactive legend to toggling lines.
*   **Expense Breakdown Pie/Donut:** Top 5 expense categories by amount. Click slice to drill down to ledger transactions.
*   **Receivables Aging Bar Chart:** Horizontal stacked bar showing invoices at 0-30, 31-60, 61-90, 90+ days. Color-coded (green→amber→red).
*   **Bank Reconciliation Progress:** Circular progress ring showing % of transactions matched. Center displays "245 of 300 matched."

**Action Queue (Prioritized Tasks)**
*   **Container:** Glass card at width 100% or right sidebar (responsive).
*   **Sections:**
    *   **"Bills Awaiting Upload":** Count badge, "Upload Bills" CTA button.
    *   **"Transactions Pending Review":** Count badge, link to review queue.
    *   **"Reconciliation Items":** Count badge, link to reconciliation workflow.
*   **Design:** Each section is a 2-line item (title, count), stacked vertically, clickable to navigate.

### 13.4 Bill & Statement Upload Interface

**Upload Zone (Drag-and-Drop)**
*   **Container:** 320px × 200px glassmorphic card with centered dashed border (2px, `rgba(255,255,255,0.2)`).
*   **Icon:** Large document icon (48px, `rgba(255,255,255,0.5)`).
*   **Text:**
    *   "Drag PDFs, images, or Excel files here" (16px/500, `--text-primary`).
    *   "or browse" (hyperlinkstyle, `System Blue`).
*   **Supported Formats Badge:** Below upload zone: "PDF, PNG, JPG, Excel, CSV" (12px/400, `--text-tertiary`).
*   **Upload Progress:** File name, upload progress bar (linear), and upload percentage (right-aligned).

**Bulk Upload List View**
*   **Columns:** Filename | Pages | Vendor (auto-detected) | Status (badge) | Action.
*   **Statuses:**
    *   **Processing:** Blue spinning icon, "Extracting details..." (12px/400) text.
    *   **Uploading:** Green icon, "Ready to sync" with checkmark.
    *   **Error:** Red icon, "OCR failed" with Retry button.
*   **Expandable Rows:** Clicking a file shows extracted details preview:
    *   Invoice Amount (large, bold).
    *   Vendor Name (with option to manually correct).
    *   Invoice Date.
    *   GL Ledger (auto-mapped with confidence %).
    *   "Confirm" / "Edit & Confirm" buttons.

**Batch Actions**
*   **Toolbar above list:**
    *   **"Select All"** checkbox.
    *   **Action buttons (disabled until items selected):**
        *   "Confirm All" (blue CTA).
        *   "Download Details" (secondary, exports CSV).
        *   "Delete Selected" (destructive, red, requires confirmation modal).

### 13.5 Transaction Mapping & Reconciliation

**Mapping Review Card**
*   **Layout:** 3-column view (responsive stacking on mobile):
    *   **Col 1 - Transaction Details:** Date, amount, description, bank account.
    *   **Col 2 - AI Suggested Mapping:** Ledger name, sub-ledger (if applicable), cost center, mapped confidence score (70%, 85%, 95%+ badges).
    *   **Col 3 - Actions:** "Approve" (green CTA), "Edit Mapping" (secondary), "Mark as Review" (tertiary).

*   **Confidence Badge:**
    *   **95%+:** Golden badge "High Confidence."
    *   **70–94%:** Amber badge "Review Suggested."
    *   **<70%:** Red badge "Manual Review Required."

**Edit Mapping Modal**
*   **Structure:** Overlay sheet with year/month selector, ledger dropdown (searchable), sub-ledger dropdown (conditional), cost center selector, notes field.
*   **Ledger Dropdown:** Auto-complete style with recent selections pinned at top, grouped by category.
*   **Confirmation:** "Save Mapping" blue button + "Cancel" secondary button.
*   **Side Note:** "This correction will apply only to this transaction. Use 'Bulk Update Rules' to apply to similar future transactions."

**Reconciliation Dashboard (Bank-to-Tally Matching)**
*   **Summary Metrics:**
    *   "Outstanding Payments: $X,XXX"
    *   "Outstanding Receipts: $X,XXX"
    *   "Matched Transactions: X Y Z"
*   **Match List:**
    *   Columns: **Statement Date | Bank Description | Tally Match | Amount | Status**
    *   **Status Badges:**
        *   "✓ Matched" (green).
        *   "⊘ Unmatched" (red).
        *   "⟳ Pending Review" (amber).
    *   **Interaction:** Clicking "Unmatched" opens a modal to manually select a matching Tally transaction or create a new entry.

### 13.6 Document Management & OCR

**Document Library View**
*   **Left Sidebar (Collapsible):**
    *   Folder tree: Bills | Invoices | Bank Statements | Tax Documents | Other.
    *   Each folder shows count of documents.
    *   Search bar for quick file lookup.

*   **Main Canvas:**
    *   **Grid View (Default):** 4-column grid of document cards (responsive).
    *   **List View (Toggle):** Table with columns: Name | Type | Date Uploaded | Vendor | Status.

**Document Card**
*   **Thumbnail:** 120px × 160px PDF/image preview or file icon overlay.
*   **Metadata Below:**
    *   Filename (12px/600).
    *   Upload date (12px/400, `--text-tertiary`).
    *   Vendor name badge.
    *   "OCR Confidence: 92%" small pill (green if >85%, amber if 70–85%, red if <70%).

**Document Detail View (Modal)**
*   **Split Pane:**
    *   **Left:** Document preview (PDF viewer with zoom/pan).
    *   **Right:** Extracted data panel (JSON-like structure showing parsed fields).
*   **Toolbar:**
    *   Download (`⬇︎`)
    *   Re-process OCR (`⟳`)
    *   Edit Extracted Data (pen icon → inline edit mode)
    *   Link to Transaction (`🔗` → dropdown search for matching transaction)
    *   Delete (`🗑︎` → confirmation)

**OCR Error Handling**
*   If extraction fails, show a card with:
    *   Error icon (red).
    *   "OCR processing failed for this document."
    *   Thumbnail with annotation arrows pointing to difficult regions.
    *   "Try again" button.
    *   "Download original & manually enter" link.

### 13.7 Compliance & Audit Reports

**Report Builder**
*   **Type Selector:** Dropdown menu (P&L Board, Balance Sheet, Cash Flow, GST Report, TDS Report, Audit Trail, Aging Reports, etc.).
*   **Period Selector:** Custom date range picker (start, end dates).
*   **Format/Export Options:** Radio buttons: PDF | Excel | CSV.
*   **CTA:** "Generate Report" (blue).

**Report Output (Glassmorphic Container)**
*   **Header:** Report title, company name, period, timestamp.
*   **Table:** Financial data in structured rows/columns.
*   **Footnotes:** Audit trail (e.g., "Last updated: 2 hrs ago by Accountant Name").
*   **Footer:** Company name, digital signature badge (if DKIM/signed).
*   **Export:** "Download as PDF/Excel" CTA button bottom-right.

**Audit Log Viewer**
*   **Table Columns:** Timestamp | User | Action | Data Changed (from → to) | IP Address | Device.
*   **Filters (Top):** Date range, user name, action type (dropdown).
*   **Expandable Rows:** Clicking shows full change details in JSON format.
*   **Download Audit Report:** CTA button exports to CSV with legal timestamp.

### 13.8 Real-Time Data Sync Configuration

**Sync Settings Panel** (Settings > Data Integration)
*   **Tally Connection Status:** Shows "Connected to [Company Name] | Last sync: 2 mins ago | Uptime: 99.8% (last 30 days)."
*   **Sync Frequency Dropdown:** Options: Real-time (every minute), Every 5 minutes, Every 15 minutes, Hourly, Manual only.
*   **Data Scope Toggle Switches:**
    *   ☑ Sync Ledgers → Synced: 45/45 ledgers
    *   ☑ Sync Invoices → Synced: 3,240 invoices
    *   ☑ Sync Bank Statements → Synced: 890 statements
    *   ☑ Sync Bills & Receipts → Synced: 2,150 items
*   **Failed Sync History:** Expandable section showing last 10 failed syncs with error messages and "Retry" button.
*   **Manual Sync CTA:** "Sync Now" button with loading state (spinner, "Syncing...").

**Webhook Configuration**
*   **Incoming Webhooks:** Table showing:
    *   Event Type (dropdown options: New Invoice, Payment Received, Bill Posted, etc.)
    *   Endpoint URL (text input)
    *   Status (Active/Inactive toggle)
    *   Last Triggered (timestamp)
    *   Delete button
*   **Add Webhook:** CTA button opens form to create new webhook.

### 13.9 Color & Badge System for Financial Data

**Transaction/Invoice Status Badges**
*   **"Synced"** → Green (`#34C759`) background, white text, checkmark icon.
*   **"Pending Sync"** → Amber (`#F59E0B`) background, dark text, hourglass icon.
*   **"Under Review"** → Blue (`#007AFF`) background, white text, eye icon.
*   **"Error/Failed"** → Red (`#FF3B30`) background, white text, X icon.
*   **"Archived"** → Slate 400 (`#94A3B8`) background, dark text, archive icon.

**Confidence & Data Quality**
*   **High Confidence (>90%):** Golden badge `#F59E0B` with star icon.
*   **Medium Confidence (70–90%):** Blue badge `#007AFF` with exclamation icon.
*   **Low Confidence (<70%):** Red badge `#FF3B30` with alert icon.

### 13.10 Mobile Optimization for Accountant

**Responsive Breakpoints**
*   **Mobile (<640px):**
    *   Single-column layout; dashboards stack vertically.
    *   Upload zone reduces to 200px × 150px.
    *   Transaction list becomes compact (2 columns: Description + Amount + Status).
    *   Sidebars collapse into bottom tabs (Documents | Transactions | Analytics | Settings).

*   **Tablet (640px–1024px):**
    *   Dual-pane layout; left sidebar pinnable (collapse/expand toggle).
    *   Cards arranged in 2-column grid.

*   **Desktop (>1024px):**
    *   Full triple-pane layout supported.
    *   Cards in 3–4 column grids.

**Touch-Friendly Interactions**
*   Larger touch targets (48px minimum for buttons).
*   Swipe gestures for navigating between sections (left/right swipe for previous/next transaction).
*   Bottom sheet modals for inline actions (edit mapping, confirm upload) instead of overlays.

**Offline Mode Indicator**
*   Top banner (amber/gold): "Offline Mode — Changes will sync when connection is restored."
*   Disable all sync-dependent features (e.g., "Sync Now" button grayed out).
*   Local-first caching ensures previously loaded data remains visible.

### 13.11 Animation & Transitions for Financial Data

**Real-Time Number Updates**
*   When a metric (e.g., "Total Receivable") updates live, the number briefly flashes (`opacity: 1 → 0.7 → 1`, 0.3s).
*   Color transition: Current color → accent color (Growth Teal), then back to original (0.6s ease-in-out).

**Progress Indicators**
*   Sync progress: Linear bar extends left to right, color-coded (blue → green when complete).
*   Upload progress: Circular progress ring (SVG-based) with central percentage text.

**Notification Toast Animations**
*   Entry: `slideInUp` (from bottom, 0.4s cubic-bezier).
*   Exit: `slideOutDown` (downward, 0.3s cubic-bezier).
*   Duration: 5s auto-dismiss for success/info, 10s for errors (or manual close).

**List Item Transitions**
*   When transactions load or new items appear: `fadeInLeft` (0.3s, staggered 50ms per item).
*   When item is selected: `scale(1.02)`, shadow deepens, border highlights.

---

## 14. Design System Maintenance & Evolution

### 18.1 Version Control & Changelogs
*   **Current Version:** 1.0 (Released Feb 2026)
*   **Update Frequency:** Quarterly design audit, with hot-fixes for critical UI bugs released ad-hoc.
*   **Change Log:** Maintain a CHANGELOG.md file tracking:
    *   New components added or deprecated.
    *   Color/typography updates.
    *   Breaking changes (e.g., spacing token value changes).
    *   Migration guides for designers/developers.

### 18.2 Consistency Audits
*   **Monthly:** Review new components/screens against the design system. Ensure all glassmorphism effects, spacing, and typography adhere to defined tokens.
*   **Quarterly:** Usability testing with actual users (doctors, accountants, pharmacists) to validate assumptions and gather feedback.
*   **Annually:** Comprehensive accessibility audit (WCAG 3.0 Bronze compliance, including Guideline 3.3 Input Assistance) and performance review.

### 18.3 Component Library & Storybook
*   Maintain a Storybook instance documenting all UI components with:
    *   Visual snapshots.
    *   Code examples (React, Vue, etc.).
    *   Accessibility notes (ARIA attributes, keyboard navigation).
    *   Common use cases and anti-patterns.
    *   Real-time preview with all design tokens applied.

### 18.4 Platform Standards
*   **Web:** React.js + Tailwind CSS + Apache ECharts.
*   **Mobile:** Flutter + Custom Slivers + Apache ECharts (SVG).
*   **AI Logic:** Genkit (Go/TS) with structured JSON outputs mapping to design system tokens.

### 18.5 Accessibility Standards
*   **Target:** WCAG 3.0 Bronze compliance (see §11 for full details, including Guideline 3.3 Input Assistance).
*   **Key Practices:**
    *   All interactive elements have visible focus indicators (3px accent ring).
    *   Color is never the only means of conveying information (use text labels + icons).
    *   All images and icons have alt-text or aria-labels.
    *   Modal/drawer modals trap focus within the container.
    *   Form errors are announced to screen readers.
    *   Keyboard navigation is fully supported (Tab, Shift+Tab, Enter, Escape).

### 18.6 Performance Optimization
*   **CSS Animations:** Prefer `transform` and `opacity` changes (GPU-accelerated) over `width`, `height`, or `left`/`right` changes.
*   **Glassmorphism:** Use `backdrop-filter: blur()` judiciously; test on lower-end devices (iOS 13+, Android 9+).
*   **Icon Loading:** Lazy-load icon libraries; prefer SVG inline or symbol references.
*   **Image Optimization:** Use WebP with fallbacks (JPG/PNG); apply responsive `srcset` for different screen sizes.

### 18.7 Design System In-Product Documentation
*   **Quick Reference Widget:** A small "?" icon in the app's header links users to:
    *   Glossary of terms (e.g., "What is a Daily Queue?").
    *   Keyboard shortcuts guide.
    *   Feature walkthroughs (video or step-by-step).
    *   Contact support (chat or email).

---

## 15. Testing & Quality Assurance

### 15.1 Visual Regression Testing
*   **Tools:** Percy, Chromatic, or BackstopJS for pixel-level screenshot comparisons.
*   **Workflow:** Every PR automatically captures screenshots on key pages/components. Team reviews visual diffs for unintended changes.
*   **Coverage:** Dashboard, calendar view, modals, mobile breakpoints.

### 15.2 Cross-Browser & Device Testing
*   **Desktop Browsers:** Chrome (latest), Safari (latest), Firefox (latest), Edge (latest).
*   **Mobile:** iOS 14+, Android 9+. Test glassmorphism effects and touch interactions.
*   **Devices:** iPhone 12, iPhone 15 Pro, iPad (gen 5+), Samsung Galaxy A12, Galaxy S23, Pixel 6 (example baseline).
*   **Tools:** BrowserStack, local device testing, or Chromatic.

### 15.3 Interaction Testing
*   **Touch Gestures:** Tap, long-press, swipe, pinch-zoom on mobile.
*   **Hover States:** Ensure all interactive elements have clear hover feedback (background change, shadow, scale).
*   **Keyboard Navigation:** Tab through all interactive elements; ensure logical tab order; test Escape key closing modals.
*   **Animation Smoothness:** Verify animations run at 60 FPS with no janky transitions.

### 15.4 Accessibility Testing
*   **Automated Tools:** Axe, Lighthouse, WAVE to flag common issues (missing alt-text, color contrast, missing labels).
*   **Manual Testing:** Screen reader testing with NVDA (Windows) or VoiceOver (macOS/iOS). Test with actual users who rely on assistive technologies.
*   **Keyboard-Only Navigation:** Ensure all features are accessible without a mouse.

### 15.5 Performance Testing
*   **Lighthouse Audits:** Target scores of 90+ on Performance and Accessibility.
*   **Core Web Vitals:** LCP <2.5s, FID <100ms, CLS <0.1.
*   **Load Time:** Dashboard initial load <3s on 4G connection.
*   **Animation Performance:** Use Timeline / DevTools to ensure 60 FPS (no dropped frames).

---

## 16. Implementation Handoff to Development

### 16.1 Design Assets & Exports
*   **Figma File:** Single source of truth for all designs. Organized with clear naming conventions, locked components, and comprehensive annotations.
*   **Exports:**
    *   **CSS Variables File:** All colors, fonts, spacing, shadows as CSS custom properties.
    *   **Icon SVG Sprite:** All 18px, 22px, 28px icon sizes in a single sprite or folder.
    *   **Gradient Definitions:** Pre-written CSS gradients for common use cases.

### 16.2 Developer Handoff Documentation
*   Create a detailed **IMPLEMENTATION_GUIDE.md** for developers including:
    *   Component naming conventions (e.g., `glass-card`, `button-primary`).
    *   How to apply design tokens in code (CSS custom properties, design tokens JSON, etc.).
    *   Common patterns (card layouts, form inputs, modals) with code snippets.
    *   Spacing grid guidelines with practical examples.
    *   Animation/transition implementation examples.

### 16.3 Ongoing Collaboration
*   **Weekly Design-Dev Sync:** 15-min standup to discuss blockers, new feature specs, and design system updates.
*   **Design Review:** Every feature implementation review with designer present to ensure fidelity to design.
*   **Feedback Loop:** Developers flag design edge cases or technical constraints early. Designers iterate quickly.

---

## 17. Key Takeaways

**MedMentor AI + AI Accountant Design System Summary:**

1. **Glassmorphism + Medical Trust = Modern Healthcare UX:** Layered translucency, living backgrounds, and glass shine effects create an immersive, premium feel while maintaining clinical authority.

2. **Precision at Every Scale:** From 4px grid spacing to 6-second animation timings, every design decision is intentional and measurable.

3. **Financial Data Visualization:** The AI Accountant module adds accounting-specific patterns (sync indicators, transaction mapping, reconciliation workflows) while retaining the core visual language.

4. **Real-Time Mindset:** Indicators, animations, and notifications are designed for live data updates without overwhelming the user.

5. **Accessibility & Inclusivity:** WCAG 3.0 Bronze compliance (including Guideline 3.3 Input Assistance — error identification, error suggestions, error prevention for financial data, redundant entry elimination, and accessible authentication) ensures the platform is usable by all users, including those with visual, motor, or cognitive impairments. iOS 26 Liquid Glass adaptations (Tinted Mode, reduced transparency, reduced motion) provide additional accessibility fallbacks.

6. **Scalability & Maintainability:** A living design system with clear documentation, component libraries, and testing protocols ensures consistency as the product grows.

---

## 18. CSS Custom Properties Reference
Complete token map for both light and glass/dark mode implementations:

### 13.1 Glass Mode Variables
```css
:root {
  /* Glass surfaces */
  --glass-bg: rgba(255, 255, 255, 0.12);
  --glass-bg-hover: rgba(255, 255, 255, 0.18);
  --glass-border: rgba(255, 255, 255, 0.2);
  --glass-border-strong: rgba(255, 255, 255, 0.35);
  --glass-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  --glass-shadow-elevated: 0 20px 60px rgba(0, 0, 0, 0.4);
  --glass-blur: blur(40px);
  --glass-blur-heavy: blur(60px);

    /* Accent colors */
    --accent: #007AFF;
    --accent-light: #5AC8FA;
    --accent-gradient: linear-gradient(135deg, #0056D2, #0F766E);      /* WCAG-safe: white text OK (5.47:1+) */
    --accent-glow: linear-gradient(135deg, #0056D2, #00E8C6);          /* Decorative only: glows, orbs, borders */

    /* Text hierarchy */
    --text-primary: rgba(255, 255, 255, 0.95);
    --text-secondary: rgba(255, 255, 255, 0.6);
    --text-tertiary: rgba(255, 255, 255, 0.35);

    /* Semantic */
    --danger: #FF3B30;
    --success: #34C759;
    --warning: #FF9500;
    --pink: #FF2D55;
    --teal: #5AC8FA;
    --teal-accent: #00E8C6;           /* Decorative glow only — FAILS WCAG with white text */
    --teal-safe: #0F766E;             /* WCAG-safe teal for text-bearing surfaces (5.47:1) */

    /* Radii */
    --radius-sm: 12px;
    --radius-md: 16px;
    --radius-lg: 22px;
    --radius-xl: 28px;
}
```

### 13.2 Light Mode Variables
```css
:root[data-theme='light'] {
  --bg-page: #F1F5F9;
  --bg-card: #FFFFFF;
  --bg-input: #F8FAFC;
  --border-default: #E2E8F0;
  --border-focus: #0056D2;
  --text-primary: #0F172A;
  --text-secondary: #334155;
  --text-tertiary: #64748B;
    --accent: #0056D2;
    --accent-light: #00E8C6;
}
```

---

## 19. Implementation Checklist
Use this checklist when building new screens or components to ensure design system compliance.

### 14.1 Core Checks
- [ ] Font loaded: Inter (400–800) + Plus Jakarta Sans (700, 800) + Cairo (400–700) + Noto Sans Arabic
- [ ] CSS custom properties defined in `:root` for both light and glass themes
- [ ] All spacing on 4px grid — no arbitrary pixel values
- [ ] Touch targets ≥ 44px on all interactive elements (WCAG 2.5.8)
- [ ] Color contrast passes WCAG 3.0 Bronze (4.5:1 text · 3:1 large text/UI · 3:1 non-text)
- [ ] Glassmorphism includes `-webkit-backdrop-filter` prefix for Safari
- [ ] High contrast mode (`prefers-contrast: high`) tested
- [ ] Animations respect `prefers-reduced-motion`
- [ ] Glass surfaces respect `prefers-reduced-transparency` (iOS 26)
- [ ] iOS 26.1 Tinted Mode fallback tested (`.liquid-glass-tinted`)
- [ ] Focus rings visible on all buttons, inputs, and interactive elements (3:1 contrast)
- [ ] Dropdowns use custom glass-elevated panels (not native selects in immersive contexts)
- [ ] Modals use bottom-sheet pattern on mobile with drag handle
- [ ] Loading, error, and empty states defined for each component
- [ ] Dark mode supported
- [ ] RTL compatible (logical CSS properties throughout)
- [ ] Visual regression test in both themes before PR merge
- [ ] No console errors · 60fps animations confirmed

### 14.2 WCAG 3.0 Guideline 3.3 Input Assistance Checks
- [ ] All form fields have visible labels (SC 3.3.2)
- [ ] Required fields marked with `*` + `aria-required="true"` (SC 3.3.2)
- [ ] Error messages are text-based, specific, and adjacent to the field (SC 3.3.1)
- [ ] Erroneous fields identified by border color + icon (not color alone) (SC 3.3.1)
- [ ] Error messages include correction suggestions where feasible (SC 3.3.3)
- [ ] Financial/legal/data-modifying forms include review & confirm step (SC 3.3.4)
- [ ] Multi-step flows carry forward previously entered data (SC 3.3.7)
- [ ] Login fields support password manager paste & autofill (SC 3.3.8)
- [ ] Alternative authentication (magic link / biometric / SSO) available (SC 3.3.8)
- [ ] No `autocomplete="off"` on authentication fields (SC 3.3.8)

### 14.3 RTL / i18n Checks (append for every new screen)
- [ ] Arabic font loaded: Cairo (400–700) + Noto Sans Arabic fallback
- [ ] No hardcoded `left`/`right` in CSS — all layout uses logical properties
- [ ] `dir="auto"` on user-generated text fields (chat input, notes, search)
- [ ] All directional icons have `rtl:rotate-180` or `rtl:scale-x-[-1]` applied
- [ ] Chart configuration includes `isRTL` branch (see §19.6)
- [ ] Arabic translation keys added to all `ar/` namespace JSON files
- [ ] Playwright RTL visual regression test added before PR merge

---

## 20. Internationalisation & RTL Design System

**Cross-ref:** [docs/i18n-architecture.md](../i18n-architecture.md) | [PRD §6.10](../PRD.md)

The MediSync design system natively supports Arabic (RTL) alongside English (LTR). This section defines the visual and interaction design principles that govern the bilingual experience.

---

### 19.1 Core RTL Principle — Mirror, Don't Translate

Layout is **fluid** — it mirrors around the vertical axis when the locale changes. The conceptual model:

```
LTR (English)                      RTL (Arabic)
┌──────────────────────────┐       ┌──────────────────────────┐
│  [Logo]  Nav  Nav  Nav   │       │   Nav  Nav  Nav  [Logo]  │
├─────────┬────────────────┤       ├────────────────┬─────────┤
│         │                │       │                │         │
│ Sidebar │   Content      │       │   Content      │ Sidebar │
│  (left) │                │       │                │ (right) │
├─────────┴────────────────┤       ├────────────────┴─────────┤
│ ← Back   Forward →       │       │       → Back   ← Forward │
└──────────────────────────┘       └──────────────────────────┘
```

**Elements that mirror:**
- Navigation sidebar: left → right
- Page/section flow: left-to-right → right-to-left
- Breadcrumbs: `Home › Settings › Profile` → `الرئيسية ‹ الإعدادات ‹ الملف الشخصي`
- Icons with inherent direction (arrows, chevrons, play/forward)
- Progress bars and loading indicators
- Chat bubble alignment and tail direction

**Elements that do NOT mirror:**
- The MediSync logo and brand mark
- Clocks and circular gauges
- Icons with no direction (star, bell, search, user)
- Data visualisation content (chart bars stay visually comparative; only axis labels translate)
- Mathematical operators and formulae

---

### 19.2 Typography Adjustments for Arabic

| Property | English (LTR) | Arabic (RTL) | Rationale |
|----------|--------------|-------------|-----------|
| Font family | Inter | Cairo / Noto Sans Arabic | Arabic requires dedicated Arabic-designed typefaces |
| Line height (body) | `1.5` | `1.8` | Arabic letterforms have more vertical complexity |
| Line height (display) | `1.1–1.2` | `1.4` | Descenders and ascenders need more breathing room |
| Letter spacing | `0 – -0.02em` | `0` (always) | Arabic typography must NOT have tracking adjustments |
| Font weight (body) | 400 | 400–500 | Cairo Regular renders slightly lighter; Medium improves legibility |
| Text alignment | `start` (left) | `start` (right) | Use logical `text-align: start` — never `left`/`right` |

**CSS Implementation:**
```css
/* Activate Arabic typography */
:lang(ar) {
  font-family: 'Cairo', 'Noto Sans Arabic', sans-serif;
  line-height: 1.8;
  letter-spacing: 0; /* Always reset — Arabic has no tracking */
}

/* Display headings in Arabic */
:lang(ar) h1, :lang(ar) h2 {
  font-weight: 700;  /* Cairo Bold */
  line-height: 1.4;
}

/* Body text */
:lang(ar) p, :lang(ar) li, :lang(ar) td {
  font-weight: 500;  /* Cairo Medium for readability */
  line-height: 1.8;
}
```

---

### 19.3 Spacing & Layout Adjustments for RTL

All components use **CSS Logical Properties**. No physical `left`/`right` in new code.

```css
/* ❌ Never:
   margin-left, padding-right, border-left, right: 0 */

/* ✅ Always:
   margin-inline-start, padding-inline-end, border-inline-start, inset-inline-end: 0 */

/* Sidebar — positions itself correctly in both LTR and RTL */
.sidebar {
  position: fixed;
  inset-block: 0;
  inset-inline-start: 0;   /* left in LTR, right in RTL */
  width: 256px;
  border-inline-end: 1px solid var(--glass-border);
}

/* Icon + label pair */
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;       /* gap is direction-agnostic */
  padding-inline: 16px;
}
```

**Tailwind utilities for RTL:**
```html
<!-- Breadcrumb chevron flips automatically -->
<svg class="rtl:rotate-180 transition-transform" .../>

<!-- Dropdown opens from correct side -->
<div class="absolute inset-inline-end-0 top-full ...">

<!-- Chat bubble —user —mirrors -->
<div class="ms-auto ltr:rounded-br-none rtl:rounded-bl-none ...">
```

---

### 19.4 Chat Interface RTL Design

```
LTR (English)                         RTL (Arabic)
┌──────────────────────────────┐      ┌──────────────────────────────┐
│ 🤖 AI bubble                 │      │                 bubble AI 🤖 │
│ ╰─ tail left                 │      │             tail right ─╯    │
│                              │      │                              │
│           User bubble 👤     │      │     👤 bubble المستخدم       │
│               tail right ─╯  │      │ ╰─ tail left                 │
└──────────────────────────────┘      └──────────────────────────────┘
```

- AI message bubbles: `inset-inline-start` aligned, tail on `inline-start` side
- User message bubbles: `inset-inline-end` aligned (pushed to opposite edge of AI)
- Chat input: `dir="auto"` — cursor and caret move to right-edge when user types Arabic
- Send button: always at `inline-end` of input (right in LTR → left in RTL)
- Pre-defined prompt buttons: flow wraps naturally with logical margin

---

### 19.5 Data Tables & Reports RTL

| Element | LTR | RTL |
|---------|-----|-----|
| Column order | First column at left | First column at right |
| Numeric cells | Right-aligned | Right-aligned (stays — numerics always right-to-left) |
| Text cells | Left-aligned | Right-aligned |
| Sort icon | Appears after header text | Appears before header text |
| Expand / drill-down `›` | Right side of row | Left side of row |
| Pagination `‹ 1 2 3 ›` | Left buttons = back | Right buttons = back |

Note: **Right-aligning numbers is identical in LTR and RTL** — use `text-align: end` which resolves to right in both directions for numeric columns (since numbers read right-to-left regardless of document direction).

---

### 19.6 Charts & Visualisations (Apache ECharts) RTL

```js
// echarts locale-aware configuration
const isRTL = locale === 'ar';

const option = {
  textStyle: {
    fontFamily: isRTL ? "'Cairo', 'Noto Sans Arabic', sans-serif" : "'Inter', sans-serif",
  },
  // X-axis: categories read right→left in Arabic
  xAxis: {
    inverse: isRTL,          // flip category order
    axisLabel: {
      align: isRTL ? 'right' : 'left',
      formatter: (val) => fmtDate.format(new Date(val)),
    },
  },
  // Legend: flip to opposite side
  legend: {
    left:  isRTL ? 'auto' : 10,
    right: isRTL ? 10 : 'auto',
  },
  // Tooltip: open to the left in RTL to avoid overflow
  tooltip: {
    confine: true,
    position: isRTL ? 'left' : 'right',
  },
};
```

**Note:** Bar chart bars, line trends, and data points are **not mirrored** — only the axis labelling and legend direction change. Time-series charts always progress left-to-right chronologically for both locales (most-recent right) as financial data is universally read this way.

---

### 19.7 Language Switcher Component

A persistent control appears in the top navigation bar and on the Profile Settings screen.

**Design spec:**
```
[ EN | ع ]   ← compact pill toggle in nav bar (always visible)
```

- Pill shows ISO codes: `EN` / `ع` (Arabic abbreviation `ع`)
- Active locale has Trust Blue background, white text
- Inactive locale has ghost style
- Switching locale triggers instant layout flip — no page reload
- The language switcher itself is **exempt from RTL mirroring** — it always appears at the same edge to avoid confusion during language transitions
- Accessible: `role="radiogroup"`, `aria-label="Language"`, `aria-checked` on active lang

**Flutter bottom sheet equivalent for mobile:**
```dart
showModalBottomSheet(
  context: context,
  builder: (_) => LanguageSelectionSheet(
    options: [Locale('en'), Locale('ar')],
    current: AppLocalizations.of(context)!.localeName,
    onSelect: (locale) => context.read<LocaleNotifier>().setLocale(locale),
  ),
);
```

---

### 19.8 Glassmorphism in RTL Context

The glassmorphism system requires no colour changes for RTL. Only the glass shine animation direction is adjusted:

```css
/* Glass shine sweeps in the reading direction */
:lang(ar) .glass::after {
  background: linear-gradient(
    -135deg,        /* Reverse direction for RTL: 135deg → -135deg */
    transparent 40%,
    rgba(255,255,255,0.04) 45%,
    rgba(255,255,255,0.08) 50%,
    rgba(255,255,255,0.04) 55%,
    transparent 60%
  );
}
```

---

### 19.9 Implementation Checklist — RTL / i18n Addition

Append these checks to the existing Implementation Checklist (§14) for every new screen or component:

11. Arabic font loaded: Cairo (400, 500, 600, 700) + Noto Sans Arabic fallback.
12. No hardcoded `left`/`right` in CSS — all layout uses logical properties.
13. `dir="auto"` on user-generated text fields (chat input, notes, search).
14. All directional icons (`→`, `‹`, `›`) have `rtl:rotate-180` or `rtl:scale-x-[-1]` applied.
15. Chart configuration includes `isRTL` branch for `inverse`, `legend`, and `tooltip`.
16. PDF report template loaded with Cairo font and `direction: rtl` for Arabic output.
17. Language switcher component visible in nav bar on every screen.
18. Playwright RTL visual regression test added for new screen before PR merge.
19. Arabic translation keys added to all relevant `ar/` namespace JSON files.
20. `line-height` set to `1.8` minimum for Arabic body text blocks.

---

## 21. Future Spatial Design (v2.0)

*(Replaces former §15)*

As MedMentor moves to "Sherlock Mode" (AR) and Vision Pro, the design system extends beyond 2D:

### 20.1 Spatial Glass (Z-Depth)
*   **Layer 0 (Reality)**: The physical world (textbook, patient model).
*   **Layer 1 (Ambient)**: "Glass" panels floating 1m away (Notes, simple Q&A).
*   **Layer 2 (Focus)**: Interactive 3D models (Hearts, Molecules) instantiated 0.5m away for manipulation.

### 20.2 Ambient Interactions
*   **Gaze-Driven**: UI elements expand slightly when looked at (Eye tracking).
*   **Air Gestures**: "Pinch to extract" text from a physical book into the digital notebook.
*   **Bio-Feedback UI**: The interface "breathes" (subtle pulse animation) to match the user's stress level, guiding them to calmness.

### 20.3 Immersive Learning UI (Student Features)
*   **OSCE Mode**: Minimalist UI. "Voice-first" interaction. The patient is a full-screen video avatar (or 3D model). Design must hide all text prompts to force active recall.
*   **Memory Palace**: World-scale AR. Labels must use **billboarding** (always face user) and maintain legibility against variable real-world backgrounds (dynamic contrast/shadows).
*   **Neuro-Flashcards**: Invisible UI. No "grade yourself" buttons. The card simply *knows* and swipes itself away when mastery is detected.
