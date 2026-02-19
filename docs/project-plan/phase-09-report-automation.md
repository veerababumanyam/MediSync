# Phase 09 — Report Automation & Distribution (Easy Reports — Part 2)

**Phase Duration:** Weeks 31–33 (3 weeks)  
**Module(s):** Module C (Easy Reports)  
**Status:** Planning  
**Depends On:** Phase 08 complete (pre-built reports library live)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.3](../ARCHITECTURE.md)

---

## 1. Objectives

Automate report delivery: users configure schedules and the system generates, formats, and emails reports automatically. This phase eliminates the manual effort of weekly/monthly report compilation — one of the core PRD business objectives ("Reduce manual reporting time by 90%").

---

## 2. Scope

### In Scope
- C-03 Report Scheduling & Distribution Agent (full implementation)
- Email delivery with multi-format attachments (PDF, Excel, CSV, HTML)
- Report portal — self-service access to historical scheduled reports
- Distribution list management
- Schedule management UI (create/edit/delete/pause schedules)
- Report delivery tracking (success/failure logs)
- Mobile-friendly report portal
- Conditional scheduling (run only if threshold breached)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | C-03 Report Scheduling Agent | AI Engineer + Backend | Cron-triggered; generates report + emails in configured format; tested with 5 schedule types |
| D-02 | SMTP Email Service Integration | Backend Engineer | SendGrid / AWS SES / Postfix configured; email delivered with PDF/Excel attachment |
| D-03 | Schedule Management UI | Frontend Engineer | Create/edit/delete schedules; cron expression builder; recipient list management |
| D-04 | Report Portal | Frontend Engineer | Self-service list of all scheduled + historical reports; filter by date/type; download |
| D-05 | Delivery Tracking | Backend + Frontend | Success/failure status per delivery; retry mechanism for failed deliveries |
| D-06 | Conditional Scheduling | AI Engineer | "Run only if revenue drops > 10% vs last period" logic supported |
| D-07 | Multi-format Attachments | Backend Engineer | Same report attached in user's choice of PDF + Excel + CSV simultaneously |

---

## 4. AI Agent: C-03 Report Scheduling & Distribution

**Type:** Proactive (L2) — cron-driven  

**Schedule types supported:**
| Type | Cron Example | Description |
|---|---|---|
| Daily | `0 7 * * *` | Every day at 7:00 AM |
| Weekly | `0 7 * * 1` | Every Monday at 7:00 AM |
| Monthly | `0 7 1 * *` | 1st of every month at 7:00 AM |
| Quarterly | `0 7 1 1,4,7,10 *` | First day of each quarter |
| Custom | Any cron expression | User-defined |
| Conditional | Cron + threshold check | e.g., Monthly but only if budget overrun > 5% |

**Execution flow:**
```
Cron trigger fires
    │
    ▼ Load schedule config from app.scheduled_reports
    │   - report_type, params, date_range_mode, locale
    │   - format list (pdf|xlsx|csv), recipients []
    │   - conditional rule (optional)
    │
    ▼ IF conditional rule: evaluate metric threshold
    │   Condition not met → skip + log "skipped: condition not met"
    │
    ▼ Execute C-01 with date_range_mode = 'previous_period'
    │   (auto-calculates prior month/week/quarter)
    │
    ▼ Render output per format:
    │   PDF  → WeasyPrint (locale-aware RTL if ar)
    │   XLSX → excelize
    │   CSV  → Go encoding/csv
    │
    ▼ Email assembly:
    │   Subject: "MediSync — {{report_name}} — {{period}} [{{locale}}]"
    │   Body: HTML summary + key metrics + download link
    │   Attachments: up to 3 format files
    │
    ▼ SMTP delivery (retry 3× on failure)
    │
    ▼ Log to app.report_deliveries (status, timestamp, size)
    │
    ▼ NATS: report.scheduled.delivered
```

**Stale report protection:** If ETL sync failed in the report generation window, report is delayed 1 hour (auto-retry) with warning note. If data is > 4 hours stale when report generates, a data-freshness warning is included in the email body.

---

## 5. Schedule Management UI

**Create Schedule Flow:**
1. Select Report Type from pre-built library
2. Configure parameters (period mode, department filter, entity)
3. Set schedule (frequency picker → cron expression preview)
4. Configure recipients (add emails; assign MediSync user roles to auto-populate their email)
5. Select output formats (PDF / Excel / CSV — multi-select)
6. Optional: add conditional trigger rule
7. Select locale (EN / AR / Bilingual)
8. Test Run: "Send Now" to verify report generates and emails correctly
9. Save & Activate

**Schedule List View:**
- Active / Paused / Disabled status badges
- Next run time, last run time, last status
- Quick pause/resume toggle
- Edit / Delete / Duplicate actions

---

## 6. Report Portal

**Features:**
- Chronological list of all generated reports
- Filter by: report type, period, entity, status (delivered/failed)
- Download any past report in original format
- "Regenerate" option to produce fresh version of any past report
- Retention: reports kept for 24 months online; archived to cold storage thereafter

**Notification on delivery:** E-06 sends in-app notification to report schedule owner when delivery completes or fails.

---

## 7. Email Configuration

```yaml
# Application config
smtp:
  host: ${SMTP_HOST}          # pulled from Vault
  port: 587
  username: ${SMTP_USERNAME}
  password: ${SMTP_PASSWORD}
  from: "MediSync Reports <reports@medisync.local>"
  tls: true

email_templates:
  scheduled_report:
    subject_en: "MediSync — {{.ReportName}} — {{.Period}}"
    subject_ar: "ميدي سينك — {{.ReportName}} — {{.Period}}"
    body_template: "email_report_body.html"
```

**Email body includes:**
- MediSync logo (inline base64)
- Report summary: up to 5 headline KPIs from the report
- "View Full Report" link to portal (JWT-authenticated)
- "Unsubscribe from this schedule" link
- Footer with locale, period, generation timestamp

---

## 8. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| C-03 cron execution | 5 schedule types (daily/weekly/monthly/quarterly/custom) | All fire at correct times; report generated |
| Conditional schedules | 5 conditional rules (3 met, 2 not met) | Correct send/skip decisions |
| Email delivery | Send to 5 test recipients in 3 formats | Received in inbox; attachment opens correctly |
| Arabic schedule | Schedule with locale=ar | PDF attachment is Arabic RTL |
| Report portal | Access last 10 historical reports | All downloadable; regeneration works |
| Stale data handling | Simulate ETL failure before scheduled run | Report delayed; warning note included |
| Retry on SMTP failure | Simulate SMTP error on first attempt | Retried 3×; final failure logged + notified |

---

## 9. Phase Exit Criteria

- [ ] C-03 scheduling agent running all 5 schedule types correctly
- [ ] Email delivery working with PDF + Excel + CSV attachments
- [ ] Arabic-locale scheduled reports delivering Arabic PDFs
- [ ] Report portal accessible with full history download
- [ ] Delivery tracking dashboard showing success/failure status
- [ ] Conditional scheduling logic tested and working
- [ ] Phase gate reviewed and signed off

---

*Phase 09 | Version 1.0 | February 19, 2026*
