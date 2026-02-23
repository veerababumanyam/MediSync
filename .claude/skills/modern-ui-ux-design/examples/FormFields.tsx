/**
 * Accessible Form Components
 *
 * Form components with WCAG 3.0 accessibility compliance including:
 * - Associated labels
 * - ARIA attributes
 * - Error states
 * - Focus indicators
 * - Proper contrast
 */

import { forwardRef, InputHTMLAttributes, TextareaHTMLAttributes, SelectHTMLAttributes, ReactNode, useId } from 'react';
import { cn } from '@/lib/utils';

// ============================================
// TEXT INPUT
// ============================================

interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  /** Input label (required for accessibility) */
  label: string;
  /** Helper text displayed below input */
  helperText?: string;
  /** Error message (makes input invalid) */
  error?: string;
  /** Input size */
  size?: 'sm' | 'md' | 'lg';
  /** Icon to display at start */
  startIcon?: ReactNode;
  /** Icon to display at end */
  endIcon?: ReactNode;
  /** Visually hide label (still accessible) */
  hideLabel?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      helperText,
      error,
      size = 'md',
      startIcon,
      endIcon,
      hideLabel = false,
      className,
      id: providedId,
      'aria-describedby': ariaDescribedBy,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const inputId = providedId || generatedId;
    const errorId = `${inputId}-error`;
    const helperId = `${inputId}-helper`;

    const sizeClasses = {
      sm: 'h-10 text-sm px-3',
      md: 'h-11 text-base px-4',
      lg: 'h-12 text-lg px-5',
    };

    const describedBy = [
      error && errorId,
      helperText && helperId,
      ariaDescribedBy,
    ]
      .filter(Boolean)
      .join(' ') || undefined;

    return (
      <div className={cn('liquid-clear-label', className)}>
        <label
          htmlFor={inputId}
          className={cn(
            'block font-semibold text-sm mb-2',
            'text-primary dark:text-primary-dark',
            hideLabel && 'sr-only'
          )}
        >
          {label}
          {props.required && (
            <span className="text-error ml-1" aria-hidden="true">*</span>
          )}
        </label>

        <div className="relative">
          {startIcon && (
            <span className="absolute start-3 top-1/2 -translate-y-1/2 text-muted pointer-events-none">
              {startIcon}
            </span>
          )}

          <input
            ref={ref}
            id={inputId}
            className={cn(
              'liquid-glass-input',
              'w-full',
              sizeClasses[size],
              startIcon && 'ps-10',
              endIcon && 'pe-10',
              error && 'state-error'
            )}
            aria-invalid={error ? 'true' : undefined}
            aria-describedby={describedBy}
            aria-errormessage={error ? errorId : undefined}
            {...props}
          />

          {endIcon && (
            <span className="absolute end-3 top-1/2 -translate-y-1/2 text-muted pointer-events-none">
              {endIcon}
            </span>
          )}
        </div>

        {helperText && !error && (
          <p
            id={helperId}
            className="mt-2 text-sm text-tertiary"
          >
            {helperText}
          </p>
        )}

        {error && (
          <p
            id={errorId}
            role="alert"
            className="mt-2 text-sm text-error font-medium"
          >
            {error}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

// ============================================
// TEXTAREA
// ============================================

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  /** Textarea label (required for accessibility) */
  label: string;
  /** Helper text */
  helperText?: string;
  /** Error message */
  error?: string;
  /** Visually hide label */
  hideLabel?: boolean;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  (
    {
      label,
      helperText,
      error,
      hideLabel = false,
      className,
      id: providedId,
      rows = 4,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const textareaId = providedId || generatedId;
    const errorId = `${textareaId}-error`;
    const helperId = `${textareaId}-helper`;

    return (
      <div className={cn('liquid-clear-label', className)}>
        <label
          htmlFor={textareaId}
          className={cn(
            'block font-semibold text-sm mb-2',
            'text-primary dark:text-primary-dark',
            hideLabel && 'sr-only'
          )}
        >
          {label}
          {props.required && (
            <span className="text-error ml-1" aria-hidden="true">*</span>
          )}
        </label>

        <textarea
          ref={ref}
          id={textareaId}
          rows={rows}
          className={cn(
            'liquid-glass-input',
            'w-full p-4 min-h-[100px]',
            'resize-y',
            error && 'state-error'
          )}
          aria-invalid={error ? 'true' : undefined}
          aria-describedby={error ? errorId : helperText ? helperId : undefined}
          {...props}
        />

        {helperText && !error && (
          <p id={helperId} className="mt-2 text-sm text-tertiary">
            {helperText}
          </p>
        )}

        {error && (
          <p id={errorId} role="alert" className="mt-2 text-sm text-error font-medium">
            {error}
          </p>
        )}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';

// ============================================
// SELECT
// ============================================

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  /** Select label (required for accessibility) */
  label: string;
  /** Options array */
  options: Array<{ value: string; label: string; disabled?: boolean }>;
  /** Placeholder text */
  placeholder?: string;
  /** Error message */
  error?: string;
  /** Helper text */
  helperText?: string;
  /** Visually hide label */
  hideLabel?: boolean;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  (
    {
      label,
      options,
      placeholder,
      error,
      helperText,
      hideLabel = false,
      className,
      id: providedId,
      children,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const selectId = providedId || generatedId;
    const errorId = `${selectId}-error`;

    return (
      <div className={cn('liquid-clear-label', className)}>
        <label
          htmlFor={selectId}
          className={cn(
            'block font-semibold text-sm mb-2',
            'text-primary dark:text-primary-dark',
            hideLabel && 'sr-only'
          )}
        >
          {label}
          {props.required && (
            <span className="text-error ml-1" aria-hidden="true">*</span>
          )}
        </label>

        <div className="relative">
          <select
            ref={ref}
            id={selectId}
            className={cn(
              'liquid-glass-input',
              'w-full h-11 px-4 pe-10',
              'appearance-none cursor-pointer',
              error && 'state-error'
            )}
            aria-invalid={error ? 'true' : undefined}
            aria-describedby={error ? errorId : undefined}
            {...props}
          >
            {placeholder && (
              <option value="" disabled>
                {placeholder}
              </option>
            )}
            {options.map((option) => (
              <option
                key={option.value}
                value={option.value}
                disabled={option.disabled}
              >
                {option.label}
              </option>
            ))}
          </select>

          {/* Chevron icon */}
          <span className="absolute end-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </span>
        </div>

        {helperText && !error && (
          <p className="mt-2 text-sm text-tertiary">{helperText}</p>
        )}

        {error && (
          <p id={errorId} role="alert" className="mt-2 text-sm text-error font-medium">
            {error}
          </p>
        )}
      </div>
    );
  }
);

Select.displayName = 'Select';

// ============================================
// CHECKBOX
// ============================================

interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'> {
  /** Checkbox label (required for accessibility) */
  label: string;
  /** Helper text */
  helperText?: string;
  /** Error state */
  error?: boolean;
}

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  (
    {
      label,
      helperText,
      error,
      className,
      id: providedId,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const checkboxId = providedId || generatedId;

    return (
      <div className={cn('flex items-start gap-3', className)}>
        <input
          ref={ref}
          type="checkbox"
          id={checkboxId}
          className={cn(
            'w-5 h-5 mt-0.5 shrink-0',
            'rounded border-2 border-default',
            'bg-surface-primary',
            'text-brand-primary',
            'focus:ring-2 focus:ring-brand-primary focus:ring-offset-2',
            'dark:focus:ring-focus-gold',
            error && 'border-error'
          )}
          aria-invalid={error ? 'true' : undefined}
          {...props}
        />

        <div className="flex flex-col">
          <label
            htmlFor={checkboxId}
            className="text-primary cursor-pointer select-none"
          >
            {label}
          </label>

          {helperText && (
            <p className="text-sm text-tertiary mt-1">{helperText}</p>
          )}
        </div>
      </div>
    );
  }
);

Checkbox.displayName = 'Checkbox';

// ============================================
// RADIO GROUP
// ============================================

interface RadioOption {
  value: string;
  label: string;
  disabled?: boolean;
}

interface RadioGroupProps {
  /** Group label (required for accessibility) */
  label: string;
  /** Radio options */
  options: RadioOption[];
  /** Selected value */
  value?: string;
  /** Change handler */
  onChange?: (value: string) => void;
  /** Group name (required for grouping radios) */
  name: string;
  /** Error state */
  error?: boolean;
  /** Helper text */
  helperText?: string;
  /** Layout direction */
  direction?: 'horizontal' | 'vertical';
  /** Additional className */
  className?: string;
}

export function RadioGroup({
  label,
  options,
  value,
  onChange,
  name,
  error,
  helperText,
  direction = 'vertical',
  className,
}: RadioGroupProps) {
  const groupId = useId();

  return (
    <fieldset className={className}>
      <legend className="font-semibold text-sm mb-3 text-primary">
        {label}
      </legend>

      <div
        role="radiogroup"
        aria-labelledby={groupId}
        className={cn(
          'flex gap-4',
          direction === 'vertical' && 'flex-col'
        )}
      >
        {options.map((option) => (
          <label
            key={option.value}
            className={cn(
              'flex items-center gap-3 cursor-pointer',
              option.disabled && 'opacity-50 cursor-not-allowed'
            )}
          >
            <input
              type="radio"
              name={name}
              value={option.value}
              checked={value === option.value}
              onChange={() => onChange?.(option.value)}
              disabled={option.disabled}
              className={cn(
                'w-5 h-5 shrink-0',
                'border-2 border-default rounded-full',
                'text-brand-primary',
                'focus:ring-2 focus:ring-brand-primary focus:ring-offset-2',
                error && 'border-error'
              )}
            />
            <span className="text-primary">{option.label}</span>
          </label>
        ))}
      </div>

      {helperText && (
        <p className="mt-2 text-sm text-tertiary">{helperText}</p>
      )}
    </fieldset>
  );
}

// ============================================
// FORM GROUP
// ============================================

interface FormGroupProps {
  children: ReactNode;
  /** Group label */
  label?: string;
  /** Additional className */
  className?: string;
}

export function FormGroup({ children, label, className }: FormGroupProps) {
  return (
    <div className={cn('space-y-4', className)}>
      {label && (
        <h3 className="font-semibold text-lg text-primary">{label}</h3>
      )}
      {children}
    </div>
  );
}

// ============================================
// FORM ROW (for horizontal layouts)
// ============================================

interface FormRowProps {
  children: ReactNode;
  /** Gap between items */
  gap?: 'sm' | 'md' | 'lg';
  /** Additional className */
  className?: string;
}

export function FormRow({ children, gap = 'md', className }: FormRowProps) {
  const gapClasses = {
    sm: 'gap-3',
    md: 'gap-4',
    lg: 'gap-6',
  };

  return (
    <div className={cn('flex flex-col sm:flex-row', gapClasses[gap], className)}>
      {children}
    </div>
  );
}
