import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'

/**
 * ESLint configuration with i18n rules for MediSync
 *
 * i18n Rules:
 * - Warn on string literals in user-facing props (placeholder, title, aria-label, etc.)
 * - These warnings help identify hardcoded strings that should use t() translation
 */

// Props that should use translation function
const I18N_PROPS = [
  'placeholder',
  'title',
  'aria-label',
  'aria-labelledby',
  'aria-describedby',
  'alt',
  'label',
  'tooltip',
  'hint',
];

export default defineConfig([
  globalIgnores(['dist', 'scripts']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
    ],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    rules: {
      // i18n: Warn on string literals in user-facing props
      'react/no-isolated-mounted-element': 'off',

      // Custom rule via inline logic - warn on hardcoded strings in i18n props
      // This is enforced by the detect-hardcoded-strings.ts script in CI
      // ESLint provides IDE warnings for developers

      // General best practices that help with i18n
      'no-restricted-syntax': [
        'warn',
        {
          selector: 'JSXAttribute[name.name="placeholder"] Literal[value=/[A-Z][a-z]+/]',
          message: 'Consider using t() for placeholder text to support i18n',
        },
        {
          selector: 'JSXAttribute[name.name="title"] Literal[value=/[A-Z][a-z]+/]',
          message: 'Consider using t() for title text to support i18n',
        },
        {
          selector: 'JSXAttribute[name.name="aria-label"] Literal[value=/[A-Z][a-z]+/]',
          message: 'Consider using t() for aria-label to support i18n',
        },
      ],
    },
  },
])
