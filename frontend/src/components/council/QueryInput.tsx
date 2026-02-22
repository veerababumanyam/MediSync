/**
 * QueryInput Component
 *
 * Input component for submitting Council deliberation queries with:
 * - Configurable consensus threshold slider
 * - Loading state during deliberation
 * - Error display
 * - RTL support
 *
 * @module components/council/QueryInput
 */

import React, { useState, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { LiquidGlassButton } from '../ui/LiquidGlassButton'
import { LiquidGlassCard } from '../ui/LiquidGlassCard'
import { cn } from '@/lib/cn'

/**
 * Props for QueryInput component
 */
export interface QueryInputProps {
  /** Handler called when query is submitted */
  onSubmit: (query: string, threshold: number) => void
  /** Whether a deliberation is in progress */
  isLoading?: boolean
  /** Error message to display */
  error?: string | null
  /** Additional CSS classes */
  className?: string
  /** Default consensus threshold (0-1) */
  defaultThreshold?: number
  /** Minimum threshold value */
  minThreshold?: number
  /** Maximum threshold value */
  maxThreshold?: number
  /** Placeholder text for input */
  placeholder?: string
  /** Disable the input */
  disabled?: boolean
}

/**
 * Query input component for Council deliberations
 */
export function QueryInput({
  onSubmit,
  isLoading = false,
  error,
  className,
  defaultThreshold = 0.80,
  minThreshold = 0.50,
  maxThreshold = 1.00,
  placeholder,
  disabled = false,
}: QueryInputProps) {
  const { t } = useTranslation('council')
  const [query, setQuery] = useState('')
  const [threshold, setThreshold] = useState(defaultThreshold)

  /**
   * Handle form submission
   */
  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault()
      if (query.trim() && !isLoading && !disabled) {
        onSubmit(query.trim(), threshold)
      }
    },
    [query, threshold, isLoading, disabled, onSubmit]
  )

  /**
   * Handle threshold slider change
   */
  const handleThresholdChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setThreshold(parseFloat(e.target.value))
    },
    []
  )

  /**
   * Clear the input
   */
  const handleClear = useCallback(() => {
    setQuery('')
  }, [])

  const isValid = query.trim().length > 0

  return (
    <LiquidGlassCard className={cn('p-6', className)}>
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Query Input */}
        <div className="space-y-2">
          <label
            htmlFor="council-query"
            className="block text-sm font-medium text-primary"
          >
            {t('queryInput.label', 'Your Question')}
          </label>
          <textarea
            id="council-query"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={
              placeholder ||
              t('queryInput.placeholder', 'Ask a medical question...')
            }
            disabled={isLoading || disabled}
            rows={3}
            className={cn(
              'w-full px-4 py-3 rounded-xl',
              'bg-surface-glass/50 backdrop-blur-sm',
              'border border-glass focus:border-primary focus:ring-1 focus:ring-primary',
              'text-primary placeholder:text-secondary/60',
              'transition-all duration-200',
              'resize-none',
              'disabled:opacity-50 disabled:cursor-not-allowed',
              error && 'border-red-500 focus:border-red-500 focus:ring-red-500'
            )}
          />
        </div>

        {/* Threshold Slider */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label
              htmlFor="consensus-threshold"
              className="text-sm font-medium text-primary"
            >
              {t('queryInput.threshold', 'Consensus Threshold')}
            </label>
            <span className="text-sm font-semibold text-primary">
              {Math.round(threshold * 100)}%
            </span>
          </div>
          <input
            id="consensus-threshold"
            type="range"
            min={minThreshold}
            max={maxThreshold}
            step={0.05}
            value={threshold}
            onChange={handleThresholdChange}
            disabled={isLoading || disabled}
            className={cn(
              'w-full h-2 rounded-full appearance-none cursor-pointer',
              'bg-surface-glass-strong',
              '[&::-webkit-slider-thumb]:appearance-none',
              '[&::-webkit-slider-thumb]:w-5 [&::-webkit-slider-thumb]:h-5',
              '[&::-webkit-slider-thumb]:rounded-full',
              '[&::-webkit-slider-thumb]:bg-primary',
              '[&::-webkit-slider-thumb]:shadow-lg',
              '[&::-webkit-slider-thumb]:transition-transform',
              '[&::-webkit-slider-thumb]:hover:scale-110',
              'disabled:opacity-50 disabled:cursor-not-allowed'
            )}
          />
          <div className="flex justify-between text-xs text-secondary">
            <span>{Math.round(minThreshold * 100)}%</span>
            <span>{t('queryInput.thresholdHint', 'Higher = more agreement required')}</span>
            <span>{Math.round(maxThreshold * 100)}%</span>
          </div>
        </div>

        {/* Error Display */}
        {error && (
          <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/30">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex items-center justify-end gap-3">
          {query && (
            <LiquidGlassButton
              type="button"
              variant="ghost"
              onClick={handleClear}
              disabled={isLoading || disabled}
            >
              {t('queryInput.clear', 'Clear')}
            </LiquidGlassButton>
          )}
          <LiquidGlassButton
            type="submit"
            variant="primary"
            disabled={!isValid || isLoading || disabled}
            className="min-w-[120px]"
          >
            {isLoading ? (
              <span className="flex items-center gap-2">
                <svg
                  className="animate-spin h-4 w-4"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                {t('queryInput.processing', 'Processing...')}
              </span>
            ) : (
              t('queryInput.submit', 'Ask Council')
            )}
          </LiquidGlassButton>
        </div>
      </form>
    </LiquidGlassCard>
  )
}

export default QueryInput
