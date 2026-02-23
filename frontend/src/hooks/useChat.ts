/**
 * useChat Hook
 *
 * React hook for managing chat state with:
 * - Message state management
 * - Streaming responses via SSE
 * - Session management
 * - Error handling
 *
 * @module hooks/useChat
 */
import { useCallback, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  apiClient,
  APIError,
} from '../services/api'
import type { ChatMessage, SSEEvent } from '../services/api'

/**
 * Internal message structure with optional streaming state
 */
export interface Message extends Omit<ChatMessage, 'createdAt'> {
  /** ISO timestamp */
  createdAt: string
  /** Whether the message is still being streamed */
  isStreaming?: boolean
  /** Intermediate content during streaming */
  partialContent?: string
}

/**
 * Hook return type
 */
export interface UseChatReturn {
  /** Array of chat messages */
  messages: Message[]
  /** Whether a message is being sent/processed */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
  /** Send a query to the chat API */
  sendMessage: (query: string) => Promise<void>
  /** Clear all messages */
  clearMessages: () => void
  /** Current session ID */
  sessionId: string
  /** Abort the current streaming request */
  abort: () => void
}

/**
 * Storage key for session persistence
 */
const SESSION_STORAGE_KEY = 'AnySync-chat-session'

/**
 * Get or create a session ID
 */
function getOrCreateSessionId(): string {
  const stored = sessionStorage.getItem(SESSION_STORAGE_KEY)
  if (stored) {
    return stored
  }
  const newId = crypto.randomUUID()
  sessionStorage.setItem(SESSION_STORAGE_KEY, newId)
  return newId
}

/**
 * Hook for managing chat interactions with streaming support
 */
export function useChat(): UseChatReturn {
  const { t } = useTranslation('chat')
  const [messages, setMessages] = useState<Message[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [sessionId] = useState<string>(() => getOrCreateSessionId())

  // Abort controller ref for cancelling requests
  const abortControllerRef = useRef<AbortController | null>(null)

  /**
   * Handle SSE events during streaming
   */
  const handleSSEEvent = useCallback((event: SSEEvent, assistantMessageId: string) => {
    setMessages(prev => prev.map(msg => {
      if (msg.id !== assistantMessageId) return msg

      switch (event.type) {
        case 'thinking':
          return {
            ...msg,
            partialContent: event.message || t('thinking'),
            isStreaming: true,
          }

        case 'sql_preview':
          return {
            ...msg,
            partialContent: event.sql || msg.partialContent,
            isStreaming: true,
          }

        case 'result':
          return {
            ...msg,
            content: event.message || '',
            chartSpec: event.chartType && event.data ? {
              type: event.chartType,
              chart: event.data,
            } : undefined,
            confidence: event.confidence,
            isStreaming: false,
            partialContent: undefined,
          }

        case 'error':
          return {
            ...msg,
            content: event.message || t('error.genericError'),
            isStreaming: false,
            partialContent: undefined,
          }

        case 'clarification':
          return {
            ...msg,
            content: event.message || t('error.clarification'),
            isStreaming: false,
            partialContent: undefined,
          }

        default:
          return msg
      }
    }))
  }, [t])

  /**
   * Send a query to the chat API with streaming
   */
  const sendMessage = useCallback(async (query: string) => {
    if (!query.trim()) {
      return
    }

    // Cancel any ongoing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    // Create new abort controller
    abortControllerRef.current = new AbortController()

    setError(null)
    setIsLoading(true)

    // Create user message
    const userMessage: Message = {
      id: crypto.randomUUID(),
      sessionId,
      role: 'user',
      content: query.trim(),
      createdAt: new Date().toISOString(),
    }

    // Create placeholder assistant message
    const assistantMessageId = crypto.randomUUID()
    const assistantMessage: Message = {
      id: assistantMessageId,
      sessionId,
      role: 'assistant',
      content: '',
      createdAt: new Date().toISOString(),
      isStreaming: true,
    }

    // Add both messages
    setMessages(prev => [...prev, userMessage, assistantMessage])

    try {
      const result = await apiClient.streamChat(
        {
          query: query.trim(),
          sessionId,
          locale: t('locale', 'en'), // Optional: send locale context if supported
        },
        (event) => handleSSEEvent(event, assistantMessageId),
        abortControllerRef.current.signal
      )

      // If we got a result, ensure the message is updated
      if (result) {
        setMessages(prev => prev.map(msg => {
          if (msg.id !== assistantMessageId) return msg
          return {
            ...msg,
            content: result.message || msg.content,
            chartSpec: result.chartSpec || msg.chartSpec,
            confidence: result.confidence,
            isStreaming: false,
            partialContent: undefined,
          }
        }))
      }
    } catch (err) {
      // Ignore abort errors
      if (err instanceof Error && err.name === 'AbortError') {
        return
      }

      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : t('error.send')

      setError(message)

      // Update assistant message with error
      setMessages(prev => prev.map(msg => {
        if (msg.id !== assistantMessageId) return msg
        return {
          ...msg,
          content: message,
          isStreaming: false,
          partialContent: undefined,
        }
      }))
    } finally {
      setIsLoading(false)
      abortControllerRef.current = null
    }
  }, [sessionId, handleSSEEvent, t])

  /**
   * Clear all messages
   */
  const clearMessages = useCallback(() => {
    setMessages([])
    setError(null)
  }, [])

  /**
   * Abort the current streaming request
   */
  const abort = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
      abortControllerRef.current = null
    }
    setIsLoading(false)
  }, [])

  return {
    messages,
    isLoading,
    error,
    sendMessage,
    clearMessages,
    sessionId,
    abort,
  }
}

export default useChat
