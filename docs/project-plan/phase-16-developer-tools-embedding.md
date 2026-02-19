# Phase 16 — Developer Tools & Embedded Analytics

**Phase Duration:** Weeks 49–50 (2 weeks)  
**Module(s):** D — Advanced Search Analytics  
**Status:** Planning  
**Depends On:** Phase 15 (Full Module D analytics intelligence live)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Unlock MediSync as a developer platform: allow technical users (data engineers, analysts, power users) to write and execute Python notebooks, generate code directly from natural language, federate queries across multiple data nodes with in-memory optimisation, and expose an embedding API surface for external portals to consume MediSync analytics inside their own UIs.

---

## 2. Scope

### In Scope
- D-11: Code Generation Agent (SpotterCode)
- D-12: Federated Query Optimisation Agent
- Analyst Studio: browser-based Python notebook environment
- REST + GraphQL embedding APIs
- SDK (JavaScript/TypeScript) for embedding charts
- Federated multi-node query (cross-entity data join without data copy)
- Zero-copy in-memory query acceleration (DuckDB)

### Out of Scope
- D-13, D-14 (already deployed in Phase 14/15)
- Governance layer (Phase 17)
- Production v2 launch (Phase 18)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | D-11: SpotterCode agent | AI Eng | NL → Python/SQL/React code; reviewed before execution |
| D-02 | D-12: Federated Query Optimisation | AI Eng + Data Eng | Federated query across 2 nodes with ≥ 40% latency reduction vs naive join |
| D-03 | Analyst Studio (Python notebooks) | Frontend + Backend | Jupyter-compatible browser IDE; execute Python cells; access MediSync datasources |
| D-04 | Embedding REST API | Backend Eng | Authenticated embed tokens; chart data served via REST |
| D-05 | Embedding GraphQL API | Backend Eng | GraphQL schema for flexible metric + dimension queries |
| D-06 | JavaScript Embed SDK | Frontend Eng | `npm install @medisync/embed`; `<MediSyncChart>` React component |
| D-07 | DuckDB in-memory layer | Data Eng | Hot-path queries served from DuckDB; measurable latency improvement |
| D-08 | Developer documentation | All Eng | SDK docs, embed API reference, Python SDK reference, notebook examples |

---

## 4. AI Agents Deployed

### D-11 — SpotterCode (Code Generation Agent)

**Purpose:** Translate natural language requests into executable code — Python data analysis scripts, SQL queries, or React dashboard components — and present them for human review before execution.

**Supported output types:**
| Output type | Example request | Output |
|---|---|---|
| Python data script | "Write a Python script to export all vendor AP balances to CSV" | Pandas DataFrame + export code |
| SQL query | "Write SQL to calculate 12-month rolling gross margin" | Valid PostgreSQL CTE |
| MetricFlow YAML | "Create a metric for net collection rate" | YAML for D-09 catalogue |
| React component | "Build a chart component showing OPD visitor trend" | JSX using ECharts |
| Meltano pipeline | "Write a Meltano config to sync new HIMS procedures table" | `meltano.yml` snippet |

**Execution model:**
```
User NL request
      │
      ▼
D-11 generates code + inline explanation
      │
      ▼
Code rendered in code viewer (syntax highlighted)
      │
User reviews: [COPY]  [RUN IN ANALYST STUDIO]  [SEND TO PIPELINE]
      │
      ├─ Run in Studio → Analyst Studio executes in sandbox
      └─ Send to Pipeline → Meltano/scheduler queue
```

**Safety controls:**
- Generated SQL: read-only. Never generates `UPDATE`, `DELETE`, `DROP`.
- Generated Python: executes in sandboxed kernel (no filesystem write outside `/tmp`)
- React code: requires developer role to deploy to production; preview-only in sandbox

**LLM context:** D-11 is given the full D-09 semantic model + current schema as context; this ensures generated SQL references real table/column names.

---

### D-12 — Federated Query Optimisation Agent

**Purpose:** Enable queries that span multiple logical data nodes (e.g., two separate clinic HIMS databases, or HIMS + external pharmacy system) without physically copying data between them.

**Problem addressed:** Multi-entity healthcare organisations often have separate instances (per hospital / branch). Cross-entity analytics (e.g., "compare revenue across all branches") would normally require a full ETL merge. D-12 enables query-time federation.

**Architecture:**

```
Client query: "Total revenue by branch — all 3 entities"
      │
      ▼
D-12 Query Analyser: decomposes query → per-node sub-queries
      │
      ├─ Node 1 (Entity A): SQL → partial result
      ├─ Node 2 (Entity B): SQL → partial result
      └─ Node 3 (Entity C): SQL → partial result
      │
      ▼
D-12 Result Merger: JOIN + UNION in DuckDB in-memory
      │
      ▼
Unified result returned to user
```

**Optimisation strategies:**
| Strategy | Description |
|---|---|
| Predicate pushdown | Push WHERE clauses to individual nodes before merging |
| Result caching | Cache per-node partial results in Redis (TTL 5 min) |
| Parallel execution | All node sub-queries launch in parallel (goroutines) |
| DuckDB in-memory join | Merge partial results in DuckDB for zero-copy performance |

**Supported node types:** PostgreSQL (primary), PostgreSQL read replica, REST API (with schema descriptor)

**Security:** Each node query runs under the requesting user's Keycloak roles; no cross-tenant data leakage possible (OPA `d12.cross_entity_gate` policy).

---

## 5. Analyst Studio

**Purpose:** Browser-based Python notebook environment for power users / data engineers to explore data, build analyses, and prototype reports.

**Technology:** JupyterLite (WASM-based) or server-side JupyterHub (configurable at deployment)  
**Python SDK:** Pre-installed `medisync` Python package with shortcuts to data access

```python
# Example: MediSync Python SDK in Analyst Studio
import medisync as ms

# Connect (uses user's session token automatically)
client = ms.Client()

# Metric query via semantic layer
df = client.metric("pharmacy_revenue", period="last_quarter", granularity="month")
print(df.head())

# Raw SQL (read-only)
df2 = client.query("SELECT * FROM tally_analytics.fact_vouchers WHERE voucher_date > '2026-01-01' LIMIT 100")

# Search
results = client.search("invoices from Al Noor over AED 50K")

# Export to charts (ECharts spec)
chart = ms.bar_chart(df, x="month", y="revenue", title="Pharmacy Revenue Q1")
chart.show()
```

**Notebook features:**
- Code cells (Python 3.11)
- Markdown cells
- Output: tables, charts, JSON
- Auto-complete (D-11 inline suggestions)
- Save/load notebooks from user library
- Export notebook as PDF report (WeasyPrint)

**Access control:** Analyst Studio accessible to: `analyst`, `developer`, `admin` roles only; `accountant` and `viewer` roles cannot access.

---

## 6. Embedding APIs

### REST Embedding API

**Embed token creation:**
```http
POST /v1/embed/token
Authorization: Bearer {admin_token}
{
  "user_id": "external_user_123",
  "allowed_metrics": ["pharmacy_revenue", "total_revenue"],
  "allowed_dashboards": ["dashboard-uuid-1"],
  "expires_at": "2026-03-01T00:00:00Z"
}
```
**Returns:** `{ "embed_token": "emb_xxx..." }` (JWT, short-lived, scoped)

**Chart data endpoint:**
```http
GET /v1/embed/metric/pharmacy_revenue?period=last_month
Authorization: Embed emb_xxx...
```
**Returns:** `{ "value": 2400000, "chart_data": [...], "locale": "en" }`

---

### GraphQL Embedding API

```graphql
query PharmacyMetrics($period: Period!) {
  metrics(names: ["pharmacy_revenue", "gross_margin_pct"], period: $period) {
    name
    value
    trend {
      direction
      pct_change
    }
    chartSeries {
      label
      data
    }
  }
  dimensions(name: "cost_centre") {
    values
  }
}
```

**Endpoint:** `POST /v1/graphql` (accepts embed token or full session token)

---

### JavaScript Embed SDK

```javascript
// npm install @medisync/embed

import { MediSyncEmbed } from '@medisync/embed';

const embed = new MediSyncEmbed({
  token: 'emb_xxx...',
  baseUrl: 'https://medisync.example.com'
});

// Render a metric card
embed.renderMetric('pharmacy_revenue', {
  container: document.getElementById('pharmacy-card'),
  period: 'last_month',
  locale: 'en'
});

// Render a full dashboard
embed.renderDashboard('dashboard-uuid-1', {
  container: document.getElementById('dashboard'),
  theme: 'light', // or 'dark'
  locale: 'ar',   // RTL automatic
});
```

---

## 7. DuckDB In-Memory Optimisation

**Hot-path queries served from DuckDB:**
```
Regular PostgreSQL path:
  Client → Go API → PostgreSQL → result           (avg 800ms)

DuckDB hot path:
  Client → Go API → DuckDB in-memory parquet → result  (avg 120ms)
```

**Population strategy:**
- DuckDB instance hydrated on startup from PostgreSQL via Parquet export
- Refresh: `fact_vouchers`, `fact_transactions` refreshed every 5 minutes
- Cache invalidation: triggered by Tally sync completion event (NATS)

**Scope:** DuckDB serves: aggregate metric queries, time-series chart data, dimension lookups  
Not used for: write operations, document queries, live transactional data

---

## 8. Database Schema Additions

```sql
-- Analyst Studio notebooks
CREATE TABLE app.notebooks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES app.users(id),
  title TEXT NOT NULL,
  content JSONB NOT NULL,       -- JupyterLab notebook JSON
  kernel_state JSONB,
  last_executed TIMESTAMPTZ,
  status TEXT DEFAULT 'draft',  -- 'draft', 'published', 'archived'
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Embed tokens
CREATE TABLE app.embed_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  token_hash TEXT UNIQUE NOT NULL,    -- SHA-256 of token (never store raw)
  created_by UUID REFERENCES app.users(id),
  allowed_metrics TEXT[],
  allowed_dashboards UUID[],
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  last_used_at TIMESTAMPTZ,
  use_count INT DEFAULT 0
);

-- Federated query nodes
CREATE TABLE app.federated_nodes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  node_name TEXT UNIQUE NOT NULL,
  entity_id UUID REFERENCES app.tally_entities(id),
  connection_type TEXT NOT NULL,    -- 'postgresql', 'rest'
  connection_config JSONB,          -- encrypted connection string
  status TEXT DEFAULT 'active',
  last_health_check TIMESTAMPTZ
);
```

---

## 9. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/code/generate` | D-11 generate code from NL |
| `POST` | `/v1/code/execute` | Execute in Analyst Studio sandbox |
| `GET` | `/v1/notebooks` | List user notebooks |
| `POST` | `/v1/notebooks` | Create/save notebook |
| `POST` | `/v1/embed/token` | Create embed token |
| `GET` | `/v1/embed/metric/{name}` | Embed-scoped metric data |
| `POST` | `/v1/graphql` | GraphQL analytics query |
| `GET` | `/v1/federated/nodes` | List federated nodes |
| `POST` | `/v1/federated/query` | Federated cross-node query |

---

## 10. Testing Requirements

| Test | Target |
|---|---|
| D-11 SQL generation | Generated SQL parseable; valid against schema; 0 mutation statements |
| D-11 Python generation | Generated code runs in sandbox without error on 15 test prompts |
| D-12 federated query | 2-node federated query returns correct merged result; ≥ 40% latency improvement |
| Analyst Studio | Python cell executes; `medisync` SDK connects; table + chart output rendered |
| Embed token scope | Embed token cannot access metrics outside `allowed_metrics` (OPA enforced) |
| GraphQL | 10 test queries return correct data; schema introspection available |
| JS SDK | `renderMetric` and `renderDashboard` render correctly in test app |
| DuckDB hot path | P95 aggregate metric < 200ms (vs 800ms cold PostgreSQL) |

---

## 11. Phase Exit Criteria

- [ ] D-11 SpotterCode: generates valid Python, SQL, MetricFlow YAML from NL
- [ ] D-12 Federated Query: 2-node federation working; DuckDB merge in-memory
- [ ] Analyst Studio: accessible to analyst/developer roles; Python SDK connected
- [ ] REST embedding API: embed tokens scoped + working; chart data endpoint
- [ ] GraphQL API: schema live; metric + dimension queries working
- [ ] JavaScript Embed SDK: `renderMetric` + `renderDashboard` functional
- [ ] DuckDB hot path: measurable latency improvement documented
- [ ] Developer documentation complete (SDK + embed API + notebook examples)
- [ ] Phase gate reviewed and signed off

---

*Phase 16 | Version 1.0 | February 19, 2026*
