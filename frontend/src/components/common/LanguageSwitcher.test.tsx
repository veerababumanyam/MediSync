import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { LanguageSwitcher } from './LanguageSwitcher'

// Mock i18next
const mockChangeLanguage = vi.fn().mockResolvedValue(undefined)

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: {
      language: 'en',
      changeLanguage: mockChangeLanguage,
    },
  }),
}))

describe('LanguageSwitcher', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders language switcher button', () => {
    render(<LanguageSwitcher />)

    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('shows next language option', () => {
    render(<LanguageSwitcher />)

    // Should show Arabic option since current is English
    expect(screen.getByText(/عربي/i)).toBeInTheDocument()
  })

  it('calls changeLanguage when clicked', async () => {
    render(<LanguageSwitcher />)

    await userEvent.click(screen.getByRole('button'))

    expect(mockChangeLanguage).toHaveBeenCalledWith('ar')
  })

  it('applies custom className', () => {
    const { container } = render(<LanguageSwitcher className="custom-class" />)

    expect(container.querySelector('.custom-class')).toBeInTheDocument()
  })

  it('renders in compact mode', () => {
    render(<LanguageSwitcher compact={true} />)

    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('has accessible aria-label', () => {
    render(<LanguageSwitcher />)

    expect(screen.getByRole('button')).toHaveAttribute('aria-label')
  })
})
