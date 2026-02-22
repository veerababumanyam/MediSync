import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { AnimatedBackground } from './AnimatedBackground'

describe('AnimatedBackground', () => {
  it('renders with mesh gradient background', () => {
    const { container } = render(<AnimatedBackground />)
    expect(container.firstChild).toHaveClass('fixed inset-0')
  })

  it('applies custom className', () => {
    const { container } = render(<AnimatedBackground className="custom-class" />)
    expect(container.firstChild).toHaveClass('custom-class')
  })

  it('renders floating orbs', () => {
    const { container } = render(<AnimatedBackground />)
    const orbs = container.querySelectorAll('.animate-float')
    expect(orbs.length).toBe(3)
  })
})
