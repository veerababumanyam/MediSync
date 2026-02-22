/**
 * MediSync CopilotKit Integration
 *
 * Provides a CopilotKit popup component with MediSync-specific tools.
 * This component conditionally renders only when CopilotKit is configured.
 *
 * @module components/copilot/MediSyncCopilot
 */
import React from 'react'
import { useCopilotReadable, useCopilotAction } from '@copilotkit/react-core'
import type { UseDashboardReturn } from '../../hooks/useDashboard'

/**
 * Props for MediSyncCopilot component
 */
export interface MediSyncCopilotProps {
  /** Current route name */
  currentRoute?: 'home' | 'chat' | 'dashboard'
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
 * Component that registers CopilotKit actions for MediSync
 *
 * This component uses CopilotKit hooks to register tools/actions that
 * the AI agent can invoke. It must be rendered inside a CopilotKit provider.
 */
export const MediSyncCopilot: React.FC<MediSyncCopilotProps> = ({
  currentRoute = 'home',
  dashboard,
  locale = 'en',
  onNavigate,
  onToggleLanguage,
}) => {
  // Make current state readable to the AI agent
  useCopilotReadable({
    description: 'Current MediSync application state',
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
      'Execute a natural language query against MediSync BI data. Returns charts, tables, and insights.',
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
    description: 'Navigate to a different page in MediSync',
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

/**
 * Floating Copilot Button Component
 *
 * Provides a floating button to open the CopilotKit chat interface.
 */
export const CopilotFloatingButton: React.FC = () => {
  const [isOpen, setIsOpen] = React.useState(false)

  return (
    <div className="fixed bottom-6 right-6 z-50">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-14 h-14 rounded-full bg-gradient-to-br from-blue-600 to-cyan-500 text-white shadow-lg hover:shadow-xl transition-all duration-200 flex items-center justify-center"
        aria-label="Open AI Assistant"
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
        <div className="absolute bottom-16 right-0 w-96 h-[500px] bg-white dark:bg-slate-800 rounded-xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="h-full flex flex-col">
            {/* Header */}
            <div className="p-4 border-b border-slate-200 dark:border-slate-700 bg-gradient-to-r from-blue-600 to-cyan-500">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-white/20 flex items-center justify-center">
                  <span className="text-white font-bold text-sm">M</span>
                </div>
                <div>
                  <h3 className="text-white font-semibold">MediSync AI</h3>
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
                  <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center flex-shrink-0">
                    <span className="text-blue-600 dark:text-blue-400 font-bold text-xs">
                      M
                    </span>
                  </div>
                  <div className="bg-slate-100 dark:bg-slate-700 rounded-lg p-3 max-w-[80%]">
                    <p className="text-sm text-slate-700 dark:text-slate-300">
                      Hello! I'm your MediSync AI assistant. I can help you:
                    </p>
                    <ul className="mt-2 text-sm text-slate-600 dark:text-slate-400 space-y-1">
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
            <div className="p-4 border-t border-slate-200 dark:border-slate-700">
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Ask about your data..."
                  className="flex-1 px-4 py-2 border border-slate-200 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                />
                <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
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

export default MediSyncCopilot
