# Agent Specification — A-09: Report Scheduling Agent

**Agent ID:** `A-09`  
**Agent Name:** Report Scheduling Agent  
**Module:** A — Conversational BI Dashboard  
**Phase:** 3  
**Priority:** P1 High  
**HITL Required:** No  
**Status:** Draft

---

## 1. Purpose

Generates scheduled reports on a cron schedule, formats output in PDF/Excel/CSV, and triggers delivery via email to configurable recipient lists.

> **Addresses:** PRD §6.2, US5 — Automated report scheduling and email delivery.

---

## 2. Trigger

| Property | Value |
|----------|-------|
| **Trigger type** | Scheduled |
| **Scheduled trigger** | User-defined cron (e.g. `0 8 * * 1` = every Monday 8 AM) |
| **Calling agent** | Scheduler daemon |

---

## 3. Inputs

| Input | Type | Source | Required |
|-------|------|--------|:--------:|
| `schedule_id` | `UUID` | Schedule config store | ✅ |
| `report_query` | `string` | Saved SQL / natural language query | ✅ |
| `format` | `enum` | `pdf / xlsx / csv` | ✅ |
| `recipients` | `[]string` | Email list | ✅ |
| `user_id` | `string` | Owner of schedule | ✅ |

---

## 4. Outputs

| Output | Type | Destination |
|--------|------|-------------|
| `report_file` | `bytes` | Email attachment + object storage |
| `delivery_status` | `enum` | `sent / failed` |
| `delivery_timestamp` | `datetime` | Audit log |

---

## 5. Tool Chain

| Step | Tool | License | Purpose |
|------|------|---------|---------|
| 1 | Go cron scheduler (`robfig/cron`) | MIT | Trigger on schedule |
| 2 | A-01 | Internal | Execute report query |
| 3 | `excelize` (Go) | BSD | Excel generation |
| 4 | `wkhtmltopdf` / `chromedp` | LGPL / MIT | PDF rendering |
| 5 | Apprise | MIT | Email delivery |

---

## 6. Guardrails

- Schedules owned by deactivated users are auto-paused.
- Delivery failures retry 3× then notify schedule owner.
- Report queries run under the owning user's role (column masking applies).

---

## 7. Evaluation Criteria

| Metric | Target |
|--------|--------|
| On-time delivery rate | ≥ 99% |
| Delivery failure rate | < 1% |

---

## 8. Deployment Notes

| Property | Value |
|----------|-------|
| **Runtime** | Go service (background worker) |
| **Queue** | Redis scheduled jobs |
| **Depends on** | A-01, Apprise |
| **Consumed by** | User-configured schedules |
