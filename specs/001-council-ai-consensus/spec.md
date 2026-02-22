# Feature Specification: Council of AIs Consensus System

**Feature Branch**: `001-council-ai-consensus`
**Created**: 2026-02-22
**Status**: Draft
**Input**: User description: "To eradicate hallucinations, the system should implement a 'Council of AIs' methodology, forcing consensus among multiple agent instances, grounded by a Graph-of-Thoughts retrieval mechanism querying a robust Medical Knowledge Graph."

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Verified Response to Medical Query (Priority: P1)

A healthcare administrator asks a complex question about patient revenue trends and medication costs. The system must provide an accurate, hallucination-free response by consulting multiple AI agents that verify facts against the Medical Knowledge Graph before reaching consensus.

**Why this priority**: This is the core value proposition - ensuring users receive accurate, trustworthy answers for critical healthcare business decisions. Without this, the entire system fails its primary purpose.

**Independent Test**: Can be fully tested by submitting a healthcare query and verifying that the response includes confidence indicators, source attribution, and passes accuracy validation against known correct answers.

**Acceptance Scenarios**:

1. **Given** a user submits a query about patient revenue trends, **When** the Council of AIs processes the request, **Then** at least 3 agent instances independently analyze the query and retrieve supporting evidence from the Knowledge Graph before a consensus response is generated.

2. **Given** multiple agents are deliberating on a response, **When** all agents reach agreement (consensus threshold met), **Then** the system presents a unified response with a consensus confidence score of at least 80%.

3. **Given** a user query requires medical domain knowledge, **When** the Graph-of-Thoughts retrieval executes, **Then** the system traces reasoning paths through the Medical Knowledge Graph and surfaces relevant evidence nodes to support the response.

---

### User Story 2 - Disagreement Detection and Uncertainty Signaling (Priority: P2)

When the Council of AIs cannot reach consensus on an answer, the system must transparently communicate uncertainty rather than presenting potentially incorrect information as fact.

**Why this priority**: User trust depends on honest communication about certainty levels. This prevents false confidence in uncertain situations and protects against decision-making based on unreliable information.

**Independent Test**: Can be tested by submitting queries with intentionally ambiguous or incomplete data and verifying the system correctly identifies and communicates uncertainty.

**Acceptance Scenarios**:

1. **Given** agents produce conflicting answers for a query, **When** the consensus threshold is not met after deliberation, **Then** the system displays an uncertainty indicator with the range of agent positions and confidence levels.

2. **Given** the Knowledge Graph lacks sufficient information for a query, **When** retrieval returns low-relevance evidence, **Then** the system signals "insufficient knowledge" rather than generating an unsupported response.

3. **Given** a partial consensus exists among agents, **When** presenting the response, **Then** the system shows which aspects have consensus and which remain contested with clear visual distinction.

---

### User Story 3 - Knowledge Graph Evidence Exploration (Priority: P2)

Users need to understand the evidence supporting AI responses, enabling them to verify claims and build trust in the system's reasoning process.

**Why this priority**: Transparency is essential for user adoption and trust, especially in healthcare contexts where decisions have significant consequences.

**Independent Test**: Can be tested by requesting evidence for any response and verifying that the system displays relevant Knowledge Graph nodes with clear connections to the answer.

**Acceptance Scenarios**:

1. **Given** a response has been generated, **When** the user requests evidence details, **Then** the system displays the Knowledge Graph paths traversed, highlighting the key evidence nodes that informed the consensus.

2. **Given** multiple evidence paths exist for a response, **When** displaying evidence, **Then** the system shows alternative reasoning routes considered by different agents during deliberation.

3. **Given** a user clicks on an evidence node, **When** the node expands, **Then** related concepts and their relationships in the Knowledge Graph become visible, allowing deeper exploration.

---

### User Story 4 - Response Accuracy Audit Trail (Priority: P3)

Administrators need to review historical queries and verify that the Council of AIs produced accurate, non-hallucinated responses for compliance and quality assurance purposes.

**Why this priority**: Supports organizational governance and continuous improvement but is not essential for day-to-day operation.

**Independent Test**: Can be tested by accessing the audit log and verifying that each query record contains complete deliberation data, evidence references, and consensus metrics.

**Acceptance Scenarios**:

1. **Given** a query has been processed, **When** an administrator reviews the audit trail, **Then** the system displays the complete deliberation record including each agent's initial response, retrieved evidence, and final consensus outcome.

2. **Given** an audit record is selected, **When** viewing details, **Then** the administrator can see the exact Knowledge Graph nodes accessed and the reasoning chain that led to the final response.

3. **Given** an audit reveals a potential hallucination, **When** flagged for review, **Then** the incident is logged with all relevant context for quality improvement analysis.

---

### Edge Cases

- What happens when the Medical Knowledge Graph is temporarily unavailable or incomplete for a specific query domain?
- How does the system handle queries that span multiple knowledge domains (e.g., clinical + financial)?
- What happens when agent instances produce syntactically different but semantically equivalent answers?
- How does the system behave when consensus requires an extended deliberation time that exceeds user response time expectations?
- What happens if an agent instance fails or times out during deliberation?

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST invoke a minimum of 3 independent AI agent instances for each user query requiring consensus validation.

- **FR-002**: System MUST require each agent to retrieve supporting evidence from the Medical Knowledge Graph using Graph-of-Thoughts traversal before submitting its response.

- **FR-003**: System MUST implement a consensus mechanism that measures agreement among agents and only releases responses when the consensus threshold (default 80%) is met.

- **FR-004**: System MUST calculate and display a confidence score for every response, derived from agent agreement levels and evidence strength.

- **FR-005**: System MUST trace reasoning paths through the Knowledge Graph and associate each response with its evidence trail.

- **FR-006**: System MUST signal uncertainty when consensus cannot be reached, displaying the range of agent positions rather than selecting a single answer.

- **FR-007**: System MUST maintain an audit trail of all deliberations, including agent responses, evidence accessed, and consensus outcomes.

- **FR-008**: System MUST allow users to explore the evidence nodes and reasoning paths that supported any response.

- **FR-009**: System MUST handle agent instance failures gracefully, continuing deliberation with remaining agents if minimum quorum (2 agents) is maintained.

- **FR-010**: System MUST validate retrieved evidence against the Knowledge Graph schema before including it in deliberation.

- **FR-011**: System MUST support configuration of consensus thresholds per query type or user role.

- **FR-012**: System MUST provide source attribution for all factual claims in responses, linking back to Knowledge Graph evidence.

### Key Entities

- **Council Deliberation**: Represents a single query processing session involving multiple agents. Contains the original query, participating agents, deliberation timeline, consensus outcome, and final response.

- **Agent Instance**: An independent AI reasoning unit participating in the Council. Each instance analyzes queries, retrieves evidence, and produces responses independently of other agents.

- **Knowledge Graph Node**: A unit of verified medical/healthcare knowledge within the Medical Knowledge Graph. Contains concept definitions, relationships to other concepts, and metadata about source and reliability.

- **Evidence Trail**: The sequence of Knowledge Graph nodes and relationships traversed during Graph-of-Thoughts retrieval for a specific query. Links deliberations to their supporting evidence.

- **Consensus Record**: Captures the agreement state among agents for a deliberation. Includes individual agent positions, agreement metrics, and threshold evaluation results.

- **Confidence Score**: A numerical value (0-100%) representing the system's certainty in a response, calculated from agent agreement and evidence strength.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Hallucination rate (verified incorrect factual claims) in system responses is reduced to less than 2% as measured by periodic human audit of sampled responses.

- **SC-002**: 95% of queries that achieve consensus return a response within 10 seconds of query submission.

- **SC-003**: Users correctly identify the confidence level of responses at least 90% of the time when surveyed, indicating clear uncertainty communication.

- **SC-004**: Evidence exploration feature is used by at least 40% of active users within 30 days of feature availability, indicating perceived value in transparency.

- **SC-005**: System correctly identifies and signals uncertainty for at least 85% of queries where ground truth is unknown or Knowledge Graph coverage is insufficient.

- **SC-006**: Audit trail completeness reaches 100% - every processed query has a retrievable record of deliberation, evidence, and consensus outcome.

- **SC-007**: User trust metric (measured via quarterly survey) shows 30% improvement in "confidence in AI response accuracy" compared to pre-implementation baseline.

- **SC-008**: System maintains 99.5% availability for consensus deliberations, with graceful degradation to single-agent mode when quorum cannot be met.

---

## Assumptions

- The Medical Knowledge Graph already exists or will be developed in parallel with this feature. It contains verified medical, pharmaceutical, and healthcare business domain knowledge.

- The underlying AI infrastructure supports running multiple agent instances concurrently without significant performance degradation.

- Users have basic familiarity with confidence indicators and understand that lower scores suggest the need for additional verification.

- The organization has established ground truth datasets for measuring hallucination rates during testing and ongoing quality monitoring.

- Network latency between agent instances and the Knowledge Graph is sufficient to support the target response times.

- The consensus algorithm's default 80% threshold is a starting point that may be tuned based on domain requirements and user feedback.

---

## Out of Scope

- Automated updates to the Medical Knowledge Graph based on new information (requires separate governance process).

- Real-time streaming of deliberation progress to users (deliberation happens behind the scenes).

- Integration with external knowledge sources beyond the Medical Knowledge Graph in the initial release.

- Custom agent instance configuration by end users (admin-level capability only).

- Multi-language support for deliberation processes (handled by separate i18n layer).
