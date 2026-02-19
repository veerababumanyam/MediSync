# MediSync ETL Operations Policy Tests
# Package: medisync.etl_test
# Description: Unit tests for ETL authorization policy
#
# Run tests with: opa test policies/ -v

package medisync.etl_test

import rego.v1

import data.medisync.etl

# =============================================================================
# TEST: ETL SERVICE EXTRACTION
# =============================================================================

test_etl_service_can_extract_from_tally if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "extract",
		"source": "tally",
	}
}

test_etl_service_can_extract_from_hims if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "extract",
		"source": "hims",
	}
}

test_etl_service_cannot_extract_from_unknown_source if {
	not etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "extract",
		"source": "unknown",
	}
}

# =============================================================================
# TEST: ETL SERVICE TRANSFORMATION
# =============================================================================

test_etl_service_can_transform if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "transform",
	}
}

# =============================================================================
# TEST: ETL SERVICE LOAD OPERATIONS
# =============================================================================

test_etl_service_can_load_to_hims_analytics if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "load",
		"target_schema": "hims_analytics",
	}
}

test_etl_service_can_load_to_tally_analytics if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "load",
		"target_schema": "tally_analytics",
	}
}

test_etl_service_can_load_to_etl_state if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "load",
		"target_schema": "app",
		"target_table": "etl_state",
	}
}

test_etl_service_can_load_to_etl_quarantine if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "load",
		"target_schema": "app",
		"target_table": "etl_quarantine",
	}
}

test_etl_service_cannot_load_to_unauthorized_app_table if {
	not etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "load",
		"target_schema": "app",
		"target_table": "users",
	}
}

# =============================================================================
# TEST: ADMIN OPERATIONS
# =============================================================================

test_admin_can_extract_with_emergency_override if {
	etl.allow with input as {
		"user": {"roles": ["admin"]},
		"action": "extract",
		"source": "tally",
		"emergency_override": true,
	}
}

test_admin_can_truncate_with_confirmation if {
	etl.allow with input as {
		"user": {"roles": ["admin"]},
		"action": "truncate",
		"target_schema": "hims_analytics",
		"confirmation_token": "valid-token-123",
	}
}

test_admin_cannot_truncate_without_confirmation if {
	not etl.allow with input as {
		"user": {"roles": ["admin"]},
		"action": "truncate",
		"target_schema": "hims_analytics",
		"confirmation_token": "",
	}
}

# =============================================================================
# TEST: ETL STATE MANAGEMENT
# =============================================================================

test_etl_service_can_read_state if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "read_state",
	}
}

test_etl_service_can_write_state if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "write_state",
	}
}

test_app_service_can_read_state if {
	etl.allow with input as {
		"user": {"service_account": "sa-app-service", "roles": []},
		"action": "read_state",
	}
}

test_app_service_cannot_write_state if {
	not etl.allow with input as {
		"user": {"service_account": "sa-app-service", "roles": []},
		"action": "write_state",
	}
}

# =============================================================================
# TEST: QUARANTINE OPERATIONS
# =============================================================================

test_etl_service_can_quarantine_write if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "quarantine_write",
	}
}

test_admin_can_reprocess_quarantine if {
	etl.allow with input as {
		"user": {"roles": ["admin"]},
		"action": "quarantine_reprocess",
	}
}

# =============================================================================
# TEST: QUALITY REPORTS
# =============================================================================

test_etl_service_can_write_quality_report if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "quality_report_write",
	}
}

test_authenticated_user_can_read_quality_report if {
	etl.allow with input as {
		"user": {"authenticated": true, "roles": ["viewer"]},
		"action": "quality_report_read",
	}
}

# =============================================================================
# TEST: AUDIT LOG OPERATIONS
# =============================================================================

test_etl_service_can_insert_audit_log if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "audit_log_insert",
	}
}

test_deny_audit_log_update if {
	etl.deny_audit_modification with input as {
		"action": "audit_log_update",
	}
}

test_deny_audit_log_delete if {
	etl.deny_audit_modification with input as {
		"action": "audit_log_delete",
	}
}

# =============================================================================
# TEST: NATS EVENT PUBLISHING
# =============================================================================

test_etl_service_can_publish_sync_completed if {
	etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "nats_publish",
		"subject": "etl.sync.completed",
	}
}

test_etl_service_cannot_publish_unauthorized_subject if {
	not etl.allow with input as {
		"user": {"service_account": "sa-etl-service", "roles": []},
		"action": "nats_publish",
		"subject": "user.created",
	}
}

# =============================================================================
# TEST: DEFAULT DENY
# =============================================================================

test_default_deny_for_unauthenticated if {
	not etl.allow with input as {
		"user": {"roles": []},
		"action": "extract",
		"source": "tally",
	}
}

test_default_deny_for_regular_user if {
	not etl.allow with input as {
		"user": {"roles": ["viewer"]},
		"action": "extract",
		"source": "tally",
	}
}

# =============================================================================
# TEST: POLICY METADATA
# =============================================================================

test_policy_version_exists if {
	etl.policy_version == "1.0.0"
}

test_policy_description_exists if {
	etl.policy_description != ""
}
