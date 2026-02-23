/**
 * Chat Input Component
 *
 * Text input area for chat with WCAG 3.0 Gold focus indicators,
 * pronounced glassmorphism, and enhanced accessibility.
 *
 * Features:
 * - Auto-resizing textarea
 * - WCAG 3.0 Gold focus indicators
 * - Pronounced glass effect (iOS 26 style)
 * - RTL support
 * - Keyboard navigation (Enter to send, Shift+Enter for newline)
 * - Loading state with spinner animation
 *
 * @module components/chat/ChatInput
 * @version 2.1.0
 */

import React, { useState, useCallback, useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { getChatSendButtonClasses } from '@/lib/cta-classes'

interface ChatInputProps {
  onSend: (message: string) => void
  disabled?: boolean
  locale: string
  placeholder?: string
}

export const ChatInput: React.FC<ChatInputProps> = ({
  onSend,
  disabled = false,
  locale,
  placeholder,
}) => {
  const { t } = useTranslation('chat')
  const [input, setInput] = useState('')
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 200)}px`
    }
  }, [input])

  const handleSend = useCallback(() => {
    if (input.trim() && !disabled) {
      onSend(input.trim())
      setInput('')
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto'
      }
    }
  }, [input, disabled, onSend])

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        handleSend()
      }
    },
    [handleSend]
  )

  const isRTL = locale === 'ar'
  const sendLabel = t('input.send', 'Send')
  const placeholderText = placeholder || t('input.placeholder', 'Type your question...')

  return (
    <div
      dir={isRTL ? 'rtl' : 'ltr'}
      className={`relative flex items-end gap-3 p-4 rounded-2xl liquid-glass-pronounced liquid-glass-hover-lift transition-all duration-300 ${isRTL ? 'flex-row-reverse' : ''}`}
    >
      <textarea
        ref={textareaRef}
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholderText}
        disabled={disabled}
        rows={1}
        aria-label={placeholderText}
        className={`flex-1 min-h-[44px] resize-none bg-transparent text-primary placeholder:text-muted focus:outline-none liquid-focus-gold text-sm disabled:opacity-50 text-start rounded-xl transition-all duration-200 ${isRTL ? 'rounded-e-xl' : 'rounded-s-xl'}`}
        dir={isRTL ? 'rtl' : 'ltr'}
      />

      <button
        type="button"
        onClick={handleSend}
        disabled={disabled || !input.trim()}
        className={`${getChatSendButtonClasses(!disabled && input.trim().length > 0)} ${isRTL ? 'rounded-s-xl' : 'rounded-e-xl'}`}
        aria-label={sendLabel}
      >
        {disabled ? (
          <svg
            className="w-5 h-5 animate-spin"
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
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        ) : (
          <svg
            className={`w-5 h-5 transition-transform duration-300 ${!disabled && input.trim() ? 'group-hover:translate-x-0.5' : ''} ${isRTL ? 'rotate-180' : ''}`}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
            />
          </svg>
        )}
      </button>
    </div>
  )
}

export default ChatInput
