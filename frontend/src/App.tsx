import { Suspense, useCallback, useEffect, useMemo, useState } from 'react'
import { CopilotKit } from '@copilotkit/react-core'
import { useTranslation } from 'react-i18next'
import './i18n'
import './styles/globals.css'

// Lazy load page components
import { ChatPage } from './pages/ChatPage'
import { CouncilPage } from './pages/CouncilPage'
import { DashboardPage } from './pages/DashboardPage'

// CopilotKit components
import { MediSyncCopilot, CopilotFloatingButton } from './components/copilot'
import { AppLogo } from './components/common'
import { GlassCard, ThemeToggle, AnimatedBackground } from './components/ui'
import { FadeIn, StaggerChildren } from './components/animations'
import { ThemeProvider } from './components/theme'
import { webMCPService } from './services/WebMCPService'

/**
 * Route type definition
 */
type Route = 'home' | 'chat' | 'council' | 'dashboard'

const sendDebugLog = (payload: {
  runId: string
  hypothesisId: string
  location: string
  message: string
  data: Record<string, unknown>
}) => {
  // #region agent log
  // Disabled to prevent ERR_CONNECTION_REFUSED in environments where the debug server is not running
  /*
  fetch('http://127.0.0.1:7583/ingest/d551f5e1-ee3b-4b67-a81f-cb9e9d8d73e8', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Debug-Session-Id': 'aeb0a5',
    },
    body: JSON.stringify({
      sessionId: 'aeb0a5',
      timestamp: Date.now(),
      ...payload,
    }),
  }).catch(() => {})
  */
  if (import.meta.env.DEV && false) {
    console.debug('[DebugLog]', payload);
  }
  // #endregion
}

/**
 * Get current route from URL
 */
function getCurrentRoute(): Route {
  const path = window.location.pathname
  if (path === '/chat') return 'chat'
  if (path === '/council') return 'council'
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

  // Update document direction and language when locale changes (W3C i18n: use i18n.dir for dir)
  useEffect(() => {
    document.documentElement.dir = i18n.dir(currentLocale)
    document.documentElement.lang = currentLocale
  }, [currentLocale, i18n])

  // Language toggle handler (await changeLanguage for consistent UI and lang/dir updates)
  const toggleLanguage = useCallback(async () => {
    const newLocale = currentLocale === 'en' ? 'ar' : 'en'
    await i18n.changeLanguage(newLocale)
  }, [currentLocale, i18n])

  // Navigation handlers
  const navigateTo = useCallback((route: Route) => {
    const path = route === 'home' ? '/' : `/${route}`
    sendDebugLog({
      runId: 'initial',
      hypothesisId: 'H1',
      location: 'frontend/src/App.tsx:navigateTo',
      message: 'navigateTo invoked',
      data: {
        currentPath: window.location.pathname,
        targetRoute: route,
        targetPath: path,
      },
    })
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
      onShowRecommendations: (category: string) => {
        console.log('WebMCP requested recommendations for category:', category)
      },
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
      case 'council':
        return <CouncilPage />
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
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const isArabicLocale = currentLocale === 'ar'

  const handleLogoNavigation = useCallback(
    (source: 'click' | 'keyboard') => {
      sendDebugLog({
        runId: 'initial',
        hypothesisId: 'H2',
        location: 'frontend/src/App.tsx:handleLogoNavigation',
        message: 'home logo navigation triggered',
        data: {
          source,
          currentPath: window.location.pathname,
          locale: currentLocale,
        },
      })
      navigateTo('home')
    },
    [currentLocale, navigateTo]
  )

  const handleChatNavigation = useCallback(() => {
    sendDebugLog({
      runId: 'initial',
      hypothesisId: 'H2',
      location: 'frontend/src/App.tsx:handleChatNavigation',
      message: 'chat button click triggered',
      data: {
        currentPath: window.location.pathname,
        locale: currentLocale,
      },
    })
    navigateTo('chat')
  }, [currentLocale, navigateTo])

  useEffect(() => {
    const header = document.querySelector('header')
    const themeToggle =
      document.querySelector('[data-debug-theme-toggle="1"]') ??
      document.querySelector('button[title*="Switch to"]')
    const logoButton = document.querySelector('button[aria-label*="MediSync"]')
    const navButtons = Array.from(document.querySelectorAll('nav button'))
    const chatButton = navButtons.find((button) => {
      const text = button.textContent?.trim().toLowerCase() ?? ''
      return text === 'chat' || text === 'دردشة'
    })

    const readSnapshot = (element: Element | null) => {
      if (!(element instanceof HTMLElement)) {
        return { found: false }
      }
      const style = window.getComputedStyle(element)
      const rect = element.getBoundingClientRect()
      const centerX = rect.left + rect.width / 2
      const centerY = rect.top + rect.height / 2
      const topElement = document.elementFromPoint(centerX, centerY)
      return {
        found: true,
        pointerEvents: style.pointerEvents,
        opacity: style.opacity,
        visibility: style.visibility,
        zIndex: style.zIndex,
        rect: {
          left: rect.left,
          top: rect.top,
          width: rect.width,
          height: rect.height,
        },
        topElementAtCenter:
          topElement instanceof HTMLElement
            ? {
              tagName: topElement.tagName,
              className: topElement.className,
            }
            : null,
      }
    }

    sendDebugLog({
      runId: 'initial',
      hypothesisId: 'H3',
      location: 'frontend/src/App.tsx:HomePage/useEffect/domSnapshot',
      message: 'header interactivity snapshot',
      data: {
        path: window.location.pathname,
        locale: currentLocale,
        isArabicLocale,
        headerClassName: header instanceof HTMLElement ? header.className : null,
        logoButton: readSnapshot(logoButton),
        chatButton: readSnapshot(chatButton ?? null),
        themeToggle: readSnapshot(themeToggle),
      },
    })
  }, [currentLocale, isArabicLocale])

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
      onShowDashboard: () => {
        navigateTo('dashboard')
      },
    })

    return () => {
      // Cleanup is handled by AppContent
    }
  }, [navigateTo])

  return (
    <div
      className={`min-h-screen ${isRTL ? 'rtl' : 'ltr'
        }`}
      // WebMCP declarative attributes for home page
      {...({
        'tool-name': 'medi-home',
        'tool-description': 'The MediSync home page with feature overview and quick access to Chat and Dashboard',
      } as React.HTMLAttributes<HTMLDivElement>)}
    >
      {/* Header - Enhanced Glass Effect with Colorful Border */}
      <header className="border-b border-glass bg-surface-glass-strong backdrop-blur-xl sticky top-0 z-50 shadow-glass-md">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <button
            type="button"
            className="flex items-center gap-3 cursor-pointer group text-start"
            onClick={() => handleLogoNavigation('click')}
            onKeyDown={(e) => e.key === 'Enter' && handleLogoNavigation('keyboard')}
            aria-label={`${t('app.name', 'MediSync')} - ${t('navigation.home', 'Home')}`}
          >
            <AppLogo size="sm" className="shadow-lg group-hover:shadow-xl group-hover:scale-105 transition-all duration-300" />
            <div>
              <h1 className="text-xl font-bold text-primary">
                {t('app.name', 'MediSync')}
              </h1>
              <p className="text-sm text-secondary">
                {t('app.tagline', 'AI-Powered Business Intelligence')}
              </p>
            </div>
          </button>

          {/* Navigation Links */}
          <nav className="flex items-center gap-2 sm:gap-4">
            {/* Mobile Menu Toggle */}
            <button
              className="md:hidden p-2 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors mr-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              aria-expanded={isMobileMenuOpen}
              aria-label={t('navigation.toggleMenu', 'Toggle mobile menu')}
              type="button"
            >
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                {isMobileMenuOpen ? (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>

            {/* Section Navigation Links */}
            <div className="hidden md:flex items-center gap-1">
              <a
                href="#features"
                className="px-3 py-2 text-sm font-semibold text-secondary hover:text-blue-700 dark:hover:text-blue-300 transition-colors rounded-lg hover:bg-surface-glass focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/80"
              >
                {t('navigation.features', 'Features')}
              </a>
              <a
                href="#pricing"
                className="px-3 py-2 text-sm font-semibold text-secondary hover:text-blue-700 dark:hover:text-blue-300 transition-colors rounded-lg hover:bg-surface-glass focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/80"
              >
                {t('navigation.pricing', 'Pricing')}
              </a>
              <a
                href="#about"
                className="px-3 py-2 text-sm font-semibold text-secondary hover:text-blue-700 dark:hover:text-blue-300 transition-colors rounded-lg hover:bg-surface-glass focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/80"
              >
                {t('navigation.about', 'About')}
              </a>
            </div>

            {/* Divider */}
            <div className="hidden md:block w-px h-6 bg-slate-200 dark:bg-slate-700 mx-2"></div>

            {/* Action Buttons */}
            <div className="flex items-center gap-2 sm:gap-4">
              <button
                type="button"
                onClick={handleChatNavigation}
                className="inline-flex items-center justify-center gap-2 min-h-[44px] px-4 py-2 bg-action-primary text-on-brand rounded-lg transition-all hover:scale-105 active:scale-95 text-sm font-medium shadow-sm hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
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
                type="button"
                onClick={() => navigateTo('dashboard')}
                className="hidden sm:inline-flex items-center justify-center gap-2 min-h-[44px] px-4 py-2 text-secondary hover:text-blue-700 dark:hover:text-blue-300 hover:bg-surface-glass rounded-lg transition-all text-sm font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/80"
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
            </div>

            {/* Action Toggles Group */}
            <div className="flex items-center gap-2 sm:gap-4 border-l border-slate-200 dark:border-slate-700 pl-2 sm:pl-4">
              {/* Theme Toggle with Icon (min 44px touch target) */}
              <ThemeToggle size="lg" className="min-h-[44px] min-w-[44px]" />

              {/* Language Toggle with Globe Icon */}
              <button
                type="button"
                onClick={toggleLanguage}
                className="inline-flex items-center justify-center gap-2 min-h-[44px] px-3 py-2 rounded-lg bg-surface-glass border border-glass hover:bg-surface-glass-strong transition-colors text-sm font-medium text-secondary"
                title={t('app.toggleLanguage', 'Toggle language between English and Arabic')}
                aria-label={t('app.toggleLanguage', 'Toggle language')}
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
                <span>{currentLocale === 'en' ? 'عربي' : 'EN'}</span>
              </button>
            </div>
          </nav>
        </div>
      </header>

      {/* Mobile Menu Dropdown */}
      {isMobileMenuOpen && (
        <div className="md:hidden bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 px-4 py-4 space-y-3 shadow-md absolute w-full z-40">
          <a href="#features" onClick={() => setIsMobileMenuOpen(false)} className="block px-3 py-2 rounded-md text-base font-medium text-slate-700 dark:text-slate-200 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
            {t('navigation.features', 'Features')}
          </a>
          <a href="#pricing" onClick={() => setIsMobileMenuOpen(false)} className="block px-3 py-2 rounded-md text-base font-medium text-slate-700 dark:text-slate-200 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
            {t('navigation.pricing', 'Pricing')}
          </a>
          <a href="#about" onClick={() => setIsMobileMenuOpen(false)} className="block px-3 py-2 rounded-md text-base font-medium text-slate-700 dark:text-slate-200 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
            {t('navigation.about', 'About')}
          </a>
        </div>
      )}

      {/* Main Content - Redesigned Healthcare BI Command Center */}
      <main className="container mx-auto px-4 py-8">
        {/* Hero Section with Glass Effect */}
        <FadeIn>
          <section className="mb-16 text-center">
            {/* Prominent AI Badge with Better Contrast */}
            <div className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-linear-to-r from-blue-700 to-teal-700 backdrop-blur-sm border border-white/30 text-white text-base font-semibold mb-6 shadow-lg shadow-blue-500/20">
              <span className="relative flex h-2.5 w-2.5">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-300 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-cyan-400"></span>
              </span>
              {t('home.hero.badge')}
            </div>

            <h1 className="text-5xl md:text-6xl font-bold text-slate-900 dark:text-white mb-6 tracking-tight drop-shadow-md">
              {t('home.hero.title')}<span className="text-transparent bg-clip-text bg-linear-to-r from-logo-blue to-logo-teal font-extrabold">{t('home.hero.titleHighlight')}</span>
            </h1>

            <p className="text-xl text-secondary max-w-3xl mx-auto mb-8 leading-relaxed drop-shadow-sm">
              {t('home.hero.subtitle')}
            </p>

            {/* CTA Buttons - Prominent Pill-Shaped Design (solid primary for readability) */}
            <div className="flex flex-col sm:flex-row gap-6 justify-center items-stretch sm:items-center mt-8">
              <button
                type="button"
                onClick={() => navigateTo('chat')}
                className="group inline-flex items-center justify-center gap-3 min-w-[200px] min-h-[56px] px-8 py-4 bg-action-primary text-on-brand drop-shadow-sm rounded-full shadow-lg shadow-blue-500/40 hover:shadow-xl hover:shadow-blue-500/50 hover:scale-105 active:scale-100 transition-all duration-300 text-lg font-semibold ring-1 ring-inset ring-white/20 focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-logo-blue/50"
              >
                <svg className="w-6 h-6 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                {t('home.hero.ctaChat')}
              </button>

              {/* Secondary CTA - View Dashboard */}
              <button
                type="button"
                onClick={() => navigateTo('dashboard')}
                className="group inline-flex items-center justify-center gap-3 min-w-[200px] min-h-[56px] px-8 py-4 liquid-glass-light text-slate-900 dark:text-white rounded-full shadow-md hover:shadow-lg hover:scale-105 active:scale-100 transition-all duration-300 text-lg font-semibold focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-slate-400/50"
              >
                <svg className="w-6 h-6 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                {t('home.hero.ctaDashboard')}
              </button>
            </div>

            {/* Trust Signals - Social Proof */}
            <div className="mt-12 pt-8 border-t border-slate-300/30 dark:border-white/10">
              <p className="text-sm font-medium text-slate-800 dark:text-white mb-6">{t('home.hero.trustBy')}</p>
              <div className="flex flex-wrap items-center justify-center gap-6 md:gap-8 opacity-90">
                {/* Trust Badges */}
                <div className="flex items-center gap-2 text-primary bg-surface-glass py-2.5 px-5 rounded-full backdrop-blur-sm border border-glass shadow-sm">
                  <svg className="w-5 h-5 text-emerald-600 dark:text-emerald-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  <span className="text-sm font-semibold">{t('home.hero.hipaa')}</span>
                </div>
                <div className="flex items-center gap-2 text-primary bg-surface-glass py-2.5 px-5 rounded-full backdrop-blur-sm border border-glass shadow-sm">
                  <svg className="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  <span className="text-sm font-semibold">{t('home.hero.soc2')}</span>
                </div>
                <div className="flex items-center gap-2 text-primary bg-surface-glass py-2.5 px-5 rounded-full backdrop-blur-sm border border-glass shadow-sm">
                  <svg className="w-5 h-5 text-purple-600 dark:text-purple-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
                  </svg>
                  <span className="text-sm font-semibold">{t('home.hero.clinics')}</span>
                </div>
                <div className="flex items-center gap-2 text-primary bg-surface-glass py-2.5 px-5 rounded-full backdrop-blur-sm border border-glass shadow-sm">
                  <svg className="w-5 h-5 text-cyan-600 dark:text-cyan-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M3 3a1 1 0 000 2v8a2 2 0 002 2h2.586l-1.293 1.293a1 1 0 101.414 1.414L10 15.414l2.293 2.293a1 1 0 001.414-1.414L12.414 15H15a2 2 0 002-2V5a1 1 0 100-2H3zm11.707 4.707a1 1 0 00-1.414-1.414L10 9.586 8.707 8.293a1 1 0 00-1.414 0l-2 2a1 1 0 101.414 1.414L8 10.414l1.293 1.293a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  <span className="text-sm font-medium">{t('home.hero.queries')}</span>
                </div>
              </div>
            </div>
          </section>
        </FadeIn>

        {/* Live Preview Section - What the User Can Actually Do */}
        <StaggerChildren className="grid lg:grid-cols-3 gap-6 mb-16">
          <GlassCard intensity="light" shadow="lg" hover="glow" className="p-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="w-12 h-12 rounded-xl bg-linear-to-br from-logo-blue to-logo-teal flex items-center justify-center shrink-0">
                <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-slate-900 dark:text-white">{t('home.preview.askAnything')}</h3>
                <p className="text-sm text-slate-500 dark:text-slate-400">{t('home.preview.askDesc')}</p>
              </div>
            </div>

            {/* Simulated Chat Interface */}
            <div className="space-y-3 mb-4">
              <div className="flex gap-3">
                <div className="w-6 h-6 rounded-full bg-slate-200 dark:bg-slate-700 shrink-0 flex items-center justify-center">
                  <svg className="w-3 h-3 text-slate-500" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" />
                  </svg>
                </div>
                <div className="flex-1 bg-slate-100 dark:bg-slate-800/50 rounded-lg px-3 py-2 text-sm text-slate-700 dark:text-slate-300">
                  "{t('home.preview.revenueQuestion')}"
                </div>
              </div>
              <div className="flex gap-3">
                <div className="w-6 h-6 rounded-full bg-logo-blue shrink-0 flex items-center justify-center">
                  <span className="text-white text-xs font-bold">M</span>
                </div>
                <div className="flex-1 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg px-3 py-2 text-sm text-blue-700 dark:text-blue-300">
                  <div className="flex items-center gap-2">
                    <span>$124,500</span>
                    <span className="text-green-500 text-xs">
                      <span className="sr-only">Increased by </span>
                      ↑ 12%
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="text-xs font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              {t('home.preview.instantResponse')}
            </div>
          </GlassCard>

          {/* Preview Card 2: Financial Insights */}
          <GlassCard intensity="light" shadow="lg" hover="cyanGlow" className="p-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="w-12 h-12 rounded-xl bg-linear-to-br from-logo-teal to-[rgba(24,146,157,0.5)] flex items-center justify-center shrink-0">
                <svg className="w-6 h-6 text-slate-900" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-slate-900 dark:text-white">{t('home.preview.financialInsights')}</h3>
                <p className="text-sm text-slate-500 dark:text-slate-400">{t('home.preview.tallySync')}</p>
              </div>
            </div>

            {/* Mini Chart Preview */}
            <div className="space-y-3 mb-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-600 dark:text-slate-400">{t('home.preview.outstanding')}</span>
                <span className="text-sm font-semibold text-slate-900 dark:text-white">$45,200</span>
              </div>
              <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
                <div className="bg-linear-to-r from-emerald-500 to-teal-400 h-2 rounded-full" style={{ width: '72%' }}></div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-slate-600 dark:text-slate-400">{t('home.preview.collected')}</span>
                <span className="text-sm font-semibold text-slate-900 dark:text-white">$118,300</span>
              </div>
              <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
                <div className="bg-linear-to-r from-emerald-500 to-teal-400 h-2 rounded-full" style={{ width: '89%' }}></div>
              </div>
            </div>

            <div className="text-xs font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
              {t('home.preview.autoSynced')}
            </div>
          </GlassCard>

          {/* Preview Card 3: Healthcare Metrics */}
          <GlassCard intensity="light" shadow="lg" hover="glow" className="p-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="w-12 h-12 rounded-xl bg-linear-to-br from-logo-blue to-logo-teal flex items-center justify-center shrink-0">
                <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-slate-900 dark:text-white">{t('home.preview.patientMetrics')}</h3>
                <p className="text-sm text-slate-500 dark:text-slate-400">{t('home.preview.himsIntegration')}</p>
              </div>
            </div>

            {/* Vital Signs Style Metrics */}
            <div className="grid grid-cols-2 gap-3 mb-4">
              <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                <div className="text-xs font-medium text-slate-500 dark:text-slate-400 mb-0.5">{t('home.preview.today')}</div>
                <div className="text-xl font-bold text-slate-900 dark:text-white">247</div>
              </div>
              <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                <div className="text-xs font-medium text-slate-500 dark:text-slate-400 mb-0.5">{t('home.preview.vsYesterday')}</div>
                <div className="text-sm font-bold text-emerald-600 dark:text-emerald-400 mt-1">
                  <span className="sr-only">Increased by </span>
                  ↑ 8%
                </div>
              </div>
              <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                <div className="text-xs font-medium text-slate-500 dark:text-slate-400 mb-0.5">{t('home.preview.thisMonth')}</div>
                <div className="text-xl font-bold text-slate-900 dark:text-white">1,842</div>
              </div>
              <div className="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3">
                <div className="text-xs font-medium text-slate-500 dark:text-slate-400 mb-0.5">{t('home.preview.depts')}</div>
                <div className="text-sm font-bold text-blue-600 dark:text-blue-400 mt-1">42</div>
              </div>
            </div>

            <div className="text-xs font-medium text-slate-700 dark:text-slate-300 flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {t('home.preview.realTime')}
            </div>
          </GlassCard>
        </StaggerChildren>

        {/* Capabilities Grid - What Makes MediSync Different */}
        <section className="mb-16">
          <div className="text-center mb-10">
            <h2 className="text-3xl font-bold text-slate-900 dark:text-white mb-6">
              {t('home.section.title')}
            </h2>
            <p className="text-slate-600 dark:text-slate-400 max-w-2xl mx-auto leading-relaxed">
              {t('home.section.subtitle')}
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6" id="features">
            <CapabilityCard
              icon="chat"
              title={t('features.conversationalBI.title')}
              description={t('features.conversationalBI.description')}
              color="blue"
            />
            <CapabilityCard
              icon="document"
              title={t('features.aiAccountant.title')}
              description={t('features.aiAccountant.description')}
              color="emerald"
            />
            <CapabilityCard
              icon="chart"
              title={t('features.smartReports.title')}
              description={t('features.smartReports.description')}
              color="violet"
            />
            <CapabilityCard
              icon="analytics"
              title={t('features.deepAnalytics.title')}
              description={t('features.deepAnalytics.description')}
              color="amber"
            />
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-200/50 dark:border-slate-700/50 mt-12 pt-8 bg-white/5 dark:bg-slate-900/50 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-8">
          {/* Footer Links */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-x-8 gap-y-10 mb-12">
            {/* Product */}
            <div className="flex flex-col space-y-4">
              <h4 id="footer-product" className="text-sm font-semibold text-slate-900 dark:text-white uppercase tracking-wider mb-1">{t('home.footer.product')}</h4>
              <ul aria-labelledby="footer-product" className="space-y-3">
                <li><a href="#features" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.features')}</a></li>
                <li><a href="#pricing" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.pricing')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.integrations')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.api')}</a></li>
              </ul>
            </div>

            {/* Company */}
            <div className="flex flex-col space-y-4 border-s border-transparent md:border-slate-200/50 md:dark:border-slate-700/50 md:ps-8">
              <h4 id="footer-company" className="text-sm font-semibold text-slate-900 dark:text-white uppercase tracking-wider mb-1">{t('home.footer.company')}</h4>
              <ul aria-labelledby="footer-company" className="space-y-3">
                <li><a href="#about" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.about')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.blog')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.careers')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.contact')}</a></li>
              </ul>
            </div>

            {/* Resources */}
            <div className="flex flex-col space-y-4 border-s border-transparent md:border-slate-200/50 md:dark:border-slate-700/50 md:ps-8">
              <h4 id="footer-resources" className="text-sm font-semibold text-slate-900 dark:text-white uppercase tracking-wider mb-1">{t('home.footer.resources')}</h4>
              <ul aria-labelledby="footer-resources" className="space-y-3">
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.documentation')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.helpCenter')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.status')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.security')}</a></li>
              </ul>
            </div>

            {/* Legal */}
            <div className="flex flex-col space-y-4 border-s border-transparent md:border-slate-200/50 md:dark:border-slate-700/50 md:ps-8">
              <h4 id="footer-legal" className="text-sm font-semibold text-slate-900 dark:text-white uppercase tracking-wider mb-1">{t('home.footer.legal')}</h4>
              <ul aria-labelledby="footer-legal" className="space-y-3">
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.privacyPolicy')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.termsOfService')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.cookiePolicy')}</a></li>
                <li><a href="#" aria-disabled="true" className="text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">{t('home.footer.compliance')}</a></li>
              </ul>
            </div>
          </div>

          {/* Footer Bottom */}
          <div className="flex flex-col md:flex-row items-center justify-between pt-6 border-t border-slate-200/50 dark:border-slate-700/50">
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-4 md:mb-0">
              {t('home.footer.copyright', { year: new Date().getFullYear() })}
            </p>

            {/* Social Links */}
            <div className="flex items-center gap-4">
              <a href="#" className="text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors" aria-label={t('social.twitter', 'Twitter')}>
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84" />
                </svg>
              </a>
              <a href="#" className="text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors" aria-label={t('social.linkedin', 'LinkedIn')}>
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                </svg>
              </a>
              <a href="#" className="text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors" aria-label={t('social.github', 'GitHub')}>
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
                </svg>
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}

/**
 * Capability Card Component - Glassmorphic variant with SVG icons
 */
function CapabilityCard({
  icon,
  title,
  description,
  color,
}: {
  icon: 'chat' | 'document' | 'chart' | 'analytics'
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

  const iconBgClasses = {
    blue: 'bg-linear-to-br from-blue-500 to-cyan-400',
    emerald: 'bg-linear-to-br from-emerald-500 to-teal-400',
    violet: 'bg-linear-to-br from-violet-500 to-purple-400',
    amber: 'bg-linear-to-br from-amber-500 to-orange-400',
  }

  // SVG Icons for each feature
  const icons = {
    chat: (
      <svg className="w-6 h-6 text-white" role="img" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <title>Chat Interface</title>
        <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
      </svg>
    ),
    document: (
      <svg className="w-6 h-6 text-white" role="img" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <title>AI Accountant</title>
        <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    ),
    chart: (
      <svg className="w-6 h-6 text-white" role="img" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <title>Smart Reports</title>
        <path strokeLinecap="round" strokeLinejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
    ),
    analytics: (
      <svg className="w-6 h-6 text-white" role="img" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <title>Deep Analytics</title>
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
    ),
  }

  const ariaIdPrefix = title.toLowerCase().replace(/\s+/g, '-')

  return (
    <div
      className="group p-8 rounded-3xl liquid-glass liquid-glass-heavy hover:-translate-y-2 transition-all duration-500 cursor-pointer flex flex-col items-center text-center h-full"
      role="region"
      aria-labelledby={`${ariaIdPrefix}-title`}
      aria-describedby={`${ariaIdPrefix}-desc`}
    >
      <div className={`w-14 h-14 rounded-full ${iconBgClasses[color]} flex items-center justify-center mb-6 group-hover:scale-110 group-hover:-translate-y-1 transition-all duration-500 shadow-xl shrink-0`} aria-hidden="true">
        {icons[icon]}
      </div>
      <h3 id={`${ariaIdPrefix}-title`} className="text-xl font-bold text-slate-900 dark:text-white mb-4">
        {title}
      </h3>
      <p id={`${ariaIdPrefix}-desc`} className="text-sm text-slate-700 dark:text-slate-300 leading-relaxed font-medium">
        {description}
      </p>
      <div className={`mt-auto pt-6 w-full h-1.5 rounded-full bg-linear-to-r ${colorClasses[color]} opacity-0 group-hover:opacity-100 transition-opacity duration-500 shadow-lg`} aria-hidden="true"></div>
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
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800">
      <div className="text-center">
        <div className="mx-auto mb-4">
          <AppLogo size="md" />
        </div>
        <p className="text-slate-600 dark:text-slate-400">Loading MediSync...</p>
      </div>
    </div>
  )
}
