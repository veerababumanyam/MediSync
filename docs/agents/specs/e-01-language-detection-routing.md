# Agent Specification — E-01: Language Detection & Routing Agent

**Agent ID:** `E-01`  
**Agent Name:** Language Detection & Routing Agent  
**Module:** E — Language & Localisation  
**Phase:** 2 (ships alongside A-01 Text-to-SQL; all chat agents depend on E-01)  
**Priority:** P0 Critical — blocks all bilingual chat functionality  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Acts as the **mandatory pre-processing gate** for all user queries entering the MediSync chat pipeline. Detects the language of the incoming query, resolves the active user locale, normalises cross-lingual intent so downstream SQL agents always receive unambiguous English-intent metadata, and injects a `locale` tag into the Genkit flow context that all downstream agents use to format their responses in the user's preferred language.

> **Addresses:** PRD §6.10.5 (AI multilingual requirements), PRD §5 AI Skill #5, User Stories US35, US36.  
> **Architecture reference:** [docs/i18n-architecture.md §5](../../i18n-architecture.md)

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Reactive (always) |
| **Manual trigger** | Any user query submitted to the BI chat interface |
| **Calling context** | Invoked by the API Gateway before routing to any domain agent (A-01 through D-14) |
| **Position in pipeline** | Step 0 — runs before every domain agent; output is injected into downstream context |

---

## 3. Inputs

| Input | Type | Source | Required | Notes |
|-------|------|--------|:--------:|-------|
| `raw_query` | `string` | Chat interface | ✅ | Original user text, unmodified |
| `user_locale_pref` | `string` | JWT claim `locale` | ✅ | User's stored language preference (`en`, `ar`); fallback: `Accept-Language` header |
| `ai_response_lang` | `string` | JWT claim / user prefs | ✅ | `en`, `ar`, or `auto` |
| `session_id` | `string` | Session store | ✅ | Trace correlation |
| `user_id` | `string` | JWT | ✅ | Audit logging |

```go
type LanguageDetectionInput struct {
    RawQuery        string `json:"raw_query"      validate:"required,max=1000"`
    UserLocalePref  string `json:"user_locale_pref" validate:"required,oneof=en ar"`
    AIResponseLang  string `json:"ai_response_lang" validate:"required,oneof=en ar auto"`
    SessionID       string `json:"session_id"     validate:"required"`
    UserID          string `json:"user_id"        validate:"required"`
}
```

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `detected_language` | `string` | Audit log; downstream context |
| `response_locale` | `string` | Injected into all downstream Genkit flows |
| `normalised_query` | `string` | Passed to domain agents (A-01, D-01, etc.) |
| `original_query` | `string` | Preserved for user display; audit logging |
| `requires_translation` | `bool` | If `true`, E-02 Query Translation Agent is triggered |
| `confidence` | `float64` | Language detection confidence (0.0–1.0) |
| `trace_id` | `string` | OpenTelemetry trace ID |

```go
type LanguageDetectionOutput struct {
    AgentID             string  `json:"agent_id"`           // "E-01"
    DetectedLanguage    string  `json:"detected_language"`   // ISO 639-1: "en" | "ar"
    ResponseLocale      string  `json:"response_locale"`     // Active locale for response
    NormalisedQuery     string  `json:"normalised_query"`    // Query ready for downstream
    OriginalQuery       string  `json:"original_query"`      // Unmodified for audit
    RequiresTranslation bool    `json:"requires_translation"` // Triggers E-02 if true
    Confidence          float64 `json:"confidence"`
    TraceID             string  `json:"trace_id"`
}
```

---

## 5. Processing Logic

### 5.1 Language Detection

```
1. If query Unicode range is predominantly Arabic script (U+0600–U+06FF) → detected = "ar"
2. Else if LLM fast-classification (model: Gemini Flash 1.5, 1-token response) → detected = "en" | "ar"
3. Confidence threshold: if < 0.75, default to user_locale_pref
```

Fast Unicode heuristic (Go):

```go
func detectLanguage(query string) (lang string, confidence float64) {
    arabicRunes := 0
    totalRunes := 0
    for _, r := range query {
        if r >= 0x0600 && r <= 0x06FF {
            arabicRunes++
        }
        totalRunes++
    }
    if totalRunes == 0 { return "en", 1.0 }
    ratio := float64(arabicRunes) / float64(totalRunes)
    if ratio > 0.3 {
        return "ar", min(ratio*1.5, 1.0)
    }
    return "en", 1.0 - ratio
}
```

### 5.2 Response Locale Resolution

```
response_locale is determined by ai_response_lang setting:
  "auto"  → use detected_language (respond in the language the user queried in)
  "en"    → always respond in English regardless of query language
  "ar"    → always respond in Arabic regardless of query language
```

### 5.3 Query Normalisation

- If `detected_language == "ar"` AND `requires_translation == true`:
  - E-01 flags `requires_translation: true`
  - E-02 (Query Translation Agent) is triggered and returns English intent representation
  - The domain agents (A-01, A-07, etc.) receive the English-intent normalised query
  - The original Arabic query is preserved in `original_query` for display and audit

- If `detected_language == "en"` OR the query is mixed (numbers + Arabic phrases):
  - `normalised_query` = `raw_query` (no translation needed)
  - `requires_translation: false`

### 5.4 Locale Injection into Genkit Flows

E-01 output is stored in the Genkit `FlowContext` and is automatically available to all downstream flows:

```go
// Genkit flow context injection (internal/ai/context.go)
ctx = genkit.WithValue(ctx, "response_locale", output.ResponseLocale)
ctx = genkit.WithValue(ctx, "detected_language", output.DetectedLanguage)
ctx = genkit.WithValue(ctx, "original_query", output.OriginalQuery)
```

All downstream flow system prompts read `response_locale` from context and apply the `ResponseLanguageInstruction` template:

```go
const ResponseLanguageInstruction = `
LANGUAGE RULE (MANDATORY):
The user's preferred language is: {{.ResponseLocale}}
All explanations, chart titles, table headers, insight narratives,
error messages, and recommendations MUST be written in {{.ResponseLocale}}.
Numbers, dates, and currency must follow {{.ResponseLocale}} formatting conventions.
SQL identifiers and column names remain in English — do NOT translate them.
`
```

---

## 6. Tool Chain (Genkit Flow)

```
User Query
  → [E-01] Unicode Heuristic (Go — 0ms)
  → [E-01] LLM Fast-Classify (Gemini Flash 1.5 — only if confidence < 0.75)
  → [E-01] Resolve response_locale (from ai_response_lang pref)
  → [E-01] Flag requires_translation
  → if requires_translation:
      → [E-02] Arabic Query → English Intent
  → Inject locale into FlowContext
  → [Domain Agent: A-01 / A-07 / D-01 / etc.]
      (reads locale from context; applies ResponseLanguageInstruction)
  → [E-03] Localised Response Formatter
      (number/date/currency formatting for response_locale)
  → Chat Response (in user's language)
```

---

## 7. OSS Components

- **Go** — Unicode range detection (built-in `unicode` package)
- **`golang.org/x/text/language`** — `language.BestMatch` for Accept-Language resolution
- **Genkit (Firebase)** — Flow context propagation
- **Gemini Flash 1.5** — Low-cost language classification for ambiguous queries
- **OpenTelemetry** — Trace ID propagation across language + domain agent hops

---

## 8. System Prompt (Fast Classification)

```
You are a language detection classifier. Your only job is to classify the language of the input text.
Output EXACTLY ONE token: "en" or "ar". No explanation, no punctuation.
Input: {{.raw_query}}
```

---

## 9. Guardrails

| Guardrail | Mechanism |
|-----------|-----------|
| **Supported locale enforcement** | If `detected_language` is neither `en` nor `ar`, default to `user_locale_pref` and log a `UNSUPPORTED_LOCALE` warning |
| **SQL column/table name preservation** | `normalised_query` MUST NOT translate database identifiers; E-02 is instructed to preserve technical terms verbatim |
| **Latency budget** | E-01 Unicode heuristic must complete in < 2ms; LLM fallback path must complete in < 500ms to keep total E-01+E-02 overhead under 1 second |
| **Audit preservation** | `original_query` in source language always stored in audit log; `normalised_query` stored separately for SQL traceability |
| **Fallback on failure** | If E-01 fails for any reason (timeout, LLM error), set `response_locale = user_locale_pref` and pass `raw_query` as `normalised_query` unmodified — never block the user query |

---

## 10. HITL Gate

**E-01 does not require HITL.** It is a fast, transparent pre-processing step. Failures are silent with safe fallback.

For HITL purposes:
- Low-confidence detection (< 0.75) is logged for monitoring dashboards but does not block the request.
- Incorrect language detection by users can always be corrected via the language preference toggle in the UI.

---

## 11. Evaluation Criteria

| Metric | Target | Measurement |
|--------|--------|-------------|
| Language detection accuracy | > 98% | Eval dataset of 500 Arabic + 500 English + 100 mixed queries |
| False Arabic detection rate (numbers-only query) | < 1% | Subset: queries with only numerals/dates |
| E-01 end-to-end latency (Unicode path) | < 5ms | p95 |
| E-01 end-to-end latency (LLM fallback path) | < 600ms | p95 |
| Response locale correctness (`auto` mode) | > 99% | `detected_language == response_locale` for `auto` users |
| Downstream response language accuracy | > 97% | Spot-check: LLM responses are actually in `response_locale` |

---

## 12. Eval Dataset Format

```jsonl
{"raw_query": "أعطني إيرادات هذا الشهر", "expected_lang": "ar", "user_locale": "ar", "expected_response_locale": "ar"}
{"raw_query": "Show me last month revenue", "expected_lang": "en", "user_locale": "ar", "ai_response_lang": "ar", "expected_response_locale": "ar"}
{"raw_query": "Show me sales for Q1", "expected_lang": "en", "user_locale": "en", "expected_response_locale": "en"}
{"raw_query": "أعطني top 10 drugs هذا الأسبوع", "expected_lang": "ar", "note": "mixed script query"}
{"raw_query": "12345", "expected_lang": "en", "note": "numbers only — must not detect as Arabic"}
```

---

## 13. Related Agents

| Agent | Relationship |
|-------|-------------|
| **E-02 Query Translation** | Triggered by E-01 when `requires_translation: true`; translates Arabic query to English intent |
| **E-03 Response Formatter** | Runs after every domain agent; reads `response_locale` injected by E-01 to format numbers/dates/currency |
| **A-01 Text-to-SQL** | Primary consumer of E-01 output; receives `normalised_query` and `response_locale` |
| **A-04 Domain Terminology Normaliser** | Runs in parallel with E-01 in the pre-processing chain; E-01 output feeds its Arabic→English synonym mapping |
| **D-01 Natural Language Search** | Reads `response_locale` from context to format search result snippets |

---

*Spec Owner: AI Platform Team | Created: February 19, 2026 | Status: Draft — Pending Phase 2 Sprint Planning*
