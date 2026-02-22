import type { Config } from 'tailwindcss'

/**
 * Tailwind CSS Configuration for MediSync
 *
 * Brand colors and design tokens from docs/DESIGN.md
 */
export default {
    content: [
        './index.html',
        './src/**/*.{ts,tsx}',
    ],
    darkMode: 'class',
    theme: {
        extend: {
            colors: {
                // Primary Brand Colors
                'trust-blue': '#0056D2',
                'growth-teal': '#00E8C6',
                'midnight-navy': '#0F172A',
                // iOS Accent Palette
                'system-blue': '#007AFF',
                'system-purple': '#5856D6',
                'system-pink': '#FF2D55',
                'system-orange': '#FF9500',
                'system-green': '#34C759',
                'system-red': '#FF3B30',
                'system-teal': '#5AC8FA',
            },
            fontFamily: {
                sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
                arabic: ['Cairo', 'Noto Sans Arabic', 'Tajawal', 'sans-serif'],
                mono: ['JetBrains Mono', 'SF Mono', 'ui-monospace', 'monospace'],
            },
            backdropBlur: {
                glass: '40px',
                'glass-elevated': '60px',
            },
            animation: {
                'float-slow': 'float 20s ease-in-out infinite',
                'float-medium': 'float 25s ease-in-out infinite',
                'float-fast': 'float 18s ease-in-out infinite',
                'shine': 'shineSlide 6s ease-in-out infinite',
                'fade-in-up': 'fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards',
            },
            keyframes: {
                float: {
                    '0%, 100%': { transform: 'translate(0, 0) scale(1)' },
                    '33%': { transform: 'translate(30px, -40px) scale(1.05)' },
                    '66%': { transform: 'translate(-20px, 20px) scale(0.95)' },
                },
                shineSlide: {
                    '0%': { backgroundPosition: '-200% 0' },
                    '100%': { backgroundPosition: '200% 0' },
                },
                fadeInUp: {
                    '0%': { opacity: '0', transform: 'translateY(20px)' },
                    '100%': { opacity: '1', transform: 'translateY(0)' },
                },
            },
        },
    },
    plugins: [],
} satisfies Config
