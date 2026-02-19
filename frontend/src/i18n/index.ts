import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

/**
 * i18n Configuration for MediSync
 *
 * Supports:
 * - English (en) - Left-to-Right (LTR)
 * - Arabic (ar) - Right-to-Left (RTL)
 *
 * Locale Detection Priority:
 * 1. localStorage user preference
 * 2. URL parameter ?lang=ar or ?lang=en
 * 3. Browser's Accept-Language header
 * 4. Default: 'en'
 */

// English translations
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
  },
}

// Arabic translations
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

// Initialize i18next
void i18n.use(initReactI18next).init({
  resources: {
    en: enResources,
    ar: arResources,
  },
  lng: detectInitialLanguage(),
  fallbackLng: 'en',
  debug: import.meta.env.DEV,

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
