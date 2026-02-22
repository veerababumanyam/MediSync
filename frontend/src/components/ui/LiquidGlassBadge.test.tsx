import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { LiquidGlassBadge } from './LiquidGlassBadge'

describe('LiquidGlassBadge', () => {
  it('renders children', () => {
    render(<LiquidGlassBadge>Badge</LiquidGlassBadge>)
    expect(screen.getByText('Badge')).toBeInTheDocument()
  })

  it('applies variant styling', () => {
    const { container } = render(<LiquidGlassBadge variant="success">Success</LiquidGlassBadge>)
    expect(container.firstChild).toHaveClass('liquid-glass-badge-green')
  })

  it('renders with icon', () => {
    render(<LiquidGlassBadge icon={<span>Icon</span>}>With Icon</LiquidGlassBadge>)
    expect(screen.getByText('Icon')).toBeInTheDocument()
  })
})
