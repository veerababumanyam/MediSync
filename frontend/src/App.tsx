import { Suspense, useCallback, useEffect, useMemo } from 'react'
import { CopilotKit } from '@copilotkit/react-core'
import { useTranslation } from 'react-i18next'
import './i18n'
import './styles/globals.css'

/**
 * MediSync Main Application Component
 *
 * Features:
 * - CopilotKit integration for generative UI
 * - i18n support for English (LTR) and Arabic (RTL)
 * - Automatic RTL layout based on locale
 * - Error boundary and loading states
 */
function AppContent() {
  const { t, i18n } = useTranslation()
  const currentLocale = i18n.language
  const isRTL = currentLocale === 'ar'

  // Update document direction and language when locale changes
  useEffect(() => {
    document.documentElement.dir = isRTL ? 'rtl' : 'ltr'
    document.documentElement.lang = currentLocale
  }, [currentLocale, isRTL])

  // Language toggle handler
  const toggleLanguage = useCallback(() => {
    const newLocale = currentLocale === 'en' ? 'ar' : 'en'
    i18n.changeLanguage(newLocale)
  }, [currentLocale, i18n])

  // CopilotKit configuration
  const copilotConfig = useMemo(() => ({
    endpoint: import.meta.env.VITE_COPILOT_API_URL || '/api/copilotkit',
  }), [])

  return (
    <CopilotKit {...copilotConfig}>
      <div
        className={`min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 ${
          isRTL ? 'rtl' : 'ltr'
        }`}
      >
        {/* Header */}
        <header className="border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm">
          <div className="container mx-auto px-4 py-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 flex items-center justify-center">
                <span className="text-white font-bold text-xl">M</span>
              </div>
              <div>
                <h1 className="text-xl font-bold text-slate-900 dark:text-white">
                  {t('app.name', 'MediSync')}
                </h1>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  {t('app.tagline', 'AI-Powered Business Intelligence')}
                </p>
              </div>
            </div>

            <button
              onClick={toggleLanguage}
              className="px-4 py-2 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 transition-colors text-sm font-medium text-slate-700 dark:text-slate-300"
              aria-label={t('app.toggleLanguage', 'Toggle language')}
            >
              {currentLocale === 'en' ? 'عربي' : 'English'}
            </button>
          </div>
        </header>

        {/* Main Content */}
        <main className="container mx-auto px-4 py-8">
          {/* Welcome Section */}
          <section className="mb-12 text-center">
            <h2 className="text-4xl font-bold text-slate-900 dark:text-white mb-4">
              {t('welcome.title', 'Welcome to MediSync')}
            </h2>
            <p className="text-lg text-slate-600 dark:text-slate-400 max-w-2xl mx-auto">
              {t('welcome.description',
                'Ask questions in plain language and get instant charts, tables, and insights from your healthcare and accounting data.'
              )}
            </p>
          </section>

          {/* Feature Cards */}
          <section className="grid md:grid-cols-3 gap-6 mb-12">
            <FeatureCard
              icon="💬"
              title={t('features.conversationalBI.title', 'Conversational BI')}
              description={t('features.conversationalBI.description',
                'Chat with your data using natural language. Get instant visualizations.'
              )}
            />
            <FeatureCard
              icon="🤖"
              title={t('features.aiAccountant.title', 'AI Accountant')}
              description={t('features.aiAccountant.description',
                'Upload documents and let AI extract, map, and sync to your accounting system.'
              )}
            />
            <FeatureCard
              icon="📊"
              title={t('features.easyReports.title', 'Easy Reports')}
              description={t('features.easyReports.description',
                'Pre-built reports and custom dashboards with automated delivery.'
              )}
            />
          </section>

          {/* Setup Status */}
          <section className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
            <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-4">
              {t('status.title', 'System Status')}
            </h3>
            <div className="space-y-3">
              <StatusItem
                name={t('status.react', 'React')}
                version="19.2.4"
                status="ready"
              />
              <StatusItem
                name={t('status.vite', 'Vite')}
                version="7.3.1"
                status="ready"
              />
              <StatusItem
                name={t('status.copilotkit', 'CopilotKit')}
                version="1.3.6"
                status="ready"
              />
              <StatusItem
                name={t('status.i18n', 'i18n')}
                version={currentLocale === 'en' ? 'English (LTR)' : 'Arabic (RTL)'}
                status="ready"
              />
            </div>
          </section>
        </main>

        {/* Footer */}
        <footer className="border-t border-slate-200 dark:border-slate-700 mt-12">
          <div className="container mx-auto px-4 py-6 text-center text-sm text-slate-500 dark:text-slate-400">
            <p>
              {t('footer.copyright', '© 2026 MediSync. AI-Powered Conversational BI & Intelligent Accounting for Healthcare.')}
            </p>
          </div>
        </footer>
      </div>
    </CopilotKit>
  )
}

/**
 * Feature Card Component
 */
function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: string
  title: string
  description: string
}) {
  return (
    <div className="bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700 hover:shadow-md transition-shadow">
      <div className="text-4xl mb-4">{icon}</div>
      <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-2">
        {title}
      </h3>
      <p className="text-slate-600 dark:text-slate-400">
        {description}
      </p>
    </div>
  )
}

/**
 * Status Item Component
 */
function StatusItem({
  name,
  version,
  status,
}: {
  name: string
  version: string
  status: 'ready' | 'loading' | 'error'
}) {
  const statusColors = {
    ready: 'bg-emerald-500',
    loading: 'bg-amber-500',
    error: 'bg-rose-500',
  }

  return (
    <div className="flex items-center justify-between py-2">
      <div className="flex items-center gap-3">
        <div className={`w-2 h-2 rounded-full ${statusColors[status]}`} />
        <span className="font-medium text-slate-700 dark:text-slate-300">{name}</span>
      </div>
      <span className="text-sm text-slate-500 dark:text-slate-400 font-mono">
        {version}
      </span>
    </div>
  )
}

/**
 * App Root with Suspense boundary for i18n
 */
export default function App() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <AppContent />
    </Suspense>
  )
}

/**
 * Loading fallback component
 */
function LoadingFallback() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800">
      <div className="text-center">
        <div className="w-16 h-16 rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 flex items-center justify-center mx-auto mb-4">
          <span className="text-white font-bold text-3xl">M</span>
        </div>
        <p className="text-slate-600 dark:text-slate-400">Loading MediSync...</p>
      </div>
    </div>
  )
}
