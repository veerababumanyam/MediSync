# WCAG 3.0 Compliance Checklist

Comprehensive checklist for ensuring AnySync components meet WCAG 3.0 Bronze minimum compliance.

## APCA Contrast Requirements

WCAG 3.0 uses APCA (Accessible Perceptual Contrast Algorithm) instead of the WCAG 2.x contrast ratios.

### Minimum Lightness Contrast (Lc) Values

| Use Case | Minimum Lc | Target Lc |
|----------|------------|-----------|
| Body text (< 24px) | 60 Lc | 75 Lc |
| Large text (≥ 24px) | 45 Lc | 60 Lc |
| Non-text elements (icons, borders) | 45 Lc | 60 Lc |
| Placeholder text | 45 Lc | 60 Lc |
| Focused UI components | 45 Lc | 60 Lc |

### AnySync Compliant Colors

#### Light Mode (on #FFFFFF background)

```css
/* Text Colors */
--text-primary:   #0F172A;  /* ~106 Lc ✅ */
--text-secondary: #475569;  /* ~68 Lc ✅ */
--text-tertiary:  #576574;  /* ~61 Lc ✅ */
--text-muted:     #6B7280;  /* ~57 Lc (large text only) */

/* Brand Colors on White */
--brand-primary:  #1B3FA0;  /* ~72 Lc ✅ */
--brand-accent:   #0F8A3D;  /* ~60 Lc ✅ (dark variant for text) */
```

#### Dark Mode (on #0F172A background)

```css
/* Text Colors */
--text-primary:   #F8FAFC;  /* ~107 Lc ✅ */
--text-secondary: #CBD5E1;  /* ~75 Lc ✅ */
--text-tertiary:  #94A3B8;  /* ~62 Lc ✅ */

/* Brand Colors on Dark */
--text-brand:     #7596F5;  /* Primary-300 for readability */
--text-accent:    #3DE878;  /* Accent-400, vibrant on dark */
```

## Perceivable

### 1.1 Text Alternatives

- [ ] All informative images have descriptive alt text
- [ ] Decorative images use empty alt="" or aria-hidden
- [ ] Complex images have extended descriptions
- [ ] Image buttons describe the action, not the image

### 1.2 Time-based Media

- [ ] Video has captions
- [ ] Video has audio description
- [ ] Audio has transcript

### 1.3 Adaptable

- [ ] Content structure uses semantic HTML
- [ ] Reading sequence is logical
- [ ] Instructions don't rely solely on sensory characteristics

### 1.4 Distinguishable

- [ ] Text can be resized to 200% without loss of content
- [ ] Text spacing can be modified without loss of content
- [ ] Content reflows at 320px width
- [ ] APCA contrast meets minimum requirements
- [ ] Text isn't presented as images of text

## Operable

### 2.1 Keyboard Accessible

- [ ] All functionality available via keyboard
- [ ] No keyboard traps
- [ ] Focus can be moved away from all components

### 2.2 Enough Time

- [ ] Time limits can be turned off or extended
- [ ] Moving content can be paused

### 2.3 Seizures and Physical Reactions

- [ ] No content flashes more than 3 times per second

### 2.4 Navigable

- [ ] Pages have descriptive titles
- [ ] Headings establish correct hierarchy
- [ ] Focus order is logical
- [ ] Focus is visible (not just color change)
- [ ] Multiple navigation methods available

### 2.5 Input Modalities

- [ ] Touch targets are at least 44x44px (WCAG 2.5.8)
- [ ] Pointer gestures have single-pointer alternatives
- [ ] Motion-activated features have alternatives

## Understandable

### 3.1 Readable

- [ ] Page language is declared
- [ ] Language changes are marked

### 3.2 Predictable

- [ ] Focus doesn't trigger context changes
- [ ] Input doesn't trigger unexpected context changes
- [ ] Navigation is consistent across pages
- [ ] Components are identified consistently

### 3.3 Input Assistance

- [ ] Form inputs have visible labels
- [ ] Form inputs have accessible names matching visible labels
- [ ] Error messages are specific and helpful
- [ ] Error prevention for important actions
- [ ] Help available for complex forms

## Robust

### 4.1 Compatible

- [ ] Valid HTML markup
- [ ] Name, Role, Value for custom components
- [ ] Status messages announced to screen readers

## WCAG 3.0 Specific

### Focus Indicators (WCAG 3.0)

- [ ] Focus ring is at least 2px thick
- [ ] Focus ring has at least 2px offset from component
- [ ] Focus ring has ≥ 3:1 contrast against adjacent colors
- [ ] Gold focus (#FFD700) used for maximum visibility
- [ ] Dark mode uses lighter gold (#FFED4E)

### Media Queries

- [ ] `prefers-reduced-motion: reduce` - Animations disabled
- [ ] `prefers-contrast: high` - High contrast mode supported
- [ ] `forced-colors: active` - Windows high contrast supported
- [ ] `prefers-reduced-transparency: reduce` - Solid backgrounds

### Cognitive Accessibility

- [ ] Clear visual hierarchy
- [ ] Consistent navigation
- [ ] Error prevention and recovery
- [ ] Help and documentation available
- [ ] No unexpected context changes

## Testing Tools

### Automated Testing

```bash
# Run accessibility tests
npm run test:a11y

# Lighthouse accessibility audit
npx lighthouse https://localhost:3000 --only-categories=accessibility
```

### Manual Testing

1. **Keyboard Navigation** - Tab through entire page
2. **Screen Reader** - Test with VoiceOver (Mac) or NVDA (Windows)
3. **Zoom** - Set browser to 200% zoom
4. **Mobile** - Test touch targets on actual device
5. **High Contrast** - Enable Windows High Contrast mode
6. **Reduced Motion** - Enable reduced motion in OS settings

## Common Failures to Avoid

1. **Placeholder text as labels** - Use actual `<label>` elements
2. **Color-only indicators** - Add icons or text
3. **Tiny touch targets** - Ensure 44px minimum
4. **Missing focus states** - Always include `:focus-visible`
5. **Auto-playing media** - Don't auto-play with sound
6. **No skip links** - Include "Skip to main content"
7. **Inaccessible modals** - Trap focus, manage aria-hidden
8. **Form errors after submit** - Show errors inline, announce to SR
