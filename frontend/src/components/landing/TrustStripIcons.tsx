import type { FC, SVGProps } from 'react'

interface TrustIconProps extends SVGProps<SVGSVGElement> {
  isDark: boolean
}

const IconHIMS: FC<TrustIconProps> = ({ isDark, ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <rect x="4" y="4" width="16" height="16" rx="3" stroke="currentColor" strokeWidth="1.6" />
    <path d="M12 8v8M8 12h8" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
  </svg>
)

const IconLIMS: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <path d="M9 4h6M10.5 4v4.8l-3.3 5.4a3.2 3.2 0 0 0 2.7 4.8h4.2a3.2 3.2 0 0 0 2.7-4.8l-3.3-5.4V4" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
    <path d="M8.8 14h6.4" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" opacity="0.7" />
  </svg>
)

const IconTally: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <rect x="3.5" y="4.5" width="17" height="15" rx="2.5" stroke="currentColor" strokeWidth="1.4" opacity="0.35" />
    <path d="M7 8h10M12 8v8" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    <circle cx="12" cy="16" r="1.3" fill="currentColor" />
  </svg>
)

const IconSQL: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <ellipse cx="12" cy="6.6" rx="6.5" ry="2.6" stroke="currentColor" strokeWidth="1.6" />
    <path d="M5.5 6.6v7.8c0 1.5 2.9 2.6 6.5 2.6s6.5-1.1 6.5-2.6V6.6" stroke="currentColor" strokeWidth="1.6" />
    <path d="M5.5 10.5c0 1.4 2.9 2.5 6.5 2.5s6.5-1.1 6.5-2.5" stroke="currentColor" strokeWidth="1.4" opacity="0.7" />
  </svg>
)

const IconCustomAPI: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <path d="M9 8.2h6M9 15.8h6M8.5 12h7" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
    <path d="M4 12a4 4 0 0 1 4-4M4 12a4 4 0 0 0 4 4M20 12a4 4 0 0 0-4-4M20 12a4 4 0 0 1-4 4" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
  </svg>
)

const IconOracle: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <rect x="4" y="7" width="16" height="10" rx="5" stroke="currentColor" strokeWidth="2.2" />
  </svg>
)

const IconSAP: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <path d="M4 8.5L8.2 5h11.8v14H4z" stroke="currentColor" strokeWidth="1.6" />
    <path d="M7 15.2c.5.4 1.1.6 1.8.6.9 0 1.5-.3 1.5-.9s-.3-.8-1.3-1.1c-1.3-.4-2.1-1.1-2.1-2.2 0-1.3 1.1-2.2 2.8-2.2.8 0 1.4.2 1.8.4M12.5 15.6V9.6h1.9c1.8 0 2.8 1 2.8 2.9 0 2-1.1 3.1-2.9 3.1h-1.8z" stroke="currentColor" strokeWidth="1.1" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
)

const IconRestGraphQL: FC<TrustIconProps> = ({ ...props }) => (
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true" {...props}>
    <path d="M5.5 7.5h2.2M5.5 16.5h2.2M5.5 7.5c-.8.8-1.2 2.2-1.2 4.5s.4 3.7 1.2 4.5M18.5 7.5h-2.2M18.5 16.5h-2.2M18.5 7.5c.8.8 1.2 2.2 1.2 4.5s-.4 3.7-1.2 4.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
    <path d="M12 7l4 2.3v5.4L12 17l-4-2.3V9.3L12 7z" stroke="currentColor" strokeWidth="1.4" opacity="0.8" />
    <circle cx="12" cy="7" r="1.1" fill="currentColor" />
    <circle cx="16" cy="9.3" r="1.1" fill="currentColor" />
    <circle cx="16" cy="14.7" r="1.1" fill="currentColor" />
    <circle cx="12" cy="17" r="1.1" fill="currentColor" />
    <circle cx="8" cy="14.7" r="1.1" fill="currentColor" />
    <circle cx="8" cy="9.3" r="1.1" fill="currentColor" />
  </svg>
)

const TRUST_STRIP_ICONS = {
  item1: IconHIMS,
  item2: IconLIMS,
  item3: IconTally,
  item4: IconSQL,
  item5: IconCustomAPI,
  item6: IconOracle,
  item7: IconSAP,
  item8: IconRestGraphQL,
} as const

interface TrustStripIconProps extends TrustIconProps {
  itemKey: keyof typeof TRUST_STRIP_ICONS
}

const TRUST_BRAND_COLORS: Record<keyof typeof TRUST_STRIP_ICONS, { light: string; dark: string }> = {
  item1: { light: '#0056D2', dark: '#60A5FA' }, // HIMS
  item2: { light: '#0E7490', dark: '#22D3EE' }, // LIMS
  item3: { light: '#1E3A8A', dark: '#93C5FD' }, // Tally
  item4: { light: '#7C3AED', dark: '#C4B5FD' }, // SQL
  item5: { light: '#0F766E', dark: '#5EEAD4' }, // Custom API
  item6: { light: '#C2410C', dark: '#FDBA74' }, // Oracle
  item7: { light: '#15803D', dark: '#86EFAC' }, // SAP
  item8: { light: '#BE185D', dark: '#F9A8D4' }, // REST/GraphQL
}

export function getTrustBrandColor(itemKey: keyof typeof TRUST_STRIP_ICONS, isDark: boolean): string {
  const palette = TRUST_BRAND_COLORS[itemKey]
  return isDark ? palette.dark : palette.light
}

export function TrustStripIcon({ itemKey, ...props }: TrustStripIconProps) {
  const IconComponent = TRUST_STRIP_ICONS[itemKey]
  return <IconComponent {...props} />
}
