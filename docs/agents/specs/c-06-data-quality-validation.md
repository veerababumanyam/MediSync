# Agent Specification — C-06: Data Quality Validation Agent

**Agent ID:** `C-06`  
**Agent Name:** Data Quality Validation Agent  
**Module:** C — Easy Reports  
**Phase:** 1  
**Priority:** P0 Critical  
**HITL Required:** No (blocks pipeline; alerts humans)  
**Status:** Draft

---

## 1. Purpose

Runs automated data quality checks (completeness, uniqueness, referential integrity, range, temporal consistency, cross-source reconciliation) on every ETL batch before data enters the warehouse. Blocks bad data from loading; alerts on anomalies.

> **Addresses:** PRD §8, §6.8.6 — Data quality as a core platform capability starting Phase 1.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Event trigger** | ETL pipeline DAG task completion (runs as a gate step) |
| **Scheduled trigger** | Every ETL run (continuous) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `batch_id` | `UUID` | ETL pipeline | ✅ |
| `source` | `enum` | `tally / hims / bank` | ✅ |
| `staging_table` | `string` | ETL staging schema | ✅ |
| `expectation_suite` | `string` | great_expectations suite name | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `validation_passed` | `bool` | Gate decision: load or block |
| `results` | `[]ExpectationResult` | Per-check pass/fail |
| `failure_count` | `int` | Critical failures |
| `warning_count` | `int` | Non-blocking warnings |
| `anomalies` | `[]Anomaly` | Statistical outliers detected |
| `data_quality_report` | `DataQualityReport` | Summary for DQ dashboard |

---

## 5. Check Categories

| Category | Examples | On Failure |
|----------|---------|-----------|
| Completeness | `patient_id NOT NULL`, `amount NOT NULL` | Block |
| Uniqueness | No duplicate `voucher_id`, `invoice_no` | Block |
| Referential Integrity | All ledger codes exist in COA | Block |
| Range Validation | `amount > 0 AND amount < 10_000_000` | Block |
| Temporal | `transaction_date` within valid range | Block |
| Cross-Source | HIMS billing ≈ Tally receipts for same period | Warn |

---

## 6. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | great_expectations | Apache-2.0 | Declarative expectation suites |
| 2 | Validation runner (Python or Go) | Internal | Execute expectations on staging table |
| 3 | Z-score anomaly detector (Go) | Internal | Statistical outlier detection |
| 4 | Apprise | MIT | Alert on critical failures |
| 5 | B-14 Audit Log Writer | Internal | Log all validation results |

```
ETL Staging Table
  → great_expectations validation run
  → Anomaly detection (Z-score on daily totals)
  → Decision: PASS → promote to warehouse | FAIL → block + alert
  → Write validation result to data_quality_log table
  → Apprise alert (on FAIL)
```

---

## 7. Guardrails

- **Hard block** on critical failures (data never promoted to warehouse on failure).
- Warnings: data loaded with `data_quality_flag=warning` tag.
- All validation results retained in `data_quality_log` for trend monitoring.

---

## 8. Evaluation Criteria

| Metric | Target |
|--------|--------|
| False positive rate (valid data blocked) | < 0.1% |
| Critical defect detection rate | ≥ 99.9% |
| Validation run latency (per batch) | < 60s |

---

## 9. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Python sidecar (great_expectations) + Go orchestrator |
| **Integration** | Runs as a mandatory gate step in ETL pipeline |
| **Depends on** | great_expectations, Apprise |
| **Consumed by** | ETL pipeline, DQ monitoring dashboard |
