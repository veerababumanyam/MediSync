import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { LiquidGlassToast } from './LiquidGlassToast'

describe('LiquidGlassToast', () => {
  it('renders with message', () => {
    render(<LiquidGlassToast message="Test toast" type="info" onClose={() => {}} />)
    expect(screen.getByText('Test toast')).toBeInTheDocument()
  })

  it('calls onClose when close button clicked', () => {
    const onClose = vi.fn()
    render(<LiquidGlassToast message="Test" type="success" onClose={onClose} />)
    fireEvent.click(screen.getByRole('button', { name: /close/i }))
    expect(onClose).toHaveBeenCalled()
  })

  it('applies type-specific styling', () => {
    const { container } = render(<LiquidGlassToast message="Error" type="error" onClose={() => {}} />)
    expect(container.firstChild).toHaveClass('border-red-400/30')
  })
})
