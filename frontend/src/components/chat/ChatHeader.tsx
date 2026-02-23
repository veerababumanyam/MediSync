/**
 * Chat Header Component
 *
 * Header for the chat interface with new session button,
 * language switcher, and pronounced glassmorphism.
 *
 * Features:
 * - Pronounced glass effect (iOS 26 style)
 * - Gradient underline for active state indication
 * - WCAG 3.0 Gold focus indicators
 * - Language switcher integration
 * - New session button with hover effects
 *
 * @module components/chat/ChatHeader
 * @version 2.1.0
 */

import React from 'react'
import { useTranslation } from 'react-i18next'
import { LanguageSwitcher } from '../common/LanguageSwitcher'
import { CTA_CLASSES } from '@/lib/cta-classes'

interface ChatHeaderProps {
  onNewSession: () => void
}

export const ChatHeader: React.FC<ChatHeaderProps> = ({
  onNewSession,
}) => {
  const { t } = useTranslation('chat')

  return (
    <header className="flex items-center justify-between px-4 py-3 liquid-glass-pronounced">
      <div className="flex items-center gap-3">
        <h2 className="text-sm font-semibold hero-gradient-text">
          {t('header.panelTitle')}
        </h2>
        <div className="liquid-glass-badge px-3 py-1 rounded-full">
          <span className="text-xs text-secondary font-medium">
            {t('header.sessionStatus', 'Session active')}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Language Switcher */}
        <LanguageSwitcher className="liquid-glass-button-prominent liquid-focus-gold !px-3 !py-2" />

        {/* New Session Button */}
        <button
          onClick={onNewSession}
          className={CTA_CLASSES.HEADER_PRIMARY}
          aria-label={t('header.newSession')}
        >
          <svg
            className="w-4 h-4 transition-transform duration-300 group-hover:rotate-90"
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
  )
}

export default ChatHeader
