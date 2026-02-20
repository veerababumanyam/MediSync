/**
 * LanguageSwitcher Component
 *
 * Toggle button for switching between English and Arabic locales.
 * Uses i18next changeLanguage and persists preference to backend API.
 *
 * @module components/common/LanguageSwitcher
 */
import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { usePreferences } from '../../hooks/usePreferences'

export interface LanguageSwitcherProps {
  /** Additional CSS classes */
  className?: string
  /** Compact mode shows shorter labels */
  compact?: boolean
  /** Whether to persist to backend (default: true) */
  persistToBackend?: boolean
}

/**
 * Language configuration
 */
const languages = {
  en: {
    code: 'en',
    nativeLabel: 'English',
    direction: 'ltr',
  },
  ar: {
    code: 'ar',
    nativeLabel: 'عربي',
    direction: 'rtl',
  },
} as const

/**
 * LanguageSwitcher - Toggle between EN/AR locales
 *
 * Features:
 * - Toggles between English and Arabic
 * - Displays native language labels
 * - Accessible with aria-label
 * - RTL-aware styling
 * - Persists preference via API and i18n
 * - Loading state while saving
 */
export function LanguageSwitcher({
  className = '',
  compact = false,
  persistToBackend = true,
}: LanguageSwitcherProps) {
  const { t, i18n } = useTranslation()
  const { updatePreferences, isUpdating } = usePreferences()
  const [isSwitching, setIsSwitching] = useState(false)

  const currentLang = i18n.language === 'ar' ? 'ar' : 'en'
  const nextLang = currentLang === 'en' ? 'ar' : 'en'
  const nextLangConfig = languages[nextLang]

  const toggleLanguage = useCallback(async () => {
    setIsSwitching(true)

    try {
      // Update i18n immediately for responsive UI
      await i18n.changeLanguage(nextLang)

      // Persist to backend if enabled
      if (persistToBackend) {
        await updatePreferences({ locale: nextLang })
      }
    } catch (error) {
      console.error('Failed to switch language:', error)
      // Revert on error
      await i18n.changeLanguage(currentLang)
    } finally {
      setIsSwitching(false)
    }
  }, [nextLang, i18n, persistToBackend, updatePreferences, currentLang])

  const isLoading = isSwitching || isUpdating

  return (
    <button
      onClick={toggleLanguage}
      disabled={isLoading}
      className={`
        inline-flex
        items-center
        justify-center
        gap-2
        px-4
        py-2
        rounded-lg
        bg-slate-100
        hover:bg-slate-200
        dark:bg-slate-800
        dark:hover:bg-slate-700
        text-sm
        font-medium
        text-slate-700
        dark:text-slate-300
        transition-colors
        duration-200
        focus:outline-none
        focus-visible:ring-2
        focus-visible:ring-blue-500
        focus-visible:ring-offset-2
        disabled:opacity-50
        disabled:cursor-not-allowed
        ${className}
      `}
      aria-label={t(
        'common.language.switchTo',
        `Switch to ${nextLangConfig.nativeLabel}`
      )}
      title={t(
        'common.language.switchTo',
        `Switch to ${nextLangConfig.nativeLabel}`
      )}
    >
      {/* Globe icon or spinner */}
      {isLoading ? (
        <svg
          className="w-4 h-4 animate-spin"
          fill="none"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      ) : (
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
          />
        </svg>
      )}

      {/* Language label */}
      <span className={compact ? 'sr-only' : ''}>
        {nextLangConfig.nativeLabel}
      </span>
    </button>
  )
}

export default LanguageSwitcher
