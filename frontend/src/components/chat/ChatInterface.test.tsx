import { describe, it, expect, vi, beforeEach, beforeAll } from 'vitest'
import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChatInterface } from './ChatInterface'
import { apiClient } from '../../services/api'

// Mock scrollIntoView for jsdom
beforeAll(() => {
  Element.prototype.scrollIntoView = vi.fn()
})

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => {
      const translations: Record<string, string> = {
        'suggestions.revenue': 'Show me revenue',
        'suggestions.topDepartments': 'Top departments',
        'suggestions.patientTrend': 'Patient trend',
        'suggestions.inventory': 'Inventory status',
        'welcome.title': 'Welcome to MediSync',
        'welcome.subtitle': 'Ask questions about your data',
        'input.placeholder': 'Type your question...',
        'error.loadMessages': 'Failed to load messages',
        'error.send': 'Failed to send message',
      }
      return translations[key] || defaultValue || key
    },
    i18n: { language: 'en' },
  }),
}))

// Mock useLocale hook
vi.mock('../../hooks/useLocale', () => ({
  useLocale: () => ({
    locale: 'en',
    isRTL: false,
    setLocale: vi.fn(),
    toggleLocale: vi.fn(),
  }),
}))

// Mock uuid
vi.mock('uuid', () => ({
  v4: vi.fn(() => 'test-uuid-1234'),
}))

// Mock apiClient
vi.mock('../../services/api', () => ({
  apiClient: {
    getChatMessages: vi.fn().mockResolvedValue({ messages: [] }),
    streamChat: vi.fn().mockResolvedValue({
      message: 'Query completed',
      chartSpec: null,
      confidence: 0.95,
    }),
  },
  SSEEvent: vi.fn(),
  ChatMessage: vi.fn(),
}))

// Mock ChatHeader
vi.mock('./ChatHeader', () => ({
  ChatHeader: ({ sessionId }: { sessionId: string }) => (
    <div data-testid="chat-header">Session: {sessionId}</div>
  ),
}))

// Mock ChatInput
vi.mock('./ChatInput', () => ({
  ChatInput: ({ onSend, disabled, placeholder }: {
    onSend: (msg: string) => void;
    disabled: boolean;
    placeholder: string;
  }) => (
    <div data-testid="chat-input">
      <input
        type="text"
        placeholder={placeholder}
        disabled={disabled}
        onKeyDown={(e) => {
          if (e.key === 'Enter' && e.currentTarget.value) {
            onSend(e.currentTarget.value)
          }
        }}
      />
      <button
        disabled={disabled}
        onClick={() => onSend('Test message')}
      >
        Send
      </button>
    </div>
  ),
}))

// Mock MessageList
vi.mock('./MessageList', () => ({
  MessageList: ({ messages }: { messages: Array<{ id: string; content: string }> }) => (
    <div data-testid="message-list">
      {messages.map(m => (
        <div key={m.id} data-testid={`message-${m.id}`}>{m.content}</div>
      ))}
    </div>
  ),
}))

// Mock StreamingMessage
vi.mock('./StreamingMessage', () => ({
  StreamingMessage: ({ events, onCancel }: {
    events: Array<{ type: string }>;
    onCancel: () => void;
  }) => (
    <div data-testid="streaming-message">
      Events: {events.length}
      <button onClick={onCancel}>Cancel</button>
    </div>
  ),
}))

// Mock QuerySuggestions
vi.mock('./QuerySuggestions', () => ({
  QuerySuggestions: ({ suggestions, onSuggestionClick }: {
    suggestions: string[];
    onSuggestionClick: (s: string) => void;
  }) => (
    <div data-testid="query-suggestions">
      {suggestions.map((s, i) => (
        <button key={i} onClick={() => onSuggestionClick(s)}>{s}</button>
      ))}
    </div>
  ),
}))

describe('ChatInterface', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(apiClient.getChatMessages).mockResolvedValue({ messages: [] })
    vi.mocked(apiClient.streamChat).mockResolvedValue({
      message: 'Query completed',
      chartSpec: null,
      confidence: 0.95,
    })
  })

  it('renders welcome screen when no messages', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Welcome to MediSync')).toBeInTheDocument()
    })
    expect(screen.getByText('Ask questions about your data')).toBeInTheDocument()
  })

  it('renders query suggestions on welcome screen', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByTestId('query-suggestions')).toBeInTheDocument()
    })
    expect(screen.getByText('Show me revenue')).toBeInTheDocument()
  })

  it('renders chat header with session ID', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByTestId('chat-header')).toBeInTheDocument()
    })
  })

  it('renders chat input', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByTestId('chat-input')).toBeInTheDocument()
    })
  })

  it('sends message when send button clicked', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Send')).toBeInTheDocument()
    })

    await act(async () => {
      await userEvent.click(screen.getByText('Send'))
    })

    expect(apiClient.streamChat).toHaveBeenCalled()
  })

  it('shows streaming message while streaming', async () => {
    // Make streamChat hang to see streaming state
    vi.mocked(apiClient.streamChat).mockImplementation(() => new Promise(() => {}))

    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Send')).toBeInTheDocument()
    })

    await act(async () => {
      await userEvent.click(screen.getByText('Send'))
    })

    await waitFor(() => {
      expect(screen.getByTestId('streaming-message')).toBeInTheDocument()
    })
  })

  it('disables input while streaming', async () => {
    vi.mocked(apiClient.streamChat).mockImplementation(() => new Promise(() => {}))

    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Send')).toBeInTheDocument()
    })

    await act(async () => {
      await userEvent.click(screen.getByText('Send'))
    })

    await waitFor(() => {
      expect(screen.getByText('Send')).toBeDisabled()
    })
  })

  it('loads messages for initial session ID', async () => {
    render(<ChatInterface initialSessionId="existing-session-123" />)

    await waitFor(() => {
      expect(apiClient.getChatMessages).toHaveBeenCalledWith('existing-session-123')
    })
  })

  it('applies custom className', async () => {
    const { container } = render(<ChatInterface className="custom-chat" />)

    await waitFor(() => {
      expect(container.querySelector('.custom-chat')).toBeInTheDocument()
    })
  })

  it('handles suggestion click', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Show me revenue')).toBeInTheDocument()
    })

    await act(async () => {
      await userEvent.click(screen.getByText('Show me revenue'))
    })

    expect(apiClient.streamChat).toHaveBeenCalledWith(
      expect.objectContaining({ query: 'Show me revenue' }),
      expect.any(Function),
      expect.any(AbortSignal)
    )
  })

  it('uses initial session ID when provided', async () => {
    render(<ChatInterface initialSessionId="my-session-id" />)

    await waitFor(() => {
      expect(screen.getByText('Session: my-session-id')).toBeInTheDocument()
    })
  })

  it('renders message list after successful send', async () => {
    render(<ChatInterface />)

    await waitFor(() => {
      expect(screen.getByText('Send')).toBeInTheDocument()
    })

    await act(async () => {
      await userEvent.click(screen.getByText('Send'))
    })

    // After message is sent, MessageList should appear
    await waitFor(() => {
      expect(screen.getByTestId('message-list')).toBeInTheDocument()
    })
  })
})
