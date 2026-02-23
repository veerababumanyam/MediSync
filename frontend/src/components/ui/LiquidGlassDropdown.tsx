import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronDown, Check, Search, X } from 'lucide-react';
import { cn } from '@/lib/cn';

export interface DropdownOption {
    value: string;
    label: string;
    icon?: React.ReactNode;
    disabled?: boolean;
}

export interface LiquidGlassDropdownProps {
    options: DropdownOption[];
    value?: string;
    onChange: (value: string) => void;
    placeholder?: string;
    label?: string;
    error?: string;
    disabled?: boolean;
    searchable?: boolean;
    className?: string;
    containerClassName?: string;
}

export const LiquidGlassDropdown: React.FC<LiquidGlassDropdownProps> = ({
    options,
    value,
    onChange,
    placeholder = 'Select an option',
    label,
    error,
    disabled = false,
    searchable = false,
    className,
    containerClassName,
}) => {
    const [isOpen, setIsOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const containerRef = useRef<HTMLDivElement>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const selectedOption = options.find((opt) => opt.value === value);

    const filteredOptions = options.filter((opt) =>
        opt.label.toLowerCase().includes(searchQuery.toLowerCase())
    );

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [isOpen]);

    const handleToggle = () => {
        if (!disabled) {
            setIsOpen(!isOpen);
            if (!isOpen) setSearchQuery('');
        }
    };

    const handleSelect = (option: DropdownOption) => {
        if (!option.disabled) {
            onChange(option.value);
            setIsOpen(false);
        }
    };

    return (
        <div className={cn('flex flex-col gap-1.5 w-full', containerClassName)} ref={containerRef}>
            {label && (
                <label className="text-sm font-medium text-secondary px-1">
                    {label}
                </label>
            )}

            <div className="relative">
                <button
                    type="button"
                    onClick={handleToggle}
                    disabled={disabled}
                    className={cn(
                        'flex items-center justify-between w-full px-4 py-3 text-start transition-all duration-300',
                        'bg-surface-glass-light backdrop-blur-md border border-glass-soft rounded-radius-input',
                        'shadow-ambient hover:shadow-raised hover:border-glass-bright',
                        'focus:outline-none focus:ring-2 focus:ring-brand-primary/40 focus:border-brand-primary/60',
                        disabled && 'opacity-50 cursor-not-allowed grayscale',
                        error && 'border-error ring-error/20',
                        className
                    )}
                    aria-haspopup="listbox"
                    aria-expanded={isOpen}
                >
                    <div className="flex items-center gap-3 overflow-hidden">
                        {selectedOption?.icon && (
                            <span className="shrink-0 text-brand-primary">
                                {selectedOption.icon}
                            </span>
                        )}
                        <span className={cn(
                            'truncate',
                            !selectedOption ? 'text-tertiary' : 'text-primary font-medium'
                        )}>
                            {selectedOption ? selectedOption.label : placeholder}
                        </span>
                    </div>
                    <ChevronDown className={cn(
                        'w-4 h-4 text-tertiary transition-transform duration-300',
                        isOpen && 'rotate-180 text-brand-primary'
                    )} />
                </button>

                <AnimatePresence>
                    {isOpen && (
                        <motion.div
                            ref={dropdownRef}
                            initial={{ opacity: 0, y: -10, scale: 0.95 }}
                            animate={{ opacity: 1, y: 4, scale: 1 }}
                            exit={{ opacity: 0, y: -10, scale: 0.95 }}
                            transition={{ duration: 0.2, ease: 'easeOut' }}
                            className={cn(
                                'absolute z-dropdown w-full mt-1 overflow-hidden',
                                'bg-surface-glass-heavy backdrop-blur-xl border border-glass-soft rounded-radius-card',
                                'shadow-glass-floating shadow-2xl animate-liquid-fade-in'
                            )}
                        >
                            {searchable && (
                                <div className="p-2 border-b border-glass-dim sticky top-0 bg-surface-glass-heavy z-10">
                                    <div className="relative flex items-center">
                                        <Search className="absolute left-3 w-4 h-4 text-tertiary" />
                                        <input
                                            autoFocus
                                            type="text"
                                            className="w-full bg-surface-glass-subtle border border-glass-dim rounded-lg pl-9 pr-3 py-2 text-sm text-primary focus:outline-none focus:ring-2 focus:ring-brand-primary/30"
                                            placeholder="Search..."
                                            value={searchQuery}
                                            onChange={(e) => setSearchQuery(e.target.value)}
                                        />
                                        {searchQuery && (
                                            <button
                                                onClick={() => setSearchQuery('')}
                                                className="absolute right-3 p-1 rounded-full hover:bg-white/10"
                                            >
                                                <X className="w-3 h-3 text-tertiary" />
                                            </button>
                                        )}
                                    </div>
                                </div>
                            )}

                            <ul
                                className="max-h-60 overflow-y-auto liquid-glass-scroll p-1.5"
                                role="listbox"
                            >
                                {filteredOptions.length > 0 ? (
                                    filteredOptions.map((option) => (
                                        <li key={option.value} role="option" aria-selected={value === option.value}>
                                            <button
                                                type="button"
                                                onClick={() => handleSelect(option)}
                                                disabled={option.disabled}
                                                className={cn(
                                                    'flex items-center justify-between w-full px-3 py-2.5 rounded-radius-md transition-all duration-200',
                                                    'hover:bg-brand-primary/10 group',
                                                    value === option.value ? 'bg-brand-primary/20' : 'bg-transparent',
                                                    option.disabled && 'opacity-40 cursor-not-allowed grayscale'
                                                )}
                                            >
                                                <div className="flex items-center gap-3 overflow-hidden">
                                                    {option.icon && (
                                                        <span className={cn(
                                                            'shrink-0 transition-colors',
                                                            value === option.value ? 'text-brand-primary' : 'text-tertiary group-hover:text-brand-primary'
                                                        )}>
                                                            {option.icon}
                                                        </span>
                                                    )}
                                                    <span className={cn(
                                                        'text-sm truncate',
                                                        value === option.value ? 'text-brand-primary font-semibold' : 'text-secondary group-hover:text-primary'
                                                    )}>
                                                        {option.label}
                                                    </span>
                                                </div>
                                                {value === option.value && (
                                                    <Check className="w-4 h-4 text-brand-primary shrink-0" />
                                                )}
                                            </button>
                                        </li>
                                    ))
                                ) : (
                                    <li className="px-4 py-8 text-center text-sm text-tertiary">
                                        No options found
                                    </li>
                                )}
                            </ul>
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>

            {error && (
                <p className="text-xs text-error font-medium px-1 flex items-center gap-1">
                    <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {error}
                </p>
            )}
        </div>
    );
};
