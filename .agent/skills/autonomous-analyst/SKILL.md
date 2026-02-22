---
name: autonomous-analyst
description: Accept high-level business questions, autonomously decompose into sub-tasks (retrieve, analyze, forecast), and return a structured insights report.
---

# Autonomous AI Analyst Skill

Guidelines for sophisticated multi-agent data analysis and reporting within the MediSync platform.

## Orchestration Patterns (CrewAI)

### Task Decomposition
Agents must break down a prompt like "Analyze pharmacy performance" into:
1. **Retrieval**: Query SQL for sales, lab costs, or operational metrics.
2. **Analysis**: Compute growth, margins, and top products.
3. **Forecasting**: Predict next month's demand.
4. **Synthesis**: Write the narrative executive summary.

### Multi-Agent Interaction
```python
from crewai import Agent, Task, Crew

data_analyst = Agent(
    role='Financial Analyst',
    goal='Identify profit leakage points in pharmacy operations',
    backstory='Expert in healthcare accounting, lab operations, and legacy supply chain data.',
    tools=[sql_query_tool, stats_tool]
)

analysis_task = Task(description='...', agent=data_analyst)
crew = Crew(agents=[data_analyst], tasks=[analysis_task])
result = crew.start()
```

## Analytical Tools

### Statistical Analysis
- Use **statsmodels** for trend analysis and significance testing.
- **Goal**: Don't just show numbers; explain if a change is statistically meaningful.

### Forecasting (Prophet)
- Use **Facebook Prophet** for 90-day time-series forecasting.
- **Guidelines**: Always include confidence intervals (Upper/Lower bounds) in the report output.

## Reporting Standards

### The "Insight Narrative"
Every report must follow a structured Pydantic model:
- `summary`: One-sentence takeaway.
- `key_findings`: Bullet points with data citations.
- `charts`: JSON definitions for Plotly/Seaborn.
- `recommendations`: Concrete actions derived from the data.

## Accuracy & Quality

- **Hallucination Check**: If `confidence_score < 0.8`, the report must be labeled as "Preliminary - Requires Human Verification".
- **Data Traceability**: Every chart and finding must reference the source table/query (e.g., "Source: `tally.ledger_entries`" or "`lims.lab_results`").

## Accessibility Checklist
- [ ] Use plain English titles (no technical jargon).
- [ ] Ensure all charts are colorblind-friendly.
- [ ] Provide a "Download as PDF" option for shared reports.
- [ ] Support voice-to-query input via mobile interface.
