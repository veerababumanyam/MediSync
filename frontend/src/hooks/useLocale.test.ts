import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useLocale } from './useLocale'

// Mock the preferences API
const mockPreferencesGet = vi.fn().mockResolvedValue({ locale: 'en' })
const mockPreferencesUpdate = vi.fn().mockResolvedValue(undefined)

vi.mock('../services/api', () => ({
  preferencesApi: {
    get: () => mockPreferencesGet(),
    update: (data: { locale: string }) => mockPreferencesUpdate(data),
  },
}))

// Mock i18next
const mockChangeLanguage = vi.fn().mockResolvedValue(undefined)

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, defaultValue?: string) => defaultValue || key,
    i18n: {
      language: 'en',
      changeLanguage: (lang: string) => mockChangeLanguage(lang),
    },
  }),
}))

describe('useLocale', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    // Reset document
    document.documentElement.dir = 'ltr'
    document.documentElement.lang = 'en'
    document.body.classList.remove('rtl', 'ltr')
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('returns default locale as English', () => {
    const { result } = renderHook(() => useLocale())

    expect(result.current.locale).toBe('en')
    expect(result.current.isRTL).toBe(false)
  })

  it('correctly identifies isRTL as false for English', () => {
    const { result } = renderHook(() => useLocale())

    expect(result.current.isRTL).toBe(false)
  })

  it('provides setLocale function', () => {
    const { result } = renderHook(() => useLocale())

    expect(typeof result.current.setLocale).toBe('function')
  })

  it('provides toggleLocale function', () => {
    const { result } = renderHook(() => useLocale())

    expect(typeof result.current.toggleLocale).toBe('function')
  })

  it('provides isLoading state', () => {
    const { result } = renderHook(() => useLocale())

    expect(typeof result.current.isLoading).toBe('boolean')
  })

  it('provides error state', () => {
    const { result } = renderHook(() => useLocale())

    expect(result.current.error).toBeNull()
  })

  it('has correct return type structure', () => {
    const { result } = renderHook(() => useLocale())

    expect(result.current).toHaveProperty('locale')
    expect(result.current).toHaveProperty('isRTL')
    expect(result.current).toHaveProperty('setLocale')
    expect(result.current).toHaveProperty('toggleLocale')
    expect(result.current).toHaveProperty('isLoading')
    expect(result.current).toHaveProperty('error')
  })

  it('loads preferences on mount', async () => {
    renderHook(() => useLocale())

    await waitFor(() => {
      expect(mockPreferencesGet).toHaveBeenCalled()
    })
  })

  it('updates document direction to ltr for English', async () => {
    renderHook(() => useLocale())

    await waitFor(() => {
      expect(document.documentElement.dir).toBe('ltr')
      expect(document.documentElement.lang).toBe('en')
    })
  })

  it('adds ltr class to body for English', async () => {
    renderHook(() => useLocale())

    await waitFor(() => {
      expect(document.body.classList.contains('ltr')).toBe(true)
      expect(document.body.classList.contains('rtl')).toBe(false)
    })
  })

  it('calls setLocale when setLocale is invoked', async () => {
    const { result } = renderHook(() => useLocale())

    await act(async () => {
      await result.current.setLocale('ar')
    })

    expect(mockChangeLanguage).toHaveBeenCalledWith('ar')
  })

  it('calls toggleLocale when toggleLocale is invoked', async () => {
    const { result } = renderHook(() => useLocale())

    await act(async () => {
      await result.current.toggleLocale()
    })

    expect(mockChangeLanguage).toHaveBeenCalled()
  })

  it('calls preferencesApi.update when locale is changed', async () => {
    const { result } = renderHook(() => useLocale())

    await act(async () => {
      await result.current.setLocale('ar')
    })

    expect(mockPreferencesUpdate).toHaveBeenCalledWith({ locale: 'ar' })
  })

  it('handles API update failure gracefully', async () => {
    mockPreferencesUpdate.mockRejectedValueOnce(new Error('API Error'))

    const { result } = renderHook(() => useLocale())

    await act(async () => {
      await result.current.setLocale('ar')
    })

    // Should still change language despite API failure
    expect(mockChangeLanguage).toHaveBeenCalledWith('ar')
    // Error should not be set (API failure is handled gracefully)
    expect(result.current.error).toBeNull()
  })

  it('sets error when changeLanguage fails', async () => {
    mockChangeLanguage.mockRejectedValueOnce(new Error('Language change failed'))

    const { result } = renderHook(() => useLocale())

    await act(async () => {
      await result.current.setLocale('ar')
    })

    expect(result.current.error).toBe('Language change failed')
  })

  it('sets isLoading to true during locale change', async () => {
    let resolveChangeLanguage: () => void
    mockChangeLanguage.mockImplementation(() => new Promise<void>(resolve => {
      resolveChangeLanguage = resolve
    }))

    const { result } = renderHook(() => useLocale())

    act(() => {
      result.current.setLocale('ar')
    })

    // Should be loading
    expect(result.current.isLoading).toBe(true)

    // Resolve the promise
    await act(async () => {
      resolveChangeLanguage!()
    })

    // Should no longer be loading
    expect(result.current.isLoading).toBe(false)
  })

  it('handles preferences load failure gracefully', async () => {
    mockPreferencesGet.mockRejectedValueOnce(new Error('Preferences not found'))

    // Should not throw
    const { result } = renderHook(() => useLocale())

    await waitFor(() => {
      expect(mockPreferencesGet).toHaveBeenCalled()
    })

    // Hook should still work
    expect(result.current.locale).toBe('en')
  })
})
