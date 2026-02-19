# Agent Specs — Index

Complete per-agent specification files for all 58 MediSync AI agents across 5 modules.

All specs use the Go + Genkit stack (see [01-oss-toolchain.md](../01-oss-toolchain.md)). Python sidecars are used only for ML workloads (OCR, forecasting, NER, voice).

---

## Module A — Conversational BI (13 agents)

| ID | Agent | Phase | Priority | HITL | Spec |
|----|-------|-------|----------|------|------|
| A-01 | Text-to-SQL Agent | 5 | P0 | No | [a-01-text-to-sql.md](a-01-text-to-sql.md) |
| A-02 | SQL Self-Correction Agent | 5 | P0 | No | [a-02-sql-self-correction.md](a-02-sql-self-correction.md) |
| A-03 | Visualization Routing Agent | 5 | P1 | No | [a-03-visualization-routing.md](a-03-visualization-routing.md) |
| A-04 | Domain Terminology Agent | 5 | P1 | No | [a-04-domain-terminology.md](a-04-domain-terminology.md) |
| A-05 | Hallucination Guard Agent | 5 | P0 | No | [a-05-hallucination-guard.md](a-05-hallucination-guard.md) |
| A-06 | Confidence Scoring Agent | 5 | P1 | No | [a-06-confidence-scoring.md](a-06-confidence-scoring.md) |
| A-07 | Drill-Down Context Agent | 6 | P1 | No | [a-07-drill-down-context.md](a-07-drill-down-context.md) |
| A-08 | Multi-Period Comparison Agent | 6 | P1 | No | [a-08-multi-period-comparison.md](a-08-multi-period-comparison.md) |
| A-09 | Report Scheduling Agent | 7 | P1 | No | [a-09-report-scheduling.md](a-09-report-scheduling.md) |
| A-10 | KPI Alert Agent | 7 | P1 | No | [a-10-kpi-alert.md](a-10-kpi-alert.md) |
| A-11 | Chart Dashboard Pin Agent | 7 | P2 | No | [a-11-chart-dashboard-pin.md](a-11-chart-dashboard-pin.md) |
| A-12 | Trend Forecasting Agent | 8 | P1 | No | [a-12-trend-forecasting.md](a-12-trend-forecasting.md) |
| A-13 | Anomaly Detection Agent | 8 | P1 | No | [a-13-anomaly-detection.md](a-13-anomaly-detection.md) |

---

## Module B — AI Accountant (16 agents)

| ID | Agent | Phase | Priority | HITL | Spec |
|----|-------|-------|----------|------|------|
| B-01 | Document Classification Agent | 2 | P0 | No | [b-01-document-classification.md](b-01-document-classification.md) |
| B-02 | OCR Extraction Agent | 2 | P0 | Yes | [b-02-ocr-extraction.md](b-02-ocr-extraction.md) |
| B-03 | Handwriting Recognition Agent | 3 | P1 | Yes (always) | [b-03-handwriting-recognition.md](b-03-handwriting-recognition.md) |
| B-04 | Vendor Matching Agent | 2 | P0 | Yes | [b-04-vendor-matching.md](b-04-vendor-matching.md) |
| B-05 | Ledger Mapping Agent | 2 | P0 | Yes | [b-05-ledger-mapping.md](b-05-ledger-mapping.md) |
| B-06 | Sub-Ledger / Cost Centre Agent | 3 | P1 | Yes | [b-06-sub-ledger-cost-centre.md](b-06-sub-ledger-cost-centre.md) |
| B-07 | Duplicate Invoice Detection Agent | 2 | P0 | Yes | [b-07-duplicate-invoice-detection.md](b-07-duplicate-invoice-detection.md) |
| B-08 | Approval Workflow Agent | 3 | P0 | Yes (always) | [b-08-approval-workflow.md](b-08-approval-workflow.md) |
| B-09 | Tally Sync Agent | 4 | P0 | Yes (always) | [b-09-tally-sync.md](b-09-tally-sync.md) |
| B-10 | Bank Reconciliation Agent | 4 | P0 | Yes | [b-10-bank-reconciliation.md](b-10-bank-reconciliation.md) |
| B-11 | Outstanding Items Agent | 5 | P1 | No | [b-11-outstanding-items.md](b-11-outstanding-items.md) |
| B-12 | Expense Categorisation Agent | 3 | P1 | Yes | [b-12-expense-categorisation.md](b-12-expense-categorisation.md) |
| B-13 | Tax Compliance Agent | 6 | P1 | No | [b-13-tax-compliance.md](b-13-tax-compliance.md) |
| B-14 | Audit Trail Logger Agent | 1 | P0 | No | [b-14-audit-trail-logger.md](b-14-audit-trail-logger.md) |
| B-15 | Cash Flow Forecasting Agent | 7 | P1 | No | [b-15-cash-flow-forecasting.md](b-15-cash-flow-forecasting.md) |
| B-16 | Multi-Entity Tally Agent | 6 | P2 | No | [b-16-multi-entity-tally.md](b-16-multi-entity-tally.md) |

---

## Module C — Easy Reports (8 agents)

| ID | Agent | Phase | Priority | HITL | Spec |
|----|-------|-------|----------|------|------|
| C-01 | Pre-Built Report Generator Agent | 8 | P1 | No | [c-01-pre-built-report-generator.md](c-01-pre-built-report-generator.md) |
| C-02 | Multi-Company Consolidation Agent | 8 | P2 | No | [c-02-multi-company-consolidation.md](c-02-multi-company-consolidation.md) |
| C-03 | Report Scheduling & Distribution Agent | 9 | P1 | No | [c-03-report-scheduling-distribution.md](c-03-report-scheduling-distribution.md) |
| C-04 | Custom Metric Formula Agent | 10 | P2 | No | [c-04-custom-metric-formula.md](c-04-custom-metric-formula.md) |
| C-05 | Row/Column Security Enforcement Agent | 10 | P0 | No | [c-05-row-column-security.md](c-05-row-column-security.md) |
| C-06 | Data Quality Validation Agent | 1 | P0 | No (blocks pipeline) | [c-06-data-quality-validation.md](c-06-data-quality-validation.md) |
| C-07 | Budget vs. Actual Variance Agent | 8 | P1 | No | [c-07-budget-vs-actual-variance.md](c-07-budget-vs-actual-variance.md) |
| C-08 | Inventory Aging & Reorder Agent | 8 | P2 | Yes (reorder) | [c-08-inventory-aging-reorder.md](c-08-inventory-aging-reorder.md) |

---

## Module D — Advanced Search Analytics (14 agents)

| ID | Agent | Phase | Priority | HITL | Spec |
|----|-------|-------|----------|------|------|
| D-01 | Natural Language Search Agent | 13 | P1 | No | [d-01-natural-language-search.md](d-01-natural-language-search.md) |
| D-02 | Entity Recognition Agent | 13 | P1 | No | [d-02-entity-recognition.md](d-02-entity-recognition.md) |
| D-03 | Multi-Step Conversational Analysis Agent | 14 | P1 | No | [d-03-multi-step-conversational.md](d-03-multi-step-conversational.md) |
| D-04 | Autonomous AI Analyst (Spotter) Agent | 14 | P1 | No (autonomous) | [d-04-autonomous-ai-analyst.md](d-04-autonomous-ai-analyst.md) |
| D-05 | Deep Research Agent | 14 | P2 | No | [d-05-deep-research.md](d-05-deep-research.md) |
| D-06 | Dashboard Auto-Generation Agent | 15 | P2 | No | [d-06-dashboard-auto-generation.md](d-06-dashboard-auto-generation.md) |
| D-07 | Insight Discovery & Prioritisation Agent | 15 | P1 | No | [d-07-insight-discovery.md](d-07-insight-discovery.md) |
| D-08 | Prescriptive Recommendations Agent | 15 | P1 | No (recommends only) | [d-08-prescriptive-recommendations.md](d-08-prescriptive-recommendations.md) |
| D-09 | Semantic Layer Management Agent | 13 | P0 | Yes (governance) | [d-09-semantic-layer-management.md](d-09-semantic-layer-management.md) |
| D-10 | Insight-to-Action Workflow Agent | 14 | P1 | Yes (always) | [d-10-insight-to-action-workflow.md](d-10-insight-to-action-workflow.md) |
| D-11 | Code Generation Agent (SpotterCode) | 16 | P3 | No | [d-11-code-generation.md](d-11-code-generation.md) |
| D-12 | Federated Query Optimisation Agent | 16 | P2 | No | [d-12-federated-query-optimisation.md](d-12-federated-query-optimisation.md) |
| D-13 | Scheduled Autonomous Monitoring Agent | 14 | P1 | No | [d-13-scheduled-autonomous-monitoring.md](d-13-scheduled-autonomous-monitoring.md) |
| D-14 | Voice/Mobile Search Agent | 15 | P2 | No | [d-14-voice-mobile-search.md](d-14-voice-mobile-search.md) |

---

## Module E — Language & Localisation (7 agents)

> Cross-cutting module: E-agents run as pre/post-processing steps for all other modules. Default language: **English**. Phase 1 ships **Arabic (RTL)**.  
> See [docs/i18n-architecture.md](../../i18n-architecture.md) and [PRD §6.10](../../PRD.md) for full specification.

| ID | Agent | Phase | Priority | HITL | Spec |
|----|-------|-------|----------|------|------|
| E-01 | Language Detection & Routing Agent | 2 | P0 | No | [e-01-language-detection-routing.md](e-01-language-detection-routing.md) |
| E-02 | Query Translation Agent | 2 | P0 | No | *(spec pending)* |
| E-03 | Localised Response Formatter | 2 | P0 | No | *(spec pending)* |
| E-04 | Multilingual Report Generator | 5 | P1 | No | *(spec pending)* |
| E-05 | Translation Coverage Guard (CI) | 1 | P0 | No | *(spec pending)* |
| E-06 | Multilingual Notification Agent | 6 | P1 | No | *(spec pending)* |
| E-07 | Bilingual Glossary Sync Agent | 3 | P1 | Yes (review) | *(spec pending)* |

---

## Summary

| Module | Agents | HITL Agents | P0 Critical |
|--------|--------|-------------|-------------|
| A — Conversational BI | 13 | 0 | 3 |
| B — AI Accountant | 16 | 11 | 7 |
| C — Easy Reports | 8 | 2 | 3 |
| D — Advanced Search Analytics | 14 | 2 | 1 |
| **E — Language & Localisation** | **7** | **1** | **4** |
| **Total** | **58** | **16** | **18** |
