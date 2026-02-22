import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useChat } from './useChat'
import { APIError } from '../services/api'

// Mock uuid - let the hook manage actual UUIDs, we'll just check structure
vi.mock('uuid', () => ({
  v4: vi.fn().mockImplementation(() => `uuid-${Date.now()}-${Math.random()}`),
}))

// Mock apiClient
const mockStreamChat = vi.fn()

vi.mock('../services/api', () => ({
  apiClient: {
    streamChat: (...args: unknown[]) => mockStreamChat(...args),
  },
  APIError: class APIError extends Error {
    status: number
    statusText: string
    constructor(status: number, statusText: string, message: string) {
      super(message)
      this.name = 'APIError'
      this.status = status
      this.statusText = statusText
    }
  },
}))

describe('useChat', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    sessionStorage.clear()
    mockStreamChat.mockReset()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('returns initial state with empty messages', () => {
    const { result } = renderHook(() => useChat())

    expect(result.current.messages).toEqual([])
    expect(result.current.isLoading).toBe(false)
    expect(result.current.error).toBeNull()
    expect(result.current.sessionId).toBeDefined()
    expect(typeof result.current.sessionId).toBe('string')
    expect(result.current.sessionId.length).toBeGreaterThan(0)
  })

  it('provides sendMessage function', () => {
    const { result } = renderHook(() => useChat())

    expect(typeof result.current.sendMessage).toBe('function')
  })

  it('provides clearMessages function', () => {
    const { result } = renderHook(() => useChat())

    expect(typeof result.current.clearMessages).toBe('function')
  })

  it('provides abort function', () => {
    const { result } = renderHook(() => useChat())

    expect(typeof result.current.abort).toBe('function')
  })

  it('generates and persists a session ID', () => {
    const { result } = renderHook(() => useChat())

    const sessionId = result.current.sessionId
    expect(sessionId).toBeDefined()
    expect(typeof sessionId).toBe('string')
    expect(sessionStorage.getItem('medisync-chat-session')).toBe(sessionId)
  })

  it('reuses existing session ID from sessionStorage', () => {
    const existingId = 'existing-session-id'
    sessionStorage.setItem('medisync-chat-session', existingId)

    const { result } = renderHook(() => useChat())

    expect(result.current.sessionId).toBe(existingId)
  })

  it('has correct return type structure', () => {
    const { result } = renderHook(() => useChat())

    expect(result.current).toHaveProperty('messages')
    expect(result.current).toHaveProperty('isLoading')
    expect(result.current).toHaveProperty('error')
    expect(result.current).toHaveProperty('sendMessage')
    expect(result.current).toHaveProperty('clearMessages')
    expect(result.current).toHaveProperty('sessionId')
    expect(result.current).toHaveProperty('abort')
  })

  describe('sendMessage', () => {
    it('does nothing for empty query', async () => {
      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('')
      })

      expect(mockStreamChat).not.toHaveBeenCalled()
      expect(result.current.messages).toHaveLength(0)
    })

    it('does nothing for whitespace-only query', async () => {
      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('   ')
      })

      expect(mockStreamChat).not.toHaveBeenCalled()
      expect(result.current.messages).toHaveLength(0)
    })

    it('adds user message immediately', async () => {
      mockStreamChat.mockResolvedValue({ message: 'Response' })

      const { result } = renderHook(() => useChat())
      const sessionId = result.current.sessionId

      act(() => {
        result.current.sendMessage('Hello')
      })

      // Check user message was added immediately
      await waitFor(() => {
        expect(result.current.messages).toHaveLength(2)
      })

      const userMessage = result.current.messages.find(m => m.role === 'user')
      expect(userMessage).toBeDefined()
      expect(userMessage?.content).toBe('Hello')
      expect(userMessage?.sessionId).toBe(sessionId)
    })

    it('adds placeholder assistant message', async () => {
      let resolveStream: (value: unknown) => void
      mockStreamChat.mockImplementation(() => new Promise(resolve => {
        resolveStream = resolve
      }))

      const { result } = renderHook(() => useChat())

      act(() => {
        result.current.sendMessage('Hello')
      })

      // Wait for messages to be added
      await waitFor(() => {
        expect(result.current.messages).toHaveLength(2)
      })

      const assistantMessage = result.current.messages.find(m => m.role === 'assistant')
      expect(assistantMessage).toBeDefined()
      expect(assistantMessage?.isStreaming).toBe(true)

      // Clean up the pending promise
      await act(async () => {
        resolveStream!({ message: 'Response' })
      })
    })

    it('sets isLoading to true while sending', async () => {
      let resolveStream: (value: unknown) => void
      mockStreamChat.mockImplementation(() => new Promise(resolve => {
        resolveStream = resolve
      }))

      const { result } = renderHook(() => useChat())

      act(() => {
        result.current.sendMessage('Hello')
      })

      expect(result.current.isLoading).toBe(true)

      await act(async () => {
        resolveStream!({ message: 'Response' })
      })

      expect(result.current.isLoading).toBe(false)
    })

    it('calls streamChat with correct parameters', async () => {
      mockStreamChat.mockResolvedValue({ message: 'Response' })

      const { result } = renderHook(() => useChat())
      const sessionId = result.current.sessionId

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(mockStreamChat).toHaveBeenCalledWith(
        expect.objectContaining({
          query: 'Hello',
          sessionId: sessionId,
        }),
        expect.any(Function),
        expect.any(AbortSignal)
      )
    })

    it('trims whitespace from query', async () => {
      mockStreamChat.mockResolvedValue({ message: 'Response' })

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('  Hello  ')
      })

      expect(mockStreamChat).toHaveBeenCalledWith(
        expect.objectContaining({
          query: 'Hello',
        }),
        expect.any(Function),
        expect.any(AbortSignal)
      )
    })

    it('handles SSE events correctly', async () => {
      let eventHandler: (event: unknown) => void
      mockStreamChat.mockImplementation(async (_req, onEvent) => {
        eventHandler = onEvent
        return { message: 'Final response' }
      })

      const { result } = renderHook(() => useChat())

      act(() => {
        result.current.sendMessage('Hello')
      })

      await waitFor(() => {
        expect(result.current.messages).toHaveLength(2)
      })

      // Simulate thinking event
      act(() => {
        eventHandler!({ type: 'thinking', message: 'Analyzing...' })
      })

      await waitFor(() => {
        const assistant = result.current.messages.find(m => m.role === 'assistant')
        expect(assistant?.partialContent).toBe('Analyzing...')
        expect(assistant?.isStreaming).toBe(true)
      })

      // Simulate sql_preview event
      act(() => {
        eventHandler!({ type: 'sql_preview', sql: 'SELECT * FROM users' })
      })

      await waitFor(() => {
        const assistant = result.current.messages.find(m => m.role === 'assistant')
        expect(assistant?.partialContent).toBe('SELECT * FROM users')
      })

      // Simulate result event
      act(() => {
        eventHandler!({
          type: 'result',
          message: 'Found 10 users',
          chartType: 'bar',
          data: { labels: [], values: [] },
          confidence: 0.95,
        })
      })

      await waitFor(() => {
        const assistant = result.current.messages.find(m => m.role === 'assistant')
        expect(assistant?.content).toBe('Found 10 users')
        expect(assistant?.chartSpec).toEqual({ type: 'bar', chart: { labels: [], values: [] } })
        expect(assistant?.confidence).toBe(0.95)
        expect(assistant?.isStreaming).toBe(false)
      })
    })

    it('handles error SSE event', async () => {
      let eventHandler: (event: unknown) => void
      mockStreamChat.mockImplementation(async (_req, onEvent) => {
        eventHandler = onEvent
        return null
      })

      const { result } = renderHook(() => useChat())

      await act(async () => {
        result.current.sendMessage('Hello')
      })

      // Simulate error event
      act(() => {
        eventHandler!({ type: 'error', message: 'Something went wrong' })
      })

      await waitFor(() => {
        const assistant = result.current.messages.find(m => m.role === 'assistant')
        expect(assistant?.content).toBe('Something went wrong')
        expect(assistant?.isStreaming).toBe(false)
      })
    })

    it('handles clarification SSE event', async () => {
      let eventHandler: (event: unknown) => void
      mockStreamChat.mockImplementation(async (_req, onEvent) => {
        eventHandler = onEvent
        return null
      })

      const { result } = renderHook(() => useChat())

      await act(async () => {
        result.current.sendMessage('Hello')
      })

      // Simulate clarification event
      act(() => {
        eventHandler!({ type: 'clarification', message: 'Did you mean X or Y?' })
      })

      await waitFor(() => {
        const assistant = result.current.messages.find(m => m.role === 'assistant')
        expect(assistant?.content).toBe('Did you mean X or Y?')
        expect(assistant?.isStreaming).toBe(false)
      })
    })

    it('handles APIError correctly', async () => {
      const apiError = new APIError(500, 'Internal Server Error', 'Server error occurred')
      mockStreamChat.mockRejectedValue(apiError)

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.error).toBe('Server error occurred')

      const assistant = result.current.messages.find(m => m.role === 'assistant')
      expect(assistant?.content).toBe('Server error occurred')
    })

    it('handles generic error correctly', async () => {
      mockStreamChat.mockRejectedValue(new Error('Network error'))

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.error).toBe('Network error')

      const assistant = result.current.messages.find(m => m.role === 'assistant')
      expect(assistant?.content).toBe('Network error')
    })

    it('handles non-Error rejection', async () => {
      mockStreamChat.mockRejectedValue('Unknown error')

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.error).toBe('Failed to send message')
    })

    it('ignores abort errors', async () => {
      const abortError = new Error('Aborted')
      abortError.name = 'AbortError'
      mockStreamChat.mockRejectedValue(abortError)

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.error).toBeNull()
    })

    it('sets isLoading to false after error', async () => {
      mockStreamChat.mockRejectedValue(new Error('Error'))

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.isLoading).toBe(false)
    })
  })

  describe('clearMessages', () => {
    it('clears all messages', async () => {
      mockStreamChat.mockResolvedValue({ message: 'Response' })

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.messages).toHaveLength(2)

      act(() => {
        result.current.clearMessages()
      })

      expect(result.current.messages).toHaveLength(0)
    })

    it('clears error state', async () => {
      mockStreamChat.mockRejectedValue(new Error('Error'))

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello')
      })

      expect(result.current.error).toBe('Error')

      act(() => {
        result.current.clearMessages()
      })

      expect(result.current.error).toBeNull()
    })
  })

  describe('abort', () => {
    it('aborts ongoing request', async () => {
      let resolveStream: (value: unknown) => void
      mockStreamChat.mockImplementation(() => new Promise(resolve => {
        resolveStream = resolve
      }))

      const { result } = renderHook(() => useChat())

      act(() => {
        result.current.sendMessage('Hello')
      })

      expect(result.current.isLoading).toBe(true)

      act(() => {
        result.current.abort()
      })

      expect(result.current.isLoading).toBe(false)

      // Resolve the promise to clean up
      await act(async () => {
        resolveStream!({ message: 'Response' })
      })
    })
  })

  describe('multiple messages', () => {
    it('appends messages correctly', async () => {
      mockStreamChat.mockResolvedValue({ message: 'Response 1' })

      const { result } = renderHook(() => useChat())

      await act(async () => {
        await result.current.sendMessage('Hello 1')
      })

      expect(result.current.messages).toHaveLength(2)

      mockStreamChat.mockResolvedValue({ message: 'Response 2' })

      await act(async () => {
        await result.current.sendMessage('Hello 2')
      })

      expect(result.current.messages).toHaveLength(4)

      // Check order
      expect(result.current.messages[0].role).toBe('user')
      expect(result.current.messages[0].content).toBe('Hello 1')
      expect(result.current.messages[1].role).toBe('assistant')
      expect(result.current.messages[2].role).toBe('user')
      expect(result.current.messages[2].content).toBe('Hello 2')
      expect(result.current.messages[3].role).toBe('assistant')
    })
  })
})
