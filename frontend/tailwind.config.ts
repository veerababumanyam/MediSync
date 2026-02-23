import type { Config } from 'tailwind-css'

/**
 * Tailwind CSS Configuration for AnySync
 *
 * All design values reference CSS custom properties from design-tokens.css
 * This ensures a single source of truth for the entire design system.
 *
 * @version 2.0.0
 * @lastUpdated 2026-02-23
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
        // ============================================
        // BRAND COLORS (from design-tokens.css)
        // ============================================
        'brand-primary': 'var(--brand-primary)',
        'brand-primary-light': 'var(--brand-primary-light)',
        'brand-primary-dark': 'var(--brand-primary-dark)',

        'brand-secondary': 'var(--brand-secondary)',
        'brand-secondary-light': 'var(--brand-secondary-light)',
        'brand-secondary-dark': 'var(--brand-secondary-dark)',

        'brand-navy': 'var(--brand-navy)',

        // Logo colors (aliases)
        'logo-blue': 'var(--brand-primary)',
        'logo-teal': 'var(--brand-secondary)',

        // Backward compatible aliases
        'trust-blue': 'var(--brand-primary)',
        'growth-teal': 'var(--brand-secondary)',
        'midnight-navy': 'var(--brand-navy)',

        // ============================================
        // SEMANTIC COLORS
        // ============================================
        'success': 'var(--color-success)',
        'success-light': 'var(--color-success-light)',
        'warning': 'var(--color-warning)',
        'warning-light': 'var(--color-warning-light)',
        'error': 'var(--color-error)',
        'error-light': 'var(--color-error-light)',
        'info': 'var(--color-info)',
        'info-light': 'var(--color-info-light)',

        // iOS System colors
        'system-blue': '#007AFF',
        'system-purple': '#5856D6',
        'system-pink': '#FF2D55',
        'system-orange': '#FF9500',
        'system-green': '#34C759',
        'system-red': '#FF3B30',
        'system-teal': '#5AC8FA',

        // ============================================
        // GLASS COLORS
        // ============================================
        'glass-bg-subtle': 'var(--glass-bg-subtle)',
        'glass-bg-light': 'var(--glass-bg-light)',
        'glass-bg-medium': 'var(--glass-bg-medium)',
        'glass-bg-heavy': 'var(--glass-bg-heavy)',
        'glass-bg-pronounced': 'var(--glass-bg-pronounced)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        arabic: ['Cairo', 'Noto Sans Arabic', 'Tajawal', 'sans-serif'],
        mono: ['JetBrains Mono', 'SF Mono', 'ui-monospace', 'monospace'],
      },
      spacing: {
        // Token-based spacing from design-tokens.css
        'section-sm': 'var(--section-padding-sm)',
        'section-md': 'var(--section-padding-md)',
        'section-lg': 'var(--section-padding-lg)',
        'section-xl': 'var(--section-padding-xl)',
      },
      boxShadow: {
        // Elevation shadows from design-tokens.css
        'ambient': 'var(--shadow-ambient)',
        'raised': 'var(--shadow-raised)',
        'floating': 'var(--shadow-floating)',
        'elevated': 'var(--shadow-elevated)',
        'glass-floating': 'var(--shadow-glass-floating)',
        'inset': 'var(--shadow-inset)',

        // Colored glow shadows
        'glow-primary': 'var(--shadow-glow-primary)',
        'glow-secondary': 'var(--shadow-glow-secondary)',
        'glow-success': 'var(--shadow-glow-success)',
        'glow-warning': 'var(--shadow-glow-warning)',
        'glow-error': 'var(--shadow-glow-error)',
      },
      borderRadius: {
        'dynamic': 'var(--radius-dynamic)',
      },
      backdropBlur: {
        'xs': '10px',
        'glass': '40px',
        'glass-elevated': '60px',
        'glass-pronounced': '80px',
      },
      animation: {
        'float-slow': 'float 20s ease-in-out infinite',
        'float-medium': 'float 25s ease-in-out infinite',
        'float-fast': 'float 18s ease-in-out infinite',
        'shine': 'shineSlide 6s ease-in-out infinite',
        'fade-in-up': 'fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards',
        'pulse-slow': 'pulseGlow 4s ease-in-out infinite',
        'liquid-fade-in': 'liquidFadeIn 0.5s ease-out forwards',
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
        pulseGlow: {
          '0%, 100%': {
            boxShadow: '0 0 0 1px rgba(255, 255, 255, 0.1), 0 8px 32px rgba(39, 80, 168, 0.1)',
          },
          '50%': {
            boxShadow: '0 0 0 1px rgba(255, 255, 255, 0.2), 0 8px 32px rgba(39, 80, 168, 0.2), 0 0 40px rgba(39, 80, 168, 0.1)',
          },
        },
        liquidFadeIn: {
          '0%': {
            opacity: '0',
            backdropFilter: 'blur(0px)',
            transform: 'translateY(10px) scale(0.98)',
          },
          '100%': {
            opacity: '1',
            backdropFilter: 'blur(60px)',
            transform: 'translateY(0) scale(1)',
          },
        },
      },
    },
  },
  plugins: [],
} satisfies Config
