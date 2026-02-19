---
name: project-references
description: Essential reference navigator for MediSync development. Provides quick access to architectural standards, PRD requirements, agent specifications, and the modular project structure to ensure implementation consistency.
---

# Project References Skill

This skill serves as the central "ground truth" navigator for agents developing the MediSync platform. It ensures all development aligns with the established architecture, security policies, and product requirements.

## Core Project Documents

When tasked with development, analysis, or debugging, start by consulting these primary references:

| Document | Purpose |
|----------|---------|
| `README.md` | High-level overview, tech stack, and project structure. |
| `docs/PRD.md` | Product requirements, user stories, and phased roadmap. |
| `docs/ARCHITECTURE.md` | System layers, data flow diagrams, and architectural principles. |
| `docs/DESIGN.md` | UI/UX standards, color palettes, and Generative UI patterns. |
| `docs/agents/BLUEPRINTS.md` | Detailed specifications for top-priority agents. |
| `docs/AI_RULES.md` | Hard constraints and quality guidelines for AI-generated code and logic. |

## Modular Architecture (A-E)

MediSync is divided into five functional modules. Always identify which module a task belongs to:

- **Module A: Conversational BI** (`internal/agents/module_a/`) - Text-to-SQL, charts, and dashboarding.
- **Module B: AI Accountant** (`internal/agents/module_b/`) - OCR, ledger mapping, Tally sync, and approvals.
- **Module C: Easy Reports** (`internal/agents/module_c/`) - Pre-built MIS reports and automated scheduling.
- **Module D: Advanced Search Analytics** (`internal/agents/module_d/`) - Autonomous analysts and prescriptive AI.
- **Module E: Language & i18n** (`internal/agents/module_e/`) - Locale routing, Arabic RTL, and translation.

## Reference Skill Repository

The project maintains a vast library of "meta-skills" in `/References/skills`. These should be used as templates and best-practice guides for building new MediSync agent skills.

### Key Knowledge Areas from References:
- **Search Strategy**: How to decompose queries (`References/skills/search-strategy`).
- **Knowledge Synthesis**: How to merge multi-source data (`References/skills/knowledge-synthesis`).
- **Source Management**: How to handle MCP and third-party data connections (`References/skills/source-management`).
- **Data Validation**: Patterns for ETL sanity checks (`References/skills/data-validation`).

## Technical Standards

### Backend (Go)
- **Runtime**: Go 1.26.
- **Patterns**: Clean architecture in `internal/`, use of `sqlx` for warehouse access, and `go-chi` for routing.
- **Security**: Mandatory OPA policy checks for every write action.

### AI Orchestration (Genkit)
- **Framework**: Google Genkit (Apache-2.0).
- **Agents**: Every agent must have a Pydantic input/output contract.
- **Protocol**: Coordinate multi-agent tasks using the Google A2A Protocol.

### Frontend (React & Flutter)
- **Web**: React 19 + Vite + CopilotKit for Generative UI.
- **Mobile**: Flutter + PowerSync for offline-first capabilities.
- **i18n**: Support for both LTR (English) and RTL (Arabic) is non-negotiable.

## Security & Governance Protocols

1. **Read-Only Data Plane**: AI agents NEVER write directly to the data warehouse. Use the `medisync_readonly` role.
2. **HITL Gateways**: All write-backs to Tally ERP MUST be human-approved via Module B's approval workflow.
3. **OPA Enforcement**: All authorization logic resides in `.rego` files in the `/policies` directory.
4. **Audit Logging**: Every financial transaction or AI decision must be logged to the immutable `audit_log` table.

## Implementation Checklist

- [ ] Does this change align with the **Phase** defined in `PRD.md`?
- [ ] Have I used the correct **Agent ID** (e.g., A-XX) for new agent logic?
- [ ] Is there a **Pydantic model** for the inputs and outputs?
- [ ] Are **i18n** considerations (strings, RTL layout) addressed?
- [ ] Has **OPA authorization** been checked for write operations?
- [ ] Is **OpenTelemetry/Langfuse** tracing implemented for this flow?
