/**
 * LiquidGlass Component
 *
 * A versatile glass morphism component following iOS 26 design principles
 * with WCAG 3.0 accessibility compliance.
 */

import { forwardRef, HTMLAttributes, ReactNode } from 'react';
import { cn } from '@/lib/utils';

// Glass variant types
type GlassVariant = 'subtle' | 'light' | 'medium' | 'heavy' | 'pronounced';
type GlassElevation = 'ambient' | 'raised' | 'floating' | 'elevated';
type GlassTint = 'none' | 'blue' | 'teal' | 'green';

interface LiquidGlassProps extends HTMLAttributes<HTMLDivElement> {
  /** Glass opacity variant */
  variant?: GlassVariant;
  /** Shadow elevation level */
  elevation?: GlassElevation;
  /** Brand color tint */
  tint?: GlassTint;
  /** Enable hover lift effect */
  hoverLift?: boolean;
  /** Enable hover glow effect */
  hoverGlow?: boolean;
  /** Enable shimmer effect on hover */
  hoverShimmer?: boolean;
  /** Make component interactive (adds focus styles) */
  interactive?: boolean;
  /** Apply gold focus ring for accessibility */
  goldFocus?: boolean;
  /** Child content */
  children: ReactNode;
}

/**
 * LiquidGlass Component
 *
 * @example
 * // Basic glass card
 * <LiquidGlass>
 *   <h3>Card Title</h3>
 *   <p>Card content</p>
 * </LiquidGlass>
 *
 * @example
 * // Interactive card with hover effects
 * <LiquidGlass
 *   interactive
 *   hoverLift
 *   hoverGlow
 *   goldFocus
 * >
 *   Clickable Card
 * </LiquidGlass>
 *
 * @example
 * // Pronounced glass modal
 * <LiquidGlass variant="pronounced" elevation="elevated">
 *   Modal Content
 * </LiquidGlass>
 */
export const LiquidGlass = forwardRef<HTMLDivElement, LiquidGlassProps>(
  (
    {
      variant = 'medium',
      elevation = 'raised',
      tint = 'none',
      hoverLift = false,
      hoverGlow = false,
      hoverShimmer = false,
      interactive = false,
      goldFocus = true,
      className,
      children,
      ...props
    },
    ref
  ) => {
    const baseClasses = [
      'liquid-glass',
      `liquid-glass-${variant}`,
      `liquid-shadow-${elevation}`,
      tint !== 'none' && `liquid-glass-${tint}`,
      hoverLift && 'liquid-glass-hover-lift',
      hoverGlow && 'liquid-glass-hover-glow',
      hoverShimmer && 'liquid-glass-hover-shimmer',
      interactive && 'cursor-pointer',
      goldFocus && 'liquid-focus-gold',
      className,
    ]
      .filter(Boolean)
      .join(' ');

    return (
      <div
        ref={ref}
        className={baseClasses}
        tabIndex={interactive ? 0 : undefined}
        role={interactive ? 'button' : undefined}
        {...props}
      >
        {children}
      </div>
    );
  }
);

LiquidGlass.displayName = 'LiquidGlass';

// Preset Components

/**
 * GlassCard - Standard content card with glass effect
 */
interface GlassCardProps extends Omit<LiquidGlassProps, 'variant' | 'elevation'> {
  /** Card title */
  title?: string;
  /** Card description */
  description?: string;
}

export const GlassCard = forwardRef<HTMLDivElement, GlassCardProps>(
  ({ title, description, children, className, ...props }, ref) => {
    return (
      <LiquidGlass
        ref={ref}
        variant="heavy"
        elevation="raised"
        className={cn('p-6', className)}
        {...props}
      >
        {title && (
          <h3 className="text-lg font-semibold text-primary mb-2">{title}</h3>
        )}
        {description && (
          <p className="text-secondary text-sm mb-4">{description}</p>
        )}
        {children}
      </LiquidGlass>
    );
  }
);

GlassCard.displayName = 'GlassCard';

/**
 * GlassPanel - Content panel with pronounced glass effect
 */
export const GlassPanel = forwardRef<HTMLDivElement, LiquidGlassProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <LiquidGlass
        ref={ref}
        variant="pronounced"
        elevation="elevated"
        className={cn('p-8', className)}
        {...props}
      >
        {children}
      </LiquidGlass>
    );
  }
);

GlassPanel.displayName = 'GlassPanel';

/**
 * InteractiveCard - Clickable card with full hover/focus effects
 */
export const InteractiveCard = forwardRef<HTMLDivElement, GlassCardProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <LiquidGlass
        ref={ref}
        variant="heavy"
        elevation="raised"
        interactive
        hoverLift
        hoverGlow
        goldFocus
        className={cn('p-6', className)}
        {...props}
      >
        {children}
      </LiquidGlass>
    );
  }
);

InteractiveCard.displayName = 'InteractiveCard';

/**
 * GlassHeader - Sticky glass header with strong blur
 */
export const GlassHeader = forwardRef<HTMLElement, HTMLAttributes<HTMLElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <header
        ref={ref}
        className={cn(
          'liquid-glass-header',
          'sticky top-0 z-sticky',
          className
        )}
        {...props}
      >
        {children}
      </header>
    );
  }
);

GlassHeader.displayName = 'GlassHeader';

/**
 * GlassModal - Modal container with maximum glass effect
 */
export const GlassModal = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        role="dialog"
        aria-modal="true"
        className={cn(
          'liquid-glass-modal',
          'liquid-glass-pronounced',
          'p-8',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

GlassModal.displayName = 'GlassModal';

export default LiquidGlass;
