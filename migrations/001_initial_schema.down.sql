-- MediSync Initial Schema Migration Rollback
-- Version: 001
-- Description: Drop all schemas and tables created by 001_initial_schema.up.sql
--
-- This rollback drops:
-- 1. All triggers
-- 2. All functions
-- 3. All tables (in reverse dependency order)
-- 4. All schemas
-- 5. Extensions (commented out - typically not removed as they may be shared)
--
-- WARNING: This will permanently delete all data in these schemas!

-- ============================================================================
-- DROP TRIGGERS
-- Note: Triggers are dropped automatically when tables are dropped,
-- but we explicitly drop them first for clarity
-- ============================================================================

-- Drop update_updated_at triggers from all schemas
DO $$
DECLARE
    t record;
BEGIN
    FOR t IN
        SELECT tgname, relname, nspname
        FROM pg_trigger
        JOIN pg_class ON pg_trigger.tgrelid = pg_class.oid
        JOIN pg_namespace ON pg_class.relnamespace = pg_namespace.oid
        WHERE nspname IN ('hims_analytics', 'tally_analytics', 'app', 'vectors')
        AND tgname LIKE 'update_%_updated_at'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON %I.%I', t.tgname, t.nspname, t.relname);
    END LOOP;
END;
$$;

-- ============================================================================
-- DROP FUNCTIONS
-- ============================================================================

DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- ============================================================================
-- SCHEMA: vectors
-- Drop in reverse order of creation, respecting foreign key dependencies
-- ============================================================================

-- Drop indexes first (optional, as they're dropped with tables)
DROP INDEX IF EXISTS vectors.idx_document_embeddings_type;
DROP INDEX IF EXISTS vectors.idx_document_embeddings_doc;
DROP INDEX IF EXISTS vectors.idx_document_embeddings_vector;
DROP INDEX IF EXISTS vectors.idx_query_history_created;
DROP INDEX IF EXISTS vectors.idx_query_history_feedback;
DROP INDEX IF EXISTS vectors.idx_query_history_success;
DROP INDEX IF EXISTS vectors.idx_query_history_user;
DROP INDEX IF EXISTS vectors.idx_query_history_vector;
DROP INDEX IF EXISTS vectors.idx_metric_embeddings_type;
DROP INDEX IF EXISTS vectors.idx_metric_embeddings_domain;
DROP INDEX IF EXISTS vectors.idx_metric_embeddings_vector;
DROP INDEX IF EXISTS vectors.idx_schema_embeddings_table;
DROP INDEX IF EXISTS vectors.idx_schema_embeddings_vector;

-- Drop tables
DROP TABLE IF EXISTS vectors.document_embeddings CASCADE;
DROP TABLE IF EXISTS vectors.query_history CASCADE;
DROP TABLE IF EXISTS vectors.metric_embeddings CASCADE;
DROP TABLE IF EXISTS vectors.schema_embeddings CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS vectors CASCADE;

-- ============================================================================
-- SCHEMA: app
-- Drop in reverse order of creation, respecting foreign key dependencies
-- ============================================================================

-- Drop indexes
DROP INDEX IF EXISTS app.idx_approval_workflows_created;
DROP INDEX IF EXISTS app.idx_approval_workflows_due_date;
DROP INDEX IF EXISTS app.idx_approval_workflows_created_by;
DROP INDEX IF EXISTS app.idx_approval_workflows_type;
DROP INDEX IF EXISTS app.idx_approval_workflows_status;


DROP INDEX IF EXISTS app.idx_notification_queue_type;
DROP INDEX IF EXISTS app.idx_notification_queue_scheduled;
DROP INDEX IF EXISTS app.idx_notification_queue_read;
DROP INDEX IF EXISTS app.idx_notification_queue_user;

DROP INDEX IF EXISTS app.idx_audit_log_session;
DROP INDEX IF EXISTS app.idx_audit_log_created;
DROP INDEX IF EXISTS app.idx_audit_log_resource;
DROP INDEX IF EXISTS app.idx_audit_log_action;
DROP INDEX IF EXISTS app.idx_audit_log_user;
DROP INDEX IF EXISTS app.idx_etl_state_status;
DROP INDEX IF EXISTS app.idx_etl_state_source;
DROP INDEX IF EXISTS app.idx_etl_quality_report_passed;
DROP INDEX IF EXISTS app.idx_etl_quality_report_created;
DROP INDEX IF EXISTS app.idx_etl_quality_report_source;
DROP INDEX IF EXISTS app.idx_etl_quality_report_batch;
DROP INDEX IF EXISTS app.idx_etl_quarantine_batch;
DROP INDEX IF EXISTS app.idx_etl_quarantine_created;
DROP INDEX IF EXISTS app.idx_etl_quarantine_status;
DROP INDEX IF EXISTS app.idx_etl_quarantine_source;

DROP INDEX IF EXISTS app.idx_users_department;
DROP INDEX IF EXISTS app.idx_users_role;
DROP INDEX IF EXISTS app.idx_users_keycloak_sub;
DROP INDEX IF EXISTS app.idx_users_email;

-- Drop tables (order matters due to FK constraints)
DROP TABLE IF EXISTS app.approval_workflows CASCADE;


DROP TABLE IF EXISTS app.notification_queue CASCADE;
DROP TABLE IF EXISTS app.etl_state CASCADE;
DROP TABLE IF EXISTS app.etl_quarantine CASCADE;
DROP TABLE IF EXISTS app.audit_log CASCADE;
DROP TABLE IF EXISTS app.users CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS app CASCADE;

-- ============================================================================
-- SCHEMA: tally_analytics
-- Drop in reverse order of creation, respecting foreign key dependencies
-- ============================================================================

-- Drop indexes
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_godown;
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_synced_at;
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_voucher;
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_type;
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_date;
DROP INDEX IF EXISTS tally_analytics.idx_fact_stock_movements_item;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_voucher_number;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_party;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_reference;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_synced_at;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_cost_centre;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_contra_ledger;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_ledger;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_type;
DROP INDEX IF EXISTS tally_analytics.idx_fact_vouchers_date;
DROP INDEX IF EXISTS tally_analytics.idx_dim_inventory_items_synced_at;
DROP INDEX IF EXISTS tally_analytics.idx_dim_inventory_items_part_number;
DROP INDEX IF EXISTS tally_analytics.idx_dim_inventory_items_stock_group;
DROP INDEX IF EXISTS tally_analytics.idx_dim_inventory_items_category;
DROP INDEX IF EXISTS tally_analytics.idx_dim_inventory_items_name;
DROP INDEX IF EXISTS tally_analytics.idx_dim_cost_centres_synced_at;
DROP INDEX IF EXISTS tally_analytics.idx_dim_cost_centres_code;
DROP INDEX IF EXISTS tally_analytics.idx_dim_cost_centres_name;
DROP INDEX IF EXISTS tally_analytics.idx_dim_ledgers_external_id;
DROP INDEX IF EXISTS tally_analytics.idx_dim_ledgers_synced_at;
DROP INDEX IF EXISTS tally_analytics.idx_dim_ledgers_type;
DROP INDEX IF EXISTS tally_analytics.idx_dim_ledgers_group;
DROP INDEX IF EXISTS tally_analytics.idx_dim_ledgers_name;

-- Drop tables (order matters due to FK constraints)
DROP TABLE IF EXISTS tally_analytics.fact_stock_movements CASCADE;
DROP TABLE IF EXISTS tally_analytics.fact_vouchers CASCADE;
DROP TABLE IF EXISTS tally_analytics.dim_inventory_items CASCADE;
DROP TABLE IF EXISTS tally_analytics.dim_cost_centres CASCADE;
DROP TABLE IF EXISTS tally_analytics.dim_ledgers CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS tally_analytics CASCADE;

-- ============================================================================
-- SCHEMA: hims_analytics
-- Drop in reverse order of creation, respecting foreign key dependencies
-- ============================================================================

-- Drop indexes
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_doctor;
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_prescription;
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_date;
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_patient;
DROP INDEX IF EXISTS hims_analytics.idx_fact_pharmacy_disp_drug;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_bill_type;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_department;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_payment_mode;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_payment_status;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_date;
DROP INDEX IF EXISTS hims_analytics.idx_fact_billing_patient;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_datetime;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_department;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_status;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_date;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_doctor;
DROP INDEX IF EXISTS hims_analytics.idx_fact_appointments_patient;
DROP INDEX IF EXISTS hims_analytics.idx_dim_departments_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_dim_departments_code;
DROP INDEX IF EXISTS hims_analytics.idx_dim_drugs_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_dim_drugs_category;
DROP INDEX IF EXISTS hims_analytics.idx_dim_drugs_name_en;
DROP INDEX IF EXISTS hims_analytics.idx_dim_drugs_external_id;
DROP INDEX IF EXISTS hims_analytics.idx_dim_doctors_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_dim_doctors_department;
DROP INDEX IF EXISTS hims_analytics.idx_dim_doctors_specialty;
DROP INDEX IF EXISTS hims_analytics.idx_dim_doctors_name_en;
DROP INDEX IF EXISTS hims_analytics.idx_dim_doctors_external_id;
DROP INDEX IF EXISTS hims_analytics.idx_dim_patients_synced_at;
DROP INDEX IF EXISTS hims_analytics.idx_dim_patients_phone;
DROP INDEX IF EXISTS hims_analytics.idx_dim_patients_name_ar;
DROP INDEX IF EXISTS hims_analytics.idx_dim_patients_name_en;
DROP INDEX IF EXISTS hims_analytics.idx_dim_patients_external_id;

-- Drop tables (order matters due to FK constraints)
-- Fact tables first (they reference dimension tables)
DROP TABLE IF EXISTS hims_analytics.fact_pharmacy_dispensations CASCADE;
DROP TABLE IF EXISTS hims_analytics.fact_billing CASCADE;
DROP TABLE IF EXISTS hims_analytics.fact_appointments CASCADE;
-- Then dimension tables
DROP TABLE IF EXISTS hims_analytics.dim_departments CASCADE;
DROP TABLE IF EXISTS hims_analytics.dim_drugs CASCADE;
DROP TABLE IF EXISTS hims_analytics.dim_doctors CASCADE;
DROP TABLE IF EXISTS hims_analytics.dim_patients CASCADE;

-- Drop schema
DROP SCHEMA IF EXISTS hims_analytics CASCADE;

-- ============================================================================
-- EXTENSIONS
-- Commented out by default as extensions may be used by other databases
-- Uncomment if you need to completely remove extensions
-- ============================================================================

-- DROP EXTENSION IF EXISTS "pg_stat_statements";
-- DROP EXTENSION IF EXISTS "vector";
-- DROP EXTENSION IF EXISTS "uuid-ossp";

-- ============================================================================
-- END OF ROLLBACK MIGRATION
-- ============================================================================
