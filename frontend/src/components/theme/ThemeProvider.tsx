/**
 * Theme Provider Component
 *
 * Wraps the application with next-themes ThemeProvider to enable
 * manual dark/light mode toggling with system preference detection
 * and localStorage persistence.
 *
 * Features:
 * - Prevents flash of wrong theme (FOUC)
 * - Automatically syncs with system preference
 * - Persists user choice to localStorage
 * - SSR-safe hydration
 *
 * @module components/theme/ThemeProvider
 */
import React from 'react'
import { ThemeProvider as NextThemesProvider, type ThemeProviderProps } from 'next-themes'

/**
 * Props for ThemeProvider component
 *
 * Extends next-themes ThemeProviderProps with any additional props
 */
export interface MediSyncThemeProviderProps extends ThemeProviderProps {
  /** Application content */
  children: React.ReactNode
}

/**
 * Theme Provider Wrapper
 *
 * Must be rendered at the root of the application to enable theme
 * switching throughout the component tree.
 *
 * @example
 * ```tsx
 * <ThemeProvider defaultTheme="system" attribute="data-theme">
 *   <App />
 * </ThemeProvider>
 * ```
 */
export function ThemeProvider({ children, ...props }: MediSyncThemeProviderProps) {
  return (
    <NextThemesProvider
      attribute="data-theme"
      defaultTheme="system"
      enableSystem={true}
      disableTransitionOnChange={false}
      {...props}
    >
      {children}
    </NextThemesProvider>
  )
}

export default ThemeProvider
