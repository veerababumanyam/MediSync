# Agent Specification — C-03: Report Scheduling & Distribution Agent

**Agent ID:** `C-03`  
**Agent Name:** Report Scheduling & Distribution Agent  
**Module:** C — Easy Reports  
**Phase:** 9  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Create, manage, and execute report schedules. Generates reports in PDF/Excel/HTML on a cron schedule and emails them to configured distribution lists. Supports department-wise distribution with role-scoped data.

> **Addresses:** PRD §6.8.5, US19 — Automated report generation and scheduled distribution.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled |
| **Scheduled trigger** | User-defined cron per schedule config |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `schedule_config` | `ReportSchedule` | Schedule store | ✅ |
| `recipient_list` | `[]RecipientConfig` | Schedule config | ✅ |

```go
type ReportSchedule struct {
    ScheduleID   UUID     `json:"schedule_id"`
    ReportType   string   `json:"report_type"`
    Period       string   `json:"period"`       // e.g. "previous_month"
    Cron         string   `json:"cron"`         // e.g. "0 8 1 * *"
    Format       string   `json:"format"`       // pdf/xlsx/html
    EntityID     string   `json:"entity_id"`
    CreatedBy    string   `json:"created_by"`
    Recipients   []RecipientConfig
}
```

---

## 4. Outputs

| Output | Type | Description |
|--------|------|-------------|
| `delivery_results` | `[]DeliveryResult` | Per-recipient send status |
| `report_archive_url` | `string` | Stored copy in object storage |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Go cron (`robfig/cron`) | MIT | Trigger schedule |
| 2 | C-01 | Internal | Generate report |
| 3 | Apprise | MIT | Email delivery per recipient |

---

## 6. Guardrails

- Each recipient receives data scoped to their own role (report regenerated per recipient if roles differ).
- Schedules auto-paused if owner's account is inactive.
- Failed delivery retried 3×; owner notified on total failure.

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| On-time delivery rate | ≥ 99% |
| Delivery failure rate | < 0.5% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go background worker |
| **Depends on** | C-01, Apprise, object storage |
| **Consumed by** | Finance Head, Dept Managers |
