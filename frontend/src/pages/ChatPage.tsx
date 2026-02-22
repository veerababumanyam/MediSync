/**
 * ChatPage Component
 *
 * Main chat interface page for MediSync's conversational BI feature.
 * Provides natural language interaction with healthcare and accounting data.
 *
 * Features:
 * - CopilotKit generative UI integration
 * - WebMCP tools for AI agent discoverability
 * - RTL support
 * - Session persistence
 *
 * @module pages/ChatPage
 */
import React, { useCallback, useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { CopilotKit } from '@copilotkit/react-core'
import { ChatInterface } from '../components/chat/ChatInterface'
import { useChat } from '../hooks/useChat'
import { useLocale } from '../hooks/useLocale'
import { AnimatedBackground, LiquidGlassHeader } from '../components/ui'
import type { Route } from '../components/ui/LiquidGlassHeader'
import { webMCPService } from '../services/WebMCPService'

/**
 * Props for the ChatPage component
 */
interface ChatPageProps {
  /** Optional CSS class name */
  className?: string
}

/**
 * ChatPage - Main conversational BI interface
 *
 * Features:
 * - Natural language queries
 * - Streaming responses with SSE
 * - Chart visualizations
 * - RTL support
 * - Session persistence
 */
export const ChatPage: React.FC<ChatPageProps> = ({ className = '' }) => {
  const { t, i18n } = useTranslation('chat')
  const { isRTL } = useLocale()
  const currentLocale = i18n.language
  const {
    error,
    clearMessages,
    sessionId,
  } = useChat()

  // Navigation handler
  const handleNavigate = useCallback((route: Route) => {
    const path = route === 'home' ? '/' : `/${route}`
    window.location.href = path
  }, [])

  // Language toggle handler
  const handleToggleLanguage = useCallback(async () => {
    const newLocale = currentLocale === 'en' ? 'ar' : 'en'
    await i18n.changeLanguage(newLocale)
  }, [currentLocale, i18n])

  // Update document title
  useEffect(() => {
    document.title = `${t('pageTitle', 'Chat')} | MediSync`
    return () => {
      document.title = 'MediSync'
    }
  }, [t])

  // Register WebMCP tools for Chat page
  useEffect(() => {
    webMCPService.registerChatTools({
      onQuery: (query: string) => {
        // The query is handled by the ChatInterface component
        console.log('WebMCP query:', query)
      },
      onSyncTally: async () => {
        console.log('Tally sync requested from chat')
        // This would trigger the Tally sync workflow
      },
      onShowDashboard: () => {
        window.location.href = '/dashboard'
      },
    })

    return () => {
      webMCPService.cleanup()
    }
  }, [])

  // CopilotKit configuration - only enable when API key is available
  const isCopilotKitEnabled = !!import.meta.env.VITE_COPILOT_PUBLIC_API_KEY

  const copilotConfig = useMemo(() => {
    if (!isCopilotKitEnabled) return null
    return {
      publicApiKey: import.meta.env.VITE_COPILOT_PUBLIC_API_KEY,
    }
  }, [isCopilotKitEnabled])

  const pageContent = (
    <div
      className={`min-h-screen relative z-0 ${isRTL ? 'rtl' : 'ltr'} ${className}`}
    >
      <AnimatedBackground />

      {/* Shared Header - Liquid Glass */}
      <LiquidGlassHeader
        currentRoute="chat"
        onNavigate={handleNavigate}
        currentLocale={currentLocale}
        onToggleLanguage={handleToggleLanguage}
      />

      {/* Main Content */}
      <main id="main-content" tabIndex={-1} className="container mx-auto px-4 pt-6 pb-6">
        {/* Chat Interface Container - iOS 26 Liquid Glass */}
        <div className="liquid-glass-content-card h-[calc(100vh-160px)] overflow-hidden">
          <ChatInterface
            initialSessionId={sessionId}
            className="h-full"
          />
        </div>

        {/* Error Display - Liquid Glass with Red Accent */}
        {error && (
          <div
            role="alert"
            aria-live="assertive"
            className="liquid-glass-content-card mt-4 p-4 border-s-4 border-s-red-500"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="text-red-600 dark:text-red-400">{error}</span>
              </div>
              <button
                onClick={clearMessages}
                aria-label={t('error.dismissLabel', 'Dismiss error message')}
                className="liquid-glass-button-prominent text-sm text-red-600! dark:text-red-400! hover:bg-red-50! dark:hover:bg-red-900/20!"
              >
                {t('error.dismiss', 'Dismiss')}
              </button>
            </div>
          </div>
        )}
      </main>
    </div>
  )

  // Wrap with CopilotKit only if configured
  if (isCopilotKitEnabled && copilotConfig) {
    return (
      <CopilotKit {...copilotConfig}>
        {pageContent}
      </CopilotKit>
    )
  }

  return pageContent
}

export default ChatPage
