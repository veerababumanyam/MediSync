import '@testing-library/jest-dom'
import { vi } from 'vitest'
import { createI18nMock } from './i18n-mock'

// Global i18n mock with proper interpolation support
vi.mock('react-i18next', () => createI18nMock())

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock ResizeObserver
;(globalThis as unknown as { ResizeObserver: unknown }).ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock IntersectionObserver
;(globalThis as unknown as { IntersectionObserver: unknown }).IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock document.elementFromPoint for tests
document.elementFromPoint = vi.fn().mockReturnValue(document.body)

// Mock getBoundingClientRect for tests
Element.prototype.getBoundingClientRect = vi.fn().mockReturnValue({
  left: 0,
  top: 0,
  width: 100,
  height: 100,
  right: 100,
  bottom: 100,
  x: 0,
  y: 0,
  toJSON: () => ({}),
})
