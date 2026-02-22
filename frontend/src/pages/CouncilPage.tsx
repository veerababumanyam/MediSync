/**
 * CouncilPage Component
 *
 * Main interface for the Council of AIs consensus system.
 * Provides multi-agent deliberation with confidence scoring.
 *
 * Features:
 * - Query submission with configurable threshold
 * - Real-time deliberation status
 * - Response visualization with agent breakdown
 * - Evidence trail display
 * - RTL support
 *
 * @module pages/CouncilPage
 */

import React, { useCallback, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useCouncil, useCouncilHealth } from '../hooks/useCouncil'
import { QueryInput, ResponseDisplay, ConfidenceIndicator } from '../components/council'
import { useLocale } from '../hooks/useLocale'
import { AppLogo, LoadingSpinner } from '../components/common'
import { ThemeToggle, AnimatedBackground } from '../components/ui'
import { cn } from '../lib/cn'

/**
 * Props for the CouncilPage component
 */
interface CouncilPageProps {
  /** Optional CSS class name */
  className?: string
}

/**
 * CouncilPage - Main Council of AIs interface
 *
 * Provides a user interface for:
 * - Submitting queries to the Council
 * - Configuring consensus thresholds
 * - Viewing deliberation results
 * - Inspecting agent responses
 * - Reviewing evidence trails
 */
export const CouncilPage: React.FC<CouncilPageProps> = ({ className = '' }) => {
  const { t } = useTranslation('council')
  const { isRTL } = useLocale()

  const {
    deliberation,
    isLoading,
    error,
    createDeliberation,
    reset,
  } = useCouncil()

  const {
    health,
    isLoading: healthLoading,
    fetchHealth,
  } = useCouncilHealth()

  // Update document title
  useEffect(() => {
    document.title = `${t('title', 'Council of AIs')} | MediSync`
    return () => {
      document.title = 'MediSync'
    }
  }, [t])

  // Fetch health status on mount
  useEffect(() => {
    fetchHealth()
  }, [fetchHealth])

  /**
   * Handle query submission
   */
  const handleSubmit = useCallback(
    async (query: string, threshold: number) => {
      await createDeliberation(query, threshold)
    },
    [createDeliberation]
  )

  /**
   * Get health status indicator color
   */
  const getHealthColor = () => {
    if (!health) return 'bg-gray-400'
    switch (health.status) {
      case 'healthy':
        return 'bg-green-500'
      case 'degraded':
        return 'bg-amber-500'
      case 'failed':
        return 'bg-red-500'
      default:
        return 'bg-gray-400'
    }
  }

  return (
    <div
      className={cn(
        'min-h-screen relative z-0',
        isRTL ? 'rtl' : 'ltr',
        className
      )}
    >
      <AnimatedBackground />

      {/* Header - Glassmorphic */}
      <header className="sticky top-0 z-50 border-b border-glass bg-surface-glass-strong backdrop-blur-[60px] shadow-glass-ios dark:shadow-glass-ios-dark">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <AppLogo size="sm" className="shadow-lg" />
            <div>
              <h1 className="text-xl font-bold text-primary">
                {t('title', 'Council of AIs')}
              </h1>
              <p className="text-sm text-secondary">
                {t('subtitle', 'Multi-Agent Consensus System')}
              </p>
            </div>
          </div>

          <nav className="flex items-center gap-4 flex-wrap">
            {/* Health Status Indicator */}
            {!healthLoading && health && (
              <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-surface-glass text-sm">
                <span
                  className={cn('w-2 h-2 rounded-full', getHealthColor())}
                  aria-label={t(`health.${health.status}`, health.status)}
                />
                <span className="text-secondary">
                  {t('health.agentsOnline', 'Agents Online')}:{' '}
                  <span className="font-medium text-primary">
                    {health.healthy_agents}/{health.total_agents}
                  </span>
                </span>
              </div>
            )}

            <a
              href="/"
              className="inline-flex items-center rounded-md px-2 py-1 text-sm text-secondary hover:bg-surface-glass hover:text-logo-blue dark:hover:text-logo-teal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 transition-colors font-medium"
            >
              {t('navigation.home', 'Home')}
            </a>
            <a
              href="/chat"
              className="inline-flex items-center rounded-md px-2 py-1 text-sm text-secondary hover:bg-surface-glass hover:text-logo-blue dark:hover:text-logo-teal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 transition-colors font-medium"
            >
              {t('navigation.chat', 'Chat')}
            </a>
            <ThemeToggle
              useSwitchStyle
              size="sm"
              className="rounded-full border border-glass p-1 bg-surface-glass"
            />
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-6 space-y-6">
        {/* Two-Column Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Query Input Section */}
          <div className="space-y-4">
            <h2 className="text-lg font-semibold text-primary">
              {t('query.title', 'Ask a Question')}
            </h2>
            <QueryInput
              onSubmit={handleSubmit}
              isLoading={isLoading}
              error={error}
              placeholder={t('queryInput.placeholder', 'Ask a medical question...')}
            />

            {/* Quick Tips */}
            <div className="p-4 rounded-xl bg-surface-glass/50 backdrop-blur-sm border border-glass">
              <h3 className="text-sm font-medium text-primary mb-2">
                {t('tips.title', 'Tips for Best Results')}
              </h3>
              <ul className="text-sm text-secondary space-y-1">
                <li>• {t('tips.tip1', 'Be specific with medical terms')}</li>
                <li>• {t('tips.tip2', 'Ask one question at a time')}</li>
                <li>• {t('tips.tip3', 'Higher thresholds require more agent agreement')}</li>
              </ul>
            </div>
          </div>

          {/* Response Display Section */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-primary">
                {t('response.title', 'Consensus Response')}
              </h2>
              {deliberation && (
                <button
                  onClick={reset}
                  className="text-sm text-secondary hover:text-primary transition-colors"
                >
                  {t('queryInput.clear', 'Clear')}
                </button>
              )}
            </div>

            {/* Deliberation Result or Empty State */}
            {deliberation ? (
              <ResponseDisplay
                deliberation={deliberation}
                showAgentResponses
                showEvidenceSummary
              />
            ) : (
              <div className="p-12 rounded-xl bg-surface-glass/30 backdrop-blur-sm border border-glass text-center">
                <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-surface-glass-strong flex items-center justify-center">
                  <svg
                    className="w-8 h-8 text-secondary"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <p className="text-secondary">
                  {t('response.noResponse', 'Submit a question to receive a Council response')}
                </p>
              </div>
            )}
          </div>
        </div>

        {/* Recent Deliberations Section (Placeholder for future) */}
        <div className="p-6 rounded-xl bg-surface-glass/30 backdrop-blur-sm border border-glass">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-primary">
              {t('audit.title', 'Audit Trail')}
            </h2>
            <button
              className="text-sm text-logo-blue dark:text-logo-teal hover:underline"
              disabled
            >
              {t('audit.viewAll', 'View All Deliberations')}
            </button>
          </div>
          <p className="text-sm text-secondary">
            {t('audit.subtitle', 'Review deliberation records and compliance')}
          </p>
          <div className="mt-4 p-4 rounded-lg bg-surface-glass/50 border border-glass text-center text-sm text-secondary">
            {t('audit.comingSoon', 'Deliberation history will appear here')}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-glass bg-surface-glass backdrop-blur-md mt-8">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between text-sm text-secondary">
          <span>
            {t('footer.copyright', '© 2026 MediSync. AI-Powered Healthcare Intelligence.')}
          </span>
          <div className="flex items-center gap-4">
            {health && (
              <span className="flex items-center gap-1.5">
                <span className={cn('w-2 h-2 rounded-full', getHealthColor())} />
                {t(`health.${health.status}`, health.status)}
              </span>
            )}
          </div>
        </div>
      </footer>
    </div>
  )
}

export default CouncilPage
