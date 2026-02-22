import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { HeroCarousel } from '../components/landing/HeroCarousel'
import { SectorsSection } from '../components/landing/SectorsSection'
import { FeatureCard } from '../components/landing/FeatureCard'
import { FinalCTA } from '../components/landing/FinalCTA'
import { LeadCaptureModal } from '../components/common/LeadCaptureModal'

export interface HomeProps {
    isDark: boolean
}

export function Home({ isDark }: HomeProps) {
    const { t } = useTranslation()
    const [isModalOpen, setIsModalOpen] = useState(false)

    const faqSchema = {
        "@context": "https://schema.org",
        "@type": "FAQPage",
        "mainEntity": [
            {
                "@type": "Question",
                "name": t('faq.q1', 'What is MediSync?'),
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": t('faq.a1', 'MediSync is an advanced AI-powered business intelligence platform specifically designed for hospitals and healthcare providers. It acts as an intelligent bridge, seamlessly syncing operational data from your HIMS with financial data in Tally ERP.')
                }
            },
            {
                "@type": "Question",
                "name": t('faq.q2', 'How does the Tally ERP integration work?'),
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": t('faq.a2', 'MediSync utilizes secure, proprietary TDL XML over HTTP to establish a bi-directional sync with Tally ERP. This ensures idempotent, auditable, and instant synchronization of all financial vouchers and ledgers without manual export/import processes.')
                }
            },
            {
                "@type": "Question",
                "name": t('faq.q3', 'Is patient data secure and compliant?'),
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": t('faq.a3', 'Absolutely. MediSync employs a robust, multi-layered PII (Personally Identifiable Information) protection architecture. Patient identities are masked or anonymized before any data is processed for financial analytics, ensuring strict compliance with healthcare data regulations.')
                }
            }
        ]
    }

    return (
        <div className="flex flex-col flex-grow">
            {/* FAQ Schema for SEO */}
            <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(faqSchema) }} />
            <main className="max-w-6xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-12 sm:py-16 relative z-10 flex-grow">
                {/* Hero Section */}
                <section className="mb-20 animate-fade-in-up">
                    <div className="text-center mb-8">
                        <div className={`inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold tracking-wider uppercase ${isDark
                            ? 'bg-blue-500/10 border border-blue-500/20 text-blue-400'
                            : 'bg-blue-50 border border-blue-200 text-blue-600'
                            }`}>
                            <div className={`w-1.5 h-1.5 rounded-full animate-pulse ${isDark ? 'bg-blue-400' : 'bg-blue-600'}`} />
                            {t('welcome.badge', '#1 Rated Healthcare BI Platform')}
                        </div>
                    </div>

                    <HeroCarousel isDark={isDark} onOpenLeadCapture={() => setIsModalOpen(true)} />
                </section>

                <SectorsSection isDark={isDark} />

                {/* Feature Cards - Expanded to 6 features for comprehensive coverage */}
                <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-20">
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 01-.825-.242m9.345-8.334a2.126 2.126 0 00-.476-.095 48.64 48.64 0 00-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0011.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155" />
                            </svg>
                        }
                        gradient="from-blue-500 to-cyan-400"
                        shadowColor="shadow-blue-500/20"
                        title={t('features.conversationalBI.title', 'Conversational BI')}
                        description={t('features.conversationalBI.description', 'Chat with your healthcare data. Ask complex financial questions in plain language and get instant, beautiful visualizations.')}
                        delay="delay-1"
                    />
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
                            </svg>
                        }
                        gradient="from-emerald-500 to-teal-400"
                        shadowColor="shadow-emerald-500/20"
                        title={t('features.tallySync.title', 'Seamless Tally Integration')}
                        description={t('features.tallySync.description', 'Bi-directional, zero-latency sync with Tally ERP via TDL XML. Eliminate manual data entry and reconciliation errors forever.')}
                        delay="delay-2"
                    />
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                            </svg>
                        }
                        gradient="from-purple-500 to-pink-400"
                        shadowColor="shadow-purple-500/20"
                        title={t('features.aiAccountant.title', 'Automated AI Accountant')}
                        description={t('features.aiAccountant.description', 'Our proprietary OCR engine digitizes invoices and receipts with 99.9% accuracy, automatically mapping to your Chart of Accounts.')}
                        delay="delay-3"
                    />
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
                            </svg>
                        }
                        gradient="from-amber-400 to-orange-500"
                        shadowColor="shadow-amber-500/20"
                        title={t('features.piiProtection.title', 'Enterprise PII Protection')}
                        description={t('features.piiProtection.description', 'Bank-grade security ensures patient privacy. Built-in anonymization guarantees full compliance with global healthcare data regulations.')}
                        delay="delay-1"
                    />
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5m.75-9l3-3 2.148 2.148A12.061 12.061 0 0116.5 7.605" />
                            </svg>
                        }
                        gradient="from-indigo-500 to-blue-400"
                        shadowColor="shadow-indigo-500/20"
                        title={t('features.prescriptiveAnalytics.title', 'Prescriptive Analytics')}
                        description={t('features.prescriptiveAnalytics.description', "Don't just see what happened—know what to do next. Get actionable recommendations to plug revenue leaks and optimize billing.")}
                        delay="delay-2"
                    />
                    <FeatureCard
                        isDark={isDark}
                        icon={
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" />
                            </svg>
                        }
                        gradient="from-rose-500 to-red-400"
                        shadowColor="shadow-rose-500/20"
                        title={t('features.himsConnectivity.title', 'Unified HIMS Connectivity')}
                        description={t('features.himsConnectivity.description', 'Plug-and-play integrations with leading Hospital Information Management Systems. Break down data silos in minutes, not months.')}
                        delay="delay-3"
                    />
                </section>

                {/* AI GEO / SEO FAQ Section */}
                <section className="mb-16 animate-fade-in-up delay-4">
                    <div className="text-center mb-10">
                        <h2 className={`text-3xl font-bold mb-4 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                            {t('faq.title', 'Frequently Asked Questions')}
                        </h2>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className={`glass glass-shine rounded-2xl p-6 ${isDark ? 'hover:bg-white/5' : 'hover:bg-slate-50'}`}>
                            <h3 className={`text-xl font-semibold mb-3 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                                {t('faq.q1', 'What is MediSync?')}
                            </h3>
                            <p className={`leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                                {t('faq.a1', 'MediSync is an advanced AI-powered business intelligence platform specifically designed for hospitals and healthcare providers. It acts as an intelligent bridge, seamlessly syncing operational data from your HIMS with financial data in Tally ERP.')}
                            </p>
                        </div>
                        <div className={`glass glass-shine rounded-2xl p-6 ${isDark ? 'hover:bg-white/5' : 'hover:bg-slate-50'}`}>
                            <h3 className={`text-xl font-semibold mb-3 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                                {t('faq.q2', 'How does the Tally ERP integration work?')}
                            </h3>
                            <p className={`leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                                {t('faq.a2', 'MediSync utilizes secure, proprietary TDL XML over HTTP to establish a bi-directional sync with Tally ERP. This ensures idempotent, auditable, and instant synchronization of all financial vouchers and ledgers without manual export/import processes.')}
                            </p>
                        </div>
                        <div className={`glass glass-shine rounded-2xl p-6 md:col-span-2 ${isDark ? 'hover:bg-white/5' : 'hover:bg-slate-50'}`}>
                            <h3 className={`text-xl font-semibold mb-3 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                                {t('faq.q3', 'Is patient data secure and compliant?')}
                            </h3>
                            <p className={`leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                                {t('faq.a3', 'Absolutely. MediSync employs a robust, multi-layered PII (Personally Identifiable Information) protection architecture. Patient identities are masked or anonymized before any data is processed for financial analytics, ensuring strict compliance with healthcare data regulations.')}
                            </p>
                        </div>
                    </div>
                </section>

                <FinalCTA isDark={isDark} onOpenLeadCapture={() => setIsModalOpen(true)} />

            </main>

            {/* Footer - Enhanced with links and compliance policies */}
            <footer className={`border-t mt-auto relative z-10 transition-colors duration-300 ${isDark ? 'border-white/10' : 'border-slate-200'
                }`}>
                <div className={`max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10 ${isDark ? 'text-slate-500' : 'text-slate-500'}`}>
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8 text-sm">
                        <div className="md:col-span-1">
                            <img src="/logo.png" alt="MediSync" className="w-8 h-8 mb-4 opacity-70" />
                            <p className="mb-4">AI-Powered Conversational BI & Intelligent Accounting for Healthcare.</p>
                        </div>
                        <div>
                            <h4 className={`font-bold mb-4 ${isDark ? 'text-white' : 'text-slate-900'}`}>Platform</h4>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Features</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Integrations</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Pricing</a></li>
                            </ul>
                        </div>
                        <div>
                            <h4 className={`font-bold mb-4 ${isDark ? 'text-white' : 'text-slate-900'}`}>Compliance</h4>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-blue-500 transition-colors">HIPAA Compliance</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Security Policy</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">System Status</a></li>
                            </ul>
                        </div>
                        <div>
                            <h4 className={`font-bold mb-4 ${isDark ? 'text-white' : 'text-slate-900'}`}>Company</h4>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-blue-500 transition-colors">About Us</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Contact</a></li>
                                <li><a href="#" className="hover:text-blue-500 transition-colors">Privacy Policy</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className={`pt-8 border-t text-center text-xs ${isDark ? 'border-white/5' : 'border-slate-200'}`}>
                        <p>{t('footer.copyright', '© 2026 MediSync. AI-Powered Conversational BI & Intelligent Accounting for Healthcare.')}</p>
                    </div>
                </div>
            </footer>

            <LeadCaptureModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                isDark={isDark}
            />
        </div>
    )
}
