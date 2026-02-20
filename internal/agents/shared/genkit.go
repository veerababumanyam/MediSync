// Package shared provides shared utilities for AI agents in MediSync.
//
// This file implements the Genkit framework initialization and flow registration
// for all agents across the 5 modules (BI, Accounting, Reports, Analytics, i18n).
//
// Genkit is used for AI flow orchestration with tracing and observability support.
// All agent flows are registered here and can be retrieved by name.
package shared

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// GenkitConfig holds configuration for the Genkit framework.
type GenkitConfig struct {
	// Provider is the LLM provider (openai, ollama).
	Provider string

	// Model is the default model to use.
	Model string

	// TracingEnabled enables distributed tracing.
	TracingEnabled bool

	// TracingEndpoint is the tracing collector endpoint.
	TracingEndpoint string

	// LogLevel sets the Genkit logging level.
	LogLevel string
}

// GenkitInstance represents the initialized Genkit framework.
type GenkitInstance struct {
	config *GenkitConfig
	logger *slog.Logger
	flows  map[string]*Flow
	mu     sync.RWMutex
}

// Flow represents a Genkit flow definition.
type Flow struct {
	// Name is the unique identifier for the flow.
	Name string

	// Description is a human-readable description.
	Description string

	// Module is the module this flow belongs to (A, B, C, D, E).
	Module string

	// Handler is the flow execution function.
	Handler FlowHandler

	// InputSchema defines the expected input structure.
	InputSchema interface{}

	// OutputSchema defines the expected output structure.
	OutputSchema interface{}
}

// FlowHandler is the function signature for flow handlers.
type FlowHandler func(ctx context.Context, input interface{}) (interface{}, error)

// FlowRegistry manages all registered flows.
type FlowRegistry struct {
	flows  map[string]*Flow
	logger *slog.Logger
	mu     sync.RWMutex
}

// globalRegistry is the singleton flow registry.
var globalRegistry *FlowRegistry
var registryOnce sync.Once

// GetRegistry returns the global flow registry.
func GetRegistry() *FlowRegistry {
	registryOnce.Do(func() {
		globalRegistry = &FlowRegistry{
			flows:  make(map[string]*Flow),
			logger: slog.Default(),
		}
	})
	return globalRegistry
}

// InitializeGenkit initializes the Genkit framework with the given configuration.
func InitializeGenkit(cfg *GenkitConfig, logger *slog.Logger) (*GenkitInstance, error) {
	if logger == nil {
		logger = slog.Default()
	}

	if cfg == nil {
		cfg = &GenkitConfig{
			Provider:       "openai",
			Model:          "gpt-4",
			TracingEnabled: false,
			LogLevel:       "info",
		}
	}

	instance := &GenkitInstance{
		config: cfg,
		logger: logger,
		flows:  make(map[string]*Flow),
	}

	// Register all agent flows
	if err := registerAllFlows(instance); err != nil {
		return nil, fmt.Errorf("failed to register flows: %w", err)
	}

	logger.Info("Genkit initialized",
		slog.String("provider", cfg.Provider),
		slog.String("model", cfg.Model),
		slog.Bool("tracing", cfg.TracingEnabled),
		slog.Int("flows_registered", len(instance.flows)),
	)

	return instance, nil
}

// registerAllFlows registers all agent flows from all modules.
func registerAllFlows(g *GenkitInstance) error {
	registry := GetRegistry()

	// Register Module A flows (Conversational BI)
	moduleAFlows := []struct {
		name        string
		description string
		handler     FlowHandler
	}{
		{"a01_text_to_sql", "Convert natural language to SQL queries", nil},
		{"a02_sql_correction", "Detect and fix SQL query errors", nil},
		{"a03_visualization_routing", "Choose optimal chart type for data", nil},
		{"a04_terminology_mapping", "Map healthcare/accounting terms", nil},
		{"a05_hallucination_detector", "Detect LLM hallucinations", nil},
		{"a06_confidence_scorer", "Calculate confidence scores", nil},
	}

	for _, f := range moduleAFlows {
		registry.Register(&Flow{
			Name:        f.name,
			Description: f.description,
			Module:      "A",
			Handler:     f.handler,
		})
	}

	// Register Module B flows (AI Accountant)
	moduleBFlows := []struct {
		name        string
		description string
		handler     FlowHandler
	}{
		{"b01_document_classifier", "Classify uploaded documents", nil},
		{"b02_ocr_extraction", "Extract fields from documents", nil},
		{"b03_handwriting_recognition", "Recognize handwritten content", nil},
		{"b04_validation_rules", "Apply validation rules", nil},
		{"b05_ledger_mapping", "AI-suggest Tally GL ledger", nil},
		{"b06_cost_centre_suggest", "Suggest cost centre allocation", nil},
		{"b07_gst_code_predictor", "Predict GST tax codes", nil},
		{"b08_approval_workflow", "Multi-level approval routing", nil},
		{"b09_tally_sync", "Push approved data to Tally", nil},
	}

	for _, f := range moduleBFlows {
		registry.Register(&Flow{
			Name:        f.name,
			Description: f.description,
			Module:      "B",
			Handler:     f.handler,
		})
	}

	// Register Module C flows (Easy Reports)
	moduleCFlows := []struct {
		name        string
		description string
		handler     FlowHandler
	}{
		{"c01_prebuilt_reports", "Generate pre-built MIS reports", nil},
		{"c02_consolidation", "Multi-company consolidation", nil},
		{"c03_drill_down", "Report drill-down navigation", nil},
		{"c04_export_generator", "Generate exports (PDF, Excel)", nil},
	}

	for _, f := range moduleCFlows {
		registry.Register(&Flow{
			Name:        f.name,
			Description: f.description,
			Module:      "C",
			Handler:     f.handler,
		})
	}

	// Register Module D flows (Search Analytics)
	moduleDFlows := []struct {
		name        string
		description string
		handler     FlowHandler
	}{
		{"d01_semantic_search", "Semantic search across data", nil},
		{"d02_trend_analysis", "Analyze data trends", nil},
		{"d03_anomaly_detector", "Detect data anomalies", nil},
		{"d04_autonomous_analyst", "Full analytical workflows", nil},
		{"d05_recommendation_engine", "Generate recommendations", nil},
	}

	for _, f := range moduleDFlows {
		registry.Register(&Flow{
			Name:        f.name,
			Description: f.description,
			Module:      "D",
			Handler:     f.handler,
		})
	}

	// Register Module E flows (i18n)
	moduleEFlows := []struct {
		name        string
		description string
		handler     FlowHandler
	}{
		{"e01_language_detection", "Detect query language", nil},
		{"e02_query_translation", "Translate Arabic to English intent", nil},
		{"e03_localized_formatter", "Format numbers/dates/currency", nil},
		{"e04_rtl_layout_generator", "Generate RTL layout hints", nil},
	}

	for _, f := range moduleEFlows {
		registry.Register(&Flow{
			Name:        f.name,
			Description: f.description,
			Module:      "E",
			Handler:     f.handler,
		})
	}

	return nil
}

// Register registers a flow in the registry.
func (r *FlowRegistry) Register(flow *Flow) error {
	if flow == nil {
		return fmt.Errorf("cannot register nil flow")
	}

	if flow.Name == "" {
		return fmt.Errorf("flow name is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flows[flow.Name]; exists {
		r.logger.Warn("overwriting existing flow",
			slog.String("flow", flow.Name),
		)
	}

	r.flows[flow.Name] = flow
	r.logger.Debug("registered flow",
		slog.String("flow", flow.Name),
		slog.String("module", flow.Module),
	)

	return nil
}

// GetFlow retrieves a flow by name.
func (r *FlowRegistry) GetFlow(name string) (*Flow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flow, exists := r.flows[name]
	if !exists {
		return nil, fmt.Errorf("flow not found: %s", name)
	}

	return flow, nil
}

// ListFlows returns all registered flows.
func (r *FlowRegistry) ListFlows() []*Flow {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flows := make([]*Flow, 0, len(r.flows))
	for _, flow := range r.flows {
		flows = append(flows, flow)
	}

	return flows
}

// ListFlowsByModule returns all flows for a specific module.
func (r *FlowRegistry) ListFlowsByModule(module string) []*Flow {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var flows []*Flow
	for _, flow := range r.flows {
		if flow.Module == module {
			flows = append(flows, flow)
		}
	}

	return flows
}

// ExecuteFlow executes a flow by name with the given input.
func (r *FlowRegistry) ExecuteFlow(ctx context.Context, name string, input interface{}) (interface{}, error) {
	flow, err := r.GetFlow(name)
	if err != nil {
		return nil, err
	}

	if flow.Handler == nil {
		return nil, fmt.Errorf("flow handler not implemented: %s", name)
	}

	return flow.Handler(ctx, input)
}

// GetFlow returns a flow from the Genkit instance.
func (g *GenkitInstance) GetFlow(name string) (*Flow, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	flow, exists := g.flows[name]
	if !exists {
		return nil, fmt.Errorf("flow not found: %s", name)
	}

	return flow, nil
}

// ExecuteFlow executes a flow by name.
func (g *GenkitInstance) ExecuteFlow(ctx context.Context, name string, input interface{}) (interface{}, error) {
	registry := GetRegistry()
	return registry.ExecuteFlow(ctx, name, input)
}

// ConfigureTracing sets up distributed tracing for observability.
func (g *GenkitInstance) ConfigureTracing(endpoint string) error {
	g.config.TracingEnabled = true
	g.config.TracingEndpoint = endpoint

	g.logger.Info("tracing configured",
		slog.String("endpoint", endpoint),
	)

	return nil
}

// Shutdown gracefully shuts down the Genkit instance.
func (g *GenkitInstance) Shutdown() error {
	g.logger.Info("shutting down Genkit")
	return nil
}
