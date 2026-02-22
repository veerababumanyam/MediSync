import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/cn'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface LiquidGlassToastProps {
  message: string
  type: ToastType
  onClose: () => void
  duration?: number
  className?: string
}

const typeStyles: Record<ToastType, string> = {
  success: 'border-emerald-400/30 bg-emerald-500/10',
  error: 'border-red-400/30 bg-red-500/10',
  warning: 'border-amber-400/30 bg-amber-500/10',
  info: 'border-blue-400/30 bg-blue-500/10',
}

const typeIcons: Record<ToastType, React.ReactNode> = {
  success: (
    <svg className="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
    </svg>
  ),
  error: (
    <svg className="w-5 h-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
  ),
  warning: (
    <svg className="w-5 h-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  ),
  info: (
    <svg className="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  ),
}

export const LiquidGlassToast: React.FC<LiquidGlassToastProps> = ({
  message,
  type,
  onClose,
  className,
}) => {
  const { t } = useTranslation('common')
  return (
    <motion.div
      className={cn(
        'liquid-glass flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border',
        typeStyles[type],
        className
      )}
      initial={{ opacity: 0, y: 50, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: 20, scale: 0.95 }}
      transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1] }}
    >
      {typeIcons[type]}
      <p className="flex-1 text-sm font-medium liquid-text-primary">{message}</p>
      <button
        onClick={onClose}
        className="p-1 hover:bg-white/10 rounded-lg transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-1"
        aria-label={t('toast.close', 'Close')}
      >
        <svg className="w-4 h-4 liquid-text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </motion.div>
  )
}

export default LiquidGlassToast
