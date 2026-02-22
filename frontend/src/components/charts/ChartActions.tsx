/**
 * ChartActions Component
 *
 * Provides actions for chart interactions:
 * - Pin to dashboard
 * - Export (PDF, CSV, Excel)
 * - Share
 *
 * @module components/charts/ChartActions
 */
import React, { useState, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocale } from '../../hooks/useLocale'
import { dashboardApi } from '../../services/api'

/**
 * Props for the ChartActions component
 */
export interface ChartActionsProps {
  /** Chart title */
  title: string
  /** Natural language query that generated this chart */
  query: string
  /** SQL query used to generate this chart */
  sqlQuery: string
  /** Chart type (line, bar, pie, kpi-card, etc.) */
  chartType: 'bar' | 'line' | 'pie' | 'table' | 'kpi'
  /** Chart specification data */
  chartSpec: Record<string, unknown>
  /** Callback when chart is pinned successfully */
  onPinned?: () => void
  /** Callback when export starts */
  onExport?: (format: 'pdf' | 'csv' | 'excel') => void
  /** Additional CSS class name */
  className?: string
  /** Whether to show compact layout */
  compact?: boolean
}

/**
 * ChartActions - Actions toolbar for chart visualizations
 *
 * Features:
 * - Pin chart to dashboard
 * - Export to PDF/CSV/Excel
 * - RTL support
 * - Dropdown menu for export options
 */
export const ChartActions: React.FC<ChartActionsProps> = ({
  title,
  query,
  sqlQuery,
  chartType,
  chartSpec,
  onPinned,
  onExport,
  className = '',
  compact = false,
}) => {
  const { t } = useTranslation('dashboard')
  const { isRTL } = useLocale()
  const [isPinning, setIsPinning] = useState(false)
  const [isPinned, setIsPinned] = useState(false)
  const [showExportMenu, setShowExportMenu] = useState(false)
  const [error, setError] = useState<string | null>(null)

  /**
   * Handle pinning chart to dashboard
   */
  const handlePin = useCallback(async () => {
    if (isPinning || isPinned) return

    setIsPinning(true)
    setError(null)

    try {
      await dashboardApi.pinChart({
        title,
        naturalLanguageQuery: query,
        sqlQuery,
        chartType,
        chartSpec,
        refreshInterval: 0, // No auto-refresh by default
        position: { row: 0, col: 0, size: 1 },
      })

      setIsPinned(true)
      onPinned?.()
    } catch (err) {
      console.error('Failed to pin chart:', err)
      setError(t('error.pinChart', 'Failed to pin chart'))
    } finally {
      setIsPinning(false)
    }
  }, [isPinning, isPinned, title, query, sqlQuery, chartType, chartSpec, onPinned, t])

  /**
   * Handle export action
   */
  const handleExport = useCallback((format: 'pdf' | 'csv' | 'excel') => {
    onExport?.(format)
    setShowExportMenu(false)
  }, [onExport])

  /**
   * Export button click handler (default to PDF)
   */
  const handleExportClick = useCallback(() => {
    handleExport('pdf')
  }, [handleExport])

  if (compact) {
    return (
      <div className={`flex items-center gap-2 ${className}`}>
        <button
          onClick={handlePin}
          disabled={isPinning || isPinned}
          className={`glass-interactive p-2 rounded-xl transition-all ${
            isPinned
              ? 'text-blue-600 ring-2 ring-blue-400/40'
              : 'text-slate-500 hover:text-blue-600'
          }`}
          title={isPinned ? t('pinned', 'Pinned to dashboard') : t('pinToDashboard', 'Pin to dashboard')}
        >
          <svg className={`w-5 h-5 ${isPinning ? 'animate-pulse' : ''}`} fill={isPinned ? 'currentColor' : 'none'} viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
          </svg>
        </button>
        <button
          onClick={handleExportClick}
          className="glass-interactive p-2 text-slate-500 hover:text-blue-600 rounded-xl transition-all"
          title={t('export', 'Export')}
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </button>
      </div>
    )
  }

  return (
    <div className={`flex items-center gap-3 ${isRTL ? 'flex-row-reverse' : ''} ${className}`}>
      {/* Error display */}
      {error && (
        <span className="text-sm text-red-600 dark:text-red-400">{error}</span>
      )}

      {/* Pin button */}
      <button
        onClick={handlePin}
        disabled={isPinning || isPinned}
        className={`glass-interactive inline-flex items-center gap-2 px-3 py-1.5 rounded-xl text-sm font-medium transition-all ${
          isPinned
            ? 'text-blue-600 ring-2 ring-blue-400/40'
            : 'text-slate-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400'
        }`}
      >
        <svg className={`w-4 h-4 ${isPinning ? 'animate-pulse' : ''}`} fill={isPinned ? 'currentColor' : 'none'} viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
        {isPinned ? t('pinned', 'Pinned') : t('pin', 'Pin')}
      </button>

      {/* Export dropdown */}
      <div className="relative">
        <button
          onClick={() => setShowExportMenu(!showExportMenu)}
          className="glass-interactive inline-flex items-center gap-2 px-3 py-1.5 text-slate-700 dark:text-gray-300 rounded-xl text-sm font-medium transition-all"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          {t('export', 'Export')}
          <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        {showExportMenu && (
          <>
            <div
              className="fixed inset-0 z-10"
              onClick={() => setShowExportMenu(false)}
            />
            <div className={`glass-elevated absolute ${isRTL ? 'left-0' : 'right-0'} mt-1 w-40 rounded-xl z-20`}>
              <button
                onClick={() => handleExport('pdf')}
                className="glass-interactive w-full flex items-center gap-2 px-4 py-2 text-sm text-slate-700 dark:text-gray-200 rounded-t-xl"
              >
                <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 3.5L18.5 11H13V5.5zM8.5 13h7v1.5h-7V13zm0 3h7v1.5h-7V16z"/>
                </svg>
                PDF
              </button>
              <button
                onClick={() => handleExport('csv')}
                className="glass-interactive w-full flex items-center gap-2 px-4 py-2 text-sm text-slate-700 dark:text-gray-200"
              >
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 3.5L18.5 11H13V5.5z"/>
                </svg>
                CSV
              </button>
              <button
                onClick={() => handleExport('excel')}
                className="glass-interactive w-full flex items-center gap-2 px-4 py-2 text-sm text-slate-700 dark:text-gray-200 rounded-b-xl"
              >
                <svg className="w-4 h-4 text-emerald-600" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 3.5L18.5 11H13V5.5z"/>
                </svg>
                Excel
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default ChartActions
