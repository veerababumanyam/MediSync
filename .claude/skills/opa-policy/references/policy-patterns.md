# Advanced Rego Policy Patterns

## Partial Evaluation and Performance

### Pre-Compiled Policies

```rego
# Use partial evaluation for performance
package medisync.optimized

# Pre-compute role permissions at compile time
precomputed_permissions[permission] {
    some role in input.roles
    permission := role_permissions[role][_]
}

# Fast allow using precomputed
allow if {
    some permission in precomputed_permissions
    permission.resource == input.resource
    permission.action == input.action
}
```

### Indexing Patterns

```rego
# Use indexing for large data sets
# OPA indexes rules with single-value heads

# Good - indexed by company_id
user_accessible_patients[patient_id] if {
    some patient in data.patients
    patient.company_id == input.user_company_id
    patient_id := patient.id
}

# Avoid - not indexed (iteration required)
user_accessible_patients if {
    some patient in data.patients
    patient.company_id == input.user_company_id
}
```

## Data Filtering

### Row-Level Security

```rego
package medisync.row_security

# Filter patients based on access rules
filter_patients[patient] {
    some patient in data.patients
    patient.company_id == input.user_company_id
    not patient.sensitive or has_sensitive_access
}

has_sensitive_access if {
    "view_sensitive" in input.permissions
}

# Filter appointments
filter_appointments[appt] {
    some appt in data.appointments
    appt.company_id == input.user_company_id

    # Non-admins can only see their own or team's appointments
    input.user_role != "admin"
    appt.provider_id == input.user_id or
    appt.provider_id in input.team_member_ids
}
```

### Field-Level Masking

```rego
package medisync.field_masking

# Define sensitive fields per role
sensitive_fields := {
    "viewer": ["ssn", "salary", "notes"],
    "analyst": ["ssn"],
    "finance_head": [],
    "admin": [],
}

# Mask sensitive fields in response
mask_patient[masked] {
    some patient in data.patients
    patient.id == input.patient_id

    # Start with full patient
    masked := patient

    # Remove sensitive fields based on role
    some role in input.roles
    fields_to_remove := sensitive_fields[role]
    some field in fields_to_remove
    not masked[field]
}
```

## Time-Based Policies

### Time Window Restrictions

```rego
package medisync.time_restrictions

import future.keywords.if

# Business hours only for certain operations
allow_during_business_hours if {
    input.action in ["sync", "approve", "export"]

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

# Allow outside business hours for read-only
allow if {
    input.action == "read"
}

# Allow sync only during business hours
allow if {
    input.action == "sync"
    allow_during_business_hours
}

# Admins can override
allow if {
    "admin" in input.roles
}
```

### Expiration Policies

```rego
package medisync.expiration

# Check if approval has expired
approval_valid if {
    now := time.now_ns()
    approval_time := input.context.approval_timestamp

    # Approval valid for 24 hours
    valid_duration := 24 * 60 * 60 * 1000000000  # nanoseconds

    now - approval_time < valid_duration
}

# Session validity
session_valid if {
    now := time.now_ns()
    issued_at := input.token.iat * 1000000000

    # Session valid for 15 minutes
    session_duration := 15 * 60 * 1000000000

    now - issued_at < session_duration
}
```

## Complex Conditions

### Multi-Factor Authorization

```rego
package medisync.mfa

# High-value operations require multiple factors
require_mfa if {
    input.action in ["sync", "approve", "delete"]
    input.context.transaction_value > 100000
    not "admin" in input.roles
}

# Tally sync requirements
allow_sync if {
    input.action == "sync"

    # Must have approval
    input.context.approved == true

    # Approval must be from different user
    input.context.approver_id != input.user

    # Must have MFA if high value
    not require_mfa or input.context.mfa_verified == true
}

# Deletion requires additional approval
allow_delete if {
    input.action == "delete"

    # Must be admin
    "admin" in input.roles

    # Must have secondary approval
    input.context.secondary_approval == true

    # Cannot delete if synced
    not input.context.synced_to_tally
}
```

### Contextual Access

```rego
package medisync.contextual

# IP-based restrictions
ip_allowed if {
    allowed_ranges := data.network.allowed_ip_ranges
    user_ip := net.cidr_expand(input.context.ip_address)

    some range in allowed_ranges
    net.cidr_contains(range, user_ip[_])
}

# Location-based access
location_allowed if {
    input.action in ["read", "export"]
    not input.context.location_restricted
}

location_allowed if {
    input.action in ["sync", "approve"]
    input.context.location_verified == true
}

# Device trust
device_trusted if {
    input.context.device_registered == true
    input.context.device_last_used > time.now_ns() - 30 * 24 * 60 * 60 * 1000000000
}

# Combined access check
allow if {
    role_allowed
    ip_allowed
    location_allowed
    device_trusted
}
```

## Audit and Compliance

### Audit Policy

```rego
package medisync.audit

# Operations that must be logged
audit_required_operations := {
    "create", "update", "delete", "sync", "approve", "export"
}

# Generate audit record
audit_record := {
    "timestamp": time.now_ns(),
    "user_id": input.user,
    "roles": input.roles,
    "action": input.action,
    "resource": input.resource,
    "resource_id": input.context.resource_id,
    "company_id": input.user_company_id,
    "ip_address": input.context.ip_address,
    "user_agent": input.context.user_agent,
    "decision": "allow",
    "reason": decision_reason,
}

# Always audit high-value operations
decision_reason := "role_authorized" if {
    role_grants_access
}

decision_reason := "admin_override" if {
    "admin" in input.roles
}

decision_reason := "approval_based" if {
    input.context.approved == true
}
```

### Compliance Checks

```rego
package medisync.compliance

# GDPR compliance
gdpr_compliant if {
    # Must have consent for data access
    input.context.consent_given == true

    # Must be within consent scope
    input.action in input.context.consent_scope

    # Data retention policy met
    data_retention_met
}

data_retention_met if {
    some data in data.subjects
    data.id == input.context.subject_id
    retention_days := 365 * 7  # 7 years for healthcare

    now := time.now_ns()
    created := data.created_at * 1000000000

    (now - created) / (24 * 60 * 60 * 1000000000) < retention_days
}

# HIPAA compliance
hipaa_compliant if {
    # PHI access requires specific roles
    if input.context.contains_phi then {
        "hipaa_viewer" in input.roles or
        "admin" in input.roles
    }

    # Must have business associate agreement
    input.context.baa_signed == true

    # Audit trail required
    audit_record.decision == "allow"
}
```

## Policy Composition

### Inheritance Pattern

```rego
package medisync.base

# Base allow rule
default allow = false

# Common conditions for all operations
base_conditions if {
    input.user != ""
    input.user_company_id == input.context.company_id
}

# Base read permission
allow_read if {
    base_conditions
    input.action == "read"
    has_any_role(["viewer", "analyst", "finance_head", "admin"])
}

# Base write permission
allow_write if {
    base_conditions
    input.action == "write"
    has_any_role(["analyst", "finance_head", "admin"])
}
```

```rego
package medisync.documents

import data.medisync.base

# Extend base policies
allow if {
    base.allow_read
    document_accessible
}

allow if {
    base.allow_write
    document_editable
}

# Document-specific rules
document_accessible if {
    input.context.document.company_id == input.user_company_id

    # Sensitive documents require elevated access
    if input.context.document.sensitive then {
        has_any_role(["finance_head", "admin"])
    }
}
```

## Testing Patterns

### Table-Driven Tests

```rego
package medisync.dashboard_test

import data.medisync.dashboard

# Test cases as data
test_cases := [
    {
        "name": "viewer_can_read",
        "input": {"user": "u1", "roles": ["viewer"], "action": "read"},
        "expected": true,
    },
    {
        "name": "viewer_cannot_write",
        "input": {"user": "u1", "roles": ["viewer"], "action": "write"},
        "expected": false,
    },
    {
        "name": "analyst_can_write",
        "input": {"user": "u2", "roles": ["analyst"], "action": "write"},
        "expected": true,
    },
]

# Run all tests
test_all if {
    some tc in test_cases
    result := allow with input as tc.input
    result == tc.expected
}
```

### Mock Data Tests

```rego
package medisync.tally_test

import data.medisync.tally

# Mock data for tests
mock_patients := [
    {"id": "p1", "company_id": "c1", "name": "Patient 1"},
    {"id": "p2", "company_id": "c2", "name": "Patient 2"},
]

# Test with mock data
test_company_isolation if {
    # User from c1 should only see p1
    result := filter_patients with input as {
        "user": "u1",
        "roles": ["viewer"],
        "user_company_id": "c1",
    } with data.patients as mock_patients

    count(result) == 1
    result[0].id == "p1"
}
```
