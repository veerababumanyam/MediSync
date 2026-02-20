import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, cleanup } from '@testing-library/react'
import { ChartRenderer } from './ChartRenderer'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: { language: 'en' },
  }),
}))

// Mock echarts
const mockChartInstance = {
  setOption: vi.fn(),
  resize: vi.fn(),
  dispose: vi.fn(),
}

vi.mock('echarts', () => ({
  default: {
    init: vi.fn(() => mockChartInstance),
  },
  init: vi.fn(() => mockChartInstance),
}))

describe('ChartRenderer', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    cleanup()
  })

  it('renders bar chart container', () => {
    const data = {
      labels: ['A', 'B', 'C'],
      series: [{ name: 'Sales', values: [100, 200, 300] }],
    }
    const { container } = render(
      <ChartRenderer chartType="barChart" data={data} locale="en" />
    )

    // Chart container should exist
    expect(container.querySelector('div[style*="height: 300px"]')).toBeInTheDocument()
  })

  it('renders line chart container', () => {
    const data = {
      labels: ['Jan', 'Feb', 'Mar'],
      series: [{ name: 'Revenue', values: [1000, 1500, 2000] }],
    }
    const { container } = render(
      <ChartRenderer chartType="lineChart" data={data} locale="en" />
    )

    expect(container.querySelector('div[style*="height: 300px"]')).toBeInTheDocument()
  })

  it('renders pie chart container', () => {
    const data = {
      labels: ['Category A', 'Category B'],
      series: [{ name: 'Distribution', values: [60, 40] }],
    }
    const { container } = render(
      <ChartRenderer chartType="pieChart" data={data} locale="en" />
    )

    expect(container.querySelector('div[style*="height: 300px"]')).toBeInTheDocument()
  })

  it('renders KPI card', () => {
    const data = {
      value: 1000000,
      formatted: '$1,000,000',
    }
    render(<ChartRenderer chartType="kpiCard" data={data} locale="en" />)

    expect(screen.getByText('$1,000,000')).toBeInTheDocument()
  })

  it('renders KPI card with raw value when no formatted', () => {
    const data = {
      value: 5000,
    }
    render(<ChartRenderer chartType="kpiCard" data={data} locale="en" />)

    expect(screen.getByText('5000')).toBeInTheDocument()
  })

  it('renders data table', () => {
    const data = {
      columns: [
        { name: 'Name', type: 'string' },
        { name: 'Value', type: 'number' },
      ],
      rows: [
        { Name: 'Item A', Value: 100 },
        { Name: 'Item B', Value: 200 },
      ],
    }
    render(<ChartRenderer chartType="dataTable" data={data} locale="en" />)

    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('Value')).toBeInTheDocument()
    expect(screen.getByText('Item A')).toBeInTheDocument()
    expect(screen.getByText('Item B')).toBeInTheDocument()
  })

  it('renders data table with number formatting', () => {
    const data = {
      columns: [
        { name: 'Item', type: 'string' },
        { name: 'Amount', type: 'number' },
      ],
      rows: [{ Item: 'Test', Amount: 1234567 }],
    }
    render(<ChartRenderer chartType="dataTable" data={data} locale="en" />)

    // Numbers should be formatted with toLocaleString
    expect(screen.getByText('1,234,567')).toBeInTheDocument()
  })

  it('renders data table with null values', () => {
    const data = {
      columns: [{ name: 'Field', type: 'string' }],
      rows: [{ Field: null }],
    }
    render(<ChartRenderer chartType="dataTable" data={data} locale="en" />)

    expect(screen.getByText('-')).toBeInTheDocument()
  })

  it('applies RTL direction for Arabic locale', () => {
    const data = {
      labels: ['A', 'B'],
      series: [{ name: 'Test', values: [1, 2] }],
    }
    const { container } = render(
      <ChartRenderer chartType="barChart" data={data} locale="ar" />
    )

    const chartDiv = container.querySelector('div[dir="rtl"]')
    expect(chartDiv).toBeInTheDocument()
  })

  it('applies LTR direction for English locale', () => {
    const data = {
      labels: ['A', 'B'],
      series: [{ name: 'Test', values: [1, 2] }],
    }
    const { container } = render(
      <ChartRenderer chartType="barChart" data={data} locale="en" />
    )

    const chartDiv = container.querySelector('div[dir="ltr"]')
    expect(chartDiv).toBeInTheDocument()
  })

  it('does not render data table without columns', () => {
    const data = {
      rows: [{ a: 1 }],
    }
    const { container } = render(
      <ChartRenderer chartType="dataTable" data={data} locale="en" />
    )

    // Should not render table
    expect(container.querySelector('table')).not.toBeInTheDocument()
  })

  it('does not render data table without rows', () => {
    const data = {
      columns: [{ name: 'Test', type: 'string' }],
    }
    const { container } = render(
      <ChartRenderer chartType="dataTable" data={data} locale="en" />
    )

    // Should not render table
    expect(container.querySelector('table')).not.toBeInTheDocument()
  })
})
