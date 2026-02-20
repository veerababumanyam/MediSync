import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useReports, Report, ReportRun, NewReport } from './useReports'
import { APIError } from '../services/api'

// Mock the reports API
const mockGetScheduled = vi.fn()
const mockCreateScheduled = vi.fn()
const mockUpdateScheduled = vi.fn()
const mockDeleteScheduled = vi.fn()
const mockRunScheduled = vi.fn()
const mockGetRuns = vi.fn()
const mockDownloadRun = vi.fn()

vi.mock('../services/api', () => ({
  reportsApi: {
    getScheduled: () => mockGetScheduled(),
    createScheduled: (report: NewReport) => mockCreateScheduled(report),
    updateScheduled: (id: string, updates: Partial<Report>) => mockUpdateScheduled(id, updates),
    deleteScheduled: (id: string) => mockDeleteScheduled(id),
    runScheduled: (id: string) => mockRunScheduled(id),
    getRuns: (id: string) => mockGetRuns(id),
    downloadRun: (runId: string) => mockDownloadRun(runId),
  },
  APIError: class APIError extends Error {
    constructor(public status: number, public statusText: string, message: string) {
      super(message)
      this.name = 'APIError'
    }
  },
}))

describe('useReports', () => {
  const mockReport: Report = {
    id: 'report-1',
    userId: 'user-1',
    name: 'Monthly Sales Report',
    description: 'Sales summary for the month',
    queryId: 'query-1',
    naturalLanguageQuery: 'Show monthly sales',
    sqlQuery: 'SELECT * FROM sales WHERE month = current_month',
    scheduleType: 'monthly',
    scheduleTime: '09:00',
    scheduleDay: 1,
    recipients: [{ email: 'admin@example.com', name: 'Admin' }],
    format: 'pdf',
    locale: 'en',
    includeCharts: true,
    lastRunAt: '2026-02-01T09:00:00Z',
    nextRunAt: '2026-03-01T09:00:00Z',
    isActive: true,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-02-01T09:00:00Z',
  }

  const mockRun: ReportRun = {
    id: 'run-1',
    reportId: 'report-1',
    status: 'completed',
    filePath: '/reports/run-1.pdf',
    startedAt: '2026-02-01T09:00:00Z',
    completedAt: '2026-02-01T09:01:00Z',
    error: null,
  }

  const newReport: NewReport = {
    name: 'Weekly Revenue Report',
    description: 'Revenue summary for the week',
    naturalLanguageQuery: 'Show weekly revenue',
    sqlQuery: 'SELECT * FROM revenue WHERE week = current_week',
    scheduleType: 'weekly',
    scheduleTime: '08:00',
    scheduleDay: 1,
    recipients: [{ email: 'finance@example.com', name: 'Finance Team' }],
    format: 'xlsx',
    includeCharts: false,
  }

  beforeEach(() => {
    vi.clearAllMocks()
    mockGetScheduled.mockResolvedValue([mockReport])
    mockCreateScheduled.mockResolvedValue({ ...mockReport, id: 'report-2' })
    mockUpdateScheduled.mockResolvedValue(mockReport)
    mockDeleteScheduled.mockResolvedValue(undefined)
    mockRunScheduled.mockResolvedValue(undefined)
    mockGetRuns.mockResolvedValue([mockRun])
    mockDownloadRun.mockResolvedValue(new Blob(['pdf content'], { type: 'application/pdf' }))
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('returns empty reports array initially', () => {
      const { result } = renderHook(() => useReports())

      expect(result.current.reports).toEqual([])
    })

    it('returns empty runs array initially', () => {
      const { result } = renderHook(() => useReports())

      expect(result.current.runs).toEqual([])
    })

    it('returns isLoading as true initially while loading', () => {
      const { result } = renderHook(() => useReports())

      // isLoading is true initially because refreshReports runs on mount
      expect(result.current.isLoading).toBe(true)
    })

    it('returns null error initially', () => {
      const { result } = renderHook(() => useReports())

      expect(result.current.error).toBeNull()
    })
  })

  describe('return type structure', () => {
    it('has correct return type structure', () => {
      const { result } = renderHook(() => useReports())

      expect(result.current).toHaveProperty('reports')
      expect(result.current).toHaveProperty('runs')
      expect(result.current).toHaveProperty('isLoading')
      expect(result.current).toHaveProperty('error')
      expect(result.current).toHaveProperty('createReport')
      expect(result.current).toHaveProperty('updateReport')
      expect(result.current).toHaveProperty('deleteReport')
      expect(result.current).toHaveProperty('runReport')
      expect(result.current).toHaveProperty('downloadReport')
    })

    it('provides createReport function', () => {
      const { result } = renderHook(() => useReports())

      expect(typeof result.current.createReport).toBe('function')
    })

    it('provides updateReport function', () => {
      const { result } = renderHook(() => useReports())

      expect(typeof result.current.updateReport).toBe('function')
    })

    it('provides deleteReport function', () => {
      const { result } = renderHook(() => useReports())

      expect(typeof result.current.deleteReport).toBe('function')
    })

    it('provides runReport function', () => {
      const { result } = renderHook(() => useReports())

      expect(typeof result.current.runReport).toBe('function')
    })

    it('provides downloadReport function', () => {
      const { result } = renderHook(() => useReports())

      expect(typeof result.current.downloadReport).toBe('function')
    })
  })

  describe('loading reports on mount', () => {
    it('loads reports on mount', async () => {
      renderHook(() => useReports())

      await waitFor(() => {
        expect(mockGetScheduled).toHaveBeenCalled()
      })
    })

    it('sets reports after loading', async () => {
      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
        expect(result.current.reports[0]).toEqual(mockReport)
      })
    })
  })

  describe('createReport', () => {
    it('creates a new report', async () => {
      const { result } = renderHook(() => useReports())

      // Wait for initial load
      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        await result.current.createReport(newReport)
      })

      expect(mockCreateScheduled).toHaveBeenCalledWith(newReport)
      expect(result.current.reports).toHaveLength(2)
    })

    it('sets isLoading to true during creation', async () => {
      let resolveCreate: () => void
      mockCreateScheduled.mockImplementation(() => new Promise(resolve => {
        resolveCreate = resolve as () => void
      }))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      act(() => {
        result.current.createReport(newReport)
      })

      expect(result.current.isLoading).toBe(true)

      await act(async () => {
        resolveCreate!()
      })

      expect(result.current.isLoading).toBe(false)
    })

    it('sets error on creation failure', async () => {
      mockCreateScheduled.mockRejectedValueOnce(new Error('Creation failed'))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        try {
          await result.current.createReport(newReport)
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Creation failed')
    })

    it('handles APIError with custom message', async () => {
      const apiError = new APIError(400, 'Bad Request', 'Invalid schedule type')
      mockCreateScheduled.mockRejectedValueOnce(apiError)

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        try {
          await result.current.createReport(newReport)
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Invalid schedule type')
    })
  })

  describe('updateReport', () => {
    it('updates an existing report', async () => {
      const updates = { name: 'Updated Report Name' }
      mockUpdateScheduled.mockResolvedValueOnce({ ...mockReport, ...updates })

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        await result.current.updateReport('report-1', updates)
      })

      expect(mockUpdateScheduled).toHaveBeenCalledWith('report-1', updates)
      expect(result.current.reports[0].name).toBe('Updated Report Name')
    })

    it('sets error on update failure', async () => {
      mockUpdateScheduled.mockRejectedValueOnce(new Error('Update failed'))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        try {
          await result.current.updateReport('report-1', { name: 'New Name' })
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Update failed')
    })
  })

  describe('deleteReport', () => {
    it('deletes a report', async () => {
      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        await result.current.deleteReport('report-1')
      })

      expect(mockDeleteScheduled).toHaveBeenCalledWith('report-1')
      expect(result.current.reports).toHaveLength(0)
    })

    it('sets error on delete failure', async () => {
      mockDeleteScheduled.mockRejectedValueOnce(new Error('Delete failed'))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        try {
          await result.current.deleteReport('report-1')
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Delete failed')
    })
  })

  describe('runReport', () => {
    it('triggers report generation', async () => {
      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        await result.current.runReport('report-1')
      })

      expect(mockRunScheduled).toHaveBeenCalledWith('report-1')
    })

    it('refreshes reports after running', async () => {
      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(mockGetScheduled).toHaveBeenCalledTimes(1)
      })

      await act(async () => {
        await result.current.runReport('report-1')
      })

      // Should have called getScheduled again for refresh
      expect(mockGetScheduled).toHaveBeenCalledTimes(2)
    })

    it('sets error on run failure', async () => {
      mockRunScheduled.mockRejectedValueOnce(new Error('Run failed'))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.reports).toHaveLength(1)
      })

      await act(async () => {
        try {
          await result.current.runReport('report-1')
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Run failed')
    })
  })

  describe('loadRuns', () => {
    it('loads runs for a specific report', async () => {
      const { result } = renderHook(() => useReports())

      await act(async () => {
        await result.current.loadRuns('report-1')
      })

      expect(mockGetRuns).toHaveBeenCalledWith('report-1')
      expect(result.current.runs).toHaveLength(1)
      expect(result.current.runs[0]).toEqual(mockRun)
    })

    it('sets error on load runs failure', async () => {
      mockGetRuns.mockRejectedValueOnce(new Error('Failed to load runs'))

      const { result } = renderHook(() => useReports())

      await act(async () => {
        await result.current.loadRuns('report-1')
      })

      expect(result.current.error).toBe('Failed to load runs')
    })
  })

  describe('downloadReport', () => {
    it('downloads a report file', async () => {
      // Mock URL methods
      const mockCreateObjectURL = vi.fn(() => 'blob:mock-url')
      const mockRevokeObjectURL = vi.fn()
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = mockRevokeObjectURL

      const { result } = renderHook(() => useReports())

      await act(async () => {
        await result.current.downloadReport('run-1', 'pdf')
      })

      expect(mockDownloadRun).toHaveBeenCalledWith('run-1')
      expect(mockCreateObjectURL).toHaveBeenCalled()
    })

    it('sets error on download failure', async () => {
      mockDownloadRun.mockRejectedValueOnce(new Error('Download failed'))

      const { result } = renderHook(() => useReports())

      await act(async () => {
        try {
          await result.current.downloadReport('run-1', 'pdf')
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Download failed')
    })
  })

  describe('refreshReports', () => {
    it('refreshes the reports list', async () => {
      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(mockGetScheduled).toHaveBeenCalledTimes(1)
      })

      // Add a new report to mock
      mockGetScheduled.mockResolvedValueOnce([mockReport, { ...mockReport, id: 'report-2' }])

      await act(async () => {
        await result.current.refreshReports()
      })

      expect(mockGetScheduled).toHaveBeenCalledTimes(2)
      expect(result.current.reports).toHaveLength(2)
    })
  })

  describe('clearError', () => {
    it('clears the error state', async () => {
      mockGetScheduled.mockRejectedValueOnce(new Error('Load failed'))

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.error).toBe('Load failed')
      })

      act(() => {
        result.current.clearError()
      })

      expect(result.current.error).toBeNull()
    })
  })

  describe('error handling', () => {
    it('handles non-Error objects', async () => {
      mockGetScheduled.mockRejectedValueOnce('Unknown error')

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.error).toBe('Failed to load reports')
      })
    })

    it('handles undefined/null errors', async () => {
      mockGetScheduled.mockRejectedValueOnce(null)

      const { result } = renderHook(() => useReports())

      await waitFor(() => {
        expect(result.current.error).toBe('Failed to load reports')
      })
    })
  })
})
