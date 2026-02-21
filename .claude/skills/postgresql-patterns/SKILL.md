---
name: postgresql-patterns
description: This skill should be used when the user asks to "write PostgreSQL queries", "create database schema", "PostgreSQL optimization", "pgvector queries", "database migrations", "SQL patterns", "PostgreSQL indexes", "full-text search", or mentions PostgreSQL-specific features like JSONB, arrays, CTEs, or window functions.
---

# PostgreSQL Patterns for MediSync

PostgreSQL 18.2 with pgvector is MediSync's data warehouse. This skill covers query patterns, schema design, vector operations, and performance optimization.

★ Insight ─────────────────────────────────────
MediSync's database architecture:
1. **Read-only role** (`medisync_readonly`) for AI agents
2. **pgvector** for semantic search and embeddings
3. **JSONB** for flexible document storage
4. **Partitioning** for time-series data

Never use the superuser role for application queries.
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Convention |
|--------|------------|
| **Schema** | `public` for core, `auth` for user data |
| **Naming** | snake_case for tables/columns |
| **Primary Keys** | UUID v7 for distributed systems |
| **Timestamps** | `timestamptz` always (not `timestamp`) |
| **Soft Delete** | `deleted_at timestamptz` column |

## Schema Design

### Standard Table Pattern

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index for common queries
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
```

### Partitioned Table (Time-Series)

```sql
CREATE TABLE events (
    id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE events_2026_02 PARTITION OF events
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE events_2026_03 PARTITION OF events
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
```

## Query Patterns

### Safe Parameterized Queries

```sql
-- Use $1, $2 placeholders (prevent SQL injection)
SELECT id, name, email
FROM users
WHERE email = $1
  AND deleted_at IS NULL;
```

### Pagination with Cursor

```sql
-- Efficient cursor-based pagination
SELECT id, name, created_at
FROM users
WHERE deleted_at IS NULL
  AND created_at < $1  -- cursor from last page
ORDER BY created_at DESC
LIMIT 20;
```

### CTEs for Complex Queries

```sql
WITH monthly_revenue AS (
    SELECT
        DATE_TRUNC('month', created_at) AS month,
        SUM(amount) AS total
    FROM transactions
    WHERE deleted_at IS NULL
    GROUP BY DATE_TRUNC('month', created_at)
),
ranked_months AS (
    SELECT
        month,
        total,
        LAG(total) OVER (ORDER BY month) AS prev_total
    FROM monthly_revenue
)
SELECT
    month,
    total,
    ROUND((total - prev_total) * 100.0 / NULLIF(prev_total, 0), 2) AS growth_pct
FROM ranked_months
ORDER BY month DESC;
```

### Window Functions

```sql
-- Running totals and rankings
SELECT
    date,
    amount,
    SUM(amount) OVER (ORDER BY date) AS running_total,
    ROW_NUMBER() OVER (PARTITION BY date_trunc('month', date) ORDER BY amount DESC) AS month_rank
FROM transactions
ORDER BY date;
```

## JSONB Operations

```sql
-- Insert JSON data
INSERT INTO users (id, email, name, preferences)
VALUES (gen_random_uuid(), 'user@example.com', 'John', '{"theme": "dark", "notifications": true}');

-- Query JSON fields
SELECT id, name, preferences->>'theme' AS theme
FROM users
WHERE preferences->>'theme' = 'dark';

-- Update nested JSON
UPDATE users
SET preferences = jsonb_set(
    preferences,
    '{notifications}',
    'false'
)
WHERE id = $1;

-- Add to JSON array
UPDATE users
SET preferences = jsonb_set(
    preferences,
    '{tags}',
    COALESCE(preferences->'tags', '[]') || '["new_tag"]'::jsonb
)
WHERE id = $1;

-- Index JSON for fast queries
CREATE INDEX idx_users_prefs_theme ON users((preferences->>'theme'));
```

## pgvector Operations

### Setup

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    embedding VECTOR(1536),  -- OpenAI ada-002 dimensions
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- HNSW index for fast approximate search
CREATE INDEX idx_documents_embedding ON documents
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);
```

### Similarity Search

```sql
-- Cosine similarity search
SELECT id, content, 1 - (embedding <=> $1) AS similarity
FROM documents
ORDER BY embedding <=> $1  -- $1 is query embedding
LIMIT 10;

-- With metadata filter
SELECT id, content, 1 - (embedding <=> $1) AS similarity
FROM documents
WHERE metadata->>'category' = 'medical'
ORDER BY embedding <=> $1
LIMIT 10;
```

### Embedding Generation

```sql
-- Insert with embedding (generated by application)
INSERT INTO documents (content, embedding, metadata)
VALUES ($1, $2::vector, $3::jsonb);

-- Batch insert
INSERT INTO documents (content, embedding)
SELECT content, generate_embedding(content)
FROM raw_documents;
```

## Performance Optimization

### Index Strategies

```sql
-- Composite index for common query patterns
CREATE INDEX idx_transactions_company_date ON transactions(company_id, created_at DESC);

-- Partial index for soft deletes
CREATE INDEX idx_active_users ON users(email) WHERE deleted_at IS NULL;

-- Covering index (include non-indexed columns)
CREATE INDEX idx_orders_covering ON orders(user_id) INCLUDE (total, status);

-- Full-text search index
CREATE INDEX idx_documents_content ON documents USING GIN(to_tsvector('english', content));
```

### Full-Text Search

```sql
-- Basic search
SELECT id, content
FROM documents
WHERE to_tsvector('english', content) @@ plainto_tsquery('healthcare analytics');

-- Ranked search
SELECT id, content,
       ts_rank(to_tsvector('english', content), plainto_tsquery('healthcare')) AS rank
FROM documents
WHERE to_tsvector('english', content) @@ plainto_tsquery('healthcare')
ORDER BY rank DESC
LIMIT 20;

-- Trigram similarity for fuzzy matching
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_users_name_trgm ON users USING GIN(name gin_trgm_ops);

SELECT id, name, similarity(name, 'Jonh') AS sim
FROM users
WHERE name % 'Jonh'  -- similarity threshold
ORDER BY sim DESC;
```

### Query Analysis

```sql
-- Analyze query plan
EXPLAIN ANALYZE
SELECT * FROM transactions WHERE company_id = $1;

-- Check for sequential scans
SELECT relname, seq_scan, idx_scan
FROM pg_stat_user_tables
WHERE seq_scan > idx_scan;

-- Find missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE n_distinct > 100 AND correlation < 0.1;
```

## Migration Pattern

```sql
-- migrations/001_create_users.up.sql
BEGIN;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

COMMIT;

-- migrations/001_create_users.down.sql
BEGIN;
DROP TABLE IF EXISTS users;
COMMIT;
```

## Read-Only Role Setup

```sql
-- Create read-only role for AI agents
CREATE ROLE medisync_readonly WITH LOGIN PASSWORD 'secure_password';

GRANT CONNECT ON DATABASE medisync TO medisync_readonly;
GRANT USAGE ON SCHEMA public TO medisync_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO medisync_readonly;

-- Default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO medisync_readonly;
```

## Additional Resources

### Reference Files
- **`references/vector-operations.md`** - Advanced pgvector patterns
- **`references/performance.md`** - Query optimization techniques

### Example Files
- **`examples/schema.sql`** - Complete schema example
- **`examples/queries.sql`** - Common query patterns
