// Package a03_visualization provides the A-03 Visualization Routing Agent.
//
// This file defines the LLM client interface for chart type routing.
// The interface allows for different LLM backends (Ollama, vLLM, etc.)
// to be used interchangeably.
package a03_visualization

import (
	"context"
)

// LLMClient defines the interface for LLM interactions.
// This interface abstracts the underlying LLM provider (Ollama, vLLM, etc.)
type LLMClient interface {
	// Generate generates a text completion from the given prompt.
	Generate(ctx context.Context, prompt string) (string, error)

	// GenerateWithJSON generates a completion and expects JSON output.
	GenerateWithJSON(ctx context.Context, prompt string) ([]byte, error)

	// IsAvailable returns true if the LLM service is healthy.
	IsAvailable(ctx context.Context) bool
}

// MockLLMClient is a mock implementation for testing purposes.
type MockLLMClient struct {
	Response string
	JSONData []byte
	Healthy  bool
}

// Generate returns a mock response.
func (m *MockLLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	if m.Response != "" {
		return m.Response, nil
	}
	return `{"chart_type": "kpiCard", "confidence": 95.0, "reasoning": "Single value query"}`, nil
}

// GenerateWithJSON returns mock JSON data.
func (m *MockLLMClient) GenerateWithJSON(ctx context.Context, prompt string) ([]byte, error) {
	if m.JSONData != nil {
		return m.JSONData, nil
	}
	return []byte(`{"chart_type": "kpiCard", "confidence": 95.0, "reasoning": "Single value query"}`), nil
}

// IsAvailable returns the mock health status.
func (m *MockLLMClient) IsAvailable(ctx context.Context) bool {
	return m.Healthy
}
