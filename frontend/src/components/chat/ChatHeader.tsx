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
    <header className={`flex items-center justify-between px-4 py-3 border-b transition-colors duration-300 ${isDark
        ? 'border-white/10 bg-white/5'
        : 'border-slate-200 bg-slate-50/80'
      }`}>
      <div className="flex items-center gap-3">
        {/* Logo/Brand */}
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 flex items-center justify-center shadow-lg shadow-blue-500/20">
            <svg
              className="w-5 h-5 text-white"
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
          <h1 className={`text-lg font-semibold ${isDark ? 'text-white' : 'text-slate-900'
            }`}>
            {t('header.title')}
          </h1>
        </div>

        {/* Session ID (shortened) */}
        <span className={`text-xs font-mono ${isDark ? 'text-slate-500' : 'text-slate-400'
          }`}>
          {sessionId.slice(0, 8)}...
        </span>
      </div>

      <div className="flex items-center gap-3">
        {/* New Session Button */}
        <button
          onClick={onNewSession}
          className={`inline-flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-xl transition-all duration-300 ${isDark
              ? 'text-slate-300 bg-white/10 border border-white/10 hover:bg-white/15 hover:text-white hover:border-white/20'
              : 'text-slate-600 bg-slate-100 border border-slate-200 hover:bg-slate-200 hover:text-slate-800'
            }`}
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
              d="M12 4v16m8-8H4"
            />
          </svg>
          {t('header.newSession')}
        </button>
      </div>
    </header>
  );
};

export default ChatHeader;
