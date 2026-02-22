import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MessageList } from './MessageList'

// i18n is mocked globally in src/test/setup.ts

// Mock ChartRenderer
vi.mock('./ChartRenderer', () => ({
  ChartRenderer: ({ chartSpec }: { chartSpec: unknown }) => (
    <div data-testid="chart-renderer">{JSON.stringify(chartSpec)}</div>
  ),
}))

describe('MessageList', () => {
  const mockMessages = [
    {
      id: '1',
      sessionId: 'session-1',
      role: 'user' as const,
      content: 'Show me revenue',
      createdAt: '2024-01-01T10:00:00Z',
    },
    {
      id: '2',
      sessionId: 'session-1',
      role: 'assistant' as const,
      content: 'Here is the revenue data',
      chartSpec: { type: 'bar', chart: { chartType: 'bar' } },
      createdAt: '2024-01-01T10:01:00Z',
    },
  ]

  it('renders messages', () => {
    render(<MessageList messages={mockMessages} locale="en" />)

    expect(screen.getByText('Show me revenue')).toBeInTheDocument()
    expect(screen.getByText('Here is the revenue data')).toBeInTheDocument()
  })

  it('renders chart when present', () => {
    render(<MessageList messages={mockMessages} locale="en" />)

    expect(screen.getByTestId('chart-renderer')).toBeInTheDocument()
  })

  it('renders empty state when no messages', () => {
    const { container } = render(<MessageList messages={[]} locale="en" />)

    expect(container.querySelector('.space-y-3')).toBeInTheDocument()
  })

  it('applies different styling for user vs assistant messages', () => {
    const { container } = render(<MessageList messages={mockMessages} locale="en" />)

    const userBubble = container.querySelector('.from-blue-600')
    const assistantBubble = container.querySelector('.rounded-bl-md')

    expect(userBubble).toBeInTheDocument()
    expect(assistantBubble).toBeInTheDocument()
  })

  it('renders message container', () => {
    const { container } = render(<MessageList messages={mockMessages} locale="en" />)

    expect(container.querySelector('.space-y-3')).toBeInTheDocument()
  })
})
