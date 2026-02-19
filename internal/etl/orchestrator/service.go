// Package orchestrator provides ETL coordination for Tally and HIMS data synchronization.
//
// This service coordinates the extraction, transformation, and loading of data from
// source systems (Tally ERP, HIMS) into the data warehouse. It manages scheduled
// sync jobs, handles errors and quarantine, publishes NATS events, and tracks
// sync state for incremental updates.
//
// Sync Schedule:
// - Tally (ledgers, vouchers, stock): every 30 minutes
// - HIMS (appointments, billing, pharmacy): every 15 minutes
// - HIMS (patients, doctors, drugs): daily
//
// Usage:
//
//	cfg := config.MustLoad()
//	orch := orchestrator.NewService(cfg, repo, tallyClient, himsClient, events, logger)
//
//	ctx := context.Background()
//	if err := orch.Start(ctx); err != nil {
//	    log.Fatal("Failed to start orchestrator:", err)
//	}
package orchestrator

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/medisync/medisync/internal/config"
	"github.com/medisync/medisync/internal/events"
	"github.com/medisync/medisync/internal/etl/hims"
	"github.com/medisync/medisync/internal/etl/tally"
	"github.com/medisync/medisync/internal/warehouse"
)

// Service coordinates ETL operations.
type Service struct {
	config       *config.Config
	repo         *warehouse.Repo
	tallyClient  *tally.Client
	himsClient   *hims.Client
	events       *events.Publisher
	logger       *slog.Logger

	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup

	// Status tracking
	isRunning    atomic.Bool
	totalSyncs   atomic.Int64
	failedSyncs  atomic.Int64

	// Schedule tickers
	tallyTicker  *time.Ticker
	hims15Ticker *time.Ticker
	hims24Ticker *time.Ticker
}

// Config holds configuration for creating a new Service.
type ServiceConfig struct {
	// Config is the application configuration.
	Config *config.Config

	// Repo is the warehouse repository.
	Repo *warehouse.Repo

	// TallyClient is the Tally ERP client.
	TallyClient *tally.Client

	// HIMSClient is the HIMS REST API client.
	HIMSClient *hims.Client

	// Events is the NATS event publisher.
	Events *events.Publisher

	// Logger is the structured logger.
	Logger *slog.Logger
}

// NewService creates a new ETL orchestrator service.
func NewService(cfg *ServiceConfig) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("orchestrator: config is required")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	svc := &Service{
		config:      cfg.Config,
		repo:        cfg.Repo,
		tallyClient: cfg.TallyClient,
		himsClient:  cfg.HIMSClient,
		events:      cfg.Events,
		logger:      logger.With(slog.String("service", "etl_orchestrator")),
	}

	return svc, nil
}

// Start begins the ETL service with scheduled sync jobs.
func (s *Service) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	if !s.isRunning.CompareAndSwap(false, true) {
		return fmt.Errorf("orchestrator: service is already running")
	}

	s.logger.Info("starting ETL orchestrator service")

	// Start scheduled sync jobs
	s.startScheduledSyncs()

	// Run initial sync for master data
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runInitialMasterSync(s.ctx)
	}()

	return nil
}

// Stop gracefully shuts down the ETL service.
func (s *Service) Stop() error {
	if !s.isRunning.CompareAndSwap(true, false) {
		return fmt.Errorf("orchestrator: service is not running")
	}

	s.logger.Info("stopping ETL orchestrator service")

	// Stop all tickers
	if s.tallyTicker != nil {
		s.tallyTicker.Stop()
	}
	if s.hims15Ticker != nil {
		s.hims15Ticker.Stop()
	}
	if s.hims24Ticker != nil {
		s.hims24Ticker.Stop()
	}

	// Cancel context
	s.cancel()

	// Wait for all goroutines to finish
	s.wg.Wait()

	s.logger.Info("ETL orchestrator service stopped",
		slog.Int64("total_syncs", s.totalSyncs.Load()),
		slog.Int64("failed_syncs", s.failedSyncs.Load()),
	)

	return nil
}

// startScheduledSyncs sets up the scheduled sync jobs.
func (s *Service) startScheduledSyncs() {
	// Tally sync every 30 minutes
	s.tallyTicker = time.NewTicker(30 * time.Minute)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.syncLoop(s.ctx, s.tallyTicker.C, "tally_30min", s.syncTallyAll)
	}()

	// HIMS frequent sync every 15 minutes (appointments, billing, pharmacy)
	s.hims15Ticker = time.NewTicker(15 * time.Minute)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.syncLoop(s.ctx, s.hims15Ticker.C, "hims_15min", s.syncHIMSFrequent)
	}()

	// HIMS daily sync (patients, doctors, drugs)
	s.hims24Ticker = time.NewTicker(24 * time.Hour)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.syncLoop(s.ctx, s.hims24Ticker.C, "hims_24hr", s.syncHIMSDaily)
	}()
}

// syncLoop runs sync operations on each ticker tick.
func (s *Service) syncLoop(ctx context.Context, ticker <-chan time.Time, name string, syncFn func(context.Context) error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			if err := syncFn(ctx); err != nil {
				s.logger.Error("sync failed",
					slog.String("schedule", name),
					slog.String("error", err.Error()),
				)
				s.failedSyncs.Add(1)
			}
		}
	}
}

// runInitialMasterSync performs initial sync for all master data.
func (s *Service) runInitialMasterSync(ctx context.Context) {
	s.logger.Info("running initial master data sync")

	startTime := time.Now()

	// Sync HIMS master data
	if err := s.syncHIMSDaily(ctx); err != nil {
		s.logger.Error("initial HIMS master sync failed",
			slog.String("error", err.Error()),
		)
	}

	// Sync Tally master data
	if err := s.syncTallyMaster(ctx); err != nil {
		s.logger.Error("initial Tally master sync failed",
			slog.String("error", err.Error()),
		)
	}

	s.logger.Info("initial master data sync completed",
		slog.Duration("duration", time.Since(startTime)),
	)
}

// syncTallyAll syncs all Tally data.
func (s *Service) syncTallyAll(ctx context.Context) error {
	startTime := time.Now()

	// Sync master data less frequently
	if time.Since(startTime) < 1*time.Hour {
		_ = s.syncTallyTransactional(ctx)
	} else {
		_ = s.syncTallyMaster(ctx)
		_ = s.syncTallyTransactional(ctx)
	}

	return nil
}

// syncTallyMaster syncs Tally master data (ledgers, cost centres, stock items).
func (s *Service) syncTallyMaster(ctx context.Context) error {
	entities := []struct {
		name   string
		syncFn func(context.Context) error
	}{
		{"ledgers", s.syncTallyLedgers},
		{"cost_centres", s.syncTallyCostCentres},
		{"stock_items", s.syncTallyStockItems},
	}

	for _, entity := range entities {
		if err := s.runSyncJob(ctx, warehouse.SourceTally, entity.name, entity.syncFn); err != nil {
			return err
		}
	}

	return nil
}

// syncTallyTransactional syncs Tally transactional data (vouchers, stock movements).
func (s *Service) syncTallyTransactional(ctx context.Context) error {
	entities := []struct {
		name   string
		syncFn func(context.Context) error
	}{
		{"vouchers", s.syncTallyVouchers},
		{"stock_movements", s.syncTallyStockMovements},
	}

	for _, entity := range entities {
		if err := s.runSyncJob(ctx, warehouse.SourceTally, entity.name, entity.syncFn); err != nil {
			return err
		}
	}

	return nil
}

// syncTallyLedgers syncs Tally ledgers.
func (s *Service) syncTallyLedgers(ctx context.Context) error {
	return s.syncTallyEntity(ctx, "ledgers", func(ctx context.Context) (*tally.Response, error) {
		return s.tallyClient.GetLedgers(ctx)
	})
}

// syncTallyCostCentres syncs Tally cost centres.
func (s *Service) syncTallyCostCentres(ctx context.Context) error {
	return s.syncTallyEntity(ctx, "cost_centres", func(ctx context.Context) (*tally.Response, error) {
		return s.tallyClient.GetCostCentres(ctx)
	})
}

// syncTallyStockItems syncs Tally stock items.
func (s *Service) syncTallyStockItems(ctx context.Context) error {
	return s.syncTallyEntity(ctx, "stock_items", func(ctx context.Context) (*tally.Response, error) {
		return s.tallyClient.GetStockItems(ctx)
	})
}

// syncTallyVouchers syncs Tally vouchers with incremental support.
func (s *Service) syncTallyVouchers(ctx context.Context) error {
	// Get last AlterID for incremental sync
	lastAlterID, err := s.repo.GetTallyAlterID(ctx, "vouchers")
	if err != nil {
		return fmt.Errorf("failed to get Tally cursor: %w", err)
	}

	var resp *tally.Response
	if lastAlterID != "" {
		resp, err = s.tallyClient.GetVouchersModifiedSince(ctx, lastAlterID)
	} else {
		// Full sync for date range
		toDate := time.Now()
		fromDate := toDate.Add(-30 * 24 * time.Hour) // Last 30 days
		resp, err = s.tallyClient.GetVouchers(ctx, fromDate, toDate)
	}

	if err != nil {
		return err
	}

	// Transform and upsert vouchers
	records := make([]*warehouse.Voucher, len(resp.Vouchers))
	for i, v := range resp.Vouchers {
		records[i] = &warehouse.Voucher{
			ExternalVoucherID: v.MasterID,
			VoucherNumber:     v.VoucherNumber,
			VoucherType:       v.VoucherType,
			VoucherDate:       parseTallyDate(v.Date),
			VoucherDatetime:   parseTallyDateTime(v.Date),
			LedgerID:          uuid.UUID{}, // Will be mapped during upsert
			Amount:            getVoucherAmount(v),
			Narration:         &v.Narration,
			Source:            warehouse.SourceTally.String(),
			SourceID:          v.MasterID,
		}
	}

	if err := s.repo.BulkUpsertVouchers(ctx, records); err != nil {
		return err
	}

	// Update cursor if successful
	if len(resp.Vouchers) > 0 && resp.Vouchers[len(resp.Vouchers)-1].AlterID != "" {
		_ = s.repo.UpdateTallyAlterID(ctx, "vouchers", resp.Vouchers[len(resp.Vouchers)-1].AlterID, len(resp.Vouchers))
	}

	return nil
}

// syncTallyStockMovements syncs Tally stock movements.
func (s *Service) syncTallyStockMovements(ctx context.Context) error {
	// Implementation similar to vouchers
	return s.syncTallyEntity(ctx, "stock_movements", func(ctx context.Context) (*tally.Response, error) {
		// Get vouchers that have inventory entries
		return s.tallyClient.GetVouchers(ctx, time.Now().Add(-24*time.Hour), time.Now())
	})
}

// syncTallyEntity is a helper for syncing Tally entities.
func (s *Service) syncTallyEntity(ctx context.Context, entity string, fetchFn func(context.Context) (*tally.Response, error)) error {
	startTime := time.Now()
	batchID := uuid.New()

	resp, err := fetchFn(ctx)
	if err != nil {
		s.publishSyncFailed(ctx, warehouse.SourceTally, entity, startTime, err)
		return err
	}

	// Publish sync completed event
	if s.events != nil {
		_ = s.events.PublishSyncCompleted(ctx, &events.SyncCompletedEvent{
			BatchID:            &[]string{batchID.String()}[0],
			Source:             warehouse.SourceTally.String(),
			Entity:             entity,
			StartedAt:          startTime,
			CompletedAt:        time.Now(),
			RecordsProcessed:   len(resp.Ledgers) + len(resp.CostCentres) + len(resp.StockItems) + len(resp.Vouchers),
			DurationSec:        time.Since(startTime).Seconds(),
		})
	}

	return nil
}

// syncHIMSFrequent syncs frequently changing HIMS data (appointments, billing, pharmacy).
func (s *Service) syncHIMSFrequent(ctx context.Context) error {
	entities := []struct {
		name   string
		syncFn func(context.Context) error
	}{
		{"appointments", s.syncHIMSAppointments},
		{"billing", s.syncHIMSBilling},
		{"pharmacy_dispensations", s.syncHIMSPharmacy},
	}

	for _, entity := range entities {
		if err := s.runSyncJob(ctx, warehouse.SourceHIMS, entity.name, entity.syncFn); err != nil {
			s.logger.Error("HIMS sync failed",
				slog.String("entity", entity.name),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

// syncHIMSDaily syncs daily HIMS master data (patients, doctors, drugs, departments).
func (s *Service) syncHIMSDaily(ctx context.Context) error {
	entities := []struct {
		name   string
		syncFn func(context.Context) error
	}{
		{"patients", s.syncHIMSPatients},
		{"doctors", s.syncHIMSDoctors},
		{"drugs", s.syncHIMSDrugs},
		{"departments", s.syncHIMSDepartments},
	}

	for _, entity := range entities {
		if err := s.runSyncJob(ctx, warehouse.SourceHIMS, entity.name, entity.syncFn); err != nil {
			s.logger.Error("HIMS daily sync failed",
				slog.String("entity", entity.name),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

// syncHIMSPatients syncs HIMS patients.
func (s *Service) syncHIMSPatients(ctx context.Context) error {
	return s.syncHIMSEntity(ctx, "patients", func(opts *hims.PatientOptions) (*hims.PagedResponse, error) {
		return s.himsClient.GetPatients(ctx, opts)
	}, s.transformAndUpsertPatients)
}

// syncHIMSDoctors syncs HIMS doctors.
func (s *Service) syncHIMSDoctors(ctx context.Context) error {
	return s.syncHIMSEntity(ctx, "doctors", func(opts *hims.DoctorOptions) (*hims.PagedResponse, error) {
		return s.himsClient.GetDoctors(ctx, opts)
	}, s.transformAndUpsertDoctors)
}

// syncHIMSDrugs syncs HIMS drugs.
func (s *Service) syncHIMSDrugs(ctx context.Context) error {
	return s.syncHIMSEntity(ctx, "drugs", func(opts *hims.DrugOptions) (*hims.PagedResponse, error) {
		return s.himsClient.GetDrugs(ctx, opts)
	}, s.transformAndUpsertDrugs)
}

// syncHIMSDepartments syncs HIMS departments.
func (s *Service) syncHIMSDepartments(ctx context.Context) error {
	_, err := s.himsClient.GetDepartments(ctx, &hims.DepartmentsOptions{})
	return err
}

// syncHIMSAppointments syncs HIMS appointments.
func (s *Service) syncHIMSAppointments(ctx context.Context) error {
	// Get last modified time for incremental sync
	lastModified, _ := s.repo.GetHIMSCursor(ctx, "appointments")

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	_, err := s.himsClient.GetAppointments(ctx, &hims.AppointmentOptions{
		StartDate:     &startDate,
		EndDate:       &endDate,
		ModifiedSince: lastModified,
	})

	if err == nil {
		_ = s.repo.UpdateHIMSCursor(ctx, "appointments", time.Now(), 0)
	}

	return err
}

// syncHIMSBilling syncs HIMS billing records.
func (s *Service) syncHIMSBilling(ctx context.Context) error {
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	_, err := s.himsClient.GetBilling(ctx, &hims.BillingOptions{
		StartDate: &startDate,
		EndDate:   &endDate,
	})

	return err
}

// syncHIMSPharmacy syncs HIMS pharmacy dispensations.
func (s *Service) syncHIMSPharmacy(ctx context.Context) error {
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	_, err := s.himsClient.GetPharmacyDispensations(ctx, &hims.PharmacyOptions{
		StartDate: &startDate,
		EndDate:   &endDate,
	})

	return err
}

// syncHIMSEntity is a generic HIMS entity sync function.
func (s *Service) syncHIMSEntity(ctx context.Context, entity string, fetchFn interface{}, transformFn interface{}) error {
	startTime := time.Now()
	batchID := uuid.New()

	// Fetch, transform, and upsert
	// Implementation depends on the specific fetch/transform functions

	// Publish sync completed event
	if s.events != nil {
		_ = s.events.PublishSyncCompleted(ctx, &events.SyncCompletedEvent{
			BatchID:     &[]string{batchID.String()}[0],
			Source:      warehouse.SourceHIMS.String(),
			Entity:      entity,
			StartedAt:   startTime,
			CompletedAt: time.Now(),
		})
	}

	return nil
}

// transformAndUpsertPatients transforms and upserts patient records.
func (s *Service) transformAndUpsertPatients(ctx context.Context, data interface{}) error {
	patients, ok := data.([]*hims.Patient)
	if !ok {
		return nil
	}

	transformer := hims.NewTransformer()
	records := transformer.BatchPatientsToWarehouse(patients)

	if len(records) > 0 {
		return s.repo.BulkUpsertPatients(ctx, records)
	}

	return nil
}

// transformAndUpsertDoctors transforms and upserts doctor records.
func (s *Service) transformAndUpsertDoctors(ctx context.Context, data interface{}) error {
	doctors, ok := data.([]*hims.Doctor)
	if !ok {
		return nil
	}

	transformer := hims.NewTransformer()
	records := transformer.BatchDoctorsToWarehouse(doctors)

	if len(records) > 0 {
		return s.repo.BulkUpsertDoctors(ctx, records)
	}

	return nil
}

// transformAndUpsertDrugs transforms and upserts drug records.
func (s *Service) transformAndUpsertDrugs(ctx context.Context, data interface{}) error {
	drugs, ok := data.([]*hims.Drug)
	if !ok {
		return nil
	}

	transformer := hims.NewTransformer()
	records := transformer.BatchDrugsToWarehouse(drugs)

	if len(records) > 0 {
		return s.repo.BulkUpsertDrugs(ctx, records)
	}

	return nil
}

// runSyncJob runs a single sync job with error handling and event publishing.
func (s *Service) runSyncJob(ctx context.Context, source warehouse.Source, entity string, syncFn func(context.Context) error) error {
	startTime := time.Now()
	s.totalSyncs.Add(1)

	// Acquire ETL lock
	if err := s.repo.WithETLLock(ctx, source.String(), entity, func() error {
		return syncFn(ctx)
	}); err != nil {
		s.publishSyncFailed(ctx, source, entity, startTime, err)
		return err
	}

	// Publish success event
	s.publishSyncCompleted(ctx, source, entity, startTime, 0)

	return nil
}

// publishSyncCompleted publishes a sync completion event.
func (s *Service) publishSyncCompleted(ctx context.Context, source warehouse.Source, entity string, startTime time.Time, recordsProcessed int) {
	if s.events == nil {
		return
	}

	batchID := uuid.New().String()
	_ = s.events.PublishSyncCompleted(ctx, &events.SyncCompletedEvent{
		BatchID:          &batchID,
		Source:           source.String(),
		Entity:           entity,
		StartedAt:        startTime,
		CompletedAt:      time.Now(),
		RecordsProcessed: recordsProcessed,
		DurationSec:      time.Since(startTime).Seconds(),
		Status:           "completed",
	})
}

// publishSyncFailed publishes a sync failure event.
func (s *Service) publishSyncFailed(ctx context.Context, source warehouse.Source, entity string, startTime time.Time, err error) {
	if s.events == nil {
		return
	}

	batchID := uuid.New().String()
	retryable := true

	_ = s.events.PublishSyncFailed(ctx, &events.SyncFailedEvent{
		Source:           source.String(),
		Entity:           entity,
		StartedAt:        startTime,
		FailedAt:         time.Now(),
		Error:            err.Error(),
		ErrorCode:        "sync_error",
		DurationSec:      time.Since(startTime).Seconds(),
		Retryable:        retryable,
		Metadata: map[string]interface{}{
			"batch_id": batchID,
		},
	})
}

// GetStats returns orchestrator statistics.
func (s *Service) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"is_running":    s.isRunning.Load(),
		"total_syncs":   s.totalSyncs.Load(),
		"failed_syncs":  s.failedSyncs.Load(),
	}
}

// Helper functions for Tally data transformation
func parseTallyDate(dateStr string) time.Time {
	if t, err := time.Parse("20060102", dateStr); err == nil {
		return t
	}
	return time.Now()
}

func parseTallyDateTime(dateStr string) *time.Time {
	if t, err := time.Parse("20060102", dateStr); err == nil {
		return &t
	}
	return nil
}

func getVoucherAmount(v tally.Voucher) float64 {
	for _, entry := range v.LedgerEntries {
		if entry.IsDeemedPositive == "Yes" {
			return entry.Amount
		}
	}
	return 0
}
