// Package tally provides tests for the Tally TDL XML client.
package tally

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/medisync/medisync/internal/config"
)

// testConfig returns a TallyConfig for testing.
func testConfig() config.TallyConfig {
	return config.TallyConfig{
		Host:       "localhost",
		Port:       9000,
		Company:    "Test Company",
		Timeout:    5 * time.Second,
		MaxRetries: 2,
		RetryDelay: 100 * time.Millisecond,
	}
}

// mockTallyServer creates a test HTTP server that simulates Tally responses.
func mockTallyServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestNewClient(t *testing.T) {
	cfg := testConfig()
	client := NewClient(cfg)

	if client == nil {
		t.Fatal("expected non-nil client")
	}

	expectedURL := "http://localhost:9000"
	if client.baseURL != expectedURL {
		t.Errorf("expected baseURL %s, got %s", expectedURL, client.baseURL)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	cfg := testConfig()
	customHTTPClient := &http.Client{Timeout: 10 * time.Second}

	client := NewClient(cfg, WithHTTPClient(customHTTPClient))

	if client.httpClient != customHTTPClient {
		t.Error("expected custom HTTP client to be used")
	}
}

func TestExecute_Success(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "application/xml") {
			t.Errorf("expected Content-Type to contain application/xml, got %s", contentType)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE>
			<STATUS>1</STATUS>
			<BODY>
				<DATA>
					<COLLECTION>
						<LEDGER NAME="Cash">
							<NAME>Cash</NAME>
							<PARENT>Current Assets</PARENT>
							<OPENINGBALANCE>10000</OPENINGBALANCE>
						</LEDGER>
					</COLLECTION>
				</DATA>
			</BODY>
		</ENVELOPE>`))
	})
	defer server.Close()

	// Create client pointing to test server
	cfg := testConfig()
	client := NewClient(cfg)
	// Override the base URL to point to test server
	client.baseURL = server.URL

	ctx := context.Background()
	xmlPayload := `<ENVELOPE><HEADER><TALLYREQUEST>Export Data</TALLYREQUEST></HEADER></ENVELOPE>`

	resp, err := client.Execute(ctx, RequestExportData, xmlPayload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != 1 {
		t.Errorf("expected status 1, got %d", resp.Status)
	}

	if len(resp.Ledgers) != 1 {
		t.Errorf("expected 1 ledger, got %d", len(resp.Ledgers))
	}
}

func TestExecute_ConnectionFailure(t *testing.T) {
	cfg := testConfig()
	cfg.Host = "invalid-host-that-does-not-exist"
	cfg.MaxRetries = 0 // Disable retries for faster test

	client := NewClient(cfg)
	client.httpClient.Timeout = 100 * time.Millisecond

	ctx := context.Background()
	xmlPayload := `<ENVELOPE><HEADER><TALLYREQUEST>Export Data</TALLYREQUEST></HEADER></ENVELOPE>`

	_, err := client.Execute(ctx, RequestExportData, xmlPayload)
	if err == nil {
		t.Fatal("expected error for invalid host")
	}
}

func TestExecute_ContextCancellation(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Delay response to trigger context cancellation
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	cfg := testConfig()
	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	xmlPayload := `<ENVELOPE><HEADER><TALLYREQUEST>Export Data</TALLYREQUEST></HEADER></ENVELOPE>`

	_, err := client.Execute(ctx, RequestExportData, xmlPayload)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestExecute_RetryOnFailure(t *testing.T) {
	attemptCount := 0
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 2 {
			// Fail first request
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Succeed on second attempt
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE><STATUS>1</STATUS><BODY><DATA><COLLECTION></COLLECTION></DATA></BODY></ENVELOPE>`))
	})
	defer server.Close()

	cfg := testConfig()
	cfg.MaxRetries = 3
	cfg.RetryDelay = 10 * time.Millisecond

	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx := context.Background()
	xmlPayload := `<ENVELOPE><HEADER><TALLYREQUEST>Export Data</TALLYREQUEST></HEADER></ENVELOPE>`

	resp, err := client.Execute(ctx, RequestExportData, xmlPayload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != 1 {
		t.Errorf("expected status 1, got %d", resp.Status)
	}

	if attemptCount < 2 {
		t.Errorf("expected at least 2 attempts, got %d", attemptCount)
	}
}

func TestGetLedgers(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Verify the request contains expected report name
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		if !strings.Contains(string(body), "List of Ledgers") {
			t.Errorf("expected request to contain 'List of Ledgers', got %s", string(body))
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE>
			<STATUS>1</STATUS>
			<BODY>
				<DATA>
					<COLLECTION>
						<LEDGER NAME="Cash">
							<NAME>Cash</NAME>
							<PARENT>Current Assets</PARENT>
							<MASTERID>1</MASTERID>
						</LEDGER>
						<LEDGER NAME="Bank">
							<NAME>Bank</NAME>
							<PARENT>Bank Accounts</PARENT>
							<MASTERID>2</MASTERID>
						</LEDGER>
					</COLLECTION>
				</DATA>
			</BODY>
		</ENVELOPE>`))
	})
	defer server.Close()

	cfg := testConfig()
	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.GetLedgers(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Ledgers) != 2 {
		t.Errorf("expected 2 ledgers, got %d", len(resp.Ledgers))
	}
}

func TestGetVouchers(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		// Verify date range is included
		if !strings.Contains(string(body), "SVFROMDATE") {
			t.Error("expected request to contain SVFROMDATE")
		}
		if !strings.Contains(string(body), "SVTODATE") {
			t.Error("expected request to contain SVTODATE")
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE>
			<STATUS>1</STATUS>
			<BODY>
				<DATA>
					<COLLECTION>
						<VOUCHER VCHTYPE="Sales">
							<VOUCHERNUMBER>001</VOUCHERNUMBER>
							<DATE>20240101</DATE>
							<MASTERID>100</MASTERID>
						</VOUCHER>
					</COLLECTION>
				</DATA>
			</BODY>
		</ENVELOPE>`))
	})
	defer server.Close()

	cfg := testConfig()
	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx := context.Background()
	fromDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	resp, err := client.GetVouchers(ctx, fromDate, toDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Vouchers) != 1 {
		t.Errorf("expected 1 voucher, got %d", len(resp.Vouchers))
	}

	if resp.Vouchers[0].VoucherNumber != "001" {
		t.Errorf("expected voucher number '001', got '%s'", resp.Vouchers[0].VoucherNumber)
	}
}

func TestGetVouchersModifiedSince(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		if !strings.Contains(string(body), "LASTVCHRID") {
			t.Error("expected request to contain LASTVCHRID")
		}
		if !strings.Contains(string(body), "12345") {
			t.Error("expected request to contain last alter ID '12345'")
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE>
			<STATUS>1</STATUS>
			<BODY>
				<DATA>
					<COLLECTION>
						<VOUCHER VCHTYPE="Payment">
							<VOUCHERNUMBER>002</VOUCHERNUMBER>
							<ALTERID>12346</ALTERID>
						</VOUCHER>
					</COLLECTION>
				</DATA>
			</BODY>
		</ENVELOPE>`))
	})
	defer server.Close()

	cfg := testConfig()
	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.GetVouchersModifiedSince(ctx, "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Vouchers) != 1 {
		t.Errorf("expected 1 voucher, got %d", len(resp.Vouchers))
	}
}

func TestPing_Success(t *testing.T) {
	server := mockTallyServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ENVELOPE><STATUS>1</STATUS><BODY><DATA></DATA></BODY></ENVELOPE>`))
	})
	defer server.Close()

	cfg := testConfig()
	client := NewClient(cfg)
	client.baseURL = server.URL

	ctx := context.Background()
	err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPing_Failure(t *testing.T) {
	cfg := testConfig()
	cfg.Host = "invalid-host"
	cfg.MaxRetries = 0

	client := NewClient(cfg)
	client.httpClient.Timeout = 100 * time.Millisecond

	ctx := context.Background()
	err := client.Ping(ctx)
	if err == nil {
		t.Fatal("expected ping to fail for invalid host")
	}
}

func TestFormatTallyDate(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "standard date",
			input:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "20240115",
		},
		{
			name:     "end of year",
			input:    time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "20241231",
		},
		{
			name:     "start of year",
			input:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "20240101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTallyDate(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseTallyDate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  time.Time
		wantError bool
	}{
		{
			name:     "standard date",
			input:    "20240115",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "invalid format",
			input:     "2024-01-15",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTallyDate(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special chars",
			input:    "Simple Text",
			expected: "Simple Text",
		},
		{
			name:     "ampersand",
			input:    "Smith & Co.",
			expected: "Smith &amp; Co.",
		},
		{
			name:     "less than",
			input:    "x < y",
			expected: "x &lt; y",
		},
		{
			name:     "greater than",
			input:    "x > y",
			expected: "x &gt; y",
		},
		{
			name:     "quotes",
			input:    `Say "Hello"`,
			expected: "Say &#34;Hello&#34;",
		},
		{
			name:     "multiple special chars",
			input:    `<tag attr="value">content & more</tag>`,
			expected: "&lt;tag attr=&#34;value&#34;&gt;content &amp; more&lt;/tag&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeXML(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "connection failed",
			err:      ErrConnectionFailed,
			expected: true,
		},
		{
			name:     "timeout",
			err:      ErrTimeout,
			expected: true,
		},
		{
			name:     "operation failed",
			err:      ErrOperationFailed,
			expected: false,
		},
		{
			name:     "invalid XML",
			err:      ErrInvalidXML,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRequestOptions(t *testing.T) {
	t.Run("WithCompany", func(t *testing.T) {
		opts := &requestOptions{}
		WithCompany("Test Corp")(opts)
		if opts.company != "Test Corp" {
			t.Errorf("expected company 'Test Corp', got '%s'", opts.company)
		}
	})

	t.Run("WithDateRange", func(t *testing.T) {
		opts := &requestOptions{}
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
		WithDateRange(from, to)(opts)

		if !opts.fromDate.Equal(from) {
			t.Errorf("expected from date %v, got %v", from, opts.fromDate)
		}
		if !opts.toDate.Equal(to) {
			t.Errorf("expected to date %v, got %v", to, opts.toDate)
		}
	})

	t.Run("WithVoucherType", func(t *testing.T) {
		opts := &requestOptions{}
		WithVoucherType(VoucherTypeSales)(opts)
		if opts.voucherType != VoucherTypeSales {
			t.Errorf("expected voucher type '%s', got '%s'", VoucherTypeSales, opts.voucherType)
		}
	})

	t.Run("WithLastAlterID", func(t *testing.T) {
		opts := &requestOptions{}
		WithLastAlterID("12345")(opts)
		if opts.lastAlterID != "12345" {
			t.Errorf("expected last alter ID '12345', got '%s'", opts.lastAlterID)
		}
	})
}

func TestLedgerXMLParsing(t *testing.T) {
	xmlData := `<LEDGER NAME="Cash Account">
		<NAME>Cash Account</NAME>
		<PARENT>Current Assets</PARENT>
		<OPENINGBALANCE>50000</OPENINGBALANCE>
		<CLOSINGBALANCE>75000</CLOSINGBALANCE>
		<MASTERID>L001</MASTERID>
		<ALTERID>A001</ALTERID>
		<PARTYGSTIN>29ABCDE1234F1Z5</PARTYGSTIN>
		<CREDITPERIOD>30</CREDITPERIOD>
		<CREDITLIMIT>100000</CREDITLIMIT>
	</LEDGER>`

	var ledger Ledger
	err := xml.Unmarshal([]byte(xmlData), &ledger)
	if err != nil {
		t.Fatalf("failed to parse ledger XML: %v", err)
	}

	if ledger.Name != "Cash Account" {
		t.Errorf("expected name attr 'Cash Account', got '%s'", ledger.Name)
	}
	if ledger.LedgerName != "Cash Account" {
		t.Errorf("expected ledger name 'Cash Account', got '%s'", ledger.LedgerName)
	}
	if ledger.Parent != "Current Assets" {
		t.Errorf("expected parent 'Current Assets', got '%s'", ledger.Parent)
	}
	if ledger.OpeningBalance != 50000 {
		t.Errorf("expected opening balance 50000, got %f", ledger.OpeningBalance)
	}
	if ledger.ClosingBalance != 75000 {
		t.Errorf("expected closing balance 75000, got %f", ledger.ClosingBalance)
	}
	if ledger.GSTRegistration != "29ABCDE1234F1Z5" {
		t.Errorf("expected GST registration '29ABCDE1234F1Z5', got '%s'", ledger.GSTRegistration)
	}
}

func TestVoucherXMLParsing(t *testing.T) {
	xmlData := `<VOUCHER VCHTYPE="Sales" ACTION="Create">
		<VOUCHERNUMBER>S-001</VOUCHERNUMBER>
		<DATE>20240115</DATE>
		<PARTYLEDGERNAME>Customer A</PARTYLEDGERNAME>
		<NARRATION>Sale of goods</NARRATION>
		<MASTERID>V001</MASTERID>
		<ALTERID>A100</ALTERID>
		<ALLLEDGERENTRIES.LIST>
			<LEDGERNAME>Sales Account</LEDGERNAME>
			<ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>
			<AMOUNT>-10000</AMOUNT>
		</ALLLEDGERENTRIES.LIST>
		<ALLLEDGERENTRIES.LIST>
			<LEDGERNAME>Customer A</LEDGERNAME>
			<ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>
			<AMOUNT>10000</AMOUNT>
		</ALLLEDGERENTRIES.LIST>
	</VOUCHER>`

	var voucher Voucher
	err := xml.Unmarshal([]byte(xmlData), &voucher)
	if err != nil {
		t.Fatalf("failed to parse voucher XML: %v", err)
	}

	if voucher.VoucherType != "Sales" {
		t.Errorf("expected voucher type 'Sales', got '%s'", voucher.VoucherType)
	}
	if voucher.VoucherNumber != "S-001" {
		t.Errorf("expected voucher number 'S-001', got '%s'", voucher.VoucherNumber)
	}
	if voucher.Date != "20240115" {
		t.Errorf("expected date '20240115', got '%s'", voucher.Date)
	}
	if voucher.PartyName != "Customer A" {
		t.Errorf("expected party name 'Customer A', got '%s'", voucher.PartyName)
	}
	if len(voucher.LedgerEntries) != 2 {
		t.Errorf("expected 2 ledger entries, got %d", len(voucher.LedgerEntries))
	}
}

func TestStockItemXMLParsing(t *testing.T) {
	xmlData := `<STOCKITEM NAME="Paracetamol 500mg">
		<NAME>Paracetamol 500mg</NAME>
		<PARENT>Medicines</PARENT>
		<CATEGORY>Analgesics</CATEGORY>
		<BASEUNITS>Tab</BASEUNITS>
		<GSTRATE>5</GSTRATE>
		<HSNCODE>30049099</HSNCODE>
		<PARTNUMBER>MED-001</PARTNUMBER>
		<STANDARDCOST>2.50</STANDARDCOST>
		<STANDARDPRICE>5.00</STANDARDPRICE>
		<MRP>5.50</MRP>
		<MASTERID>I001</MASTERID>
	</STOCKITEM>`

	var item StockItem
	err := xml.Unmarshal([]byte(xmlData), &item)
	if err != nil {
		t.Fatalf("failed to parse stock item XML: %v", err)
	}

	if item.Name != "Paracetamol 500mg" {
		t.Errorf("expected name 'Paracetamol 500mg', got '%s'", item.Name)
	}
	if item.Category != "Analgesics" {
		t.Errorf("expected category 'Analgesics', got '%s'", item.Category)
	}
	if item.GSTRate != 5 {
		t.Errorf("expected GST rate 5, got %f", item.GSTRate)
	}
	if item.HSNCode != "30049099" {
		t.Errorf("expected HSN code '30049099', got '%s'", item.HSNCode)
	}
	if item.StandardPrice != 5.00 {
		t.Errorf("expected standard price 5.00, got %f", item.StandardPrice)
	}
}

func TestBuildExportRequest(t *testing.T) {
	cfg := testConfig()
	client := NewClient(cfg)

	opts := &requestOptions{
		fromDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		toDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	xml := client.buildExportRequest(ReportLedgerMasters, "Test Company", opts)

	if !strings.Contains(xml, "ENVELOPE") {
		t.Error("expected XML to contain ENVELOPE")
	}
	if !strings.Contains(xml, "Export Data") {
		t.Error("expected XML to contain 'Export Data'")
	}
	if !strings.Contains(xml, "Test Company") {
		t.Error("expected XML to contain company name")
	}
	if !strings.Contains(xml, "20240101") {
		t.Error("expected XML to contain from date")
	}
	if !strings.Contains(xml, "20240131") {
		t.Error("expected XML to contain to date")
	}
	if !strings.Contains(xml, "List of Ledgers") {
		t.Error("expected XML to contain report name")
	}
}

func TestExtractErrorFromXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "line error",
			input:    `<ENVELOPE><LINEERROR>Invalid voucher type</LINEERROR></ENVELOPE>`,
			expected: "Invalid voucher type",
		},
		{
			name:     "error tag",
			input:    `<ENVELOPE><ERROR>Connection failed</ERROR></ENVELOPE>`,
			expected: "Connection failed",
		},
		{
			name:     "no error found",
			input:    `<ENVELOPE><STATUS>1</STATUS></ENVELOPE>`,
			expected: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractErrorFromXML(tt.input)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
