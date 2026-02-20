/**
 * LoadingSpinner Component
 *
 * Animated loading spinner with RTL support and size variants.
 * Uses Tailwind CSS for styling and CSS animations.
 *
 * @module components/common/LoadingSpinner
 */
import { useTranslation } from 'react-i18next'

export interface LoadingSpinnerProps {
  /** Spinner size variant */
  size?: 'sm' | 'md' | 'lg'
  /** Optional label text below spinner */
  label?: string
  /** Additional CSS classes */
  className?: string
}

/**
 * Size configurations for the spinner
 */
const sizeConfig = {
  sm: {
    spinner: 'w-4 h-4',
    border: 'border-2',
    text: 'text-xs mt-1',
  },
  md: {
    spinner: 'w-8 h-8',
    border: 'border-[3px]',
    text: 'text-sm mt-2',
  },
  lg: {
    spinner: 'w-12 h-12',
    border: 'border-4',
    text: 'text-base mt-3',
  },
} as const

/**
 * LoadingSpinner - Animated loading indicator
 *
 * Features:
 * - Three size variants (sm, md, lg)
 * - Optional localized label
 * - RTL-compatible animation
 * - Dark mode support
 */
export function LoadingSpinner({
  size = 'md',
  label,
  className = '',
}: LoadingSpinnerProps) {
  const { t } = useTranslation()
  const config = sizeConfig[size]

  const labelText = label || t('common.loading', 'Loading...')

  return (
    <div
      className={`inline-flex flex-col items-center justify-center ${className}`}
      role="status"
      aria-live="polite"
      aria-busy="true"
    >
      {/* Spinner element */}
      <div
        className={`
          ${config.spinner}
          ${config.border}
          border-slate-200
          border-t-blue-600
          dark:border-slate-600
          dark:border-t-blue-400
          rounded-full
          animate-spin
        `}
        style={{
          // Ensure consistent spin direction in RTL
          animationDirection: 'normal',
        }}
        aria-hidden="true"
      />

      {/* Label text */}
      {labelText && (
        <span
          className={`
            ${config.text}
            text-slate-600
            dark:text-slate-400
            font-medium
          `}
        >
          {labelText}
        </span>
      )}

      {/* Screen reader only text */}
      <span className="sr-only">{labelText}</span>
    </div>
  )
}

export default LoadingSpinner
