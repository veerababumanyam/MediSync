# Phase 04 — Document Processing (AI Accountant — Part 1)

**Phase Duration:** Weeks 12–15 (4 weeks)  
**Module(s):** Module B (AI Accountant), Module E  
**Status:** Planning  
**Milestone:** M4 — OCR Pipeline Live  
**Depends On:** Phase 03 complete (UI framework in place)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.2](../ARCHITECTURE.md) | [agents/specs/b-01-document-classification.md](../agents/specs/b-01-document-classification.md)

---

## 1. Objectives

Launch the first half of the AI Accountant module: the document ingestion and OCR pipeline. Users can upload bank statements, vendor bills, invoices (including handwritten), and the system automatically classifies and extracts structured data from them with confidence scoring. This phase transforms raw source documents into machine-readable records ready for the transaction intelligence pipeline (Phase 5).

---

## 2. Scope

### In Scope
- Document upload drag-and-drop UI (React + Flutter)
- Bulk upload batch processing pipeline
- B-01 Document Classification Agent
- B-02 OCR Extraction Agent (PaddleOCR)
- B-03 Handwriting Recognition Agent
- Low-confidence human review queue UI
- Document library and search
- NATS event topics for document pipeline
- Secure document storage (AES-256 encrypted)
- Arabic invoice OCR support
- OCR for: PDF, JPEG, PNG, TIFF, Excel/CSV (bank statements)

### Out of Scope
- Ledger mapping (Phase 5)
- Tally sync (Phase 6)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | Document Upload UI (Web) | Frontend Engineer | Drag-and-drop; multi-file; format validation; real-time upload progress bar |
| D-02 | Document Upload API | Backend Engineer | `POST /v1/documents/upload`; stores encrypted file; emits `document.uploaded` NATS event |
| D-03 | Batch Processing Queue | Backend Engineer | Consumes `document.uploaded`; queues documents for classification; handles 100+ simultaneous uploads without bottleneck |
| D-04 | B-01 Document Classification Agent | AI Engineer | Classifies Invoice/Bill/Bank Statement/Receipt/Tax Doc/Other with ≥ 95% accuracy on 100-doc test set |
| D-05 | B-02 OCR Extraction Agent | AI Engineer | Extracts amount, vendor, date, invoice#, tax_amount with ≥ 95% accuracy on printed docs; confidence score per field |
| D-06 | B-03 Handwriting Recognition Agent | AI Engineer | ≥ 90% accuracy on handwritten invoices test set (30 docs); Arabic handwriting supported |
| D-07 | HITL Review Queue UI | Frontend Engineer | Low-confidence fields highlighted; accountant can correct inline; correction feeds back to system |
| D-08 | Document Library UI | Frontend Engineer | List view with filter by type/date/vendor; search; preview; link to transaction |
| D-09 | Secure Document Storage | DevOps | AES-256 per-file encryption; stored on on-premises file server; encryption key in Vault |
| D-10 | Arabic OCR Support | AI Engineer | PaddleOCR Arabic model loaded; Arabic invoices extracted correctly in test set |
| D-11 | Upload Mobile UI (Flutter) | Frontend Engineer | Photo capture upload from mobile camera; batch selection from photo library |

---

## 4. AI Agents Deployed

### B-01 Document Classification Agent

**Trigger:** NATS subscription to `document.uploaded`  
**Model approach:**  
1. Load first page as image
2. Extract text via lightweight OCR pass
3. Classify on: document structure, field patterns, key phrases ("Invoice #", "Bank Statement", "Credit", "Purchase")
4. Emit NATS `document.classified` with class and confidence

**Classes:**
| Class | Key Identifiers |
|---|---|
| `invoice` | Invoice number, vendor GST, line items |
| `purchase_bill` | "Bill To", vendor address, PO reference |
| `bank_statement` | Account number, bank name, balance, debit/credit columns |
| `receipt` | "Received from", single-line amount, payment mode |
| `tax_document` | GST number, tax period, GSTIN references |
| `other` | Fallback; flagged for human review |

**Accuracy target:** ≥ 95% on 100-document test set covering all classes

### B-02 OCR Extraction Agent

**Trigger:** NATS subscription to `document.classified`  
**Engine:** PaddleOCR (Apache-2.0) — supports English + Arabic scripts  
**Go service wrapper:** PaddleOCR Python service called as microservice via gRPC from Go

**Extraction pipeline:**
```
Document (PDF/Image)
    │
    ▼ Pre-processing: deskew, denoise, contrast enhancement
    │
    ▼ PaddleOCR: full-page text + bounding boxes
    │
    ▼ Field extractor (regex + LLM hybrid):
    │     - amount        (regex: currency patterns)
    │     - vendor_name   (LLM: entity detection)
    │     - invoice_date  (regex: date patterns)
    │     - invoice_number (regex: INV-###, #/##/####)
    │     - tax_amount    (regex: GST/VAT patterns)
    │
    ▼ Confidence score per field (0–100%)
    │
    ▼ Low-confidence fields (< 70%) → HITL queue
    │
    ▼ Store extraction in app.extracted_documents
    │
    ▼ Emit NATS: document.extracted
```

**Output schema:**
```json
{
  "document_id": "uuid",
  "vendor_name": {"value": "ABC Medical Supplies", "confidence": 92},
  "invoice_date": {"value": "2026-01-15", "confidence": 88},
  "invoice_number": {"value": "INV-2026-0042", "confidence": 97},
  "amount": {"value": 15750.00, "confidence": 99},
  "tax_amount": {"value": 2362.50, "confidence": 95},
  "currency": "INR",
  "needs_human_review": false,
  "low_confidence_fields": []
}
```

**Accuracy targets:**
- Standard printed invoices: ≥ 95% per-field accuracy
- Arabic printed invoices: ≥ 90% per-field accuracy
- Handwritten invoices: handled by B-03

### B-03 Handwriting Recognition Agent

**Trigger:** B-02 detects handwritten document (detected via ink stroke pattern analysis or low OCR confidence across all fields)  
**Approach:**
1. Enhanced image pre-processing (binarization, morphological operations)
2. PaddleOCR handwriting model (PP-OCRv4 handwriting variant)
3. LLM post-processing: "Given these OCR text fragments from a handwritten invoice, extract the amount, vendor, date, and invoice number. If a field is illegible, return null."
4. Confidence score based on LLM certainty + OCR character confidence
5. Low confidence (< 70% on any required field) → always HITL

**Arabic handwriting:** Supported via PaddleOCR Arabic handwriting model  
**Target accuracy:** ≥ 90% on test set of 30 handwritten invoices (English + Arabic)  
**HITL:** ALL handwritten documents that have at least one field < 90% confidence go to human review queue

---

## 5. Database Schema Additions

```sql
-- Schema: app
CREATE TABLE documents (
    doc_id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID REFERENCES users(user_id),
    file_path       TEXT NOT NULL,            -- encrypted on-disk path
    file_type       VARCHAR(20),              -- pdf, jpeg, png, csv, xlsx
    doc_class       VARCHAR(50),              -- invoice, bank_statement, etc.
    class_confidence NUMERIC(5,2),
    upload_status   VARCHAR(30) DEFAULT 'uploaded', -- uploaded|classifying|extracting|extracted|reviewed|linked
    uploaded_at     TIMESTAMPTZ DEFAULT NOW(),
    extracted_at    TIMESTAMPTZ,
    reviewed_at     TIMESTAMPTZ,
    reviewed_by     UUID REFERENCES users(user_id)
);

CREATE TABLE extracted_documents (
    extraction_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doc_id          UUID REFERENCES documents(doc_id),
    vendor_name     TEXT,
    vendor_confidence NUMERIC(5,2),
    invoice_date    DATE,
    date_confidence NUMERIC(5,2),
    invoice_number  TEXT,
    invoice_number_confidence NUMERIC(5,2),
    amount          NUMERIC(15,2),
    amount_confidence NUMERIC(5,2),
    tax_amount      NUMERIC(15,2),
    tax_confidence  NUMERIC(5,2),
    currency        VARCHAR(5) DEFAULT 'INR',
    raw_ocr_text    TEXT,
    needs_review    BOOLEAN DEFAULT FALSE,
    low_conf_fields TEXT[],
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE review_queue (
    review_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doc_id          UUID REFERENCES documents(doc_id),
    extraction_id   UUID REFERENCES extracted_documents(extraction_id),
    review_type     VARCHAR(30), -- 'ocr_low_confidence' | 'handwriting' | 'ai_query'
    assigned_to     UUID REFERENCES users(user_id),
    status          VARCHAR(20) DEFAULT 'pending', -- pending|in_review|resolved
    resolution_json JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    resolved_at     TIMESTAMPTZ
);
```

---

## 6. NATS Event Topics (Phase 04)

```
document.uploaded        → consumed by: B-01 Document Classifier
document.classified      → consumed by: B-02 OCR Extraction
document.handwritten     → consumed by: B-03 Handwriting Recognition
document.extracted       → consumed by: Phase 5 agents (B-04, B-05)
document.review.needed   → consumed by: Notification Dispatcher (alert accountant)
document.reviewed        → consumed by: Phase 5 pipeline resume
```

---

## 7. HITL Review Queue UI

**Screens:**
1. **Review Queue List** — shows all documents requiring review, sorted by priority (oldest first)
2. **Document Review Detail** — side-by-side: original document image + extracted fields
   - Extracted fields displayed with confidence badge (green ≥ 90%, amber 70–89%, red < 70%)
   - Accountant can click any field to correct value
   - "Confirm" button marks extraction as reviewed; NATS `document.reviewed` emitted
   - "Reject Document" option for illegible/corrupt files

**Notification:** When a document enters review queue, assigned accountant receives in-app notification + email

---

## 8. Document Storage Architecture

```
/data/documents/                          ← On-premises file server
    {year}/{month}/{doc_id}.enc           ← AES-256 encrypted
    
Encryption:
    - Per-file AES-256-GCM key
    - Key stored in HashiCorp Vault with policy: only sa-ocr-agent and sa-document-service can access
    - Key ID stored in documents.encryption_key_id

Access:
    - Documents served via `GET /v1/documents/{doc_id}/file` with JWT + OPA check
    - Streaming download; never cached to disk on app server
    - Access logged to audit_log
```

---

## 9. API Endpoints (Phase 04)

| Method | Path | Description |
|---|---|---|
| `POST` | `/v1/documents/upload` | Upload single or multiple documents |
| `GET` | `/v1/documents` | List documents with filter/search |
| `GET` | `/v1/documents/{id}` | Get document metadata + extraction |
| `GET` | `/v1/documents/{id}/file` | Stream document file (encrypted) |
| `PATCH` | `/v1/documents/{id}/extraction` | Submit human correction of extraction |
| `GET` | `/v1/review-queue` | List pending review items |
| `POST` | `/v1/review-queue/{id}/resolve` | Resolve a review item with corrections |

---

## 10. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| B-01 classification | 100 documents (all 6 classes, EN + AR) | ≥ 95% accuracy |
| B-02 OCR printed | 50 printed invoices (EN + AR) | ≥ 95% field-level accuracy |
| B-03 handwriting | 30 handwritten invoices (EN + AR) | ≥ 90% field-level accuracy |
| HITL queue routing | 20 low-confidence extractions | 100% correctly queued |
| Encrypted storage | Upload + retrieve 10 files | Decryption produces bit-exact original |
| Batch upload | 100 simultaneous file uploads | No failures, queue processes in < 5 min |
| Arabic OCR | 20 Arabic invoices (printed) | ≥ 90% field accuracy |

---

## 11. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| PaddleOCR handwriting accuracy below 90% | High | Early test with 30-doc Arabic + English handwritten sample; tune pre-processing; HITL as safety net |
| Arabic OCR model not available for PaddleOCR production use | Medium | Evaluate ArabicOCR, EasyOCR as fallback; consider commercial API (Google Vision, AWS Textract) for Arabic specifically |
| Document storage disk space at scale | Medium | Monitor growth; archive docs older than 3 years to cold storage |
| High HITL queue volume slowing accountants | Medium | Target automatic resolution ≥ 80% of docs; HITL only for genuinely ambiguous cases |

---

## 12. Phase Exit Criteria

- [ ] Document upload (web + mobile) functional; bulk upload tested with 100 files
- [ ] B-01 classification at ≥ 95% accuracy on test set
- [ ] B-02 OCR extraction at ≥ 95% printed, ≥ 90% Arabic
- [ ] B-03 handwriting recognition at ≥ 90% on test set
- [ ] HITL review queue UI operational; accountants can review and correct
- [ ] Document library searchable and filterable
- [ ] Encrypted storage verified; access logs in `audit_log`
- [ ] NATS document pipeline events all flowing correctly
- [ ] Phase gate reviewed and signed off

---

*Phase 04 | Version 1.0 | February 19, 2026*
