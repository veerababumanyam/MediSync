/**
 * TableRenderer Component Tests
 *
 * Tests for tabular data rendering with sorting and RTL support.
 */
import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { TableRenderer } from './TableRenderer'

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (_key: string, fallback: string) => fallback,
    i18n: {
      language: 'en',
    },
  }),
}))

// Mock useLocale hook
vi.mock('../../hooks/useLocale', () => ({
  useLocale: () => ({
    locale: 'en',
    isRTL: false,
  }),
}))

describe('TableRenderer', () => {
  const mockData = [
    { id: 1, name: 'Alice', age: 30, city: 'New York' },
    { id: 2, name: 'Bob', age: 25, city: 'Los Angeles' },
    { id: 3, name: 'Charlie', age: 35, city: 'Chicago' },
  ]

  describe('rendering', () => {
    it('renders table with data', () => {
      render(<TableRenderer data={mockData} />)

      // Check headers are present
      expect(screen.getByText('id')).toBeInTheDocument()
      expect(screen.getByText('name')).toBeInTheDocument()
      expect(screen.getByText('age')).toBeInTheDocument()
      expect(screen.getByText('city')).toBeInTheDocument()

      // Check data is rendered
      expect(screen.getByText('Alice')).toBeInTheDocument()
      expect(screen.getByText('Bob')).toBeInTheDocument()
      expect(screen.getByText('Charlie')).toBeInTheDocument()
    })

    it('renders with custom columns', () => {
      render(
        <TableRenderer
          data={mockData}
          columns={['name', 'age']}
        />
      )

      // Only specified columns should be present
      expect(screen.getByText('name')).toBeInTheDocument()
      expect(screen.getByText('age')).toBeInTheDocument()
      expect(screen.queryByText('id')).not.toBeInTheDocument()
      expect(screen.queryByText('city')).not.toBeInTheDocument()
    })

    it('renders with custom column labels', () => {
      render(
        <TableRenderer
          data={mockData}
          columnLabels={{ name: 'Full Name', age: 'Age (years)' }}
        />
      )

      expect(screen.getByText('Full Name')).toBeInTheDocument()
      expect(screen.getByText('Age (years)')).toBeInTheDocument()
    })

    it('renders with row numbers when showRowNumbers is true', () => {
      render(
        <TableRenderer
          data={mockData}
          showRowNumbers
        />
      )

      // Check header has # symbol
      expect(screen.getByText('#')).toBeInTheDocument()

      // Verify rows have row numbers by checking first td of each row
      const rows = screen.getAllByRole('row')
      const dataRows = rows.slice(1) // Skip header

      // Check that each data row has a row number cell
      expect(dataRows).toHaveLength(3)
      dataRows.forEach((row, index) => {
        const cells = row.querySelectorAll('td')
        expect(cells[0]).toHaveTextContent(String(index + 1))
      })
    })

    it('limits rows when maxRows is set', () => {
      render(
        <TableRenderer
          data={mockData}
          maxRows={2}
        />
      )

      expect(screen.getByText('Alice')).toBeInTheDocument()
      expect(screen.getByText('Bob')).toBeInTheDocument()
      expect(screen.queryByText('Charlie')).not.toBeInTheDocument()
    })

    it('displays empty state when no data', () => {
      render(<TableRenderer data={[]} />)

      expect(screen.getByText('No data to display')).toBeInTheDocument()
    })

    it('formats numbers with locale separators', () => {
      const dataWithLargeNumbers = [
        { name: 'Revenue', value: 1000000 },
      ]

      render(<TableRenderer data={dataWithLargeNumbers} />)

      // Should have comma-separated number
      expect(screen.getByText('1,000,000')).toBeInTheDocument()
    })

    it('formats null values as dash', () => {
      const dataWithNull = [
        { name: 'Test', value: null },
      ]

      render(<TableRenderer data={dataWithNull} />)

      expect(screen.getAllByText('-')).toHaveLength(1)
    })

    it('formats boolean values', () => {
      const dataWithBoolean = [
        { name: 'Active', isActive: true },
        { name: 'Inactive', isActive: false },
      ]

      render(<TableRenderer data={dataWithBoolean} />)

      expect(screen.getByText('Yes')).toBeInTheDocument()
      expect(screen.getByText('No')).toBeInTheDocument()
    })
  })

  describe('sorting', () => {
    it('sorts by column when header is clicked', () => {
      render(<TableRenderer data={mockData} sortable />)

      const nameHeader = screen.getByText('name')
      fireEvent.click(nameHeader)

      // Check rows are in ascending order by name
      const rows = screen.getAllByRole('row')
      // Skip header row
      const dataRows = rows.slice(1)
      expect(dataRows[0]).toHaveTextContent('Alice')
      expect(dataRows[1]).toHaveTextContent('Bob')
      expect(dataRows[2]).toHaveTextContent('Charlie')
    })

    it('toggles sort direction on repeated clicks', () => {
      render(<TableRenderer data={mockData} sortable />)

      const nameHeader = screen.getByText('name')

      // First click: ascending
      fireEvent.click(nameHeader)
      let rows = screen.getAllByRole('row').slice(1)
      expect(rows[0]).toHaveTextContent('Alice')

      // Second click: descending
      fireEvent.click(nameHeader)
      rows = screen.getAllByRole('row').slice(1)
      expect(rows[0]).toHaveTextContent('Charlie')

      // Third click: no sort
      fireEvent.click(nameHeader)
      // Should return to original order
      rows = screen.getAllByRole('row').slice(1)
      expect(rows[0]).toHaveTextContent('Alice')
    })

    it('sorts numbers numerically', () => {
      render(<TableRenderer data={mockData} sortable />)

      const ageHeader = screen.getByText('age')
      fireEvent.click(ageHeader)

      const rows = screen.getAllByRole('row').slice(1)
      expect(rows[0]).toHaveTextContent('25')
      expect(rows[1]).toHaveTextContent('30')
      expect(rows[2]).toHaveTextContent('35')
    })

    it('disables sorting when sortable is false', () => {
      render(<TableRenderer data={mockData} sortable={false} />)

      const nameHeader = screen.getByText('name')
      fireEvent.click(nameHeader)

      // Order should remain unchanged
      const rows = screen.getAllByRole('row').slice(1)
      expect(rows[0]).toHaveTextContent('Alice')
    })
  })

  describe('interactivity', () => {
    it('calls onRowClick when row is clicked', () => {
      const handleRowClick = vi.fn()

      render(
        <TableRenderer
          data={mockData}
          onRowClick={handleRowClick}
        />
      )

      const bobRow = screen.getByText('Bob').closest('tr')
      if (bobRow) {
        fireEvent.click(bobRow)
      }

      expect(handleRowClick).toHaveBeenCalledTimes(1)
      expect(handleRowClick).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'Bob' }),
        1
      )
    })
  })

  describe('accessibility', () => {
    it('has proper table structure', () => {
      render(<TableRenderer data={mockData} />)

      expect(screen.getByRole('table')).toBeInTheDocument()
      expect(screen.getAllByRole('columnheader')).toHaveLength(4)
      expect(screen.getAllByRole('row')).toHaveLength(4) // 1 header + 3 data rows
    })

    it('has aria-sort attribute on sortable columns', () => {
      render(<TableRenderer data={mockData} sortable />)

      const nameHeader = screen.getByText('name').closest('th')
      expect(nameHeader).toHaveAttribute('aria-sort', 'none')

      fireEvent.click(screen.getByText('name'))
      expect(nameHeader).toHaveAttribute('aria-sort', 'ascending')

      fireEvent.click(screen.getByText('name'))
      expect(nameHeader).toHaveAttribute('aria-sort', 'descending')
    })
  })
})
