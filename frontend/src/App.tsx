import { Suspense, useCallback, useEffect, useMemo, useState } from 'react'
import { CopilotKit } from '@copilotkit/react-core'
import { useTranslation } from 'react-i18next'
import './i18n'
import './styles/globals.css'

// Lazy load page components
import { ChatPage } from './pages/ChatPage'
import { DashboardPage } from './pages/DashboardPage'

// CopilotKit components
import { MediSyncCopilot, CopilotFloatingButton } from './components/copilot'
import { GlassCard, ThemeToggle, AnimatedBackground } from './components/ui'
import { FadeIn, StaggerChildren } from './components/animations'
import { ThemeProvider } from './components/theme'
import { webMCPService } from './services/WebMCPService'

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
 *
 * Features:
 * - CopilotKit integration for generative UI
 * - WebMCP integration for AI agent discoverability
 * - i18n support for English (LTR) and Arabic (RTL)
 * - Automatic RTL layout based on locale
 * - Client-side routing for /chat and /dashboard
 * - Error boundary and loading states
 */
function AppContent() {
  const { i18n } = useTranslation()
  const currentLocale = i18n.language
  const isRTL = currentLocale === 'ar'
  const [currentRoute, setCurrentRoute] = useState<Route>(() => getCurrentRoute())

  // Handle browser navigation (back/forward buttons)
  useEffect(() => {
    const handlePopState = () => {
      setCurrentRoute(getCurrentRoute())
    }

    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [])

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

  // Navigation handlers
  const navigateTo = useCallback((route: Route) => {
    const path = route === 'home' ? '/' : `/${route}`
    window.history.pushState({}, '', path)
    setCurrentRoute(route)
  }, [])

  // Register WebMCP navigation tools (available on all pages)
  useEffect(() => {
    webMCPService.registerNavigationTools({
      onNavigate: (route: string) => {
        navigateTo(route as Route)
      },
      onToggleLanguage: toggleLanguage,
    })

    return () => {
      // Don't cleanup navigation tools on route change
    }
  }, [navigateTo, toggleLanguage])

  // CopilotKit configuration
  // For CopilotKit 1.51+, we need runtimeUrl or publicApiKey
  // During development without backend, we can disable CopilotKit features
  const copilotConfig = useMemo(() => {
    const publicApiKey = import.meta.env.VITE_COPILOT_PUBLIC_API_KEY

    if (publicApiKey) {
      return { publicApiKey }
    }
    // For local development without CopilotKit backend, use a self-hosted runtime
    // This prevents the "Missing required prop" error while allowing the app to load
    return {
      runtimeUrl: '/api/copilotkit',
      // Disable agent features when no backend is available
      agent: undefined,
    }
  }, [])

  // Check if CopilotKit should be enabled
  const isCopilotKitEnabled = !!import.meta.env.VITE_COPILOT_PUBLIC_API_KEY

  // Render route content
  const renderRoute = () => {
    switch (currentRoute) {
      case 'chat':
        return <ChatPage />
      case 'dashboard':
        return <DashboardPage />
      default:
        return (
          <HomePage
            isRTL={isRTL}
            currentLocale={currentLocale}
            toggleLanguage={toggleLanguage}
            navigateTo={navigateTo}
            copilotConfig={copilotConfig}
          />
        )
    }
  }

  // CopilotKit content with MediSync tools
  const copilotContent = (
    <>
      {renderRoute()}
      {/* Floating AI assistant button - visible when CopilotKit is enabled */}
      {isCopilotKitEnabled && <CopilotFloatingButton />}
    </>
  )

  // Wrap with CopilotKit only if enabled
  if (isCopilotKitEnabled) {
    return (
      <CopilotKit {...copilotConfig}>
        <MediSyncCopilot
          currentRoute={currentRoute}
          locale={currentLocale}
          onNavigate={(route) => navigateTo(route as Route)}
          onToggleLanguage={toggleLanguage}
        />
        {copilotContent}
      </CopilotKit>
    )
  }

  // Render without CopilotKit when not configured
  return copilotContent
}

/**
 * Home Page Component
 */
interface HomePageProps {
  isRTL: boolean
  currentLocale: string
  toggleLanguage: () => void
  navigateTo: (route: Route) => void
  copilotConfig: { publicApiKey?: string; runtimeUrl?: string; agent?: undefined }
}

function HomePage({ isRTL, currentLocale, toggleLanguage, navigateTo }: HomePageProps) {
  const { t } = useTranslation()

  // Register WebMCP tools for Home page
  useEffect(() => {
    webMCPService.registerChatTools({
      onQuery: (query: string) => {
        // Navigate to chat with the query
        const encodedQuery = encodeURIComponent(query)
        window.location.href = `/chat?query=${encodedQuery}`
      },
      onSyncTally: async () => {
        console.log('Tally sync requested from home page')
      },
      onShowDashboard: (_id: string) => {
        navigateTo('dashboard')
      },
    })

    return () => {
      // Cleanup is handled by AppContent
    }
  }, [navigateTo])

  return (
    <div
      className={`min-h-screen animated-gradient-bg mesh-gradient-overlay ${
        isRTL ? 'rtl' : 'ltr'
      }`}
      // WebMCP declarative attributes for home page
      {...({
        'tool-name': 'medi-home',
        'tool-description': 'The MediSync home page with feature overview and quick access to Chat and Dashboard',
      } as any)}
    >
        {/* Header - Enhanced Glass Effect with Colorful Border */}
        <header className="border-b border-white/20 dark:border-white/10 bg-white/10 dark:bg-slate-900/30 backdrop-blur-xl sticky top-0 z-50 shadow-glass-md">
          <div className="container mx-auto px-4 py-4 flex items-center justify-between">
            <div
              className="flex items-center gap-3 cursor-pointer group"
              onClick={() => navigateTo('home')}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => e.key === 'Enter' && navigateTo('home')}
            >
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-105 transition-all duration-300">
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

            {/* Navigation Links */}
            <nav className="flex items-center gap-4">
              <button
                onClick={() => navigateTo('chat')}
                className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                  />
                </svg>
                {t('navigation.chat', 'Chat')}
              </button>
              <button
                onClick={() => navigateTo('dashboard')}
                className="inline-flex items-center gap-2 px-4 py-2 text-slate-600 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm font-medium"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
                {t('navigation.dashboard', 'Dashboard')}
              </button>
              <ThemeToggle />
              <button
                onClick={toggleLanguage}
                className="px-4 py-2 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 transition-colors text-sm font-medium text-slate-700 dark:text-slate-300"
                aria-label={t('app.toggleLanguage', 'Toggle language')}
              >
                {currentLocale === 'en' ? 'عربي' : 'English'}
              </button>
            </nav>
          </div>
        </header>

        {/* Main Content - Redesigned Healthcare BI Command Center */}
        <main className="container mx-auto px-4 py-8">
          {/* Hero Section with Glass Effect */}
          <FadeIn>
            <section className="mb-16 text-center">
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/20 dark:bg-white/10 backdrop-blur-md border border-white/30 text-white text-sm font-medium mb-6 shadow-glass-md">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
                </span>
                AI-Powered Healthcare Intelligence
              </div>

              <h1 className="text-5xl md:text-6xl font-bold text-white mb-6 tracking-tight drop-shadow-lg">
                Your Data, <span className="text-transparent bg-clip-text bg-gradient-to-r from-cyan-300 to-pink-300 font-extrabold">Understood</span>
              </h1>

              <p className="text-xl text-white/90 max-w-3xl mx-auto mb-8 leading-relaxed drop-shadow-md">
                Ask questions in plain language. Get instant insights from your HIMS and Tally data.
                No SQL required. No spreadsheets.
              </p>

              {/* CTA Buttons */}
              <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
                <button
                  onClick={() => navigateTo('chat')}
                  className="group inline-flex items-center gap-2 px-8 py-4 bg-gradient-to-r from-blue-600 to-cyan-500 text-white rounded-xl hover:shadow-lg hover:shadow-blue-500/25 hover:-translate-y-0.5 transition-all duration-300 text-lg font-semibold"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  Start Chatting
                  <svg className="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                </button>
                <button
                  onClick={() => navigateTo('dashboard')}
                  className="inline-flex items-center gap-2 px-8 py-4 bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-200 rounded-xl hover:bg-white/80 dark:hover:bg-slate-800/80 transition-all duration-300 text-lg font-semibold"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                  View Dashboard
                </button>
              </div>
            </section>
          </FadeIn>

          {/* Live Preview Section - What the User Can Actually Do */}
          <StaggerChildren className="grid lg:grid-cols-3 gap-6 mb-16">
            {/* Preview Card 1: Conversational BI Query */}
            <GlassCard intensity="light" shadow="lg" hover="blueGlow" className="p-6 bg-white/20 dark:bg-white/5 backdrop-blur-xl border border-blue-400/30 dark:border-blue-400/20">
              <div className="flex items-start gap-4 mb-4">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-400 flex items-center justify-center flex-shrink-0">
                  <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Ask Anything</h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400">Natural language queries</p>
                </div>
              </div>

              {/* Simulated Chat Interface */}
              <div className="space-y-3 mb-4">
                <div className="flex gap-2">
                  <div className="w-6 h-6 rounded-full bg-slate-200 dark:bg-slate-700 flex-shrink-0 flex items-center justify-center">
                    <svg className="w-3 h-3 text-slate-500" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" />
                    </svg>
                  </div>
                  <div className="flex-1 bg-slate-100 dark:bg-slate-800/50 rounded-lg px-3 py-2 text-sm text-slate-700 dark:text-slate-300">
                    "What's today's revenue?"
                  </div>
                </div>
                <div className="flex gap-2">
                  <div className="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-cyan-400 flex-shrink-0 flex items-center justify-center">
                    <span className="text-white text-xs font-bold">M</span>
                  </div>
                  <div className="flex-1 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg px-3 py-2 text-sm text-blue-700 dark:text-blue-300">
                    <div className="flex items-center gap-2">
                      <span>$124,500</span>
                      <span className="text-green-500 text-xs">↑ 12%</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Instant AI response
              </div>
            </GlassCard>

            {/* Preview Card 2: Financial Insights */}
            <GlassCard intensity="light" shadow="lg" hover="cyanGlow" className="p-6 bg-white/20 dark:bg-white/5 backdrop-blur-xl border border-emerald-400/30 dark:border-emerald-400/20">
              <div className="flex items-start gap-4 mb-4">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-400 flex items-center justify-center flex-shrink-0">
                  <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Financial Insights</h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400">Tally ERP sync</p>
                </div>
              </div>

              {/* Mini Chart Preview */}
              <div className="space-y-3 mb-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-600 dark:text-slate-400">Outstanding</span>
                  <span className="text-sm font-semibold text-slate-900 dark:text-white">$45,200</span>
                </div>
                <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
                  <div className="bg-gradient-to-r from-emerald-500 to-teal-400 h-2 rounded-full" style={{width: '72%'}}></div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-600 dark:text-slate-400">Collected</span>
                  <span className="text-sm font-semibold text-slate-900 dark:text-white">$118,300</span>
                </div>
                <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
                  <div className="bg-gradient-to-r from-emerald-500 to-teal-400 h-2 rounded-full" style={{width: '89%'}}></div>
                </div>
              </div>

              <div className="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Auto-synced with Tally
              </div>
            </GlassCard>

            {/* Preview Card 3: Healthcare Metrics */}
            <GlassCard intensity="light" shadow="lg" hover="glow" className="p-6 bg-white/20 dark:bg-white/5 backdrop-blur-xl border border-purple-400/30 dark:border-purple-400/20">
              <div className="flex items-start gap-4 mb-4">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-violet-500 to-purple-400 flex items-center justify-center flex-shrink-0">
                  <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                  </svg>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-slate-900 dark:text-white">Patient Metrics</h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400">HIMS integration</p>
                </div>
              </div>

              {/* Vital Signs Style Metrics */}
              <div className="grid grid-cols-2 gap-3 mb-4">
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                  <div className="text-2xl font-bold text-slate-900 dark:text-white">247</div>
                  <div className="text-xs text-slate-500 dark:text-slate-400">Today</div>
                </div>
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                  <div className="text-2xl font-bold text-green-500">↑ 8%</div>
                  <div className="text-xs text-slate-500 dark:text-slate-400">vs Yesterday</div>
                </div>
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                  <div className="text-2xl font-bold text-slate-900 dark:text-white">1,842</div>
                  <div className="text-xs text-slate-500 dark:text-slate-400">This Month</div>
                </div>
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                  <div className="text-2xl font-bold text-blue-500">42</div>
                  <div className="text-xs text-slate-500 dark:text-slate-400">Depts</div>
                </div>
              </div>

              <div className="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Real-time from HIMS
              </div>
            </GlassCard>
          </StaggerChildren>

          {/* Capabilities Grid - What Makes MediSync Different */}
          <section className="mb-16">
            <div className="text-center mb-10">
              <h2 className="text-3xl font-bold text-slate-900 dark:text-white mb-3">
                Everything You Need
              </h2>
              <p className="text-slate-600 dark:text-slate-400 max-w-2xl mx-auto">
                From conversational queries to automated accounting, MediSync connects your healthcare data in ways you never thought possible.
              </p>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-4">
              <CapabilityCard
                icon="💬"
                title="Conversational BI"
                description="Ask questions in plain English or Arabic. Get instant charts."
                color="blue"
              />
              <CapabilityCard
                icon="📄"
                title="AI Accountant"
                description="OCR extracts ledger data. Auto-maps to Tally GL."
                color="emerald"
              />
              <CapabilityCard
                icon="📊"
                title="Smart Reports"
                description="Pre-built MIS reports. Custom dashboards."
                color="violet"
              />
              <CapabilityCard
                icon="🔬"
                title="Deep Analytics"
                description="Autonomous AI analyst. Prescriptive insights."
                color="amber"
              />
            </div>
          </section>

          {/* System Status - Now with Glass Effect */}
          <FadeIn>
            <section>
              <GlassCard intensity="medium" shadow="lg" className="p-6 bg-white/20 dark:bg-white/5 backdrop-blur-xl border border-white/30 dark:border-white/10">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-lg font-semibold text-white drop-shadow-md flex items-center gap-2">
                    <span className="relative flex h-3 w-3">
                      <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                      <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                    </span>
                    System Status
                  </h3>
                  <span className="text-sm text-green-500 font-medium">All Systems Operational</span>
                </div>

                <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
                  <StatusItem
                    name="Frontend"
                    version="React 19.2.4"
                    status="ready"
                  />
                  <StatusItem
                    name="Build Tool"
                    version="Vite 7.3.1"
                    status="ready"
                  />
                  <StatusItem
                    name="AI SDK"
                    version="CopilotKit 1.3.6"
                    status="ready"
                  />
                  <StatusItem
                    name="Language"
                    version={currentLocale === 'en' ? 'English (LTR)' : 'Arabic (RTL)'}
                    status="ready"
                  />
                </div>
              </GlassCard>
            </section>
          </FadeIn>
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
  )
}

/**
 * Capability Card Component - Glassmorphic variant
 */
function CapabilityCard({
  icon,
  title,
  description,
  color,
}: {
  icon: string
  title: string
  description: string
  color: 'blue' | 'emerald' | 'violet' | 'amber'
}) {
  const colorClasses = {
    blue: 'from-blue-500 to-cyan-400',
    emerald: 'from-emerald-500 to-teal-400',
    violet: 'from-violet-500 to-purple-400',
    amber: 'from-amber-500 to-orange-400',
  }

  return (
    <div className="group p-5 rounded-xl bg-white/20 dark:bg-white/5 backdrop-blur-xl border border-white/30 dark:border-white/10 hover:bg-white/30 dark:hover:bg-white/10 hover:-translate-y-1 hover:shadow-glass-lg transition-all duration-300 cursor-pointer">
      <div className={`text-3xl mb-3 group-hover:scale-110 transition-transform duration-300 drop-shadow-lg`}>{icon}</div>
      <h3 className="text-base font-semibold text-white drop-shadow-md mb-2">
        {title}
      </h3>
      <p className="text-sm text-white/80 leading-relaxed">
        {description}
      </p>
      <div className={`mt-4 h-1 rounded-full bg-gradient-to-r ${colorClasses[color]} opacity-0 group-hover:opacity-100 transition-opacity duration-300 shadow-lg`}></div>
    </div>
  )
}

/**
 * Status Item Component - Enhanced with glass effect
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
  const statusConfig = {
    ready: { color: 'bg-green-500', label: 'Ready', pulse: true },
    loading: { color: 'bg-amber-500', label: 'Loading', pulse: true },
    error: { color: 'bg-rose-500', label: 'Error', pulse: false },
  }

  const config = statusConfig[status]

  return (
    <div className="flex items-center justify-between p-3 rounded-lg bg-white/20 dark:bg-white/5 backdrop-blur-sm border border-white/20 dark:border-white/10">
      <div className="flex items-center gap-3">
        <div className="relative">
          <div className={`w-2 h-2 rounded-full ${config.color} shadow-lg`}></div>
          {config.pulse && (
            <span className={`absolute inset-0 rounded-full ${config.color} animate-ping opacity-75`}></span>
          )}
        </div>
        <div>
          <div className="text-sm font-medium text-white drop-shadow-sm">{name}</div>
          <div className="text-xs text-white/70">{config.label}</div>
        </div>
      </div>
      <span className="text-xs text-white/90 font-mono bg-white/10 dark:bg-white/5 px-2 py-1 rounded border border-white/20">
        {version}
      </span>
    </div>
  )
}

/**
 * App Root with Suspense boundary for i18n
 * Includes ThemeProvider and global AnimatedBackground
 */
export default function App() {
  return (
    <ThemeProvider>
      <AnimatedBackground />
      <Suspense fallback={<LoadingFallback />}>
        <AppContent />
      </Suspense>
    </ThemeProvider>
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
