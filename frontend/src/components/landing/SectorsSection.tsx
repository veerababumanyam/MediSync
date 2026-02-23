
import { useTranslation } from 'react-i18next'
import { SectorIcon } from './icons'

export interface SectorsSectionProps {
    isDark: boolean
}

export function SectorsSection({ isDark }: SectorsSectionProps) {
    const { t } = useTranslation()
    const sectors = ['hospitals', 'labs', 'pharmacies', 'clinics']

    return (
        <section id="sectors" className="pt-24 mb-24 animate-fade-in-up delay-2" aria-labelledby="sectors-heading">
            <div className="text-center mb-12">
                <h2 id="sectors-heading" className="text-3xl sm:text-4xl font-extrabold mb-5 text-primary">
                    {t('sectors.title', 'Dominating Complexity Across Every Healthcare Sector')}
                </h2>
                <p className="text-base sm:text-lg max-w-3xl mx-auto leading-relaxed text-secondary">
                    {t('sectors.description', "We don't just understand data; we understand the business of healthcare. Our tailored Agentic AI bridges seamlessly adapt to the unique reporting, compliance, and velocity requirements of your specific vertical.")}
                </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 items-start">
                {sectors.map((sector) => (
                    <div key={sector} className="liquid-glass-content-card rounded-2xl p-6 h-full" role="article">
                        <div className={`w-12 h-12 rounded-xl flex items-center justify-center mb-5 transition-all duration-300
              ${isDark ? 'bg-gradient-to-br from-brand-primary-alpha to-brand-secondary-alpha text-brand-secondary border border-brand-secondary/20' : 'bg-gradient-to-br from-brand-primary-alpha to-brand-secondary-alpha text-brand-primary border-2 border-brand-primary-alpha shadow-md shadow-brand-primary-alpha'}`}
                            aria-hidden="true"
                        >
                            <SectorIcon type={sector} className="w-6 h-6" />
                        </div>
                        <h3 className="text-lg font-semibold mb-3 text-primary">
                            {t(`sectors.${sector}.title`)}
                        </h3>
                        <p className="text-sm leading-relaxed text-secondary">
                            {t(`sectors.${sector}.description`)}
                        </p>
                    </div>
                ))}
            </div>
        </section>
    )
}
