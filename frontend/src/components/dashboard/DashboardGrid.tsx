import React, { useState, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { PinnedChartCard } from './PinnedChartCard';
import { ChartPinDialog } from './ChartPinDialog';
import { useLocale } from '../../hooks/useLocale';
import { LoadingSkeleton } from '../ui';
import { LiquidGlassCard, GlassBrandCard } from '../ui/LiquidGlassCard';
import { ButtonPrimary } from '../ui/LiquidGlassButton';
import { StaggerChildren, FadeIn } from '../animations';
import { dashboardApi } from '../../services/api';
import type { PinnedChart } from '../../services/api';

interface DashboardGridProps {
  onChartClick?: (chart: PinnedChart) => void;
  className?: string;
}

export const DashboardGrid: React.FC<DashboardGridProps> = ({
  onChartClick,
  className = '',
}) => {
  const { t } = useTranslation('dashboard');
  const { locale } = useLocale();

  const [charts, setCharts] = useState<PinnedChart[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showPinDialog, setShowPinDialog] = useState(false);

  // Load charts on mount
  useEffect(() => {
    loadCharts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadCharts = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await dashboardApi.getCharts();
      setCharts(response.filter(c => c.isActive));
    } catch (err) {
      console.error('Failed to load charts:', err);
      setError(t('error.loadCharts'));
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteChart = useCallback(async (chartId: string) => {
    try {
      await dashboardApi.deleteChart(chartId);
      setCharts(prev => prev.filter(c => c.id !== chartId));
    } catch (err) {
      console.error('Failed to delete chart:', err);
      setError(t('error.deleteChart'));
    }
  }, [t]);

  const handleRefreshChart = useCallback(async (chartId: string) => {
    try {
      const updated = await dashboardApi.refreshChart(chartId);
      setCharts(prev => prev.map(c => c.id === chartId ? updated : c));
    } catch (err) {
      console.error('Failed to refresh chart:', err);
      setError(t('error.refreshChart'));
    }
  }, [t]);

  const handleToggleChart = useCallback(async (chartId: string, isActive: boolean) => {
    try {
      await dashboardApi.updateChart(chartId, { isActive });
      if (!isActive) {
        setCharts(prev => prev.filter(c => c.id !== chartId));
      }
    } catch (err) {
      console.error('Failed to toggle chart:', err);
    }
  }, []);

  const handlePinNewChart = useCallback(async (chartData: Partial<PinnedChart>) => {
    try {
      const newChart = await dashboardApi.pinChart(chartData);
      setCharts(prev => [...prev, newChart]);
      setShowPinDialog(false);
    } catch (err) {
      console.error('Failed to pin chart:', err);
      setError(t('error.pinChart'));
    }
  }, [t]);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <LoadingSkeleton variant="card" className="h-24" />
        <div className="grid gap-4" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))' }}>
          <LoadingSkeleton variant="chart" />
          <LoadingSkeleton variant="chart" />
          <LoadingSkeleton variant="chart" />
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <FadeIn>
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold liquid-text-primary">
              {t('title')}
            </h2>
            <p className="text-sm liquid-text-secondary mt-1">
              Your pinned business insights
            </p>
          </div>
          <ButtonPrimary
            onClick={() => setShowPinDialog(true)}
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            {t('pinChart')}
          </ButtonPrimary>
        </div>
      </FadeIn>

      {/* Error Message */}
      {error && (
        <FadeIn>
          <LiquidGlassCard intensity="light" className="p-4 border-l-4 border-l-red-500">
            <div className="flex items-center gap-3">
              <svg className="w-5 h-5 text-red-500 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span className="text-red-400 flex-1">{error}</span>
            </div>
          </LiquidGlassCard>
        </FadeIn>
      )}

      {/* Empty State */}
      {charts.length === 0 && !error && (
        <FadeIn>
          <LiquidGlassCard intensity="medium" elevation="raised" className="text-center py-16">
            <div className="w-20 h-20 rounded-2xl liquid-glass flex items-center justify-center mx-auto mb-6">
              <svg
                className="w-10 h-5 liquid-text-secondary"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                />
              </svg>
            </div>
            <h3 className="mt-4 text-lg font-semibold liquid-text-primary">
              {t('empty.title')}
            </h3>
            <p className="mt-2 liquid-text-secondary max-w-sm mx-auto">
              {t('empty.description')}
            </p>
            <ButtonPrimary
              onClick={() => setShowPinDialog(true)}
              className="mt-6"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
              {t('pinFirstChart')}
            </ButtonPrimary>
          </LiquidGlassCard>
        </FadeIn>
      )}

      {/* Charts Grid */}
      {charts.length > 0 && (
        <StaggerChildren
          className="grid gap-4"
          style={{
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
          }}
        >
          {charts.map(chart => (
            <PinnedChartCard
              key={chart.id}
              chart={chart}
              locale={locale}
              onDelete={() => handleDeleteChart(chart.id)}
              onRefresh={() => handleRefreshChart(chart.id)}
              onToggle={(active) => handleToggleChart(chart.id, active)}
              onClick={() => onChartClick?.(chart)}
            />
          ))}
        </StaggerChildren>
      )}

      {/* Pin Dialog */}
      {showPinDialog && (
        <ChartPinDialog
          onClose={() => setShowPinDialog(false)}
          onPin={handlePinNewChart}
          locale={locale}
        />
      )}
    </div>
  );
};

export default DashboardGrid;
