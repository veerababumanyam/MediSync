#!/usr/bin/env node
/**
 * MediSync Hardcoded String Detector
 *
 * Scans .tsx files for hardcoded strings that should be internationalized:
 * - JSX text content: <div>Hello</div>
 * - String props: placeholder="Search", aria-label="Close", title="Info"
 * - Ignores: CSS classes, URLs, numbers, technical strings, imports
 *
 * Exit codes:
 *   0 = No hardcoded strings found
 *   1 = Hardcoded strings detected
 *
 * Note: This script uses only Node.js fs module for file operations.
 * No shell command execution is used.
 */

import { readdirSync, readFileSync, statSync } from 'fs';
import { join, extname, relative } from 'path';

interface HardcodedString {
  file: string;
  line: number;
  type: 'jsx_text' | 'string_prop';
  content: string;
  suggestion: string;
}

const SRC_DIR = join(import.meta.dirname, '../src');

// Props that should be internationalized
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
  'caption',
  'description',
  'heading',
  'subheading',
];

// Patterns to ignore (not user-facing strings)
const IGNORE_PATTERNS = [
  /^https?:\/\//, // URLs
  /^\/[\w/-]+$/, // Routes/paths
  /^#[\w-]+$/, // IDs and hex colors
  /^data-/, // Data attributes
  /^class$/, // class names
  /^className$/, // className values (handled separately)
  /^\d+$/, // Pure numbers
  /^[A-Z_]+$/, // Constants (e.g., ENUM_VALUE)
  /^[a-z-]+$/, // CSS class names, prop names
  /^grid$/, /^flex$/, /^block$/, /^inline/, // CSS values
  /^auto$/, /^none$/, /^inherit$/, // CSS keywords
  /^submit$/, /^reset$/, /^button$/, // Input types
  /^GET$/, /^POST$/, /^PUT$/, /^DELETE$/, // HTTP methods
  /^json$/, /^text$/, /^html$/, // Content types
  /^id$/, /^key$/, /^name$/, // Technical props
  /^type$/, /^role$/, /^tabIndex$/, // ARIA/system props
];

// Strings that are considered technical/valid
const TECHNICAL_STRINGS = new Set([
  // ARIA roles
  'row',
  'column',
  'cell',
  'gridcell',
  'listitem',
  'menuitem',
  'option',
  'button',
  'link',
  'navigation',
  'main',
  'header',
  'footer',
  'article',
  'section',
  'form',
  'search',
  'dialog',
  'alertdialog',
  'tooltip',
  'menu',
  'listbox',
  'combobox',
  'textbox',
  'checkbox',
  'radio',
  'switch',
  'slider',
  'progressbar',
  'spinbutton',
  'tab',
  'tablist',
  'tabpanel',
  'tree',
  'treeitem',
  'group',
  'region',
  // TypeScript/JavaScript keywords
  'Promise',
  'void',
  'null',
  'undefined',
  'string',
  'number',
  'boolean',
  'object',
  'any',
  'never',
  'unknown',
  'async',
  'await',
  'return',
  'const',
  'let',
  'var',
  'function',
  'class',
  'interface',
  'type',
  'enum',
  // React keywords
  'React',
  'Component',
  'FC',
  'Props',
  'State',
  'Ref',
  // Common type names
  'Partial',
  'Required',
  'Readonly',
  'Record',
  'Pick',
  'Omit',
  'Exclude',
  'Extract',
  'NonNullable',
  'ReturnType',
  'Parameters',
  'InstanceType',
]);

const findings: HardcodedString[] = [];

/**
 * Check if a string should be ignored
 */
function shouldIgnore(str: string): boolean {
  const trimmed = str.trim();

  // Empty or whitespace only
  if (trimmed.length === 0) return true;

  // Single character (often icons or punctuation)
  if (trimmed.length === 1) return true;

  // Numbers
  if (/^\d+(\.\d+)?$/.test(trimmed)) return true;

  // Match against ignore patterns
  for (const pattern of IGNORE_PATTERNS) {
    if (pattern.test(trimmed)) return true;
  }

  // Technical ARIA roles
  if (TECHNICAL_STRINGS.has(trimmed.toLowerCase())) return true;

  // CSS-like strings (contain dashes, likely class names)
  if (/^[a-z]+(-[a-z]+)+$/.test(trimmed) && !trimmed.includes(' ')) return true;

  // File extensions and imports
  if (/\.(tsx?|jsx?|css|scss|json|svg|png|jpg|gif)$/.test(trimmed)) return true;

  // Variable-like strings (camelCase starting with lowercase)
  if (/^[a-z][a-zA-Z0-9]*$/.test(trimmed) && trimmed.length < 20) return true;

  return false;
}

/**
 * Check if a string looks like user-facing text
 */
function isUserFacingText(str: string): boolean {
  const trimmed = str.trim();

  // Contains spaces (likely a phrase)
  if (trimmed.includes(' ')) return true;

  // Starts with capital letter (likely a label/title)
  if (/^[A-Z]/.test(trimmed) && trimmed.length > 3) return true;

  // Contains punctuation (likely a sentence)
  if (/[.!?,;:]$/.test(trimmed)) return true;

  // Longer strings are more likely to be user-facing
  if (trimmed.length > 15) return true;

  return false;
}

/**
 * Scan a single file for hardcoded strings
 */
function scanFile(filePath: string, srcDir: string): void {
  const content = readFileSync(filePath, 'utf-8');
  const lines = content.split('\n');
  const relPath = relative(srcDir, filePath);

  lines.forEach((line, index) => {
    const lineNum = index + 1;

    // Skip import statements
    if (line.trim().startsWith('import ') || line.includes('from ')) return;

    // Skip comments
    if (line.trim().startsWith('//') || line.trim().startsWith('*')) return;

    // Skip lines that already use translation function
    if (line.includes('t(') || line.includes('useTranslation')) return;

    // Skip type definitions
    if (line.includes('interface ') || line.includes('type ')) return;

    // Skip test files specific patterns
    if (line.includes('describe(') || line.includes('it(') || line.includes('expect(')) return;

    // Detect JSX text content: >Text<
    const jsxTextMatches = line.matchAll(/>\s*([A-Z][a-zA-Z\s,.!?']+?)\s*</g);
    for (const match of jsxTextMatches) {
      const text = match[1].trim();
      if (!shouldIgnore(text) && isUserFacingText(text)) {
        findings.push({
          file: relPath,
          line: lineNum,
          type: 'jsx_text',
          content: text,
          suggestion: `{t('your.key.here')}`,
        });
      }
    }

    // Detect string props: placeholder="Text"
    for (const prop of I18N_PROPS) {
      // Pattern: prop="value" or prop={'value'}
      const directPattern = new RegExp(`${prop}=["']([^"']+)["']`, 'g');
      const bracePattern = new RegExp(`${prop}=\\{["']([^"']+)["']\\}`, 'g');

      let match;
      while ((match = directPattern.exec(line)) !== null) {
        const text = match[1];
        if (!shouldIgnore(text)) {
          findings.push({
            file: relPath,
            line: lineNum,
            type: 'string_prop',
            content: `${prop}="${text}"`,
            suggestion: `${prop}={t('your.key.here')}`,
          });
        }
      }

      while ((match = bracePattern.exec(line)) !== null) {
        const text = match[1];
        if (!shouldIgnore(text)) {
          findings.push({
            file: relPath,
            line: lineNum,
            type: 'string_prop',
            content: `${prop}="${text}"`,
            suggestion: `${prop}={t('your.key.here')}`,
          });
        }
      }
    }
  });
}

/**
 * Recursively find all .tsx files
 */
function findTsxFiles(dir: string): string[] {
  const files: string[] = [];

  // Skip node_modules, test directories, and other non-source directories
  const skipDirs = new Set(['node_modules', '__tests__', '__mocks__', '.git', 'dist', 'build']);

  const items = readdirSync(dir);

  for (const item of items) {
    const fullPath = join(dir, item);

    try {
      const stat = statSync(fullPath);

      if (stat.isDirectory()) {
        if (!skipDirs.has(item)) {
          files.push(...findTsxFiles(fullPath));
        }
      } else if (stat.isFile() && extname(item) === '.tsx') {
        // Skip test files
        if (!item.endsWith('.test.tsx') && !item.endsWith('.spec.tsx')) {
          files.push(fullPath);
        }
      }
    } catch {
      // Skip files/dirs that can't be accessed
      continue;
    }
  }

  return files;
}

/**
 * Main function
 */
function main(): void {
  console.log('\n🔍 MediSync Hardcoded String Detector');
  console.log('='.repeat(40));

  const files = findTsxFiles(SRC_DIR);
  console.log(`\n📁 Scanning ${files.length} .tsx files in src/...`);

  // Scan each file
  for (const file of files) {
    scanFile(file, SRC_DIR);
  }

  // Report results
  console.log('\n' + '='.repeat(40));

  if (findings.length === 0) {
    console.log('\n✅ No hardcoded strings found!\n');
    process.exit(0);
  }

  console.log(`\n⚠️  Found ${findings.length} potential hardcoded string(s):\n`);

  // Group by file
  const byFile: Record<string, HardcodedString[]> = {};
  for (const f of findings) {
    if (!byFile[f.file]) byFile[f.file] = [];
    byFile[f.file].push(f);
  }

  // Print findings
  for (const [file, items] of Object.entries(byFile)) {
    console.log(`\n📄 ${file}`);
    for (const item of items) {
      console.log(`   Line ${item.line}: ${item.content}`);
      console.log(`   Suggestion: ${item.suggestion}`);
    }
  }

  console.log('\n💡 Tip: Add translations to the appropriate namespace file and use useTranslation() hook.\n');
  process.exit(1);
}

main();
