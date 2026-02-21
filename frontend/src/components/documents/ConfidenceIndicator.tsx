/**
 * ConfidenceIndicator Component
 *
 * Visual confidence meter showing extraction confidence level.
 * Color coding: green (>=95%), yellow (70-94%), red (<70%)
 */

import { useTranslation } from 'react-i18next'

interface ConfidenceIndicatorProps {
  confidence: number // 0-1 range
  showLabel?: boolean
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function ConfidenceIndicator({
  confidence,
  showLabel = true,
  size = 'md',
  className = '',
}: ConfidenceIndicatorProps) {
  const { t, i18n } = useTranslation()
  const isRTL = i18n.language === 'ar'

  const percentage = Math.round(confidence * 100)

  const getColorClass = () => {
    if (percentage >= 95) return 'bg-green-500'
    if (percentage >= 70) return 'bg-yellow-500'
    return 'bg-red-500'
  }

  const getTextColorClass = () => {
    if (percentage >= 95) return 'text-green-600 dark:text-green-400'
    if (percentage >= 70) return 'text-yellow-600 dark:text-yellow-400'
    return 'text-red-600 dark:text-red-400'
  }

  const getStatusLabel = () => {
    if (percentage >= 95) return t('documents.confidence.high', 'High')
    if (percentage >= 70) return t('documents.confidence.medium', 'Medium')
    return t('documents.confidence.low', 'Low')
  }

  const sizeClasses = {
    sm: 'h-1.5',
    md: 'h-2',
    lg: 'h-3',
  }

  return (
    <div
      className={`confidence-indicator ${className}`}
      dir={isRTL ? 'rtl' : 'ltr'}
      title={`${percentage}% ${getStatusLabel()}`}
    >
      {/* Progress Bar */}
      <div className={`w-full bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden ${sizeClasses[size]}`}>
        <div
          className={`${sizeClasses[size]} ${getColorClass()} transition-all duration-300 rounded-full`}
          style={{ width: `${percentage}%` }}
          role="progressbar"
          aria-valuenow={percentage}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={t('documents.confidence.label', 'Confidence score')}
        />
      </div>

      {/* Label */}
      {showLabel && (
        <div className="flex justify-between items-center mt-1">
          <span className={`text-xs font-medium ${getTextColorClass()}`}>
            {getStatusLabel()}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            {percentage}%
          </span>
        </div>
      )}
    </div>
  )
}

/**
 * ConfidenceBadge - A compact badge version for table cells
 */
export function ConfidenceBadge({ confidence }: { confidence: number }) {
  const { t } = useTranslation()
  const percentage = Math.round(confidence * 100)

  const getClasses = () => {
    if (percentage >= 95) {
      return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
    }
    if (percentage >= 70) {
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
    }
    return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
  }

  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${getClasses()}`}
      title={t('documents.confidence.tooltip', '{{percent}}% confidence', { percent: percentage })}
    >
      {percentage}%
    </span>
  )
}

export default ConfidenceIndicator
