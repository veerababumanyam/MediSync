# Phase 10 — Customization & RBAC Security (Easy Reports — Part 3)

**Phase Duration:** Weeks 34–36 (3 weeks)  
**Module(s):** Module C (Easy Reports)  
**Status:** Planning  
**Milestone:** M7 — Reports Module (complete)  
**Depends On:** Phase 09 complete (scheduling and delivery live)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §8](../ARCHITECTURE.md)

---

## 1. Objectives

Complete the Easy Reports module with the zero-code drag-and-drop report builder, custom metric formula creation, and enterprise-grade row/column security enforcement. After this phase, non-technical business users can build their own reports, and all data access is strictly governed by role.

---

## 2. Scope

### In Scope
- C-04 Custom Metric Formula Agent
- C-05 Row/Column Security Enforcement Agent
- Zero-code drag-and-drop report builder UI
- Formula builder for calculated fields
- User-defined dimensions and custom fields (Tally UDFs)
- Cost-centre and department analytics module
- Role-based access control for reports (view/export/schedule permissions)
- Column masking (cost/margin fields by role)
- Row-level filtering (department/region scoping)
- Conditional formatting (RAG thresholds on cells)
- Custom KPI dashboard creator for end users

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | C-04 Custom Metric Formula Agent | AI Engineer | Users define formula in natural language → agent generates MetricFlow metric definition; tested with 10 formulas |
| D-02 | C-05 Row/Column Security Agent | AI Engineer + DevOps | OPA row-filter per department; column masking for cost/margin fields; verified for all 8 user roles |
| D-03 | Zero-Code Report Builder UI | Frontend Engineer | Drag fields from data model → configure report; save and share; no SQL knowledge required |
| D-04 | Formula Builder | Frontend Engineer | GUI for calculated fields: operations (+, -, *, /, %), functions (SUM, AVG, COUNT, IF); preview with sample data |
| D-05 | Custom KPI Dashboard Creator | Frontend Engineer | Users add/remove/rearrange KPI cards and charts; save as personal dashboard |
| D-06 | Tally UDF Support | Backend Engineer | User-defined fields from Tally extracted during ETL; available in report builder |
| D-07 | Cost-Centre Analytics Module | Frontend + AI | Department-wise P&L; overhead allocation; cost-centre profitability dashboard |
| D-08 | Report Permission Management | Frontend + Backend | Per-report: view/export/schedule/share permissions assignable by admin |

---

## 4. AI Agents Deployed

### C-04 Custom Metric Formula Agent

**Purpose:** Let non-technical users create governed metrics by describing them in plain language.

**Examples:**
| User Input | Agent Output (MetricFlow) |
|---|---|
| "Revenue minus cost of goods sold divided by revenue as a percentage" | `gross_margin_pct = (revenue - cogs) / revenue * 100` |
| "Number of unique patients who visited in the last 30 days" | `active_patients_30d = COUNT(DISTINCT patient_id WHERE appt_date >= NOW()-30)` |
| "Total salary expense as a percentage of total revenue" | `salary_cost_ratio = salary_ledger_sum / total_revenue * 100` |

**Workflow:**
1. User describes metric in natural language
2. Agent generates MetricFlow YAML definition
3. Agent shows preview with sample data ("This metric would return: 23.4%")
4. User can tweak or approve
5. On approve: metric registered to `semantic.metric_definitions` by D-09 (Phase 13) or directly in Phase 10 via MetricFlow CLI

**Governance:** New custom metrics from non-admin users are flagged for data governance team review before publishing to shared metric library.

### C-05 Row/Column Security Enforcement Agent

**Implementation:** OPA Rego policies applied as a post-processing filter on every report/dashboard query.

**Row-level filtering:**

```rego
package medisync.data

# Apply department filter for non-admin users
row_filter[filter] if {
    input.user.role != "admin"
    input.user.role != "finance_head"
    input.resource_type in ["report", "dashboard"]
    filter := {
        "column": "department",
        "value": input.user.department
    }
}

# Pharmacy manager sees only pharmacy data
row_filter[filter] if {
    input.user.role == "pharmacy_manager"
    filter := {"column": "department", "value": "pharmacy"}
}
```

**Column masking:**

```rego
# Mask cost and margin columns for viewer/manager roles
masked_columns[col] if {
    input.user.role in ["viewer", "manager"]
    col := {"cost_price", "gross_margin_pct", "cogs", "net_profit_amount"}[_]
}

# Mask patient PII for non-clinical roles
masked_columns[col] if {
    not input.user.role in ["admin", "clinic_manager", "analyst"]
    col := {"patient_name", "patient_phone", "patient_dob"}[_]
}
```

**Application:** C-05 runs after C-01 / report builder generates SQL; before results are returned to the UI.

---

## 5. Zero-Code Report Builder

### Data Model Panel (Left)
Hierarchical tree of available fields:
```
Tally Data
  └── Ledgers
       ├── Ledger Name
       ├── Closing Balance
       └── ...
  └── Vouchers
       ├── Date
       ├── Amount
       └── ...
HIMS Data
  └── Appointments
       ├── Patient Name (masked for some roles)
       ├── Doctor
       └── ...
  └── Pharmacy
       └── ...
```

### Report Canvas (Centre)
- Drop fields from Data Model panel onto canvas
- Rows: dimension fields (group-by)
- Values: measure fields (sum/avg/count/min/max)
- Filters: date range, department, entity
- Sort: click column header

### Configuration Panel (Right)
- Report title (EN + AR)
- Chart toggle: Table | Bar | Line | Pie | Mixed
- Conditional formatting: set threshold + colour rules
- Subtotals and grand totals toggle
- Pagination (rows per page)

### Saved Report Actions
- Save as personal report
- Share with team (read-only link or assign to role)
- Schedule via C-03 (link to Phase 9 scheduler)
- Export (PDF / Excel / CSV)

---

## 6. Custom KPI Dashboard Creator

**User flow:**
1. Click "New Dashboard" → name + description
2. Add KPI card: select metric, display name, target value, RAG thresholds
3. Add chart widget: select report type, chart type, period
4. Drag to reorder, resize grid cells
5. Save → appears in "My Dashboards" sidebar

**Dashboard sharing:** Share with team members (read access) or keep private.

---

## 7. Cost-Centre Analytics Module

**Screens:**
1. **Cost Centre P&L** — revenue, expenses, gross margin per cost centre with period selector
2. **Overhead Allocation** — configurable allocation methods: fixed %, headcount-based, revenue-based
3. **Cost Centre Comparison** — multi-cost-centre bar chart comparison
4. **Drill-down to Vouchers** — click any cost centre metric → transaction list

---

## 8. Tally User-Defined Fields (UDFs)

**ETL enhancement:** During Phase 10, the Tally connector is extended to extract UDF metadata:
- List of custom field definitions from Tally schema
- UDF values per voucher / ledger / stock item

**Report builder:** UDFs appear in Data Model panel under "Custom Fields (Tally)"; fully usable in report design.

---

## 9. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| C-04 formula generation | 10 user-described metrics | All 10 generate correct MetricFlow definitions; preview values match manual |
| C-05 row filtering | 6 user roles × 3 report types | Correct row scope enforced for each role/report combination |
| C-05 column masking | Cost/margin columns for viewer + manager roles | Masked in report output; not accessible via API |
| Report builder | Build 5 custom reports without SQL | All 5 render correctly with live data |
| Formula builder | 5 computed fields | Correct results on sample data |
| OPA bypass attempt | Direct API call to bypass C-05 | Request blocked by OPA policy |
| Cost-centre P&L | 3 cost centres with known values | Figures match manual calculation |

---

## 10. Phase Exit Criteria

- [ ] Zero-code report builder functional for non-technical users (validated by accountant persona UAT)
- [ ] C-04 Custom Metric Formula Agent generating correct MetricFlow definitions
- [ ] C-05 Row/Column Security enforced for all 8 user roles
- [ ] OPA bypass attempts blocked
- [ ] Custom KPI dashboard creator working
- [ ] Tally UDF fields available in report builder
- [ ] Cost-centre analytics module live
- [ ] M7 Reports Module complete milestone declared
- [ ] Phase gate reviewed and signed off

---

*Phase 10 | Version 1.0 | February 19, 2026*
