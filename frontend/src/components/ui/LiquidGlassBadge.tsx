import React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

const badgeVariants = cva(
  'liquid-glass-badge inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold',
  {
    variants: {
      variant: {
        default: '',
        blue: 'liquid-glass-badge-blue',
        teal: 'liquid-glass-badge-teal',
        green: 'liquid-glass-badge-green',
        red: 'liquid-glass-badge-red',
        success: 'liquid-glass-badge-green',
        warning: 'liquid-glass-badge-red',
        error: 'liquid-glass-badge-red',
      },
      size: {
        sm: 'text-[10px] px-2 py-0.5',
        md: 'text-xs px-2.5 py-1',
        lg: 'text-sm px-3 py-1.5',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  }
)

export interface LiquidGlassBadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {
  icon?: React.ReactNode
}

export const LiquidGlassBadge: React.FC<LiquidGlassBadgeProps> = ({
  className,
  variant,
  size,
  icon,
  children,
  ...props
}) => {
  return (
    <span className={cn(badgeVariants({ variant, size }), className)} {...props}>
      {icon}
      {children}
    </span>
  )
}

export default LiquidGlassBadge
