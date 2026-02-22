/**
 * MediSync Document API Client
 *
 * HTTP client for document processing API endpoints.
 *
 * @module services/documents
 */

import { api, APIError } from './api'
import type {
  Document,
  DocumentListResponse,
  DocumentStatsResponse,
  ExtractedField,
  FieldListResponse,
  LineItem,
  AuditLogEntry,
  DocumentListFilter,
  UpdateDocumentRequest,
  UpdateFieldRequest,
  ApproveDocumentRequest,
  RejectDocumentRequest,
  BulkUploadResponse,
} from '../types/documents'

const DOCUMENTS_PATH = '/documents'

/**
 * Document API client
 */
export const documentApi = {
  // ===========================================================================
  // Document Operations
  // ===========================================================================

  /**
   * List documents with optional filtering
   */
  list: async (filter?: DocumentListFilter): Promise<DocumentListResponse> => {
    const params: Record<string, string | number | boolean> = {}

    if (filter?.page) params.page = filter.page
    if (filter?.pageSize) params.pageSize = filter.pageSize
    if (filter?.search) params.search = filter.search
    if (filter?.dateFrom) params.dateFrom = filter.dateFrom
    if (filter?.dateTo) params.dateTo = filter.dateTo
    if (filter?.status?.length) params.status = filter.status.join(',')
    if (filter?.type?.length) params.type = filter.type.join(',')

    return api.get<DocumentListResponse>(DOCUMENTS_PATH, { params })
  },

  /**
   * Upload a document
   */
  upload: async (file: File, uploadId?: string): Promise<Document> => {
    const formData = new FormData()
    formData.append('file', file)
    if (uploadId) formData.append('uploadId', uploadId)

    const response = await fetch(`${import.meta.env.VITE_API_URL || '/api/v1'}${DOCUMENTS_PATH}`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('medisync-token')}`,
      },
      body: formData,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: response.statusText }))
      throw new APIError(response.status, response.statusText, error.message, error)
    }

    return response.json()
  },

  /**
   * Get a single document
   */
  get: async (id: string): Promise<Document> => {
    return api.get<Document>(`${DOCUMENTS_PATH}/${id}`)
  },

  /**
   * Update a document
   */
  update: async (id: string, data: UpdateDocumentRequest): Promise<Document> => {
    return api.patch<Document>(`${DOCUMENTS_PATH}/${id}`, data)
  },

  /**
   * Delete a document
   */
  delete: async (id: string): Promise<void> => {
    return api.delete(`${DOCUMENTS_PATH}/${id}`)
  },

  /**
   * Get document statistics
   */
  getStats: async (): Promise<DocumentStatsResponse> => {
    return api.get<DocumentStatsResponse>(`${DOCUMENTS_PATH}/stats`)
  },

  // ===========================================================================
  // Document Lock Operations
  // ===========================================================================

  /**
   * Lock a document for review
   */
  lock: async (id: string): Promise<Document> => {
    return api.post<Document>(`${DOCUMENTS_PATH}/${id}/lock`)
  },

  /**
   * Unlock a document
   */
  unlock: async (id: string): Promise<void> => {
    return api.delete(`${DOCUMENTS_PATH}/${id}/lock`)
  },

  // ===========================================================================
  // Field Operations
  // ===========================================================================

  /**
   * Get all fields for a document
   */
  getFields: async (documentId: string): Promise<FieldListResponse> => {
    return api.get<FieldListResponse>(`${DOCUMENTS_PATH}/${documentId}/fields`)
  },

  /**
   * Update a field
   */
  updateField: async (
    documentId: string,
    fieldId: string,
    data: UpdateFieldRequest
  ): Promise<ExtractedField> => {
    return api.patch<ExtractedField>(`${DOCUMENTS_PATH}/${documentId}/fields/${fieldId}`, data)
  },

  /**
   * Verify a field without changing value
   */
  verifyField: async (documentId: string, fieldId: string): Promise<ExtractedField> => {
    return api.post<ExtractedField>(`${DOCUMENTS_PATH}/${documentId}/fields/${fieldId}/verify`)
  },

  // ===========================================================================
  // Approval Operations
  // ===========================================================================

  /**
   * Approve a document
   */
  approve: async (id: string, data?: ApproveDocumentRequest): Promise<Document> => {
    return api.post<Document>(`${DOCUMENTS_PATH}/${id}/approve`, data)
  },

  /**
   * Reject a document
   */
  reject: async (id: string, data: RejectDocumentRequest): Promise<Document> => {
    return api.post<Document>(`${DOCUMENTS_PATH}/${id}/reject`, data)
  },

  /**
   * Reprocess a rejected/failed document
   */
  reprocess: async (id: string): Promise<Document> => {
    return api.post<Document>(`${DOCUMENTS_PATH}/${id}/reprocess`)
  },

  // ===========================================================================
  // Line Items
  // ===========================================================================

  /**
   * Get line items for a document
   */
  getLineItems: async (documentId: string): Promise<LineItem[]> => {
    return api.get<LineItem[]>(`${DOCUMENTS_PATH}/${documentId}/line-items`)
  },

  // ===========================================================================
  // Audit Log
  // ===========================================================================

  /**
   * Get audit log for a document
   */
  getAuditLog: async (documentId: string): Promise<AuditLogEntry[]> => {
    return api.get<AuditLogEntry[]>(`${DOCUMENTS_PATH}/${documentId}/audit-log`)
  },

  // ===========================================================================
  // Bulk Upload
  // ===========================================================================

  /**
   * Bulk upload documents
   */
  bulkUpload: async (files: File[]): Promise<BulkUploadResponse> => {
    const formData = new FormData()
    files.forEach((file) => {
      formData.append('files', file)
    })

    const response = await fetch(`${import.meta.env.VITE_API_URL || '/api/v1'}${DOCUMENTS_PATH}/bulk-upload`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('medisync-token')}`,
      },
      body: formData,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: response.statusText }))
      throw new APIError(response.status, response.statusText, error.message, error)
    }

    return response.json()
  },

  // ===========================================================================
  // Helper Methods
  // ===========================================================================

  /**
   * Get download URL for a document
   */
  getDownloadUrl: (document: Document): string => {
    return document.uploadUrl || `${import.meta.env.VITE_API_URL || '/api/v1'}${DOCUMENTS_PATH}/${document.id}/download`
  },

  /**
   * Check if a document can be reviewed
   */
  canReview: (document: Document): boolean => {
    return document.status === 'ready_for_review' || document.status === 'under_review'
  },

  /**
   * Check if a document can be approved
   */
  canApprove: (document: Document): boolean => {
    return document.status === 'reviewed' || document.status === 'under_review'
  },

  /**
   * Check if a document can be rejected
   */
  canReject: (document: Document): boolean => {
    return document.canApprove || document.status === 'ready_for_review'
  },

  /**
   * Check if a document can be reprocessed
   */
  canReprocess: (document: Document): boolean => {
    return document.status === 'rejected' || document.status === 'failed'
  },

  /**
   * Get status display text
   */
  getStatusText: (status: string, locale: string = 'en'): string => {
    const statusTexts: Record<string, Record<string, string>> = {
      en: {
        uploading: 'Uploading',
        uploaded: 'Uploaded',
        classifying: 'Classifying',
        extracting: 'Extracting',
        ready_for_review: 'Ready for Review',
        under_review: 'Under Review',
        reviewed: 'Reviewed',
        approved: 'Approved',
        rejected: 'Rejected',
        failed: 'Failed',
      },
      ar: {
        uploading: 'جاري الرفع',
        uploaded: 'تم الرفع',
        classifying: 'جاري التصنيف',
        extracting: 'جاري الاستخراج',
        ready_for_review: 'جاهز للمراجعة',
        under_review: 'قيد المراجعة',
        reviewed: 'تمت المراجعة',
        approved: 'تمت الموافقة',
        rejected: 'مرفوض',
        failed: 'فشل',
      },
    }
    return statusTexts[locale]?.[status] || status
  },

  /**
   * Get verification status display text
   */
  getVerificationStatusText: (status: string, locale: string = 'en'): string => {
    const statusTexts: Record<string, Record<string, string>> = {
      en: {
        pending: 'Pending',
        auto_accepted: 'Auto Accepted',
        needs_review: 'Needs Review',
        high_priority: 'High Priority',
        manually_verified: 'Manually Verified',
        manually_corrected: 'Manually Corrected',
        rejected: 'Rejected',
      },
      ar: {
        pending: 'قيد الانتظار',
        auto_accepted: 'مقبول تلقائياً',
        needs_review: 'يحتاج مراجعة',
        high_priority: 'أولوية عالية',
        manually_verified: 'تم التحقق يدوياً',
        manually_corrected: 'تم التصحيح يدوياً',
        rejected: 'مرفوض',
      },
    }
    return statusTexts[locale]?.[status] || status
  },
}

export default documentApi
