-- Council of AIs Consensus System: Audit Trail Partitioned Table
-- Migration: 004_audit_partition.up.sql
-- Purpose: Create partitioned audit table for 7-year HIPAA compliance retention

-- Audit Action Types
CREATE TYPE audit_action_type AS ENUM (
    'query',
    'review',
    'flag',
    'export',
    'access'
);

-- Audit Entries Table (Partitioned by Month)
CREATE TABLE audit_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deliberation_id UUID NOT NULL REFERENCES council_deliberations(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL,
    action audit_action_type NOT NULL,
    details JSONB NOT NULL DEFAULT '{}',
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    partition_date DATE NOT NULL DEFAULT CURRENT_DATE
) PARTITION BY RANGE (partition_date);

-- Create partitions for the next 12 months (template for ongoing maintenance)
-- Note: In production, use pg_partman for automatic partition management

-- 2026 partitions
CREATE TABLE audit_entries_2026_02 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE audit_entries_2026_03 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE TABLE audit_entries_2026_04 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

CREATE TABLE audit_entries_2026_05 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');

CREATE TABLE audit_entries_2026_06 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

CREATE TABLE audit_entries_2026_07 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');

CREATE TABLE audit_entries_2026_08 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');

CREATE TABLE audit_entries_2026_09 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE audit_entries_2026_10 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE audit_entries_2026_11 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

CREATE TABLE audit_entries_2026_12 PARTITION OF audit_entries
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

CREATE TABLE audit_entries_2027_01 PARTITION OF audit_entries
    FOR VALUES FROM ('2027-01-01') TO ('2027-02-01');

-- Default partition for safety (catches any dates outside defined ranges)
CREATE TABLE audit_entries_default PARTITION OF audit_entries DEFAULT;

-- Indexes for Audit Table
CREATE INDEX idx_audit_deliberation ON audit_entries(deliberation_id);
CREATE INDEX idx_audit_user ON audit_entries(user_id);
CREATE INDEX idx_audit_action ON audit_entries(action);
CREATE INDEX idx_audit_created_at ON audit_entries(created_at);
CREATE INDEX idx_audit_partition_date ON audit_entries(partition_date);

-- Function for automatic partition maintenance
CREATE OR REPLACE FUNCTION maintain_audit_partitions()
RETURNS void AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    -- Create partition for 2 months ahead if it doesn't exist
    partition_date := DATE_TRUNC('month', CURRENT_DATE + INTERVAL '2 months');
    partition_name := 'audit_entries_' || TO_CHAR(partition_date, 'YYYY_MM');
    start_date := TO_CHAR(partition_date, 'YYYY-MM-DD');
    end_date := TO_CHAR(partition_date + INTERVAL '1 month', 'YYYY-MM-DD');

    -- Check if partition exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_class WHERE relname = partition_name
    ) THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF audit_entries FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE audit_entries IS 'Immutable audit log for HIPAA compliance (7-year retention)';
COMMENT ON COLUMN audit_entries.partition_date IS 'Partition key for monthly partitioning';
COMMENT ON FUNCTION maintain_audit_partitions() IS 'Creates future partitions for audit_entries table';

-- Grant read-only access to medisync_readonly role
GRANT SELECT ON ALL TABLES IN SCHEMA public TO medisync_readonly;
