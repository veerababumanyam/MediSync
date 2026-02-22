/**
 * DashboardPage Component
 *
 * Main dashboard page for viewing pinned charts and quick insights.
 * Displays user-saved visualizations with auto-refresh capability.
 *
 * Features:
 * - WebMCP integration for AI agent discoverability
 * - CopilotKit Generative UI ready
 * - Pinned chart display with drill-down
 * - RTL support
 * - Liquid Glass design system
 *
 * @module pages/DashboardPage
 */
import React, { useCallback, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { DashboardGrid } from '../components/dashboard/DashboardGrid'
import { useDashboard } from '../hooks/useDashboard'
import { useLocale } from '../hooks/useLocale'
import { LoadingSpinner } from '../components/common/LoadingSpinner'
import { AnimatedBackground } from '../components/ui/AnimatedBackground'
import { LiquidGlassNavbar } from '../components/ui/LiquidGlassNavbar'
import { LiquidGlassCard, GlassBrandCard, GlassTealCard, GlassBlueCard, GlassGreenCard } from '../components/ui/LiquidGlassCard'
import { ButtonPrimary } from '../components/ui/LiquidGlassButton'
import { ThemeToggle } from '../components/ui'
import { FadeIn } from '../components/animations'
import { webMCPService } from '../services/WebMCPService'
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
 * - WebMCP tools for AI agent integration
 */
export const DashboardPage: React.FC<DashboardPageProps> = ({ className = '' }) => {
  const { t } = useTranslation('dashboard')
  const { isRTL } = useLocale()
  const {
    charts,
    isLoading,
    error,
    refreshChart,
    refreshAll,
    pinChart,
  } = useDashboard()

  // Update document title
  useEffect(() => {
    document.title = `${t('pageTitle', 'Dashboard')} | MediSync`
    return () => {
      document.title = 'MediSync'
    }
  }, [t])

  // Register WebMCP tools for Dashboard
  useEffect(() => {
    webMCPService.registerDashboardTools({
      onRefreshChart: async (chartId: string) => {
        try {
          await refreshChart(chartId)
        } catch (err) {
          console.error('Failed to refresh chart:', err)
        }
      },
      onPinChart: async (query: string, title: string) => {
        try {
          await pinChart({
            naturalLanguageQuery: query,
            title,
            chartType: 'bar', // default
          })
        } catch (err) {
          console.error('Failed to pin chart:', err)
        }
      },
      onNavigateToChat: (query: string) => {
        const encodedQuery = encodeURIComponent(query)
        window.location.href = `/chat?query=${encodedQuery}`
      },
      onRefreshAll: async () => {
        try {
          await refreshAll()
        } catch (err) {
          console.error('Failed to refresh dashboard:', err)
        }
      },
      onNavigateToSettings: (section: string) => {
        console.log(`Navigating to settings section: ${section}`)
      },
      onOpenDashboardModal: (id: string) => {
        console.log(`Opening dashboard modal: ${id}`)
      }
    })

    // Register navigation tools
    webMCPService.registerNavigationTools({
      onNavigate: (route: string) => {
        window.location.href = route === 'home' ? '/' : `/${route}`
      },
      onToggleLanguage: () => {
        // This would be connected to the i18n context
        console.log('Toggle language requested')
      },
      onShowRecommendations: (category: string) => {
        console.log(`Showing recommendations for category: ${category}`)
      }
    })

    return () => {
      webMCPService.cleanup()
    }
  }, [refreshChart, refreshAll, pinChart])

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
      className={`min-h-screen ${isRTL ? 'rtl' : 'ltr'} ${className}`}
      // WebMCP declarative attributes for dashboard container
      {...({
        'tool-name': 'medi-dashboard',
        'tool-description': 'The main dashboard for viewing pinned charts and business insights in MediSync',
      } as { 'tool-name': string; 'tool-description': string })}
    >
      {/* Animated Background */}
      <AnimatedBackground />

      {/* Header - Liquid Glass Navbar */}
      <LiquidGlassNavbar
        left={
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl liquid-glass-brand flex items-center justify-center shadow-lg">
              <span className="text-white font-bold text-xl">M</span>
            </div>
            <div>
              <h1 className="text-xl font-semibold liquid-text-primary">
                {t('header.title', 'Dashboard')}
              </h1>
              <p className="text-sm liquid-text-secondary">
                {t('header.subtitle', 'Your pinned insights')}
              </p>
            </div>
          </div>
        }
        right={
          <div className="flex items-center gap-3">
            <a
              href="/"
              className="text-sm liquid-text-secondary hover:text-blue-400 transition-colors font-medium"
            >
              {t('nav.home', 'Home')}
            </a>
            <ButtonPrimary
              onClick={() => window.location.href = '/chat'}
              // WebMCP declarative attribute for chat navigation
              {...({
                'tool-name': 'medi-nav-chat',
                'tool-description': 'Navigate to the conversational BI chat interface',
              } as { 'tool-name': string; 'tool-description': string })}
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
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0 4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
              {t('nav.chat', 'Chat')}
            </ButtonPrimary>
            <ThemeToggle />
          </div>
        }
      />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto p-4 md:p-6 space-y-6">
        {/* KPI Cards Section */}
        <FadeIn>
          <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Revenue Card */}
            <GlassBrandCard interactive className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-linear-to-br from-logo-blue to-logo-teal flex items-center justify-center text-xl">
                  💰
                </div>
                <div>
                  <p className="text-sm liquid-text-secondary">Revenue</p>
                  <p className="text-2xl font-bold liquid-text-primary">$124,500</p>
                </div>
              </div>
              <div className="mt-2 text-sm text-emerald-400"><span aria-hidden="true">↑</span> Up 12% from last month</div>
            </GlassBrandCard>

            {/* Patients Card */}
            <GlassTealCard interactive className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-linear-to-br from-logo-teal to-[rgba(24,146,157,0.5)] flex items-center justify-center text-xl">
                  👥
                </div>
                <div>
                  <p className="text-sm liquid-text-secondary">Patients Today</p>
                  <p className="text-2xl font-bold liquid-text-primary">248</p>
                </div>
              </div>
              <div className="mt-2 text-sm text-emerald-400"><span aria-hidden="true">↑</span> Up 8% from yesterday</div>
            </GlassTealCard>

            {/* Appointments Card */}
            <GlassBlueCard interactive className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-linear-to-br from-logo-blue to-[rgba(39,80,168,0.5)] flex items-center justify-center text-xl">
                  📅
                </div>
                <div>
                  <p className="text-sm liquid-text-secondary">Appointments</p>
                  <p className="text-2xl font-bold liquid-text-primary">56</p>
                </div>
              </div>
              <div className="mt-2 text-sm text-blue-400">12 pending confirmation</div>
            </GlassBlueCard>

            {/* Inventory Card */}
            <GlassGreenCard interactive className="p-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-linear-to-br from-green-500 to-emerald-400 flex items-center justify-center text-xl">
                  📦
                </div>
                <div>
                  <p className="text-sm liquid-text-secondary">Low Stock Items</p>
                  <p className="text-2xl font-bold liquid-text-primary">7</p>
                </div>
              </div>
              <div className="mt-2 text-sm text-amber-400">Requires attention</div>
            </GlassGreenCard>
          </section>
        </FadeIn>

        {/* Error Display - Glass Effect (WCAG 4.1.3: status message, role=alert) */}
        {error && (
          <FadeIn>
            <LiquidGlassCard
              intensity="light"
              className="p-4 border-s-4 border-s-red-500"
              role="alert"
              aria-live="assertive"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <svg className="w-5 h-5 text-red-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span className="text-red-400">{t(error)}</span>
                </div>
                <button
                  onClick={() => window.location.reload()}
                  className="text-sm font-medium text-red-400 hover:text-red-300 underline underline-offset-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2"
                >
                  {t('error.retry', 'Retry')}
                </button>
              </div>
            </LiquidGlassCard>
          </FadeIn>
        )}

        {/* Loading State - Glass Skeleton */}
        {isLoading && charts.length === 0 && (
          <div className="flex items-center justify-center h-64">
            <LoadingSpinner size="lg" label={t('loading', 'Loading dashboard...')} />
          </div>
        )}

        {/* Dashboard Grid - Wrapped with animations */}
        <div
          // WebMCP declarative attribute for charts grid
          {...({
            'tool-name': 'medi-dashboard-grid',
            'tool-description': 'The grid of pinned charts displaying business insights and visualizations',
          } as { 'tool-name': string; 'tool-description': string })}
        >
          {(!isLoading || charts.length > 0) ? (
            <DashboardGrid
              onChartClick={handleChartClick}
              className="pb-4"
            />
          ) : null}
        </div>

        {/* Quick Actions Section - Glass Effect */}
        <FadeIn>
          <section
            className="mt-6"
            // WebMCP declarative attribute for quick actions
            {...({
              'tool-name': 'medi-quick-actions',
              'tool-description': 'Quick action buttons to explore common business queries',
            } as { 'tool-name': string; 'tool-description': string })}
          >
            <LiquidGlassCard intensity="medium" elevation="raised" className="p-6">
              <div className="flex items-center gap-3 mb-6">
                <div className="w-10 h-10 rounded-xl liquid-glass-brand flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                  </svg>
                </div>
                <div>
                  <h2 className="text-lg font-semibold liquid-text-primary">
                    {t('quickActions.title', 'Quick Actions')}
                  </h2>
                  <p className="text-sm liquid-text-secondary">
                    Common queries to get you started
                  </p>
                </div>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <QuickActionButton
                  icon="💰"
                  label={t('quickActions.revenue', "Today's Revenue")}
                  href="/chat?query=What%20is%20today's%20revenue%3F"
                  toolName="quick-revenue"
                  toolDescription="View today's revenue metrics"
                />
                <QuickActionButton
                  icon="👥"
                  label={t('quickActions.patients', 'Patient Count')}
                  href="/chat?query=How%20many%20patients%20today%3F"
                  toolName="quick-patients"
                  toolDescription="View today's patient count"
                />
                <QuickActionButton
                  icon="📦"
                  label={t('quickActions.inventory', 'Low Stock')}
                  href="/chat?query=Show%20low%20stock%20items"
                  toolName="quick-inventory"
                  toolDescription="Check low stock inventory items"
                />
                <QuickActionButton
                  icon="📊"
                  label={t('quickActions.trends', 'Weekly Trends')}
                  href="/chat?query=Show%20weekly%20revenue%20trend"
                  toolName="quick-trends"
                  toolDescription="View weekly revenue trends"
                />
              </div>
            </LiquidGlassCard>
          </section>
        </FadeIn>
      </main>

      {/* Footer */}
      <footer className="border-t border-white/10 mt-8">
        <div className="max-w-7xl mx-auto px-4 py-6 text-center text-sm liquid-text-secondary">
          <p>
            {t('footer.hint', 'Pin charts from chat conversations to see them here.')}
          </p>
        </div>
      </footer>
    </div>
  )
}

/**
 * Quick Action Button Component - Liquid Glass styling
 */
interface QuickActionButtonProps {
  icon: string
  label: string
  href: string
  toolName?: string
  toolDescription?: string
}

const QuickActionButton: React.FC<QuickActionButtonProps> = ({
  icon,
  label,
  href,
  toolName,
  toolDescription,
}) => {
  return (
    <a
      href={href}
      className="liquid-glass flex items-center gap-3 p-4 rounded-xl hover:shadow-lg hover:shadow-blue-500/10 transition-all duration-200 group"
      // WebMCP declarative attributes
      {...(toolName ? {
        'tool-name': `medi-${toolName}`,
        'tool-description': toolDescription || label,
      } as { 'tool-name': string; 'tool-description': string } : {})}
    >
      <span className="text-2xl">{icon}</span>
      <span className="text-sm font-medium liquid-text-primary group-hover:text-blue-400 transition-colors">
        {label}
      </span>
    </a>
  )
}

export default DashboardPage
