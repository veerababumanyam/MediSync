---
name: ledger-mapping
description: Automatically map extracted transactions to the correct Tally Chart of Accounts (GL Ledgers) using vector similarity and historical patterns.
---

# Ledger Mapping Skill

Guidelines for categorizing financial transactions into the appropriate accounting ledgers within MediSync.

## Mapping Mechanisms

### Vector Similarity (RAG Pattern)
- **Tool**: Chroma DB.
- **Embedding**: BAAI/bge-small.
- **Process**:
    1. Generate an embedding for the `transaction_description` and `vendor_name`.
    2. Search Chroma for the top 5 most similar historical mappings.
    3. Pass matches as context to the LLM to decide the best current mapping.

### Historical Learning Loop
Any user correction to a mapping MUST be saved back to Chroma:
```python
# Pseudo-code for feedback loop
def update_mapping_memory(transaction_text, final_ledger_name):
    embedding = model.encode(transaction_text)
    chroma_collection.add(
        embeddings=[embedding],
        documents=[transaction_text],
        metadatas=[{"ledger": final_ledger_name}]
    )
```

## Prompt Design for Mapping
```
Given the transaction: "{{ vendor }} - {{ description }}"
And the following Tally Ledgers:
{{ chart_of_accounts_subset }}

And historical matches:
{{ historical_mappings }}

Select the most appropriate Ledger. If uncertain, suggest the top 3 alternatives.
```

## Classification Thresholds

| Confidence | Action taken | UI Badge |
|---|---|---|
| **> 0.95** | Auto-suggest (pre-selected) | High Confidence (Green) |
| **0.70 - 0.94** | Require explicit confirmation | Review Needed (Amber) |
| **< 0.70** | Leave blank, force manual input | Action Required (Red) |

## Implementation Principles

- **Separation of Concerns**: Mapping should only suggest. Writing only happens in the Tally Sync agent.
- **Multi-Entity Awareness**: Ensure the `chart_of_accounts` context is filtered by the active `company_id`.
- **Cost Centers**: Suggest both the Primary Ledger and the Cost Center/Department if available in the transaction context (e.g., "Pharmacy Supplies" -> "Pharmacy Dept").

## Accessibility Checklist
- [ ] Include "Reasoning" explaining why a particular ledger was chosen.
- [ ] Show alternative mappings for quick selection in the UI.
- [ ] Support bulk selection of mappings for similar transactions.
