# Feature Specification: AI Agent & Core Analytics

**Feature Branch**: `001-ai-agent-core`
**Created**: 2026-02-19
**Status**: Draft
**Input**: Phase 02 plan - AI Agent & Core Analytics layer for conversational BI

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Natural Language Query to Data (Priority: P1)

A clinic manager types "Show me total revenue for January 2026" in English or Arabic and receives an accurate, formatted response with a chart within 5 seconds—without knowing SQL or database structure.

**Why this priority**: This is the foundational capability that defines the product value. All other features depend on users being able to ask questions and get answers.

**Independent Test**: Can be fully tested by submitting natural language queries through the `/v1/chat` API and verifying SQL accuracy, response formatting, and latency. Delivers immediate BI value.

**Acceptance Scenarios**:

1. **Given** an authenticated user submits "Show clinic revenue last month", **When** the query is processed, **Then** a KPI card displays total revenue with correct SQL generated and 95%+ confidence score
2. **Given** an authenticated user submits an Arabic query "أظهر إيرادات العيادة في يناير", **When** the query is processed, **Then** the response is in Arabic with correctly formatted numbers (Eastern Arabic numerals)
3. **Given** a user submits an ambiguous query like "Show revenue", **When** the system detects ambiguity, **Then** a clarifying question is returned asking for time period and revenue type

---

### User Story 2 - Multi-Chart Visualization (Priority: P2)

A pharmacy owner asks "Compare drug sales by category this quarter" and receives a bar chart; asks "What's the trend of patient visits?" and receives a line chart—automatically routed to the optimal visualization.

**Why this priority**: Visualization significantly enhances comprehension. After users get answers (P1), they need to see patterns visually.

**Independent Test**: Can be tested by submitting queries with different intent patterns (trend, comparison, breakdown) and verifying correct chart type assignment.

**Acceptance Scenarios**:

1. **Given** a user asks "What's the trend of patient visits over the last 6 months?", **When** the query is processed, **Then** a line chart is returned with time on X-axis
2. **Given** a user asks "Compare revenue across departments", **When** the query is processed, **Then** a bar chart is returned with departments as categories
3. **Given** a user asks "What percentage of revenue comes from pharmacy vs clinic?", **When** the query is processed, **Then** a pie chart is returned showing proportions

---

### User Story 3 - Error Self-Correction (Priority: P2)

A user asks a complex question that generates SQL with an error. The system automatically detects the error, corrects the query, and retries—user sees the result without knowing an error occurred.

**Why this priority**: Improves reliability and user trust. Prevents frustration from failed queries.

**Independent Test**: Can be tested by submitting queries that deliberately trigger SQL errors and verifying the self-correction loop works.

**Acceptance Scenarios**:

1. **Given** a generated SQL query fails with "column does not exist", **When** A-02 agent corrects it, **Then** the corrected query executes successfully within 3 retry attempts
2. **Given** a generated SQL query has a syntax error, **When** A-02 agent analyzes the error, **Then** a valid corrected SQL is generated and executed
3. **Given** all 3 retry attempts fail, **When** the system exhausts self-correction, **Then** a graceful error message is returned with partial results if available

---

### User Story 4 - Off-Topic Query Rejection (Priority: P3)

A user asks "What's the weather today?" or "Write me a poem"—the system politely declines and redirects to its purpose as a healthcare/finance data analyst.

**Why this priority**: Prevents misuse and maintains focus. Important for professional credibility but not core functionality.

**Independent Test**: Can be tested by submitting non-business queries and verifying rejection responses.

**Acceptance Scenarios**:

1. **Given** a user asks a general knowledge question unrelated to healthcare/finance, **When** the hallucination guard processes it, **Then** a polite redirect message is returned
2. **Given** a user asks a valid business question, **When** the hallucination guard processes it, **Then** it passes through without false positive rejection

---

### User Story 5 - Confidence-Based Routing (Priority: P3)

A user asks a question with multiple possible interpretations. The system detects low confidence and either asks for clarification or flags the response for human review.

**Why this priority**: Enhances trust and accuracy but is a refinement over the core query-response flow.

**Independent Test**: Can be tested by submitting ambiguous queries and verifying confidence scores and routing behavior.

**Acceptance Scenarios**:

1. **Given** a query has confidence score below 50%, **When** A-06 evaluates it, **Then** a clarifying question is returned instead of executing SQL
2. **Given** a query has confidence score between 50-69%, **When** A-06 evaluates it, **Then** the result is shown with a "Low confidence" warning and added to review queue
3. **Given** a query has confidence score 70% or above, **When** A-06 evaluates it, **Then** the result is shown normally with confidence badge visible

---

### Edge Cases

- **Empty query**: What happens when user submits an empty or whitespace-only query?
  - System returns a helpful prompt asking for a question
- **Extremely long query**: What happens when query exceeds reasonable length?
  - System truncates with a warning or asks user to simplify
- **Mixed language query**: What happens when user mixes English and Arabic in same query?
  - Language detection picks the dominant language; intent extraction handles mixed terms
- **Query with no matching data**: What happens when SQL executes but returns zero rows?
  - System returns "No data found" message with suggestions to adjust filters
- **Query timeout**: What happens when SQL execution takes too long?
  - System cancels after configurable timeout (default 30s) and returns timeout message
- **User lacks permission for requested data**: What happens when OPA blocks column access?
  - System returns data with restricted columns masked or omitted, with explanation

---

## Requirements *(mandatory)*

### Functional Requirements

#### Core Query Processing

- **FR-001**: System MUST accept natural language queries in English and Arabic through `/v1/chat` API endpoint
- **FR-002**: System MUST convert natural language queries to valid, read-only SQL statements
- **FR-003**: System MUST enforce SELECT-only queries—all DML operations (INSERT, UPDATE, DELETE, DROP) must be blocked
- **FR-004**: System MUST execute SQL queries using a read-only database role with restricted permissions
- **FR-005**: System MUST return query results with confidence scores (0-100%)

#### Language Support

- **FR-006**: System MUST detect query language (English or Arabic) with ≥99% accuracy
- **FR-007**: System MUST normalize Arabic queries to English intent before SQL generation
- **FR-008**: System MUST format all responses in the user's detected or preferred locale
- **FR-009**: System MUST format numbers according to locale (English: 1,234.56 | Arabic: ١٬٢٣٤٫٥٦)

#### Domain Intelligence

- **FR-010**: System MUST normalize healthcare and accounting domain terminology to canonical terms (e.g., "footfall" → patient_visits)
- **FR-011**: System MUST maintain a bilingual glossary of domain terms for synonym mapping
- **FR-012**: System MUST use schema context from vector embeddings to inform SQL generation

#### Visualization

- **FR-013**: System MUST automatically classify query results and route to appropriate visualization type
- **FR-014**: System MUST support the following visualization types: line chart, bar chart, pie chart, KPI card, data table

#### Error Handling & Recovery

- **FR-015**: System MUST detect SQL execution errors and attempt self-correction (up to 3 retries)
- **FR-016**: System MUST reject off-topic queries with a helpful redirect message
- **FR-017**: System MUST return graceful error messages when queries cannot be answered

#### Security & Access Control

- **FR-018**: System MUST validate all user authentication via JWT before processing queries
- **FR-019**: System MUST apply column-level masking for PII fields based on user role
- **FR-020**: System MUST block access to cost/margin columns for non-manager roles
- **FR-021**: System MUST use parameterized queries to prevent SQL injection

#### Streaming & Performance

- **FR-022**: System MUST stream responses via Server-Sent Events (SSE) with progress updates
- **FR-023**: System MUST complete standard queries within 5 seconds (P95 latency)

### Key Entities

- **Query Session**: Represents a user's chat session containing query history, locale preference, and context
- **Natural Language Query**: The user's raw question in English or Arabic
- **SQL Statement**: The generated read-only SQL query with parameterized bindings
- **Query Result**: The data returned from query execution with metadata (row count, columns)
- **Visualization Specification**: The chart type configuration and data mapping for rendering
- **Confidence Score**: A numerical assessment (0-100) of result accuracy
- **Schema Embedding**: Vector representations of table/column descriptions for semantic search
- **Domain Term**: A mapping between user vocabulary and canonical database terminology
- **Audit Log Entry**: Record of query submission, execution, and response for compliance

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users receive accurate answers to 95% of business queries on the 50-query test set
- **SC-002**: System responds to standard queries in under 5 seconds (P95 latency)
- **SC-003**: SQL self-correction successfully recovers from 80% of SQL errors
- **SC-004**: Visualization routing achieves 98% accuracy on the 40-sample test set
- **SC-005**: Language detection achieves 99% accuracy on 100 mixed-language queries
- **SC-006**: Arabic response formatting is 100% correct for numbers, dates, and currency
- **SC-007**: Zero false positives in off-topic query detection (valid queries never blocked)
- **SC-008**: Zero false negatives in off-topic query detection (all off-topic queries rejected)
- **SC-009**: Domain terminology mapping covers 100% of the 30-term glossary
- **SC-010**: Confidence scores correlate positively with actual query accuracy
- **SC-011**: All generated SQL passes read-only policy validation (no DML)
- **SC-012**: PII column masking works correctly for viewer role in 100% of test cases

---

## Assumptions

- Users have valid Keycloak JWT tokens with appropriate role claims
- PostgreSQL data warehouse is populated with seed data from Phase 01
- LLM API keys (GPT-5.2 or local Ollama) are configured and available
- Bilingual reviewer is available to validate Arabic translation quality
- Standard healthcare and accounting domain knowledge is applicable
- Users have basic familiarity with BI concepts but not SQL

---

## Dependencies

- Phase 01 must be complete (data warehouse seeded with tables and sample data)
- Keycloak identity provider running and configured
- PostgreSQL with pgvector extension installed
- LLM API access (cloud or local) configured
- OPA policy engine running with `bi.read_only` policy loaded

---

## Out of Scope

- Frontend UI components (Phase 03)
- Dashboard creation and pinning (Phase 03)
- Scheduled reports (Phase 03)
- Document processing and OCR (Phase 04+)
- Write-back operations to Tally or HIMS (Phase 04+)
- Additional languages beyond English and Arabic

---

*Specification Version 1.0 | February 19, 2026*
