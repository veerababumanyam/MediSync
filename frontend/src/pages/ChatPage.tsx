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
import React, { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { CopilotKit } from '@copilotkit/react-core'
import { ChatInterface } from '../components/chat/ChatInterface'
import { useChat } from '../hooks/useChat'
import { useLocale } from '../hooks/useLocale'
import { AppLogo, LoadingSpinner } from '../components/common'
import { GlassCard, ThemeToggle, AnimatedBackground } from '../components/ui'
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
  const { t } = useTranslation('chat')
  const { isRTL } = useLocale()
  const {
    isLoading,
    error,
    clearMessages,
    sessionId,
    abort,
  } = useChat()

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

      {/* Header - Glassmorphic */}
      <header className="sticky top-0 z-50 border-b border-glass bg-surface-glass-strong backdrop-blur-[60px] shadow-glass-ios dark:shadow-glass-ios-dark">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <AppLogo size="sm" className="shadow-lg" />
            <div>
              <h1 className="text-xl font-bold text-primary">
                {t('header.title', 'MediSync Chat')}
              </h1>
              <p className="text-sm text-secondary">
                {t('header.subtitle', 'AI-Powered Business Intelligence')}
              </p>
            </div>
          </div>

          <nav className="flex items-center gap-4 flex-wrap">
            <a
              href="/"
              className="inline-flex items-center rounded-md px-2 py-1 text-sm text-secondary hover:bg-surface-glass hover:text-logo-blue dark:hover:text-logo-teal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 transition-colors font-medium"
            >
              {t('nav.home', 'Home')}
            </a>
            <a
              href="/dashboard"
              className="inline-flex items-center rounded-md px-2 py-1 text-sm text-secondary hover:bg-surface-glass hover:text-logo-blue dark:hover:text-logo-teal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 transition-colors font-medium"
            >
              {t('nav.dashboard', 'Dashboard')}
            </a>
            <ThemeToggle useSwitchStyle size="sm" className="rounded-full border border-glass p-1 bg-surface-glass" />
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-6">
        {/* Chat Interface Container - Glassmorphic */}
        <GlassCard intensity="medium" shadow="lg" className="h-[calc(100vh-180px)] overflow-hidden">
          <ChatInterface
            initialSessionId={sessionId}
            className="h-full"
          />
        </GlassCard>

        {/* Error Display - Glass Effect */}
        {error && (
          <GlassCard intensity="light" shadow="sm" className="mt-4 p-4 border-s-4 border-s-red-500">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <svg className="w-5 h-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="text-red-600 dark:text-red-400">{error}</span>
              </div>
              <button
                onClick={clearMessages}
                className="text-sm font-medium text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 underline underline-offset-2"
              >
                {t('error.dismiss', 'Dismiss')}
              </button>
            </div>
          </GlassCard>
        )}
      </main>

      {/* Footer - Glassmorphic */}
      <footer className="fixed bottom-0 left-0 right-0 border-t border-glass bg-surface-glass backdrop-blur-md">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between text-sm text-secondary">
          <span>
            {t('footer.sessionStatus', 'Session active')}
          </span>
          <div className="flex items-center gap-4">
            {isLoading && (
              <div className="flex items-center gap-2">
                <LoadingSpinner size="sm" />
                <span>{t('footer.processing', 'Processing...')}</span>
              </div>
            )}
            {isLoading && (
              <button
                onClick={abort}
                className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 font-medium hover:underline underline-offset-2"
              >
                {t('footer.cancel', 'Cancel')}
              </button>
            )}
            <button
              onClick={clearMessages}
              className="hover:text-slate-700 dark:hover:text-slate-300 font-medium transition-colors"
            >
              {t('footer.newChat', 'New Chat')}
            </button>
          </div>
        </div>
      </footer>
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
