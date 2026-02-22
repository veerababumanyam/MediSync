/**
 * FadeIn Animation Component
 *
 * A wrapper component that adds a fade-in with vertical slide animation
 * to its children. Respects prefers-reduced-motion for accessibility.
 *
 * Features:
 * - Configurable delay for staggered animations
 * - Respects prefers-reduced-motion (WCAG 3.0)
 * - Customizable duration and easing
 * - Works with any React element
 *
 * @module components/animations/FadeIn
 */
import React from 'react'
import { motion } from 'framer-motion'

/**
 * Props for FadeIn component
 */
export interface FadeInProps {
  /** Content to animate */
  children: React.ReactNode
  /** Delay before animation starts (in seconds) */
  delay?: number
  /** Animation duration (in seconds) */
  duration?: number
  /** Custom CSS classes */
  className?: string
  /** Animation axis (vertical or horizontal) */
  axis?: 'y' | 'x'
  /** Distance to move from (in pixels) */
  distance?: number
  /** Whether to disable animation */
  disabled?: boolean
}

/**
 * FadeIn Component
 *
 * Wraps children with a fade-in and slide animation.
 *
 * @example
 * ```tsx
 * // Basic fade-in
 * <FadeIn>
 *   <h1>Hello World</h1>
 * </FadeIn>
 *
 * // With delay for staggered effect
 * <FadeIn delay={0.2}>
 *   <p>Appears after 200ms</p>
 * </FadeIn>
 *
 * // Horizontal slide from left
 * <FadeIn axis="x" distance={-20}>
 *   <p>Slides from left</p>
 * </FadeIn>
 * ```
 */
export const FadeIn = React.forwardRef<HTMLDivElement, FadeInProps>(
  (
    {
      children,
      delay = 0,
      duration = 0.4,
      className,
      axis = 'y',
      distance = 10,
      disabled = false,
    },
    ref
  ) => {
    const [prefersReducedMotion, setPrefersReducedMotion] = React.useState(false)

    React.useEffect(() => {
      const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
      setPrefersReducedMotion(mediaQuery.matches)

      const handleChange = () => setPrefersReducedMotion(mediaQuery.matches)
      mediaQuery.addEventListener('change', handleChange)

      return () => mediaQuery.removeEventListener('change', handleChange)
    }, [])

    if (disabled || prefersReducedMotion) {
      return <div ref={ref} className={className}>{children}</div>
    }

    // Create animation variants based on axis
    const initial = { opacity: 0, ...(axis === 'x' ? { x: distance } : { y: distance }) }
    const animate = { opacity: 1, ...(axis === 'x' ? { x: 0 } : { y: 0 }) }

    return (
      <motion.div
        ref={ref}
        className={className}
        initial={initial}
        animate={animate}
        transition={{
          duration,
          delay,
          ease: [0.4, 0, 0.2, 1],
        }}
      >
        {children}
      </motion.div>
    )
  }
)

FadeIn.displayName = 'FadeIn'

export default FadeIn
