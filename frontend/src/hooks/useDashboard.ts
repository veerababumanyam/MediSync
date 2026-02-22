/**
 * useDashboard Hook
 *
 * React hook for managing dashboard pinned charts with:
 * - State management for pinned charts
 * - CRUD operations (pin, update, delete, reorder)
 * - Chart refresh functionality
 * - Loading and error states
 *
 * @module hooks/useDashboard
 */
import { useCallback, useEffect, useState } from 'react'
import { dashboardApi } from '../services/api'
import type { PinnedChart } from '../services/api'

/**
 * Type for creating a new chart (excludes auto-generated fields)
 */
export type NewChart = Omit<Partial<PinnedChart>, 'id' | 'userId' | 'createdAt' | 'updatedAt'>

/**
 * Return type for useDashboard hook
 */
export interface UseDashboardReturn {
  /** List of pinned charts */
  charts: PinnedChart[]
  /** Whether charts are currently loading */
  isLoading: boolean
  /** Any error that occurred during operations */
  error: string | null
  /** Pin a new chart to the dashboard */
  pinChart: (chart: NewChart) => Promise<void>
  /** Update an existing chart */
  updateChart: (id: string, updates: Partial<PinnedChart>) => Promise<void>
  /** Delete a chart from the dashboard */
  deleteChart: (id: string) => Promise<void>
  /** Refresh a chart's data */
  refreshChart: (id: string) => Promise<void>
  /** Refresh all charts on the dashboard */
  refreshAll: () => Promise<void>
  /** Reorder charts on the dashboard */
  reorderCharts: (order: string[]) => Promise<void>
}

/**
 * Hook for managing dashboard pinned charts
 *
 * @example
 * ```tsx
 * function Dashboard() {
 *   const { charts, isLoading, error, pinChart, deleteChart } = useDashboard()
 *
 *   if (isLoading) return <Spinner />
 *   if (error) return <Error message={error} />
 *
 *   return (
 *     <div>
 *       {charts.map(chart => (
 *         <ChartCard key={chart.id} chart={chart} onDelete={deleteChart} />
 *       ))}
 *     </div>
 *   )
 * }
 * ```
 */
export function useDashboard(): UseDashboardReturn {
  const [charts, setCharts] = useState<PinnedChart[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  /**
   * Load charts on mount
   */
  useEffect(() => {
    let mounted = true

    async function loadCharts() {
      setIsLoading(true)
      setError(null)

      try {
        const data = await dashboardApi.getCharts()
        if (mounted) {
          setCharts(data)
        }
      } catch (err) {
        if (mounted) {
          const message = err instanceof Error ? err.message : 'Failed to load charts'
          setError(message)
          console.error('Failed to load charts:', err)
        }
      } finally {
        if (mounted) {
          setIsLoading(false)
        }
      }
    }

    loadCharts()

    return () => {
      mounted = false
    }
  }, [])

  /**
   * Pin a new chart to the dashboard
   */
  const pinChart = useCallback(async (chart: NewChart) => {
    setError(null)

    try {
      const newChart = await dashboardApi.pinChart(chart)
      setCharts(prev => [...prev, newChart])
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to pin chart'
      setError(message)
      console.error('Failed to pin chart:', err)
      throw err
    }
  }, [])

  /**
   * Update an existing chart
   */
  const updateChart = useCallback(async (id: string, updates: Partial<PinnedChart>) => {
    setError(null)

    // Optimistic update
    setCharts(prev =>
      prev.map(chart =>
        chart.id === id
          ? { ...chart, ...updates, updatedAt: new Date().toISOString() }
          : chart
      )
    )

    try {
      const updatedChart = await dashboardApi.updateChart(id, updates)
      setCharts(prev =>
        prev.map(chart =>
          chart.id === id ? updatedChart : chart
        )
      )
    } catch (err) {
      // Revert on error - refetch from server
      const message = err instanceof Error ? err.message : 'Failed to update chart'
      setError(message)
      console.error('Failed to update chart:', err)

      // Refetch to get correct state
      try {
        const data = await dashboardApi.getCharts()
        setCharts(data)
      } catch {
        // Ignore refetch error, keep the optimistic update reverted
      }

      throw err
    }
  }, [])

  /**
   * Delete a chart from the dashboard
   */
  const deleteChart = useCallback(async (id: string) => {
    setError(null)

    // Optimistic delete
    const previousCharts = charts
    setCharts(prev => prev.filter(chart => chart.id !== id))

    try {
      await dashboardApi.deleteChart(id)
    } catch (err) {
      // Revert on error
      const message = err instanceof Error ? err.message : 'Failed to delete chart'
      setError(message)
      console.error('Failed to delete chart:', err)
      setCharts(previousCharts)
      throw err
    }
  }, [charts])

  /**
   * Refresh a chart's data
   */
  const refreshChart = useCallback(async (id: string) => {
    setError(null)

    try {
      const updatedChart = await dashboardApi.refreshChart(id)
      setCharts(prev =>
        prev.map(chart =>
          chart.id === id ? updatedChart : chart
        )
      )
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to refresh chart'
      setError(message)
      console.error('Failed to refresh chart:', err)
      throw err
    }
  }, [])

  /**
   * Refresh all charts on the dashboard
   */
  const refreshAll = useCallback(async () => {
    setError(null)
    setIsLoading(true)

    try {
      // Refresh each chart sequentially
      const refreshedCharts = await Promise.all(
        charts.map(chart => dashboardApi.refreshChart(chart.id))
      )
      setCharts(refreshedCharts)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to refresh charts'
      setError(message)
      console.error('Failed to refresh all charts:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [charts])

  /**
   * Reorder charts on the dashboard
   */
  const reorderCharts = useCallback(async (order: string[]) => {
    setError(null)

    // Optimistic reorder - create a map for quick lookup
    const chartMap = new Map(charts.map(chart => [chart.id, chart]))
    const reorderedCharts = order
      .map(id => chartMap.get(id))
      .filter((chart): chart is PinnedChart => chart !== undefined)

    // Add any charts not in the order array at the end
    const orderedIds = new Set(order)
    const remainingCharts = charts.filter(chart => !orderedIds.has(chart.id))
    const newOrder = [...reorderedCharts, ...remainingCharts]

    setCharts(newOrder)

    try {
      // Build positions array for API
      const positions = newOrder.map((chart, index) => ({
        id: chart.id,
        position: {
          row: Math.floor(index / 3),
          col: index % 3,
          size: chart.position.size,
        },
      }))

      await dashboardApi.reorderCharts(positions)
    } catch (err) {
      // Revert on error
      const message = err instanceof Error ? err.message : 'Failed to reorder charts'
      setError(message)
      console.error('Failed to reorder charts:', err)

      // Refetch to get correct state
      try {
        const data = await dashboardApi.getCharts()
        setCharts(data)
      } catch {
        // Ignore refetch error
      }

      throw err
    }
  }, [charts])

  return {
    charts,
    isLoading,
    error,
    pinChart,
    updateChart,
    deleteChart,
    refreshChart,
    refreshAll,
    reorderCharts,
  }
}

export default useDashboard
