// Package shared provides shared utilities for AI agents in MediSync.
//
// This file implements the LLM plugin configuration with support for
// OpenAI and Ollama as switchable providers, configured via LLM_PROVIDER env var.
package shared

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// ProviderType represents the type of LLM provider.
type ProviderType string

const (
	// ProviderOpenAI uses OpenAI's API.
	ProviderOpenAI ProviderType = "openai"
	// ProviderOllama uses local Ollama server.
	ProviderOllama ProviderType = "ollama"
	// ProviderGemini uses Google Gemini API.
	ProviderGemini ProviderType = "gemini"
)

// LLMConfig holds configuration for the LLM provider.
type LLMConfig struct {
	// Provider is the LLM provider type.
	Provider ProviderType

	// Model is the model identifier to use.
	Model string

	// APIKey is the API key for cloud providers.
	APIKey string

	// BaseURL is the base URL for the provider (for Ollama or custom endpoints).
	BaseURL string

	// Temperature controls randomness in output (0.0-2.0).
	Temperature float64

	// MaxTokens is the maximum tokens to generate.
	MaxTokens int

	// TopP controls diversity via nucleus sampling.
	TopP float64

	// FrequencyPenalty reduces repetition of tokens.
	FrequencyPenalty float64

	// PresencePenalty encourages talking about new topics.
	PresencePenalty float64

	// Timeout for API requests.
	TimeoutSeconds int
}

// Default configurations for each provider.
var defaultConfigs = map[ProviderType]*LLMConfig{
	ProviderOpenAI: {
		Provider:         ProviderOpenAI,
		Model:            "gpt-4-turbo-preview",
		Temperature:      0.7,
		MaxTokens:        4096,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		TimeoutSeconds:   60,
	},
	ProviderOllama: {
		Provider:         ProviderOllama,
		Model:            "llama3.1:8b",
		BaseURL:          "http://localhost:11434",
		Temperature:      0.7,
		MaxTokens:        4096,
		TopP:             0.9,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		TimeoutSeconds:   120,
	},
	ProviderGemini: {
		Provider:         ProviderGemini,
		Model:            "gemini-pro",
		Temperature:      0.7,
		MaxTokens:        8192,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		TimeoutSeconds:   60,
	},
}

// LLMClient provides a unified interface for LLM operations.
type LLMClient struct {
	config *LLMConfig
	logger *slog.Logger

	// Provider-specific clients
	openai *OpenAIClient
	ollama *OllamaClient
	gemini *GeminiClient
}

// NewLLMClient creates a new LLM client based on environment configuration.
func NewLLMClient(logger *slog.Logger) (*LLMClient, error) {
	if logger == nil {
		logger = slog.Default()
	}

	// Determine provider from environment
	providerStr := os.Getenv("LLM_PROVIDER")
	if providerStr == "" {
		providerStr = "openai"
	}

	provider := ProviderType(providerStr)
	config, exists := defaultConfigs[provider]
	if !exists {
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}

	// Override with environment variables
	config.APIKey = os.Getenv("LLM_API_KEY")
	if model := os.Getenv("LLM_MODEL"); model != "" {
		config.Model = model
	}
	if baseURL := os.Getenv("LLM_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	client := &LLMClient{
		config: config,
		logger: logger,
	}

	// Initialize provider-specific client
	switch provider {
	case ProviderOpenAI:
		if config.APIKey == "" {
			return nil, fmt.Errorf("LLM_API_KEY is required for OpenAI provider")
		}
		client.openai = NewOpenAIClient(config)
	case ProviderOllama:
		if config.BaseURL == "" {
			config.BaseURL = "http://localhost:11434"
		}
		client.ollama = NewOllamaClient(config)
	case ProviderGemini:
		if config.APIKey == "" {
			return nil, fmt.Errorf("LLM_API_KEY is required for Gemini provider")
		}
		client.gemini = NewGeminiClient(config)
	}

	logger.Info("LLM client initialized",
		slog.String("provider", string(provider)),
		slog.String("model", config.Model),
		slog.Float64("temperature", config.Temperature),
		slog.Int("max_tokens", config.MaxTokens),
	)

	return client, nil
}

// NewLLMClientWithConfig creates a new LLM client with explicit configuration.
func NewLLMClientWithConfig(cfg *LLMConfig, logger *slog.Logger) (*LLMClient, error) {
	if logger == nil {
		logger = slog.Default()
	}

	if cfg == nil {
		return nil, fmt.Errorf("LLM config is required")
	}

	client := &LLMClient{
		config: cfg,
		logger: logger,
	}

	// Initialize provider-specific client
	switch cfg.Provider {
	case ProviderOpenAI:
		client.openai = NewOpenAIClient(cfg)
	case ProviderOllama:
		client.ollama = NewOllamaClient(cfg)
	case ProviderGemini:
		client.gemini = NewGeminiClient(cfg)
	}

	return client, nil
}

// Generate generates a completion from the LLM.
func (c *LLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	switch c.config.Provider {
	case ProviderOpenAI:
		return c.openai.Generate(ctx, prompt)
	case ProviderOllama:
		return c.ollama.Generate(ctx, prompt)
	case ProviderGemini:
		return c.gemini.Generate(ctx, prompt)
	default:
		return "", fmt.Errorf("unsupported provider: %s", c.config.Provider)
	}
}

// GenerateWithSystem generates a completion with a system prompt.
func (c *LLMClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	switch c.config.Provider {
	case ProviderOpenAI:
		return c.openai.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	case ProviderOllama:
		return c.ollama.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	case ProviderGemini:
		return c.gemini.GenerateWithSystem(ctx, systemPrompt, userPrompt)
	default:
		return "", fmt.Errorf("unsupported provider: %s", c.config.Provider)
	}
}

// GenerateJSON generates a JSON completion from the LLM.
func (c *LLMClient) GenerateJSON(ctx context.Context, prompt string, schema interface{}) (interface{}, error) {
	switch c.config.Provider {
	case ProviderOpenAI:
		return c.openai.GenerateJSON(ctx, prompt, schema)
	case ProviderOllama:
		return c.ollama.GenerateJSON(ctx, prompt, schema)
	case ProviderGemini:
		return c.gemini.GenerateJSON(ctx, prompt, schema)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", c.config.Provider)
	}
}

// Embed generates embeddings for the given text.
func (c *LLMClient) Embed(ctx context.Context, text string) ([]float64, error) {
	switch c.config.Provider {
	case ProviderOpenAI:
		return c.openai.Embed(ctx, text)
	case ProviderOllama:
		return c.ollama.Embed(ctx, text)
	case ProviderGemini:
		return c.gemini.Embed(ctx, text)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", c.config.Provider)
	}
}

// GetConfig returns the current LLM configuration.
func (c *LLMClient) GetConfig() *LLMConfig {
	return c.config
}

// GetProvider returns the current provider type.
func (c *LLMClient) GetProvider() ProviderType {
	return c.config.Provider
}

// ============================================================================
// Provider-specific client stubs (to be implemented with actual API calls)
// ============================================================================

// OpenAIClient implements the OpenAI API client.
type OpenAIClient struct {
	config *LLMConfig
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(config *LLMConfig) *OpenAIClient {
	return &OpenAIClient{config: config}
}

// Generate generates a completion.
func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement actual OpenAI API call
	return "", fmt.Errorf("OpenAI client not yet implemented")
}

// GenerateWithSystem generates a completion with system prompt.
func (c *OpenAIClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// TODO: Implement actual OpenAI API call
	return "", fmt.Errorf("OpenAI client not yet implemented")
}

// GenerateJSON generates a JSON completion.
func (c *OpenAIClient) GenerateJSON(ctx context.Context, prompt string, schema interface{}) (interface{}, error) {
	// TODO: Implement actual OpenAI API call
	return nil, fmt.Errorf("OpenAI client not yet implemented")
}

// Embed generates embeddings.
func (c *OpenAIClient) Embed(ctx context.Context, text string) ([]float64, error) {
	// TODO: Implement actual OpenAI API call
	return nil, fmt.Errorf("OpenAI client not yet implemented")
}

// OllamaClient implements the Ollama API client.
type OllamaClient struct {
	config *LLMConfig
}

// NewOllamaClient creates a new Ollama client.
func NewOllamaClient(config *LLMConfig) *OllamaClient {
	return &OllamaClient{config: config}
}

// Generate generates a completion.
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement actual Ollama API call
	return "", fmt.Errorf("Ollama client not yet implemented")
}

// GenerateWithSystem generates a completion with system prompt.
func (c *OllamaClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// TODO: Implement actual Ollama API call
	return "", fmt.Errorf("Ollama client not yet implemented")
}

// GenerateJSON generates a JSON completion.
func (c *OllamaClient) GenerateJSON(ctx context.Context, prompt string, schema interface{}) (interface{}, error) {
	// TODO: Implement actual Ollama API call
	return nil, fmt.Errorf("Ollama client not yet implemented")
}

// Embed generates embeddings.
func (c *OllamaClient) Embed(ctx context.Context, text string) ([]float64, error) {
	// TODO: Implement actual Ollama API call
	return nil, fmt.Errorf("Ollama client not yet implemented")
}

// GeminiClient implements the Google Gemini API client.
type GeminiClient struct {
	config *LLMConfig
}

// NewGeminiClient creates a new Gemini client.
func NewGeminiClient(config *LLMConfig) *GeminiClient {
	return &GeminiClient{config: config}
}

// Generate generates a completion.
func (c *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement actual Gemini API call
	return "", fmt.Errorf("Gemini client not yet implemented")
}

// GenerateWithSystem generates a completion with system prompt.
func (c *GeminiClient) GenerateWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// TODO: Implement actual Gemini API call
	return "", fmt.Errorf("Gemini client not yet implemented")
}

// GenerateJSON generates a JSON completion.
func (c *GeminiClient) GenerateJSON(ctx context.Context, prompt string, schema interface{}) (interface{}, error) {
	// TODO: Implement actual Gemini API call
	return nil, fmt.Errorf("Gemini client not yet implemented")
}

// Embed generates embeddings.
func (c *GeminiClient) Embed(ctx context.Context, text string) ([]float64, error) {
	// TODO: Implement actual Gemini API call
	return nil, fmt.Errorf("Gemini client not yet implemented")
}

// ============================================================================
// Helper functions
// ============================================================================

// GetDefaultConfig returns the default configuration for a provider.
func GetDefaultConfig(provider ProviderType) (*LLMConfig, error) {
	config, exists := defaultConfigs[provider]
	if !exists {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
	// Return a copy to avoid mutations
	copy := *config
	return &copy, nil
}

// ValidateConfig validates an LLM configuration.
func ValidateConfig(config *LLMConfig) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if config.Model == "" {
		return fmt.Errorf("model is required")
	}

	if config.Temperature < 0 || config.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if config.MaxTokens < 1 {
		return fmt.Errorf("max_tokens must be positive")
	}

	return nil
}
