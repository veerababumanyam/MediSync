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
      <div className="flex items-center gap-3">
        <h2 className="text-sm font-semibold hero-gradient-text">
          {t('header.panelTitle')}
        </h2>
        <div className="liquid-glass-badge">
          <span className="text-xs text-slate-600 dark:text-slate-400">
            {t('header.sessionStatus', 'Session active')}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Language Switcher */}
        <LanguageSwitcher className="liquid-glass-button-prominent !px-3 !py-2" />

        {/* New Session Button */}
        <button
          onClick={onNewSession}
          className="liquid-glass-button-primary inline-flex items-center gap-2 px-4 py-2 text-sm font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#18929d] transition-colors"
          aria-label={t('header.newSession')}
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
