import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { LiquidGlassSidebar } from './LiquidGlassSidebar'

describe('LiquidGlassSidebar', () => {
  it('renders with glass styling', () => {
    const { container } = render(<LiquidGlassSidebar>Content</LiquidGlassSidebar>)
    expect(container.firstChild).toHaveClass('liquid-glass')
  })

  it('applies collapsed state', () => {
    const { container } = render(<LiquidGlassSidebar collapsed>Content</LiquidGlassSidebar>)
    expect(container.firstChild).toHaveClass('w-16')
  })

  it('applies expanded width by default', () => {
    const { container } = render(<LiquidGlassSidebar>Content</LiquidGlassSidebar>)
    expect(container.firstChild).toHaveClass('w-64')
  })

  it('renders children content', () => {
    render(<LiquidGlassSidebar><span>Sidebar Item</span></LiquidGlassSidebar>)
    expect(screen.getByText('Sidebar Item')).toBeInTheDocument()
  })

  it('applies custom className', () => {
    const { container } = render(<LiquidGlassSidebar className="custom-class">Content</LiquidGlassSidebar>)
    expect(container.firstChild).toHaveClass('custom-class')
  })

  it('applies transition classes', () => {
    const { container } = render(<LiquidGlassSidebar>Content</LiquidGlassSidebar>)
    expect(container.firstChild).toHaveClass('transition-all duration-300')
  })
})
