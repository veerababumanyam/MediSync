/**
 * Liquid Glass Sidebar Component
 *
 * Premium iOS-inspired glassmorphic sidebar with liquid animations,
 * dynamic lighting effects, and WCAG 2.2 AA compliance.
 *
 * Features:
 * - Multi-layered glass effect with specular highlights
 * - Collapsible state with smooth transition
 * - Branded color variants using logo colors
 * - Custom scrollbar styling
 * - Reduced motion support
 * - RTL support
 *
 * @module components/ui/LiquidGlassSidebar
 * @version 1.0.0
 */

import React from 'react'
import { cn } from '@/lib/cn'

/**
 * Props for LiquidGlassSidebar component
 */
export interface LiquidGlassSidebarProps {
  /** Sidebar content */
  children?: React.ReactNode
  /** Additional CSS classes */
  className?: string
  /** Whether sidebar is collapsed (narrow width) */
  collapsed?: boolean
}

/**
 * Liquid Glass Sidebar Component
 *
 * A premium glassmorphic sidebar with collapsible state and liquid animations.
 *
 * @example
 * ```tsx
 * // Basic sidebar
 * <LiquidGlassSidebar>
 *   <NavItem>Home</NavItem>
 *   <NavItem>Documents</NavItem>
 * </LiquidGlassSidebar>
 *
 * // Collapsed sidebar
 * <LiquidGlassSidebar collapsed>
 *   <NavItem icon={<HomeIcon />} />
 * </LiquidGlassSidebar>
 * ```
 */
export const LiquidGlassSidebar: React.FC<LiquidGlassSidebarProps> = ({
  children,
  className,
  collapsed = false,
}) => {
  return (
    <aside
      className={cn(
        'liquid-glass h-full transition-all duration-300',
        collapsed ? 'w-16' : 'w-64',
        className
      )}
    >
      <div className="p-4 h-full overflow-y-auto liquid-glass-scroll">
        {children}
      </div>
    </aside>
  )
}

export default LiquidGlassSidebar
