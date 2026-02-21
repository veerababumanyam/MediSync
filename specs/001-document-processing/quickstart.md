# Quickstart: Document Processing Pipeline

**Feature**: 001-document-processing
**Date**: 2026-02-21

This guide helps developers quickly understand and start working on the document processing feature.

---

## Prerequisites

1. **Go 1.26** installed
2. **Node.js 20+** and **npm** for frontend
3. **Docker** and **docker-compose** for local infrastructure
4. **PaddleOCR** (run via Docker container)
5. **PostgreSQL 18.2** with pgvector extension
6. **Redis** for queue status caching
7. **MinIO** (S3-compatible local storage)

---

## Quick Setup

### 1. Start Infrastructure

```bash
# From repository root
docker-compose up -d postgres redis minio nats

# Start PaddleOCR service
docker run -d --name paddleocr \
  -p 8868:8868 \
  --gpus all \
  paddlepaddle/paddleocr:latest
```

### 2. Run Database Migrations

```bash
go run ./cmd/migrate up
```

### 3. Start Backend API

```bash
# Set required environment variables
export DATABASE_URL="postgres://medisync:password@localhost:5432/medisync?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export NATS_URL="nats://localhost:4222"
export S3_ENDPOINT="http://localhost:9000"
export S3_ACCESS_KEY="minioadmin"
export S3_SECRET_KEY="minioadmin"
export S3_BUCKET="medisync-documents"
export OCR_SERVICE_URL="http://localhost:8868"

# Start API server
go run ./cmd/api
```

### 4. Start OCR Worker

```bash
# In a separate terminal
go run ./cmd/ocr-worker
```

### 5. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

---

## Key Files to Understand

### Backend

| File | Purpose |
|------|---------|
| `internal/agents/module_b/ocr_extraction.go` | Agent B-02: Main OCR extraction flow |
| `internal/etl/ocr/paddle.go` | PaddleOCR client wrapper |
| `internal/etl/ocr/classifier.go` | Document type classification |
| `internal/etl/ocr/confidence.go` | Confidence scoring logic |
| `internal/api/handlers/documents.go` | Document upload/list API handlers |
| `internal/api/handlers/review.go` | Review queue API handlers |
| `internal/warehouse/documents.go` | Document database operations |
| `internal/storage/objects.go` | S3-compatible storage client |

### Frontend

| File | Purpose |
|------|---------|
| `frontend/src/pages/Documents/UploadPage.tsx` | Document upload UI |
| `frontend/src/pages/Documents/QueuePage.tsx` | Review queue list |
| `frontend/src/pages/Documents/ReviewPage.tsx` | Single document review |
| `frontend/src/components/DocumentUploader/` | Upload component with drag-drop |
| `frontend/src/components/ReviewQueue/` | Queue table with filters |
| `frontend/src/components/FieldEditor/` | Field value editor |
| `frontend/src/components/ConfidenceIndicator/` | Confidence score visualization |

---

## API Endpoints

### Document Operations

```bash
# Upload a document
curl -X POST http://localhost:8080/v1/documents \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@invoice.pdf"

# List documents
curl http://localhost:8080/v1/documents?status=ready_for_review \
  -H "Authorization: Bearer $TOKEN"

# Get document details
curl http://localhost:8080/v1/documents/{documentId} \
  -H "Authorization: Bearer $TOKEN"

# Get review queue
curl http://localhost:8080/v1/review-queue \
  -H "Authorization: Bearer $TOKEN"
```

### Review Operations

```bash
# Start review
curl -X POST http://localhost:8080/v1/documents/{documentId}/review \
  -H "Authorization: Bearer $TOKEN"

# Edit a field
curl -X PATCH http://localhost:8080/v1/documents/{documentId}/fields/{fieldId} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action": "edit", "value": "New Value"}'

# Approve document
curl -X POST http://localhost:8080/v1/documents/{documentId}/approve \
  -H "Authorization: Bearer $TOKEN"
```

---

## Processing Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Upload    │────▶│   Store     │────▶│   Publish   │
│   (API)     │     │   (S3)      │     │   (NATS)    │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Update    │◀────│   Extract   │◀────│   OCR       │
│   (DB)      │     │   Fields    │     │   Worker    │
└─────────────┘     └─────────────┘     └─────────────┘
       │
       ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Review    │────▶│   Approve   │────▶│   Ready for │
│   Queue     │     │   (User)    │     │   Ledger    │
└─────────────┘     └─────────────┘     └─────────────┘
```

---

## Testing

### Run Backend Tests

```bash
# All tests
go test ./internal/... -v

# OCR-specific tests
go test ./internal/etl/ocr/... -v

# Integration tests (requires Docker)
go test ./tests/integration/... -v -tags=integration
```

### Run Frontend Tests

```bash
cd frontend
npm test

# With coverage
npm test -- --coverage
```

### Test Documents

Sample test documents are available in `tests/fixtures/documents/`:
- `invoice_en.pdf` - English invoice
- `invoice_ar.pdf` - Arabic invoice
- `bank_statement.csv` - Bank statement
- `handwritten_receipt.jpg` - Handwritten receipt

---

## Common Tasks

### Add a New Document Type

1. Add type to `document_type_enum` in migration
2. Add field definitions to `internal/etl/ocr/classifier.go`
3. Add expected fields to classification logic
4. Update frontend i18n keys in `en/documents.json` and `ar/documents.json`

### Add a New Extracted Field

1. Add field validation to `internal/etl/ocr/confidence.go`
2. Update field type enum if needed
3. Add to document type field definitions
4. Update API schemas in `contracts/documents-api.yaml`

### Modify Confidence Thresholds

Edit `internal/etl/ocr/confidence.go`:

```go
const (
    AutoAcceptThreshold    = 0.95  // Auto-accept at 95%
    NeedsReviewThreshold   = 0.70  // Flag for review below 70%
    HandwritingCap         = 0.85  // Cap handwriting at 85%
)
```

---

## Troubleshooting

### OCR Returns Empty Results

1. Check PaddleOCR container is running: `docker ps | grep paddleocr`
2. Check logs: `docker logs paddleocr`
3. Verify file format is supported
4. Check file isn't password-protected

### Upload Fails with 413

1. Check file size is under 25MB
2. Verify nginx/client body size limits if behind proxy

### Fields Not Appearing in Review Queue

1. Check OCR worker logs: `tail -f logs/ocr-worker.log`
2. Verify NATS connection: `nats server info`
3. Check document status in database

### Arabic Text Not Displaying Correctly

1. Verify frontend i18n locale is set correctly
2. Check `dir="rtl"` is applied to container
3. Verify font supports Arabic glyphs

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |
| `NATS_URL` | NATS connection string | `nats://localhost:4222` |
| `S3_ENDPOINT` | S3-compatible endpoint | Required |
| `S3_ACCESS_KEY` | S3 access key | Required |
| `S3_SECRET_KEY` | S3 secret key | Required |
| `S3_BUCKET` | Document storage bucket | `medisync-documents` |
| `OCR_SERVICE_URL` | PaddleOCR service URL | `http://localhost:8868` |
| `MAX_FILE_SIZE_MB` | Maximum upload size | `25` |
| `MAX_PAGES` | Maximum pages per document | `20` |
| `MAX_BATCH_SIZE` | Maximum bulk upload files | `50` |

---

## Related Documentation

- [spec.md](./spec.md) - Feature specification
- [data-model.md](./data-model.md) - Database schema
- [research.md](./research.md) - Technical decisions
- [contracts/documents-api.yaml](./contracts/documents-api.yaml) - API contract
