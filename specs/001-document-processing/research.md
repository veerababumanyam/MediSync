# Research: Document Processing Pipeline

**Feature**: 001-document-processing
**Date**: 2026-02-21
**Purpose**: Resolve technical unknowns and establish best practices for OCR pipeline implementation

---

## 1. OCR Engine Selection

### Decision: PaddleOCR

**Rationale**:
- Apache-2.0 license (constitution-compliant)
- Native Arabic language support with pre-trained models
- Supports 80+ languages including mixed-language documents
- Handwriting recognition capability
- High accuracy on printed text (95%+ field-level)
- Active community and regular updates

**Alternatives Considered**:

| Engine | License | Arabic | Handwriting | Rejected Because |
|--------|---------|--------|-------------|------------------|
| Tesseract | Apache-2.0 | Limited | No | Poor Arabic accuracy, no handwriting |
| EasyOCR | Apache-2.0 | Yes | Limited | Slower inference, higher memory |
| Google Vision | Proprietary | Yes | Yes | Not OSI-approved, cloud dependency |
| AWS Textract | Proprietary | Limited | No | Not OSI-approved, cloud-only |

### Configuration

```yaml
# PaddleOCR configuration for MediSync
ocr:
  use_angle_cls: true        # Handle rotated documents
  lang: "en,ar"              # English + Arabic
  det_model: "PP-OCRv4"      # Latest detection model
  rec_model: "PP-OCRv4"      # Latest recognition model
  use_gpu: true              # GPU acceleration for batch
  gpu_mem: 4096              # 4GB GPU memory
```

---

## 2. Arabic OCR Support

### Decision: PP-OCRv4 Arabic Model with RTL Post-Processing

**Rationale**:
- PP-OCRv4 includes pre-trained Arabic model
- Bidirectional text handling via ICU library
- RTL display handled at frontend layer (i18next + Tailwind logical properties)

**Implementation Approach**:

1. **Detection**: PaddleOCR detects text regions regardless of language
2. **Recognition**: Arabic model recognizes Arabic characters
3. **Direction Detection**: ICU BiDi algorithm determines text direction per field
4. **Storage**: Store text in logical order (not visual order)
5. **Display**: Frontend applies RTL via `dir="rtl"` and Tailwind logical properties

**Arabic-Specific Considerations**:
- Ligature handling (لا, لله, etc.) - PaddleOCR handles natively
- Numerals: Arabic numerals (٠١٢٣) vs Arabic-Indic (0123) - store as-is, format at display
- Date formats: Support both Hijri and Gregorian calendars

---

## 3. Handwriting Recognition

### Decision: PaddleOCR + Confidence Capping at 85%

**Rationale**:
- Handwriting accuracy is inherently lower than printed text
- Capping confidence ensures human review for all handwritten fields
- Using same engine simplifies architecture

**Implementation Approach**:

1. **Handwriting Detection**: Analyze stroke patterns and consistency
2. **Confidence Adjustment**: Apply 0.85 multiplier to OCR confidence
3. **Threshold**: Cap maximum confidence at 85% for handwritten fields
4. **Flagging**: Mark fields as `is_handwritten: true` in extraction results

**Confidence Thresholds**:
| Field Type | Auto-Accept | Flag for Review | High Priority |
|------------|-------------|-----------------|---------------|
| Printed | ≥95% | 70-94% | <70% |
| Handwritten | N/A (capped at 85%) | All flagged | <70% |

---

## 4. Confidence Scoring Algorithm

### Decision: Multi-Factor Confidence Score

**Formula**:
```
final_confidence = ocr_confidence * (1 - validation_penalty) * handwriting_multiplier
```

**Factors**:

| Factor | Weight | Description |
|--------|--------|-------------|
| OCR Confidence | 60% | Raw PaddleOCR confidence score |
| Format Validation | 20% | Regex match for dates, amounts, invoice numbers |
| Cross-Field Validation | 20% | Subtotal + Tax = Total, date consistency |
| Handwriting Multiplier | 0.85 | Applied if field detected as handwritten |

**Validation Rules**:

| Field Type | Validation | Penalty if Invalid |
|------------|------------|-------------------|
| Invoice Number | Alphanumeric pattern | +0.1 |
| Date | Valid date format, not future | +0.15 |
| Amount | Numeric, positive | +0.1 |
| Total | Sum validation | +0.2 |
| GST/VAT Number | Country-specific format | +0.05 |

---

## 5. Document Storage Architecture

### Decision: PostgreSQL Metadata + S3-Compatible Object Storage

**Rationale**:
- Large files (up to 25MB) should not be stored in database
- S3-compatible storage allows flexible deployment (MinIO on-prem, AWS S3 cloud)
- PostgreSQL stores metadata, extracted fields, and embeddings for search

**Storage Layout**:

```
S3 Bucket: medisync-documents/
├── {tenant_id}/
│   ├── {document_id}/
│   │   ├── original.pdf      # Original uploaded file
│   │   ├── page_001.png      # Rendered pages for preview
│   │   ├── page_002.png
│   │   └── thumbnails/
│   │       ├── page_001_thumb.png
│   │       └── page_002_thumb.png
```

**Encryption**:
- At rest: AES-256 via S3 server-side encryption
- In transit: TLS 1.3
- Key management: Tenant-specific KMS keys

---

## 6. Async Processing with NATS JetStream

### Decision: NATS JetStream with Document Processing Queue

**Rationale**:
- NATS already in constitution stack
- JetStream provides durable queues with replay
- Supports retry with backoff for failed OCR jobs

**Queue Design**:

```yaml
# JetStream stream configuration
stream:
  name: DOCUMENT_PROCESSING
  subjects: documents.>
  retention: limits
  max_msgs: 100000
  max_age: 7d
  duplicates: 5m

consumers:
  - name: ocr-worker
    durable: true
    ack_policy: explicit
    max_deliver: 3
    backoff: [10s, 30s, 60s]
```

**Processing Flow**:

```
1. Upload → API validates → Store in S3 → Publish "documents.uploaded"
2. OCR Worker subscribes → Retrieves file → Runs OCR → Publishes "documents.extracted"
3. API receives "documents.extracted" → Updates database → Notifies frontend via WebSocket
```

**Error Handling**:

| Error Type | Action | Retry |
|------------|--------|-------|
| Transient OCR failure | Retry with backoff | Yes (3x) |
| Corrupted file | Mark as failed, notify user | No |
| Unsupported format | Reject at upload | No |
| Processing timeout | Retry once | Yes (1x) |

---

## 7. Multi-Page Document Processing

### Decision: Page-by-Page Extraction with Consolidation

**Rationale**:
- Each page processed independently for parallelization
- Results consolidated into single document view
- Page boundaries preserved for reference

**Consolidation Rules**:

1. **Header fields** (supplier, invoice number): Take from first page
2. **Line items**: Concatenate from all pages with page reference
3. **Totals**: Take from last page (summary page)
4. **Conflicts**: Flag for human review if pages disagree on header fields

---

## 8. Document Classification

### Decision: Rule-Based Classification with AI Fallback

**Rationale**:
- Rule-based is fast and deterministic for common document types
- AI fallback handles edge cases
- Classification accuracy target: 98%

**Classification Signals**:

| Document Type | Primary Signals |
|---------------|-----------------|
| Invoice | "Invoice", "Bill No", supplier name, line items |
| Receipt | "Receipt", "Paid", payment method |
| Bank Statement | "Bank", "Account", transaction table, "Balance" |
| Expense Report | Employee name, expense categories, approvals |

**AI Fallback**:
- When rule-based confidence < 80%
- Use Genkit flow with document preview
- Returns classification with confidence

---

## 9. Audit Trail Requirements

### Decision: Immutable Audit Log with Full Context

**Audit Events**:

| Event | Data Captured |
|-------|---------------|
| document.uploaded | user_id, file_name, size, document_id |
| document.classified | document_type, confidence |
| document.extracted | field_count, avg_confidence |
| field.edited | field_name, old_value, new_value, user_id |
| document.approved | user_id, all field values snapshot |
| document.rejected | user_id, reason |

**Storage**: PostgreSQL `audit_log` table with append-only permissions

---

## Summary of Decisions

| Area | Decision | License |
|------|----------|---------|
| OCR Engine | PaddleOCR PP-OCRv4 | Apache-2.0 |
| Arabic Support | PP-OCRv4 Arabic model + ICU BiDi | Apache-2.0 / ICU |
| Handwriting | PaddleOCR with 85% confidence cap | Apache-2.0 |
| Confidence Scoring | Multi-factor with validation penalties | N/A |
| Document Storage | S3-compatible + PostgreSQL | Apache-2.0 / PostgreSQL |
| Async Processing | NATS JetStream | Apache-2.0 |
| Classification | Rule-based + AI fallback | N/A |
| Audit Trail | PostgreSQL append-only table | PostgreSQL |

All selected technologies comply with constitution requirements (OSI-approved licenses).
