# Agent Specification — B-02: OCR Extraction Agent

**Agent ID:** `B-02`  
**Agent Name:** OCR Extraction Agent  
**Module:** B — AI Accountant  
**Phase:** 4  
**Priority:** P0 Critical  
**HITL Required:** Yes — overall confidence < 0.85 or any required field < 0.70  
**Status:** Draft

---

## 1. Purpose

Extracts structured financial fields (amount, vendor, date, invoice number, tax) from uploaded documents (PDFs, images, Excel) with per-field confidence scoring. Flags low-confidence items for human review before any data proceeds to the ledger mapping pipeline.

> **Addresses:** PRD §6.7.1, US9, US13 — Document ingestion and structured data extraction.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Event-driven |
| **Event trigger** | B-01 output: `document.classified` event |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `file` | `bytes` | Upload store | ✅ |
| `file_type` | `enum` | B-01 | ✅ |
| `document_type` | `enum` | B-01 | ✅ |
| `extraction_schema` | `JSON` | Config: required fields per doc type | ✅ |
| `upload_id` | `UUID` | Upload service | ✅ |
| `user_id` | `string` | JWT | ✅ |
| `session_id` | `string` | Session store | ✅ |

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `extracted_fields` | `map[string]any` | `{vendor, amount, invoice_date, invoice_no, tax_amount, currency}` |
| `confidence_scores` | `map[string]float64` | Per-field confidence |
| `overall_confidence` | `float64` | Aggregated score |
| `low_confidence_flags` | `[]string` | Fields below threshold |
| `raw_text` | `string` | Full OCR text |
| `hitl_required` | `bool` | Trigger review if true |
| `trace_id` | `string` | OTel trace |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | PyMuPDF | AGPL-3.0 | Digital PDF text layer extraction |
| 2 | PaddleOCR | Apache-2.0 | Scanned / image OCR + layout |
| 3 | Tesseract | Apache-2.0 | Fallback OCR for simple docs |
| 4 | Unstructured.io | Apache-2.0 | Table + header segmentation |
| 5 | Genkit Flow (`field-extract`) | Apache-2.0 | LLM structured field extraction |
| 6 | A-06 Confidence Scorer | Internal | Per-field + aggregated score |
| 7 | HITL router | Internal | Queue if below threshold |

```
File
  → B-01 classification (doc type known)
  → Router: digital PDF → PyMuPDF | scanned → PaddleOCR | handwritten → B-03
  → Unstructured.io segmenter
  → Genkit extraction flow
  → Confidence scoring
  → HITL gate (if overall < 0.85)
  → ExtractionResult struct → B-05
```

---

## 6. System Prompt

```
You are a financial document extraction specialist. Extract the following fields from the document text:
- vendor_name
- invoice_number
- invoice_date (ISO 8601)
- subtotal_amount
- tax_amount
- total_amount
- currency (ISO 4217)
- line_items (array of {description, quantity, unit_price, total})

Document type: {{ document_type }}
Document text:
{{ raw_text }}

For each field, provide a confidence score (0.0–1.0) based on how clearly it appears in the text.
Respond ONLY with valid JSON matching the extraction schema.
```

---

## 7. Guardrails

| # | Guard | Trigger | Action |
|---|-------|---------|--------|
| 1 | File size limit | > 50 MB | Reject with user error |
| 2 | MIME type whitelist | Invalid type | Reject |
| 3 | Required field check | `amount` or `vendor` or `date` missing | Set `hitl_required=true` |
| 4 | Confidence gate | `overall_confidence < 0.85` | Queue for human review |
| 5 | Data not forwarded to Tally | Always | Only proceeds via B-05 → B-08 → B-09 |
| 6 | Encryption at rest | Upload complete | Files stored AES-256 |

---

## 8. HITL Gate

| Property | Value |
|----------|-------|
| **Trigger** | `overall_confidence < 0.85` OR `confidence_scores[required_field] < 0.70` |
| **Notified role** | `accountant` |
| **Notification** | In-app: flagged fields highlighted in extraction review UI |
| **SLA** | 24h |
| **Approval actions** | Confirm / Edit + Confirm |
| **On confirmation** | ExtractionResult proceeds to B-05 |

---

## 9. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Field extraction accuracy (printed) | ≥ 95% |
| Field extraction accuracy (handwritten, via B-03) | ≥ 90% |
| Wrong document type routing | < 3% |
| P95 processing time per page | < 10s |
| HITL escalation rate | < 20% |

---

## 10. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service + Python OCR sidecar (PaddleOCR) |
| **Queue** | Redis `ocr-extraction-queue` |
| **Depends on** | B-01, PaddleOCR service, Unstructured.io |
| **Consumed by** | B-03 (handwriting), B-05 (ledger mapping) |
| **Env vars** | `PADDLEOCR_SERVICE_URL`, `UNSTRUCTURED_API_URL` |
