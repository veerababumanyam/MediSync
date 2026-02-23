import { useState, useEffect, useMemo, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { HeroIllustration } from './icons'
import { TrustStripIcon, getTrustBrandColor } from './TrustStripIcons'
import { CTA_CLASSES } from '@/lib/cta-classes'

type TrustItemKey = 'item1' | 'item2' | 'item3' | 'item4' | 'item5' | 'item6' | 'item7' | 'item8'

export interface HeroCarouselProps {
    isDark: boolean
    onOpenLeadCapture?: () => void
}

export function HeroCarousel({ isDark, onOpenLeadCapture }: HeroCarouselProps) {
    const { t } = useTranslation()
    const [currentIndex, setCurrentIndex] = useState(0)
    const [slideKey, setSlideKey] = useState(0)
    const [isHovered, setIsHovered] = useState(false)

    const carouselItems = useMemo(() => ['slide1', 'slide2', 'slide3'], [])
    const trustItems = useMemo<TrustItemKey[]>(
        () => ['item1', 'item2', 'item3', 'item4', 'item5', 'item6', 'item7', 'item8'],
        []
    )

    // Auto-rotate effect, paused on hover
    useEffect(() => {
        if (isHovered) return

        const timer = setInterval(() => {
            setCurrentIndex((prev) => (prev + 1) % carouselItems.length)
            setSlideKey((prev) => prev + 1)
        }, 6000)
        return () => clearInterval(timer)
    }, [carouselItems.length, isHovered])

    const goToSlide = useCallback((index: number) => {
        setCurrentIndex(index)
        setSlideKey((prev) => prev + 1)
    }, [])

    return (
        <div
            className="relative max-w-7xl mx-auto mb-24"
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            onFocus={() => setIsHovered(true)}
            onBlur={() => setIsHovered(false)}
            role="region"
            aria-roledescription="carousel"
            aria-label="Hero Features Carousel"
        >
            {/* Visually hidden H1 for SEO, prevents multiple H1s */}
            <h1 className="sr-only">AnySync - Turn Legacy Healthcare System into AI Systems</h1>

            {/* Main Carousel Area - Using Pronounced Liquid Glass Design System */}
            <div
                className="liquid-glass-pronounced hero-carousel-container relative overflow-hidden p-1 animate-liquid-fade-in"
                aria-live={isHovered ? 'polite' : 'off'}
            >
                {/* Animated gradient border accent - using design tokens */}
                <div className="absolute inset-0 -z-10 rounded-[28px] opacity-50 will-change-transform" style={{
                    background: 'var(--gradient-brand-animated)',
                    backgroundSize: '300% 300%',
                    animation: 'gradientBorderShift 8s ease infinite',
                }} />

                <div className="relative grid">
                    {carouselItems.map((item, index) => {
                        const isActive = index === currentIndex
                        return (
                            <div
                                key={`${item}-${slideKey}`}
                                className={`col-start-1 row-start-1 w-full transition-all duration-700 ease-out will-change-transform will-change-opacity
                  ${isActive
                                        ? 'opacity-100 z-10 pointer-events-auto'
                                        : 'opacity-0 z-0 pointer-events-none'
                                    }`}
                                style={isActive ? { animation: 'heroSlideIn 0.7s ease-out forwards' } : undefined}
                                role="tabpanel"
                                id={`slide-panel-${index}`}
                                aria-labelledby={`slide-tab-${index}`}
                                aria-hidden={!isActive}
                            >
                                <div className="hero-slide-layout px-6 sm:px-10 md:px-14 py-10 sm:py-14 lg:py-16">
                                    {/* Inner content card with subtle separation */}
                                    <div className="relative rounded-2xl p-4 sm:p-6" style={{
                                        boxShadow: isDark
                                            ? 'inset 0 1px 0 0 rgba(255,255,255,0.05)'
                                            : 'inset 0 1px 0 0 rgba(255,255,255,0.5)',
                                    }}>
                                    {/* LEFT: Text content */}
                                    <div className="flex-1 min-w-0 text-center lg:text-left">
                                        {/* Logo + badge row */}
                                        <div className="flex items-center justify-center lg:justify-start gap-3 mb-5 animate-fade-in-down" style={{ animationDelay: '100ms' }}>
                                            <img src="/logo.png" fetchPriority="high" alt="AnySync logo - AI-powered healthcare BI platform" className="w-10 h-10 object-contain" />
                                            <span className={`text-sm font-bold tracking-wide ${isDark ? 'text-white/70' : 'text-slate-500'}`}>AnySync</span>
                                        </div>

                                        {/* Gradient headline as H2 for SEO */}
                                        <h2 className="hero-gradient-text text-3xl sm:text-4xl lg:text-[2.75rem] xl:text-5xl font-extrabold leading-[1.15] tracking-tight mb-5 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                                            {t(`heroCarousel.${item}.title`)}
                                        </h2>

                                        {/* Subtitle / Tagline */}
                                        <p className={`text-sm sm:text-base font-medium tracking-wide mb-4 ${isDark ? 'text-slate-400' : 'text-slate-500'} animate-fade-in-up`}
                                           style={{ animationDelay: '250ms' }}>
                                            {t(`heroCarousel.${item}.subtitle`)}
                                        </p>

                                        {/* Description */}
                                        <p className={`text-base sm:text-lg leading-relaxed mb-4 max-w-xl mx-auto lg:mx-0 ${isDark ? 'text-slate-400' : 'text-slate-600'} animate-fade-in-up`} style={{ animationDelay: '300ms' }}>
                                            {t(`heroCarousel.${item}.description`)}
                                        </p>

                                        {/* Capability strip - slide 2 only */}
                                        {item === 'slide2' && (
                                            <div className="flex flex-wrap justify-center lg:justify-start gap-3 mb-6 animate-fade-in-up" style={{ animationDelay: '350ms' }}>
                                                {(['bi', 'accountant', 'reports', 'analytics'] as const).map((cap) => (
                                                    <div
                                                        key={cap}
                                                        className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold transition-transform hover:scale-105 ${isDark ? 'bg-white/5 text-slate-300' : 'bg-slate-100 text-slate-600'}`}
                                                    >
                                                        <div className={`w-2 h-2 rounded-full ${isDark ? 'bg-teal-400' : 'bg-blue-500'}`} />
                                                        {t(`heroCarousel.capabilities.${cap}`)}
                                                    </div>
                                                ))}
                                            </div>
                                        )}

                                        {/* CTA Button - with pronounced glass and hover glow */}
                                        <div className="mb-8 animate-fade-in-up" style={{ animationDelay: '400ms' }}>
                                            <button
                                                className={CTA_CLASSES.HERO}
                                                type="button"
                                                onClick={onOpenLeadCapture}
                                            >
                                                {t(`heroCarousel.${item}.cta`)}
                                            </button>
                                        </div>

                                        {/* Stat counters with Top Rated Badge */}
                                        <div className="flex flex-wrap justify-center lg:justify-start gap-4 sm:gap-6 items-center animate-fade-in-up" style={{ animationDelay: '500ms' }}>
                                            <div className="liquid-glass-badge flex items-center gap-1 px-3 py-1.5" aria-label="Rated 4.9 out of 5 stars based on 154 reviews">
                                                <div className="flex text-amber-400">
                                                    {'★'.repeat(5)}
                                                </div>
                                                <span className={`text-xs ml-1 font-bold ${isDark ? 'text-slate-200' : 'text-slate-700'}`}>4.9/5</span>
                                            </div>

                                            {['stat1', 'stat2', 'stat3'].map((stat, i) => (
                                                <div key={stat} className={`flex items-center gap-2 text-sm font-semibold ${isDark ? 'text-slate-200' : 'text-slate-700'}`} style={{ animationDelay: `${600 + i * 100}ms` }}>
                                                    <div className={`w-2 h-2 rounded-full ${isDark ? 'bg-teal-400' : 'bg-blue-500'}`} />
                                                    {t(`heroCarousel.${item}.${stat}`)}
                                                </div>
                                            ))}
                                        </div>
                                    </div>

                                    {/* RIGHT: SVG Illustration - with floating animation and backdrop glow */}
                                    <div className="shrink-0 w-full max-w-[260px] sm:max-w-[300px] lg:max-w-[340px] animate-float-medium relative">
                                        {/* Backdrop glow */}
                                        <div
                                            className="absolute inset-0 -z-10 rounded-full opacity-30"
                                            style={{
                                                background: `radial-gradient(circle, ${isDark ? 'rgba(24, 146, 157, 0.2)' : 'rgba(39, 80, 168, 0.15)'} 0%, transparent 70%)`,
                                                transform: 'scale(1.3)',
                                            }}
                                        />
                                        <HeroIllustration slide={item} isDark={isDark} />
                                    </div>
                                    </div>{/* End inner content card */}
                                </div>
                            </div>
                        )
                    })}
                </div>

                {/* Progress bar - with gradient shine using tokens */}
                <div className={`h-1 ${isDark ? 'bg-white/5' : 'bg-slate-100'} rounded-full overflow-hidden`}>
                    <div
                        key={`progress-${slideKey}`}
                        className="h-full rounded-full animate-shine"
                        style={{
                            background: 'var(--gradient-brand-shift)',
                            backgroundSize: '200% 100%',
                            animation: isHovered ? 'none' : 'heroProgress 6s linear forwards',
                            width: isHovered ? '100%' : '0%',
                        }}
                    />
                </div>
            </div>

            {/* Dot indicators - visually hidden, kept for keyboard and screen reader navigation (Option B) */}
            <div className="sr-only flex justify-center gap-3 mt-4" role="tablist" aria-label="Carousel slide controls">
                {carouselItems.map((_, index) => (
                    <button
                        key={index}
                        type="button"
                        role="tab"
                        id={`slide-tab-${index}`}
                        aria-selected={index === currentIndex}
                        aria-controls={`slide-panel-${index}`}
                        aria-label={`Go to slide ${index + 1}`}
                        className={`hero-dot ${index === currentIndex ? 'active' : ''} transition-all duration-300 liquid-focus-gold`}
                        style={{
                            width: index === currentIndex ? '24px' : '10px',
                            height: '10px',
                            borderRadius: 'var(--radius-full)',
                            backgroundColor: index === currentIndex
                                ? 'var(--brand-secondary)'
                                : isDark ? 'rgba(255,255,255,0.2)' : 'rgba(0,0,0,0.1)',
                            boxShadow: index === currentIndex
                                ? 'var(--shadow-glow-secondary)'
                                : 'none',
                        }}
                        onClick={() => goToSlide(index)}
                    />
                ))}
            </div>

            {/* Trust Marquee Strip - Using Pronounced Liquid Glass */}
            <div
                className="liquid-glass-pronounced mt-10 rounded-2xl py-4 overflow-hidden"
                style={{ WebkitMaskImage: 'linear-gradient(to right, transparent, black 10%, black 90%, transparent)', maskImage: 'linear-gradient(to right, transparent, black 10%, black 90%, transparent)' }}
            >
                <p className={`text-center text-[10px] uppercase tracking-[0.2em] font-semibold mb-3 ${isDark ? 'text-slate-500' : 'text-slate-400'}`}>
                    Trusted Integrations & Partners
                </p>
                <div className="relative overflow-hidden w-full">
                    <div className="hero-marquee-track">
                        {/* Duplicate items for seamless loop */}
                        {[...trustItems, ...trustItems].map((item, i) => {
                            const brandColor = getTrustBrandColor(item, isDark)
                            return (
                                <div key={`${item}-${i}`} className="shrink-0 px-3 sm:px-4">
                                    <div
                                        className="liquid-glass-button-prominent flex items-center gap-2 rounded-full px-4 py-2"
                                    >
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
        </div>
    )
}
