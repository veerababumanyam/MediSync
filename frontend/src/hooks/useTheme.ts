import { useState, useEffect, useCallback } from 'react'

/**
 * Hook for managing dark/light mode
 *
 * Persists preference to localStorage.
 * Falls back to system preference on first visit.
 */
export function useTheme() {
    const [isDark, setIsDark] = useState<boolean>(() => {
        // Check localStorage first
        const stored = localStorage.getItem('medisync-theme')
        if (stored !== null) {
            return stored === 'dark'
        }
        // Force dark mode as the default for new users
        return true
    })

    // Apply the theme class to <html>
    useEffect(() => {
        const root = document.documentElement
        if (isDark) {
            root.classList.add('dark')
        } else {
            root.classList.remove('dark')
        }
        localStorage.setItem('medisync-theme', isDark ? 'dark' : 'light')
    }, [isDark])

    const toggleTheme = useCallback(() => {
        setIsDark((prev) => !prev)
    }, [])

    return { isDark, toggleTheme }
}
