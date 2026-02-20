import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useDashboard } from './useDashboard'
import { PinnedChart } from '../services/api'

// Mock the dashboard API
const mockGetCharts = vi.fn()
const mockPinChart = vi.fn()
const mockUpdateChart = vi.fn()
const mockDeleteChart = vi.fn()
const mockRefreshChart = vi.fn()
const mockReorderCharts = vi.fn()

vi.mock('../services/api', () => ({
  dashboardApi: {
    getCharts: () => mockGetCharts(),
    pinChart: (chart: Partial<PinnedChart>) => mockPinChart(chart),
    updateChart: (id: string, chart: Partial<PinnedChart>) => mockUpdateChart(id, chart),
    deleteChart: (id: string) => mockDeleteChart(id),
    refreshChart: (id: string) => mockRefreshChart(id),
    reorderCharts: (positions: Array<{ id: string; position: { row: number; col: number; size: number } }>) =>
      mockReorderCharts(positions),
  },
}))

// Sample chart data
const createMockChart = (id: string, overrides?: Partial<PinnedChart>): PinnedChart => ({
  id,
  userId: 'user-1',
  title: `Chart ${id}`,
  queryId: `query-${id}`,
  naturalLanguageQuery: `Show me data for chart ${id}`,
  sqlQuery: `SELECT * FROM data WHERE id = '${id}'`,
  chartSpec: { type: 'bar', data: [] },
  chartType: 'bar',
  refreshInterval: 300,
  locale: 'en',
  position: { row: 0, col: 0, size: 1 },
  lastRefreshedAt: new Date().toISOString(),
  isActive: true,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  ...overrides,
})

describe('useDashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Default: return empty array
    mockGetCharts.mockResolvedValue([])
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('returns empty charts array initially', () => {
      const { result } = renderHook(() => useDashboard())

      expect(result.current.charts).toEqual([])
    })

    it('starts with isLoading set to true', () => {
      const { result } = renderHook(() => useDashboard())

      expect(result.current.isLoading).toBe(true)
    })

    it('starts with error set to null', () => {
      const { result } = renderHook(() => useDashboard())

      expect(result.current.error).toBeNull()
    })
  })

  describe('loading charts', () => {
    it('loads charts on mount', async () => {
      const mockCharts = [createMockChart('1'), createMockChart('2')]
      mockGetCharts.mockResolvedValue(mockCharts)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => {
        expect(mockGetCharts).toHaveBeenCalled()
        expect(result.current.charts).toHaveLength(2)
        expect(result.current.isLoading).toBe(false)
      })
    })

    it('sets isLoading to false after loading', async () => {
      mockGetCharts.mockResolvedValue([])

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })
    })

    it('sets error when loading fails', async () => {
      mockGetCharts.mockRejectedValue(new Error('Failed to load'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => {
        expect(result.current.error).toBe('Failed to load')
        expect(result.current.isLoading).toBe(false)
      })
    })
  })

  describe('pinChart', () => {
    it('adds a new chart to the list', async () => {
      const newChart = createMockChart('new')
      mockGetCharts.mockResolvedValue([])
      mockPinChart.mockResolvedValue(newChart)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.isLoading).toBe(false))

      await act(async () => {
        await result.current.pinChart({
          title: 'New Chart',
          naturalLanguageQuery: 'Show me data',
          sqlQuery: 'SELECT * FROM data',
          chartType: 'bar',
        })
      })

      expect(mockPinChart).toHaveBeenCalled()
      expect(result.current.charts).toHaveLength(1)
      expect(result.current.charts[0].id).toBe('new')
    })

    it('sets error when pinning fails', async () => {
      mockGetCharts.mockResolvedValue([])
      mockPinChart.mockRejectedValue(new Error('Failed to pin'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.isLoading).toBe(false))

      await act(async () => {
        try {
          await result.current.pinChart({ title: 'New Chart' })
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Failed to pin')
    })

    it('clears previous error on successful pin', async () => {
      const newChart = createMockChart('new')
      mockGetCharts.mockResolvedValue([])
      mockPinChart.mockRejectedValueOnce(new Error('First error'))
      mockPinChart.mockResolvedValueOnce(newChart)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.isLoading).toBe(false))

      // First attempt fails
      await act(async () => {
        try {
          await result.current.pinChart({ title: 'Chart 1' })
        } catch {
          // Expected
        }
      })

      expect(result.current.error).toBe('First error')

      // Second attempt succeeds
      await act(async () => {
        await result.current.pinChart({ title: 'Chart 2' })
      })

      expect(result.current.error).toBeNull()
    })
  })

  describe('updateChart', () => {
    it('updates an existing chart', async () => {
      const chart1 = createMockChart('1')
      const updatedChart = { ...chart1, title: 'Updated Title' }
      mockGetCharts.mockResolvedValue([chart1])
      mockUpdateChart.mockResolvedValue(updatedChart)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      await act(async () => {
        await result.current.updateChart('1', { title: 'Updated Title' })
      })

      expect(mockUpdateChart).toHaveBeenCalledWith('1', { title: 'Updated Title' })
      expect(result.current.charts[0].title).toBe('Updated Title')
    })

    it('optimistically updates the chart', async () => {
      const chart1 = createMockChart('1')
      mockGetCharts.mockResolvedValue([chart1])

      // Delay the API response
      let resolveUpdate: (value: PinnedChart) => void
      mockUpdateChart.mockImplementation(
        () =>
          new Promise(resolve => {
            resolveUpdate = resolve
          })
      )

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      act(() => {
        result.current.updateChart('1', { title: 'Updated Title' })
      })

      // Optimistic update should be visible immediately
      expect(result.current.charts[0].title).toBe('Updated Title')

      // Resolve the promise
      await act(async () => {
        resolveUpdate!({ ...chart1, title: 'Updated Title' })
      })
    })

    it('reverts update on error', async () => {
      const chart1 = createMockChart('1')
      mockGetCharts.mockResolvedValue([chart1])
      mockUpdateChart.mockRejectedValue(new Error('Update failed'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      await act(async () => {
        try {
          await result.current.updateChart('1', { title: 'Updated Title' })
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Update failed')
    })
  })

  describe('deleteChart', () => {
    it('removes a chart from the list', async () => {
      const chart1 = createMockChart('1')
      const chart2 = createMockChart('2')
      mockGetCharts.mockResolvedValue([chart1, chart2])
      mockDeleteChart.mockResolvedValue(undefined)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(2))

      await act(async () => {
        await result.current.deleteChart('1')
      })

      expect(mockDeleteChart).toHaveBeenCalledWith('1')
      expect(result.current.charts).toHaveLength(1)
      expect(result.current.charts[0].id).toBe('2')
    })

    it('optimistically removes the chart', async () => {
      const chart1 = createMockChart('1')
      mockGetCharts.mockResolvedValue([chart1])

      // Delay the API response
      let resolveDelete: () => void
      mockDeleteChart.mockImplementation(
        () =>
          new Promise(resolve => {
            resolveDelete = resolve
          })
      )

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      act(() => {
        result.current.deleteChart('1')
      })

      // Optimistic delete should be visible immediately
      expect(result.current.charts).toHaveLength(0)

      // Resolve the promise
      await act(async () => {
        resolveDelete!()
      })
    })

    it('reverts delete on error', async () => {
      const chart1 = createMockChart('1')
      mockGetCharts.mockResolvedValue([chart1])
      mockDeleteChart.mockRejectedValue(new Error('Delete failed'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      await act(async () => {
        try {
          await result.current.deleteChart('1')
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Delete failed')
    })
  })

  describe('refreshChart', () => {
    it('refreshes a chart', async () => {
      const chart1 = createMockChart('1', { lastRefreshedAt: '2024-01-01T00:00:00Z' })
      const refreshedChart = {
        ...chart1,
        lastRefreshedAt: '2024-01-02T00:00:00Z',
      }
      mockGetCharts.mockResolvedValue([chart1])
      mockRefreshChart.mockResolvedValue(refreshedChart)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      await act(async () => {
        await result.current.refreshChart('1')
      })

      expect(mockRefreshChart).toHaveBeenCalledWith('1')
      expect(result.current.charts[0].lastRefreshedAt).toBe('2024-01-02T00:00:00Z')
    })

    it('sets error when refresh fails', async () => {
      const chart1 = createMockChart('1')
      mockGetCharts.mockResolvedValue([chart1])
      mockRefreshChart.mockRejectedValue(new Error('Refresh failed'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(1))

      await act(async () => {
        try {
          await result.current.refreshChart('1')
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Refresh failed')
    })
  })

  describe('reorderCharts', () => {
    it('reorders charts according to new order', async () => {
      const chart1 = createMockChart('1')
      const chart2 = createMockChart('2')
      const chart3 = createMockChart('3')
      mockGetCharts.mockResolvedValue([chart1, chart2, chart3])
      mockReorderCharts.mockResolvedValue(undefined)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(3))

      await act(async () => {
        await result.current.reorderCharts(['3', '1', '2'])
      })

      expect(result.current.charts[0].id).toBe('3')
      expect(result.current.charts[1].id).toBe('1')
      expect(result.current.charts[2].id).toBe('2')
    })

    it('calls API with correct positions', async () => {
      const chart1 = createMockChart('1', { position: { row: 0, col: 0, size: 1 } })
      const chart2 = createMockChart('2', { position: { row: 0, col: 1, size: 1 } })
      mockGetCharts.mockResolvedValue([chart1, chart2])
      mockReorderCharts.mockResolvedValue(undefined)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(2))

      await act(async () => {
        await result.current.reorderCharts(['2', '1'])
      })

      expect(mockReorderCharts).toHaveBeenCalledWith([
        { id: '2', position: { row: 0, col: 0, size: 1 } },
        { id: '1', position: { row: 0, col: 1, size: 1 } },
      ])
    })

    it('optimistically reorders charts', async () => {
      const chart1 = createMockChart('1')
      const chart2 = createMockChart('2')
      mockGetCharts.mockResolvedValue([chart1, chart2])

      // Delay the API response
      let resolveReorder: () => void
      mockReorderCharts.mockImplementation(
        () =>
          new Promise(resolve => {
            resolveReorder = resolve
          })
      )

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(2))

      act(() => {
        result.current.reorderCharts(['2', '1'])
      })

      // Optimistic reorder should be visible immediately
      expect(result.current.charts[0].id).toBe('2')
      expect(result.current.charts[1].id).toBe('1')

      // Resolve the promise
      await act(async () => {
        resolveReorder!()
      })
    })

    it('handles charts not in order array by appending them', async () => {
      const chart1 = createMockChart('1')
      const chart2 = createMockChart('2')
      const chart3 = createMockChart('3')
      mockGetCharts.mockResolvedValue([chart1, chart2, chart3])
      mockReorderCharts.mockResolvedValue(undefined)

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(3))

      await act(async () => {
        // Only specify 2 of 3 charts
        await result.current.reorderCharts(['2', '1'])
      })

      // Chart 3 should be appended at the end
      expect(result.current.charts[0].id).toBe('2')
      expect(result.current.charts[1].id).toBe('1')
      expect(result.current.charts[2].id).toBe('3')
    })

    it('sets error when reorder fails', async () => {
      const chart1 = createMockChart('1')
      const chart2 = createMockChart('2')
      mockGetCharts.mockResolvedValue([chart1, chart2])
      mockReorderCharts.mockRejectedValue(new Error('Reorder failed'))

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.charts).toHaveLength(2))

      await act(async () => {
        try {
          await result.current.reorderCharts(['2', '1'])
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Reorder failed')
    })
  })

  describe('return type structure', () => {
    it('has correct return type structure', async () => {
      mockGetCharts.mockResolvedValue([])

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.isLoading).toBe(false))

      expect(result.current).toHaveProperty('charts')
      expect(result.current).toHaveProperty('isLoading')
      expect(result.current).toHaveProperty('error')
      expect(result.current).toHaveProperty('pinChart')
      expect(result.current).toHaveProperty('updateChart')
      expect(result.current).toHaveProperty('deleteChart')
      expect(result.current).toHaveProperty('refreshChart')
      expect(result.current).toHaveProperty('reorderCharts')
    })

    it('provides functions for CRUD operations', async () => {
      mockGetCharts.mockResolvedValue([])

      const { result } = renderHook(() => useDashboard())

      await waitFor(() => expect(result.current.isLoading).toBe(false))

      expect(typeof result.current.pinChart).toBe('function')
      expect(typeof result.current.updateChart).toBe('function')
      expect(typeof result.current.deleteChart).toBe('function')
      expect(typeof result.current.refreshChart).toBe('function')
      expect(typeof result.current.reorderCharts).toBe('function')
    })
  })

  describe('cleanup', () => {
    it('does not update state after unmount', async () => {
      let resolveGetCharts: (value: PinnedChart[]) => void
      mockGetCharts.mockImplementation(
        () =>
          new Promise(resolve => {
            resolveGetCharts = resolve
          })
      )

      const { result, unmount } = renderHook(() => useDashboard())

      // Unmount before promise resolves
      unmount()

      // Resolve the promise after unmount
      await act(async () => {
        resolveGetCharts!([createMockChart('1')])
      })

      // Should not throw or update state
      expect(result.current.charts).toEqual([])
    })
  })
})
