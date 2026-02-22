// Package warehouse provides data access interfaces for the Council of AIs system.
//
// This package implements repository patterns for Knowledge Graph access,
// following MediSync's read-only data plane architecture.
package warehouse

import (
	"context"
	"time"

	"github.com/pgvector/pgvector-go"
)

// KGNodeType represents the type of a Knowledge Graph node.
type KGNodeType string

const (
	NodeTypeConcept      KGNodeType = "concept"
	NodeTypeMedication   KGNodeType = "medication"
	NodeTypeProcedure    KGNodeType = "procedure"
	NodeTypeCondition    KGNodeType = "condition"
	NodeTypeOrganization KGNodeType = "organization"
)

// KGEdgeType represents the type of relationship between Knowledge Graph nodes.
type KGEdgeType string

const (
	EdgeTreats         KGEdgeType = "TREATS"
	EdgeCauses         KGEdgeType = "CAUSES"
	EdgeContraindicates KGEdgeType = "CONTRAINDICATES"
	EdgeRelatedTo      KGEdgeType = "RELATED_TO"
	EdgeSubsumes       KGEdgeType = "SUBSUMES"
	EdgePartOf         KGEdgeType = "PART_OF"
)

// KnowledgeGraphNode represents a unit of verified medical knowledge.
type KnowledgeGraphNode struct {
	ID           string          `json:"id"`
	NodeType     KGNodeType      `json:"node_type"`
	Concept      string          `json:"concept"`
	Definition   string          `json:"definition"`
	Embedding    pgvector.Vector `json:"embedding"`
	Source       string          `json:"source"`
	SourceID     string          `json:"source_id,omitempty"`
	Confidence   float64         `json:"confidence"`
	LastVerified time.Time       `json:"last_verified"`
	Edges        []string        `json:"edges"`
	EdgeTypes    []KGEdgeType    `json:"edge_types"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// TraversalStep represents a single step in the Graph-of-Thoughts traversal.
type TraversalStep struct {
	FromNodeID string     `json:"from_node_id"`
	ToNodeID   string     `json:"to_node_id"`
	EdgeType   KGEdgeType `json:"edge_type"`
	Weight     float64    `json:"weight"`
}

// KnowledgeGraphRepository defines the interface for Knowledge Graph data access.
//
// All methods use the medisync_readonly role to ensure read-only access
// to the Knowledge Graph, following MediSync's security architecture.
type KnowledgeGraphRepository interface {
	// FindSimilar finds nodes with embeddings similar to the query embedding.
	// Returns up to limit nodes sorted by similarity (descending).
	FindSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]*KnowledgeGraphNode, error)

	// GetNode retrieves a single Knowledge Graph node by ID.
	GetNode(ctx context.Context, id string) (*KnowledgeGraphNode, error)

	// GetNodes retrieves multiple Knowledge Graph nodes by IDs.
	GetNodes(ctx context.Context, ids []string) ([]*KnowledgeGraphNode, error)

	// GetRelatedNodes retrieves nodes related to a given node via specific edge types.
	GetRelatedNodes(ctx context.Context, nodeID string, edgeTypes []KGEdgeType, limit int) ([]*KnowledgeGraphNode, error)

	// TraverseMultiHop performs multi-hop traversal starting from initial nodes.
	// Returns the complete traversal path with relevance scores.
	TraverseMultiHop(ctx context.Context, initialNodeIDs []string, maxHops int) (*TraversalResult, error)

	// HealthCheck verifies Knowledge Graph availability.
	HealthCheck(ctx context.Context) error
}

// TraversalResult represents the result of a multi-hop graph traversal.
type TraversalResult struct {
	Nodes          []*KnowledgeGraphNode
	TraversalPath  []*TraversalStep
	RelevanceScore map[string]float64 // nodeID → relevance score
	HopCount       int
}

// KnowledgeGraphRepo implements KnowledgeGraphRepository using PostgreSQL + pgvector.
type KnowledgeGraphRepo struct {
	db DBTX
}

// DBTX represents a database connection or transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (any, error)
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) Row
}

// Rows represents database query results.
type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
	Err() error
}

// Row represents a single database row.
type Row interface {
	Scan(dest ...any) error
}

// NewKnowledgeGraphRepo creates a new Knowledge Graph repository.
func NewKnowledgeGraphRepo(db DBTX) *KnowledgeGraphRepo {
	return &KnowledgeGraphRepo{db: db}
}

// FindSimilar finds nodes with embeddings similar to the query embedding.
func (r *KnowledgeGraphRepo) FindSimilar(ctx context.Context, embedding pgvector.Vector, limit int) ([]*KnowledgeGraphNode, error) {
	query := `
		SELECT id, node_type, concept, definition, embedding, source, source_id,
		       confidence, last_verified, edges, edge_types, created_at, updated_at
		FROM knowledge_graph_nodes
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, embedding, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*KnowledgeGraphNode
	for rows.Next() {
		node := &KnowledgeGraphNode{}
		var edges, edgeTypes []string
		err := rows.Scan(
			&node.ID, &node.NodeType, &node.Concept, &node.Definition,
			&node.Embedding, &node.Source, &node.SourceID,
			&node.Confidence, &node.LastVerified,
			&edges, &edgeTypes,
			&node.CreatedAt, &node.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Convert string arrays to typed slices
		node.Edges = edges
		for _, et := range edgeTypes {
			node.EdgeTypes = append(node.EdgeTypes, KGEdgeType(et))
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetNode retrieves a single Knowledge Graph node by ID.
func (r *KnowledgeGraphRepo) GetNode(ctx context.Context, id string) (*KnowledgeGraphNode, error) {
	query := `
		SELECT id, node_type, concept, definition, embedding, source, source_id,
		       confidence, last_verified, edges, edge_types, created_at, updated_at
		FROM knowledge_graph_nodes
		WHERE id = $1
	`

	node := &KnowledgeGraphNode{}
	var edges, edgeTypes []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&node.ID, &node.NodeType, &node.Concept, &node.Definition,
		&node.Embedding, &node.Source, &node.SourceID,
		&node.Confidence, &node.LastVerified,
		&edges, &edgeTypes,
		&node.CreatedAt, &node.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	node.Edges = edges
	for _, et := range edgeTypes {
		node.EdgeTypes = append(node.EdgeTypes, KGEdgeType(et))
	}

	return node, nil
}

// GetNodes retrieves multiple Knowledge Graph nodes by IDs.
func (r *KnowledgeGraphRepo) GetNodes(ctx context.Context, ids []string) ([]*KnowledgeGraphNode, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `
		SELECT id, node_type, concept, definition, embedding, source, source_id,
		       confidence, last_verified, edges, edge_types, created_at, updated_at
		FROM knowledge_graph_nodes
		WHERE id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*KnowledgeGraphNode
	for rows.Next() {
		node := &KnowledgeGraphNode{}
		var edges, edgeTypes []string
		err := rows.Scan(
			&node.ID, &node.NodeType, &node.Concept, &node.Definition,
			&node.Embedding, &node.Source, &node.SourceID,
			&node.Confidence, &node.LastVerified,
			&edges, &edgeTypes,
			&node.CreatedAt, &node.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		node.Edges = edges
		for _, et := range edgeTypes {
			node.EdgeTypes = append(node.EdgeTypes, KGEdgeType(et))
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// GetRelatedNodes retrieves nodes related to a given node via specific edge types.
func (r *KnowledgeGraphRepo) GetRelatedNodes(ctx context.Context, nodeID string, edgeTypes []KGEdgeType, limit int) ([]*KnowledgeGraphNode, error) {
	query := `
		WITH node_edges AS (
			SELECT edges, edge_types
			FROM knowledge_graph_nodes
			WHERE id = $1
		)
		SELECT n.id, n.node_type, n.concept, n.definition, n.embedding, n.source, n.source_id,
		       n.confidence, n.last_verified, n.edges, n.edge_types, n.created_at, n.updated_at
		FROM knowledge_graph_nodes n, node_edges ne
		WHERE n.id = ANY(ne.edges)
		  AND (
			SELECT COUNT(*) FROM unnest(ne.edge_types) WITH ORDINALITY AS et(type, idx)
			WHERE et.type = ANY($2) AND ne.edges[et.idx] = n.id
		  ) > 0
		LIMIT $3
	`

	// Convert edge types to strings
	typeStrs := make([]string, len(edgeTypes))
	for i, et := range edgeTypes {
		typeStrs[i] = string(et)
	}

	rows, err := r.db.QueryContext(ctx, query, nodeID, typeStrs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*KnowledgeGraphNode
	for rows.Next() {
		node := &KnowledgeGraphNode{}
		var edges, edgeTypes []string
		err := rows.Scan(
			&node.ID, &node.NodeType, &node.Concept, &node.Definition,
			&node.Embedding, &node.Source, &node.SourceID,
			&node.Confidence, &node.LastVerified,
			&edges, &edgeTypes,
			&node.CreatedAt, &node.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		node.Edges = edges
		for _, et := range edgeTypes {
			node.EdgeTypes = append(node.EdgeTypes, KGEdgeType(et))
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

// TraverseMultiHop performs multi-hop traversal starting from initial nodes.
func (r *KnowledgeGraphRepo) TraverseMultiHop(ctx context.Context, initialNodeIDs []string, maxHops int) (*TraversalResult, error) {
	// Use recursive CTE for multi-hop traversal
	query := `
		WITH RECURSIVE traversal AS (
			-- Base case: initial nodes
			SELECT
				id, node_type, concept, definition, embedding, source, source_id,
				confidence, last_verified, edges, edge_types, created_at, updated_at,
				0 AS hop, ARRAY[id] AS path, NULL::uuid AS from_node, NULL::text AS edge_type
			FROM knowledge_graph_nodes
			WHERE id = ANY($1)

			UNION ALL

			-- Recursive case: follow edges
			SELECT
				n.id, n.node_type, n.concept, n.definition, n.embedding, n.source, n.source_id,
				n.confidence, n.last_verified, n.edges, n.edge_types, n.created_at, n.updated_at,
				t.hop + 1, t.path || n.id, t.id, ne.edge_types[array_position(t.edges, n.id)]
			FROM knowledge_graph_nodes n
			JOIN traversal t ON n.id = ANY(t.edges)
			CROSS JOIN LATERAL (
				SELECT edge_types FROM knowledge_graph_nodes WHERE id = t.id
			) ne
			WHERE t.hop < $2
			  AND NOT (n.id = ANY(t.path))  -- Prevent cycles
		)
		SELECT * FROM traversal ORDER BY hop, path
	`

	rows, err := r.db.QueryContext(ctx, query, initialNodeIDs, maxHops)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &TraversalResult{
		Nodes:          make([]*KnowledgeGraphNode, 0),
		TraversalPath:  make([]*TraversalStep, 0),
		RelevanceScore: make(map[string]float64),
	}

	seenNodes := make(map[string]bool)

	for rows.Next() {
		node := &KnowledgeGraphNode{}
		var hop int
		var path []string
		var fromNode *string
		var edgeType *string
		var edges, edgeTypes []string

		err := rows.Scan(
			&node.ID, &node.NodeType, &node.Concept, &node.Definition,
			&node.Embedding, &node.Source, &node.SourceID,
			&node.Confidence, &node.LastVerified,
			&edges, &edgeTypes,
			&node.CreatedAt, &node.UpdatedAt,
			&hop, &path, &fromNode, &edgeType,
		)
		if err != nil {
			return nil, err
		}

		node.Edges = edges
		for _, et := range edgeTypes {
			node.EdgeTypes = append(node.EdgeTypes, KGEdgeType(et))
		}

		// Add node if not already seen
		if !seenNodes[node.ID] {
			result.Nodes = append(result.Nodes, node)
			seenNodes[node.ID] = true
			// Calculate relevance based on hop distance
			result.RelevanceScore[node.ID] = 1.0 / float64(hop+1)
		}

		// Record traversal step (except for initial nodes)
		if fromNode != nil && edgeType != nil {
			result.TraversalPath = append(result.TraversalPath, &TraversalStep{
				FromNodeID: *fromNode,
				ToNodeID:   node.ID,
				EdgeType:   KGEdgeType(*edgeType),
				Weight:     1.0 / float64(hop+1),
			})
		}

		if hop > result.HopCount {
			result.HopCount = hop
		}
	}

	return result, rows.Err()
}

// HealthCheck verifies Knowledge Graph availability.
func (r *KnowledgeGraphRepo) HealthCheck(ctx context.Context) error {
	query := `SELECT 1 FROM knowledge_graph_nodes LIMIT 1`
	var dummy int
	return r.db.QueryRowContext(ctx, query).Scan(&dummy)
}

// Compile-time interface compliance check
var _ KnowledgeGraphRepository = (*KnowledgeGraphRepo)(nil)
