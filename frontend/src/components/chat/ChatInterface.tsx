import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { v4 as uuidv4 } from 'uuid';
import { MessageList } from './MessageList';
import { ChatInput } from './ChatInput';
import { StreamingMessage } from './StreamingMessage';
import { QuerySuggestions } from './QuerySuggestions';
import { ChatHeader } from './ChatHeader';
import { useLocale } from '../../hooks/useLocale';
import { FadeIn, StaggerChildren } from '../animations';
import type { SSEEvent, ChatMessage } from '../../services/api';
import { apiClient, api } from '../../services/api';
import { webMCPService } from '../../services/WebMCPService';

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

  // Register WebMCP tools
  useEffect(() => {
    webMCPService.registerChatTools({
      onQuery: (query: string) => {
        handleSendMessage(query);
      },
      onSyncTally: async () => {
        try {
          await api.post('/sync/now');
        } catch (err) {
          console.error('WebMCP syncTally failed:', err);
          throw err;
        }
      },
      onShowDashboard: () => {
        console.log(`WebMCP: Navigating to dashboard`);
        window.location.href = '/dashboard';
      }
    });

    return () => {
      webMCPService.cleanup();
    };
  }, [handleSendMessage]);

  // Default query suggestions
  const suggestions = [
    t('suggestions.revenue'),
    t('suggestions.topDepartments'),
    t('suggestions.patientTrend'),
    t('suggestions.inventory'),
  ];

  return (
    <div className={`flex flex-col h-full bg-transparent ${className}`}>
      {/* Visually hidden H1 for proper heading hierarchy */}
      <h1 className="sr-only">{t('pageTitle')}</h1>

      {/* Utility Row - Liquid Glass Subtle */}
      <div className="liquid-glass-subtle border-b border-white/10 dark:border-white/5">
        <ChatHeader
          onNewSession={handleNewSession}
        />
      </div>

      {/* Messages Area */}
      <div id="chat-messages" className="flex-1 overflow-y-auto px-4 py-6 liquid-glass-scroll">
        {messages.length === 0 && !streaming.isStreaming ? (
          <FadeIn>
            <div className="flex flex-col items-center justify-center h-full text-center">
              {/* Welcome Card - iOS 26 Liquid Glass */}
              <div className="liquid-glass-content-card p-8 max-w-2xl mb-8">
                {/* AI Icon with Brand Gradient */}
                <div className="w-12 h-12 rounded-xl flex items-center justify-center mx-auto mb-6 shadow-lg" style={{
                  background: 'linear-gradient(135deg, #2750a8 0%, #18929d 100%)',
                  boxShadow: '0 8px 24px rgba(39, 80, 168, 0.3)'
                }}>
                  <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                </div>
                <h2 className="hero-gradient-text text-xl font-semibold mb-4">
                  {t('welcome.title')}
                </h2>
                <p className="liquid-text-secondary mb-8 max-w-md mx-auto">
                  {t('welcome.subtitle')}
                </p>
                <StaggerChildren className="w-full">
                  <QuerySuggestions
                    suggestions={suggestions}
                    onSuggestionClick={handleSuggestionClick}
                  />
                </StaggerChildren>
              </div>
            </div>
          </FadeIn>
        ) : (
          <>
            <MessageList messages={messages} locale={locale} />
            {streaming.isStreaming && (
              <FadeIn>
                <StreamingMessage
                  events={streaming.events}
                  locale={locale}
                  onCancel={handleCancel}
                />
              </FadeIn>
            )}
            {error && (
              <FadeIn>
                <div
                  className="liquid-glass-content-card mt-4 p-4 border-s-4 border-s-red-500"
                  role="alert"
                  aria-live="assertive"
                >
                  <div className="flex items-center gap-3">
                    <svg className="w-5 h-5 text-red-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="text-red-600 dark:text-red-400 flex-1">{error}</span>
                  </div>
                </div>
              </FadeIn>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Input Area - Liquid Glass */}
      <div className="liquid-glass-light border-t border-white/20 dark:border-white/10 p-4">
        <ChatInput
          onSend={handleSendMessage}
          disabled={streaming.isStreaming}
          locale={locale}
          placeholder={t('input.placeholder')}
          {...({
            'tool-name': 'medi-chat-input',
            'tool-description': 'The text input for natural language BI queries'
          } as React.HTMLAttributes<HTMLDivElement>)}
        />
      </div>
    </div>
  );
};

export default ChatInterface;
