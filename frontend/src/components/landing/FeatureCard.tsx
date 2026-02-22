import React from 'react'

export interface FeatureCardProps {
    icon: React.ReactNode
    gradient: string
    shadowColor: string
    title: string
    description: string
    delay?: string
    isDark: boolean
}

export function FeatureCard({
    icon,
    gradient,
    shadowColor,
    title,
    description,
    delay = '',
    isDark,
}: FeatureCardProps) {
    return (
        <div className={`glass glass-shine rounded-2xl p-6 hover:-translate-y-1 hover:scale-[1.02] transition-all duration-300 group animate-fade-in-up ${delay}`}>
            <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${gradient} flex items-center justify-center mb-5 text-white shadow-lg ${shadowColor} group-hover:scale-110 transition-transform duration-300`}>
                {icon}
            </div>
            <h3 className={`text-lg font-semibold mb-2 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                {title}
            </h3>
            <p className={`text-sm leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                {description}
            </p>
        </div>
    )
}
