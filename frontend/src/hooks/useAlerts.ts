/**
 * useAlerts Hook
 *
 * React hook for managing alert rules and notifications with:
 * - State management for alert rules and notifications
 * - CRUD operations for alert rules
 * - Toggle alert enabled/disabled
 * - Notification polling (WebSocket placeholder)
 * - Loading and error states
 *
 * @module hooks/useAlerts
 */
import { useCallback, useEffect, useState, useRef } from 'react'
import { alertsApi, type AlertRule, type Notification } from '../services/api'

/**
 * Type for creating a new alert rule
 */
export interface NewAlert {
  name: string
  description?: string
  metricId: string
  operator: 'gt' | 'gte' | 'lt' | 'lte' | 'eq'
  threshold: number
  checkInterval: number
  channels: string[]
  locale?: 'en' | 'ar'
  cooldownPeriod?: number
}

/**
 * Return type for useAlerts hook
 */
export interface UseAlertsReturn {
  /** List of alert rules */
  alerts: AlertRule[]
  /** List of notifications */
  notifications: Notification[]
  /** Whether data is being loaded */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
  /** Create a new alert rule */
  createAlert: (alert: NewAlert) => Promise<void>
  /** Update an existing alert rule */
  updateAlert: (id: string, updates: Partial<AlertRule>) => Promise<void>
  /** Delete an alert rule */
  deleteAlert: (id: string) => Promise<void>
  /** Toggle alert enabled/disabled state */
  toggleAlert: (id: string, enabled: boolean) => Promise<void>
  /** Refresh alerts and notifications */
  refresh: () => Promise<void>
  /** Mark notification as read */
  markNotificationRead: (id: string) => Promise<void>
  /** Mark all notifications as read */
  markAllNotificationsRead: () => Promise<void>
  /** Unread notification count */
  unreadCount: number
}

/**
 * Default polling interval for notifications (in milliseconds)
 */
const NOTIFICATION_POLL_INTERVAL = 30000 // 30 seconds

/**
 * Hook for managing alert rules and notifications
 */
export function useAlerts(pollNotifications = true): UseAlertsReturn {
  const [alerts, setAlerts] = useState<AlertRule[]>([])
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Ref for polling interval cleanup
  const pollingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  // Ref to track if WebSocket is connected (placeholder for future implementation)
  const webSocketRef = useRef<WebSocket | null>(null)

  /**
   * Fetch alerts and notifications from API
   */
  const fetchData = useCallback(async () => {
    try {
      setError(null)
      const [alertsData, notificationsData] = await Promise.all([
        alertsApi.getRules(),
        alertsApi.getNotifications(),
      ])
      setAlerts(alertsData)
      setNotifications(notificationsData)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load alerts'
      setError(message)
      console.error('Failed to fetch alerts:', err)
    }
  }, [])

  /**
   * Initial data load
   */
  useEffect(() => {
    async function loadInitialData() {
      setIsLoading(true)
      await fetchData()
      setIsLoading(false)
    }

    loadInitialData()
  }, [fetchData])

  /**
   * Notification polling setup
   * This is a placeholder for WebSocket implementation
   */
  useEffect(() => {
    if (!pollNotifications) return

    // Start polling for notifications
    pollingIntervalRef.current = setInterval(async () => {
      try {
        const notificationsData = await alertsApi.getNotifications()
        setNotifications(notificationsData)
      } catch (err) {
        // Silent fail for polling - don't update error state
        console.debug('Notification poll failed:', err)
      }
    }, NOTIFICATION_POLL_INTERVAL)

    // Cleanup polling on unmount
    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current)
        pollingIntervalRef.current = null
      }
    }
  }, [pollNotifications])

  /**
   * WebSocket connection placeholder
   * TODO: Implement WebSocket for real-time notifications
   */
  useEffect(() => {
    // Placeholder for WebSocket implementation
    // In the future, this will connect to the notification WebSocket endpoint
    // For now, we use polling as a fallback

    return () => {
      // Cleanup WebSocket on unmount
      if (webSocketRef.current) {
        webSocketRef.current.close()
        webSocketRef.current = null
      }
    }
  }, [])

  /**
   * Create a new alert rule
   */
  const createAlert = useCallback(async (alert: NewAlert) => {
    setIsLoading(true)
    setError(null)

    try {
      const newAlert = await alertsApi.createRule(alert)
      setAlerts(prev => [...prev, newAlert])
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create alert'
      setError(message)
      console.error('Failed to create alert:', err)
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Update an existing alert rule
   */
  const updateAlert = useCallback(async (id: string, updates: Partial<AlertRule>) => {
    setError(null)

    // Optimistic update
    const previousAlerts = alerts
    setAlerts(prev =>
      prev.map(alert =>
        alert.id === id ? { ...alert, ...updates } : alert
      )
    )

    try {
      const updatedAlert = await alertsApi.updateRule(id, updates)
      setAlerts(prev =>
        prev.map(alert =>
          alert.id === id ? updatedAlert : alert
        )
      )
    } catch (err) {
      // Revert optimistic update on error
      setAlerts(previousAlerts)
      const message = err instanceof Error ? err.message : 'Failed to update alert'
      setError(message)
      console.error('Failed to update alert:', err)
      throw err
    }
  }, [alerts])

  /**
   * Delete an alert rule
   */
  const deleteAlert = useCallback(async (id: string) => {
    setError(null)

    // Optimistic delete
    const previousAlerts = alerts
    setAlerts(prev => prev.filter(alert => alert.id !== id))

    try {
      await alertsApi.deleteRule(id)
    } catch (err) {
      // Revert optimistic delete on error
      setAlerts(previousAlerts)
      const message = err instanceof Error ? err.message : 'Failed to delete alert'
      setError(message)
      console.error('Failed to delete alert:', err)
      throw err
    }
  }, [alerts])

  /**
   * Toggle alert enabled/disabled state
   */
  const toggleAlert = useCallback(async (id: string, enabled: boolean) => {
    setError(null)

    // Optimistic update
    const previousAlerts = alerts
    setAlerts(prev =>
      prev.map(alert =>
        alert.id === id ? { ...alert, isActive: enabled } : alert
      )
    )

    try {
      await alertsApi.toggleRule(id, enabled)
    } catch (err) {
      // Revert optimistic update on error
      setAlerts(previousAlerts)
      const message = err instanceof Error ? err.message : 'Failed to toggle alert'
      setError(message)
      console.error('Failed to toggle alert:', err)
      throw err
    }
  }, [alerts])

  /**
   * Refresh all data
   */
  const refresh = useCallback(async () => {
    setIsLoading(true)
    await fetchData()
    setIsLoading(false)
  }, [fetchData])

  /**
   * Mark a notification as read
   */
  const markNotificationRead = useCallback(async (id: string) => {
    // Optimistic update
    const previousNotifications = notifications
    setNotifications(prev =>
      prev.map(notification =>
        notification.id === id
          ? { ...notification, readAt: new Date().toISOString() }
          : notification
      )
    )

    try {
      await alertsApi.markNotificationRead(id)
    } catch (err) {
      // Revert on error
      setNotifications(previousNotifications)
      const message = err instanceof Error ? err.message : 'Failed to mark notification as read'
      setError(message)
      console.error('Failed to mark notification as read:', err)
      throw err
    }
  }, [notifications])

  /**
   * Mark all notifications as read
   */
  const markAllNotificationsRead = useCallback(async () => {
    // Optimistic update
    const previousNotifications = notifications
    const now = new Date().toISOString()
    setNotifications(prev =>
      prev.map(notification => ({ ...notification, readAt: now }))
    )

    try {
      await alertsApi.markAllNotificationsRead()
    } catch (err) {
      // Revert on error
      setNotifications(previousNotifications)
      const message = err instanceof Error ? err.message : 'Failed to mark all notifications as read'
      setError(message)
      console.error('Failed to mark all notifications as read:', err)
      throw err
    }
  }, [notifications])

  /**
   * Calculate unread notification count
   */
  const unreadCount = notifications.filter(n => !n.readAt).length

  return {
    alerts,
    notifications,
    isLoading,
    error,
    createAlert,
    updateAlert,
    deleteAlert,
    toggleAlert,
    refresh,
    markNotificationRead,
    markAllNotificationsRead,
    unreadCount,
  }
}

export default useAlerts
