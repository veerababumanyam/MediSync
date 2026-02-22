/**
 * MediSync CopilotKit Tools
 *
 * Defines tools for CopilotKit generative UI with render functions.
 * These tools allow AI agents to render React components dynamically.
 *
 * @module components/copilot/MediSyncTools
 */
import React from 'react'

/**
 * Tool parameter types
 */
export interface QueryBIParams {
  query: string
}

export interface SyncTallyParams {
  entryIds?: string[]
}

export interface PinChartParams {
  query: string
  title: string
  chartType?: 'bar' | 'line' | 'pie' | 'table' | 'kpi'
}

export interface NavigateParams {
  route: 'home' | 'chat' | 'dashboard'
}

export interface CreateAlertParams {
  metric: string
  threshold: number
  operator: 'gt' | 'gte' | 'lt' | 'lte' | 'eq'
}

export interface CreateReportParams {
  query: string
  schedule: 'daily' | 'weekly' | 'monthly'
}

export interface ExportParams {
  format: 'pdf' | 'xlsx' | 'csv'
}

/**
 * Result types
 */
export interface ToolResult {
  success: boolean
  message: string
  data?: unknown
}

/**
 * Query result for BI queries
 */
export interface QueryResult extends ToolResult {
  data?: {
    columns: string[]
    rows: Record<string, unknown>[]
    chartType: string
    sql: string
    confidence: number
  }
}

/**
 * Component for rendering query results
 */
export const QueryResultComponent: React.FC<{
  result: QueryResult
  onExplore?: (query: string) => void
}> = ({ result, onExplore }) => {
  if (!result.success || !result.data) {
    return (
      <div className="p-4 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg">
        {result.message}
      </div>
    )
  }

  const { columns, rows, chartType, sql, confidence } = result.data

  return (
    <div
      className="bg-white dark:bg-slate-800 rounded-xl p-4 shadow-sm border border-slate-200 dark:border-slate-700"
      data-tool-name="query-result"
      data-tool-description="Displays the results of a BI query with chart visualization"
    >
      {/* Confidence indicator */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs text-slate-500 dark:text-slate-400">
          Chart: {chartType.toUpperCase()}
        </span>
        <div className="flex items-center gap-2">
          <span className="text-xs text-slate-500 dark:text-slate-400">
            Confidence:
          </span>
          <div className="w-20 h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
            <div
              className={`h-full ${
                confidence >= 90
                  ? 'bg-emerald-500'
                  : confidence >= 70
                    ? 'bg-amber-500'
                    : 'bg-red-500'
              }`}
              style={{ width: `${confidence}%` }}
            />
          </div>
          <span className="text-xs font-medium text-slate-700 dark:text-slate-300">
            {confidence}%
          </span>
        </div>
      </div>

      {/* Data table preview */}
      <div className="overflow-x-auto mb-3">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-200 dark:border-slate-700">
              {columns.map((col) => (
                <th
                  key={col}
                  className="px-3 py-2 text-left font-medium text-slate-700 dark:text-slate-300"
                >
                  {col}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.slice(0, 5).map((row, i) => (
              <tr
                key={i}
                className="border-b border-slate-100 dark:border-slate-800"
              >
                {columns.map((col) => (
                  <td
                    key={col}
                    className="px-3 py-2 text-slate-600 dark:text-slate-400"
                  >
                    {String(row[col])}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
        {rows.length > 5 && (
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-2">
            Showing 5 of {rows.length} rows
          </p>
        )}
      </div>

      {/* SQL preview (collapsible) */}
      <details className="mb-3">
        <summary className="text-xs text-slate-500 dark:text-slate-400 cursor-pointer hover:text-slate-700 dark:hover:text-slate-300">
          View generated SQL
        </summary>
        <pre className="mt-2 p-2 bg-slate-100 dark:bg-slate-900 rounded text-xs overflow-x-auto">
          {sql}
        </pre>
      </details>

      {/* Actions */}
      <div className="flex items-center gap-2">
        <button
          onClick={() => onExplore?.(sql)}
          className="px-3 py-1.5 text-xs bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Explore in Chat
        </button>
        <button
          className="px-3 py-1.5 text-xs bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-600 transition-colors"
        >
          Pin to Dashboard
        </button>
      </div>
    </div>
  )
}

/**
 * Component for rendering sync status
 */
export const SyncStatusComponent: React.FC<{
  result: ToolResult
  status?: 'pending' | 'syncing' | 'completed' | 'failed'
}> = ({ result, status = 'completed' }) => {
  const statusConfig = {
    pending: { color: 'bg-amber-500', label: 'Pending' },
    syncing: { color: 'bg-blue-500', label: 'Syncing...' },
    completed: { color: 'bg-emerald-500', label: 'Completed' },
    failed: { color: 'bg-red-500', label: 'Failed' },
  }

  const config = statusConfig[status]

  return (
    <div className="p-4 bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700">
      <div className="flex items-center gap-3">
        <div className={`w-3 h-3 rounded-full ${config.color}`} />
        <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
          Tally Sync {config.label}
        </span>
      </div>
      <p className="mt-2 text-sm text-slate-600 dark:text-slate-400">
        {result.message}
      </p>
    </div>
  )
}

/**
 * Component for rendering navigation actions
 */
export const NavigationComponent: React.FC<{
  result: ToolResult
  currentRoute: string
}> = ({ result, currentRoute }) => {
  return (
    <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200 dark:border-blue-800">
      <div className="flex items-center gap-2 mb-2">
        <svg
          className="w-5 h-5 text-blue-600 dark:text-blue-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M13 7l5 5m0 0l-5 5m5-5H6"
          />
        </svg>
        <span className="text-sm font-medium text-blue-700 dark:text-blue-300">
          Navigating to {currentRoute}
        </span>
      </div>
      <p className="text-sm text-blue-600 dark:text-blue-400">
        {result.message}
      </p>
    </div>
  )
}

/**
 * Component for rendering alert creation
 */
export const AlertCreatedComponent: React.FC<{
  result: ToolResult
  params: CreateAlertParams
}> = ({ result, params }) => {
  const operatorSymbols: Record<string, string> = {
    gt: '>',
    gte: '>=',
    lt: '<',
    lte: '<=',
    eq: '=',
  }

  return (
    <div className="p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
      <div className="flex items-center gap-2 mb-2">
        <svg
          className="w-5 h-5 text-amber-600 dark:text-amber-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
          />
        </svg>
        <span className="text-sm font-medium text-amber-700 dark:text-amber-300">
          Alert Created
        </span>
      </div>
      <p className="text-sm text-amber-600 dark:text-amber-400">
        {params.metric} {operatorSymbols[params.operator]} {params.threshold}
      </p>
      <p className="mt-2 text-xs text-amber-500 dark:text-amber-500">
        {result.message}
      </p>
    </div>
  )
}

/**
 * Component for rendering report creation
 */
export const ReportCreatedComponent: React.FC<{
  result: ToolResult
  params: CreateReportParams
}> = ({ result, params }) => {
  return (
    <div className="p-4 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200 dark:border-emerald-800">
      <div className="flex items-center gap-2 mb-2">
        <svg
          className="w-5 h-5 text-emerald-600 dark:text-emerald-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
        <span className="text-sm font-medium text-emerald-700 dark:text-emerald-300">
          Scheduled Report Created
        </span>
      </div>
      <p className="text-sm text-emerald-600 dark:text-emerald-400">
        Query: "{params.query}"
      </p>
      <p className="text-xs text-emerald-500 dark:text-emerald-500 mt-1">
        Schedule: {params.schedule}
      </p>
      <p className="mt-2 text-xs text-emerald-500 dark:text-emerald-500">
        {result.message}
      </p>
    </div>
  )
}

/**
 * Component for rendering export status
 */
export const ExportStatusComponent: React.FC<{
  result: ToolResult
  format: string
}> = ({ result, format }) => {
  const formatIcons: Record<string, string> = {
    pdf: '📄',
    xlsx: '📊',
    csv: '📋',
  }

  return (
    <div className="p-4 bg-slate-50 dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700">
      <div className="flex items-center gap-2 mb-2">
        <span className="text-xl">{formatIcons[format] || '📄'}</span>
        <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
          Export as {format.toUpperCase()}
        </span>
      </div>
      <p className="text-sm text-slate-600 dark:text-slate-400">
        {result.message}
      </p>
      {result.success && (
        <button className="mt-3 px-3 py-1.5 text-xs bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
          Download
        </button>
      )}
    </div>
  )
}

/**
 * Tool definitions for CopilotKit
 * These are passed to CopilotKit provider
 */
export const medisyncTools = [
  {
    name: 'queryBI',
    description:
      'Execute a natural language query against MediSync BI data. Returns charts, tables, and insights.',
    parameters: {
      type: 'object' as const,
      properties: {
        query: {
          type: 'string',
          description: 'The natural language query to execute',
        },
      },
      required: ['query'],
    },
    handler: async (params: QueryBIParams): Promise<QueryResult> => {
      // This would call the actual API
      console.log('queryBI called with:', params)
      return {
        success: true,
        message: `Query "${params.query}" executed successfully`,
        data: {
          columns: ['Metric', 'Value'],
          rows: [{ Metric: 'Revenue', Value: '$125,000' }],
          chartType: 'bar',
          sql: 'SELECT metric, value FROM analytics WHERE ...',
          confidence: 95,
        },
      }
    },
    render: (result: QueryResult) => <QueryResultComponent result={result} />,
  },
  {
    name: 'syncTally',
    description: 'Synchronize approved entries to Tally ERP',
    parameters: {
      type: 'object' as const,
      properties: {
        entryIds: {
          type: 'array',
          items: { type: 'string' },
          description: 'Optional specific entry IDs to sync',
        },
      },
    },
    handler: async (params: SyncTallyParams): Promise<ToolResult> => {
      console.log('syncTally called with:', params)
      return {
        success: true,
        message: 'Synchronization completed successfully',
      }
    },
    render: (result: ToolResult) => (
      <SyncStatusComponent result={result} status="completed" />
    ),
  },
  {
    name: 'navigate',
    description: 'Navigate to a different page in MediSync',
    parameters: {
      type: 'object' as const,
      properties: {
        route: {
          type: 'string',
          enum: ['home', 'chat', 'dashboard'],
          description: 'The route to navigate to',
        },
      },
      required: ['route'],
    },
    handler: async (params: NavigateParams): Promise<ToolResult> => {
      console.log('navigate called with:', params)
      const path = params.route === 'home' ? '/' : `/${params.route}`
      window.location.href = path
      return {
        success: true,
        message: `Navigating to ${params.route}`,
      }
    },
    render: (result: ToolResult, params: NavigateParams) => (
      <NavigationComponent result={result} currentRoute={params.route} />
    ),
  },
  {
    name: 'createAlert',
    description: 'Create an alert for a specific metric threshold',
    parameters: {
      type: 'object' as const,
      properties: {
        metric: {
          type: 'string',
          description: 'The metric to monitor (e.g., revenue, patient_count)',
        },
        threshold: {
          type: 'number',
          description: 'The threshold value',
        },
        operator: {
          type: 'string',
          enum: ['gt', 'gte', 'lt', 'lte', 'eq'],
          description: 'The comparison operator',
        },
      },
      required: ['metric', 'threshold'],
    },
    handler: async (params: CreateAlertParams): Promise<ToolResult> => {
      console.log('createAlert called with:', params)
      return {
        success: true,
        message: `Alert created for ${params.metric}`,
      }
    },
    render: (result: ToolResult, params: CreateAlertParams) => (
      <AlertCreatedComponent result={result} params={params} />
    ),
  },
  {
    name: 'createReport',
    description: 'Create a scheduled report from a natural language query',
    parameters: {
      type: 'object' as const,
      properties: {
        query: {
          type: 'string',
          description: 'The natural language query for the report',
        },
        schedule: {
          type: 'string',
          enum: ['daily', 'weekly', 'monthly'],
          description: 'The schedule type',
        },
      },
      required: ['query'],
    },
    handler: async (params: CreateReportParams): Promise<ToolResult> => {
      console.log('createReport called with:', params)
      return {
        success: true,
        message: 'Report scheduled successfully',
      }
    },
    render: (result: ToolResult, params: CreateReportParams) => (
      <ReportCreatedComponent result={result} params={params} />
    ),
  },
  {
    name: 'exportView',
    description: 'Export the current view or report to a file',
    parameters: {
      type: 'object' as const,
      properties: {
        format: {
          type: 'string',
          enum: ['pdf', 'xlsx', 'csv'],
          description: 'The export format',
        },
      },
      required: ['format'],
    },
    handler: async (params: ExportParams): Promise<ToolResult> => {
      console.log('exportView called with:', params)
      return {
        success: true,
        message: `Exported as ${params.format.toUpperCase()}`,
      }
    },
    render: (result: ToolResult, params: ExportParams) => (
      <ExportStatusComponent result={result} format={params.format} />
    ),
  },
]

export default medisyncTools
