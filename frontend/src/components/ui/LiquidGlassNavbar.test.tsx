import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { LiquidGlassNavbar } from './LiquidGlassNavbar'

describe('LiquidGlassNavbar', () => {
  it('renders with glass styling', () => {
    render(<LiquidGlassNavbar>Content</LiquidGlassNavbar>)
    expect(screen.getByText('Content')).toBeInTheDocument()
  })

  it('applies sticky positioning by default', () => {
    const { container } = render(<LiquidGlassNavbar>Content</LiquidGlassNavbar>)
    expect(container.firstChild).toHaveClass('sticky top-0')
  })

  it('renders left and right sections', () => {
    render(
      <LiquidGlassNavbar
        left={<span>Logo</span>}
        right={<span>Actions</span>}
      />
    )
    expect(screen.getByText('Logo')).toBeInTheDocument()
    expect(screen.getByText('Actions')).toBeInTheDocument()
  })
})
