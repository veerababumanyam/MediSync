---
name: echarts-visualization
description: This skill should be used when the user asks to "create charts", "build visualizations", "implement ECharts", "add data visualizations", "dashboard charts", "graph components", or mentions Apache ECharts chart types (line, bar, pie, scatter, heatmap).
---

# Apache ECharts Visualization for MediSync

Apache ECharts 5.6 provides rich, interactive visualizations for MediSync dashboards, supporting real-time streaming, responsive design, and RTL layouts.

★ Insight ─────────────────────────────────────
MediSync visualization architecture:
1. **React-ECharts** - Declarative wrapper for ECharts
2. **Streaming charts** - Real-time data via WebSocket
3. **Responsive** - Auto-resize with container
4. **RTL support** - Chart mirroring for Arabic
5. **Confidence bands** - Visualize AI confidence scores
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Package** | echarts (5.6.x) + echarts-for-react (3.x) |
| **Themes** | Light/dark with MediSync brand colors |
| **Responsive** | Auto-resize enabled |
| **Export** | PNG, SVG, PDF via echarts extension |
| **i18n** | Locale-aware number/date formatting |

## Basic Setup

### Install

```bash
npm install echarts echarts-for-react
```

### Simple Chart Component

```typescript
import ReactECharts from 'echarts-for-react';
import { useTheme } from '../hooks/useTheme';

interface ChartProps {
  data: DataPoint[];
  title?: string;
}

export function LineChart({ data, title }: ChartProps) {
  const { isDark } = useTheme();

  const option = {
    title: {
      text: title,
      textStyle: { color: isDark ? '#fff' : '#333' },
    },
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      data: data.map(d => d.label),
    },
    yAxis: {
      type: 'value',
    },
    series: [{
      data: data.map(d => d.value),
      type: 'line',
      smooth: true,
    }],
  };

  return (
    <ReactECharts
      option={option}
      style={{ height: '400px', width: '100%' }}
      opts={{ renderer: 'canvas' }}
    />
  );
}
```

## MediSync Chart Types

### Revenue Trend Chart

```typescript
export function RevenueTrendChart({ data }: RevenueChartProps) {
  const { t, i18n } = useTranslation();
  const { isRTL } = useRTL();

  const option = {
    title: {
      text: t('dashboard.revenueTrend'),
      left: isRTL ? 'right' : 'left',
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const point = params[0];
        return `${point.name}<br/>
                ${t('common.revenue')}: ${formatCurrency(point.value, i18n.language)}`;
      },
    },
    grid: {
      left: isRTL ? '10%' : '3%',
      right: isRTL ? '3%' : '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.map(d => d.month),
      inverse: isRTL,
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => formatCompactCurrency(value, i18n.language),
      },
    },
    series: [
      {
        name: t('dashboard.actual'),
        type: 'line',
        smooth: true,
        data: data.map(d => d.actual),
        areaStyle: { opacity: 0.3 },
        itemStyle: { color: '#3b82f6' },
      },
      {
        name: t('dashboard.forecast'),
        type: 'line',
        smooth: true,
        data: data.map(d => d.forecast),
        lineStyle: { type: 'dashed' },
        itemStyle: { color: '#94a3b8' },
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: '350px' }} />;
}
```

### Patient Distribution Pie Chart

```typescript
export function PatientDistributionChart({ data }: PieChartProps) {
  const { t, i18n } = useTranslation();
  const { isRTL } = useRTL();

  const option = {
    title: {
      text: t('dashboard.patientDistribution'),
      left: 'center',
    },
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)',
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      rtl: isRTL,
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#fff',
          borderWidth: 2,
        },
        label: {
          show: true,
          formatter: '{b}: {d}%',
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold',
          },
        },
        data: data.map(d => ({
          value: d.count,
          name: t(`patientTypes.${d.type}`),
        })),
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: '300px' }} />;
}
```

### Appointment Heatmap

```typescript
export function AppointmentHeatmap({ data }: HeatmapProps) {
  const hours = Array.from({ length: 24 }, (_, i) => `${i}:00`);
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

  const option = {
    title: {
      text: 'Appointment Density',
    },
    tooltip: {
      position: 'top',
      formatter: (params: any) =>
        `${days[params.data[1]]} ${hours[params.data[0]]}: ${params.data[2]} appointments`,
    },
    grid: {
      left: '10%',
      right: '10%',
      top: '15%',
      bottom: '15%',
    },
    xAxis: {
      type: 'category',
      data: hours,
      splitArea: { show: true },
    },
    yAxis: {
      type: 'category',
      data: days,
      splitArea: { show: true },
    },
    visualMap: {
      min: 0,
      max: Math.max(...data.map(d => d[2])),
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '0%',
      inRange: {
        color: ['#f0f9ff', '#0ea5e9', '#0369a1'],
      },
    },
    series: [
      {
        type: 'heatmap',
        data: data,
        label: { show: false },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: '400px' }} />;
}
```

### Confidence Score Gauge (AI Feature)

```typescript
export function ConfidenceGauge({ confidence, label }: GaugeProps) {
  const option = {
    series: [
      {
        type: 'gauge',
        startAngle: 180,
        endAngle: 0,
        min: 0,
        max: 100,
        splitNumber: 4,
        radius: '100%',
        center: ['50%', '70%'],
        axisLine: {
          lineStyle: {
            width: 20,
            color: [
              [0.5, '#ef4444'],  // Low confidence: red
              [0.75, '#f59e0b'], // Medium: amber
              [1, '#22c55e'],    // High: green
            ],
          },
        },
        pointer: {
          icon: 'path://M12,2C13.1,2 14,2.9 14,4V8C14,9.1 13.1,10 12,10C10.9,10 10,9.1 10,8V4C10,2.9 10.9,2 12,2Z',
          length: '60%',
          width: 8,
          offsetCenter: [0, '-10%'],
        },
        axisTick: { show: false },
        splitLine: { show: false },
        axisLabel: { show: false },
        title: {
          offsetCenter: [0, '20%'],
          fontSize: 14,
        },
        detail: {
          fontSize: 24,
          offsetCenter: [0, '0%'],
          formatter: '{value}%',
        },
        data: [
          {
            value: confidence,
            name: label,
          },
        ],
      },
    ],
  };

  return <ReactECharts option={option} style={{ height: '200px' }} />;
}
```

## Real-Time Streaming Chart

```typescript
export function StreamingChart({ stream }: StreamingChartProps) {
  const chartRef = useRef<any>(null);
  const [data, setData] = useState<DataPoint[]>([]);

  useEffect(() => {
    const subscription = stream.subscribe((point) => {
      setData(prev => {
        const newData = [...prev, point];
        // Keep last 100 points
        return newData.slice(-100);
      });

      // Update chart without full re-render
      const chart = chartRef.current?.getEchartsInstance();
      if (chart) {
        chart.setOption({
          series: [{ data: data.map(d => d.value) }],
          xAxis: { data: data.map(d => d.label) },
        });
      }
    });

    return () => subscription.unsubscribe();
  }, [stream]);

  return (
    <ReactECharts
      ref={chartRef}
      option={{
        animation: false, // Disable animation for performance
        xAxis: { type: 'category', data: data.map(d => d.label) },
        yAxis: { type: 'value' },
        series: [{ type: 'line', data: data.map(d => d.value) }],
      }}
      style={{ height: '300px' }}
    />
  );
}
```

## Responsive Charts

```typescript
export function ResponsiveChart({ option }: ResponsiveChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<any>(null);

  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      chartRef.current?.getEchartsInstance()?.resize();
    });

    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }

    return () => resizeObserver.disconnect();
  }, []);

  return (
    <div ref={containerRef} className="w-full h-full">
      <ReactECharts
        ref={chartRef}
        option={option}
        style={{ width: '100%', height: '100%' }}
        opts={{ renderer: 'canvas' }}
      />
    </div>
  );
}
```

## RTL Support

```typescript
function applyRTLOption(option: EChartsOption, isRTL: boolean): EChartsOption {
  if (!isRTL) return option;

  return {
    ...option,
    // Mirror grid margins
    grid: {
      left: option.grid?.right ?? '3%',
      right: option.grid?.left ?? '4%',
      bottom: option.grid?.bottom ?? '3%',
      top: option.grid?.top ?? '10%',
      containLabel: option.grid?.containLabel ?? true,
    },
    // RTL legend
    legend: {
      ...option.legend,
      rtl: true,
      align: 'right',
    },
    // Invert axis if needed
    xAxis: option.xAxis?.inverse ? { ...option.xAxis, inverse: true } : option.xAxis,
  };
}
```

## Export Functionality

```typescript
export function useChartExport() {
  const exportToPNG = useCallback((chart: ECharts, filename: string) => {
    const url = chart.getDataURL({
      type: 'png',
      pixelRatio: 2,
      backgroundColor: '#fff',
    });

    const link = document.createElement('a');
    link.download = `${filename}.png`;
    link.href = url;
    link.click();
  }, []);

  const exportToSVG = useCallback((chart: ECharts, filename: string) => {
    // Requires svg renderer
    const url = chart.getDataURL({ type: 'svg' });
    const link = document.createElement('a');
    link.download = `${filename}.svg`;
    link.href = url;
    link.click();
  }, []);

  return { exportToPNG, exportToSVG };
}
```

## Additional Resources

### Reference Files
- **`references/chart-types.md`** - Complete chart type reference
- **`references/theming.md`** - Custom themes and styling

### Example Files
- **`examples/dashboard-charts.tsx`** - Complete dashboard chart components
- **`examples/ai-confidence-charts.tsx`** - Confidence visualization patterns
