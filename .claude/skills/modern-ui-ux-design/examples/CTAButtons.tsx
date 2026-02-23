/**
 * CTA Button Components
 *
 * Centralized CTA button variants following iOS 26 HIG and WCAG 3.0 standards.
 * All buttons include:
 * - min-height: 44px (WCAG 2.5.5 touch target)
 * - Gold focus indicators for accessibility
 * - Consistent spacing and typography
 * - Theme-aware via .dark selector
 */

import { forwardRef, ButtonHTMLAttributes, ReactNode } from 'react';
import { cn } from '@/lib/utils';

// Base CTA props
interface CTAButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Button content */
  children: ReactNode;
  /** Icon to display before text */
  startIcon?: ReactNode;
  /** Icon to display after text */
  endIcon?: ReactNode;
  /** Loading state */
  isLoading?: boolean;
  /** Full width */
  fullWidth?: boolean;
}

/**
 * HeroCTA - Most prominent CTA for hero sections
 *
 * @example
 * <HeroCTA onClick={handleClick}>
 *   Get Started Free
 * </HeroCTA>
 */
export const HeroCTA = forwardRef<HTMLButtonElement, CTAButtonProps>(
  (
    {
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        className={cn(
          'cta-hero',
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-5 h-5" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span className="relative z-10">{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

HeroCTA.displayName = 'HeroCTA';

/**
 * LandingPrimaryCTA - Primary action on landing pages
 *
 * @example
 * <LandingPrimaryCTA href="/signup">
 *   Start Your Free Trial
 * </LandingPrimaryCTA>
 */
export const LandingPrimaryCTA = forwardRef<HTMLButtonElement, CTAButtonProps>(
  (
    {
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        className={cn(
          'cta-landing-primary',
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-5 h-5" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span>{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

LandingPrimaryCTA.displayName = 'LandingPrimaryCTA';

/**
 * LandingSecondaryCTA - Secondary action on landing pages
 *
 * @example
 * <LandingSecondaryCTA href="/demo">
 *   Watch Demo
 * </LandingSecondaryCTA>
 */
export const LandingSecondaryCTA = forwardRef<HTMLButtonElement, CTAButtonProps>(
  (
    {
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        className={cn(
          'cta-landing-secondary',
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-5 h-5 border-white/30 border-t-white" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span>{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

LandingSecondaryCTA.displayName = 'LandingSecondaryCTA';

/**
 * HeaderPrimaryCTA - Compact CTA for headers/navigation
 *
 * @example
 * <HeaderPrimaryCTA onClick={handleSubscribe}>
 *   Subscribe
 * </HeaderPrimaryCTA>
 */
export const HeaderPrimaryCTA = forwardRef<HTMLButtonElement, CTAButtonProps>(
  (
    {
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        className={cn(
          'cta-header-primary',
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-4 h-4" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span>{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

HeaderPrimaryCTA.displayName = 'HeaderPrimaryCTA';

/**
 * GlassButton - Glass morphism button with brand color
 *
 * @example
 * <GlassButton variant="blue">
 *   <Icon />
 *   Action
 * </GlassButton>
 */
interface GlassButtonProps extends CTAButtonProps {
  /** Glass tint color */
  variant?: 'default' | 'blue' | 'teal' | 'green';
  /** Size variant */
  size?: 'sm' | 'md' | 'lg';
}

export const GlassButton = forwardRef<HTMLButtonElement, GlassButtonProps>(
  (
    {
      variant = 'default',
      size = 'md',
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      sm: 'px-3 py-2 text-sm min-h-10',
      md: 'px-4 py-2.5 text-base min-h-11',
      lg: 'px-6 py-3 text-lg min-h-12',
    };

    return (
      <button
        ref={ref}
        className={cn(
          'liquid-glass-button',
          'liquid-focus-gold',
          variant !== 'default' && `liquid-glass-${variant}`,
          sizeClasses[size],
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-5 h-5" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span>{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

GlassButton.displayName = 'GlassButton';

/**
 * ProminentGlassButton - iOS 26 style prominent button
 *
 * @example
 * <ProminentGlassButton>
 *   Get Started
 * </ProminentGlassButton>
 */
export const ProminentGlassButton = forwardRef<HTMLButtonElement, CTAButtonProps>(
  (
    {
      children,
      startIcon,
      endIcon,
      isLoading,
      fullWidth,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        className={cn(
          'liquid-glass-button-prominent',
          'liquid-focus-gold',
          fullWidth && 'w-full',
          disabled && 'opacity-50 cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <span className="liquid-spinner w-5 h-5" />
        ) : (
          <>
            {startIcon && <span className="shrink-0">{startIcon}</span>}
            <span>{children}</span>
            {endIcon && <span className="shrink-0">{endIcon}</span>}
          </>
        )}
      </button>
    );
  }
);

ProminentGlassButton.displayName = 'ProminentGlassButton';

/**
 * IconButton - Compact button with icon only
 *
 * @example
 * <IconButton aria-label="Settings">
 *   <SettingsIcon />
 * </IconButton>
 */
interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Icon element */
  children: ReactNode;
  /** Button size */
  size?: 'sm' | 'md' | 'lg';
  /** Glass variant */
  variant?: 'glass' | 'ghost' | 'filled';
}

export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (
    {
      children,
      size = 'md',
      variant = 'glass',
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      sm: 'w-8 h-8',
      md: 'w-10 h-10',
      lg: 'w-12 h-12',
    };

    const variantClasses = {
      glass: 'liquid-glass-button liquid-glass-subtle',
      ghost: 'bg-transparent hover:bg-hover',
      filled: 'bg-brand-primary text-white',
    };

    return (
      <button
        ref={ref}
        className={cn(
          'inline-flex items-center justify-center',
          'rounded-lg',
          'liquid-focus-gold',
          'transition-all duration-fast ease-liquid',
          'disabled:opacity-50 disabled:cursor-not-allowed',
          sizeClasses[size],
          variantClasses[variant],
          className
        )}
        disabled={disabled}
        {...props}
      >
        {children}
      </button>
    );
  }
);

IconButton.displayName = 'IconButton';
