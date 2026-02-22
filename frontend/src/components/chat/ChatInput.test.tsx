import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChatInput } from './ChatInput'

// i18n is mocked globally in src/test/setup.ts

describe('ChatInput', () => {
  const defaultProps = {
    onSend: vi.fn(),
    locale: 'en',
  }

  it('renders textarea field', () => {
    render(<ChatInput {...defaultProps} />)

    expect(screen.getByRole('textbox')).toBeInTheDocument()
  })

  it('renders send button', () => {
    render(<ChatInput {...defaultProps} />)

    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('calls onSend with input value', async () => {
    const mockSend = vi.fn()
    render(<ChatInput {...defaultProps} onSend={mockSend} />)

    const textarea = screen.getByRole('textbox')
    await userEvent.type(textarea, 'Show me revenue')
    await userEvent.click(screen.getByRole('button'))

    expect(mockSend).toHaveBeenCalledWith('Show me revenue')
  })

  it('clears input after send', async () => {
    const mockSend = vi.fn()
    render(<ChatInput {...defaultProps} onSend={mockSend} />)

    const textarea = screen.getByRole('textbox') as HTMLTextAreaElement
    await userEvent.type(textarea, 'Show me revenue')
    await userEvent.click(screen.getByRole('button'))

    expect(textarea.value).toBe('')
  })

  it('does not send empty input', async () => {
    const mockSend = vi.fn()
    render(<ChatInput {...defaultProps} onSend={mockSend} />)

    await userEvent.click(screen.getByRole('button'))

    expect(mockSend).not.toHaveBeenCalled()
  })

  it('submits on Enter key', async () => {
    const mockSend = vi.fn()
    render(<ChatInput {...defaultProps} onSend={mockSend} />)

    const textarea = screen.getByRole('textbox')
    await userEvent.type(textarea, 'Show me revenue{enter}')

    expect(mockSend).toHaveBeenCalledWith('Show me revenue')
  })

  it('can be disabled', () => {
    render(<ChatInput {...defaultProps} disabled={true} />)

    expect(screen.getByRole('textbox')).toBeDisabled()
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('applies RTL direction for Arabic locale', () => {
    render(<ChatInput {...defaultProps} locale="ar" />)

    const textarea = screen.getByRole('textbox')
    expect(textarea).toHaveAttribute('dir', 'rtl')
  })

  it('applies LTR direction for English locale', () => {
    render(<ChatInput {...defaultProps} locale="en" />)

    const textarea = screen.getByRole('textbox')
    expect(textarea).toHaveAttribute('dir', 'ltr')
  })
})
