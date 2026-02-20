import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useAlerts, NewAlert } from './useAlerts'
import { AlertRule, Notification } from '../services/api'

// Mock data
const mockAlerts: AlertRule[] = [
  {
    id: 'alert-1',
    userId: 'user-1',
    name: 'High Revenue Alert',
    description: 'Alert when revenue exceeds threshold',
    metricId: 'metric-1',
    metricName: 'Daily Revenue',
    operator: 'gt',
    threshold: 10000,
    checkInterval: 300,
    channels: ['email', 'in_app'],
    locale: 'en',
    cooldownPeriod: 3600,
    lastTriggeredAt: null,
    lastValue: null,
    isActive: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'alert-2',
    userId: 'user-1',
    name: 'Low Stock Alert',
    description: 'Alert when stock is low',
    metricId: 'metric-2',
    metricName: 'Inventory Level',
    operator: 'lt',
    threshold: 100,
    checkInterval: 600,
    channels: ['in_app'],
    locale: 'en',
    cooldownPeriod: 1800,
    lastTriggeredAt: '2024-01-15T10:30:00Z',
    lastValue: 50,
    isActive: false,
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-15T10:30:00Z',
  },
]

const mockNotifications: Notification[] = [
  {
    id: 'notif-1',
    alertRuleId: 'alert-1',
    userId: 'user-1',
    type: 'in_app',
    status: 'delivered',
    content: {
      title: 'High Revenue Alert Triggered',
      message: 'Revenue has exceeded $10,000',
      actionUrl: '/dashboard/revenue',
    },
    locale: 'en',
    metricValue: 12500,
    threshold: 10000,
    errorMessage: null,
    sentAt: '2024-01-15T10:00:00Z',
    deliveredAt: '2024-01-15T10:00:01Z',
    readAt: null,
    createdAt: '2024-01-15T10:00:00Z',
  },
  {
    id: 'notif-2',
    alertRuleId: 'alert-2',
    userId: 'user-1',
    type: 'in_app',
    status: 'delivered',
    content: {
      title: 'Low Stock Alert',
      message: 'Inventory level is below 100 units',
    },
    locale: 'en',
    metricValue: 50,
    threshold: 100,
    errorMessage: null,
    sentAt: '2024-01-14T08:00:00Z',
    deliveredAt: '2024-01-14T08:00:01Z',
    readAt: '2024-01-14T09:00:00Z',
    createdAt: '2024-01-14T08:00:00Z',
  },
]

// Mock the alerts API
const mockGetRules = vi.fn().mockResolvedValue(mockAlerts)
const mockCreateRule = vi.fn()
const mockUpdateRule = vi.fn()
const mockDeleteRule = vi.fn()
const mockToggleRule = vi.fn()
const mockGetNotifications = vi.fn().mockResolvedValue(mockNotifications)
const mockMarkNotificationRead = vi.fn()
const mockMarkAllNotificationsRead = vi.fn()

vi.mock('../services/api', () => ({
  alertsApi: {
    getRules: () => mockGetRules(),
    createRule: (rule: Partial<AlertRule>) => mockCreateRule(rule),
    updateRule: (id: string, rule: Partial<AlertRule>) => mockUpdateRule(id, rule),
    deleteRule: (id: string) => mockDeleteRule(id),
    toggleRule: (id: string, isActive: boolean) => mockToggleRule(id, isActive),
    getNotifications: (unreadOnly?: boolean) => mockGetNotifications(unreadOnly),
    markNotificationRead: (id: string) => mockMarkNotificationRead(id),
    markAllNotificationsRead: () => mockMarkAllNotificationsRead(),
  },
}))

describe('useAlerts', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset mock implementations to default values
    mockGetRules.mockResolvedValue(mockAlerts)
    mockGetNotifications.mockResolvedValue(mockNotifications)
    mockCreateRule.mockResolvedValue({ id: 'new-alert', ...mockAlerts[0] })
    mockUpdateRule.mockResolvedValue(mockAlerts[0])
    mockDeleteRule.mockResolvedValue(undefined)
    mockToggleRule.mockResolvedValue(mockAlerts[0])
    mockMarkNotificationRead.mockResolvedValue(undefined)
    mockMarkAllNotificationsRead.mockResolvedValue(undefined)
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('returns empty alerts initially', () => {
      const { result } = renderHook(() => useAlerts())

      expect(result.current.alerts).toEqual([])
    })

    it('returns empty notifications initially', () => {
      const { result } = renderHook(() => useAlerts())

      expect(result.current.notifications).toEqual([])
    })

    it('returns isLoading as true initially', () => {
      const { result } = renderHook(() => useAlerts())

      expect(result.current.isLoading).toBe(true)
    })

    it('returns null error initially', () => {
      const { result } = renderHook(() => useAlerts())

      expect(result.current.error).toBeNull()
    })
  })

  describe('data loading', () => {
    it('loads alerts on mount', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(mockGetRules).toHaveBeenCalled()
      })

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })
    })

    it('loads notifications on mount', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(mockGetNotifications).toHaveBeenCalled()
      })

      await waitFor(() => {
        expect(result.current.notifications).toEqual(mockNotifications)
      })
    })

    it('sets isLoading to false after loading', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })
    })

    it('calculates unreadCount correctly', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.notifications).toEqual(mockNotifications)
      })

      // One notification is unread (notif-1)
      expect(result.current.unreadCount).toBe(1)
    })

    it('sets error when loading fails', async () => {
      mockGetRules.mockRejectedValueOnce(new Error('Failed to load'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.error).toBe('Failed to load')
      })
    })
  })

  describe('createAlert', () => {
    it('creates an alert successfully', async () => {
      const newAlert: NewAlert = {
        name: 'New Alert',
        metricId: 'metric-3',
        operator: 'gte',
        threshold: 5000,
        checkInterval: 300,
        channels: ['email'],
      }

      const createdAlert: AlertRule = {
        id: 'alert-3',
        userId: 'user-1',
        name: 'New Alert',
        description: null,
        metricId: 'metric-3',
        metricName: 'Test Metric',
        operator: 'gte',
        threshold: 5000,
        checkInterval: 300,
        channels: ['email'],
        locale: 'en',
        cooldownPeriod: 3600,
        lastTriggeredAt: null,
        lastValue: null,
        isActive: true,
        createdAt: '2024-01-20T00:00:00Z',
        updatedAt: '2024-01-20T00:00:00Z',
      }

      mockCreateRule.mockResolvedValueOnce(createdAlert)

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      await act(async () => {
        await result.current.createAlert(newAlert)
      })

      expect(mockCreateRule).toHaveBeenCalledWith(newAlert)
      expect(result.current.alerts).toHaveLength(3)
      expect(result.current.alerts).toContainEqual(createdAlert)
    })

    it('sets error when create fails', async () => {
      mockCreateRule.mockRejectedValueOnce(new Error('Create failed'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      await act(async () => {
        try {
          await result.current.createAlert({
            name: 'Test',
            metricId: 'metric-1',
            operator: 'gt',
            threshold: 100,
            checkInterval: 300,
            channels: ['email'],
          })
        } catch {
          // Expected to throw
        }
      })

      expect(result.current.error).toBe('Create failed')
    })
  })

  describe('updateAlert', () => {
    it('updates an alert optimistically', async () => {
      mockUpdateRule.mockImplementation(async (id: string, updates: Partial<AlertRule>) => ({
        ...mockAlerts[0],
        ...updates,
      }))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        await result.current.updateAlert('alert-1', { threshold: 15000 })
      })

      expect(mockUpdateRule).toHaveBeenCalledWith('alert-1', { threshold: 15000 })
      expect(result.current.alerts[0].threshold).toBe(15000)
    })

    it('reverts optimistic update on error', async () => {
      mockUpdateRule.mockRejectedValueOnce(new Error('Update failed'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        try {
          await result.current.updateAlert('alert-1', { threshold: 15000 })
        } catch {
          // Expected to throw
        }
      })

      // Should revert to original value
      expect(result.current.alerts[0].threshold).toBe(10000)
      expect(result.current.error).toBe('Update failed')
    })
  })

  describe('deleteAlert', () => {
    it('deletes an alert optimistically', async () => {
      mockDeleteRule.mockResolvedValueOnce(undefined)

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        await result.current.deleteAlert('alert-1')
      })

      expect(mockDeleteRule).toHaveBeenCalledWith('alert-1')
      expect(result.current.alerts).toHaveLength(1)
      expect(result.current.alerts.find(a => a.id === 'alert-1')).toBeUndefined()
    })

    it('reverts optimistic delete on error', async () => {
      mockDeleteRule.mockRejectedValueOnce(new Error('Delete failed'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        try {
          await result.current.deleteAlert('alert-1')
        } catch {
          // Expected to throw
        }
      })

      // Should revert to original state
      expect(result.current.alerts).toHaveLength(2)
      expect(result.current.alerts.find(a => a.id === 'alert-1')).toBeDefined()
      expect(result.current.error).toBe('Delete failed')
    })
  })

  describe('toggleAlert', () => {
    it('toggles alert enabled state', async () => {
      mockToggleRule.mockResolvedValueOnce(undefined)

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        await result.current.toggleAlert('alert-1', false)
      })

      expect(mockToggleRule).toHaveBeenCalledWith('alert-1', false)
      expect(result.current.alerts[0].isActive).toBe(false)
    })

    it('reverts toggle on error', async () => {
      mockToggleRule.mockRejectedValueOnce(new Error('Toggle failed'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      await act(async () => {
        try {
          await result.current.toggleAlert('alert-1', false)
        } catch {
          // Expected to throw
        }
      })

      // Should revert to original state
      expect(result.current.alerts[0].isActive).toBe(true)
      expect(result.current.error).toBe('Toggle failed')
    })
  })

  describe('markNotificationRead', () => {
    it('marks notification as read optimistically', async () => {
      mockMarkNotificationRead.mockResolvedValueOnce(undefined)

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.notifications).toEqual(mockNotifications)
      })

      expect(result.current.unreadCount).toBe(1)

      await act(async () => {
        await result.current.markNotificationRead('notif-1')
      })

      expect(mockMarkNotificationRead).toHaveBeenCalledWith('notif-1')
      expect(result.current.unreadCount).toBe(0)
    })

    it('reverts on error', async () => {
      mockMarkNotificationRead.mockRejectedValueOnce(new Error('Mark failed'))

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.notifications).toEqual(mockNotifications)
      })

      await act(async () => {
        try {
          await result.current.markNotificationRead('notif-1')
        } catch {
          // Expected to throw
        }
      })

      // Should revert
      expect(result.current.unreadCount).toBe(1)
      expect(result.current.error).toBe('Mark failed')
    })
  })

  describe('markAllNotificationsRead', () => {
    it('marks all notifications as read', async () => {
      mockMarkAllNotificationsRead.mockResolvedValueOnce(undefined)

      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.notifications).toEqual(mockNotifications)
      })

      expect(result.current.unreadCount).toBe(1)

      await act(async () => {
        await result.current.markAllNotificationsRead()
      })

      expect(mockMarkAllNotificationsRead).toHaveBeenCalled()
      expect(result.current.unreadCount).toBe(0)
    })
  })

  describe('refresh', () => {
    it('refreshes all data', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.alerts).toEqual(mockAlerts)
      })

      // Clear mocks to verify refresh calls
      mockGetRules.mockClear()
      mockGetNotifications.mockClear()

      await act(async () => {
        await result.current.refresh()
      })

      expect(mockGetRules).toHaveBeenCalled()
      expect(mockGetNotifications).toHaveBeenCalled()
    })
  })

  describe('notification polling', () => {
    it('sets up polling interval when enabled', async () => {
      const setIntervalSpy = vi.spyOn(global, 'setInterval')
      const clearIntervalSpy = vi.spyOn(global, 'clearInterval')

      const { unmount } = renderHook(() => useAlerts(true))

      // Wait for initial load
      await waitFor(() => {
        expect(mockGetNotifications).toHaveBeenCalled()
      })

      // Should have set up an interval
      expect(setIntervalSpy).toHaveBeenCalled()

      unmount()

      // Should clear interval on unmount
      expect(clearIntervalSpy).toHaveBeenCalled()

      setIntervalSpy.mockRestore()
      clearIntervalSpy.mockRestore()
    })

    it('does not set up polling interval when disabled', async () => {
      // Track intervals with our specific duration
      const pollingIntervals: number[] = []
      const originalSetInterval = global.setInterval

      global.setInterval = vi.fn((callback: TimerHandler, delay?: number) => {
        // Track intervals with the polling duration (30000ms)
        if (typeof delay === 'number' && delay === 30000) {
          pollingIntervals.push(delay)
        }
        return originalSetInterval.call(global, callback, delay)
      })

      const { unmount } = renderHook(() => useAlerts(false))

      // Wait for initial load
      await waitFor(() => {
        expect(mockGetNotifications).toHaveBeenCalled()
      })

      // Should not have set up a polling interval (30000ms)
      expect(pollingIntervals).toHaveLength(0)

      unmount()

      global.setInterval = originalSetInterval
    })

    it('cleans up polling on unmount', async () => {
      const clearIntervalSpy = vi.spyOn(global, 'clearInterval')

      const { unmount } = renderHook(() => useAlerts(true))

      // Wait for initial load
      await waitFor(() => {
        expect(mockGetNotifications).toHaveBeenCalled()
      })

      unmount()

      // Should clear interval on unmount
      expect(clearIntervalSpy).toHaveBeenCalled()

      clearIntervalSpy.mockRestore()
    })
  })

  describe('return type', () => {
    it('has correct return type structure', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      expect(result.current).toHaveProperty('alerts')
      expect(result.current).toHaveProperty('notifications')
      expect(result.current).toHaveProperty('isLoading')
      expect(result.current).toHaveProperty('error')
      expect(result.current).toHaveProperty('createAlert')
      expect(result.current).toHaveProperty('updateAlert')
      expect(result.current).toHaveProperty('deleteAlert')
      expect(result.current).toHaveProperty('toggleAlert')
      expect(result.current).toHaveProperty('refresh')
      expect(result.current).toHaveProperty('markNotificationRead')
      expect(result.current).toHaveProperty('markAllNotificationsRead')
      expect(result.current).toHaveProperty('unreadCount')
    })

    it('provides functions with correct types', async () => {
      const { result } = renderHook(() => useAlerts())

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false)
      })

      expect(typeof result.current.createAlert).toBe('function')
      expect(typeof result.current.updateAlert).toBe('function')
      expect(typeof result.current.deleteAlert).toBe('function')
      expect(typeof result.current.toggleAlert).toBe('function')
      expect(typeof result.current.refresh).toBe('function')
      expect(typeof result.current.markNotificationRead).toBe('function')
      expect(typeof result.current.markAllNotificationsRead).toBe('function')
    })
  })
})
