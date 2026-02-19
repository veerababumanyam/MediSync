# Tally Developer Reference (MediSync)

This document serves as the primary technical reference for integrating MediSync with TallyPrime.

## 1. Integration Architecture

TallyPrime acts as an HTTP server (Gateway) capable of receiving and responding to XML/JSON requests.

- **Transport**: HTTP POST
- **Default Port**: 9000
- **Endpoint**: `http://<Tally-IP>:9000`
- **Format**: XML (Native), JSON (Native/Non-Native)

### Gateway Configuration
To enable integration in TallyPrime:
1. **F1 (Help)** → **Settings** → **Advanced Configuration**.
2. Set **Enable HTTP Server** to **Yes**.
3. Set **Port** to `{desired_port}` (Default 9000).

---

## 2. XML Integration (Core)

XML is the most mature integration method. Every message is wrapped in an `<ENVELOPE>`.

### Base Structure
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
                    <!-- Data Payload (Voucher/Master) -->
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

### Key Header Elements
- `<TALLYREQUEST>`: `Import`, `Export`, or `Execute`.
- `<ID>`: Name of Report, Collection, or Function.
- `<VERSION>`: Version of message format (usually 1).

### Key Body Sections
- `<DESC>`: Description/Static variables for the request.
- `<DATA>`: The actual payload being transferred.

---

## 3. JSON Integration

TallyPrime Release 3.0+ supports JSON integration for modern web-style interactions.

- **Capabilities**: Fetching, modifying, and creating data objects.
- **Conversion**: Developers can use "Convert to JSON" TDL capabilities to export complex data structures.
- **REST-like**: Can be used for CRUD operations on Tally Objects.

---

## 4. TDL (Tally Definition Language) Basics

TDL is the language used to customize Tally. Understanding its hierarchy is crucial for deep integrations.

| Component | Responsibility |
| :--- | :--- |
| **Object** | The data container (e.g., Voucher, Ledger, Company). |
| **Collection** | A list of Objects (e.g., All Vouchers). |
| **Report** | Defines how data is viewed or formatted for output. |
| **UDF** | User Defined Fields—essential for storing extra platform metadata (like MediSync IDs). |

---

## 5. Implementation Best Practices

### A. Idempotency (Duplicate Prevention)
Tally lacks a native unique constraint on some fields.
- **Action**: Always generate and store a `RemoteID` (guid) for every transaction pushed.
- **Strategy**: Before pushing, check if a voucher with that `RemoteID` or a specific `Reference Number` already exists.

### B. Error Handling
Responses from Tally are also in XML.
- `<STATUS>1</STATUS>`: Success.
- `<STATUS>0</STATUS>`: Failure (check internal error tags for details).

### C. Performance
- Use **Collections** with filters to fetch only the data you need.
- Batch multiple `<TALLYMESSAGE>` items inside a single `<ENVELOPE>` for efficient imports.

---

## 6. Official Resources (Internal Links)
- [TDL Fundamentals](https://help.tallysolutions.com/developer-reference/tally-definition-language/)
- [XML Integration Guide](https://help.tallysolutions.com/xml-integration/)
- [Object & Collection Reference](https://help.tallysolutions.com/developer-reference/tally-definition-language/objects-and-collections/)
