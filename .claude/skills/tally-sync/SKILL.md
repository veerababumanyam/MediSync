---
name: tally-sync
description: Push approved transactions from the MediSync platform into Tally ERP via TDL XML API. Safe, idempotent, and auditable synchronization.
---

# Tally Sync Skill

Guidelines for synchronizing approved financial data from MediSync to Tally Prime/ERP 9.

> [!NOTE]
> For detailed XML/JSON API structures and Tally Gateway configuration, refer to the [Tally Developer Reference](file:///Users/v13478/Desktop/MediSync/References/tally-reference.md).

## Integration Architecture

### TDL XML Communication
- **Format**: All data must be wrapped in `<TALLYMESSAGE>` XML tags.
- **Transport**: HTTP POST to Tally Gateway (usually port 9000).
- **Template Engine**: Use **Jinja2** to generate XML payloads dynamically.

### Idempotency & Safety
- **Duplicate Guard**: Generate a unique `RemoteID` or hash for every voucher (e.g., `hash(vendor_id + date + invoice_no + amount)`).
- **Check-Before-Push**: Always verify that the Tally Ledger and Group exist before attempting to push a voucher.

## Code Patterns

### Jinja2 TDL Template (Journal Voucher)
```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDESC>
                <REPORTNAME>All Masters</REPORTNAME>
            </REQUESTDESC>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="{{ voucher_type }}" ACTION="Create" OBJVIEW="Accounting Voucher View">
                        <DATE>{{ date }}</DATE>
                        <VOUCHERNUMBER>{{ invoice_no }}</VOUCHERNUMBER>
                        {% for entry in entries %}
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>{{ entry.ledger_name }}</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>{{ entry.is_debit }}</ISDEEMEDPOSITIVE>
                            <AMOUNT>{{ entry.amount }}</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                        {% endfor %}
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Error Handling

| Error Code | Meaning | Agent Action |
|---|---|---|
| **Ledger Missing** | Master does not exist in Tally | Flag as "Sync Error" and prompt user to create master. |
| **Timeout** | Tally Gateway unreachable | Exponential backoff retry (3x). |
| **Auth Failure** | Invalid API configuration | Alert Finance Head immediately. |

## Sync Standards

- **Atomic Batches**: If pushing 10 vouchers, ensure the UI reflects status for each individual sub-task.
- **Audit Logging**: Write every sync result (Success/Failure + Tally Response XML) to the platform's immutable audit trail.
- **Read-Back Verification**: (Optional) After sync, verify the voucher exists in Tally by querying it back by `VoucherNumber`.

## Accessibility Checklist
- [ ] Provide a "Sync Log" view showing detailed history of successful and failed pushes.
- [ ] Enable "Export as XML" for manual sync as a fallback.
- [ ] Show "One-click Sync all approved" button in the summary dashboard.
