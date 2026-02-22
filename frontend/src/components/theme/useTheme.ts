/**
 * useTheme Hook
 *
 * Custom hook for accessing and manipulating the application theme.
 * Wraps next-themes useTheme with type safety and guaranteed values.
 *
 * Features:
 * - Type-safe theme values ('light' | 'dark' | 'system')
 * - Guaranteed non-undefined theme value
 * - Simplified setTheme API
 *
 * @module components/theme/useTheme
 */
import { useTheme as useNextTheme } from 'next-themes'

type ThemeValue = 'light' | 'dark' | 'system'

/**
 * Return type for useTheme hook
 *
 * Extends next-themes useTheme return type with guaranteed values
 */
export interface UseThemeReturn {
  /** Current theme value (never undefined) */
  theme: ThemeValue
  /** Set the theme to 'light', 'dark', or 'system' */
  setTheme: (theme: ThemeValue) => void
  /** All available themes */
  themes: ThemeValue[]
  /** The actual theme in use (resolved from 'system' preference) */
  resolvedTheme: 'light' | 'dark' | undefined
}

/**
 * Theme Access Hook
 *
 * Provides access to the current theme and functions to change it.
 *
 * @example
 * ```tsx
 * function MyComponent() {
 *   const { theme, setTheme } = useTheme()
 *
 *   return (
 *     <button onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>
 *       Toggle Theme
 *     </button>
 *   )
 * }
 * ```
 *
 * @returns Theme state and setter functions
 */
export function useTheme(): UseThemeReturn {
  const theme = useNextTheme()

  return {
    // Ensure theme is never undefined by defaulting to 'system'
    theme: (theme.theme || 'system') as ThemeValue,
    // Simplified setTheme with type safety
    setTheme: (themeValue: ThemeValue) => theme.setTheme(themeValue),
    themes: theme.themes as ThemeValue[],
    resolvedTheme: theme.resolvedTheme as 'light' | 'dark' | undefined,
  }
}

export default useTheme
