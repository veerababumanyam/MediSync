/**
 * ReviewQueue Component
 *
 * Table display of documents pending review with filtering and sorting.
 */

import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { documentApi } from '../../services/documents'
import { ConfidenceBadge } from './ConfidenceIndicator'
import type { Document, DocumentListFilter, DocumentStatus } from '../../types/documents'

interface ReviewQueueProps {
  onSelectDocument?: (document: Document) => void
  refreshTrigger?: number
  className?: string
}

export function ReviewQueue({ onSelectDocument, refreshTrigger, className = '' }: ReviewQueueProps) {
  const { t, i18n } = useTranslation()
  const isRTL = i18n.language === 'ar'

  const [documents, setDocuments] = useState<Document[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)

  const [filter, setFilter] = useState<{
    status?: DocumentStatus[]
    search?: string
  }>({})

  const fetchDocuments = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const listFilter: DocumentListFilter = {
        page,
        pageSize,
        status: filter.status,
        search: filter.search,
      }

      const response = await documentApi.list(listFilter)
      setDocuments(response.documents)
      setTotal(response.total)
    } catch {
      setError(t('documents.queue.loadError', 'Failed to load documents'))
    } finally {
      setIsLoading(false)
    }
  }, [page, pageSize, filter, t])

  useEffect(() => {
    fetchDocuments()
  }, [fetchDocuments, refreshTrigger])

  const getStatusBadgeClasses = (status: DocumentStatus) => {
    switch (status) {
      case 'ready_for_review':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
      case 'under_review':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
      case 'approved':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
      case 'rejected':
        return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400'
    }
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString(i18n.language === 'ar' ? 'ar-SA' : 'en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className={`review-queue ${className}`} dir={isRTL ? 'rtl' : 'ltr'}>
      {/* Stats Summary */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <StatCard
          label={t('documents.queue.pending', 'Pending Review')}
          value={documents.filter((d) => d.status === 'ready_for_review').length}
          color="yellow"
        />
        <StatCard
          label={t('documents.queue.underReview', 'Under Review')}
          value={documents.filter((d) => d.status === 'under_review').length}
          color="blue"
        />
        <StatCard
          label={t('documents.queue.highPriority', 'High Priority')}
          value={documents.filter((d) => d.overallConfidence < 0.7).length}
          color="red"
        />
        <StatCard
          label={t('documents.queue.total', 'Total')}
          value={total}
          color="gray"
        />
      </div>

      {/* Search and Filter */}
      <div className="flex flex-col sm:flex-row gap-4 mb-4">
        <div className="flex-1">
          <input
            type="text"
            placeholder={t('documents.queue.searchPlaceholder', 'Search documents...')}
            value={filter.search || ''}
            onChange={(e) => {
              setFilter((prev) => ({ ...prev, search: e.target.value }))
              setPage(1)
            }}
            className="w-full px-4 py-2 border rounded-lg bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="flex gap-2">
          <select
            value={filter.status?.[0] || ''}
            onChange={(e) => {
              const status = e.target.value as DocumentStatus
              setFilter((prev) => ({
                ...prev,
                status: status ? [status] : undefined,
              }))
              setPage(1)
            }}
            className="px-4 py-2 border rounded-lg bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600"
          >
            <option value="">{t('documents.queue.allStatus', 'All Status')}</option>
            <option value="ready_for_review">{t('documents.status.readyForReview', 'Ready for Review')}</option>
            <option value="under_review">{t('documents.status.underReview', 'Under Review')}</option>
            <option value="approved">{t('documents.status.approved', 'Approved')}</option>
            <option value="rejected">{t('documents.status.rejected', 'Rejected')}</option>
          </select>
        </div>
      </div>

      {/* Table */}
      {isLoading ? (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto" />
          <p className="mt-2 text-gray-500">{t('common.loading', 'Loading...')}</p>
        </div>
      ) : error ? (
        <div className="text-center py-8 text-red-500">{error}</div>
      ) : documents.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          {t('documents.queue.empty', 'No documents found')}
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b dark:border-gray-700">
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.filename', 'Filename')}
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.type', 'Type')}
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.status', 'Status')}
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.confidence', 'Confidence')}
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.date', 'Date')}
                </th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  {t('documents.queue.actions', 'Actions')}
                </th>
              </tr>
            </thead>
            <tbody className="divide-y dark:divide-gray-700">
              {documents.map((doc) => (
                <tr
                  key={doc.id}
                  className="hover:bg-gray-50 dark:hover:bg-gray-800/50 cursor-pointer"
                  onClick={() => onSelectDocument?.(doc)}
                >
                  <td className="px-4 py-4">
                    <div>
                      <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        {doc.originalFilename}
                      </p>
                      <p className="text-xs text-gray-500">{formatFileSize(doc.fileSizeBytes)}</p>
                    </div>
                  </td>
                  <td className="px-4 py-4 text-sm text-gray-500 dark:text-gray-400">
                    {doc.documentType || '-'}
                  </td>
                  <td className="px-4 py-4">
                    <span
                      className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${getStatusBadgeClasses(
                        doc.status
                      )}`}
                    >
                      {documentApi.getStatusText(doc.status, i18n.language)}
                    </span>
                  </td>
                  <td className="px-4 py-4">
                    <ConfidenceBadge confidence={doc.overallConfidence} />
                  </td>
                  <td className="px-4 py-4 text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(doc.createdAt)}
                  </td>
                  <td className="px-4 py-4 text-right">
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        onSelectDocument?.(doc)
                      }}
                      className="text-blue-600 dark:text-blue-400 hover:underline text-sm"
                    >
                      {t('documents.queue.review', 'Review')}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm text-gray-500">
            {t('documents.queue.showing', 'Showing {{from}}-{{to}} of {{total}}', {
              from: (page - 1) * pageSize + 1,
              to: Math.min(page * pageSize, total),
              total,
            })}
          </p>
          <div className="flex gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 text-sm border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              {t('common.previous', 'Previous')}
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="px-3 py-1 text-sm border rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              {t('common.next', 'Next')}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

function StatCard({
  label,
  value,
  color,
}: {
  label: string
  value: number
  color: 'yellow' | 'blue' | 'red' | 'gray'
}) {
  const colorClasses = {
    yellow: 'bg-yellow-50 dark:bg-yellow-900/20 text-yellow-600 dark:text-yellow-400',
    blue: 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400',
    red: 'bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400',
    gray: 'bg-gray-50 dark:bg-gray-800 text-gray-600 dark:text-gray-400',
  }

  return (
    <div className={`p-4 rounded-lg ${colorClasses[color]}`}>
      <p className="text-2xl font-bold">{value}</p>
      <p className="text-sm opacity-75">{label}</p>
    </div>
  )
}

export default ReviewQueue
