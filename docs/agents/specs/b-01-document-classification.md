# Agent Specification — B-01: Document Classification Agent

**Agent ID:** `B-01`  
**Agent Name:** Document Classification Agent  
**Module:** B — AI Accountant  
**Phase:** 4  
**Priority:** P0 Critical  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Auto-categorises uploaded financial documents into one of the known types (Invoice, Bank Statement, Bill, Receipt, Tax Document, Other) before routing to the appropriate extraction pipeline.

> **Addresses:** PRD §6.7.1 — Smart document classification on upload.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | `document.uploaded` event on upload queue |
| **Calling agent** | User (file upload UI) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `file_bytes` | `[]byte` | Upload queue | ✅ |
| `file_name` | `string` | Upload metadata | ✅ |
| `mime_type` | `string` | Upload metadata | ✅ |
| `user_id` | `string` | JWT | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `document_type` | `enum` | `invoice / bill / bank_statement / receipt / tax_doc / other` |
| `confidence_score` | `float64` | Classification confidence |
| `routing_target` | `string` | Next agent ID (always B-02) |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | MIME type validator (Go) | Internal | Security: reject disallowed types |
| 2 | PaddleOCR layout analysis | Apache-2.0 | Detect document structure |
| 3 | Go HTTP client → classification sidecar | Internal | Call ML classifier |
| 4 | Fine-tuned LayoutLM (ONNX) | Apache-2.0 | Document layout classification |

---

## 6. Guardrails

- Allowed MIME types: `application/pdf`, `image/png`, `image/jpeg`, `image/tiff`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`, `text/csv`
- Max file size: 50 MB
- `document_type=other` → user prompted to manually confirm document type before continuing.
- All uploads encrypted (AES-256) at rest on object storage.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Classification accuracy | ≥ 97% |
| P95 Latency | < 3s |
| `other` rate (unrecognised) | < 2% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service |
| **Queue** | Redis `document.upload` queue |
| **Depends on** | PaddleOCR sidecar, LayoutLM ONNX model |
| **Consumed by** | B-02 |
