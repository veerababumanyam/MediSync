// Package hims provides a REST client for the Healthcare Information Management System.
//
// This client handles communication with HIMS REST APIs for extracting
// patient, appointment, billing, and pharmacy data. It implements rate limiting,
// exponential backoff retry, and incremental sync via modified_since parameter.
//
// Usage:
//
//	cfg := config.MustLoad()
//	client := hims.NewClient(cfg.HIMS, logger)
//
//	ctx := context.Background()
//	patients, err := client.GetPatients(ctx, &hims.PatientOptions{
//	    ModifiedSince: time.Now().Add(-24 * time.Hour),
//	})
package hims

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/medisync/medisync/internal/config"
)

var (
	// ErrConnectionFailed indicates the HIMS server could not be reached.
	ErrConnectionFailed = errors.New("hims: connection failed")
	// ErrRequestFailed indicates the HIMS request failed.
	ErrRequestFailed = errors.New("hims: request failed")
	// ErrResponseParseFailed indicates the HIMS response could not be parsed.
	ErrResponseParseFailed = errors.New("hims: response parse failed")
	// ErrTimeout indicates the request timed out.
	ErrTimeout = errors.New("hims: request timeout")
	// ErrRetryExhausted indicates all retry attempts have been exhausted.
	ErrRetryExhausted = errors.New("hims: retry attempts exhausted")
	// ErrUnauthorized indicates authentication failed.
	ErrUnauthorized = errors.New("hims: unauthorized")
	// ErrRateLimited indicates the rate limit was exceeded.
	ErrRateLimited = errors.New("hims: rate limited")
)

// Client provides methods to interact with HIMS REST API.
type Client struct {
	config     config.HIMSConfig
	httpClient *http.Client
	logger     *slog.Logger
	baseURL    string
	rateLimiter <-chan time.Time
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

// NewClient creates a new HIMS client with the given configuration.
func NewClient(cfg config.HIMSConfig, opts ...ClientOption) *Client {
	c := &Client{
		config:  cfg,
		baseURL: strings.TrimSuffix(cfg.URL, "/"),
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

	// Set up rate limiter using ticker
	rateLimit := cfg.RateLimitPerSecond
	if rateLimit <= 0 {
		rateLimit = 10 // Default to 10 requests per second
	}
	c.rateLimiter = time.NewTicker(time.Second / time.Duration(rateLimit)).C

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request with retry logic and rate limiting.
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := c.config.RetryDelay * time.Duration(1<<uint(attempt-1))
			c.logger.Debug("retrying hims request",
				slog.Int("attempt", attempt),
				slog.Duration("backoff", backoff),
			)

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
			case <-time.After(backoff):
			}
		}

		// Rate limiting
		select {
		case <-c.rateLimiter:
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
		}

		// Create request
		var reqBody io.Reader
		if body != nil {
			reqBody = bytes.NewReader(body)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			lastErr = fmt.Errorf("%w: failed to create request: %v", ErrRequestFailed, err)
			continue
		}

		// Set headers
		if c.config.APIKey != "" {
			req.Header.Set("X-API-Key", c.config.APIKey)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Execute request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, fmt.Errorf("%w: %v", ErrTimeout, ctx.Err())
			}
			lastErr = fmt.Errorf("%w: %v", ErrConnectionFailed, err)
			if isRetryableError(err) {
				continue
			}
			return nil, lastErr
		}

		// Check status code
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = ErrRateLimited
			// Retry after rate limit
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, ErrUnauthorized
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("%w: server returned %d", ErrRequestFailed, resp.StatusCode)
			continue // Retry server errors
		}

		if resp.StatusCode >= 400 {
			// Parse error response
			var apiErr APIError
			if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil {
				resp.Body.Close()
				return nil, &apiErr
			}
			resp.Body.Close()
			return nil, fmt.Errorf("%w: HTTP %d", ErrRequestFailed, resp.StatusCode)
		}

		return resp, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrRetryExhausted, lastErr)
}

// PatientOptions provides optional parameters for patient queries.
type PatientOptions struct {
	ModifiedSince *time.Time
	IsActive      *bool
	Limit         int
	Offset        int
}

// GetPatients retrieves patients from HIMS.
func (c *Client) GetPatients(ctx context.Context, opts *PatientOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &PatientOptions{}
	}

	path := "/patients"
	params := buildQueryParams(map[string]interface{}{
		"modified_since": opts.ModifiedSince,
		"is_active":      opts.IsActive,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var patients []Patient
	result.Data = &patients

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = patients
	return &result, nil
}

// DoctorOptions provides optional parameters for doctor queries.
type DoctorOptions struct {
	ModifiedSince *time.Time
	IsActive      *bool
	DepartmentID  *string
	Limit         int
	Offset        int
}

// GetDoctors retrieves doctors from HIMS.
func (c *Client) GetDoctors(ctx context.Context, opts *DoctorOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &DoctorOptions{}
	}

	path := "/doctors"
	params := buildQueryParams(map[string]interface{}{
		"modified_since": opts.ModifiedSince,
		"is_active":      opts.IsActive,
		"department_id":  opts.DepartmentID,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var doctors []Doctor
	result.Data = &doctors

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = doctors
	return &result, nil
}

// DrugOptions provides optional parameters for drug queries.
type DrugOptions struct {
	ModifiedSince *time.Time
	IsActive      *bool
	Category      *string
	Limit         int
	Offset        int
}

// GetDrugs retrieves drugs from HIMS.
func (c *Client) GetDrugs(ctx context.Context, opts *DrugOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &DrugOptions{}
	}

	path := "/drugs"
	params := buildQueryParams(map[string]interface{}{
		"modified_since": opts.ModifiedSince,
		"is_active":      opts.IsActive,
		"category":       opts.Category,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var drugs []Drug
	result.Data = &drugs

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = drugs
	return &result, nil
}

// AppointmentOptions provides optional parameters for appointment queries.
type AppointmentOptions struct {
	StartDate     *time.Time
	EndDate       *time.Time
	ModifiedSince *time.Time
	Status        *string
	PatientID     *string
	DoctorID      *string
	DepartmentID  *string
	Limit         int
	Offset        int
}

// GetAppointments retrieves appointments from HIMS.
func (c *Client) GetAppointments(ctx context.Context, opts *AppointmentOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &AppointmentOptions{}
	}

	path := "/appointments"
	params := buildQueryParams(map[string]interface{}{
		"start_date":     opts.StartDate,
		"end_date":       opts.EndDate,
		"modified_since": opts.ModifiedSince,
		"status":         opts.Status,
		"patient_id":     opts.PatientID,
		"doctor_id":      opts.DoctorID,
		"department_id":  opts.DepartmentID,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var appointments []Appointment
	result.Data = &appointments

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = appointments
	return &result, nil
}

// BillingOptions provides optional parameters for billing queries.
type BillingOptions struct {
	StartDate     *time.Time
	EndDate       *time.Time
	ModifiedSince *time.Time
	PaymentStatus *string
	PatientID     *string
	BillType      *string
	Limit         int
	Offset        int
}

// GetBilling retrieves billing records from HIMS.
func (c *Client) GetBilling(ctx context.Context, opts *BillingOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &BillingOptions{}
	}

	path := "/billing"
	params := buildQueryParams(map[string]interface{}{
		"start_date":     opts.StartDate,
		"end_date":       opts.EndDate,
		"modified_since": opts.ModifiedSince,
		"payment_status": opts.PaymentStatus,
		"patient_id":     opts.PatientID,
		"bill_type":      opts.BillType,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var billing []Billing
	result.Data = &billing

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = billing
	return &result, nil
}

// PharmacyOptions provides optional parameters for pharmacy dispensation queries.
type PharmacyOptions struct {
	StartDate     *time.Time
	EndDate       *time.Time
	ModifiedSince *time.Time
	PatientID     *string
	DrugID        *string
	PrescriptionID *string
	Limit         int
	Offset        int
}

// GetPharmacyDispensations retrieves pharmacy dispensations from HIMS.
func (c *Client) GetPharmacyDispensations(ctx context.Context, opts *PharmacyOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &PharmacyOptions{}
	}

	path := "/pharmacy/dispensations"
	params := buildQueryParams(map[string]interface{}{
		"start_date":     opts.StartDate,
		"end_date":       opts.EndDate,
		"modified_since": opts.ModifiedSince,
		"patient_id":     opts.PatientID,
		"drug_id":        opts.DrugID,
		"prescription_id": opts.PrescriptionID,
		"limit":          opts.Limit,
		"offset":         opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var dispensations []PharmacyDispensation
	result.Data = &dispensations

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = dispensations
	return &result, nil
}

// DepartmentsOptions provides optional parameters for department queries.
type DepartmentsOptions struct {
	IsActive *bool
	Limit    int
	Offset   int
}

// GetDepartments retrieves departments from HIMS.
func (c *Client) GetDepartments(ctx context.Context, opts *DepartmentsOptions) (*PagedResponse, error) {
	if opts == nil {
		opts = &DepartmentsOptions{}
	}

	path := "/departments"
	params := buildQueryParams(map[string]interface{}{
		"is_active": opts.IsActive,
		"limit":     opts.Limit,
		"offset":    opts.Offset,
	})

	if params != "" {
		path += "?" + params
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PagedResponse
	var departments []Department
	result.Data = &departments

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResponseParseFailed, err)
	}

	result.Data = departments
	return &result, nil
}

// Ping checks if the HIMS server is reachable.
func (c *Client) Ping(ctx context.Context) error {
	resp, err := c.doRequest(ctx, http.MethodGet, "/health", nil)
	if err != nil {
		// Even if error, if we got a response the server is reachable
		if errors.Is(err, ErrResponseParseFailed) || errors.Is(err, ErrUnauthorized) {
			return nil
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetAllPatients retrieves all patients using pagination.
func (c *Client) GetAllPatients(ctx context.Context, opts *PatientOptions) ([]Patient, error) {
	if opts == nil {
		opts = &PatientOptions{}
	}

	// Set a reasonable page size if not specified
	if opts.Limit == 0 {
		opts.Limit = 500
	}

	var allPatients []Patient
	offset := 0

	for {
		opts.Offset = offset
		result, err := c.GetPatients(ctx, opts)
		if err != nil {
			return nil, err
		}

		patients, ok := result.Data.([]Patient)
		if !ok {
			return nil, fmt.Errorf("unexpected data type")
		}

		allPatients = append(allPatients, patients...)

		if len(patients) < opts.Limit {
			break
		}

		offset += len(patients)
	}

	return allPatients, nil
}

// GetAllDoctors retrieves all doctors using pagination.
func (c *Client) GetAllDoctors(ctx context.Context, opts *DoctorOptions) ([]Doctor, error) {
	if opts == nil {
		opts = &DoctorOptions{}
	}

	if opts.Limit == 0 {
		opts.Limit = 500
	}

	var allDoctors []Doctor
	offset := 0

	for {
		opts.Offset = offset
		result, err := c.GetDoctors(ctx, opts)
		if err != nil {
			return nil, err
		}

		doctors, ok := result.Data.([]Doctor)
		if !ok {
			return nil, fmt.Errorf("unexpected data type")
		}

		allDoctors = append(allDoctors, doctors...)

		if len(doctors) < opts.Limit {
			break
		}

		offset += len(doctors)
	}

	return allDoctors, nil
}

// GetAllDrugs retrieves all drugs using pagination.
func (c *Client) GetAllDrugs(ctx context.Context, opts *DrugOptions) ([]Drug, error) {
	if opts == nil {
		opts = &DrugOptions{}
	}

	if opts.Limit == 0 {
		opts.Limit = 500
	}

	var allDrugs []Drug
	offset := 0

	for {
		opts.Offset = offset
		result, err := c.GetDrugs(ctx, opts)
		if err != nil {
			return nil, err
		}

		drugs, ok := result.Data.([]Drug)
		if !ok {
			return nil, fmt.Errorf("unexpected data type")
		}

		allDrugs = append(allDrugs, drugs...)

		if len(drugs) < opts.Limit {
			break
		}

		offset += len(drugs)
	}

	return allDrugs, nil
}

// buildQueryParams builds a URL query string from options.
func buildQueryParams(params map[string]interface{}) string {
	var parts []string

	for key, value := range params {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case bool:
			if v {
				parts = append(parts, fmt.Sprintf("%s=true", key))
			} else {
				parts = append(parts, fmt.Sprintf("%s=false", key))
			}
		case int:
			if v > 0 {
				parts = append(parts, fmt.Sprintf("%s=%d", key, v))
			}
		case string:
			if v != "" {
				parts = append(parts, fmt.Sprintf("%s=%s", key, v))
			}
		case *string:
			if v != nil && *v != "" {
				parts = append(parts, fmt.Sprintf("%s=%s", key, *v))
			}
		case *bool:
			if v != nil {
				if *v {
					parts = append(parts, fmt.Sprintf("%s=true", key))
				} else {
					parts = append(parts, fmt.Sprintf("%s=false", key))
				}
			}
		case *int:
			if v != nil && *v > 0 {
				parts = append(parts, fmt.Sprintf("%s=%d", key, *v))
			}
		case time.Time:
			if !v.IsZero() {
				parts = append(parts, fmt.Sprintf("%s=%s", key, v.Format(time.RFC3339)))
			}
		case *time.Time:
			if v != nil && !v.IsZero() {
				parts = append(parts, fmt.Sprintf("%s=%s", key, v.Format(time.RFC3339)))
			}
		}
	}

	return strings.Join(parts, "&")
}

// isRetryableError checks if an error should trigger a retry.
func isRetryableError(err error) bool {
	if errors.Is(err, ErrConnectionFailed) {
		return true
	}
	if errors.Is(err, ErrTimeout) {
		return true
	}
	if errors.Is(err, ErrRateLimited) {
		return true
	}
	// Check for common retryable HTTP errors
	errStr := err.Error()
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary failure") {
		return true
	}
	return false
}
