---
name: language-routing
description: Acts as the pre-processing gate for all chat queries. Detects query language, resolves locale, and ensures multi-lingual response generation.
---

# Language Detection & Routing Skill

Guidelines for handling multi-lingual (Arabic/English) user interactions and ensuring locale consistency across MediSync agents.

## Detection Logic

### Hierarchical Detection path
1. **Unicode Heuristic (Fast Path)**:
    - Check for Arabic character ranges (U+0600–U+06FF).
    - If detected, set `detected_language = "ar"`.
    - **Latency**: < 2ms.
2. **LLM Fallback (Accuracy Path)**:
    - If Unicode is ambiguous (only numbers or emojis), use **Gemini 1.5 Flash** for fast classification.
    - Prompt: `Identify language: "en" or "ar". Input: {{ query }}`.

## Locale Propagation

### Genkit FlowContext
- Every downstream agent (Text-to-SQL, Analyst, etc.) must receive the `locale` from the `FlowContext`.
- **Instruction**: "Respond in {{ locale }}. Technical terms/SQL IDs must remain in English/Verbatim."

### Translation Layer (E-02)
- If the input is Arabic but the target agent executes in English (e.g., SQL generation), use the **Translation Agent (E-02)** to convert the intent while preserving identifiers.

## Standard Patterns

### Unicode Detection (Go Engine)
```go
func detectLanguage(text string) string {
    for _, r := range text {
        if unicode.Is(unicode.Arabic, r) {
            return "ar"
        }
    }
    return "en"
}
```

### Prompt for Localized Response
```
You are a MediSync assistant. 
The user's preferred locale is: {{ response_locale }}
Tone: Professional and helpful.
Instruction: Translate the final business insight into {{ response_locale }}, but keep column names and table data exactly as they appear in the database.
```

## Accuracy & Quality

- **No-Block Policy**: If detection fails, default to English. Never block a user query due to language detection errors.
- **Numbers-Only Guard**: A query like "100234" should default to the user's preferred profile locale, not trigger a language error.

## Accessibility Checklist
- [ ] Support easy language switching in the chat UI.
- [ ] Ensure Right-to-Left (RTL) support for Arabic responses in the frontend.
- [ ] Audit language accuracy weekly using native speakers.
