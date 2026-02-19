# MediSync Testing Guide

Comprehensive testing strategies for the MediSync platform.

## Testing Pyramid

```
        ┌─────────┐
       /   E2E    \     ← 10% - Critical user flows
      /───────────\
     /             \
    /   Integration  \   ← 30% - Agent flows, API endpoints
   /─────────────────\
  /                     \
 /      Unit Tests       \ ← 60% - Individual functions, agents
/─────────────────────────\
```

## Unit Testing

### Go Unit Tests

Use the standard `testing` package with `testify/assert`:

```go
package module_a

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock LLM for testing
type MockLLM struct {
    mock.Mock
}

func (m *MockLLM) Generate(ctx context.Context, prompt string) (string, error) {
    args := m.Called(ctx, prompt)
    return args.String(0), args.Error(1)
}

func TestTextToSQLAgent_GenerateSQL(t *testing.T) {
    // Arrange
    mockLLM := new(MockLLM)
    agent := NewTextToSQLAgent(mockLLM)

    mockLLM.On("Generate", mock.Anything, mock.Anything).
        Return("SELECT * FROM patients WHERE city = 'Dubai'", nil)

    req := TextToSQLRequest{
        Query:  "Show me patients in Dubai",
        Locale: "en",
    }

    // Act
    resp, err := agent.GenerateSQL(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "SELECT * FROM patients WHERE city = 'Dubai'", resp.SQL)
    assert.Greater(t, resp.Confidence, 0.0)
    assert.LessOrEqual(t, resp.Confidence, 1.0)
    mockLLM.AssertExpectations(t)
}

func TestIsSelectOnlyQuery(t *testing.T) {
    tests := []struct {
        name     string
        sql      string
        expected bool
    }{
        {"valid select", "SELECT * FROM patients", true},
        {"select with join", "SELECT p.* FROM patients p JOIN doctors d", true},
        {"invalid insert", "INSERT INTO patients VALUES", false},
        {"invalid update", "UPDATE patients SET", false},
        {"invalid delete", "DELETE FROM patients", false},
        {"select lowercase", "select from patients", false},  // Must be uppercase
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isSelectOnlyQuery(tt.sql)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### React Unit Tests

Use Vitest with React Testing Library:

```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import QueryInput from './QueryInput';

const mocki18n = {
    t: (key: string) => key,
    language: 'en',
};

describe('QueryInput', () => {
    it('renders input placeholder', () => {
        render(
            <I18nextProvider i18n={mocki18n as any}>
                <QueryInput />
            </I18nextProvider>
        );

        expect(screen.getByPlaceholderText('agents.text2sql.placeholder')).toBeInTheDocument();
    });

    it('calls onSubmit with query text', async () => {
        const handleSubmit = vi.fn();

        render(
            <I18nextProvider i18n={mocki18n as any}>
                <QueryInput onSubmit={handleSubmit} />
            </I18nextProvider>
        );

        const input = screen.getByRole('textbox');
        const button = screen.getByRole('button', { name: /search/i });

        fireEvent.change(input, { target: { value: 'Show me patients' } });
        fireEvent.click(button);

        await waitFor(() => {
            expect(handleSubmit).toHaveBeenCalledWith('Show me patients');
        });
    });

    it('renders in RTL for Arabic locale', () => {
        const arI18n = { ...mocki18n, language: 'ar' };

        render(
            <I18nextProvider i18n={arI18n as any}>
                <QueryInput />
            </I18nextProvider>
        );

        expect(document.documentElement.dir).toBe('rtl');
    });
});
```

### Flutter Unit Tests

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:medisync/models/query_result.dart';

void main() {
  group('QueryResult', () {
    test('calculates confidence score correctly', () {
      final result = QueryResult(
        sql: 'SELECT * FROM patients',
        rows: [],
        confidence: 0.85,
      );

      expect(result.confidence, 0.85);
      expect(result.isHighConfidence, true);
    });

    test('serializes to JSON correctly', () {
      final result = QueryResult(
        sql: 'SELECT * FROM patients',
        rows: [{'id': 1, 'name': 'John'}],
        confidence: 0.9,
      );

      final json = result.toJson();

      expect(json['sql'], 'SELECT * FROM patients');
      expect(json['confidence'], 0.9);
    });
  });
}
```

## Integration Testing

### Agent Flow Integration Tests

Test agent flows with deterministic LLM responses (temperature=0):

```go
func TestAgentFlow_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup
    ctx := context.Background()
    config := LoadTestConfig()

    // Use real components with test database
    db := setupTestDB(t)
    defer db.Close()

    llm := setupTestLLM(t)  // Configured with temperature=0
    agent := NewTextToSQLAgent(llm, db)

    // Test cases with known inputs
    tests := []struct {
        name          string
        query         string
        expectedTable string
    }{
        {
            name:          "count patients",
            query:         "How many patients do we have?",
            expectedTable: "patients",
        },
        {
            name:          "list doctors",
            query:         "Show all doctors",
            expectedTable: "doctors",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := TextToSQLRequest{
                Query:  tt.query,
                Locale: "en",
            }

            resp, err := agent.GenerateSQL(ctx, req)

            assert.NoError(t, err)
            assert.Contains(t, strings.ToLower(resp.SQL), tt.expectedTable)
            assert.Contains(t, strings.ToUpper(resp.SQL), "SELECT")
        })
    }
}
```

### API Endpoint Tests

```go
func TestAPI_QueryEndpoint(t *testing.T) {
    // Setup test server
    router := setupTestRouter(t)

    // Create request
    body := `{"query": "Show me patients", "locale": "en"}`
    req := httptest.NewRequest("POST", "/api/query", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+testToken)

    // Record response
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var resp QueryResponse
    json.Unmarshal(w.Body.Bytes(), &resp)

    assert.NotEmpty(t, resp.SQL)
    assert.Greater(t, resp.Confidence, 0.0)
}
```

### ETL Pipeline Tests

```go
func TestETL_TallyImport(t *testing.T) {
    // Mock Tally gateway
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        xmlData := loadTestFile("tally-response.xml")
        w.Header().Set("Content-Type", "application/xml")
        w.Write([]byte(xmlData))
    }))
    defer server.Close()

    // Run ETL
    etl := NewTallyETL(server.URL, testDB)
    err := etl.Import(context.Background(), time.Now().Add(-24*time.Hour))

    assert.NoError(t, err)

    // Verify data in warehouse
    var count int
    testDB.QueryRow("SELECT COUNT(*) FROM vouchers WHERE source = 'tally'").Scan(&count)
    assert.Greater(t, count, 0)
}
```

## End-to-End Testing

### E2E Test Framework (Playwright)

```typescript
import { test, expect } from '@playwright/test';

test.describe('Conversational BI Flow', () => {
    test.beforeEach(async ({ page }) => {
        // Login
        await page.goto('http://localhost:3000');
        await page.fill('[name="email"]', 'test@example.com');
        await page.fill('[name="password"]', 'testpass');
        await page.click('button[type="submit"]');

        // Wait for dashboard
        await expect(page).toHaveURL(/.*dashboard/);
    });

    test('English query flow', async ({ page }) => {
        // Enter query
        await page.fill('[data-testid="query-input"]', 'How many patients do we have?');
        await page.click('[data-testid="submit-query"]');

        // Wait for results
        await expect(page.locator('[data-testid="query-result"]')).toBeVisible();

        // Verify chart is displayed
        await expect(page.locator('[data-testid="chart"]')).toBeVisible();

        // Verify confidence score
        const confidence = await page.textContent('[data-testid="confidence-score"]');
        expect(parseFloat(confidence)).toBeGreaterThan(0.5);
    });

    test('Arabic query flow', async ({ page }) => {
        // Switch to Arabic
        await page.click('[data-testid="lang-switcher"]');
        await page.click('button:has-text("العربية")');

        // Verify RTL
        expect(await page.evaluate(() => document.documentElement.dir)).toBe('rtl');

        // Enter Arabic query
        await page.fill('[data-testid="query-input"]', 'كم عدد المرضى لدينا؟');
        await page.click('[data-testid="submit-query"]');

        // Wait for results
        await expect(page.locator('[data-testid="query-result"]')).toBeVisible();
    });
});
```

### Document Flow E2E

```typescript
test.describe('AI Accountant Document Flow', () => {
    test('upload and sync invoice', async ({ page }) => {
        // Navigate to documents
        await page.click('a:has-text("Documents")');

        // Upload invoice
        const fileInput = await page.locator('input[type="file"]');
        await fileInput.setInputFiles('test-data/invoice.pdf');

        // Wait for OCR
        await expect(page.locator('[data-testid="ocr-result"]')).toBeVisible({ timeout: 10000 });

        // Verify extracted data
        await expect(page.locator('text=Invoice Amount:')).toBeVisible();

        // Submit for approval
        await page.click('[data-testid="submit-approval"]');

        // Wait for approval workflow
        await expect(page.locator('[data-testid="approval-pending"]')).toBeVisible();

        // Login as finance head and approve
        await page.click('[data-testid="logout"]');
        await loginAs(page, 'finance@example.com', 'finpass');

        // Approve the transaction
        await page.click('a:has-text("Approvals")');
        await page.click('[data-testid="approve-first"]');

        // Verify Tally sync initiated
        await expect(page.locator('[data-testid="sync-status"]')).toContainText('Syncing');
    });
});
```

## OPA Policy Testing

```go
func TestOPA_Policies(t *testing.T) {
    opa, err := NewOPAClient("http://localhost:8181")
    require.NoError(t, err)

    tests := []struct {
        name     string
        user     User
        action   string
        resource string
        allowed  bool
    }{
        {
            name: "finance head can approve",
            user: User{Roles: []string{"finance_head"}},
            action:   "approve",
            resource: "journal_entry",
            allowed:  true,
        },
        {
            name: "regular user cannot approve",
            user: User{Roles: []string{"user"}},
            action:   "approve",
            resource: "journal_entry",
            allowed:  false,
        },
        {
            name: "anyone can query dashboard",
            user: User{Roles: []string{"user"}},
            action:   "query",
            resource: "dashboard",
            allowed:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := opa.Allow(context.Background(), tt.action, map[string]interface{}{
                "user":     tt.user,
                "resource": tt.resource,
            })

            assert.NoError(t, err)
            assert.Equal(t, tt.allowed, result)
        })
    }
}
```

## Tally Integration Testing

### Mock Tally Server

```go
func TestTallySync_MockServer(t *testing.T) {
    // Start mock Tally gateway
    mockTally := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request format
        contentType := r.Header.Get("Content-Type")
        assert.Equal(t, "application/xml", contentType)

        // Parse request body
        body, _ := io.ReadAll(r.Body)
        assert.Contains(t, string(body), "<TALLYMESSAGE>")

        // Return success response
        w.Header().Set("Content-Type", "application/xml")
        w.Write([]byte(`<ENVELOPE><BODY><DATA><STATUS>1</STATUS></DATA></BODY></ENVELOPE>`))
    }))
    defer mockTally.Close()

    // Test sync
    sync := NewTallySync(mockTally.URL)
    entry := JournalEntry{
        ID:       uuid.New().String(),
        Date:     time.Now(),
        Entries:  []LedgerEntry{{...}},
    }

    err := sync.Push(context.Background(), entry)
    assert.NoError(t, err)
}
```

## Test Data Management

### Fixtures

```go
// testutil/fixtures.go
package testutil

func LoadTestFixtures(db *sql.DB) error {
    fixtures := []string{
        "testdata/companies.sql",
        "testdata/patients.sql",
        "testdata/vouchers.sql",
    }

    for _, f := range fixtures {
        data, err := os.ReadFile(f)
        if err != nil {
            return err
        }
        if _, err := db.Exec(string(data)); err != nil {
            return err
        }
    }
    return nil
}
```

### Factories

```go
// testutil/factory.go
package testutil

type UserFactory struct{}

func (f *UserFactory) Create(overrides ...func(*User)) User {
    user := User{
        ID:       uuid.New().String(),
        Email:    "test@example.com",
        Name:     "Test User",
        Role:     "user",
        Locale:   "en",
    }

    for _, override := range overrides {
        override(&user)
    }

    return user
}

// Usage:
user := UserFactory{}.Create(func(u *User) {
    u.Role = "finance_head"
    u.Locale = "ar"
})
```

## Running Tests

### Go Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests only
go test -tags=integration ./...

# Run specific test
go test -v ./internal/agents/module_a -run TestTextToSQL

# Race detection
go test -race ./...
```

### React Tests

```bash
# Run all tests
npm test

# Run in watch mode
npm test -- --watch

# Run with coverage
npm test -- --coverage

# Run specific file
npm test -- QueryInput.test.tsx
```

### Flutter Tests

```bash
# Run all tests
flutter test

# Run with coverage
flutter test --coverage

# Run specific test
flutter test test/models/query_result_test.dart
```

## CI/CD Integration

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:18.2
        env:
          POSTGRES_DB: medisync_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.26'

      - name: Run Go tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...

      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci

      - name: Run frontend tests
        working-directory: ./frontend
        run: npm test -- --coverage --run
```

## Test Checklist

Before committing code:

- [ ] Unit tests cover new functions
- [ ] Integration tests cover new endpoints/flows
- [ ] E2E tests cover critical user paths
- [ ] Tests pass in both English and Arabic
- [ ] OPA policy tests updated for new actions
- [ ] Mock tests for external dependencies (Tally, HIMS)
- [ ] Coverage above 80% for new code
- [ ] No race conditions detected
