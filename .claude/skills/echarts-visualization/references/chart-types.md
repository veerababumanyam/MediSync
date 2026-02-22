# Apache ECharts Chart Types Reference

## Line Charts

### Basic Line

```typescript
const option = {
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
  },
  yAxis: {
    type: 'value',
  },
  series: [{
    data: [150, 230, 224, 218, 135, 147, 260],
    type: 'line',
  }],
};
```

### Multi-Line with Legend

```typescript
const option = {
  legend: {
    data: ['Actual', 'Forecast', 'Target'],
  },
  xAxis: {
    type: 'category',
    data: months,
  },
  yAxis: {
    type: 'value',
  },
  series: [
    {
      name: 'Actual',
      type: 'line',
      data: actualData,
      smooth: true,
      areaStyle: { opacity: 0.3 },
    },
    {
      name: 'Forecast',
      type: 'line',
      data: forecastData,
      lineStyle: { type: 'dashed' },
    },
    {
      name: 'Target',
      type: 'line',
      data: targetData,
      lineStyle: { width: 2, type: 'solid' },
    },
  ],
};
```

### Line with Confidence Band

```typescript
const option = {
  xAxis: { type: 'category', data: dates },
  yAxis: { type: 'value' },
  series: [
    {
      name: 'Upper Bound',
      type: 'line',
      data: upperBound,
      lineStyle: { opacity: 0 },
      areaStyle: { opacity: 0 },
      stack: 'confidence',
      symbol: 'none',
    },
    {
      name: 'Lower Bound',
      type: 'line',
      data: lowerBound.map((v, i) => upperBound[i] - v),
      lineStyle: { opacity: 0 },
      areaStyle: { color: 'rgba(59, 130, 246, 0.2)' },
      stack: 'confidence',
      symbol: 'none',
    },
    {
      name: 'Actual',
      type: 'line',
      data: actual,
      itemStyle: { color: '#3b82f6' },
    },
  ],
};
```

## Bar Charts

### Vertical Bar

```typescript
const option = {
  xAxis: {
    type: 'category',
    data: ['Q1', 'Q2', 'Q3', 'Q4'],
  },
  yAxis: {
    type: 'value',
  },
  series: [{
    data: [120, 200, 150, 80],
    type: 'bar',
    itemStyle: {
      color: '#3b82f6',
      borderRadius: [4, 4, 0, 0],
    },
  }],
};
```

### Horizontal Bar

```typescript
const option = {
  xAxis: {
    type: 'value',
  },
  yAxis: {
    type: 'category',
    data: ['Category A', 'Category B', 'Category C', 'Category D'],
  },
  series: [{
    data: [120, 200, 150, 80],
    type: 'bar',
    itemStyle: {
      borderRadius: [0, 4, 4, 0],
    },
  }],
};
```

### Grouped Bar

```typescript
const option = {
  legend: {
    data: ['2024', '2025'],
  },
  xAxis: {
    type: 'category',
    data: ['Jan', 'Feb', 'Mar', 'Apr'],
  },
  yAxis: {
    type: 'value',
  },
  series: [
    {
      name: '2024',
      type: 'bar',
      data: [320, 302, 301, 334],
      itemStyle: { color: '#3b82f6' },
    },
    {
      name: '2025',
      type: 'bar',
      data: [120, 132, 101, 134],
      itemStyle: { color: '#10b981' },
    },
  ],
};
```

### Stacked Bar

```typescript
const option = {
  legend: {
    data: ['Direct', 'Referral', 'Organic'],
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed'],
  },
  yAxis: {
    type: 'value',
  },
  series: [
    {
      name: 'Direct',
      type: 'bar',
      stack: 'total',
      data: [320, 302, 301],
    },
    {
      name: 'Referral',
      type: 'bar',
      stack: 'total',
      data: [120, 132, 101],
    },
    {
      name: 'Organic',
      type: 'bar',
      stack: 'total',
      data: [220, 182, 191],
    },
  ],
};
```

## Pie Charts

### Basic Pie

```typescript
const option = {
  series: [{
    type: 'pie',
    radius: '50%',
    data: [
      { value: 1048, name: 'Search Engine' },
      { value: 735, name: 'Direct' },
      { value: 580, name: 'Email' },
      { value: 484, name: 'Union Ads' },
      { value: 300, name: 'Video Ads' },
    ],
  }],
};
```

### Donut Chart

```typescript
const option = {
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: {
      borderRadius: 10,
      borderColor: '#fff',
      borderWidth: 2,
    },
    label: {
      show: false,
      position: 'center',
    },
    emphasis: {
      label: {
        show: true,
        fontSize: 20,
        fontWeight: 'bold',
      },
    },
    labelLine: {
      show: false,
    },
    data: pieData,
  }],
};
```

### Rose Chart (Nightingale)

```typescript
const option = {
  series: [{
    type: 'pie',
    radius: ['20%', '70%'],
    roseType: 'area',
    itemStyle: {
      borderRadius: 5,
    },
    data: [
      { value: 30, name: 'rose 1' },
      { value: 28, name: 'rose 2' },
      { value: 26, name: 'rose 3' },
    ],
  }],
};
```

## Scatter Charts

### Basic Scatter

```typescript
const option = {
  xAxis: {},
  yAxis: {},
  series: [{
    symbolSize: 20,
    data: [
      [10.0, 8.04],
      [8.07, 6.95],
      [13.0, 7.58],
      [9.05, 8.81],
    ],
    type: 'scatter',
  }],
};
```

### Bubble Chart

```typescript
const option = {
  xAxis: {},
  yAxis: {},
  series: [{
    data: [
      [10.0, 8.04, 10],  // x, y, size
      [8.07, 6.95, 20],
      [13.0, 7.58, 30],
    ],
    type: 'scatter',
    symbolSize: (data: number[]) => data[2],
  }],
};
```

## Heatmap

```typescript
const option = {
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
    max: 10,
    calculable: true,
    orient: 'horizontal',
    left: 'center',
    inRange: {
      color: ['#f0f9ff', '#0ea5e9', '#0369a1'],
    },
  },
  series: [{
    type: 'heatmap',
    data: heatmapData,
    label: {
      show: true,
    },
  }],
};
```

## Gauge

```typescript
const option = {
  series: [{
    type: 'gauge',
    startAngle: 180,
    endAngle: 0,
    min: 0,
    max: 100,
    splitNumber: 4,
    axisLine: {
      lineStyle: {
        width: 20,
        color: [
          [0.5, '#ef4444'],
          [0.75, '#f59e0b'],
          [1, '#22c55e'],
        ],
      },
    },
    pointer: {
      length: '60%',
      width: 8,
    },
    detail: {
      valueAnimation: true,
      formatter: '{value}%',
    },
    data: [{ value: 85, name: 'Confidence' }],
  }],
};
```

## Radar Chart

```typescript
const option = {
  radar: {
    indicator: [
      { name: 'Sales', max: 6500 },
      { name: 'Administration', max: 16000 },
      { name: 'Information Technology', max: 30000 },
      { name: 'Customer Support', max: 38000 },
      { name: 'Development', max: 52000 },
    ],
  },
  series: [{
    type: 'radar',
    data: [
      {
        value: [4200, 3000, 20000, 35000, 50000],
        name: 'Budget',
      },
    ],
  }],
};
```

## Treemap

```typescript
const option = {
  series: [{
    type: 'treemap',
    data: [
      {
        name: 'Node A',
        value: 10,
        children: [
          { name: 'Node Aa', value: 4 },
          { name: 'Node Ab', value: 6 },
        ],
      },
      {
        name: 'Node B',
        value: 20,
      },
    ],
  }],
};
```

## Sunburst

```typescript
const option = {
  series: [{
    type: 'sunburst',
    data: [
      {
        name: 'Grandpa',
        children: [
          {
            name: 'Uncle Leo',
            value: 15,
            children: [
              { name: 'Cousin Ben', value: 2 },
            ],
          },
        ],
      },
    ],
    radius: ['15%', '80%'],
    label: {
      rotate: 'radial',
    },
  }],
};
```

## Sankey

```typescript
const option = {
  series: [{
    type: 'sankey',
    layout: 'none',
    emphasis: {
      focus: 'adjacency',
    },
    data: [
      { name: 'a' },
      { name: 'b' },
      { name: 'c' },
    ],
    links: [
      { source: 'a', target: 'b', value: 5 },
      { source: 'a', target: 'c', value: 3 },
    ],
  }],
};
```

## Funnel

```typescript
const option = {
  series: [{
    type: 'funnel',
    left: '10%',
    top: 60,
    bottom: 60,
    width: '80%',
    min: 0,
    max: 100,
    minSize: '0%',
    maxSize: '100%',
    sort: 'descending',
    gap: 2,
    data: [
      { value: 100, name: 'Show' },
      { value: 80, name: 'Click' },
      { value: 60, name: 'Visit' },
      { value: 40, name: 'Inquiry' },
      { value: 20, name: 'Order' },
    ],
  }],
};
```
