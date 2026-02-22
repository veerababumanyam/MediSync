import React from 'react';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../common/LanguageSwitcher';

interface ChatHeaderProps {
  onNewSession: () => void;
}

export const ChatHeader: React.FC<ChatHeaderProps> = ({
  onNewSession,
}) => {
  const { t } = useTranslation('chat');

  return (
    <header className="flex items-center justify-between px-4 py-3">
      <div className="flex items-center gap-2">
        <h2 className="text-sm font-semibold text-primary">
          {t('header.panelTitle')}
        </h2>
        <span className="text-xs text-secondary">
          {t('header.sessionStatus', 'Session active')}
        </span>
      </div>

      <div className="flex items-center gap-3">
        {/* Language Switcher */}
        <LanguageSwitcher className="shadow-sm" />

        {/* New Session Button */}
        <button
          onClick={onNewSession}
          className="inline-flex items-center gap-2 px-3 py-2 text-sm font-medium text-primary bg-surface-glass border border-glass rounded-lg shadow-sm hover:bg-surface-glass-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 transition-colors"
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
