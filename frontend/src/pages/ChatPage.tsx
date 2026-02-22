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
    error,
    clearMessages,
    sessionId,
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
      {/* Main Content — full viewport, edge-to-edge on mobile */}
      <main className="flex-1 min-h-0 flex flex-col w-full max-w-4xl mx-auto px-0 sm:px-4 sm:py-3 lg:px-6 lg:py-4 relative z-10">
        <div className="glass glass-shine flex flex-col flex-1 min-h-0 overflow-hidden rounded-none sm:rounded-2xl">
          <ChatInterface
            initialSessionId={sessionId}
            className="flex-1 min-h-0 flex flex-col"
            isDark={isDark}
          />
        </div>

        {/* Error Display */}
        {error && (
          <div className={`glass-subtle mt-2 sm:mt-4 mx-2 sm:mx-0 p-3 sm:p-4 rounded-xl flex items-center justify-between ${isDark
            ? 'border-red-500/30 text-red-400'
            : 'border-red-300 text-red-600'
            }`}>
            <span className="text-sm">{error}</span>
            <button
              onClick={clearMessages}
              className="text-sm underline hover:no-underline ml-3 shrink-0"
            >
              {t('error.dismiss', 'Dismiss')}
            </button>
          </div>
        )}
      </main>
    </CopilotKit>
  )
}

export default ChatPage
