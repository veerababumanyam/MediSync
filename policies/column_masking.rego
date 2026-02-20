# MediSync Column Masking Policy
# Package: medisync.column_masking
# Version: 1.0
# Description: OPA policy for column-level access control and data masking
#
# This policy determines which columns should be masked based on user roles:
# - PII columns (patient_name, phone) are masked for non-admin roles
# - Cost columns are masked for non-finance roles
# - Different mask types: partial, full, hash, redacted
#
# Cross-reference: internal/api/opa.go, docs/agents/03-governance-security.md

package medisync.column_masking

import rego.v1

# Default: no masking
default masks := {}

# =============================================================================
# COLUMN CLASSIFICATIONS
# =============================================================================

# PII (Personally Identifiable Information) columns
pii_columns := {
	"patient_name",
	"patient_name_en",
	"patient_name_ar",
	"phone",
	"mobile",
	"email",
	"address",
	"address_en",
	"address_ar",
	"national_id",
	"emergency_contact_name",
	"emergency_contact_phone",
	"date_of_birth",
}

# Sensitive financial columns
cost_columns := {
	"unit_price",
	"purchase_price",
	"selling_price",
	"mrp",
	"subtotal_amount",
	"discount_amount",
	"tax_amount",
	"total_amount",
	"paid_amount",
	"insurance_amount",
	"consultation_fee",
	"opening_balance",
	"closing_balance",
	"amount",
	"base_currency_amount",
	"cost",
	"revenue",
	"profit",
	"margin",
}

# Highly sensitive columns (always masked for non-admins)
highly_sensitive_columns := {
	"national_id",
	"passport_number",
	"insurance_policy_number",
	"bank_account_number",
	"credit_card_number",
}

# =============================================================================
# ROLE DEFINITIONS
# =============================================================================

# Check if user has admin role
is_admin if {
	"admin" in input.user.roles
}

# Check if user has finance role
is_finance if {
	"finance" in input.user.roles
}

# Check if user has medical role (doctors, nurses)
is_medical if {
	some role in {"doctor", "nurse", "medical_staff"}
	role in input.user.roles
}

# Check if user has reporting role
is_reporting if {
	"reporting" in input.user.roles
}

# =============================================================================
# MASK TYPE DETERMINATION
# =============================================================================

# Determine mask type for a column based on user roles and column type
get_mask_type(column) := "none" if {
	is_admin
}

get_mask_type(column) := "full" if {
	not is_admin
	column in highly_sensitive_columns
}

get_mask_type(column) := "partial" if {
	not is_admin
	column in pii_columns
	column not in highly_sensitive_columns
}

get_mask_type(column) := "full" if {
	not is_finance
	column in cost_columns
}

get_mask_type(column) := "none" if {
	true  # Default case
}

# =============================================================================
# MASKING RULES
# =============================================================================

# Build the masks map for all requested columns
masks[mask_column] := mask_type if {
	some column in input.columns
	mask_column := column
	mask_type := get_mask_type(column)
	mask_type != "none"
}

# =============================================================================
# MASKING FUNCTIONS (for reference)
# =============================================================================

# Partial mask: show first 2 and last 2 characters
# Example: "John Smith" -> "Jo******th"
partial_mask_pattern := "show_ends"

# Full mask: replace with asterisks
# Example: "John Smith" -> "**********"
full_mask_pattern := "asterisks"

# Hash mask: one-way hash of the value
# Example: "john@example.com" -> "a1b2c3d4e5f6..."
hash_mask_pattern := "sha256"

# Redacted: replace with [REDACTED]
# Example: "123-45-6789" -> "[REDACTED]"
redacted_pattern := "placeholder"

# =============================================================================
# ADDITIONAL MASKING RULES BY TABLE
# =============================================================================

# Patient table specific rules
patient_table_masks := {
	"patient_name":     "partial",
	"phone":           "partial",
	"email":           "partial",
	"address":         "full",
	"national_id":     "full",
	"date_of_birth":   "partial",
}

# Billing table specific rules
billing_table_masks := {
	"total_amount":    "full",
	"paid_amount":     "full",
	"insurance_amount": "full",
}

# Apply table-specific masks when appropriate
table_specific_mask(column, table) := mask_type if {
	table == "dim_patients"
	mask_type := patient_table_masks[column]
}

table_specific_mask(column, table) := mask_type if {
	table == "fact_billing"
	mask_type := billing_table_masks[column]
}

# =============================================================================
# DECISION METADATA
# =============================================================================

# Reason for masking decisions
reason := "No masking required" if {
	len(masks) == 0
}

reason := sprintf("Masking %d columns based on user roles: %v", [len(masks), input.user.roles]) if {
	len(masks) > 0
}

# Summary of masked columns by type
mask_summary := {
	"partial":  count_partial_masks,
	"full":     count_full_masks,
	"redacted": count_redacted_masks,
}

count_partial_masks if {
	count({column |
		some column, mask_type in masks
		mask_type == "partial"
	})
}

count_full_masks if {
	count({column |
		some column, mask_type in masks
		mask_type == "full"
	})
}

count_redacted_masks if {
	count({column |
		some column, mask_type in masks
		mask_type == "redacted"
	})
}

# =============================================================================
# COMPLIANCE HELPERS
# =============================================================================

# Check if query accesses any PII columns
accesses_pii if {
	some column in input.columns
	column in pii_columns
}

# Check if query accesses any cost columns
accesses_costs if {
	some column in input.columns
	column in cost_columns
}

# Warning if accessing sensitive data without appropriate role
warnings contains msg if {
	accesses_pii
	not is_admin
	not is_medical
	msg := "Accessing PII columns without medical role"
}

warnings contains msg if {
	accesses_costs
	not is_finance
	not is_admin
	msg := "Accessing cost columns without finance role"
}

# =============================================================================
# AUDIT TRAIL
# =============================================================================

# Audit information for logging
audit_info := {
	"user_id":       input.user.id,
	"roles":         input.user.roles,
	"columns":       input.columns,
	"masked_count":  len(masks),
	"accesses_pii":  accesses_pii,
	"accesses_cost": accesses_costs,
}

# =============================================================================
# POLICY METADATA
# =============================================================================

policy_version := "1.0.0"
policy_description := "Column-level masking policy for data access control"

# Required input fields
required_input_fields := {"columns", "user"}
