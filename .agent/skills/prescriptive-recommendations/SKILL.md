---
name: prescriptive-recommendations
description: Go beyond descriptive insights to deliver specific, quantified, actionable recommendations — including root-cause analysis and expected business impact.
---

# Prescriptive Recommendations Skill

Guidelines for generating high-impact business recommendations that drive action in MediSync.

## Reasoning Framework (ReAct)

### The "Why" and the "What"
Recommendations must follow the **Insight → Root Cause → Action → Impact** chain:
- **Insight**: "Pharmacy margin dropped 5%."
- **Root Cause**: "Supplier X increased drug price by 12%."
- **Action**: "Switch to Supplier Y for the top 5 high-volume drugs."
- **Impact**: "Expected savings of ₹45,000 per month."

### Root Cause Analysis Tooling
Use the **LangChain ReAct** loop to dig deeper:
```python
def run_rca(target_metric):
    # Tool 1: Get sub-category breakdown
    # Tool 2: Check supplier price changes
    # Tool 3: Check inventory wastage levels
    pass
```

## Impact Quantification

### Scenario Forecasting
Run "What-If" simulations using **Prophet**:
- "If we implement recommendation A, the cash flow forecast improves by X%."
- "If we do nothing, we will stock out of Item Z in 3 days."

## Guardrails & Ethics

- **No Auto-Execution**: Recommendations never execute write actions (e.g., creating a Purchase Order) automatically. They must route to the **Approval Workflow Agent (B-08)**.
- **Explainability**: Always show the data behind the impact estimate.
- **Biased Source Guard**: Ensure recommendations don't unfairly penalize suppliers without sufficient data points.

## Accuracy & Quality

- **Confidence Score**: Recommendations with `< 0.85` confidence must be presented with "Potential Scenarios" rather than a single prescriptive path.
- **Stakeholder Alignment**: Tag recommendations with the required role (e.g., "Requires Pharmacy Manager Approval").

## Accessibility Checklist
- [ ] Use "Impact Badges" (e.g., Low Effort / High Impact).
- [ ] Link recommendations directly to the "Action" button in the ERP.
- [ ] Record all rejected recommendations for "Learning" (don't suggest the same rejected action twice without changing context).
