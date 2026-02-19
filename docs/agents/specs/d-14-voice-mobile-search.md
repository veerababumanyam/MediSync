# Agent Specification — D-14: Voice/Mobile Search Agent

**Agent ID:** `D-14`  
**Agent Name:** Voice/Mobile Search Agent  
**Module:** D — Advanced Search Analytics  
**Phase:** 15  
**Priority:** P2 Medium  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Enables voice-driven and mobile-optimised analytics queries via the Flutter mobile app. Transcribes speech, interprets analytical intent, executes the query, and returns a voice-friendly summarised response alongside a mobile chart.

> **Addresses:** PRD §6.9.10, US31 — Voice and mobile-first analytics for field users.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Manual |
| **Manual trigger** | Voice button or typed query in Flutter mobile app |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `audio_bytes` | `[]byte` | Mobile microphone | ⬜ (voice) |
| `text_query` | `string` | Mobile keyboard | ⬜ (text) |
| `user_id` | `string` | JWT | ✅ |
| `session_id` | `string` | Mobile app | ✅ |
| `locale` | `string` | Device locale | ✅ |

> Either `audio_bytes` or `text_query` must be provided.

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `transcription` | `string` | Transcribed text (voice input) |
| `response_text` | `string` | Brief, voice-friendly answer (≤60 words) |
| `tts_audio` | `[]byte` | Text-to-speech response audio |
| `mobile_chart` | `MobileChart` | Simplified chart for mobile screen |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Whisper (OpenAI, MIT) | MIT | Speech-to-text transcription |
| 2 | D-01 Natural Language Search | Internal | Route query |
| 3 | D-03 Conversational | Internal | Context-aware handling |
| 4 | Genkit Flow (`mobile-response-format`) | Apache-2.0 | Compress response to ≤60 words |
| 5 | Coqui TTS (MPL-2.0) or Piper TTS (MIT) | Open | Text-to-speech |

---

## 6. Guardrails

- All data access filtered through C-05 OPA row/column security.
- Voice responses never contain raw financial figures without context.
- Audio input retained for max 5 minutes for debugging; deleted automatically.
- Mobile chart simplified to max 10 data points.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Speech-to-text WER (Word Error Rate) | < 8% |
| Query success rate | ≥ 90% |
| P95 Latency (voice input → response) | < 10s |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service + Whisper Python sidecar |
| **Depends on** | Whisper sidecar, D-01, D-03, TTS sidecar |
| **Consumed by** | Flutter mobile app |
