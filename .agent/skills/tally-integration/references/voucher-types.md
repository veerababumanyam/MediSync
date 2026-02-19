# Tally Voucher Types Reference

Complete reference for common voucher types used in MediSync integration.

## Overview

Tally uses different voucher types for different kinds of transactions. This guide covers the most common types used in healthcare/pharmacy operations.

## Journal Voucher (Journal)

**Use Case**: Accounting entries without cash/bank involvement, adjustments, accruals.

### Complete Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Journal" ACTION="Create" OBJVIEW="Accounting Voucher View">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>JV-2026-0219-001</VOUCHERNUMBER>
                        <REFERENCE>
                            <NUMBER>REF-001</NUMBER>
                            <DATE>20260219</DATE>
                        </REFERENCE>
                        <NARRATION>Monthly depreciation adjustment for medical equipment</NARRATION>
                        <REMOTEID>medisync-depr-20260219</REMOTEID>

                        <!-- Debit Entry (Asset Debit) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Depreciation - Equipment</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
                            <AMOUNT>5000</AMOUNT>
                            <BANKALLOCTYPE>NotApplicable</BANKALLOCTYPE>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Credit Entry (Accumulated Depreciation) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Accumulated Depreciation</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <AMOUNT>5000</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>

                        <UDF:MSYNC.SYNCEDBY>system@medisync.com</UDF:MSYNC.SYNCEDBY>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

### Field Descriptions

| Field | Required | Description |
|-------|----------|-------------|
| `VCHTYPE` | Yes | Must be "Journal" |
| `DATE` | Yes | YYYYMMDD format |
| `VOUCHERNUMBER` | Yes | Unique identifier |
| `ALLLEDGERENTRIES.LIST` | Yes | At least 2 entries (debit + credit) |

## Payment Voucher (Payment)

**Use Case**: Cash/Bank payments (salaries, rent, supplier payments, expenses).

### Complete Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Payment" ACTION="Create" OBJVIEW="Accounting Voucher View">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>PAY-2026-0219-001</VOUCHERNUMBER>

                        <!-- Bank Account (Credit - Money Going Out) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>ADCB Bank - Current Account</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
                            <AMOUNT>25000</AMOUNT>
                            <BANKALLOCTYPE>NotApplicable</BANKALLOCTYPE>
                            <transactionid>BANKTXN-12345</transactionid>
                            <paymentmethodslink />
                            <BANKACCOUNTNAME>ADCB-001</BANKACCOUNTNAME>
                            <instrument>cheque</instrument>
                            <CHEQUENUMBER>001234</CHEQUENUMBER>
                            <CHEQUEDATE>20260219</CHEQUEDATE>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Salary Expense (Debit) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Salaries - Nursing Staff</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <AMOUNT>15000</AMOUNT>
                            <costallocationlist>
                                <COSTCENTRE.LIST>
                                    <NAME>Main Branch</NAME>
                                    <COSTALLOCATION>15000</COSTALLOCATION>
                                </COSTCENTRE.LIST>
                            </costallocationlist>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Salary Expense (Debit) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Salaries - Administrative Staff</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <AMOUNT>10000</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>

                        <NARRATION>Monthly salary payment - February 2026</NARRATION>
                        <REMOTEID>medisync-pay-salary-20260219</REMOTEID>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

### Payment-Specific Fields

| Field | Required | Description |
|-------|----------|-------------|
| `instrument` | No | `cheque`, `neft`, `rtgs`, `card`, `cash` |
| `CHEQUENUMBER` | If instrument=cheque | Cheque number |
| `CHEQUEDATE` | If instrument=cheque | Cheque date (YYYYMMDD) |
| `transactionid` | For NEFT/RTGS | Bank transaction reference |

## Receipt Voucher (Receipt)

**Use Case**: Cash/Bank receipts (patient payments, insurance receipts).

### Complete Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Receipt" ACTION="Create" OBJVIEW="Accounting Voucher View">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>RCT-2026-0219-001</VOUCHERNUMBER>

                        <!-- Bank Account (Debit - Money Coming In) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>ADCB Bank - Current Account</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <AMOUNT>500</AMOUNT>
                            <BANKALLOCTYPE>NotApplicable</BANKALLOCTYPE>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Patient Account (Credit - Patient Liability Decreases) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Patient - Ahmed Al Mansouri</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
                            <AMOUNT>500</AMOUNT>
                            <BILLALLOCATIONS.LIST>
                                <NAME>BILL-2026-001</NAME>
                                <BILLTYPE>NewRef</BILLTYPE>
                                <AMOUNT>500</AMOUNT>
                                <INTERESTCOLLECTION>0</INTERESTCOLLECTION>
                            </BILLALLOCATIONS.LIST>
                        </ALLLEDGERENTRIES.LIST>

                        <NARRATION>Outpatient consultation payment</NARRATION>
                        <REMOTEID>medisync-rct-ahmed-20260219</REMOTEID>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Purchase Voucher (Purchase)

**Use Case**: Purchase of goods (medicines, medical supplies).

### Complete Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Purchase" ACTION="Create" OBJVIEW="Inventory Voucher View">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>PUR-2026-0219-001</VOUCHERNUMBER>
                        <PARTYLEDGERNAME>Pharma Supplies LLC</PARTYLEDGERNAME>
                        <PARTYNAME>Pharma Supplies LLC</PARTYNAME>
                        <VCHTYPE1>Purchase</VCHTYPE1>
                        <CSTFORMTYPE />
                        <CONSIGNEEACCOUNT />

                        <!-- Ledger Entries (Purchase Account) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Purchase Account - Medicines</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
                            <ISLASTDEEMEDPOSITIVE>No</ISLASTDEEMEDPOSITIVE>
                            <AMOUNT>50000</AMOUNT>
                            <SERVICETAXDETAILS />
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Inventory Entries -->
                        <ALLINVENTORYENTRIES.LIST>
                            <STOCKITEMNAME>Paracetamol 500mg - Tablet</STOCKITEMNAME>
                            <RATE>50</RATE>
                            <AMOUNT>25000</AMOUNT>
                            <ACTUALQTY>500 nos</ACTUALQTY>
                            <BATCHALLOCATIONS.LIST>
                                <BATCHNAME>BATCH-PAR-2026-001</BATCHNAME>
                                <GODOWNNAME>Main Location</GODOWNNAME>
                                <INDENTNO />
                                <ORDERNO />
                                <ORDEREDQUANTITY>500</ORDEREDQUANTITY>
                                <DESTINATIONGODOWN />
                                <BATCHDATE>20260219</BATCHDATE>
                                <BASEPARTYNAME />
                                <SERIALNOTRACKING>No</SERIALNOTRACKING>
                                <NARRATION />
                            </BATCHALLOCATIONS.LIST>
                            <ACCOUNTINGALLOCATIONS.LIST>
                                <LEDGERFROM />
                                <LEDGERTO />
                                <GODOWNNAME />
                                <QUANTITY>500</QUANTITY>
                                <RATE>50</RATE>
                                <PER>100</PER>
                                <AMOUNT>25000</AMOUNT>
                            </ACCOUNTINGALLOCATIONS.LIST>
                        </ALLINVENTORYENTRIES.LIST>

                        <ALLINVENTORYENTRIES.LIST>
                            <STOCKITEMNAME>Amoxicillin 250mg - Capsule</STOCKITEMNAME>
                            <RATE>25</RATE>
                            <AMOUNT>25000</AMOUNT>
                            <ACTUALQTY>1000 nos</ACTUALQTY>
                            <BATCHALLOCATIONS.LIST>
                                <BATCHNAME>BATCH-AMX-2026-001</BATCHNAME>
                                <GODOWNNAME>Main Location</GODOWNNAME>
                                <BATCHDATE>20260219</BATCHDATE>
                            </BATCHALLOCATIONS.LIST>
                        </ALLINVENTORYENTRIES.LIST>

                        <REMOTEID>medisync-pur-pharma-20260219</REMOTEID>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Sales Voucher (Sales)

**Use Case**: Sale of goods (pharmacy retail sales).

### Complete Example

```xml
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Sales" ACTION="Create" OBJVIEW="Inventory Voucher View">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>SALE-2026-0219-001</VOUCHERNUMBER>
                        <PARTYLEDGERNAME>Cash Sales - Pharmacy</PARTYLEDGERNAME>

                        <!-- Ledger Entries (Sales) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Sales Account - Pharmacy</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <ISLASTDEEMEDPOSITIVE>No</ISLASTDEEMEDPOSITIVE>
                            <AMOUNT>530</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Tax (VAT 5%) -->
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>VAT Output 5%</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <ISLASTDEEMEDPOSITIVE>No</ISLASTDEEMEDPOSITIVE>
                            <AMOUNT>25.24</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>

                        <!-- Inventory Entries -->
                        <ALLINVENTORYENTRIES.LIST>
                            <STOCKITEMNAME>Paracetamol 500mg - Strip</STOCKITEMNAME>
                            <RATE>10</RATE>
                            <AMOUNT>100</AMOUNT>
                            <ACTUALQTY>10 nos</ACTUALQTY>
                            <TRACKINGNUMBER />
                        </ALLINVENTORYENTRIES.LIST>

                        <ALLINVENTORYENTRIES.LIST>
                            <STOCKITEMNAME>Vitamin C 500mg - Bottle</STOCKITEMNAME>
                            <RATE>43</RATE>
                            <AMOUNT>430</AMOUNT>
                            <ACTUALQTY>10 nos</ACTUALQTY>
                        </ALLINVENTORYENTRIES.LIST>

                        <REMOTEID>medisync-sale-pharma-20260219</REMOTEID>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Credit Note Voucher

**Use Case**: Sales returns, credit to customers.

### Example

```xml
<VOUCHER VCHTYPE="Credit Note" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>CN-2026-0219-001</VOUCHERNUMBER>

    <!-- Customer Account (Debit - Liability Returns) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Patient - Fatima Hassan</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>500</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <!-- Sales Account (Credit - Sales Reduced) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Sales Account - Pharmacy</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>500</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <NARRATION>Medicine return - expired item</NARRATION>
</VOUCHER>
```

## Debit Note Voucher

**Use Case**: Purchase returns, debit from suppliers.

### Example

```xml
<VOUCHER VCHTYPE="Debit Note" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>DN-2026-0219-001</VOUCHERNUMBER>

    <!-- Supplier Account (Credit - Liability Reduced) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Pharma Supplies LLC</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>1000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <!-- Purchase Account (Debit - Purchase Reduced) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Purchase Account - Medicines</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>1000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <NARRATION>Goods return - damaged during transit</NARRATION>
</VOUCHER>
```

## Contra Voucher

**Use Case**: Cash deposits/withdrawals, bank transfers.

### Example

```xml
<VOUCHER VCHTYPE="Contra" ACTION="Create">
    <DATE>20260219</DATE>
    <VOUCHERNUMBER>CNT-2026-0219-001</VOUCHERNUMBER>

    <!-- Cash (Credit) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>Cash</LEDGERNAME>
        <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
        <AMOUNT>10000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <!-- Bank (Debit) -->
    <ALLLEDGERENTRIES.LIST>
        <LEDGERNAME>ADCB Bank - Current Account</LEDGERNAME>
        <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
        <AMOUNT>10000</AMOUNT>
    </ALLLEDGERENTRIES.LIST>

    <NARRATION>Cash deposit to bank</NARRATION>
</VOUCHER>
```

## Voucher Type Summary

| Voucher Type | Arabic Name | Typical Use |
|--------------|-------------|-------------|
| Journal | يومية | Adjustments, accruals, depreciation |
| Payment | دفع | Payments out (salaries, rent, expenses) |
| Receipt | قبض | Receipts in (patient payments) |
| Purchase | مشتريات | Buying inventory |
| Sales | مبيعات | Selling inventory |
| Credit Note | إشعار دائن | Sales returns |
| Debit Note | إشعار مدين | Purchase returns |
| Contra | مقابلة | Bank transfers, cash movements |

## Common Entry Patterns

### Debit/Credit Rules

```
For Payment Voucher:
├── Bank/Cash Account: ISDEEMEDPOSITIVE = No (Credit, money going out)
└── Expense Account: ISDEEMEDPOSITIVE = Yes (Debit, expense increases)

For Receipt Voucher:
├── Bank/Cash Account: ISDEEMEDPOSITIVE = Yes (Debit, money coming in)
└── Income/Patient Account: ISDEEMEDPOSITIVE = No (Credit, income increases)

For Journal Voucher:
├── Debit Accounts: ISDEEMEDPOSITIVE = No
└── Credit Accounts: ISDEEMEDPOSITIVE = Yes
```

### Double Entry Verification

Every voucher must balance:
- Sum of `ISDEEMEDPOSITIVE=Yes` amounts = Sum of `ISDEEMEDPOSITIVE=No` amounts
- In other words: Total Debits = Total Credits
