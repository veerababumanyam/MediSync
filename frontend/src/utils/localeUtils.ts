import { getBcp47Locale } from '../i18n'

/**
 * Format a number using BCP 47 locale (e.g. en-US, ar-SA).
 * Use for consistent number formatting across the app.
 */
export function formatNumber(
  value: number,
  appLocale: string,
  options?: Intl.NumberFormatOptions
): string {
  const locale = getBcp47Locale(appLocale)
  return new Intl.NumberFormat(locale, options).format(value)
}

/**
 * Format a date using BCP 47 locale.
 * Use for consistent date/time formatting across the app.
 */
export function formatDate(
  date: Date,
  appLocale: string,
  options?: Intl.DateTimeFormatOptions
): string {
  const locale = getBcp47Locale(appLocale)
  return new Intl.DateTimeFormat(locale, options).format(date)
}

/**
 * Format a time string using BCP 47 locale.
 */
export function formatTime(
  date: Date,
  appLocale: string,
  options?: Intl.DateTimeFormatOptions
): string {
  const locale = getBcp47Locale(appLocale)
  return new Intl.DateTimeFormat(locale, { ...options, hour: '2-digit', minute: '2-digit' }).format(date)
}
