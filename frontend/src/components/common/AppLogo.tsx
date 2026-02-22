/**
 * AppLogo Component
 *
 * Renders an always-available inline MediSync mark.
 * Used in headers and loading states across the application.
 *
 * @module components/common/AppLogo
 */
import React from 'react'
import { useTranslation } from 'react-i18next'

export interface AppLogoProps {
  /** Size variant; 'sm' = 40px, 'md' = 64px */
  size?: 'sm' | 'md'
  /** Optional wrapper className (e.g. for hover/transition) */
  className?: string
  /** Optional img className */
  imgClassName?: string
}

/**
 * AppLogo - Always-available MediSync logo mark
 */
export const AppLogo: React.FC<AppLogoProps> = ({
  size = 'sm',
  className = '',
  imgClassName = '',
}) => {
  const { t } = useTranslation()
  const sizePx = size === 'sm' ? 40 : 64
  const alt = t('app.name', 'MediSync')
  const [sourceIndex, setSourceIndex] = React.useState(0)
  const base = import.meta.env.BASE_URL || '/'
  const primary = size === 'sm' ? 'icons/logo-64x64.png' : 'icons/logo-128x128.png'
  const candidates = React.useMemo(
    () => [`${base}${primary}`, `${base}logo.png`],
    [base, primary]
  )
  const activeSource = candidates[sourceIndex]
  const hasImageError = sourceIndex >= candidates.length

  return (
    <div
      className={`flex items-center justify-center overflow-hidden rounded-lg bg-surface-glass border border-glass shadow-lg ${className}`}
      style={{ width: sizePx, height: sizePx }}
    >
      {!hasImageError ? (
        <img
          src={activeSource}
          alt={alt}
          width={sizePx}
          height={sizePx}
          className={`object-contain ${imgClassName}`}
          decoding="async"
          onError={() => setSourceIndex((prev) => prev + 1)}
        />
      ) : (
        <span
          className={`inline-flex items-center justify-center w-full h-full text-primary text-xs font-bold ${imgClassName}`}
          aria-label={alt}
        >
          MS
        </span>
      )}
    </div>
  )
}

export default AppLogo
