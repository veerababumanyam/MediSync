

export interface FinalCTAProps {
    isDark: boolean
    onOpenLeadCapture?: () => void
}

export function FinalCTA({ isDark, onOpenLeadCapture }: FinalCTAProps) {
    return (
        <section id="cta" className="mb-24 mt-28 relative z-0 animate-fade-in-up" aria-labelledby="cta-heading">
            <div className="liquid-glass-cta p-8 sm:p-12 lg:p-16 text-center">
                {/* Decorative background orbs */}
                <div className="absolute top-0 left-0 w-64 h-64 bg-white/10 rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2 pointer-events-none" aria-hidden="true" />
                <div className="absolute bottom-0 right-0 w-64 h-64 bg-teal-400/20 rounded-full blur-3xl translate-x-1/3 translate-y-1/3 pointer-events-none" aria-hidden="true" />

                <div className="relative z-10 max-w-3xl mx-auto">
                    <h2 id="cta-heading" className="text-3xl sm:text-4xl lg:text-5xl font-extrabold mb-6 tracking-tight text-white">
                        Ready to Modernize Your Healthcare Data?
                    </h2>
                    <p className={`text-base sm:text-lg mb-10 leading-relaxed ${isDark ? 'text-slate-300' : 'text-blue-100'}`}>
                        Join leading hospitals and clinics turning their legacy systems into intelligent, conversational engines in days—not years.
                    </p>
                    <div className="flex flex-col sm:flex-row justify-center items-center gap-4">
                        <button
                            type="button"
                            onClick={onOpenLeadCapture}
                            aria-label="Get started with MediSync for free"
                            className={`liquid-glass-button-primary inline-flex items-center justify-center px-8 py-4 rounded-xl font-bold text-lg shadow-xl hover:-translate-y-1 transition-all w-full sm:w-auto min-h-[44px] focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-trust-blue)] dark:focus-visible:outline-cyan-400 ${isDark
                                ? 'bg-gradient-to-r from-blue-500 to-cyan-400 text-slate-900 hover:shadow-cyan-500/25'
                                : 'bg-white text-blue-600 hover:bg-slate-50 hover:shadow-white/25'
                                }`}
                        >
                            Get Started Free
                        </button>
                        <button
                            type="button"
                            onClick={onOpenLeadCapture}
                            aria-label="Book a demo with MediSync team"
                            className={`liquid-glass-button inline-flex items-center justify-center px-8 py-4 rounded-xl font-bold text-lg transition-all w-full sm:w-auto min-h-[44px] focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-trust-blue)] dark:focus-visible:outline-cyan-400 ${isDark
                                ? 'border border-white/20 text-white hover:bg-white/5'
                                : 'border border-white/30 text-white hover:bg-white/10'
                                }`}
                        >
                            Book a Demo
                        </button>
                    </div>
                </div>
            </div>
        </section>
    )
}
