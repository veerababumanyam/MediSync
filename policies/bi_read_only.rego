# MediSync BI Read-Only Policy
# Package: medisync.bi_read_only
# Version: 1.0
# Description: OPA policy that ensures only SELECT statements are allowed for BI queries
#
# This policy blocks all DML operations (INSERT, UPDATE, DELETE, DROP, CREATE, ALTER, TRUNCATE)
# and only allows SELECT statements for the AI agents.
#
# Cross-reference: internal/agents/module_a/a01_text_to_sql/, docs/agents/03-governance-security.md

package medisync.bi_read_only

import rego.v1

# Default deny - only explicit SELECT statements are allowed
default allow := false

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

# Normalize query by removing leading/trailing whitespace and converting to uppercase
normalize_query(query) := trimmed if {
	lower_query := lower(query)
	trimmed := trim(lower_query)
}

# Check if query starts with SELECT
is_select_query(query) if {
	normalized := normalize_query(query)
	startswith(normalized, "select")
}

# Check if query contains forbidden keywords
contains_dml_keywords(query) if {
	normalized := normalize_query(query)
	keywords := {"insert", "update", "delete", "drop", "create", "alter", "truncate", "grant", "revoke"}

	some keyword in keywords
	contains(normalized, keyword)
}

# Check if query contains common SQL injection patterns
contains_suspicious_patterns(query) if {
	normalized := normalize_query(query)
	patterns := {
		";--",           # Comment injection
		"/*",            # Block comment start
		"*/",            # Block comment end
		"xp_",           # Extended stored procedure
		"sp_",           # System stored procedure
		"exec(",         # Execute command
		"execute(",      # Execute command
		"union select",  # Union-based injection
		"or 1=1",        # Always true injection
		"or '1'='1",     # String injection
		"waitfor delay", # Time-based injection
		"benchmark(",    # MySQL time-based
		"sleep(",        # MySQL sleep
		"pg_sleep(",     # PostgreSQL sleep
	}

	some pattern in patterns
	contains(normalized, pattern)
}

# Check if query is a valid SELECT with proper structure
is_valid_select(query) if {
	is_select_query(query)
	not contains_dml_keywords(query)
	not contains_suspicious_patterns(query)
}

# =============================================================================
# ALLOW RULES
# =============================================================================

# Allow valid SELECT statements from authenticated users
allow if {
	input.user.authenticated == true
	is_valid_select(input.query)
}

# =============================================================================
# DENY RULES (for violation reporting)
# =============================================================================

# DML operation detected
deny_dml if {
	contains_dml_keywords(input.query)
}

# Suspicious pattern detected
deny_injection if {
	contains_suspicious_patterns(input.query)
}

# Not a SELECT query
deny_not_select if {
	not is_select_query(input.query)
}

# =============================================================================
# DECISION METADATA
# =============================================================================

# Provide reason for the decision
reason := "Valid SELECT query from authenticated user" if {
	allow
	input.user.authenticated == true
	is_valid_select(input.query)
}

reason := "Query blocked: DML operation not allowed" if {
	deny_dml
}

reason := "Query blocked: suspicious pattern detected" if {
	deny_injection
}

reason := "Query blocked: not a SELECT statement" if {
	deny_not_select
}

reason := "Query blocked: user not authenticated" if {
	not input.user.authenticated
}

reason := "Query blocked: failed validation" if {
	not allow
}

# =============================================================================
# VIOLATION REPORTING
# =============================================================================

# Collect all violations for the query
violations contains msg if {
	deny_dml
	msg := sprintf("DML keyword detected in query: %v", [extract_dml_keywords(input.query)])
}

violations contains msg if {
	deny_injection
	msg := "Suspicious SQL injection pattern detected"
}

violations contains msg if {
	deny_not_select
	msg := "Query is not a SELECT statement"
}

violations contains msg if {
	not input.user.authenticated
	msg := "User is not authenticated"
}

# Extract which DML keywords were found
extract_dml_keywords(query) := keywords if {
	normalized := normalize_query(query)
	keywords := {keyword |
		keyword in {"insert", "update", "delete", "drop", "create", "alter", "truncate", "grant", "revoke"}
		contains(normalized, keyword)
	}
}

# =============================================================================
# QUERY ANALYSIS
# =============================================================================

# Query type classification
query_type := "select" if {
	is_select_query(input.query)
	not contains_dml_keywords(input.query)
}

query_type := "dml" if {
	contains_dml_keywords(input.query)
}

query_type := "unknown" if {
	not is_select_query(input.query)
	not contains_dml_keywords(input.query)
}

# Security level assessment
security_level := "safe" if {
	is_valid_select(input.query)
	input.user.authenticated
}

security_level := "dangerous" if {
	contains_dml_keywords(input.query)
}

security_level := "suspicious" if {
	contains_suspicious_patterns(input.query)
}

security_level := "blocked" if {
	not allow
}

# =============================================================================
# POLICY METADATA
# =============================================================================

policy_version := "1.0.0"
policy_description := "BI read-only policy - enforces SELECT-only queries for AI agents"

# Required input fields
required_input_fields := {"user", "query"}
