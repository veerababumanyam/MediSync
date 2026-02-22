import React from 'react';
import { useTranslation } from 'react-i18next';

interface ChatHeaderProps {
  sessionId: string;
  onNewSession: () => void;
  isDark?: boolean;
}

export const ChatHeader: React.FC<ChatHeaderProps> = ({
  sessionId,
  onNewSession,
  isDark = true,
}) => {
  const { t } = useTranslation('chat');

  return (
    <header className="glass-subtle glass-border flex items-center justify-between px-3 py-2 sm:px-4 sm:py-3 border-b">
      <div className="flex items-center gap-2 sm:gap-3 min-w-0">
        {/* Logo/Brand */}
        <div className="flex items-center gap-2 shrink-0">
          <div className="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 flex items-center justify-center shadow-lg shadow-blue-500/20">
            <svg
              className="w-4 h-4 sm:w-5 sm:h-5 text-white"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
            </svg>
          </div>
          <h1 className={`text-base sm:text-lg font-semibold truncate ${isDark ? 'text-white' : 'text-slate-900'
            }`}>
            {t('header.title')}
          </h1>
        </div>

        {/* Session ID (hidden on mobile) */}
        <span className={`hidden sm:inline text-xs font-mono ${isDark ? 'text-slate-500' : 'text-slate-400'
          }`}>
          {sessionId.slice(0, 8)}...
        </span>
      </div>

      <div className="flex items-center gap-2 sm:gap-3 shrink-0">
        {/* New Session Button — icon-only on mobile */}
        <button
          onClick={onNewSession}
          className="glass-interactive inline-flex items-center justify-center gap-1.5 sm:gap-2 p-2 sm:px-3 sm:py-2 text-sm font-medium rounded-xl transition-all duration-300 min-w-9 min-h-9 sm:min-w-0 sm:min-h-0"
          aria-label={t('header.newSession')}
        >
          <svg
            className="w-4 h-4 shrink-0"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 4v16m8-8H4"
            />
          </svg>
          <span className="hidden sm:inline">{t('header.newSession')}</span>
        </button>
      </div>
    </header>
  );
};

export default ChatHeader;
