import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

interface ChatInputProps {
  onSend: (message: string) => void;
  disabled?: boolean;
  locale: string;
  placeholder?: string;
  isDark?: boolean;
}

export const ChatInput: React.FC<ChatInputProps> = ({
  onSend,
  disabled = false,
  locale,
  placeholder,
  isDark = true,
}) => {
  const { t } = useTranslation('chat');
  const [input, setInput] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 200)}px`;
    }
  }, [input]);

  const handleSend = useCallback(() => {
    if (input.trim() && !disabled) {
      onSend(input.trim());
      setInput('');
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    }
  }, [input, disabled, onSend]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSend();
      }
    },
    [handleSend]
  );

  const isRTL = locale === 'ar';

  return (
    <div
      className={`flex items-end gap-3 rounded-xl border p-3 transition-colors duration-300 ${isDark
          ? 'bg-white/10 border-white/15 backdrop-blur-sm'
          : 'bg-white border-slate-200 shadow-sm'
        } ${isRTL ? 'flex-row-reverse' : ''}`}
    >
      <textarea
        ref={textareaRef}
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        placeholder={placeholder || t('input.placeholder')}
        rows={1}
        className={`flex-1 resize-none bg-transparent focus:outline-none text-sm ${isDark
            ? 'text-white placeholder-slate-500'
            : 'text-slate-900 placeholder-slate-400'
          } ${isRTL ? 'text-right' : 'text-left'}`}
        dir={isRTL ? 'rtl' : 'ltr'}
      />

      <button
        onClick={handleSend}
        disabled={disabled || !input.trim()}
        className={`flex-shrink-0 p-2 rounded-lg transition-all duration-300 ${disabled || !input.trim()
            ? isDark
              ? 'bg-white/10 text-slate-600 cursor-not-allowed'
              : 'bg-slate-100 text-slate-400 cursor-not-allowed'
            : 'bg-gradient-to-r from-blue-600 to-cyan-500 text-white shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 hover:scale-105 active:scale-95'
          }`}
        aria-label={t('input.send')}
      >
        {disabled ? (
          <svg
            className="w-5 h-5 animate-spin"
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
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        ) : (
          <svg
            className={`w-5 h-5 ${isRTL ? 'rotate-180' : ''}`}
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
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
  );
};

export default ChatInput;
