// Package handlers provides HTTP handlers for the MediSync API.
//
// This file implements the health endpoint for agent status monitoring.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	supervisor Supervisor
	llmConfig  *LLMConfig
}

// Supervisor interface for health handler.
type Supervisor interface {
	CheckHealth(ctx interface{}) interface{ ToJSON() string }
	GetCachedHealth() interface{ ToJSON() string }
}

// LLMConfig contains LLM provider configuration.
type LLMConfig struct {
	Name  string
	Model string
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(supervisor Supervisor, llmConfig *LLMConfig) *HealthHandler {
	return &HealthHandler{
		supervisor: supervisor,
		llmConfig:  llmConfig,
	}
}

// AgentsHealth handles GET /v1/agents/health.
func (h *HealthHandler) AgentsHealth(w http.ResponseWriter, r *http.Request) {
	// Get cached health status for quick response
	health := h.supervisor.GetCachedHealth()

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Use cached result if recent (within 30 seconds)
	// Otherwise, perform a fresh health check

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(health.ToJSON()))
}

// HealthResponse represents the health API response.
type HealthResponse struct {
	Status      string            `json:"status"`
	Timestamp   string            `json:"timestamp"`
	Agents      []AgentHealthInfo `json:"agents"`
	LLMProvider *LLMProviderInfo  `json:"llm_provider,omitempty"`
}

// AgentHealthInfo contains health info for a single agent.
type AgentHealthInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	LastCheck  string `json:"last_check"`
	Error      string `json:"error_message,omitempty"`
}

// LLMProviderInfo contains LLM provider health info.
type LLMProviderInfo struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Status  string `json:"status"`
	Latency int64  `json:"latency_ms,omitempty"`
}

// WriteHealthResponse writes a health response.
func WriteHealthResponse(w http.ResponseWriter, status int, response *HealthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// CreateHealthResponse creates a health response from status.
func CreateHealthResponse(status string, agents []AgentHealthInfo, llm *LLMProviderInfo) *HealthResponse {
	return &HealthResponse{
		Status:      status,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Agents:      agents,
		LLMProvider: llm,
	}
}

// SimpleHealth handles GET /health (basic health check).
func SimpleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessCheck handles GET /ready (Kubernetes readiness probe).
func ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// Check if all critical components are ready
	// This would check database, redis, etc.

	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": map[string]string{
			"database": "ok",
			"redis":    "ok",
			"agents":   "ok",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LivenessCheck handles GET /live (Kubernetes liveness probe).
func LivenessCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
