import React from 'react';
import { useTranslation } from 'react-i18next';

interface QuerySuggestionsProps {
  suggestions: string[];
  onSuggestionClick: (suggestion: string) => void;
}

export const QuerySuggestions: React.FC<QuerySuggestionsProps> = ({
  suggestions,
  onSuggestionClick,
}) => {
  const { t } = useTranslation('chat');

  return (
    <div className="w-full max-w-2xl" role="list" aria-label={t('suggestions.ariaLabel', 'Suggested queries')}>
      <h3 className="text-sm font-medium text-slate-600 dark:text-slate-400 mb-3">
        {t('suggestions.title')}
      </h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {suggestions.map((suggestion, index) => (
          <button
            key={index}
            onClick={() => onSuggestionClick(suggestion)}
            className="liquid-glass-button-prominent group flex items-center gap-3 px-4 py-3 min-h-[44px] text-start liquid-glass-hover-lift focus-visible:outline-[3px] focus-visible:outline-offset-[3px] focus-visible:outline-[var(--color-trust-blue)] dark:focus-visible:outline-cyan-400"
            role="listitem"
            aria-label={t('suggestions.suggestionAriaLabel', 'Suggestion: {{suggestion}}', { suggestion })}
          >
            <div className="shrink-0 w-8 h-8 rounded-lg flex items-center justify-center transition-colors" style={{
              background: 'linear-gradient(135deg, rgba(39,80,168,0.1) 0%, rgba(24,146,157,0.1) 100%)'
            }}>
              <SuggestionIcon index={index} />
            </div>
            <span className="text-sm font-medium text-slate-700 dark:text-slate-300 group-hover:text-[#2750a8] dark:group-hover:text-[#18929d] transition-colors">
              {suggestion}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
};

// Different icons for different suggestion types with brand colors
const SuggestionIcon: React.FC<{ index: number }> = ({ index }) => {
  const iconColor = '#2750a8'; // Brand deep blue
  const icons = [
    // Revenue/Money
    <svg key="0" className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke={iconColor} strokeWidth={2} aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>,
    // Chart/Departments
    <svg key="1" className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="#18929d" strokeWidth={2} aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
    </svg>,
    // Trend/Patients
    <svg key="2" className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke={iconColor} strokeWidth={2} aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
    </svg>,
    // Inventory/Box
    <svg key="3" className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="#18929d" strokeWidth={2} aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
    </svg>,
  ];

  return icons[index % icons.length];
};

export default QuerySuggestions;
