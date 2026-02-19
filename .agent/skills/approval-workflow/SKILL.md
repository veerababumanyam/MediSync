---
name: approval-workflow
description: Route transactions through a configurable multi-step approval chain before any data is written to Tally. Enforce separation of duties.
---

# Approval Workflow Skill

Guidelines for implementing secure, auditable, and multi-layered approval workflows for financial transactions in MediSync.

## Workflow Principles

### Separation of Duties
- **Principle**: The user who uploads/maps a transaction should not be the one who gives final approval.
- **Enforcement**: Check JWT `user_id` against the `created_by` field of the workflow instance.

### Role-Based Gates
Define the standard chain:
1. **Accountant Level**: Review for classification accuracy.
2. **Manager Level**: Review for budget compliance (Required if > ₹50,000).
3. **Finance Head Level**: Final sign-off (Required if > ₹5,00,000).

## Tool Chain Documentation

### OPA Policy Enforcement
Use Open Policy Agent (Rego) to define approval logic:
```rego
package medisync.approval

default allow = false

# Allow Accountant to review if they didn't create it
allow {
    input.user_role == "accountant"
    input.action == "review"
    input.user_id != input.transaction_owner
}

# Require Finance Head for large transactions
allow {
    input.user_role == "finance_head"
    input.amount > 500000
}
```

### State Machine Management
Implement persistent state transitions:
```python
# Valid transitions
TRANSITIONS = {
    "draft": ["pending_accountant"],
    "pending_accountant": ["pending_manager", "rejected"],
    "pending_manager": ["pending_finance", "approved", "rejected"],
    "pending_finance": ["approved", "rejected"]
}
```

## Communication & Auditing

- **Notifications**: Use **Apprise** to send real-time alerts to approvers via Email and In-App notifications.
- **Immutable Log**: Every approval or rejection action must be logged with:
    - `timestamp`
    - `user_id`
    - `action` (Approve/Reject/Comment)
    - `rejection_reason` (mandatory if rejected)

## Accessibility Checklist
- [ ] Support bulk approvals for low-value transactions.
- [ ] Show "SLA status" (time pending in current stage) in the dashboard.
- [ ] Enable "Mobile Approval" via a simplified responsive interface.
- [ ] Attach original document (B-02 output) to every approval request.
