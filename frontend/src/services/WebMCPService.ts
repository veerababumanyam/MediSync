/**
 * WebMCP Service
 *
 * Handles registration of MediSync as a WebMCP-enabled application.
 * WebMCP allows browser-based AI agents to discover and invoke tools
 * directly within the web application.
 *
 * @see https://github.com/webmachinelearning/webmcp
 */

// Basic types for the experimental WebMCP API
interface WebMCPTool {
    name: string;
    description: string;
    parameters: {
        type: 'object';
        properties: Record<string, any>;
        required?: string[];
    };
    handler: (args: any) => Promise<any> | any;
}

interface ModelContext {
    registerTool: (tool: WebMCPTool) => void;
    unregisterTool: (name: string) => void;
}

declare global {
    interface Navigator {
        modelContext?: ModelContext;
    }
}

/**
 * Tool categories for different pages
 */
export type ToolCategory = 'chat' | 'dashboard' | 'navigation' | 'reports' | 'alerts';

class WebMCPService {
    private registeredTools: Set<string> = new Set();

    /**
     * Check if WebMCP is supported by the browser
     */
    public isSupported(): boolean {
        return typeof window !== 'undefined' && !!navigator.modelContext;
    }

    /**
     * Get browser support status message
     */
    public getSupportMessage(): string {
        if (this.isSupported()) {
            return 'WebMCP is enabled';
        }
        return 'WebMCP requires Chrome 146+ with chrome://flags/#web-mcp enabled';
    }

    /**
     * Register core MediSync tools for Chat functionality
     */
    public registerChatTools(callbacks: {
        onQuery: (query: string) => void;
        onSyncTally: () => Promise<void>;
        onShowDashboard: (id: string) => void;
    }) {
        if (!this.isSupported()) {
            console.warn('WebMCP is not supported in this browser. Ensure you are using Chrome 146+ with #web-mcp enabled.');
            return;
        }

        this.unregisterCategory('chat');

        // 1. Tool for querying BI
        this.registerTool({
            name: 'queryBI',
            description: 'Execute a natural language query against MediSync BI data (e.g., "What was the revenue last month?")',
            parameters: {
                type: 'object',
                properties: {
                    query: {
                        type: 'string',
                        description: 'The natural language query to execute'
                    }
                },
                required: ['query']
            },
            handler: async ({ query }) => {
                callbacks.onQuery(query);
                return { success: true, message: `Query "${query}" submitted to MediSync.` };
            }
        }, 'chat');

        // 2. Tool for Tally Sync
        this.registerTool({
            name: 'syncTally',
            description: 'Trigger a manual synchronization with Tally ERP',
            parameters: {
                type: 'object',
                properties: {}
            },
            handler: async () => {
                await callbacks.onSyncTally();
                return { success: true, message: 'Tally synchronization triggered.' };
            }
        }, 'chat');

        // 3. Tool for navigating to dashboards
        this.registerTool({
            name: 'showDashboard',
            description: 'Navigate to a specific MediSync dashboard',
            parameters: {
                type: 'object',
                properties: {
                    dashboardId: {
                        type: 'string',
                        description: 'The ID of the dashboard to display'
                    }
                },
                required: ['dashboardId']
            },
            handler: async ({ dashboardId }) => {
                callbacks.onShowDashboard(dashboardId);
                return { success: true, message: `Navigating to dashboard ${dashboardId}.` };
            }
        }, 'chat');

        console.log('MediSync WebMCP Chat tools registered successfully.');
    }

    /**
     * Register tools for Dashboard functionality
     */
    public registerDashboardTools(callbacks: {
        onRefreshChart: (chartId: string) => Promise<void>;
        onPinChart: (query: string, title: string) => Promise<void>;
        onNavigateToChat: (query: string) => void;
        onRefreshAll: () => Promise<void>;
    }) {
        if (!this.isSupported()) {
            console.warn('WebMCP is not supported in this browser.');
            return;
        }

        this.unregisterCategory('dashboard');

        // 1. Tool for refreshing a specific chart
        this.registerTool({
            name: 'refreshChart',
            description: 'Refresh a specific pinned chart on the dashboard',
            parameters: {
                type: 'object',
                properties: {
                    chartId: {
                        type: 'string',
                        description: 'The ID of the chart to refresh'
                    }
                },
                required: ['chartId']
            },
            handler: async ({ chartId }) => {
                await callbacks.onRefreshChart(chartId);
                return { success: true, message: `Chart ${chartId} refreshed.` };
            }
        }, 'dashboard');

        // 2. Tool for pinning a new chart
        this.registerTool({
            name: 'pinChart',
            description: 'Pin a new chart to the dashboard from a natural language query',
            parameters: {
                type: 'object',
                properties: {
                    query: {
                        type: 'string',
                        description: 'The natural language query for the chart'
                    },
                    title: {
                        type: 'string',
                        description: 'The title for the pinned chart'
                    }
                },
                required: ['query']
            },
            handler: async ({ query, title }) => {
                await callbacks.onPinChart(query, title || 'Untitled Chart');
                return { success: true, message: `Chart "${title}" pinned to dashboard.` };
            }
        }, 'dashboard');

        // 3. Tool for navigating to chat with a query
        this.registerTool({
            name: 'exploreInChat',
            description: 'Open the chat interface with a specific query to explore data further',
            parameters: {
                type: 'object',
                properties: {
                    query: {
                        type: 'string',
                        description: 'The query to explore in chat'
                    }
                },
                required: ['query']
            },
            handler: async ({ query }) => {
                callbacks.onNavigateToChat(query);
                return { success: true, message: `Navigating to chat with query: ${query}` };
            }
        }, 'dashboard');

        // 4. Tool for refreshing all charts
        this.registerTool({
            name: 'refreshDashboard',
            description: 'Refresh all charts on the dashboard',
            parameters: {
                type: 'object',
                properties: {}
            },
            handler: async () => {
                await callbacks.onRefreshAll();
                return { success: true, message: 'All dashboard charts refreshed.' };
            }
        }, 'dashboard');

        console.log('MediSync WebMCP Dashboard tools registered successfully.');
    }

    /**
     * Register navigation tools (available on all pages)
     */
    public registerNavigationTools(callbacks: {
        onNavigate: (route: string) => void;
        onToggleLanguage: () => void;
    }) {
        if (!this.isSupported()) {
            return;
        }

        this.unregisterCategory('navigation');

        // 1. Tool for navigation
        this.registerTool({
            name: 'navigate',
            description: 'Navigate to a different page in MediSync',
            parameters: {
                type: 'object',
                properties: {
                    route: {
                        type: 'string',
                        description: 'The route to navigate to (home, chat, dashboard)',
                        enum: ['home', 'chat', 'dashboard']
                    }
                },
                required: ['route']
            },
            handler: async ({ route }) => {
                callbacks.onNavigate(route);
                return { success: true, message: `Navigating to ${route}.` };
            }
        }, 'navigation');

        // 2. Tool for language toggle
        this.registerTool({
            name: 'toggleLanguage',
            description: 'Toggle between English and Arabic (عربي) languages',
            parameters: {
                type: 'object',
                properties: {}
            },
            handler: async () => {
                callbacks.onToggleLanguage();
                return { success: true, message: 'Language toggled.' };
            }
        }, 'navigation');

        console.log('MediSync WebMCP Navigation tools registered successfully.');
    }

    /**
     * Register reports tools
     */
    public registerReportsTools(callbacks: {
        onCreateReport: (query: string, schedule: string) => Promise<void>;
        onExportReport: (format: string) => Promise<void>;
    }) {
        if (!this.isSupported()) {
            return;
        }

        this.unregisterCategory('reports');

        this.registerTool({
            name: 'createScheduledReport',
            description: 'Create a scheduled report from a natural language query',
            parameters: {
                type: 'object',
                properties: {
                    query: {
                        type: 'string',
                        description: 'The natural language query for the report'
                    },
                    schedule: {
                        type: 'string',
                        description: 'The schedule type (daily, weekly, monthly)',
                        enum: ['daily', 'weekly', 'monthly']
                    }
                },
                required: ['query']
            },
            handler: async ({ query, schedule }) => {
                await callbacks.onCreateReport(query, schedule || 'daily');
                return { success: true, message: `Scheduled report created.` };
            }
        }, 'reports');

        this.registerTool({
            name: 'exportCurrentView',
            description: 'Export the current view or report to a file',
            parameters: {
                type: 'object',
                properties: {
                    format: {
                        type: 'string',
                        description: 'The export format',
                        enum: ['pdf', 'xlsx', 'csv']
                    }
                },
                required: ['format']
            },
            handler: async ({ format }) => {
                await callbacks.onExportReport(format);
                return { success: true, message: `Exported as ${format}.` };
            }
        }, 'reports');

        console.log('MediSync WebMCP Reports tools registered successfully.');
    }

    /**
     * Register alerts tools
     */
    public registerAlertsTools(callbacks: {
        onCreateAlert: (metric: string, threshold: number, operator: string) => Promise<void>;
        onViewAlerts: () => void;
    }) {
        if (!this.isSupported()) {
            return;
        }

        this.unregisterCategory('alerts');

        this.registerTool({
            name: 'createAlert',
            description: 'Create an alert for a specific metric threshold',
            parameters: {
                type: 'object',
                properties: {
                    metric: {
                        type: 'string',
                        description: 'The metric to monitor (e.g., revenue, patient_count)'
                    },
                    threshold: {
                        type: 'number',
                        description: 'The threshold value'
                    },
                    operator: {
                        type: 'string',
                        description: 'The comparison operator',
                        enum: ['gt', 'gte', 'lt', 'lte', 'eq']
                    }
                },
                required: ['metric', 'threshold']
            },
            handler: async ({ metric, threshold, operator }) => {
                await callbacks.onCreateAlert(metric, threshold, operator || 'gt');
                return { success: true, message: `Alert created for ${metric}.` };
            }
        }, 'alerts');

        this.registerTool({
            name: 'viewAlerts',
            description: 'View all active alerts and notifications',
            parameters: {
                type: 'object',
                properties: {}
            },
            handler: async () => {
                callbacks.onViewAlerts();
                return { success: true, message: 'Viewing alerts.' };
            }
        }, 'alerts');

        console.log('MediSync WebMCP Alerts tools registered successfully.');
    }

    /**
     * Helper to register a tool with category tracking
     */
    private registerTool(tool: WebMCPTool, category: ToolCategory) {
        const fullName = `${category}_${tool.name}`;
        if (this.registeredTools.has(fullName)) {
            navigator.modelContext?.unregisterTool(fullName);
        }
        // Prefix tool name with category for uniqueness
        const prefixedTool = { ...tool, name: fullName };
        navigator.modelContext?.registerTool(prefixedTool);
        this.registeredTools.add(fullName);
    }

    /**
     * Unregister all tools in a category
     */
    private unregisterCategory(category: ToolCategory) {
        const prefix = `${category}_`;
        this.registeredTools.forEach(name => {
            if (name.startsWith(prefix)) {
                navigator.modelContext?.unregisterTool(name);
                this.registeredTools.delete(name);
            }
        });
    }

    /**
     * Unregister all tools
     */
    public cleanup() {
        if (!this.isSupported()) return;
        this.registeredTools.forEach(name => {
            navigator.modelContext?.unregisterTool(name);
        });
        this.registeredTools.clear();
    }

    /**
     * Get list of currently registered tools
     */
    public getRegisteredTools(): string[] {
        return Array.from(this.registeredTools);
    }
}

export const webMCPService = new WebMCPService();
export default webMCPService;
