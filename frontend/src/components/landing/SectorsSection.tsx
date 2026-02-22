
import { useTranslation } from 'react-i18next'
import { SectorIcon } from './icons'

export interface SectorsSectionProps {
    isDark: boolean
}

export function SectorsSection({ isDark }: SectorsSectionProps) {
    const { t } = useTranslation()
    const sectors = ['hospitals', 'labs', 'pharmacies', 'clinics']

    return (
        <section className="mb-24 animate-fade-in-up delay-2">
            <div className="text-center mb-12">
                <h2 className={`text-3xl sm:text-4xl font-extrabold mb-5 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                    {t('sectors.title', 'Dominating Complexity Across Every Healthcare Sector')}
                </h2>
                <p className={`text-lg sm:text-xl max-w-3xl mx-auto leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                    {t('sectors.description', "We don't just understand data; we understand the business of healthcare. Our tailored Agentic AI bridges seamlessly adapt to the unique reporting, compliance, and velocity requirements of your specific vertical.")}
                </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {sectors.map((sector) => (
                    <div key={sector} className="glass rounded-2xl p-6 transition-all duration-300 hover:-translate-y-2">
                        <div className={`w-14 h-14 rounded-xl flex items-center justify-center mb-5
              ${isDark ? 'bg-gradient-to-br from-blue-500/20 to-teal-400/20 text-teal-400 border border-teal-500/20' : 'bg-gradient-to-br from-blue-100 to-teal-100 text-blue-600 border-2 border-blue-200 shadow-md shadow-blue-500/15'}`}
                        >
                            <SectorIcon type={sector} />
                        </div>
                        <h3 className={`text-xl font-bold mb-3 ${isDark ? 'text-white' : 'text-slate-900'}`}>
                            {t(`sectors.${sector}.title`)}
                        </h3>
                        <p className={`text-sm leading-relaxed ${isDark ? 'text-slate-400' : 'text-slate-600'}`}>
                            {t(`sectors.${sector}.description`)}
                        </p>
                    </div>
                ))}
            </div>
        </section>
    )
}
