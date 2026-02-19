# MediSync Agent Development Guide

Comprehensive guide for developing AI agents in the MediSync platform using Genkit, Agent ADK, and A2A protocol.

## Genkit Flow Pattern

All MediSync agents follow the Genkit Flow pattern for orchestration.

### Basic Agent Structure

```go
package module_a

import (
    "context"
    "fmt"

    "github.com/firebase/genkit/go/genkit"
    "github.com/firebase/genkit/go/ai"
)

type TextToSQLRequest struct {
    Query      string `json:"query"`
    Locale     string `json:"locale"`
    UserID     string `json:"user_id"`
    CompanyID  string `json:"company_id"`
}

type TextToSQLResponse struct {
    SQL        string  `json:"sql"`
    Confidence float64 `json:"confidence"`
    Explanation string `json:"explanation"`
    ChartHint  string  `json:"chart_hint,omitempty"`
}

func (s *AgentService) TextToSQLFlow(ctx context.Context, req TextToSQLRequest) (*TextToSQLResponse, error) {
    // 1. Validate input
    if req.Query == "" {
        return nil, fmt.Errorf("query cannot be empty")
    }

    // 2. Detect language (E-01)
    locale := req.Locale
    if locale == "" {
        locale = s.detectLocale(ctx, req.Query)
    }

    // 3. Normalize domain terms (A-04)
    normalizedQuery, err := s.terminologyAgent.Normalize(ctx, req.Query, locale)
    if err != nil {
        return nil, fmt.Errorf("normalization failed: %w", err)
    }

    // 4. Generate SQL with LLM
    prompt := s.buildSQLPrompt(ctx, normalizedQuery, locale)
    model := s.geminiModel // or s.ollamaModel

    resp, err := genkit.Generate(ctx, model, ai.WithPrompt(prompt))
    if err != nil {
        return nil, fmt.Errorf("LLM generation failed: %w", err)
    }

    sql := extractSQL(resp.Text)

    // 5. Validate SQL is SELECT-only
    if !isSelectOnlyQuery(sql) {
        return nil, fmt.Errorf("generated SQL is not SELECT-only")
    }

    // 6. Execute via readonly role
    rows, err := s.warehouse.Query(ctx, sql, "medisync_readonly")
    if err != nil {
        return nil, fmt.Errorf("query execution failed: %w", err)
    }

    // 7. Score confidence (A-06)
    confidence := s.confidenceAgent.Score(ctx, sql, rows, resp.Text)

    // 8. Route to visualization (A-03)
    chartHint := s.vizAgent.SuggestChart(ctx, sql, rows)

    // 9. Format localized response (E-03)
    explanation := s.formatAgent.Explain(ctx, sql, rows, locale)

    return &TextToSQLResponse{
        SQL:        sql,
        Confidence: confidence,
        Explanation: explanation,
        ChartHint:  chartHint,
    }, nil
}

func isSelectOnlyQuery(sql string) bool {
    trimmed := strings.ToUpper(strings.TrimSpace(sql))
    return strings.HasPrefix(trimmed, "SELECT")
}
```

### Multi-Agent Orchestration with A2A

The Google A2A (Agent-to-Agent) protocol enables structured inter-agent communication.

```go
// A2A Message Structure
type A2AMessage struct {
    FromAgent    string                 `json:"from_agent"`
    ToAgent      string                 `json:"to_agent"`
    MessageID    string                 `json:"message_id"`
    ConversationID string               `json:"conversation_id"`
    Payload      map[string]interface{} `json:"payload"`
    Timestamp    time.Time              `json:"timestamp"`
}

// Agent Supervisor coordinates multi-agent workflows
func (s *Supervisor) ProcessQuery(ctx context.Context, query string) (*AgentResponse, error) {
    // Create conversation context
    convID := uuid.New().String()

    // Step 1: Language Detection (E-01)
    locale, err := s.agentE01.DetectLanguage(ctx, query)
    if err != nil {
        return nil, err
    }

    // Step 2: Domain Terminology (A-04)
    normalized, err := s.agentA04.Normalize(ctx, query, locale)
    if err != nil {
        return nil, err
    }

    // Step 3: Text to SQL (A-01)
    sqlResp, err := s.agentA01.GenerateSQL(ctx, normalized, locale)
    if err != nil {
        // Step 4: Self-Correction (A-02)
        sqlResp, err = s.agentA02.CorrectAndRetry(ctx, normalized, err, locale)
        if err != nil {
            return nil, err
        }
    }

    // Step 5: Execute Query
    data, err := s.warehouse.Query(ctx, sqlResp.SQL)
    if err != nil {
        return nil, err
    }

    // Step 6: Visualization Routing (A-03)
    chartType := s.agentA03.SelectChart(ctx, sqlResp.SQL, data)

    // Step 7: Confidence Scoring (A-06)
    confidence := s.agentA06.Score(ctx, query, sqlResp.SQL, data)

    // Step 8: Format Response (E-03)
    response := s.agentE03.Format(ctx, data, chartType, locale)

    return response, nil
}
```

## Complete Agent Catalog (58 Agents)

### Module A: Conversational BI (13 agents)

| ID | Name | Input | Output | Description |
|----|------|-------|--------|-------------|
| A-01 | Text-to-SQL | Natural language query | SQL SELECT statement | Converts user questions to safe SQL |
| A-02 | SQL Self-Correction | Failed SQL + error | Corrected SQL | Detects and fixes query errors |
| A-03 | Visualization Routing | SQL + result data | Chart type hint | Suggests optimal visualization |
| A-04 | Domain Terminology | Raw query | Normalized query | Maps healthcare/accounting terms |
| A-05 | Schema Explorer | User intent | Relevant tables | Identifies tables for query context |
| A-06 | Confidence Scorer | Query + SQL + result | 0-100% score | Evaluates answer reliability |
| A-07 | Drill-Down Advisor | Chart + user click | Next-level query | Suggests drill-down paths |
| A-08 | Query Simplifier | Complex SQL | Explanation | Explains SQL in natural language |
| A-09 | Join Optimizer | Multi-table query | Optimized SQL | Suggests efficient join strategies |
| A-10 | Time-Series Handler | Temporal query | Date-filtered SQL | Handles date range queries |
| A-11 | Aggregation Router | Question type | Aggregation type | Routes to SUM/COUNT/AVG/etc |
| A-12 | Fuzzy Matcher | Misspelled terms | Corrected terms | Handles typos in entity names |
| A-13 | Result Summarizer | Query result | Text summary | Generates insights from data |

### Module B: AI Accountant (16 agents)

| ID | Name | Input | Output | Description |
|----|------|-------|--------|-------------|
| B-01 | Document Classifier | Document image | Document type | Classifies invoice/bill/statement |
| B-02 | OCR Extraction | Document image | Extracted fields | Field-level text extraction |
| B-03 | Handwriting OCR | Handwritten text | Transcribed text | Specialized OCR for script |
| B-04 | Field Validator | Extracted fields | Validation result | Checks field completeness |
| B-05 | Ledger Mapping | Invoice line item | Suggested ledger | AI-suggests Tally GL account |
| B-06 | Tax Code Detector | Invoice data | Tax code | Identifies applicable tax codes |
| B-07 | Duplicate Finder | Invoice data | Duplicate status | Detects duplicate invoices |
| B-08 | Approval Workflow | Journal entry | Approval decision | Multi-level approval routing |
| B-09 | Tally Sync | Approved entry | Sync result | Pushes data to Tally ERP |
| B-10 | Bank Reconciliation | Bank statement | Reconciliation report | Matches bank to Tally entries |
| B-11 | Expense Categorizer | Expense description | Category | Classifies expenses by type |
| B-12 | Currency Converter | Multi-currency amount | Base currency amount | Handles currency conversion |
| B-13 | GST/VAT Calculator | Invoice amount | Tax breakdown | Calculates indirect taxes |
| B-14 | Payment Matcher | Payment vs invoice | Matching status | Matches payments to invoices |
| B-15 | Aging Report Generator | Outstanding invoices | Aging buckets | Generates aged receivables |
| B-16 | Cash Flow Predictor | Historical data | Forecast | Predicts cash flow trends |

### Module C: Easy Reports (8 agents)

| ID | Name | Input | Output | Description |
|----|------|-------|--------|-------------|
| C-01 | Template Resolver | Report name | SQL template | Maps reports to queries |
| C-02 | Parameter Validator | Report parameters | Validation result | Validates report inputs |
| C-03 | Multi-Company Consolidator | Company IDs | Combined data | Aggregates across companies |
| C-04 | Scheduled Report Runner | Schedule config | Execution result | Runs reports on schedule |
| C-05 | Export Formatter | Report data | PDF/Excel | Formats report exports |
| C-06 | Dashboard Builder | User config | Dashboard JSON | Zero-code dashboard creation |
| C-07 | Alert Monitor | Report output | Alert trigger | Sends threshold alerts |
| C-08 | Share Link Generator | Report instance | Shareable link | Creates secure report links |

### Module D: Search Analytics (14 agents)

| ID | Name | Input | Output | Description |
|----|------|-------|--------|-------------|
| D-01 | Query Intent Analyzer | User query | Intent classification | Identifies analysis type |
| D-02 | Data Source Router | Analysis intent | Source list | Identifies relevant data sources |
| D-03 | Hypothesis Generator | Research question | Testable hypotheses | Creates analysis hypotheses |
| D-04 | Autonomous AI Analyst | Research question | Full analysis | End-to-end analytical workflow |
| D-05 | Deep Research | Question + data | Research report | Pattern discovery |
| D-06 | Statistical Tester | Hypothesis + data | Test results | Runs statistical tests |
| D-07 | Anomaly Detector | Time series data | Anomaly list | Identifies outliers |
| D-08 | Prescriptive AI | Analysis results | Action recommendations | Quantified recommendations |
| D-09 | Trend Forecaster | Historical data | Forecast | Predicts future trends |
| D-10 | Correlation Finder | Variables | Correlation matrix | Finds relationships |
| D-11 | Code Generation | Analysis spec | Query code | Generates analysis code |
| D-12 | Visualization Planner | Analysis type | Viz recommendations | Suggests multi-chart layouts |
| D-13 | Insight Extractor | Analysis results | Key insights | Summarizes findings |
| D-14 | Report Writer | Analysis + insights | Narrative report | Generates narrative report |

### Module E: i18n (7 agents)

| ID | Name | Input | Output | Description |
|----|------|-------|--------|-------------|
| E-01 | Language Detection | Text | Language code | Detects en/ar with confidence |
| E-02 | Query Translation | Arabic query | English intent | Translates Arabic to English |
| E-03 | Localized Formatter | Data + locale | Formatted output | Formats numbers/dates/currency |
| E-04 | RTL Layout Handler | Component | RTL properties | Adjusts layouts for Arabic |
| E-05 | Terminology Translator | English term | Arabic equivalent | Domain-specific translations |
| E-06 | Number Formatter | Number + locale | Localized string | Handles Arabic numeral variants |
| E-07 | Date Formatter | Date + locale | Localized string | Hijri/Gregorian calendar support |

## Agent Development Patterns

### Error Handling Pattern

```go
type AgentError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    Confidence float64 `json:"confidence"`
    Suggestion string `json:"suggestion,omitempty"`
    Retryable  bool   `json:"retryable"`
}

func (e *AgentError) Error() string {
    return e.Message
}

// Usage in agent
func (s *MyAgent) Process(ctx context.Context, input string) (*Result, error) {
    if input == "" {
        return nil, &AgentError{
            Code:       "EMPTY_INPUT",
            Message:    "Input cannot be empty",
            Confidence: 1.0,
            Suggestion: "Please provide a valid query",
            Retryable:  false,
        }
    }

    // ... agent logic

    if err != nil {
        return nil, &AgentError{
            Code:       "LLM_ERROR",
            Message:    "Failed to generate response",
            Confidence: 0.0,
            Suggestion: "Please try rephrasing your query",
            Retryable:  true,
        }
    }

    return result, nil
}
```

### Confidence Scoring Pattern

```go
func (s *ConfidenceAgent) Score(ctx context.Context, query, sql string, result any) float64 {
    score := 1.0

    // Deduct for low result count
    if count := len(result); count == 0 {
        score -= 0.5
    } else if count < 3 {
        score -= 0.2
    }

    // Deduct for complex joins without explanation
    if strings.Count(strings.ToLower(sql), "join") > 2 {
        score -= 0.1
    }

    // Increase for exact matches
    if s.hasExactTermMatch(query, sql) {
        score += 0.1
    }

    // Deduct for SQL containing functions that may hallucinate
    if strings.Contains(strings.ToUpper(sql), "COALESCE") {
        score -= 0.15
    }

    return math.Max(0, math.Min(1, score))
}
```

### HITL Routing Pattern

```go
const ConfidenceThreshold = 0.75

func (s *AgentService) ProcessWithHITL(ctx context.Context, req Request) (*Response, error) {
    resp, err := s.flow(ctx, req)
    if err != nil {
        return nil, err
    }

    // Route to HITL if confidence is low
    if resp.Confidence < ConfidenceThreshold {
        return s.hitlGateway.RequestVerification(ctx, resp)
    }

    return resp, nil
}
```

### Logging Pattern

```go
import "log/slog"

func (s *AgentService) MyFlow(ctx context.Context, req Request) (*Response, error) {
    logger := slog.With(
        "agent", "A-01",
        "user_id", req.UserID,
        "query", req.Query,
    )

    logger.InfoContext(ctx, "Starting agent flow")

    // ... processing

    logger.InfoContext(ctx, "Agent flow completed",
        "confidence", resp.Confidence,
        "sql", resp.SQL,
    )

    return resp, nil
}
```

## Testing Agents

### Unit Testing with Mock LLM

```go
func TestTextToSQLFlow(t *testing.T) {
    // Mock the LLM
    mockLLM := &MockModel{
        Response: "SELECT * FROM patients WHERE city = 'Dubai'",
    }

    service := &AgentService{
        geminiModel: mockLLM,
        warehouse:   &MockWarehouse{},
    }

    req := TextToSQLRequest{
        Query: "Show me patients in Dubai",
        Locale: "en",
    }

    resp, err := service.TextToSQLFlow(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, "SELECT * FROM patients WHERE city = 'Dubai'", resp.SQL)
    assert.Greater(t, resp.Confidence, 0.5)
}
```

### Integration Testing with Deterministic LLM

```go
// Use temperature=0 for deterministic responses
func TestAgentIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Use real LLM with temperature=0
    model := genkit.DefineModel("gemini-1.5-flash", ai.WithTemperature(0))

    service := NewAgentService(model, realWarehouse)

    // Test with known inputs
    testCases := []struct {
        query     string
        expectedSQL string
    }{
        {"count patients", "SELECT COUNT(*) FROM patients"},
        {"list doctors", "SELECT * FROM doctors"},
    }

    for _, tc := range testCases {
        resp, err := service.TextToSQLFlow(ctx, TextToSQLRequest{Query: tc.query})
        assert.NoError(t, err)
        assert.Contains(t, resp.SQL, "SELECT")
    }
}
```

## Agent Registration

```go
// Register all agents in the supervisor
func NewSupervisor(config Config) *Supervisor {
    return &Supervisor{
        agentA01: NewTextToSQLAgent(config),
        agentA02: NewSQLSelfCorrectionAgent(config),
        agentA03: NewVisualizationRoutingAgent(config),
        // ... all 58 agents
    }
}
```

## OPA Policy Integration

Every agent that accesses data must check OPA policies:

```go
func (s *AgentService) Query(ctx context.Context, sql string) (*Result, error) {
    // Check OPA policy before executing
    allowed, err := s.opaClient.Allow(ctx, "warehouse_query", map[string]interface{}{
        "user_id": s.userID,
        "sql":     sql,
    })
    if err != nil {
        return nil, fmt.Errorf("OPA check failed: %w", err)
    }
    if !allowed {
        return nil, fmt.Errorf("not authorized for this query")
    }

    return s.warehouse.Query(ctx, sql)
}
```
