import React from 'react';
import { useTranslation } from 'react-i18next';

interface QuerySuggestionsProps {
  suggestions: string[];
  onSuggestionClick: (suggestion: string) => void;
  isDark?: boolean;
}

export const QuerySuggestions: React.FC<QuerySuggestionsProps> = ({
  suggestions,
  onSuggestionClick,
  isDark = true,
}) => {
  const { t } = useTranslation('chat');

  return (
    <div className="w-full max-w-2xl">
      <h3 className={`text-sm font-medium mb-3 ${isDark ? 'text-slate-400' : 'text-slate-500'
        }`}>
        {t('suggestions.title')}
      </h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {suggestions.map((suggestion, index) => (
          <button
            key={index}
            onClick={() => onSuggestionClick(suggestion)}
            className={`group flex items-center gap-3 p-4 rounded-xl transition-all duration-300 hover:-translate-y-0.5 text-left ${isDark
                ? 'bg-white/5 border border-white/10 hover:bg-white/10 hover:border-white/20'
                : 'bg-white border border-slate-200 hover:border-blue-200 hover:shadow-md'
              }`}
          >
            <div className={`flex-shrink-0 w-10 h-10 rounded-xl flex items-center justify-center transition-all duration-300 ${isDark
                ? 'bg-white/10 group-hover:bg-gradient-to-br group-hover:from-blue-500/20 group-hover:to-cyan-400/20'
                : 'bg-slate-100 group-hover:bg-blue-50'
              }`}>
              <SuggestionIcon index={index} isDark={isDark} />
            </div>
            <span className={`text-sm font-medium transition-colors ${isDark
                ? 'text-slate-300 group-hover:text-white'
                : 'text-slate-700 group-hover:text-blue-600'
              }`}>
              {suggestion}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
};

// Different icons for different suggestion types
const SuggestionIcon: React.FC<{ index: number; isDark?: boolean }> = ({ index, isDark = true }) => {
  const colorClass = isDark
    ? 'text-slate-400 group-hover:text-blue-400'
    : 'text-slate-500 group-hover:text-blue-500';

  const icons = [
    // Revenue/Money
    <svg key="0" className={`w-5 h-5 ${colorClass} transition-colors`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>,
    // Chart/Departments
    <svg key="1" className={`w-5 h-5 ${colorClass} transition-colors`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
    </svg>,
    // Trend/Patients
    <svg key="2" className={`w-5 h-5 ${colorClass} transition-colors`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
    </svg>,
    // Inventory/Box
    <svg key="3" className={`w-5 h-5 ${colorClass} transition-colors`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
    </svg>,
  ];

  return icons[index % icons.length];
};

export default QuerySuggestions;
