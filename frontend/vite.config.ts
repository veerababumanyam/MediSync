import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

import tailwindcss from '@tailwindcss/vite'

/**
 * Vite configuration for MediSync Frontend
 *
 * Features:
 * - Hot Module Replacement (HMR) for development
 * - Path aliases for clean imports
 * - Production optimizations (chunk splitting, minification)
 * - Environment variable handling
 */
export default defineConfig(({ mode }) => ({
  plugins: [react(), tailwindcss()],

  // Path aliases for cleaner imports
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@pages': path.resolve(__dirname, './src/pages'),
      '@hooks': path.resolve(__dirname, './src/hooks'),
      '@lib': path.resolve(__dirname, './src/lib'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@i18n': path.resolve(__dirname, './src/i18n'),
      '@types': path.resolve(__dirname, './src/types'),
    },
  },

  // Development server configuration
  server: {
    host: true, // Listen on all addresses
    port: 3000,
    strictPort: false,
    open: false, // Don't auto-open browser
    // Proxy API requests to backend in development
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      '/ws': {
        target: process.env.VITE_WS_URL || 'ws://localhost:8080',
        ws: true,
      },
    },
    hmr: {
      overlay: true,
    },
  },

  // Preview server configuration (for production builds)
  preview: {
    host: true,
    port: 4173,
    strictPort: false,
    // Proxy for preview mode
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },

  // Build optimizations for production
  build: {
    target: 'esnext',
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: mode === 'development',
    minify: 'terser',
    // Split vendor code for better caching
    rollupOptions: {
      output: {
        manualChunks: {
          // React and React DOM
          'react-vendor': ['react', 'react-dom'],
          // CopilotKit
          'copilotkit': ['@copilotkit/react-core', '@copilotkit/react-ui'],
          // i18n
          'i18n': ['i18next', 'react-i18next'],
          // Charts
          'echarts': ['echarts', 'echarts-for-react'],
        },
        // Asset naming for cache busting
        chunkFileNames: 'assets/js/[name]-[hash].js',
        entryFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
      },
    },
    // Chunk size warning threshold
    chunkSizeWarningLimit: 1000,
  },

  // Environment variables prefix
  envPrefix: 'VITE_',

  // Optimize dependency pre-bundling
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      '@copilotkit/react-core',
      '@copilotkit/react-ui',
      'i18next',
      'react-i18next',
      'echarts',
    ],
  },
}))
