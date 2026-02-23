/**
 * AnySync Design Tokens - TypeScript/JavaScript Export
 *
 * Provides typed access to design tokens for use in React components.
 * All values reference the CSS custom properties defined in design-tokens.css.
 *
 * @example
 * ```tsx
 * import { tokens } from '@/tailwind-theme'
 *
 * <div style={{ background: tokens.gradients.brand }} />
 * <div style={{ boxShadow: tokens.shadows.floating }} />
 * ```
 *
 * @version 2.0.0
 * @lastUpdated 2026-02-23
 */

// ============================================
// COLOR TOKENS
// ============================================
export const colors = {
  // Brand colors
  brandPrimary: '#2750a8',
  brandPrimaryLight: '#3a6fc4',
  brandPrimaryDark: '#1e3d7a',

  brandSecondary: '#18929d',
  brandSecondaryLight: '#1ab5c2',
  brandSecondaryDark: '#0f6b73',

  brandNavy: '#0f172a',

  // Semantic colors
  success: '#10b981',
  successLight: '#34d399',
  warning: '#f59e0b',
  warningLight: '#fbbf24',
  error: '#ef4444',
  errorLight: '#f87171',
  info: '#3b82f6',
  infoLight: '#60a5fa',

  // WCAG 3.0 Gold
  focusGold: '#FFD700',
} as const

// ============================================
// GRADIENT TOKENS
// ============================================
export const gradients = {
  brand: 'linear-gradient(135deg, #2750a8, #18929d)',
  brandShift: 'linear-gradient(90deg, #2750a8, #18929d, #18929d)',
  brandReverse: 'linear-gradient(135deg, #18929d, #2750a8)',
  brandAnimated: 'linear-gradient(135deg, #0056d2, #00e8c6, #0056d2, #00e8c6)',
  shine: 'linear-gradient(90deg, transparent, rgba(255,255,255,0.4), transparent)',
} as const

// ============================================
// SPACING TOKENS (4px base unit)
// ============================================
export const spacing = {
  xs: '4px',
  sm: '8px',
  md: '12px',
  lg: '16px',
  xl: '20px',
  '2xl': '24px',
  '3xl': '32px',
  '4xl': '40px',
  '5xl': '48px',
  '6xl': '64px',
  '7xl': '80px',
  '8xl': '96px',

  // Component spacing
  sectionSm: '48px',
  sectionMd: '80px',
  sectionLg: '112px',
  sectionXl: '128px',

  gapSm: '16px',
  gapMd: '24px',
  gapLg: '32px',
} as const

// ============================================
// SHADOW TOKENS
// ============================================
export const shadows = {
  ambient: '0 0 0 1px rgba(0,0,0,0.04), 0 2px 4px rgba(0,0,0,0.02), 0 8px 16px rgba(0,0,0,0.04)',
  raised: '0 0 0 1px rgba(0,0,0,0.06), 0 4px 8px rgba(0,0,0,0.04), 0 16px 32px rgba(0,0,0,0.08)',
  floating: '0 0 0 1px rgba(0,0,0,0.08), 0 8px 16px rgba(0,0,0,0.06), 0 32px 64px rgba(0,0,0,0.12)',
  elevated: '0 0 0 1px rgba(255,255,255,0.3), 0 20px 60px rgba(0,0,0,0.15), 0 40px 80px rgba(0,0,0,0.10)',
  glassFloating: '0 0 0 1px rgba(255,255,255,0.4), 0 32px 80px rgba(0,0,0,0.18), 0 60px 120px rgba(0,0,0,0.14)',

  // Colored glows
  glowPrimary: '0 0 20px rgba(39, 80, 168, 0.3)',
  glowSecondary: '0 0 20px rgba(24, 146, 157, 0.3)',
  glowSuccess: '0 0 20px rgba(16, 185, 129, 0.3)',
  glowWarning: '0 0 20px rgba(245, 158, 11, 0.3)',
  glowError: '0 0 20px rgba(239, 68, 68, 0.3)',
} as const

// ============================================
// BORDER RADIUS TOKENS
// ============================================
export const radius = {
  none: '0',
  sm: '8px',
  md: '12px',
  lg: '16px',
  xl: '20px',
  '2xl': '28px',
  '3xl': '32px',
  dynamic: '40px', // iOS Dynamic Island
  full: '9999px',
} as const

// ============================================
// TYPOGRAPHY TOKENS
// ============================================
export const fontSize = {
  xs: '0.75rem',   // 12px
  sm: '0.875rem',  // 14px
  base: '1rem',    // 16px
  lg: '1.125rem',  // 18px
  xl: '1.25rem',   // 20px
  '2xl': '1.5rem', // 24px
  '3xl': '1.875rem', // 30px
  '4xl': '2.25rem',  // 36px
  '5xl': '3rem',      // 48px
} as const

export const fontWeight = {
  regular: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
} as const

// ============================================
// ANIMATION TOKENS
// ============================================
export const duration = {
  instant: '100ms',
  fast: '150ms',
  base: '250ms',
  moderate: '350ms',
  slow: '500ms',
  slower: '750ms',
} as const

export const easing = {
  liquid: 'cubic-bezier(0.4, 0, 0.2, 1)',
  liquidOut: 'cubic-bezier(0, 0, 0.2, 1)',
  liquidIn: 'cubic-bezier(0.4, 0, 1, 1)',
  elastic: 'cubic-bezier(0.68, -0.6, 0.32, 1.6)',
  bounce: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
} as const

// ============================================
// Z-INDEX LAYERS
// ============================================
export const zIndex = {
  background: -1,
  base: 0,
  raised: 10,
  dropdown: 100,
  sticky: 200,
  fixed: 300,
  modalBackdrop: 400,
  modal: 500,
  popover: 600,
  tooltip: 700,
  toast: 800,
  max: 9999,
} as const

// ============================================
// GLASSMORPHISM TOKENS
// ============================================
export const glass = {
  bgSubtle: 'rgba(255, 255, 255, 0.40)',
  bgLight: 'rgba(255, 255, 255, 0.85)',
  bgMedium: 'rgba(255, 255, 255, 0.70)',
  bgHeavy: 'rgba(255, 255, 255, 0.55)',
  bgPronounced: 'rgba(255, 255, 255, 0.90)',

  borderBright: 'rgba(255, 255, 255, 0.6)',
  borderSoft: 'rgba(255, 255, 255, 0.35)',
  borderDim: 'rgba(255, 255, 255, 0.15)',

  blur: {
    xs: 'blur(10px)',
    sm: 'blur(20px)',
    md: 'blur(45px)',
    lg: 'blur(60px)',
    xl: 'blur(80px)',
    xxl: 'blur(100px)',
  },
} as const

// ============================================
// NAMED EXPORTS FOR CONVENIENCE
// ============================================
export const tokens = {
  colors,
  gradients,
  spacing,
  shadows,
  radius,
  fontSize,
  fontWeight,
  duration,
  easing,
  zIndex,
  glass,
} as const

/**
 * Hook to get a CSS custom property value
 * @param key - The CSS variable name (e.g., '--brand-primary')
 * @returns The computed value of the CSS variable
 */
export function getTokenValue(key: string): string {
  if (typeof document === 'undefined') return ''
  return getComputedStyle(document.documentElement).getPropertyValue(key).trim()
}

/**
 * Hook to get all design tokens as CSS variable references
 * Useful for inline styles that need to respect theme changes
 */
export const cssVars = {
  // Colors
  brandPrimary: 'var(--brand-primary)',
  brandSecondary: 'var(--brand-secondary)',
  brandNavy: 'var(--brand-navy)',

  // Gradients
  gradientBrand: 'var(--gradient-brand)',
  gradientShine: 'var(--gradient-shine)',

  // Spacing
  space4: 'var(--space-4)',
  space6: 'var(--space-6)',
  space8: 'var(--space-8)',

  // Shadows
  shadowAmbient: 'var(--shadow-ambient)',
  shadowRaised: 'var(--shadow-raised)',
  shadowFloating: 'var(--shadow-floating)',

  // Radius
  radiusMd: 'var(--radius-md)',
  radiusLg: 'var(--radius-lg)',
  radiusXl: 'var(--radius-xl)',

  // Glass
  glassPronounced: 'var(--glass-bg-pronounced)',

  // Focus
  focusGold: 'var(--focus-gold)',
} as const

export default tokens
