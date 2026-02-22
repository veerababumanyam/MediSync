

export interface FinalCTAProps {
    isDark: boolean
    onOpenLeadCapture?: () => void
}

export function FinalCTA({ isDark, onOpenLeadCapture }: FinalCTAProps) {
    return (
        <section className="mb-24 mt-32 animate-fade-in-up">
            <div className={`relative overflow-hidden rounded-[2rem] p-8 sm:p-12 lg:p-16 text-center shadow-2xl ${isDark
                ? 'bg-gradient-to-br from-blue-900/40 via-slate-800 to-teal-900/40 border border-white/10 shadow-black/50'
                : 'bg-gradient-to-br from-blue-600 via-blue-700 to-teal-600 text-white shadow-blue-500/20'
                }`}
            >
                {/* Decorative background orbs */}
                <div className="absolute top-0 left-0 w-64 h-64 bg-white/10 rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2 pointer-events-none" />
                <div className="absolute bottom-0 right-0 w-64 h-64 bg-teal-400/20 rounded-full blur-3xl translate-x-1/3 translate-y-1/3 pointer-events-none" />

                <div className="relative z-10 max-w-3xl mx-auto">
                    <h2 className={`text-3xl sm:text-4xl lg:text-5xl font-extrabold mb-6 tracking-tight ${isDark ? 'text-white' : 'text-white'}`}>
                        Ready to Modernize Your Healthcare Data?
                    </h2>
                    <p className={`text-lg sm:text-xl mb-10 leading-relaxed ${isDark ? 'text-slate-300' : 'text-blue-100'}`}>
                        Join leading hospitals and clinics turning their legacy systems into intelligent, conversational engines in days—not years.
                    </p>
                    <div className="flex flex-col sm:flex-row justify-center items-center gap-4">
                        <button
                            onClick={onOpenLeadCapture}
                            className={`px-8 py-4 rounded-xl font-bold text-lg shadow-xl hover:-translate-y-1 transition-all w-full sm:w-auto ${isDark
                                ? 'bg-gradient-to-r from-blue-500 to-cyan-400 text-slate-900 hover:shadow-cyan-500/25'
                                : 'bg-white text-blue-600 hover:bg-slate-50 hover:shadow-white/25'
                                }`}
                        >
                            Get Started Free
                        </button>
                        <button
                            onClick={onOpenLeadCapture}
                            className={`px-8 py-4 rounded-xl font-bold text-lg border-2 transition-all w-full sm:w-auto ${isDark
                                ? 'border-white/20 text-white hover:bg-white/5'
                                : 'border-white/30 text-white hover:bg-white/10'
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
