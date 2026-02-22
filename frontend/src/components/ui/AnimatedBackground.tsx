import React from 'react'
import { cn } from '@/lib/cn'

export interface AnimatedBackgroundProps {
  className?: string
  children?: React.ReactNode
}

/**
 * AnimatedBackground Component
 *
 * Creates a dynamic mesh gradient background with floating animated orbs.
 * Part of the Liquid Glass design system for MediSync.
 *
 * Features:
 * - Layered mesh gradient base with deep blues and purples
 * - Three floating animated orbs with blur effects
 * - Smooth float animation for depth and movement
 * - Respects prefers-reduced-motion for accessibility
 *
 * @see docs/DESIGN.md Section 5.3 - Animated Background Orbs
 */
export const AnimatedBackground: React.FC<AnimatedBackgroundProps> = ({
  className,
  children,
}) => {
  return (
    <div className={cn('fixed inset-0 -z-10 overflow-hidden pointer-events-none transition-colors duration-700 bg-slate-50 dark:bg-black', className)}>
      {/* Orb 1 (Logo Blue) */}
      <div
        className="absolute w-[600px] h-[600px] rounded-full animate-float opacity-30 dark:opacity-60 mix-blend-multiply dark:mix-blend-screen"
        style={{
          background: 'radial-gradient(circle at center, rgba(39, 80, 168, 0.8) 0%, transparent 70%)',
          filter: 'blur(120px)',
          top: '10%',
          left: '5%',
          animationDuration: '25s',
        }}
      />

      {/* Orb 2 (Logo Teal) */}
      <div
        className="absolute w-[500px] h-[500px] rounded-full animate-float opacity-30 dark:opacity-60 mix-blend-multiply dark:mix-blend-screen"
        style={{
          background: 'radial-gradient(circle at center, rgba(24, 146, 157, 0.6) 0%, transparent 70%)',
          filter: 'blur(100px)',
          bottom: '15%',
          right: '5%',
          animationDuration: '30s',
          animationDelay: '-5s',
        }}
      />

      {children}
    </div>
  )
}

export default AnimatedBackground
