/**
 * Card Component Patterns
 *
 * Reusable card patterns with WCAG 3.0 accessibility and iOS 26 liquid glass.
 */

import { forwardRef, HTMLAttributes, ReactNode } from 'react';
import { cn } from '@/lib/utils';

// ============================================
// BASE CARD
// ============================================

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  /** Card variant */
  variant?: 'default' | 'pronounced' | 'interactive';
  /** Enable hover effects */
  hoverable?: boolean;
  /** Card padding */
  padding?: 'none' | 'sm' | 'md' | 'lg';
  /** Child content */
  children: ReactNode;
}

export const Card = forwardRef<HTMLDivElement, CardProps>(
  (
    {
      variant = 'default',
      hoverable = false,
      padding = 'md',
      className,
      children,
      ...props
    },
    ref
  ) => {
    const paddingClasses = {
      none: '',
      sm: 'p-4',
      md: 'p-6',
      lg: 'p-8',
    };

    const variantClasses = {
      default: 'liquid-container-card',
      pronounced: 'liquid-container-card liquid-container-card--pronounced',
      interactive: 'liquid-container-card liquid-container-card--interactive',
    };

    return (
      <div
        ref={ref}
        className={cn(
          variantClasses[variant],
          paddingClasses[padding],
          hoverable && !variant.includes('interactive') && 'liquid-glass-hover-lift',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';

// ============================================
// CARD HEADER
// ============================================

interface CardHeaderProps extends HTMLAttributes<HTMLDivElement> {
  /** Header title */
  title: string;
  /** Optional subtitle */
  subtitle?: string;
  /** Optional action (button, link, etc.) */
  action?: ReactNode;
  /** Header level for SEO */
  as?: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';
}

export function CardHeader({
  title,
  subtitle,
  action,
  as: Component = 'h3',
  className,
  ...props
}: CardHeaderProps) {
  return (
    <div className={cn('flex items-start justify-between gap-4 mb-4', className)} {...props}>
      <div className="flex-1 min-w-0">
        <Component className="text-lg font-semibold text-primary truncate">
          {title}
        </Component>
        {subtitle && (
          <p className="text-sm text-secondary mt-1">{subtitle}</p>
        )}
      </div>
      {action && (
        <div className="shrink-0">{action}</div>
      )}
    </div>
  );
}

// ============================================
// CARD CONTENT
// ============================================

interface CardContentProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}

export function CardContent({ children, className, ...props }: CardContentProps) {
  return (
    <div className={cn('text-secondary', className)} {...props}>
      {children}
    </div>
  );
}

// ============================================
// CARD FOOTER
// ============================================

interface CardFooterProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
  /** Align items */
  align?: 'start' | 'center' | 'end' | 'between';
}

export function CardFooter({
  children,
  align = 'end',
  className,
  ...props
}: CardFooterProps) {
  const alignClasses = {
    start: 'justify-start',
    center: 'justify-center',
    end: 'justify-end',
    between: 'justify-between',
  };

  return (
    <div
      className={cn(
        'flex items-center gap-3 mt-6 pt-4 border-t border-subtle',
        alignClasses[align],
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

// ============================================
// STAT CARD
// ============================================

interface StatCardProps extends Omit<CardProps, 'children'> {
  /** Stat label */
  label: string;
  /** Stat value */
  value: string | number;
  /** Change indicator (e.g., "+12%" or "-5%") */
  change?: string;
  /** Whether change is positive */
  isPositive?: boolean;
  /** Optional icon */
  icon?: ReactNode;
}

export function StatCard({
  label,
  value,
  change,
  isPositive,
  icon,
  className,
  ...props
}: StatCardProps) {
  return (
    <Card className={cn('', className)} {...props}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-sm text-secondary font-medium">{label}</p>
          <p className="text-3xl font-bold text-primary mt-2">{value}</p>
          {change && (
            <p
              className={cn(
                'text-sm font-medium mt-2',
                isPositive ? 'text-success' : 'text-error'
              )}
            >
              {isPositive ? '↑' : '↓'} {change}
            </p>
          )}
        </div>
        {icon && (
          <div className="shrink-0 p-3 rounded-xl bg-brand-primary/10 text-brand-primary">
            {icon}
          </div>
        )}
      </div>
    </Card>
  );
}

// ============================================
// CLICKABLE CARD
// ============================================

interface ClickableCardProps extends Omit<CardProps, 'variant'> {
  /** Click handler */
  onClick?: () => void;
  /** Link href (renders as anchor) */
  href?: string;
  /** Card title */
  title: string;
  /** Card description */
  description?: string;
  /** Optional icon */
  icon?: ReactNode;
  /** Arrow indicator */
  showArrow?: boolean;
}

export const ClickableCard = forwardRef<HTMLDivElement, ClickableCardProps>(
  (
    {
      onClick,
      href,
      title,
      description,
      icon,
      showArrow = true,
      className,
      children,
      ...props
    },
    ref
  ) => {
    const content = (
      <>
        <div className="flex items-start gap-4">
          {icon && (
            <div className="shrink-0 p-3 rounded-xl bg-brand-primary/10 text-brand-primary">
              {icon}
            </div>
          )}
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-primary text-start">{title}</h3>
            {description && (
              <p className="text-sm text-secondary mt-1 text-start">{description}</p>
            )}
          </div>
          {showArrow && (
            <svg
              className="w-5 h-5 text-muted shrink-0 rtl:rotate-180"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          )}
        </div>
        {children}
      </>
    );

    const cardClassName = cn(
      'liquid-container-card',
      'liquid-container-card--interactive',
      'liquid-focus-gold',
      className
    );

    if (href) {
      return (
        <a
          ref={ref as any}
          href={href}
          className={cn(cardClassName, 'block no-underline')}
          {...(props as any)}
        >
          {content}
        </a>
      );
    }

    return (
      <div
        ref={ref}
        role="button"
        tabIndex={0}
        onClick={onClick}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            onClick?.();
          }
        }}
        className={cardClassName}
        {...props}
      >
        {content}
      </div>
    );
  }
);

ClickableCard.displayName = 'ClickableCard';

// ============================================
// FEATURE CARD
// ============================================

interface FeatureCardProps extends CardProps {
  /** Feature icon */
  icon: ReactNode;
  /** Feature title */
  title: string;
  /** Feature description */
  description: string;
  /** Optional badge */
  badge?: string;
}

export function FeatureCard({
  icon,
  title,
  description,
  badge,
  className,
  ...props
}: FeatureCardProps) {
  return (
    <Card className={cn('text-center', className)} {...props}>
      {badge && (
        <span className="inline-block px-3 py-1 mb-4 text-xs font-semibold rounded-full bg-brand-accent/10 text-brand-accent-700 dark:text-brand-accent-400">
          {badge}
        </span>
      )}
      <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-brand-primary/10 text-brand-primary mb-4">
        {icon}
      </div>
      <h3 className="text-lg font-semibold text-primary mb-2">{title}</h3>
      <p className="text-secondary text-sm">{description}</p>
    </Card>
  );
}

// ============================================
// EMPTY STATE CARD
// ============================================

interface EmptyStateCardProps extends CardProps {
  /** Empty state icon */
  icon?: ReactNode;
  /** Title */
  title: string;
  /** Description */
  description?: string;
  /** Action button */
  action?: ReactNode;
}

export function EmptyStateCard({
  icon,
  title,
  description,
  action,
  className,
  ...props
}: EmptyStateCardProps) {
  return (
    <Card className={cn('text-center py-12', className)} {...props}>
      {icon && (
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-tertiary/20 text-tertiary mb-4">
          {icon}
        </div>
      )}
      <h3 className="text-lg font-semibold text-primary mb-2">{title}</h3>
      {description && (
        <p className="text-secondary text-sm max-w-sm mx-auto mb-4">{description}</p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </Card>
  );
}

// ============================================
// SKELETON CARD
// ============================================

export function SkeletonCard({ className }: { className?: string }) {
  return (
    <Card className={cn('animate-pulse', className)}>
      <div className="liquid-skeleton h-4 w-3/4 rounded mb-4" />
      <div className="liquid-skeleton h-8 w-1/2 rounded mb-4" />
      <div className="liquid-skeleton h-4 w-full rounded mb-2" />
      <div className="liquid-skeleton h-4 w-2/3 rounded" />
    </Card>
  );
}

export default Card;
