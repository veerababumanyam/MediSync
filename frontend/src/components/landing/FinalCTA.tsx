import { useTranslation } from 'react-i18next'
import { CTA_CLASSES } from '@/lib/cta-classes'

export interface FinalCTAProps {
    isDark: boolean
    onOpenLeadCapture?: () => void
}

export function FinalCTA({ onOpenLeadCapture }: Omit<FinalCTAProps, 'isDark'>) {
    const { t } = useTranslation()
    return (
        <section id="cta" className="mb-24 mt-28 relative z-0 animate-fade-in-up" aria-labelledby="cta-heading">
            <div className="liquid-glass-cta p-8 sm:p-12 lg:p-16 text-center">
                {/* Decorative background orbs */}
                <div className="absolute top-0 left-0 w-64 h-64 bg-brand-primary-alpha rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2 pointer-events-none" aria-hidden="true" />
                <div className="absolute bottom-0 right-0 w-64 h-64 bg-brand-secondary-alpha rounded-full blur-3xl translate-x-1/3 translate-y-1/3 pointer-events-none" aria-hidden="true" />

                <div className="relative z-10 max-w-3xl mx-auto">
                    <h2 id="cta-heading" className="text-3xl sm:text-4xl lg:text-5xl font-extrabold mb-6 tracking-tight text-primary">
                        Ready to Modernize Your Healthcare Data?
                    </h2>
                    <p className={`text-base sm:text-lg mb-10 leading-relaxed text-secondary`}>
                        {t('cta.description', 'Join leading hospitals and clinics turning their legacy systems into intelligent, conversational engines in days—not years.')}
                    </p>
                    <div className="flex flex-col sm:flex-row justify-center items-center gap-4">
                        <button
                            type="button"
                            onClick={onOpenLeadCapture}
                            aria-label="Get started with AnySync for free"
                            className={CTA_CLASSES.LANDING_PRIMARY}
                        >
                            Get Started Free
                        </button>
                        <button
                            type="button"
                            onClick={onOpenLeadCapture}
                            aria-label="Book a demo with AnySync team"
                            className={CTA_CLASSES.LANDING_SECONDARY}
                        >
                            Book a Demo
                        </button>
                    </div>
                </div>
            </div>
        </section>
    )
}
