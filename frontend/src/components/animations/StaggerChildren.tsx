/**
 * StaggerChildren Animation Component
 *
 * Wraps a list of children and animates them in sequence with a
 * staggered delay. Perfect for lists, grids, and card layouts.
 *
 * Features:
 * - Configurable stagger delay between children
 * - Respects prefers-reduced-motion (WCAG 2.2)
 * - Works with direct children or nested arrays
 * - Customizable animation duration and distance
 *
 * @module components/animations/StaggerChildren
 */
import React from 'react'
import { motion } from 'framer-motion'

/**
 * Props for StaggerChildren component
 */
export interface StaggerChildrenProps {
  /** Children to animate in sequence */
  children: React.ReactNode
  /** Delay between each child (in seconds) */
  staggerDelay?: number
  /** Animation duration (in seconds) */
  duration?: number
  /** Custom CSS classes for container */
  className?: string
  /** Animation axis (vertical or horizontal) */
  axis?: 'y' | 'x'
  /** Distance to move from (in pixels) */
  distance?: number
  /** Whether to disable animation */
  disabled?: boolean
}

/**
 * StaggerChildren Component
 *
 * Animates children in sequence with a staggered delay.
 *
 * @example
 * ```tsx
 * // Stagger a list of items
 * <StaggerChildren>
 *   <Item>First</Item>
 *   <Item>Second</Item>
 *   <Item>Third</Item>
 * </StaggerChildren>
 *
 * // With custom stagger delay
 * <StaggerChildren staggerDelay={0.2}>
 *   <Item>Slow stagger</Item>
 * </StaggerChildren>
 *
 * // Horizontal stagger
 * <StaggerChildren axis="x">
 *   <Item>Slides horizontally</Item>
 * </StaggerChildren>
 * ```
 */
export const StaggerChildren = React.forwardRef<HTMLDivElement, StaggerChildrenProps>(
  (
    {
      children,
      staggerDelay = 0.1,
      duration = 0.4,
      className,
      axis = 'y',
      distance = 20,
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

    // Animation variants for container
    const containerVariants = {
      hidden: { opacity: 0 },
      visible: {
        opacity: 1,
        transition: {
          staggerChildren: staggerDelay,
          delayChildren: 0,
        },
      },
    }

    // Animation variants for each child - simpler version without custom function
    const itemVariants = {
      hidden: {
        opacity: 0,
        [axis]: distance,
      } as any,
      visible: {
        opacity: 1,
        [axis]: 0,
        transition: {
          duration,
          ease: [0.4, 0, 0.2, 1],
        },
      } as any,
    }

    if (disabled || prefersReducedMotion) {
      return <div ref={ref} className={className}>{children}</div>
    }

    return (
      <motion.div
        ref={ref}
        className={className}
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        {React.Children.map(children, (child) => (
          <motion.div
            variants={itemVariants}
          >
            {child}
          </motion.div>
        ))}
      </motion.div>
    )
  }
)

StaggerChildren.displayName = 'StaggerChildren'

export default StaggerChildren
