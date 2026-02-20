// Package shared provides AgentCard metadata for ADK (Agent Development Kit) integration.
//
// AgentCard follows the Google A2A Protocol specification for agent discovery
// and capability advertisement. Each agent registers its capabilities, endpoint,
// and metadata to enable inter-agent communication.
//
// Usage:
//
//	card := shared.GetAgentCard("a-01-text-to-sql")
//	if card != nil {
//	    fmt.Printf("Agent: %s\n", card.Name)
//	}
package shared

import "time"

// AgentCard represents metadata about an agent following ADK/A2A specification.
type AgentCard struct {
	// ID is the unique identifier for the agent (e.g., "a-01-text-to-sql")
	ID string `json:"id"`

	// Name is the human-readable name of the agent
	Name string `json:"name"`

	// Description provides a detailed description of the agent's purpose
	Description string `json:"description"`

	// Capabilities lists the agent's capabilities
	Capabilities []string `json:"capabilities"`

	// Endpoint is the HTTP/gRPC endpoint for the agent
	Endpoint string `json:"endpoint"`

	// Version is the agent's version string
	Version string `json:"version"`

	// Module identifies which module the agent belongs to (A, B, C, D, E)
	Module string `json:"module"`

	// Metadata contains additional agent-specific metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// CreatedAt is when the agent was registered
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the agent was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// AgentCards contains all registered agent metadata.
// This map serves as the central registry for agent discovery.
var AgentCards = map[string]AgentCard{
	// ============================================================================
	// Module A: Conversational BI Agents
	// ============================================================================

	"a-01-text-to-sql": {
		ID:          "a-01-text-to-sql",
		Name:        "Text-to-SQL Agent",
		Description: "Converts natural language queries to safe, validated SQL queries using schema context and domain terminology.",
		Capabilities: []string{
			"natural_language_to_sql",
			"schema_context_retrieval",
			"parameterized_query_generation",
			"sql_validation",
		},
		Endpoint: "/api/v1/agents/a-01",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"model":        "llama4",
			"max_tokens":   "2048",
			"temperature":  "0.1",
			"supports_rag": "true",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"a-02-sql-correction": {
		ID:          "a-02-sql-correction",
		Name:        "SQL Self-Correction Agent",
		Description: "Detects SQL errors, analyzes root causes, and attempts automatic correction with retry logic.",
		Capabilities: []string{
			"error_detection",
			"root_cause_analysis",
			"automatic_correction",
			"retry_orchestration",
		},
		Endpoint: "/api/v1/agents/a-02",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"max_retries":      "3",
			"timeout_seconds":  "30",
			"fallback_enabled": "true",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"a-03-visualization-routing": {
		ID:          "a-03-visualization-routing",
		Name:        "Visualization Routing Agent",
		Description: "Analyzes query results and data characteristics to recommend optimal visualization types.",
		Capabilities: []string{
			"chart_type_selection",
			"data_pattern_analysis",
			"accessibility_optimization",
			"multi_series_detection",
		},
		Endpoint: "/api/v1/agents/a-03",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"supported_charts": "lineChart,barChart,pieChart,kpiCard,dataTable,scatterChart",
			"min_confidence":   "0.7",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"a-04-terminology-normalizer": {
		ID:          "a-04-terminology-normalizer",
		Name:        "Domain Terminology Normalizer",
		Description: "Maps healthcare and accounting domain terms to canonical database schema references.",
		Capabilities: []string{
			"term_recognition",
			"synonym_mapping",
			"domain_context_enrichment",
			"multi_language_support",
		},
		Endpoint: "/api/v1/agents/a-04",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"domains":          "healthcare,accounting,pharmacy",
			"languages":        "en,ar",
			"embedding_source": "pgvector",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"a-05-hallucination-guard": {
		ID:          "a-05-hallucination-guard",
		Name:        "Hallucination Guard Agent",
		Description: "Validates AI outputs against schema constraints and detects fabricated information.",
		Capabilities: []string{
			"schema_alignment_check",
			"fabrication_detection",
			"data_consistency_validation",
			"risk_scoring",
		},
		Endpoint: "/api/v1/agents/a-05",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"max_risk_score":   "30",
			"check_schema":     "true",
			"check_data":       "true",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"a-06-confidence-scorer": {
		ID:          "a-06-confidence-scorer",
		Name:        "Confidence Scoring Agent",
		Description: "Aggregates confidence signals from all pipeline stages and determines if clarification is needed.",
		Capabilities: []string{
			"multi_factor_scoring",
			"clarification_triggering",
			"uncertainty_quantification",
			"quality_assessment",
		},
		Endpoint: "/api/v1/agents/a-06",
		Version:  "1.0.0",
		Module:   "A",
		Metadata: map[string]string{
			"clarification_threshold": "70",
			"min_overall_confidence":  "60",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	// ============================================================================
	// Module E: i18n Agents
	// ============================================================================

	"e-01-language-detection": {
		ID:          "e-01-language-detection",
		Name:        "Language Detection Agent",
		Description: "Detects the language of user queries (English or Arabic) using script analysis and ML models.",
		Capabilities: []string{
			"language_detection",
			"script_identification",
			"mixed_language_handling",
			"confidence_scoring",
		},
		Endpoint: "/api/v1/agents/e-01",
		Version:  "1.0.0",
		Module:   "E",
		Metadata: map[string]string{
			"supported_languages": "en,ar",
			"detection_method":    "hybrid",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"e-02-query-translation": {
		ID:          "e-02-query-translation",
		Name:        "Query Translation Agent",
		Description: "Translates Arabic queries to English for SQL generation while preserving semantic intent.",
		Capabilities: []string{
			"arabic_to_english",
			"semantic_preservation",
			"domain_aware_translation",
			"back_translation_validation",
		},
		Endpoint: "/api/v1/agents/e-02",
		Version:  "1.0.0",
		Module:   "E",
		Metadata: map[string]string{
			"source_languages":  "ar",
			"target_language":   "en",
			"validation_method": "back_translation",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},

	"e-03-localized-formatter": {
		ID:          "e-03-localized-formatter",
		Name:        "Localized Formatter Agent",
		Description: "Formats numbers, dates, and currency according to user locale preferences with RTL support.",
		Capabilities: []string{
			"number_formatting",
			"currency_formatting",
			"date_formatting",
			"rtl_layout_support",
		},
		Endpoint: "/api/v1/agents/e-03",
		Version:  "1.0.0",
		Module:   "E",
		Metadata: map[string]string{
			"locales":          "en-IN,ar-SA",
			"calendars":        "gregorian,islamic",
			"currency_codes":   "INR,SAR",
		},
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
	},
}

// GetAgentCard retrieves an agent's card by ID.
// Returns nil if the agent ID is not found.
func GetAgentCard(id string) *AgentCard {
	card, exists := AgentCards[id]
	if !exists {
		return nil
	}
	return &card
}

// GetAllAgentCards returns a slice of all registered agent cards.
func GetAllAgentCards() []AgentCard {
	cards := make([]AgentCard, 0, len(AgentCards))
	for _, card := range AgentCards {
		cards = append(cards, card)
	}
	return cards
}

// GetAgentsByModule returns all agents belonging to a specific module.
func GetAgentsByModule(module string) []AgentCard {
	var cards []AgentCard
	for _, card := range AgentCards {
		if card.Module == module {
			cards = append(cards, card)
		}
	}
	return cards
}

// GetAgentCapabilities returns the capabilities of a specific agent.
func GetAgentCapabilities(id string) []string {
	card := GetAgentCard(id)
	if card == nil {
		return nil
	}
	return card.Capabilities
}

// HasCapability checks if an agent has a specific capability.
func HasCapability(agentID, capability string) bool {
	capabilities := GetAgentCapabilities(agentID)
	if capabilities == nil {
		return false
	}
	for _, c := range capabilities {
		if c == capability {
			return true
		}
	}
	return false
}

// AgentIDs returns all registered agent IDs.
func AgentIDs() []string {
	ids := make([]string, 0, len(AgentCards))
	for id := range AgentCards {
		ids = append(ids, id)
	}
	return ids
}

// AgentCount returns the number of registered agents.
func AgentCount() int {
	return len(AgentCards)
}
