/**
 * MediSync API Client Service
 *
 * Centralized HTTP client for all API requests with:
 * - Authentication header injection
 * - Error handling
 * - TypeScript types for responses
 *
 * @module services/api
 */

// API Configuration
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

/**
 * API Error class for structured error handling
 */
export class APIError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    message: string,
    public data?: unknown
  ) {
    super(message)
    this.name = 'APIError'
  }
}

/**
 * Request options type
 */
interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown
  params?: Record<string, string | number | boolean>
}

/**
 * Get auth token from storage
 */
function getAuthToken(): string | null {
  return localStorage.getItem('medisync-token')
}

/**
 * Build URL with query parameters
 */
function buildURL(path: string, params?: Record<string, string | number | boolean>): string {
  const url = new URL(`${API_BASE_URL}${path}`, window.location.origin)
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        url.searchParams.set(key, String(value))
      }
    })
  }
  return url.toString()
}

/**
 * Make an API request
 */
async function request<T>(
  method: string,
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const { body, params, headers: customHeaders, ...init } = options

  const url = buildURL(path, params)
  const token = getAuthToken()

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...customHeaders,
  }

  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    ...init,
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  // Handle non-JSON responses
  const contentType = response.headers.get('content-type')
  let data: unknown

  if (contentType?.includes('application/json')) {
    data = await response.json()
  } else if (contentType?.includes('text/')) {
    data = await response.text()
  } else {
    data = await response.blob()
  }

  if (!response.ok) {
    const message = typeof data === 'object' && data !== null && 'message' in data
      ? (data as { message: string }).message
      : response.statusText

    throw new APIError(response.status, response.statusText, message, data)
  }

  return data as T
}

/**
 * API client object with HTTP methods
 */
export const api = {
  get: <T>(path: string, options?: RequestOptions) =>
    request<T>('GET', path, options),

  post: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>('POST', path, { ...options, body }),

  patch: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>('PATCH', path, { ...options, body }),

  put: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>('PUT', path, { ...options, body }),

  delete: <T>(path: string, options?: RequestOptions) =>
    request<T>('DELETE', path, options),
}

// ============================================================================
// API Response Types
// ============================================================================

/**
 * User preferences response
 */
export interface UserPreferences {
  id: string
  userId: string
  locale: 'en' | 'ar'
  numeralSystem: 'western' | 'eastern_arabic'
  calendarSystem: 'gregorian' | 'hijri'
  reportLanguage: 'en' | 'ar'
  timezone: string
  createdAt: string
  updatedAt: string
}

/**
 * Chat request
 */
export interface ChatRequest {
  query: string
  sessionId?: string
  locale?: string
}

/**
 * Pinned chart response
 */
export interface PinnedChart {
  id: string
  userId: string
  title: string
  queryId: string | null
  naturalLanguageQuery: string
  sqlQuery: string
  chartSpec: Record<string, unknown>
  chartType: 'bar' | 'line' | 'pie' | 'table' | 'kpi'
  refreshInterval: number
  locale: 'en' | 'ar'
  position: { row: number; col: number; size: number }
  lastRefreshedAt: string | null
  isActive: boolean
  createdAt: string
  updatedAt: string
}

/**
 * Alert rule response
 */
export interface AlertRule {
  id: string
  userId: string
  name: string
  description: string | null
  metricId: string
  metricName: string
  operator: 'gt' | 'gte' | 'lt' | 'lte' | 'eq'
  threshold: number
  checkInterval: number
  channels: string[]
  locale: 'en' | 'ar'
  cooldownPeriod: number
  lastTriggeredAt: string | null
  lastValue: number | null
  isActive: boolean
  createdAt: string
  updatedAt: string
}

/**
 * Notification response
 */
export interface Notification {
  id: string
  alertRuleId: string
  userId: string
  type: 'in_app' | 'email'
  status: 'pending' | 'sent' | 'delivered' | 'failed'
  content: {
    title: string
    message: string
    actionUrl?: string
  }
  locale: 'en' | 'ar'
  metricValue: number
  threshold: number
  errorMessage: string | null
  sentAt: string | null
  deliveredAt: string | null
  readAt: string | null
  createdAt: string
}

/**
 * Scheduled report response
 */
export interface ScheduledReport {
  id: string
  userId: string
  name: string
  description: string | null
  queryId: string | null
  naturalLanguageQuery: string
  sqlQuery: string
  scheduleType: 'daily' | 'weekly' | 'monthly' | 'quarterly'
  scheduleTime: string
  scheduleDay: number | null
  recipients: Array<{ email: string; name: string }>
  format: 'pdf' | 'xlsx' | 'csv'
  locale: 'en' | 'ar'
  includeCharts: boolean
  lastRunAt: string | null
  nextRunAt: string | null
  isActive: boolean
  createdAt: string
  updatedAt: string
}

/**
 * Chat message response
 */
export interface ChatMessage {
  id: string
  sessionId: string
  role: 'user' | 'assistant' | 'system'
  content: string
  chartSpec?: {
    type: string
    chart: unknown
  }
  tableData?: Record<string, unknown>
  drilldownQuery?: string
  confidence?: number
  createdAt: string
}

/**
 * SSE Event types
 */
export type SSEEventType = 'thinking' | 'sql_preview' | 'result' | 'error' | 'clarification'

/**
 * SSE Event response
 */
export interface SSEEvent {
  type: SSEEventType
  message?: string
  sql?: string
  chartType?: string
  data?: unknown
  confidence?: number
  options?: string[]
}

/**
 * Stream chat result
 */
export interface StreamChatResult {
  message?: string
  chartSpec?: {
    type: string
    chart: unknown
  }
  confidence?: number
}

/**
 * Chat messages list response
 */
export interface ChatMessagesResponse {
  messages: ChatMessage[]
  total: number
  hasMore: boolean
}

/**
 * Chat sessions response
 */
export interface ChatSessionsResponse {
  sessions: Array<{
    id: string
    title?: string
    lastMessage?: string
    messageCount: number
    createdAt: string
    lastActivityAt: string
  }>
  total: number
}

// ============================================================================
// API Endpoints
// ============================================================================

/**
 * Preferences API
 */
export const preferencesApi = {
  get: () => api.get<UserPreferences>('/preferences'),
  update: (data: Partial<UserPreferences>) =>
    api.patch<UserPreferences>('/preferences', data),
}

/**
 * Chat API
 */
export const chatApi = {
  query: (data: ChatRequest) =>
    api.post<Response>('/chat/query', data),
  drilldown: (data: { queryId: string; filters: Record<string, unknown> }) =>
    api.post<Response>('/chat/drilldown', data),
  export: (format: string, data: { query: string; results: unknown[] }) =>
    api.post<Blob>(`/chat/export/${format}`, data),
}

/**
 * Dashboard API
 */
export const dashboardApi = {
  getCharts: () => api.get<PinnedChart[]>('/dashboard/charts'),
  pinChart: (chart: Partial<PinnedChart>) =>
    api.post<PinnedChart>('/dashboard/charts', chart),
  updateChart: (id: string, chart: Partial<PinnedChart>) =>
    api.patch<PinnedChart>(`/dashboard/charts/${id}`, chart),
  deleteChart: (id: string) =>
    api.delete(`/dashboard/charts/${id}`),
  refreshChart: (id: string) =>
    api.post<PinnedChart>(`/dashboard/charts/${id}/refresh`),
  reorderCharts: (positions: Array<{ id: string; position: { row: number; col: number; size: number } }>) =>
    api.post('/dashboard/charts/reorder', { positions }),
  getQuickActions: () =>
    api.get<Array<{ id: string; query: string; label: string }>>('/dashboard/quick-actions'),
}

/**
 * Alerts API
 */
export const alertsApi = {
  getRules: () => api.get<AlertRule[]>('/alerts/rules'),
  createRule: (rule: Partial<AlertRule>) =>
    api.post<AlertRule>('/alerts/rules', rule),
  updateRule: (id: string, rule: Partial<AlertRule>) =>
    api.patch<AlertRule>(`/alerts/rules/${id}`, rule),
  deleteRule: (id: string) =>
    api.delete(`/alerts/rules/${id}`),
  toggleRule: (id: string, isActive: boolean) =>
    api.post(`/alerts/rules/${id}/toggle`, { isActive }),
  testRule: (id: string) =>
    api.post(`/alerts/rules/${id}/test`),
  getMetrics: () =>
    api.get<Array<{ id: string; name: string; description: string }>>('/alerts/metrics'),
  getNotifications: (unreadOnly = false) =>
    api.get<Notification[]>('/notifications', { params: { unreadOnly } }),
  markNotificationRead: (id: string) =>
    api.post(`/notifications/${id}/read`),
  markAllNotificationsRead: () =>
    api.post('/notifications/read-all'),
}

/**
 * Reports API
 */
export const reportsApi = {
  getScheduled: () => api.get<ScheduledReport[]>('/reports/scheduled'),
  createScheduled: (report: Partial<ScheduledReport>) =>
    api.post<ScheduledReport>('/reports/scheduled', report),
  updateScheduled: (id: string, report: Partial<ScheduledReport>) =>
    api.patch<ScheduledReport>(`/reports/scheduled/${id}`, report),
  deleteScheduled: (id: string) =>
    api.delete(`/reports/scheduled/${id}`),
  toggleScheduled: (id: string, isActive: boolean) =>
    api.post(`/reports/scheduled/${id}/toggle`, { isActive }),
  runScheduled: (id: string) =>
    api.post(`/reports/scheduled/${id}/run`),
  getRuns: (id: string) =>
    api.get<Array<{ id: string; status: string; filePath: string; startedAt: string; completedAt: string }>>(
      `/reports/scheduled/${id}/runs`
    ),
  downloadRun: (runId: string) =>
    api.get<Blob>(`/reports/runs/${runId}/download`),
  getTemplates: () =>
    api.get<Array<{ id: string; name: string; description: string }>>('/reports/templates'),
}

export default api

/**
 * API Client with streaming support
 */
export const apiClient = {
  ...api,

  /**
   * Stream chat with SSE
   */
  async streamChat(
    request: ChatRequest,
    onEvent: (event: SSEEvent) => void,
    signal?: AbortSignal
  ): Promise<StreamChatResult | null> {
    const url = buildURL('/chat/stream')
    const token = getAuthToken()

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream',
    }

    if (token) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(request),
      signal,
    })

    if (!response.ok) {
      throw new APIError(response.status, response.statusText, 'Stream request failed')
    }

    const reader = response.body?.getReader()
    if (!reader) {
      throw new Error('No response body')
    }

    const decoder = new TextDecoder()
    let result: StreamChatResult | null = null
    let buffer = ''

    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')

        // Keep the last incomplete line in the buffer
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)

            if (data === '[DONE]') {
              return result
            }

            try {
              const event = JSON.parse(data) as SSEEvent
              onEvent(event)

              // Capture the result event for return value
              if (event.type === 'result') {
                result = {
                  message: event.message,
                  chartSpec: event.chartType && event.data ? {
                    type: event.chartType,
                    chart: event.data,
                  } : undefined,
                  confidence: event.confidence,
                }
              }
            } catch {
              // Skip invalid JSON
            }
          }
        }
      }
    } finally {
      reader.releaseLock()
    }

    return result
  },

  /**
   * Get chat messages for a session
   */
  async getChatMessages(sessionId: string, limit = 50): Promise<ChatMessagesResponse> {
    return api.get<ChatMessagesResponse>(`/chat/sessions/${sessionId}/messages`, {
      params: { limit },
    })
  },

  /**
   * Get recent chat sessions
   */
  async getChatSessions(): Promise<ChatSessionsResponse> {
    return api.get<ChatSessionsResponse>('/chat/sessions')
  },

  /**
   * Create a new chat session
   */
  async createChatSession(locale?: string): Promise<{ id: string; createdAt: string }> {
    return api.post('/chat/sessions', { locale })
  },

  /**
   * Get recent messages across all sessions
   */
  async getRecentMessages(limit = 20): Promise<ChatMessagesResponse> {
    return api.get<ChatMessagesResponse>('/chat/recent', { params: { limit } })
  },
}
