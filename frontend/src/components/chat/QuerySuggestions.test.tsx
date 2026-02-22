import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QuerySuggestions } from './QuerySuggestions'

// i18n is mocked globally in src/test/setup.ts

describe('QuerySuggestions', () => {
  const mockSuggestions = [
    'Show me revenue',
    'Top departments by sales',
    'Patient count this month',
  ]

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders suggestions', () => {
    const mockClick = vi.fn()
    render(<QuerySuggestions suggestions={mockSuggestions} onSuggestionClick={mockClick} />)

    expect(screen.getByText('Show me revenue')).toBeInTheDocument()
    expect(screen.getByText('Top departments by sales')).toBeInTheDocument()
    expect(screen.getByText('Patient count this month')).toBeInTheDocument()
  })

  it('calls onSuggestionClick when suggestion clicked', async () => {
    const mockClick = vi.fn()
    render(<QuerySuggestions suggestions={mockSuggestions} onSuggestionClick={mockClick} />)

    await userEvent.click(screen.getByText('Show me revenue'))

    expect(mockClick).toHaveBeenCalledWith('Show me revenue')
  })

  it('renders empty state when no suggestions', () => {
    const mockClick = vi.fn()
    const { container } = render(<QuerySuggestions suggestions={[]} onSuggestionClick={mockClick} />)

    // Should render container but no suggestion buttons
    expect(container.querySelector('div')).toBeInTheDocument()
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('renders all suggestion buttons', () => {
    const mockClick = vi.fn()
    render(<QuerySuggestions suggestions={mockSuggestions} onSuggestionClick={mockClick} />)

    const buttons = screen.getAllByRole('button')
    expect(buttons).toHaveLength(3)
  })
})
