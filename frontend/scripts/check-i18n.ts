#!/usr/bin/env node
/**
 * MediSync i18n Validation Script
 *
 * Validates translation files for:
 * - Key parity: All EN keys must exist in AR
 * - Empty values: No empty translation strings
 * - Interpolation consistency: {{variable}} must match between EN and AR
 *
 * Exit codes:
 *   0 = All checks passed
 *   1 = Errors found
 */

import { readdirSync, readFileSync, existsSync } from 'fs';
import { join } from 'path';

interface TranslationValue {
  [key: string]: string | TranslationValue;
}

interface ValidationError {
  type: 'missing_key' | 'empty_value' | 'mismatched_interpolation' | 'file_error';
  namespace: string;
  key: string;
  message: string;
}

const LOCALES_DIR = join(import.meta.dirname, '../src/i18n/locales');
const SOURCE_LOCALE = 'en';
const TARGET_LOCALE = 'ar';

const errors: ValidationError[] = [];

/**
 * Extract interpolation variables from a string (e.g., {{count}}, {{name}})
 */
function extractInterpolationVars(str: string): string[] {
  const matches = str.match(/\{\{(\w+)\}\}/g);
  return matches ? matches.map(m => m.replace(/[{}]/g, '')) : [];
}

/**
 * Recursively collect all keys with their paths and values
 */
function collectKeys(
  obj: TranslationValue,
  prefix = ''
): Map<string, string | TranslationValue> {
  const keys = new Map<string, string | TranslationValue>();

  for (const [key, value] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;

    if (typeof value === 'string') {
      keys.set(fullKey, value);
    } else if (typeof value === 'object' && value !== null) {
      const nested = collectKeys(value, fullKey);
      for (const [k, v] of nested) {
        keys.set(k, v);
      }
    }
  }

  return keys;
}

/**
 * Validate a single namespace file pair
 */
function validateNamespace(namespace: string): void {
  const enPath = join(LOCALES_DIR, SOURCE_LOCALE, `${namespace}.json`);
  const arPath = join(LOCALES_DIR, TARGET_LOCALE, `${namespace}.json`);

  // Check if target file exists
  if (!existsSync(arPath)) {
    errors.push({
      type: 'file_error',
      namespace,
      key: '',
      message: `Missing translation file: ${TARGET_LOCALE}/${namespace}.json`,
    });
    return;
  }

  let enData: TranslationValue;
  let arData: TranslationValue;

  try {
    enData = JSON.parse(readFileSync(enPath, 'utf-8'));
  } catch {
    errors.push({
      type: 'file_error',
      namespace,
      key: '',
      message: `Failed to parse ${SOURCE_LOCALE}/${namespace}.json`,
    });
    return;
  }

  try {
    arData = JSON.parse(readFileSync(arPath, 'utf-8'));
  } catch {
    errors.push({
      type: 'file_error',
      namespace,
      key: '',
      message: `Failed to parse ${TARGET_LOCALE}/${namespace}.json`,
    });
    return;
  }

  const enKeys = collectKeys(enData);
  const arKeys = collectKeys(arData);

  // Check for missing keys in target locale
  for (const [key, enValue] of enKeys) {
    if (!arKeys.has(key)) {
      errors.push({
        type: 'missing_key',
        namespace,
        key,
        message: `Missing key in ${TARGET_LOCALE}: ${key}`,
      });
      continue;
    }

    const arValue = arKeys.get(key);

    // Check for empty values
    if (typeof arValue === 'string' && arValue.trim() === '') {
      errors.push({
        type: 'empty_value',
        namespace,
        key,
        message: `Empty value in ${TARGET_LOCALE}: ${key}`,
      });
    }

    // Check interpolation variable consistency
    if (typeof enValue === 'string' && typeof arValue === 'string') {
      const enVars = extractInterpolationVars(enValue);
      const arVars = extractInterpolationVars(arValue);

      if (enVars.length > 0 || arVars.length > 0) {
        const enSet = new Set(enVars);
        const arSet = new Set(arVars);

        // Check for mismatched variables
        const missing = enVars.filter(v => !arSet.has(v));
        const extra = arVars.filter(v => !enSet.has(v));

        if (missing.length > 0 || extra.length > 0) {
          errors.push({
            type: 'mismatched_interpolation',
            namespace,
            key,
            message: `Interpolation mismatch in ${key}: EN has [${enVars.join(', ')}], AR has [${arVars.join(', ')}]`,
          });
        }
      }
    }
  }
}

/**
 * Main validation function
 */
function main(): void {
  console.log('\n🌐 MediSync i18n Validation');
  console.log('='.repeat(40));

  // Get all namespace files from source locale
  const enDir = join(LOCALES_DIR, SOURCE_LOCALE);

  if (!existsSync(enDir)) {
    console.error(`\n❌ Source locale directory not found: ${enDir}`);
    process.exit(1);
  }

  const namespaces = readdirSync(enDir)
    .filter(f => f.endsWith('.json'))
    .map(f => f.replace('.json', ''));

  console.log(`\n📁 Found ${namespaces.length} namespaces: ${namespaces.join(', ')}`);

  // Validate each namespace
  for (const ns of namespaces) {
    console.log(`\n🔍 Checking ${ns}...`);
    validateNamespace(ns);
  }

  // Report results
  console.log('\n' + '='.repeat(40));

  if (errors.length === 0) {
    console.log('\n✅ All i18n checks passed!\n');
    process.exit(0);
  }

  console.log(`\n❌ Found ${errors.length} error(s):\n`);

  // Group errors by type
  const byType: Record<string, ValidationError[]> = {};
  for (const err of errors) {
    if (!byType[err.type]) byType[err.type] = [];
    byType[err.type].push(err);
  }

  // Print errors by type
  for (const [type, errs] of Object.entries(byType)) {
    console.log(`\n${type.toUpperCase().replace(/_/g, ' ')} (${errs.length}):`);
    for (const err of errs) {
      console.log(`  • [${err.namespace}] ${err.message}`);
    }
  }

  console.log('\n');
  process.exit(1);
}

main();
