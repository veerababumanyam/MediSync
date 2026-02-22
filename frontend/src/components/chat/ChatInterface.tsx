import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { v4 as uuidv4 } from 'uuid';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { StreamingMessage } from './StreamingMessage';
import { QuerySuggestions } from './QuerySuggestions';
import { ChatHeader } from './ChatHeader';
import { useLocale } from '../../hooks/useLocale';
import { apiClient, type SSEEvent, type ChatMessage } from '../../services/api';

interface ChatInterfaceProps {
  initialSessionId?: string;
  onSessionChange?: (sessionId: string) => void;
  className?: string;
  isDark?: boolean;
}

interface StreamingState {
  isStreaming: boolean;
  events: SSEEvent[];
  currentMessage: string;
}

export const ChatInterface: React.FC<ChatInterfaceProps> = ({
  initialSessionId,
  onSessionChange,
  className = '',
  isDark = true,
}) => {
  const { t } = useTranslation('chat');
  const { locale } = useLocale();

  const [sessionId, setSessionId] = useState(initialSessionId || uuidv4());
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [streaming, setStreaming] = useState<StreamingState>({
    isStreaming: false,
    events: [],
    currentMessage: '',
  });
  const [error, setError] = useState<string | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, streaming.events]);

  // Load existing messages when session changes
  useEffect(() => {
    if (sessionId) {
      loadMessages(sessionId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sessionId]);

  const loadMessages = async (sid: string) => {
    setError(null); try {
      const response = await apiClient.getChatMessages(sid);
      setMessages(response.messages);
    } catch (err) {
      console.error('Failed to load messages:', err);
      setError(t('error.loadMessages'));
    }
  };

  const handleSendMessage = useCallback(async (query: string) => {
    if (!query.trim() || streaming.isStreaming) return;

    setError(null);

    // Add user message immediately
    const userMessage: ChatMessage = {
      id: uuidv4(),
      sessionId,
      role: 'user',
      content: query,
      createdAt: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, userMessage]);

    // Start streaming
    setStreaming({
      isStreaming: true,
      events: [],
      currentMessage: '',
    });

    // Create abort controller for this request
    abortControllerRef.current = new AbortController();

    try {
      const eventHandler = (event: SSEEvent) => {
        setStreaming((prev) => ({
          ...prev,
          events: [...prev.events, event],
          currentMessage: event.message || prev.currentMessage,
        }));
      };

      const result = await apiClient.streamChat(
        {
          query,
          sessionId,
          locale,
        },
        eventHandler,
        abortControllerRef.current.signal
      );

      // Add assistant message from result
      if (result) {
        const assistantMessage: ChatMessage = {
          id: uuidv4(),
          sessionId,
          role: 'assistant',
          content: result.message || 'Query completed',
          chartSpec: result.chartSpec,
          confidence: result.confidence,
          createdAt: new Date().toISOString(),
        };
        setMessages((prev) => [...prev, assistantMessage]);
      }
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        console.log('Request aborted');
      } else {
        console.error('Chat error:', err);
        setError(t('error.send'));
      }
    } finally {
      setStreaming({
        isStreaming: false,
        events: [],
        currentMessage: '',
      });
      abortControllerRef.current = null;
    }
  }, [sessionId, locale, streaming.isStreaming, t]);

  const handleSuggestionClick = useCallback((suggestion: string) => {
    handleSendMessage(suggestion);
  }, [handleSendMessage]);

  const handleNewSession = useCallback(() => {
    const newSessionId = uuidv4();
    setSessionId(newSessionId);
    setMessages([]);
    setError(null);
    onSessionChange?.(newSessionId);
  }, [onSessionChange]);

  const handleCancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }, []);

  // Default query suggestions
  const suggestions = [
    t('suggestions.revenue'),
    t('suggestions.topDepartments'),
    t('suggestions.patientTrend'),
    t('suggestions.inventory'),
  ];

  return (
    <div className={`flex flex-col min-h-0 ${className}`}>
      {/* Header */}
      <ChatHeader
        sessionId={sessionId}
        onNewSession={handleNewSession}
        isDark={isDark}
      />

      {/* Messages Area — takes all space; input stays at bottom */}
      <div className="flex-1 min-h-0 overflow-y-auto chat-scrollbar px-3 py-3 sm:px-6 sm:py-4">
        {messages.length === 0 && !streaming.isStreaming ? (
          <div className="flex flex-col items-center justify-center h-full text-center px-2">
            <h2 className={`text-xl sm:text-2xl font-semibold mb-3 sm:mb-4 ${isDark ? 'text-white' : 'text-slate-900'
              }`}>
              {t('welcome.title')}
            </h2>
            <p className={`mb-6 sm:mb-8 max-w-md text-sm sm:text-base ${isDark ? 'text-slate-400' : 'text-slate-600'
              }`}>
              {t('welcome.subtitle')}
            </p>
            <QuerySuggestions
              suggestions={suggestions}
              onSuggestionClick={handleSuggestionClick}
              isDark={isDark}
            />
          </div>
        ) : (
          <>
            <MessageList messages={messages} locale={locale} isDark={isDark} />
            {streaming.isStreaming && (
              <StreamingMessage
                events={streaming.events}
                locale={locale}
                onCancel={handleCancel}
                isDark={isDark}
              />
            )}
            {error && (
              <div className={`glass-subtle mt-3 sm:mt-4 p-3 sm:p-4 rounded-xl text-sm ${isDark
                ? 'border-red-500/30 text-red-400'
                : 'border-red-300 text-red-600'
                }`}>
                {error}
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Input Area — pinned to bottom, never shrinks */}
      <div className="flex-shrink-0 glass-border border-t p-2 sm:p-3 lg:p-4 safe-bottom">
        <ChatInput
          onSend={handleSendMessage}
          disabled={streaming.isStreaming}
          locale={locale}
          placeholder={t('input.placeholder')}
          isDark={isDark}
        />
      </div>
    </div>
  );
};

export default ChatInterface;
