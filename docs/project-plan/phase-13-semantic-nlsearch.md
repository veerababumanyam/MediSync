# Phase 13 — Semantic Layer & Natural Language Search

**Phase Duration:** Weeks 41–43 (3 weeks)  
**Module(s):** D — Advanced Search Analytics  
**Status:** Planning  
**Milestone:** M9 — Module D search core live  
**Depends On:** Phase 12 (Production v1 launched; stable platform)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md](../ARCHITECTURE.md)

---

## 1. Objectives

Launch Module D (Advanced Search Analytics) with its foundational layer: a fully governed semantic model, natural language full-text search across all data in the warehouse, a structured entity recognition pipeline, and a user-facing search interface with autocomplete and typed health/finance results. This phase establishes the data intelligence backbone that Phases 14–18 will build on.

---

## 2. Scope

### In Scope
- D-01: Natural Language Search Agent
- D-02: Entity Recognition Agent
- D-09: Semantic Layer Management Agent
- Full semantic model implementation (MetricFlow extensions)
- pgvector-powered NL search across all schemas
- Search UI: autocomplete, entity highlighting, search history, facets
- Query history and saved searches
- Semantic model versioning and governance
- Expanded metric catalogue (≥ 30 metrics)
- Integration with Module A chat (richer context from semantic model)

### Out of Scope
- D-03 to D-08, D-10 to D-14 (Phases 14–16)
- Autonomous agent actions (Phase 14)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | D-01: NL Search Agent production-ready | AI Engineer | Full-text + vector search; P95 < 2s; bilingual (EN/AR) |
| D-02 | D-02: Entity Recognition pipeline | AI Engineer | ≥ 90% entity precision on test set of 500 queries |
| D-03 | D-09: Semantic Layer Management Agent | AI Engineer + Data Eng | Manages MetricFlow YAML; versioning; governance tier; API accessible |
| D-04 | Semantic model ≥ 30 metrics | Data Engineer | Metrics validated, documented, and linked to source tables |
| D-05 | pgvector NL search index | Data Engineer | All fact tables embedded; HNSW index; search returns ranked results |
| D-06 | Search UI | Frontend Engineer | Autocomplete, entity highlighting, faceted results, query history; RTL |
| D-07 | MetricFlow extension | Data Engineer | Covers HIMS (clinical) + Tally (finance) business domains |
| D-08 | Chat integration (Module A enrichment) | AI + Backend | A-01 uses D-09 semantic model for improved metric resolution |

---

## 4. AI Agents Deployed

### D-01 — Natural Language Search Agent

**Purpose:** Translate free-text queries into structured multi-modal search across relational, vector, and JSONB data.

**Pipeline:**

```
User query (text)
      │ E-01 Language Detection
      ▼
E-02 Bilingual Normalisation (EN/AR stemming)
      │
      ▼
D-02 Entity Recognition
      │ Extracts: {entity_type, value, timeframe, filters}
      ▼
Search dispatcher:
  ├── Vector search (pgvector cosine similarity)   ← unstructured
  ├── Full-text search (PostgreSQL tsvector)        ← structured fields
  └── Semantic metric lookup (MetricFlow)           ← KPI queries
      │
      ▼
Result merger + ranking (BM25 + cosine blend)
      │
      ▼
Result presentation (entity cards, data tables, metric cards)
```

**Input:** `{ "query": string, "locale": "en"|"ar", "context_filters"?: object }`  
**Output:** `{ "results": [...], "entities_detected": [...], "query_interpreted": string, "search_type": "vector|fulltext|metric" }`  
**SLOs:** P95 latency < 2s; recall@10 ≥ 90% on test set

**Scoped search domains:**
| Domain | Tables searched | Result type |
|---|---|---|
| Transactions | `fact_vouchers`, `fact_transactions` | Transaction card |
| Vendors | `dim_vendors` | Entity card |
| Patients | `dim_patients` | Patient card (masked PII) |
| Metrics | MetricFlow metric catalogue | KPI card |
| Documents | `app.documents` (text chunks) | Document card |
| Reports | `app.report_deliveries` | Report card |

---

### D-02 — Entity Recognition Agent

**Purpose:** Extract typed entities from NL queries to enable precise structured queries.

**Entity types recognised:**
| Entity Type | Examples |
|---|---|
| `VENDOR` | "Al Noor Pharmacy", "Gulf Medical" |
| `PATIENT` | "patient ID 12345", "Ahmed Al Mansouri" |
| `LEDGER` | "sundry creditors", "salary expense" |
| `METRIC` | "revenue", "gross margin", "DSO" |
| `TIMEFRAME` | "last month", "Q3 2025", "يناير" (January in Arabic) |
| `COST_CENTRE` | "ICU", "pharmacy department" |
| `AMOUNT` | "over AED 5,000", "between 10k and 50k" |
| `DOCUMENT_TYPE` | "invoice", "receipt", "credit note" |

**Model:** Fine-tuned NER on domain corpus (healthcare + finance Arabic+English)  
**Output:** `{ "entities": [{ "type", "value", "span", "confidence" }] }`  
**Minimum confidence threshold:** 0.75 (lower: entity passed to user for confirmation)

**Arabic entity handling:**
- Normalise Arabic date terms: يناير→January, الربع الثالث→Q3
- Vendor names: transliterated name lookup in `dim_vendors` (both scripts)

---

### D-09 — Semantic Layer Management Agent

**Purpose:** Own and govern the MetricFlow YAML metric catalogue; enable versioned, discoverable metrics across all business domains.

**Capabilities:**
1. **Metric discovery:** Auto-scans source schemas and suggests new metrics
2. **YAML generation:** Natural language description → MetricFlow YAML definition
3. **Validation:** Validates metric YAML against source tables before commit
4. **Versioning:** Git-backed metric catalogue; change history; diff view
5. **Governance:** Requires approver sign-off for production metric changes
6. **Propagation:** Pushes approved changes to MetricFlow server; invalidates metric cache

**Metric YAML example (auto-generated):**

```yaml
metrics:
  - name: pharmacy_gross_margin
    description: "Gross margin percentage for pharmacy department"
    type: ratio
    label: "Pharmacy Gross Margin %"
    numerator:
      name: pharmacy_revenue_less_cogs
    denominator:
      name: total_pharmacy_revenue
    filter: |
      {{ Dimension('cost_centre__department') }} = 'pharmacy'
```

**API:**
- `GET /v1/semantic/metrics` — List all metrics with metadata
- `GET /v1/semantic/metrics/{name}` — Full definition + lineage
- `POST /v1/semantic/metrics` — Create (requires D-09 validation)
- `PATCH /v1/semantic/metrics/{name}` — Update (requires approver)
- `GET /v1/semantic/metrics/{name}/history` — Version history

---

## 5. Semantic Model: Metric Catalogue (v1 for Phase 13)

**Financial Metrics (Tally domain):**
| Metric Name | Description | Formula |
|---|---|---|
| `total_revenue` | Total revenue across all ledgers | sum(credit_vouchers) |
| `gross_profit` | Revenue minus direct costs | revenue - direct_cost |
| `operating_expenses` | all operating expense ledgers | sum(opex_vouchers) |
| `ebitda` | Earnings before interest, tax, depreciation, amortisation | operating_profit + D&A |
| `net_profit` | Net profit after all | gross_profit - opex - tax |
| `gross_margin_pct` | Gross profit / revenue × 100 | ratio |
| `accounts_receivable_balance` | Outstanding AR | sum(AR_ledger) |
| `accounts_payable_balance` | Outstanding AP | sum(AP_ledger) |
| `dso` | Days Sales Outstanding | AR / (revenue/days) |
| `dpo` | Days Payable Outstanding | AP / (cogs/days) |
| `cash_flow_from_operations` | Operating cash flow | derived from vouchers |
| `budget_variance_pct` | Actual vs budget | (actual-budget)/budget×100 |
| `tax_liability_gst` | GST liability | sum(GST_output - GST_input) |
| `tax_liability_vat` | VAT liability | sum(VAT_output - VAT_input) |

**Clinical / HIMS Metrics:**
| Metric Name | Description |
|---|---|
| `total_patient_visits` | Total OPD + IPD admissions |
| `avg_revenue_per_patient` | Revenue / unique patients |
| `pharmacy_revenue` | Revenue from pharmacy department |
| `laboratory_revenue` | Revenue from lab department |
| `bed_occupancy_rate` | IPD beds occupied / total beds |
| `average_length_of_stay` | IPD avg LOS in days |
| `new_vs_returning_patients` | Ratio new/returning |
| `top_procedures_by_revenue` | Top N procedures ranked |
| `outstanding_patient_invoices` | Unpaid patient bills |
| `collection_efficiency` | Collections / billings |

**Cross-domain:**
| Metric Name | Description |
|---|---|
| `pharmacy_gross_margin` | Pharmacy revenue - pharmacy COGS |
| `cost_per_patient_visit` | Total opex / total visits |
| `revenue_per_department` | Breakdown by cost-centre |

---

## 6. Database Schema Additions

```sql
-- Search index for full-text + vector search
CREATE TABLE vectors.search_documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_table TEXT NOT NULL,        -- 'fact_vouchers', 'app.documents', etc.
  source_id UUID NOT NULL,
  entity_type TEXT NOT NULL,         -- 'transaction', 'vendor', 'document', etc.
  content_text TEXT NOT NULL,        -- denormalised searchable text
  content_vector vector(1536),       -- OpenAI/local embedding
  locale TEXT DEFAULT 'en',
  metadata JSONB,
  indexed_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(source_table, source_id)
);

CREATE INDEX idx_search_docs_vector ON vectors.search_documents
  USING hnsw (content_vector vector_cosine_ops)
  WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_search_docs_fts ON vectors.search_documents
  USING gin (to_tsvector('english', content_text));

-- Semantic metric catalogue (metadata layer — YAML stored in Git, this is the index)
CREATE TABLE semantic.metrics_catalogue (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  metric_name TEXT UNIQUE NOT NULL,
  display_name_en TEXT NOT NULL,
  display_name_ar TEXT,
  description TEXT,
  domain TEXT,                       -- 'finance', 'clinical', 'cross'
  yaml_version INT DEFAULT 1,
  yaml_hash TEXT,                    -- SHA-256 of current YAML
  status TEXT DEFAULT 'active',      -- 'active', 'deprecated', 'draft'
  approved_by UUID,
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Search query history
CREATE TABLE app.search_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES app.users(id),
  query_text TEXT NOT NULL,
  locale TEXT,
  entities_detected JSONB,
  result_count INT,
  search_type TEXT,
  executed_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 7. Search UI Specification

**Search interface components:**

```
┌─────────────────────────────────────────────────────────────┐
│  🔍  [Search bar — autocomplete — placeholder: "Search      │
│       patients, invoices, metrics, reports..."]             │
│                                                             │
│  Filters:  [All ▼]  [Date range ▼]  [Cost centre ▼]       │
├─────────────────────────────────────────────────────────────┤
│  RECENT SEARCHES                                            │
│  • "Pharmacy revenue last quarter"                          │
│  • "Outstanding invoices Al Noor Pharmacy"                  │
│                                                             │
│  SAVED SEARCHES                                             │
│  • "Monthly ICU cost-centre report"                         │
├─────────────────────────────────────────────────────────────┤
│  RESULTS (124 found)                                        │
│                                                             │
│  📊 METRICS (3)                                             │
│  ┌────────────────────────────────┐                         │
│  │ Pharmacy Revenue  AED 2.4M ↑   │                         │
│  │ Last 30 days | +12% vs prev    │                         │
│  └────────────────────────────────┘                         │
│                                                             │
│  🧾 TRANSACTIONS (42)                                       │
│  Invoice #1234 | Al Noor | AED 50,000 | 15 Jan 2026 ✅     │
│  Invoice #1235 | Gulf Med | AED 12,500 | 16 Jan 2026 ⏳     │
│                                                             │
│  📄 DOCUMENTS (15)                                          │
│  Invoice_AlNoor_Jan2026.pdf | Classified: Invoice | Jan 15  │
└─────────────────────────────────────────────────────────────┘
```

**Autocomplete sources:** metric names, vendor names, patient IDs, ledger names, report names  
**RTL:** Full RTL layout when locale = Arabic; Arabic autocomplete terms from D-09 catalogue

---

## 8. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/search` | Execute NL search; returns ranked results |
| `GET` | `/v1/search/autocomplete?q=` | Autocomplete suggestions |
| `GET` | `/v1/search/history` | User's query history |
| `POST` | `/v1/search/save` | Save a search query |
| `GET` | `/v1/semantic/metrics` | List all metrics from catalogue |
| `GET` | `/v1/semantic/metrics/{name}` | Metric detail + YAML definition |
| `POST` | `/v1/semantic/metrics` | Create new metric (D-09 governed) |
| `GET` | `/v1/semantic/entities` | List recognised entity types + examples |

---

## 9. Testing Requirements

| Test | Target |
|---|---|
| D-01 NL search recall@10 | ≥ 90% on 500-query test set (EN) |
| D-01 Arabic search recall@10 | ≥ 85% on 200-query Arabic test set |
| D-02 entity precision | ≥ 90% precision on all 8 entity types |
| D-09 metric YAML validation | No invalid YAML ever committed to catalogue |
| Search P95 latency | < 2 seconds |
| pgvector index | HNSW index created; ANN recall ≥ 0.95 |
| Chat integration | A-01 references metric from D-09 catalogue correctly ≥ 95% |

---

## 10. Phase Exit Criteria

- [ ] D-01 NL Search Agent live: recall@10 ≥ 90% (EN), ≥ 85% (AR)
- [ ] D-02 Entity Recognition: precision ≥ 90%
- [ ] D-09 Semantic Layer: ≥ 30 metrics in catalogue; versioning and governance active
- [ ] pgvector HNSW index built on all source data
- [ ] Search UI live: autocomplete, facets, entity highlighting, history, RTL
- [ ] MetricFlow extended to cover all HIMS + Tally domains
- [ ] Module A chat enriched by D-09 semantic model
- [ ] All search API endpoints documented and passing integration tests
- [ ] Milestone M9 gate signed off

---

*Phase 13 | Version 1.0 | February 19, 2026*
