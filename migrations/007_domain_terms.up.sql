-- MediSync AI Agent Core - Domain Terms Migration
-- Version: 007
-- Description: Create domain terms table for terminology mapping in AI Agent Core
-- Task: T010
--
-- This migration establishes:
-- 1. app.domain_terms - Mapping between user vocabulary and canonical database terminology
--
-- Purpose:
-- - Maps synonyms to canonical terms (e.g., "footfall" -> "patient_visits")
-- - Stores SQL fragments for quick reference
-- - Supports locale-specific variants for bilingual (EN/AR) queries
--
-- Example:
--   synonym: "footfall"
--   canonical_term: "patient_visits"
--   category: "healthcare"
--   sql_fragment: "fact_appointments.count"
--   locale_variants: {"en": ["walk-ins", "visits", "footfall"], "ar": ["زيارات", "حضور"]}

-- ============================================================================
-- TABLE: app.domain_terms
-- Purpose: Mapping between user vocabulary and canonical database terminology
-- ============================================================================

CREATE TABLE IF NOT EXISTS app.domain_terms (
    id SERIAL PRIMARY KEY,
    synonym VARCHAR(255) NOT NULL,
    canonical_term VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    sql_fragment TEXT,
    locale_variants JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT uq_domain_terms_synonym UNIQUE (synonym),
    CONSTRAINT ck_domain_terms_category CHECK (category IN ('healthcare', 'accounting', 'general', 'finance', 'operations', 'pharmacy'))
);

-- Indexes for domain_terms
CREATE INDEX IF NOT EXISTS idx_domain_terms_synonym ON app.domain_terms(synonym);
CREATE INDEX IF NOT EXISTS idx_domain_terms_category ON app.domain_terms(category);
CREATE INDEX IF NOT EXISTS idx_domain_terms_canonical ON app.domain_terms(canonical_term);

COMMENT ON TABLE app.domain_terms IS 'Maps user-facing synonyms to canonical database terminology for AI query understanding';

COMMENT ON COLUMN app.domain_terms.synonym IS 'User-facing term that needs to be mapped (e.g., "footfall")';
COMMENT ON COLUMN app.domain_terms.canonical_term IS 'Canonical database term (e.g., "patient_visits")';
COMMENT ON COLUMN app.domain_terms.category IS 'Domain category: healthcare, accounting, general, finance, operations, pharmacy';
COMMENT ON COLUMN app.domain_terms.sql_fragment IS 'SQL mapping hint or fragment (e.g., "fact_appointments.count")';
COMMENT ON COLUMN app.domain_terms.locale_variants IS 'Locale-specific variants: {"en": ["walk-ins"], "ar": ["زيارات"]}';

-- ============================================================================
-- FUNCTION: Find canonical term
-- Purpose: Look up canonical term by synonym
-- ============================================================================

CREATE OR REPLACE FUNCTION app.find_canonical_term(p_synonym VARCHAR(255))
RETURNS TABLE (
    canonical_term VARCHAR(255),
    category VARCHAR(50),
    sql_fragment TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        dt.canonical_term,
        dt.category,
        dt.sql_fragment
    FROM app.domain_terms dt
    WHERE LOWER(dt.synonym) = LOWER(p_synonym);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.find_canonical_term(VARCHAR) IS 'Look up canonical term, category, and SQL fragment by synonym';

-- ============================================================================
-- FUNCTION: Search domain terms
-- Purpose: Search for domain terms by pattern
-- ============================================================================

CREATE OR REPLACE FUNCTION app.search_domain_terms(
    p_search_term VARCHAR(255),
    p_category VARCHAR(50) DEFAULT NULL,
    p_limit INTEGER DEFAULT 10
)
RETURNS TABLE (
    id INTEGER,
    synonym VARCHAR(255),
    canonical_term VARCHAR(255),
    category VARCHAR(50),
    sql_fragment TEXT,
    locale_variants JSONB,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        dt.id,
        dt.synonym,
        dt.canonical_term,
        dt.category,
        dt.sql_fragment,
        dt.locale_variants,
        similarity(LOWER(dt.synonym), LOWER(p_search_term)) AS sim
    FROM app.domain_terms dt
    WHERE
        (p_category IS NULL OR dt.category = p_category)
        AND (
            LOWER(dt.synonym) LIKE '%' || LOWER(p_search_term) || '%'
            OR LOWER(dt.canonical_term) LIKE '%' || LOWER(p_search_term) || '%'
        )
    ORDER BY sim DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.search_domain_terms(VARCHAR, VARCHAR, INTEGER) IS 'Search domain terms by pattern with optional category filter. Returns matches ordered by similarity.';

-- ============================================================================
-- FUNCTION: Get all terms for locale
-- Purpose: Get all domain terms with locale-specific variants
-- ============================================================================

CREATE OR REPLACE FUNCTION app.get_domain_terms_for_locale(
    p_locale VARCHAR(2) DEFAULT 'en',
    p_category VARCHAR(50) DEFAULT NULL
)
RETURNS TABLE (
    synonym VARCHAR(255),
    canonical_term VARCHAR(255),
    category VARCHAR(50),
    sql_fragment TEXT,
    locale_variants JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        dt.synonym,
        dt.canonical_term,
        dt.category,
        dt.sql_fragment,
        dt.locale_variants
    FROM app.domain_terms dt
    WHERE
        (p_category IS NULL OR dt.category = p_category)
        AND (
            -- Include if has variant for requested locale
            dt.locale_variants ? p_locale
            -- Or if no locale variants specified (available in all locales)
            OR dt.locale_variants = '{}'::jsonb
        );
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.get_domain_terms_for_locale(VARCHAR, VARCHAR) IS 'Get domain terms filtered by locale availability and optional category.';

-- ============================================================================
-- FUNCTION: Upsert domain term
-- Purpose: Insert or update a domain term mapping
-- ============================================================================

CREATE OR REPLACE FUNCTION app.upsert_domain_term(
    p_synonym VARCHAR(255),
    p_canonical_term VARCHAR(255),
    p_category VARCHAR(50),
    p_sql_fragment TEXT DEFAULT NULL,
    p_locale_variants JSONB DEFAULT '{}'
)
RETURNS INTEGER AS $$
DECLARE
    v_id INTEGER;
BEGIN
    -- Check if term exists
    SELECT id INTO v_id
    FROM app.domain_terms
    WHERE synonym = p_synonym;

    IF v_id IS NOT NULL THEN
        -- Update existing record
        UPDATE app.domain_terms
        SET
            canonical_term = p_canonical_term,
            category = p_category,
            sql_fragment = p_sql_fragment,
            locale_variants = p_locale_variants
        WHERE id = v_id;
        RETURN v_id;
    ELSE
        -- Insert new record
        INSERT INTO app.domain_terms (
            synonym, canonical_term, category, sql_fragment, locale_variants, created_at
        ) VALUES (
            p_synonym, p_canonical_term, p_category, p_sql_fragment, p_locale_variants, NOW()
        )
        RETURNING id INTO v_id;
        RETURN v_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION app.upsert_domain_term(VARCHAR, VARCHAR, VARCHAR, TEXT, JSONB) IS 'Insert or update a domain term mapping. Returns the ID of the upserted record.';

-- ============================================================================
-- SEED DATA: Common healthcare and accounting domain terms
-- ============================================================================

INSERT INTO app.domain_terms (synonym, canonical_term, category, sql_fragment, locale_variants, created_at) VALUES
-- Healthcare terms
('footfall', 'patient_visits', 'healthcare', 'hims_analytics.fact_appointments', '{"en": ["walk-ins", "visits", "footfall", "patient traffic"], "ar": ["زيارات", "حضور"]}', NOW()),
('outpatient', 'opd_visits', 'healthcare', 'hims_analytics.fact_appointments WHERE appt_type = ''consultation''', '{"en": ["opd", "outpatient visits", "clinic visits"], "ar": ["عيادات خارجية", "زيارات خارجية"]}', NOW()),
('revenue', 'total_revenue', 'healthcare', 'hims_analytics.fact_billing.total_amount', '{"en": ["income", "earnings", "revenue", "sales"], "ar": ["إيرادات", "دخل"]}', NOW()),
('collections', 'payment_received', 'healthcare', 'hims_analytics.fact_billing.paid_amount', '{"en": ["collections", "payments received", "cash collected"], "ar": ["تحصيلات", "مدفوعات"]}', NOW()),
('pending', 'outstanding_amount', 'healthcare', 'hims_analytics.fact_billing.outstanding_amount', '{"en": ["pending", "receivables", "outstanding", "dues"], "ar": ["معلق", "مستحقات"]}', NOW()),

-- Accounting terms
('sales', 'sales_vouchers', 'accounting', 'tally_analytics.fact_vouchers WHERE voucher_type = ''sales''', '{"en": ["sales", "revenue from sales", "turnover"], "ar": ["مبيعات"]}', NOW()),
('purchases', 'purchase_vouchers', 'accounting', 'tally_analytics.fact_vouchers WHERE voucher_type = ''purchase''', '{"en": ["purchases", "procurement", "buying"], "ar": ["مشتريات"]}', NOW()),
('receivables', 'accounts_receivable', 'accounting', 'tally_analytics.dim_ledgers WHERE ledger_type = ''asset'' AND ledger_name LIKE ''%receivable%''', '{"en": ["receivables", "debtors", "money owed"], "ar": ["ذمم مدينة", "مستحقات"]}', NOW()),
('payables', 'accounts_payable', 'accounting', 'tally_analytics.dim_ledgers WHERE ledger_type = ''liability'' AND ledger_name LIKE ''%payable%''', '{"en": ["payables", "creditors", "money owed to"], "ar": ["ذمم دائنة"]}', NOW()),
('profit', 'net_profit', 'accounting', 'tally_analytics.dim_ledgers WHERE ledger_type = ''income''', '{"en": ["profit", "earnings", "net income"], "ar": ["ربح", "صافي الربح"]}', NOW()),
('expense', 'expenses', 'accounting', 'tally_analytics.fact_vouchers WHERE voucher_type IN (''payment'', ''journal'')', '{"en": ["expense", "expenditure", "cost"], "ar": ["مصروفات", "نفقات"]}', NOW()),

-- Finance terms
('cash', 'cash_balance', 'finance', 'tally_analytics.dim_ledgers WHERE is_bank_account = FALSE AND ledger_name ILIKE ''%cash%''', '{"en": ["cash", "cash in hand", "petty cash"], "ar": ["نقدية", "كاش"]}', NOW()),
('bank', 'bank_balance', 'finance', 'tally_analytics.dim_ledgers WHERE is_bank_account = TRUE', '{"en": ["bank", "bank balance", "bank account"], "ar": ["بنك", "رصيد بنكي"]}', NOW()),

-- Pharmacy terms
('stock', 'inventory_stock', 'pharmacy', 'tally_analytics.fact_stock_movements', '{"en": ["stock", "inventory", "medicine stock"], "ar": ["مخزون", "مستودع"]}', NOW()),
('dispensed', 'pharmacy_dispensations', 'pharmacy', 'hims_analytics.fact_pharmacy_dispensations', '{"en": ["dispensed", "medicines given", "prescriptions filled"], "ar": ["صرف", "أدوية مصروفة"]}', NOW())
ON CONFLICT (synonym) DO NOTHING;

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

-- Grant SELECT to medisync_readonly role for AI agents
GRANT SELECT ON app.domain_terms TO medisync_readonly;

-- Grant full CRUD to medisync_app role for application operations
GRANT SELECT, INSERT, UPDATE, DELETE ON app.domain_terms TO medisync_app;
GRANT USAGE, SELECT ON SEQUENCE app.domain_terms_id_seq TO medisync_app;

-- Grant execute on functions
GRANT EXECUTE ON FUNCTION app.find_canonical_term(VARCHAR) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.find_canonical_term(VARCHAR) TO medisync_readonly;
GRANT EXECUTE ON FUNCTION app.search_domain_terms(VARCHAR, VARCHAR, INTEGER) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.search_domain_terms(VARCHAR, VARCHAR, INTEGER) TO medisync_readonly;
GRANT EXECUTE ON FUNCTION app.get_domain_terms_for_locale(VARCHAR, VARCHAR) TO medisync_app;
GRANT EXECUTE ON FUNCTION app.get_domain_terms_for_locale(VARCHAR, VARCHAR) TO medisync_readonly;
GRANT EXECUTE ON FUNCTION app.upsert_domain_term(VARCHAR, VARCHAR, VARCHAR, TEXT, JSONB) TO medisync_app;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
