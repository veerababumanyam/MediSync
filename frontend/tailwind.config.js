/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    darkMode: 'class',
    theme: {
        extend: {
            colors: {
                'logo-blue': '#2750a8',
                'logo-teal': '#18929d',
                'midnight-navy': '#0f172a',
                'text-primary': 'var(--color-text-primary)',
                'text-secondary': 'var(--color-text-secondary)',
                'text-tertiary': 'var(--color-text-tertiary)',
                'text-muted': 'var(--color-text-muted)',
                'on-brand': 'var(--color-text-on-brand)',
                'surface-glass': 'var(--color-surface-glass)',
                'surface-glass-strong': 'var(--color-surface-glass-strong)',
                'border-default': 'var(--color-border-default)',
                'border-glass': 'var(--color-border-glass)',
                'action-primary': 'var(--color-action-primary-bg)',
                'action-primary-hover': 'var(--color-action-primary-bg-hover)',
            },
            boxShadow: {
                'glass-sm': '0 8px 32px rgba(0, 0, 0, 0.3)',
                'glass-elevated': '0 20px 60px rgba(0, 0, 0, 0.4)',
                'glass-ios': '0 4px 24px rgba(0, 0, 0, 0.05)',
                'glass-ios-dark': '0 4px 24px rgba(0, 0, 0, 0.3)',
            }
        },
    },
    plugins: [],
}
