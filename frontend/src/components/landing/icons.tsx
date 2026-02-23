/**
 * Hero SVG Illustrations — unique per slide
 *
 * Uses design tokens from design-tokens.css for consistent theming
 */
export function HeroIllustration({ slide, isDark }: { slide: string; isDark: boolean }) {
    // Use design tokens for colors
    const accent = isDark ? 'var(--brand-secondary)' : 'var(--brand-primary)'
    const accentLight = isDark ? 'var(--brand-secondary)' : 'var(--brand-primary)'
    const nodeColor = isDark ? '#1e293b' : '#f1f5f9'
    const strokeColor = isDark ? 'rgba(255,255,255,0.15)' : 'rgba(0,0,0,0.1)'

    // Semantic colors from design tokens
    const errorColor = 'var(--color-error)'
    const warningColor = 'var(--color-warning)'
    const successColor = 'var(--color-success)'

    if (slide === 'slide1') {
        // Conversation with Legacy — AI brain connecting to legacy systems
        return (
            <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
                {/* Outer glow ring */}
                <circle cx="160" cy="140" r="120" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 6s ease-in-out infinite' }} />

                {/* Central AI node */}
                <g style={{ animation: 'floatY 4s ease-in-out infinite' }}>
                    <circle cx="160" cy="100" r="35" fill={accent} fillOpacity="0.15" stroke={accent} strokeWidth="2" />
                    <text x="160" y="105" textAnchor="middle" fill={accent} fontSize="14" fontWeight="700">AI</text>
                </g>

                {/* Connection lines with dash animation */}
                <line x1="160" y1="135" x2="80" y2="190" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />
                <line x1="160" y1="135" x2="160" y2="200" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />
                <line x1="160" y1="135" x2="240" y2="190" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.4" />

                {/* Legacy system nodes */}
                <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite' }}>
                    <rect x="45" y="175" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                    <text x="80" y="200" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">HIMS</text>
                </g>

                <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite 0.3s' }}>
                    <rect x="125" y="195" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                    <text x="160" y="220" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">Tally</text>
                </g>

                <g style={{ animation: 'floatYReverse 4.5s ease-in-out infinite 0.6s' }}>
                    <rect x="205" y="175" rx="8" width="70" height="40" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />
                    <text x="240" y="200" textAnchor="middle" fill={accent} fontSize="10" fontWeight="600">SQL</text>
                </g>

                {/* Chat bubbles — showing conversation */}
                <g style={{ animation: 'floatY 3.5s ease-in-out infinite' }}>
                    <rect x="50" y="50" rx="12" width="80" height="32" fill={nodeColor} stroke={strokeColor} strokeWidth="1" />
                    <rect x="62" y="60" width="45" height="4" rx="2" fill={strokeColor} />
                    <rect x="62" y="68" width="30" height="4" rx="2" fill={strokeColor} />
                </g>

                <g style={{ animation: 'floatYReverse 3s ease-in-out infinite 0.5s' }}>
                    <rect x="190" y="45" rx="12" width="80" height="32" fill={accent} fillOpacity="0.15" stroke={accent} strokeWidth="1" />
                    <rect x="202" y="57" width="55" height="4" rx="2" fill={accent} fillOpacity="0.6" />
                    <rect x="202" y="65" width="35" height="4" rx="2" fill={accent} fillOpacity="0.4" />
                </g>

                {/* Sparkle accents */}
                <circle cx="55" cy="140" r="3" fill={accent} fillOpacity="0.5" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
                <circle cx="265" cy="130" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.7s' }} />
                <circle cx="160" cy="260" r="2" fill={accent} fillOpacity="0.3" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1.2s' }} />

                {/* Data flow particles */}
                <circle cx="120" cy="165" r="3" fill={accent} fillOpacity="0.6" style={{ animation: 'floatY 2s ease-in-out infinite' }} />
                <circle cx="200" cy="165" r="3" fill={accent} fillOpacity="0.6" style={{ animation: 'floatYReverse 2.3s ease-in-out infinite 0.4s' }} />
            </svg>
        )
    }

    if (slide === 'slide2') {
        // Platform Capabilities — central hub with 4 capability nodes
        return (
            <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
                {/* Outer glow ring */}
                <circle cx="160" cy="140" r="115" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 5s ease-in-out infinite' }} />

                {/* Central hub */}
                <g style={{ animation: 'floatY 5s ease-in-out infinite' }}>
                    <circle cx="160" cy="140" r="40" fill={accent} fillOpacity="0.2" stroke={accent} strokeWidth="2.5" />
                    <text x="160" y="136" textAnchor="middle" fill={accent} fontSize="10" fontWeight="700">AnySync</text>
                    <text x="160" y="150" textAnchor="middle" fill={accent} fontSize="8" fillOpacity="0.7">PLATFORM</text>
                </g>

                {/* Connection lines to nodes */}
                <line x1="160" y1="100" x2="160" y2="60" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
                <line x1="200" y1="140" x2="240" y2="140" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
                <line x1="160" y1="180" x2="160" y2="220" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />
                <line x1="120" y1="140" x2="80" y2="140" stroke={accent} strokeWidth="1.5" strokeDasharray="4 4" strokeOpacity="0.5" />

                {/* Top node — Conversational BI */}
                <g style={{ animation: 'floatYReverse 3.5s ease-in-out infinite' }}>
                    <circle cx="160" cy="45" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                    <rect x="148" y="36" width="16" height="10" rx="4" fill={accent} fillOpacity="0.6" />
                    <circle cx="156" cy="50" r="4" fill={accent} fillOpacity="0.4" />
                    <text x="160" y="62" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">BI</text>
                </g>

                {/* Right node — AI Accountant */}
                <g style={{ animation: 'floatY 3.8s ease-in-out infinite 0.3s' }}>
                    <circle cx="255" cy="140" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                    <rect x="245" y="130" width="14" height="18" rx="3" fill={accent} fillOpacity="0.5" />
                    <path d="M248 140 l3 3 l6 -6" stroke={accent} strokeWidth="2" fill="none" />
                    <text x="255" y="158" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">ACC</text>
                </g>

                {/* Bottom node — Smart Reports */}
                <g style={{ animation: 'floatYReverse 3.5s ease-in-out infinite 0.6s' }}>
                    <circle cx="160" cy="235" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                    <rect x="148" y="225" width="6" height="15" rx="2" fill={accent} fillOpacity="0.4" />
                    <rect x="157" y="230" width="6" height="10" rx="2" fill={accent} fillOpacity="0.6" />
                    <rect x="166" y="222" width="6" height="18" rx="2" fill={accent} fillOpacity="0.8" />
                    <text x="160" y="253" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">RPT</text>
                </g>

                {/* Left node — Deep Analytics */}
                <g style={{ animation: 'floatY 3.8s ease-in-out infinite 0.9s' }}>
                    <circle cx="65" cy="140" r="28" fill={nodeColor} stroke={accent} strokeWidth="1.5" />
                    <circle cx="65" cy="135" r="10" stroke={accent} strokeWidth="1.5" fill="none" />
                    <circle cx="65" cy="135" r="4" fill={accent} fillOpacity="0.5" />
                    <circle cx="60" cy="130" r="2" fill={accent} fillOpacity="0.6" />
                    <circle cx="72" cy="138" r="2" fill={accent} fillOpacity="0.6" />
                    <text x="65" y="158" textAnchor="middle" fill={accent} fontSize="7" fontWeight="600">ANL</text>
                </g>

                {/* Data flow particles */}
                <circle cx="160" cy="80" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatY 2s ease-in-out infinite' }} />
                <circle cx="215" cy="140" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatYReverse 2.3s ease-in-out infinite' }} />
                <circle cx="160" cy="200" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatY 2.1s ease-in-out infinite 0.2s' }} />
                <circle cx="105" cy="140" r="3" fill={accent} fillOpacity="0.7" style={{ animation: 'floatYReverse 2.5s ease-in-out infinite 0.4s' }} />

                {/* Sparkle accents */}
                <circle cx="40" cy="60" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
                <circle cx="280" cy="70" r="2" fill={accent} fillOpacity="0.35" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.5s' }} />
                <circle cx="280" cy="220" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1s' }} />
                <circle cx="40" cy="210" r="2" fill={accent} fillOpacity="0.3" style={{ animation: 'pulseGlow 3.2s ease-in-out infinite 0.8s' }} />
            </svg>
        )
    }

    // slide3 — Trust Dashboard with outcomes and compliance
    return (
        <svg viewBox="0 0 320 280" fill="none" className="w-full h-auto max-w-[320px]" aria-hidden="true">
            {/* Outer glow ring */}
            <circle cx="160" cy="140" r="115" fill={accentLight} fillOpacity="0.08" style={{ animation: 'pulseGlow 5s ease-in-out infinite' }} />

            {/* Dashboard frame */}
            <g style={{ animation: 'floatY 5s ease-in-out infinite' }}>
                <rect x="50" y="55" rx="14" width="220" height="165" fill={nodeColor} stroke={strokeColor} strokeWidth="1.5" />

                {/* Title bar */}
                <rect x="50" y="55" rx="14" width="220" height="26" fill={accent} fillOpacity="0.1" />
                <circle cx="68" cy="68" r="4" fill={errorColor} fillOpacity="0.7" />
                <circle cx="82" cy="68" r="4" fill={warningColor} fillOpacity="0.7" />
                <circle cx="96" cy="68" r="4" fill={successColor} fillOpacity="0.7" />
                <text x="160" y="72" textAnchor="middle" fill={accent} fontSize="9" fontWeight="600" fillOpacity="0.8">AnySync Dashboard</text>
            </g>

            {/* Floating metric cards */}
            <g style={{ animation: 'floatYReverse 3s ease-in-out infinite' }}>
                <rect x="30" y="225" rx="8" width="65" height="30" fill={nodeColor} stroke={accent} strokeWidth="1" />
                <text x="62" y="238" textAnchor="middle" fill={successColor} fontSize="10" fontWeight="700">↑ +127%</text>
                <text x="62" y="250" textAnchor="middle" fill={accent} fontSize="7" fillOpacity="0.7">Growth</text>
            </g>

            <g style={{ animation: 'floatY 3.5s ease-in-out infinite 0.5s' }}>
                <rect x="225" y="35" rx="8" width="65" height="30" fill={nodeColor} stroke={accent} strokeWidth="1" />
                <text x="257" y="48" textAnchor="middle" fill={accent} fontSize="10" fontWeight="700">₹2.1Cr</text>
                <text x="257" y="60" textAnchor="middle" fill={accent} fontSize="7" fillOpacity="0.7">Saved</text>
            </g>

            {/* Bar chart */}
            <g>
                <rect x="70" y="175" width="18" height="30" rx="3" fill={accent} fillOpacity="0.25" />
                <rect x="95" y="160" width="18" height="45" rx="3" fill={accent} fillOpacity="0.35" />
                <rect x="120" y="145" width="18" height="60" rx="3" fill={accent} fillOpacity="0.5" />
                <rect x="145" y="130" width="18" height="75" rx="3" fill={accent} fillOpacity="0.65" />
                <rect x="170" y="115" width="18" height="90" rx="3" fill={accent} fillOpacity="0.8" />
                <rect x="195" y="100" width="18" height="105" rx="3" fill={accent} fillOpacity="0.95" />
                <rect x="220" y="90" width="18" height="115" rx="3" fill={accent} />
            </g>

            {/* Trend line */}
            <polyline points="79,165 104,150 129,135 154,120 179,105 204,92 229,82" stroke={successColor} strokeWidth="2.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />

            {/* Trust badges */}
            <g style={{ animation: 'floatY 4s ease-in-out infinite 0.3s' }}>
                <rect x="105" y="240" rx="6" width="50" height="22" fill={nodeColor} stroke={successColor} strokeWidth="1" />
                <text x="130" y="255" textAnchor="middle" fill={successColor} fontSize="8" fontWeight="600">HIPAA</text>
            </g>

            <g style={{ animation: 'floatYReverse 4s ease-in-out infinite 0.6s' }}>
                <rect x="165" y="240" rx="6" width="50" height="22" fill={nodeColor} stroke={accent} strokeWidth="1" />
                <text x="190" y="255" textAnchor="middle" fill={accent} fontSize="8" fontWeight="600">99.9%</text>
            </g>

            {/* Sparkle accents */}
            <circle cx="35" cy="90" r="2.5" fill={accent} fillOpacity="0.4" style={{ animation: 'pulseGlow 3s ease-in-out infinite' }} />
            <circle cx="285" cy="170" r="2" fill={accent} fillOpacity="0.35" style={{ animation: 'pulseGlow 3.5s ease-in-out infinite 0.5s' }} />
            <circle cx="160" cy="270" r="2" fill={successColor} fillOpacity="0.5" style={{ animation: 'pulseGlow 4s ease-in-out infinite 1s' }} />
        </svg>
    )
}

export function SectorIcon({ type, className = 'w-7 h-7' }: { type: string; className?: string }) {
    switch (type) {
        case 'hospitals':
            return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5} aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 3.75h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008z" /></svg>
        case 'labs':
            return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5} aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23-.693L5 14.5m14.8.8l1.402 1.402c1.232 1.232.65 3.318-1.067 3.611A48.309 48.309 0 0112 21c-2.773 0-5.491-.235-8.135-.687-1.718-.293-2.3-2.379-1.067-3.61L5 14.5" /></svg>
        case 'pharmacies':
            return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5} aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 00-3.7-3.7 48.678 48.678 0 00-7.324 0 4.006 4.006 0 00-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3l-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 003.7 3.7 48.656 48.656 0 007.324 0 4.006 4.006 0 003.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3l-3 3" /></svg>
        case 'clinics':
            return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5} aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>
        default:
            return null
    }
}
