/**
 * MediSync Document Processing Types
 *
 * TypeScript types for document processing API requests and responses.
 *
 * @module types/documents
 */

// ============================================================================
// Document Status Types
// ============================================================================

export type DocumentStatus =
  | 'uploading'
  | 'uploaded'
  | 'classifying'
  | 'extracting'
  | 'ready_for_review'
  | 'under_review'
  | 'reviewed'
  | 'approved'
  | 'rejected'
  | 'failed'

export type DocumentType =
  | 'invoice'
  | 'receipt'
  | 'bank_statement'
  | 'expense_report'
  | 'credit_note'
  | 'debit_note'
  | 'other'

export type FileFormat = 'pdf' | 'jpeg' | 'png' | 'tiff' | 'xlsx' | 'csv'

export type FieldType =
  | 'string'
  | 'number'
  | 'currency'
  | 'date'
  | 'percentage'
  | 'identifier'
  | 'tax_id'

export type VerificationStatus =
  | 'pending'
  | 'auto_accepted'
  | 'needs_review'
  | 'high_priority'
  | 'manually_verified'
  | 'manually_corrected'
  | 'rejected'

// ============================================================================
// Document Types
// ============================================================================

/**
 * Bounding box for field location
 */
export interface BoundingBox {
  x: number
  y: number
  width: number
  height: number
  page?: number
}

/**
 * Document response
 */
export interface Document {
  id: string
  tenantId: string
  uploadedBy: string
  status: DocumentStatus
  documentType?: DocumentType
  originalFilename: string
  fileSizeBytes: number
  fileFormat: FileFormat
  pageCount: number
  detectedLanguage: string
  processingStartedAt?: string
  processingCompletedAt?: string
  classificationConfidence: number
  overallConfidence: number
  rejectionReason?: string
  isLocked: boolean
  lockedBy?: string
  lockedAt?: string
  uploadUrl?: string
  createdAt: string
  updatedAt: string
}

/**
 * Document list response
 */
export interface DocumentListResponse {
  documents: Document[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

/**
 * Document stats response
 */
export interface DocumentStatsResponse {
  totalDocuments: number
  pendingReview: number
  underReview: number
  approved: number
  rejected: number
  processing: number
  highPriorityQueue: number
  averageConfidence: number
}

// ============================================================================
// Extracted Field Types
// ============================================================================

/**
 * Extracted field response
 */
export interface ExtractedField {
  id: string
  documentId: string
  pageNumber: number
  fieldName: string
  fieldType: FieldType
  extractedValue: string
  confidenceScore: number
  boundingBox?: BoundingBox
  isHandwritten: boolean
  verificationStatus: VerificationStatus
  verifiedBy?: string
  verifiedAt?: string
  originalValue?: string
  wasEdited: boolean
  createdAt: string
  updatedAt: string
}

/**
 * Field list response
 */
export interface FieldListResponse {
  fields: ExtractedField[]
  totalFields: number
  fieldsNeedingReview: number
  highPriorityCount: number
  autoAcceptedCount: number
}

// ============================================================================
// Line Item Types
// ============================================================================

/**
 * Line item response
 */
export interface LineItem {
  id: string
  documentId: string
  extractedFieldId: string
  lineNumber: number
  description: string
  quantity: number
  unitPrice: number
  amount: number
  taxRate: number
  transactionDate?: string
  reference: string
  debitAmount: number
  creditAmount: number
  balance: number
  createdAt: string
}

// ============================================================================
// Audit Log Types
// ============================================================================

export type AuditAction =
  | 'uploaded'
  | 'classified'
  | 'extracted'
  | 'review_started'
  | 'field_edited'
  | 'field_verified'
  | 'approved'
  | 'rejected'
  | 'reprocessed'

export type ActorType = 'user' | 'system'

/**
 * Audit log response
 */
export interface AuditLogEntry {
  id: string
  documentId: string
  action: AuditAction
  actorId: string
  actorType: ActorType
  fieldName?: string
  oldValue?: unknown
  newValue?: unknown
  notes?: string
  createdAt: string
}

// ============================================================================
// Request Types
// ============================================================================

/**
 * Document upload request
 */
export interface UploadDocumentRequest {
  uploadId?: string
}

/**
 * Document update request
 */
export interface UpdateDocumentRequest {
  documentType?: DocumentType
}

/**
 * Field update request
 */
export interface UpdateFieldRequest {
  value: string
  isVerified?: boolean
}

/**
 * Document approval request
 */
export interface ApproveDocumentRequest {
  notes?: string
}

/**
 * Document rejection request
 */
export interface RejectDocumentRequest {
  reason: string
}

/**
 * Document list filter
 */
export interface DocumentListFilter {
  page?: number
  pageSize?: number
  status?: DocumentStatus[]
  type?: DocumentType[]
  search?: string
  dateFrom?: string
  dateTo?: string
}

// ============================================================================
// Bulk Upload Types
// ============================================================================

/**
 * Bulk upload response
 */
export interface BulkUploadResponse {
  uploadId: string
  totalFiles: number
  uploadedFiles: string[]
  failedFiles: BulkUploadFailure[]
}

/**
 * Bulk upload failure
 */
export interface BulkUploadFailure {
  filename: string
  error: string
}
