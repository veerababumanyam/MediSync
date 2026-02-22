/**
 * GlassCard Component
 *
 * A glassmorphic card container with frosted glass effect, backdrop blur,
 * and subtle borders. Supports multiple intensity levels, shadow depths,
 * and hover effects.
 *
 * Features:
 * - Configurable glass intensity (light, medium, heavy)
 * - Multiple shadow depths
 * - Hover effects (lift, glow)
 * - Built-in entrance/exit animations
 * - Full TypeScript support with CVA variants
 * - WCAG 2.2 AA accessible
 *
 * @module components/ui/GlassCard
 */
import React from 'react'
import { motion } from 'framer-motion'
import type { MotionProps } from 'framer-motion'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

/**
 * Glass card variant definitions using CVA (Class Variance Authority)
 *
 * Provides type-safe variant props for:
 * - intensity: Background and border opacity
 * - shadow: Depth of shadow
 * - hover: Interactive hover effects
 */
const glassCardVariants = cva(
  // Base classes: glass effect, backdrop blur, border, rounded corners, transition
  'backdrop-blur-xl border rounded-2xl transition-all duration-300',
  {
    variants: {
      // Glass intensity levels - affects background and border opacity
      intensity: {
        // Light: Most opaque, good for headers and navigation
        light: 'bg-white/30 dark:bg-white/10 border-white/40 dark:border-white/20 shadow-glass-lg',
        // Medium: Balanced opacity, default for most cards
        medium: 'bg-white/20 dark:bg-white/5 border-white/30 dark:border-white/15 shadow-glass-md',
        // Heavy: Most transparent, for overlays and floating elements
        heavy: 'bg-white/10 dark:bg-white/5 border-white/20 dark:border-white/10 shadow-glass-sm',
      },
      // Shadow depth - creates visual hierarchy
      shadow: {
        none: '',
        sm: 'shadow-glass-sm',
        md: 'shadow-glass-md',
        lg: 'shadow-glass-lg',
      },
      // Hover effects for interactive cards
      hover: {
        none: '',
        lift: 'hover:-translate-y-2 hover:shadow-glass-lg hover:bg-white/40 dark:hover:bg-white/15',
        glow: 'hover:shadow-[0_0_48px_rgba(139,92,246,0.5)] hover:border-purple-400/50',
        blueGlow: 'hover:shadow-[0_0_48px_rgba(59,130,246,0.5)] hover:border-blue-400/50',
        cyanGlow: 'hover:shadow-[0_0_48px_rgba(6,182,212,0.5)] hover:border-cyan-400/50',
      },
    },
    defaultVariants: {
      intensity: 'medium',
      shadow: 'md',
      hover: 'none',
    },
  }
)

/**
 * Props for GlassCard component
 *
 * Extends framer-motion's MotionProps for full animation control
 * and CVA's VariantProps for type-safe variants.
 */
export interface GlassCardProps
  extends Omit<React.PropsWithChildren<MotionProps<'div'>>, 'variants'>,
    VariantProps<typeof glassCardVariants> {
  /** Optional HTML tag to render (defaults to div) */
  as?: React.ElementType
  /** Whether to disable entrance animation */
  disableAnimation?: boolean
}

/**
 * GlassCard Component
 *
 * A glassmorphic container with frosted glass effect.
 *
 * @example
 * ```tsx
 * // Basic glass card
 * <GlassCard>
 *   <h2>Card Title</h2>
 *   <p>Card content</p>
 * </GlassCard>
 *
 * // Light intensity with lift hover effect
 * <GlassCard intensity="light" hover="lift">
 *   Interactive content
 * </GlassCard>
 *
 * // Heavy glass with glow for modals
 * <GlassCard intensity="heavy" shadow="lg" hover="glow">
 *   Modal content
 * </GlassCard>
 * ```
 */
export const GlassCard = React.forwardRef<HTMLDivElement, GlassCardProps>(
  (
    {
      className,
      intensity,
      shadow,
      hover,
      as = motion.div,
      disableAnimation = false,
      children,
      ...props
    },
    ref
  ) => {
    const Component = as

    return (
      <Component
        ref={ref}
        className={cn(glassCardVariants({ intensity, shadow, hover }), className)}
        initial={disableAnimation ? undefined : { opacity: 0, y: 20 }}
        animate={disableAnimation ? undefined : { opacity: 1, y: 0 }}
        exit={disableAnimation ? undefined : { opacity: 0, y: -20 }}
        transition={{
          duration: 0.3,
          ease: [0.4, 0, 0.2, 1],
        }}
        {...props}
      >
        {children}
      </Component>
    )
  }
)

GlassCard.displayName = 'GlassCard'

export default GlassCard
