/**
 * ThemeToggle Component
 *
 * An animated toggle for switching between light and dark themes.
 * Styled as a glass icon button as part of the Liquid Glass design system.
 *
 * Features:
 * - Glass icon button design
 * - Animated sun/moon icons
 * - WCAG 2.2 AA accessible with proper ARIA attributes
 * - Keyboard navigation support
 * - Focus-visible indicators
 * - System theme detection support
 *
 * @module components/ui/ThemeToggle
 * @version 2.0.0 - Updated for Liquid Glass design
 */
import React from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useTheme } from '@/components/theme'
import { IconButton } from './LiquidGlassButton'
import { cn } from '@/lib/cn'

const sendDebugLog = (payload: {
  runId: string
  hypothesisId: string
  location: string
  message: string
  data: Record<string, unknown>
}) => {
  // #region agent log
  // Disabled to prevent ERR_CONNECTION_REFUSED in environments where the debug server is not running
  /*
  fetch('http://127.0.0.1:7583/ingest/d551f5e1-ee3b-4b67-a81f-cb9e9d8d73e8', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Debug-Session-Id': 'aeb0a5',
    },
    body: JSON.stringify({
      sessionId: 'aeb0a5',
      timestamp: Date.now(),
      ...payload,
    }),
  }).catch(() => { })
  */
  if (import.meta.env.DEV && false) {
    console.debug('[DebugLog]', payload);
  }
  // #endregion
}

/**
 * Props for ThemeToggle component
 */
export interface ThemeToggleProps {
  /** Additional CSS classes */
  className?: string
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
  /** Whether to show labels */
  showLabels?: boolean
  /** Use legacy toggle switch style instead of icon button */
  useSwitchStyle?: boolean
}

/**
 * Sun Icon Component
 */
const SunIcon = () => (
  <svg
    className="w-5 h-5"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    strokeWidth={2}
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
    />
  </svg>
)

/**
 * Moon Icon Component
 */
const MoonIcon = () => (
  <svg
    className="w-5 h-5"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
    strokeWidth={2}
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
    />
  </svg>
)

/**
 * ThemeToggle Component
 *
 * A glass icon button for toggling between light and dark themes.
 *
 * @example
 * ```tsx
 * // Basic usage - glass icon button
 * <ThemeToggle />
 *
 * // With custom size
 * <ThemeToggle size="lg" />
 *
 * // Legacy switch style
 * <ThemeToggle useSwitchStyle />
 * ```
 */
export const ThemeToggle = React.forwardRef<HTMLButtonElement, ThemeToggleProps>(
  ({ className, size = 'md', showLabels = false, useSwitchStyle = false }, ref) => {
    const { theme, resolvedTheme, setTheme } = useTheme()
    const [mounted, setMounted] = React.useState(false)

    // Prevent hydration mismatch
    React.useEffect(() => {
      setMounted(true)
    }, [])

    const effectiveTheme = theme === 'system' ? resolvedTheme : theme
    const isDark = effectiveTheme === 'dark'

    React.useEffect(() => {
      sendDebugLog({
        runId: 'initial',
        hypothesisId: 'H5',
        location: 'frontend/src/components/ui/ThemeToggle.tsx:stateSnapshot',
        message: 'theme toggle state snapshot',
        data: {
          mounted,
          theme,
          resolvedTheme,
          effectiveTheme: effectiveTheme ?? null,
          isDark,
          useSwitchStyle,
          size,
        },
      })
    }, [effectiveTheme, isDark, mounted, resolvedTheme, size, theme, useSwitchStyle])

    React.useEffect(() => {
      const button = document.querySelector('[data-debug-theme-toggle="1"]')
      if (!(button instanceof HTMLElement)) {
        sendDebugLog({
          runId: 'initial',
          hypothesisId: 'H6',
          location: 'frontend/src/components/ui/ThemeToggle.tsx:domProbe',
          message: 'theme toggle button not found in DOM probe',
          data: {
            mounted,
            useSwitchStyle,
          },
        })
        return
      }
      const style = window.getComputedStyle(button)
      const rect = button.getBoundingClientRect()
      sendDebugLog({
        runId: 'initial',
        hypothesisId: 'H6',
        location: 'frontend/src/components/ui/ThemeToggle.tsx:domProbe',
        message: 'theme toggle button DOM probe',
        data: {
          mounted,
          useSwitchStyle,
          tagName: button.tagName,
          className: button.className,
          title: button.getAttribute('title'),
          ariaLabel: button.getAttribute('aria-label'),
          pointerEvents: style.pointerEvents,
          opacity: style.opacity,
          visibility: style.visibility,
          color: style.color,
          backgroundColor: style.backgroundColor,
          rect: {
            left: rect.left,
            top: rect.top,
            width: rect.width,
            height: rect.height,
          },
        },
      })
    }, [mounted, isDark, useSwitchStyle])

    if (!mounted) {
      // Return placeholder during SSR to prevent hydration mismatch
      return (
        <div
          className={cn(
            'rounded-full bg-slate-200 dark:bg-slate-700 animate-pulse',
            useSwitchStyle
              ? size === 'sm' ? 'w-11 h-6' : size === 'lg' ? 'w-16 h-9' : 'w-14 h-8'
              : size === 'sm' ? 'w-8 h-8' : size === 'lg' ? 'w-12 h-12' : 'w-10 h-10',
            className
          )}
        />
      )
    }

    // Handle theme toggle
    const toggleTheme = () => {
      const targetTheme = isDark ? 'light' : 'dark'
      sendDebugLog({
        runId: 'initial',
        hypothesisId: 'H7',
        location: 'frontend/src/components/ui/ThemeToggle.tsx:toggleTheme',
        message: 'theme toggle click handler fired',
        data: {
          theme,
          resolvedTheme,
          effectiveTheme: effectiveTheme ?? null,
          isDark,
          targetTheme,
          useSwitchStyle,
        },
      })
      setTheme(targetTheme)
    }

    // Glass icon button style (default for Liquid Glass design)
    if (!useSwitchStyle) {
      return (
        <div className={cn('flex items-center gap-2 flex-wrap', className)}>
          <IconButton
            ref={ref}
            type="button"
            data-debug-theme-toggle="1"
            icon={
              <AnimatePresence mode="wait" initial={false}>
                <motion.div
                  key={isDark ? 'moon' : 'sun'}
                  initial={{ scale: 0, rotate: -90, opacity: 0 }}
                  animate={{ scale: 1, rotate: 0, opacity: 1 }}
                  exit={{ scale: 0, rotate: 90, opacity: 0 }}
                  transition={{ duration: 0.2 }}
                  className={cn(
                    'relative z-10',
                    isDark ? 'text-slate-100' : 'text-slate-700'
                  )}
                >
                  {isDark ? <MoonIcon /> : <SunIcon />}
                </motion.div>
              </AnimatePresence>
            }
            onClick={toggleTheme}
            size={size}
            aria-label={`Switch to ${isDark ? 'light' : 'dark'} mode`}
            title={`Switch to ${isDark ? 'light' : 'dark'} mode`}
          />

          {/* Optional labels */}
          {showLabels && (
            <span className="text-sm text-slate-600 dark:text-slate-400">
              {isDark ? 'Dark' : 'Light'}
            </span>
          )}
        </div>
      )
    }

    // Legacy switch style (for backward compatibility)
    const sizeConfig = {
      sm: { width: 'w-11', height: 'h-6', thumb: 'w-4 h-4', translate: 16 },
      md: { width: 'w-14', height: 'h-8', thumb: 'w-6 h-6', translate: 24 },
      lg: { width: 'w-16', height: 'h-9', thumb: 'w-7 h-7', translate: 28 },
    }

    const config = sizeConfig[size]

    return (
      <div className={cn('flex items-center gap-2 flex-wrap', className)}>
        <motion.button
          ref={ref}
          data-debug-theme-toggle="1"
          className={cn(
            'relative rounded-full p-1 transition-colors duration-300',
            'bg-slate-200 dark:bg-slate-700',
            'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500',
            config.width,
            config.height
          )}
          onClick={toggleTheme}
          aria-label={`Switch to ${isDark ? 'light' : 'dark'} mode`}
          role="switch"
          aria-checked={isDark}
          type="button"
        >
          {/* Animated toggle thumb */}
          <motion.span
            className={cn(
              'block rounded-full bg-white shadow-md',
              config.thumb
            )}
            animate={{ x: isDark ? config.translate : 0 }}
            transition={{
              type: 'spring',
              stiffness: 500,
              damping: 30,
            }}
            aria-hidden="true"
          >
            {/* Sun icon for light mode */}
            <motion.svg
              className="absolute inset-0 m-auto w-4 h-4 text-amber-500"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              initial={{ scale: 0, rotate: -90 }}
              animate={{ scale: isDark ? 0 : 1, rotate: isDark ? -90 : 0 }}
              transition={{ duration: 0.2 }}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
              />
            </motion.svg>

            {/* Moon icon for dark mode */}
            <motion.svg
              className="absolute inset-0 m-auto w-4 h-4 text-slate-700"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              initial={{ scale: 0, rotate: 90 }}
              animate={{ scale: isDark ? 1 : 0, rotate: isDark ? 0 : 90 }}
              transition={{ duration: 0.2 }}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
              />
            </motion.svg>
          </motion.span>
        </motion.button>

        {/* Optional labels */}
        {showLabels && (
          <span className="text-sm text-slate-600 dark:text-slate-400">
            {isDark ? 'Dark' : 'Light'}
          </span>
        )}
      </div>
    )
  }
)

ThemeToggle.displayName = 'ThemeToggle'

export default ThemeToggle
