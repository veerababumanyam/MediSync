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
    return (
        <div className={`liquid-glass-content-card rounded-2xl overflow-hidden group animate-fade-in-up p-6 h-full ${delay}`} role="article">
            <div className={`w-12 h-12 rounded-xl flex items-center justify-center mb-5 group-hover:scale-110 transition-transform duration-300
                ${isDark
                    ? `bg-gradient-to-br ${gradientDark} ${iconColorDark} ${borderDark}`
                    : `bg-gradient-to-br ${gradientLight} ${iconColorLight} ${shadowLight} ${borderLight}`
                }`} aria-hidden="true">
                {icon}
            </div>
            <h3 className={`text-lg font-semibold mb-3 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                {title}
            </h3>
            <p className={`text-sm leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                {description}
            </p>
        </div>
    )
}
