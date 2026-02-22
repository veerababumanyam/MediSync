/**
 * Council of AIs API Service
 *
 * Service for interacting with the Council deliberation API including:
 * - Creating deliberations
 * - Fetching deliberation results
 * - Listing deliberations with RBAC filtering
 * - Health monitoring
 *
 * @module services/councilService
 */

import { api, APIError } from './api'

// ============================================================================
// Types
// ============================================================================

/**
 * Deliberation status
 */
export type DeliberationStatus =
  | 'pending'
  | 'deliberating'
  | 'consensus'
  | 'uncertain'
  | 'failed'

/**
 * Agent health status
 */
export type AgentHealthStatus = 'healthy' | 'degraded' | 'failed'

/**
 * Agent response from a single AI instance
 */
export interface AgentResponse {
  id: string
  deliberation_id: string
  agent_id: string
  response_text: string
  evidence_ids: string[]
  confidence: number
  created_at: string
}

/**
 * Consensus record capturing agreement calculation
 */
export interface ConsensusRecord {
  id: string
  deliberation_id: string
  agreement_score: number
  equivalence_groups: EquivalenceGroup[]
  threshold_met: boolean
  dissenting_agents: string[]
  consensus_method: string
  created_at: string
}

/**
 * Group of semantically equivalent agent responses
 */
export interface EquivalenceGroup {
  agent_ids: string[]
  canonical_response: string
  similarity: number
}

/**
 * Traversal step in the Knowledge Graph
 */
export interface TraversalStep {
  from_node_id: string
  to_node_id: string
  edge_type: string
  weight: number
}

/**
 * Evidence trail from Knowledge Graph traversal
 */
export interface EvidenceTrail {
  id: string
  deliberation_id: string
  node_ids: string[]
  traversal_path: TraversalStep[]
  relevance_scores: Record<string, number>
  hop_count: number
  cached_at: string
  expires_at: string
}

/**
 * Full deliberation result
 */
export interface Deliberation {
  id: string
  query: string
  status: DeliberationStatus
  final_response?: string
  confidence_score?: number
  consensus_record?: ConsensusRecord
  evidence_trail?: EvidenceTrail
  agent_responses?: AgentResponse[]
  created_at: string
  completed_at?: string
}

/**
 * Paginated list of deliberations
 */
export interface DeliberationListResponse {
  deliberations: Deliberation[]
  total: number
  limit: number
  offset: number
}

/**
 * Request to create a new deliberation
 */
export interface CreateDeliberationRequest {
  query: string
  consensus_threshold?: number // Default: 0.80
  metadata?: Record<string, unknown>
}

/**
 * Council health response
 */
export interface CouncilHealthResponse {
  status: AgentHealthStatus
  total_agents: number
  healthy_agents: number
  degraded_agents: number
  failed_agents: number
  agent_statuses: Record<string, string>
  last_checked: string
}

/**
 * List deliberations query parameters
 */
export interface ListDeliberationsParams {
  status?: DeliberationStatus
  from?: string // ISO date
  to?: string // ISO date
  flagged?: boolean
  limit?: number
  offset?: number
}

// ============================================================================
// Council API Service
// ============================================================================

const COUNCIL_BASE_PATH = '/council'

export const councilService = {
  /**
   * Create a new deliberation and get the result
   */
  async createDeliberation(request: CreateDeliberationRequest): Promise<Deliberation> {
    return api.post<Deliberation>(`${COUNCIL_BASE_PATH}/deliberations`, request)
  },

  /**
   * Get a specific deliberation by ID
   */
  async getDeliberation(id: string): Promise<Deliberation> {
    return api.get<Deliberation>(`${COUNCIL_BASE_PATH}/deliberations/${id}`)
  },

  /**
   * List deliberations with optional filtering
   */
  async listDeliberations(params?: ListDeliberationsParams): Promise<DeliberationListResponse> {
    return api.get<DeliberationListResponse>(`${COUNCIL_BASE_PATH}/deliberations`, {
      params: params as Record<string, string | number | boolean>,
    })
  },

  /**
   * Get evidence trail for a deliberation
   */
  async getEvidenceTrail(deliberationId: string): Promise<EvidenceTrail> {
    return api.get<EvidenceTrail>(
      `${COUNCIL_BASE_PATH}/deliberations/${deliberationId}/evidence`
    )
  },

  /**
   * Get Council system health
   */
  async getHealth(): Promise<CouncilHealthResponse> {
    return api.get<CouncilHealthResponse>(`${COUNCIL_BASE_PATH}/health`)
  },
}

// Re-export APIError for convenience
export { APIError }
