# Liquid Glass Complete Redesign — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Apply the Liquid Glass design system to ALL user-facing pages, creating a cohesive, modern, premium experience.

**Architecture:** Build new UI components (Modal, Badge, Toast, Navbar, Sidebar, Table, etc.), update existing components, redesign pages with animated mesh backgrounds, and ensure WCAG 2.2 AA compliance throughout.

**Tech Stack:** React 19, TypeScript, Tailwind CSS, Framer Motion, class-variance-authority, Apache ECharts

---

## Phase 1: Foundation — Animated Background & Layout

### Task 1.1: Create AnimatedBackground Component

**Files:**
- Create: `frontend/src/components/ui/AnimatedBackground.tsx`
- Test: `frontend/src/components/ui/AnimatedBackground.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/AnimatedBackground.test.tsx
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
```

**Step 2: Run test to verify it fails**

Run: `cd frontend && npm test -- AnimatedBackground.test.tsx`
Expected: FAIL — component doesn't exist

**Step 3: Write implementation**

```tsx
// frontend/src/components/ui/AnimatedBackground.tsx
import React from 'react'
import { cn } from '@/lib/cn'

export interface AnimatedBackgroundProps {
  className?: string
  children?: React.ReactNode
}

export const AnimatedBackground: React.FC<AnimatedBackgroundProps> = ({
  className,
  children,
}) => {
  return (
    <div className={cn('fixed inset-0 -z-10 overflow-hidden', className)}>
      {/* Mesh gradient base */}
      <div
        className="absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse 80% 60% at 10% 20%, rgba(88, 86, 214, 0.4) 0%, transparent 60%),
            radial-gradient(ellipse 60% 80% at 80% 80%, rgba(0, 122, 255, 0.3) 0%, transparent 60%),
            radial-gradient(ellipse 50% 50% at 50% 50%, rgba(175, 82, 222, 0.15) 0%, transparent 50%),
            #0A0A1A
          `,
        }}
      />

      {/* Floating orbs */}
      <div
        className="absolute w-[500px] h-[500px] rounded-full animate-float opacity-35"
        style={{
          background: 'radial-gradient(circle, rgba(0, 122, 255, 0.35) 0%, transparent 70%)',
          filter: 'blur(80px)',
          top: '10%',
          left: '5%',
          animationDuration: '20s',
        }}
      />
      <div
        className="absolute w-[400px] h-[400px] rounded-full animate-float opacity-30"
        style={{
          background: 'radial-gradient(circle, rgba(175, 82, 222, 0.30) 0%, transparent 70%)',
          filter: 'blur(80px)',
          bottom: '20%',
          right: '10%',
          animationDuration: '25s',
          animationDelay: '-5s',
        }}
      />
      <div
        className="absolute w-[350px] h-[350px] rounded-full animate-float opacity-20"
        style={{
          background: 'radial-gradient(circle, rgba(255, 45, 85, 0.20) 0%, transparent 70%)',
          filter: 'blur(80px)',
          top: '50%',
          left: '50%',
          animationDuration: '18s',
          animationDelay: '-10s',
        }}
      />

      {children}
    </div>
  )
}

export default AnimatedBackground
```

**Step 4: Run test to verify it passes**

Run: `cd frontend && npm test -- AnimatedBackground.test.tsx`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/src/components/ui/AnimatedBackground.tsx frontend/src/components/ui/AnimatedBackground.test.tsx
git commit -m "feat(ui): add AnimatedBackground component with mesh gradient and floating orbs"
```

---

### Task 1.2: Add Float Animation to globals.css

**Files:**
- Modify: `frontend/src/styles/globals.css`

**Step 1: Add float animation keyframes**

Add to `frontend/src/styles/globals.css`:

```css
/* Floating animation for background orbs */
@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -30px) scale(1.05);
  }
  50% {
    transform: translate(-20px, 20px) scale(0.95);
  }
  75% {
    transform: translate(-30px, -20px) scale(1.02);
  }
}

.animate-float {
  animation: float 20s ease-in-out infinite;
}

/* Respect reduced motion */
@media (prefers-reduced-motion: reduce) {
  .animate-float {
    animation: none;
  }
}
```

**Step 2: Commit**

```bash
git add frontend/src/styles/globals.css
git commit -m "feat(styles): add float animation for background orbs"
```

---

### Task 1.3: Create LiquidGlassNavbar Component

**Files:**
- Create: `frontend/src/components/ui/LiquidGlassNavbar.tsx`
- Test: `frontend/src/components/ui/LiquidGlassNavbar.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/LiquidGlassNavbar.test.tsx
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
```

**Step 2: Run test to verify it fails**

Run: `cd frontend && npm test -- LiquidGlassNavbar.test.tsx`
Expected: FAIL — component doesn't exist

**Step 3: Write implementation**

```tsx
// frontend/src/components/ui/LiquidGlassNavbar.tsx
import React from 'react'
import { cn } from '@/lib/cn'

export interface LiquidGlassNavbarProps {
  children?: React.ReactNode
  left?: React.ReactNode
  center?: React.ReactNode
  right?: React.ReactNode
  className?: string
  sticky?: boolean
}

export const LiquidGlassNavbar: React.FC<LiquidGlassNavbarProps> = ({
  children,
  left,
  center,
  right,
  className,
  sticky = true,
}) => {
  return (
    <nav
      className={cn(
        'liquid-glass-header z-50 px-4 py-3 md:px-6',
        sticky && 'sticky top-0',
        className
      )}
    >
      <div className="flex items-center justify-between max-w-7xl mx-auto">
        {/* Left section */}
        {left && (
          <div className="flex items-center gap-4 flex-shrink-0">
            {left}
          </div>
        )}

        {/* Center section */}
        {center && (
          <div className="flex items-center justify-center flex-1">
            {center}
          </div>
        )}

        {/* Right section */}
        {right && (
          <div className="flex items-center gap-3 flex-shrink-0">
            {right}
          </div>
        )}

        {/* Default children */}
        {children}
      </div>
    </nav>
  )
}

export default LiquidGlassNavbar
```

**Step 4: Run test to verify it passes**

Run: `cd frontend && npm test -- LiquidGlassNavbar.test.tsx`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/src/components/ui/LiquidGlassNavbar.tsx frontend/src/components/ui/LiquidGlassNavbar.test.tsx
git commit -m "feat(ui): add LiquidGlassNavbar component with glass styling"
```

---

### Task 1.4: Create LiquidGlassModal Component

**Files:**
- Create: `frontend/src/components/ui/LiquidGlassModal.tsx`
- Test: `frontend/src/components/ui/LiquidGlassModal.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/LiquidGlassModal.test.tsx
import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { LiquidGlassModal } from './LiquidGlassModal'

describe('LiquidGlassModal', () => {
  it('renders when open', () => {
    render(
      <LiquidGlassModal isOpen onClose={() => {}}>
        Modal Content
      </LiquidGlassModal>
    )
    expect(screen.getByText('Modal Content')).toBeInTheDocument()
  })

  it('does not render when closed', () => {
    render(
      <LiquidGlassModal isOpen={false} onClose={() => {}}>
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
      <LiquidGlassModal isOpen onClose={() => {}} title="Test Title" actions={<button>Action</button>}>
        Content
      </LiquidGlassModal>
    )
    expect(screen.getByText('Test Title')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Action' })).toBeInTheDocument()
  })
})
```

**Step 2: Run test to verify it fails**

Run: `cd frontend && npm test -- LiquidGlassModal.test.tsx`
Expected: FAIL — component doesn't exist

**Step 3: Write implementation**

```tsx
// frontend/src/components/ui/LiquidGlassModal.tsx
import React, { useEffect, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/lib/cn'

export interface LiquidGlassModalProps {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
  title?: string
  actions?: React.ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  className?: string
  closeOnBackdrop?: boolean
  closeOnEscape?: boolean
}

export const LiquidGlassModal: React.FC<LiquidGlassModalProps> = ({
  isOpen,
  onClose,
  children,
  title,
  actions,
  size = 'md',
  className,
  closeOnBackdrop = true,
  closeOnEscape = true,
}) => {
  // Handle escape key
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape' && closeOnEscape) {
        onClose()
      }
    },
    [onClose, closeOnEscape]
  )

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown)
      document.body.style.overflow = 'hidden'
    }
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      document.body.style.overflow = ''
    }
  }, [isOpen, handleKeyDown])

  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    full: 'max-w-4xl',
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <div
          role="presentation"
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
        >
          {/* Backdrop */}
          <motion.div
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={closeOnBackdrop ? onClose : undefined}
          />

          {/* Modal */}
          <motion.div
            className={cn(
              'liquid-glass-modal relative w-full p-6',
              sizeClasses[size],
              className
            )}
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ duration: 0.2, ease: [0.4, 0, 0.2, 1] }}
          >
            {/* Title */}
            {title && (
              <h2 className="text-xl font-semibold mb-4 liquid-text-primary">
                {title}
              </h2>
            )}

            {/* Content */}
            <div className="mb-4">{children}</div>

            {/* Actions */}
            {actions && (
              <div className="flex justify-end gap-3 pt-4 border-t border-white/10">
                {actions}
              </div>
            )}
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  )
}

export default LiquidGlassModal
```

**Step 4: Run test to verify it passes**

Run: `cd frontend && npm test -- LiquidGlassModal.test.tsx`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/src/components/ui/LiquidGlassModal.tsx frontend/src/components/ui/LiquidGlassModal.test.tsx
git commit -m "feat(ui): add LiquidGlassModal component with backdrop blur and animations"
```

---

### Task 1.5: Create LiquidGlassBadge Component

**Files:**
- Create: `frontend/src/components/ui/LiquidGlassBadge.tsx`
- Test: `frontend/src/components/ui/LiquidGlassBadge.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/LiquidGlassBadge.test.tsx
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
```

**Step 2: Run test to verify it fails**

Run: `cd frontend && npm test -- LiquidGlassBadge.test.tsx`
Expected: FAIL — component doesn't exist

**Step 3: Write implementation**

```tsx
// frontend/src/components/ui/LiquidGlassBadge.tsx
import React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/cn'

const badgeVariants = cva(
  'liquid-glass-badge inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold',
  {
    variants: {
      variant: {
        default: '',
        blue: 'liquid-glass-badge-blue',
        teal: 'liquid-glass-badge-teal',
        green: 'liquid-glass-badge-green',
        red: 'liquid-glass-badge-red',
        success: 'liquid-glass-badge-green',
        warning: 'liquid-glass-badge-red',
        error: 'liquid-glass-badge-red',
      },
      size: {
        sm: 'text-[10px] px-2 py-0.5',
        md: 'text-xs px-2.5 py-1',
        lg: 'text-sm px-3 py-1.5',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  }
)

export interface LiquidGlassBadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {
  icon?: React.ReactNode
}

export const LiquidGlassBadge: React.FC<LiquidGlassBadgeProps> = ({
  className,
  variant,
  size,
  icon,
  children,
  ...props
}) => {
  return (
    <span className={cn(badgeVariants({ variant, size }), className)} {...props}>
      {icon}
      {children}
    </span>
  )
}

export default LiquidGlassBadge
```

**Step 4: Run test to verify it passes**

Run: `cd frontend && npm test -- LiquidGlassBadge.test.tsx`
Expected: PASS

**Step 5: Commit**

```bash
git add frontend/src/components/ui/LiquidGlassBadge.tsx frontend/src/components/ui/LiquidGlassBadge.test.tsx
git commit -m "feat(ui): add LiquidGlassBadge component with semantic variants"
```

---

### Task 1.6: Create LiquidGlassToast Component

**Files:**
- Create: `frontend/src/components/ui/LiquidGlassToast.tsx`
- Create: `frontend/src/hooks/useToast.ts`
- Test: `frontend/src/components/ui/LiquidGlassToast.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/LiquidGlassToast.test.tsx
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
```

**Step 2: Run test to verify it fails**

Run: `cd frontend && npm test -- LiquidGlassToast.test.tsx`
Expected: FAIL — component doesn't exist

**Step 3: Write implementation**

```tsx
// frontend/src/components/ui/LiquidGlassToast.tsx
import React from 'react'
import { motion } from 'framer-motion'
import { cn } from '@/lib/cn'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface LiquidGlassToastProps {
  message: string
  type: ToastType
  onClose: () => void
  duration?: number
  className?: string
}

const typeStyles: Record<ToastType, string> = {
  success: 'border-emerald-400/30 bg-emerald-500/10',
  error: 'border-red-400/30 bg-red-500/10',
  warning: 'border-amber-400/30 bg-amber-500/10',
  info: 'border-blue-400/30 bg-blue-500/10',
}

const typeIcons: Record<ToastType, React.ReactNode> = {
  success: (
    <svg className="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
    </svg>
  ),
  error: (
    <svg className="w-5 h-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
  ),
  warning: (
    <svg className="w-5 h-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  ),
  info: (
    <svg className="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  ),
}

export const LiquidGlassToast: React.FC<LiquidGlassToastProps> = ({
  message,
  type,
  onClose,
  className,
}) => {
  return (
    <motion.div
      className={cn(
        'liquid-glass flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg',
        typeStyles[type],
        className
      )}
      initial={{ opacity: 0, y: 50, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: 20, scale: 0.95 }}
      transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1] }}
    >
      {typeIcons[type]}
      <p className="flex-1 text-sm font-medium liquid-text-primary">{message}</p>
      <button
        onClick={onClose}
        className="p-1 hover:bg-white/10 rounded-lg transition-colors"
        aria-label="Close"
      >
        <svg className="w-4 h-4 liquid-text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </motion.div>
  )
}

export default LiquidGlassToast
```

**Step 4: Create useToast hook**

```tsx
// frontend/src/hooks/useToast.ts
import { useState, useCallback } from 'react'
import type { ToastType } from '@/components/ui/LiquidGlassToast'

export interface Toast {
  id: string
  message: string
  type: ToastType
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([])

  const addToast = useCallback((message: string, type: ToastType = 'info') => {
    const id = crypto.randomUUID()
    setToasts((prev) => [...prev, { id, message, type }])
    return id
  }, [])

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const success = useCallback((message: string) => addToast(message, 'success'), [addToast])
  const error = useCallback((message: string) => addToast(message, 'error'), [addToast])
  const warning = useCallback((message: string) => addToast(message, 'warning'), [addToast])
  const info = useCallback((message: string) => addToast(message, 'info'), [addToast])

  return { toasts, addToast, removeToast, success, error, warning, info }
}
```

**Step 5: Run test to verify it passes**

Run: `cd frontend && npm test -- LiquidGlassToast.test.tsx`
Expected: PASS

**Step 6: Commit**

```bash
git add frontend/src/components/ui/LiquidGlassToast.tsx frontend/src/hooks/useToast.ts frontend/src/components/ui/LiquidGlassToast.test.tsx
git commit -m "feat(ui): add LiquidGlassToast component and useToast hook"
```

---

## Phase 2: Chat Page Redesign

### Task 2.1: Update ChatInterface with Liquid Glass

**Files:**
- Modify: `frontend/src/components/chat/ChatInterface.tsx`

**Step 1: Read current implementation**

Run: Read the current `ChatInterface.tsx` file to understand structure

**Step 2: Update imports and add glass styling**

Add imports:
```tsx
import { LiquidGlassCard } from '@/components/ui/LiquidGlassCard'
import { AnimatedBackground } from '@/components/ui/AnimatedBackground'
```

Wrap the interface with:
```tsx
<LiquidGlassCard intensity="medium" className="flex flex-col h-full">
  {/* Existing content */}
</LiquidGlassCard>
```

**Step 3: Commit**

```bash
git add frontend/src/components/chat/ChatInterface.tsx
git commit -m "feat(chat): apply Liquid Glass styling to ChatInterface"
```

---

### Task 2.2: Update MessageList with Glass Message Bubbles

**Files:**
- Modify: `frontend/src/components/chat/MessageList.tsx`

**Step 1: Update message bubbles to use glass styling**

For AI messages:
```tsx
<LiquidGlassCard
  brand="teal"
  hover="glow"
  intensity="medium"
  className="max-w-[80%]"
>
  {/* Message content */}
</LiquidGlassCard>
```

For user messages:
```tsx
<LiquidGlassCard
  brand="blue"
  intensity="medium"
  className="max-w-[80%] ml-auto"
>
  {/* Message content */}
</LiquidGlassCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/chat/MessageList.tsx
git commit -m "feat(chat): apply glass styling to message bubbles"
```

---

### Task 2.3: Update ChatInput with Glass Container

**Files:**
- Modify: `frontend/src/components/chat/ChatInput.tsx`

**Step 1: Wrap input in glass container**

```tsx
<LiquidGlassCard intensity="light" radius="xl" className="p-2">
  <div className="flex items-center gap-2">
    {/* Attachment button */}
    <LiquidGlassInput
      placeholder="Ask about your data..."
      className="flex-1 border-0 bg-transparent"
    />
    <LiquidGlassButton variant="primary" radius="lg">
      Send
    </LiquidGlassButton>
  </div>
</LiquidGlassCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/chat/ChatInput.tsx
git commit -m "feat(chat): apply glass styling to ChatInput"
```

---

### Task 2.4: Update ChatPage with Animated Background

**Files:**
- Modify: `frontend/src/pages/ChatPage.tsx`

**Step 1: Add AnimatedBackground and LiquidGlassNavbar**

```tsx
import { AnimatedBackground } from '@/components/ui/AnimatedBackground'
import { LiquidGlassNavbar } from '@/components/ui/LiquidGlassNavbar'
import { LanguageSwitcher } from '@/components/common/LanguageSwitcher'
import { ThemeToggle } from '@/components/ui/ThemeToggle'

export function ChatPage() {
  return (
    <div className="min-h-screen">
      <AnimatedBackground />

      <LiquidGlassNavbar
        left={<Logo />}
        right={
          <>
            <LanguageSwitcher />
            <ThemeToggle />
          </>
        }
      />

      <main className="h-[calc(100vh-64px)]">
        <ChatInterface />
      </main>
    </div>
  )
}
```

**Step 2: Commit**

```bash
git add frontend/src/pages/ChatPage.tsx
git commit -m "feat(chat): add animated background and glass navbar to ChatPage"
```

---

## Phase 3: Dashboard Page Redesign

### Task 3.1: Update DashboardPage with Glass Layout

**Files:**
- Modify: `frontend/src/pages/DashboardPage.tsx`

**Step 1: Add AnimatedBackground and glass structure**

```tsx
import { AnimatedBackground } from '@/components/ui/AnimatedBackground'
import { LiquidGlassNavbar } from '@/components/ui/LiquidGlassNavbar'
import { GlassBrandCard, GlassTealCard, GlassBlueCard, GlassGreenCard } from '@/components/ui/LiquidGlassCard'

export function DashboardPage() {
  return (
    <div className="min-h-screen">
      <AnimatedBackground />

      <LiquidGlassNavbar
        left={<h1 className="text-xl font-semibold">Dashboard</h1>}
        right={
          <>
            <DateRangeSelector />
            <LanguageSwitcher />
            <ThemeToggle />
          </>
        }
      />

      <main className="max-w-7xl mx-auto p-4 md:p-6 space-y-6">
        {/* KPI Cards */}
        <DashboardGrid />

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <ChartWidget title="Revenue Trend" />
          <ChartWidget title="Expense Breakdown" />
        </div>
      </main>
    </div>
  )
}
```

**Step 2: Commit**

```bash
git add frontend/src/pages/DashboardPage.tsx
git commit -m "feat(dashboard): add animated background and glass layout to DashboardPage"
```

---

### Task 3.2: Update DashboardGrid with Glass Cards

**Files:**
- Modify: `frontend/src/components/dashboard/DashboardGrid.tsx`

**Step 1: Update KPI cards to use glass variants**

```tsx
import { GlassBrandCard, GlassTealCard, GlassBlueCard, GlassGreenCard } from '@/components/ui/LiquidGlassCard'

// Revenue card
<GlassBrandCard hover="lift-glow" interactive className="p-4">
  <div className="flex items-center gap-3">
    <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-teal-400 flex items-center justify-center text-xl">
      💰
    </div>
    <div>
      <p className="text-sm liquid-text-secondary">Revenue</p>
      <p className="text-2xl font-bold liquid-text-primary">$124,500</p>
    </div>
  </div>
  <div className="mt-2 text-sm text-emerald-400">↑ 12% from last month</div>
</GlassBrandCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/dashboard/DashboardGrid.tsx
git commit -m "feat(dashboard): apply glass card variants to KPI metrics"
```

---

### Task 3.3: Update PinnedChartCard with Floating Glass

**Files:**
- Modify: `frontend/src/components/dashboard/PinnedChartCard.tsx`

**Step 1: Wrap chart in floating glass card**

```tsx
<LiquidGlassCard
  elevation="floating"
  hover="lift"
  className="p-4"
>
  <div className="flex items-center justify-between mb-4">
    <h3 className="font-semibold liquid-text-primary">{title}</h3>
    <div className="flex gap-2">
      <IconButton icon={<PinIcon />} onClick={handlePin} />
      <IconButton icon={<ExportIcon />} onClick={handleExport} />
    </div>
  </div>
  <div className="h-64">
    <ChartRenderer spec={chartSpec} />
  </div>
</LiquidGlassCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/dashboard/PinnedChartCard.tsx
git commit -m "feat(dashboard): apply floating glass styling to chart cards"
```

---

### Task 3.4: Update ChartPinDialog with Glass Modal

**Files:**
- Modify: `frontend/src/components/dashboard/ChartPinDialog.tsx`

**Step 1: Replace dialog with LiquidGlassModal**

```tsx
import { LiquidGlassModal } from '@/components/ui/LiquidGlassModal'
import { ButtonPrimary, ButtonSecondary } from '@/components/ui/LiquidGlassButton'

<LiquidGlassModal
  isOpen={isOpen}
  onClose={onClose}
  title="Pin Chart to Dashboard"
  actions={
    <>
      <ButtonSecondary onClick={onClose}>Cancel</ButtonSecondary>
      <ButtonPrimary onClick={handlePin}>Pin Chart</ButtonPrimary>
    </>
  }
>
  {/* Dialog content */}
</LiquidGlassModal>
```

**Step 2: Commit**

```bash
git add frontend/src/components/dashboard/ChartPinDialog.tsx
git commit -m "feat(dashboard): replace dialog with LiquidGlassModal for pin dialog"
```

---

## Phase 4: Documents Module Redesign

### Task 4.1: Create LiquidGlassSidebar Component

**Files:**
- Create: `frontend/src/components/ui/LiquidGlassSidebar.tsx`
- Test: `frontend/src/components/ui/LiquidGlassSidebar.test.tsx`

**Step 1: Write the failing test**

```tsx
// frontend/src/components/ui/LiquidGlassSidebar.test.tsx
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
})
```

**Step 2: Run test, then implement**

```tsx
// frontend/src/components/ui/LiquidGlassSidebar.tsx
import React from 'react'
import { cn } from '@/lib/cn'

export interface LiquidGlassSidebarProps {
  children?: React.ReactNode
  className?: string
  collapsed?: boolean
}

export const LiquidGlassSidebar: React.FC<LiquidGlassSidebarProps> = ({
  children,
  className,
  collapsed = false,
}) => {
  return (
    <aside
      className={cn(
        'liquid-glass h-full transition-all duration-300',
        collapsed ? 'w-16' : 'w-64',
        className
      )}
    >
      <div className="p-4 h-full overflow-y-auto liquid-glass-scroll">
        {children}
      </div>
    </aside>
  )
}

export default LiquidGlassSidebar
```

**Step 3: Commit**

```bash
git add frontend/src/components/ui/LiquidGlassSidebar.tsx frontend/src/components/ui/LiquidGlassSidebar.test.tsx
git commit -m "feat(ui): add LiquidGlassSidebar component"
```

---

### Task 4.2: Update DocumentUploader with Glass Styling

**Files:**
- Modify: `frontend/src/components/documents/DocumentUploader.tsx`

**Step 1: Wrap upload zone in glass card**

```tsx
<LiquidGlassCard
  intensity="medium"
  hover="shimmer"
  className="border-2 border-dashed border-white/20 p-8 text-center cursor-pointer"
>
  <div className="text-4xl mb-4">📄</div>
  <p className="liquid-text-primary font-medium">
    Drag PDFs, images, or Excel files here
  </p>
  <p className="liquid-text-secondary text-sm mt-1">or browse</p>
</LiquidGlassCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/documents/DocumentUploader.tsx
git commit -m "feat(documents): apply glass styling to DocumentUploader"
```

---

### Task 4.3: Update ReviewQueue with Glass Cards

**Files:**
- Modify: `frontend/src/components/documents/ReviewQueue.tsx`

**Step 1: Style review items with glass cards**

```tsx
<LiquidGlassCard hover="lift" className="p-4">
  <div className="flex items-center justify-between">
    <div>
      <p className="font-medium liquid-text-primary">{document.name}</p>
      <LiquidGlassBadge variant="warning" className="mt-1">
        {document.status}
      </LiquidGlassBadge>
    </div>
    <LiquidGlassButton variant="primary" size="sm">Review</LiquidGlassButton>
  </div>
</LiquidGlassCard>
```

**Step 2: Commit**

```bash
git add frontend/src/components/documents/ReviewQueue.tsx
git commit -m "feat(documents): apply glass styling to ReviewQueue"
```

---

## Phase 5: Common Components & Final Integration

### Task 5.1: Update LanguageSwitcher with Glass Pill

**Files:**
- Modify: `frontend/src/components/common/LanguageSwitcher.tsx`

**Step 1: Style as glass pill toggle**

```tsx
<div className="liquid-glass flex items-center rounded-full p-1">
  <button
    className={cn(
      'px-3 py-1 rounded-full text-sm font-medium transition-all',
      locale === 'en' ? 'bg-blue-500 text-white' : 'liquid-text-secondary'
    )}
    onClick={() => setLocale('en')}
  >
    EN
  </button>
  <button
    className={cn(
      'px-3 py-1 rounded-full text-sm font-medium transition-all',
      locale === 'ar' ? 'bg-blue-500 text-white' : 'liquid-text-secondary'
    )}
    onClick={() => setLocale('ar')}
  >
    ع
  </button>
</div>
```

**Step 2: Commit**

```bash
git add frontend/src/components/common/LanguageSwitcher.tsx
git commit -m "feat(common): apply glass pill styling to LanguageSwitcher"
```

---

### Task 5.2: Update ThemeToggle with Glass Button

**Files:**
- Modify: `frontend/src/components/ui/ThemeToggle.tsx`

**Step 1: Style as glass icon button**

```tsx
import { IconButton } from '@/components/ui/LiquidGlassButton'

<IconButton
  icon={isDark ? <SunIcon /> : <MoonIcon />}
  onClick={toggleTheme}
  aria-label="Toggle theme"
/>
```

**Step 2: Commit**

```bash
git add frontend/src/components/ui/ThemeToggle.tsx
git commit -m "feat(ui): apply glass styling to ThemeToggle"
```

---

### Task 5.3: Update App.tsx with Global Layout

**Files:**
- Modify: `frontend/src/App.tsx`

**Step 1: Ensure AnimatedBackground is available globally**

```tsx
import { AnimatedBackground } from '@/components/ui/AnimatedBackground'

function App() {
  return (
    <ThemeProvider>
      <AnimatedBackground />
      <Routes>
        {/* routes */}
      </Routes>
    </ThemeProvider>
  )
}
```

**Step 2: Commit**

```bash
git add frontend/src/App.tsx
git commit -m "feat(app): add global AnimatedBackground to app shell"
```

---

### Task 5.4: Export All UI Components from Index

**Files:**
- Create: `frontend/src/components/ui/index.ts`

**Step 1: Create barrel export**

```tsx
// frontend/src/components/ui/index.ts
export * from './LiquidGlassCard'
export * from './LiquidGlassButton'
export * from './LiquidGlassInput'
export * from './LiquidGlassModal'
export * from './LiquidGlassNavbar'
export * from './LiquidGlassSidebar'
export * from './LiquidGlassBadge'
export * from './LiquidGlassToast'
export * from './AnimatedBackground'
export * from './ThemeToggle'
export * from './LoadingSkeleton'
```

**Step 2: Commit**

```bash
git add frontend/src/components/ui/index.ts
git commit -m "feat(ui): create barrel export for all UI components"
```

---

### Task 5.5: Final Verification — Run All Tests

**Step 1: Run full test suite**

Run: `cd frontend && npm test`
Expected: All tests pass

**Step 2: Run linter**

Run: `cd frontend && npm run lint`
Expected: No errors

**Step 3: Run type check**

Run: `cd frontend && npm run typecheck`
Expected: No errors

**Step 4: Commit any fixes**

```bash
git add -A
git commit -m "fix: resolve test and lint issues after Liquid Glass redesign"
```

---

### Task 5.6: Create Final Summary Commit

**Step 1: Create summary**

```bash
git add -A
git commit -m "feat: complete Liquid Glass redesign of all user-facing pages

- Add AnimatedBackground with mesh gradient and floating orbs
- Create LiquidGlassNavbar, Modal, Badge, Toast, Sidebar components
- Redesign ChatPage with glass message bubbles and input
- Redesign DashboardPage with glass KPI cards and chart containers
- Redesign Documents module with glass uploader and review queue
- Update common components (LanguageSwitcher, ThemeToggle)
- Ensure WCAG 2.2 AA compliance throughout
- Add reduced motion support for accessibility

BREAKING CHANGE: All pages now use dark glass mode by default"
```

---

## Summary

| Phase | Tasks | Description |
|-------|-------|-------------|
| **Phase 1** | 1.1 - 1.6 | Foundation: AnimatedBackground, Navbar, Modal, Badge, Toast |
| **Phase 2** | 2.1 - 2.4 | Chat: ChatInterface, MessageList, ChatInput, ChatPage |
| **Phase 3** | 3.1 - 3.4 | Dashboard: Page, Grid, ChartCard, PinDialog |
| **Phase 4** | 4.1 - 4.3 | Documents: Sidebar, Uploader, ReviewQueue |
| **Phase 5** | 5.1 - 5.6 | Integration: Common components, exports, verification |

---

## References

- Design Doc: `docs/plans/2026-02-22-liquid-glass-redesign.md`
- Design System: `docs/DESIGN.md`
- Component Specs: `docs/LIQUID-GLASS-DESIGN-SYSTEM.md`
