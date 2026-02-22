# Tasks: Document Processing Pipeline

**Input**: Design documents from `/specs/001-document-processing/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Included per constitution requirement for TDD. Write tests first, ensure they FAIL before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `internal/` for Go backend, `frontend/src/` for React frontend
- **Migrations**: `migrations/` for database schema changes

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create migration file for documents schema in migrations/001_documents.sql
- [ ] T002 [P] Create Document model struct in internal/warehouse/models/document.go
- [ ] T003 [P] Create ExtractedField model struct in internal/warehouse/models/extracted_field.go
- [ ] T004 [P] Create LineItem model struct in internal/warehouse/models/line_item.go
- [ ] T005 [P] Create DocumentAuditLog model struct in internal/warehouse/models/audit_log.go
- [ ] T006 [P] Create enums file for document types/status in internal/warehouse/models/enums.go
- [ ] T007 Create S3-compatible storage client in internal/storage/objects.go
- [ ] T008 [P] Create i18n translation files in frontend/public/locales/en/documents.json
- [ ] T009 [P] Create i18n translation files in frontend/public/locales/ar/documents.json

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T010 Create document repository interface in internal/warehouse/documents.go
- [ ] T011 Create audit log repository in internal/warehouse/audit_log.go
- [ ] T012 Create upload validation middleware in internal/api/middleware/upload.go
- [ ] T013 [P] Create API error types and helpers in internal/api/errors.go
- [ ] T014 Create base router with document routes in internal/api/routes.go
- [ ] T015 [P] Create frontend API client base in frontend/src/services/api.ts
- [ ] T016 [P] Create frontend documents API client in frontend/src/services/documents.ts
- [ ] T017 Run database migrations to create tables

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Upload and Extract Invoice Data (Priority: P1)

**Goal**: Users can upload PDF/image invoices and have OCR automatically extract financial fields with confidence scores

**Independent Test**: Upload a sample invoice PDF, verify it appears in processing queue, then verify extracted fields with confidence scores appear in the system

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T018 [P] [US1] Create unit tests for document upload handler in internal/api/handlers/documents_test.go
- [ ] T019 [P] [US1] Create unit tests for PaddleOCR client in internal/etl/ocr/paddle_test.go
- [ ] T020 [P] [US1] Create unit tests for confidence scoring in internal/etl/ocr/confidence_test.go
- [ ] T021 [P] [US1] Create integration test for upload-extract flow in tests/integration/document_upload_test.go
- [ ] T022 [P] [US1] Create frontend tests for DocumentUploader in frontend/tests/DocumentUpload.test.tsx

### Backend Implementation for User Story 1

- [ ] T023 [P] [US1] Implement PaddleOCR client wrapper in internal/etl/ocr/paddle.go
- [ ] T024 [P] [US1] Implement document classifier in internal/etl/ocr/classifier.go
- [ ] T025 [US1] Implement confidence scoring logic in internal/etl/ocr/confidence.go (depends on T023, T024)
- [ ] T026 [US1] Implement OCR extraction agent B-02 in internal/agents/module_b/ocr_extraction.go
- [ ] T027 [US1] Implement document upload handler in internal/api/handlers/documents.go
- [ ] T028 [US1] Implement document list handler in internal/api/handlers/documents.go
- [ ] T029 [US1] Implement document get handler in internal/api/handlers/documents.go
- [ ] T030 [US1] Implement document download handler in internal/api/handlers/documents.go
- [ ] T031 [US1] Create NATS JetStream consumer for OCR processing in cmd/ocr-worker/main.go
- [ ] T032 [US1] Implement document status transitions in internal/warehouse/documents.go
- [ ] T033 [US1] Add audit logging for document upload in internal/warehouse/audit_log.go

### Frontend Implementation for User Story 1

- [ ] T034 [P] [US1] Create DocumentUploader component in frontend/src/components/DocumentUploader/DocumentUploader.tsx
- [ ] T035 [P] [US1] Create FileDropZone subcomponent in frontend/src/components/DocumentUploader/FileDropZone.tsx
- [ ] T036 [P] [US1] Create UploadProgress subcomponent in frontend/src/components/DocumentUploader/UploadProgress.tsx
- [ ] T037 [US1] Create Documents upload page in frontend/src/pages/Documents/UploadPage.tsx (depends on T034-T036)
- [ ] T038 [US1] Create Documents list page in frontend/src/pages/Documents/Documents.tsx
- [ ] T039 [US1] Add routing for /documents/* in frontend/src/App.tsx

**Checkpoint**: At this point, User Story 1 should be fully functional - users can upload documents and see extracted fields

---

## Phase 4: User Story 2 - Review and Correct Extracted Data (Priority: P1)

**Goal**: Finance team can review extracted fields, edit values, verify correctness, and approve/reject documents

**Independent Test**: Open a processed document in review queue, edit a field value, verify the change is saved, then approve the document

### Tests for User Story 2

- [ ] T040 [P] [US2] Create unit tests for review queue handler in internal/api/handlers/review_test.go
- [ ] T041 [P] [US2] Create unit tests for field update handler in internal/api/handlers/extraction_test.go
- [ ] T042 [P] [US2] Create integration test for review-approve flow in tests/integration/document_review_test.go
- [ ] T043 [P] [US2] Create frontend tests for ReviewQueue in frontend/tests/ReviewQueue.test.tsx
- [ ] T044 [P] [US2] Create frontend tests for FieldEditor in frontend/tests/FieldEditor.test.tsx

### Backend Implementation for User Story 2

- [ ] T045 [US2] Implement review queue list handler in internal/api/handlers/review.go
- [ ] T046 [US2] Implement review queue stats handler in internal/api/handlers/review.go
- [ ] T047 [US2] Implement start review handler in internal/api/handlers/review.go
- [ ] T048 [US2] Implement field update/verify handler in internal/api/handlers/extraction.go
- [ ] T049 [US2] Implement document approve handler in internal/api/handlers/review.go
- [ ] T050 [US2] Implement document reject handler in internal/api/handlers/review.go
- [ ] T051 [US2] Add field verification status logic in internal/warehouse/documents.go
- [ ] T052 [US2] Add audit logging for review actions in internal/warehouse/audit_log.go
- [ ] T053 [US2] Implement page image rendering handler in internal/api/handlers/documents.go

### Frontend Implementation for User Story 2

- [ ] T054 [P] [US2] Create ConfidenceIndicator component in frontend/src/components/ConfidenceIndicator/ConfidenceIndicator.tsx
- [ ] T055 [P] [US2] Create FieldEditor component in frontend/src/components/FieldEditor/FieldEditor.tsx
- [ ] T056 [P] [US2] Create FieldRow subcomponent in frontend/src/components/FieldEditor/FieldRow.tsx
- [ ] T057 [P] [US2] Create ReviewQueue table component in frontend/src/components/ReviewQueue/ReviewQueue.tsx
- [ ] T058 [P] [US2] Create QueueStats subcomponent in frontend/src/components/ReviewQueue/QueueStats.tsx
- [ ] T059 [US2] Create Review queue page in frontend/src/pages/Documents/QueuePage.tsx (depends on T057, T058)
- [ ] T060 [US2] Create Document review page with field editing in frontend/src/pages/Documents/ReviewPage.tsx (depends on T054-T056)
- [ ] T061 [US2] Create DocumentPreview component for original file display in frontend/src/components/DocumentPreview/DocumentPreview.tsx

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - complete upload → review → approve workflow

---

## Phase 5: User Story 3 - Process Multi-Page Documents (Priority: P2)

**Goal**: Users can upload multi-page PDFs and bank statements with all pages/transactions extracted correctly

**Independent Test**: Upload a 5-page invoice PDF, verify all line items from all pages are extracted and consolidated

### Tests for User Story 3

- [ ] T062 [P] [US3] Create unit tests for multi-page processing in internal/etl/ocr/paddle_test.go
- [ ] T063 [P] [US3] Create integration test for multi-page document in tests/integration/multipage_document_test.go

### Backend Implementation for User Story 3

- [ ] T064 [US3] Add PDF page splitting logic in internal/etl/ocr/paddle.go
- [ ] T065 [US3] Implement page consolidation logic in internal/etl/ocr/paddle.go
- [ ] T066 [US3] Add bank statement transaction extraction in internal/etl/ocr/classifier.go
- [ ] T067 [US3] Update line item repository for bulk insert in internal/warehouse/documents.go
- [ ] T068 [US3] Add page count tracking in document model in internal/warehouse/models/document.go

### Frontend Implementation for User Story 3

- [ ] T069 [US3] Add multi-page preview navigation in frontend/src/components/DocumentPreview/PageNavigator.tsx
- [ ] T070 [US3] Update ReviewPage to show page references in frontend/src/pages/Documents/ReviewPage.tsx

**Checkpoint**: Multi-page documents now process correctly with all pages consolidated

---

## Phase 6: User Story 4 - Handle Handwritten Documents (Priority: P2)

**Goal**: System detects handwritten portions, applies lower confidence caps, and flags for manual verification

**Independent Test**: Upload a receipt with handwritten notes, verify fields are flagged with confidence <85% and highlighted for review

### Tests for User Story 4

- [ ] T071 [P] [US4] Create unit tests for handwriting detection in internal/etl/ocr/handwriting_test.go
- [ ] T072 [P] [US4] Create integration test for handwritten document in tests/integration/handwritten_document_test.go

### Backend Implementation for User Story 4

- [ ] T073 [US4] Implement handwriting detection in internal/etl/ocr/handwriting.go
- [ ] T074 [US4] Add handwriting confidence capping (85% max) in internal/etl/ocr/confidence.go
- [ ] T075 [US4] Create handwriting agent B-03 in internal/agents/module_b/handwriting.go
- [ ] T076 [US4] Add is_handwritten field handling in internal/warehouse/models/extracted_field.go

### Frontend Implementation for User Story 4

- [ ] T077 [US4] Add handwriting indicator badge in frontend/src/components/FieldEditor/FieldRow.tsx
- [ ] T078 [US4] Add high-priority styling for <70% confidence in frontend/src/components/ConfidenceIndicator/ConfidenceIndicator.tsx

**Checkpoint**: Handwritten documents are now properly detected and flagged

---

## Phase 7: User Story 5 - Arabic Document Support (Priority: P2)

**Goal**: Arabic invoices are extracted correctly with RTL display in the review interface

**Independent Test**: Upload an Arabic invoice, verify text is extracted correctly and displayed in RTL layout

### Tests for User Story 5

- [ ] T079 [P] [US5] Create unit tests for Arabic OCR in internal/etl/ocr/paddle_test.go
- [ ] T080 [P] [US5] Create integration test for Arabic document in tests/integration/arabic_document_test.go
- [ ] T081 [P] [US5] Create frontend RTL layout tests in frontend/tests/RTLLayout.test.tsx

### Backend Implementation for User Story 5

- [ ] T082 [US5] Configure PaddleOCR for Arabic language in internal/etl/ocr/paddle.go
- [ ] T083 [US5] Add BiDi text direction detection in internal/etl/ocr/bidi.go
- [ ] T084 [US5] Store detected language in document model in internal/warehouse/models/document.go

### Frontend Implementation for User Story 5

- [ ] T085 [US5] Add RTL layout support to DocumentPreview in frontend/src/components/DocumentPreview/DocumentPreview.tsx
- [ ] T086 [US5] Add RTL layout support to FieldEditor in frontend/src/components/FieldEditor/FieldEditor.tsx
- [ ] T087 [US5] Add Arabic translations for all document strings in frontend/public/locales/ar/documents.json

**Checkpoint**: Arabic documents are fully supported with proper RTL display

---

## Phase 8: User Story 6 - Bulk Document Upload (Priority: P3)

**Goal**: Users can upload up to 50 documents at once for efficient batch processing

**Independent Test**: Select 10 files for upload, verify all are uploaded and appear in the processing queue

### Tests for User Story 6

- [ ] T088 [P] [US6] Create unit tests for bulk upload in internal/api/handlers/documents_test.go
- [ ] T089 [P] [US6] Create integration test for bulk upload in tests/integration/bulk_upload_test.go
- [ ] T090 [P] [US6] Create frontend tests for bulk upload in frontend/tests/BulkUpload.test.tsx

### Backend Implementation for User Story 6

- [ ] T091 [US6] Implement bulk upload endpoint in internal/api/handlers/documents.go
- [ ] T092 [US6] Add batch processing optimization in internal/etl/ocr/paddle.go
- [ ] T093 [US6] Add upload_id tracking for bulk operations in internal/warehouse/models/document.go

### Frontend Implementation for User Story 6

- [ ] T094 [US6] Add multi-file selection to DocumentUploader in frontend/src/components/DocumentUploader/DocumentUploader.tsx
- [ ] T095 [US6] Create BulkUploadProgress component in frontend/src/components/DocumentUploader/BulkUploadProgress.tsx
- [ ] T096 [US6] Update UploadPage to show bulk status in frontend/src/pages/Documents/UploadPage.tsx

**Checkpoint**: Bulk upload now works efficiently for up to 50 documents

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T097 [P] Add API rate limiting for upload endpoints in internal/api/middleware/rate_limit.go
- [ ] T098 [P] Add Prometheus metrics for document processing in internal/api/metrics.go
- [ ] T099 [P] Create test fixture documents in tests/fixtures/documents/
- [ ] T100 Run full integration test suite and verify all acceptance scenarios
- [ ] T101 [P] Add error handling for edge cases (password-protected, corrupted files) in internal/api/handlers/documents.go
- [ ] T102 [P] Add WebSocket notifications for processing status in internal/api/handlers/websocket.go
- [ ] T103 Run quickstart.md validation to verify developer setup works
- [ ] T104 Update CLAUDE.md with document processing agent documentation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-8)**: All depend on Foundational phase completion
  - US1 and US2 are both P1 - implement US1 first (upload required before review)
  - US3, US4, US5 are all P2 - can proceed in parallel or any order
  - US6 (P3) - can proceed after US1 is complete
- **Polish (Phase 9)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Depends on US1 (need documents to review)
- **User Story 3 (P2)**: Can start after US1 (extends single-page capability)
- **User Story 4 (P2)**: Can start after US1 (extends OCR pipeline)
- **User Story 5 (P2)**: Can start after US1 (extends OCR pipeline)
- **User Story 6 (P3)**: Can start after US1 (extends upload capability)

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Backend models before services
- Backend services before handlers
- Frontend components before pages
- Core implementation before integration

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T002-T006, T008-T009)
- All Foundational tasks marked [P] can run in parallel (T013, T015-T016)
- All test tasks for a user story marked [P] can run in parallel
- Backend and frontend tasks can run in parallel once tests are written
- US3, US4, US5 can all be worked on in parallel after US1+US2 complete

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Create unit tests for document upload handler in internal/api/handlers/documents_test.go"
Task: "Create unit tests for PaddleOCR client in internal/etl/ocr/paddle_test.go"
Task: "Create unit tests for confidence scoring in internal/etl/ocr/confidence_test.go"
Task: "Create integration test for upload-extract flow in tests/integration/document_upload_test.go"
Task: "Create frontend tests for DocumentUploader in frontend/tests/DocumentUpload.test.tsx"

# After tests fail, launch parallel backend implementation:
Task: "Implement PaddleOCR client wrapper in internal/etl/ocr/paddle.go"
Task: "Implement document classifier in internal/etl/ocr/classifier.go"

# Launch parallel frontend implementation (independent of backend):
Task: "Create DocumentUploader component in frontend/src/components/DocumentUploader/DocumentUploader.tsx"
Task: "Create FileDropZone subcomponent in frontend/src/components/DocumentUploader/FileDropZone.tsx"
Task: "Create UploadProgress subcomponent in frontend/src/components/DocumentUploader/UploadProgress.tsx"
```

---

## Implementation Strategy

**Scope**: Production-ready release - ALL 6 user stories must be complete.

### Sequential Execution (All Stories)

1. **Phase 1**: Setup (9 tasks) - Database, models, i18n
2. **Phase 2**: Foundational (8 tasks) - Repositories, middleware, API base
3. **Phase 3**: User Story 1 - Upload & Extract (22 tasks)
4. **Phase 4**: User Story 2 - Review & Correct (22 tasks)
5. **Phase 5**: User Story 3 - Multi-Page (9 tasks)
6. **Phase 6**: User Story 4 - Handwriting (8 tasks)
7. **Phase 7**: User Story 5 - Arabic Support (9 tasks)
8. **Phase 8**: User Story 6 - Bulk Upload (9 tasks)
9. **Phase 9**: Polish & Production Hardening (8 tasks)

### Production Readiness Checklist

Before deployment, verify ALL of the following:

- [ ] All 104 tasks complete
- [ ] All 22 test suites passing (unit + integration + frontend)
- [ ] OCR accuracy targets met (95% English, 90% Arabic, 80% handwriting)
- [ ] Processing time targets met (<30s single doc, <2min multi-page, <10min bulk)
- [ ] All edge cases handled (password-protected, corrupted, oversized files)
- [ ] Arabic RTL interface fully functional
- [ ] Audit logging complete for all actions
- [ ] Rate limiting and metrics in place
- [ ] quickstart.md validated by fresh developer setup

### Parallel Team Strategy

With multiple developers, maximize parallelization:

**After Phase 2 (Foundational) Complete:**

| Developer | Primary Assignment | Secondary Assignment |
|-----------|-------------------|----------------------|
| Dev A | US1 Backend | US1 Frontend |
| Dev B | US2 Backend | US2 Frontend |

**After US1 + US2 Complete:**

| Developer | Assignment |
|-----------|------------|
| Dev A | US3 (Multi-Page) |
| Dev B | US4 (Handwriting) |
| Dev C | US5 (Arabic) |

**After US3, US4, US5 Complete:**

| Developer | Assignment |
|-----------|------------|
| Dev A | US6 (Bulk Upload) |
| Dev B | Phase 9 (Polish) |

### Critical Path

```
Phase 1 (Setup) → Phase 2 (Foundational) → US1 → US2 → [US3, US4, US5 parallel] → US6 → Phase 9
                                    ↓
                              BLOCKS ALL STORIES
```

**Minimum Timeline**: US1 and US2 are sequential (US2 requires US1). US3, US4, US5 can run in parallel. US6 depends on US1.

---

## Task Summary

| Phase | Tasks | Parallelizable | Description |
|-------|-------|----------------|-------------|
| Phase 1: Setup | 9 | 6 | Database migrations, models, i18n |
| Phase 2: Foundational | 8 | 2 | Repositories, middleware, API base |
| Phase 3: US1 (P1) | 22 | 10 | Upload & OCR extraction |
| Phase 4: US2 (P1) | 22 | 8 | Review queue & field editing |
| Phase 5: US3 (P2) | 9 | 2 | Multi-page document support |
| Phase 6: US4 (P2) | 8 | 2 | Handwriting detection |
| Phase 7: US5 (P2) | 9 | 3 | Arabic RTL support |
| Phase 8: US6 (P3) | 9 | 3 | Bulk upload (50 files) |
| Phase 9: Polish | 8 | 4 | Rate limiting, metrics, edge cases |
| **Production Total** | **104** | **40** | All features production-ready |

### Production Scope (Required)

**All 104 tasks must be completed for production release:**
- US1 + US2: Core workflow (upload → review → approve)
- US3: Multi-page documents (healthcare invoices, bank statements)
- US4: Handwriting recognition (receipts, expense notes)
- US5: Arabic language support (constitution requirement for i18n)
- US6: Bulk upload (efficiency for high-volume clinics)
- Polish: Security, observability, error handling

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Constitution requires TDD - tests are included for all stories
