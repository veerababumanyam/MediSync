/**
 * ChatPage Component
 *
 * Main chat interface page for MediSync's conversational BI feature.
 * Uses the shared AppHeader and glassmorphism design from the landing page.
 *
 * @module pages/ChatPage
 */
import React, { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { CopilotKit } from '@copilotkit/react-core'
import { ChatInterface } from '../components/chat/ChatInterface'
import { useChat } from '../hooks/useChat'
import { LoadingSpinner } from '../components/common/LoadingSpinner'

/**
 * Props for the ChatPage component
 */
interface ChatPageProps {
  /** Whether dark mode is active */
  isDark: boolean
  /** Optional CSS class name */
  className?: string
}

/**
 * ChatPage - Main conversational BI interface
 */
export const ChatPage: React.FC<ChatPageProps> = ({ isDark }) => {
  const { t } = useTranslation('chat')
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
    runtimeUrl: import.meta.env.VITE_COPILOT_API_URL || '/api/copilotkit',
  }

  return (
    <CopilotKit {...copilotConfig}>
      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 relative z-10">
        <div className={`rounded-2xl h-[calc(100vh-200px)] overflow-hidden ${isDark
          ? 'glass glass-shine'
          : 'bg-white/90 backdrop-blur-sm border border-slate-200 shadow-lg'
          }`}>
          <ChatInterface
            initialSessionId={sessionId}
            className="h-full"
            isDark={isDark}
          />
        </div>

        {/* Error Display */}
        {error && (
          <div className={`mt-4 p-4 rounded-xl flex items-center justify-between ${isDark
            ? 'bg-red-500/10 border border-red-500/20 text-red-400'
            : 'bg-red-50 border border-red-200 text-red-600'
            }`}>
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

      {/* Status Footer */}
      <div className={`fixed bottom-0 left-0 right-0 z-40 border-t transition-colors duration-300 ${isDark
        ? 'border-white/10 bg-white/5 backdrop-blur-xl'
        : 'border-slate-200 bg-white/80 backdrop-blur-xl'
        }`}>
        <div className={`max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex items-center justify-between text-sm ${isDark ? 'text-slate-500' : 'text-slate-400'
          }`}>
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
                className="text-red-400 hover:underline"
              >
                {t('footer.cancel', 'Cancel')}
              </button>
            )}
            <button
              onClick={clearMessages}
              className={`transition-colors ${isDark ? 'hover:text-slate-300' : 'hover:text-slate-600'
                }`}
            >
              {t('footer.newChat', 'New Chat')}
            </button>
          </div>
        </div>
      </div>
    </CopilotKit>
  )
}

export default ChatPage
