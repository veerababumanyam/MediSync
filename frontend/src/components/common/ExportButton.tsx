/**
 * ExportButton Component
 *
 * Dropdown button for exporting data in multiple formats (PDF, XLSX, CSV).
 * Handles API calls, loading states, and file downloads.
 *
 * @module components/common/ExportButton
 */
import { useCallback, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { LoadingSpinner } from './LoadingSpinner'

export interface ExportData {
  /** The query string used to generate the data */
  query: string
  /** The result data to export */
  results: Record<string, unknown>[]
  /** Optional chart configuration for PDF exports */
  chartConfig?: Record<string, unknown>
}

export interface ExportButtonProps {
  /** Data to export */
  data: ExportData
  /** Base filename for downloaded file (without extension) */
  filename?: string
  /** Additional CSS classes */
  className?: string
  /** Disable the export button */
  disabled?: boolean
  /** Callback when export starts */
  onExportStart?: () => void
  /** Callback when export completes or fails */
  onExportEnd?: (success: boolean, error?: string) => void
}

/**
 * Export format configuration
 */
const exportFormats = [
  {
    id: 'pdf',
    extension: 'pdf',
    mimeType: 'application/pdf',
  },
  {
    id: 'xlsx',
    extension: 'xlsx',
    mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  },
  {
    id: 'csv',
    extension: 'csv',
    mimeType: 'text/csv',
  },
] as const

type ExportFormat = typeof exportFormats[number]['id']

/**
 * ExportButton - Export data to PDF, XLSX, or CSV
 *
 * Features:
 * - Dropdown with format options
 * - Loading state during export
 * - Automatic file download
 * - RTL-aware dropdown positioning
 * - Accessible keyboard navigation
 */
export function ExportButton({
  data,
  filename = 'export',
  className = '',
  disabled = false,
  onExportStart,
  onExportEnd,
}: ExportButtonProps) {
  const { t, i18n } = useTranslation()
  const [isOpen, setIsOpen] = useState(false)
  const [isExporting, setIsExporting] = useState(false)
  const [exportingFormat, setExportingFormat] = useState<ExportFormat | null>(null)
  const [error, setError] = useState<string | null>(null)

  const dropdownRef = useRef<HTMLDivElement>(null)
  const buttonRef = useRef<HTMLButtonElement>(null)

  const isRTL = i18n.language === 'ar'

  // Close dropdown on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  // Close dropdown on Escape key
  useEffect(() => {
    function handleEscape(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
        buttonRef.current?.focus()
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
    }

    return () => {
      document.removeEventListener('keydown', handleEscape)
    }
  }, [isOpen])

  const handleExport = useCallback(
    async (format: ExportFormat) => {
      if (isExporting || disabled) return

      setIsExporting(true)
      setExportingFormat(format)
      setError(null)
      setIsOpen(false)
      onExportStart?.()

      try {
        const response = await fetch(`/api/v1/chat/export/${format}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            ...data,
            filename,
            locale: i18n.language,
          }),
        })

        if (!response.ok) {
          throw new Error(
            `Export failed: ${response.status} ${response.statusText}`
          )
        }

        // Get blob from response
        const blob = await response.blob()

        // Create download link
        const url = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `${filename}.${format}`

        // Trigger download
        document.body.appendChild(link)
        link.click()

        // Cleanup
        document.body.removeChild(link)
        URL.revokeObjectURL(url)

        onExportEnd?.(true)
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'Export failed'
        console.error('Export error:', err)
        setError(errorMessage)
        onExportEnd?.(false, errorMessage)
      } finally {
        setIsExporting(false)
        setExportingFormat(null)
      }
    },
    [data, filename, disabled, isExporting, i18n.language, onExportStart, onExportEnd]
  )

  const toggleDropdown = useCallback(() => {
    if (!isExporting) {
      setIsOpen((prev) => !prev)
    }
  }, [isExporting])

  return (
    <div ref={dropdownRef} className={`relative inline-block ${className}`}>
      {/* Main button */}
      <button
        ref={buttonRef}
        type="button"
        onClick={toggleDropdown}
        disabled={disabled || isExporting}
        className={`
          inline-flex
          items-center
          justify-center
          gap-2
          px-4
          py-2
          rounded-lg
          text-sm
          font-medium
          transition-colors
          duration-200
          focus:outline-none
          focus-visible:ring-2
          focus-visible:ring-blue-500
          focus-visible:ring-offset-2
          disabled:opacity-50
          disabled:cursor-not-allowed
          ${
            isOpen
              ? 'bg-blue-600 text-white'
              : 'bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300'
          }
        `}
        aria-expanded={isOpen}
        aria-haspopup="listbox"
        aria-label={t('common.export.label', 'Export data')}
        aria-busy={isExporting}
      >
        {isExporting ? (
          <>
            <LoadingSpinner size="sm" className="!flex" />
            <span>
              {t('common.export.exporting', 'Exporting {{format}}...', {
                format: exportingFormat?.toUpperCase(),
              })}
            </span>
          </>
        ) : (
          <>
            {/* Download icon */}
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
              />
            </svg>
            <span>{t('common.export.label', 'Export')}</span>
            {/* Dropdown arrow */}
            <svg
              className={`w-4 h-4 transition-transform duration-200 ${
                isOpen ? 'rotate-180' : ''
              }`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </>
        )}
      </button>

      {/* Dropdown menu */}
      {isOpen && (
        <div
          className={`
            absolute
            z-50
            mt-2
            w-48
            rounded-lg
            bg-white
            dark:bg-slate-800
            shadow-lg
            border
            border-slate-200
            dark:border-slate-700
            py-1
          `}
          style={{
            // RTL-aware positioning
            [isRTL ? 'left' : 'right']: 0,
            [isRTL ? 'right' : 'left']: 'auto',
          }}
          role="listbox"
          aria-label={t('common.export.formatOptions', 'Export format options')}
        >
          {exportFormats.map((format) => (
            <button
              key={format.id}
              type="button"
              onClick={() => handleExport(format.id)}
              className={`
                w-full
                px-4
                py-2
                text-start
                text-sm
                text-slate-700
                dark:text-slate-300
                hover:bg-slate-100
                dark:hover:bg-slate-700
                transition-colors
                duration-150
                flex
                items-center
                gap-3
              `}
              role="option"
              aria-selected={false}
            >
              {/* Format icon */}
              <span className="w-6 text-center font-mono text-xs font-bold uppercase">
                {format.extension}
              </span>
              <span>
                {t(`common.export.format.${format.id}`, format.extension.toUpperCase())}
              </span>
            </button>
          ))}
        </div>
      )}

      {/* Error message */}
      {error && (
        <div
          className="absolute top-full mt-2 start-0 end-0 text-xs text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 rounded px-2 py-1"
          role="alert"
        >
          {error}
        </div>
      )}
    </div>
  )
}

export default ExportButton
