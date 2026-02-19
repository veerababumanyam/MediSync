# MediSync — Internationalisation & Localisation Architecture

**Version:** 1.0 | **Created:** February 19, 2026  
**Status:** Approved — Design & Development Baseline  
**Cross-ref:** [PRD.md §6.10](./PRD.md) | [DESIGN.md §16](./DESIGN.md) | [agents/00-agent-backlog.md §2.5](./agents/00-agent-backlog.md) | [agents/specs/e-01-language-detection-routing.md](./agents/specs/e-01-language-detection-routing.md)

---

## 0. Motivation & Scope

MediSync operates in a healthcare and accounting context where **language is a clinical concern**. Misread Arabic-script drug names, RTL-broken report layouts, or English-only AI responses force bilingual staff to context-switch, introducing errors. This document architects full product i18n (internationalisation) and l10n (localisation) across every surface.

**Phase 0 (this document) — supported locales:**

| Code | Language | Script | Direction | Status |
|------|----------|--------|-----------|--------|
| `en` | English | Latin | LTR | ✅ Default |
| `ar` | Arabic | Arabic | **RTL** | ✅ Phase 1 |

Architecture is designed for easy addition of further locales (e.g., `ur`, `hi`, `fr`) without structural changes.

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        USER LAYER                                       │
│  React Web (i18next + Tailwind RTL)   Flutter Mobile (flutter_intl ARB) │
└────────────────────────┬────────────────────────┬───────────────────────┘
                         │ locale pref (HTTP header│ / JWT claim)
                         ▼                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     API GATEWAY (Go / go-chi)                           │
│  Accept-Language negotiation → inject locale into request context       │
└────────────────────────┬────────────────────────────────────────────────┘
                         │
          ┌──────────────┴──────────────────┐
          ▼                                 ▼
┌──────────────────┐            ┌───────────────────────┐
│  Business Logic  │            │   AI Orchestration     │
│  Layer (Go)      │            │   (Genkit Flows)       │
│                  │            │                        │
│ • golang.org/x/  │            │ • Multilingual system  │
│   text/language  │            │   prompt injection     │
│ • Locale-aware   │            │ • E-01 Language        │
│   date/number/   │            │   Detection & Routing  │
│   currency fmt   │            │   Agent                │
└──────────────────┘            └───────────────────────┘
          │                                 │
          ▼                                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    DATA & RENDERING LAYER                               │
│  • PostgreSQL: locale column in user_preferences table                  │
│  • Report Engine: WeasyPrint/puppeteer with Arabic font + RTL PDF       │
│  • Translation Store: /locales/{lang}/translation.json (flat namespace) │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Locale Detection & Preference Chain

Priority is evaluated top-down. First match wins.

| Priority | Source | Mechanism |
|----------|--------|-----------|
| 1 | **User Profile** | `user_preferences.locale` stored in Postgres; loaded into JWT claim `locale` at login. |
| 2 | **Browser / OS** | `Accept-Language` HTTP header sent by browser or Flutter OS locale. |
| 3 | **URL Parameter** | `?lang=ar` query param (for report share links and email report links). |
| 4 | **Default** | `en` (English, LTR). |

Locale preference is settable on the **Profile Settings** screen and persists to the database immediately. The JWT is refreshed on save so all tabs and devices update within one session tick.

---

## 3. Translation File Format & Structure

### 3.1 React Web — i18next JSON

Translation files live under `frontend/public/locales/{lang}/`:

```
frontend/public/locales/
  en/
    common.json         ← shared UI (buttons, labels, errors)
    dashboard.json      ← dashboard screen namespace
    chat.json           ← chat interface namespace
    reports.json        ← reports module namespace
    accountant.json     ← AI accountant module namespace
    ai-responses.json   ← AI reply templates and status messages
    notifications.json  ← alert and notification text
    validation.json     ← form validation messages
  ar/
    common.json
    dashboard.json
    chat.json
    reports.json
    accountant.json
    ai-responses.json
    notifications.json
    validation.json
```

**Naming convention:** Keys use `snake_case`, nested max 2 levels deep, English as canonical source-of-truth.

```json
// en/common.json
{
  "nav": {
    "dashboard": "Dashboard",
    "chat": "Chat",
    "reports": "Reports",
    "accountant": "AI Accountant",
    "settings": "Settings"
  },
  "actions": {
    "save": "Save",
    "cancel": "Cancel",
    "export": "Export",
    "pin_to_dashboard": "Pin to Dashboard",
    "sync_now": "Sync Now"
  }
}

// ar/common.json
{
  "nav": {
    "dashboard": "لوحة التحكم",
    "chat": "المحادثة",
    "reports": "التقارير",
    "accountant": "المحاسب الذكي",
    "settings": "الإعدادات"
  },
  "actions": {
    "save": "حفظ",
    "cancel": "إلغاء",
    "export": "تصدير",
    "pin_to_dashboard": "تثبيت في لوحة التحكم",
    "sync_now": "مزامنة الآن"
  }
}
```

### 3.2 Flutter Mobile — ARB (Application Resource Bundle)

```
mobile/lib/l10n/
  app_en.arb          ← English strings
  app_ar.arb          ← Arabic strings
```

ARB keys mirror the i18next namespaced keys for consistency during copy management:

```json
// app_ar.arb
{
  "@@locale": "ar",
  "navDashboard": "لوحة التحكم",
  "@navDashboard": { "description": "Bottom nav: Dashboard tab" },
  "chatInputPlaceholder": "اسأل عن بياناتك...",
  "@chatInputPlaceholder": { "description": "Chat input placeholder text" }
}
```

---

## 4. RTL Layout Architecture

### 4.1 Principle: Logical Properties First

All layout code must use **CSS logical properties** (not physical left/right) to enable automatic mirroring:

| ❌ Physical (avoid) | ✅ Logical (use) |
|---------------------|-----------------|
| `margin-left: 16px` | `margin-inline-start: 16px` |
| `padding-right: 8px` | `padding-inline-end: 8px` |
| `border-left: ...` | `border-inline-start: ...` |
| `text-align: left` | `text-align: start` |
| `float: right` | Float replaced with Flexbox `flex-end` |

### 4.2 HTML `dir` Attribute

The root `<html>` element receives `dir` based on locale:

```tsx
// frontend/src/App.tsx
import { useTranslation } from 'react-i18next';

export function App() {
  const { i18n } = useTranslation();
  const isRTL = i18n.dir() === 'rtl';

  useEffect(() => {
    document.documentElement.dir = i18n.dir();
    document.documentElement.lang = i18n.language;
  }, [i18n.language]);

  return <RouterProvider router={router} />;
}
```

### 4.3 Tailwind CSS RTL Plugin

Add `tailwindcss-rtl` (MIT) or use Tailwind v3.3+ built-in `rtl:` and `ltr:` variants:

```js
// tailwind.config.js
module.exports = {
  content: ['./src/**/*.{tsx,ts}'],
  theme: { /* design tokens */ },
  plugins: [],
  // Tailwind v3.3+ supports dir-aware utilities natively:
  // rtl:mr-4 → margin-right: 1rem only in RTL mode
}
```

Usage in components:

```tsx
// A sidebar that slides from the correct edge
<aside className="
  fixed top-0 start-0        // logical: left in LTR, right in RTL
  w-64 h-full
  ltr:border-r rtl:border-l  // fallback for browsers without logical border
  border-glass
">
```

### 4.4 Flutter RTL

Flutter's `Directionality` widget and `TextDirection` are driven by the active `Locale`:

```dart
// mobile/lib/main.dart
MaterialApp(
  locale: _currentLocale,
  supportedLocales: const [Locale('en'), Locale('ar')],
  localizationsDelegates: const [
    AppLocalizations.delegate,
    GlobalMaterialLocalizations.delegate,    // RTL-aware Material widgets
    GlobalWidgetsLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
  ],
  // Flutter automatically sets TextDirection.rtl for 'ar' locale
);
```

Use `EdgeInsetsDirectional` and `AlignmentDirectional` everywhere:

```dart
// ✅ Correct — mirrors automatically
padding: EdgeInsetsDirectional.only(start: 16, end: 8);
alignment: AlignmentDirectional.centerStart;

// ❌ Incorrect — breaks RTL
padding: EdgeInsets.only(left: 16, right: 8);
```

### 4.5 RTL-Specific UI Adjustments

| Element | LTR Behaviour | RTL (Arabic) Behaviour |
|---------|--------------|------------------------|
| Navigation sidebar | Left edge | Right edge |
| Chat bubbles — user | Right-aligned, tail right | Left-aligned, tail left |
| Chat bubbles — AI | Left-aligned, tail left | Right-aligned, tail right |
| Data tables | Text left, numbers right | Text right, numbers left |
| Breadcrumb chevron | `›` pointing right | `‹` pointing left |
| Progress bars | Fill left→right | Fill right→left |
| Drill-down arrows | `→` | `←` |
| Carousel / slide | Swipe left for next | Swipe right for next |
| Icons with direction | Forward-facing arrows flip | Apply `rtl:scale-x-[-1]` CSS class |

### 4.6 Bidirectional Text (BiDi) in Chat

Chat input and AI responses can contain mixed-language content (English numbers, Arabic text). Use the Unicode Bidirectional Algorithm correctly:

```tsx
// ChatMessage.tsx
<p
  dir="auto"          // browser auto-detects per paragraph
  className="text-body leading-relaxed"
>
  {message.content}
</p>
```

For structured AI responses embedding numbers and brand names in Arabic text, wrap inline LTR segments:

```tsx
<span dir="ltr" className="inline">{currencyValue}</span>
```

---

## 5. AI Layer — Multilingual Responses

### 5.1 System Prompt Language Injection

Every Genkit flow that produces user-visible text receives a `response_language` instruction appended to the system prompt:

```go
// internal/ai/prompts/base_prompt.go

const ResponseLanguageInstruction = `
LANGUAGE RULE (MANDATORY):
- The user's preferred language is: {{.locale}}
- All explanations, labels, chart titles, table headers, insight narratives,
  error messages, and recommendations MUST be written in {{.locale}}.
- Numbers, dates, and currency must follow {{.locale}} formatting conventions.
- For Arabic (ar): use right-to-left natural sentence structure; use Eastern
  Arabic-Indic numerals (٠١٢٣٤٥٦٧٨٩) only if the user has explicitly enabled
  them; default to Western numerals (0-9) for interoperability with data systems.
- SQL identifiers, column names, and database values remain in their source
  language (typically English) — do NOT translate table/column names.
- "Low confidence" and audit notices remain bilingual (English + {{.locale}})
  for compliance traceability.
`
```

### 5.2 LLM Selection for Arabic

| Capability | Recommended Model | Notes |
|------------|------------------|-------|
| Text-to-SQL + Arabic explanation | **GPT-4o** or **Claude 3.5 Sonnet** | Strong Arabic reasoning |
| Arabic OCR post-processing | **Gemini 1.5 Pro** | Best Arabic document understanding |
| Short UI copy generation | **Gemini Flash 1.5** | Cost-efficient for bulk strings |
| Local/offline (sensitive data) | **Llama 3.1 70B** (via Ollama) | Arabic capability reasonable at 70B+ |

### 5.3 E-01 Language Detection & Routing Agent

All user queries pass through the **E-01 Language Detection & Routing Agent** before reaching the domain agents (A-01 Text-to-SQL, etc.). See [agents/specs/e-01-language-detection-routing.md](./agents/specs/e-01-language-detection-routing.md).

```
User Query (any language)
  → [E-01] Language Detection
  → [E-01] Query Normalisation (if needed)
  → [A-01 / B-xx / C-xx / D-xx] Domain Agent
      (receives locale-tagged request)
  → Response Generator
      (applies locale formatting + ResponseLanguageInstruction)
  → Chat Response (in user's language)
```

### 5.4 Arabic Text-to-SQL Considerations

Arabic user queries are understood but SQL is always generated in English to interact with the English-schema data warehouse. The agent pipeline handles this transparently:

```
Arabic query: "أعطني إيرادات هذا الشهر مقارنةً بالشهر الماضي"
    ↓ E-01 normalises intent
Intent (internal): compare_revenue(period=current_month, vs=previous_month)
    ↓ A-01 generates SQL (in English)
SQL: SELECT ... FROM revenue WHERE ...
    ↓ Execute → results
    ↓ Response formatted in Arabic with Arabic number formatting
Arabic response: "إيرادات هذا الشهر: ١٢٣,٤٥٦ ريال — بزيادة ٨٪ عن الشهر الماضي"
```

---

## 6. Number, Date, Currency & Calendar Formatting

### 6.1 Number Formatting

```go
// Go backend — golang.org/x/text/message
import "golang.org/x/text/message"

func FormatNumber(val float64, locale string) string {
    p := message.NewPrinter(language.Make(locale))
    return p.Sprintf("%v", val)
}
// en: 1,234,567.89
// ar: ١٬٢٣٤٬٥٦٧٫٨٩  (or 1,234,567.89 with Western digits per user pref)
```

```tsx
// React frontend
const fmt = new Intl.NumberFormat(locale, {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});
fmt.format(1234567.89);
// en-US → "1,234,567.89"
// ar-SA → "١٬٢٣٤٬٥٦٧٫٨٩"
```

### 6.2 Currency Formatting

```tsx
const fmtCurrency = new Intl.NumberFormat(locale, {
  style: 'currency',
  currency: 'AED',          // configurable per tenant
  currencyDisplay: 'symbol',
});
// en → "AED 1,234.50"
// ar → "١٬٢٣٤٫٥٠ د.إ"   (symbol on the right in Arabic)
```

### 6.3 Date Formatting

```tsx
// Accept locale + calendar pref from user_preferences
const fmtDate = new Intl.DateTimeFormat(locale, {
  year: 'numeric', month: 'long', day: 'numeric',
  calendar: userPreferences.calendar ?? 'gregory', // 'gregory' | 'islamic-umalqura'
});
// en → "February 19, 2026"
// ar + gregory → "١٩ فبراير ٢٠٢٦"
// ar + islamic-umalqura → "٢١ شعبان ١٤٤٧"
```

### 6.4 Report Date Ranges

Scheduling agents (A-09, C-03) store date ranges in UTC ISO 8601 internally. Presentation layer converts to locale-aware display. Hijri calendar support is **opt-in per user** — not default — to avoid confusion in financial reporting.

---

## 7. Report Generation — RTL PDF & Excel

### 7.1 PDF (WeasyPrint / Puppeteer)

Both PDF renderers support HTML/CSS as input, enabling full RTL support via the same CSS logical properties used in the React app.

**Font requirements for Arabic PDF:**

```css
/* report-base.css — loaded by PDF renderer */
@font-face {
  font-family: 'NotoSansArabic';
  src: url('/fonts/NotoSansArabic-VF.ttf') format('truetype');
}

@font-face {
  font-family: 'Cairo';
  src: url('/fonts/Cairo-VF.ttf') format('truetype');
}

body[lang="ar"] {
  font-family: 'Cairo', 'NotoSansArabic', sans-serif;
  direction: rtl;
  unicode-bidi: bidi-override;
}
```

**Recommended fonts (OFL licensed — free for commercial use):**

| Font | Style | Use |
|------|-------|-----|
| **Cairo** | Sans-serif, professional | Report body, table cells |
| **Noto Sans Arabic** | Neutral, multi-weight | Fallback, technical labels |
| **Scheherazade New** | Traditional | Legal/formal documents |

### 7.2 Excel (xlsx)

```go
// Go — excelize library
f := excelize.NewFile()
// Set sheet RTL
f.SetSheetView("Sheet1", 0, &excelize.ViewOptions{
    RightToLeft: boolPtr(locale == "ar"),
})
// Write Arabic strings with correct encoding (excelize handles UTF-8)
f.SetCellValue("Sheet1", "A1", "اسم المورد")
```

### 7.3 Chart Axis Labels

Apache ECharts supports RTL axis:

```js
// echarts RTL configuration
const option = {
  textStyle: { fontFamily: locale === 'ar' ? 'Cairo, NotoSansArabic' : 'Inter' },
  xAxis: {
    axisLabel: {
      formatter: (val) => fmtDate.format(new Date(val)),
      align: locale === 'ar' ? 'right' : 'left',
    }
  },
  // Flip legend and tooltip positions for RTL
  legend: { right: locale === 'ar' ? 'auto' : 10, left: locale === 'ar' ? 10 : 'auto' },
};
```

---

## 8. Backend Go — Locale-Aware Middleware

### 8.1 Locale Context Middleware

```go
// middleware/locale.go
func LocaleMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        locale := resolveLocale(r)
        ctx := context.WithValue(r.Context(), ctxKeyLocale, locale)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func resolveLocale(r *http.Request) string {
    // 1. JWT claim
    if claims, ok := jwtFromContext(r.Context()); ok {
        if l := claims.Locale; l != "" { return normalise(l) }
    }
    // 2. Accept-Language header (golang.org/x/text/language matching)
    supported := language.NewMatcher([]language.Tag{
        language.English, language.Arabic,
    })
    t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
    tag, _, _ := supported.Match(t...)
    return tag.String()
}
```

### 8.2 Localised Error Messages

```go
// errors/messages.go — key-based, loaded from YAML per locale
var errMessages = map[string]map[string]string{
    "en": {
        "query.ambiguous":   "Your question is ambiguous. Please clarify which time period you mean.",
        "query.off_topic":   "I can only answer questions about your business data.",
        "sync.tally_error":  "Tally sync failed. Please check the connection and retry.",
    },
    "ar": {
        "query.ambiguous":   "سؤالك غير واضح. يرجى توضيح الفترة الزمنية المقصودة.",
        "query.off_topic":   "يمكنني فقط الإجابة عن أسئلة تتعلق ببيانات عملك.",
        "sync.tally_error":  "فشلت مزامنة Tally. يرجى التحقق من الاتصال والمحاولة مرة أخرى.",
    },
}

func LocalisedError(locale, key string) string {
    if msgs, ok := errMessages[locale]; ok {
        if msg, ok := msgs[key]; ok { return msg }
    }
    return errMessages["en"][key] // English fallback
}
```

---

## 9. User Settings Screen — Language Switcher

### 9.1 Settings Fields

| Field | Type | Default | Notes |
|-------|------|---------|-------|
| `display_language` | enum: `en`, `ar` | `en` | Drives all UI text |
| `number_format` | enum: `western`, `eastern_arabic` | `western` | `western` = 0–9, `eastern_arabic` = ٠–٩ |
| `calendar_system` | enum: `gregorian`, `hijri` | `gregorian` | Affects date display only |
| `report_language` | enum: `en`, `ar`, `both` | inherits from `display_language` | Reports can be bilingual |
| `ai_response_language` | enum: `en`, `ar`, `auto` | `auto` | `auto` = match query language |

### 9.2 Instant Locale Switch (no page reload)

```tsx
// LanguageSwitch.tsx
const { i18n } = useTranslation();

async function switchLocale(code: 'en' | 'ar') {
  await i18n.changeLanguage(code);                        // loads namespace bundle
  document.documentElement.dir = i18n.dir(code);          // flip layout
  document.documentElement.lang = code;
  await updateUserPreference({ display_language: code }); // persist to DB
}
```

---

## 10. Database Schema Changes

```sql
-- Add locale preference to user table
ALTER TABLE users
  ADD COLUMN locale              VARCHAR(10)  NOT NULL DEFAULT 'en',
  ADD COLUMN number_format       VARCHAR(20)  NOT NULL DEFAULT 'western',
  ADD COLUMN calendar_system     VARCHAR(20)  NOT NULL DEFAULT 'gregorian',
  ADD COLUMN report_language     VARCHAR(10)  NOT NULL DEFAULT 'en',
  ADD COLUMN ai_response_lang    VARCHAR(10)  NOT NULL DEFAULT 'auto';

-- Locale constraint
ALTER TABLE users
  ADD CONSTRAINT chk_locale CHECK (locale IN ('en', 'ar'));

-- Audit log — record locale at time of action
ALTER TABLE audit_log
  ADD COLUMN user_locale VARCHAR(10) NOT NULL DEFAULT 'en';

-- Scheduled reports — language for distributed report
ALTER TABLE scheduled_reports
  ADD COLUMN report_locale VARCHAR(10) NOT NULL DEFAULT 'en';
```

---

## 11. Translation Workflow & Governance

### 11.1 Source of Truth

English (`en`) is the canonical source. All keys are defined in English first.

### 11.2 Translation Management

| Stage | Tool | Owner |
|-------|------|-------|
| Key extraction | `i18next-parser` (auto-scans source) | Dev |
| Missing key detection | CI check: `node scripts/check-missing-keys.js` | CI/CD |
| Translation | Human translators + AI pre-fill (GPT-4o) for draft | Localisation team |
| Review | Native Arabic speaker reviews all AI-generated drafts | QA |
| Medical term validation | Clinical terminology check for healthcare strings (drug names, diagnoses) | Medical advisor |
| Deployment | Translation files bundled at build time; lazy-loaded per namespace | DevOps |

### 11.3 Missing Translation Fallback

```tsx
// i18next config
i18n.init({
  fallbackLng: 'en',        // show English if Arabic key is missing
  ns: ['common', 'dashboard', 'chat', 'reports', 'accountant', 'ai-responses'],
  defaultNS: 'common',
  interpolation: { escapeValue: false },
  saveMissing: process.env.NODE_ENV === 'development',  // log missing keys in dev
});
```

### 11.4 Medical & Accounting Terminology Glossary

A governed bilingual glossary is maintained at `docs/i18n-glossary.md` and fed into:
- AI system prompts (domain terminology normalisation — agent A-04)
- Human translator style guides
- OCR post-processing validation (agent B-02)

Sample entries:

| English | Arabic | Context |
|---------|--------|---------|
| Patient Footfall | زيارات المرضى | BI dashboard, HIMS |
| Accounts Receivable | المدينون / ذمم مدينة | Accounting, Tally |
| Outstanding Invoices | الفواتير المستحقة | Accountant module |
| Ledger Mapping | تعيين دفتر الأستاذ | AI Accountant |
| Inventory Aging | تقادم المخزون | Reports, Pharmacy |
| Bank Reconciliation | مطابقة الحسابات البنكية | AI Accountant |
| Days Sales Outstanding | أيام المبيعات المعلقة | Reporting KPIs |
| Cost of Goods Sold | تكلفة البضاعة المباعة | P&L reports |
| Gross Profit Margin | هامش الربح الإجمالي | Financial analytics |

---

## 12. Testing Strategy

### 12.1 Unit Tests

- Translation key coverage test: ensures all `en` keys exist in `ar`.
- String interpolation test: variables (`{{count}}`, `{{name}}`) render correctly in both locales.
- Number/date formatter tests with locale snapshots.

### 12.2 Visual Regression Tests (RTL)

Use **Playwright** with `locale: 'ar'` configuration to screenshot every screen and compare against RTL baselines:

```ts
// tests/rtl-visual.spec.ts
test.use({ locale: 'ar', timezoneId: 'Asia/Dubai' });

test('Dashboard renders correctly in RTL', async ({ page }) => {
  await page.goto('/dashboard');
  await expect(page).toHaveScreenshot('dashboard-rtl.png');
});
```

### 12.3 BiDi Text Tests

- Validate mixed LTR/RTL paragraphs (numbers in Arabic text) render without overflow.
- Test chat input accepts and displays Arabic text correctly.
- Test report PDF export produces valid RTL layout with correct Arabic font rendering.

### 12.4 AI Response Quality Tests

Evaluation dataset in `tests/i18n/ar-response-eval.jsonl`:

```jsonl
{"query_ar": "ما هي أعلى ١٠ أدوية مبيعاً هذا الأسبوع؟", "expected_lang": "ar", "expected_contains_sql": true}
{"query_en": "Show me last month revenue", "user_locale": "ar", "expected_response_lang": "ar"}
```

---

## 13. Phase Roadmap

| Phase | Scope | Target |
|-------|-------|--------|
| **Phase 1** | Core UI (web + mobile), AI responses, chat in `en` + `ar` | Sprint 3 |
| **Phase 2** | Reports (PDF/Excel) in `en` + `ar`, bilingual report option | Sprint 5 |
| **Phase 3** | Email notifications, scheduled reports, alert messages in locale | Sprint 6 |
| **Phase 4** | Hijri calendar opt-in, Eastern Arabic numeral opt-in | Sprint 8 |
| **Phase 5** | Additional locale scaffolding (`ur`, `hi`) — translation content TBD | Post-v1 |

---

## 14. Open Questions & Decisions

| # | Question | Decision | Owner |
|---|----------|----------|-------|
| Q1 | Eastern Arabic-Indic numerals (٠١٢٣) vs Western (012) as default in Arabic reports? | **Western digits default** — avoids confusion in financial documents; Eastern opt-in | Product |
| Q2 | Should the Hijri calendar be offered for scheduling financial report periods? | **No** — financial periods follow Gregorian to align with Tally/HIMS. Hijri display opt-in only | Finance Advisor |
| Q3 | Bilingual PDF reports (Arabic + English on same page)? | **Supported** — via `report_language: "both"` user setting; columns duplicate with both scripts | Design |
| Q4 | AI query language vs response language — should `ar` query always yield `ar` response? | **`auto` default**: respond in the language the user queried in; override via `ai_response_language` setting | AI Lead |
| Q5 | Should alert SMS/email notifications respect locale? | **Yes** — locale stored on `scheduled_reports` and `alert_rules` tables | Backend |

---

*Document Owner: Platform Architecture Team | Next Review: Sprint 5 Review*
