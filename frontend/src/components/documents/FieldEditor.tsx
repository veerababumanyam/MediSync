/**
 * FieldEditor Component
 *
 * Edit and verify individual extracted fields from documents.
 * Supports inline editing, validation, and verification status display.
 */

import { useState, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { documentApi } from '../../services/documents'
import { ConfidenceBadge } from './ConfidenceIndicator'
import type { ExtractedField, FieldType } from '../../types/documents'

interface FieldEditorProps {
  field: ExtractedField
  documentId: string
  isLocked: boolean
  onFieldUpdate?: (field: ExtractedField) => void
  className?: string
}

export function FieldEditor({
  field,
  documentId,
  isLocked,
  onFieldUpdate,
  className = '',
}: FieldEditorProps) {
  const { t, i18n } = useTranslation()
  const isRTL = i18n.language === 'ar'

  const [isEditing, setIsEditing] = useState(false)
  const [value, setValue] = useState(field.extractedValue)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const validateValue = useCallback(
    (val: string, type: FieldType): string | null => {
      if (!val.trim()) {
        return t('documents.field.required', 'Value is required')
      }

      switch (type) {
        case 'currency':
          if (isNaN(parseFloat(val)) || parseFloat(val) < 0) {
            return t('documents.field.invalidCurrency', 'Invalid currency amount')
          }
          break
        case 'number':
        case 'percentage':
          if (isNaN(parseFloat(val))) {
            return t('documents.field.invalidNumber', 'Invalid number')
          }
          break
        case 'date':
          if (isNaN(Date.parse(val))) {
            return t('documents.field.invalidDate', 'Invalid date format')
          }
          break
        case 'tax_id':
          if (!/^\d{2,15}$/.test(val.replace(/[-\s]/g, ''))) {
            return t('documents.field.invalidTaxId', 'Invalid tax ID format')
          }
          break
      }

      return null
    },
    [t]
  )

  const handleSave = useCallback(async () => {
    const validationError = validateValue(value, field.fieldType)
    if (validationError) {
      setError(validationError)
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      const updatedField = await documentApi.updateField(documentId, field.id, { value })
      setIsEditing(false)
      onFieldUpdate?.(updatedField)
    } catch {
      setError(t('documents.field.saveError', 'Failed to save changes'))
    } finally {
      setIsLoading(false)
    }
  }, [value, field.fieldType, field.id, documentId, validateValue, onFieldUpdate, t])

  const handleVerify = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const updatedField = await documentApi.verifyField(documentId, field.id)
      onFieldUpdate?.(updatedField)
    } catch {
      setError(t('documents.field.verifyError', 'Failed to verify field'))
    } finally {
      setIsLoading(false)
    }
  }, [documentId, field.id, onFieldUpdate, t])

  const handleCancel = useCallback(() => {
    setValue(field.extractedValue)
    setIsEditing(false)
    setError(null)
  }, [field.extractedValue])

  const getStatusClasses = () => {
    switch (field.verificationStatus) {
      case 'auto_accepted':
        return 'border-green-300 dark:border-green-700 bg-green-50 dark:bg-green-900/20'
      case 'manually_verified':
      case 'manually_corrected':
        return 'border-blue-300 dark:border-blue-700 bg-blue-50 dark:bg-blue-900/20'
      case 'high_priority':
        return 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-900/20'
      case 'needs_review':
        return 'border-yellow-300 dark:border-yellow-700 bg-yellow-50 dark:bg-yellow-900/20'
      default:
        return 'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800'
    }
  }

  const getStatusLabel = () => {
    switch (field.verificationStatus) {
      case 'auto_accepted':
        return t('documents.status.autoAccepted', 'Auto Accepted')
      case 'manually_verified':
        return t('documents.status.manuallyVerified', 'Manually Verified')
      case 'manually_corrected':
        return t('documents.status.manuallyCorrected', 'Manually Corrected')
      case 'high_priority':
        return t('documents.status.highPriority', 'High Priority')
      case 'needs_review':
        return t('documents.status.needsReview', 'Needs Review')
      case 'rejected':
        return t('documents.status.rejected', 'Rejected')
      default:
        return t('documents.status.pending', 'Pending')
    }
  }

  const getInputType = (): string => {
    switch (field.fieldType) {
      case 'currency':
      case 'number':
      case 'percentage':
        return 'number'
      case 'date':
        return 'date'
      default:
        return 'text'
    }
  }

  return (
    <div
      className={`field-editor border rounded-lg p-4 ${getStatusClasses()} ${className}`}
      dir={isRTL ? 'rtl' : 'ltr'}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-2">
        <div>
          <h4 className="text-sm font-medium text-gray-900 dark:text-gray-100">
            {field.fieldName}
          </h4>
          <div className="flex items-center gap-2 mt-1">
            <ConfidenceBadge confidence={field.confidenceScore} />
            {field.isHandwritten && (
              <span className="text-xs text-purple-600 dark:text-purple-400">
                {t('documents.field.handwritten', 'Handwritten')}
              </span>
            )}
          </div>
        </div>
        <span
          className={`text-xs px-2 py-0.5 rounded ${
            field.verificationStatus === 'high_priority'
              ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
              : field.verificationStatus === 'needs_review'
              ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
              : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400'
          }`}
        >
          {getStatusLabel()}
        </span>
      </div>

      {/* Value / Edit Input */}
      <div className="mt-2">
        {isEditing ? (
          <div>
            <input
              type={getInputType()}
              value={value}
              onChange={(e) => setValue(e.target.value)}
              className={`
                w-full px-3 py-2 border rounded-lg
                bg-white dark:bg-gray-900
                text-gray-900 dark:text-gray-100
                border-gray-300 dark:border-gray-600
                focus:ring-2 focus:ring-blue-500 focus:border-transparent
                ${isRTL ? 'text-right' : 'text-left'}
              `}
              disabled={isLoading}
              autoFocus
            />

            {/* Edit Actions */}
            <div className="flex justify-end gap-2 mt-2">
              <button
                onClick={handleCancel}
                disabled={isLoading}
                className="px-3 py-1.5 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
              >
                {t('common.cancel', 'Cancel')}
              </button>
              <button
                onClick={handleSave}
                disabled={isLoading}
                className="px-3 py-1.5 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
              >
                {isLoading ? t('common.saving', 'Saving...') : t('common.save', 'Save')}
              </button>
            </div>
          </div>
        ) : (
          <div className="flex items-center justify-between">
            <p className="text-gray-900 dark:text-gray-100">
              {field.extractedValue}
              {field.wasEdited && (
                <span className="ml-2 text-xs text-blue-500">
                  ({t('documents.field.edited', 'edited')})
                </span>
              )}
            </p>

            {/* View Actions */}
            {!isLocked && (
              <div className="flex gap-2">
                <button
                  onClick={() => setIsEditing(true)}
                  disabled={isLoading}
                  className="px-2 py-1 text-xs text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400"
                  title={t('documents.field.edit', 'Edit field')}
                >
                  {t('common.edit', 'Edit')}
                </button>
                <button
                  onClick={handleVerify}
                  disabled={isLoading}
                  className="px-2 py-1 text-xs text-gray-600 dark:text-gray-400 hover:text-green-600 dark:hover:text-green-400"
                  title={t('documents.field.verify', 'Verify field')}
                >
                  {t('documents.field.verifyBtn', 'Verify')}
                </button>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <p className="mt-2 text-sm text-red-500">{error}</p>
      )}

      {/* Original Value (if edited) */}
      {field.wasEdited && field.originalValue && !isEditing && (
        <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
          {t('documents.field.original', 'Original')}: {field.originalValue}
        </p>
      )}
    </div>
  )
}

export default FieldEditor
