# Data Model: Document Processing Pipeline

**Feature**: 001-document-processing
**Date**: 2026-02-21
**Source**: [spec.md](./spec.md) Key Entities section

---

## Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                           TENANT                                │
│  (Existing - multi-tenant isolation)                           │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                          DOCUMENT                               │
├─────────────────────────────────────────────────────────────────┤
│ id: UUID (PK)                                                   │
│ tenant_id: UUID (FK → tenants.id)                              │
│ uploaded_by: UUID (FK → users.id)                              │
│ status: document_status_enum                                   │
│ document_type: document_type_enum                              │
│ original_filename: VARCHAR(255)                                │
│ storage_path: VARCHAR(500)                                     │
│ file_size_bytes: BIGINT                                        │
│ file_format: file_format_enum                                  │
│ page_count: INTEGER                                            │
│ processing_started_at: TIMESTAMP                               │
│ processing_completed_at: TIMESTAMP                             │
│ classification_confidence: DECIMAL(5,4)                        │
│ overall_confidence: DECIMAL(5,4)                               │
│ rejection_reason: TEXT                                         │
│ created_at: TIMESTAMP                                          │
│ updated_at: TIMESTAMP                                          │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ 1:N
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                       EXTRACTED_FIELD                           │
├─────────────────────────────────────────────────────────────────┤
│ id: UUID (PK)                                                   │
│ document_id: UUID (FK → documents.id)                          │
│ page_number: INTEGER                                           │
│ field_name: VARCHAR(100)                                       │
│ field_type: field_type_enum                                    │
│ extracted_value: TEXT                                          │
│ confidence_score: DECIMAL(5,4)                                 │
│ bounding_box: JSONB {x, y, width, height}                      │
│ is_handwritten: BOOLEAN                                        │
│ verification_status: verification_status_enum                  │
│ verified_by: UUID (FK → users.id, nullable)                    │
│ verified_at: TIMESTAMP                                         │
│ original_value: TEXT (for audit if edited)                     │
│ created_at: TIMESTAMP                                          │
│ updated_at: TIMESTAMP                                          │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            │ N:1 (for line items)
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                         LINE_ITEM                               │
├─────────────────────────────────────────────────────────────────┤
│ id: UUID (PK)                                                   │
│ document_id: UUID (FK → documents.id)                          │
│ extracted_field_id: UUID (FK → extracted_fields.id)            │
│ line_number: INTEGER                                           │
│ description: TEXT                                              │
│ quantity: DECIMAL(12,4)                                        │
│ unit_price: DECIMAL(15,4)                                      │
│ amount: DECIMAL(15,4)                                          │
│ tax_rate: DECIMAL(5,4)                                         │
│ created_at: TIMESTAMP                                          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      DOCUMENT_AUDIT_LOG                         │
├─────────────────────────────────────────────────────────────────┤
│ id: UUID (PK)                                                   │
│ tenant_id: UUID (FK → tenants.id)                              │
│ document_id: UUID (FK → documents.id)                          │
│ action: audit_action_enum                                      │
│ actor_id: UUID (FK → users.id)                                 │
│ actor_type: actor_type_enum (user/system)                      │
│ field_name: VARCHAR(100) (nullable)                            │
│ old_value: JSONB (nullable)                                    │
│ new_value: JSONB (nullable)                                    │
│ notes: TEXT                                                    │
│ created_at: TIMESTAMP                                          │
└─────────────────────────────────────────────────────────────────┘
```

---

## Enumerations

### document_status_enum

| Value | Description | Transitions From |
|-------|-------------|------------------|
| `uploading` | File upload in progress | Initial |
| `uploaded` | Upload complete, awaiting processing | `uploading` |
| `classifying` | Document type detection in progress | `uploaded` |
| `extracting` | OCR extraction in progress | `classifying` |
| `ready_for_review` | Extraction complete, awaiting human review | `extracting` |
| `under_review` | User has opened document for review | `ready_for_review` |
| `reviewed` | All fields verified by user | `under_review` |
| `approved` | Document approved for ledger mapping | `reviewed` |
| `rejected` | Document rejected by user | `ready_for_review`, `under_review` |
| `failed` | Processing failed (corrupted, unreadable) | Any processing state |

### document_type_enum

| Value | Description | Expected Fields |
|-------|-------------|-----------------|
| `invoice` | Supplier invoice | supplier_name, invoice_number, date, line_items, total |
| `receipt` | Payment receipt | supplier_name, receipt_number, date, amount, payment_method |
| `bank_statement` | Bank account statement | bank_name, account_number, period, transactions |
| `expense_report` | Employee expense claim | employee_name, period, expenses, total |
| `credit_note` | Credit note from supplier | supplier_name, credit_note_number, date, amount |
| `debit_note` | Debit note to supplier | supplier_name, debit_note_number, date, amount |
| `other` | Unclassified document | Varies |

### file_format_enum

| Value | MIME Type |
|-------|-----------|
| `pdf` | application/pdf |
| `jpeg` | image/jpeg |
| `png` | image/png |
| `tiff` | image/tiff |
| `xlsx` | application/vnd.openxmlformats-officedocument.spreadsheetml.sheet |
| `csv` | text/csv |

### field_type_enum

| Value | Description | Validation |
|-------|-------------|------------|
| `string` | Free text | None |
| `number` | Numeric value | Is numeric |
| `currency` | Monetary amount | Positive number, max 2 decimals |
| `date` | Date value | Valid date, not in future |
| `percentage` | Percentage | 0-100 range |
| `identifier` | Reference number | Alphanumeric pattern |
| `tax_id` | GST/VAT number | Country-specific format |

### verification_status_enum

| Value | Description | Confidence Threshold |
|-------|-------------|---------------------|
| `pending` | Not yet reviewed | N/A |
| `auto_accepted` | Accepted automatically | ≥95% (printed), N/A (handwritten) |
| `needs_review` | Flagged for human review | 70-94% |
| `high_priority` | Requires immediate attention | <70% |
| `manually_verified` | User confirmed value | N/A |
| `manually_corrected` | User corrected value | N/A |
| `rejected` | User rejected field | N/A |

### audit_action_enum

| Value | Description |
|-------|-------------|
| `uploaded` | Document uploaded |
| `classified` | Document type detected |
| `extracted` | OCR extraction completed |
| `review_started` | User opened for review |
| `field_edited` | Field value modified |
| `field_verified` | Field marked as correct |
| `approved` | Document approved |
| `rejected` | Document rejected |
| `reprocessed` | Document reprocessed |

---

## Field Definitions by Document Type

### Invoice Fields

| Field Name | Type | Required | Validation |
|------------|------|----------|------------|
| supplier_name | string | Yes | Min 2 chars |
| supplier_tax_id | tax_id | No | GST format |
| invoice_number | identifier | Yes | Alphanumeric |
| invoice_date | date | Yes | Not in future |
| due_date | date | No | After invoice_date |
| subtotal | currency | Yes | Positive |
| tax_amount | currency | No | ≥0 |
| tax_rate | percentage | No | 0-100 |
| total | currency | Yes | = subtotal + tax_amount |
| currency | string | No | ISO 4217 code |

### Bank Statement Fields

| Field Name | Type | Required | Validation |
|------------|------|----------|------------|
| bank_name | string | Yes | Min 2 chars |
| account_number | identifier | Yes | Numeric |
| account_name | string | No | Min 2 chars |
| statement_date | date | Yes | Valid date |
| opening_balance | currency | Yes | Any value |
| closing_balance | currency | Yes | Any value |

### Bank Transaction Fields (Line Items)

| Field Name | Type | Required | Validation |
|------------|------|----------|------------|
| transaction_date | date | Yes | Valid date |
| description | string | Yes | Min 1 char |
| reference | identifier | No | Alphanumeric |
| debit | currency | No | ≥0 |
| credit | currency | No | ≥0 |
| balance | currency | Yes | Any value |

---

## Indexes

```sql
-- Primary lookups
CREATE INDEX idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX idx_documents_status ON documents(status);
CREATE INDEX idx_documents_uploaded_by ON documents(uploaded_by);
CREATE INDEX idx_documents_created_at ON documents(created_at DESC);

-- Review queue queries
CREATE INDEX idx_documents_tenant_status ON documents(tenant_id, status)
    WHERE status IN ('ready_for_review', 'under_review');
CREATE INDEX idx_documents_priority ON documents(tenant_id, overall_confidence)
    WHERE status = 'ready_for_review';

-- Extracted fields
CREATE INDEX idx_extracted_fields_document_id ON extracted_fields(document_id);
CREATE INDEX idx_extracted_fields_verification ON extracted_fields(document_id, verification_status);

-- Line items
CREATE INDEX idx_line_items_document_id ON line_items(document_id);

-- Audit log
CREATE INDEX idx_audit_log_document_id ON document_audit_log(document_id);
CREATE INDEX idx_audit_log_tenant_created ON document_audit_log(tenant_id, created_at DESC);
```

---

## Row-Level Security

```sql
-- Enable RLS on all document tables
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE extracted_fields ENABLE ROW LEVEL SECURITY;
ALTER TABLE line_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE document_audit_log ENABLE ROW LEVEL SECURITY;

-- Users can only see documents from their tenant
CREATE POLICY documents_tenant_isolation ON documents
    USING (tenant_id = current_setting('app.current_tenant')::uuid);

CREATE POLICY extracted_fields_tenant_isolation ON extracted_fields
    USING (document_id IN (
        SELECT id FROM documents WHERE tenant_id = current_setting('app.current_tenant')::uuid
    ));

-- Similar policies for line_items and document_audit_log
```

---

## Constraints

### Document Status Transitions

```sql
-- Valid status transitions enforced via trigger or application logic
-- See verification_status_enum table for transition rules
```

### Cross-Field Validation

| Rule | Fields | Condition |
|------|--------|-----------|
| Total must equal sum | subtotal, tax_amount, total | total = subtotal + tax_amount |
| Due date after invoice | invoice_date, due_date | due_date >= invoice_date |
| Debit/credit exclusive | debit, credit (transactions) | (debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0) |

---

## Storage Estimates

| Entity | Est. Records/Month | Est. Size/Record | Monthly Storage |
|--------|-------------------|------------------|-----------------|
| documents | 5,000 | 1 KB | 5 MB |
| extracted_fields | 100,000 | 500 bytes | 50 MB |
| line_items | 50,000 | 300 bytes | 15 MB |
| document_audit_log | 150,000 | 400 bytes | 60 MB |

**Total Database Growth**: ~130 MB/month per tenant (metadata only)

**Object Storage**: 25 MB avg × 5,000 docs = 125 GB/month per tenant
