import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { LoadingSpinner } from './LoadingSpinner'

describe('LoadingSpinner', () => {
  it('renders with default props', () => {
    const { container } = render(<LoadingSpinner />)

    // Check for spinner element
    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()

    // Check for loading text (appears twice: visible + sr-only)
    const loadingTexts = screen.getAllByText('Loading...')
    expect(loadingTexts).toHaveLength(2)
  })

  it('renders with custom label', () => {
    render(<LoadingSpinner label="Custom loading text" />)

    // Custom label appears in both visible and sr-only spans
    const customTexts = screen.getAllByText('Custom loading text')
    expect(customTexts).toHaveLength(2)
  })

  it('renders small size', () => {
    const { container } = render(<LoadingSpinner size="sm" />)

    const spinner = container.querySelector('.w-4.h-4')
    expect(spinner).toBeInTheDocument()
  })

  it('renders large size', () => {
    const { container } = render(<LoadingSpinner size="lg" />)

    const spinner = container.querySelector('.w-12.h-12')
    expect(spinner).toBeInTheDocument()
  })

  it('has correct accessibility attributes', () => {
    const { container } = render(<LoadingSpinner />)

    const wrapper = container.firstChild as HTMLElement
    expect(wrapper).toHaveAttribute('role', 'status')
    expect(wrapper).toHaveAttribute('aria-live', 'polite')
    expect(wrapper).toHaveAttribute('aria-busy', 'true')
  })

  it('applies custom className', () => {
    const { container } = render(<LoadingSpinner className="custom-class" />)

    const wrapper = container.firstChild as HTMLElement
    expect(wrapper).toHaveClass('custom-class')
  })

  it('renders screen reader text', () => {
    render(<LoadingSpinner label="Loading data..." />)

    // Check for sr-only class on one of the text elements
    const allTexts = screen.getAllByText('Loading data...')
    const srOnlyElement = allTexts.find((el) => el.classList.contains('sr-only'))
    expect(srOnlyElement).toBeInTheDocument()
  })

  it('renders medium size by default', () => {
    const { container } = render(<LoadingSpinner />)

    const spinner = container.querySelector('.w-8.h-8')
    expect(spinner).toBeInTheDocument()
  })

  it('renders spinner with correct border styling', () => {
    const { container } = render(<LoadingSpinner />)

    const spinner = container.querySelector('.animate-spin')
    expect(spinner).toHaveClass('rounded-full')
    expect(spinner).toHaveClass('border-slate-200')
    expect(spinner).toHaveClass('border-t-blue-600')
  })

  it('has normal animation direction for RTL support', () => {
    const { container } = render(<LoadingSpinner />)

    const spinner = container.querySelector('.animate-spin') as HTMLElement
    expect(spinner?.style.animationDirection).toBe('normal')
  })
})
