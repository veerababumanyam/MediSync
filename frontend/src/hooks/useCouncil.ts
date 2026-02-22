/**
 * useCouncil Hook
 *
 * React hook for managing Council of AIs deliberations with:
 * - Deliberation creation and state management
 * - Loading and error states
 * - Deliberation listing with pagination
 * - Health status monitoring
 *
 * @module hooks/useCouncil
 */

import { useCallback, useState } from 'react'
import {
  councilService,
  APIError,
  type Deliberation,
  type DeliberationListResponse,
  type DeliberationStatus,
  type AgentResponse,
  type ConsensusRecord,
  type CreateDeliberationRequest,
  type ListDeliberationsParams,
  type CouncilHealthResponse,
  type EvidenceTrail,
} from '../services/councilService'

/**
 * Hook state for current deliberation
 */
export interface CouncilState {
  /** Current deliberation result (if any) */
  deliberation: Deliberation | null
  /** Whether a request is in progress */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
}

/**
 * Hook return type
 */
export interface UseCouncilReturn extends CouncilState {
  /** Create a new deliberation */
  createDeliberation: (query: string, threshold?: number) => Promise<Deliberation | null>
  /** Get a specific deliberation by ID */
  getDeliberation: (id: string) => Promise<Deliberation | null>
  /** Clear the current deliberation and error */
  reset: () => void
}

/**
 * Hook for managing Council deliberations
 */
export function useCouncil(): UseCouncilReturn {
  const [state, setState] = useState<CouncilState>({
    deliberation: null,
    isLoading: false,
    error: null,
  })

  /**
   * Create a new deliberation
   */
  const createDeliberation = useCallback(
    async (query: string, threshold?: number): Promise<Deliberation | null> => {
      setState((prev) => ({ ...prev, isLoading: true, error: null }))

      try {
        const request: CreateDeliberationRequest = {
          query,
          consensus_threshold: threshold,
        }

        const deliberation = await councilService.createDeliberation(request)
        setState({ deliberation, isLoading: false, error: null })
        return deliberation
      } catch (err) {
        const message =
          err instanceof APIError ? err.message : 'Failed to create deliberation'
        setState((prev) => ({ ...prev, isLoading: false, error: message }))
        return null
      }
    },
    []
  )

  /**
   * Get a specific deliberation by ID
   */
  const getDeliberation = useCallback(async (id: string): Promise<Deliberation | null> => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }))

    try {
      const deliberation = await councilService.getDeliberation(id)
      setState({ deliberation, isLoading: false, error: null })
      return deliberation
    } catch (err) {
      const message =
        err instanceof APIError ? err.message : 'Failed to get deliberation'
      setState((prev) => ({ ...prev, isLoading: false, error: message }))
      return null
    }
  }, [])

  /**
   * Reset state
   */
  const reset = useCallback(() => {
    setState({ deliberation: null, isLoading: false, error: null })
  }, [])

  return {
    ...state,
    createDeliberation,
    getDeliberation,
    reset,
  }
}

// ============================================================================
// useCouncilList Hook
// ============================================================================

/**
 * List hook state
 */
export interface CouncilListState {
  /** List of deliberations */
  deliberations: Deliberation[]
  /** Total count for pagination */
  total: number
  /** Whether a request is in progress */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
}

/**
 * List hook return type
 */
export interface UseCouncilListReturn extends CouncilListState {
  /** Fetch deliberations with optional filters */
  fetchDeliberations: (params?: ListDeliberationsParams) => Promise<void>
  /** Load next page */
  loadMore: () => Promise<void>
  /** Refresh the current list */
  refresh: () => Promise<void>
}

/**
 * Hook for listing deliberations with pagination
 */
export function useCouncilList(): UseCouncilListReturn {
  const [state, setState] = useState<CouncilListState>({
    deliberations: [],
    total: 0,
    isLoading: false,
    error: null,
  })

  const [lastParams, setLastParams] = useState<ListDeliberationsParams>({
    limit: 20,
    offset: 0,
  })

  /**
   * Fetch deliberations with optional filters
   */
  const fetchDeliberations = useCallback(async (params?: ListDeliberationsParams) => {
    const mergedParams = { limit: 20, offset: 0, ...params }
    setLastParams(mergedParams)
    setState((prev) => ({ ...prev, isLoading: true, error: null }))

    try {
      const response: DeliberationListResponse = await councilService.listDeliberations(
        mergedParams
      )
      setState({
        deliberations: response.deliberations,
        total: response.total,
        isLoading: false,
        error: null,
      })
    } catch (err) {
      const message =
        err instanceof APIError ? err.message : 'Failed to fetch deliberations'
      setState((prev) => ({ ...prev, isLoading: false, error: message }))
    }
  }, [])

  /**
   * Load next page
   */
  const loadMore = useCallback(async () => {
    const { offset, limit } = lastParams
    const currentOffset = offset ?? 0
    const pageSize = limit ?? 20
    const newOffset = currentOffset + (state.deliberations.length || pageSize)

    setState((prev) => ({ ...prev, isLoading: true, error: null }))

    try {
      const response = await councilService.listDeliberations({
        ...lastParams,
        offset: newOffset,
      })
      setState((prev) => ({
        deliberations: [...prev.deliberations, ...response.deliberations],
        total: response.total,
        isLoading: false,
        error: null,
      }))
      setLastParams((prev) => ({ ...prev, offset: newOffset }))
    } catch (err) {
      const message =
        err instanceof APIError ? err.message : 'Failed to load more'
      setState((prev) => ({ ...prev, isLoading: false, error: message }))
    }
  }, [lastParams, state.deliberations.length])

  /**
   * Refresh the current list
   */
  const refresh = useCallback(async () => {
    await fetchDeliberations({ ...lastParams, offset: 0 })
  }, [fetchDeliberations, lastParams])

  return {
    ...state,
    fetchDeliberations,
    loadMore,
    refresh,
  }
}

// ============================================================================
// useCouncilHealth Hook
// ============================================================================

/**
 * Health hook state
 */
export interface CouncilHealthState {
  /** Health status */
  health: CouncilHealthResponse | null
  /** Whether a request is in progress */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
}

/**
 * Health hook return type
 */
export interface UseCouncilHealthReturn extends CouncilHealthState {
  /** Fetch health status */
  fetchHealth: () => Promise<void>
}

/**
 * Hook for monitoring Council health
 */
export function useCouncilHealth(): UseCouncilHealthReturn {
  const [state, setState] = useState<CouncilHealthState>({
    health: null,
    isLoading: false,
    error: null,
  })

  /**
   * Fetch health status
   */
  const fetchHealth = useCallback(async () => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }))

    try {
      const health = await councilService.getHealth()
      setState({ health, isLoading: false, error: null })
    } catch (err) {
      const message =
        err instanceof APIError ? err.message : 'Failed to fetch health'
      setState((prev) => ({ ...prev, isLoading: false, error: message }))
    }
  }, [])

  return {
    ...state,
    fetchHealth,
  }
}

// ============================================================================
// useEvidenceTrail Hook
// ============================================================================

/**
 * Evidence trail hook state
 */
export interface EvidenceTrailState {
  /** Evidence trail data */
  trail: EvidenceTrail | null
  /** Whether a request is in progress */
  isLoading: boolean
  /** Any error that occurred */
  error: string | null
}

/**
 * Evidence trail hook return type
 */
export interface UseEvidenceTrailReturn extends EvidenceTrailState {
  /** Fetch evidence trail for a deliberation */
  fetchTrail: (deliberationId: string) => Promise<void>
}

/**
 * Hook for fetching evidence trails
 */
export function useEvidenceTrail(): UseEvidenceTrailReturn {
  const [state, setState] = useState<EvidenceTrailState>({
    trail: null,
    isLoading: false,
    error: null,
  })

  /**
   * Fetch evidence trail
   */
  const fetchTrail = useCallback(async (deliberationId: string) => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }))

    try {
      const trail = await councilService.getEvidenceTrail(deliberationId)
      setState({ trail, isLoading: false, error: null })
    } catch (err) {
      const message =
        err instanceof APIError ? err.message : 'Failed to fetch evidence trail'
      setState((prev) => ({ ...prev, isLoading: false, error: message }))
    }
  }, [])

  return {
    ...state,
    fetchTrail,
  }
}

export type {
  Deliberation,
  DeliberationStatus,
  AgentResponse,
  ConsensusRecord,
  EvidenceTrail,
  CouncilHealthResponse,
}
