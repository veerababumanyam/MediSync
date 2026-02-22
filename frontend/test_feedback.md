# Fixing Design Feedback

## 1. Contrast Issues in Dark Mode (App.tsx and liquid-glass.css)
The user screenshot showed light-grey cards with white text on a dark background.
**Root Cause**: The `.liquid-glass-light` and `.liquid-glass-medium` classes apply `var(--glass-bg-light)` and `var(--glass-bg-medium)`. In dark mode, these were mapped to `rgba(60,60,60,0.8)` and `rgba(45,45,45,0.6)`. These dense, grey opacities ruined the dark mode aesthetic, appearing like flat muddy grey shapes.
**Fix**: Update the Dark Mode variables in `liquid-glass.css` to use Apple's true dark mode material rendering—typically extremely low-opacity whites or translucent solid blacks—letting the background show through.

## 2. Flat Animated Background
The background appeared completely flat `#0f172a`.
**Root Cause**: The blur of `100px` combined with `.45` opacity on a dark background made the orbs practically invisible.
**Fix**: Increase the opacity of the orbs or use more vibrant variants of the colors for the background elements so they cut through the heavy blur. E.g., `.60` or `.80` opacity.

## 3. Dull Primary Button
The primary button ("Start Chatting") lacked the "pop" expected.
**Root Cause**: `bg-gradient-to-r from-logo-blue to-logo-teal` on dark mode looks slightly subdued.
**Fix**: We should enhance the button with a lighter gradient overlay or a more vibrant text contrast, perhaps a subtle white border, or a brighter shadow `shadow-logo-blue/50`, and a `bg-gradient-to-r from-blue-500 to-logo-teal`. Actually, let's keep the real brand colors but ensure the text is a crisp white and the shadow makes it float.
