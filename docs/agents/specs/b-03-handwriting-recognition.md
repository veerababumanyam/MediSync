# Agent Specification — B-03: Handwriting Recognition Agent

**Agent ID:** `B-03`  
**Agent Name:** Handwriting Recognition Agent  
**Module:** B — AI Accountant  
**Phase:** 4  
**Priority:** P2 Medium  
**HITL Required:** Yes — always (handwritten accuracy lower than print)  
**Status:** Draft

---

## 1. Purpose

Specialized sub-agent for handwritten invoices and receipts. Applies handwriting-specific OCR (HTR model via PaddleOCR) plus LLM post-processing to extract structured fields, then passes results to B-02's confidence scorer. Given inherent accuracy limitations of handwriting, HITL review is always triggered.

> **Addresses:** PRD §6.7.1, US13 — Handwritten document support.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Upstream-agent-output |
| **Calling agent** | B-02 (when handwriting detected in layout analysis) |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `file` | `bytes` | B-02 | ✅ |
| `layout_segments` | `[]Segment` | PaddleOCR layout from B-02 | ✅ |
| `document_type` | `enum` | B-01 | ✅ |

---

## 4. Outputs

Same shape as B-02 `ExtractionResult`, passed back into B-02 pipeline.

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | PaddleOCR HTR model | Apache-2.0 | Handwriting transcription |
| 2 | Image pre-processor (Go/OpenCV) | MIT | Deskew, binarize, denoise |
| 3 | Genkit Flow (`htr-post-process`) | Apache-2.0 | LLM cleanup of HTR output |

---

## 6. Guardrails

- `hitl_required` is **always true** for handwritten documents regardless of confidence score.
- Images pre-processed (deskewed, binarized) before HTR to improve accuracy.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| Character error rate (CER) | < 5% |
| Field extraction accuracy | ≥ 85% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Python sidecar (PaddleOCR HTR) + Go service |
| **Depends on** | B-02, PaddleOCR HTR model |
| **Consumed by** | B-02 (merges result) |
