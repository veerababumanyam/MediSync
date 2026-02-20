import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import I18nextHttpBackend from 'i18next-http-backend'

/**
 * i18n Configuration for MediSync
 *
 * Supports:
 * - English (en) - Left-to-Right (LTR)
 * - Arabic (ar) - Right-to-Left (RTL)
 *
 * Namespaces:
 * - translation: Common app-wide translations
 * - chat: Chat interface translations
 * - dashboard: Dashboard widget translations
 * - alerts: Alert and notification translations
 * - reports: Scheduled report translations
 *
 * Locale Detection Priority:
 * 1. localStorage user preference
 * 2. URL parameter ?lang=ar or ?lang=en
 * 3. Browser's Accept-Language header
 * 4. Default: 'en'
 */

// English translations (inline for core app)
const enResources = {
  translation: {
    app: {
      name: 'MediSync',
      tagline: 'AI-Powered Business Intelligence',
      toggleLanguage: 'Toggle language',
    },
    welcome: {
      title: 'Welcome to MediSync',
      description:
        'Ask questions in plain language and get instant charts, tables, and insights from your healthcare and accounting data.',
    },
    features: {
      conversationalBI: {
        title: 'Conversational BI',
        description: 'Chat with your data using natural language. Get instant visualizations.',
      },
      aiAccountant: {
        title: 'AI Accountant',
        description:
          'Upload documents and let AI extract, map, and sync to your accounting system.',
      },
      easyReports: {
        title: 'Easy Reports',
        description: 'Pre-built reports and custom dashboards with automated delivery.',
      },
    },
    status: {
      title: 'System Status',
      react: 'React',
      vite: 'Vite',
      copilotkit: 'CopilotKit',
      i18n: 'i18n',
    },
    footer: {
      copyright:
        '© 2026 MediSync. AI-Powered Conversational BI & Intelligent Accounting for Healthcare.',
    },
    common: {
      loading: 'Loading...',
      error: 'An error occurred',
      retry: 'Retry',
      cancel: 'Cancel',
      save: 'Save',
      delete: 'Delete',
      edit: 'Edit',
      close: 'Close',
      confirm: 'Confirm',
      yes: 'Yes',
      no: 'No',
    },
    navigation: {
      dashboard: 'Dashboard',
      chat: 'Chat',
      alerts: 'Alerts',
      reports: 'Reports',
      settings: 'Settings',
    },
  },
}

// Arabic translations (inline for core app)
const arResources = {
  translation: {
    app: {
      name: 'ميدي سنك',
      tagline: 'ذكاء الأعمال المدعوم بالذكاء الاصطناعي',
      toggleLanguage: 'تبديل اللغة',
    },
    welcome: {
      title: 'مرحبًا بك في ميدي سنك',
      description:
        'اطرح أسئلة بلغة بسيطة واحصل على رسوم بيانية وجداول ورؤى فورية من بيانات الرعاية الصحية والمحاسبة الخاصة بك.',
    },
    features: {
      conversationalBI: {
        title: 'الذكاء التجاري المحادثي',
        description: 'تحدث مع بياناتك باستخدام اللغة الطبيعية. احصل على تصورات فورية.',
      },
      aiAccountant: {
        title: 'المحاسب الذكي',
        description:
          'قم بتحميل المستندات ودع الذكاء الاصطناعي يستخرج ويربط ويزامن مع نظام المحاسبة الخاص بك.',
      },
      easyReports: {
        title: 'التقارير السهلة',
        description: 'تقارير جاهزة ولوحات معلومات مخصصة مع التوصيل التلقائي.',
      },
    },
    status: {
      title: 'حالة النظام',
      react: 'رياكت',
      vite: 'فايت',
      copilotkit: 'كوبيلوكت كيت',
      i18n: 'دعم اللغات',
    },
    footer: {
      copyright:
        '© 2026 ميدي سنك. ذكاء الأعمال المحادثي والمحاسبة الذكية للرعاية الصحية المدعومة بالذكاء الاصطناعي.',
    },
    common: {
      loading: 'جاري التحميل...',
      error: 'حدث خطأ',
      retry: 'إعادة المحاولة',
      cancel: 'إلغاء',
      save: 'حفظ',
      delete: 'حذف',
      edit: 'تعديل',
      close: 'إغلاق',
      confirm: 'تأكيد',
      yes: 'نعم',
      no: 'لا',
    },
    navigation: {
      dashboard: 'لوحة التحكم',
      chat: 'المحادثة',
      alerts: 'التنبيهات',
      reports: 'التقارير',
      settings: 'الإعدادات',
    },
  },
}

// Detect initial language
const detectInitialLanguage = (): string => {
  // Check URL parameter first
  const urlParams = new URLSearchParams(window.location.search)
  const urlLang = urlParams.get('lang')
  if (urlLang === 'ar' || urlLang === 'en') {
    return urlLang
  }

  // Check localStorage
  const storedLang = localStorage.getItem('medisync-locale')
  if (storedLang === 'ar' || storedLang === 'en') {
    return storedLang
  }

  // Check browser language
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('ar')) {
    return 'ar'
  }

  // Default to English
  return 'en'
}

// Initialize i18next with HTTP backend for lazy-loaded namespaces
void i18n
  .use(I18nextHttpBackend)
  .use(initReactI18next)
  .init({
    resources: {
      en: enResources,
      ar: arResources,
    },
    lng: detectInitialLanguage(),
    fallbackLng: 'en',
    debug: import.meta.env.DEV,

    // Namespace configuration
    ns: ['translation', 'common', 'chat', 'dashboard', 'alerts', 'reports'],
    defaultNS: 'translation',
    fallbackNS: 'translation',

    // Backend configuration for lazy loading
    backend: {
      loadPath: '/src/i18n/locales/{{lng}}/{{ns}}.json',
    },

    interpolation: {
      escapeValue: false, // React already escapes values
    },

    react: {
      useSuspense: true,
    },

    // Save language preference to localStorage on change
    saveMissing: false,
  })

// Save language preference on change
i18n.on('languageChanged', (lng) => {
  localStorage.setItem('medisync-locale', lng)
})

export default i18n
