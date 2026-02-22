# CopilotKit Agent Integration Patterns

This reference provides detailed patterns for integrating AI agents with CopilotKit.

## Agent Configuration

### Basic Setup

```typescript
// copilotkit.config.ts
import { CopilotKit } from "@copilotkit/react-core";

export const copilotConfig = {
  // Backend endpoint for agent communication
  agent: "medisync-agent",

  // Default instructions for the agent
  instructions: `You are a helpful AI assistant for MediSync,
  a healthcare business intelligence platform. You can help users
  query data, run reports, and manage accounting tasks.`,

  // Model configuration
  model: {
    provider: "openai",
    model: "gpt-4-turbo"
  },

  // Context providers
  contextProviders: [
    "user-preferences",
    "dashboard-state",
    "active-company"
  ]
};
```

### Provider Wrapper

```typescript
// App.tsx
import { CopilotKit } from "@copilotkit/react-core";
import { CopilotPopup } from "@copilotkit/react-ui";
import { copilotConfig } from "./copilotkit.config";

function App() {
  return (
    <CopilotKit
      agent={copilotConfig.agent}
      instructions={copilotConfig.instructions}
    >
      <Dashboard />
      <CopilotPopup />
    </CopilotKit>
  );
}
```

## Tool Definition Patterns

### Query Tool (Text-to-SQL)

```typescript
import { tool } from "@copilotkit/react-core";

export const textToSQLTool = tool({
  name: "run_sql_query",
  description: `Execute a read-only SQL query against the data warehouse.
  Use this to retrieve business data when users ask questions about
  revenue, patients, appointments, or any operational metrics.`,

  parameters: {
    type: "object",
    properties: {
      sql: {
        type: "string",
        description: "SELECT-only SQL query"
      },
      explanation: {
        type: "string",
        description: "Human-readable explanation of the query"
      }
    },
    required: ["sql"]
  },

  handler: async ({ sql, explanation }, { user }) => {
    // Validate SQL is read-only
    if (!isSelectOnly(sql)) {
      throw new Error("Only SELECT queries are allowed");
    }

    // Execute query
    const results = await warehouse.query(sql);

    return {
      sql,
      explanation,
      results,
      rowCount: results.length,
      executedAt: new Date().toISOString()
    };
  },

  // UI component to render results
  render: ({ results, explanation }) => (
    <QueryResultPreview
      data={results}
      explanation={explanation}
    />
  )
});
```

### Report Generation Tool

```typescript
export const generateReportTool = tool({
  name: "generate_report",
  description: `Generate a formatted business report.
  Use when users ask for summaries, comparisons, or periodic reports.`,

  parameters: {
    type: "object",
    properties: {
      reportType: {
        type: "string",
        enum: ["revenue", "patients", "appointments", "inventory"]
      },
      period: {
        type: "string",
        enum: ["daily", "weekly", "monthly", "quarterly", "yearly"]
      },
      format: {
        type: "string",
        enum: ["table", "chart", "pdf"],
        default: "table"
      }
    },
    required: ["reportType", "period"]
  },

  handler: async ({ reportType, period, format }) => {
    const report = await reportService.generate({
      type: reportType,
      period,
      format
    });

    return report;
  },

  render: (report) => {
    switch (report.format) {
      case "chart":
        return <ChartReport data={report} />;
      case "pdf":
        return <PDFDownloadLink report={report} />;
      default:
        return <TableReport data={report} />;
    }
  }
});
```

### Tally Sync Tool (with HITL)

```typescript
export const tallySyncTool = tool({
  name: "sync_to_tally",
  description: `Sync approved transactions to Tally ERP.
  Requires finance head approval before execution.
  Use when users want to push journal entries to accounting.`,

  parameters: {
    type: "object",
    properties: {
      entryIds: {
        type: "array",
        items: { type: "string" },
        description: "List of journal entry IDs to sync"
      }
    },
    required: ["entryIds"]
  },

  handler: async ({ entryIds }, context) => {
    // Check approval status
    const entries = await journalRepo.findByIds(entryIds);

    const unapproved = entries.filter(e => !e.isApproved);
    if (unapproved.length > 0) {
      return {
        success: false,
        requiresApproval: true,
        pendingEntries: unapproved.map(e => e.id)
      };
    }

    // Request human confirmation
    const confirmed = await context.requestConfirmation({
      title: "Confirm Tally Sync",
      message: `Sync ${entries.length} journal entries to Tally?`,
      details: entries.map(e => ({
        date: e.date,
        amount: e.amount,
        ledger: e.ledgerName
      }))
    });

    if (!confirmed) {
      return { success: false, reason: "User cancelled" };
    }

    // Execute sync
    const result = await tallyService.sync(entries);

    return {
      success: true,
      syncedCount: result.synced.length,
      failedCount: result.failed.length,
      failures: result.failed
    };
  },

  render: (result) => (
    <TallySyncResult status={result} />
  )
});
```

## State Management with useAgent

### Reading Agent State

```typescript
import { useAgent } from "@copilotkit/react-core";

function QueryBuilder() {
  const { agent } = useAgent({ agentId: "medisync-agent" });

  // Access current query state
  const currentQuery = agent.state.currentQuery;
  const queryHistory = agent.state.queryHistory || [];

  return (
    <div>
      <h3>Current Query</h3>
      <pre>{currentQuery?.sql}</pre>

      <h3>History</h3>
      <ul>
        {queryHistory.map((q, i) => (
          <li key={i}>{q.explanation}</li>
        ))}
      </ul>
    </div>
  );
}
```

### Updating Agent State

```typescript
function DashboardFilters() {
  const { agent } = useAgent({ agentId: "medisync-agent" });

  const handleFilterChange = (filters) => {
    // Update agent state - this is visible to the agent
    agent.setState({
      filters,
      lastUpdated: Date.now()
    });
  };

  return (
    <FilterPanel
      currentFilters={agent.state.filters || {}}
      onChange={handleFilterChange}
    />
  );
}
```

### Two-Way State Sync

```typescript
function SharedStateExample() {
  const { agent } = useAgent({ agentId: "medisync-agent" });

  // Local state synced with agent
  const [localState, setLocalState] = useState(agent.state);

  // Sync local changes to agent
  const updateSharedState = (updates) => {
    const newState = { ...localState, ...updates };
    setLocalState(newState);
    agent.setState(newState);
  };

  // Listen for agent updates
  useEffect(() => {
    const unsubscribe = agent.onStateChange((newState) => {
      setLocalState(newState);
    });
    return unsubscribe;
  }, [agent]);

  return (
    <div>
      <input
        value={localState.searchQuery || ""}
        onChange={(e) => updateSharedState({ searchQuery: e.target.value })}
      />
      <button onClick={() => agent.run({ action: "search" })}>
        Search
      </button>
    </div>
  );
}
```

## Human-in-the-Loop Patterns

### Confirmation Dialog

```typescript
// In tool handler
const confirmed = await context.requestConfirmation({
  title: "Confirm Action",
  message: "Are you sure?",
  type: "warning", // 'info' | 'warning' | 'danger'
  confirmText: "Yes, proceed",
  cancelText: "Cancel"
});
```

### Form Input Request

```typescript
// Agent requests user to fill a form
const formData = await context.requestInput({
  title: "Complete the Entry",
  fields: [
    { name: "date", type: "date", label: "Transaction Date", required: true },
    { name: "amount", type: "number", label: "Amount", required: true },
    { name: "description", type: "text", label: "Description" }
  ]
});
```

### Selection Choice

```typescript
// Agent presents options
const choice = await context.requestSelection({
  title: "Choose an Action",
  message: "How would you like to proceed?",
  options: [
    { id: "save", label: "Save Draft", icon: "save" },
    { id: "submit", label: "Submit for Approval", icon: "send" },
    { id: "discard", label: "Discard Changes", icon: "delete", destructive: true }
  ]
});
```

## Error Handling

### Graceful Degradation

```typescript
handler: async (params, context) => {
  try {
    const result = await riskyOperation(params);
    return { success: true, result };
  } catch (error) {
    // Log for debugging
    console.error("Tool execution failed:", error);

    // Return user-friendly error
    return {
      success: false,
      error: error.message,
      suggestion: "Please try again or contact support if the issue persists."
    };
  }
}
```

### Fallback UI

```typescript
render: (result) => {
  if (!result.success) {
    return (
      <ErrorFallback
        error={result.error}
        suggestion={result.suggestion}
        onRetry={() => {/* trigger retry */}}
      />
    );
  }

  return <SuccessView data={result.result} />;
}
```

## Streaming Results

```typescript
handler: async function* ({ sql }, context) {
  // Yield progress updates
  yield { status: "validating", message: "Validating query..." };

  await validateSQL(sql);
  yield { status: "executing", message: "Running query..." };

  // Stream results
  for await (const chunk of streamQueryResults(sql)) {
    yield {
      status: "streaming",
      rowsReceived: chunk.length,
      partialResults: chunk
    };
  }

  yield { status: "complete", message: "Query complete!" };
}
```
