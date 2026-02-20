/**
 * WidgetSettings Component
 *
 * Settings dialog for configuring pinned chart widgets:
 * - Title editing
 * - Refresh interval
 * - Position/size on dashboard grid
 * - Delete option
 *
 * @module components/dashboard/WidgetSettings
 */
import React, { useState, useCallback, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocale } from '../../hooks/useLocale'
import type { PinnedChart } from '../../services/api'

/**
 * Props for the WidgetSettings component
 */
export interface WidgetSettingsProps {
  /** The chart to configure */
  chart: PinnedChart
  /** Callback when settings are saved */
  onSave: (updates: Partial<PinnedChart>) => Promise<void>
  /** Callback when widget is deleted */
  onDelete: () => Promise<void>
  /** Callback to close the settings dialog */
  onClose: () => void
  /** Whether the dialog is open */
  isOpen: boolean
}

/**
 * Refresh interval options in minutes
 */
const REFRESH_INTERVALS = [
  { value: 0, label: 'never' },
  { value: 5, label: '5min' },
  { value: 15, label: '15min' },
  { value: 30, label: '30min' },
  { value: 60, label: '1hour' },
  { value: 360, label: '6hours' },
  { value: 1440, label: 'daily' },
]

/**
 * WidgetSettings - Configuration dialog for dashboard widgets
 *
 * Features:
 * - Title editing
 * - Refresh interval selection
 * - Delete confirmation
 * - RTL support
 * - Loading states
 */
export const WidgetSettings: React.FC<WidgetSettingsProps> = ({
  chart,
  onSave,
  onDelete,
  onClose,
  isOpen,
}) => {
  const { t } = useTranslation('dashboard')
  const { isRTL } = useLocale()

  const [title, setTitle] = useState(chart.title)
  const [refreshInterval, setRefreshInterval] = useState(chart.refreshInterval)
  const [isSaving, setIsSaving] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Reset form when chart changes
  useEffect(() => {
    setTitle(chart.title)
    setRefreshInterval(chart.refreshInterval)
    setShowDeleteConfirm(false)
    setError(null)
  }, [chart])

  /**
   * Handle save action
   */
  const handleSave = useCallback(async () => {
    if (!title.trim()) {
      setError(t('settings.titleRequired', 'Title is required'))
      return
    }

    setIsSaving(true)
    setError(null)

    try {
      await onSave({
        title: title.trim(),
        refreshInterval,
      })
      onClose()
    } catch (err) {
      console.error('Failed to save settings:', err)
      setError(t('settings.saveError', 'Failed to save settings'))
    } finally {
      setIsSaving(false)
    }
  }, [title, refreshInterval, onSave, onClose, t])

  /**
   * Handle delete action
   */
  const handleDelete = useCallback(async () => {
    setIsDeleting(true)
    setError(null)

    try {
      await onDelete()
      onClose()
    } catch (err) {
      console.error('Failed to delete widget:', err)
      setError(t('settings.deleteError', 'Failed to delete widget'))
      setShowDeleteConfirm(false)
    } finally {
      setIsDeleting(false)
    }
  }, [onDelete, onClose, t])

  /**
   * Handle keyboard events
   */
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      onClose()
    } else if (e.key === 'Enter' && !e.shiftKey) {
      handleSave()
    }
  }, [onClose, handleSave])

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      onClick={onClose}
      onKeyDown={handleKeyDown}
    >
      <div
        className={`bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-md overflow-hidden ${
          isRTL ? 'rtl' : 'ltr'
        }`}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            {t('settings.title', 'Widget Settings')}
          </h2>
          <button
            onClick={onClose}
            className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Body */}
        <div className="p-4 space-y-4">
          {/* Error display */}
          {error && (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg text-sm">
              {error}
            </div>
          )}

          {/* Title field */}
          <div>
            <label
              htmlFor="widget-title"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              {t('settings.titleLabel', 'Title')}
            </label>
            <input
              id="widget-title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder={t('settings.titlePlaceholder', 'Enter widget title')}
            />
          </div>

          {/* Refresh interval */}
          <div>
            <label
              htmlFor="refresh-interval"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              {t('settings.refreshLabel', 'Auto-refresh interval')}
            </label>
            <select
              id="refresh-interval"
              value={refreshInterval}
              onChange={(e) => setRefreshInterval(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {REFRESH_INTERVALS.map((interval) => (
                <option key={interval.value} value={interval.value}>
                  {t(`settings.refreshIntervals.${interval.label}`, interval.label)}
                </option>
              ))}
            </select>
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {t('settings.refreshHint', 'How often to automatically update this chart')}
            </p>
          </div>

          {/* Query preview (read-only) */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('settings.queryLabel', 'Query')}
            </label>
            <div className="p-3 bg-gray-50 dark:bg-gray-900 rounded-lg text-sm text-gray-600 dark:text-gray-400 font-mono overflow-x-auto">
              {chart.naturalLanguageQuery}
            </div>
          </div>

          {/* Delete section */}
          {!showDeleteConfirm ? (
            <button
              onClick={() => setShowDeleteConfirm(true)}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors text-sm"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              {t('settings.deleteWidget', 'Delete widget')}
            </button>
          ) : (
            <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg space-y-3">
              <p className="text-sm text-red-800 dark:text-red-200">
                {t('settings.deleteConfirm', 'Are you sure you want to delete this widget?')}
              </p>
              <div className="flex items-center gap-2">
                <button
                  onClick={handleDelete}
                  disabled={isDeleting}
                  className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 text-sm font-medium"
                >
                  {isDeleting ? t('deleting', 'Deleting...') : t('settings.confirmDelete', 'Yes, delete')}
                </button>
                <button
                  onClick={() => setShowDeleteConfirm(false)}
                  disabled={isDeleting}
                  className="flex-1 px-4 py-2 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 text-sm"
                >
                  {t('cancel', 'Cancel')}
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        {!showDeleteConfirm && (
          <div className="flex items-center justify-end gap-2 px-4 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-sm"
            >
              {t('cancel', 'Cancel')}
            </button>
            <button
              onClick={handleSave}
              disabled={isSaving || !title.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 text-sm font-medium"
            >
              {isSaving ? t('saving', 'Saving...') : t('save', 'Save')}
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

export default WidgetSettings
