// MediSync CopilotKit Integration Example
// This example shows how to integrate CopilotKit with MediSync

import React from 'react';
import { CopilotKit, useAgent, tool } from "@copilotkit/react-core";
import { CopilotPopup, CopilotSidebar } from "@copilotkit/react-ui";
import "@copilotkit/react-ui/styles.css";

// ============================================
// TOOL DEFINITIONS
// ============================================

// Tool 1: Text-to-SQL Query
const textToSQLTool = tool({
  name: "run_sql_query",
  description: `Execute a read-only SQL query against the MediSync data warehouse.
  Use this when users ask questions about business metrics, revenue, patients,
  appointments, inventory, or any operational data.

  Examples of when to use:
  - "What was our revenue last month?"
  - "How many patients visited this week?"
  - "Show me top selling products"`,

  parameters: {
    type: "object",
    properties: {
      sql: {
        type: "string",
        description: "SELECT-only SQL query"
      },
      explanation: {
        type: "string",
        description: "Plain English explanation of what the query does"
      }
    },
    required: ["sql", "explanation"]
  },

  handler: async ({ sql, explanation }) => {
    // Security: Validate SQL is SELECT-only
    const normalizedSQL = sql.trim().toUpperCase();
    if (!normalizedSQL.startsWith("SELECT")) {
      throw new Error("Only SELECT queries are allowed");
    }

    // Execute via warehouse API
    const response = await fetch('/api/warehouse/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sql })
    });

    if (!response.ok) {
      throw new Error(`Query failed: ${response.statusText}`);
    }

    const { results, executionTime } = await response.json();

    return {
      sql,
      explanation,
      results,
      rowCount: results.length,
      executionTime,
      executedAt: new Date().toISOString()
    };
  },

  // Render component for the results
  render: ({ results, explanation, rowCount, executionTime }) => (
    <div className="query-result">
      <div className="query-header">
        <p className="explanation">{explanation}</p>
        <div className="meta">
          <span>{rowCount} rows</span>
          <span>{executionTime}ms</span>
        </div>
      </div>
      <ResultTable data={results} />
    </div>
  )
});

// Tool 2: Generate Report
const generateReportTool = tool({
  name: "generate_report",
  description: `Generate a business report with visualizations.
  Use when users want formatted summaries or periodic reports.`,

  parameters: {
    type: "object",
    properties: {
      reportType: {
        type: "string",
        enum: ["revenue", "patients", "appointments", "inventory", "custom"]
      },
      period: {
        type: "string",
        enum: ["daily", "weekly", "monthly", "quarterly", "yearly"]
      },
      format: {
        type: "string",
        enum: ["chart", "table", "pdf"],
        default: "chart"
      },
      title: {
        type: "string",
        description: "Optional custom title for the report"
      }
    },
    required: ["reportType", "period"]
  },

  handler: async ({ reportType, period, format, title }) => {
    const response = await fetch('/api/reports/generate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ reportType, period, format, title })
    });

    const report = await response.json();
    return report;
  },

  render: (report) => {
    switch (report.format) {
      case 'chart':
        return <ChartReport report={report} />;
      case 'pdf':
        return (
          <div className="pdf-download">
            <a href={report.downloadUrl} download>
              Download PDF Report
            </a>
          </div>
        );
      default:
        return <TableReport report={report} />;
    }
  }
});

// Tool 3: Sync to Tally (with HITL)
const tallySyncTool = tool({
  name: "sync_to_tally",
  description: `Sync approved journal entries to Tally ERP.
  IMPORTANT: This requires finance head approval before execution.
  Use when users want to push accounting entries to Tally.`,

  parameters: {
    type: "object",
    properties: {
      entryIds: {
        type: "array",
        items: { type: "string" },
        description: "List of journal entry IDs to sync"
      },
      company: {
        type: "string",
        description: "Tally company name"
      }
    },
    required: ["entryIds"]
  },

  handler: async ({ entryIds, company }, context) => {
    // 1. Fetch entries and check approval status
    const response = await fetch('/api/journal/entries', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ids: entryIds })
    });

    const { entries } = await response.json();

    // 2. Check for unapproved entries
    const unapproved = entries.filter(e => !e.isApproved);
    if (unapproved.length > 0) {
      return {
        success: false,
        requiresApproval: true,
        message: `${unapproved.length} entries require approval before syncing`,
        pendingEntries: unapproved.map(e => ({
          id: e.id,
          date: e.date,
          amount: e.amount,
          status: "pending_approval"
        }))
      };
    }

    // 3. Request human confirmation
    const confirmed = await context.requestConfirmation({
      title: "Confirm Tally Sync",
      message: `Sync ${entries.length} journal entries to Tally?`,
      details: entries.map(e => ({
        date: e.date,
        description: e.description,
        amount: `${e.currency} ${e.amount}`,
        ledger: e.ledgerName
      }))
    });

    if (!confirmed) {
      return {
        success: false,
        reason: "User cancelled the sync operation"
      };
    }

    // 4. Execute sync
    const syncResponse = await fetch('/api/tally/sync', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ entryIds, company })
    });

    const result = await syncResponse.json();

    return {
      success: true,
      syncedCount: result.synced.length,
      failedCount: result.failed.length,
      tallyVoucherNos: result.synced.map(s => s.voucherNo),
      failures: result.failed
    };
  },

  render: (result) => (
    <TallySyncResult
      success={result.success}
      syncedCount={result.syncedCount}
      failedCount={result.failedCount}
      requiresApproval={result.requiresApproval}
      failures={result.failures}
    />
  )
});

// Tool 4: OCR Document Extraction
const ocrExtractTool = tool({
  name: "extract_document_data",
  description: `Extract structured data from uploaded documents (invoices, receipts, bank statements).
  Use when users want to digitize documents for accounting.`,

  parameters: {
    type: "object",
    properties: {
      documentId: {
        type: "string",
        description: "ID of the uploaded document"
      },
      documentType: {
        type: "string",
        enum: ["invoice", "receipt", "bank_statement", "bill"],
        description: "Type of document for optimized extraction"
      }
    },
    required: ["documentId"]
  },

  handler: async ({ documentId, documentType }) => {
    const response = await fetch('/api/ocr/extract', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ documentId, documentType })
    });

    const extraction = await response.json();

    return {
      documentId,
      documentType: extraction.documentType,
      fields: extraction.fields,
      confidence: extraction.confidence,
      rawText: extraction.rawText,
      needsReview: extraction.confidence < 0.85
    };
  },

  render: (result) => (
    <DocumentExtractionResult
      fields={result.fields}
      confidence={result.confidence}
      needsReview={result.needsReview}
    />
  )
});

// ============================================
// COPILOTKIT PROVIDER SETUP
// ============================================

interface CopilotProviderProps {
  children: React.ReactNode;
}

export function MediSyncCopilotProvider({ children }: CopilotProviderProps) {
  return (
    <CopilotKit
      agent="medisync-agent"
      instructions={`You are MediSync AI, an intelligent assistant for healthcare business intelligence.

      Your capabilities:
      - Query business data using natural language (converted to SQL)
      - Generate formatted reports with visualizations
      - Extract data from documents using OCR
      - Sync accounting entries to Tally ERP (with user approval)

      Important rules:
      1. Always explain what you're doing before executing queries
      2. Show confidence scores for AI-generated suggestions
      3. Never push data to Tally without explicit user confirmation
      4. Support both English and Arabic - respond in the user's language
      5. For low-confidence results, suggest human review

      Current context:
      - User's company: {companyId}
      - User's role: {userRole}
      - Active period: {activePeriod}`}
      tools={[textToSQLTool, generateReportTool, tallySyncTool, ocrExtractTool]}
    >
      {children}
      <CopilotSidebar
        defaultOpen={false}
        labels={{
          title: "MediSync AI",
          initial: "How can I help you today?",
          placeholder: "Ask about revenue, patients, reports..."
        }}
      />
    </CopilotKit>
  );
}

// ============================================
// COMPONENTS WITH AGENT STATE
// ============================================

function Dashboard() {
  const { agent } = useAgent({ agentId: "medisync-agent" });

  // Access agent's current context
  const currentQuery = agent.state?.lastQuery;
  const filters = agent.state?.filters || {};

  // Update agent state when user changes filters
  const handleFilterChange = (newFilters) => {
    agent.setState({ filters: newFilters });
  };

  return (
    <div className="dashboard">
      <header>
        <h1>MediSync Dashboard</h1>
        <FilterBar filters={filters} onChange={handleFilterChange} />
      </header>

      <main>
        {currentQuery && (
          <QueryPreview
            sql={currentQuery.sql}
            results={currentQuery.results}
          />
        )}
        <DashboardGrid />
      </main>
    </div>
  );
}

// ============================================
// MAIN APP
// ============================================

export default function App() {
  return (
    <MediSyncCopilotProvider>
      <Dashboard />
    </MediSyncCopilotProvider>
  );
}

// ============================================
// HELPER COMPONENTS (simplified)
// ============================================

function ResultTable({ data }: { data: any[] }) {
  if (!data.length) return <p>No results</p>;

  const columns = Object.keys(data[0]);

  return (
    <table className="result-table">
      <thead>
        <tr>
          {columns.map(col => <th key={col}>{col}</th>)}
        </tr>
      </thead>
      <tbody>
        {data.map((row, i) => (
          <tr key={i}>
            {columns.map(col => <td key={col}>{String(row[col])}</td>)}
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function TallySyncResult({ success, syncedCount, failedCount, requiresApproval, failures }) {
  return (
    <div className={`tally-sync-result ${success ? 'success' : 'error'}`}>
      {requiresApproval ? (
        <p>Some entries require approval before syncing</p>
      ) : success ? (
        <p>Successfully synced {syncedCount} entries to Tally</p>
      ) : (
        <p>Sync failed: {failures?.map(f => f.error).join(', ')}</p>
      )}
    </div>
  );
}

function DocumentExtractionResult({ fields, confidence, needsReview }) {
  return (
    <div className="ocr-result">
      <div className="confidence-bar">
        <span>Confidence: {(confidence * 100).toFixed(0)}%</span>
        {needsReview && <span className="warning">Needs Review</span>}
      </div>
      <table>
        {Object.entries(fields).map(([key, value]) => (
          <tr key={key}>
            <th>{key}</th>
            <td>{String(value)}</td>
          </tr>
        ))}
      </table>
    </div>
  );
}

function ChartReport({ report }) {
  // Would use Apache ECharts in real implementation
  return <div className="chart-report">Chart: {report.title}</div>;
}

function TableReport({ report }) {
  return <ResultTable data={report.data} />;
}

function QueryPreview({ sql, results }) {
  return (
    <div className="query-preview">
      <pre>{sql}</pre>
      <ResultTable data={results} />
    </div>
  );
}

function FilterBar({ filters, onChange }) {
  return (
    <div className="filter-bar">
      {/* Filter controls */}
    </div>
  );
}

function DashboardGrid() {
  return (
    <div className="dashboard-grid">
      {/* Dashboard widgets */}
    </div>
  );
}
