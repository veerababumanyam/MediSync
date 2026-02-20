/**
 * useReports Hook
 *
 * React hook for managing scheduled reports with:
 * - State management for reports and report runs
 * - CRUD operations for scheduled reports
 * - Report generation and download functionality
 * - Loading and error state management
 *
 * @module hooks/useReports
 */
import { useCallback, useEffect, useState } from 'react'
import { reportsApi, ScheduledReport, APIError } from '../services/api'

/**
 * Report type matching ScheduledReport from API
 */
export type Report = ScheduledReport

/**
 * Report run history entry
 */
export interface ReportRun {
  id: string
  reportId: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  filePath: string | null
  startedAt: string
  completedAt: string | null
  error: string | null
}

/**
 * New report creation payload
 */
export interface NewReport {
  name: string
  description?: string | null
  queryId?: string | null
  naturalLanguageQuery: string
  sqlQuery: string
  scheduleType: 'daily' | 'weekly' | 'monthly' | 'quarterly'
  scheduleTime: string
  scheduleDay?: number | null
  recipients: Array<{ email: string; name: string }>
  format: 'pdf' | 'xlsx' | 'csv'
  locale?: 'en' | 'ar'
  includeCharts?: boolean
}

/**
 * Return type for useReports hook
 */
export interface UseReportsReturn {
  /** List of scheduled reports */
  reports: Report[]
  /** List of report runs */
  runs: ReportRun[]
  /** Loading state */
  isLoading: boolean
  /** Error state */
  error: string | null
  /** Create a new scheduled report */
  createReport: (report: NewReport) => Promise<void>
  /** Update an existing report */
  updateReport: (id: string, updates: Partial<Report>) => Promise<void>
  /** Delete a report */
  deleteReport: (id: string) => Promise<void>
  /** Trigger report generation */
  runReport: (id: string) => Promise<void>
  /** Download a completed report run */
  downloadReport: (runId: string, format: string) => Promise<void>
  /** Refresh reports list */
  refreshReports: () => Promise<void>
  /** Load runs for a specific report */
  loadRuns: (reportId: string) => Promise<void>
  /** Clear error state */
  clearError: () => void
}

/**
 * Hook for managing scheduled reports
 */
export function useReports(): UseReportsReturn {
  const [reports, setReports] = useState<Report[]>([])
  const [runs, setRuns] = useState<ReportRun[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  /**
   * Refresh the list of scheduled reports
   */
  const refreshReports = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await reportsApi.getScheduled()
      setReports(data)
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to load reports'
      setError(message)
      console.error('Failed to load reports:', err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Load all scheduled reports on mount
   */
  useEffect(() => {
    refreshReports()
  }, [refreshReports])

  /**
   * Create a new scheduled report
   */
  const createReport = useCallback(async (report: NewReport) => {
    setIsLoading(true)
    setError(null)

    try {
      const newReport = await reportsApi.createScheduled(report)
      setReports(prev => [...prev, newReport])
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to create report'
      setError(message)
      console.error('Failed to create report:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Update an existing report
   */
  const updateReport = useCallback(async (id: string, updates: Partial<Report>) => {
    setIsLoading(true)
    setError(null)

    try {
      const updatedReport = await reportsApi.updateScheduled(id, updates)
      setReports(prev =>
        prev.map(report =>
          report.id === id ? updatedReport : report
        )
      )
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to update report'
      setError(message)
      console.error('Failed to update report:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Delete a report
   */
  const deleteReport = useCallback(async (id: string) => {
    setIsLoading(true)
    setError(null)

    try {
      await reportsApi.deleteScheduled(id)
      setReports(prev => prev.filter(report => report.id !== id))
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to delete report'
      setError(message)
      console.error('Failed to delete report:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Trigger report generation
   */
  const runReport = useCallback(async (id: string) => {
    setIsLoading(true)
    setError(null)

    try {
      await reportsApi.runScheduled(id)
      // Refresh reports to get updated lastRunAt/nextRunAt
      await refreshReports()
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to run report'
      setError(message)
      console.error('Failed to run report:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [refreshReports])

  /**
   * Load runs for a specific report
   */
  const loadRuns = useCallback(async (reportId: string) => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await reportsApi.getRuns(reportId)
      setRuns(data as ReportRun[])
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to load report runs'
      setError(message)
      console.error('Failed to load report runs:', err)
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Download a completed report run
   */
  const downloadReport = useCallback(async (runId: string, format: string) => {
    setIsLoading(true)
    setError(null)

    try {
      const blob = await reportsApi.downloadRun(runId)

      // Create download link
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `report-${runId}.${format}`

      // Trigger download
      document.body.appendChild(link)
      link.click()

      // Cleanup
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (err) {
      const message = err instanceof APIError
        ? err.message
        : err instanceof Error
          ? err.message
          : 'Failed to download report'
      setError(message)
      console.error('Failed to download report:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setError(null)
  }, [])

  return {
    reports,
    runs,
    isLoading,
    error,
    createReport,
    updateReport,
    deleteReport,
    runReport,
    downloadReport,
    refreshReports,
    loadRuns,
    clearError,
  }
}

export default useReports
