/**
 * AnnouncementBanner Component
 *
 * A simple announcement banner that appears above the hero section.
 * Features dismissible state with localStorage persistence, RTL support.
 *
 * @module components/landing/AnnouncementBanner
 */

import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'

const STORAGE_KEY = 'medisync-announcement-dismissed'

export interface AnnouncementBannerProps {
  /** Whether dark mode is enabled */
  isDark: boolean
  /** Announcement message to display */
  message?: string
}

/**
 * AnnouncementBanner Component
 *
 * Displays a dismissible announcement banner.
 * Persists dismissed state to localStorage.
 *
 * @example
 * ```tsx
 * <AnnouncementBanner
 *   isDark={isDark}
 *   message="🎉 New feature available!"
 * />
 * ```
 */
export function AnnouncementBanner({
  isDark,
  message,
}: AnnouncementBannerProps) {
  const { t } = useTranslation()
  const [isDismissed, setIsDismissed] = useState(true)
  const [isAnimating, setIsAnimating] = useState(false)

  const displayMessage = message ?? t('announcement.message', '🎉 New: AI Accountant module now available!')

  useEffect(() => {
    const dismissed = localStorage.getItem(STORAGE_KEY)
    if (dismissed !== 'true') {
      setIsDismissed(false)
      requestAnimationFrame(() => {
        setIsAnimating(true)
      })
    }
  }, [])

  if (isDismissed) {
    return null
  }

  return (
    <div
      className={`
        mx-4 mt-2 rounded-xl transition-all duration-300 ease-out
        ${isAnimating ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-2'}
        ${isDark
          ? 'bg-linear-to-r from-violet-600/20 via-purple-500/20 to-fuchsia-500/20 border-white/10'
          : 'bg-linear-to-r from-violet-100 via-purple-50 to-fuchsia-100 border-purple-200/60'
        }
        border backdrop-blur-sm
      `}
      role="banner"
      aria-label="Announcement"
    >
      <div className="flex items-center justify-center gap-3 px-4 py-2.5">
        {/* Message */}
        <p
          className={`
            text-sm font-medium text-center
            ${isDark ? 'text-purple-200' : 'text-purple-800'}
          `}
        >
          {displayMessage}
        </p>
      </div>
    </div>
  )
}

export default AnnouncementBanner
