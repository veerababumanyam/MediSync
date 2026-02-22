import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { LiquidGlassCard } from '../ui/LiquidGlassCard';
import { FadeIn } from '../animations';
import type { PinnedChart } from '../../services/api';
import { ChartRenderer } from '../chat/ChartRenderer';

interface PinnedChartCardProps {
  chart: PinnedChart;
  locale: string;
  onDelete: () => void;
  onRefresh: () => void;
  onToggle: (active: boolean) => void;
  onClick?: () => void;
}

export const PinnedChartCard: React.FC<PinnedChartCardProps> = ({
  chart,
  locale,
  onDelete,
  onRefresh,
  onToggle,
  onClick,
}) => {
  const { t } = useTranslation('dashboard');
  const [showMenu, setShowMenu] = useState(false);
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await onRefresh();
    setIsRefreshing(false);
    setShowMenu(false);
  };

  const formatLastRefreshed = (timestamp: string | null) => {
    if (!timestamp) return t('never');
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) return t('justNow');
    if (diffMins < 60) return t('minutesAgo', { count: diffMins });
    if (diffMins < 1440) return t('hoursAgo', { count: Math.floor(diffMins / 60) });
    return date.toLocaleDateString(locale === 'ar' ? 'ar-SA' : 'en-US');
  };

  return (
    <FadeIn>
      <LiquidGlassCard
        elevation="floating"
        hover="lift"
        className={`overflow-hidden ${onClick ? 'cursor-pointer' : ''}`}
        onClick={onClick}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-white/10">
          <h3 className="font-medium liquid-text-primary truncate">
            {chart.title}
          </h3>

          {/* Actions Menu */}
          <div className="relative">
            <button
              onClick={(e) => {
                e.stopPropagation();
                setShowMenu(!showMenu);
              }}
              className="p-1.5 rounded-lg hover:bg-white/10 transition-colors"
            >
              <svg className="w-5 h-5 liquid-text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </button>

            {showMenu && (
              <>
                <div
                  className="fixed inset-0 z-10"
                  onClick={(e) => {
                    e.stopPropagation();
                    setShowMenu(false);
                  }}
                />
                <LiquidGlassCard intensity="heavy" elevation="floating" className="absolute right-0 mt-1 w-48 rounded-lg z-20 p-1">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleRefresh();
                    }}
                    disabled={isRefreshing}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm liquid-text-primary hover:bg-white/10 rounded-lg transition-colors disabled:opacity-50"
                  >
                    <svg className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    {t('refresh')}
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onToggle(false);
                      setShowMenu(false);
                    }}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm liquid-text-primary hover:bg-white/10 rounded-lg transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    {t('unpin')}
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onDelete();
                      setShowMenu(false);
                    }}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    {t('delete')}
                  </button>
                </LiquidGlassCard>
              </>
            )}
          </div>
        </div>

        {/* Chart Content */}
        <div className="p-4">
          <ChartRenderer
            chartType={chart.chartType}
            data={chart.chartSpec}
            locale={locale}
          />
        </div>

        {/* Footer */}
        <div className="px-4 py-2 bg-white/5 backdrop-blur-sm border-t border-white/10 text-xs liquid-text-secondary">
          <div className="flex items-center justify-between">
            <span>{t('lastRefreshed')}: {formatLastRefreshed(chart.lastRefreshedAt)}</span>
            {chart.refreshInterval > 0 && (
              <span className="flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                {t('refreshesEvery', { minutes: chart.refreshInterval })}
              </span>
            )}
          </div>
        </div>
      </LiquidGlassCard>
    </FadeIn>
  );
};

export default PinnedChartCard;
