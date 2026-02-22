import React, { useEffect, useRef, useMemo, useState } from 'react';
import * as echarts from 'echarts';
import { useTranslation } from 'react-i18next';
import { formatNumber, formatDate } from '../../utils/localeUtils';
import { LoadingSkeleton } from '../ui/LoadingSkeleton';

interface ChartRendererProps {
  chartType: string;
  data: unknown;
  locale: string;
}

interface ChartData {
  labels?: string[];
  series?: { name: string; values: unknown[] }[];
  columns?: { name: string; type: string }[];
  rows?: Record<string, unknown>[];
  value?: unknown;
  formatted?: string;
}

// Chart color palette
const CHART_COLORS = [
  '#3b82f6', // primary-500
  '#10b981', // green-500
  '#f59e0b', // amber-500
  '#ef4444', // red-500
  '#8b5cf6', // violet-500
  '#06b6d4', // cyan-500
];

function getChartColor(index: number): string {
  return CHART_COLORS[index % CHART_COLORS.length];
}

function formatCellValue(value: unknown, type: string, appLocale?: string): string {
  if (value === null || value === undefined) return '-';
  const locale = appLocale ?? 'en';

  switch (type) {
    case 'number':
      if (typeof value === 'number') {
        return formatNumber(value, locale);
      }
      return String(value);
    case 'date':
      if (value instanceof Date) {
        return formatDate(value, locale);
      }
      return String(value);
    default:
      return String(value);
  }
}

// Chart option creators (moved outside component)
function createLineChartOption(data: ChartData, isRTL: boolean): echarts.EChartsOption {
  return {
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: data.series?.map((s) => s.name) || [],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: data.labels || [],
      inverse: isRTL,
    },
    yAxis: {
      type: 'value',
      inverse: isRTL,
    },
    series:
      data.series?.map((s, idx) => ({
        name: s.name,
        type: 'line',
        data: s.values as number[],
        smooth: true,
        itemStyle: {
          color: getChartColor(idx),
        },
      })) || [],
  };
}

function createBarChartOption(data: ChartData, isRTL: boolean): echarts.EChartsOption {
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
    legend: {
      data: data.series?.map((s) => s.name) || [],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: data.labels || [],
      inverse: isRTL,
    },
    yAxis: {
      type: 'value',
      inverse: isRTL,
    },
    series:
      data.series?.map((s, idx) => ({
        name: s.name,
        type: 'bar',
        data: s.values as number[],
        itemStyle: {
          color: getChartColor(idx),
        },
      })) || [],
  };
}

function createPieChartOption(data: ChartData, isRTL: boolean): echarts.EChartsOption {
  const pieData =
    data.labels?.map((label, idx) => ({
      name: label,
      value: (data.series?.[0]?.values?.[idx] as number) || 0,
    })) || [];

  return {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      left: isRTL ? 'right' : 'left',
    },
    series: [
      {
        name: data.series?.[0]?.name || 'Distribution',
        type: 'pie',
        radius: ['40%', '70%'],
        center: [isRTL ? '35%' : '65%', '50%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold',
          },
        },
        labelLine: {
          show: false,
        },
        data: pieData,
      },
    ],
    color: CHART_COLORS,
  };
}

// KPI Card component
const KPICard: React.FC<{ chartData: ChartData }> = ({ chartData }) => {
  const { t } = useTranslation('dashboard');
  return (
    <div className="text-center py-6 p-4">
      <div className="text-4xl font-bold text-primary-600 dark:text-primary-400">
        {chartData.formatted || String(chartData.value || '')}
      </div>
      <div className="text-sm text-gray-500 dark:text-gray-400 mt-2">
        {t('kpi.totalValue')}
      </div>
    </div>
  );
};

// Data Table component
const DataTable: React.FC<{ chartData: ChartData; locale: string }> = ({ chartData, locale }) => {
  if (!chartData.columns || !chartData.rows) return null;

  return (
    <div className="overflow-x-auto p-4">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead>
          <tr>
            {chartData.columns.map((col, idx) => (
              <th
                key={idx}
                className="px-4 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
              >
                {col.name}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
          {chartData.rows.map((row, rowIdx) => (
            <tr key={rowIdx}>
              {chartData.columns?.map((col, colIdx) => (
                <td
                  key={colIdx}
                  className="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 whitespace-nowrap"
                >
                  {formatCellValue(row[col.name], col.type, locale)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

// Error state component
const ChartError: React.FC<{ message: string }> = ({ message }) => {
  const { t } = useTranslation('dashboard');
  return (
    <div className="liquid-glass-badge bg-red-50! dark:bg-red-900/20! border-red-200! dark:border-red-800! text-red-600 dark:text-red-400 rounded-xl p-4 text-sm" role="alert">
      <div className="flex items-center gap-2">
        <svg className="w-5 h-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{message || t('chart.error', 'Failed to load chart')}</span>
      </div>
    </div>
  );
};

export const ChartRenderer: React.FC<ChartRendererProps> = ({
  chartType,
  data,
  locale,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const chartData = useMemo(() => data as ChartData, [data]);
  const isRTL = locale === 'ar';

  // Handle ECharts rendering
  useEffect(() => {
    // Skip for non-chart types
    if (chartType === 'kpiCard' || chartType === 'dataTable') {
      setIsLoading(false);
      return;
    }

    if (!chartRef.current || !chartData) {
      setIsLoading(false);
      return;
    }

    // Initialize chart
    if (!chartInstance.current) {
      try {
        chartInstance.current = echarts.init(chartRef.current);
      } catch (err) {
        setError('Failed to initialize chart');
        setIsLoading(false);
        return;
      }
    }

    let option: echarts.EChartsOption;

    switch (chartType) {
      case 'lineChart':
        option = createLineChartOption(chartData, isRTL);
        break;
      case 'pieChart':
        option = createPieChartOption(chartData, isRTL);
        break;
      case 'barChart':
      default:
        option = createBarChartOption(chartData, isRTL);
    }

    try {
      chartInstance.current.setOption(option);
      setIsLoading(false);
      setError(null);
    } catch (err) {
      setError('Failed to render chart');
      setIsLoading(false);
    }

    // Handle resize
    const handleResize = () => {
      chartInstance.current?.resize();
    };
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [chartType, chartData, isRTL]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
    };
  }, []);

  // Render KPI Card
  if (chartType === 'kpiCard') {
    return <KPICard chartData={chartData} />;
  }

  // Render Data Table
  if (chartType === 'dataTable') {
    return <DataTable chartData={chartData} locale={locale} />;
  }

  // Render loading state
  if (isLoading) {
    return (
      <div className="p-4">
        <LoadingSkeleton variant="chart" height="300px" />
      </div>
    );
  }

  // Render error state
  if (error) {
    return <ChartError message={error} />;
  }

  // Render ECharts container
  return (
    <div
      ref={chartRef}
      style={{ width: '100%', height: '300px' }}
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
      role="img"
      aria-label={`${chartType} visualization`}
    />
  );
};

export default ChartRenderer;
