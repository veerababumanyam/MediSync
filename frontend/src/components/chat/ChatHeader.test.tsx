import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { ChatHeader } from './ChatHeader'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: { language: 'en' },
  }),
}))

// Mock LanguageSwitcher
vi.mock('../common/LanguageSwitcher', () => ({
  LanguageSwitcher: () => <div data-testid="language-switcher">Language Switcher</div>,
}))

describe('ChatHeader', () => {
  const defaultProps = {
    sessionId: '12345678-1234-1234-1234-123456789abc',
    onNewSession: vi.fn(),
    locale: 'en',
  }

  it('renders header element', () => {
    const { container } = render(<ChatHeader {...defaultProps} />)

    expect(container.querySelector('header')).toBeInTheDocument()
  })

  it('renders language switcher', () => {
    render(<ChatHeader {...defaultProps} />)

    expect(screen.getByTestId('language-switcher')).toBeInTheDocument()
  })

  it('renders new session button', () => {
    render(<ChatHeader {...defaultProps} />)

    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('displays shortened session ID', () => {
    render(<ChatHeader {...defaultProps} />)

    expect(screen.getByText(/12345678\.\.\./)).toBeInTheDocument()
  })

  it('applies correct styling classes', () => {
    const { container } = render(<ChatHeader {...defaultProps} />)

    const header = container.querySelector('header')
    expect(header).toHaveClass('flex')
    expect(header).toHaveClass('items-center')
  })
})
