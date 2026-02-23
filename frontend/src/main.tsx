import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { ThemeProvider } from './components/theme'
// Import design system (order matters - tokens first, then components)
import './styles/design-tokens.css'  // Single source of truth for all design values
import './styles/liquid-glass.css'    // Liquid glass component classes
import './styles/animations.css'      // Unified animation library

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="system" attribute="data-theme">
      <App />
    </ThemeProvider>
  </StrictMode>,
)
