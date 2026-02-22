import { useState } from 'react'

export interface LeadCaptureModalProps {
    isOpen: boolean
    onClose: () => void
    isDark: boolean
}

export function LeadCaptureModal({ isOpen, onClose, isDark }: LeadCaptureModalProps) {
    const [email, setEmail] = useState('')
    const [submitted, setSubmitted] = useState(false)

    if (!isOpen) return null

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        // Simulated submission
        setSubmitted(true)
        setTimeout(() => {
            onClose()
            setSubmitted(false)
            setEmail('')
        }, 2000)
    }

    return (
        <div className="fixed inset-0 z-[200] flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm transition-opacity"
                onClick={onClose}
                aria-hidden="true"
            />

            {/* Modal */}
            <div className={`relative w-full max-w-md p-8 rounded-3xl shadow-2xl transition-all elevate ${isDark ? 'bg-slate-900 border border-white/10' : 'bg-white border border-slate-200'}`}>
                <button
                    onClick={onClose}
                    className={`absolute top-4 right-4 p-2 rounded-full transition-colors ${isDark ? 'text-slate-400 hover:bg-white/10' : 'text-slate-500 hover:bg-slate-100'}`}
                    aria-label="Close modal"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>

                {submitted ? (
                    <div className="text-center py-8">
                        <div className="w-16 h-16 bg-emerald-100 text-emerald-500 rounded-full flex items-center justify-center mx-auto mb-4">
                            <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                        </div>
                        <h3 className={`text-2xl font-bold mb-2 ${isDark ? 'text-white' : 'text-slate-900'}`}>Request Received</h3>
                        <p className={isDark ? 'text-slate-400' : 'text-slate-600'}>We'll be in touch shortly to schedule your demo.</p>
                    </div>
                ) : (
                    <>
                        <h3 className={`text-2xl font-bold mb-2 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                            Get Started Free
                        </h3>
                        <p className={`mb-6 ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                            See how MediSync can transform your legacy systems. Enter your email to begin.
                        </p>

                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div>
                                <label htmlFor="email" className="sr-only">Email address</label>
                                <input
                                    id="email"
                                    type="email"
                                    required
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    placeholder="Enter your work email"
                                    className={`w-full px-4 py-3 rounded-xl border focus:ring-2 focus:ring-blue-500 outline-none transition-all ${isDark
                                        ? 'bg-slate-800 border-slate-700 text-white placeholder-slate-500'
                                        : 'bg-slate-50 border-slate-200 text-slate-900 placeholder-slate-400'
                                        }`}
                                />
                            </div>
                            <button
                                type="submit"
                                className="w-full py-3 px-6 rounded-xl font-bold text-white bg-gradient-to-r from-blue-600 to-cyan-500 hover:from-blue-500 hover:to-cyan-400 shadow-lg shadow-blue-500/25 transition-all outline-none"
                            >
                                Request Demo
                            </button>
                        </form>
                    </>
                )}
            </div>
        </div>
    )
}
