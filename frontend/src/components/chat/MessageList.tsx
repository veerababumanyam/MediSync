import React from 'react';
import { useTranslation } from 'react-i18next';
import type { ChatMessage } from '../../services/api';
import { ChartRenderer } from './ChartRenderer';

interface MessageListProps {
  messages: ChatMessage[];
  locale: string;
  isDark?: boolean;
}

export const MessageList: React.FC<MessageListProps> = ({ messages, locale, isDark = true }) => {
  const { t } = useTranslation('chat');

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString(locale === 'ar' ? 'ar-SA' : 'en-US', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatConfidence = (confidence?: number) => {
    if (confidence === undefined) return null;
    const percentage = Math.round(confidence);
    const colorClass =
      percentage >= 90
        ? 'text-emerald-400'
        : percentage >= 70
          ? 'text-amber-400'
          : 'text-rose-400';
    return (
      <span className={`text-xs ${colorClass}`}>
        {t('message.confidence', { value: percentage })}
      </span>
    );
  };

  return (
    <div className="space-y-6">
      {messages.map((message) => (
        <div
          key={message.id}
          className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'
            }`}
        >
          <div
            className={`max-w-[80%] ${message.role === 'user'
                ? 'bg-gradient-to-r from-blue-600 to-cyan-500 text-white rounded-2xl rounded-br-md shadow-lg shadow-blue-500/20'
                : isDark
                  ? 'bg-white/10 backdrop-blur-md border border-white/15 text-white rounded-2xl rounded-bl-md'
                  : 'bg-white border border-slate-200 shadow-sm text-slate-900 rounded-2xl rounded-bl-md'
              } px-4 py-3`}
          >
            {/* Message Content */}
            <div className="prose prose-sm max-w-none">
              {message.role === 'assistant' && (
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-6 h-6 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-full flex items-center justify-center shadow-sm">
                    <svg
                      className="w-3.5 h-3.5 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                      />
                    </svg>
                  </div>
                  <span className={`text-xs font-medium ${isDark ? 'text-slate-400' : 'text-slate-500'
                    }`}>
                    MediSync BI
                  </span>
                </div>
              )}

              <p className="whitespace-pre-wrap">{message.content}</p>

              {/* Chart Visualization */}
              {message.role === 'assistant' && message.chartSpec && (
                <div className={`mt-4 rounded-xl p-4 border ${isDark
                    ? 'bg-white/5 border-white/10'
                    : 'bg-slate-50 border-slate-200'
                  }`}>
                  <ChartRenderer
                    chartType={message.chartSpec.type}
                    data={message.chartSpec.chart}
                    locale={locale}
                  />
                </div>
              )}

              {/* Confidence Score */}
              {message.role === 'assistant' && (
                <div className="mt-2 flex items-center justify-between">
                  {formatConfidence(message.confidence)}
                  <span className={`text-xs ${isDark ? 'text-slate-500' : 'text-slate-400'
                    }`}>
                    {formatTime(message.createdAt)}
                  </span>
                </div>
              )}

              {message.role === 'user' && (
                <span className="text-xs text-white/60 mt-2 block text-right">
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
