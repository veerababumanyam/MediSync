/**
 * AnySync CopilotKit Integration
 *
 * Provides a CopilotKit popup component with AnySync-specific tools.
 * This component conditionally renders only when CopilotKit is configured.
 *
 * @module components/copilot/AnySyncCopilot
 */
import React from 'react'
import { useCopilotReadable, useCopilotAction } from '@copilotkit/react-core'
import type { UseDashboardReturn } from '../../hooks/useDashboard'

/**
 * Props for AnySyncCopilot component
 */
export interface AnySyncCopilotProps {
  /** Current route name */
  currentRoute?: 'home' | 'chat' | 'council' | 'dashboard'
  /** Dashboard hook for actions */
  dashboard?: UseDashboardReturn
  /** Locale for i18n */
  locale?: string
  /** Navigate callback */
  onNavigate?: (route: string) => void
  /** Toggle language callback */
  onToggleLanguage?: () => void
}

/**
 * Component that registers CopilotKit actions for AnySync
 *
 * This component uses CopilotKit hooks to register tools/actions that
 * the AI agent can invoke. It must be rendered inside a CopilotKit provider.
 */
export const AnySyncCopilot: React.FC<AnySyncCopilotProps> = ({
  currentRoute = 'home',
  dashboard,
  locale = 'en',
  onNavigate,
  onToggleLanguage,
}) => {
  // Make current state readable to the AI agent
  useCopilotReadable({
    description: 'Current AnySync application state',
    value: {
      route: currentRoute,
      locale,
      chartsCount: dashboard?.charts?.length || 0,
      isLoading: dashboard?.isLoading || false,
    },
  })

  // Register queryBI action
  useCopilotAction({
    name: 'queryBI',
    description:
      'Execute a natural language query against AnySync BI data. Returns charts, tables, and insights.',
    parameters: [
      {
        name: 'query',
        type: 'string',
        description: 'The natural language query to execute',
        required: true,
      },
    ],
    handler: async ({ query }: { query: string }) => {
      console.log('CopilotKit queryBI:', query)
      // Navigate to chat with the query
      if (onNavigate) {
        onNavigate('chat')
      }
      return { success: true, message: `Query "${query}" submitted` }
    },
  })

  // Register navigate action
  useCopilotAction({
    name: 'navigate',
    description: 'Navigate to a different page in AnySync',
    parameters: [
      {
        name: 'route',
        type: 'string',
        description: 'The route to navigate to (home, chat, dashboard)',
        required: true,
      },
    ],
    handler: async ({ route }: { route: string }) => {
      if (onNavigate && ['home', 'chat', 'dashboard'].includes(route)) {
        onNavigate(route)
        return { success: true, message: `Navigating to ${route}` }
      }
      return { success: false, message: 'Invalid route' }
    },
  })

  // Register refreshDashboard action
  useCopilotAction({
    name: 'refreshDashboard',
    description: 'Refresh all charts on the dashboard',
    parameters: [],
    handler: async () => {
      if (dashboard?.refreshAll) {
        await dashboard.refreshAll()
        return { success: true, message: 'Dashboard refreshed' }
      }
      return { success: false, message: 'Dashboard not available' }
    },
  })

  // Register pinChart action
  useCopilotAction({
    name: 'pinChart',
    description: 'Pin a new chart to the dashboard from a natural language query',
    parameters: [
      {
        name: 'query',
        type: 'string',
        description: 'The natural language query for the chart',
        required: true,
      },
      {
        name: 'title',
        type: 'string',
        description: 'The title for the pinned chart',
        required: false,
      },
    ],
    handler: async ({ query, title }: { query: string; title?: string }) => {
      if (dashboard?.pinChart) {
        await dashboard.pinChart({
          naturalLanguageQuery: query,
          title: title || 'Untitled Chart',
          chartType: 'bar',
        })
        return { success: true, message: `Chart "${title}" pinned` }
      }
      return { success: false, message: 'Dashboard not available' }
    },
  })

  // Register toggleLanguage action
  useCopilotAction({
    name: 'toggleLanguage',
    description: 'Toggle between English and Arabic languages',
    parameters: [],
    handler: async () => {
      if (onToggleLanguage) {
        onToggleLanguage()
        return {
          success: true,
          message: `Language toggled to ${locale === 'en' ? 'Arabic' : 'English'}`,
        }
      }
      return { success: false, message: 'Language toggle not available' }
    },
  })

  // This component doesn't render anything visible
  return null
}

import { useTranslation } from 'react-i18next'

/**
 * Floating Copilot Button Component
 *
 * Provides a floating button to open the CopilotKit chat interface.
 */
export const CopilotFloatingButton: React.FC = () => {
  const [isOpen, setIsOpen] = React.useState(false)
  const { t } = useTranslation('copilot')

  return (
    <div className="fixed bottom-6 right-6 z-50">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-14 h-14 rounded-full bg-gradient-brand text-white shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300 flex items-center justify-center liquid-focus-gold"
        aria-label={isOpen ? t('button.close', 'Close AI Assistant') : t('button.open', 'Open AI Assistant')}
      >
        {isOpen ? (
          <svg
            className="w-6 h-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        ) : (
          <svg
            className="w-6 h-6"
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
        )}
      </button>

      {isOpen && (
        <div className="absolute bottom-16 right-0 w-96 h-[500px] liquid-glass-pronounced rounded-2xl shadow-2xl border border-glass-soft overflow-hidden animate-fade-in-up">
          <div className="h-full flex flex-col">
            {/* Header */}
            <div className="p-4 border-b border-glass-soft bg-gradient-brand">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-white/20 flex items-center justify-center">
                  <span className="text-white font-bold text-sm">M</span>
                </div>
                <div>
                  <h3 className="text-white font-semibold">AnySync AI</h3>
                  <p className="text-white/80 text-xs">
                    Ask anything about your data
                  </p>
                </div>
              </div>
            </div>

            {/* Chat area placeholder */}
            <div className="flex-1 p-4 overflow-y-auto">
              <div className="space-y-4">
                <div className="flex gap-3">
                  <div className="w-8 h-8 rounded-full bg-brand-primary-alpha flex items-center justify-center shrink-0">
                    <span className="text-brand-primary font-bold text-xs uppercase">
                      M
                    </span>
                  </div>
                  <div className="bg-surface-glass-subtle rounded-2xl p-3 max-w-[80%] border border-glass-soft">
                    <p className="text-sm text-primary">
                      Hello! I'm your AnySync AI assistant. I can help you:
                    </p>
                    <ul className="mt-2 text-sm text-secondary space-y-1">
                      <li>• Query your business data</li>
                      <li>• Create and manage dashboard charts</li>
                      <li>• Sync data to Tally ERP</li>
                      <li>• Generate reports</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            {/* Input area */}
            <div className="p-4 border-t border-glass-soft bg-surface-glass-subtle">
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder={t('input.placeholder', 'Ask about your data...')}
                  className="flex-1 px-4 py-2 liquid-glass-input rounded-xl text-primary focus-ring text-sm"
                />
                <button
                  type="button"
                  className="px-4 py-2 bg-brand-primary text-white rounded-xl hover:bg-brand-primary/90 transition-colors liquid-focus-gold"
                  aria-label={t('input.send', 'Send')}
                >
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                    />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default AnySyncCopilot
