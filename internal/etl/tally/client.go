// Package tally provides a client for communicating with Tally ERP via TDL XML over HTTP.
//
// This package implements the Tally Gateway interface for extracting and syncing
// data between MediSync and Tally Prime/ERP. It supports both master data (ledgers,
// cost centres, stock items) and transactional data (vouchers, stock movements).
//
// Usage:
//
//	cfg := config.MustLoad()
//	client := tally.NewClient(cfg.Tally, logger)
//
//	ctx := context.Background()
//	response, err := client.Execute(ctx, tally.RequestExportData, "<ENVELOPE>...</ENVELOPE>")
//	if err != nil {
//	    log.Error("failed to execute request", "error", err)
//	}
package tally

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/medisync/medisync/internal/config"
)

// Request types for Tally API.
const (
	// RequestImportData is used to import data into Tally.
	RequestImportData = "Import Data"
	// RequestExportData is used to export data from Tally.
	RequestExportData = "Export Data"
	// RequestExecute is used to execute TDL functions in Tally.
	RequestExecute = "Execute"
)

// Report names for common Tally reports.
const (
	// ReportAllMasters retrieves all master data.
	ReportAllMasters = "All Masters"
	// ReportLedgerMasters retrieves ledger masters.
	ReportLedgerMasters = "List of Ledgers"
	// ReportCostCentres retrieves cost centres.
	ReportCostCentres = "List of Cost Centres"
	// ReportStockItems retrieves stock items.
	ReportStockItems = "List of Stock Items"
	// ReportDayBook retrieves vouchers for a date range.
	ReportDayBook = "Day Book"
	// ReportVouchers retrieves all vouchers.
	ReportVouchers = "List of Vouchers"
	// ReportStockSummary retrieves stock summary.
	ReportStockSummary = "Stock Summary"
)

// Voucher types in Tally.
const (
	VoucherTypeSales       = "Sales"
	VoucherTypePurchase    = "Purchase"
	VoucherTypeReceipt     = "Receipt"
	VoucherTypePayment     = "Payment"
	VoucherTypeJournal     = "Journal"
	VoucherTypeContra      = "Contra"
	VoucherTypeCreditNote  = "Credit Note"
	VoucherTypeDebitNote   = "Debit Note"
	VoucherTypeDeliveryNote = "Delivery Note"
	VoucherTypeReceiptNote = "Receipt Note"
)

// Error types for Tally operations.
var (
	// ErrConnectionFailed indicates the Tally server could not be reached.
	ErrConnectionFailed = errors.New("tally: connection failed")
	// ErrRequestFailed indicates the Tally request failed.
	ErrRequestFailed = errors.New("tally: request failed")
	// ErrResponseParseFailed indicates the Tally response could not be parsed.
	ErrResponseParseFailed = errors.New("tally: response parse failed")
	// ErrOperationFailed indicates the Tally operation returned an error status.
	ErrOperationFailed = errors.New("tally: operation failed")
	// ErrInvalidXML indicates the provided XML is invalid.
	ErrInvalidXML = errors.New("tally: invalid XML")
	// ErrTimeout indicates the request timed out.
	ErrTimeout = errors.New("tally: request timeout")
	// ErrRetryExhausted indicates all retry attempts have been exhausted.
	ErrRetryExhausted = errors.New("tally: retry attempts exhausted")
	// ErrCompanyNotFound indicates the specified company was not found.
	ErrCompanyNotFound = errors.New("tally: company not found")
)

// Envelope represents the TDL XML envelope structure for Tally communication.
type Envelope struct {
	XMLName xml.Name `xml:"ENVELOPE"`
	Header  Header   `xml:"HEADER"`
	Body    Body     `xml:"BODY"`
}

// Header contains the request metadata.
type Header struct {
	TallyRequest string `xml:"TALLYREQUEST"`
	Type         string `xml:"TYPE,omitempty"`
	ID           string `xml:"ID,omitempty"`
	Version      int    `xml:"VERSION,omitempty"`
}

// Body contains the request or response data.
type Body struct {
	Desc       *Desc       `xml:"DESC,omitempty"`
	Data       *Data       `xml:"DATA,omitempty"`
	ImportData *ImportData `xml:"IMPORTDATA,omitempty"`
	ExportData *ExportData `xml:"EXPORTDATA,omitempty"`
}

// Desc contains static variables for the request.
type Desc struct {
	StaticVariables *StaticVariables `xml:"STATICVARIABLES,omitempty"`
	TDL             *TDL             `xml:"TDL,omitempty"`
	FetchList       []string         `xml:"FETCHLIST>FETCH,omitempty"`
}

// StaticVariables contains variables passed to Tally reports.
type StaticVariables struct {
	Company           string `xml:"SVCURRENTCOMPANY,omitempty"`
	FromDate          string `xml:"SVFROMDATE,omitempty"`
	ToDate            string `xml:"SVTODATE,omitempty"`
	ExportFormat      string `xml:"SVEXPORTFORMAT,omitempty"`
	DayBookDate       string `xml:"SVCDSSVDAYBOOKDATE,omitempty"`
	LedgerName        string `xml:"LEDGERNAME,omitempty"`
	VoucherTypeName   string `xml:"VOUCHERTYPENAME,omitempty"`
	StockItemName     string `xml:"STOCKITEMNAME,omitempty"`
	AlterID           string `xml:"ALTERID,omitempty"`
	LastVoucherAlterID string `xml:"LASTVCHRID,omitempty"`
}

// TDL contains custom TDL definitions.
type TDL struct {
	TDLMessage string `xml:"TDLMESSAGE,omitempty"`
}

// Data contains the response data.
type Data struct {
	Collection   *Collection   `xml:"COLLECTION,omitempty"`
	TallyMessage *TallyMessage `xml:"TALLYMESSAGE,omitempty"`
}

// Collection represents a collection of Tally objects.
type Collection struct {
	Ledgers    []Ledger    `xml:"LEDGER,omitempty"`
	CostCentres []CostCentre `xml:"COSTCENTRE,omitempty"`
	StockItems []StockItem `xml:"STOCKITEM,omitempty"`
	Vouchers   []Voucher   `xml:"VOUCHER,omitempty"`
}

// ImportData contains the import request structure.
type ImportData struct {
	RequestDesc *RequestDesc `xml:"REQUESTDESC,omitempty"`
	RequestData *RequestData `xml:"REQUESTDATA,omitempty"`
}

// ExportData contains the export request structure.
type ExportData struct {
	RequestDesc *RequestDesc `xml:"REQUESTDESC,omitempty"`
}

// RequestDesc describes what data to export.
type RequestDesc struct {
	ReportName      string           `xml:"REPORTNAME,omitempty"`
	StaticVariables *StaticVariables `xml:"STATICVARIABLES,omitempty"`
}

// RequestData contains the data to import.
type RequestData struct {
	TallyMessage *TallyMessage `xml:"TALLYMESSAGE,omitempty"`
}

// TallyMessage wraps the actual data payload.
type TallyMessage struct {
	XMLNs    string    `xml:"xmlns:UDF,attr,omitempty"`
	Ledgers  []Ledger  `xml:"LEDGER,omitempty"`
	Vouchers []Voucher `xml:"VOUCHER,omitempty"`
}

// Ledger represents a Tally ledger (account).
type Ledger struct {
	XMLName        xml.Name `xml:"LEDGER"`
	Name           string   `xml:"NAME,attr,omitempty"`
	LedgerName     string   `xml:"NAME,omitempty"`
	Parent         string   `xml:"PARENT,omitempty"`
	IsBillWiseOn   string   `xml:"ISBILLWISEON,omitempty"`
	OpeningBalance float64  `xml:"OPENINGBALANCE,omitempty"`
	ClosingBalance float64  `xml:"CLOSINGBALANCE,omitempty"`
	MasterID       string   `xml:"MASTERID,omitempty"`
	AlterID        string   `xml:"ALTERID,omitempty"`
	GSTRegistration string  `xml:"PARTYGSTIN,omitempty"`
	PANNumber      string   `xml:"INCOMETAXNUMBER,omitempty"`
	CreditPeriod   int      `xml:"CREDITPERIOD,omitempty"`
	CreditLimit    float64  `xml:"CREDITLIMIT,omitempty"`
	BankDetails    *BankDetails `xml:"BANKDETAILS,omitempty"`
	Address        []string `xml:"ADDRESS.LIST>ADDRESS,omitempty"`
}

// BankDetails contains bank account information.
type BankDetails struct {
	BankName      string `xml:"BANKNAME,omitempty"`
	AccountNumber string `xml:"ACCOUNTNUMBER,omitempty"`
	IFSCCode      string `xml:"IFSCODE,omitempty"`
}

// CostCentre represents a Tally cost centre.
type CostCentre struct {
	XMLName        xml.Name `xml:"COSTCENTRE"`
	Name           string   `xml:"NAME,attr,omitempty"`
	CostCentreName string   `xml:"NAME,omitempty"`
	Parent         string   `xml:"PARENT,omitempty"`
	Category       string   `xml:"CATEGORY,omitempty"`
	IsRevenue      string   `xml:"ISREVENUE,omitempty"`
	MasterID       string   `xml:"MASTERID,omitempty"`
	AlterID        string   `xml:"ALTERID,omitempty"`
}

// StockItem represents a Tally stock item.
type StockItem struct {
	XMLName          xml.Name `xml:"STOCKITEM"`
	Name             string   `xml:"NAME,attr,omitempty"`
	ItemName         string   `xml:"NAME,omitempty"`
	Parent           string   `xml:"PARENT,omitempty"`
	Category         string   `xml:"CATEGORY,omitempty"`
	BaseUnit         string   `xml:"BASEUNITS,omitempty"`
	OpeningBalance   float64  `xml:"OPENINGBALANCE,omitempty"`
	OpeningRate      float64  `xml:"OPENINGRATE,omitempty"`
	OpeningValue     float64  `xml:"OPENINGVALUE,omitempty"`
	GSTRate          float64  `xml:"GSTRATE,omitempty"`
	HSNCode          string   `xml:"HSNCODE,omitempty"`
	PartNumber       string   `xml:"PARTNUMBER,omitempty"`
	BatchWiseOn      string   `xml:"ISBATCHWISEON,omitempty"`
	ExpiryDateOn     string   `xml:"HASMFGDATE,omitempty"`
	StandardCost     float64  `xml:"STANDARDCOST,omitempty"`
	StandardPrice    float64  `xml:"STANDARDPRICE,omitempty"`
	MRP              float64  `xml:"MRP,omitempty"`
	ReorderLevel     int      `xml:"REORDERLEVEL,omitempty"`
	MinOrderQty      int      `xml:"MINORDERQTY,omitempty"`
	MasterID         string   `xml:"MASTERID,omitempty"`
	AlterID          string   `xml:"ALTERID,omitempty"`
}

// Voucher represents a Tally voucher (transaction).
type Voucher struct {
	XMLName          xml.Name        `xml:"VOUCHER"`
	VoucherType      string          `xml:"VCHTYPE,attr,omitempty"`
	Action           string          `xml:"ACTION,attr,omitempty"`
	ObjView          string          `xml:"OBJVIEW,attr,omitempty"`
	VoucherNumber    string          `xml:"VOUCHERNUMBER,omitempty"`
	Date             string          `xml:"DATE,omitempty"`
	EffectiveDate    string          `xml:"EFFECTIVEDATE,omitempty"`
	PartyName        string          `xml:"PARTYLEDGERNAME,omitempty"`
	Narration        string          `xml:"NARRATION,omitempty"`
	Reference        string          `xml:"REFERENCE,omitempty"`
	ReferenceDate    string          `xml:"REFERENCEDATE,omitempty"`
	GSTRegistration  string          `xml:"PARTYGSTIN,omitempty"`
	IsCancelled      string          `xml:"ISCANCELLED,omitempty"`
	IsOptional       string          `xml:"ISOPTIONAL,omitempty"`
	MasterID         string          `xml:"MASTERID,omitempty"`
	AlterID          string          `xml:"ALTERID,omitempty"`
	VoucherKey       string          `xml:"VCHKEY,omitempty"`
	LedgerEntries    []LedgerEntry   `xml:"ALLLEDGERENTRIES.LIST,omitempty"`
	InventoryEntries []InventoryEntry `xml:"ALLINVENTORYENTRIES.LIST,omitempty"`
	BillAllocations  []BillAllocation `xml:"BILLALLOCATIONS.LIST,omitempty"`
}

// LedgerEntry represents a ledger line in a voucher.
type LedgerEntry struct {
	LedgerName      string  `xml:"LEDGERNAME,omitempty"`
	IsDeemedPositive string  `xml:"ISDEEMEDPOSITIVE,omitempty"`
	IsPartyledger   string  `xml:"ISPARTYLEDGER,omitempty"`
	Amount          float64 `xml:"AMOUNT,omitempty"`
	CostCentre      string  `xml:"COSTCENTRE,omitempty"`
}

// InventoryEntry represents an inventory line in a voucher.
type InventoryEntry struct {
	StockItemName   string  `xml:"STOCKITEMNAME,omitempty"`
	IsDeemedPositive string  `xml:"ISDEEMEDPOSITIVE,omitempty"`
	ActualQty       float64 `xml:"ACTUALQTY,omitempty"`
	BilledQty       float64 `xml:"BILLEDQTY,omitempty"`
	Rate            float64 `xml:"RATE,omitempty"`
	Amount          float64 `xml:"AMOUNT,omitempty"`
	GodownName      string  `xml:"GODOWNNAME,omitempty"`
	BatchName       string  `xml:"BATCHNAME,omitempty"`
}

// BillAllocation represents a bill allocation in a voucher.
type BillAllocation struct {
	Name      string  `xml:"NAME,omitempty"`
	BillType  string  `xml:"BILLTYPE,omitempty"`
	Amount    float64 `xml:"AMOUNT,omitempty"`
}

// Response represents a parsed Tally response.
type Response struct {
	// Status is 1 for success, 0 for failure.
	Status int `xml:"STATUS,omitempty"`
	// Error contains error details if status is 0.
	Error string
	// RawXML contains the raw response XML.
	RawXML string
	// Envelope contains the parsed response envelope.
	Envelope *Envelope
	// Duration is the time taken for the request.
	Duration time.Duration
	// Ledgers contains extracted ledger data.
	Ledgers []Ledger
	// CostCentres contains extracted cost centre data.
	CostCentres []CostCentre
	// StockItems contains extracted stock item data.
	StockItems []StockItem
	// Vouchers contains extracted voucher data.
	Vouchers []Voucher
}

// ResponseEnvelope is used to parse Tally responses.
type ResponseEnvelope struct {
	XMLName xml.Name `xml:"ENVELOPE"`
	Body    struct {
		Data struct {
			Collection struct {
				Ledgers     []Ledger    `xml:"LEDGER"`
				CostCentres []CostCentre `xml:"COSTCENTRE"`
				StockItems  []StockItem `xml:"STOCKITEM"`
				Vouchers    []Voucher   `xml:"VOUCHER"`
			} `xml:"COLLECTION"`
			TallyMessage struct {
				Ledgers  []Ledger  `xml:"LEDGER"`
				Vouchers []Voucher `xml:"VOUCHER"`
			} `xml:"TALLYMESSAGE"`
		} `xml:"DATA"`
		ImportResult struct {
			Created    int    `xml:"CREATED"`
			Altered    int    `xml:"ALTERED"`
			Deleted    int    `xml:"DELETED"`
			LastError  string `xml:"LASTERROR"`
			LastReason string `xml:"LASTREASON"`
		} `xml:"IMPORTRESULT"`
	} `xml:"BODY"`
	Status   int    `xml:"STATUS"`
	LastStatus int  `xml:"LASTSTATUS"`
	Error    string `xml:"LINEERROR"`
}

// Client provides methods to interact with Tally ERP via HTTP.
type Client struct {
	config     config.TallyConfig
	httpClient *http.Client
	logger     *slog.Logger
	baseURL    string
}

// ClientOption is a function that configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Tally client with the given configuration.
func NewClient(cfg config.TallyConfig, opts ...ClientOption) *Client {
	c := &Client{
		config:  cfg,
		baseURL: fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		logger:  slog.Default(),
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Execute sends a raw XML request to Tally and returns the response.
// It handles retries with exponential backoff on transient failures.
func (c *Client) Execute(ctx context.Context, requestType, xmlPayload string) (*Response, error) {
	start := time.Now()

	c.logger.Debug("executing tally request",
		slog.String("request_type", requestType),
		slog.String("base_url", c.baseURL),
	)

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := c.config.RetryDelay * time.Duration(1<<uint(attempt-1))
			c.logger.Debug("retrying tally request",
				slog.Int("attempt", attempt),
				slog.Duration("backoff", backoff),
			)

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
			case <-time.After(backoff):
			}
		}

		resp, err := c.doRequest(ctx, xmlPayload)
		if err != nil {
			lastErr = err
			// Check if error is retryable
			if isRetryableError(err) {
				continue
			}
			return nil, err
		}

		resp.Duration = time.Since(start)
		return resp, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrRetryExhausted, lastErr)
}

// doRequest performs the actual HTTP request to Tally.
func (c *Client) doRequest(ctx context.Context, xmlPayload string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBufferString(xmlPayload))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %v", ErrRequestFailed, err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Accept", "application/xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
		}
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}
	defer resp.Body.Close()

	// Check for HTTP status errors - 5xx errors should be retryable
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: HTTP %d %s", ErrConnectionFailed, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response body: %v", ErrResponseParseFailed, err)
	}

	response := &Response{
		RawXML: string(body),
	}

	// Parse the response
	var envelope ResponseEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		// Try to extract error from raw response
		if strings.Contains(string(body), "ERROR") {
			response.Status = 0
			response.Error = extractErrorFromXML(string(body))
			return response, fmt.Errorf("%w: %s", ErrOperationFailed, response.Error)
		}
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	// Check for explicit errors - only fail if STATUS is explicitly 0 AND there's an error message
	// If STATUS field is missing (defaults to 0), we should still check for data
	hasError := envelope.Error != "" || envelope.Body.ImportResult.LastError != ""
	hasData := len(envelope.Body.Data.Collection.Ledgers) > 0 ||
		len(envelope.Body.Data.Collection.CostCentres) > 0 ||
		len(envelope.Body.Data.Collection.StockItems) > 0 ||
		len(envelope.Body.Data.Collection.Vouchers) > 0 ||
		len(envelope.Body.Data.TallyMessage.Ledgers) > 0 ||
		len(envelope.Body.Data.TallyMessage.Vouchers) > 0

	if hasError && !hasData {
		response.Status = 0
		response.Error = envelope.Error
		if response.Error == "" && envelope.Body.ImportResult.LastError != "" {
			response.Error = fmt.Sprintf("%s: %s",
				envelope.Body.ImportResult.LastError,
				envelope.Body.ImportResult.LastReason,
			)
		}
		return response, fmt.Errorf("%w: %s", ErrOperationFailed, response.Error)
	}

	// Extract data from response
	response.Status = 1
	response.Ledgers = envelope.Body.Data.Collection.Ledgers
	if len(response.Ledgers) == 0 {
		response.Ledgers = envelope.Body.Data.TallyMessage.Ledgers
	}
	response.CostCentres = envelope.Body.Data.Collection.CostCentres
	response.StockItems = envelope.Body.Data.Collection.StockItems
	response.Vouchers = envelope.Body.Data.Collection.Vouchers
	if len(response.Vouchers) == 0 {
		response.Vouchers = envelope.Body.Data.TallyMessage.Vouchers
	}

	c.logger.Debug("tally request completed",
		slog.Int("status", response.Status),
		slog.Int("ledgers_count", len(response.Ledgers)),
		slog.Int("cost_centres_count", len(response.CostCentres)),
		slog.Int("stock_items_count", len(response.StockItems)),
		slog.Int("vouchers_count", len(response.Vouchers)),
	)

	return response, nil
}

// GetLedgers retrieves all ledgers from Tally.
func (c *Client) GetLedgers(ctx context.Context, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildExportRequest(ReportLedgerMasters, company, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// GetCostCentres retrieves all cost centres from Tally.
func (c *Client) GetCostCentres(ctx context.Context, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildExportRequest(ReportCostCentres, company, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// GetStockItems retrieves all stock items from Tally.
func (c *Client) GetStockItems(ctx context.Context, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildExportRequest(ReportStockItems, company, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// GetVouchers retrieves vouchers from Tally within the specified date range.
func (c *Client) GetVouchers(ctx context.Context, fromDate, toDate time.Time, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{
		fromDate: fromDate,
		toDate:   toDate,
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildVoucherExportRequest(company, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// GetVouchersModifiedSince retrieves vouchers modified since the given AlterID.
// This is used for incremental sync.
func (c *Client) GetVouchersModifiedSince(ctx context.Context, lastAlterID string, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{
		lastAlterID: lastAlterID,
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildIncrementalVoucherRequest(company, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// GetDayBook retrieves the day book for a specific date.
func (c *Client) GetDayBook(ctx context.Context, date time.Time, opts ...RequestOption) (*Response, error) {
	reqOpts := &requestOptions{
		fromDate: date,
		toDate:   date,
	}
	for _, opt := range opts {
		opt(reqOpts)
	}

	company := c.config.Company
	if reqOpts.company != "" {
		company = reqOpts.company
	}

	xml := c.buildDayBookRequest(company, date, reqOpts)
	return c.Execute(ctx, RequestExportData, xml)
}

// Ping checks if the Tally server is reachable and responsive.
func (c *Client) Ping(ctx context.Context) error {
	// Simple request to check connectivity
	xmlPayload := `<ENVELOPE>
		<HEADER>
			<TALLYREQUEST>Export Data</TALLYREQUEST>
		</HEADER>
		<BODY>
			<EXPORTDATA>
				<REQUESTDESC>
					<STATICVARIABLES>
						<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>
					</STATICVARIABLES>
					<REPORTNAME>List of Companies</REPORTNAME>
				</REQUESTDESC>
			</EXPORTDATA>
		</BODY>
	</ENVELOPE>`

	_, err := c.Execute(ctx, RequestExportData, xmlPayload)
	if err != nil {
		// Even if the response has an error, if we got a response, the server is reachable
		if errors.Is(err, ErrOperationFailed) {
			return nil
		}
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// GetCompanies retrieves the list of companies loaded in Tally.
func (c *Client) GetCompanies(ctx context.Context) ([]string, error) {
	xmlPayload := `<ENVELOPE>
		<HEADER>
			<TALLYREQUEST>Export Data</TALLYREQUEST>
		</HEADER>
		<BODY>
			<EXPORTDATA>
				<REQUESTDESC>
					<STATICVARIABLES>
						<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>
					</STATICVARIABLES>
					<REPORTNAME>List of Companies</REPORTNAME>
				</REQUESTDESC>
			</EXPORTDATA>
		</BODY>
	</ENVELOPE>`

	resp, err := c.Execute(ctx, RequestExportData, xmlPayload)
	if err != nil {
		return nil, err
	}

	// Parse companies from response
	var companies []string
	// Note: This is a simplified extraction. The actual response structure may vary.
	// In a production implementation, we'd parse the COMPANY elements.
	if resp.RawXML != "" {
		companies = extractCompaniesFromXML(resp.RawXML)
	}

	return companies, nil
}

// requestOptions holds optional parameters for requests.
type requestOptions struct {
	company      string
	fromDate     time.Time
	toDate       time.Time
	voucherType  string
	ledgerName   string
	lastAlterID  string
	fetchList    []string
}

// RequestOption configures a request.
type RequestOption func(*requestOptions)

// WithCompany sets the company name for the request.
func WithCompany(company string) RequestOption {
	return func(o *requestOptions) {
		o.company = company
	}
}

// WithDateRange sets the date range for voucher requests.
func WithDateRange(from, to time.Time) RequestOption {
	return func(o *requestOptions) {
		o.fromDate = from
		o.toDate = to
	}
}

// WithVoucherType filters by voucher type.
func WithVoucherType(voucherType string) RequestOption {
	return func(o *requestOptions) {
		o.voucherType = voucherType
	}
}

// WithLedgerName filters by ledger name.
func WithLedgerName(ledgerName string) RequestOption {
	return func(o *requestOptions) {
		o.ledgerName = ledgerName
	}
}

// WithLastAlterID sets the last alter ID for incremental sync.
func WithLastAlterID(alterID string) RequestOption {
	return func(o *requestOptions) {
		o.lastAlterID = alterID
	}
}

// WithFetchList specifies which fields to fetch.
func WithFetchList(fields ...string) RequestOption {
	return func(o *requestOptions) {
		o.fetchList = fields
	}
}

// buildExportRequest builds an export request XML for master data.
func (c *Client) buildExportRequest(reportName, company string, opts *requestOptions) string {
	var sb strings.Builder
	sb.WriteString(`<ENVELOPE>
	<HEADER>
		<TALLYREQUEST>Export Data</TALLYREQUEST>
	</HEADER>
	<BODY>
		<EXPORTDATA>
			<REQUESTDESC>
				<STATICVARIABLES>
					<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>`)

	if company != "" {
		sb.WriteString(fmt.Sprintf(`
					<SVCURRENTCOMPANY>%s</SVCURRENTCOMPANY>`, escapeXML(company)))
	}

	if !opts.fromDate.IsZero() {
		sb.WriteString(fmt.Sprintf(`
					<SVFROMDATE>%s</SVFROMDATE>`, formatTallyDate(opts.fromDate)))
	}

	if !opts.toDate.IsZero() {
		sb.WriteString(fmt.Sprintf(`
					<SVTODATE>%s</SVTODATE>`, formatTallyDate(opts.toDate)))
	}

	sb.WriteString(`
				</STATICVARIABLES>
				<REPORTNAME>`)
	sb.WriteString(escapeXML(reportName))
	sb.WriteString(`</REPORTNAME>
			</REQUESTDESC>
		</EXPORTDATA>
	</BODY>
</ENVELOPE>`)

	return sb.String()
}

// buildVoucherExportRequest builds an export request for vouchers.
func (c *Client) buildVoucherExportRequest(company string, opts *requestOptions) string {
	var sb strings.Builder
	sb.WriteString(`<ENVELOPE>
	<HEADER>
		<TALLYREQUEST>Export Data</TALLYREQUEST>
	</HEADER>
	<BODY>
		<EXPORTDATA>
			<REQUESTDESC>
				<STATICVARIABLES>
					<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>`)

	if company != "" {
		sb.WriteString(fmt.Sprintf(`
					<SVCURRENTCOMPANY>%s</SVCURRENTCOMPANY>`, escapeXML(company)))
	}

	if !opts.fromDate.IsZero() {
		sb.WriteString(fmt.Sprintf(`
					<SVFROMDATE>%s</SVFROMDATE>`, formatTallyDate(opts.fromDate)))
	}

	if !opts.toDate.IsZero() {
		sb.WriteString(fmt.Sprintf(`
					<SVTODATE>%s</SVTODATE>`, formatTallyDate(opts.toDate)))
	}

	if opts.voucherType != "" {
		sb.WriteString(fmt.Sprintf(`
					<VOUCHERTYPENAME>%s</VOUCHERTYPENAME>`, escapeXML(opts.voucherType)))
	}

	sb.WriteString(`
				</STATICVARIABLES>
				<REPORTNAME>`)
	sb.WriteString(ReportVouchers)
	sb.WriteString(`</REPORTNAME>
			</REQUESTDESC>
		</EXPORTDATA>
	</BODY>
</ENVELOPE>`)

	return sb.String()
}

// buildIncrementalVoucherRequest builds a request for vouchers modified after a given AlterID.
func (c *Client) buildIncrementalVoucherRequest(company string, opts *requestOptions) string {
	var sb strings.Builder
	sb.WriteString(`<ENVELOPE>
	<HEADER>
		<TALLYREQUEST>Export Data</TALLYREQUEST>
	</HEADER>
	<BODY>
		<EXPORTDATA>
			<REQUESTDESC>
				<STATICVARIABLES>
					<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>`)

	if company != "" {
		sb.WriteString(fmt.Sprintf(`
					<SVCURRENTCOMPANY>%s</SVCURRENTCOMPANY>`, escapeXML(company)))
	}

	if opts.lastAlterID != "" {
		sb.WriteString(fmt.Sprintf(`
					<LASTVCHRID>%s</LASTVCHRID>`, escapeXML(opts.lastAlterID)))
	}

	sb.WriteString(`
				</STATICVARIABLES>
				<REPORTNAME>Altered Vouchers</REPORTNAME>
			</REQUESTDESC>
		</EXPORTDATA>
	</BODY>
</ENVELOPE>`)

	return sb.String()
}

// buildDayBookRequest builds a request for day book data.
func (c *Client) buildDayBookRequest(company string, date time.Time, opts *requestOptions) string {
	var sb strings.Builder
	sb.WriteString(`<ENVELOPE>
	<HEADER>
		<TALLYREQUEST>Export Data</TALLYREQUEST>
	</HEADER>
	<BODY>
		<EXPORTDATA>
			<REQUESTDESC>
				<STATICVARIABLES>
					<SVEXPORTFORMAT>$$SysName:XML</SVEXPORTFORMAT>`)

	if company != "" {
		sb.WriteString(fmt.Sprintf(`
					<SVCURRENTCOMPANY>%s</SVCURRENTCOMPANY>`, escapeXML(company)))
	}

	sb.WriteString(fmt.Sprintf(`
					<SVCDSSVDAYBOOKDATE>%s</SVCDSSVDAYBOOKDATE>`, formatTallyDate(date)))

	sb.WriteString(`
				</STATICVARIABLES>
				<REPORTNAME>`)
	sb.WriteString(ReportDayBook)
	sb.WriteString(`</REPORTNAME>
			</REQUESTDESC>
		</EXPORTDATA>
	</BODY>
</ENVELOPE>`)

	return sb.String()
}

// Helper functions

// formatTallyDate formats a time.Time as Tally date string (YYYYMMDD).
func formatTallyDate(t time.Time) string {
	return t.Format("20060102")
}

// parseTallyDate parses a Tally date string (YYYYMMDD) to time.Time.
func parseTallyDate(s string) (time.Time, error) {
	return time.Parse("20060102", s)
}

// escapeXML escapes special XML characters.
func escapeXML(s string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

// isRetryableError checks if an error should trigger a retry.
func isRetryableError(err error) bool {
	if errors.Is(err, ErrConnectionFailed) {
		return true
	}
	if errors.Is(err, ErrTimeout) {
		return true
	}
	// Check for common retryable HTTP errors
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "temporary failure") {
		return true
	}
	return false
}

// extractErrorFromXML attempts to extract error message from raw XML response.
func extractErrorFromXML(xml string) string {
	// Look for common error patterns
	patterns := []string{
		"<LINEERROR>", "<ERROR>", "<ERRORMSG>",
		"<LASTMSG>", "<LASTREASON>",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(xml, pattern); idx != -1 {
			endPattern := "</" + pattern[1:]
			if endIdx := strings.Index(xml[idx:], endPattern); endIdx != -1 {
				return strings.TrimSpace(xml[idx+len(pattern) : idx+endIdx])
			}
		}
	}

	return "unknown error"
}

// extractCompaniesFromXML extracts company names from the response XML.
func extractCompaniesFromXML(xmlData string) []string {
	var companies []string
	// Simple extraction - look for COMPANY NAME tags
	start := 0
	for {
		startTag := strings.Index(xmlData[start:], "<COMPANY NAME=\"")
		if startTag == -1 {
			// Try alternative format
			startTag = strings.Index(xmlData[start:], "<NAME>")
			if startTag == -1 {
				break
			}
			endTag := strings.Index(xmlData[start+startTag:], "</NAME>")
			if endTag == -1 {
				break
			}
			name := xmlData[start+startTag+6 : start+startTag+endTag]
			if name != "" && !containsString(companies, name) {
				companies = append(companies, name)
			}
			start = start + startTag + endTag
		} else {
			endQuote := strings.Index(xmlData[start+startTag+15:], "\"")
			if endQuote == -1 {
				break
			}
			name := xmlData[start+startTag+15 : start+startTag+15+endQuote]
			if name != "" && !containsString(companies, name) {
				companies = append(companies, name)
			}
			start = start + startTag + 15 + endQuote
		}
	}
	return companies
}

// containsString checks if a slice contains a string.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ParseTallyDate is exported for use by other packages.
func ParseTallyDate(s string) (time.Time, error) {
	return parseTallyDate(s)
}

// FormatTallyDate is exported for use by other packages.
func FormatTallyDate(t time.Time) string {
	return formatTallyDate(t)
}
