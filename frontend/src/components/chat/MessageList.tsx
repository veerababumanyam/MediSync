import React from 'react';
import { useTranslation } from 'react-i18next';
import type { ChatMessage } from '../../services/api';
import { ChartRenderer } from './ChartRenderer';
import { formatTime as formatTimeLocale } from '../../utils/localeUtils';

interface MessageListProps {
  messages: ChatMessage[];
  locale: string;
}

export const MessageList: React.FC<MessageListProps> = ({ messages, locale }) => {
  const { t } = useTranslation('chat');

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return formatTimeLocale(date, locale);
  };

  const formatConfidence = (confidence?: number) => {
    if (confidence === undefined) return null;
    const percentage = Math.round(confidence);
    const colorClass =
      percentage >= 90
        ? 'text-green-600 dark:text-green-400'
        : percentage >= 70
          ? 'text-yellow-600 dark:text-yellow-400'
          : 'text-red-600 dark:text-red-400';
    return (
      <span className={`text-xs ${colorClass}`}>
        {t('message.confidence', { value: percentage })}
      </span>
    );
  };

  return (
    <div className="space-y-6" role="list" aria-label={t('messageList.ariaLabel', 'Chat messages')}>
      {messages.map((message) => (
        <div
          key={message.id}
          role="listitem"
          className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'
            }`}
        >
          <div
            className={`max-w-[80%] p-4 ${message.role === 'user'
                ? 'liquid-glass-cta text-white rounded-2xl rounded-br-md'
                : 'liquid-glass-content-card text-slate-900 dark:text-white rounded-2xl rounded-bl-md'
              }`}
          >
            {/* Message Content */}
            <div className="prose prose-sm dark:prose-invert max-w-none">
              {message.role === 'assistant' && (
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-5 h-5 rounded-full flex items-center justify-center" style={{
                    background: 'linear-gradient(135deg, #2750a8 0%, #18929d 100%)'
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
                  <span className="text-xs font-medium liquid-text-secondary">
                    MediSync BI
                  </span>
                </div>
              )}

              <p className="text-sm leading-relaxed whitespace-pre-wrap">{message.content}</p>

              {/* Chart Visualization - Liquid Glass Card */}
              {message.role === 'assistant' && message.chartSpec && (
                <div className="mt-4 liquid-glass-subtle rounded-xl p-4">
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
                  <span className="text-xs liquid-text-secondary">
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
  );
};

export default MessageList;
