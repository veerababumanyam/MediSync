/**
 * ReviewQueue Component
 *
 * Table display of documents pending review with filtering and sorting.
 * Styled with Liquid Glass design system for premium glassmorphic aesthetics.
 */

import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { documentApi } from '../../services/documents'
import { ConfidenceBadge } from './ConfidenceIndicator'
import { LiquidGlassCard } from '../ui/LiquidGlassCard'
import { LiquidGlassBadge } from '../ui/LiquidGlassBadge'
import { LiquidGlassButton } from '../ui/LiquidGlassButton'
import { LiquidGlassInput } from '../ui/LiquidGlassInput'
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
          variant="warning"
        />
        <StatCard
          label={t('documents.queue.underReview', 'Under Review')}
          value={documents.filter((d) => d.status === 'under_review').length}
          variant="blue"
        />
        <StatCard
          label={t('documents.queue.highPriority', 'High Priority')}
          value={documents.filter((d) => d.overallConfidence < 0.7).length}
          variant="error"
        />
        <StatCard
          label={t('documents.queue.total', 'Total')}
          value={total}
          variant="default"
        />
      </div>

      {/* Search and Filter */}
      <div className="flex flex-col sm:flex-row gap-4 mb-4">
        <div className="flex-1">
          <LiquidGlassInput
            placeholder={t('documents.queue.searchPlaceholder', 'Search documents...')}
            value={filter.search || ''}
            onChange={(e) => {
              setFilter((prev) => ({ ...prev, search: e.target.value }))
              setPage(1)
            }}
            className="w-full"
          />
        </div>
        <div className="flex gap-2">
          <LiquidGlassCard
            intensity="subtle"
            className="px-4 py-2"
          >
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
              className="bg-transparent liquid-text-primary focus:outline-none cursor-pointer"
            >
              <option value="">{t('documents.queue.allStatus', 'All Status')}</option>
              <option value="ready_for_review">{t('documents.status.readyForReview', 'Ready for Review')}</option>
              <option value="under_review">{t('documents.status.underReview', 'Under Review')}</option>
              <option value="approved">{t('documents.status.approved', 'Approved')}</option>
              <option value="rejected">{t('documents.status.rejected', 'Rejected')}</option>
            </select>
          </LiquidGlassCard>
        </div>
      </div>

      {/* Table */}
      {isLoading ? (
        <LiquidGlassCard intensity="subtle" className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto" />
          <p className="mt-2 liquid-text-secondary">{t('common.loading', 'Loading...')}</p>
        </LiquidGlassCard>
      ) : error ? (
        <LiquidGlassCard intensity="subtle" className="text-center py-8 text-red-400">{error}</LiquidGlassCard>
      ) : documents.length === 0 ? (
        <LiquidGlassCard intensity="subtle" className="text-center py-8 liquid-text-muted">
          {t('documents.queue.empty', 'No documents found')}
        </LiquidGlassCard>
      ) : (
        <LiquidGlassCard intensity="subtle" className="overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-white/10">
                  <th className="px-4 py-3 text-left text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.filename', 'Filename')}
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.type', 'Type')}
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.status', 'Status')}
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.confidence', 'Confidence')}
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.date', 'Date')}
                  </th>
                  <th className="px-4 py-3 text-right text-xs font-medium liquid-text-muted uppercase tracking-wider">
                    {t('documents.queue.actions', 'Actions')}
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {documents.map((doc) => (
                  <tr
                    key={doc.id}
                    className="hover:bg-white/5 cursor-pointer transition-colors"
                    onClick={() => onSelectDocument?.(doc)}
                  >
                    <td className="px-4 py-4">
                      <div>
                        <p className="text-sm font-medium liquid-text-primary">
                          {doc.originalFilename}
                        </p>
                        <p className="text-xs liquid-text-muted">{formatFileSize(doc.fileSizeBytes)}</p>
                      </div>
                    </td>
                    <td className="px-4 py-4 text-sm liquid-text-secondary">
                      {doc.documentType || '-'}
                    </td>
                    <td className="px-4 py-4">
                      <LiquidGlassBadge variant={getStatusBadgeVariant(doc.status)}>
                        {documentApi.getStatusText(doc.status, i18n.language)}
                      </LiquidGlassBadge>
                    </td>
                    <td className="px-4 py-4">
                      <ConfidenceBadge confidence={doc.overallConfidence} />
                    </td>
                    <td className="px-4 py-4 text-sm liquid-text-secondary">
                      {formatDate(doc.createdAt)}
                    </td>
                    <td className="px-4 py-4 text-right">
                      <LiquidGlassButton
                        variant="primary"
                        size="sm"
                        onClick={(e) => {
                          e.stopPropagation()
                          onSelectDocument?.(doc)
                        }}
                      >
                        {t('documents.queue.review', 'Review')}
                      </LiquidGlassButton>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </LiquidGlassCard>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm liquid-text-muted">
            {t('documents.queue.showing', 'Showing {{from}}-{{to}} of {{total}}', {
              from: (page - 1) * pageSize + 1,
              to: Math.min(page * pageSize, total),
              total,
            })}
          </p>
          <div className="flex gap-2 flex-wrap">
            <LiquidGlassButton
              variant="glass"
              size="sm"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              {t('common.previous', 'Previous')}
            </LiquidGlassButton>
            <LiquidGlassButton
              variant="glass"
              size="sm"
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
            >
              {t('common.next', 'Next')}
            </LiquidGlassButton>
          </div>
        </div>
      )}
    </div>
  )
}

/**
 * Map document status to badge variant
 */
function getStatusBadgeVariant(status: DocumentStatus): 'warning' | 'blue' | 'success' | 'error' | 'default' {
  switch (status) {
    case 'ready_for_review':
      return 'warning'
    case 'under_review':
      return 'blue'
    case 'approved':
      return 'success'
    case 'rejected':
      return 'error'
    default:
      return 'default'
  }
}

function StatCard({
  label,
  value,
  variant,
}: {
  label: string
  value: number
  variant: 'warning' | 'blue' | 'error' | 'default'
}) {
  const variantClasses = {
    warning: 'liquid-glass-badge-yellow',
    blue: 'liquid-glass-badge-blue',
    error: 'liquid-glass-badge-red',
    default: '',
  }

  return (
    <LiquidGlassCard
      intensity="subtle"
      hover="lift"
      className={`p-4 ${variantClasses[variant]}`}
    >
      <p className="text-2xl font-bold liquid-text-primary">{value}</p>
      <p className="text-sm liquid-text-secondary">{label}</p>
    </LiquidGlassCard>
  )
}

export default ReviewQueue
