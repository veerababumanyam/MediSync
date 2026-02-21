---
name: i18next-i18n
description: This skill should be used when the user asks to "add translations", "implement i18n", "internationalize", "add language support", "RTL layout", "Arabic translations", "i18next setup", or mentions multi-language support for MediSync.
---

# i18next Internationalization for MediSync

i18next 24.2 provides comprehensive internationalization for MediSync's React web app, supporting English (LTR) and Arabic (RTL) as first-class languages from Phase 1.

★ Insight ─────────────────────────────────────
MediSync i18n architecture:
1. **Dual language** - English (en) and Arabic (ar) from day one
2. **RTL support** - Logical properties + react-i18next
3. **Namespace per feature** - Lazy-loaded translations
4. **Backend AI responses** - Locale instruction in prompts
5. **Validation** - CI checks for missing keys
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Package** | i18next (24.2.x) + react-i18next (14.x) |
| **Languages** | en (English), ar (Arabic/RTL) |
| **Format** | JSON translation files |
| **Loading** | Lazy by namespace |
| **Detection** | JWT preference → Accept-Language → URL param → Default |

## Project Structure

```
frontend/
├── public/
│   └── locales/
│       ├── en/
│       │   ├── common.json
│       │   ├── dashboard.json
│       │   ├── chat.json
│       │   ├── reports.json
│       │   └── errors.json
│       └── ar/
│           ├── common.json
│           ├── dashboard.json
│           ├── chat.json
│           ├── reports.json
│           └── errors.json
├── src/
│   ├── i18n/
│   │   ├── config.ts
│   │   └── RTLProvider.tsx
│   └── ...
```

## Configuration

### i18n Setup

```typescript
// src/i18n/config.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'en',
    supportedLngs: ['en', 'ar'],

    ns: ['common', 'dashboard', 'chat', 'reports', 'errors'],
    defaultNS: 'common',

    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },

    detection: {
      order: ['querystring', 'localStorage', 'navigator'],
      lookupQuerystring: 'lang',
      lookupLocalStorage: 'i18nextLng',
      caches: ['localStorage'],
    },

    interpolation: {
      escapeValue: false, // React escapes by default
    },

    react: {
      useSuspense: true,
    },
  });

export default i18n;

// Export RTL helper
export const isRTL = (lng: string): boolean => lng === 'ar';

// Export direction helper
export const getDirection = (lng: string): 'ltr' | 'rtl' =>
  isRTL(lng) ? 'rtl' : 'ltr';
```

### RTL Provider

```typescript
// src/i18n/RTLProvider.tsx
import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { isRTL, getDirection } from './config';

export function RTLProvider({ children }: { children: React.ReactNode }) {
  const { i18n } = useTranslation();

  useEffect(() => {
    const dir = getDirection(i18n.language);
    const html = document.documentElement;

    html.setAttribute('dir', dir);
    html.setAttribute('lang', i18n.language);

    // Add RTL class for Tailwind
    if (isRTL(i18n.language)) {
      html.classList.add('rtl');
    } else {
      html.classList.remove('rtl');
    }
  }, [i18n.language]);

  return <>{children}</>;
}
```

## Translation Files

### English (en/common.json)

```json
{
  "app": {
    "name": "MediSync",
    "tagline": "AI-Powered Healthcare Intelligence"
  },
  "navigation": {
    "dashboard": "Dashboard",
    "chat": "AI Assistant",
    "reports": "Reports",
    "settings": "Settings",
    "profile": "Profile"
  },
  "actions": {
    "save": "Save",
    "cancel": "Cancel",
    "delete": "Delete",
    "edit": "Edit",
    "view": "View",
    "download": "Download",
    "export": "Export",
    "search": "Search",
    "filter": "Filter",
    "apply": "Apply",
    "reset": "Reset",
    "retry": "Retry"
  },
  "status": {
    "loading": "Loading...",
    "saving": "Saving...",
    "error": "An error occurred",
    "success": "Operation successful",
    "noData": "No data available"
  },
  "dates": {
    "today": "Today",
    "yesterday": "Yesterday",
    "thisWeek": "This Week",
    "thisMonth": "This Month",
    "thisYear": "This Year"
  },
  "numbers": {
    "thousand": "K",
    "million": "M",
    "billion": "B"
  }
}
```

### Arabic (ar/common.json)

```json
{
  "app": {
    "name": "ميديسينك",
    "tagline": "الذكاء الاصطناعي للرعاية الصحية"
  },
  "navigation": {
    "dashboard": "لوحة التحكم",
    "chat": "المساعد الذكي",
    "reports": "التقارير",
    "settings": "الإعدادات",
    "profile": "الملف الشخصي"
  },
  "actions": {
    "save": "حفظ",
    "cancel": "إلغاء",
    "delete": "حذف",
    "edit": "تعديل",
    "view": "عرض",
    "download": "تحميل",
    "export": "تصدير",
    "search": "بحث",
    "filter": "تصفية",
    "apply": "تطبيق",
    "reset": "إعادة تعيين",
    "retry": "إعادة المحاولة"
  },
  "status": {
    "loading": "جاري التحميل...",
    "saving": "جاري الحفظ...",
    "error": "حدث خطأ",
    "success": "تمت العملية بنجاح",
    "noData": "لا توجد بيانات"
  },
  "dates": {
    "today": "اليوم",
    "yesterday": "أمس",
    "thisWeek": "هذا الأسبوع",
    "thisMonth": "هذا الشهر",
    "thisYear": "هذه السنة"
  },
  "numbers": {
    "thousand": "ألف",
    "million": "مليون",
    "billion": "مليار"
  }
}
```

## Usage Patterns

### Basic Translation

```typescript
import { useTranslation } from 'react-i18next';

export function Navigation() {
  const { t } = useTranslation();

  return (
    <nav>
      <a href="/dashboard">{t('navigation.dashboard')}</a>
      <a href="/chat">{t('navigation.chat')}</a>
      <a href="/reports">{t('navigation.reports')}</a>
    </nav>
  );
}
```

### With Interpolation

```typescript
// Translation: "Welcome back, {{name}}!"
const { t } = useTranslation();

<h1>{t('welcome.message', { name: user.name })}</h1>
```

### Pluralization

```json
{
  "items": "{{count}} item",
  "items_plural": "{{count}} items"
}
```

```typescript
<p>{t('items', { count: itemCount })}</p>
```

### Namespaces

```typescript
// Load specific namespace
const { t } = useTranslation('dashboard');

// Multiple namespaces
const { t } = useTranslation(['dashboard', 'common']);

// Access with namespace prefix
<h1>{t('dashboard:title')}</h1>
```

### Lazy Loading Namespaces

```typescript
// Component-level namespace loading
const ChatView = lazy(() => import('./ChatView'));

function ChatPage() {
  const { t } = useTranslation('chat');

  return (
    <Suspense fallback={<Loading />}>
      <ChatView />
    </Suspense>
  );
}
```

## RTL Layout Patterns

### Logical CSS Properties (Tailwind)

```tsx
// Use logical properties for RTL support
<div className="ms-4 me-8 ps-6 pe-2">
  {/* ms = margin-start (left in LTR, right in RTL) */}
  {/* me = margin-end (right in LTR, left in RTL) */}
  {/* ps = padding-start */}
  {/* pe = padding-end */}
</div>

// Text alignment
<div className="text-start">...</div>  // Use instead of text-left

// Border radius
<div className="rounded-s-lg">...</div>  // Use instead of rounded-l-lg
<div className="rounded-e-lg">...</div>  // Use instead of rounded-r-lg
```

### Directional Components

```typescript
import { useTranslation } from 'react-i18next';
import { isRTL } from '../i18n/config';

export function DirectionalIcon({ icon: Icon, className }: Props) {
  const { i18n } = useTranslation();
  const rtl = isRTL(i18n.language);

  return (
    <Icon
      className={cn(
        className,
        rtl && 'scale-x-[-1]' // Mirror icon in RTL
      )}
    />
  );
}

// Icons that should NOT be mirrored
const NO_MIRROR_ICONS = ['language', 'settings', 'help'];

export function SmartIcon({ name }: Props) {
  const { i18n } = useTranslation();
  const shouldMirror = isRTL(i18n.language) && !NO_MIRROR_ICONS.includes(name);

  return (
    <Icon
      name={name}
      className={shouldMirror ? 'transform -scale-x-100' : ''}
    />
  );
}
```

## Number and Date Formatting

### Localized Formatters

```typescript
// src/utils/formatters.ts
import { getLocale } from './locale';

export function formatNumber(value: number, lng: string): string {
  return new Intl.NumberFormat(getLocale(lng)).format(value);
}

export function formatCurrency(
  value: number,
  lng: string,
  currency: string = 'USD'
): string {
  return new Intl.NumberFormat(getLocale(lng), {
    style: 'currency',
    currency,
  }).format(value);
}

export function formatDate(date: Date | string, lng: string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat(getLocale(lng), {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(d);
}

export function formatRelativeTime(date: Date, lng: string): string {
  const rtf = new Intl.RelativeTimeFormat(getLocale(lng), {
    numeric: 'auto',
  });
  const diff = (date.getTime() - Date.now()) / 1000;

  if (Math.abs(diff) < 60) return rtf.format(Math.round(diff), 'seconds');
  if (Math.abs(diff) < 3600) return rtf.format(Math.round(diff / 60), 'minutes');
  if (Math.abs(diff) < 86400) return rtf.format(Math.round(diff / 3600), 'hours');
  return rtf.format(Math.round(diff / 86400), 'days');
}

// Arabic locale for numbers uses Eastern Arabic numerals automatically
// To use Western numerals in Arabic, use { nu: 'latn' }
export function formatNumberWestern(value: number, lng: string): string {
  if (lng === 'ar') {
    return new Intl.NumberFormat('ar', { numberingSystem: 'latn' }).format(value);
  }
  return formatNumber(value, lng);
}
```

## Backend AI Locale Integration

### Passing Locale to AI Prompts

```typescript
// src/services/ai/chat.ts
export function buildPrompt(query: string, user: User): string {
  const locale = user.preferences.locale || 'en';
  const isRTL = locale === 'ar';

  return `
${query}

ResponseLanguageInstruction: Respond entirely in ${isRTL ? 'Arabic' : 'English'}.
Format numbers according to ${isRTL ? 'Arabic (Egypt)' : 'English (US)'} locale.
Use ${isRTL ? 'Arabic' : 'Western'} numerals.
Format dates as ${isRTL ? 'DD/MM/YYYY' : 'MM/DD/YYYY'}.
`;
}
```

## Testing and Validation

### CI Translation Check

```bash
#!/bin/bash
# scripts/check-translations.sh

EN_DIR="public/locales/en"
AR_DIR="public/locales/ar"

# Get all English keys
en_keys=$(cat $EN_DIR/*.json | jq -s 'add | keys')

# Get all Arabic keys
ar_keys=$(cat $AR_DIR/*.json | jq -s 'add | keys')

# Compare keys
missing=$(jq -n --argjson en "$en_keys" --argjson ar "$ar_keys" \
  '$en - $ar')

if [ "$missing" != "[]" ]; then
  echo "Missing Arabic translations for keys:"
  echo "$missing"
  exit 1
fi

echo "All translations present!"
```

### Unit Tests

```typescript
// src/i18n/config.test.ts
import i18n from './config';

describe('i18n configuration', () => {
  it('should have English translations', async () => {
    await i18n.changeLanguage('en');
    expect(i18n.t('app.name')).toBe('MediSync');
  });

  it('should have Arabic translations', async () => {
    await i18n.changeLanguage('ar');
    expect(i18n.t('app.name')).toBe('ميديسينك');
  });

  it('should detect RTL for Arabic', () => {
    expect(isRTL('ar')).toBe(true);
    expect(isRTL('en')).toBe(false);
  });
});
```

## Additional Resources

### Reference Files
- **`references/rtl-patterns.md`** - Comprehensive RTL layout patterns
- **`references/arabic-typography.md`** - Arabic typography guidelines

### Example Files
- **`examples/TranslationProvider.tsx`** - Complete provider setup
- **`examples/RTLComponents.tsx`** - RTL-aware component examples

### Scripts
- **`scripts/check-translations.sh`** - Validate translation completeness
- **`scripts/sort-translations.sh`** - Sort translation files alphabetically
