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
      tagline: 'Turn Any Legacy Healthcare System into Conversational AI',
      toggleLanguage: 'Toggle language',
    },
    welcome: {
      badge: 'The Agentic AI Bridge for Legacy Healthcare IT',
    },
    heroCarousel: {
      slide1: {
        title: "Don't Replace Your Software. Make It Speak.",
        description: 'Transform your legacy HIMS, LIMS, Accounting, and Custom Databases into a conversational AI interface — without changing a single line of code.',
        cta: 'Get Started Free',
        stat1: '10x Faster Insights',
        stat2: 'Zero Migration',
        stat3: '99.9% Uptime',
      },
      slide2: {
        title: 'Zero Code. Zero Migration. Instant AI.',
        description: 'Our Agentic AI layers securely over your existing infrastructure. Ask questions in plain English and get instant answers from any legacy database.',
        cta: 'See It In Action',
        stat1: '50+ Integrations',
        stat2: '< 2 Min Setup',
        stat3: '24/7 AI Agents',
      },
      slide3: {
        title: 'Your Legacy Systems. Supercharged.',
        description: 'Stop ripping and replacing. Keep the systems your teams already know — and let our AI unlock instant analytics, smart accounting, and prescriptive insights.',
        cta: 'Book a Demo',
        stat1: '₹2Cr+ Saved Avg.',
        stat2: 'HIPAA Compliant',
        stat3: 'Custom Integrations',
      },
      trustStrip: {
        item1: 'HIMS',
        item2: 'LIMS',
        item3: 'Tally ERP',
        item4: 'SQL Databases',
        item5: 'Custom APIs',
        item6: 'Oracle',
        item7: 'SAP',
        item8: 'REST / GraphQL',
      },
    },
    sectors: {
      title: 'Dominating Complexity Across Every Healthcare Sector',
      description: 'We don\'t just understand data; we understand the business of healthcare. Our tailored Agentic AI bridges seamlessly adapt to the unique reporting, compliance, and velocity requirements of your specific vertical.',
      hospitals: {
        title: 'Large Hospital Networks',
        description: 'Unify fragmented data lakes, consolidate multi-branch P&L, and visualize system-wide operational efficiency without altering your legacy HIMS.',
      },
      labs: {
        title: 'Clinical Laboratories',
        description: 'Drive margins by correlating test volumes with reagent costs, tracking turnaround times, and optimizing your LIMS supply chains.',
      },
      pharmacies: {
        title: 'Enterprise Pharmacies',
        description: 'Analyze prescription trends, automate inventory forecasting, and instantly map supplier invoices to Tally ERP.',
      },
      clinics: {
        title: 'Specialized Clinics',
        description: 'Maximize chair time, track provider revenue generation, and automatically sync daily receipts to your accounting software.',
      },
    },
    features: {
      conversationalBI: {
        title: 'Seamless Natural Language to Tech',
        description: 'Stop digging through clunky menus. Just type "Show me last month\'s lab revenue" and our AI agents instantly write the SQL, pull the data, and render beautiful reports from your old systems.',
      },
      tallySync: {
        title: 'Universal Legacy Connectivity',
        description: 'HIMS, LIMS, Tally ERP, or a custom 20-year-old database? We integrate with them all. We build custom APIs based on your exact requirements so no system is left behind.',
      },
      aiAccountant: {
        title: '100% Zero Rip-and-Replace',
        description: 'The ultimate USP: Keep the software you already know and own. We act as an invisible, intelligent brain on top of your current stack, saving you millions in migration costs.',
      },
      piiProtection: {
        title: 'Bank-Grade Agentic Security',
        description: 'Our AI translates your questions into code, not the other way around. Built-in PII scrubbers ensure your patient data never leaves your secure environment during processing.',
      },
      prescriptiveAnalytics: {
        title: 'Instantly Modernized Insights',
        description: 'Turn a dusty, traditional database into a proactive advisor. Ask your old software what to do next, and get prescriptive AI recommendations to plug revenue leaks.',
      },
      himsConnectivity: {
        title: 'Custom Integrations on Demand',
        description: 'Need to connect a specialized lab machine or an ancient ERP? Our team continuously develops custom integrations so our AI agents can converse with any data source on Earth.',
      },
    },
    faq: {
      title: 'Frequently Asked Questions',
      q1: 'Do I need to change or upgrade my current software?',
      a1: 'Absolutely not. This is the entire purpose of MediSync. Whether you use HIMS, LIMS, or older accounting systems, we plug into what you already have. You keep your systems; we just make them smart, conversational, and instantly accessible.',
      q2: 'How does it actually work?',
      a2: 'You type a question in natural language. Our specialized backend AI Agents instantly translate your question into technical code (like SQL or API calls), query your old databases in real-time, and return a beautiful visual report. It\'s like having a senior data engineer on staff 24/7.',
      q3: 'Can it connect to my highly specific, custom-built system?',
      a3: 'Yes. We pride ourselves on universal connectivity. We actively develop custom integrations tailored strictly to our customers\' requirements. If you have data, our AI can talk to it.',
    },
    status: {
      title: 'Platform Real-Time Status',
      react: 'React',
      vite: 'Vite',
      copilotkit: 'CopilotKit',
      i18n: 'i18n',
    },
    footer: {
      copyright:
        '© 2026 MediSync. The World\'s Smartest AI-Powered Conversational BI & Intelligent Accounting Platform for Healthcare.',
    },
    common: {
      loading: 'Loading Excellence...',
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
      chat: 'AI Chat',
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
      tagline: 'حول أي نظام صحي قديم إلى ذكاء اصطناعي محادثي',
      toggleLanguage: 'تبديل اللغة',
    },
    welcome: {
      badge: 'جسر الذكاء الاصطناعي لتكنولوجيا الرعاية الصحية القديمة',
    },
    heroCarousel: {
      slide1: {
        title: 'لا تستبدل برامجك. اجعلها تتحدث.',
        description: 'حوّل أنظمة HIMS وLIMS والمحاسبة وقواعد البيانات المخصصة إلى واجهة ذكاء اصطناعي محادثة — بدون تغيير سطر برمجي واحد.',
        cta: 'ابدأ مجاناً',
        stat1: 'رؤى أسرع 10 مرات',
        stat2: 'بدون ترحيل',
        stat3: '99.9% وقت تشغيل',
      },
      slide2: {
        title: 'بدون برمجة. بدون ترحيل. ذكاء اصطناعي فوري.',
        description: 'يتصل ذكاؤنا الاصطناعي بأمان فوق بنيتك التحتية الحالية. اطرح أسئلة بلغتك واحصل على إجابات فورية من أي قاعدة بيانات.',
        cta: 'شاهد العرض',
        stat1: '+50 تكامل',
        stat2: 'إعداد < دقيقتين',
        stat3: 'وكلاء ذكاء 24/7',
      },
      slide3: {
        title: 'أنظمتك القديمة. بقوة خارقة.',
        description: 'توقف عن الاستبدال. احتفظ بالأنظمة التي يعرفها فريقك — ودع ذكاءنا الاصطناعي يطلق التحليلات الفورية والمحاسبة الذكية.',
        cta: 'احجز عرضاً',
        stat1: 'توفير +2 كرور متوسط',
        stat2: 'متوافق مع HIPAA',
        stat3: 'تكاملات مخصصة',
      },
      trustStrip: {
        item1: 'HIMS',
        item2: 'LIMS',
        item3: 'Tally ERP',
        item4: 'قواعد SQL',
        item5: 'واجهات مخصصة',
        item6: 'Oracle',
        item7: 'SAP',
        item8: 'REST / GraphQL',
      },
    },
    sectors: {
      title: 'السيطرة على التعقيد في كل قطاع صحي',
      description: 'نحن لا نفهم البيانات فحسب؛ بل نفهم أعمال الرعاية الصحية. تتكيف جسور الذكاء الاصطناعي الخاصة بنا بسلاسة مع المتطلبات الفريدة لقطاعك.',
      hospitals: {
        title: 'شبكات المستشفيات الكبيرة',
        description: 'توحيد بحيرات البيانات المجزأة، ودمج الأرباح والخسائر للفروع، وتصور الكفاءة التشغيلية دون تغيير أنظمتك القديمة.',
      },
      labs: {
        title: 'المختبرات السريرية',
        description: 'زيادة هوامش الربح من خلال ربط أحجام الاختبارات بتكاليف الكواشف، وتتبع أوقات الإنجاز، وتحسين سلاسل التوريد الخاصة بالمختبر.',
      },
      pharmacies: {
        title: 'صيدليات المؤسسات',
        description: 'تحليل اتجاهات الوصفات الطبية، وأتمتة التنبؤ بالمخزون، وربط فواتير الموردين فوراً ببرنامج المحاسبة.',
      },
      clinics: {
        title: 'العيادات المتخصصة',
        description: 'زيادة كفاءة المواعيد، وتتبع الإيرادات التي يولدها مزود الخدمة، ومزامنة الإيصالات اليومية تلقائيًا.',
      },
    },
    features: {
      conversationalBI: {
        title: 'من لغة طبيعية إلى تقنية بسلاسة',
        description: 'توقف عن البحث في القوائم المعقدة. فقط اكتب "أرني إيرادات المختبر للشهر الماضي" وسيقوم وكلاء الذكاء الاصطناعي لدينا بكتابة الكود فوراً واستخراج البيانات من أنظمتك القديمة.',
      },
      tallySync: {
        title: 'اتصال شامل بالأنظمة القديمة',
        description: 'سواء كان HIMS، LIMS، Tally ERP، أو قاعدة بيانات عمرها 20 عاماً، نحن نتكامل معها جميعاً. نبني واجهات برمجة مخصصة بناءً على متطلباتك الدقيقة.',
      },
      aiAccountant: {
        title: '100% بدون استبدال',
        description: 'الميزة التنافسية الكبرى: احتفظ بالبرامج التي تعرفها وتملكها. نحن نعمل كعقل ذكي غير مرئي مهيمن على نظامك الحالي، مما يوفر لك الملايين في تكاليف الترحيل.',
      },
      piiProtection: {
        title: 'أمان بذكاء اصطناعي بمستوى البنوك',
        description: 'يقوم الذكاء الاصطناعي لدينا بترجمة أسئلتك إلى رموز برمجية. تضمن أدوات إخفاء الهوية المدمجة بقاء بيانات مرضاك داخل بيئتك الآمنة.',
      },
      prescriptiveAnalytics: {
        title: 'رؤى حديثة فورية',
        description: 'حول قاعدة البيانات التقليدية إلى مستشار استباقي. اسأل برامجك القديمة عما يجب فعله تالياً، واحصل على توصيات ذكية لسد تسرب الإيرادات.',
      },
      himsConnectivity: {
        title: 'تكاملات مخصصة عند الطلب',
        description: 'هل تحتاج إلى توصيل آلة مختبر متخصصة أو نظام تخطيط قديم؟ يقوم فريقنا بتطوير واجهات مخصصة باستمرار لتلبية أي متطلب لعملائنا.',
      },
    },
    faq: {
      title: 'الأسئلة الشائعة',
      q1: 'هل أحتاج إلى تغيير أو ترقية برامجي الحالية؟',
      a1: 'بالتأكيد لا. هذا هو الغرض الأساسي من ميدي سنك. سواء كنت تستخدم HIMS أو LIMS أو أنظمة قديمة، نحن نتصل بما لديك بالفعل لنجعله ذكياً وسهل الوصول.',
      q2: 'كيف يعمل النظام بالفعل؟',
      a2: 'تكتب سؤالك بلغة طبيعية، فيقوم وكلاء الذكاء الاصطناعي لدينا بترجمته فوراً إلى استعلامات برمجية لقواعد بياناتك القديمة وإعادة الإجابة بصرياً في ثوانٍ. إنه كأن تملك مهندس بيانات خبير يعمل 24/7.',
      q3: 'هل يمكنكم الاتصال بنظامي الخاص والمعقد جداً؟',
      a3: 'نعم. نحن فخورون بقدرتنا على الاتصال الشامل. نقوم بتطوير تكاملات مخصصة تتناسب تماماً مع متطلبات عملائنا. إذا كانت لديك بيانات، يمكن لذكائنا الاصطناعي التحدث إليها.',
    },
    status: {
      title: 'حالة المنصة في الوقت الفعلي',
      react: 'رياكت',
      vite: 'فايت',
      copilotkit: 'كوبيلوكت كيت',
      i18n: 'دعم اللغات',
    },
    footer: {
      copyright:
        '© 2026 ميدي سنك. منصة الذكاء التجاري والمحاسبة الأذكى في العالم للرعاية الصحية.',
    },
    common: {
      loading: 'جاري تحميل التميز...',
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
      chat: 'محادثة الذكاء الاصطناعي',
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
