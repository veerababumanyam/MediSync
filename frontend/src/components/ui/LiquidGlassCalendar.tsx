import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/cn';

export interface LiquidGlassCalendarProps {
    selectedDate?: Date;
    onDateSelect: (date: Date) => void;
    minDate?: Date;
    maxDate?: Date;
    className?: string;
    classNameLabel?: string;
}

export const LiquidGlassCalendar: React.FC<LiquidGlassCalendarProps> = ({
    selectedDate,
    onDateSelect,
    minDate,
    maxDate,
    className,
}) => {
    const [viewDate, setViewDate] = useState(selectedDate || new Date());

    const currentMonth = viewDate.getMonth();
    const currentYear = viewDate.getFullYear();

    const daysInMonth = new Date(currentYear, currentMonth + 1, 0).getDate();
    const firstDayOfMonth = new Date(currentYear, currentMonth, 1).getDay();

    const days = Array.from({ length: daysInMonth }, (_, i) => i + 1);
    const paddingBefore = Array.from({ length: firstDayOfMonth }, (_, i) => i);

    const prevMonth = () => setViewDate(new Date(currentYear, currentMonth - 1, 1));
    const nextMonth = () => setViewDate(new Date(currentYear, currentMonth + 1, 1));

    const isToday = (day: number) => {
        const today = new Date();
        return (
            today.getDate() === day &&
            today.getMonth() === currentMonth &&
            today.getFullYear() === currentYear
        );
    };

    const isSelected = (day: number) => {
        return (
            selectedDate?.getDate() === day &&
            selectedDate?.getMonth() === currentMonth &&
            selectedDate?.getFullYear() === currentYear
        );
    };

    const isDisabled = (day: number) => {
        const date = new Date(currentYear, currentMonth, day);
        if (minDate && date < minDate) return true;
        if (maxDate && date > maxDate) return true;
        return false;
    };

    const monthNames = [
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
    ];

    const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

    return (
        <div className={cn(
            'w-full max-w-sm p-5',
            'bg-surface-glass-heavy backdrop-blur-xl border border-glass-soft rounded-radius-card',
            'shadow-glass-floating shadow-2xl animate-liquid-fade-in',
            className
        )}>
            {/* Calendar Header */}
            <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-bold text-primary flex items-center gap-2">
                    <CalendarIcon className="w-5 h-5 text-brand-primary" />
                    {monthNames[currentMonth]} {currentYear}
                </h3>
                <div className="flex items-center gap-1">
                    <button
                        onClick={prevMonth}
                        className="p-1.5 rounded-full hover:bg-brand-primary/10 text-secondary hover:text-brand-primary transition-all active:scale-90"
                        aria-label="Previous month"
                    >
                        <ChevronLeft className="w-5 h-5" />
                    </button>
                    <button
                        onClick={nextMonth}
                        className="p-1.5 rounded-full hover:bg-brand-primary/10 text-secondary hover:text-brand-primary transition-all active:scale-90"
                        aria-label="Next month"
                    >
                        <ChevronRight className="w-5 h-5" />
                    </button>
                </div>
            </div>

            {/* Day Labels */}
            <div className="grid grid-cols-7 mb-2">
                {dayLabels.map((label) => (
                    <div key={label} className="text-center text-[10px] font-bold uppercase tracking-widest text-tertiary">
                        {label}
                    </div>
                ))}
            </div>

            {/* Days Grid */}
            <div className="grid grid-cols-7 gap-1">
                {paddingBefore.map((i) => (
                    <div key={`pad-${i}`} className="aspect-square" />
                ))}
                {days.map((day) => {
                    const disabled = isDisabled(day);
                    const selected = isSelected(day);
                    const today = isToday(day);

                    return (
                        <button
                            key={day}
                            onClick={() => !disabled && onDateSelect(new Date(currentYear, currentMonth, day))}
                            disabled={disabled}
                            className={cn(
                                'relative aspect-square rounded-radius-md flex items-center justify-center text-sm transition-all duration-300',
                                'hover:bg-brand-primary/10 group',
                                today && 'font-bold text-brand-secondary',
                                selected && 'bg-brand-primary text-white shadow-lg shadow-brand-primary/30 font-bold',
                                disabled && 'opacity-20 cursor-not-allowed grayscale',
                                !selected && !disabled && 'text-secondary group-hover:text-brand-primary'
                            )}
                        >
                            <span className="relative z-10">{day}</span>
                            {today && !selected && (
                                <div className="absolute bottom-1.5 w-1 h-1 rounded-full bg-brand-secondary" />
                            )}
                            {selected && (
                                <motion.div
                                    layoutId="selectedDay"
                                    className="absolute inset-0 bg-brand-primary rounded-radius-md shadow-lg"
                                    transition={{ type: 'spring', bounce: 0.3, duration: 0.5 }}
                                />
                            )}
                        </button>
                    );
                })}
            </div>

            {/* Legend / Footer */}
            <div className="mt-6 pt-4 border-t border-glass-dim flex items-center justify-between text-[11px] text-tertiary">
                <button
                    onClick={() => {
                        const today = new Date();
                        setViewDate(today);
                        onDateSelect(today);
                    }}
                    className="hover:text-brand-primary transition-colors font-medium"
                >
                    Today
                </button>
                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-1">
                        <div className="w-1.5 h-1.5 rounded-full bg-brand-primary" />
                        Selected
                    </div>
                    <div className="flex items-center gap-1">
                        <div className="w-1.5 h-1.5 rounded-full bg-brand-secondary" />
                        Today
                    </div>
                </div>
            </div>
        </div>
    );
};
