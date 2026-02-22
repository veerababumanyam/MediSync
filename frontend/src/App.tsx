import { Suspense, useCallback, useEffect, useMemo, useState } from 'react'
import { CopilotKit } from '@copilotkit/react-core'
import { useTranslation } from 'react-i18next'
import { useTheme } from './hooks/useTheme'
import './i18n'
import './styles/globals.css'

// Components
import { AppHeader } from './components/common/AppHeader'

// Lazy load page components
import { ChatPage } from './pages/ChatPage'
import { DashboardPage } from './pages/DashboardPage'

/**
 * Route type definition
 */
type Route = 'home' | 'chat' | 'dashboard'

/**
 * Get current route from URL
 */
function getCurrentRoute(): Route {
  const path = window.location.pathname
  if (path === '/chat') return 'chat'
  if (path === '/dashboard') return 'dashboard'
  return 'home'
}

/**
 * MediSync Main Application Component
 */
function AppContent() {
  const { i18n, t } = useTranslation()
  const { isDark, toggleTheme } = useTheme()
  const currentLocale = i18n.language
  const isRTL = currentLocale === 'ar'
  const [currentRoute, setCurrentRoute] = useState<Route>(() => getCurrentRoute())

  useEffect(() => {
    const handlePopState = () => setCurrentRoute(getCurrentRoute())
    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [])

  useEffect(() => {
    document.documentElement.dir = isRTL ? 'rtl' : 'ltr'
    document.documentElement.lang = currentLocale
  }, [currentLocale, isRTL])

  const toggleLanguage = useCallback(() => {
    const newLocale = currentLocale === 'en' ? 'ar' : 'en'
    i18n.changeLanguage(newLocale)
  }, [currentLocale, i18n])

  const navigateTo = useCallback((route: Route) => {
    const path = route === 'home' ? '/' : `/${route}`
    window.history.pushState({}, '', path)
    setCurrentRoute(route)
  }, [])

  const copilotConfig = useMemo(() => ({
    runtimeUrl: import.meta.env.VITE_COPILOT_API_URL || '/api/copilotkit',
  }), [])

  const renderRoute = () => {
    switch (currentRoute) {
      case 'chat':
        return <ChatPage isDark={isDark} />
      case 'dashboard':
        return <DashboardPage isDark={isDark} navigateTo={navigateTo} />
      default:
        return (
          <HomePageContent
            isRTL={isRTL}
            currentLocale={currentLocale}
            navigateTo={navigateTo}
            isDark={isDark}
          />
        )
    }
  }

  return (
    <CopilotKit {...copilotConfig}>
      <div
        className={`min-h-screen relative overflow-hidden transition-colors duration-500 ${isRTL ? 'rtl' : 'ltr'
          }`}
      >
        {/* Background — adapts to theme */}
        <div className={`fixed inset-0 -z-20 transition-colors duration-500 ${isDark
          ? 'bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900'
          : 'bg-gradient-to-br from-slate-50 via-blue-50/30 to-slate-100'
          }`} />
        {/* Animated Orbs (both modes, different opacity) */}
        <div className="fixed inset-0 -z-10 overflow-hidden">
          <div className={`orb orb-blue top-[-10%] left-[-5%] ${isDark ? 'opacity-100' : 'opacity-20'}`} />
          <div className={`orb orb-purple bottom-[-15%] right-[-10%] ${isDark ? 'opacity-100' : 'opacity-15'}`} />
          <div className={`orb orb-pink top-[40%] left-[60%] ${isDark ? 'opacity-100' : 'opacity-10'}`} />
        </div>

        {/* Shared Header */}
        <AppHeader
          isDark={isDark}
          toggleTheme={toggleTheme}
          toggleLanguage={toggleLanguage}
          currentLocale={currentLocale}
          navigateTo={navigateTo}
          currentRoute={currentRoute}
        />

        {/* Page Content */}
        {renderRoute()}

        {/* Footer */}
        <footer className={`border-t mt-16 relative z-10 transition-colors duration-300 ${isDark ? 'border-white/10' : 'border-slate-200'
          }`}>
          <div className={`max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 text-center text-sm ${isDark ? 'text-slate-500' : 'text-slate-400'
            }`}>
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
 * Home Page Content — Premium Glassmorphism Design with Dark/Light support
 * (Only the content below the header, since AppHeader is rendered in the shared layout)
 */
interface HomePageContentProps {
  isRTL: boolean
  currentLocale: string
  navigateTo: (route: Route) => void
  isDark: boolean
}

function HomePageContent({ currentLocale, isDark }: HomePageContentProps) {
  const { t } = useTranslation()

  return (
    <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16 relative z-10">
      {/* Hero Section */}
      <section className="mb-16 text-center animate-fade-in-up">
        <div className={`inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold tracking-wider uppercase mb-6 ${isDark
          ? 'bg-blue-500/10 border border-blue-500/20 text-blue-400'
          : 'bg-blue-50 border border-blue-200 text-blue-600'
          }`}>
          <div className={`w-1.5 h-1.5 rounded-full animate-pulse ${isDark ? 'bg-blue-400' : 'bg-blue-600'}`} />
          {t('welcome.badge', 'AI-Powered Platform')}
        </div>

        {/* Logo in hero */}
        <div className="flex justify-center mb-6">
          <img
            src="/logo.png"
            alt="MediSync"
            className="w-16 h-16 sm:w-20 sm:h-20 object-contain"
          />
        </div>

        <h2 className={`text-4xl sm:text-5xl lg:text-6xl font-extrabold mb-6 leading-tight tracking-tight ${isDark ? 'text-white' : 'text-slate-900'
          }`}>
          {t('welcome.title', 'Welcome to MediSync')}
        </h2>
        <p className={`text-lg sm:text-xl max-w-2xl mx-auto leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'
          }`}>
          {t('welcome.description',
            'Ask questions in plain language and get instant charts, tables, and insights from your healthcare and accounting data.'
          )}
        </p>
      </section>

      {/* Feature Cards */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
        <FeatureCard
          isDark={isDark}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 01-.825-.242m9.345-8.334a2.126 2.126 0 00-.476-.095 48.64 48.64 0 00-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0011.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155" />
            </svg>
          }
          gradient="from-blue-500 to-cyan-400"
          shadowColor="shadow-blue-500/20"
          title={t('features.conversationalBI.title', 'Conversational BI')}
          description={t('features.conversationalBI.description',
            'Chat with your data using natural language. Get instant visualizations.'
          )}
          delay="delay-1"
        />
        <FeatureCard
          isDark={isDark}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23.693L5 14.5m14.8.8l1.402 1.402c1.232 1.232.65 3.318-1.067 3.611A48.309 48.309 0 0112 21c-2.773 0-5.491-.235-8.135-.687-1.718-.293-2.3-2.379-1.067-3.611L5 14.5" />
            </svg>
          }
          gradient="from-emerald-500 to-teal-400"
          shadowColor="shadow-emerald-500/20"
          title={t('features.aiAccountant.title', 'AI Accountant')}
          description={t('features.aiAccountant.description',
            'Upload documents and let AI extract, map, and sync to your accounting system.'
          )}
          delay="delay-2"
        />
        <FeatureCard
          isDark={isDark}
          icon={
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
            </svg>
          }
          gradient="from-purple-500 to-pink-400"
          shadowColor="shadow-purple-500/20"
          title={t('features.easyReports.title', 'Easy Reports')}
          description={t('features.easyReports.description',
            'Pre-built reports and custom dashboards with automated delivery.'
          )}
          delay="delay-3"
        />
      </section>

      {/* System Status */}
      <section className={`rounded-2xl p-6 sm:p-8 animate-fade-in-up delay-4 transition-colors duration-300 ${isDark
        ? 'glass glass-shine'
        : 'bg-white border border-slate-200 shadow-sm'
        }`}>
        <h3 className={`text-lg font-semibold mb-6 flex items-center gap-2 ${isDark ? 'text-white' : 'text-slate-900'
          }`}>
          <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          {t('status.title', 'System Status')}
        </h3>
        <div className="space-y-1">
          <StatusItem isDark={isDark} name={t('status.react', 'React')} version="19.2.4" status="ready" />
          <StatusItem isDark={isDark} name={t('status.vite', 'Vite')} version="7.3.1" status="ready" />
          <StatusItem isDark={isDark} name={t('status.copilotkit', 'CopilotKit')} version="1.3.6" status="ready" />
          <StatusItem isDark={isDark} name={t('status.i18n', 'i18n')} version={currentLocale === 'en' ? 'English (LTR)' : 'Arabic (RTL)'} status="ready" />
        </div>
      </section>
    </main>
  )
}

/**
 * Feature Card Component
 */
function FeatureCard({
  icon,
  gradient,
  shadowColor,
  title,
  description,
  delay = '',
  isDark,
}: {
  icon: React.ReactNode
  gradient: string
  shadowColor: string
  title: string
  description: string
  delay?: string
  isDark: boolean
}) {
  return (
    <div className={`rounded-2xl p-6 hover:-translate-y-1 transition-all duration-300 group animate-fade-in-up ${delay} ${isDark
      ? 'glass glass-shine hover:border-white/30'
      : 'bg-white border border-slate-200 shadow-sm hover:shadow-md hover:border-blue-200'
      }`}>
      <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${gradient} flex items-center justify-center mb-5 text-white shadow-lg ${shadowColor} group-hover:scale-110 transition-transform duration-300`}>
        {icon}
      </div>
      <h3 className={`text-lg font-semibold mb-2 ${isDark ? 'text-white' : 'text-slate-900'}`}>
        {title}
      </h3>
      <p className={`text-sm leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
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
  isDark,
}: {
  name: string
  version: string
  status: 'ready' | 'loading' | 'error'
  isDark: boolean
}) {
  const statusConfig = {
    ready: { color: 'bg-emerald-400', ring: 'ring-emerald-400/20' },
    loading: { color: 'bg-amber-400', ring: 'ring-amber-400/20' },
    error: { color: 'bg-rose-400', ring: 'ring-rose-400/20' },
  }
  const { color, ring } = statusConfig[status]

  return (
    <div className={`flex items-center justify-between py-3 px-4 rounded-xl transition-colors duration-200 ${isDark ? 'hover:bg-white/5' : 'hover:bg-slate-50'
      }`}>
      <div className="flex items-center gap-3">
        <div className={`w-2.5 h-2.5 rounded-full ${color} ring-4 ${ring}`} />
        <span className={`font-medium ${isDark ? 'text-slate-300' : 'text-slate-700'}`}>{name}</span>
      </div>
      <span className={`text-sm font-mono ${isDark ? 'text-slate-500' : 'text-slate-400'}`}>
        {version}
      </span>
    </div>
  )
}

/**
 * App Root with Suspense boundary
 */
export default function App() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <AppContent />
    </Suspense>
  )
}

/**
 * Loading fallback
 */
function LoadingFallback() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="text-center">
        <img src="/logo.png" alt="MediSync" className="w-16 h-16 mx-auto mb-4 animate-pulse" />
        <p className="text-slate-400 text-sm font-medium">Loading MediSync...</p>
      </div>
    </div>
  )
}
