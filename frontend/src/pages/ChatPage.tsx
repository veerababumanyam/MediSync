/**
 * ChatPage Component
 *
 * Main chat interface page for MediSync's conversational BI feature.
 * Provides natural language interaction with healthcare and accounting data.
 *
 * @module pages/ChatPage
 */
import React, { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { CopilotKit } from '@copilotkit/react-core'
import { ChatInterface } from '../components/chat/ChatInterface'
import { useChat } from '../hooks/useChat'
import { useLocale } from '../hooks/useLocale'
import { LoadingSpinner } from '../components/common/LoadingSpinner'

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

  // CopilotKit configuration
  const copilotConfig = {
    endpoint: import.meta.env.VITE_COPILOT_API_URL || '/api/copilotkit',
  }

  return (
    <CopilotKit {...copilotConfig}>
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
                  {t('header.title', 'MediSync Chat')}
                </h1>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  {t('header.subtitle', 'AI-Powered Business Intelligence')}
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
                href="/dashboard"
                className="text-sm text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
              >
                {t('nav.dashboard', 'Dashboard')}
              </a>
            </nav>
          </div>
        </header>

        {/* Main Content */}
        <main className="container mx-auto px-4 py-6">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 h-[calc(100vh-180px)]">
            <ChatInterface
              initialSessionId={sessionId}
              className="h-full"
            />
          </div>

          {/* Error Display */}
          {error && (
            <div className="mt-4 p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg flex items-center justify-between">
              <span>{error}</span>
              <button
                onClick={clearMessages}
                className="text-sm underline hover:no-underline"
              >
                {t('error.dismiss', 'Dismiss')}
              </button>
            </div>
          )}
        </main>

        {/* Footer */}
        <footer className="fixed bottom-0 left-0 right-0 border-t border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm">
          <div className="container mx-auto px-4 py-3 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
            <span>
              {t('footer.session', 'Session')}: {sessionId.slice(0, 8)}...
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
                  className="text-red-600 dark:text-red-400 hover:underline"
                >
                  {t('footer.cancel', 'Cancel')}
                </button>
              )}
              <button
                onClick={clearMessages}
                className="hover:text-slate-700 dark:hover:text-slate-300"
              >
                {t('footer.newChat', 'New Chat')}
              </button>
            </div>
          </div>
        </footer>
      </div>
    </CopilotKit>
  )
}

export default ChatPage
