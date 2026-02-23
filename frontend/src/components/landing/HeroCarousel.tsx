import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { HeroIllustration } from './icons'
import { TrustStripIcon, getTrustBrandColor } from './TrustStripIcons'
import { CTA_CLASSES } from '../../lib/cta-classes'
import { cn } from '../../lib/cn'

type TrustItemKey = 'item1' | 'item2' | 'item3' | 'item4' | 'item5' | 'item6' | 'item7' | 'item8'

export interface HeroCarouselProps {
    isDark: boolean
    onOpenLeadCapture?: () => void
}

export function HeroCarousel({ isDark, onOpenLeadCapture }: HeroCarouselProps) {
    const { t } = useTranslation()
    const trustItems = useMemo<TrustItemKey[]>(
        () => ['item1', 'item2', 'item3', 'item4', 'item5', 'item6', 'item7', 'item8'],
        []
    )

    return (
        <section
            className="relative mb-24"
            aria-label={t('heroCarousel.ariaLabel', 'AnySync hero section')}
        >
            {/* Mesh gradient background */}
            <div
                className={`absolute inset-0 -z-20 overflow-hidden ${isDark ? 'bg-neutral-950' : 'bg-surface'
                    }`}
            >
                <div
                    className="pointer-events-none absolute -top-40 -left-32 h-96 w-96 rounded-full blur-3xl opacity-60"
                    style={{
                        background:
                            'radial-gradient(circle at 20% 0%, rgba(var(--brand-primary-rgb), 0.55), transparent 60%)',
                    }}
                />
                <div
                    className="pointer-events-none absolute -bottom-40 -right-32 h-[28rem] w-[28rem] rounded-full blur-3xl opacity-60"
                    style={{
                        background:
                            'radial-gradient(circle at 80% 100%, rgba(var(--brand-secondary-rgb), 0.6), transparent 60%)',
                    }}
                />
                <div
                    className="pointer-events-none absolute inset-x-0 top-1/3 h-64 blur-3xl opacity-40"
                    style={{
                        background:
                            'radial-gradient(circle at 50% 50%, rgba(var(--brand-primary-rgb), 0.45), transparent 65%)',
                    }}
                />
            </div>

            <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <h1 className="sr-only">
                    AnySync – Turn legacy healthcare systems into AI systems
                </h1>

                <div
                    className={cn(
                        'relative overflow-hidden rounded-3xl border transition-all duration-300 shadow-glass-lg',
                        isDark
                            ? 'border-glass-soft bg-neutral-900/80 shadow-[0_24px_80px_rgba(0,0,0,0.4)]'
                            : 'border-glass-soft bg-surface/90 shadow-[0_24px_80px_rgba(30,64,175,0.1)]'
                    )}
                >
                    {/* Subtle inner mesh / vignette */}
                    <div
                        className="pointer-events-none absolute inset-0 -z-10"
                        style={{
                            background:
                                'radial-gradient(circle at 0% 0%, rgba(var(--brand-primary-rgb), 0.16), transparent 55%), radial-gradient(circle at 100% 0%, rgba(var(--brand-secondary-rgb), 0.18), transparent 55%), radial-gradient(circle at 0% 100%, rgba(var(--brand-accent-rgb), 0.18), transparent 55%)',
                        }}
                    />

                    <div
                        className="relative grid grid-cols-1 lg:grid-cols-[3fr_2fr] gap-10 lg:gap-14 px-6 sm:px-10 md:px-14 py-10 sm:py-14 lg:py-16 items-center"
                    >
                        {/* Left: Text content */}
                        <div className="flex flex-col gap-6 lg:gap-7 text-center lg:text-left">
                            {/* Headline */}
                            <div className="space-y-3 pt-5">
                                <h2
                                    className={`text-4xl sm:text-5xl lg:text-6xl font-extrabold leading-tight tracking-tight`}
                                >
                                    <span
                                        className={`block bg-gradient-to-r from-white via-sky-100 to-slate-300 bg-clip-text text-transparent ${isDark ? '' : 'from-slate-900 via-sky-700 to-slate-900'
                                            }`}
                                    >
                                        {t('heroCarousel.slide1.title')}
                                    </span>
                                </h2>
                            </div>

                            {/* Body copy */}
                            <div className="space-y-3 max-w-xl mx-auto lg:mx-0">
                                <p
                                    className={`text-base sm:text-lg leading-relaxed ${isDark ? 'text-neutral-300' : 'text-secondary'
                                        }`}
                                >
                                    {t('heroCarousel.slide1.description')}
                                </p>
                                <p
                                    className={`text-sm sm:text-base ${isDark ? 'text-neutral-400' : 'text-neutral-500'
                                        }`}
                                >
                                    {t(
                                        'heroCarousel.slide1.subtitle',
                                        'Unify HIMS and Tally data into one AI-first analytics workspace.'
                                    )}
                                </p>
                            </div>

                            {/* CTAs */}
                            <div className="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-3 sm:gap-4 pt-2">
                                <button
                                    type="button"
                                    onClick={onOpenLeadCapture}
                                    className={`${CTA_CLASSES.HERO} px-6 sm:px-7 py-3 text-sm sm:text-base`}
                                >
                                    {t('heroCarousel.slide1.cta')}
                                </button>
                                <button
                                    type="button"
                                    className={`inline-flex items-center gap-2 rounded-full border px-4 sm:px-5 py-2.5 text-sm font-semibold transition-all duration-300 ${isDark
                                        ? 'border-neutral-700 bg-neutral-900/40 text-neutral-100 hover:border-brand-primary/50 hover:bg-neutral-800/70'
                                        : 'border-glass-soft bg-surface text-secondary hover:border-brand-primary/30 hover:bg-neutral-50'
                                        }`}
                                >
                                    {t('heroCarousel.secondaryCta', 'Watch 3‑minute demo')}
                                    <span aria-hidden="true">↗</span>
                                </button>
                            </div>

                            {/* Social proof */}
                            <div className="flex flex-wrap items-center justify-center lg:justify-start gap-x-6 gap-y-3 pt-2">
                                <div
                                    className={`flex items-center gap-2 rounded-full px-3 py-1.5 text-xs font-semibold ${isDark ? 'bg-white/5 text-amber-300' : 'bg-amber-50 text-amber-700'
                                        }`}
                                    aria-label="Rated 4.9 out of 5 stars"
                                >
                                    <span className="text-amber-400 text-base leading-none">
                                        {'★★★★★'}
                                    </span>
                                    <span>4.9/5 from finance teams</span>
                                </div>
                                <div
                                    className={`flex items-center gap-2 text-xs sm:text-sm font-medium ${isDark ? 'text-neutral-300' : 'text-secondary'
                                        }`}
                                >
                                    <span className="h-1.5 w-1.5 rounded-full bg-color-success" />
                                    {t('heroCarousel.slide1.stat1')}
                                </div>
                                <div
                                    className={`flex items-center gap-2 text-xs sm:text-sm font-medium ${isDark ? 'text-neutral-300' : 'text-secondary'
                                        }`}
                                >
                                    <span className="h-1.5 w-1.5 rounded-full bg-brand-primary" />
                                    {t('heroCarousel.slide1.stat2')}
                                </div>
                            </div>
                        </div>

                        {/* Right: Product visual */}
                        <div className="relative">
                            <div className="pointer-events-none absolute -inset-10 -z-10 opacity-60">
                                <div
                                    className="absolute inset-0 rounded-[2.75rem] blur-3xl"
                                    style={{
                                        background:
                                            'conic-gradient(from 140deg at 20% 20%, rgba(var(--brand-secondary-rgb), 0.18), transparent 40%, rgba(var(--brand-primary-rgb), 0.24), transparent 70%, rgba(var(--brand-secondary-rgb), 0.2))',
                                    }}
                                />
                            </div>
                            <div className="flex justify-center lg:justify-end">
                                <div className="w-full max-w-[320px] sm:max-w-[360px] lg:max-w-[400px]">
                                    <HeroIllustration slide="slide1" isDark={isDark} />
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Trust Marquee Strip */}
            <div
                className={`mt-10 rounded-2xl py-4 overflow-hidden border transition-all duration-300 ${isDark
                    ? 'border-glass-soft bg-neutral-900/80'
                    : 'border-glass-soft bg-neutral-50/80 shadow-sm'
                    }`}
                style={{
                    WebkitMaskImage: 'linear-gradient(to right, transparent, black 10%, black 90%, transparent)',
                    maskImage: 'linear-gradient(to right, transparent, black 10%, black 90%, transparent)'
                }}
            >
                <p className={`text-center text-[10px] uppercase tracking-[0.2em] font-semibold mb-3 ${isDark ? 'text-neutral-500' : 'text-neutral-400'}`}>
                    Trusted Integrations & Partners
                </p>
                <div className="relative overflow-hidden w-full">
                    <div className="hero-marquee-track">
                        {[...trustItems, ...trustItems].map((item, i) => {
                            const brandColor = getTrustBrandColor(item, isDark)
                            return (
                                <div key={`${item}-${i}`} className="shrink-0 px-3 sm:px-4">
                                    <div className="liquid-glass-button-prominent flex items-center gap-2 rounded-full px-4 py-2">
                                        <TrustStripIcon
                                            itemKey={item}
                                            isDark={isDark}
                                            className="w-5 h-5 shrink-0"
                                            style={{ color: brandColor }}
                                        />
                                        <span className="text-sm font-semibold whitespace-nowrap" style={{ color: brandColor }}>
                                            {t(`heroCarousel.trustStrip.${item}`)}
                                        </span>
                                    </div>
                                </div>
                            )
                        })}
                    </div>
                </div>
            </div>
        </section>
    )
}
