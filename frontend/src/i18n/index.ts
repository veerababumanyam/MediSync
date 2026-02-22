import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import enChat from './locales/en/chat.json'
import enDashboard from './locales/en/dashboard.json'
import enCommon from './locales/en/common.json'
import enAlerts from './locales/en/alerts.json'
import enReports from './locales/en/reports.json'
import enCouncil from './locales/en/council.json'
import enCopilot from './locales/en/copilot.json'

import arChat from './locales/ar/chat.json'
import arDashboard from './locales/ar/dashboard.json'
import arCommon from './locales/ar/common.json'
import arAlerts from './locales/ar/alerts.json'
import arReports from './locales/ar/reports.json'
import arCouncil from './locales/ar/council.json'
import arCopilot from './locales/ar/copilot.json'

/**
 * i18n Configuration for MediSync
 *
 * Standards: W3C i18n (lang/dir), BCP 47 for Intl formatting, i18next best practices.
 * See docs/i18n-architecture.md for locale chain, RTL, and number/date formatting.
 *
 * Supports:
 * - English (en) - Left-to-Right (LTR)
 * - Arabic (ar) - Right-to-Left (RTL)
 *
 * Namespaces (bundled at build time; no HTTP backend):
 * - translation: Common app-wide translations
 * - chat, dashboard, common, alerts, reports, council, copilot
 *
 * Locale Detection (initial load): localStorage → URL ?lang= → Accept-Language → default 'en'
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
      smartReports: {
        title: 'Smart Reports',
        description: 'Pre-built MIS reports for healthcare. Create custom dashboards without code.',
      },
      deepAnalytics: {
        title: 'Deep Analytics',
        description: 'Autonomous AI analyst conducts research. Get prescriptive, actionable insights.',
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
      home: 'Home',
      dashboard: 'Dashboard',
      chat: 'Chat',
      alerts: 'Alerts',
      reports: 'Reports',
      settings: 'Settings',
      features: 'Features',
      pricing: 'Pricing',
      about: 'About',
      toggleMenu: 'Toggle mobile menu',
    },
    home: {
      hero: {
        badge: 'AI-Powered Healthcare Intelligence',
        title: 'Your Data, ',
        titleHighlight: 'Understood',
        subtitle: 'Ask questions in plain language. Get instant insights from your HIMS and Tally data. No SQL required. No spreadsheets.',
        ctaChat: 'Start Chatting',
        ctaDashboard: 'View Dashboard',
        trustBy: 'Trusted by healthcare organizations worldwide',
        hipaa: 'HIPAA Compliant',
        soc2: 'SOC 2 Certified',
        clinics: '500+ Clinics',
        queries: '10M+ Queries',
      },
      preview: {
        askAnything: 'Ask Anything',
        askDesc: 'Natural language queries',
        revenueQuestion: "What's today's revenue?",
        instantResponse: 'Instant AI response',
        financialInsights: 'Financial Insights',
        tallySync: 'Tally ERP sync',
        outstanding: 'Outstanding',
        collected: 'Collected',
        autoSynced: 'Auto-synced with Tally',
        patientMetrics: 'Patient Metrics',
        himsIntegration: 'HIMS integration',
        today: 'Today',
        vsYesterday: 'vs Yesterday',
        thisMonth: 'This Month',
        depts: 'Depts',
        realTime: 'Real-time from HIMS',
      },
      section: {
        title: 'Everything You Need',
        subtitle: 'From conversational queries to automated accounting, MediSync connects your healthcare data in ways you never thought possible.',
      },
      footer: {
        product: 'Product',
        company: 'Company',
        resources: 'Resources',
        legal: 'Legal',
        features: 'Features',
        pricing: 'Pricing',
        about: 'About',
        integrations: 'Integrations',
        api: 'API',
        blog: 'Blog',
        careers: 'Careers',
        contact: 'Contact',
        documentation: 'Documentation',
        helpCenter: 'Help Center',
        status: 'Status',
        security: 'Security',
        privacyPolicy: 'Privacy Policy',
        termsOfService: 'Terms of Service',
        cookiePolicy: 'Cookie Policy',
        compliance: 'Compliance',
        copyright: '© {{year}} MediSync. AI-Powered Conversational BI & Intelligent Accounting for Healthcare.',
      },
    },
    social: {
      twitter: 'Twitter',
      linkedin: 'LinkedIn',
      github: 'GitHub',
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
      smartReports: {
        title: 'التقارير الذكية',
        description: 'تقارير MIS جاهزة للرعاية الصحية. أنشئ لوحات معلومات مخصصة بدون برمجة.',
      },
      deepAnalytics: {
        title: 'التحليلات المتقدمة',
        description: 'محلل ذكاء اصطناعي مستقل يجري الأبحاث. احصل على رؤى قابلة للتطبيق.',
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
      home: 'الرئيسية',
      dashboard: 'لوحة التحكم',
      chat: 'المحادثة',
      alerts: 'التنبيهات',
      reports: 'التقارير',
      settings: 'الإعدادات',
      features: 'المميزات',
      pricing: 'الأسعار',
      about: 'من نحن',
      toggleMenu: 'قائمة الجوال',
    },
    home: {
      hero: {
        badge: 'ذكاء الرعاية الصحية المدعوم بالذكاء الاصطناعي',
        title: 'بياناتك، ',
        titleHighlight: 'مفهومة',
        subtitle: 'اطرح أسئلة بلغة بسيطة. احصل على رؤى فورية من بيانات HIMS وTally. بدون SQL. بدون جداول.',
        ctaChat: 'ابدأ المحادثة',
        ctaDashboard: 'عرض لوحة التحكم',
        trustBy: 'موثوق من منظمات الرعاية الصحية حول العالم',
        hipaa: 'متوافق مع HIPAA',
        soc2: 'شهادة SOC 2',
        clinics: 'أكثر من 500 عيادة',
        queries: 'أكثر من 10 ملايين استعلام',
      },
      preview: {
        askAnything: 'اسأل أي شيء',
        askDesc: 'استعلامات بلغة طبيعية',
        revenueQuestion: 'ما إيرادات اليوم؟',
        instantResponse: 'استجابة فورية من الذكاء الاصطناعي',
        financialInsights: 'رؤى مالية',
        tallySync: 'مزامنة Tally ERP',
        outstanding: 'مستحق',
        collected: 'محصّل',
        autoSynced: 'مزامنة تلقائية مع Tally',
        patientMetrics: 'مقاييس المرضى',
        himsIntegration: 'تكامل HIMS',
        today: 'اليوم',
        vsYesterday: 'مقارنة بالأمس',
        thisMonth: 'هذا الشهر',
        depts: 'الأقسام',
        realTime: 'مباشر من HIMS',
      },
      section: {
        title: 'كل ما تحتاجه',
        subtitle: 'من الاستعلامات المحادثية إلى المحاسبة الآلية، يربط ميدي سنك بيانات الرعاية الصحية بطرق لم تتخيلها.',
      },
      footer: {
        product: 'المنتج',
        company: 'الشركة',
        resources: 'الموارد',
        legal: 'قانوني',
        features: 'المميزات',
        pricing: 'الأسعار',
        about: 'من نحن',
        integrations: 'التكاملات',
        api: 'واجهة برمجة التطبيقات',
        blog: 'المدونة',
        careers: 'الوظائف',
        contact: 'اتصل بنا',
        documentation: 'التوثيق',
        helpCenter: 'مركز المساعدة',
        status: 'الحالة',
        security: 'الأمان',
        privacyPolicy: 'سياسة الخصوصية',
        termsOfService: 'شروط الخدمة',
        cookiePolicy: 'سياسة ملفات التعريف',
        compliance: 'الامتثال',
        copyright: '© {{year}} ميدي سنك. ذكاء الأعمال المحادثي والمحاسبة الذكية للرعاية الصحية.',
      },
    },
    social: {
      twitter: 'تويتر',
      linkedin: 'لينكد إن',
      github: 'جيت هب',
    },
  },
}

// Detect initial language (order: stored preference → URL → Accept-Language → default)
const detectInitialLanguage = (): string => {
  const storedLang = localStorage.getItem('medisync-locale')
  if (storedLang === 'ar' || storedLang === 'en') {
    return storedLang
  }
  const urlParams = new URLSearchParams(window.location.search)
  const urlLang = urlParams.get('lang')
  if (urlLang === 'ar' || urlLang === 'en') {
    return urlLang
  }
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('ar')) {
    return 'ar'
  }
  return 'en'
}

// Merge bundled namespaces into full resources (no HTTP backend; Vite bundles JSON)
const enFull = {
  ...enResources,
  chat: enChat as Record<string, unknown>,
  dashboard: enDashboard as Record<string, unknown>,
  common: enCommon as Record<string, unknown>,
  alerts: enAlerts as Record<string, unknown>,
  reports: enReports as Record<string, unknown>,
  council: enCouncil as Record<string, unknown>,
  copilot: enCopilot as Record<string, unknown>,
}

const arFull = {
  ...arResources,
  chat: arChat as Record<string, unknown>,
  dashboard: arDashboard as Record<string, unknown>,
  common: arCommon as Record<string, unknown>,
  alerts: arAlerts as Record<string, unknown>,
  reports: arReports as Record<string, unknown>,
  council: arCouncil as Record<string, unknown>,
  copilot: arCopilot as Record<string, unknown>,
}

void i18n.use(initReactI18next).init({
  resources: {
    en: enFull,
    ar: arFull,
  },
  lng: detectInitialLanguage(),
  fallbackLng: 'en',
  debug: import.meta.env.DEV,

  ns: ['translation', 'common', 'chat', 'dashboard', 'alerts', 'reports', 'council', 'copilot'],
  defaultNS: 'translation',
  fallbackNS: 'translation',

  interpolation: {
    escapeValue: false, // React already escapes values
  },
  parseMissingKeyHandler: (key: string) => (import.meta.env.DEV ? `[missing:${key}]` : key),

  react: {
    useSuspense: true,
  },

  saveMissing: false,
})

// Save language preference on change
i18n.on('languageChanged', (lng) => {
  localStorage.setItem('medisync-locale', lng)
})

/** Map app language code to BCP 47 locale tag for Intl (number, date, currency) formatting. */
export const APP_LOCALE_TO_BCP47: Record<string, string> = {
  en: 'en-US',
  ar: 'ar-SA',
}

/** Resolve BCP 47 locale from app language code (e.g. en → en-US, ar → ar-SA). */
export function getBcp47Locale(appLocale: string): string {
  return APP_LOCALE_TO_BCP47[appLocale] ?? appLocale
}

export default i18n
