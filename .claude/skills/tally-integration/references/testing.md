# Tally Integration Testing Guide

Testing strategies for Tally ERP integration.

## Testing Pyramid

```
        ┌─────────────┐
       /   E2E Test   \     ← Tally sandbox environment
      /────────────────\
     /   Integration    \   ← Mock Tally server
    /────────────────────\
   /     Unit Tests       \ ← Individual functions
  /────────────────────────\
```

## Unit Testing

### Testing TDL Generation

```go
func TestGenerateTDLXML(t *testing.T) {
    entry := &JournalEntry{
        ID:          "test-001",
        Date:        time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
        VoucherType: "Journal",
        Number:      "JV-001",
        Narration:   "Test entry",
        Entries: []LedgerEntry{
            {LedgerName: "Rent Expense", IsDebit: true, Amount: 1000},
            {LedgerName: "Cash", IsDebit: false, Amount: 1000},
        },
    }

    generator := NewTDLGenerator()
    xml, err := generator.Generate(entry)

    require.NoError(t, err)
    assert.Contains(t, xml, `<VOUCHER VCHTYPE="Journal"`)
    assert.Contains(t, xml, `<DATE>20260219</DATE>`)
    assert.Contains(t, xml, `<LEDGERNAME>Rent Expense</LEDGERNAME>`)
    assert.Contains(t, xml, `<ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>`)  // Debit
    assert.Contains(t, xml, `<ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>`) // Credit
}

func TestRemoteIDGeneration(t *testing.T) {
    entry := &JournalEntry{
        CompanyID:   "company-123",
        Date:        time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
        VoucherType: "Journal",
        TotalAmount: 5000,
    }

    id1 := generateRemoteID(entry)
    id2 := generateRemoteID(entry)

    // Same input should generate same ID
    assert.Equal(t, id1, id2)
    assert.HasPrefix(t, id1, "medisync-")
    assert.Len(t, id1, 12+8)  // "medisync-" + 12 char hash
}

func TestIsDeemedPositive(t *testing.T) {
    tests := []struct {
        name     string
        isDebit  bool
        expected string
    }{
        {"debit entry", true, "No"},
        {"credit entry", false, "Yes"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isDeemedPositive(tt.isDebit)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Testing Response Parsing

```go
func TestParseTallyResponse(t *testing.T) {
    tests := []struct {
        name      string
        xml       string
        expectErr bool
        status    int
        errorMsg  string
    }{
        {
            name: "success response",
            xml: `<ENVELOPE>
                <HEADER>
                    <VERSION>1</VERSION>
                    <STATUS>1</STATUS>
                </HEADER>
                <BODY>
                    <DATA>
                        <LINEERROR>0</LINEERROR>
                    </DATA>
                </BODY>
            </ENVELOPE>`,
            expectErr: false,
            status:    1,
        },
        {
            name: "error response",
            xml: `<ENVELOPE>
                <HEADER>
                    <STATUS>0</STATUS>
                </HEADER>
                <BODY>
                    <DATA>
                        <LINEERROR>1</LINEERROR>
                        <ERRORMESSAGE>Ledger does not exist</ERRORMESSAGE>
                    </DATA>
                </BODY>
            </ENVELOPE>`,
            expectErr: true,
            status:    0,
            errorMsg:  "Ledger does not exist",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := parseTallyResponse([]byte(tt.xml))

            if tt.expectErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.status, resp.Status)
            }
        })
    }
}
```

## Integration Testing

### Mock Tally Server

```go
// mock_tally_server.go
package testing

import (
    "encoding/xml"
    "net/http"
    "net/http/httptest"
    "strings"
)

type MockTallyServer struct {
    server *httptest.Server
}

type TallyResponseConfig struct {
    Status       int
    LineError    int
    ErrorMessage string
    Delay        time.Duration
}

func NewMockTallyServer(config TallyResponseConfig) *MockTallyServer {
    m := &MockTallyServer{}

    m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Simulate network delay
        if config.Delay > 0 {
            time.Sleep(config.Delay)
        }

        // Verify request format
        contentType := r.Header.Get("Content-Type")
        if contentType != "application/xml" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Read request for logging
        body, _ := io.ReadAll(r.Body)

        // Log received XML
        t.Logf("Received Tally request:\n%s", string(body))

        // Build response
        resp := map[string]interface{}{
            "HEADER": map[string]interface{}{
                "VERSION": "1",
                "STATUS":  config.Status,
            },
            "BODY": map[string]interface{}{
                "DATA": map[string]interface{}{
                    "LINEERROR":     config.LineError,
                    "ERRORMESSAGE":  config.ErrorMessage,
                },
            },
        }

        xmlData, _ := xml.Marshal(resp)
        xmlEnvelope := fmt.Sprintf(`<ENVELOPE>%s</ENVELOPE>`, xmlData)

        w.Header().Set("Content-Type", "application/xml")
        w.Write([]byte(xmlEnvelope))
    }))

    return m
}

func (m *MockTallyServer) URL() string {
    return m.server.URL
}

func (m *MockTallyServer) Close() {
    m.server.Close()
}
```

### Integration Test Example

```go
func TestTallySync_Integration(t *testing.T) {
    // Setup mock Tally
    mockTally := NewMockTallyServer(TallyResponseConfig{
        Status:    1,
        LineError: 0,
    })
    defer mockTally.Close()

    // Create sync service with mock
    sync := NewTallySyncService(mockTally.URL())

    // Create test entry
    entry := &JournalEntry{
        ID:          uuid.New().String(),
        Date:        time.Now(),
        VoucherType: "Journal",
        Number:      fmt.Sprintf("JV-%d", time.Now().Unix()),
        Entries: []LedgerEntry{
            {LedgerName: "Test Debit", IsDebit: true, Amount: 100},
            {LedgerName: "Test Credit", IsDebit: false, Amount: 100},
        },
    }

    // Execute sync
    err := sync.Sync(context.Background(), entry)

    // Assert
    assert.NoError(t, err)

    // Verify status in database
    updated, err := repo.Get(context.Background(), entry.ID)
    assert.NoError(t, err)
    assert.Equal(t, "synced", updated.Status)
}

func TestTallySync_LedgerNotFound(t *testing.T) {
    // Mock Tally returns ledger error
    mockTally := NewMockTallyServer(TallyResponseConfig{
        Status:       0,
        LineError:    1,
        ErrorMessage: "Ledger 'Unknown Ledger' does not exist",
    })
    defer mockTally.Close()

    sync := NewTallySyncService(mockTally.URL())

    entry := &JournalEntry{
        ID:          uuid.New().String(),
        Date:        time.Now(),
        VoucherType: "Journal",
        Entries: []LedgerEntry{
            {LedgerName: "Unknown Ledger", IsDebit: true, Amount: 100},
            {LedgerName: "Cash", IsDebit: false, Amount: 100},
        },
    }

    err := sync.Sync(context.Background(), entry)

    // Should return error with Tally message
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "does not exist")

    // Status should be failed
    updated, _ := repo.Get(context.Background(), entry.ID)
    assert.Equal(t, "failed", updated.Status)
}
```

### Testing Retry Logic

```go
func TestTallySync_RetryOnTimeout(t *testing.T) {
    attemptCount := 0

    mockTally := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attemptCount++

        // First two attempts timeout, third succeeds
        if attemptCount < 3 {
            time.Sleep(2 * time.Second)
            return
        }

        // Success response
        w.Write([]byte(`<ENVELOPE><HEADER><STATUS>1</STATUS></HEADER></ENVELOPE>`))
    }))
    defer mockTally.Close()

    // Create client with short timeout
    sync := NewTallySyncService(mockTally.URL())
    sync.client.Timeout = 500 * time.Millisecond
    sync.maxRetries = 3

    entry := createTestEntry()

    // Should succeed after retries
    err := sync.Sync(context.Background(), entry)

    assert.NoError(t, err)
    assert.Equal(t, 3, attemptCount)
}
```

## OPA Policy Testing

### Unit Test OPA Policies

```rego
# policies/tally_sync_test.rego
package medisync.tally

test_allow_sync_with_all_approvals {
    allow with input as {
        "action": "sync_to_tally",
        "entry_id": "entry-001",
        "entry": {
            "total_amount": 5000,
            "approvals_completed": 2,
            "company_id": "company-123",
        },
        "user_roles": {"permissions": ["tally.sync"]},
        "user": {"max_sync_amount": 10000},
    }
}

test_deny_sync_insufficient_approvals {
    not allow with input as {
        "action": "sync_to_tally",
        "entry": {
            "total_amount": 15000,  # Requires 2 approvals
            "approvals_completed": 1,  # Only 1 given
        },
    }
}

test_deny_sync_over_user_limit {
    not allow with input as {
        "action": "sync_to_tally",
        "entry": {"total_amount": 50000},
        "user": {"max_sync_amount": 10000},
    }
}

test_auto_approve_small_amount {
    get_required_levels(500) == 0
}

test_single_approval_medium_amount {
    get_required_levels(5000) == 1
}

test_dual_approval_large_amount {
    get_required_levels(25000) == 2
}

test_tri_approval_very_large_amount {
    get_required_levels(100000) == 3
}
```

### Go OPA Test

```go
func TestOPA_TallySync(t *testing.T) {
    // Setup OPA
    ctx := context.Background()
    opa, err := opa.NewClient(opa.Config{
        URL:    "http://localhost:8181",
        Policies: []string{"policies/rego/tally_sync.rego"},
    })
    require.NoError(t, err)

    tests := []struct {
        name     string
        input    map[string]interface{}
        expected bool
    }{
        {
            name: "allow with proper approvals",
            input: map[string]interface{}{
                "action": "sync_to_tally",
                "entry": map[string]interface{}{
                    "total_amount":        5000,
                    "approvals_completed": 2,
                },
                "user_roles": map[string]interface{}{
                    "permissions": []string{"tally.sync"},
                },
            },
            expected: true,
        },
        {
            name: "deny without approvals",
            input: map[string]interface{}{
                "action": "sync_to_tally",
                "entry": map[string]interface{}{
                    "total_amount":        15000,
                    "approvals_completed": 0,
                },
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := opa.Decision(ctx, opa.DecisionOptions{
                Path:  "/medisync/tally/allow",
                Input: tt.input,
            })

            require.NoError(t, err)
            assert.Equal(t, tt.expected, result.Result)
        })
    }
}
```

## End-to-End Testing

### Test Data Setup

```sql
-- test/fixtures/tally_test_data.sql
INSERT INTO companies (id, name, tally_company_name) VALUES
('test-company-1', 'Test Clinic', 'Test Clinic LLC');

INSERT INTO ledgers (id, company_id, medisync_name, tally_name) VALUES
('ledger-1', 'test-company-1', 'Cash', 'Cash'),
('ledger-2', 'test-company-1', 'Rent Expense', 'Rent Expense'),
('ledger-3', 'test-company-1', 'Sales Revenue', 'Sales Revenue');

INSERT INTO users (id, email, role, company_id) VALUES
('user-1', 'test@example.com', 'manager', 'test-company-1'),
('user-2', 'finance@example.com', 'finance_head', 'test-company-1');
```

### E2E Test Flow

```go
func TestE2E_DocumentToTallySync(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping E2E test")
    }

    ctx := context.Background()

    // 1. Upload invoice document
    doc := uploadTestDocument(ctx, t, "test-data/invoice.pdf")

    // 2. Wait for OCR processing
    extracted := waitForOCR(ctx, t, doc.ID)

    // 3. Verify extracted data
    assert.Equal(t, "Purchase Invoice", extracted.Type)
    assert.NotEmpty(t, extracted.VendorName)
    assert.Greater(t, extracted.TotalAmount, 0.0)

    // 4. Create draft entry
    entry := createDraftEntry(ctx, t, extracted)
    assert.Equal(t, "draft", entry.Status)

    // 5. Submit for approval
    request := submitForApproval(ctx, t, entry.ID)
    assert.Equal(t, "pending", request.Status)

    // 6. Approve as manager
    approveAs(ctx, t, "user-1", request.ID, "Looks good")

    // 7. Approve as finance head (amount > threshold)
    approveAs(ctx, t, "user-2", request.ID, "Approved")

    // 8. Wait for Tally sync
    syncedEntry := waitForSync(ctx, t, entry.ID)
    assert.Equal(t, "synced", syncedEntry.Status)

    // 9. Verify audit log
    logs := getAuditLogs(ctx, t, entry.ID)
    assert.Len(t, logs, 4)  // create, submit, 2 approvals
}
```

## Test Fixtures

### Sample TDL Files

```xml
<!-- test/fixtures/success_response.xml -->
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

```xml
<!-- test/fixtures/ledger_error_response.xml -->
<ENVELOPE>
    <HEADER>
        <VERSION>1</VERSION>
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

### Sample Voucher

```xml
<!-- test/fixtures/sample_journal_voucher.xml -->
<ENVELOPE>
    <HEADER>
        <TALLYREQUEST>Import Data</TALLYREQUEST>
    </HEADER>
    <BODY>
        <IMPORTDATA>
            <REQUESTDATA>
                <TALLYMESSAGE xmlns:UDF="TallyUDF">
                    <VOUCHER VCHTYPE="Journal" ACTION="Create">
                        <DATE>20260219</DATE>
                        <VOUCHERNUMBER>TEST-001</VOUCHERNUMBER>
                        <NARRATION>Test journal entry</NARRATION>
                        <REMOTEID>medisync-test-001</REMOTEID>
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Cash</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
                            <AMOUNT>100</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                        <ALLLEDGERENTRIES.LIST>
                            <LEDGERNAME>Rent Expense</LEDGERNAME>
                            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
                            <AMOUNT>100</AMOUNT>
                        </ALLLEDGERENTRIES.LIST>
                    </VOUCHER>
                </TALLYMESSAGE>
            </REQUESTDATA>
        </IMPORTDATA>
    </BODY>
</ENVELOPE>
```

## Running Tests

```bash
# Unit tests only
go test -short ./internal/tally/...

# Integration tests (requires mock server)
go test -run Integration ./internal/tally/...

# E2E tests (requires Tally sandbox)
go test -run E2E ./internal/tally/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# OPA policy tests
opa test policies/rego/tally_sync.rego -v

# Run all tests with verbose output
go test -v ./...
```

## CI/CD Integration

```yaml
# .github/workflows/tally-integration-tests.yml
name: Tally Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      opa:
        image: openpolicyagent/opa:latest
        ports:
          - 8181:8181

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.26'

      - name: Load OPA policies
        run: |
          opa eval -d policies/rego -i input.json "data.medisync.tally.allow"

      - name: Run unit tests
        run: go test -short -v ./internal/tally/...

      - name: Run integration tests
        run: go test -run Integration -v ./internal/tally/...
```

## Testing Checklist

- [ ] Unit tests cover TDL generation
- [ ] Unit tests cover response parsing
- [ ] Mock server tests sync flow
- [ ] Retry logic tested
- [ ] OPA policies tested
- [ ] Approval workflow tested
- [ ] Idempotency tested
- [ ] Error handling tested
- [ ] E2E test in sandbox environment
- [ ] Audit logging verified
