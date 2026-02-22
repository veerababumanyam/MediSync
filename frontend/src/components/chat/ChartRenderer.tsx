import React, { useEffect, useRef, useMemo } from 'react';
import * as echarts from 'echarts';
import { useTranslation } from 'react-i18next';

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

function formatCellValue(value: unknown, type: string): string {
  if (value === null || value === undefined) return '-';

  switch (type) {
    case 'number':
      if (typeof value === 'number') {
        return value.toLocaleString();
      }
      return String(value);
    case 'date':
      if (value instanceof Date) {
        return value.toLocaleDateString();
      }
      return String(value);
    default:
      return String(value);
  }
}

// Chart option creators (moved outside component)
function createLineChartOption(data: ChartData): echarts.EChartsOption {
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
    },
    yAxis: {
      type: 'value',
    },
    series:
      data.series?.map((s, idx) => ({
        name: s.name,
        type: 'line' as const,
        data: s.values as number[],
        smooth: true,
        itemStyle: {
          color: getChartColor(idx),
        },
      })) || [],
  };
}

function createBarChartOption(data: ChartData): echarts.EChartsOption {
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
    },
    yAxis: {
      type: 'value',
    },
    series:
      data.series?.map((s, idx) => ({
        name: s.name,
        type: 'bar' as const,
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
        type: 'pie' as const,
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
    <div className="text-center py-6">
      <div className="text-4xl font-bold text-primary-600 dark:text-primary-400">
        {chartData.formatted || String(chartData.value ?? '')}
      </div>
      <div className="text-sm text-gray-500 dark:text-gray-400 mt-2">
        {t('kpi.totalValue')}
      </div>
    </div>
  );
};

// Data Table component
const DataTable: React.FC<{ chartData: ChartData }> = ({ chartData }) => {
  if (!chartData.columns || !chartData.rows) return null;

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead>
          <tr>
            {chartData.columns.map((col, idx) => (
              <th
                key={idx}
                className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
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
                  {formatCellValue(row[col.name], col.type)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
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

  const chartData = useMemo(() => data as ChartData, [data]);
  const isRTL = locale === 'ar';

  // Handle ECharts rendering
  useEffect(() => {
    // Skip for non-chart types
    if (chartType === 'kpiCard' || chartType === 'dataTable') {
      return;
    }

    if (!chartRef.current || !chartData) return;

    // Initialize chart
    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    let option: echarts.EChartsOption = {};

    switch (chartType) {
      case 'lineChart':
        option = createLineChartOption(chartData);
        break;
      case 'pieChart':
        option = createPieChartOption(chartData, isRTL);
        break;
      case 'barChart':
      default:
        option = createBarChartOption(chartData);
    }

    chartInstance.current.setOption(option);

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
    return <DataTable chartData={chartData} />;
  }

  // Render ECharts container
  return (
    <div
      ref={chartRef}
      style={{ width: '100%', height: '300px' }}
      dir={locale === 'ar' ? 'rtl' : 'ltr'}
    />
  );
};

export default ChartRenderer;
