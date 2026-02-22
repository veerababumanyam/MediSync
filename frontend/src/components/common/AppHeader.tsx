/**
 * AppHeader Component
 *
 * Shared header used across all pages for design consistency.
 * Extracted from the landing page header in App.tsx.
 *
 * Features:
 * - Logo with app name & tagline
 * - Navigation with active-state highlighting
 * - Dark/Light mode toggle
 * - Language toggle (EN/AR)
 * - Responsive design
 *
 * @module components/common/AppHeader
 */
import React from 'react'
import { useTranslation } from 'react-i18next'

type Route = 'home' | 'chat' | 'dashboard'

export interface AppHeaderProps {
    isDark: boolean
    toggleTheme: () => void
    toggleLanguage: () => void
    currentLocale: string
    navigateTo: (route: Route) => void
    currentRoute: Route
}

export const AppHeader: React.FC<AppHeaderProps> = ({
    isDark,
    toggleTheme,
    toggleLanguage,
    currentLocale,
    navigateTo,
    currentRoute,
}) => {
    const { t } = useTranslation()

    return (
        <header
            className={`glass-subtle glass-border sticky top-0 z-100 transition-shadow duration-300 backdrop-blur-md bg-white/70 dark:bg-slate-900/70 border-b flex items-center h-14 sm:h-16 lg:h-[72px] ${!isDark ? 'shadow-sm shadow-slate-200/50' : 'border-white/10'
                }`}
        >
            <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-2 flex items-center justify-between">
                {/* Logo - WCAG 2.5.5 compliant with 44px minimum touch target */}
                <button
                    className="flex items-center gap-2 sm:gap-3 cursor-pointer group shrink-0 min-h-11 px-1 -mx-1 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                    onClick={() => navigateTo('home')}
                    aria-label="MediSync - Home"
                    type="button"
                >
                    <img
                        src="/logo.png"
                        alt=""
                        aria-hidden="true"
                        className="w-8 h-8 sm:w-10 sm:h-10 rounded-xl object-contain group-hover:scale-105 transition-transform duration-300"
                    />
                    <div className="flex flex-col justify-center">
                        <span
                            className={`text-base sm:text-lg font-bold leading-tight ${isDark ? 'text-white' : 'text-slate-900'
                                }`}
                        >
                            {t('app.name', 'MediSync')}
                        </span>
                        <p
                            className={`text-xs hidden md:block pt-0.5 ${isDark ? 'text-slate-400' : 'text-slate-500'
                                }`}
                        >
                            {t('app.tagline', 'AI-Powered Business Intelligence')}
                        </p>
                    </div>
                </button>

                {/* Navigation */}
                <nav className="flex items-center gap-1 sm:gap-2">
                    {/* Chat CTA */}
                    <button
                        onClick={() => navigateTo('chat')}
                        className={
                            currentRoute === 'chat'
                                ? 'inline-flex items-center justify-center gap-1.5 sm:gap-2 px-2.5 sm:px-3.5 py-1.5 sm:py-2 bg-linear-to-r from-blue-600 to-cyan-500 text-white rounded-xl transition-all duration-300 text-xs sm:text-sm font-semibold shadow-lg shadow-blue-500/25 ring-2 ring-blue-400/40 min-w-10 min-h-10'
                                : 'inline-flex items-center justify-center gap-1.5 sm:gap-2 px-2.5 sm:px-3.5 py-1.5 sm:py-2 bg-linear-to-r from-blue-600 to-cyan-500 text-white rounded-xl hover:from-blue-500 hover:to-cyan-400 transition-all duration-300 text-xs sm:text-sm font-semibold shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 hover:-translate-y-0.5 active:scale-95 min-w-10 min-h-10'
                        }
                    >
                        <svg
                            className="w-4 h-4 sm:w-4 sm:h-4 shrink-0"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            strokeWidth={2}
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                            />
                        </svg>
                        <span className="hidden sm:inline whitespace-nowrap">
                            {t('navigation.chat', 'Chat')}
                        </span>
                    </button>

                    {/* Dashboard */}
                    <button
                        onClick={() => navigateTo('dashboard')}
                        className={`glass-interactive inline-flex items-center justify-center gap-1.5 sm:gap-2 px-2.5 sm:px-3.5 py-1.5 sm:py-2 transition-all duration-300 text-xs sm:text-sm font-medium rounded-xl min-w-10 min-h-10 ${currentRoute === 'dashboard'
                            ? 'text-blue-600 dark:text-white ring-1 ring-blue-400/40'
                            : 'text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'
                            }`}
                    >
                        <svg
                            className="w-4 h-4 sm:w-4 sm:h-4 shrink-0"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            strokeWidth={2}
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                            />
                        </svg>
                        <span className="hidden sm:inline whitespace-nowrap">
                            {t('navigation.dashboard', 'Dashboard')}
                        </span>
                    </button>

                    {/* Dark/Light Mode Toggle */}
                    <button
                        onClick={toggleTheme}
                        className="glass-interactive p-1.5 sm:p-2 transition-all duration-300 text-slate-500 hover:text-amber-500 dark:text-amber-400 dark:hover:text-amber-300 rounded-xl flex items-center justify-center min-w-10 min-h-10"
                        aria-label={
                            isDark
                                ? t('app.lightMode', 'Switch to light mode')
                                : t('app.darkMode', 'Switch to dark mode')
                        }
                        title={isDark ? 'Light mode' : 'Dark mode'}
                    >
                        {isDark ? (
                            <svg
                                className="w-4 h-4 sm:w-5 sm:h-5 shrink-0"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                strokeWidth={2}
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                                />
                            </svg>
                        ) : (
                            <svg
                                className="w-4 h-4 sm:w-5 sm:h-5 shrink-0"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                strokeWidth={2}
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"
                                />
                            </svg>
                        )}
                    </button>

                    {/* Language Toggle */}
                    <button
                        onClick={toggleLanguage}
                        className="glass-interactive px-2 sm:px-3 py-1.5 sm:py-2 transition-all duration-300 text-xs sm:text-sm font-medium text-slate-600 hover:text-slate-800 dark:text-slate-300 dark:hover:text-white rounded-xl flex items-center justify-center min-w-10 min-h-10"
                        aria-label={t('app.toggleLanguage', 'Toggle language')}
                    >
                        {currentLocale === 'en' ? 'عربي' : 'EN'}
                    </button>
                </nav>
            </div>
        </header>
    )
}

export default AppHeader
