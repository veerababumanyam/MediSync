// Package shared provides shared utilities for AI agents in MediSync.
//
// This file implements the AgentContext which carries user context through
// all agent operations, including locale-aware prompting helpers.
package shared

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// contextKey is a type for context keys.
type contextKey string

const (
	// AgentContextKey is the context key for AgentContext.
	AgentContextKey contextKey = "agent_context"
)

// AgentContext carries user context through agent operations.
type AgentContext struct {
	// UserID is the unique identifier for the user.
	UserID string `json:"user_id"`

	// TenantID is the tenant/organization identifier.
	TenantID string `json:"tenant_id"`

	// Locale is the user's preferred language (en, ar).
	Locale string `json:"locale"`

	// SessionID is the current session identifier.
	SessionID string `json:"session_id"`

	// Roles are the user's assigned roles.
	Roles []string `json:"roles"`

	// CostCentres are the user's accessible cost centres.
	CostCentres []string `json:"cost_centres"`

	// Timezone is the user's timezone.
	Timezone string `json:"timezone"`

	// CalendarSystem is the user's calendar system (gregorian, hijri).
	CalendarSystem string `json:"calendar_system"`

	// RequestID is the unique request identifier for tracing.
	RequestID string `json:"request_id"`

	// TraceID is the distributed tracing identifier.
	TraceID string `json:"trace_id"`
}

// NewAgentContext creates a new AgentContext from a standard context.Context.
// It extracts values set by the API middleware.
func NewAgentContext(ctx context.Context) *AgentContext {
	agentCtx := &AgentContext{
		Locale:         "en",
		Roles:          []string{},
		CostCentres:    []string{},
		SessionID:      uuid.New().String(),
		CalendarSystem: "gregorian",
	}

	// Extract values from context (set by middleware)
	if userID := ctx.Value(contextKey("user_id")); userID != nil {
		agentCtx.UserID = userID.(string)
	}

	if tenantID := ctx.Value(contextKey("tenant_id")); tenantID != nil {
		agentCtx.TenantID = tenantID.(string)
	}

	if locale := ctx.Value(contextKey("locale")); locale != nil {
		agentCtx.Locale = locale.(string)
	}

	if roles := ctx.Value(contextKey("roles")); roles != nil {
		agentCtx.Roles = roles.([]string)
	}

	if requestID := ctx.Value(contextKey("request_id")); requestID != nil {
		agentCtx.RequestID = requestID.(string)
	}

	return agentCtx
}

// NewAgentContextWithValues creates a new AgentContext with explicit values.
func NewAgentContextWithValues(userID, tenantID, locale string, roles []string) *AgentContext {
	return &AgentContext{
		UserID:         userID,
		TenantID:       tenantID,
		Locale:         locale,
		Roles:          roles,
		SessionID:      uuid.New().String(),
		CalendarSystem: "gregorian",
	}
}

// WithContext returns a new context.Context with the AgentContext embedded.
func (a *AgentContext) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, AgentContextKey, a)
}

// FromContext retrieves the AgentContext from a context.Context.
func FromContext(ctx context.Context) *AgentContext {
	if agentCtx, ok := ctx.Value(AgentContextKey).(*AgentContext); ok {
		return agentCtx
	}
	// Return a default context if none found
	return NewAgentContext(ctx)
}

// ResponseLanguageInstruction returns the language instruction for LLM prompts.
// This should be appended to every agent prompt to ensure localized responses.
func (a *AgentContext) ResponseLanguageInstruction() string {
	switch a.Locale {
	case "ar":
		return `ResponseLanguageInstruction: Respond in Arabic (العربية). Format all numbers using Arabic numerals or Eastern Arabic numerals as appropriate. Format dates in the Arabic style (e.g., ١٩ فبراير ٢٠٢٦). Use right-to-left (RTL) text direction. Maintain formal Arabic business language.`
	default:
		return `ResponseLanguageInstruction: Respond in English. Format numbers using standard Western numerals. Format dates in ISO 8601 or locale-friendly format (e.g., February 19, 2026). Use professional business English.`
	}
}

// FormatInstruction returns formatting instructions based on locale.
func (a *AgentContext) FormatInstruction() string {
	switch a.Locale {
	case "ar":
		return `Format: Use Arabic numerals (٠١٢٣٤٥٦٧٨٩) for numbers. Format currency as SAR with Arabic text. Use Hijri calendar references if relevant.`
	default:
		return `Format: Use Western numerals (0-9) for numbers. Format currency as SAR or AED with English text. Use Gregorian calendar.`
	}
}

// Dateformat returns the date format string for the locale.
func (a *AgentContext) DateFormat() string {
	switch a.Locale {
	case "ar":
		return "02/01/2006" // DD/MM/YYYY (Arabic style)
	default:
		return "2006-01-02" // ISO 8601
	}
}

// NumberFormat returns number formatting instructions.
func (a *AgentContext) NumberFormat() string {
	switch a.Locale {
	case "ar":
		return "Use Arabic numerals and Arabic thousands separators."
	default:
		return "Use Western numerals with comma thousands separators."
	}
}

// CurrencyFormat returns currency formatting instructions.
func (a *AgentContext) CurrencyFormat() string {
	switch a.Locale {
	case "ar":
		return "ر.س" // Saudi Riyal in Arabic
	default:
		return "SAR" // Saudi Riyal in English
	}
}

// IsRTL returns true if the locale is a right-to-left language.
func (a *AgentContext) IsRTL() bool {
	return a.Locale == "ar"
}

// HasRole checks if the user has a specific role.
func (a *AgentContext) HasRole(role string) bool {
	for _, r := range a.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin returns true if the user has admin role.
func (a *AgentContext) IsAdmin() bool {
	return a.HasRole("admin")
}

// HasFinanceRole returns true if the user has finance-related roles.
func (a *AgentContext) HasFinanceRole() bool {
	financeRoles := []string{"finance", "accountant", "finance_head"}
	for _, role := range financeRoles {
		if a.HasRole(role) {
			return true
		}
	}
	return false
}

// HasMedicalRole returns true if the user has medical-related roles.
func (a *AgentContext) HasMedicalRole() bool {
	medicalRoles := []string{"doctor", "nurse", "medical_staff", "pharmacist"}
	for _, role := range medicalRoles {
		if a.HasRole(role) {
			return true
		}
	}
	return false
}

// CanAccessCostCentres checks if user can access specified cost centres.
func (a *AgentContext) CanAccessCostCentres(centres []string) bool {
	if a.IsAdmin() {
		return true
	}

	centreMap := make(map[string]bool)
	for _, cc := range a.CostCentres {
		centreMap[cc] = true
	}

	for _, centre := range centres {
		if !centreMap[centre] {
			return false
		}
	}

	return true
}

// SystemPromptAdditions returns locale-specific additions to system prompts.
func (a *AgentContext) SystemPromptAdditions() string {
	return fmt.Sprintf(`
=== USER CONTEXT ===
User ID: %s
Tenant ID: %s
Locale: %s
Timezone: %s
Calendar System: %s
Roles: %v

=== RESPONSE REQUIREMENTS ===
%s
%s
====================
`,
		a.UserID,
		a.TenantID,
		a.Locale,
		a.Timezone,
		a.CalendarSystem,
		a.Roles,
		a.ResponseLanguageInstruction(),
		a.FormatInstruction(),
	)
}

// PromptWithLocale wraps a prompt with locale instructions.
func (a *AgentContext) PromptWithLocale(prompt string) string {
	return fmt.Sprintf("%s\n\n%s", a.ResponseLanguageInstruction(), prompt)
}

// FullPromptWithSystem generates a complete prompt with system context.
func (a *AgentContext) FullPromptWithSystem(systemPrompt, userPrompt string) string {
	return fmt.Sprintf(`%s

%s

User Request: %s`,
		systemPrompt,
		a.SystemPromptAdditions(),
		userPrompt,
	)
}

// Clone creates a copy of the AgentContext.
func (a *AgentContext) Clone() *AgentContext {
	return &AgentContext{
		UserID:         a.UserID,
		TenantID:       a.TenantID,
		Locale:         a.Locale,
		SessionID:      a.SessionID,
		Roles:          append([]string{}, a.Roles...),
		CostCentres:    append([]string{}, a.CostCentres...),
		Timezone:       a.Timezone,
		CalendarSystem: a.CalendarSystem,
		RequestID:      a.RequestID,
		TraceID:        a.TraceID,
	}
}

// Validate validates the AgentContext has required fields.
func (a *AgentContext) Validate() error {
	if a.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if a.Locale != "en" && a.Locale != "ar" {
		return fmt.Errorf("invalid locale: %s (must be 'en' or 'ar')", a.Locale)
	}

	return nil
}

// ============================================================================
// Locale-specific formatting helpers
// ============================================================================

// LocalizedYes returns "yes" in the user's locale.
func (a *AgentContext) LocalizedYes() string {
	if a.Locale == "ar" {
		return "نعم"
	}
	return "Yes"
}

// LocalizedNo returns "no" in the user's locale.
func (a *AgentContext) LocalizedNo() string {
	if a.Locale == "ar" {
		return "لا"
	}
	return "No"
}

// LocalizedError returns an error prefix in the user's locale.
func (a *AgentContext) LocalizedError() string {
	if a.Locale == "ar" {
		return "خطأ"
	}
	return "Error"
}

// LocalizedSuccess returns a success message in the user's locale.
func (a *AgentContext) LocalizedSuccess() string {
	if a.Locale == "ar" {
		return "تم بنجاح"
	}
	return "Success"
}

// LocalizedPleaseWait returns a waiting message in the user's locale.
func (a *AgentContext) LocalizedPleaseWait() string {
	if a.Locale == "ar" {
		return "يرجى الانتظار..."
	}
	return "Please wait..."
}
