import React, { useState, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { PinnedChartCard } from './PinnedChartCard';
import { ChartPinDialog } from './ChartPinDialog';
import { useLocale } from '../../hooks/useLocale';
import { dashboardApi, PinnedChart } from '../../services/api';
import { LoadingSpinner } from '../common/LoadingSpinner';

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
      <div className="flex items-center justify-center h-64">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
          {t('title')}
        </h2>
        <button
          onClick={() => setShowPinDialog(true)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          {t('pinChart')}
        </button>
      </div>

      {/* Error Message */}
      {error && (
        <div className="p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg">
          {error}
        </div>
      )}

      {/* Empty State */}
      {charts.length === 0 && !error && (
        <div className="text-center py-12">
          <svg
            className="w-16 h-16 mx-auto text-gray-400 dark:text-gray-600"
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
          <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
            {t('empty.title')}
          </h3>
          <p className="mt-2 text-gray-500 dark:text-gray-400">
            {t('empty.description')}
          </p>
          <button
            onClick={() => setShowPinDialog(true)}
            className="mt-4 inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
          >
            {t('pinFirstChart')}
          </button>
        </div>
      )}

      {/* Charts Grid */}
      {charts.length > 0 && (
        <div
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
        </div>
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
