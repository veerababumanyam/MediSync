const fs = require('fs');

const appPath = './src/App.tsx';
let content = fs.readFileSync(appPath, 'utf8');

// 1. Remove manual overrides from GlassCard wrappers in Live Preview cards
content = content.replace(/className="p-6 bg-white\/60 dark:bg-slate-900\/60 backdrop-blur-xl border border-slate-200 dark:border-white\/10"/g, 'className="p-6"');
content = content.replace(/className="p-6 bg-white\/20 dark:bg-white\/5 backdrop-blur-xl border border-emerald-400\/30 dark:border-emerald-400\/20"/g, 'className="p-6"');
content = content.replace(/className="p-6 bg-white\/40 dark:bg-indigo-900\/20 backdrop-blur-xl border border-indigo-200 dark:border-indigo-800\/30"/g, 'className="p-6"');

// 2. Fix CapabilityCard manual overrides
content = content.replace(/className="group p-6 rounded-2xl bg-white\/60 dark:bg-white\/5 backdrop-blur-xl border border-white\/50 dark:border-white\/10 hover:bg-white\/80 dark:hover:bg-white\/10 hover:-translate-y-1 hover:shadow-glass-lg transition-all duration-300 cursor-pointer flex flex-col items-start text-left h-full"/g, 'className="group p-8 rounded-2xl liquid-glass liquid-glass-heavy hover:-translate-y-1 transition-all duration-500 cursor-pointer flex flex-col items-start text-left h-full"');

// 3. Fix the "Everything You Need" container spacing
content = content.replace(/className="grid md:grid-cols-2 lg:grid-cols-4 gap-4"/g, 'className="grid md:grid-cols-2 lg:grid-cols-4 gap-6"');

// 4. Update the Button to have true Apple gradient and glow
content = content.replace(/bg-gradient-to-r from-logo-blue to-logo-teal text-white rounded-full shadow-lg shadow-logo-blue\/30/g, 'bg-gradient-to-b from-blue-500 to-logo-blue text-white rounded-full shadow-lg shadow-blue-500/40 border-t border-white/20');
// Also, the button font weight from font-bold to font-semibold for a more premium look
content = content.replace(/text-lg font-bold border-2 border-transparent/g, 'text-lg font-semibold border border-transparent');

// 5. Update the "View Dashboard" button for better dark mode glass contrast
content = content.replace(/bg-white\/50 dark:bg-white\/10 backdrop-blur-md border border-slate-300 dark:border-white\/20/g, 'liquid-glass-light');

// 6. Fix Live Preview Alignment ("Patient Metrics" etc)
content = content.replace(/<div className="text-xs text-slate-500 mb-1">Today<\/div>\s*<div className="text-2xl font-bold text-slate-900 dark:text-white">142<\/div>/g, '<div className="text-xs text-slate-500 dark:text-slate-400 mb-0.5">Today</div>\n                      <div className="text-xl font-bold text-slate-900 dark:text-white">142</div>');
content = content.replace(/<div className="text-xs text-slate-500 mb-1">vs Yesterday<\/div>\s*<div className="text-sm font-bold text-green-500">\+12%<\/div>/g, '<div className="text-xs text-slate-500 dark:text-slate-400 mb-0.5">vs Yesterday</div>\n                      <div className="text-sm font-bold text-emerald-500">+12%</div>');
content = content.replace(/<div className="text-xs text-slate-500 mb-1">This Month<\/div>\s*<div className="text-2xl font-bold text-slate-900 dark:text-white">3,842<\/div>/g, '<div className="text-xs text-slate-500 dark:text-slate-400 mb-0.5">This Month</div>\n                      <div className="text-xl font-bold text-slate-900 dark:text-white">3k</div>');
content = content.replace(/<div className="text-xs text-slate-500 mb-1">Departments<\/div>\s*<div className="text-sm font-bold text-blue-600 dark:text-blue-400">8 Active<\/div>/g, '<div className="text-xs text-slate-500 dark:text-slate-400 mb-0.5">Depts</div>\n                      <div className="text-sm font-bold text-logo-blue dark:text-blue-400">8 Active</div>');

fs.writeFileSync(appPath, content, 'utf8');
