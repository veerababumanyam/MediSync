---
name: copilotkit-generative-ui
description: This skill should be used when the user asks to "build agent-native applications", "create generative UI", "implement CopilotKit", "add AI chat to React app", "useAgent hook", "AG-UI protocol", "human-in-the-loop workflows", "frontend", "UI", "UX", "UI/UX" or mentions CopilotKit features like chat interfaces, tool rendering, or shared state between agents and UI.
---

# CopilotKit Generative UI

CopilotKit is an SDK for building agent-native applications with generative UI, shared state, and human-in-the-loop workflows. It enables AI agents to dynamically render UI components, call tools, and maintain synchronized state with the frontend.

★ Insight ─────────────────────────────────────
CopilotKit's three core capabilities:
1. **Generative UI** - Agents render React/Angular components at runtime
2. **Backend Tool Rendering** - Tools return UI that renders client-side
3. **Shared State** - Real-time sync between agent and UI via `useAgent` hook

This enables truly collaborative human-AI workflows.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Purpose** | Build agent-native apps with generative UI |
| **License** | MIT |
| **Frameworks** | React, Angular |
| **Key Features** | Chat UI, Tool Rendering, Generative UI, Shared State, HITL |
| **Protocol** | AG-UI (Agent-User Interaction Protocol) |

## Installation

### New Projects
```bash
npx copilotkit@latest create -f <framework>
```

### Existing Projects
```bash
npx copilotkit@latest init
```

This provides:
- Core packages configured
- Provider with context and hooks
- Agent-to-UI connection ready
- Deployment-ready setup

## Core Concepts

### 1. Chat UI
React-based chat interface supporting:
- Message streaming
- Tool calls
- Agent responses
- Markdown rendering

### 2. Backend Tool Rendering
Agents call backend tools that return UI components:
```typescript
// Tool returns a React component
const myTool = {
  name: "show_chart",
  description: "Display a chart with data",
  parameters: { ... },
  render: (result) => <ChartComponent data={result} />
};
```

### 3. Generative UI Types

| Type | Description | Use Case |
|------|-------------|----------|
| **Static (AG-UI)** | Pre-defined components rendered by agent | Known UI patterns |
| **Declarative (A2UI)** | Agent declares intent, UI interprets | Flexible layouts |
| **Open-Ended** | Agent generates arbitrary JSON/structures | Custom experiences |

### 4. Shared State
Real-time synchronized state between agent and UI:
```typescript
const { agent } = useAgent({ agentId: "my_agent" });

// Read state
const city = agent.state.city;

// Update state
agent.setState({ city: "NYC" });
```

### 5. Human-in-the-Loop
Agents can pause and request user input:
- Confirmation dialogs
- Form inputs
- Selection choices
- Edit approvals

## The useAgent Hook

The `useAgent` hook provides programmatic control over agents:

```typescript
import { useAgent } from "@copilotkit/react-core";

function MyComponent() {
  const { agent } = useAgent({ agentId: "assistant" });

  // Access agent state
  const isLoading = agent.isLoading;
  const state = agent.state;

  // Control agent
  const handleSubmit = () => {
    agent.run({ prompt: userInput });
  };

  // Update state
  const updateContext = () => {
    agent.setState({ context: newContext });
  };

  return (
    <div>
      <h1>{agent.state.title}</h1>
      <button onClick={() => agent.setState({ title: "New Title" })}>
        Update
      </button>
    </div>
  );
}
```

## Architecture Flow

```
User Input → Agent Processing → Tool Calls → UI Rendering
     ↑                                              ↓
     └──────────── Shared State Sync ←─────────────┘
```

CopilotKit connects:
1. **UI** - React/Angular components
2. **Agents** - Backend AI logic
3. **Tools** - Actionable functions
4. **State** - Synchronized data layer

## Common Patterns

### Pattern 1: Chat with Tool Rendering

```typescript
// Define tool with UI rendering
const tools = [
  {
    name: "search_products",
    description: "Search product catalog",
    parameters: {
      type: "object",
      properties: {
        query: { type: "string" }
      }
    },
    handler: async ({ query }) => {
      const results = await searchAPI(query);
      return results;
    },
    render: (results) => <ProductList products={results} />
  }
];
```

### Pattern 2: Human-in-the-Loop Confirmation

```typescript
const confirmTool = {
  name: "delete_item",
  description: "Delete an item after confirmation",
  handler: async (params) => {
    // Agent pauses here, waits for user
    const confirmed = await waitForConfirmation({
      message: `Delete ${params.itemName}?`,
      options: ["Confirm", "Cancel"]
    });

    if (confirmed) {
      await deleteAPI(params.id);
    }
  }
};
```

### Pattern 3: Shared State Workflow

```typescript
function WorkflowUI() {
  const { agent } = useAgent({ agentId: "workflow" });

  return (
    <div>
      {/* UI reads agent state */}
      <ProgressBar value={agent.state.progress} />
      <StepList steps={agent.state.steps} />

      {/* User actions update agent state */}
      <button onClick={() => agent.setState({ paused: true })}>
        Pause
      </button>
    </div>
  );
}
```

## Integration with MediSync

CopilotKit is the frontend layer for MediSync's AI agents:

```typescript
// In MediSync frontend
import { CopilotKit } from "@copilotkit/react-core";
import { CopilotPopup } from "@copilotkit/react-ui";

function App() {
  return (
    <CopilotKit agent="medisync-agent">
      <Dashboard />
      <CopilotPopup
        instructions="You are MediSync AI assistant..."
        tools={[textToSQL, runReport, syncToTally]}
      />
    </CopilotKit>
  );
}
```

### MediSync-Specific Tools

| Tool | Description | UI Render |
|------|-------------|-----------|
| `textToSQL` | Convert NL to SQL | QueryPreview |
| `runReport` | Generate reports | ReportViewer |
| `syncToTally` | Push to Tally ERP | SyncConfirmation |
| `ocrExtract` | Extract document data | DocumentPreview |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Agent not responding | Check agent registration, verify WebSocket connection |
| UI not rendering | Ensure tool has `render` function, check component imports |
| State not syncing | Verify `useAgent` hook usage, check state structure |
| Tool calls failing | Validate tool schema, check backend handler |

## Best Practices

1. **Keep tools focused** - Single responsibility per tool
2. **Render progressively** - Stream UI updates as data arrives
3. **Handle errors gracefully** - Provide fallback UI for failures
4. **Use TypeScript** - Type-safe tool definitions and state
5. **Implement HITL gates** - User confirmation for destructive actions

## Key Resources

| Resource | URL |
|----------|-----|
| Documentation | https://docs.copilotkit.ai |
| GitHub | https://github.com/CopilotKit/CopilotKit |
| Discord | https://discord.gg/6WbDupDckr |
| AG-UI Protocol | https://github.com/CopilotKit/AG-UI |

## Related Skills

- **webmcp** - For client-side MCP tools
- **medisync-dev** - For MediSync-specific development
- **tally-integration** - For Tally ERP integration
