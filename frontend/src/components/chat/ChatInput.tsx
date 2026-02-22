import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

interface ChatInputProps {
  onSend: (message: string) => void;
  disabled?: boolean;
  locale: string;
  placeholder?: string;
}

export const ChatInput: React.FC<ChatInputProps> = ({
  onSend,
  disabled = false,
  locale,
  placeholder,
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
  const sendLabel = t('input.send', 'Send');
  const placeholderText = placeholder || t('input.placeholder', 'Type your question...');

  return (
    <div
      dir={isRTL ? 'rtl' : 'ltr'}
      className={`relative flex items-end gap-3 p-4 rounded-2xl focus-within:ring-2 focus-within:ring-[#18929d]/50 dark:focus-within:ring-[#18929d]/30 bg-gradient-to-b from-white/80 to-white/60 dark:from-slate-800/80 dark:to-slate-800/60 shadow-[0_2px_8px_rgba(0,0,0,0.04),inset_0_1px_0_rgba(255,255,255,0.8)] dark:shadow-[0_2px_8px_rgba(0,0,0,0.2),inset_0_1px_0_rgba(255,255,255,0.08)] border border-white/60 dark:border-white/10 ${isRTL ? 'flex-row-reverse' : ''}`}
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
        className={`flex-1 min-h-[44px] resize-none bg-transparent text-slate-900 dark:text-slate-100 placeholder:text-slate-500 dark:placeholder:text-slate-400 focus:outline-none text-sm disabled:opacity-50 text-start ${isRTL ? 'rounded-e-xl' : 'rounded-s-xl'}`}
        dir={isRTL ? 'rtl' : 'ltr'}
      />

      <button
        type="button"
        onClick={handleSend}
        disabled={disabled || !input.trim()}
        className={`shrink-0 flex items-center justify-center min-h-[44px] min-w-[44px] p-3 rounded-xl transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#18929d] focus-visible:ring-offset-2 ${isRTL ? 'rounded-s-xl' : 'rounded-e-xl'} ${disabled || !input.trim()
            ? 'bg-slate-200/80 dark:bg-slate-700/80 text-slate-400 dark:text-slate-500 cursor-not-allowed'
            : 'liquid-glass-button-primary'
          }`}
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
            className={`w-5 h-5 ${isRTL ? 'rotate-180' : ''}`}
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
  );
};

export default ChatInput;

