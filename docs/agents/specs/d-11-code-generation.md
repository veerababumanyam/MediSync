# Agent Specification — D-11: Code Generation Agent (SpotterCode)

**Agent ID:** `D-11`  
**Agent Name:** Code Generation Agent (SpotterCode)  
**Module:** D — Advanced Search Analytics  
**Phase:** 16  
**Priority:** P3 Low  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Generates custom analytics code (Python/Go data scripts, SQL queries, Excel formulas, report templates) from natural language specifications, enabling power users to extend MediSync without traditional coding.

> **Addresses:** PRD §6.9.8, US33 — No-code / low-code extensions using AI code generation.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | "SpotterCode" tab in Analytics UI |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `specification` | `string` | User input (natural language) | ✅ |
| `target_language` | `enum` | `python / sql / go / excel` | ✅ |
| `context_schema` | `string` | D-09 Semantic Layer | ⬜ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `generated_code` | `string` | Generated code artifact |
| `explanation` | `string` | Line-by-line explanation |
| `safety_check` | `SafetyResult` | SAST scan result |

---

## 5. Tool Chain

| Step | Tool | Purpose |
|------|------|---------|
| 1 | Genkit Flow (`code-gen`) | Code generation from specification |
| 2 | D-09 Semantic Layer | Provide schema context |
| 3 | Static analysis (Go AST / sqlfluff) | Syntax + safety checks |
| 4 | Sandbox runner (isolated) | Optional test execution |

---

## 6. Guardrails

- Generated SQL always validated by A-05 for read-only compliance.
- No generated code runs in production automatically — user must explicitly execute.
- Code execution (if enabled) happens in an isolated sandbox with no network access.
- Only `analyst`, `admin` roles can access SpotterCode.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Code correctness (syntax valid) | ≥ 99% |
| User satisfaction | ≥ 4/5 rating |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Depends on** | D-09, Genkit |
| **Consumed by** | Power users (Analyst, Admin) |
