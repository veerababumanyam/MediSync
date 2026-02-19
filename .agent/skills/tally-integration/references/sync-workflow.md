# Tally Sync Workflow Reference

Detailed workflow for synchronizing data from MediSync to Tally ERP.

## Sync Pipeline Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         MediSync Platform                           │
├─────────────────────────────────────────────────────────────────────┤
│  1. Data Entry (OCR / Manual / HIMS)                               │
│                    ↓                                                │
│  2. Draft Record Creation (Pending Approval)                        │
│                    ↓                                                │
│  3. Agent B-05: Ledger Mapping (AI-suggests Tally ledgers)         │
│                    ↓                                                │
│  4. Agent B-08: Approval Workflow (HITL Gate)                       │
│     ├─ Level 1: Creator Review                                     │
│     ├─ Level 2: Manager Approval (if amount > threshold)           │
│     └─ Level 3: Finance Head (for Tally sync)                      │
│                    ↓                                                │
│  5. OPA Policy Check                                               │
│     └─ Verify user has tally.sync permission                       │
│                    ↓                                                │
│  6. Generate TDL XML (Jinja2 Template)                             │
│                    ↓                                                │
│  7. POST to Tally Gateway (port 9000)                              │
│                    ↓                                                │
│  8. Parse Response & Verify                                        │
│                    ↓                                                │
│  9. Update Status & Audit Log                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Phase 1: Draft Record Creation

### Data Structure

```go
type JournalEntry struct {
    ID           string            `json:"id"`
    CompanyID    string            `json:"company_id"`
    Date         time.Time         `json:"date"`
    VoucherType  string            `json:"voucher_type"` // Journal, Payment, etc.
    Number       string            `json:"number"`
    Narration    string            `json:"narration"`
    Entries      []LedgerEntry     `json:"entries"`
    TotalAmount  float64           `json:"total_amount"`
    Currency     string            `json:"currency"`
    Status       string            `json:"status"` // draft, pending_approval, approved, synced, failed
    CreatedBy    string            `json:"created_by"`
    CreatedAt    time.Time         `json:"created_at"`
    UpdatedAt    time.Time         `json:"updated_at"`
}

type LedgerEntry struct {
    LedgerID     string  `json:"ledger_id"`      // MediSync ledger reference
    LedgerName   string  `json:"ledger_name"`    // Mapped Tally ledger name
    IsDebit      bool    `json:"is_debit"`
    Amount       float64 `json:"amount"`
    CostCenter   string  `json:"cost_center,omitempty"`
}
```

### Creating a Draft

```go
func (s *Service) CreateDraftEntry(ctx context.Context, req CreateEntryRequest) (*JournalEntry, error) {
    entry := &JournalEntry{
        ID:          uuid.New().String(),
        CompanyID:   req.CompanyID,
        Date:        req.Date,
        VoucherType: req.VoucherType,
        Number:      s.generateVoucherNumber(ctx, req.VoucherType, req.Date),
        Entries:     req.Entries,
        Status:      "draft",
        CreatedBy:   req.UserID,
        CreatedAt:   time.Now(),
    }

    // Validate debits = credits
    if !s.isBalanced(entry.Entries) {
        return nil, fmt.Errorf("entry is not balanced: debits != credits")
    }

    // Save to database
    if err := s.repo.Create(ctx, entry); err != nil {
        return nil, err
    }

    return entry, nil
}
```

## Phase 2: Ledger Mapping (Agent B-05)

### AI Suggestion

```go
func (a *LedgerMappingAgent) SuggestLedger(ctx context.Context, description string, locale string) (*LedgerSuggestion, error) {
    // Build prompt with Tally ledger catalog
    catalog := a.getLedgerCatalog(ctx, description)

    prompt := fmt.Sprintf(`
Given this expense description: "%s"

Suggest the most appropriate Tally ledger from this catalog:
%s

Return the ledger name and confidence score.
`, description, catalog)

    resp, err := a.llm.Generate(ctx, prompt)
    if err != nil {
        return nil, err
    }

    return parseSuggestion(resp)
}
```

### User Confirmation

```go
type LedgerMappingConfirmation struct {
    EntryID       string            `json:"entry_id"`
    EntryIndex    int               `json:"entry_index"`
    Suggested     string            `json:"suggested_ledger"`
    Confirmed     string            `json:"confirmed_ledger"`
    Confidence    float64           `json:"confidence"`
    Overridden    bool              `json:"overridden"`
}
```

## Phase 3: Approval Workflow (Agent B-08)

### Approval Chain

```go
type ApprovalChain struct {
    EntryID        string         `json:"entry_id"`
    CurrentLevel   int            `json:"current_level"`
    RequiredLevels int            `json:"required_levels"`
    Approvals      []Approval     `json:"approvals"`
    Status         string         `json:"status"` // pending, approved, rejected
}

type Approval struct {
    Level      int       `json:"level"`
    UserID     string    `json:"user_id"`
    Role       string    `json:"role"`
    Decision   string    `json:"decision"` // approved, rejected
    Comment    string    `json:"comment"`
    Timestamp  time.Time `json:"timestamp"`
}
```

### Approval Logic

```go
func (a *ApprovalAgent) ProcessEntry(ctx context.Context, entryID string) (*ApprovalResult, error) {
    entry, err := a.repo.Get(ctx, entryID)
    if err != nil {
        return nil, err
    }

    // Determine required approval levels
    requiredLevels := a.getRequiredLevels(entry.TotalAmount)

    chain := &ApprovalChain{
        EntryID:        entryID,
        RequiredLevels: requiredLevels,
        CurrentLevel:   0,
        Status:         "pending",
    }

    // Check if auto-approve (small amounts)
    if entry.TotalAmount < a.autoApproveThreshold {
        chain.Status = "approved"
        // Auto-approve on behalf of system
        chain.Approvals = []Approval{{
            Level:     1,
            UserID:    "system",
            Decision:  "approved",
            Comment:   "Auto-approved: amount below threshold",
            Timestamp: time.Now(),
        }}
        return a.finalizeApproval(ctx, chain)
    }

    // Route to human approvers
    return a.routeForApproval(ctx, entry, chain)
}

func (a *ApprovalAgent) getRequiredLevels(amount float64) int {
    switch {
    case amount < 1000:
        return 1 // Manager only
    case amount < 10000:
        return 2 // Manager + Finance Manager
    default:
        return 3 // Manager + Finance Manager + Finance Head
    }
}
```

## Phase 4: OPA Policy Check

### Authorization

```go
func (s *SyncService) checkOPA(ctx context.Context, userID, entryID string, action string) error {
    input := map[string]interface{}{
        "user_id":   userID,
        "entry_id":  entryID,
        "action":    action,
        "timestamp": time.Now().Unix(),
    }

    allowed, err := s.opa.Allow(ctx, "tally_sync", input)
    if err != nil {
        return fmt.Errorf("OPA check failed: %w", err)
    }

    if !allowed {
        return fmt.Errorf("user %s not authorized for %s", userID, action)
    }

    return nil
}
```

### Example OPA Policy (Rego)

```rego
package medisync.tally

default allow = false

allow {
    input.action == "sync_journal_entry"
    has_role(input.user_roles, "finance_head")
    entry_amount_within_limit(input.entry_id, input.user_limits)
}

has_role(roles, role) {
    role in roles
}

entry_amount_within_limit(entry_id, limits) {
    amount := data.entries[entry_id].amount
    amount <= limits.max_sync_amount
}
```

## Phase 5: Generate TDL XML

### Jinja2 Template

```python
from jinja2 import Template

TEMPLATE = """<ENVELOPE>
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
                        <NARRATION>{{ narration }}</NARRATION>
                        <REMOTEID>{{ remote_id }}</REMOTEID>
                        {% for entry in ledger_entries %}
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>{{ entry.ledger_name }}</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>{{ 'Yes' if entry.is_debit else 'No' }}</ISDEEMEDPOSITIVE>
                            <AMOUNT>{{ entry.amount|round(2) }}</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                        {% endfor %}
                        <UDF:MSYNC.ENTRYID>{{ entry_id }}</UDF:MSYNC.ENTRYID>
                        <UDF:MSYNC.SYNCEDBY>{{ synced_by }}</UDF:MSYNC.SYNCEDBY>
                        <UDF:MSYNC.SYNCDATETIME>{{ sync_datetime|strftime('%Y-%m-%d %H:%M:%S') }}</UDF:MSYNC.SYNCDATETIME>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
"""

def generate_xml(entry: JournalEntry) -> str:
    template = Template(TEMPLATE)
    return template.render(
        voucher_type=entry.VoucherType,
        date=entry.Date,
        voucher_number=entry.Number,
        narration=entry.Narration,
        remote_id=generate_remote_id(entry),
        ledger_entries=entry.Entries,
        entry_id=entry.ID,
        synced_by=entry.UpdatedBy,
        sync_datetime=time.now()
    )
```

### Go Template Alternative

```go
const tallyXMLTemplate = `<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="{{.VoucherType}}" ACTION="Create">
                        <DATE>{{.Date}}</DATE>
                        <VOUCHERNUMBER>{{.Number}}</VOUCHERNUMBER>
                        <REMOTEID>{{.RemoteID}}</REMOTEID>
                        {{range .Entries}}
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>{{.LedgerName}}</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>{{if .IsDebit}}Yes{{else}}No{{end}}</ISDEEMEDPOSITIVE>
                            <AMOUNT>{{.Amount}}</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                        {{end}}
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
`

func (s *SyncService) generateXML(entry *JournalEntry) (string, error) {
    tmpl, err := template.New("tally").Parse(tallyXMLTemplate)
    if err != nil {
        return "", err
    }

    data := struct {
        VoucherType string
        Date        string // YYYYMMDD format
        Number      string
        RemoteID    string
        Entries     []LedgerEntry
    }{
        VoucherType: entry.VoucherType,
        Date:        entry.Date.Format("20060102"),
        Number:      entry.Number,
        RemoteID:    s.generateRemoteID(entry),
        Entries:     entry.Entries,
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}
```

## Phase 6: POST to Tally Gateway

### HTTP Client

```go
type TallyClient struct {
    client    *http.Client
    gatewayURL string
    timeout   time.Duration
}

func NewTallyClient(gatewayURL string) *TallyClient {
    return &TallyClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
        gatewayURL: gatewayURL,
    }
}

func (c *TallyClient) Send(ctx context.Context, xml string) (*TallyResponse, error) {
    req, err := http.NewRequestWithContext(ctx, "POST", c.gatewayURL, strings.NewReader(xml))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/xml")

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("tally gateway error: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return parseTallyResponse(body)
}
```

### Retry Logic

```go
func (c *TallyClient) SendWithRetry(ctx context.Context, xml string, maxRetries int) (*TallyResponse, error) {
    var lastErr error

    for attempt := 0; attempt < maxRetries; attempt++ {
        // Exponential backoff
        if attempt > 0 {
            waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(waitTime):
            }
        }

        resp, err := c.Send(ctx, xml)
        if err == nil {
            return resp, nil
        }

        lastErr = err

        // Don't retry on validation errors
        if strings.Contains(err.Error(), "Ledger") && strings.Contains(err.Error(), "does not exist") {
            return nil, err
        }
    }

    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Phase 7: Parse Response

```go
type TallyResponse struct {
    Status     int    `xml:"HEADER>STATUS"`
    LineError  int    `xml:"BODY>DATA>LINEERROR"`
    ErrorLine  int    `xml:"BODY>DATA>ERRORLINE"`
    ErrorMessage string `xml:"BODY>DATA>ERRORMESSAGE"`
    Collection string `xml:"BODY>DATA>COLLECTION>TYPE"`
}

func parseTallyResponse(data []byte) (*TallyResponse, error) {
    var resp TallyResponse
    if err := xml.Unmarshal(data, &resp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Tally uses STATUS=1 for success, 0 for failure
    if resp.Status != 1 {
        return &resp, fmt.Errorf("tally error: %s (line %d)", resp.ErrorMessage, resp.ErrorLine)
    }

    return &resp, nil
}
```

## Phase 8: Update Status & Audit Log

```go
func (s *SyncService) CompleteSync(ctx context.Context, entryID string, resp *TallyResponse) error {
    // Update entry status
    status := "synced"
    if resp.Status != 1 {
        status = "failed"
    }

    err := s.repo.UpdateStatus(ctx, entryID, status)
    if err != nil {
        return err
    }

    // Create audit log
    audit := &AuditLog{
        ID:        uuid.New().String(),
        UserID:    s.userID,
        Action:    "tally_sync",
        Resource:  entryID,
        Status:    status,
        Details:   resp.ErrorMessage,
        CreatedAt: time.Now(),
    }

    return s.auditRepo.Create(ctx, audit)
}
```

## Error Handling Matrix

| Error Type | Retry? | Action |
|------------|--------|--------|
| Connection timeout | Yes (3x) | Exponential backoff |
| Ledger not found | No | Flag for manual fix |
| Duplicate voucher | No | Skip or renumber |
| Malformed XML | No | Fix template |
| 500 Internal Error | Yes (3x) | Exponential backoff |
| 401 Unauthorized | No | Check credentials |

## Idempotency

### RemoteID Generation

```go
func generateRemoteID(entry *JournalEntry) string {
    data := fmt.Sprintf("%s|%s|%s|%.2f",
        entry.CompanyID,
        entry.Date.Format("2006-01-02"),
        entry.VoucherType,
        entry.TotalAmount,
    )

    // Add entry details for uniqueness
    for _, e := range entry.Entries {
        data += fmt.Sprintf("|%s:%.2f", e.LedgerName, e.Amount)
    }

    h := sha256.Sum256([]byte(data))
    return fmt.Sprintf("medisync-%x", h[:12])
}
```

### Duplicate Check Before Sync

```go
func (s *SyncService) CheckExists(ctx context.Context, remoteID string) (bool, error) {
    // Query Tally for existing voucher with same RemoteID
    query := fmt.Sprintf(`
        <ENVELOPE>
            <HEADER>
                <TALLYREQUEST>Export Data</TALLYREQUEST>
            </HEADER>
            <BODY>
                <EXPORTDATA>
                    <REQUESTDESC>
                        <REPORTNAME>Voucher Register</REPORTNAME>
                        <STATICVARIABLES>
                            <SVFROMDATE>20260101</SVFROMDATE>
                            <SVtodate>20261231</SVtodate>
                        </STATICVARIABLES>
                    </REQUESTDESC>
                </EXPORTDATA>
            </BODY>
        </ENVELOPE>
    `)

    resp, err := s.tally.Send(ctx, query)
    if err != nil {
        return false, err
    }

    return strings.Contains(resp.RawData, remoteID), nil
}
```

## Monitoring

### Sync Metrics

```go
type SyncMetrics struct {
    TotalAttempts   int64
    SuccessfulSyncs int64
    FailedSyncs     int64
    AvgLatency      time.Duration
   LastError       string
    LastSyncTime    time.Time
}

func (s *SyncService) RecordAttempt(success bool, latency time.Duration) {
    s.metrics.TotalAttempts++
    if success {
        s.metrics.SuccessfulSyncs++
    } else {
        s.metrics.FailedSyncs++
    }

    // Update moving average
    s.metrics.AvgLatency = (s.metrics.AvgLatency + latency) / 2
    s.metrics.LastSyncTime = time.Now()
}
```
