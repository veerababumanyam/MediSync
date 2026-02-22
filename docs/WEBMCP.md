# WebMCP (Web Machine Context & Procedure)

MediSync is integrated with the **WebMCP** standard, allowing browser-based AI agents to discover and interact with the platform's features as structured tools. This enhances the "agent-readiness" of the MediSync UI.

## Overview

WebMCP provides two primary ways for AI agents to interact with MediSync:
1.  **Declarative API**: Standard HTML attributes (`tool-name`, `tool-description`) that allow browsers to discover interactive elements.
2.  **Imperative API**: A JavaScript service (`navigator.modelContext`) that allows MediSync to register complex business logic as callable tools for agents.

## Implementation Details

### WebMCP Service
The `WebMCPService` handles the registration of core MediSync actions. It is initialized in the `ChatInterface` component.

**Registered Tools:**
*   `queryBI`: Execute a natural language query against MediSync BI data.
*   `syncTally`: Trigger a manual synchronization with Tally ERP.
*   `showDashboard`: Navigate to a specific MediSync dashboard.

### Declarative Tags
Interactive components, such as the chat input area, are tagged with WebMCP attributes for automatic discovery:

```html
<div tool-name="medi-chat-input-area" tool-description="The main interaction area for sending queries to MediSync AI">
  ...
</div>
```

## How to Test

WebMCP is currently an experimental standard and requires specific browser support.

1.  **Browser**: Use Chrome 146+ or Chrome Canary.
2.  **Enable Flag**: Navigate to `chrome://flags/#web-mcp` and set it to **Enabled**.
3.  **Discovery**: Use a WebMCP-capable agent (like the Chrome experimental AI panel or a supported extension) to "see" the tools exposed by MediSync.

## Developer Quick Start

To add new WebMCP tools:
1.  Open `frontend/src/services/WebMCPService.ts`.
2.  Add a new tool definition in `registerMediSyncTools`.
3.  Define the `parameters` schema and the `handler` callback.

To add declarative tags:
1.  Add `tool-name` and `tool-description` attributes to any interactive HTML element.
2.  In React/JSX, use the prop spread or `@ts-ignore` for custom attributes if necessary.

## Reference
- [WebMCP Explainer](https://github.com/web-mcp/explainer)
- [MediSync AI Agent Ecosystem](../README.md#ai-agent-ecosystem)
