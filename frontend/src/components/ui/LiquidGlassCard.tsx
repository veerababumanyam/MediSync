/**
 * Liquid Glass Card Component
 *
 * Premium iOS-inspired glassmorphic card with liquid animations,
 * dynamic lighting effects, and WCAG 2.2 AA compliance.
 *
 * Features:
 * - Multi-layered glass effect with specular highlights
 * - Liquid hover animations (lift, glow, shimmer)
 * - Branded color variants using logo colors
 * - Focus indicators for keyboard navigation
 * - Reduced motion support
 * - RTL support
 *
 * @module components/ui/LiquidGlassCard
 * @version 2.0.0
 */

import React, { forwardRef, useState, useCallback, type ComponentProps } from 'react'
import { motion } from 'framer-motion'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

// Type for motion div props - compatible with framer-motion v12
type MotionDivProps = ComponentProps<typeof motion.div>

/**
 * Liquid glass card variant definitions
 *
 * Provides type-safe variant props for:
 * - intensity: Background and blur opacity
 * - elevation: Shadow depth for visual hierarchy
 * - hover: Interactive hover effects
 * - brand: Logo color integration
 */
const liquidGlassVariants = cva(
  // Base classes: multi-layered glass, backdrop blur, border, radius
  'liquid-glass relative overflow-hidden',
  {
    variants: {
      // Glass intensity - affects background opacity and blur strength
      intensity: {
        // Subtle: Minimal glass effect (for backgrounds, cards within cards)
        subtle: 'liquid-glass-subtle',
        // Light: Most opaque (for headers, navigation, cards needing high contrast)
        light: 'liquid-glass-light',
        // Medium: Balanced opacity (default for most cards)
        medium: '',
        // Heavy: Most transparent (for overlays, floating panels, modals)
        heavy: 'liquid-glass-heavy',
      },
      // Elevation level - shadow depth for visual hierarchy
      elevation: {
        none: '',
        base: 'liquid-shadow-ambient',
        raised: 'liquid-shadow-elevation',
        floating: 'liquid-shadow-float',
      },
      // Hover effects for interactive elements
      hover: {
        none: '',
        // Gentle lift on hover
        lift: 'liquid-glass-hover-lift',
        // Colored glow (teal by default)
        glow: 'liquid-glass-hover-glow',
        // Blue glow
        'glow-blue': 'liquid-glass-hover-glow-blue',
        // Green glow
        'glow-green': 'liquid-glass-hover-glow-green',
        // Moving shimmer effect
        shimmer: 'liquid-glass-hover-shimmer',
        // Combined lift + glow
        'lift-glow': 'liquid-glass-hover-lift liquid-glass-hover-glow',
      },
      // Brand color variants using logo colors
      brand: {
        none: '',
        // Deep blue accent
        blue: 'liquid-glass-blue',
        // Teal accent
        teal: 'liquid-glass-teal',
        // Green accent
        green: 'liquid-glass-green',
        // Full gradient (logo colors)
        brand: 'liquid-glass-brand',
      },
      // Border radius variants
      radius: {
        sm: 'liquid-radius-sm',
        md: 'liquid-radius-md',
        lg: 'liquid-radius-lg',
        xl: 'liquid-radius-xl',
        '2xl': 'liquid-radius-2xl',
        full: 'liquid-radius-full',
      },
    },
    defaultVariants: {
      intensity: 'medium',
      elevation: 'raised',
      hover: 'none',
      brand: 'none',
      radius: 'lg',
    },
  }
)

/**
 * Props for LiquidGlassCard component
 */
export interface LiquidGlassCardProps
  extends Omit<MotionDivProps, 'variants' | 'transition'>,
  VariantProps<typeof liquidGlassVariants> {
  /** Optional HTML tag to render (defaults to div) */
  as?: React.ElementType
  /** Whether to disable entrance animation */
  disableAnimation?: boolean
  /** Animation delay in milliseconds */
  animationDelay?: number
  /** Whether to add gradient overlay effect */
  gradientOverlay?: boolean
  /** Whether card is clickable/interactive */
  interactive?: boolean
  /** Whether to show pulse glow animation */
  pulseGlow?: boolean
  /** Whether to show floating animation */
  float?: boolean
  /** Custom click handler */
  onClick?: () => void
}

/**
 * Liquid Glass Card Component
 *
 * A premium glassmorphic container with liquid animations and
 * iOS-inspired aesthetics.
 *
 * @example
 * ```tsx
 * // Basic card
 * <LiquidGlassCard>
 *   <h2>Card Title</h2>
 *   <p>Card content</p>
 * </LiquidGlassCard>
 *
 * // Interactive card with teal glow
 * <LiquidGlassCard
 *   hover="glow"
 *   brand="teal"
 *   interactive
 *   onClick={handleClick}
 * >
 *   Interactive content
 * </LiquidGlassCard>
 *
 * // Brand gradient card with shimmer
 * <LiquidGlassCard
 *   brand="brand"
 *   hover="shimmer"
 *   elevation="floating"
 * >
 *   Premium content
 * </LiquidGlassCard>
 *
 * // Header card (light intensity)
 * <LiquidGlassCard
 *   intensity="light"
 *   radius="xl"
 * >
 *   Header content
 * </LiquidGlassCard>
 * ```
 */
export const LiquidGlassCard = forwardRef<HTMLDivElement, LiquidGlassCardProps>(
  (
    {
      className,
      intensity,
      elevation,
      hover,
      brand,
      radius,
      as = motion.div,
      disableAnimation = false,
      animationDelay = 0,
      gradientOverlay = false,
      interactive = false,
      pulseGlow = false,
      float = false,
      onClick,
      children,
      ...props
    },
    ref
  ) => {
    const [isPressed, setIsPressed] = useState(false)
    const [ripplePosition, setRipplePosition] = useState({ x: '50%', y: '50%' })

    // Handle ripple effect on click
    const handleMouseDown = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
      const rect = e.currentTarget.getBoundingClientRect()
      const x = ((e.clientX - rect.left) / rect.width) * 100
      const y = ((e.clientY - rect.top) / rect.height) * 100
      setRipplePosition({ x: `${x}%`, y: `${y}%` })
      setIsPressed(true)
    }, [])

    const handleMouseUp = useCallback(() => {
      setIsPressed(false)
    }, [])

    const handleMouseLeave = useCallback(() => {
      setIsPressed(false)
    }, [])

    const Component = as

    // Build class names
    const cardClasses = cn(
      liquidGlassVariants({
        intensity,
        elevation,
        hover: interactive && !onClick ? hover : hover,
        brand,
        radius,
      }),
      // Interactive classes
      interactive && 'cursor-pointer',
      onClick && 'liquid-ripple',
      onClick && interactive && 'liquid-glass-interactive',
      isPressed && 'liquid-glass-pressed',
      pulseGlow && 'liquid-pulse-glow',
      float && 'liquid-float',
      gradientOverlay && 'liquid-gradient-overlay',
      className
    )

    // Animation variants
    const animationVariants = {
      initial: { opacity: 0, y: 20, scale: 0.96 },
      animate: { opacity: 1, y: 0, scale: 1 },
      exit: { opacity: 0, y: -10, scale: 0.98 },
    }

    // Set CSS custom properties for ripple
    const style = {
      '--ripple-x': ripplePosition.x,
      '--ripple-y': ripplePosition.y,
      animationDelay: animationDelay ? `${animationDelay}ms` : undefined,
    } as React.CSSProperties

    return (
      <Component
        ref={ref}
        className={cardClasses}
        style={style}
        initial={disableAnimation ? undefined : 'initial'}
        animate={disableAnimation ? undefined : 'animate'}
        exit={disableAnimation ? undefined : 'exit'}
        variants={animationVariants}
        transition={{
          duration: 0.4,
          ease: [0.4, 0, 0.2, 1],
          delay: animationDelay / 1000,
        }}
        onMouseDown={onClick ? handleMouseDown : undefined}
        onMouseUp={onClick ? handleMouseUp : undefined}
        onMouseLeave={onClick ? handleMouseLeave : undefined}
        onClick={onClick}
        role={onClick ? 'button' : undefined}
        tabIndex={onClick ? 0 : undefined}
        onKeyDown={onClick ? (e: React.KeyboardEvent<HTMLDivElement>) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault()
            onClick()
          }
        } : undefined}
        {...props}
      >
        {children}
      </Component>
    )
  }
)

LiquidGlassCard.displayName = 'LiquidGlassCard'

/**
 * Preset card configurations for common use cases
 */

/**
 * Default card - balanced glass effect
 */
export const GlassCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'intensity' | 'elevation'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      intensity="medium"
      elevation="raised"
      {...props}
    />
  )
)
GlassCard.displayName = 'GlassCard'

/**
 * Header card - more opaque, sits at top
 */
export const GlassHeader = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'intensity' | 'elevation' | 'radius'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      intensity="light"
      elevation="base"
      radius="xl"
      {...props}
    />
  )
)
GlassHeader.displayName = 'GlassHeader'

/**
 * Modal card - highest blur and elevation
 */
export const GlassModal = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'intensity' | 'elevation' | 'radius'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      intensity="light"
      elevation="floating"
      radius="2xl"
      {...props}
    />
  )
)
GlassModal.displayName = 'GlassModal'

/**
 * Interactive card with hover effects
 */
export const GlassInteractiveCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'hover' | 'interactive'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      intensity="medium"
      elevation="raised"
      hover="lift-glow"
      interactive
      {...props}
    />
  )
)
GlassInteractiveCard.displayName = 'GlassInteractiveCard'

/**
 * Brand card with logo gradient
 */
export const GlassBrandCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'brand' | 'hover'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      brand="brand"
      hover="shimmer"
      gradientOverlay
      {...props}
    />
  )
)
GlassBrandCard.displayName = 'GlassBrandCard'

/**
 * Blue-themed card
 */
export const GlassBlueCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'brand' | 'hover'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      brand="blue"
      hover="glow-blue"
      {...props}
    />
  )
)
GlassBlueCard.displayName = 'GlassBlueCard'

/**
 * Teal-themed card
 */
export const GlassTealCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'brand' | 'hover'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      brand="teal"
      hover="glow"
      {...props}
    />
  )
)
GlassTealCard.displayName = 'GlassTealCard'

/**
 * Green-themed card (success, growth)
 */
export const GlassGreenCard = forwardRef<HTMLDivElement, Omit<LiquidGlassCardProps, 'brand' | 'hover'>>(
  (props, ref) => (
    <LiquidGlassCard
      ref={ref}
      brand="green"
      hover="glow-green"
      {...props}
    />
  )
)
GlassGreenCard.displayName = 'GlassGreenCard'

export default LiquidGlassCard
