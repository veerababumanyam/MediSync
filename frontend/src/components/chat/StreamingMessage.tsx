import React from 'react';
import { useTranslation } from 'react-i18next';
import type { SSEEvent } from '../../services/api';
import { ChartRenderer } from './ChartRenderer';

interface StreamingMessageProps {
  events: SSEEvent[];
  locale: string;
  onCancel: () => void;
}

export const StreamingMessage: React.FC<StreamingMessageProps> = ({
  events,
  locale,
  onCancel,
}) => {
  const { t } = useTranslation('chat');

  const renderEvent = (event: SSEEvent, index: number) => {
    switch (event.type) {
      case 'thinking':
        return (
          <div
            key={index}
            className="flex items-center gap-2 liquid-text-secondary text-sm leading-relaxed"
          >
            <div className="shrink-0">
              <svg
                className="w-4 h-4 animate-spin text-logo-teal"
                fill="none"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                />
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                />
              </svg>
            </div>
            <span>{event.message}</span>
          </div>
        );

      case 'sql_preview':
        return (
          <div
            key={index}
            className="liquid-glass-subtle rounded-xl p-4 my-2 font-mono"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs text-slate-400 font-mono">SQL</span>
              <button
                onClick={() => navigator.clipboard.writeText(event.sql || '')}
                className="liquid-glass-button-prominent px-3 py-1.5 min-h-[36px] text-xs text-slate-400! hover:text-slate-200!"
                aria-label={t('streaming.copySql', 'Copy SQL to clipboard')}
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
              </button>
            </div>
            <pre className="text-sm text-logo-teal font-mono overflow-x-auto">
              {event.sql}
            </pre>
          </div>
        );

      case 'result':
        return (
          <div key={index} className="my-2">
            {event.data != null && event.chartType && (
              <div className="liquid-glass-subtle rounded-xl p-4">
                <ChartRenderer
                  chartType={event.chartType}
                  data={event.data}
                  locale={locale}
                />
              </div>
            )}
            {event.confidence !== undefined && (
              <div className="liquid-glass-badge mt-2 text-xs">
                {t('streaming.confidence', {
                  value: Math.round(event.confidence),
                })}
              </div>
            )}
          </div>
        );

      case 'error':
        return (
          <div
            key={index}
            className="liquid-glass-badge bg-red-50! dark:bg-red-900/20! border-red-200! dark:border-red-800! text-red-600 dark:text-red-400 rounded-xl p-4 my-2"
          >
            {event.message}
          </div>
        );

      case 'clarification':
        return (
          <div
            key={index}
            className="liquid-glass-badge bg-amber-50! dark:bg-amber-900/20! border-amber-200! dark:border-amber-800! text-amber-700 dark:text-amber-400 rounded-xl p-4 my-2"
          >
            <p className="mb-2 text-sm leading-relaxed">{event.message}</p>
            {event.options && event.options.length > 0 && (
              <div className="space-y-2">
                {event.options.map((option, optIndex) => (
                  <button
                    key={optIndex}
                    className="liquid-glass-button-prominent block w-full text-start px-4 py-3 min-h-[44px] text-sm font-medium"
                  >
                    {option}
                  </button>
                ))}
              </div>
            )}
          </div>
        );

      default:
        return null;
    }
  };

  // Get the latest thinking message for the status
  const latestThinking = [...events]
    .reverse()
    .find((e) => e.type === 'thinking');

  return (
    <div className="flex justify-start" aria-live="polite" aria-busy="true">
      <div className="liquid-glass-content-card max-w-[80%] rounded-2xl rounded-bl-md p-4">
        {/* AI Header */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full flex items-center justify-center" style={{
              background: 'linear-gradient(135deg, #2750a8 0%, #18929d 100%)'
            }}>
              <svg
                className="w-4 h-4 text-white animate-pulse"
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

          {/* Cancel Button */}
          <button
            onClick={onCancel}
            className="liquid-glass-button-prominent px-3 py-1.5 min-h-[36px] text-xs text-slate-400! hover:text-red-500!"
            aria-label={t('streaming.cancel')}
          >
            {t('streaming.cancel')}
          </button>
        </div>

        {/* Streaming Status */}
        {latestThinking && (
          <div className="mb-2 text-sm leading-relaxed liquid-text-secondary">
            {latestThinking.message}
          </div>
        )}

        {/* Event Stream */}
        <div className="space-y-2">
          {events.map((event, index) => renderEvent(event, index))}
        </div>

        {/* Loading Indicator - Brand Colors */}
        <div className="flex items-center gap-1 mt-3" aria-hidden="true">
          <div className="w-2 h-2 rounded-full animate-bounce" style={{ background: '#2750a8', animationDelay: '0ms' }} />
          <div className="w-2 h-2 rounded-full animate-bounce" style={{ background: '#18929d', animationDelay: '150ms' }} />
          <div className="w-2 h-2 rounded-full animate-bounce" style={{ background: '#2750a8', animationDelay: '300ms' }} />
        </div>
      </div>
    </div>
  );
};

export default StreamingMessage;
