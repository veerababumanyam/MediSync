import React from 'react'
import { cn } from '@/lib/cn'

export interface LiquidGlassNavbarProps {
  children?: React.ReactNode
  left?: React.ReactNode
  center?: React.ReactNode
  right?: React.ReactNode
  className?: string
  sticky?: boolean
}

export const LiquidGlassNavbar: React.FC<LiquidGlassNavbarProps> = ({
  children,
  left,
  center,
  right,
  className,
  sticky = true,
}) => {
  return (
    <nav
      className={cn(
        'liquid-glass-header z-50 px-4 py-3 md:px-6',
        sticky && 'sticky top-0',
        className
      )}
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        {/* Left section */}
        {left && (
          <div className="flex items-center gap-4 flex-shrink-0">
            {left}
          </div>
        )}

        {/* Center section */}
        {center && (
          <div className="flex items-center justify-center flex-1">
            {center}
          </div>
        )}

        {/* Right section */}
        {right && (
          <div className="flex items-center gap-3 flex-shrink-0">
            {right}
          </div>
        )}

        {/* Default children */}
        {children}
      </div>
    </nav>
  )
}

export default LiquidGlassNavbar
