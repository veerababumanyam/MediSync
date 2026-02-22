/**
 * ResponseDisplay Component
 *
 * Displays Council deliberation results with:
 * - Consensus response with confidence indicator
 * - Status badge (consensus, uncertain, failed)
 * - Agent response breakdown
 * - Evidence trail summary
 * - RTL support
 *
 * @module components/council/ResponseDisplay
 */

import { useTranslation } from 'react-i18next'
import { LiquidGlassCard } from '../ui/LiquidGlassCard'
import { LiquidGlassBadge } from '../ui/LiquidGlassBadge'
import { ConfidenceIndicator } from './ConfidenceIndicator'
import { cn } from '@/lib/cn'
import type {
  Deliberation,
  DeliberationStatus,
  AgentResponse,
} from '@/services/councilService'

/**
 * Props for ResponseDisplay component
 */
export interface ResponseDisplayProps {
  /** Deliberation result to display */
  deliberation: Deliberation
  /** Show detailed agent responses */
  showAgentResponses?: boolean
  /** Show evidence trail summary */
  showEvidenceSummary?: boolean
  /** Additional CSS classes */
  className?: string
}

/**
 * Status badge color mapping
 */
const statusColors: Record<DeliberationStatus, string> = {
  pending: 'bg-gray-500/20 text-gray-700',
  deliberating: 'bg-blue-500/20 text-blue-700',
  consensus: 'bg-green-500/20 text-green-700',
  uncertain: 'bg-amber-500/20 text-amber-700',
  failed: 'bg-red-500/20 text-red-700',
}

/**
 * Response display component for Council deliberations
 */
export function ResponseDisplay({
  deliberation,
  showAgentResponses = true,
  showEvidenceSummary = true,
  className,
}: ResponseDisplayProps) {
  const { t } = useTranslation('council')

  const {
    query,
    status,
    final_response: finalResponse,
    confidence_score: confidenceScore,
    consensus_record: consensusRecord,
    agent_responses: agentResponses,
    evidence_trail: evidenceTrail,
    created_at: createdAt,
    completed_at: completedAt,
  } = deliberation

  const isComplete = status === 'consensus' || status === 'uncertain' || status === 'failed'

  return (
    <LiquidGlassCard className={cn('p-6', className)}>
      {/* Header with Status */}
      <div className="flex items-start justify-between gap-4 mb-6">
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-primary truncate">
            {t('responseDisplay.title', 'Council Response')}
          </h3>
          <p className="text-sm text-secondary mt-1 line-clamp-2">{query}</p>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <LiquidGlassBadge className={statusColors[status]}>
            {t(`status.${status}`, status)}
          </LiquidGlassBadge>
        </div>
      </div>

      {/* Loading State */}
      {status === 'deliberating' && (
        <div className="flex items-center justify-center py-12">
          <div className="flex flex-col items-center gap-4">
            <div className="relative w-16 h-16">
              <div className="absolute inset-0 rounded-full border-4 border-surface-glass-strong animate-pulse" />
              <div className="absolute inset-2 rounded-full border-4 border-primary border-t-transparent animate-spin" />
            </div>
            <p className="text-secondary animate-pulse">
              {t('responseDisplay.deliberating', 'Council is deliberating...')}
            </p>
          </div>
        </div>
      )}

      {/* Consensus Response */}
      {isComplete && finalResponse && (
        <div className="space-y-6">
          {/* Confidence & Agreement */}
          <div className="flex items-center justify-between gap-4 pb-4 border-b border-glass">
            {confidenceScore !== undefined && (
              <ConfidenceIndicator
                value={confidenceScore}
                showLabel
                size="lg"
              />
            )}
            {consensusRecord && (
              <div className="text-right">
                <p className="text-xs text-secondary">
                  {t('responseDisplay.agreement', 'Agreement Score')}
                </p>
                <p className="text-lg font-semibold text-primary">
                  {Math.round(consensusRecord.agreement_score * 100)}%
                </p>
              </div>
            )}
          </div>

          {/* Final Response */}
          <div className="space-y-2">
            <h4 className="text-sm font-medium text-secondary">
              {t('responseDisplay.response', 'Consensus Response')}
            </h4>
            <div className="p-4 rounded-xl bg-surface-glass/50 backdrop-blur-sm">
              <p className="text-primary whitespace-pre-wrap">{finalResponse}</p>
            </div>
          </div>

          {/* Agent Responses */}
          {showAgentResponses && agentResponses && agentResponses.length > 0 && (
            <div className="space-y-3">
              <h4 className="text-sm font-medium text-secondary">
                {t('responseDisplay.agentResponses', {
                  count: agentResponses.length,
                  defaultValue: '{{count}} Agent Responses',
                })}
              </h4>
              <div className="grid gap-3">
                {agentResponses.map((response) => (
                  <AgentResponseCard
                    key={response.id}
                    response={response}
                    isCanonical={
                      consensusRecord?.equivalence_groups[0]?.agent_ids.includes(
                        response.agent_id
                      ) ?? false
                    }
                  />
                ))}
              </div>
            </div>
          )}

          {/* Evidence Summary */}
          {showEvidenceSummary && evidenceTrail && (
            <div className="space-y-2">
              <h4 className="text-sm font-medium text-secondary">
                {t('responseDisplay.evidence', 'Supporting Evidence')}
              </h4>
              <div className="flex items-center gap-4 text-sm">
                <span className="flex items-center gap-1.5">
                  <svg
                    className="w-4 h-4 text-primary"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                  {evidenceTrail.node_ids.length} {t('responseDisplay.nodes', 'knowledge nodes')}
                </span>
                <span className="flex items-center gap-1.5">
                  <svg
                    className="w-4 h-4 text-primary"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                    />
                  </svg>
                  {evidenceTrail.hop_count} {t('responseDisplay.hops', 'hops')}
                </span>
              </div>
            </div>
          )}

          {/* Timestamps */}
          <div className="flex items-center justify-between text-xs text-secondary pt-4 border-t border-glass">
            <span>
              {t('responseDisplay.created', 'Created')}: {new Date(createdAt).toLocaleString()}
            </span>
            {completedAt && (
              <span>
                {t('responseDisplay.completed', 'Completed')}:{' '}
                {new Date(completedAt).toLocaleString()}
              </span>
            )}
          </div>
        </div>
      )}

      {/* Error State */}
      {status === 'failed' && !finalResponse && (
        <div className="p-4 rounded-xl bg-red-500/10 border border-red-500/30">
          <p className="text-red-600">
            {t('responseDisplay.failed', 'Deliberation failed. Please try again.')}
          </p>
        </div>
      )}

      {/* Uncertain State */}
      {status === 'uncertain' && consensusRecord && !consensusRecord.threshold_met && (
        <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/30">
          <p className="text-amber-600">
            {t(
              'responseDisplay.uncertain',
              'The Council could not reach consensus. Review the individual agent responses below.'
            )}
          </p>
        </div>
      )}
    </LiquidGlassCard>
  )
}

/**
 * Agent response card subcomponent
 */
function AgentResponseCard({
  response,
  isCanonical,
}: {
  response: AgentResponse
  isCanonical: boolean
}) {
  const { t } = useTranslation('council')

  return (
    <div
      className={cn(
        'p-3 rounded-lg transition-all',
        'bg-surface-glass/30 backdrop-blur-sm',
        'border border-glass',
        isCanonical && 'border-primary/50 bg-primary/5'
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm font-medium text-primary">
              {response.agent_id}
            </span>
            {isCanonical && (
              <LiquidGlassBadge className="text-xs bg-primary/20 text-primary">
                {t('responseDisplay.canonical', 'Canonical')}
              </LiquidGlassBadge>
            )}
          </div>
          <p className="text-sm text-secondary line-clamp-2">
            {response.response_text}
          </p>
        </div>
        <div className="shrink-0 text-right">
          <span
            className={cn(
              'text-sm font-semibold',
              response.confidence >= 80
                ? 'text-green-600'
                : response.confidence >= 50
                  ? 'text-amber-600'
                  : 'text-red-600'
            )}
          >
            {Math.round(response.confidence)}%
          </span>
          <p className="text-xs text-secondary">
            {t('responseDisplay.confidence', 'confidence')}
          </p>
        </div>
      </div>
    </div>
  )
}

export default ResponseDisplay
