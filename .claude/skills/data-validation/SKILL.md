---
name: data-validation
description: Run automated data quality checks on every ETL batch before data lands in the warehouse; block bad data; alert on anomalies.
---

# Data Quality Validation Skill

Guidelines for ensuring data integrity across the MediSync platform using declarative validation frameworks.

## Validation Framework: Great Expectations

### Declarative Assertions
- **Completeness**: `expect_column_values_to_not_be_null` for critical financial IDs.
- **Uniqueness**: `expect_column_values_to_be_unique` for invoice numbers and voucher IDs.
- **Referential Integrity**: `expect_column_values_to_be_in_set` (e.g., Ledger names must exist in the Master list).
- **Logical Bounds**: `expect_column_values_to_be_between` (e.g., Transaction amount > 0).

## Tool Chain Patterns

### Running Expectations in Airflow
```python
import great_expectations as ge

def validate_batch(df):
    ge_df = ge.from_pandas(df)
    results = ge_df.expect_column_values_to_be_between("amount", 0, 10000000)
    
    if not results.success:
        raise DataQualityError("Amount out of expected range detected.")
```

### Anomaly Detection (Z-Score)
Monitor daily transaction volumes to detect spikes or drops:
```python
def check_volume_anomaly(current_volume, historical_avg, historical_std):
    z_score = (current_volume - historical_avg) / historical_std
    if abs(z_score) > 3:
        return True # Anomaly detected
    return False
```

## Response Actions

| Validation Result | Action |
|---|---|
| **Critical Failure** | Block ETL load, trigger PagerDuty/Apprise alert. |
| **Warning** | Load data but tag as "Low Quality" in BI dashboards. |
| **Pass** | Append results to the Data Quality Audit Log and proceed. |

## Accuracy & Quality

- **Source Reconciliation**: Periodically compare HIMS billing totals with Tally ledger totals. Differences > 1% should trigger an investigation ticket.
- **Schema Validation**: Ensure incoming JSON/CSV files match the Pydantic models used by downstream agents.

## Accessibility Checklist
- [ ] Provide "Data Health Dashboard" showing daily pass/fail rates.
- [ ] Enable users to define custom validation rules via a simple UI.
- [ ] Send weekly "Data Integrity Report" to the CTO/Admin role.
