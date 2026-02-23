/**
 * CTA Class Names Module
 *
 * Central CTA (Call-to-Action) button class names and helper functions.
 * This module provides a single source of truth for all CTA styling
 * across the AnySync frontend application.
 *
 * USAGE:
 * - For most CTAs: Use CTA_CLASSES constants directly with CSS classes
 * - For special cases: Use helper functions that return theme-aware strings
 *
 * @module lib/cta-classes
 * @version 1.0.0
 * @since 2026-02-23
 */

/**
 * CTA CSS class name constants
 *
 * Use these constants for components that can rely on CSS-driven theming
 * (i.e., theme is controlled by the root .dark class). The corresponding
 * CSS classes in liquid-glass.css handle light/dark mode automatically.
 *
 * @example
 * ```tsx
 * import { CTA_CLASSES } from '@/lib/cta-classes'
 *
 * <button className={CTA_CLASSES.HERO}>Get Started</button>
 * ```
 */
export const CTA_CLASSES = {
  /**
   * Hero carousel CTA - largest, most prominent
   * Features: pronounced glass, brand gradient, shimmer on hover
   * Used in: HeroCarousel.tsx
   */
  HERO: 'cta-hero',

  /**
   * Landing page primary CTA
   * Features: white glass in light mode, blue-cyan gradient in dark mode
   * Used in: FinalCTA.tsx (via helper for theme control)
   */
  LANDING_PRIMARY: 'cta-landing-primary',

  /**
   * Landing page secondary CTA
   * Features: bordered glass style, transparent background
   * Used in: FinalCTA.tsx (via helper for theme control)
   */
  LANDING_SECONDARY: 'cta-landing-secondary',

  /**
   * Header/chat primary CTA - compact variant
   * Features: brand gradient, compact size, shimmer on hover
   * Used in: ChatHeader.tsx, ChatInput.tsx
   */
  HEADER_PRIMARY: 'cta-header-primary',

  /**
   * Legacy alias for backward compatibility
   * Maps to the same class as HEADER_PRIMARY
   * @deprecated Use HEADER_PRIMARY instead
   */
  BUTTON_PRIMARY: 'cta-header-primary',
} as const

/**
 * CTA class name type for type safety
 */
export type CTAClassName = typeof CTA_CLASSES[keyof typeof CTA_CLASSES]

/**
 * Get FinalCTA primary button classes with theme-specific styling
 *
 * @deprecated Use CTA_CLASSES.LANDING_PRIMARY directly instead.
 * The CSS classes handle theming automatically via .dark selector.
 *
 * @param _isDark - Unused parameter (kept for backward compatibility)
 * @returns The CSS class name for landing primary CTA
 *
 * @example
 * ```tsx
 * // OLD way (deprecated):
 * import { getFinalCTAPrimaryClasses } from '@/lib/cta-classes'
 * <button className={getFinalCTAPrimaryClasses(isDark)}>Get Started</button>
 *
 * // NEW way (recommended):
 * import { CTA_CLASSES } from '@/lib/cta-classes'
 * <button className={CTA_CLASSES.LANDING_PRIMARY}>Get Started</button>
 * ```
 */
export function getFinalCTAPrimaryClasses(_isDark: boolean): string {
  return CTA_CLASSES.LANDING_PRIMARY
}

/**
 * Get FinalCTA secondary button classes with theme-specific styling
 *
 * @deprecated Use CTA_CLASSES.LANDING_SECONDARY directly instead.
 * The CSS classes handle theming automatically via .dark selector.
 *
 * @param _isDark - Unused parameter (kept for backward compatibility)
 * @returns The CSS class name for landing secondary CTA
 *
 * @example
 * ```tsx
 * // OLD way (deprecated):
 * import { getFinalCTASecondaryClasses } from '@/lib/cta-classes'
 * <button className={getFinalCTASecondaryClasses(isDark)}>Book Demo</button>
 *
 * // NEW way (recommended):
 * import { CTA_CLASSES } from '@/lib/cta-classes'
 * <button className={CTA_CLASSES.LANDING_SECONDARY}>Book Demo</button>
 * ```
 */
export function getFinalCTASecondaryClasses(_isDark: boolean): string {
  return CTA_CLASSES.LANDING_SECONDARY
}

/**
 * Get chat input send button classes
 *
 * Special case for ChatInput.tsx where the button needs a disabled state
 * and different styling based on whether input is valid.
 *
 * @param isEnabled - Whether the button should be enabled
 * @returns Complete className string for the send button
 *
 * @example
 * ```tsx
 * import { getChatSendButtonClasses } from '@/lib/cta-classes'
 *
 * <button className={getChatSendButtonClasses(!!input.trim())}>
 *   <SendIcon />
 * </button>
 * ```
 */
export function getChatSendButtonClasses(isEnabled: boolean): string {
  const base = [
    // Layout
    'shrink-0',
    'flex',
    'items-center',
    'justify-center',
    // Size
    'min-h-[44px]',
    'min-w-[44px]',
    'p-3',
    // Shape
    'rounded-xl',
    // Transitions
    'transition-all',
    'duration-300',
    // Focus
    'liquid-focus-gold',
  ].join(' ')

  if (isEnabled) {
    return [
      base,
      // Enabled state - primary CTA style
      'cta-header-primary',
    ].join(' ')
  }

  // Disabled state - gray, no interaction
  return [
    base,
    'bg-slate-200/80',
    'dark:bg-slate-700/80',
    'text-slate-400',
    'dark:text-slate-500',
    'cursor-not-allowed',
  ].join(' ')
}

/**
 * Type guard to check if a string is a valid CTA class name
 *
 * @param className - The class name to check
 * @returns True if the class name is a valid CTA constant
 *
 * @example
 * ```tsx
 * import { isCTAClass } from '@/lib/cta-classes'
 *
 * if (isCTAClass(someClass)) {
 *   // Use the class safely
 * }
 * ```
 */
export function isCTAClass(className: string): className is CTAClassName {
  return Object.values(CTA_CLASSES).includes(className as CTAClassName)
}

/**
 * Get all available CTA class names
 *
 * @returns Array of all CTA class name constants
 *
 * @example
 * ```tsx
 * import { getAllCTAClasses } from '@/lib/cta-classes'
 *
 * const ctaClasses = getAllCTAClasses()
 * // ['cta-hero', 'cta-landing-primary', 'cta-landing-secondary', 'cta-header-primary']
 * ```
 */
export function getAllCTAClasses(): readonly CTAClassName[] {
  return Object.values(CTA_CLASSES)
}
