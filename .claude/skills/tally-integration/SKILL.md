---
name: tally-integration
description: Guides developers integrating Tally ERP with MediSync using TDL XML over HTTP. Covers XML request/response formats, voucher types, sync workflows, and security gates. Use when creating TDL templates, mapping ledgers, syncing transactions, or testing Tally integration.
---

# Tally Integration Guide

Guidelines for integrating MediSync with TallyPrime/ERP 9 via TDL XML.

★ Insight ─────────────────────────────────────
Tally integration is the critical "action plane" component. All data movement
from MediSync to Tally MUST pass through human approval gates (HITL) to prevent
unintended financial ledger modifications.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Protocol** | HTTP POST with TDL XML |
| **Port** | 9000 (default, configurable) |
| **Format** | XML wrapped in `<ENVELOPE>` |
| **Idempotency** | Use `RemoteID` hash to prevent duplicates |
| **Security** | OPA policy + HITL approval required |

## Enabling Tally Gateway

1. Open TallyPrime
2. Press **F1** → **Settings** → **Advanced Configuration**
3. Set **Enable HTTP Server** to **Yes**
4. Set **Port** to desired port (default: 9000)

## TDL XML Structure

All Tally communication uses the envelope structure:

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
        <VERSION>1</VERSION>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDESC>
                <REPORTNAME>All Masters</REPORTNAME>
            </REQUESTDESC>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <!-- Voucher or Master data -->
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Basic Voucher Types

### Journal Voucher

Used for accounting entries without cash/bank involvement:

```xml
<VOUCHER VCHTYPE="Journal" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>JV-001</VOUCHERNUMBER>
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Sales Account</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>5000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Customer - Al Futtaim</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>5000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>
</VOUCHER>
```

### Purchase Voucher

Used for recording purchases:

```xml
<VOUCHER VCHTYPE="Purchase" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>PUR-001</VOUCHERNUMBER>
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Pharmacy Supplies</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>10000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>
    <ALLINVENTORYENTRIES.LIST>
        <ITEMNAME>Paracetamol 500mg</ITEMNAME>
        <RATE>50</RATE>
        <AMOUNT>5000</AMOUNT>
    </ALLINVENTORYENTRIES.LIST>
</VOUCHER>
```

### Payment Voucher

Used for cash/bank payments:

```xml
<VOUCHER VCHTYPE="Payment" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>PAY-001</VOUCHERNUMBER>
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>ICICI Bank</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>25000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Salary Expense</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>25000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>
</VOUCHER>
```

## Debit/Credit Pattern

Understanding `ISDEEMEDPOSITIVE`:

| Entry Type | ISDEEMEDPOSITIVE | Effect |
|------------|------------------|--------|
| Income/Credit | Yes | Increases the account |
| Income/Credit | No | Decreases the account |
| Expense/Debit | Yes | Decreases the account |
| Expense/Debit | No | Increases the account |

**Quick Rule**: For ledger entries in a voucher:
- If receiving = `Yes`
- If giving = `No`

## Sync Workflow

```
┌─────────────────────────────────────────────────────────────┐
│  1. MediSync: Generate Entry (Draft)                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Module B-08: Approval Workflow (HITL Gate)               │
│     - Verify ledger mappings                                 │
│     - Check amounts and dates                                │
│     - Finance head approval                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  3. OPA Policy Check                                        │
│     - Verify user has tally.sync permission                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Generate TDL XML (Jinja2 Template)                      │
│     - Create ENVELOPE structure                              │
│     - Include RemoteID for idempotency                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  5. POST to Tally Gateway (port 9000)                        │
│     - Content-Type: application/xml                          │
│     - Timeout: 30 seconds                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  6. Parse Response                                           │
│     - STATUS 1 = Success                                     │
│     - STATUS 0 = Failure (check error tags)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  7. Verify & Audit Log                                       │
│     - Record sync result                                     │
│     - Store Tally response XML                               │
└─────────────────────────────────────────────────────────────┘
```

## Error Handling

| Error | Meaning | Action |
|-------|---------|--------|
| **Ledger Missing** | Master doesn't exist in Tally | Flag as "Sync Error", prompt user to create master |
| **Timeout** | Tally Gateway unreachable | Exponential backoff retry (max 3 attempts) |
| **Duplicate Voucher** | VoucherNumber already exists | Check before sync, use different number |
| **Auth Failure** | Invalid credentials | Alert Finance Head immediately |
| **Invalid Date** | Date format incorrect | Use YYYYMMDD format |

## Idempotency Pattern

Prevent duplicate entries:

```go
func GenerateRemoteID(entry JournalEntry) string {
    data := fmt.Sprintf("%s|%s|%s|%f",
        entry.CompanyID,
        entry.Date.Format("2006-01-02"),
        entry.VoucherType,
        entry.TotalAmount,
    )
    h := sha256.Sum256([]byte(data))
    return fmt.Sprintf("medisync-%x", h[:8])
}
```

## Jinja2 Template Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="{{ voucher_type }}" ACTION="Create">
                        <DATE>{{ date|strftime('%Y%m%d') }}</DATE>
                        <VOUCHERNUMBER>{{ voucher_number }}</VOUCHERNUMBER>
                        <REMOTEID>{{ remote_id }}</REMOTEID>
                        {% for entry in ledger_entries %}
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>{{ entry.ledger_name }}</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>{{ 'Yes' if entry.is_debit else 'No' }}</ISDEEMEDPOSITIVE>
                            <AMOUNT>{{ entry.amount|round(2) }}</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                        {% endfor %}
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Security Gates

All sync operations must pass through:

1. **HITL Approval** (Module B-08)
   - Finance head sign-off
   - Multi-level for amounts > threshold

2. **OPA Policy Check**
   ```go
   allowed, _ := opa.Allow(ctx, "tally_sync", map[string]interface{}{
       "user_id": userID,
       "company_id": companyID,
       "action": "sync_journal_entry",
       "amount": entry.TotalAmount,
   })
   ```

3. **Audit Logging**
   ```go
   audit.Log(ctx, AuditEntry{
       UserID: userID,
       Action: "tally_sync",
       Resource: entry.ID,
       Status: "approved",
   })
   ```

## Testing Tally Integration

### Unit Test with Mock Server

```go
func TestTallySync_MockServer(t *testing.T) {
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<ENVELOPE><STATUS>1</STATUS></ENVELOPE>`))
    }))
    defer mockServer.Close()

    sync := NewTallySync(mockServer.URL)
    err := sync.Push(ctx, entry)
    assert.NoError(t, err)
}
```

### Integration Test Checklist

- [ ] Mock Tally responds with valid XML
- [ ] Ledger exists before pushing voucher
- [ ] RemoteID is generated correctly
- [ ] OPA policy denies unauthorized sync
- [ ] Audit log entry is created
- [ ] Timeout is handled correctly

## Common Tasks

### Create a Journal Entry
See `references/voucher-types.md` for complete examples.

### Map Ledgers
See `examples/` for Jinja2 templates.

### Test Sync Connection
```bash
curl -X POST http://localhost:9000 \
  -H "Content-Type: application/xml" \
  -d @examples/test-request.xml
```

### View Sync Log
```sql
SELECT * FROM audit_log WHERE action = 'tally_sync' ORDER BY created_at DESC;
```

## Detailed References

| Reference | Content |
|-----------|---------|
| `references/xml-formats.md` | Complete XML request/response schemas |
| `references/voucher-types.md` | Journal, Purchase, Sales, Payment voucher formats |
| `references/sync-workflow.md` | Detailed sync pipeline and retry logic |
| `references/security.md` | OPA policies and HITL gate configuration |
| `references/testing.md` | Mock responses and integration test patterns |
| `examples/journal-entry.xml` | Sample journal voucher TDL |
| `examples/purchase-bill.xml` | Sample purchase voucher TDL |
| `scripts/validate-tally-xml.sh` | XML validation script |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Connection refused | Verify Tally Gateway is enabled on port 9000 |
| Ledger not found | Create master in Tally before syncing voucher |
| Duplicate voucher | Check VoucherNumber exists before pushing |
| XML parse error | Verify date format is YYYYMMDD |
| Timeout | Increase timeout for large batches |

> **Note**: For the complete Tally developer reference, see the official Tally integration documentation in `/References/Tally/` PDF files.
