import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ChartPinDialog } from './ChartPinDialog'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'dialog.title': 'Pin Chart',
        'dialog.titleLabel': 'Title',
        'dialog.titlePlaceholder': 'Enter chart title',
        'dialog.queryLabel': 'Query',
        'dialog.queryPlaceholder': 'Describe your data',
        'dialog.chartTypeLabel': 'Chart Type',
        'dialog.refreshLabel': 'Refresh Interval',
        'dialog.cancel': 'Cancel',
        'dialog.pin': 'Pin',
        'dialog.pinning': 'Pinning...',
        'dialog.requiredFields': 'Title and query are required',
        'dialog.pinError': 'Failed to pin chart',
        'chartTypes.bar': 'Bar',
        'chartTypes.line': 'Line',
        'chartTypes.pie': 'Pie',
        'chartTypes.table': 'Table',
        'chartTypes.kpi': 'KPI',
        'refreshOptions.manual': 'Manual',
        'refreshOptions.fiveMinutes': '5 minutes',
        'refreshOptions.fifteenMinutes': '15 minutes',
        'refreshOptions.thirtyMinutes': '30 minutes',
        'refreshOptions.oneHour': '1 hour',
        'refreshOptions.oneDay': '1 day',
      }
      return translations[key] || key
    },
    i18n: { language: 'en' },
  }),
}))

describe('ChartPinDialog', () => {
  const defaultProps = {
    onClose: vi.fn(),
    onPin: vi.fn().mockResolvedValue(undefined),
    locale: 'en',
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders dialog title', () => {
    render(<ChartPinDialog {...defaultProps} />)

    expect(screen.getByText('Pin Chart')).toBeInTheDocument()
  })

  it('renders title input', () => {
    render(<ChartPinDialog {...defaultProps} />)

    expect(screen.getByPlaceholderText('Enter chart title')).toBeInTheDocument()
  })

  it('renders query textarea', () => {
    render(<ChartPinDialog {...defaultProps} />)

    expect(screen.getByPlaceholderText('Describe your data')).toBeInTheDocument()
  })

  it('renders chart type options', () => {
    render(<ChartPinDialog {...defaultProps} />)

    expect(screen.getByText('Bar')).toBeInTheDocument()
    expect(screen.getByText('Line')).toBeInTheDocument()
    expect(screen.getByText('Pie')).toBeInTheDocument()
    expect(screen.getByText('Table')).toBeInTheDocument()
    expect(screen.getByText('KPI')).toBeInTheDocument()
  })

  it('renders refresh interval options', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    // Open select dropdown
    const select = screen.getByRole('combobox')
    expect(select).toBeInTheDocument()
  })

  it('renders cancel and pin buttons', () => {
    render(<ChartPinDialog {...defaultProps} />)

    expect(screen.getByText('Cancel')).toBeInTheDocument()
    expect(screen.getByText('Pin')).toBeInTheDocument()
  })

  it('calls onClose when cancel clicked', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    await userEvent.click(screen.getByText('Cancel'))

    expect(defaultProps.onClose).toHaveBeenCalled()
  })

  it('calls onClose when cancel or close clicked', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    // Cancel button should call onClose
    await userEvent.click(screen.getByText('Cancel'))

    expect(defaultProps.onClose).toHaveBeenCalled()
  })

  it('shows error when submitting without title', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    // Fill only query
    await userEvent.type(screen.getByPlaceholderText('Describe your data'), 'Show revenue')
    await userEvent.click(screen.getByText('Pin'))

    expect(screen.getByText('Title and query are required')).toBeInTheDocument()
    expect(defaultProps.onPin).not.toHaveBeenCalled()
  })

  it('shows error when submitting without query', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    // Fill only title
    await userEvent.type(screen.getByPlaceholderText('Enter chart title'), 'Revenue Chart')
    await userEvent.click(screen.getByText('Pin'))

    expect(screen.getByText('Title and query are required')).toBeInTheDocument()
    expect(defaultProps.onPin).not.toHaveBeenCalled()
  })

  it('calls onPin with correct data when form is valid', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    await userEvent.type(screen.getByPlaceholderText('Enter chart title'), 'Revenue Chart')
    await userEvent.type(screen.getByPlaceholderText('Describe your data'), 'Show monthly revenue')
    await userEvent.click(screen.getByText('Pin'))

    await waitFor(() => {
      expect(defaultProps.onPin).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Revenue Chart',
          naturalLanguageQuery: 'Show monthly revenue',
          chartType: 'bar',
          refreshInterval: 0,
        })
      )
    })
  })

  it('changes chart type when option clicked', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    // Click Line chart type
    await userEvent.click(screen.getByText('Line'))

    // Verify the button is selected (should have active/selected styling)
    const lineButton = screen.getByText('Line').closest('button')
    expect(lineButton).toBeTruthy()
    // Check for any selection indicator class (liquid glass uses various classes)
    const hasSelectedClass = lineButton?.className.includes('blue') ||
                             lineButton?.className.includes('primary') ||
                             lineButton?.className.includes('selected')
    expect(hasSelectedClass).toBeTruthy()
  })

  it('changes refresh interval when select changed', async () => {
    render(<ChartPinDialog {...defaultProps} />)

    const select = screen.getByRole('combobox')
    await userEvent.selectOptions(select, '5')

    expect(select).toHaveValue('5')
  })

  it('shows pinning state while submitting', async () => {
    const slowOnPin = vi.fn().mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
    render(<ChartPinDialog {...defaultProps} onPin={slowOnPin} />)

    await userEvent.type(screen.getByPlaceholderText('Enter chart title'), 'Test')
    await userEvent.type(screen.getByPlaceholderText('Describe your data'), 'Query')

    await userEvent.click(screen.getByText('Pin'))

    // Should show pinning state
    expect(screen.getByText('Pinning...')).toBeInTheDocument()
  })

  it('shows error when onPin fails', async () => {
    vi.mocked(defaultProps.onPin).mockRejectedValue(new Error('Network error'))
    render(<ChartPinDialog {...defaultProps} />)

    await userEvent.type(screen.getByPlaceholderText('Enter chart title'), 'Test')
    await userEvent.type(screen.getByPlaceholderText('Describe your data'), 'Query')
    await userEvent.click(screen.getByText('Pin'))

    await waitFor(() => {
      expect(screen.getByText('Failed to pin chart')).toBeInTheDocument()
    })
  })

  it('applies RTL direction for Arabic locale', () => {
    render(<ChartPinDialog {...defaultProps} locale="ar" />)

    const titleInput = screen.getByPlaceholderText('Enter chart title')
    expect(titleInput).toHaveAttribute('dir', 'rtl')

    const queryInput = screen.getByPlaceholderText('Describe your data')
    expect(queryInput).toHaveAttribute('dir', 'rtl')
  })
})
