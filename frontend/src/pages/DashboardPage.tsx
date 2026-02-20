/**
 * DashboardPage Component
 *
 * Main dashboard page for viewing pinned charts and quick insights.
 * Displays user-saved visualizations with auto-refresh capability.
 *
 * @module pages/DashboardPage
 */
import React, { useCallback, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { DashboardGrid } from '../components/dashboard/DashboardGrid'
import { useDashboard } from '../hooks/useDashboard'
import { useLocale } from '../hooks/useLocale'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import type { PinnedChart } from '../services/api'

/**
 * Props for the DashboardPage component
 */
interface DashboardPageProps {
  /** Optional CSS class name */
  className?: string
}

/**
 * DashboardPage - Main dashboard for pinned charts
 *
 * Features:
 * - Pinned chart display
 * - Chart management (add, remove, refresh)
 * - Grid layout
 * - RTL support
 * - Quick actions
 */
export const DashboardPage: React.FC<DashboardPageProps> = ({ className = '' }) => {
  const { t } = useTranslation('dashboard')
  const { isRTL } = useLocale()
  const {
    charts,
    isLoading,
    error,
  } = useDashboard()

  // Update document title
  useEffect(() => {
    document.title = `${t('pageTitle', 'Dashboard')} | MediSync`
    return () => {
      document.title = 'MediSync'
    }
  }, [t])

  /**
   * Handle chart click for drill-down
   */
  const handleChartClick = useCallback((chart: PinnedChart) => {
    // Navigate to chat with the chart's query
    const encodedQuery = encodeURIComponent(chart.naturalLanguageQuery)
    window.location.href = `/chat?query=${encodedQuery}`
  }, [])

  return (
    <div
      className={`min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 ${
        isRTL ? 'rtl' : 'ltr'
      } ${className}`}
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
                {t('header.title', 'Dashboard')}
              </h1>
              <p className="text-sm text-slate-500 dark:text-slate-400">
                {t('header.subtitle', 'Your pinned insights')}
              </p>
            </div>
          </div>

          <nav className="flex items-center gap-4">
            <a
              href="/"
              className="text-sm text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
            >
              {t('nav.home', 'Home')}
            </a>
            <a
              href="/chat"
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
              {t('nav.chat', 'Chat')}
            </a>
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-6">
        {/* Error Display */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg flex items-center justify-between">
            <span>{error}</span>
            <button
              onClick={() => window.location.reload()}
              className="text-sm underline hover:no-underline"
            >
              {t('error.retry', 'Retry')}
            </button>
          </div>
        )}

        {/* Loading State */}
        {isLoading && charts.length === 0 && (
          <div className="flex items-center justify-center h-64">
            <LoadingSpinner size="lg" label={t('loading', 'Loading dashboard...')} />
          </div>
        )}

        {/* Dashboard Grid */}
        {!isLoading || charts.length > 0 ? (
          <DashboardGrid
            onChartClick={handleChartClick}
            className="pb-20"
          />
        ) : null}

        {/* Quick Actions Section */}
        <section className="mt-8 bg-white dark:bg-slate-800 rounded-xl p-6 shadow-sm border border-slate-200 dark:border-slate-700">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4">
            {t('quickActions.title', 'Quick Actions')}
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
            <QuickActionButton
              icon="💰"
              label={t('quickActions.revenue', 'Today\'s Revenue')}
              href="/chat?query=What%20is%20today's%20revenue%3F"
            />
            <QuickActionButton
              icon="👥"
              label={t('quickActions.patients', 'Patient Count')}
              href="/chat?query=How%20many%20patients%20today%3F"
            />
            <QuickActionButton
              icon="📦"
              label={t('quickActions.inventory', 'Low Stock')}
              href="/chat?query=Show%20low%20stock%20items"
            />
            <QuickActionButton
              icon="📊"
              label={t('quickActions.trends', 'Weekly Trends')}
              href="/chat?query=Show%20weekly%20revenue%20trend"
            />
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-200 dark:border-slate-700 mt-12">
        <div className="container mx-auto px-4 py-6 text-center text-sm text-slate-500 dark:text-slate-400">
          <p>
            {t('footer.hint', 'Pin charts from chat conversations to see them here.')}
          </p>
        </div>
      </footer>
    </div>
  )
}

/**
 * Quick Action Button Component
 */
interface QuickActionButtonProps {
  icon: string
  label: string
  href: string
}

const QuickActionButton: React.FC<QuickActionButtonProps> = ({
  icon,
  label,
  href,
}) => {
  return (
    <a
      href={href}
      className="flex items-center gap-3 p-4 bg-slate-50 dark:bg-slate-700/50 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors group"
    >
      <span className="text-2xl">{icon}</span>
      <span className="text-sm font-medium text-slate-700 dark:text-slate-300 group-hover:text-blue-600 dark:group-hover:text-blue-400">
        {label}
      </span>
    </a>
  )
}

export default DashboardPage
