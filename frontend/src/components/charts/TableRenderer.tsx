/**
 * TableRenderer Component
 *
 * Renders tabular data from query results with sorting and RTL support.
 * Used to display SQL query results in a user-friendly format.
 *
 * @module components/charts/TableRenderer
 */
import React, { useMemo, useState, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocale } from '../../hooks/useLocale'

/**
 * Sort direction type
 */
type SortDirection = 'asc' | 'desc' | null

/**
 * Sort state for a column
 */
interface SortState {
  column: string | null
  direction: SortDirection
}

/**
 * Props for the TableRenderer component
 */
export interface TableRendererProps {
  /** Data to render as an array of objects */
  data: Record<string, unknown>[]
  /** Optional column headers (defaults to object keys) */
  columns?: string[]
  /** Optional custom column labels */
  columnLabels?: Record<string, string>
  /** Maximum rows to display (0 for all) */
  maxRows?: number
  /** Whether to show row numbers */
  showRowNumbers?: boolean
  /** Whether to enable sorting */
  sortable?: boolean
  /** Optional CSS class name */
  className?: string
  /** Callback when a row is clicked */
  onRowClick?: (row: Record<string, unknown>, index: number) => void
}

/**
 * TableRenderer - Renders tabular data with sorting capability
 *
 * Features:
 * - Automatic column detection
 * - Sortable columns
 * - RTL support
 * - Custom column labels
 * - Row click handling
 * - Responsive design
 * - Dark mode support
 */
export const TableRenderer: React.FC<TableRendererProps> = ({
  data,
  columns: propColumns,
  columnLabels = {},
  maxRows = 0,
  showRowNumbers = false,
  sortable = true,
  className = '',
  onRowClick,
}) => {
  const { t } = useTranslation('common')
  const { isRTL } = useLocale()
  const [sortState, setSortState] = useState<SortState>({
    column: null,
    direction: null,
  })

  // Extract columns from data if not provided
  const columns = useMemo(() => {
    if (propColumns && propColumns.length > 0) {
      return propColumns
    }
    if (data.length === 0) {
      return []
    }
    return Object.keys(data[0])
  }, [propColumns, data])

  // Sort data based on current sort state
  const sortedData = useMemo(() => {
    if (!sortState.column || !sortState.direction) {
      return data
    }

    return [...data].sort((a, b) => {
      const aVal = a[sortState.column!]
      const bVal = b[sortState.column!]

      // Handle null/undefined
      if (aVal == null && bVal == null) return 0
      if (aVal == null) return sortState.direction === 'asc' ? 1 : -1
      if (bVal == null) return sortState.direction === 'asc' ? -1 : 1

      // Compare values
      let comparison = 0
      if (typeof aVal === 'number' && typeof bVal === 'number') {
        comparison = aVal - bVal
      } else {
        comparison = String(aVal).localeCompare(String(bVal))
      }

      return sortState.direction === 'asc' ? comparison : -comparison
    })
  }, [data, sortState])

  // Limit rows if maxRows is set
  const displayData = useMemo(() => {
    if (maxRows > 0) {
      return sortedData.slice(0, maxRows)
    }
    return sortedData
  }, [sortedData, maxRows])

  // Handle sort toggle
  const handleSort = useCallback((column: string) => {
    if (!sortable) return

    setSortState((prev) => {
      if (prev.column !== column) {
        return { column, direction: 'asc' }
      }

      if (prev.direction === 'asc') {
        return { column, direction: 'desc' }
      }

      return { column: null, direction: null }
    })
  }, [sortable])

  // Format cell value for display
  const formatValue = useCallback((value: unknown): string => {
    if (value == null) return '-'
    if (typeof value === 'number') {
      // Format numbers with locale-appropriate separators
      return value.toLocaleString(isRTL ? 'ar-SA' : 'en-US')
    }
    if (typeof value === 'boolean') {
      return value ? t('yes', 'Yes') : t('no', 'No')
    }
    if (typeof value === 'object') {
      return JSON.stringify(value)
    }
    return String(value)
  }, [isRTL, t])

  // Get sort icon for column header
  const getSortIcon = useCallback((column: string) => {
    if (sortState.column !== column) {
      return (
        <svg
          className="w-4 h-4 text-slate-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
          />
        </svg>
      )
    }

    if (sortState.direction === 'asc') {
      return (
        <svg
          className="w-4 h-4 text-blue-600"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M5 15l7-7 7 7"
          />
        </svg>
      )
    }

    return (
      <svg
        className="w-4 h-4 text-blue-600"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M19 9l-7 7-7-7"
        />
      </svg>
    )
  }, [sortState])

  // Empty state
  if (data.length === 0) {
    return (
      <div
        className={`flex items-center justify-center py-12 px-4 bg-slate-50 dark:bg-slate-800/50 rounded-lg ${className}`}
      >
        <div className="text-center">
          <svg
            className="w-12 h-12 mx-auto text-slate-400 dark:text-slate-600"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
            />
          </svg>
          <p className="mt-4 text-sm text-slate-500 dark:text-slate-400">
            {t('table.noData', 'No data to display')}
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className={`overflow-x-auto ${className}`}>
      <table className="min-w-full divide-y divide-slate-200 dark:divide-slate-700">
        <thead className="bg-slate-50 dark:bg-slate-800/50">
          <tr>
            {/* Row number column */}
            {showRowNumbers && (
              <th
                className="px-4 py-3 text-start text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider"
                scope="col"
              >
                #
              </th>
            )}

            {/* Data columns */}
            {columns.map((column) => (
              <th
                key={column}
                className={`px-4 py-3 text-start text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider ${
                  sortable ? 'cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors select-none' : ''
                }`}
                scope="col"
                onClick={() => handleSort(column)}
                aria-sort={
                  sortState.column === column
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  <span>{columnLabels[column] || column}</span>
                  {sortable && getSortIcon(column)}
                </div>
              </th>
            ))}
          </tr>
        </thead>

        <tbody className="bg-white dark:bg-slate-900 divide-y divide-slate-200 dark:divide-slate-700">
          {displayData.map((row, rowIndex) => (
            <tr
              key={rowIndex}
              className={`${
                onRowClick
                  ? 'cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors'
                  : ''
              }`}
              onClick={() => onRowClick?.(row, rowIndex)}
            >
              {/* Row number cell */}
              {showRowNumbers && (
                <td className="px-4 py-3 text-sm text-slate-500 dark:text-slate-400 font-mono">
                  {rowIndex + 1}
                </td>
              )}

              {/* Data cells */}
              {columns.map((column) => (
                <td
                  key={column}
                  className="px-4 py-3 text-sm text-slate-700 dark:text-slate-300 whitespace-nowrap"
                >
                  {formatValue(row[column])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>

        {/* Footer with row count */}
        {maxRows > 0 && sortedData.length > maxRows && (
          <tfoot className="bg-slate-50 dark:bg-slate-800/50">
            <tr>
              <td
                colSpan={columns.length + (showRowNumbers ? 1 : 0)}
                className="px-4 py-2 text-sm text-slate-500 dark:text-slate-400 text-center"
              >
                {t('table.showingRows', 'Showing {{count}} of {{total}} rows', {
                  count: maxRows,
                  total: sortedData.length,
                })}
              </td>
            </tr>
          </tfoot>
        )}
      </table>
    </div>
  )
}

export default TableRenderer
