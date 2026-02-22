import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { ThemeProvider } from './components/theme'
// Import liquid glassmorphism design system
import './styles/liquid-glass.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="system" attribute="data-theme">
      <App />
    </ThemeProvider>
  </StrictMode>,
)
