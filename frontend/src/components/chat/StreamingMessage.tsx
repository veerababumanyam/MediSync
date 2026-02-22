import React from 'react';
import { useTranslation } from 'react-i18next';
import type { SSEEvent } from '../../services/api';
import { ChartRenderer } from './ChartRenderer';

interface StreamingMessageProps {
  events: SSEEvent[];
  locale: string;
  onCancel: () => void;
  isDark?: boolean;
}

export const StreamingMessage: React.FC<StreamingMessageProps> = ({
  events,
  locale,
  onCancel,
  isDark = true,
}) => {
  const { t } = useTranslation('chat');

  const renderEvent = (event: SSEEvent, index: number) => {
    switch (event.type) {
      case 'thinking':
        return (
          <div
            key={index}
            className={`flex items-center gap-2 text-sm ${isDark ? 'text-slate-400' : 'text-slate-500'
              }`}
          >
            <div className="flex-shrink-0">
              <svg
                className="w-4 h-4 animate-spin"
                fill="none"
                viewBox="0 0 24 24"
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
            className={`rounded-xl p-3 my-2 ${isDark
                ? 'bg-slate-900/80 border border-white/10'
                : 'bg-slate-900 border border-slate-700'
              }`}
          >
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs text-slate-400 font-mono">SQL</span>
              <button
                onClick={() => navigator.clipboard.writeText(event.sql || '')}
                className="text-slate-400 hover:text-slate-200 transition-colors"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
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
            <pre className="text-sm text-emerald-400 font-mono overflow-x-auto">
              {event.sql}
            </pre>
          </div>
        );

      case 'result':
        return (
          <div key={index} className="my-2">
            {event.data && event.chartType ? (
              <div className={`rounded-xl p-4 border ${isDark
                  ? 'bg-white/5 border-white/10'
                  : 'bg-white border-slate-200'
                }`}>
                <ChartRenderer
                  chartType={event.chartType}
                  data={event.data}
                  locale={locale}
                />
              </div>
            ) : null}
            {event.confidence !== undefined && (
              <div className={`mt-2 text-xs ${isDark ? 'text-slate-500' : 'text-slate-400'
                }`}>
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
            className={`rounded-xl p-3 my-2 ${isDark
                ? 'bg-red-500/10 border border-red-500/20 text-red-400'
                : 'bg-red-50 border border-red-200 text-red-600'
              }`}
          >
            {event.message}
          </div>
        );

      case 'clarification':
        return (
          <div
            key={index}
            className={`rounded-xl p-3 my-2 ${isDark
                ? 'bg-amber-500/10 border border-amber-500/20 text-amber-400'
                : 'bg-amber-50 border border-amber-200 text-amber-700'
              }`}
          >
            <p className="mb-2">{event.message}</p>
            {event.options && event.options.length > 0 && (
              <div className="space-y-2">
                {event.options.map((option, optIndex) => (
                  <button
                    key={optIndex}
                    className={`block w-full text-left px-3 py-2 rounded-lg border transition-all duration-200 ${isDark
                        ? 'bg-white/5 border-amber-500/20 hover:bg-white/10'
                        : 'bg-white border-amber-200 hover:bg-amber-50'
                      }`}
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
    <div className="flex justify-start">
      <div className={`max-w-[80%] rounded-2xl rounded-bl-md px-4 py-3 ${isDark
          ? 'bg-white/10 backdrop-blur-md border border-white/15'
          : 'bg-white border border-slate-200 shadow-sm'
        }`}>
        {/* AI Header */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-full flex items-center justify-center shadow-sm">
              <svg
                className="w-3.5 h-3.5 text-white animate-pulse"
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

          {/* Cancel Button */}
          <button
            onClick={onCancel}
            className="text-xs text-slate-400 hover:text-red-400 transition-colors"
          >
            {t('streaming.cancel')}
          </button>
        </div>

        {/* Streaming Status */}
        {latestThinking && (
          <div className={`mb-2 text-sm ${isDark ? 'text-slate-400' : 'text-slate-500'
            }`}>
            {latestThinking.message}
          </div>
        )}

        {/* Event Stream */}
        <div className="space-y-2">
          {events.map((event, index) => renderEvent(event, index))}
        </div>

        {/* Loading Indicator */}
        <div className="flex items-center gap-1 mt-2">
          <div className="w-2 h-2 bg-blue-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
          <div className="w-2 h-2 bg-cyan-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
          <div className="w-2 h-2 bg-blue-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
        </div>
      </div>
    </div>
  );
};

export default StreamingMessage;
