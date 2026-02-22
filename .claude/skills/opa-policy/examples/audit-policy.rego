# Comprehensive Audit Logging Policy for MediSync
# This policy demonstrates audit logging, compliance, and security patterns

package medisync.audit

import future.keywords.if
import future.keywords.in

# =============================================================================
# CONFIGURATION
# =============================================================================

# Operations that must be logged
audit_required_operations := {
    "create",
    "update",
    "delete",
    "sync",
    "approve",
    "reject",
    "export",
    "import",
    "login",
    "logout",
}

# Sensitive resources that require enhanced logging
sensitive_resources := {
    "patients",
    "journal_entries",
    "documents",
    "users",
    "companies",
}

# High-value thresholds for additional logging
high_value_thresholds := {
    "sync": 100000,      # Amount in currency
    "approve": 50000,
    "export": 1000,      # Number of records
}

# =============================================================================
# MAIN AUDIT DECISION
# =============================================================================

# Determine if operation should be logged
should_audit if {
    input.action in audit_required_operations
}

# Determine if enhanced audit is required
requires_enhanced_audit if {
    should_audit
    input.resource in sensitive_resources
}

requires_enhanced_audit if {
    should_audit
    is_high_value_operation
}

# =============================================================================
# AUDIT RECORD GENERATION
# =============================================================================

# Generate complete audit record
audit_record := {
    "id": sprintf("%d-%s", [time.now_ns(), input.user]),
    "timestamp": time.now_ns(),
    "user": {
        "id": input.user,
        "email": input.user_email,
        "roles": input.roles,
        "company_id": input.user_company_id,
    },
    "action": {
        "type": input.action,
        "resource": input.resource,
        "resource_id": input.context.resource_id,
        "details": input.context.action_details,
    },
    "context": {
        "ip_address": input.context.ip_address,
        "user_agent": input.context.user_agent,
        "session_id": input.context.session_id,
        "request_id": input.context.request_id,
    },
    "result": {
        "decision": input.decision,
        "reason": decision_reason,
        "error": input.context.error_message,
    },
    "compliance": compliance_metadata,
    "enhanced": requires_enhanced_audit,
}

# =============================================================================
# DECISION REASONS
# =============================================================================

decision_reason := "role_authorized" if {
    role_based_access
}

decision_reason := "permission_granted" if {
    permission_based_access
}

decision_reason := "approval_based" if {
    approval_based_access
}

decision_reason := "admin_override" if {
    admin_access
}

decision_reason := "denied_unauthorized" if {
    not authorized
}

decision_reason := "denied_company_mismatch" if {
    authorized
    company_mismatch
}

decision_reason := "denied_approval_required" if {
    authorized
    requires_approval
    not has_approval
}

decision_reason := "denied_time_restricted" if {
    authorized
    time_restricted
    not within_allowed_time
}

# Default reason
decision_reason := "unknown" if {
    true
}

# =============================================================================
# COMPLIANCE METADATA
# =============================================================================

compliance_metadata := {
    "hipaa_applicable": hipaa_applicable,
    "gdpr_applicable": gdpr_applicable,
    "contains_phi": contains_phi,
    "contains_pii": contains_pii,
    "data_classification": data_classification,
    "retention_period_days": retention_period,
}

# HIPAA applies to patient data
hipaa_applicable if {
    input.resource == "patients"
}

hipaa_applicable if {
    input.resource == "documents"
    input.context.contains_phi == true
}

# GDPR applies to EU data subjects
gdpr_applicable if {
    input.context.subject_location == "EU"
}

# Check for PHI (Protected Health Information)
contains_phi if {
    input.resource == "patients"
    input.action in ["read", "export"]
}

contains_phi if {
    input.resource == "documents"
    input.context.document_type in ["medical_record", "prescription", "lab_result"]
}

# Check for PII (Personally Identifiable Information)
contains_pii if {
    input.resource in ["patients", "users"]
}

# Data classification
data_classification := "highly_sensitive" if {
    contains_phi
}

data_classification := "sensitive" if {
    contains_pii
    not contains_phi
}

data_classification := "internal" if {
    input.resource in ["journal_entries", "transactions"]
}

data_classification := "public" if {
    input.resource == "reports"
    input.context.report_type == "public"
}

# Retention period based on classification
retention_period := 2555 if {  # 7 years for healthcare
    contains_phi
}

retention_period := 1825 if {  # 5 years for financial
    input.resource in ["journal_entries", "transactions"]
}

retention_period := 365 if {   # 1 year for general
    true
}

# =============================================================================
# ACCESS HELPERS
# =============================================================================

authorized if {
    input.user != ""
    "user" in input.roles
}

role_based_access if {
    authorized
    some role in input.roles
    role_allows_action(role, input.resource, input.action)
}

permission_based_access if {
    authorized
    input.action in input.permissions
}

approval_based_access if {
    authorized
    requires_approval
    has_approval
}

admin_access if {
    authorized
    "admin" in input.roles
}

requires_approval if {
    input.action == "sync"
    input.context.transaction_value > high_value_thresholds["sync"]
}

requires_approval if {
    input.action == "approve"
    input.context.approval_level > 1
}

has_approval if {
    input.context.approved == true
    input.context.approver_id != input.user
}

company_mismatch if {
    input.context.company_id != input.user_company_id
}

# Time restrictions for certain operations
time_restricted if {
    input.action in ["sync", "approve", "export"]
    not "admin" in input.roles
}

within_allowed_time if {
    now := time.now_ns()
    hour := time.clock(now)[0]
    weekday := time.weekday(now)

    # 9 AM to 6 PM
    hour >= 9
    hour < 18

    # Monday to Friday
    weekday != "Saturday"
    weekday != "Sunday"
}

# =============================================================================
# ROLE PERMISSIONS
# =============================================================================

role_permissions := {
    "admin": {
        "patients": ["read", "write", "delete"],
        "journal_entries": ["read", "write", "delete", "approve", "sync"],
        "documents": ["read", "write", "delete", "approve"],
        "reports": ["read", "export"],
        "users": ["read", "write", "delete"],
    },
    "finance_head": {
        "patients": ["read"],
        "journal_entries": ["read", "write", "approve", "sync"],
        "documents": ["read", "write", "approve"],
        "reports": ["read", "export"],
    },
    "analyst": {
        "patients": ["read"],
        "journal_entries": ["read", "write"],
        "documents": ["read", "write"],
        "reports": ["read", "export"],
    },
    "viewer": {
        "patients": ["read"],
        "journal_entries": ["read"],
        "documents": ["read"],
        "reports": ["read"],
    },
}

role_allows_action(role, resource, action) if {
    resource_perms := role_permissions[role][resource]
    action in resource_perms
}

# =============================================================================
# HIGH VALUE DETECTION
# =============================================================================

is_high_value_operation if {
    threshold := high_value_thresholds[input.action]
    input.context.transaction_value > threshold
}

is_high_value_operation if {
    input.action == "export"
    input.context.record_count > high_value_thresholds["export"]
}

# =============================================================================
# NOTIFICATION RULES
# =============================================================================

# Operations that should trigger notifications
should_notify if {
    requires_enhanced_audit
    input.decision == "allow"
}

# Who should be notified
notification_recipients := recipients if {
    should_notify

    recipients := []

    # Notify compliance for PHI access
    if contains_phi then {
        recipients := array.concat(recipients, ["compliance@medisync.io"])
    }

    # Notify finance for high-value syncs
    if input.action == "sync" and is_high_value_operation then {
        recipients := array.concat(recipients, ["finance@medisync.io"])
    }

    # Notify security for suspicious patterns
    if input.context.suspicious_activity == true then {
        recipients := array.concat(recipients, ["security@medisync.io"])
    }
}
