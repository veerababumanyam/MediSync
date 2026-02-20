import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ExportButton, ExportData } from './ExportButton'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: { language: 'en' },
  }),
}))

// Mock LoadingSpinner
vi.mock('./LoadingSpinner', () => ({
  LoadingSpinner: () => <div data-testid="loading-spinner">Loading...</div>,
}))

describe('ExportButton', () => {
  const mockData: ExportData = {
    query: 'Show me revenue',
    results: [{ month: 'Jan', revenue: 1000 }],
  }

  beforeEach(() => {
    vi.clearAllMocks()
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      blob: () => Promise.resolve(new Blob()),
    })
    global.URL.createObjectURL = vi.fn(() => 'blob:test')
    global.URL.revokeObjectURL = vi.fn()
  })

  it('renders without crashing', () => {
    const { container } = render(<ExportButton data={mockData} />)
    expect(container).toBeDefined()
  })

  it('renders Export button', () => {
    render(<ExportButton data={mockData} />)
    expect(screen.getByText('Export')).toBeInTheDocument()
  })

  it('applies disabled prop', () => {
    render(<ExportButton data={mockData} disabled />)
    const button = screen.getByText('Export').closest('button')
    expect(button).toBeDisabled()
  })

  it('applies custom className', () => {
    const { container } = render(<ExportButton data={mockData} className="test-class" />)
    expect(container.querySelector('.test-class')).toBeTruthy()
  })

  it('opens dropdown when clicked', async () => {
    render(<ExportButton data={mockData} />)
    await userEvent.click(screen.getByText('Export'))
    expect(screen.getByText('PDF')).toBeInTheDocument()
  })

  it('has correct aria attributes', () => {
    render(<ExportButton data={mockData} />)
    const button = screen.getByText('Export').closest('button')
    expect(button).toHaveAttribute('aria-haspopup', 'listbox')
  })
})
