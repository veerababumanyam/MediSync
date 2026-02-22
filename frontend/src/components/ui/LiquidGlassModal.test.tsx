import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { LiquidGlassModal } from './LiquidGlassModal'

describe('LiquidGlassModal', () => {
  it('renders when open', () => {
    render(
      <LiquidGlassModal isOpen onClose={() => { }}>
        Modal Content
      </LiquidGlassModal>
    )
    expect(screen.getByText('Modal Content')).toBeInTheDocument()
  })

  it('does not render when closed', () => {
    render(
      <LiquidGlassModal isOpen={false} onClose={() => { }}>
        Modal Content
      </LiquidGlassModal>
    )
    expect(screen.queryByText('Modal Content')).not.toBeInTheDocument()
  })

  it('calls onClose when backdrop is clicked', () => {
    const onClose = vi.fn()
    render(
      <LiquidGlassModal isOpen onClose={onClose}>
        Modal Content
      </LiquidGlassModal>
    )
    fireEvent.click(screen.getByRole('presentation').firstChild!)
    expect(onClose).toHaveBeenCalled()
  })

  it('renders title and actions', () => {
    render(
      // eslint-disable-next-line no-restricted-syntax
      <LiquidGlassModal isOpen onClose={() => { }} title="Test Title" actions={<button>Action</button>}>
        Content
      </LiquidGlassModal>
    )
    expect(screen.getByText('Test Title')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Action' })).toBeInTheDocument()
  })
})
