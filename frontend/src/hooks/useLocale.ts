/**
 * useLocale Hook
 *
 * React hook for managing locale with:
 * - Current locale access from i18n
 * - Locale setting via API
 * - LocalStorage persistence
 * - Document direction updates (RTL/LTR)
 *
 * @module hooks/useLocale
 */
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { preferencesApi } from '../services/api'

export type Locale = 'en' | 'ar'

export interface UseLocaleReturn {
  /** Current locale */
  locale: Locale
  /** Whether the current locale is RTL */
  isRTL: boolean
  /** Set locale (updates i18n, API, and localStorage) */
  setLocale: (locale: Locale) => Promise<void>
  /** Toggle between EN and AR */
  toggleLocale: () => Promise<void>
  /** Whether a locale update is in progress */
  isLoading: boolean
  /** Any error that occurred during locale update */
  error: string | null
}

/**
 * Hook for managing locale with i18n, API, and persistence
 */
export function useLocale(): UseLocaleReturn {
  const { i18n } = useTranslation()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Get current locale from i18n
  const locale = (i18n.language === 'ar' ? 'ar' : 'en') as Locale
  const isRTL = locale === 'ar'

  // Load preferences on mount
  useEffect(() => {
    async function loadPreferences() {
      try {
        const prefs = await preferencesApi.get()

        // Sync i18n with server preference if different
        if (prefs.locale && prefs.locale !== i18n.language) {
          await i18n.changeLanguage(prefs.locale)
        }
      } catch (err) {
        // Preferences might not exist yet - that's ok
        console.debug('Could not load preferences:', err)
      }
    }

    loadPreferences()
  }, [i18n])

  // Update document direction when locale changes
  useEffect(() => {
    document.documentElement.dir = isRTL ? 'rtl' : 'ltr'
    document.documentElement.lang = locale

    // Update body class for RTL-specific styling
    document.body.classList.toggle('rtl', isRTL)
    document.body.classList.toggle('ltr', !isRTL)
  }, [locale, isRTL])

  /**
   * Set locale and persist to API
   */
  const setLocale = useCallback(async (newLocale: Locale) => {
    setIsLoading(true)
    setError(null)

    try {
      // Update i18n immediately for responsive UI
      await i18n.changeLanguage(newLocale)

      // Persist to API
      try {
        await preferencesApi.update({ locale: newLocale })
      } catch (apiErr) {
        // API update failed, but i18n is already updated
        // The preference will sync on next load
        console.warn('Failed to persist locale preference:', apiErr)
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to update locale'
      setError(message)
      console.error('Locale update error:', err)
    } finally {
      setIsLoading(false)
    }
  }, [i18n])

  /**
   * Toggle between EN and AR
   */
  const toggleLocale = useCallback(async () => {
    const newLocale = locale === 'en' ? 'ar' : 'en'
    await setLocale(newLocale)
  }, [locale, setLocale])

  return {
    locale,
    isRTL,
    setLocale,
    toggleLocale,
    isLoading,
    error,
  }
}

export default useLocale
