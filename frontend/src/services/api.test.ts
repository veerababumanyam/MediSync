import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { APIError, preferencesApi, chatApi, dashboardApi, alertsApi, reportsApi } from './api'

// Mock fetch
const mockFetch = vi.fn()
global.fetch = mockFetch

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => { store[key] = value }),
    removeItem: vi.fn((key: string) => { delete store[key] }),
    clear: vi.fn(() => { store = {} }),
  }
})()
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

describe('APIError', () => {
  it('creates error with status and message', () => {
    const error = new APIError(404, 'Not Found', 'Resource not found')

    expect(error.status).toBe(404)
    expect(error.statusText).toBe('Not Found')
    expect(error.message).toBe('Resource not found')
    expect(error.name).toBe('APIError')
  })

  it('creates error with data', () => {
    const data = { code: 'USER_NOT_FOUND' }
    const error = new APIError(404, 'Not Found', 'User not found', data)

    expect(error.data).toEqual(data)
  })

  it('is instance of Error', () => {
    const error = new APIError(500, 'Server Error', 'Internal error')

    expect(error).toBeInstanceOf(Error)
    expect(error).toBeInstanceOf(APIError)
  })
})

describe('preferencesApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('gets preferences', async () => {
    const mockPrefs = { locale: 'en', theme: 'light' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockPrefs),
    })

    const result = await preferencesApi.get()

    expect(mockFetch).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockPrefs)
  })

  it('updates preferences', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ locale: 'ar' }),
    })

    await preferencesApi.update({ locale: 'ar' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('PATCH')
    expect(call[1].body).toBe(JSON.stringify({ locale: 'ar' }))
  })
})

describe('chatApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('sends a query', async () => {
    const mockResponse = {
      id: 'msg-123',
      role: 'assistant',
      content: 'Here is the data you requested',
    }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockResponse),
    })

    const result = await chatApi.query({
      query: 'Show me revenue',
      locale: 'en',
      sessionId: 'sess-1',
    })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockResponse)
  })

  it('sends drilldown request', async () => {
    const mockResponse = { id: 'drill-123', data: [] }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockResponse),
    })

    await chatApi.drilldown({
      queryId: 'query-1',
      filters: { department: 'Cardiology' },
    })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('POST')
  })

  it('exports data', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }),
      blob: () => Promise.resolve(new Blob()),
    })

    await chatApi.export('xlsx', { query: 'test', results: [] })

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })
})

describe('dashboardApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('gets charts', async () => {
    const mockCharts = [{ id: 'chart-1', title: 'Revenue Chart' }]
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockCharts),
    })

    const result = await dashboardApi.getCharts()

    expect(mockFetch).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockCharts)
  })

  it('pins a chart', async () => {
    const mockChart = { id: 'chart-1', title: 'New Chart' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockChart),
    })

    const result = await dashboardApi.pinChart({
      title: 'New Chart',
      chartType: 'bar',
    })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockChart)
  })

  it('deletes a chart', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 204,
      headers: new Headers(),
      blob: () => Promise.resolve(new Blob()),
      text: () => Promise.resolve(''),
      json: () => Promise.resolve(null),
    })

    await dashboardApi.deleteChart('chart-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('DELETE')
  })

  it('updates a chart', async () => {
    const mockChart = { id: 'chart-1', title: 'Updated Chart' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockChart),
    })

    await dashboardApi.updateChart('chart-1', { title: 'Updated Chart' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('PATCH')
  })

  it('reorders charts', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await dashboardApi.reorderCharts([
      { id: 'chart-1', position: { row: 0, col: 0, size: 1 } },
    ])

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('refreshes a chart', async () => {
    const mockChart = { id: 'chart-1', title: 'Refreshed Chart' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockChart),
    })

    const result = await dashboardApi.refreshChart('chart-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
    expect(result).toEqual(mockChart)
  })

  it('gets quick actions', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve([]),
    })

    await dashboardApi.getQuickActions()

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })
})

describe('alertsApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('gets alert rules', async () => {
    const mockRules = [{ id: 'rule-1', name: 'Revenue Alert' }]
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockRules),
    })

    const result = await alertsApi.getRules()

    expect(result).toEqual(mockRules)
  })

  it('creates an alert rule', async () => {
    const mockRule = { id: 'rule-1', name: 'New Alert' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockRule),
    })

    await alertsApi.createRule({ name: 'New Alert' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('POST')
  })

  it('gets notifications', async () => {
    const mockNotifications = [{ id: 'notif-1', title: 'Alert!' }]
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockNotifications),
    })

    const result = await alertsApi.getNotifications(true)

    expect(result).toEqual(mockNotifications)
  })
})

describe('reportsApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('gets scheduled reports', async () => {
    const mockReports = [{ id: 'report-1', name: 'Weekly Report' }]
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockReports),
    })

    const result = await reportsApi.getScheduled()

    expect(result).toEqual(mockReports)
  })

  it('creates a scheduled report', async () => {
    const mockReport = { id: 'report-1', name: 'New Report' }
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockReport),
    })

    await reportsApi.createScheduled({ name: 'New Report' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('POST')
  })
})

describe('Authentication', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('includes auth token in request headers', async () => {
    localStorageMock.setItem('medisync-token', 'test-token-123')

    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await preferencesApi.get()

    const call = mockFetch.mock.calls[0]
    expect(call[1].headers.Authorization).toBe('Bearer test-token-123')
  })

  it('does not include auth header when no token', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await preferencesApi.get()

    const call = mockFetch.mock.calls[0]
    expect(call[1].headers.Authorization).toBeUndefined()
  })
})

describe('API Error Handling', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('throws APIError on non-OK response with message in body', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 400,
      statusText: 'Bad Request',
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ message: 'Invalid input data' }),
    })

    try {
      await preferencesApi.get()
      expect.fail('Should have thrown')
    } catch (error) {
      expect(error).toBeInstanceOf(APIError)
      expect((error as APIError).status).toBe(400)
      expect((error as APIError).message).toBe('Invalid input data')
    }
  })

  it('throws APIError on non-OK response using statusText', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await expect(preferencesApi.get()).rejects.toThrow('Internal Server Error')
  })

  it('handles text response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 200,
      statusText: 'OK',
      headers: new Headers({ 'content-type': 'text/plain' }),
      text: () => Promise.resolve('Success'),
    })

    // Text responses should work
    await expect(preferencesApi.get()).resolves.toBeDefined()
  })
})

describe('api object methods', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('api.get makes GET request', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ data: 'test' }),
    })

    await api.get('/test')

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('GET')
  })

  it('api.post makes POST request with body', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: 1 }),
    })

    await api.post('/test', { name: 'test' })

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('POST')
    expect(call[1].body).toBe(JSON.stringify({ name: 'test' }))
  })

  it('api.patch makes PATCH request', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await api.patch('/test', { name: 'updated' })

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('PATCH')
  })

  it('api.put makes PUT request', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await api.put('/test', { name: 'replaced' })

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('PUT')
  })

  it('api.delete makes DELETE request', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await api.delete('/test/1')

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('DELETE')
  })

  it('includes query params in URL', async () => {
    const { api } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve([]),
    })

    await api.get('/test', { params: { page: 1, limit: 10 } })

    const call = mockFetch.mock.calls[0]
    expect(call[0]).toContain('page=1')
    expect(call[0]).toContain('limit=10')
  })
})

describe('Additional alertsApi methods', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('updates an alert rule', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: 'rule-1', name: 'Updated' }),
    })

    await alertsApi.updateRule('rule-1', { name: 'Updated' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('deletes an alert rule', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 204,
      headers: new Headers(),
      blob: () => Promise.resolve(new Blob()),
      text: () => Promise.resolve(''),
      json: () => Promise.resolve(null),
    })

    await alertsApi.deleteRule('rule-1')

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('DELETE')
  })

  it('toggles an alert rule', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await alertsApi.toggleRule('rule-1', false)

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('tests an alert rule', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ triggered: true }),
    })

    await alertsApi.testRule('rule-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('gets alert metrics', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve([]),
    })

    await alertsApi.getMetrics()

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('marks notification as read', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await alertsApi.markNotificationRead('notif-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('marks all notifications as read', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await alertsApi.markAllNotificationsRead()

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })
})

describe('Additional reportsApi methods', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('updates a scheduled report', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: 'report-1' }),
    })

    await reportsApi.updateScheduled('report-1', { name: 'Updated' })

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('deletes a scheduled report', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      status: 204,
      headers: new Headers(),
      blob: () => Promise.resolve(new Blob()),
      text: () => Promise.resolve(''),
      json: () => Promise.resolve(null),
    })

    await reportsApi.deleteScheduled('report-1')

    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('DELETE')
  })

  it('toggles a scheduled report', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({}),
    })

    await reportsApi.toggleScheduled('report-1', true)

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('runs a scheduled report', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: 'run-1' }),
    })

    await reportsApi.runScheduled('report-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('gets report runs', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve([]),
    })

    await reportsApi.getRuns('report-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('downloads a report run', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/pdf' }),
      blob: () => Promise.resolve(new Blob()),
    })

    await reportsApi.downloadRun('run-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })

  it('gets report templates', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve([]),
    })

    await reportsApi.getTemplates()

    expect(mockFetch).toHaveBeenCalledTimes(1)
  })
})

describe('apiClient methods', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
  })

  it('getChatMessages fetches messages for session', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ messages: [] }),
    })

    await apiClient.getChatMessages('session-1')

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[0]).toContain('/chat/sessions/session-1/messages')
  })

  it('getChatSessions fetches sessions', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ sessions: [] }),
    })

    await apiClient.getChatSessions()

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[0]).toContain('/chat/sessions')
  })

  it('createChatSession creates new session', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: 'new-session', createdAt: '2024-01-01' }),
    })

    await apiClient.createChatSession('en')

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[1].method).toBe('POST')
  })

  it('getRecentMessages fetches recent messages', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ messages: [] }),
    })

    await apiClient.getRecentMessages(10)

    expect(mockFetch).toHaveBeenCalledTimes(1)
    const call = mockFetch.mock.calls[0]
    expect(call[0]).toContain('/chat/recent')
    expect(call[0]).toContain('limit=10')
  })

  it('streamChat throws APIError on failed request', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 401,
      statusText: 'Unauthorized',
      headers: new Headers(),
    })

    await expect(apiClient.streamChat({ query: 'test', sessionId: 's1', locale: 'en' }, vi.fn()))
      .rejects.toThrow('Stream request failed')
  })

  it('streamChat throws error when no body', async () => {
    const { apiClient } = await import('./api')
    mockFetch.mockResolvedValueOnce({
      ok: true,
      headers: new Headers({ 'content-type': 'text/event-stream' }),
      body: null,
    })

    await expect(apiClient.streamChat({ query: 'test', sessionId: 's1', locale: 'en' }, vi.fn()))
      .rejects.toThrow('No response body')
  })
})
