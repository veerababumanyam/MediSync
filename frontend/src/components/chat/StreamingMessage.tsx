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
            className="flex items-center gap-2 text-gray-600 dark:text-gray-400 text-sm"
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
            className="bg-gray-900 dark:bg-gray-950 rounded-lg p-3 my-2"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs text-gray-400 font-mono">SQL</span>
              <button
                onClick={() => navigator.clipboard.writeText(event.sql || '')}
                className="text-gray-400 hover:text-gray-200"
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
            <pre className="text-sm text-green-400 font-mono overflow-x-auto">
              {event.sql}
            </pre>
          </div>
        );

      case 'result':
        return (
          <div key={index} className="my-2">
            {event.data && event.chartType && (
              <div className="bg-white dark:bg-gray-900 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
                <ChartRenderer
                  chartType={event.chartType}
                  data={event.data}
                  locale={locale}
                />
              </div>
            )}
            {event.confidence !== undefined && (
              <div className="mt-2 text-xs text-gray-500">
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
            className="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg p-3 my-2"
          >
            {event.message}
          </div>
        );

      case 'clarification':
        return (
          <div
            key={index}
            className="bg-yellow-50 dark:bg-yellow-900/20 text-yellow-700 dark:text-yellow-400 rounded-lg p-3 my-2"
          >
            <p className="mb-2">{event.message}</p>
            {event.options && event.options.length > 0 && (
              <div className="space-y-2">
                {event.options.map((option, optIndex) => (
                  <button
                    key={optIndex}
                    className="block w-full text-left px-3 py-2 bg-white dark:bg-gray-800 rounded border border-yellow-200 dark:border-yellow-800 hover:bg-yellow-100 dark:hover:bg-gray-700 transition-colors"
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
      <div className="max-w-[80%] bg-gray-100 dark:bg-gray-800 rounded-2xl rounded-bl-md px-4 py-3">
        {/* AI Header */}
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center">
              <svg
                className="w-4 h-4 text-primary-600 dark:text-primary-400 animate-pulse"
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
            <span className="text-xs font-medium text-gray-500 dark:text-gray-400">
              MediSync BI
            </span>
          </div>

          {/* Cancel Button */}
          <button
            onClick={onCancel}
            className="text-xs text-gray-400 hover:text-red-500 transition-colors"
          >
            {t('streaming.cancel')}
          </button>
        </div>

        {/* Streaming Status */}
        {latestThinking && (
          <div className="mb-2 text-sm text-gray-600 dark:text-gray-400">
            {latestThinking.message}
          </div>
        )}

        {/* Event Stream */}
        <div className="space-y-2">
          {events.map((event, index) => renderEvent(event, index))}
        </div>

        {/* Loading Indicator */}
        <div className="flex items-center gap-1 mt-2">
          <div className="w-2 h-2 bg-primary-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
          <div className="w-2 h-2 bg-primary-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
          <div className="w-2 h-2 bg-primary-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
        </div>
      </div>
    </div>
  );
};

export default StreamingMessage;
