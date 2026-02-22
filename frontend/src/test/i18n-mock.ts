import { vi } from 'vitest'

/**
 * Production-ready i18n mock with proper interpolation support
 *
 * This utility provides consistent i18n mocking across all test files.
 * It properly handles interpolation patterns like {{year}}, {{count}}, etc.
 *
 * @example
 * // In your test file:
 * import { mockTranslation, createMockI18n } from '../test/i18n-mock'
 *
 * vi.mock('react-i18next', () => ({
 *   useTranslation: () => ({
 *     t: mockTranslation,
 *     i18n: createMockI18n(),
 *   }),
 * }))
 *
 * @example
 * // With custom language:
 * vi.mock('react-i18next', () => ({
 *   useTranslation: () => ({
 *     t: mockTranslation,
 *     i18n: createMockI18n('ar'),
 *   }),
 * }))
 */

/**
 * Mock translation function that handles:
 * - Simple key lookups: t('key') → 'key'
 * - Default values: t('key', 'default') → 'default'
 * - Interpolation: t('key', { year: 2026 }) → 'value with 2026' (if key in translations)
 * - Interpolation with defaults: t('key', { defaultValue: 'Hi {{name}}', name: 'John' }) → 'Hi John'
 */
export function mockTranslation(
  key: string,
  options?: string | Record<string, unknown>
): string {
  // Handle string second argument (default value)
  if (typeof options === 'string') {
    return options
  }

  // Handle object with interpolation values
  if (typeof options === 'object' && options !== null) {
    // Start with defaultValue or the key itself
    let text = (options.defaultValue as string) ?? key

    // Replace {{placeholder}} patterns with actual values
    Object.entries(options).forEach(([optKey, optValue]) => {
      if (optKey !== 'defaultValue') {
        text = text.replace(
          new RegExp(`\\{\\{\\s*${optKey}\\s*\\}\\}`, 'g'),
          String(optValue)
        )
      }
    })

    return text
  }

  return key
}

/**
 * Create a mock i18n object with configurable language
 */
export function createMockI18n(language: string = 'en') {
  return {
    language,
    dir: () => (language === 'ar' ? 'rtl' : 'ltr'),
    changeLanguage: vi.fn().mockResolvedValue(undefined),
  }
}

/**
 * Complete mock factory for react-i18next
 *
 * @example
 * vi.mock('react-i18next', () => createI18nMock())
 */
export function createI18nMock(language: string = 'en') {
  return {
    useTranslation: () => ({
      t: mockTranslation,
      i18n: createMockI18n(language),
    }),
    Trans: ({ children }: { children: React.ReactNode }) => children,
  }
}
