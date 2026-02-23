/**
 * Message List Component
 *
 * Displays chat messages with pronounced glassmorphism,
 * entrance animations, and enhanced accessibility.
 *
 * Features:
 * - WCAG 3.0 APCA compliant text
 * - Pronounced glass effect (iOS 26 style)
 * - Staggered entrance animations
 * - Hover lift effects on user messages
 * - Chart visualization with glass cards
 * - Confidence score color coding
 * - Time formatting with locale support
 *
 * @module components/chat/MessageList
 * @version 2.1.0
 */

import React from 'react'
import { useTranslation } from 'react-i18next'
import type { ChatMessage } from '../../services/api'
import { ChartRenderer } from './ChartRenderer'
import { formatTime as formatTimeLocale } from '../../utils/localeUtils'

interface MessageListProps {
  messages: ChatMessage[]
  locale: string
}

export const MessageList: React.FC<MessageListProps> = ({ messages, locale }) => {
  const { t } = useTranslation('chat')

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp)
    return formatTimeLocale(date, locale)
  }

  const formatConfidence = (confidence?: number) => {
    if (confidence === undefined) return null
    const percentage = Math.round(confidence)
    const colorClass =
      percentage >= 90
        ? 'text-success'
        : percentage >= 70
          ? 'text-warning'
          : 'text-error'
    return (
      <span className={`text-xs ${colorClass} font-medium`}>
        {t('message.confidence', { value: percentage })}
      </span>
    )
  }

  return (
    <div className="space-y-6" role="list" aria-label={t('messageList.ariaLabel', 'Chat messages')}>
      {messages.map((message, index) => (
        <div
          key={message.id}
          role="listitem"
          className={`flex animate-liquid-fade-in ${message.role === 'user' ? 'justify-end' : 'justify-start'
            }`}
          style={{ animationDelay: `${index * 50}ms` }}
        >
          <div
            className={`max-w-[80%] p-4 transition-all duration-300 ${message.role === 'user'
              ? 'liquid-glass-pronounced liquid-glass-brand liquid-glass-hover-lift liquid-glass-hover-glow text-white rounded-2xl rounded-br-md shadow-lg'
              : 'liquid-glass-pronounced liquid-text-primary rounded-2xl rounded-bl-md'
              }`}
          >
            {/* Message Content */}
            <div className="prose prose-sm dark:prose-invert max-w-none">
              {message.role === 'assistant' && (
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-6 h-6 rounded-lg flex items-center justify-center shadow-md" style={{
                    background: 'var(--gradient-logo-spectrum)',
                  }}>
                    <svg
                      className="w-4 h-4 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                      />
                    </svg>
                  </div>
                  <span className="text-xs font-semibold text-secondary">
                    AnySync BI
                  </span>
                </div>
              )}

              <p className="text-sm leading-relaxed whitespace-pre-wrap">{message.content}</p>

              {/* Chart Visualization - Pronounced Glass Card */}
              {message.role === 'assistant' && message.chartSpec && (
                <div className="mt-4 liquid-glass-pronounced liquid-glass-hover-lift rounded-xl p-4 shadow-md">
                  <ChartRenderer
                    chartType={message.chartSpec.type}
                    data={message.chartSpec.chart}
                    locale={locale}
                  />
                </div>
              )}

              {/* Confidence Score */}
              {message.role === 'assistant' && (
                <div className="mt-3 flex items-center justify-between">
                  {formatConfidence(message.confidence)}
                  <span className="text-secondary">
                    {formatTime(message.createdAt)}
                  </span>
                </div>
              )}

              {message.role === 'user' && (
                <span className="text-xs text-white/70 mt-3 block text-end">
                  {formatTime(message.createdAt)}
                </span>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default MessageList
