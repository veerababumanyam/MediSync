/**
 * LoadingSkeleton Component
 *
 * A glassmorphic loading skeleton with shimmer animation.
 * Used to indicate content is being loaded while maintaining visual consistency.
 *
 * Features:
 * - Glassmorphic background with blur
 * - Shimmer animation using framer-motion
 * - Multiple variants (default, card, chart, text)
 * - Respects prefers-reduced-motion
 * - WCAG 2.2 AA accessible
 *
 * @module components/ui/LoadingSkeleton
 */
import React from 'react'
import { motion } from 'framer-motion'
import { cn } from '@/lib/cn'

/**
 * Animation variants for shimmer effect
 */
const shimmerVariants = {
  initial: { x: '-100%' },
  animate: { x: '100%' },
}

/**
 * Props for LoadingSkeleton component
 */
export interface LoadingSkeletonProps {
  /** Additional CSS classes */
  className?: string
  /** Skeleton variant for different use cases */
  variant?: 'default' | 'card' | 'chart' | 'text' | 'circle'
  /** Custom height */
  height?: string
  /** Custom width */
  width?: string
  /** Whether animation is disabled */
  disableAnimation?: boolean
}

/**
 * LoadingSkeleton Component
 *
 * A placeholder skeleton that shows while content is loading.
 *
 * @example
 * ```tsx
 * // Default text skeleton
 * <LoadingSkeleton />
 *
 * // Card placeholder
 * <LoadingSkeleton variant="card" />
 *
 * // Chart placeholder
 * <LoadingSkeleton variant="chart" />
 *
 * // Custom size
 * <LoadingSkeleton width="100px" height="20px" />
 * ```
 */
export const LoadingSkeleton = React.forwardRef<HTMLDivElement, LoadingSkeletonProps>(
  (
    {
      className,
      variant = 'default',
      height,
      width,
      disableAnimation = false,
    },
    ref
  ) => {
    // Check for reduced motion preference
    const [prefersReducedMotion, setPrefersReducedMotion] = React.useState(false)

    React.useEffect(() => {
      const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
      setPrefersReducedMotion(mediaQuery.matches)

      const handleChange = () => setPrefersReducedMotion(mediaQuery.matches)
      mediaQuery.addEventListener('change', handleChange)

      return () => mediaQuery.removeEventListener('change', handleChange)
    }, [])

    const shouldAnimate = !disableAnimation && !prefersReducedMotion

    // Base classes for glass effect
    const baseClasses = cn(
      'relative overflow-hidden',
      'bg-slate-200/50 dark:bg-slate-700/50',
      'backdrop-blur-sm',
      'rounded',
      className
    )

    // Variant-specific classes
    const variantClasses = {
      default: 'h-4',
      card: 'p-6 rounded-2xl',
      chart: 'h-64 rounded-2xl',
      text: 'h-4 w-3/4',
      circle: 'rounded-full',
    }

    return (
      <div
        ref={ref}
        className={cn(baseClasses, variantClasses[variant])}
        style={{ height, width }}
        role="status"
        aria-label="Loading content"
      >
        {/* Shimmer animation overlay */}
        {shouldAnimate && (
          <motion.div
            className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent"
            variants={shimmerVariants}
            initial="initial"
            animate="animate"
            transition={{
              repeat: Infinity,
              duration: 1.5,
              ease: 'linear',
            }}
            aria-hidden="true"
          />
        )}

        {/* Screen reader only text */}
        <span className="sr-only">Loading...</span>
      </div>
    )
  }
)

LoadingSkeleton.displayName = 'LoadingSkeleton'

export default LoadingSkeleton
