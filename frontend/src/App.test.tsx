import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

// Mock CopilotKit
vi.mock('@copilotkit/react-core', () => ({
  CopilotKit: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))

// Mock i18next with proper interpolation support
const mockChangeLanguage = vi.fn().mockResolvedValue(undefined)

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: string | Record<string, unknown>) => {
      if (typeof options === 'string') return options
      if (typeof options === 'object' && options !== null) {
        let text = (options.defaultValue as string) ?? key
        Object.entries(options).forEach(([k, v]) => {
          if (k !== 'defaultValue') {
            text = text.replace(new RegExp(`\\{\\{\\s*${k}\\s*\\}\\}`, 'g'), String(v))
          }
        })
        return text
      }
      return key
    },
    i18n: {
      language: 'en',
      dir: () => 'ltr',
      changeLanguage: mockChangeLanguage,
    },
  }),
  Trans: ({ children }: { children: React.ReactNode }) => children,
}))

// Mock i18n initialization
vi.mock('./i18n', () => ({}))

// Mock useTheme hook
vi.mock('./hooks/useTheme', () => ({
  useTheme: () => ({
    isDark: true,
    toggleTheme: vi.fn(),
  }),
}))

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    document.documentElement.dir = 'ltr'
    document.documentElement.lang = 'en'
  })

  it('renders without crashing', () => {
    const { container } = render(<App />)
    expect(container).toBeDefined()
  })

  it('renders app name', () => {
    render(<App />)
    expect(screen.getByText('MediSync')).toBeInTheDocument()
  })

  it('renders app tagline', () => {
    render(<App />)
    expect(screen.getByText('AI-Powered Business Intelligence')).toBeInTheDocument()
  })

  it('renders hero section', () => {
    render(<App />)
    expect(screen.getByText(/Your Data/i)).toBeInTheDocument()
  })

  it('renders language toggle button', () => {
    render(<App />)
    expect(screen.getByText('عربي')).toBeInTheDocument()
  })

  it('calls changeLanguage when language toggle clicked', async () => {
    render(<App />)

    await userEvent.click(screen.getByText('عربي'))

    expect(mockChangeLanguage).toHaveBeenCalledWith('ar')
  })

  it('renders capability cards', () => {
    render(<App />)

    expect(screen.getByText('Conversational BI')).toBeInTheDocument()
    expect(screen.getByText('AI Accountant')).toBeInTheDocument()
    expect(screen.getByText('Smart Reports')).toBeInTheDocument()
  })

  it('renders capability card descriptions', () => {
    render(<App />)

    expect(screen.getByText(/Ask questions in plain English/)).toBeInTheDocument()
    expect(screen.getByText(/OCR extracts ledger data/)).toBeInTheDocument()
    expect(screen.getByText(/Pre-built MIS reports/)).toBeInTheDocument()
  })

  it('sets document direction to ltr for English', async () => {
    render(<App />)

    await waitFor(() => {
      expect(document.documentElement.dir).toBe('ltr')
    })
  })

  it('sets document language to en', async () => {
    render(<App />)

    await waitFor(() => {
      expect(document.documentElement.lang).toBe('en')
    })
  })
})
