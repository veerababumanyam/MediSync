/**
 * ConfidenceIndicator Component
 *
 * Visual indicator for Council confidence scores with:
 * - 0-100% progress bar with color coding
 * - Multiple sizes (sm, md, lg)
 * - Animated transitions
 * - RTL support
 * - Accessibility support
 *
 * @module components/council/ConfidenceIndicator
 */

import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/cn'

/**
 * Props for ConfidenceIndicator component
 */
export interface ConfidenceIndicatorProps {
  /** Confidence value (0-100) */
  value: number
  /** Show percentage label */
  showLabel?: boolean
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
  /** Additional CSS classes */
  className?: string
  /** Custom label text */
  label?: string
  /** Show status text */
  showStatus?: boolean
}

/**
 * Size configurations
 */
const sizes = {
  sm: {
    bar: 'h-1.5',
    text: 'text-xs',
    container: 'gap-1.5',
  },
  md: {
    bar: 'h-2.5',
    text: 'text-sm',
    container: 'gap-2',
  },
  lg: {
    bar: 'h-4',
    text: 'text-base',
    container: 'gap-3',
  },
}

/**
 * Get color class based on confidence value
 */
function getConfidenceColor(value: number): {
  bg: string
  text: string
  status: string
} {
  if (value >= 80) {
    return {
      bg: 'bg-green-500',
      text: 'text-green-600',
      status: 'high',
    }
  }
  if (value >= 60) {
    return {
      bg: 'bg-emerald-500',
      text: 'text-emerald-600',
      status: 'good',
    }
  }
  if (value >= 40) {
    return {
      bg: 'bg-amber-500',
      text: 'text-amber-600',
      status: 'moderate',
    }
  }
  if (value >= 20) {
    return {
      bg: 'bg-orange-500',
      text: 'text-orange-600',
      status: 'low',
    }
  }
  return {
    bg: 'bg-red-500',
    text: 'text-red-600',
    status: 'very-low',
  }
}

/**
 * Confidence indicator component
 */
export function ConfidenceIndicator({
  value,
  showLabel = false,
  size = 'md',
  className,
  label,
  showStatus = false,
}: ConfidenceIndicatorProps) {
  const { t } = useTranslation('council')

  // Clamp value to 0-100 range
  const clampedValue = Math.max(0, Math.min(100, value))
  const { bg, text, status } = getConfidenceColor(clampedValue)
  const sizeConfig = sizes[size]

  return (
    <div
      className={cn('flex items-center', sizeConfig.container, className)}
      role="progressbar"
      aria-valuenow={clampedValue}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={t('confidence.label', 'Confidence score')}
    >
      {/* Label */}
      {showLabel && (
        <span className={cn('font-medium text-secondary', sizeConfig.text)}>
          {label || t('confidence.label', 'Confidence')}
        </span>
      )}

      {/* Progress Bar Container */}
      <div className="flex-1 min-w-[80px]">
        <div
          className={cn(
            'w-full rounded-full bg-surface-glass-strong overflow-hidden',
            sizeConfig.bar
          )}
        >
          {/* Progress Fill */}
          <div
            className={cn(
              'h-full rounded-full transition-all duration-500 ease-out',
              bg
            )}
            style={{ width: `${clampedValue}%` }}
          />
        </div>
      </div>

      {/* Value Display */}
      <div className="flex items-baseline gap-1 shrink-0">
        <span className={cn('font-bold tabular-nums', text, sizeConfig.text)}>
          {Math.round(clampedValue)}
        </span>
        <span className={cn('text-secondary', size === 'sm' ? 'text-xs' : 'text-sm')}>
          %
        </span>
      </div>

      {/* Status Text */}
      {showStatus && (
        <span
          className={cn(
            'px-2 py-0.5 rounded-full text-xs font-medium',
            status === 'high' && 'bg-green-100 text-green-700',
            status === 'good' && 'bg-emerald-100 text-emerald-700',
            status === 'moderate' && 'bg-amber-100 text-amber-700',
            status === 'low' && 'bg-orange-100 text-orange-700',
            status === 'very-low' && 'bg-red-100 text-red-700'
          )}
        >
          {t(`confidence.status.${status}`, status)}
        </span>
      )}
    </div>
  )
}

/**
 * Mini confidence badge for compact displays
 */
export function ConfidenceBadge({
  value,
  size = 'md',
  className,
}: {
  value: number
  size?: 'sm' | 'md' | 'lg'
  className?: string
}) {
  const clampedValue = Math.max(0, Math.min(100, value))
  const { text } = getConfidenceColor(clampedValue)
  const textSize = size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-lg' : 'text-sm'

  return (
    <span
      className={cn(
        'inline-flex items-center gap-0.5 font-semibold tabular-nums',
        textSize,
        text,
        className
      )}
    >
      {Math.round(clampedValue)}%
    </span>
  )
}

/**
 * Circular confidence indicator for alternative display
 */
export function ConfidenceCircle({
  value,
  size = 60,
  strokeWidth = 4,
  className,
  showValue = true,
}: {
  value: number
  size?: number
  strokeWidth?: number
  className?: string
  showValue?: boolean
}) {
  const clampedValue = Math.max(0, Math.min(100, value))
  const { bg, text } = getConfidenceColor(clampedValue)

  const radius = (size - strokeWidth) / 2
  const circumference = radius * 2 * Math.PI
  const offset = circumference - (clampedValue / 100) * circumference

  return (
    <div
      className={cn('relative inline-flex items-center justify-center', className)}
      style={{ width: size, height: size }}
    >
      {/* Background Circle */}
      <svg className="absolute" width={size} height={size}>
        <circle
          className="text-surface-glass-strong"
          strokeWidth={strokeWidth}
          stroke="currentColor"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
      </svg>

      {/* Progress Circle */}
      <svg
        className="absolute transform -rotate-90"
        width={size}
        height={size}
      >
        <circle
          className={cn('transition-all duration-500 ease-out', bg.replace('bg-', 'stroke-'))}
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
      </svg>

      {/* Value Text */}
      {showValue && (
        <span className={cn('font-bold tabular-nums', text, 'text-sm')}>
          {Math.round(clampedValue)}%
        </span>
      )}
    </div>
  )
}

export default ConfidenceIndicator
