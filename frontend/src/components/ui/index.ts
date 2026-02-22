/**
 * UI Components Barrel Export
 *
 * Exports all reusable UI components for easy importing.
 * Includes both legacy components and the new Liquid Glass design system.
 *
 * @example
 * ```tsx
 * import { GlassCard, ThemeToggle, LoadingSkeleton } from '@/components/ui'
 * import { LiquidGlassCard, LiquidGlassButton } from '@/components/ui'
 * import { AnimatedBackground } from '@/components/ui'
 * ```
 *
 * @module components/ui
 * @version 3.0.0 - Complete Liquid Glass design system with all components
 */

// ============================================
// LEGACY COMPONENTS (for backward compatibility)
// ============================================

export { GlassCard } from './GlassCard'
export type { GlassCardProps } from './GlassCard'

export { ThemeToggle } from './ThemeToggle'
export type { ThemeToggleProps } from './ThemeToggle'

export { LoadingSkeleton } from './LoadingSkeleton'
export type { LoadingSkeletonProps } from './LoadingSkeleton'

// ============================================
// LIQUID GLASS COMPONENTS (iOS-inspired)
// ============================================

/**
 * Animated Background
 *
 * Dynamic mesh gradient background with floating animated orbs.
 */
export { AnimatedBackground } from './AnimatedBackground'
export type { AnimatedBackgroundProps } from './AnimatedBackground'

/**
 * Liquid Glass Cards
 *
 * Premium glassmorphic cards with liquid animations.
 * Use these for new features requiring iOS-style aesthetics.
 */
export {
  LiquidGlassCard,
  GlassCard as GlassCardLite,
  GlassHeader,
  GlassModal,
  GlassInteractiveCard,
  GlassBrandCard,
  GlassBlueCard,
  GlassTealCard,
  GlassGreenCard,
  type LiquidGlassCardProps,
} from './LiquidGlassCard'

/**
 * Liquid Glass Buttons
 *
 * Premium glassmorphic buttons with liquid hover effects.
 */
export {
  LiquidGlassButton,
  ButtonPrimary,
  ButtonSecondary,
  ButtonGhost,
  ButtonDanger,
  IconButton,
  type LiquidGlassButtonProps,
} from './LiquidGlassButton'

/**
 * Liquid Glass Inputs
 *
 * Premium glassmorphic inputs with liquid focus states.
 */
export {
  LiquidGlassInput,
  LiquidGlassTextarea,
  LiquidGlassSearch,
  type LiquidGlassInputProps,
  type LiquidGlassTextareaProps,
  type LiquidGlassSearchProps,
} from './LiquidGlassInput'

/**
 * Liquid Glass Modal
 *
 * Premium glassmorphic modal with backdrop blur and animations.
 */
export {
  LiquidGlassModal,
  type LiquidGlassModalProps,
} from './LiquidGlassModal'

/**
 * Liquid Glass Navbar
 *
 * Premium glassmorphic navigation bar with sticky positioning.
 */
export {
  LiquidGlassNavbar,
  type LiquidGlassNavbarProps,
} from './LiquidGlassNavbar'

/**
 * Liquid Glass Sidebar
 *
 * Premium glassmorphic sidebar with collapsible functionality.
 */
export {
  LiquidGlassSidebar,
  type LiquidGlassSidebarProps,
} from './LiquidGlassSidebar'

/**
 * Liquid Glass Badge
 *
 * Premium glassmorphic badges with color variants.
 */
export {
  LiquidGlassBadge,
  type LiquidGlassBadgeProps,
} from './LiquidGlassBadge'

/**
 * Liquid Glass Toast
 *
 * Premium glassmorphic toast notifications with animations.
 */
export {
  LiquidGlassToast,
  LiquidGlassToastContainer,
  useLiquidGlassToast,
  type LiquidGlassToastProps,
  type ToastOptions,
} from './LiquidGlassToast'
