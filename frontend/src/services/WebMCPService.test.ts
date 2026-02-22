import { describe, it, expect, vi, beforeEach } from 'vitest';
import { webMCPService } from './WebMCPService';

describe('WebMCPService', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset navigator mock structure
        Object.defineProperty(window, 'navigator', {
            value: {
                modelContext: undefined
            },
            writable: true
        });
        // Reset service state
        (webMCPService as unknown as { registeredTools: Set<string> }).registeredTools.clear();
    });

    it('should return false for isSupported when navigator.modelContext is missing', () => {
        expect(webMCPService.isSupported()).toBe(false);
    });

    it('should return true for isSupported when navigator.modelContext is present', () => {
        const nav = window.navigator as typeof window.navigator & { modelContext: unknown };
        nav.modelContext = {
            registerTool: vi.fn(),
            unregisterTool: vi.fn()
        };
        expect(webMCPService.isSupported()).toBe(true);
    });

    it('should register chat tools correctly when supported', () => {
        const mockRegister = vi.fn();
        const mockUnregister = vi.fn();
        const nav = window.navigator as typeof window.navigator & { modelContext: unknown };
        nav.modelContext = {
            registerTool: mockRegister,
            unregisterTool: mockUnregister
        };

        const callbacks = {
            onQuery: vi.fn(),
            onSyncTally: vi.fn(),
            onShowDashboard: vi.fn()
        };

        webMCPService.registerChatTools(callbacks);

        // Should have registered 3 tools (queryBI, syncTally, showDashboard)
        expect(mockRegister).toHaveBeenCalledTimes(3);

        // Check first tool (chat_queryBI)
        expect(mockRegister).toHaveBeenCalledWith(expect.objectContaining({
            name: 'chat_queryBI'
        }));
    });

    it('should register dashboard tools correctly when supported', () => {
        const mockRegister = vi.fn();
        const nav = window.navigator as typeof window.navigator & { modelContext: unknown };
        nav.modelContext = {
            registerTool: mockRegister,
            unregisterTool: vi.fn()
        };

        webMCPService.registerDashboardTools({
            onRefreshChart: vi.fn(),
            onPinChart: mockRegister,
            onNavigateToChat: mockRegister,
            onRefreshAll: mockRegister,
            onNavigateToSettings: mockRegister,
            onOpenDashboardModal: mockRegister
        });

        // Should have registered 5 tools
        expect(mockRegister).toHaveBeenCalledTimes(5);
    });

    it('should register navigation tools correctly when supported', () => {
        const mockRegister = vi.fn();
        const nav = window.navigator as typeof window.navigator & { modelContext: unknown };
        nav.modelContext = {
            registerTool: mockRegister,
            unregisterTool: vi.fn()
        };

        webMCPService.registerNavigationTools({
            onNavigate: mockRegister,
            onToggleLanguage: mockRegister,
            onShowRecommendations: mockRegister
        });

        // Should have registered 3 tools (navigate, toggleLanguage, showRecommendations)
        expect(mockRegister).toHaveBeenCalledTimes(3);
    });

    it('should cleanup tools correctly', () => {
        const mockUnregister = vi.fn();
        const nav = window.navigator as typeof window.navigator & { modelContext: unknown };
        nav.modelContext = {
            registerTool: vi.fn(),
            unregisterTool: mockUnregister
        };

        // Pre-register some tools
        webMCPService.registerChatTools({
            onQuery: vi.fn(),
            onSyncTally: vi.fn(),
            onShowDashboard: vi.fn()
        });

        webMCPService.cleanup();

        // Should have unregistered 3 tools
        expect(mockUnregister).toHaveBeenCalledTimes(3);
    });
});
