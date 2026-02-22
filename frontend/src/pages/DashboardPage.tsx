/**
 * DashboardPage Component
 *
 * Main dashboard page for viewing pinned charts and quick insights.
 * Uses the shared AppHeader and glassmorphism design from the landing page.
 *
 * @module pages/DashboardPage
 */
import React, { useCallback, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { DashboardGrid } from '../components/dashboard/DashboardGrid'
import { useDashboard } from '../hooks/useDashboard'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import type { PinnedChart } from '../services/api'

type Route = 'home' | 'chat' | 'dashboard'

/**
 * Props for the DashboardPage component
 */
interface DashboardPageProps {
  /** Whether dark mode is active */
  isDark: boolean
  /** Navigation function */
  navigateTo: (route: Route) => void
  /** Optional CSS class name */
  className?: string
}

/**
 * DashboardPage - Main dashboard for pinned charts
 */
export const DashboardPage: React.FC<DashboardPageProps> = ({ isDark, navigateTo }) => {
  const { t } = useTranslation('dashboard')
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
    const encodedQuery = encodeURIComponent(chart.naturalLanguageQuery)
    window.location.href = `/chat?query=${encodedQuery}`
  }, [])

  return (
    <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 relative z-10">
      {/* Error Display */}
      {error && (
        <div className={`mb-6 p-4 rounded-xl flex items-center justify-between ${isDark
          ? 'bg-red-500/10 border border-red-500/20 text-red-400'
          : 'bg-red-50 border border-red-200 text-red-600'
          }`}>
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
          className="pb-8"
          isDark={isDark}
        />
      ) : null}

      {/* Quick Actions Section */}
      <section className={`mt-8 rounded-2xl p-6 sm:p-8 transition-colors duration-300 ${isDark
        ? 'glass glass-shine'
        : 'bg-white border border-slate-200 shadow-sm'
        }`}>
        <h2 className={`text-lg font-semibold mb-6 ${isDark ? 'text-white' : 'text-slate-900'
          }`}>
          {t('quickActions.title', 'Quick Actions')}
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
          <QuickActionButton
            isDark={isDark}
            icon="💰"
            label={t('quickActions.revenue', "Today's Revenue")}
            onClick={() => navigateTo('chat')}
          />
          <QuickActionButton
            isDark={isDark}
            icon="👥"
            label={t('quickActions.patients', 'Patient Count')}
            onClick={() => navigateTo('chat')}
          />
          <QuickActionButton
            isDark={isDark}
            icon="📦"
            label={t('quickActions.inventory', 'Low Stock')}
            onClick={() => navigateTo('chat')}
          />
          <QuickActionButton
            isDark={isDark}
            icon="📊"
            label={t('quickActions.trends', 'Weekly Trends')}
            onClick={() => navigateTo('chat')}
          />
        </div>
      </section>
    </main>
  )
}

/**
 * Quick Action Button Component
 */
interface QuickActionButtonProps {
  isDark: boolean
  icon: string
  label: string
  onClick: () => void
}

const QuickActionButton: React.FC<QuickActionButtonProps> = ({
  isDark,
  icon,
  label,
  onClick,
}) => {
  return (
    <button
      onClick={onClick}
      className={`flex items-center gap-3 p-4 rounded-xl transition-all duration-300 group hover:-translate-y-0.5 text-left ${isDark
        ? 'bg-white/5 border border-white/10 hover:bg-white/10 hover:border-white/20'
        : 'bg-slate-50 border border-slate-200 hover:bg-slate-100 hover:border-blue-200 hover:shadow-sm'
        }`}
    >
      <span className="text-2xl">{icon}</span>
      <span className={`text-sm font-medium transition-colors ${isDark
        ? 'text-slate-300 group-hover:text-white'
        : 'text-slate-700 group-hover:text-blue-600'
        }`}>
        {label}
      </span>
    </button>
  )
}

export default DashboardPage
