import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { v4 as uuidv4 } from 'uuid';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { StreamingMessage } from './StreamingMessage';
import { QuerySuggestions } from './QuerySuggestions';
import { ChatHeader } from './ChatHeader';
import { useLocale } from '../../hooks/useLocale';
import { apiClient, SSEEvent, ChatMessage } from '../../services/api';

interface ChatInterfaceProps {
  initialSessionId?: string;
  onSessionChange?: (sessionId: string) => void;
  className?: string;
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
    setError(null);    try {
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
    <div className={`flex flex-col h-full bg-white dark:bg-gray-900 ${className}`}>
      {/* Header */}
      <ChatHeader
        sessionId={sessionId}
        onNewSession={handleNewSession}
        locale={locale}
      />

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto px-4 py-6">
        {messages.length === 0 && !streaming.isStreaming ? (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <h2 className="text-2xl font-semibold text-gray-800 dark:text-white mb-4">
              {t('welcome.title')}
            </h2>
            <p className="text-gray-600 dark:text-gray-400 mb-8 max-w-md">
              {t('welcome.subtitle')}
            </p>
            <QuerySuggestions
              suggestions={suggestions}
              onSuggestionClick={handleSuggestionClick}
            />
          </div>
        ) : (
          <>
            <MessageList messages={messages} locale={locale} />
            {streaming.isStreaming && (
              <StreamingMessage
                events={streaming.events}
                locale={locale}
                onCancel={handleCancel}
              />
            )}
            {error && (
              <div className="mt-4 p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg">
                {error}
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Input Area */}
      <div className="border-t border-gray-200 dark:border-gray-700 p-4">
        <ChatInput
          onSend={handleSendMessage}
          disabled={streaming.isStreaming}
          locale={locale}
          placeholder={t('input.placeholder')}
        />
      </div>
    </div>
  );
};

export default ChatInterface;
