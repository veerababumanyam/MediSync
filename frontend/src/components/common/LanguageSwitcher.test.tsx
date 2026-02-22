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

// Mock usePreferences hook
vi.mock('../../hooks/usePreferences', () => ({
  usePreferences: () => ({
    updatePreferences: vi.fn().mockResolvedValue(undefined),
    isUpdating: false,
  }),
}))

describe('LanguageSwitcher', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders language switcher with pill style (default)', () => {
    render(<LanguageSwitcher />)

    // Should render as a radiogroup with EN and ع buttons
    expect(screen.getByRole('radiogroup')).toBeInTheDocument()
    expect(screen.getByText('EN')).toBeInTheDocument()
    expect(screen.getByText('ع')).toBeInTheDocument()
  })

  it('shows Arabic option', () => {
    render(<LanguageSwitcher />)

    // Should show Arabic character
    expect(screen.getByText('ع')).toBeInTheDocument()
  })

  it('calls changeLanguage when Arabic button clicked', async () => {
    render(<LanguageSwitcher />)

    await userEvent.click(screen.getByText('ع'))

    expect(mockChangeLanguage).toHaveBeenCalledWith('ar')
  })

  it('applies custom className', () => {
    const { container } = render(<LanguageSwitcher className="custom-class" />)

    expect(container.querySelector('.custom-class')).toBeInTheDocument()
  })

  it('renders in legacy button mode when pillStyle is false', () => {
    render(<LanguageSwitcher pillStyle={false} />)

    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  it('has accessible aria-label on radiogroup', () => {
    render(<LanguageSwitcher />)

    expect(screen.getByRole('radiogroup')).toHaveAttribute('aria-label')
  })
})
