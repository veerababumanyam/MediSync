# Hero Carousel Redesign Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Redesign the landing page Hero Carousel with new copy, 3 animated SVG illustrations, capability strip, and enhanced visual hierarchy to increase demo bookings.

**Architecture:** React functional component with i18n translations, CSS animations using design tokens, and SVG illustrations. Follows existing Liquid Glass design system with progressive narrative across 3 slides.

**Tech Stack:** React 19.2, TypeScript, Tailwind CSS 4.2, i18next, CSS keyframe animations, SVG

---

## Prerequisites

- Read the design document: `docs/plans/2026-02-23-hero-carousel-redesign-design.md`
- Existing files to understand:
  - `frontend/src/components/landing/HeroCarousel.tsx`
  - `frontend/src/components/landing/icons.tsx`
  - `frontend/src/i18n/index.ts`
  - `docs/DESIGN.md` (§5 Glassmorphism System)

---

## Task 1: Update i18n Copy for All 3 Slides (English)

**Files:**
- Modify: `frontend/src/i18n/index.ts:51-86`

**Step 1: Update slide 1 copy**

Replace the `heroCarousel.slide1` object in `frontend/src/i18n/index.ts`:

```typescript
slide1: {
  title: "Don't Replace Your Software. Make It Speak.",
  subtitle: 'Any System. Any Database. Zero Migration.',
  description: 'Ask questions in plain language. Get instant answers from HIMS, Tally, SQL, or any legacy system. No rip-and-replace required.',
  cta: 'Start Free Trial',
  stat1: '50+ Integrations',
  stat2: '< 2 Min Setup',
  stat3: 'Zero Code Changes',
},
```

**Step 2: Update slide 2 copy**

Replace the `heroCarousel.slide2` object:

```typescript
slide2: {
  title: 'One Platform. Every Capability.',
  subtitle: 'Conversational BI · AI Accountant · Smart Reports · Deep Analytics',
  description: 'From natural language queries to automated accounting sync — all powered by 58 specialized AI agents working together.',
  cta: 'See It In Action',
  stat1: '58 AI Agents',
  stat2: 'NL → SQL → Charts',
  stat3: 'Auto Ledger Mapping',
},
```

**Step 3: Update slide 3 copy**

Replace the `heroCarousel.slide3` object:

```typescript
slide3: {
  title: 'Built for Healthcare & Finance.',
  subtitle: 'HIPAA Compliant · SOC 2 Certified · 99.9% Uptime',
  description: 'Trusted by clinics, labs, and hospitals across 12 countries. Average savings: ₹2Cr+ annually with full compliance.',
  cta: 'Book a Demo',
  stat1: '₹2Cr+ Avg Savings',
  stat2: '12 Countries',
  stat3: '500+ Clinics',
},
```

**Step 4: Add capabilities object**

Add after `slide3` and before `trustStrip`:

```typescript
capabilities: {
  bi: 'Conversational BI',
  accountant: 'AI Accountant',
  reports: 'Smart Reports',
  analytics: 'Deep Analytics',
},
```

**Step 5: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no errors

**Step 6: Commit**

```bash
git add frontend/src/i18n/index.ts
git commit -m "feat(hero): update carousel copy for conversion optimization

- Add subtitles to all 3 slides
- Refine descriptions for clarity
- Update stats to show platform breadth
- Add capabilities i18n keys

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Add Arabic Translations for New Copy

**Files:**
- Modify: `frontend/src/i18n/index.ts` (Arabic section)

**Step 1: Add Arabic heroCarousel translations**

Find the Arabic resources section (around line 200+) and add/update the `heroCarousel` object:

```typescript
heroCarousel: {
  slide1: {
    title: 'لا تستبدل برنامجك. اجعله يتحدث.',
    subtitle: 'أي نظام. أي قاعدة بيانات. بدون ترحيل.',
    description: 'اسأل بلغة طبيعية واحصل على إجابات فورية من HIMS أو Tally أو SQL أو أي نظام قديم. بدون استبدال.',
    cta: 'ابدأ التجربة المجانية',
    stat1: '+50 تكامل',
    stat2: 'أقل من دقيقتين',
    stat3: 'بدون تغيير код',
  },
  slide2: {
    title: 'منصة واحدة. كل القدرات.',
    subtitle: 'ذكاء اصطناعي محادثي · محاسب ذكي · تقارير ذكية · تحليلات متقدمة',
    description: 'من الاستعلامات بلغة طبيعية إلى مزامنة المحاسبة الآلية — مدعوم بـ 58 وكيل ذكاء اصطناعي متخصص.',
    cta: 'شاهده أثناء العمل',
    stat1: '58 وكيل ذكاء اصطناعي',
    stat2: 'لغة ← SQL ← رسوم',
    stat3: 'تعيين دفتر آلي',
  },
  slide3: {
    title: 'مبني للرعاية الصحية والمالية.',
    subtitle: 'متوافق مع HIPAA · معتمد SOC 2 · 99.9% وقت تشغيل',
    description: 'موثوق به من العيادات والمختبرات والمستشفيات في 12 دولة. متوسط التوفير: +2 كرور روبية سنوياً مع الامتثال الكامل.',
    cta: 'احجز عرضاً توضيحياً',
    stat1: '+2 كرور توفير',
    stat2: '12 دولة',
    stat3: '+500 عيادة',
  },
  capabilities: {
    bi: 'ذكاء اصطناعي محادثي',
    accountant: 'محاسب ذكي',
    reports: 'تقارير ذكية',
    analytics: 'تحليلات متقدمة',
  },
  trustStrip: {
    item1: 'HIMS',
    item2: 'LIMS',
    item3: 'Tally ERP',
    item4: 'SQL',
    item5: 'واجهات برمجة',
    item6: 'Oracle',
    item7: 'SAP',
    item8: 'REST / GraphQL',
  },
},
```

**Step 2: Verify i18n loads correctly**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/i18n/index.ts
git commit -m "feat(i18n): add Arabic translations for hero carousel redesign

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Create New Slide 1 Illustration (Conversation with Legacy)

**Files:**
- Modify: `frontend/src/components/landing/icons.tsx:19-51`

**Step 1: Replace the slide 1 SVG**

Replace the entire `if (slide === 'slide1')` block with:

```tsx
if (slide === 'slide1') {
    // Conversation with Legacy — AI brain connecting to legacy systems
    return (
        <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
            {/* Outer glow ring */}
            <circle cx="160" cy="140" r="120" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 6s ease-in-out infinite' }} />

            {/* Central AI node */}
            <g style={{ animation: 'floatY 4s ease-in-out infinite' }}>
                <circle cx="160" cy="100" r="35" fill={accent} fillOpacity="0.15" stroke={accent} strokeWidth="2" />
                <text x="160" y="105" textAnchor="middle" fill={accent} fontSize="14" fontWeight="700">AI</text>
            </g>

            {/* Connection lines with dash animation */}
            <line x1="160" y1="135" x2="80" y2="190" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />
            <line x1="160" y1="135" x2="160" y2="200" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />
            <line x1="160" y1="135" x2="240" y2="190" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />

            {/* Legacy system nodes */}
            <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite' }}>
                <rect x="45" y="175" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                <text x="80" y="200" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">HIMS</text>
            </g>

            <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite 0.3s' }}>
                <rect x="125" y="195" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                <text x="160" y="220" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">Tally</text>
            </g>

            <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite 0.6s' }}>
                <rect x="205" y="175" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                <text x="240" y="200" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">SQL</text>
            </g>

            {/* Chat bubbles — showing conversation */}
            <g style={{ animation: 'floatY 3.5s ease-in-out infinite' }}>
                <rect x="50" y="50" rx="12" width="80" height="32" fill={nodeColor} stroke={strokeColor} strokeWidth="1" />
                <rect x="62" y="60" width="45" height="4" rx="2" fill={strokeColor} />
                <rect x="62" y="68" width="30" height="4" rx="2" fill={strokeColor} />
            </g>

            <g style={{ animation: 'floatYReverse 3s ease-in-out infinite 0.5s' }}>
                <rect x="190" y="45" rx="12" width="80" height="32" fill={accent} fillOpacity="0.15" stroke={accent} strokeWidth="1" />
                <rect x="202" y="57" width="55" height="4" rx="2" fill={accent} fillOpacity="0.6" />
                <rect x="202" y="65" width="35" height="4" rx="2" fill={accent} fillOpacity="0.4" />
            </g>

            {/* Sparkle accents */}
            <circle cx="55" cy="140" r="3" fill={accent} fillOpacity="0.5" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
            <circle cx="265" cy="130" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.7s' }} />
            <circle cx="160" cy="260" r="2" fill={accent} fillOpacity="0.3" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1.2s' }} />

            {/* Data flow particles */}
            <circle cx="120" cy="165" r="3" fill={accent} fillOpacity="0.6" style={{ animation: 'floatY 2s ease-in-out infinite' }} />
            <circle cx="200" cy="165" r="3" fill={accent} fillOpacity="0.6" style={{ animation: 'floatYReverse 2.3s ease-in-out infinite 0.4s' }} />
        </svg>
    )
}
```

**Step 2: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/components/landing/icons.tsx
git commit -m "feat(hero): redesign slide 1 illustration - conversation with legacy

- Central AI node connecting to HIMS, Tally, SQL
- Chat bubbles showing conversation metaphor
- Data flow particles and sparkle accents
- Staggered float animations

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Create New Slide 2 Illustration (Platform Capabilities Hub)

**Files:**
- Modify: `frontend/src/components/landing/icons.tsx:53-85`

**Step 1: Replace the slide 2 SVG**

Replace the entire `if (slide === 'slide2')` block with:

```tsx
if (slide === 'slide2') {
    // Platform Capabilities — central hub with 4 capability nodes
    return (
        <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
            {/* Outer glow ring */}
            <circle cx="160" cy="140" r="115" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 5s ease-in-out infinite' }} />

            {/* Central hub */}
            <g style={{ animation: 'floatY 5s ease-in-out infinite' }}>
                <circle cx="160" cy="140" r="40" fill={accent} fillOpacity="0.2" stroke={accent} strokeWidth="2.5" />
                <text x="160" y="136" textAnchor="middle" fill={accent} fontSize="10" fontWeight="700">AnySync</text>
                <text x="160" y="150" textAnchor="middle" fill={accent} fontSize="8" fillOpacity="0.7">PLATFORM</text>
            </g>

            {/* Connection lines to nodes */}
            <line x1="160" y1="100" x2="160" y2="60" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
            <line x1="200" y1="140" x2="240" y2="140" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
            <line x1="160" y1="180" x2="160" y2="220" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
            <line x1="120" y1="140" x2="80" y2="140" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />

            {/* Top node — Conversational BI */}
            <g style={{ animation: 'floatYReverse 3.5s ease-in-out infinite' }}>
                <circle cx="160" cy="45" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                <rect x="148" y="36" width="16" height="10" rx="4" fill={accent} fillOpacity="0.6" />
                <circle cx="156" cy="50" r="4" fill={accent} fillOpacity="0.4" />
                <text x="160" y="62" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">BI</text>
            </g>

            {/* Right node — AI Accountant */}
            <g style={{ animation: 'floatY 3.8s ease-in-out infinite 0.3s' }}>
                <circle cx="255" cy="140" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                <rect x="245" y="130" width="14" height="18" rx="3" fill={accent} fillOpacity="0.5" />
                <path d="M248 140 l3 3 l6 -6" stroke={accent} strokeWidth="2" fill="none" />
                <text x="255" y="158" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">ACC</text>
            </g>

            {/* Bottom node — Smart Reports */}
            <g style={{ animation: 'floatYReverse 3.5s ease-in-out infinite 0.6s' }}>
                <circle cx="160" cy="235" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                <rect x="148" y="225" width="6" height="15" rx="2" fill={accent} fillOpacity="0.4" />
                <rect x="157" y="230" width="6" height="10" rx="2" fill={accent} fillOpacity="0.6" />
                <rect x="166" y="222" width="6" height="18" rx="2" fill={accent} fillOpacity="0.8" />
                <text x="160" y="253" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">RPT</text>
            </g>

            {/* Left node — Deep Analytics */}
            <g style={{ animation: 'floatY 3.8s ease-in-out infinite 0.9s' }}>
                <circle cx="65" cy="140" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                <circle cx="65" cy="135" r="10" stroke={accent} strokeWidth="1.5" fill="none" />
                <circle cx="65" cy="135" r="4" fill={accent} fillOpacity="0.5" />
                <circle cx="60" cy="130" r="2" fill={accent} fillOpacity="0.6" />
                <circle cx="72" cy="138" r="2" fill={accent} fillOpacity="0.6" />
                <text x="65" y="158" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">ANL</text>
            </g>

            {/* Data flow particles */}
            <circle cx="160" cy="80" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatY 2s ease-in-out infinite' }} />
            <circle cx="215" cy="140" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatYReverse 2.3s ease-in-out infinite' }} />
            <circle cx="160" cy="200" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatY 2.1s ease-in-out infinite 0.2s' }} />
            <circle cx="105" cy="140" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatYReverse 2.5s ease-in-out infinite 0.4s' }} />

            {/* Sparkle accents */}
            <circle cx="40" cy="60" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
            <circle cx="280" cy="70" r="2" fill={accent} fillOpacity="0.35" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.5s' }} />
            <circle cx="280" cy="220" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1s' }} />
            <circle cx="40" cy="210" r="2" fill={accent} fillOpacity="0.3" style={{ animation: 'pulseGlow 3.2s ease-in-out infinite 0.8s' }} />
        </svg>
    )
}
```

**Step 2: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/components/landing/icons.tsx
git commit -m "feat(hero): redesign slide 2 illustration - capabilities hub

- Central AnySync hub with 4 capability nodes
- Icons for BI, Accountant, Reports, Analytics
- Radial connection lines with data particles
- Staggered float animations for each node

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 5: Create New Slide 3 Illustration (Trust Dashboard)

**Files:**
- Modify: `frontend/src/components/landing/icons.tsx:87-114`

**Step 1: Replace the slide 3 SVG**

Replace the entire slide 3 return block (after the comment) with:

```tsx
// slide3 — Trust Dashboard with outcomes and compliance
return (
    <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
        {/* Outer glow ring */}
        <circle cx="160" cy="140" r="115" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 5s ease-in-out infinite' }} />

        {/* Dashboard frame */}
        <g style={{ animation: 'floatY 5s ease-in-out infinite' }}>
            <rect x="50" y="55" rx="14" width="220" height="165" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />

            {/* Title bar */}
            <rect x="50" y="55" rx="14" width="220" height="26" fill={accent} fillOpacity="0.1" />
            <circle cx="68" cy="68" r="4" fill={errorColor} fillOpacity="0.7" />
            <circle cx="82" cy="68" r="4" fill={warningColor} fillOpacity="0.7" />
            <circle cx="96" cy="68" r="4" fill={successColor} fillOpacity="0.7" />
            <text x="160" y="72" textAnchor="middle" fill={accent} fontSize="9" fontWeight="600" fillOpacity="0.8">AnySync Dashboard</text>
        </g>

        {/* Floating metric cards */}
        <g style={{ animation: 'floatYReverse 3s ease-in-out infinite' }}>
            <rect x="30" y="225" rx="8" width="65" height="30" fill={nodeColor} stroke={accent} strokeWidth="1" />
            <text x="62" y="238" textAnchor="middle" fill={successColor} fontSize="10" fontWeight="700">↑ +127%</text>
            <text x="62" y="250" textAnchor="middle" fill={accent} fontSize="7" fillOpacity="0.7">Growth</text>
        </g>

        <g style={{ animation: 'floatY 3.5s ease-in-out infinite 0.5s' }}>
            <rect x="225" y="35" rx="8" width="65" height="30" fill={nodeColor} stroke={accent} strokeWidth="1" />
            <text x="257" y="48" textAnchor="middle" fill={accent} fontSize="10" fontWeight="700">₹2.1Cr</text>
            <text x="257" y="60" textAnchor="middle" fill={accent} fontSize="7" fillOpacity="0.7">Saved</text>
        </g>

        {/* Bar chart */}
        <g>
            <rect x="70" y="175" width="18" height="30" rx="3" fill={accent} fillOpacity="0.25" />
            <rect x="95" y="160" width="18" height="45" rx="3" fill={accent} fillOpacity="0.35" />
            <rect x="120" y="145" width="18" height="60" rx="3" fill={accent} fillOpacity="0.5" />
            <rect x="145" y="130" width="18" height="75" rx="3" fill={accent} fillOpacity="0.65" />
            <rect x="170" y="115" width="18" height="90" rx="3" fill={accent} fillOpacity="0.8" />
            <rect x="195" y="100" width="18" height="105" rx="3" fill={accent} fillOpacity="0.95" />
            <rect x="220" y="90" width="18" height="115" rx="3" fill={accent} />
        </g>

        {/* Trend line */}
        <polyline points="79,165 104,150 129,135 154,120 179,105 204,92 229,82" stroke={successColor} strokeWidth="2.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />

        {/* Trust badges */}
        <g style={{ animation: 'floatY 4s ease-in-out infinite 0.3s' }}>
            <rect x="105" y="240" rx="6" width="50" height="22" fill={nodeColor} stroke={successColor} strokeWidth="1" />
            <text x="130" y="255" textAnchor="middle" fill={successColor} fontSize="8" fontWeight="600">HIPAA</text>
        </g>

        <g style={{ animation: 'floatYReverse 4s ease-in-out infinite 0.6s' }}>
            <rect x="165" y="240" rx="6" width="50" height="22" fill={nodeColor} stroke={accent} strokeWidth="1" />
            <text x="190" y="255" textAnchor="middle" fill={accent} fontSize="8" fontWeight="600">99.9%</text>
        </g>

        {/* Sparkle accents */}
        <circle cx="35" cy="90" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
        <circle cx="285" cy="170" r="2" fill={accent} fillOpacity="0.35" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.5s' }} />
        <circle cx="160" cy="270" r="2" fill={successColor} fillOpacity="0.5" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1s' }} />
    </svg>
)
```

**Step 2: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/components/landing/icons.tsx
git commit -m "feat(hero): redesign slide 3 illustration - trust dashboard

- Premium dashboard frame with traffic lights
- Growing bar chart with trend line
- Floating metric cards (growth, savings)
- Trust badges (HIPAA, 99.9% uptime)
- Staggered float animations

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 6: Add Subtitle Support and Capability Strip to HeroCarousel

**Files:**
- Modify: `frontend/src/components/landing/HeroCarousel.tsx`

**Step 1: Add subtitle rendering after headline**

Find line ~97 (after the `</h2>` closing tag for the headline) and add:

```tsx
                                        {/* Subtitle / Tagline */}
                                        <p className={`text-sm sm:text-base font-medium tracking-wide mb-4 ${isDark ? 'text-slate-400' : 'text-slate-500'} animate-fade-in-up`}
                                           style={{ animationDelay: '250ms' }}>
                                            {t(`heroCarousel.${item}.subtitle`)}
                                        </p>
```

**Step 2: Add capability strip component (slide 2 only)**

Find the description paragraph (~line 100-102) and add after it:

```tsx
                                        {/* Capability strip - slide 2 only */}
                                        {item === 'slide2' && (
                                            <div className="flex flex-wrap justify-center lg:justify-start gap-3 mb-6 animate-fade-in-up" style={{ animationDelay: '400ms' }}>
                                                {(['bi', 'accountant', 'reports', 'analytics'] as const).map((cap) => (
                                                    <div
                                                        key={cap}
                                                        className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold transition-transform hover:scale-105 ${isDark ? 'bg-white/5 text-slate-300' : 'bg-slate-100 text-slate-600'}`}
                                                    >
                                                        <div className={`w-2 h-2 rounded-full ${isDark ? 'bg-teal-400' : 'bg-blue-500'}`} />
                                                        {t(`heroCarousel.capabilities.${cap}`)}
                                                    </div>
                                                ))}
                                            </div>
                                        )}
```

**Step 3: Update description margin**

Change the description paragraph's `mb-7` to `mb-4` to accommodate the capability strip:

```tsx
<p className={`text-base sm:text-lg leading-relaxed mb-4 max-w-xl mx-auto lg:mx-0 ...
```

**Step 4: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 5: Commit**

```bash
git add frontend/src/components/landing/HeroCarousel.tsx
git commit -m "feat(hero): add subtitle support and capability strip

- Subtitle rendered below headline on all slides
- Capability strip on slide 2 showing 4 modules
- Updated description margin for new layout

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 7: Add Illustration Backdrop Glow and Inner Card Treatment

**Files:**
- Modify: `frontend/src/components/landing/HeroCarousel.tsx`

**Step 1: Add backdrop glow around illustration**

Find the illustration container (~line 134) and wrap it:

```tsx
                                    {/* RIGHT: SVG Illustration - with backdrop glow and floating animation */}
                                    <div className="shrink-0 w-full max-w-[260px] sm:max-w-[300px] lg:max-w-[340px] animate-float-medium relative">
                                        {/* Backdrop glow */}
                                        <div
                                            className="absolute inset-0 -z-10 rounded-full opacity-30"
                                            style={{
                                                background: `radial-gradient(circle, ${isDark ? 'rgba(24, 146, 157, 0.2)' : 'rgba(39, 80, 168, 0.15)'} 0%, transparent 70%)`,
                                                transform: 'scale(1.3)',
                                            }}
                                        />
                                        <HeroIllustration slide={item} isDark={isDark} />
                                    </div>
```

**Step 2: Add inner card subtle border to slide content**

Find the `hero-slide-layout` div and add an inner wrapper:

```tsx
                                <div className="hero-slide-layout px-6 sm:px-10 md:px-14 py-10 sm:py-14 lg:py-16">
                                    {/* Inner content card with subtle separation */}
                                    <div className="relative rounded-2xl p-4 sm:p-6" style={{
                                        boxShadow: isDark
                                            ? 'inset 0 1px 0 0 rgba(255,255,255,0.05)'
                                            : 'inset 0 1px 0 0 rgba(255,255,255,0.5)',
                                    }}>
                                        {/* LEFT: Text content */}
                                        <div className="flex-1 min-w-0 text-center lg:text-left">
```

Make sure to close the inner div before the closing `</div>` of `hero-slide-layout`.

**Step 3: Verify no TypeScript errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 4: Commit**

```bash
git add frontend/src/components/landing/HeroCarousel.tsx
git commit -m "feat(hero): add illustration backdrop glow and inner card treatment

- Radial gradient glow behind illustration
- Inner content card with subtle inset border
- Improved visual hierarchy

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 8: Add Reduced Motion CSS Support

**Files:**
- Modify: `frontend/src/styles/globals.css`

**Step 1: Add reduced motion media query**

Add at the end of `frontend/src/styles/globals.css`:

```css
/* ============================================
 * Hero Carousel - Reduced Motion Support
 * ============================================ */

@media (prefers-reduced-motion: reduce) {
  .hero-carousel-container *,
  .hero-carousel-container *::before,
  .hero-carousel-container *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }

  .hero-carousel-container .animate-float-medium,
  .hero-carousel-container [style*="animation"] {
    animation: none !important;
  }

  /* Keep opacity transitions for accessibility */
  .hero-carousel-container [role="tabpanel"] {
    transition: opacity 0.15s ease !important;
  }
}
```

**Step 2: Verify CSS compiles**

Run: `cd frontend && npm run build`
Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/styles/globals.css
git commit -m "feat(a11y): add reduced motion support for hero carousel

- Disable animations for prefers-reduced-motion users
- Keep minimal opacity transitions for accessibility
- WCAG 2.3.3 compliance

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 9: Verify and Test the Complete Implementation

**Files:**
- None (verification task)

**Step 1: Run full build**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no errors

**Step 2: Run linter**

Run: `cd frontend && npm run lint`
Expected: No linting errors (warnings acceptable)

**Step 3: Start dev server and visually verify**

Run: `cd frontend && npm run dev`
Then open browser to `http://localhost:5173`

Verify:
- [ ] All 3 slides show new copy
- [ ] Subtitles appear below headlines
- [ ] Slide 2 shows capability strip
- [ ] All 3 new illustrations render correctly
- [ ] Animations play smoothly
- [ ] CTA buttons work
- [ ] Trust strip marquee works
- [ ] Progress bar animates
- [ ] Auto-rotate works (pauses on hover)

**Step 4: Test RTL (Arabic)**

In browser dev tools or by adding `?lang=ar` to URL:
- [ ] Text flips to RTL correctly
- [ ] Layout mirrors appropriately
- [ ] Illustrations remain correctly positioned

**Step 5: Test reduced motion**

In browser dev tools, emulate `prefers-reduced-motion: reduce`:
- [ ] All animations stop
- [ ] Content is still readable
- [ ] Transitions are instant

**Step 6: Final commit if any fixes needed**

```bash
git add -A
git commit -m "fix(hero): final polish for carousel redesign

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 10: Update Design Document with Implementation Status

**Files:**
- Modify: `docs/plans/2026-02-23-hero-carousel-redesign-design.md`

**Step 1: Mark implementation as complete**

Find the Approval section at the end and update:

```markdown
## Approval

- [x] Design approved by user
- [x] Implementation completed
- [ ] Accessibility verified
- [ ] Cross-browser testing completed
- [ ] RTL testing completed
```

**Step 2: Commit**

```bash
git add docs/plans/2026-02-23-hero-carousel-redesign-design.md
git commit -m "docs: mark hero carousel redesign implementation as complete

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Summary

| Task | Description | Complexity |
|------|-------------|------------|
| 1 | Update i18n copy (EN) | Low |
| 2 | Add Arabic translations | Low |
| 3 | Slide 1 illustration | Medium |
| 4 | Slide 2 illustration | Medium |
| 5 | Slide 3 illustration | Medium |
| 6 | Subtitle + capability strip | Medium |
| 7 | Backdrop glow + inner card | Low |
| 8 | Reduced motion CSS | Low |
| 9 | Verification & testing | Medium |
| 10 | Update design doc | Low |

**Total Estimated Scope:** Medium-High

---

*Plan generated by Claude AI on 2026-02-23*
