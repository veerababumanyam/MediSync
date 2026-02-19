# Tally XML Formats Reference

Complete reference for TallyPrime XML request and response formats.

## Base Envelope Structure

All Tally XML messages are wrapped in an `<ENVELOPE>` element:

```xml
<ENVELOPE>
    <HEADER>
        <!-- Request metadata -->
    </HEADER>
    <BODY>
        <!-- Request data -->
    </BODY>
</ENVELOPE>
```

## HEADER Section

### Import Request

```xml
<HEADER>
    <TALLYREQUEST>Import Data</TALLYREQUEST>
    <VERSION>1</VERSION>
</HEADER>
```

### Export Request

```xml
<HEADER>
    <TALLYREQUEST>Export Data</TALLYREQUEST>
    <VERSION>1</VERSION>
</HEADER>
```

### Execute Request (for functions)

```xml
<HEADER>
    <TALLYREQUEST>Execute</TALLYREQUEST>
    <ID>ReportName</ID>
    <VERSION>1</VERSION>
</HEADER>
```

## BODY Section Structure

### IMPORTDATA Structure

Used for creating vouchers and masters:

```xml
<BODY>
    <IMPORTDATA>
        <REQUESTDESC>
            <REPORTNAME>All Masters</REPORTNAME>
            <STATICVARIABLES>
                <SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>
            </STATICVARIABLES>
        </REQUESTDESC>
        <REQUESTDATA>
            <!-- One or more TALLYMESSAGE elements -->
            <TALLYMESSAGE xmlns:UDF="TallyUDF">
                <!-- Voucher or Master content -->
            </TALLYMESSAGE>
        </REQUESTDATA>
    </IMPORTDATA>
</BODY>
```

### EXPORTDATA Structure

Used for fetching data from Tally:

```xml
<BODY>
    <EXPORTDATA>
        <REQUESTDESC>
            <REPORTNAME>All Masters</REPORTNAME>
            <STATICVARIABLES>
                <SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>
            </STATICVARIABLES>
        </REQUESTDESC>
        <REQUESTDATA>
            <!-- Filter criteria -->
        </REQUESTDATA>
    </EXPORTDATA>
</BODY>
```

## TALLYMESSAGE Structure

The `TALLYMESSAGE` element contains the actual data:

```xml
<TALLYMESSAGE xmlns:UDF="TallyUDF">
    <VOUCHER VCHTYPE="Journal" ACTION="Create" OBJVIEW="Accounting Voucher View">
        <!-- Voucher fields -->
    </VOUCHER>
</TALLYMESSAGE>
```

### Common Attributes

| Attribute | Values | Description |
|-----------|--------|-------------|
| `xmlns:UDF` | `TallyUDF` | Namespace for User Defined Fields |
| `ACTION` | `Create`, `Alter`, `Delete` | Operation to perform |
| `OBJVIEW` | View name | Tally view context |

## Voucher XML Schema

### Complete Voucher Structure

```xml
<VOUCHER VCHTYPE="Journal" ACTION="Create" OBJVIEW="Accounting Voucher View">
    <!-- Basic Fields -->
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>JV-001</VOUCHERNUMBER>
    <NARRATION>Monthly rent payment</NARRATION>

    <!-- Reference Fields -->
    <REFERENCE>
        <NUMBER>REF-001</NUMBER>
        <DATE>20260219</DATE>
    </REFERENCE>

    <!-- MediSync Extension -->
    <REMOTEID>medisync-a1b2c3d4</REMOTEID>

    <!-- Ledger Entries -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Rent Expense</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>15000</AMOUNT>
        <BANKALLOCTYPE>NotApplicable</BANKALLOCTYPE>
    </ALLLEDGERENTRIES.LIST>

    <!-- Cost Center (Optional) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>ICICI Bank</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>15000</AMOUNT>
        <COSTCENTRE.LIST>
            <NAME>Main Branch</NAME>
            <COSTALLOCATION>15000</COSTALLOCATION>
        </COSTCENTRE.LIST>
    </ALLLEDGERENTRIES.LIST>

    <!-- UDF Fields -->
    <UDF:MSYNC.SYNCEDBY>john.doe@medisync.com</UDF:MSYNC.SYNCEDBY>
    <UDF:MSYNC.SYNCDATETIME>2026-02-19 10:30:00</UDF:MSYNC.SYNCDATETIME>
</VOUCHER>
```

## Master XML Schema

### Ledger Master

```xml
<LEDGER NAME="Patient Advances - XYZ Clinic" ACTION="Create">
    <PARENT>Sundry Debtors</PARENT>
    <LANGUAGENAME.LIST>
        <LANGUAGE.LIST>
            <NAME.LIST>
                <NAME>Patient Advances - XYZ Clinic</NAME>
            </NAME.LIST>
            <LANGUAGEID>English</LANGUAGEID>
        </LANGUAGE.LIST>
        <LANGUAGE.LIST>
            <NAME.LIST>
                <NAME>مقدمات المرضى - عيادة XYZ</NAME>
            </NAME.LIST>
            <LANGUAGEID>Arabic</LANGUAGEID>
        </LANGUAGE.LIST>
    </LANGUAGENAME.LIST>

    <!-- Ledger Properties -->
    <ISBILLWISENO>Yes</ISBILLWISENO>
    <ISCOSTCENTRESON>Yes</ISCOSTCENTRESON>
    <ISINTERESTON>No</ISINTERESTON>
    <ISPOSTVALUATIONON>No</ISPOSTVALUATIONON>

    <!-- GST (India) or VAT (UAE/KSA) -->
    <GSTTYPE>Regular</GSTTYPE>
    <GSTAPPLICABILITY>Applicable</GSTAPPLICABILITY>

    <!-- Banking (for bank accounts) -->
    <ISBANKPAYMENTON>No</ISBANKPAYMENTON>
    <ISBANKRECONON>No</ISBANKRECONON>

    <!-- MediSync UDF -->
    <UDF:MSYNC.LEDGERID>LED-12345</UDF:MSYNC.LEDGERID>
</LEDGER>
```

### Cost Center Master

```xml
<COSTCENTRE NAME="Dubai Branch" ACTION="Create">
    <PARENT>Primary Cost Centre</PARENT>
    <COSTCENTRECLASS>Dubai Operations</COSTCENTRECLASS>
    <ISGROUPSON>No</ISGROUPSON>
</COSTCENTRE>
```

## Response XML Schema

### Success Response

```xml
<ENVELOPE>
    <HEADER>
        <VERSION>1</VERSION>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
        <ID>All Masters</ID>
        <STATUS>1</STATUS>
    </HEADER>
    <BODY>
        <DATA>
            <LINEERROR>0</LINEERROR>
            <COLLECTION>
                <TYPE>Vouchers</TYPE>
            </COLLECTION>
        </DATA>
    </BODY>
</ENVELOPE>
```

### Error Response

```xml
<ENVELOPE>
    <HEADER>
        <VERSION>1</VERSION>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
        <STATUS>0</STATUS>
    </HEADER>
    <BODY>
        <DATA>
            <LINEERROR>1</LINEERROR>
            <ERRORLINE>5</ERRORLINE>
            <ERRORMESSAGE>Ledger 'Unknown Ledger' does not exist</ERRORMESSAGE>
        </DATA>
    </BODY>
</ENVELOPE>
```

### Partial Success Response

```xml
<ENVELOPE>
    <HEADER>
        <VERSION>1</VERSION>
        <STATUS>1</STATUS>
    </HEADER>
    <BODY>
        <DATA>
            <LINEERROR>2</LINEERROR>
            <LINEMESSAGE>
                <LINEINFO>
                    <LINEINDEX>1</LINEINDEX>
                    <STATUS>1</STATUS>
                </LINEINFO>
                <LINEINFO>
                    <LINEINDEX>2</LINEINDEX>
                    <STATUS>0</STATUS>
                    <ERRORMESSAGE>Voucher number already exists</ERRORMESSAGE>
                </LINEINFO>
            </LINEMESSAGE>
        </DATA>
    </BODY>
</ENVELOPE>
```

## Date Formats

| Context | Format | Example |
|---------|--------|---------|
| XML DATE element | YYYYMMDD | 20260219 |
| XML attributes | YYYYMMDD | 20260219 |
| UDF datetime | YYYY-MM-DD HH:MM:SS | 2026-02-19 10:30:00 |
| Hijri dates (Arabic) | Depends on Tally locale | 1440-08-20 |

## Number Formats

| Type | Format | Example |
|------|--------|---------|
| Amount | Decimal with 2 places | 15000.00 |
| Quantity | Decimal with 3 places | 100.500 |
| Rate | Decimal with 2 places | 50.75 |

**Note**: Tally does not use thousand separators in XML. Always use plain numbers.

## Special Characters in XML

| Character | Escape Sequence |
|-----------|-----------------|
| `<` | `&lt;` |
| `>` | `&gt;` |
| `&` | `&amp;` |
| `'` | `&apos;` |
| `"` | `&quot;` |

For Arabic text, ensure XML encoding is UTF-8:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<ENVELOPE>
    <TALLYMESSAGE>
        <UDF:MSYNC.DESCRIPTION>دفع إيجار شهري</UDF:MSYNC.DESCRIPTION>
    </TALLYMESSAGE>
</ENVELOPE>
```

## Batch Requests

Multiple vouchers in a single request:

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <!-- Voucher 1 -->
                <TALLYMESSAGE>
                    <VOUCHER VCHTYPE="Journal">...</VOUCHER>
                </TALLYMESSAGE>

                <!-- Voucher 2 -->
                <TALLYMESSAGE>
                    <VOUCHER VCHTYPE="Payment">...</VOUCHER>
                </TALLYMESSAGE>

                <!-- Voucher N -->
                <TALLYMESSAGE>
                    <VOUCHER VCHTYPE="Receipt">...</VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

**Limit**: Keep batch size under 50 vouchers per request for optimal performance.

## UDF (User Defined Fields)

Custom fields for MediSync tracking:

```xml
<VOUCHER xmlns:UDF="TallyUDF">
    <!-- Standard Tally fields -->

    <!-- MediSync custom fields -->
    <UDF:MSYNC.ENTRYID>ENT-12345</UDF:MSYNC.ENTRYID>
    <UDF:MSYNC.SYNCSTATUS>synced</UDF:MSYNC.SYNCSTATUS>
    <UDF:MSYNC.SYNCDATETIME>2026-02-19 10:30:00</UDF:MSYNC.SYNCDATETIME>
    <UDF:MSYNC.SYNCEDBY>john.doe@medisync.com</UDF:MSYNC.SYNCEDBY>
    <UDF:MSYNC.SOURCE>ocr</UDF:MSYNC.SOURCE>
    <UDF:MSYNC.DOCUMENTHASH>a1b2c3d4e5f6</UDF:MSYNC.DOCUMENTHASH>
</VOUCHER>
```

**Note**: UDF fields must be defined in Tally TDL before use.

## Validation Checklist

Before sending XML to Tally:

- [ ] XML is well-formed (matching tags)
- [ ] All required fields are present
- [ ] Date format is YYYYMMDD
- [ ] Amounts are positive numbers
- [ ] Ledger names match Tally masters exactly
- [ ] Voucher number is unique
- [ ] Debits = Credits (in voucher entries)
- [ ] Encoding is UTF-8
