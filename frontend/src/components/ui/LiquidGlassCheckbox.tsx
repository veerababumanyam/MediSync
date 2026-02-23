import React from 'react';
import { motion } from 'framer-motion';
import { Check } from 'lucide-react';
import { cn } from '@/lib/cn';

export interface LiquidGlassCheckboxProps {
    checked: boolean;
    onChange: (checked: boolean) => void;
    disabled?: boolean;
    label?: string;
    className?: string;
}

export const LiquidGlassCheckbox: React.FC<LiquidGlassCheckboxProps> = ({
    checked,
    onChange,
    disabled = false,
    label,
    className,
}) => {
    return (
        <label className={cn(
            'flex items-center gap-3 cursor-pointer select-none group',
            disabled && 'opacity-50 cursor-not-allowed grayscale',
            className
        )}>
            <div className="relative">
                <input
                    type="checkbox"
                    className="sr-only"
                    checked={checked}
                    onChange={(e) => !disabled && onChange(e.target.checked)}
                    disabled={disabled}
                />
                <motion.div
                    animate={{
                        backgroundColor: checked ? 'var(--brand-primary)' : 'var(--surface-glass-light)',
                        borderColor: checked ? 'var(--brand-primary)' : 'var(--glass-border-soft)'
                    }}
                    className={cn(
                        'w-6 h-6 rounded-radius-xs border backdrop-blur-md transition-all duration-300',
                        'flex items-center justify-center shadow-ambient group-hover:shadow-raised',
                        checked && 'shadow-glow-primary'
                    )}
                >
                    <AnimatePresence>
                        {checked && (
                            <motion.div
                                initial={{ scale: 0, opacity: 0 }}
                                animate={{ scale: 1, opacity: 1 }}
                                exit={{ scale: 0, opacity: 0 }}
                                transition={{ type: 'spring', damping: 20, stiffness: 300 }}
                            >
                                <Check className="w-4 h-4 text-white stroke-[3px]" />
                            </motion.div>
                        )}
                    </AnimatePresence>
                </motion.div>
            </div>
            {label && (
                <span className="text-sm font-medium text-secondary group-hover:text-primary transition-colors">
                    {label}
                </span>
            )}
        </label>
    );
};

// Internal AnimatePresence wrapper for ease of use
import { AnimatePresence } from 'framer-motion';
