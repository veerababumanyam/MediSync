// Package a05_hallucination provides the hallucination guard agent.
//
// This agent detects and rejects off-topic queries that are outside
// the healthcare and financial analytics domain.
package a05_hallucination

import (
	"context"
	"encoding/json"
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

// AgentID is the unique identifier for this agent.
const AgentID = "a-05-hallucination"

// Agent implements the hallucination guard agent.
type Agent struct {
	id          string
	logger      *slog.Logger
	categories  *OffTopicCategories
	threshold   float64
	cache       map[string]bool
	cacheMu     sync.RWMutex
}

// AgentConfig holds configuration for the agent.
type AgentConfig struct {
	Logger    *slog.Logger
	Threshold float64
}

// New creates a new hallucination guard agent.
func New(cfg AgentConfig) *Agent {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.Threshold == 0 {
		cfg.Threshold = 0.7
	}

	return &Agent{
		id:         AgentID,
		logger:     cfg.Logger.With("agent", AgentID),
		categories: NewOffTopicCategories(),
		threshold:  cfg.Threshold,
		cache:      make(map[string]bool),
	}
}

// GuardRequest contains the request for topic validation.
type GuardRequest struct {
	Query       string            `json:"query"`
	Locale      string            `json:"locale"`
	UserRoles   []string          `json:"user_roles,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// GuardResponse contains the result of topic validation.
type GuardResponse struct {
	IsOnTopic        bool     `json:"is_on_topic"`
	Category         string   `json:"category"`
	Confidence       float64  `json:"confidence"`
	RejectionMessage string   `json:"rejection_message,omitempty"`
	NeedsClarification bool   `json:"needs_clarification"`
	ClarificationMessage string `json:"clarification_message,omitempty"`
	Suggestions      []string `json:"suggestions,omitempty"`
}

// AgentCard returns the ADK agent card for discovery.
func (a *Agent) AgentCard() map[string]interface{} {
	return map[string]interface{}{
		"id":          AgentID,
		"name":        "Hallucination Guard Agent",
		"description": "Detects and rejects off-topic queries outside healthcare and financial analytics",
		"capabilities": []string{
			"topic-classification",
			"off-topic-detection",
			"query-validation",
		},
		"version": "1.0.0",
	}
}

// Guard validates whether a query is on-topic.
func (a *Agent) Guard(ctx context.Context, req GuardRequest) (*GuardResponse, error) {
	a.logger.Debug("guarding query", "query_length", len(req.Query))

	// Check cache
	cacheKey := req.Query
	if cached, ok := a.getFromCache(cacheKey); ok {
		a.logger.Debug("cache hit for query")
		return &GuardResponse{
			IsOnTopic: cached,
		}, nil
	}

	// Classify the query
	response := a.classify(req.Query, req.Locale)

	// Add to cache
	a.addToCache(cacheKey, response.IsOnTopic)

	a.logger.Info("query classified",
		"is_on_topic", response.IsOnTopic,
		"category", response.Category,
		"confidence", response.Confidence)

	return response, nil
}

// classify determines if a query is on-topic.
func (a *Agent) classify(query string, locale string) *GuardResponse {
	queryLower := strings.ToLower(query)
	response := &GuardResponse{
		IsOnTopic: true,
		Confidence: 1.0,
	}

	// Check for off-topic categories
	for _, category := range a.categories.GetOffTopicCategories() {
		if a.matchesCategory(queryLower, category) {
			response.IsOnTopic = false
			response.Category = category.Name
			response.Confidence = category.Confidence
			response.RejectionMessage = a.getRejectionMessage(locale)

			a.logger.Warn("off-topic query detected",
				"category", category.Name,
				"confidence", category.Confidence)

			return response
		}
	}

	// Check for on-topic categories
	maxOnTopicScore := 0.0
	matchedCategory := ""

	for _, category := range a.categories.GetOnTopicCategories() {
		if score := a.scoreCategory(queryLower, category); score > maxOnTopicScore {
			maxOnTopicScore = score
			matchedCategory = category.Name
		}
	}

	if maxOnTopicScore < a.threshold {
		// Low confidence in on-topic classification
		response.Confidence = maxOnTopicScore
		response.Category = "ambiguous"
		response.NeedsClarification = true
		response.ClarificationMessage = a.getClarificationMessage(locale)
		response.Suggestions = a.getSuggestions(queryLower)
	} else {
		response.Confidence = maxOnTopicScore
		response.Category = matchedCategory
	}

	return response
}

// matchesCategory checks if a query matches an off-topic category.
func (a *Agent) matchesCategory(query string, category Category) bool {
	for _, pattern := range category.Patterns {
		matched, err := regexp.MatchString(pattern, query)
		if err == nil && matched {
			return true
		}
	}

	for _, keyword := range category.Keywords {
		if strings.Contains(query, keyword) {
			return true
		}
	}

	return false
}

// scoreCategory calculates how well a query matches an on-topic category.
func (a *Agent) scoreCategory(query string, category Category) float64 {
	score := 0.0
	totalKeywords := len(category.Keywords)

	if totalKeywords == 0 {
		return 0
	}

	matches := 0
	for _, keyword := range category.Keywords {
		if strings.Contains(query, keyword) {
			matches++
		}
	}

	score = float64(matches) / float64(totalKeywords)

	// Bonus for pattern matches
	for _, pattern := range category.Patterns {
		matched, _ := regexp.MatchString(pattern, query)
		if matched {
			score += 0.1
		}
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// getRejectionMessage returns a localized rejection message.
func (a *Agent) getRejectionMessage(locale string) string {
	messages := map[string]string{
		"en": "I can only help with healthcare and financial data queries. Please ask about patient visits, revenue, appointments, or other business metrics.",
		"ar": "يمكنني فقط المساعدة في استعلامات البيانات الصحية والمالية. يرجى السؤال عن زيارات المرضى أو الإيرادات أو المواعيد أو مقاييس الأعمال الأخرى.",
	}

	if msg, ok := messages[locale]; ok {
		return msg
	}
	return messages["en"]
}

// getClarificationMessage returns a localized clarification request.
func (a *Agent) getClarificationMessage(locale string) string {
	messages := map[string]string{
		"en": "Could you please specify what data you'd like to see? For example: patient visits, revenue, appointments, or inventory?",
		"ar": "هل يمكنك تحديد البيانات التي تريد رؤيتها؟ على سبيل المثال: زيارات المرضى أو الإيرادات أو المواعيد أو المخزون؟",
	}

	if msg, ok := messages[locale]; ok {
		return msg
	}
	return messages["en"]
}

// getSuggestions returns query suggestions based on partial matches.
func (a *Agent) getSuggestions(query string) []string {
	suggestions := []string{}

	// Check for partial keyword matches
	keywords := []string{
		"revenue", "patients", "appointments", "visits",
		"billing", "pharmacy", "clinic", "doctors",
	}

	for _, kw := range keywords {
		if strings.Contains(query, kw) || len(suggestions) < 3 {
			suggestions = append(suggestions, kw)
		}
	}

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

// getFromCache retrieves a cached result.
func (a *Agent) getFromCache(key string) (bool, bool) {
	a.cacheMu.RLock()
	defer a.cacheMu.RUnlock()
	result, ok := a.cache[key]
	return result, ok
}

// addToCache adds a result to the cache.
func (a *Agent) addToCache(key string, result bool) {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()

	// Simple cache eviction
	if len(a.cache) >= 1000 {
		for k := range a.cache {
			delete(a.cache, k)
			break
		}
	}

	a.cache[key] = result
}

// ToJSON serializes the guard response.
func (r *GuardResponse) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}
