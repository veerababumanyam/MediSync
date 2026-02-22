import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { PinnedChartCard } from './PinnedChartCard'
import type { PinnedChart } from '../../services/api'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) => {
      const translations: Record<string, string> = {
        never: 'Never',
        justNow: 'Just now',
        minutesAgo: `${options?.count || 0} minutes ago`,
        hoursAgo: `${options?.count || 0} hours ago`,
        refresh: 'Refresh',
        unpin: 'Unpin',
        delete: 'Delete',
        lastRefreshed: 'Last refreshed',
        refreshesEvery: `Refreshes every ${options?.minutes || 0} min`,
      }
      return translations[key] || key
    },
    i18n: { language: 'en' },
  }),
}))

// Mock ChartRenderer
vi.mock('../chat/ChartRenderer', () => ({
  ChartRenderer: ({ chartType }: { chartType: string }) => (
    <div data-testid={`chart-${chartType}`}>Chart: {chartType}</div>
  ),
}))

describe('PinnedChartCard', () => {
  const mockChart: PinnedChart = {
    id: '1',
    userId: 'user-1',
    title: 'Revenue Overview',
    queryId: 'query-1',
    naturalLanguageQuery: 'Show revenue',
    sqlQuery: 'SELECT * FROM revenue',
    chartType: 'bar',
    chartSpec: { labels: ['A', 'B'], series: [{ name: 'Revenue', values: [100, 200] }] },
    refreshInterval: 5,
    locale: 'en',
    position: { row: 0, col: 0, size: 1 },
    lastRefreshedAt: new Date().toISOString(),
    isActive: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  }

  const defaultProps = {
    chart: mockChart,
    locale: 'en',
    onDelete: vi.fn(),
    onRefresh: vi.fn().mockResolvedValue(undefined),
    onToggle: vi.fn(),
    onClick: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders chart title', () => {
    render(<PinnedChartCard {...defaultProps} />)

    expect(screen.getByText('Revenue Overview')).toBeInTheDocument()
  })

  it('renders ChartRenderer with correct type', () => {
    render(<PinnedChartCard {...defaultProps} />)

    expect(screen.getByTestId('chart-bar')).toBeInTheDocument()
  })

  it('shows menu when menu button clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)

    expect(screen.getByText('Refresh')).toBeInTheDocument()
    expect(screen.getByText('Unpin')).toBeInTheDocument()
    expect(screen.getByText('Delete')).toBeInTheDocument()
  })

  it('calls onDelete when delete clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)
    await userEvent.click(screen.getByText('Delete'))

    expect(defaultProps.onDelete).toHaveBeenCalled()
  })

  it('calls onRefresh when refresh clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)
    await userEvent.click(screen.getByText('Refresh'))

    await waitFor(() => {
      expect(defaultProps.onRefresh).toHaveBeenCalled()
    })
  })

  it('calls onToggle with false when unpin clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)
    await userEvent.click(screen.getByText('Unpin'))

    expect(defaultProps.onToggle).toHaveBeenCalledWith(false)
  })

  it('calls onClick when card clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    // Click on the card (not the menu button)
    await userEvent.click(screen.getByText('Revenue Overview'))

    expect(defaultProps.onClick).toHaveBeenCalled()
  })

  it('does not call onClick when menu button clicked', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)

    // onClick should not be called when clicking the menu button
    expect(defaultProps.onClick).not.toHaveBeenCalled()
  })

  it('disables refresh button while refreshing', async () => {
    const slowRefresh = vi.fn().mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
    render(<PinnedChartCard {...defaultProps} onRefresh={slowRefresh} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)

    const refreshButton = screen.getByText('Refresh')
    await userEvent.click(refreshButton)

    // Button should be disabled during refresh
    expect(refreshButton).toBeDisabled()
  })

  it('closes menu after action', async () => {
    render(<PinnedChartCard {...defaultProps} />)

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)
    await userEvent.click(screen.getByText('Delete'))

    // Menu should be closed
    expect(screen.queryByText('Unpin')).not.toBeInTheDocument()
  })

  it('displays "just now" for recent refresh', () => {
    const chart = {
      ...mockChart,
      lastRefreshedAt: new Date().toISOString(),
    }
    render(<PinnedChartCard {...defaultProps} chart={chart} />)

    expect(screen.getByText(/Just now/i)).toBeInTheDocument()
  })

  it('displays "never" for null lastRefreshedAt', () => {
    const chart = {
      ...mockChart,
      lastRefreshedAt: null,
    }
    render(<PinnedChartCard {...defaultProps} chart={chart} />)

    expect(screen.getByText(/Never/i)).toBeInTheDocument()
  })

  it('displays refresh interval when set', () => {
    render(<PinnedChartCard {...defaultProps} />)

    expect(screen.getByText(/5 min/)).toBeInTheDocument()
  })

  it('does not display refresh interval when zero', () => {
    const chart = {
      ...mockChart,
      refreshInterval: 0,
    }
    render(<PinnedChartCard {...defaultProps} chart={chart} />)

    expect(screen.queryByText(/min/)).not.toBeInTheDocument()
  })

  it('has cursor-pointer class when onClick provided', () => {
    const { container } = render(<PinnedChartCard {...defaultProps} />)

    expect(container.querySelector('.cursor-pointer')).toBeInTheDocument()
  })

  it('does not have cursor-pointer class when onClick not provided', () => {
    const { container } = render(<PinnedChartCard {...defaultProps} onClick={undefined} />)

    expect(container.querySelector('.cursor-pointer')).not.toBeInTheDocument()
  })

  it('closes menu when clicking outside', async () => {
    render(
      <div>
        <PinnedChartCard {...defaultProps} />
        <div data-testid="outside">Outside</div>
      </div>
    )

    const menuButton = screen.getByRole('button')
    await userEvent.click(menuButton)

    expect(screen.getByText('Unpin')).toBeInTheDocument()

    // Click on the fixed overlay (which is part of the menu close mechanism)
    const overlay = document.querySelector('.fixed.inset-0')
    if (overlay) {
      await userEvent.click(overlay)
    }

    // Menu should be closed
    expect(screen.queryByRole('button', { name: /unpin/i })).not.toBeInTheDocument()
  })
})
