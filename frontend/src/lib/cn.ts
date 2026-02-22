/**
 * Class Name Utility
 *
 * Combines clsx and tailwind-merge to conditionally merge Tailwind classes
 * without conflicts. This is essential for component variants and theme states.
 *
 * @example
 * cn('px-4 py-2', isActive && 'bg-blue-500', 'hover:bg-blue-600')
 * // Returns: 'px-4 py-2 bg-blue-500 hover:bg-blue-600'
 *
 * @module lib/cn
 */
import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Merge Tailwind CSS classes without conflicts
 *
 * @param inputs - Class values to merge (strings, objects, arrays)
 * @returns Merged class string with later classes taking precedence
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
