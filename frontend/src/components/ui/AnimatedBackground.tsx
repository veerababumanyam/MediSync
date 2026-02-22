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
    <div className={cn('fixed inset-0 -z-10 overflow-hidden pointer-events-none', className)}>
      {/* Mesh gradient base */}
      <div
        className="absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse 80% 60% at 10% 20%, rgba(88, 86, 214, 0.4) 0%, transparent 60%),
            radial-gradient(ellipse 60% 80% at 80% 80%, rgba(0, 122, 255, 0.3) 0%, transparent 60%),
            radial-gradient(ellipse 50% 50% at 50% 50%, rgba(175, 82, 222, 0.15) 0%, transparent 50%),
            #0A0A1A
          `,
        }}
      />

      {/* Floating orbs - Reduced size and opacity for subtlety */}
      <div
        className="absolute w-[280px] h-[280px] rounded-full animate-float opacity-20"
        style={{
          background: 'radial-gradient(circle, rgba(0, 122, 255, 0.30) 0%, transparent 70%)',
          filter: 'blur(100px)',
          top: '10%',
          left: '5%',
          animationDuration: '20s',
        }}
      />
      <div
        className="absolute w-[220px] h-[220px] rounded-full animate-float opacity-18"
        style={{
          background: 'radial-gradient(circle, rgba(175, 82, 222, 0.25) 0%, transparent 70%)',
          filter: 'blur(100px)',
          bottom: '20%',
          right: '10%',
          animationDuration: '25s',
          animationDelay: '-5s',
        }}
      />
      <div
        className="absolute w-[180px] h-[180px] rounded-full animate-float opacity-15"
        style={{
          background: 'radial-gradient(circle, rgba(255, 45, 85, 0.15) 0%, transparent 70%)',
          filter: 'blur(100px)',
          top: '50%',
          left: '50%',
          animationDuration: '18s',
          animationDelay: '-10s',
        }}
      />

      {children}
    </div>
  )
}

export default AnimatedBackground
