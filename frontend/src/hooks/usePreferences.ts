/**
 * usePreferences Hook
 *
 * React hook for managing user preferences with:
 * - State management for all preference fields
 * - API integration for persistence
 * - Loading and error states
 * - Optimistic updates
 *
 * @module hooks/usePreferences
 */
import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  preferencesApi,
  APIError,
} from '../services/api'
import type { UserPreferences } from '../services/api'

/**
 * Type for updating preferences (excludes auto-generated fields)
 */
export type PreferenceUpdates = Partial<
  Omit<UserPreferences, 'id' | 'userId' | 'createdAt' | 'updatedAt'>
>

/**
 * Storage key for offline caching
 */
const PREFERENCES_STORAGE_KEY = 'medisync-preferences-cache'

/**
 * Return type for usePreferences hook
 */
export interface UsePreferencesReturn {
  /** Current preferences (null if not loaded) */
  preferences: UserPreferences | null
  /** Whether preferences are currently loading */
  isLoading: boolean
  /** Whether an update is in progress */
  isUpdating: boolean
  /** Any error that occurred during operations */
  error: string | null
  /** Update preferences */
  updatePreferences: (updates: PreferenceUpdates) => Promise<void>
  /** Reload preferences from server */
  reload: () => Promise<void>
  /** Reset preferences to defaults */
  resetToDefaults: () => Promise<void>
}

/**
 * Default preferences for new users
 */
const DEFAULT_PREFERENCES: Omit<UserPreferences, 'id' | 'userId' | 'createdAt' | 'updatedAt'> = {
  locale: 'en',
  numeralSystem: 'western',
  calendarSystem: 'gregorian',
  reportLanguage: 'en',
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
}

/**
 * Get cached preferences from localStorage
 */
function getCachedPreferences(): UserPreferences | null {
  try {
    const cached = localStorage.getItem(PREFERENCES_STORAGE_KEY)
    if (cached) {
      return JSON.parse(cached)
    }
  } catch (err) {
    console.debug('Failed to parse cached preferences:', err)
  }
  return null
}

/**
 * Cache preferences to localStorage
 */
function cachePreferences(preferences: UserPreferences): void {
  try {
    localStorage.setItem(PREFERENCES_STORAGE_KEY, JSON.stringify(preferences))
  } catch (err) {
    console.debug('Failed to cache preferences:', err)
  }
}

/**
 * Hook for managing user preferences
 *
 * @example
 * ```tsx
 * function PreferencesForm() {
 *   const { preferences, isLoading, error, updatePreferences } = usePreferences()
 *
 *   if (isLoading) return <Spinner />
 *   if (error) return <Error message={error} />
 *
 *   return (
 *     <form onSubmit={(e) => {
 *       e.preventDefault()
 *       updatePreferences({ locale: 'ar' })
 *     }}>
 *       // form fields
 *     </form>
 *   )
 * }
 * ```
 */
export function usePreferences(): UsePreferencesReturn {
  const { i18n } = useTranslation()
  const [preferences, setPreferences] = useState<UserPreferences | null>(
    () => getCachedPreferences()
  )
  const [isLoading, setIsLoading] = useState(!getCachedPreferences())
  const [isUpdating, setIsUpdating] = useState(false)
  const [error, setError] = useState<string | null>(null)

  /**
   * Load preferences on mount
   */
  useEffect(() => {
    let mounted = true

    async function loadPreferences() {
      setIsLoading(true)
      setError(null)

      try {
        const data = await preferencesApi.get()
        if (mounted) {
          setPreferences(data)
          cachePreferences(data)

          // Sync i18n with server preference
          if (data.locale && data.locale !== i18n.language) {
            await i18n.changeLanguage(data.locale)
          }
        }
      } catch (err) {
        if (mounted) {
          // If 404, preferences don't exist yet - use defaults
          if (err instanceof APIError && err.status === 404) {
            const defaults: UserPreferences = {
              id: 'temp',
              userId: 'temp',
              ...DEFAULT_PREFERENCES,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            }
            setPreferences(defaults)
            cachePreferences(defaults)
          } else {
            const message = err instanceof Error
              ? err.message
              : 'Failed to load preferences'
            setError(message)
            console.error('Failed to load preferences:', err)
          }
        }
      } finally {
        if (mounted) {
          setIsLoading(false)
        }
      }
    }

    loadPreferences()

    return () => {
      mounted = false
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []) // Only run on mount - i18n changeLanguage is called internally

  /**
   * Update preferences
   */
  const updatePreferences = useCallback(async (updates: PreferenceUpdates) => {
    setError(null)
    setIsUpdating(true)

    // Optimistic update
    const previousPreferences = preferences
    const updatedPreferences = preferences
      ? { ...preferences, ...updates, updatedAt: new Date().toISOString() }
      : null

    if (updatedPreferences) {
      setPreferences(updatedPreferences)
      cachePreferences(updatedPreferences)
    }

    try {
      const result = await preferencesApi.update(updates)
      setPreferences(result)
      cachePreferences(result)

      // Sync i18n if locale changed
      if (updates.locale && updates.locale !== i18n.language) {
        await i18n.changeLanguage(updates.locale)
      }
    } catch (err) {
      // Revert on error
      if (previousPreferences) {
        setPreferences(previousPreferences)
        cachePreferences(previousPreferences)
      }

      const message = err instanceof Error
        ? err.message
        : 'Failed to update preferences'
      setError(message)
      console.error('Failed to update preferences:', err)
      throw err
    } finally {
      setIsUpdating(false)
    }
  }, [preferences, i18n])

  /**
   * Reload preferences from server
   */
  const reload = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await preferencesApi.get()
      setPreferences(data)
      cachePreferences(data)
    } catch (err) {
      const message = err instanceof Error
        ? err.message
        : 'Failed to reload preferences'
      setError(message)
      console.error('Failed to reload preferences:', err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Reset preferences to defaults
   */
  const resetToDefaults = useCallback(async () => {
    await updatePreferences(DEFAULT_PREFERENCES)
  }, [updatePreferences])

  return {
    preferences,
    isLoading,
    isUpdating,
    error,
    updatePreferences,
    reload,
    resetToDefaults,
  }
}

export default usePreferences
