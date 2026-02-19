# Agent Specification — A-04: Domain Terminology Normalisation Agent

**Agent ID:** `A-04`  
**Agent Name:** Domain Terminology Normalisation Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 2  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Maps healthcare and accounting synonyms in user queries to the canonical column and table names used in the data warehouse before the query reaches the SQL agent. Prevents A-01 from hallucinating column names due to terminology mismatch.

> **Addresses:** PRD §5.4 — Domain synonym normalisation (e.g. "footfall" → `patient_visits`, "outstanding" → `accounts_receivable`).

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | A-01 (first step in pipeline) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `raw_query` | `string` | User input | ✅ |
| `synonym_registry` | `map[string]string` | Schema Context Cache | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `normalised_query` | `string` | A-01 LLM prompt |
| `substitutions_made` | `[]Substitution` | Debug log |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Go string matcher | Internal | Fast exact/fuzzy term lookup |
| 2 | Synonym registry (YAML config) | Internal | Domain-specific term mappings |
| 3 | Genkit Flow (optional) | Apache-2.0 | LLM fallback for unknown terms |

### Synonym Registry Examples
```yaml
synonyms:
  "footfall": "patient_visits"
  "outstanding": "accounts_receivable"
  "pharmacy sales": "pharmacy_dispensations.amount"
  "bed occupancy": "inpatient_bed_utilisation_rate"
  "cost of goods": "cost_of_sales"
  "outstanding dues": "accounts_payable"
```

---

## 6. Guardrails

- If a term is not in the registry and LLM confidence < 0.80, pass the raw term through unchanged (don't guess).
- All substitutions logged so they can be reviewed and added to registry.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Term match rate on known synonyms | 100% |
| False substitution rate | < 0.5% |
| Latency (registry lookup) | < 10ms |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (embedded in A-01 flow) |
| **Config** | `config/synonym_registry.yaml` |
| **Depends on** | Schema Context Cache |
| **Consumed by** | A-01 |
