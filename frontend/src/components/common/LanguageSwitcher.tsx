/**
 * LanguageSwitcher Component
 *
 * Toggle button for switching between English and Arabic locales.
 * Uses i18next changeLanguage and persists preference to backend API.
 *
 * Styled as a glass pill toggle as part of the Liquid Glass design system.
 *
 * @module components/common/LanguageSwitcher
 * @version 2.0.0 - Updated for Liquid Glass design
 */
import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { usePreferences } from '../../hooks/usePreferences'
import { cn } from '@/lib/cn'

export interface LanguageSwitcherProps {
  /** Additional CSS classes */
  className?: string
  /** Compact mode shows shorter labels (default: true) */
  compact?: boolean
  /** Whether to persist to backend (default: true) */
  persistToBackend?: boolean
  /** Pill style variant (default: true, uses glass pill toggle) */
  pillStyle?: boolean
}

/**
 * LanguageSwitcher - Toggle between EN/AR locales
 *
 * Features:
 * - Toggles between English and Arabic
 * - Glass pill toggle design for modern aesthetics
 * - Accessible with aria-label
 * - RTL-aware styling
 * - Persists preference via API and i18n
 * - Loading state while saving
 */
export function LanguageSwitcher({
  className = '',
  compact = true,
  persistToBackend = true,
  pillStyle = true,
}: LanguageSwitcherProps) {
  const { t, i18n } = useTranslation()
  const { updatePreferences, isUpdating } = usePreferences()
  const [isSwitching, setIsSwitching] = useState(false)

  const currentLocale = i18n.language === 'ar' ? 'ar' : 'en'

  const setLocale = useCallback(async (locale: 'en' | 'ar') => {
    setIsSwitching(true)

    try {
      // Update i18n immediately for responsive UI
      await i18n.changeLanguage(locale)

      // Persist to backend if enabled
      if (persistToBackend) {
        await updatePreferences({ locale })
      }
    } catch (error) {
      console.error('Failed to switch language:', error)
      // Revert on error
      await i18n.changeLanguage(currentLocale)
    } finally {
      setIsSwitching(false)
    }
  }, [i18n, persistToBackend, updatePreferences, currentLocale])

  const isLoading = isSwitching || isUpdating

  // Glass pill toggle style (default)
  if (pillStyle) {
    return (
      <div
        className={cn(
          'liquid-glass flex items-center rounded-full p-1',
          isLoading && 'opacity-50 pointer-events-none',
          className
        )}
        role="radiogroup"
        aria-label={t('common.language.selectLanguage', 'Select language')}
      >
        <button
          className={cn(
            'px-3 py-1 rounded-full text-sm font-medium transition-all duration-200',
            'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500',
            currentLocale === 'en'
              ? 'bg-blue-500 text-white shadow-md'
              : 'liquid-text-secondary hover:text-slate-900 dark:hover:text-white'
          )}
          onClick={() => setLocale('en')}
          disabled={isLoading}
          role="radio"
          aria-checked={currentLocale === 'en'}
        >
          EN
        </button>
        <button
          className={cn(
            'px-3 py-1 rounded-full text-sm font-medium transition-all duration-200',
            'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500',
            currentLocale === 'ar'
              ? 'bg-blue-500 text-white shadow-md'
              : 'liquid-text-secondary hover:text-slate-900 dark:hover:text-white'
          )}
          onClick={() => setLocale('ar')}
          disabled={isLoading}
          role="radio"
          aria-checked={currentLocale === 'ar'}
        >
          ع
        </button>
      </div>
    )
  }

  // Legacy button style (for backward compatibility)
  const nextLang = currentLocale === 'en' ? 'ar' : 'en'

  return (
    <button
      onClick={() => setLocale(nextLang)}
      disabled={isLoading}
      className={cn(
        'inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg',
        'bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700',
        'text-sm font-medium text-slate-700 dark:text-slate-300',
        'transition-colors duration-200',
        'focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2',
        'disabled:opacity-50 disabled:cursor-not-allowed',
        className
      )}
      aria-label={t(
        'common.language.switchTo',
        `Switch to ${nextLang === 'ar' ? 'Arabic' : 'English'}`
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
        {nextLang === 'ar' ? 'عربي' : 'English'}
      </span>
    </button>
  )
}

export default LanguageSwitcher
