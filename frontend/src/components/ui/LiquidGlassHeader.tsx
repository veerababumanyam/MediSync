/**
 * LiquidGlassHeader Component
 *
 * Shared header used across all pages for design consistency.
 * Implements iOS 26-style liquid glassmorphism design system.
 *
 * Features:
 * - Logo with app name & tagline
 * - Navigation buttons (Chat, Dashboard)
 * - Dark/Light mode toggle with liquid glass container
 * - Language toggle (EN/AR)
 * - Responsive design with mobile menu support
 *
 * @module components/ui/LiquidGlassHeader
 */
import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AppLogo } from '../common'
import { ThemeToggle, LiquidGlassButton, ButtonPrimary } from './'

export type Route = 'home' | 'chat' | 'council' | 'dashboard'

export interface LiquidGlassHeaderProps {
  /** Current active route for highlighting */
  currentRoute?: Route
  /** Handler for navigation */
  onNavigate: (route: Route) => void
  /** Current locale (en/ar) */
  currentLocale: string
  /** Handler for language toggle */
  onToggleLanguage: () => void
  /** Optional CSS class name */
  className?: string
}

/**
 * LiquidGlassHeader - Consistent navigation header with iOS 26 liquid glass design
 *
 * Used across all pages (Home, Chat, Dashboard, Council) for visual consistency.
 */
export const LiquidGlassHeader: React.FC<LiquidGlassHeaderProps> = ({
  currentRoute = 'home',
  onNavigate,
  currentLocale,
  onToggleLanguage,
  className = '',
}) => {
  const { t } = useTranslation()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  return (
    <>
      <header className={`border-b border-glass bg-surface-glass-strong backdrop-blur-xl sticky top-0 z-50 shadow-glass-md ${className}`}>
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          {/* Logo Section */}
          <button
            type="button"
            className="flex items-center gap-3 cursor-pointer group text-start"
            onClick={() => onNavigate('home')}
            onKeyDown={(e) => e.key === 'Enter' && onNavigate('home')}
            aria-label={`${t('app.name', 'MediSync')} - ${t('navigation.home', 'Home')}`}
          >
            <AppLogo size="sm" className="shadow-lg group-hover:shadow-xl group-hover:scale-105 transition-all duration-300" />
            <div>
              <h1 className="text-xl font-bold text-primary">
                {t('app.name', 'MediSync')}
              </h1>
              <p className="text-sm text-secondary hidden sm:block">
                {t('app.tagline', 'Turn Any Legacy Healthcare System into Conversational AI')}
              </p>
            </div>
          </button>

          {/* Navigation Links */}
          <nav className="flex items-center gap-2 sm:gap-3">
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

            {/* Chat Button */}
            <ButtonPrimary
              onClick={() => onNavigate('chat')}
              icon={(
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              )}
              radius="full"
              hover="scale"
              className={`min-h-[44px] shadow-sm hover:shadow-md ${currentRoute === 'chat' ? 'ring-2 ring-blue-400/40' : ''}`}
              aria-current={currentRoute === 'chat' ? 'page' : undefined}
            >
              {t('navigation.chat', 'Chat')}
            </ButtonPrimary>

            {/* Dashboard Button */}
            <LiquidGlassButton
              variant="prominent"
              onClick={() => onNavigate('dashboard')}
              icon={(
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              )}
              radius="full"
              className={`min-h-[44px] ${currentRoute === 'dashboard' ? 'ring-2 ring-blue-400/40' : ''}`}
              aria-current={currentRoute === 'dashboard' ? 'page' : undefined}
            >
              {t('navigation.dashboard', 'Dashboard')}
            </LiquidGlassButton>

            {/* Theme Toggle with iOS 26-style glass container */}
            <div className="liquid-glass-button-prominent min-h-[44px] min-w-[44px] p-2 flex items-center justify-center">
              <ThemeToggle
                size="md"
                className="w-full h-full"
              />
            </div>

            {/* Language Toggle */}
            <LiquidGlassButton
              variant="prominent"
              onClick={onToggleLanguage}
              radius="full"
              className="min-h-[44px]"
              title={t('app.toggleLanguage', 'Toggle language between English and Arabic')}
              aria-label={t('app.toggleLanguage', 'Toggle language')}
            >
              {currentLocale === 'en' ? 'عربي' : 'EN'}
            </LiquidGlassButton>
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
    </>
  )
}

export default LiquidGlassHeader
