---
name: bank-reconciliation
description: Automatically match bank statement entries to existing Tally ledger entries; assign confidence scores; surface unmatched items for manual resolution.
---

# Bank Reconciliation Skill

Guidelines for reconciling bank statements with internal Tally accounting records in MediSync.

## Matching Algorithm Logic

### Tiered Matching Strategy
1. **Exact Match (High Confidence)**:
    - Amount matches exactly.
    - Date within +/- 1 day.
    - Description contains the same reference number or vendor name.
2. **Soft Match (Medium Confidence)**:
    - Amount matches exactly.
    - Date within +/- 7 days (checks for weekend processing lag).
3. **Fuzzy Match (Low Confidence)**:
    - Amount matches ± 10 INR (handling small bank fees/rounding).
    - Description similarity > 80% using `RapidFuzz`.

## Tool Chain Patterns

### Fuzzy String Matching
```python
from rapidfuzz import fuzz

def calculate_match_score(bank_desc, tally_desc):
    # Ratio calculation (Levenshtein distance based)
    return fuzz.token_sort_ratio(bank_desc, tally_desc)
```

### SQL Date-Window Query
```sql
SELECT * FROM tally.vouchers
WHERE amount = :bank_amount
  AND voucher_date BETWEEN :stmt_date - INTERVAL '7 days' AND :stmt_date + INTERVAL '7 days'
  AND reconciled = FALSE;
```

## Reconciliation Workflow

### Handling Splits
- If one bank entry matches multiple small vouchers (Batch Payment), the agent must suggest a "Split Match".
- Use the **LangChain Reasoning Chain** to verify if `sum(vouchers) == bank_amount`.

### Outcome Actions
- **Reconciled**: Mark both entries as matched in the DB.
- **Partial**: User must manually adjust one side.
- **No Match**: Create a "Suspense" entry or flag for manual search.

## Accuracy & Quality

- **Confidence Score**: Aggregated from date proximity, amount exactness, and description similarity.
- **Report**: Generate an "Outstanding Items Report" for any entry older than 30 days without a match.

## Accessibility Checklist
- [ ] Side-by-side view: Bank on left, Tally on right.
- [ ] Color-code match confidence (Green/Amber/Red).
- [ ] Provide "Manual Match" drag-and-drop interface for users.
- [ ] Support large CSV/PDF bank statement uploads.
