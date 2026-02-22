# MediSync Council of AIs Consensus Policy
# Package: medisync.council
# Version: 1.0
# Description: OPA policy for Council of AIs consensus system with role-based access control
#
# This policy implements the following access model:
# - Admin users: Full access to all deliberations, audit trail, flag operations
# - Regular users: Access only to their own deliberations
#
# Cross-reference: internal/agents/council/, docs/agents/03-governance-security.md
# Feature: specs/001-council-ai-consensus/spec.md

package medisync.council

import rego.v1

# Default deny all access
default allow := false

# =============================================================================
# ROLE DEFINITIONS
# =============================================================================

# Admin roles that have full access
admin_roles := {"admin", "superadmin", "system"}

# Roles that can create and view deliberations
user_roles := {"user", "analyst", "finance", "operations"}

# All valid roles
valid_roles := admin_roles | user_roles

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

# Check if user has admin role
is_admin(user) if {
	some role in user.roles
	admin_roles[role]
}

# Check if user has a valid role
has_valid_role(user) if {
	some role in user.roles
	valid_roles[role]
}

# Check if user is authenticated
is_authenticated(user) if {
	user.authenticated == true
	user.id != ""
}

# Check if user owns the resource
owns_resource(user, resource) if {
	user.id == resource.user_id
}

# Check if action is a read operation
is_read_action(action) if {
	action in {"view", "list", "get", "search"}
}

# Check if action is an admin-only operation
is_admin_action(action) if {
	action in {"flag", "export_audit", "view_all", "admin_list"}
}

# =============================================================================
# ALLOW RULES - ADMIN ACCESS
# =============================================================================

# Admin users can perform any action on any deliberation
allow if {
	is_authenticated(input.user)
	is_admin(input.user)
	has_valid_role(input.user)
}

# Admin users can access audit trail
allow if {
	is_authenticated(input.user)
	is_admin(input.user)
	input.action == "view_audit"
}

# Admin users can flag deliberations
allow if {
	is_authenticated(input.user)
	is_admin(input.user)
	input.action == "flag"
}

# =============================================================================
# ALLOW RULES - USER ACCESS
# =============================================================================

# Users can create new deliberations
allow if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	input.action == "create"
}

# Users can view their own deliberations
allow if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	is_read_action(input.action)
	owns_resource(input.user, input.resource)
}

# Users can view their own deliberation results
allow if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	input.action == "view_result"
	owns_resource(input.user, input.resource)
}

# Users can view evidence for their own deliberations
allow if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	input.action == "view_evidence"
	owns_resource(input.user, input.resource)
}

# =============================================================================
# DENY RULES (for violation reporting)
# =============================================================================

# User is not authenticated
deny_not_authenticated if {
	not is_authenticated(input.user)
}

# User lacks valid role
deny_invalid_role if {
	is_authenticated(input.user)
	not has_valid_role(input.user)
}

# User trying to access another user's resource
deny_not_owner if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	not is_admin(input.user)
	is_read_action(input.action)
	not owns_resource(input.user, input.resource)
}

# Non-admin trying admin action
deny_admin_required if {
	is_authenticated(input.user)
	has_valid_role(input.user)
	not is_admin(input.user)
	is_admin_action(input.action)
}

# =============================================================================
# DECISION METADATA
# =============================================================================

# Provide reason for allow decisions
reason := "Admin access granted" if {
	allow
	is_admin(input.user)
}

reason := "User accessing own resource" if {
	allow
	not is_admin(input.user)
	owns_resource(input.user, input.resource)
}

reason := "User creating deliberation" if {
	allow
	input.action == "create"
}

# Provide reason for deny decisions
reason := "Access denied: user not authenticated" if {
	deny_not_authenticated
}

reason := "Access denied: invalid role" if {
	deny_invalid_role
}

reason := "Access denied: not resource owner" if {
	deny_not_owner
}

reason := "Access denied: admin role required" if {
	deny_admin_required
}

reason := "Access denied: policy violation" if {
	not allow
}

# =============================================================================
# VIOLATION REPORTING
# =============================================================================

# Collect all violations
violations contains msg if {
	deny_not_authenticated
	msg := "User is not authenticated"
}

violations contains msg if {
	deny_invalid_role
	msg := sprintf("User has no valid role. Roles: %v", [input.user.roles])
}

violations contains msg if {
	deny_not_owner
	msg := sprintf("User %v cannot access resource owned by %v", [input.user.id, input.resource.user_id])
}

violations contains msg if {
	deny_admin_required
	msg := sprintf("Action '%v' requires admin role", [input.action])
}

# =============================================================================
# SCOPE FILTERING
# =============================================================================

# Determine the scope filter for list operations
# Admin users see all, regular users see only their own
list_scope := "all" if {
	is_authenticated(input.user)
	is_admin(input.user)
	input.action == "list"
}

list_scope := "own" if {
	is_authenticated(input.user)
	not is_admin(input.user)
	has_valid_role(input.user)
	input.action == "list"
}

list_scope := "none" if {
	not allow
	input.action == "list"
}

# SQL filter clause for repository queries
list_filter := "1=1" if {
	list_scope == "all"
}

list_filter := sprintf("user_id = '%s'", [input.user.id]) if {
	list_scope == "own"
}

list_filter := "1=0" if {
	list_scope == "none"
}

# =============================================================================
# RESOURCE ACCESS LEVELS
# =============================================================================

# Determine access level for a resource
access_level := "full" if {
	allow
	is_admin(input.user)
}

access_level := "own" if {
	allow
	not is_admin(input.user)
	owns_resource(input.user, input.resource)
}

access_level := "create_only" if {
	allow
	input.action == "create"
	not owns_resource(input.user, input.resource)
}

access_level := "none" if {
	not allow
}

# =============================================================================
# AUDIT REQUIREMENTS
# =============================================================================

# Determine if action should be audited
requires_audit if {
	input.action in {"create", "flag", "view_audit", "export_audit"}
}

# Audit entry details
audit_details := {
	"action": input.action,
	"user_id": input.user.id,
	"resource_id": input.resource.id,
	"access_level": access_level,
} if {
	allow
	requires_audit
}

# =============================================================================
# MASKING RULES
# =============================================================================

# Fields that should be masked for non-admin users viewing other users' data
masked_fields := {"ip_address", "user_email", "session_id"}

# Apply field masking based on role and ownership
should_mask_field(field) if {
	not is_admin(input.user)
	not owns_resource(input.user, input.resource)
	masked_fields[field]
}

# =============================================================================
# RATE LIMITING HINTS
# =============================================================================

# Suggest rate limit tier based on action
rate_limit_tier := "high" if {
	input.action == "create"
}

rate_limit_tier := "medium" if {
	input.action in {"view", "list", "get"}
}

rate_limit_tier := "low" if {
	input.action in {"flag", "view_audit"}
}

# =============================================================================
# POLICY METADATA
# =============================================================================

policy_version := "1.0.0"
policy_description := "Council of AIs consensus system - role-based access control with admin/user separation"

# Required input fields
required_input_fields := {"user", "action"}

# Optional input fields
optional_input_fields := {"resource", "query"}
