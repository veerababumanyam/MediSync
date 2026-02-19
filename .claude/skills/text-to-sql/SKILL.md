---
name: text-to-sql
description: Convert natural language business questions into safe, read-only SQL queries against MediSync data warehouse (HIMS & Tally). Use for business intelligence, metric retrieval, and data exploration.
---

# Text-to-SQL Skill

Guidelines for generating precise, secure, and context-aware SQL from natural language intents, specifically for healthcare (HIMS) and accounting (Tally) domains.

## Query Construction Principles

### Domain-Specific Mapping
- **HIMS Data**: Map patient, appointment, and billing queries to the `hims` schema.
- **Tally Data**: Map ledger, voucher, and inventory queries to the `tally` schema.
- **Cross-Domain**: Join HIMS and Tally data using shared identifiers (e.g., `doctor_id`, `invoice_no`) when analyzing clinic performance vs. financial receipts.

### Safety & Security
- **Read-Only**: Always prepend/validate with `SELECT`. Block all DML/DDL (`INSERT`, `UPDATE`, etc.).
- **Row-Level Security**: Apply `tenant_id` and `user_role` filters in the `WHERE` clause automatically.
- **PII Masking**: Avoid selecting columns like `patient_phone`, `patient_address` unless explicitly required and role-authorized.

## SQL Implementation Patterns

### Metric Resolution (MetricFlow)
When a user asks for "Pharmacy Margin", resolve it to the underlying SQL fragment:
```sql
SELECT 
    (SUM(sales_amount) - SUM(cost_amount)) / NULLIF(SUM(sales_amount), 0) as margin
FROM tally.pharmacy_sales
WHERE {{ date_filter }}
```

### Complex Join (BI Scenario)
"Show revenue per doctor for the last month":
```sql
SELECT 
    d.doctor_name, 
    SUM(v.amount) as total_revenue
FROM hims.doctors d
JOIN tally.vouchers v ON d.billing_id = v.ledger_id
WHERE v.voucher_date >= CURRENT_DATE - INTERVAL '1 month'
GROUP BY d.doctor_name
ORDER BY total_revenue DESC;
```

## Tool Chain Documentation

### LangChain SQLDatabaseChain
- **Setup**: Use `SQLDatabase.from_uri()` with a read-only service account.
- **Prompt**: Inject `schema_context` and `semantic_context` into the system prompt.
- **Validation**: Use a custom `SQLValidator` tool to inspect the generated SQL string before execution.

### LlamaIndex for Schema Retrieval
- **Indexing**: Create a Vector Store index of schema metadata (table descriptions, column comments).
- **Retrieval**: Use semantic search to fetch relevant table schemas based on the user's natural language question instead of passing the entire DB schema.

## Accuracy & Quality Standards

- **Ambiguity Gate**: If the intent is unclear (e.g., "Show sales" - does it mean HIMS or Tally?), the agent must ask: "Would you like to see sales from the HIMS pharmacy or Tally accounting?"
- **Chart Hinting**: Suggest an optimal `chart_type` based on the results (e.g., 2 columns with numeric value → `bar`, date column → `line`).
- **Confidence Scoring**: Return a score (0-1). If `< 0.75`, add a "Verify Results" disclaimer.

## Accessibility Checklist
- [ ] Queries use meaningful column aliases (e.g., `AS "Monthly Revenue"`).
- [ ] Results include a "Data Source" attribution.
- [ ] Explanation field describes the query in plain English for non-technical users.
- [ ] Large result sets are paginated or summarized.
