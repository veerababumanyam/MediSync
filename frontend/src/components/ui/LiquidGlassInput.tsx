/**
 * Liquid Glass Input Component
 *
 * Premium iOS-inspired glassmorphic input field with liquid animations,
 * dynamic focus states, and WCAG 3.0 Bronze compliance.
 *
 * Features:
 * - Multi-layered glass effect with specular highlights
 * - Liquid focus animations with glow effect
 * - Branded color variants using logo colors
 * - Error and success states
 * - Icon support (prefix, suffix)
 * - Character counter
 * - Reduced motion support
 * - RTL support
 *
 * @module components/ui/LiquidGlassInput
 * @version 2.0.0
 */

import React, { forwardRef, useState, useCallback, useRef, useEffect, useId } from 'react'
import { motion, type HTMLMotionProps } from 'framer-motion'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

// Type for motion input props - compatible with framer-motion v12
type MotionInputProps = Omit<HTMLMotionProps<'input'>, 'size' | 'disabled' | 'variants' | 'onAnimationStart'>

/**
 * Liquid glass input variant definitions
 */
const liquidInputVariants = cva(
  // Base classes
  'liquid-glass-input w-full transition-all duration-200',
  {
    variants: {
      // Size variants
      size: {
        sm: 'px-3 py-2 text-sm',
        md: 'px-4 py-2.5 text-base',
        lg: 'px-5 py-3 text-lg',
      },
      // State variants
      state: {
        default: '',
        error: 'border-red-400 focus:border-red-500 focus:shadow-[0_0_0_4px_rgba(239,68,68,0.15)]',
        success: 'border-emerald-400 focus:border-emerald-500 focus:shadow-[0_0_0_4px_rgba(16,185,129,0.15)]',
        warning: 'border-amber-400 focus:border-amber-500 focus:shadow-[0_0_0_4px_rgba(245,158,11,0.15)]',
      },
      // Border radius
      radius: {
        sm: 'rounded-md',
        md: 'rounded-lg',
        lg: 'rounded-xl',
        full: 'rounded-full',
      },
    },
    defaultVariants: {
      size: 'md',
      state: 'default',
      radius: 'md',
    },
  }
)

/**
 * Props for LiquidGlassInput component
 */
export interface LiquidGlassInputProps
  extends MotionInputProps,
  VariantProps<typeof liquidInputVariants> {
  /** Label text */
  label?: string
  /** Error message */
  error?: string
  /** Helper text */
  helperText?: string
  /** Success message */
  success?: string
  /** Warning message */
  warning?: string
  /** Icon to display before input */
  prefixIcon?: React.ReactNode
  /** Icon to display after input */
  suffixIcon?: React.ReactNode
  /** Whether to show character count */
  showCount?: boolean
  /** Maximum character length */
  maxLength?: number
  /** Whether input is disabled */
  disabled?: boolean
  /** Whether input is loading (shows spinner) */
  isLoading?: boolean
  /** Container className */
  containerClassName?: string
}

/**
 * Liquid Glass Input Component
 *
 * A premium glassmorphic input field with liquid animations.
 *
 * @example
 * ```tsx
 * // Basic input
 * <LiquidGlassInput placeholder="Enter text..." />
 *
 * // With label and error
 * <LiquidGlassInput
 *   label="Email"
 *   error="Invalid email format"
 *   placeholder="your@email.com"
 * />
 *
 * // With icons
 * <LiquidGlassInput
 *   prefixIcon={<SearchIcon />}
 *   placeholder="Search..."
 * />
 *
 * // With character count
 * <LiquidGlassInput
 *   maxLength={200}
 *   showCount
 *   placeholder="Enter description..."
 * />
 * ```
 */
export const LiquidGlassInput = forwardRef<HTMLInputElement, LiquidGlassInputProps>(
  (
    {
      className,
      containerClassName,
      size,
      state,
      radius,
      label,
      error,
      helperText,
      success,
      warning,
      prefixIcon,
      suffixIcon,
      showCount = false,
      maxLength,
      disabled = false,
      isLoading = false,
      value,
      onChange,
      ...props
    },
    ref
  ) => {
    const [isFocused, setIsFocused] = useState(false)
    const [internalValue, setInternalValue] = useState('')
    const inputRef = useRef<HTMLInputElement>(null)

    // WCAG 3.0 Bronze: Generate stable IDs for aria-describedby linkage
    const reactId = useId()
    const inputId = props.id || `lg-input-${reactId}`
    const errorId = `${inputId}-error`
    const helperId = `${inputId}-helper`

    // Handle ref forwarding
    useEffect(() => {
      if (ref) {
        if (typeof ref === 'function') {
          ref(inputRef.current)
        } else {
          ref.current = inputRef.current
        }
      }
    }, [ref])

    // Determine current value (controlled or uncontrolled)
    const currentValue = value !== undefined ? value : internalValue
    const currentLength = String(currentValue || '').length

    // Determine state based on props
    const inputState = error ? 'error' : success ? 'success' : warning ? 'warning' : state

    // Handle change
    const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
      if (value === undefined) {
        setInternalValue(e.target.value)
      }
      onChange?.(e)
    }, [onChange, value])

    // Build class names
    const inputClasses = cn(
      liquidInputVariants({ size, state: inputState, radius }),
      isFocused && 'scale-[1.01]',
      prefixIcon && 'pl-11',
      suffixIcon && 'pr-11',
      isLoading && 'pr-11',
      disabled && 'opacity-60 cursor-not-allowed',
      className
    )

    // Label classes for WCAG 3.0 Bronze contrast
    const labelClasses = cn(
      'block text-sm font-medium mb-1.5 transition-colors',
      disabled ? 'text-slate-400 dark:text-slate-500' : 'text-slate-700 dark:text-slate-300',
      inputState === 'error' && 'text-red-600 dark:text-red-400',
      inputState === 'success' && 'text-emerald-600 dark:text-emerald-400',
      inputState === 'warning' && 'text-amber-600 dark:text-amber-400'
    )

    // Helper text classes
    const helperTextClasses = cn(
      'mt-1.5 text-sm transition-colors',
      inputState === 'error' && 'text-red-600 dark:text-red-400',
      inputState === 'success' && 'text-emerald-600 dark:text-emerald-400',
      inputState === 'warning' && 'text-amber-600 dark:text-amber-400',
      inputState === 'default' && 'text-slate-500 dark:text-slate-400'
    )

    // Character count classes
    const countClasses = cn(
      'text-xs mt-1.5 text-end transition-colors',
      maxLength && currentLength > maxLength * 0.9
        ? 'text-red-500 dark:text-red-400'
        : 'text-slate-400 dark:text-slate-500'
    )

    return (
      <div className={containerClassName}>
        {label && (
          <label className={labelClasses} htmlFor={inputId}>
            {label}
            {props.required && <span className="text-red-500 ml-0.5" aria-hidden="true">*</span>}
          </label>
        )}

        <div className="relative">
          {/* Prefix Icon */}
          {prefixIcon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-500 pointer-events-none">
              {prefixIcon}
            </div>
          )}

          {/* Input Field */}
          <motion.input
            ref={inputRef}
            id={inputId}
            className={inputClasses}
            value={currentValue}
            onChange={handleChange}
            onFocus={(e) => {
              setIsFocused(true)
              props.onFocus?.(e)
            }}
            onBlur={(e) => {
              setIsFocused(false)
              props.onBlur?.(e)
            }}
            disabled={disabled || isLoading}
            maxLength={maxLength}
            aria-invalid={inputState === 'error' ? true : undefined}
            aria-describedby={
              error ? errorId : helperText ? helperId : undefined
            }
            aria-required={props.required || undefined}
            animate={{
              borderColor: isFocused ? 'var(--ms-teal)' : undefined,
            }}
            transition={{ duration: 0.2 }}
            {...props}
          />

          {/* Suffix Icon or Loading Spinner */}
          {(suffixIcon || isLoading) && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-500 pointer-events-none">
              {isLoading ? (
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
              ) : (
                suffixIcon
              )}
            </div>
          )}

          {/* Focus ring effect */}
          {isFocused && (
            <motion.div
              className="absolute inset-0 rounded-lg pointer-events-none"
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1.02 }}
              exit={{ opacity: 0, scale: 0.95 }}
              style={{
                boxShadow: inputState === 'error'
                  ? '0 0 0 4px rgba(239, 68, 68, 0.15)'
                  : inputState === 'success'
                    ? '0 0 0 4px rgba(16, 185, 129, 0.15)'
                    : inputState === 'warning'
                      ? '0 0 0 4px rgba(245, 158, 11, 0.15)'
                      : '0 0 0 4px rgba(0, 191, 165, 0.15)',
              }}
              transition={{ duration: 0.2 }}
            />
          )}
        </div>

        {/* Helper text, error, success, or warning message — WCAG 3.0 SC 3.3.1 */}
        {(error || success || warning || helperText) && (
          <p
            id={error ? errorId : helperId}
            className={helperTextClasses}
            role={error ? 'alert' : undefined}
            aria-live={error ? 'assertive' : undefined}
          >
            {error && <span aria-hidden="true">⚠ </span>}
            {error || success || warning || helperText}
          </p>
        )}

        {/* Character count */}
        {showCount && maxLength && (
          <p className={countClasses}>
            {currentLength} / {maxLength}
          </p>
        )}
      </div>
    )
  }
)

LiquidGlassInput.displayName = 'LiquidGlassInput'

/**
 * Textarea variant of Liquid Glass Input
 */
type MotionTextareaProps = Omit<HTMLMotionProps<'textarea'>, 'size' | 'disabled' | 'variants' | 'onAnimationStart'>

export interface LiquidGlassTextareaProps
  extends MotionTextareaProps,
  VariantProps<typeof liquidInputVariants> {
  rows?: number
  maxLength?: number
  showCount?: boolean
  label?: string
  error?: string
  helperText?: string
  success?: string
  warning?: string
  disabled?: boolean
  isLoading?: boolean
  containerClassName?: string
}

export const LiquidGlassTextarea = forwardRef<HTMLTextAreaElement, LiquidGlassTextareaProps>(
  (
    {
      className,
      size = 'md',
      state = 'default',
      radius = 'md',
      label,
      error,
      helperText,
      success,
      warning,
      disabled = false,
      isLoading = false,
      rows = 4,
      maxLength,
      showCount = false,
      value,
      onChange,
      ...props
    },
    ref
  ) => {
    const [isFocused, setIsFocused] = useState(false)
    const [internalValue, setInternalValue] = useState('')
    const textareaRef = useRef<HTMLTextAreaElement>(null)

    // Handle ref forwarding
    useEffect(() => {
      if (ref) {
        if (typeof ref === 'function') {
          ref(textareaRef.current)
        } else {
          ref.current = textareaRef.current
        }
      }
    }, [ref])

    // Determine current value
    const currentValue = value !== undefined ? value : internalValue
    const currentLength = String(currentValue || '').length

    // Determine state
    const inputState = error ? 'error' : success ? 'success' : warning ? 'warning' : state

    // Handle change
    const handleChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
      if (value === undefined) {
        setInternalValue(e.target.value)
      }
      onChange?.(e)
    }, [onChange, value])

    // Build classes
    const textareaClasses = cn(
      'liquid-glass-input w-full resize-none transition-all duration-200',
      size === 'sm' && 'px-3 py-2 text-sm',
      size === 'md' && 'px-4 py-2.5 text-base',
      size === 'lg' && 'px-5 py-3 text-lg',
      radius === 'sm' && 'rounded-md',
      radius === 'md' && 'rounded-lg',
      radius === 'lg' && 'rounded-xl',
      inputState === 'error' && 'border-red-400 focus:border-red-500',
      inputState === 'success' && 'border-emerald-400 focus:border-emerald-500',
      inputState === 'warning' && 'border-amber-400 focus:border-amber-500',
      isFocused && 'scale-[1.005]',
      isLoading && 'pr-11',
      disabled && 'opacity-60 cursor-not-allowed',
      className
    )

    return (
      <div className={props.containerClassName}>
        {label && (
          <label className={cn(
            'block text-sm font-medium mb-1.5',
            disabled ? 'text-slate-400' : 'text-slate-700 dark:text-slate-300',
            inputState === 'error' && 'text-red-600 dark:text-red-400'
          )} htmlFor={props.id}>
            {label}
          </label>
        )}

        <div className="relative">
          <motion.textarea
            ref={textareaRef}
            className={textareaClasses}
            value={currentValue}
            onChange={handleChange}
            onFocus={(e) => {
              setIsFocused(true)
              props.onFocus?.(e)
            }}
            onBlur={(e) => {
              setIsFocused(false)
              props.onBlur?.(e)
            }}
            disabled={disabled || isLoading}
            rows={rows}
            maxLength={maxLength}
            {...props}
          />

          {/* Loading spinner */}
          {isLoading && (
            <div className="absolute right-3 top-3 text-slate-400 dark:text-slate-500 pointer-events-none">
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
            </div>
          )}
        </div>

        {/* Messages */}
        {(error || success || warning || helperText) && (
          <p className={cn(
            'mt-1.5 text-sm',
            inputState === 'error' && 'text-red-600 dark:text-red-400',
            inputState === 'success' && 'text-emerald-600 dark:text-emerald-400',
            inputState === 'warning' && 'text-amber-600 dark:text-amber-400',
            inputState === 'default' && 'text-slate-500 dark:text-slate-400'
          )}>
            {error || success || warning || helperText}
          </p>
        )}

        {/* Character count */}
        {showCount && maxLength && (
          <p className={cn(
            'text-xs mt-1.5 text-end',
            currentLength > maxLength * 0.9
              ? 'text-red-500 dark:text-red-400'
              : 'text-slate-400 dark:text-slate-500'
          )}>
            {currentLength} / {maxLength}
          </p>
        )}
      </div>
    )
  }
)

LiquidGlassTextarea.displayName = 'LiquidGlassTextarea'

/**
 * Search input with glass styling
 */
export interface LiquidGlassSearchProps extends Omit<LiquidGlassInputProps, 'prefixIcon'> {
  onClear?: () => void
  showClearButton?: boolean
}

export const LiquidGlassSearch = forwardRef<HTMLInputElement, LiquidGlassSearchProps>(
  ({ onClear, showClearButton = true, value, ...props }, ref) => {
    const hasValue = value !== undefined ? String(value).length > 0 : false

    return (
      <LiquidGlassInput
        ref={ref}
        prefixIcon={
          <svg
            className="w-4 h-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        }
        suffixIcon={
          showClearButton && hasValue && onClear ? (
            <button
              type="button"
              onClick={onClear}
              className="hover:text-slate-600 dark:hover:text-slate-300 transition-colors"
              tabIndex={-1}
              aria-label="Clear search"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          ) : undefined
        }
        value={value}
        {...props}
      />
    )
  }
)

LiquidGlassSearch.displayName = 'LiquidGlassSearch'

export default LiquidGlassInput
