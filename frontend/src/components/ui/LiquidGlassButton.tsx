/**
 * Liquid Glass Button Component
 *
 * Premium iOS-inspired glassmorphic button with liquid animations,
 * dynamic lighting effects, and WCAG 2.2 AA compliance.
 *
 * Features:
 * - Multiple visual variants (glass, primary, secondary, ghost)
 * - Liquid hover animations (lift, glow, ripple)
 * - Branded color variants using logo colors
 * - Focus indicators for keyboard navigation
 * - Loading state with spinner
 * - Icon support with positioning
 * - Reduced motion support
 * - RTL support
 *
 * @module components/ui/LiquidGlassButton
 * @version 2.0.0
 */

import React, { forwardRef, useState, useCallback, type ComponentProps } from 'react'
import { motion } from 'framer-motion'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

// Type for motion button props - compatible with framer-motion v12
type MotionButtonProps = ComponentProps<typeof motion.button>

/**
 * Liquid glass button variant definitions
 */
const liquidButtonVariants = cva(
  // Base classes
  'liquid-glass-button inline-flex items-center justify-center gap-2 font-medium transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-1 disabled:opacity-50 disabled:cursor-not-allowed disabled:pointer-events-none',
  {
    variants: {
      // Button style variants
      variant: {
        // Glass button with transparent background
        glass: 'liquid-glass px-4 py-2',
        // Primary button with brand gradient
        primary: 'liquid-glass-button-primary px-5 py-2.5 text-white',
        // Secondary button with border
        secondary: 'px-4 py-2 border-2 border-glass text-primary bg-surface-glass hover:bg-surface-glass-strong rounded-lg',
        // Ghost button (minimal styling)
        ghost: 'px-4 py-2 text-secondary hover:bg-surface-glass rounded-lg',
        // Danger button
        danger: 'px-5 py-2.5 bg-red-500 text-white hover:bg-red-600 rounded-lg shadow-lg shadow-red-500/25',
      },
      // Size variants
      size: {
        xs: 'text-xs px-2.5 py-1.5 gap-1.5',
        sm: 'text-sm px-3 py-2 gap-2',
        md: 'text-base px-4 py-2.5 gap-2',
        lg: 'text-lg px-6 py-3 gap-2.5',
        xl: 'text-xl px-8 py-4 gap-3',
      },
      // Border radius
      radius: {
        sm: 'rounded-md',
        md: 'rounded-lg',
        lg: 'rounded-xl',
        xl: 'rounded-2xl',
        full: 'rounded-full',
      },
      // Hover effect
      hover: {
        none: '',
        lift: 'hover:-translate-y-0.5',
        glow: 'hover:shadow-lg',
        scale: 'hover:scale-105',
      },
    },
    defaultVariants: {
      variant: 'glass',
      size: 'md',
      radius: 'lg',
      hover: 'lift',
    },
  }
)

/**
 * Icon positions
 */
type IconPosition = 'left' | 'right' | 'only'

/**
 * Props for LiquidGlassButton component
 */
export interface LiquidGlassButtonProps
  extends Omit<MotionButtonProps, 'disabled' | 'variants'>,
  VariantProps<typeof liquidButtonVariants> {
  /** Optional icon element to display */
  icon?: React.ReactNode
  /** Position of the icon (default: 'left') */
  iconPosition?: IconPosition
  /** Whether button is in loading state */
  isLoading?: boolean
  /** Whether button is disabled */
  disabled?: boolean
  /** Optional HTML tag to render (defaults to button) */
  as?: React.ElementType
  /** Custom click handler */
  onClick?: React.MouseEventHandler<HTMLButtonElement>
}

/**
 * Liquid Glass Button Component
 *
 * A premium glassmorphic button with liquid animations.
 *
 * @example
 * ```tsx
 * // Basic glass button
 * <LiquidGlassButton>Click me</LiquidGlassButton>
 *
 * // Primary button with icon
 * <LiquidGlassButton variant="primary" icon={<Icon />}>
 *   Save Changes
 * </LiquidGlassButton>
 *
 * // Loading state
 * <LiquidGlassButton variant="primary" isLoading>
 *   Processing...
 * </LiquidGlassButton>
 *
 * // Icon only button
 * <LiquidGlassButton variant="glass" icon={<Icon />} iconPosition="only" />
 * ```
 */
export const LiquidGlassButton = forwardRef<HTMLButtonElement, LiquidGlassButtonProps>(
  (
    {
      className,
      variant,
      size,
      radius,
      hover,
      icon,
      iconPosition = 'left',
      isLoading = false,
      disabled = false,
      as = motion.button,
      onClick,
      children,
      ...props
    },
    ref
  ) => {
    const [isPressed, setIsPressed] = useState(false)
    const [ripplePosition, setRipplePosition] = useState({ x: '50%', y: '50%' })

    // Handle ripple effect on click
    const handleMouseDown = useCallback((e: React.MouseEvent<HTMLButtonElement>) => {
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
    const buttonClasses = cn(
      liquidButtonVariants({ variant, size, radius, hover }),
      isPressed && 'scale-95',
      isLoading && 'relative',
      className
    )

    // Loading spinner
    const spinner = (
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
    )

    // Render icon based on position
    const renderIcon = () => {
      if (!icon && !isLoading) return null

      const iconToRender = isLoading ? spinner : icon

      if (iconPosition === 'only') {
        return iconToRender
      }

      return <span className="shrink-0">{iconToRender}</span>
    }

    const style = {
      '--ripple-x': ripplePosition.x,
      '--ripple-y': ripplePosition.y,
    } as React.CSSProperties

    return (
      <Component
        ref={ref}
        className={buttonClasses}
        style={style}
        disabled={disabled || isLoading}
        onClick={onClick}
        onMouseDown={handleMouseDown}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseLeave}
        whileTap={{ scale: 0.97 }}
        transition={{ duration: 0.1 }}
        {...props}
      >
        {iconPosition === 'left' && renderIcon()}
        {iconPosition !== 'only' && <span>{children as React.ReactNode}</span>}
        {iconPosition === 'right' && renderIcon()}
      </Component>
    )
  }
)

LiquidGlassButton.displayName = 'LiquidGlassButton'

/**
 * Preset button configurations for common use cases
 */

/**
 * Primary action button
 */
export const ButtonPrimary = forwardRef<HTMLButtonElement, Omit<LiquidGlassButtonProps, 'variant'>>(
  (props, ref) => (
    <LiquidGlassButton ref={ref} variant="primary" {...props} />
  )
)
ButtonPrimary.displayName = 'ButtonPrimary'

/**
 * Secondary action button
 */
export const ButtonSecondary = forwardRef<HTMLButtonElement, Omit<LiquidGlassButtonProps, 'variant'>>(
  (props, ref) => (
    <LiquidGlassButton ref={ref} variant="secondary" {...props} />
  )
)
ButtonSecondary.displayName = 'ButtonSecondary'

/**
 * Ghost button (minimal styling)
 */
export const ButtonGhost = forwardRef<HTMLButtonElement, Omit<LiquidGlassButtonProps, 'variant'>>(
  (props, ref) => (
    <LiquidGlassButton ref={ref} variant="ghost" {...props} />
  )
)
ButtonGhost.displayName = 'ButtonGhost'

/**
 * Danger button
 */
export const ButtonDanger = forwardRef<HTMLButtonElement, Omit<LiquidGlassButtonProps, 'variant'>>(
  (props, ref) => (
    <LiquidGlassButton ref={ref} variant="danger" {...props} />
  )
)
ButtonDanger.displayName = 'ButtonDanger'

/**
 * Icon button (circular, icon only)
 */
export const IconButton = forwardRef<HTMLButtonElement, Omit<LiquidGlassButtonProps, 'variant' | 'radius' | 'iconPosition'>>(
  ({ icon, size = 'md', ...props }, ref) => {
    const sizeClasses = {
      xs: 'p-1.5',
      sm: 'p-2',
      md: 'p-2.5',
      lg: 'p-3',
      xl: 'p-4',
    }

    return (
      <LiquidGlassButton
        ref={ref}
        variant="glass"
        icon={icon}
        iconPosition="only"
        radius="full"
        className={cn(sizeClasses[size || 'md'])}
        {...props}
      />
    )
  }
)
IconButton.displayName = 'IconButton'

export default LiquidGlassButton
