# MediSync Development Standards

## Language Specifics

### Go (Backend)
- All shared logic goes into `internal/`.
- Use `go-chi` for HTTP routing.
- Context propagation is mandatory for tracing.
- Error messages must be localized for the user.

### React (Web)
- Use functional components and hooks.
- Styling: Tailwind CSS with RTL logical properties.
- Components must be "Generative UI" ready (CopilotKit).
- I18n: Use `useTranslation` hook for every string.

## AI Flow Standards

1. **Pydantic Contracts**: Every `genkit.defineFlow` must have an input and output schema defined via Pydantic or JSON Schema equivalent.
2. **Confidentiality**: Sanitise logs and traces. No PII should ever persist in `Langfuse` or `Loki`.
3. **Determinism**: Use temperature `0.0` for SQL generation and data extraction tasks.
4. **Resilience**: Implement `tenacity` retry patterns for all external API calls (Gemini, Tally, HIMS).

## Documentation Pattern

When creating a new agent or feature:
1. Create a spec in `docs/agents/specs/[id]-[name].md`.
2. Update `docs/agents/00-agent-backlog.md`.
3. Link the spec in `docs/agents/BLUEPRINTS.md` if high priority.
4. Create the corresponding skill in `.agent/skills/[name]/SKILL.md`.
