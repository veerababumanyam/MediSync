import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { motion, AnimatePresence } from 'framer-motion'
import { Sparkles, ArrowUpRight } from 'lucide-react'
import { useTheme } from '../../hooks/useTheme'

export interface AnnouncementBannerProps {
  /** Announcement message to display */
  message?: string
  /** Optional label for a call-to-action */
  ctaLabel?: string
  /** Optional href for CTA; when provided renders an anchor element */
  ctaHref?: string
  /** Optional click handler for CTA; used when no href is provided */
  onCtaClick?: () => void
}

/**
 * AnnouncementBanner Component
 *
 * A premium, organic announcement banner inspired by iOS "Dynamic Island".
 * Features liquid glassmorphism, Framer Motion animations, and persistent dismissal.
 *
 * @module components/landing/AnnouncementBanner
 */
export function AnnouncementBanner({
  message,
  ctaLabel,
  ctaHref,
  onCtaClick,
}: AnnouncementBannerProps) {
  const { t } = useTranslation()
  const { isDark } = useTheme()
  const [isVisible, setIsVisible] = useState(false)

  const displayMessage = message ?? t('announcement.message', 'Turn Legacy System into AI Systems')

  useEffect(() => {
    // Show banner with small delay for organic entrance
    const timer = setTimeout(() => setIsVisible(true), 100)
    return () => clearTimeout(timer)
  }, [])


  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          initial={{ opacity: 0, y: -20, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: -20, scale: 0.95 }}
          transition={{
            duration: 0.4,
            ease: [0.16, 1, 0.3, 1], // iOS-style curve
          }}
          style={{ zIndex: 'var(--z-sticky)' }}
          className="relative mx-auto mt-6 px-4 sm:px-6 lg:px-8"
          role="banner"
          aria-label={t('announcement.ariaLabel', 'Announcement')}
        >
          <div
            className={`
              relative mx-auto flex max-w-fit items-center gap-4 rounded-full px-5 py-2.5
              transition-all duration-300 shadow-glass-lg
              ${isDark
                ? 'bg-gradient-to-r from-brand-primary to-brand-secondary text-white border-t border-white/20'
                : 'bg-surface/60 backdrop-blur-md border border-brand-primary/20 text-brand-primary shadow-xl shadow-brand-primary-alpha'}
            `}
          >
            {/* Shimmer effect overlay */}
            <div className="absolute inset-0 pointer-events-none overflow-hidden rounded-full">
              <motion.div
                animate={{
                  left: ['-100%', '200%'],
                }}
                transition={{
                  duration: 4,
                  repeat: Infinity,
                  ease: "linear",
                  repeatDelay: 6
                }}
                className={`absolute top-0 w-1/4 h-full skew-x-[-20deg] bg-gradient-to-r from-transparent 
                  ${isDark ? 'via-white/10' : 'via-brand-primary-alpha'} to-transparent`}
              />
            </div>

            {/* Icon */}
            <div className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full backdrop-blur-sm shadow-inner
              ${isDark ? 'bg-white/10' : 'bg-brand-primary-alpha'}`}>
              <Sparkles className={`h-4 w-4 ${isDark ? 'text-white' : 'text-brand-primary'}`} aria-hidden="true" />
            </div>

            {/* Message */}
            <p className={`max-w-[180px] xs:max-w-[280px] sm:max-w-md truncate text-sm font-semibold tracking-tight sm:text-[0.95rem] text-center
              ${isDark ? 'text-white' : 'text-brand-primary'}`}>
              {displayMessage}
            </p>

            {/* Actions Section */}
            {ctaLabel && (
              <div className="flex items-center gap-2 border-l border-white/20 pl-4">
                {ctaHref ? (
                  <a
                    href={ctaHref}
                    className="flex items-center gap-1.5 rounded-full bg-white px-3 py-1.5 text-xs font-bold text-brand-primary shadow-md transition-all hover:scale-105 hover:shadow-lg active:scale-95"
                  >
                    {ctaLabel}
                    <ArrowUpRight className="h-3.5 w-3.5" />
                  </a>
                ) : (
                  <button
                    type="button"
                    onClick={onCtaClick}
                    className="flex items-center gap-1.5 rounded-full bg-white px-3 py-1.5 text-xs font-bold text-brand-primary shadow-md transition-all hover:scale-105 hover:shadow-lg active:scale-95"
                  >
                    {ctaLabel}
                    <ArrowUpRight className="h-3.5 w-3.5" />
                  </button>
                )}
              </div>
            )}
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}

export default AnnouncementBanner
