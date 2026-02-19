# MediSync ETL Operations Policy
# Package: medisync.etl
# Version: 1.0
# Description: OPA policy for ETL service authorization
#
# This policy controls access to ETL operations including:
# - Data extraction from Tally ERP and HIMS
# - Data transformation and loading into the warehouse
# - ETL state management (cursors, checkpoints)
# - Quarantine operations for failed records
# - Quality report generation
#
# Authorized principals:
# - medisync_etl service account (full ETL operations)
# - medisync_app for read-only access to ETL status
# - admin role for emergency operations
#
# Cross-reference: migrations/002_roles.up.sql, docs/agents/03-governance-security.md

package medisync.etl

import rego.v1

# Default deny - explicit allow required for all ETL operations
default allow := false

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

# Check if user has a specific role
user_has_role(role) if {
	role in input.user.roles
}

# Check if request is from the ETL service account
is_etl_service if {
	input.user.service_account == "sa-etl-service"
}

# Check if request is from an admin
is_admin if {
	user_has_role("admin")
}

# =============================================================================
# ETL EXTRACTION OPERATIONS
# Allowed: ETL service account, admin
# =============================================================================

# Allow extraction from Tally ERP
allow if {
	input.action == "extract"
	input.source in {"tally", "hims"}
	is_etl_service
}

# Allow extraction with admin override (emergency)
allow if {
	input.action == "extract"
	input.source in {"tally", "hims"}
	is_admin
	input.emergency_override == true
}

# =============================================================================
# ETL TRANSFORMATION OPERATIONS
# Allowed: ETL service account, admin
# =============================================================================

# Allow data transformation
allow if {
	input.action == "transform"
	is_etl_service
}

# Allow transformation with admin override
allow if {
	input.action == "transform"
	is_admin
	input.emergency_override == true
}

# =============================================================================
# ETL LOAD OPERATIONS
# Allowed: ETL service account, admin
# Target schemas: hims_analytics, tally_analytics
# =============================================================================

# List of schemas the ETL service can write to
etl_writable_schemas := {"hims_analytics", "tally_analytics"}

# Allow loading data into analytics schemas
allow if {
	input.action == "load"
	input.target_schema in etl_writable_schemas
	is_etl_service
}

# Allow load with admin override
allow if {
	input.action == "load"
	input.target_schema in etl_writable_schemas
	is_admin
	input.emergency_override == true
}

# Deny loading into app schema (except specific ETL tables)
deny_load_app_schema if {
	input.action == "load"
	input.target_schema == "app"
	not input.target_table in etl_app_tables
}

# ETL is allowed to write to these specific app schema tables
etl_app_tables := {
	"etl_state",
	"etl_quarantine",
	"etl_quality_report",
	"audit_log",
	"notification_queue",
}

# Allow loading into specific app tables
allow if {
	input.action == "load"
	input.target_schema == "app"
	input.target_table in etl_app_tables
	is_etl_service
}

# =============================================================================
# ETL STATE MANAGEMENT
# Allowed: ETL service account, admin (read), app service (read)
# =============================================================================

# Allow ETL service to manage state (read/write)
allow if {
	input.action in {"read_state", "write_state", "update_cursor"}
	is_etl_service
}

# Allow app service to read ETL state (for status monitoring)
allow if {
	input.action == "read_state"
	input.user.service_account == "sa-app-service"
}

# Allow admin to read ETL state
allow if {
	input.action == "read_state"
	is_admin
}

# =============================================================================
# QUARANTINE OPERATIONS
# Allowed: ETL service (write), admin (read/write), app (read)
# =============================================================================

# Allow ETL service to quarantine failed records
allow if {
	input.action in {"quarantine_write", "quarantine_read"}
	is_etl_service
}

# Allow admin to manage quarantine (reprocess, delete)
allow if {
	input.action in {"quarantine_read", "quarantine_reprocess", "quarantine_delete"}
	is_admin
}

# Allow app service to read quarantine for monitoring
allow if {
	input.action == "quarantine_read"
	input.user.service_account == "sa-app-service"
}

# =============================================================================
# QUALITY REPORT OPERATIONS
# Allowed: ETL service (write), all authenticated users (read)
# =============================================================================

# Allow ETL service to write quality reports
allow if {
	input.action == "quality_report_write"
	is_etl_service
}

# Allow all authenticated users to read quality reports
allow if {
	input.action == "quality_report_read"
	input.user.authenticated == true
}

# =============================================================================
# AUDIT LOG OPERATIONS
# ETL service can only INSERT (append-only per 003_audit_log_security.up.sql)
# =============================================================================

# Allow ETL service to insert audit log entries
allow if {
	input.action == "audit_log_insert"
	is_etl_service
}

# Deny any audit log modification (enforced at DB level too, but defense in depth)
deny_audit_modification if {
	input.action in {"audit_log_update", "audit_log_delete"}
}

# =============================================================================
# TRUNCATE OPERATIONS
# Allowed: admin only with explicit confirmation
# Used for full refresh scenarios
# =============================================================================

# Allow admin to truncate analytics tables with explicit confirmation
allow if {
	input.action == "truncate"
	input.target_schema in etl_writable_schemas
	is_admin
	input.confirmation_token != ""
	valid_truncate_confirmation(input.confirmation_token)
}

# Validate truncate confirmation token (stub - implement actual validation)
valid_truncate_confirmation(token) if {
	# In production: verify token was issued within last 5 minutes
	# and matches the operation being requested
	token != ""
}

# =============================================================================
# SYNC SCHEDULING OPERATIONS
# Allowed: ETL service, admin
# =============================================================================

# Allow ETL service to manage sync schedules
allow if {
	input.action in {"schedule_sync", "pause_sync", "resume_sync"}
	is_etl_service
}

# Allow admin to manage sync schedules
allow if {
	input.action in {"schedule_sync", "pause_sync", "resume_sync", "force_sync"}
	is_admin
}

# =============================================================================
# NATS EVENT PUBLISHING
# ETL service can publish ETL-related events
# =============================================================================

# Allowed NATS subjects for ETL service
etl_nats_subjects := {
	"etl.sync.started",
	"etl.sync.completed",
	"etl.sync.failed",
	"etl.quarantine.added",
	"etl.quality.alert",
}

# Allow ETL service to publish to ETL-related NATS subjects
allow if {
	input.action == "nats_publish"
	input.subject in etl_nats_subjects
	is_etl_service
}

# =============================================================================
# DECISION METADATA
# Provides additional context for allow/deny decisions
# =============================================================================

# Decision reason for logging and debugging
reason := "ETL service authorized for operation" if {
	allow
	is_etl_service
}

reason := "Admin authorized with emergency override" if {
	allow
	is_admin
	input.emergency_override == true
}

reason := "Admin authorized for operation" if {
	allow
	is_admin
	not input.emergency_override
}

reason := "App service authorized for read-only access" if {
	allow
	input.user.service_account == "sa-app-service"
}

reason := "Authenticated user authorized for quality report read" if {
	allow
	input.action == "quality_report_read"
}

reason := "Operation denied: insufficient privileges" if {
	not allow
}

# =============================================================================
# VIOLATION RULES
# Used for policy testing and compliance reporting
# =============================================================================

# Collect all policy violations for reporting
violations contains msg if {
	deny_load_app_schema
	msg := sprintf("ETL attempted to load into unauthorized app table: %v", [input.target_table])
}

violations contains msg if {
	deny_audit_modification
	msg := sprintf("Attempted audit log modification: %v", [input.action])
}

violations contains msg if {
	input.action == "extract"
	not input.source in {"tally", "hims"}
	msg := sprintf("Unknown extraction source: %v", [input.source])
}

# =============================================================================
# POLICY METADATA
# =============================================================================

# Policy version for tracking
policy_version := "1.0.0"

# Policy description
policy_description := "ETL operations authorization policy for MediSync"

# Required input fields for policy evaluation
required_input_fields := {
	"user",
	"action",
}
