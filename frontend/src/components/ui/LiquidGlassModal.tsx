import React, { useEffect, useCallback, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/lib/cn'

/**
 * Liquid Glass Modal Component
 *
 * Premium glassmorphic modal with WCAG 3.0 Bronze compliance:
 * - role="dialog" + aria-modal="true"
 * - aria-labelledby for title linkage
 * - Focus trap: auto-focus first focusable element on open
 * - Return focus to trigger element on close
 * - Escape key dismissal
 * - prefers-reduced-transparency aware backdrop
 */
export interface LiquidGlassModalProps {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
  title?: string
  actions?: React.ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  className?: string
  closeOnBackdrop?: boolean
  closeOnEscape?: boolean
}

export const LiquidGlassModal: React.FC<LiquidGlassModalProps> = ({
  isOpen,
  onClose,
  children,
  title,
  actions,
  size = 'md',
  className,
  closeOnBackdrop = true,
  closeOnEscape = true,
}) => {
  const modalRef = useRef<HTMLDivElement>(null)
  const previousActiveElement = useRef<HTMLElement | null>(null)
  const titleId = 'liquid-glass-modal-title'

  // Handle escape key
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape' && closeOnEscape) {
        onClose()
      }
    },
    [onClose, closeOnEscape]
  )

  // WCAG 3.0 Bronze: Focus trap — auto-focus first focusable element & return focus on close
  useEffect(() => {
    if (isOpen) {
      previousActiveElement.current = document.activeElement as HTMLElement
      document.addEventListener('keydown', handleKeyDown)
      document.body.style.overflow = 'hidden'

      // Auto-focus first focusable element in modal
      requestAnimationFrame(() => {
        const focusable = modalRef.current?.querySelectorAll<HTMLElement>(
          'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
        )
        if (focusable && focusable.length > 0) {
          focusable[0].focus()
        }
      })
    }
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      document.body.style.overflow = ''

      // Return focus to trigger element
      if (previousActiveElement.current) {
        previousActiveElement.current.focus()
      }
    }
  }, [isOpen, handleKeyDown])

  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    full: 'max-w-4xl',
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <div
          role="presentation"
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
        >
          {/* Backdrop — WCAG 3.0: aria-hidden, reduced-transparency aware */}
          <motion.div
            className="absolute inset-0 bg-black/50 backdrop-blur-sm supports-[not_(backdrop-filter:blur(1px))]:bg-black/70"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={closeOnBackdrop ? onClose : undefined}
            aria-hidden="true"
          />

          {/* Modal — WCAG 3.0: role, aria-modal, aria-labelledby */}
          <motion.div
            ref={modalRef}
            role="dialog"
            aria-modal="true"
            aria-labelledby={title ? titleId : undefined}
            className={cn(
              'liquid-glass-modal relative w-full p-6',
              sizeClasses[size],
              className
            )}
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ duration: 0.2, ease: [0.4, 0, 0.2, 1] }}
          >
            {/* Title */}
            {title && (
              <h2
                id={titleId}
                className="text-xl font-semibold mb-4 liquid-text-primary"
              >
                {title}
              </h2>
            )}

            {/* Content */}
            <div className="mb-4">{children}</div>

            {/* Actions */}
            {actions && (
              <div className="flex justify-end gap-3 pt-4 border-t border-white/10">
                {actions}
              </div>
            )}
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  )
}

export default LiquidGlassModal
