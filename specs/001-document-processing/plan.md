# Implementation Plan: Document Processing Pipeline

**Branch**: `001-document-processing` | **Date**: 2026-02-21 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-document-processing/spec.md`

## Summary

Implement the OCR pipeline for document ingestion as Part 1 of the AI Accountant module. Users upload financial documents (invoices, receipts, bank statements) which are automatically classified, processed through OCR with confidence scoring, and queued for human review. The pipeline supports multi-page documents, handwritten content, and Arabic text extraction. All extracted data requires human verification before proceeding to ledger mapping.

## Technical Context

**Language/Version**: Go 1.26 (backend), TypeScript 5.9 (frontend)
**Primary Dependencies**: go-chi/chi (HTTP), PaddleOCR (OCR engine), NATS JetStream (async processing), Google Genkit (AI flows), React 19 (web), i18next (i18n), Apache ECharts (confidence visualization)
**Storage**: PostgreSQL 18.2 + pgvector (document metadata, embeddings), Redis (queue status, session), S3-compatible object storage (original documents)
**Testing**: Go `testing` + `testify` (backend), Vitest + React Testing Library (frontend)
**Target Platform**: Linux server (backend), modern web browsers (frontend)
**Project Type**: Web application (backend + frontend)
**Performance Goals**: Single document processing <30s, multi-page (20 pages) <2min, bulk upload (50 docs) <10min
**Constraints**: 25MB max file size, 20 pages per document, 50 files per batch, encrypted storage at rest
**Scale/Scope**: Multi-tenant healthcare clinics, thousands of documents per day

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Evidence |
|-----------|--------|----------|
| **I. Security First & HITL Gates** | ✅ PASS | Document storage uses encryption at rest; all extracted fields require human verification before approval; confidence thresholds enforce review |
| **II. Read-Only Intelligence Plane** | ✅ PASS | OCR agents read documents for extraction only; no AI writes to Tally; document writes are user-initiated via API |
| **III. i18n by Default** | ✅ PASS | Arabic OCR support specified (FR-013); RTL interface for Arabic documents; i18next for all UI strings |
| **IV. Open Source Only** | ✅ PASS | PaddleOCR (Apache-2.0) is approved; all dependencies from constitution stack |
| **V. Test-Driven Development** | ✅ PASS | Test plan includes OCR accuracy tests, review queue tests, multi-format tests |

**Gate Status**: ✅ All gates passed. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/001-document-processing/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (OpenAPI specs)
│   ├── documents-api.yaml
│   └── ocr-api.yaml
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── agents/
│   └── module_b/              # AI Accountant agents
│       ├── ocr_extraction.go  # Agent B-02: OCR extraction flow
│       ├── handwriting.go     # Agent B-03: Handwriting recognition
│       └── document_classify.go
├── api/
│   ├── handlers/
│   │   ├── documents.go       # Upload, list, get documents
│   │   ├── review.go          # Review queue operations
│   │   └── extraction.go      # Field extraction endpoints
│   └── middleware/
│       └── upload.go          # File validation middleware
├── warehouse/
│   └── documents.go           # Document repository (read/write)
├── etl/
│   └── ocr/                   # OCR processing pipeline
│       ├── paddle.go          # PaddleOCR client
│       ├── classifier.go      # Document type classification
│       └── confidence.go      # Confidence scoring logic
└── storage/
    └── objects.go             # S3-compatible storage client

frontend/
├── src/
│   ├── pages/
│   │   ├── Documents/
│   │   │   ├── UploadPage.tsx
│   │   │   ├── QueuePage.tsx
│   │   │   └── ReviewPage.tsx
│   │   └── Documents.tsx
│   ├── components/
│   │   ├── DocumentUploader/
│   │   ├── ReviewQueue/
│   │   ├── FieldEditor/
│   │   └── ConfidenceIndicator/
│   └── services/
│       └── documents.ts       # API client
├── public/
│   └── locales/
│       ├── en/
│       │   └── documents.json
│       └── ar/
│           └── documents.json
└── tests/
    ├── DocumentUpload.test.tsx
    ├── ReviewQueue.test.tsx
    └── FieldEditor.test.tsx

migrations/
└── 001_documents.sql          # Document tables
```

**Structure Decision**: Follows established MediSync monorepo pattern with `internal/` for backend Go code and `frontend/` for React web app. Document processing agents reside in `module_b/` (AI Accountant) as per agent catalog in CLAUDE.md.

## Complexity Tracking

> No constitution violations. Table not required.
