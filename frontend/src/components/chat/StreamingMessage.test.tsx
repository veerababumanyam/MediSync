import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { StreamingMessage } from './StreamingMessage'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: { language: 'en' },
  }),
}))

// Mock ChartRenderer
vi.mock('./ChartRenderer', () => ({
  ChartRenderer: ({ chartSpec }: { chartSpec: unknown }) => (
    <div data-testid="chart-renderer">{JSON.stringify(chartSpec)}</div>
  ),
}))

describe('StreamingMessage', () => {
  const mockOnCancel = vi.fn()

  it('renders thinking event', () => {
    const events = [{ type: 'thinking' as const, message: 'Analyzing your query...' }]
    render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    // Thinking message appears twice (status + event), use getAllByText
    expect(screen.getAllByText('Analyzing your query...').length).toBeGreaterThan(0)
  })

  it('renders sql_preview event', () => {
    const events = [{ type: 'sql_preview' as const, sql: 'SELECT * FROM revenue' }]
    render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    expect(screen.getByText('SELECT * FROM revenue')).toBeInTheDocument()
  })

  it('renders result event with chart', () => {
    const events = [
      {
        type: 'result' as const,
        message: 'Here are the results',
        chartType: 'bar',
        data: [{ name: 'A', value: 100 }],
      },
    ]
    render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    expect(screen.getByTestId('chart-renderer')).toBeInTheDocument()
  })

  it('renders error event', () => {
    const events = [{ type: 'error' as const, message: 'Something went wrong' }]
    render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
  })

  it('renders empty events gracefully', () => {
    const { container } = render(<StreamingMessage events={[]} locale="en" onCancel={mockOnCancel} />)

    expect(container).toBeDefined()
  })

  it('renders multiple events', () => {
    const events = [
      { type: 'thinking' as const, message: 'Analyzing data...' },
      { type: 'sql_preview' as const, sql: 'SELECT * FROM users' },
    ]
    const { container } = render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    // Verify container renders with events
    expect(container.firstChild).toBeInTheDocument()
    // Check for SQL preview which should be unique
    expect(screen.getByText('SELECT * FROM users')).toBeInTheDocument()
  })

  it('renders cancel button', () => {
    const events = [{ type: 'thinking' as const, message: 'Processing...' }]
    render(<StreamingMessage events={events} locale="en" onCancel={mockOnCancel} />)

    // Cancel button should be present
    const buttons = screen.getAllByRole('button')
    expect(buttons.length).toBeGreaterThan(0)
  })
})
