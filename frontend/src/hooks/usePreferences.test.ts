/**
 * usePreferences Hook Tests
 *
 * Tests for user preferences management with API integration.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'

// i18n is mocked globally in src/test/setup.ts

// Mock the API module
vi.mock('../services/api', () => ({
  preferencesApi: {
    get: vi.fn(),
    update: vi.fn(),
  },
  APIError: class APIError extends Error {
    status: number
    statusText: string
    data: unknown
    constructor(status: number, statusText: string, message: string, data?: unknown) {
      super(message)
      this.status = status
      this.statusText = statusText
      this.data = data
    }
  },
}))

import { usePreferences } from './usePreferences'
import { preferencesApi, APIError } from '../services/api'

const mockPrefs = {
  id: 'pref-1',
  userId: 'user-1',
  locale: 'en' as const,
  numeralSystem: 'western' as const,
  calendarSystem: 'gregorian' as const,
  reportLanguage: 'en' as const,
  timezone: 'America/New_York',
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-01-01T00:00:00Z',
}

describe('usePreferences', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('loads preferences successfully', async () => {
    vi.mocked(preferencesApi.get).mockResolvedValue(mockPrefs)

    const { result } = renderHook(() => usePreferences())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    }, { timeout: 15000 })

    expect(result.current.preferences).toEqual(mockPrefs)
    expect(result.current.error).toBeNull()
  }, 20000)

  it('uses defaults on 404', async () => {
    vi.mocked(preferencesApi.get).mockRejectedValue(
      new APIError(404, 'Not Found', 'Not found')
    )

    const { result } = renderHook(() => usePreferences())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    }, { timeout: 15000 })

    expect(result.current.preferences).not.toBeNull()
    expect(result.current.preferences?.locale).toBe('en')
  }, 20000)

  it('sets error on other API failures', async () => {
    vi.mocked(preferencesApi.get).mockRejectedValue(
      new APIError(500, 'Error', 'Server error')
    )

    const { result } = renderHook(() => usePreferences())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    }, { timeout: 15000 })

    expect(result.current.error).toBe('Server error')
  }, 20000)

  it('updates preferences successfully', async () => {
    // Mock both APIs before rendering
    vi.mocked(preferencesApi.get).mockResolvedValue(mockPrefs)
    const updated = { ...mockPrefs, locale: 'ar' as const }
    vi.mocked(preferencesApi.update).mockResolvedValue(updated)

    const { result } = renderHook(() => usePreferences())

    // Wait for initial load
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    }, { timeout: 15000 })

    // Ensure preferences are loaded
    expect(result.current.preferences).not.toBeNull()

    // Update preferences
    await act(async () => {
      await result.current.updatePreferences({ locale: 'ar' })
    })

    expect(preferencesApi.update).toHaveBeenCalledWith({ locale: 'ar' })
    expect(result.current.preferences?.locale).toBe('ar')
  }, 20000)

  it('resets to default values', async () => {
    vi.mocked(preferencesApi.get).mockResolvedValue(mockPrefs)
    vi.mocked(preferencesApi.update).mockResolvedValue(mockPrefs)

    const { result } = renderHook(() => usePreferences())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    }, { timeout: 15000 })

    await act(async () => {
      await result.current.resetToDefaults()
    })

    expect(preferencesApi.update).toHaveBeenCalledWith(
      expect.objectContaining({
        locale: 'en',
        numeralSystem: 'western',
        calendarSystem: 'gregorian',
      })
    )
  }, 20000)
})
