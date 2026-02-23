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
import { AnySyncCopilot, CopilotFloatingButton } from './components/copilot'
import { AppLogo } from './components/common'
import { AnimatedBackground, LiquidGlassHeader, LiquidGlassFooter } from './components/ui'
import { ThemeProvider, useTheme } from './components/theme'
import { webMCPService } from './services/WebMCPService'

import { HeroCarousel } from './components/landing/HeroCarousel'
import { SectorsSection } from './components/landing/SectorsSection'
import { FinalCTA } from './components/landing/FinalCTA'
import { FeatureCard } from './components/landing/FeatureCard'
import { AnnouncementBanner } from './components/landing/AnnouncementBanner'

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
 * AnySync Main Application Component
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

  // CopilotKit content with AnySync tools
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
        <AnySyncCopilot
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
  const { resolvedTheme } = useTheme()
  const isDark = resolvedTheme === 'dark'
  const isArabicLocale = currentLocale === 'ar'

  useEffect(() => {
    const header = document.querySelector('header')
    const themeToggle =
      document.querySelector('[data-debug-theme-toggle="1"]') ??
      document.querySelector('button[title*="Switch to"]')
    const logoButton = document.querySelector('button[aria-label*="AnySync"]')
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
        'tool-description': 'The AnySync home page with feature overview and quick access to Chat and Dashboard',
      } as React.HTMLAttributes<HTMLDivElement>)}
    >
      {/* Skip-to-main link (WCAG 2.4.1 Bypass Blocks) */}
      <a href="#main-content" className="skip-to-main">Skip to main content</a>

      {/* Shared Header - Liquid Glass */}
      <LiquidGlassHeader
        currentRoute="home"
        onNavigate={navigateTo}
        currentLocale={currentLocale}
        onToggleLanguage={toggleLanguage}
      />

      {/* Main Content - Modular Landing Page */}
      <main id="main-content" className="w-full main-content-top">
        {/* Announcement Banner */}
        <AnnouncementBanner />

        <HeroCarousel
          isDark={isDark}
          onOpenLeadCapture={() => navigateTo('chat')}
        />

        <SectorsSection isDark={isDark} />

        <section id="features" className="py-24 relative z-10 overflow-hidden bg-surface-glass-subtle" aria-labelledby="features-heading">
          <div className="container mx-auto px-4 relative z-10">
            <div className="text-center max-w-3xl mx-auto mb-14 pt-6">
              <h2 id="features-heading" className="text-4xl md:text-5xl font-bold mb-8 text-primary tracking-tight drop-shadow-sm">
                Everything You Need
              </h2>
              <p className="text-xl text-secondary leading-relaxed">
                From conversational queries to automated accounting, AnySync connects your healthcare data in ways you never thought possible.
              </p>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-x-8 gap-y-10">
              <FeatureCard
                isDark={isDark}
                title={t('features.conversationalBI.title')}
                description={t('features.conversationalBI.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                )}
                gradientLight="from-blue-100 to-cyan-100"
                gradientDark="from-blue-500/20 to-cyan-400/20"
                iconColorLight="text-blue-600"
                iconColorDark="text-cyan-400"
                shadowLight="shadow-md shadow-blue-500/15"
                borderLight="border-2 border-blue-200"
                borderDark="border border-cyan-500/20"
                delay="0ms"
              />
              <FeatureCard
                isDark={isDark}
                title={t('features.tallySync.title')}
                description={t('features.tallySync.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                  </svg>
                )}
                gradientLight="from-emerald-100 to-teal-100"
                gradientDark="from-emerald-500/20 to-teal-400/20"
                iconColorLight="text-emerald-600"
                iconColorDark="text-teal-400"
                shadowLight="shadow-md shadow-emerald-500/15"
                borderLight="border-2 border-emerald-200"
                borderDark="border border-teal-500/20"
                delay="100ms"
              />
              <FeatureCard
                isDark={isDark}
                title={t('features.aiAccountant.title')}
                description={t('features.aiAccountant.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                )}
                gradientLight="from-violet-100 to-purple-100"
                gradientDark="from-violet-500/20 to-purple-400/20"
                iconColorLight="text-violet-600"
                iconColorDark="text-violet-400"
                shadowLight="shadow-md shadow-violet-500/15"
                borderLight="border-2 border-violet-200"
                borderDark="border border-violet-500/20"
                delay="200ms"
              />
              <FeatureCard
                isDark={isDark}
                title={t('features.piiProtection.title')}
                description={t('features.piiProtection.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                )}
                gradientLight="from-amber-100 to-orange-100"
                gradientDark="from-amber-500/20 to-orange-400/20"
                iconColorLight="text-amber-600"
                iconColorDark="text-amber-400"
                shadowLight="shadow-md shadow-amber-500/15"
                borderLight="border-2 border-amber-200"
                borderDark="border border-amber-500/20"
                delay="300ms"
              />
              <FeatureCard
                isDark={isDark}
                title={t('features.prescriptiveAnalytics.title')}
                description={t('features.prescriptiveAnalytics.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                  </svg>
                )}
                gradientLight="from-pink-100 to-rose-100"
                gradientDark="from-pink-500/20 to-rose-400/20"
                iconColorLight="text-pink-600"
                iconColorDark="text-pink-400"
                shadowLight="shadow-md shadow-pink-500/15"
                borderLight="border-2 border-pink-200"
                borderDark="border border-pink-500/20"
                delay="400ms"
              />
              <FeatureCard
                isDark={isDark}
                title={t('features.himsConnectivity.title')}
                description={t('features.himsConnectivity.description')}
                icon={(
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                  </svg>
                )}
                gradientLight="from-indigo-100 to-blue-100"
                gradientDark="from-indigo-500/20 to-blue-400/20"
                iconColorLight="text-indigo-600"
                iconColorDark="text-indigo-400"
                shadowLight="shadow-md shadow-indigo-500/15"
                borderLight="border-2 border-indigo-200"
                borderDark="border border-indigo-500/20"
                delay="500ms"
              />
            </div>
          </div>
        </section>

        <FinalCTA
          onOpenLeadCapture={() => navigateTo('chat')}
        />
      </main>

      <LiquidGlassFooter />
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
    <div className="min-h-screen flex items-center justify-center bg-surface">
      <div className="text-center">
        <div className="mx-auto mb-4">
          <AppLogo size="md" />
        </div>
        <p className="text-secondary">Loading AnySync...</p>
      </div>
    </div>
  )
}
