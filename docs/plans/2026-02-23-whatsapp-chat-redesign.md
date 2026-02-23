# WhatsApp-Inspired Chat Experience Design

**Document ID:** 2026-02-23-whatsapp-chat-redesign
**Status:** Approved
**Created:** 2026-02-23
**Author:** Claude (Design Session)

---

## 1. Executive Summary

Transform MediSync's chat interface from a single-panel layout to a production-grade, WhatsApp-inspired two-panel experience. This is the core of the application where users spend 90% of their time, requiring exceptional usability, elegance, responsiveness, and robustness.

### Key Outcomes
- Conversations sidebar for session management
- Responsive layout: two-panel desktop, single-panel mobile with slide navigation
- Full-height, full-bleed chat experience on mobile
- Keyboard-aware input with scroll-to-bottom FAB
- Message grouping by date with status indicators
- WCAG 3.0 Gold accessibility compliance

---

## 2. Architecture Overview

### 2.1 Component Structure

```
ChatPage/
├── ChatLayout (NEW - responsive layout controller)
│   ├── ConversationsSidebar (NEW - left panel)
│   │   ├── SidebarHeader
│   │   │   ├── SearchInput
│   │   │   └── NewChatButton
│   │   └── ConversationList
│   │       └── ConversationItem (multiple)
│   │           ├── SessionTitle/Preview
│   │           ├── Timestamp
│   │           └── UnreadIndicator
│   │
│   └── ChatPanel (enhanced ChatInterface)
│       ├── ChatHeader (responsive)
│       │   ├── BackButton (mobile only)
│       │   ├── SessionInfo
│       │   └── ActionMenu (overflow)
│       ├── MessagesArea
│       │   ├── WelcomeScreen (empty state)
│       │   ├── MessageList
│       │   ├── StreamingIndicator
│       │   └── ScrollToBottomFab
│       └── ChatInput
│           ├── AttachmentButton (future)
│           ├── TextInput
│           ├── VoiceButton (future placeholder)
│           └── SendButton
```

### 2.2 Layout Diagrams

#### Desktop (≥ 1024px)
```
┌──────────────────────────────────────────────────────────────────────────┐
│  ┌──────────────────────┬───────────────────────────────────────────────┐│
│  │  Conversations Panel │  Chat Panel                                   ││
│  │  (width: 360px)      │  (flex: 1, max-width: 900px)                  ││
│  │  ┌────────────────┐  │  ┌─────────────────────────────────────────┐ ││
│  │  │ Search/Header  │  │  │ ChatHeader                              │ ││
│  │  ├────────────────┤  │  ├─────────────────────────────────────────┤ ││
│  │  │                │  │  │                                         │ ││
│  │  │  Conversation  │  │  │  Messages Area (full height)            │ ││
│  │  │  List          │  │  │                                         │ ││
│  │  │  (scrollable)  │  │  │                                         │ ││
│  │  │                │  │  │                                         │ ││
│  │  │                │  │  ├─────────────────────────────────────────┤ ││
│  │  │                │  │  │ ChatInput (sticky bottom)               │ ││
│  │  └────────────────┘  │  └─────────────────────────────────────────┘ ││
│  └──────────────────────┴───────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────┘
```

#### Mobile (< 640px)
```
┌─────────────────────────────────────┐
│  ┌─────────────────────────────────┐│
│  │ Header (with menu/back)         ││
│  ├─────────────────────────────────┤│
│  │                                 ││
│  │  Either:                        ││
│  │  - Conversations List (default) ││
│  │  - OR Active Chat (when open)   ││
│  │                                 ││
│  │  Full-bleed, max height         ││
│  │                                 ││
│  ├─────────────────────────────────┤│
│  │ ChatInput (only in active chat) ││
│  └─────────────────────────────────┘│
│  [Scroll-to-bottom FAB]             │
└─────────────────────────────────────┘
```

---

## 3. Responsive Breakpoints

| Breakpoint | Layout | Sidebar Width | Chat Panel |
|------------|--------|---------------|------------|
| Mobile (< 640px) | Single panel, toggle | Full-screen overlay | Full-screen |
| Tablet (640-1024px) | Single panel, slide | Slide-in (280px) | Full-width |
| Desktop (≥ 1024px) | Two-panel | Fixed (360px) | Flex (max 800px) |
| Large (≥ 1440px) | Two-panel centered | Fixed (380px) | Centered (max 900px) |

### Height Calculations
```css
--chat-viewport-height: calc(100dvh - var(--header-height, 64px));
--chat-content-height: calc(var(--chat-viewport-height) - var(--chat-input-height, 80px));
```

---

## 4. Component Specifications

### 4.1 ChatLayout (NEW)

**Purpose:** Root layout component that manages responsive sidebar/chat positioning.

```tsx
interface ChatLayoutProps {
  sessions: ChatSession[]
  activeSessionId: string | null
  onSelectSession: (id: string) => void
  onNewSession: () => void
}
```

**Behavior:**
- Desktop: Always shows both panels
- Mobile/Tablet: Manages sidebar open/close state
- Handles responsive breakpoint detection

### 4.2 ConversationsSidebar (NEW)

**Purpose:** Left panel with session list and search.

```tsx
interface ConversationsSidebarProps {
  sessions: ChatSession[]
  activeSessionId: string | null
  onSelectSession: (id: string) => void
  onNewSession: () => void
  onDeleteSession: (id: string) => void
  onSearch: (query: string) => void
  isOpen: boolean
  onClose: () => void
}
```

**Features:**
- Search with instant filtering
- Session preview (first 80 chars of last message)
- Relative timestamps ("2m ago", "Yesterday", "Feb 20")
- Active session highlight with brand gradient border
- Empty state with illustration and CTA
- Swipe-to-delete on mobile (with confirmation)

### 4.3 ConversationItem (NEW)

**Purpose:** Single session row in the list.

```tsx
interface ConversationItemProps {
  session: ChatSession
  isActive: boolean
  onClick: () => void
  onDelete: () => void
}
```

**Visual Design:**
```
┌────────────────────────────────────────┐
│ 🔵  Revenue query for February...      │
│     "Show me the revenue trends..."    │  ← Preview text
│                           2m ago  ●   │  ← Timestamp + unread
└────────────────────────────────────────┘
```

### 4.4 ChatPanel (Enhanced ChatInterface)

**Purpose:** Main chat area with messages and input.

**New Features:**
- Full-height layout with `100dvh`
- Keyboard avoidance on mobile (visualViewport API)
- Message grouping by date ("Today", "Yesterday", "Feb 20")
- Typing indicator animation
- Message status indicators (sending, sent, error)

### 4.5 ScrollToBottomFab (NEW)

**Purpose:** Floating action button to scroll to latest messages.

```tsx
interface ScrollToBottomFabProps {
  visible: boolean
  onClick: () => void
  unreadCount?: number
}
```

**Behavior:**
- Appears when scrolled > 100px from bottom
- Shows unread count badge if messages arrived while scrolled up
- Smooth fade-in/out animation
- Haptic feedback on tap (if available)

### 4.6 ChatHeader (Enhanced)

**Changes:**
- Add BackButton (visible on mobile only)
- Add overflow menu for actions (New Chat, Language, Theme)
- Compact mode for narrow screens
- Session title display

---

## 5. User Experience Patterns

### 5.1 Navigation Flow

```
Open App → See Recent Sessions → Select Session → View Messages
                                                         ↓
                          Scroll History ← ← Type Query → Receive Response
                                                                   ↓
                                                           View Charts
```

### 5.2 Interaction Matrix

| Action | Mobile | Desktop |
|--------|--------|---------|
| Open session | Tap item → slides to chat | Click item → loads in panel |
| Back to list | Swipe right OR tap back | N/A (always visible) |
| New chat | FAB or header button | Header button |
| Search | Tap search → full-screen | Inline search |
| Delete session | Swipe left → delete | Right-click → context menu |
| Scroll to bottom | FAB appears when scrolled up | Same |

---

## 6. Accessibility (WCAG 3.0 Gold)

| Feature | Implementation |
|---------|---------------|
| Focus management | Focus moves to chat when session opens |
| Keyboard navigation | Arrow keys navigate list, Enter selects |
| Screen reader | Live region announces new messages |
| Touch targets | Minimum 44×44px for all interactive elements |
| Motion | Respects `prefers-reduced-motion` |
| Contrast | APCA-compliant text on all backgrounds |
| Safe areas | Uses `env(safe-area-inset-*)` for notch/home indicator |

---

## 7. Animation & Polish

| Animation | Duration | Easing |
|-----------|----------|--------|
| Sidebar slide | 300ms | `--ease-liquid` |
| Message entrance | 200ms | `--ease-liquid-out` |
| FAB appear/disappear | 150ms | `--ease-liquid` |
| Typing dots | 1.4s | infinite |
| Shimmer effect | 2s | infinite |

### Micro-interactions
- Button press: `scale(0.98)`
- Session hover: subtle background shift
- Send button: color transition + slight scale
- Message bubbles: staggered entrance (50ms delay each)

---

## 8. Files to Create/Modify

### New Files
| File | Purpose |
|------|---------|
| `components/chat/ChatLayout.tsx` | Root layout controller |
| `components/chat/ConversationsSidebar/` | Sidebar components directory |
| `components/chat/ConversationsSidebar/ConversationsSidebar.tsx` | Main sidebar |
| `components/chat/ConversationsSidebar/SidebarHeader.tsx` | Search + new chat |
| `components/chat/ConversationsSidebar/ConversationList.tsx` | Session list |
| `components/chat/ConversationsSidebar/ConversationItem.tsx` | Single session |
| `components/chat/ConversationsSidebar/EmptyConversations.tsx` | Empty state |
| `components/chat/ChatPanel/` | Panel components directory |
| `components/chat/ChatPanel/ChatPanel.tsx` | Renamed ChatInterface |
| `components/chat/ChatPanel/MessagesArea.tsx` | Extracted messages area |
| `components/chat/ChatPanel/MessageGroup.tsx` | Date grouping |
| `components/chat/ChatPanel/ScrollToBottomFab.tsx` | Scroll FAB |
| `components/chat/ChatPanel/TypingIndicator.tsx` | Typing animation |
| `hooks/useChatLayout.ts` | Responsive layout logic |
| `hooks/useKeyboardAvoidance.ts` | Mobile keyboard handling |
| `hooks/useScrollToBottom.ts` | Scroll FAB visibility |
| `styles/chat.css` | Chat-specific styles |

### Modified Files
| File | Changes |
|------|---------|
| `pages/ChatPage.tsx` | Simplified to render ChatLayout |
| `components/chat/ChatInterface.tsx` | Rename to ChatPanel, enhance |
| `components/chat/ChatHeader.tsx` | Add back button, overflow menu |
| `components/chat/ChatInput.tsx` | Keyboard avoidance integration |
| `components/chat/MessageList.tsx` | Responsive bubble width |
| `styles/globals.css` | Container padding updates |
| `styles/tokens.css` | New chat-specific tokens |
| `styles/liquid-glass.css` | Chat panel glass variants |

---

## 9. Implementation Phases

### Phase 1: Layout Foundation (HIGH)
- Create ChatLayout component
- Create ConversationsSidebar structure
- Implement responsive breakpoint logic
- Basic two-panel layout on desktop

### Phase 2: Sidebar Features (HIGH)
- ConversationList with session items
- Search functionality
- New session button
- Empty state

### Phase 3: Chat Panel Enhancements (HIGH)
- Dynamic viewport height (100dvh)
- Keyboard avoidance hook
- Scroll-to-bottom FAB
- Message responsive padding

### Phase 4: Polish & UX (MEDIUM)
- Message grouping by date
- Typing indicator
- Message status indicators
- Animation refinements

### Phase 5: Future Enhancements (LOW)
- Voice input button (placeholder)
- Attachment button (placeholder)
- Reply-to preview
- Session rename/delete

---

## 10. Design Tokens

### New Tokens to Add

```css
:root {
  /* Chat Layout */
  --chat-sidebar-width: 360px;
  --chat-sidebar-width-tablet: 280px;
  --chat-panel-max-width: 900px;
  --chat-input-min-height: 80px;

  /* Chat Offsets */
  --chat-offset-top: var(--header-height, 64px);
  --chat-viewport-height: calc(100dvh - var(--chat-offset-top));

  /* Animation */
  --chat-slide-duration: var(--duration-moderate);
  --chat-fab-duration: var(--duration-fast);
}
```

---

## 11. Success Criteria

- [ ] Sidebar shows all user sessions with search
- [ ] Two-panel layout works on desktop (≥ 1024px)
- [ ] Single-panel toggle works on mobile (< 1024px)
- [ ] Chat uses full viewport height with 100dvh
- [ ] Keyboard avoidance works on iOS Safari and Chrome Android
- [ ] Scroll-to-bottom FAB appears when scrolled up
- [ ] All WCAG 3.0 Gold accessibility requirements met
- [ ] Animations respect prefers-reduced-motion
- [ ] RTL layout works correctly
- [ ] No regression in existing chat functionality

---

## 12. References

- [DESIGN.md](../DESIGN.md) - Design system reference
- [globals.css](../../frontend/src/styles/globals.css) - Current global styles
- [tokens.css](../../frontend/src/styles/tokens.css) - Design tokens
- [liquid-glass.css](../../frontend/src/styles/liquid-glass.css) - Glassmorphism utilities
