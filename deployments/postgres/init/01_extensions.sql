-- MediSync PostgreSQL Initialization Script
-- This script runs automatically on first container startup
-- Creates required extensions and initial database configuration

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgvector";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Verify extensions are installed
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector') THEN
        RAISE EXCEPTION 'pgvector extension is not installed';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN
        RAISE EXCEPTION 'uuid-ossp extension is not installed';
    END IF;
    RAISE NOTICE 'All required extensions are installed successfully';
END;
$$;

-- Grant schema creation privileges to medisync user
-- (The user is created automatically from POSTGRES_USER env var)
ALTER DATABASE medisync SET timezone TO 'Asia/Riyadh';
