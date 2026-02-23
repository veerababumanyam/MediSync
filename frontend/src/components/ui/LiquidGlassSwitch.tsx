import React from 'react';
import { motion } from 'framer-motion';
import { cn } from '@/lib/cn';

export interface LiquidGlassSwitchProps {
    checked: boolean;
    onChange: (checked: boolean) => void;
    disabled?: boolean;
    label?: string;
    className?: string;
}

export const LiquidGlassSwitch: React.FC<LiquidGlassSwitchProps> = ({
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
            {label && (
                <span className="text-sm font-medium text-secondary group-hover:text-primary transition-colors">
                    {label}
                </span>
            )}
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
                        backgroundColor: checked
                            ? 'var(--brand-accent)'
                            : 'var(--surface-glass-heavy)'
                    }}
                    className={cn(
                        'w-12 h-6.5 rounded-full border border-glass-soft transition-shadow duration-300',
                        'shadow-ambient',
                        checked && 'border-brand-accent/50 shadow-glow-accent'
                    )}
                >
                    <motion.div
                        animate={{
                            x: checked ? 22 : 2,
                            backgroundColor: checked ? '#FFFFFF' : 'var(--text-tertiary)'
                        }}
                        transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                        className="absolute top-1 left-1 w-4.5 h-4.5 rounded-full shadow-md flex items-center justify-center"
                    >
                        {checked && (
                            <motion.div
                                initial={{ scale: 0 }}
                                animate={{ scale: 1 }}
                                className="w-1 h-1 rounded-full bg-brand-accent"
                            />
                        )}
                    </motion.div>
                </motion.div>
            </div>
        </label>
    );
};
