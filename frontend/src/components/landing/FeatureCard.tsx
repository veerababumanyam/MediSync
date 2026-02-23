import React from 'react'

export interface FeatureCardProps {
    icon: React.ReactNode
    gradientLight: string
    gradientDark: string
    iconColorLight: string
    iconColorDark: string
    shadowLight?: string
    borderLight?: string
    borderDark?: string
    title: string
    description: string
    delay?: string
    isDark: boolean
}

export function FeatureCard({
    icon,
    gradientLight,
    gradientDark,
    iconColorLight,
    iconColorDark,
    shadowLight = '',
    borderLight = '',
    borderDark = '',
    title,
    description,
    delay = '',
    isDark,
}: FeatureCardProps) {
    const delayMs = delay === 'delay-1' ? '0ms' :
        delay === 'delay-2' ? '100ms' :
            delay === 'delay-3' ? '200ms' : '300ms'

    return (
        <div
            className={`liquid-glass-pronounced rounded-2xl overflow-hidden group animate-fade-in-up p-6 h-full liquid-glass-hover-lift liquid-focus-gold ${delay}`}
            style={{ animationDelay: delayMs }}
            role="article"
        >
            <div className={`w-12 h-12 rounded-xl flex items-center justify-center mb-5 group-hover:scale-110 transition-transform duration-300 relative z-10
                ${isDark
                    ? `bg-gradient-to-br ${gradientDark} ${iconColorDark} ${borderDark}`
                    : `bg-gradient-to-br ${gradientLight} ${iconColorLight} ${shadowLight} ${borderLight}`
                }`}
                aria-hidden="true"
            >
                {icon}
                {/* Add subtle glow effect to icon */}
                <div className="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                    style={{
                        background: 'inherit',
                        filter: 'blur(8px)',
                        zIndex: -1
                    }}
                />
            </div>
            <h3 className="text-lg font-semibold mb-3 text-primary">
                {title}
            </h3>
            <p className="text-sm leading-relaxed text-secondary">
                {description}
            </p>
        </div>
    )
}
