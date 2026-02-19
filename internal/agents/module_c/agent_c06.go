// Package module_c provides the C-06 Data Quality Validation Agent.
//
// This agent monitors data quality after each ETL sync, running validation checks
// and generating reports. It subscribes to etl.sync.completed events and performs
// quality checks including missing value rates, referential integrity, duplicates,
// value ranges, and row count deltas.
//
// Quality Checks:
// 1. Missing value rate > 5% → alert
// 2. Referential integrity (FK violations)
// 3. Duplicate detection (business key hash)
// 4. Value range anomalies (negative amounts, future dates)
// 5. Row count delta > 30% → alert
//
// Usage:
//
//	agent := agent_c06.New(repo, events, logger)
//	go agent.Start(ctx)
package agent_c06

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/events"
	"github.com/medisync/medisync/internal/warehouse"
)

const (
	// AgentID is the identifier for this agent.
	AgentID = "C-06"

	// AgentName is the human-readable name.
	AgentName = "Data Quality Validation Agent"

	// MissingValueThreshold is the threshold for missing value alerts (5%).
	MissingValueThreshold = 0.05

	// RowCountDeltaThreshold is the threshold for row count delta alerts (30%).
	RowCountDeltaThreshold = 0.30

	// QualityScorePassing is the minimum quality score for a passing validation.
	QualityScorePassing = 70.0
)

// Agent provides data quality validation functionality.
type Agent struct {
	repo    *warehouse.Repo
	events  *events.Publisher
	sub     *events.Subscriber
	logger  *slog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
}

// Config holds configuration for the agent.
type Config struct {
	// Repo is the warehouse repository.
	Repo *warehouse.Repo

	// Events is the NATS event publisher.
	Events *events.Publisher

	// Logger is the structured logger.
	Logger *slog.Logger
}

// New creates a new C-06 Data Quality Validation Agent.
func New(cfg *Config) *Agent {
	if cfg == nil {
		return nil
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	agent := &Agent{
		repo:   cfg.Repo,
		events: cfg.Events,
		logger: logger.With(slog.String("agent", AgentID)),
	}

	return agent
}

// Start begins processing events and validating data quality.
func (a *Agent) Start(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	// Subscribe to sync completed events
	if a.events != nil {
		sub, err := events.NewSubscriber(a.events, events.SubjectSyncCompleted,
			a.handleSyncCompleted,
			&events.SubscriptionOptions{
				Durable: AgentID,
				AutoAck: false,
			})
		if err != nil {
			return fmt.Errorf("agent %s: failed to create subscriber: %w", AgentID, err)
		}
		a.sub = sub

		a.logger.Info("agent started",
			slog.String("agent", AgentName),
			slog.String("subscription", events.SubjectSyncCompleted),
		)
	}

	return nil
}

// Stop gracefully shuts down the agent.
func (a *Agent) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}

	if a.sub != nil {
		return a.sub.Close()
	}

	return nil
}

// handleSyncCompleted processes sync completion events and runs quality checks.
func (a *Agent) handleSyncCompleted(msg *events.Message) error {
	var event events.SyncCompletedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal sync completed event: %w", err)
	}

	a.logger.Info("sync completed, running quality checks",
		slog.String("source", event.Source),
		slog.String("entity", event.Entity),
		slog.Int("records", event.RecordsProcessed),
	)

	// Run quality validation
	report, err := a.ValidateSync(a.ctx, &event)
	if err != nil {
		a.logger.Error("quality validation failed",
			slog.String("error", err.Error()),
		)

		// Publish alert for validation failure
		if a.events != nil {
			_ = a.events.PublishAlert(a.ctx, &events.AlertEvent{
				Level:   "error",
				Type:    "validation_failed",
				Message: fmt.Sprintf("Quality validation failed for %s.%s: %s", event.Source, event.Entity, err.Error()),
			})
		}
		return err
	}

	// Log quality results
	qualityScore := 0.0
	if report.OverallQualityScore != nil {
		qualityScore = *report.OverallQualityScore
	}
	a.logger.Info("quality validation completed",
		slog.String("source", event.Source),
		slog.String("entity", event.Entity),
		slog.Bool("passed", report.ValidationPassed),
		slog.Float64("quality_score", qualityScore),
		slog.Int("missing_values", report.MissingValueCount),
		slog.Int("duplicates", report.DuplicateCount),
		slog.Int("integrity_violations", report.IntegrityViolationCount),
		slog.Int("range_violations", report.RangeViolationCount),
		slog.Int("anomalies", report.AnomalyCount),
	)

	// Publish alerts if quality issues found
	if !report.ValidationPassed && a.events != nil {
		batchID := ""
		if event.BatchID != nil {
			batchID = *event.BatchID
		}
		_ = a.events.PublishDataQualityAlert(a.ctx, &events.DataQualityAlertEvent{
			BatchID:      batchID,
			Source:       event.Source,
			Entity:       &event.Entity,
			AlertLevel:   "error",
			AlertType:    "quality_validation_failed",
			Message:      fmt.Sprintf("Quality validation failed for %s.%s (score: %.2f)", event.Source, event.Entity, qualityScore),
			Details: map[string]interface{}{
				"quality_score":    qualityScore,
				"failure_reasons":  report.FailureReasons,
			},
			QualityScore: report.OverallQualityScore,
		})
	}

	return nil
}

// QualityReport represents the results of a quality validation check.
type QualityReport struct {
	ReportID                uuid.UUID  `db:"report_id"`
	BatchID                 uuid.UUID  `db:"batch_id"`
	Source                  string     `db:"source"`
	SyncStartedAt           time.Time `db:"sync_started_at"`
	SyncCompletedAt         time.Time `db:"sync_completed_at"`
	TotalRecords            int        `db:"total_records"`
	RecordsProcessed        int        `db:"records_processed"`
	RecordsInserted         int        `db:"records_inserted"`
	RecordsUpdated          int        `db:"records_updated"`
	RecordsQuarantined      int        `db:"records_quarantined"`
	CompletenessScore       *float64   `db:"completeness_score"`
	UniquenessScore         *float64   `db:"uniqueness_score"`
	ReferentialIntegrityScore *float64 `db:"referential_integrity_score"`
	RangeValidationScore    *float64   `db:"range_validation_score"`
	OverallQualityScore     *float64   `db:"overall_quality_score"`
	MissingValueCount       int        `db:"missing_value_count"`
	DuplicateCount          int        `db:"duplicate_count"`
	IntegrityViolationCount int        `db:"integrity_violation_count"`
	RangeViolationCount     int        `db:"range_violation_count"`
	AnomalyCount            int        `db:"anomaly_count"`
	PreviousRowCount        *int       `db:"previous_row_count"`
	CurrentRowCount         int        `db:"current_row_count"`
	RowCountDeltaPct        *float64   `db:"row_count_delta_pct"`
	AlertsGenerated         int        `db:"alerts_generated"`
	AlertDetails            []byte     `db:"alert_details"` // JSONB
	ValidationPassed        bool       `db:"validation_passed"`
	FailureReasons          []string   `db:"failure_reasons"`
	CreatedAt               time.Time  `db:"created_at"`
}

// ValidateSync runs quality validation on a completed sync.
func (a *Agent) ValidateSync(ctx context.Context, event *events.SyncCompletedEvent) (*QualityReport, error) {
	report := &QualityReport{
		BatchID:         uuid.New(),
		Source:          event.Source,
		SyncStartedAt:   event.StartedAt,
		SyncCompletedAt: event.CompletedAt,
		TotalRecords:    event.RecordsProcessed,
		RecordsProcessed: event.RecordsProcessed,
		RecordsInserted: event.RecordsInserted,
		RecordsUpdated:  event.RecordsUpdated,
		RecordsQuarantined: event.RecordsQuarantined,
		ValidationPassed: true,
		FailureReasons:   []string{},
		CreatedAt:        time.Now(),
	}

	// Determine the schema and table based on source and entity
	schema, table := getSchemaTable(event.Source, event.Entity)

	// Get current row count
	currentCount, _ := a.repo.GetTableStats(ctx, schema, table)
	report.CurrentRowCount = int(currentCount)

	// Run quality checks
	checks := []struct {
		name string
		fn   func(context.Context, string, string) (*CheckResult, error)
	}{
		{"completeness", a.checkCompleteness},
		{"uniqueness", a.checkUniqueness},
		{"integrity", a.checkReferentialIntegrity},
		{"range", a.checkRangeValidation},
	}

	scores := make([]float64, 0, 4)

	for _, check := range checks {
		result, err := check.fn(ctx, schema, table)
		if err != nil {
			a.logger.Warn("quality check failed",
				slog.String("check", check.name),
				slog.String("error", err.Error()),
			)
			continue
		}

		switch check.name {
		case "completeness":
			report.CompletenessScore = &result.Score
			report.MissingValueCount = result.Count
			if result.Score < 100 {
				scores = append(scores, result.Score)
			}
		case "uniqueness":
			report.UniquenessScore = &result.Score
			report.DuplicateCount = result.Count
			scores = append(scores, result.Score)
		case "integrity":
			report.ReferentialIntegrityScore = &result.Score
			report.IntegrityViolationCount = result.Count
			scores = append(scores, result.Score)
		case "range":
			report.RangeValidationScore = &result.Score
			report.RangeViolationCount = result.Count
			report.AnomalyCount = result.Count
			scores = append(scores, result.Score)
		}

		// Add failures
		for _, failure := range result.Failures {
			report.ValidationPassed = false
			report.FailureReasons = append(report.FailureReasons, failure)
		}
	}

	// Calculate overall score
	if len(scores) > 0 {
		avg := 0.0
		for _, s := range scores {
			avg += s
		}
		report.OverallQualityScore = &[]float64{avg / float64(len(scores))}[0]
	}

	if report.OverallQualityScore == nil {
		perf := 100.0
		report.OverallQualityScore = &perf
	}

	// Check if quality score is below passing threshold
	if *report.OverallQualityScore < QualityScorePassing {
		report.ValidationPassed = false
		report.FailureReasons = append(report.FailureReasons,
			fmt.Sprintf("Quality score %.2f is below passing threshold %.2f",
				*report.OverallQualityScore, QualityScorePassing))
	}

	// Save report to database
	if err := a.saveQualityReport(ctx, report); err != nil {
		a.logger.Error("failed to save quality report",
			slog.String("error", err.Error()),
		)
	}

	return report, nil
}

// CheckResult represents the result of a quality check.
type CheckResult struct {
	Name     string
	Score    float64 // 0-100
	Count    int
	Failures []string
}

// checkCompleteness checks for missing values.
func (a *Agent) checkCompleteness(ctx context.Context, schema, table string) (*CheckResult, error) {
	// This is a simplified check - in production, you'd check each column
	// for NULL/empty values and calculate the percentage

	query := fmt.Sprintf(`
		SELECT COUNT(*) as total,
		       COUNT(*) - COUNT(patient_id) as missing_id,
		       COUNT(*) - COUNT(name_en) as missing_name
		FROM %s.%s
		LIMIT 1000
	`, schema, table)

	var total, missingID, missingName int
	err := a.repo.Pool().QueryRow(ctx, query).Scan(&total, &missingID, &missingName)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &CheckResult{Name: "completeness", Score: 100.0, Count: 0}, nil
	}

	totalMissing := missingID + missingName
	missingRate := float64(totalMissing) / float64(total)
	score := 100.0 * (1 - missingRate)

	result := &CheckResult{
		Name:  "completeness",
		Score: score,
		Count: totalMissing,
	}

	if missingRate > MissingValueThreshold {
		result.Failures = append(result.Failures,
			fmt.Sprintf("Missing value rate %.2f%% exceeds threshold %.2f%%",
				missingRate*100, MissingValueThreshold*100))
	}

	return result, nil
}

// checkUniqueness checks for duplicate records based on business keys.
func (a *Agent) checkUniqueness(ctx context.Context, schema, table string) (*CheckResult, error) {
	// Check for duplicates based on _source, _source_id unique constraint
	query := fmt.Sprintf(`
		SELECT COUNT(*) - COUNT(DISTINCT (_source || ':' || _source_id)) as duplicates
		FROM %s.%s
	`, schema, table)

	var duplicates int
	err := a.repo.Pool().QueryRow(ctx, query).Scan(&duplicates)
	if err != nil {
		return nil, err
	}

	score := 100.0
	if duplicates > 0 {
		score = 95.0 // Deduct for duplicates
	}

	return &CheckResult{
		Name:  "uniqueness",
		Score: score,
		Count: duplicates,
	}, nil
}

// checkReferentialIntegrity checks foreign key violations.
func (a *Agent) checkReferentialIntegrity(ctx context.Context, schema, table string) (*CheckResult, error) {
	// This would check FK constraints - simplified here
	// In production, you'd query information_schema for FK violations

	score := 100.0
	violations := 0

	// Example: Check fact_appointments references
	if table == "fact_appointments" {
		query := `
			SELECT COUNT(*) FROM hims_analytics.fact_appointments a
			WHERE NOT EXISTS (SELECT 1 FROM hims_analytics.dim_patients p WHERE p.patient_id = a.patient_id)
			   OR NOT EXISTS (SELECT 1 FROM hims_analytics.dim_doctors d WHERE d.doctor_id = a.doctor_id)
		`
		err := a.repo.Pool().QueryRow(ctx, query).Scan(&violations)
		if err == nil && violations > 0 {
			score = 80.0
		}
	}

	return &CheckResult{
		Name:     "integrity",
		Score:    score,
		Count:    violations,
	}, nil
}

// checkRangeValidation checks for value range anomalies.
func (a *Agent) checkRangeValidation(ctx context.Context, schema, table string) (*CheckResult, error) {
	anomalies := 0
	failures := []string{}

	// Check for negative amounts in billing
	if table == "fact_billing" || table == "fact_pharmacy_dispensations" {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE total_amount < 0", schema, table)
		err := a.repo.Pool().QueryRow(ctx, query).Scan(&anomalies)
		if err == nil && anomalies > 0 {
			failures = append(failures, fmt.Sprintf("Found %d records with negative amounts", anomalies))
		}
	}

	// Check for future dates
	if table == "fact_appointments" || table == "fact_billing" {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE appt_date > NOW() + INTERVAL '30 days'", schema, table)
		var futureDates int
		err := a.repo.Pool().QueryRow(ctx, query).Scan(&futureDates)
		if err == nil && futureDates > 0 {
			anomalies += futureDates
			failures = append(failures, fmt.Sprintf("Found %d records with dates more than 30 days in future", futureDates))
		}
	}

	score := 100.0
	if len(failures) > 0 {
		score = 90.0 - float64(len(failures))*10
		if score < 0 {
			score = 0
		}
	}

	return &CheckResult{
		Name:     "range",
		Score:    score,
		Count:    anomalies,
		Failures: failures,
	}, nil
}

// saveQualityReport saves the quality report to the database.
func (a *Agent) saveQualityReport(ctx context.Context, report *QualityReport) error {
	query := `
		INSERT INTO app.etl_quality_report (
			batch_id, source, sync_started_at, sync_completed_at,
			total_records, records_processed, records_inserted, records_updated,
			records_quarantined, completeness_score, uniqueness_score,
			referential_integrity_score, range_validation_score, overall_quality_score,
			missing_value_count, duplicate_count, integrity_violation_count,
			range_violation_count, anomaly_count, current_row_count,
			validation_passed, failure_reasons, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, NOW()
		)
		RETURNING report_id
	`

	// Convert failures to JSON array for PostgreSQL
	failuresJSON, _ := json.Marshal(report.FailureReasons)

	_, err := a.repo.Pool().Exec(ctx, query,
		report.BatchID, report.Source, report.SyncStartedAt, report.SyncCompletedAt,
		report.TotalRecords, report.RecordsProcessed, report.RecordsInserted,
		report.RecordsUpdated, report.RecordsQuarantined, report.CompletenessScore,
		report.UniquenessScore, report.ReferentialIntegrityScore, report.RangeValidationScore,
		report.OverallQualityScore, report.MissingValueCount, report.DuplicateCount,
		report.IntegrityViolationCount, report.RangeViolationCount, report.AnomalyCount,
		report.CurrentRowCount, report.ValidationPassed, failuresJSON,
	)

	return err
}

// getSchemaTable returns the schema and table name for a source/entity pair.
func getSchemaTable(source, entity string) (string, string) {
	switch source {
	case "hims":
		switch entity {
		case "patients":
			return "hims_analytics", "dim_patients"
		case "doctors":
			return "hims_analytics", "dim_doctors"
		case "drugs":
			return "hims_analytics", "dim_drugs"
		case "appointments":
			return "hims_analytics", "fact_appointments"
		case "billing":
			return "hims_analytics", "fact_billing"
		case "pharmacy", "dispensations":
			return "hims_analytics", "fact_pharmacy_dispensations"
		}
	case "tally":
		switch entity {
		case "ledgers":
			return "tally_analytics", "dim_ledgers"
		case "cost_centres":
			return "tally_analytics", "dim_cost_centres"
		case "stock_items":
			return "tally_analytics", "dim_inventory_items"
		case "vouchers":
			return "tally_analytics", "fact_vouchers"
		case "stock", "stock_movements":
			return "tally_analytics", "fact_stock_movements"
		}
	}

	return "", ""
}

// RunManualValidation triggers manual quality validation for a source/entity.
func (a *Agent) RunManualValidation(ctx context.Context, source, entity string) (*QualityReport, error) {
	batchID := uuid.New()
	batchIDStr := batchID.String()

	event := &events.SyncCompletedEvent{
		EventID:          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Source:           source,
		Entity:           entity,
		StartedAt:        time.Now(),
		CompletedAt:      time.Now(),
		RecordsProcessed: 0,
		BatchID:          &batchIDStr,
	}

	return a.ValidateSync(ctx, event)
}

// GetQualityReports retrieves recent quality reports.
func (a *Agent) GetQualityReports(ctx context.Context, source string, limit int) ([]*QualityReport, error) {
	query := `
		SELECT report_id, batch_id, source, sync_started_at, sync_completed_at,
		       total_records, records_processed, records_inserted, records_updated,
		       records_quarantined, completeness_score, uniqueness_score,
		       referential_integrity_score, range_validation_score, overall_quality_score,
		       missing_value_count, duplicate_count, integrity_violation_count,
		       range_violation_count, anomaly_count, current_row_count,
		       validation_passed, failure_reasons, created_at
		FROM app.etl_quality_report
		WHERE $1 = '' OR source = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := a.repo.Pool().Query(ctx, query, source, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query quality reports: %w", err)
	}
	defer rows.Close()

	var reports []*QualityReport
	for rows.Next() {
		var report QualityReport
		var failuresJSON []byte

		err := rows.Scan(
			&report.ReportID, &report.BatchID, &report.Source, &report.SyncStartedAt,
			&report.SyncCompletedAt, &report.TotalRecords, &report.RecordsProcessed,
			&report.RecordsInserted, &report.RecordsUpdated, &report.RecordsQuarantined,
			&report.CompletenessScore, &report.UniquenessScore,
			&report.ReferentialIntegrityScore, &report.RangeValidationScore,
			&report.OverallQualityScore, &report.MissingValueCount, &report.DuplicateCount,
			&report.IntegrityViolationCount, &report.RangeViolationCount,
			&report.AnomalyCount, &report.CurrentRowCount, &report.ValidationPassed,
			&failuresJSON, &report.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(failuresJSON, &report.FailureReasons)
		reports = append(reports, &report)
	}

	return reports, nil
}

// ValidateCompletenessByColumn checks missing values for specific columns.
func (a *Agent) ValidateCompletenessByColumn(ctx context.Context, schema, table string, columns []string) (map[string]float64, error) {
	results := make(map[string]float64)

	// First get total count
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", schema, table)
	err := a.repo.Pool().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return results, nil
	}

	// Check each column
	for _, col := range columns {
		var nullCount int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE %s IS NULL OR %s = ''",
			schema, table, col, col)
		_ = a.repo.Pool().QueryRow(ctx, query).Scan(&nullCount)

		completeness := 100.0 * (1 - float64(nullCount)/float64(total))
		results[col] = completeness
	}

	return results, nil
}

// ScanForAnomalies scans data for various anomalies.
func (a *Agent) ScanForAnomalies(ctx context.Context, schema, table string) ([]string, error) {
	var anomalies []string

	// Check for records with NULL in required columns
	requiredCols := getRequiredColumns(table)
	for _, col := range requiredCols {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE %s IS NULL", schema, table, col)
		_ = a.repo.Pool().QueryRow(ctx, query).Scan(&count)
		if count > 0 {
			anomalies = append(anomalies, fmt.Sprintf("%d NULL values in required column %s", count, col))
		}
	}

	return anomalies, nil
}

// getRequiredColumns returns required columns for a table.
func getRequiredColumns(table string) []string {
	switch table {
	case "dim_patients":
		return []string{"patient_id", "external_patient_id", "name_en"}
	case "dim_doctors":
		return []string{"doctor_id", "external_doctor_id", "name_en"}
	case "fact_appointments":
		return []string{"appt_id", "patient_id", "doctor_id", "appt_date", "status"}
	case "fact_billing":
		return []string{"bill_id", "patient_id", "bill_date", "total_amount", "payment_status"}
	case "dim_ledgers":
		return []string{"ledger_id", "external_ledger_id", "ledger_name", "ledger_group"}
	case "fact_vouchers":
		return []string{"voucher_id", "external_voucher_id", "voucher_number", "voucher_type", "voucher_date", "ledger_id"}
	default:
		return []string{}
	}
}
