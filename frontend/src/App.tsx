import { Suspense, useCallback, useEffect, useMemo } from 'react'
import { CopilotKit } from '@copilotkit/react-core'
import { useTranslation } from 'react-i18next'
import { useTheme } from './hooks/useTheme'
import { BrowserRouter, Routes, Route, useNavigate, useLocation } from 'react-router-dom'
import './i18n'
import './styles/globals.css'

// Components
import { AppHeader } from './components/common/AppHeader'

// Lazy load page components
import { Home } from './pages/Home'
import { ChatPage } from './pages/ChatPage'
import { DashboardPage } from './pages/DashboardPage'

/**
 * MediSync Main Application Component
 */
function AppContent() {
  const { i18n } = useTranslation()
  const { isDark, toggleTheme } = useTheme()
  const currentLocale = i18n.language
  const isRTL = currentLocale === 'ar'
  const navigate = useNavigate()
  const location = useLocation()

  // Map route path to Route type for AppHeader
  const currentRoute = location.pathname === '/chat' ? 'chat' : location.pathname === '/dashboard' ? 'dashboard' : 'home'

  useEffect(() => {
    document.documentElement.dir = isRTL ? 'rtl' : 'ltr'
    document.documentElement.lang = currentLocale
  }, [currentLocale, isRTL])

  const toggleLanguage = useCallback(() => {
    const newLocale = currentLocale === 'en' ? 'ar' : 'en'
    i18n.changeLanguage(newLocale)
  }, [currentLocale, i18n])

  const navigateTo = useCallback((route: 'home' | 'chat' | 'dashboard') => {
    const path = route === 'home' ? '/' : `/${route}`
    navigate(path)
  }, [navigate])

  const copilotConfig = useMemo(() => ({
    runtimeUrl: import.meta.env.VITE_COPILOT_API_URL || '/api/copilotkit',
  }), [])

  return (
    <CopilotKit {...copilotConfig}>
      <div
        className={`min-h-screen transition-colors duration-500 flex flex-col ${isRTL ? 'rtl' : 'ltr'
          }`}
      >
        {/* Background — adapts to theme */}
        <div className={`fixed inset-0 -z-20 transition-colors duration-500 ${isDark
          ? 'bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900'
          : 'bg-gradient-to-br from-blue-50/80 via-white to-teal-50/80'
          }`} />
        {/* Animated Orbs (both modes, different opacity) */}
        <div className="fixed inset-0 -z-10 overflow-hidden pointer-events-none" aria-hidden="true" style={{ '@media (prefers-reduced-motion: reduce)': { display: 'none' } } as any}>
          <div className={`orb orb-primary top-[-10%] left-[-5%] ${isDark ? 'opacity-100' : 'opacity-60'}`} />
          <div className={`orb orb-secondary bottom-[-15%] right-[-10%] ${isDark ? 'opacity-100' : 'opacity-50'}`} />
          <div className={`orb orb-tertiary top-[40%] left-[60%] ${isDark ? 'opacity-100' : 'opacity-40'}`} />
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

        {/* Page Content Routes */}
        <Suspense fallback={<LoadingFallback />}>
          <Routes>
            <Route path="/" element={<Home isDark={isDark} />} />
            <Route path="/chat" element={<ChatPage isDark={isDark} />} />
            <Route path="/dashboard" element={<DashboardPage isDark={isDark} navigateTo={navigateTo} />} />
          </Routes>
        </Suspense>
      </div>
    </CopilotKit>
  )
}

/**
 * App Root with Suspense boundary
 */
export default function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback />}>
        <AppContent />
      </Suspense>
    </BrowserRouter>
  )
}

/**
 * Loading fallback
 */
function LoadingFallback() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50/80 via-white to-teal-50/80 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      <div className="text-center relative z-10">
        <img src="/logo.png" alt="MediSync" className="w-16 h-16 mx-auto mb-4 animate-pulse relative z-10" />
        <p className="text-slate-500 dark:text-slate-400 text-sm font-medium">Loading MediSync...</p>
      </div>
      {/* Animated Orbs (both modes, different opacity) */}
      <div className="fixed inset-0 -z-10 overflow-hidden pointer-events-none" aria-hidden="true">
        <div className="orb orb-primary top-[-10%] left-[-5%] opacity-60 dark:opacity-100" />
        <div className="orb orb-secondary bottom-[-15%] right-[-10%] opacity-50 dark:opacity-100" />
        <div className="orb orb-tertiary top-[40%] left-[60%] opacity-40 dark:opacity-100" />
      </div>
    </div>
  )
}
