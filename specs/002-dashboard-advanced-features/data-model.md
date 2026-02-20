# Data Model: Dashboard, Chat UI & Advanced Features with i18n

**Feature**: 002-dashboard-advanced-features
**Date**: 2026-02-20
**Version**: 1.0

---

## Entity Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER (existing)                                 │
│  • id: UUID                                                                  │
│  • email: string                                                             │
│  • roles: string[]                                                           │
│  • organization_id: UUID                                                     │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │
        ┌───────────────────────┼───────────────────────┐
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────────┐
│ UserPreference│     │  PinnedChart  │     │    AlertRule      │
├───────────────┤     ├───────────────┤     ├───────────────────┤
│ user_id (FK)  │     │ user_id (FK)  │     │ user_id (FK)      │
│ locale        │     │ query_id      │     │ metric_id         │
│ numeral_system│     │ chart_spec    │     │ operator          │
│ calendar      │     │ refresh_interval│   │ threshold         │
│ report_lang   │     │ locale        │     │ channels          │
└───────────────┘     │ position      │     │ locale            │
                      └───────────────┘     │ last_triggered_at │
                                            │ is_active         │
                                            └─────────┬─────────┘
                                                      │
                                                      ▼
                                            ┌───────────────────┐
                                            │   Notification    │
                                            ├───────────────────┤
                                            │ alert_rule_id (FK)│
                                            │ user_id (FK)      │
                                            │ type              │
                                            │ status            │
                                            │ content           │
                                            │ locale            │
                                            │ sent_at           │
                                            └───────────────────┘

┌───────────────────┐     ┌───────────────────┐
│  ScheduledReport  │     │   ChatMessage     │
├───────────────────┤     ├───────────────────┤
│ user_id (FK)      │     │ session_id        │
│ query_id          │     │ user_id (FK)      │
│ schedule_type     │     │ role              │
│ schedule_time     │     │ content           │
│ recipients        │     │ chart_spec        │
│ format            │     │ table_data        │
│ locale            │     │ locale            │
│ last_run_at       │     │ created_at        │
│ next_run_at       │     └───────────────────┘
│ is_active         │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ ScheduledReportRun│
├───────────────────┤
│ report_id (FK)    │
│ status            │
│ file_path         │
│ error_message     │
│ started_at        │
│ completed_at      │
└───────────────────┘
```

---

## Entity Definitions

### 1. UserPreference

User-specific display and formatting preferences that persist across sessions and devices.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `user_id` | UUID | FK → users.id, UNIQUE, NOT NULL | Owner reference |
| `locale` | enum | 'en' \| 'ar', DEFAULT 'en' | Display language |
| `numeral_system` | enum | 'western' \| 'eastern_arabic', DEFAULT 'western' | Number rendering |
| `calendar_system` | enum | 'gregorian' \| 'hijri', DEFAULT 'gregorian' | Date display calendar |
| `report_language` | enum | 'en' \| 'ar', DEFAULT 'en' | Language for generated reports |
| `timezone` | string | DEFAULT 'Asia/Dubai' | User timezone for scheduling |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Record creation time |
| `updated_at` | timestamp | NOT NULL, DEFAULT NOW() | Last modification time |

**Indexes**:
- `idx_user_preference_user_id` UNIQUE on `user_id`

**Validation Rules**:
- `locale` must match `report_language` or user must explicitly choose different
- `timezone` must be valid IANA timezone identifier

---

### 2. PinnedChart

A saved visualization displayed on the user's personal dashboard.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `user_id` | UUID | FK → users.id, NOT NULL | Owner reference |
| `title` | string | NOT NULL, max 200 chars | Display title (localized) |
| `query_id` | UUID | NULLABLE | Reference to saved query (optional) |
| `natural_language_query` | text | NOT NULL | Original question text |
| `sql_query` | text | NOT NULL | Generated SQL (for re-execution) |
| `chart_spec` | jsonb | NOT NULL | ECharts configuration object |
| `chart_type` | enum | 'bar' \| 'line' \| 'pie' \| 'table' \| 'kpi', NOT NULL | Chart category |
| `refresh_interval` | integer | DEFAULT 300 | Auto-refresh seconds (0 = disabled) |
| `locale` | enum | 'en' \| 'ar', NOT NULL | Locale when chart was created |
| `position` | jsonb | NOT NULL, DEFAULT '{"row":0,"col":0,"size":1}' | Grid position |
| `last_refreshed_at` | timestamp | NULLABLE | Last data refresh time |
| `is_active` | boolean | DEFAULT true | Soft delete flag |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Pin creation time |
| `updated_at` | timestamp | NOT NULL, DEFAULT NOW() | Last modification time |

**Indexes**:
- `idx_pinned_chart_user_id` on `user_id`
- `idx_pinned_chart_user_active` on `(user_id, is_active)`

**Position JSON Schema**:
```json
{
  "row": 0,       // Grid row (0-indexed)
  "col": 0,       // Grid column (0-indexed)
  "size": 1       // Widget size multiplier (1, 2, or 3 columns)
}
```

**Chart Spec JSON Schema** (ECharts):
```json
{
  "title": { "text": "...", "left": "center" },
  "xAxis": { "type": "category", "data": [...] },
  "yAxis": { "type": "value" },
  "series": [{ "type": "bar", "data": [...] }],
  "tooltip": { "trigger": "axis" }
}
```

**Validation Rules**:
- `title` must not be empty
- `sql_query` must be SELECT-only (validated by existing agent layer)
- `refresh_interval` must be 0 or between 60 and 3600 seconds
- `position.size` must be 1, 2, or 3

**State Transitions**:
```
[Created] → active=true → [Active] → active=false → [Archived]
                              │
                              ▼ (refresh_interval > 0)
                         [Refreshing] → [Active]
```

---

### 3. AlertRule

User-defined condition that triggers notifications when a metric crosses a threshold.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `user_id` | UUID | FK → users.id, NOT NULL | Owner reference |
| `name` | string | NOT NULL, max 100 chars | Alert rule name |
| `description` | text | NULLABLE | User notes about the alert |
| `metric_id` | string | NOT NULL | Identifier for the tracked metric |
| `metric_name` | string | NOT NULL | Human-readable metric name (localized) |
| `operator` | enum | 'gt' \| 'gte' \| 'lt' \| 'lte' \| 'eq', NOT NULL | Comparison operator |
| `threshold` | decimal | NOT NULL | Threshold value |
| `check_interval` | integer | NOT NULL, DEFAULT 300 | Evaluation interval in seconds |
| `channels` | jsonb | NOT NULL, DEFAULT '["in_app"]' | Notification delivery channels |
| `locale` | enum | 'en' \| 'ar', NOT NULL | Language for alert messages |
| `cooldown_period` | integer | DEFAULT 3600 | Min seconds between alerts |
| `last_triggered_at` | timestamp | NULLABLE | Last trigger time |
| `last_value` | decimal | NULLABLE | Last observed metric value |
| `is_active` | boolean | DEFAULT true | Rule enabled flag |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Rule creation time |
| `updated_at` | timestamp | NOT NULL, DEFAULT NOW() | Last modification time |

**Indexes**:
- `idx_alert_rule_user_id` on `user_id`
- `idx_alert_rule_user_active` on `(user_id, is_active)`
- `idx_alert_rule_next_check` on `(is_active, last_triggered_at)` for scheduler

**Channels JSON Schema**:
```json
["in_app", "email"]
```

**Operator Definitions**:
| Operator | Meaning | Example |
|----------|---------|---------|
| `gt` | Greater than | value > threshold |
| `gte` | Greater than or equal | value >= threshold |
| `lt` | Less than | value < threshold |
| `lte` | Less than or equal | value <= threshold |
| `eq` | Equal to | value == threshold |

**Validation Rules**:
- `metric_id` must reference a valid metric in the metrics registry
- `threshold` must be a valid decimal number
- `check_interval` must be between 60 and 86400 seconds (1 min to 1 day)
- `cooldown_period` must be between 0 and 86400 seconds
- At least one channel must be specified

**State Transitions**:
```
[Created] → is_active=true → [Monitoring]
                              │
                              ▼ (condition met)
                         [Triggered] → [Cooldown] → [Monitoring]
                              │
                              ▼ is_active=false
                         [Paused]
```

---

### 4. Notification

Record of an alert delivery attempt.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `alert_rule_id` | UUID | FK → alert_rules.id, NOT NULL | Associated alert rule |
| `user_id` | UUID | FK → users.id, NOT NULL | Recipient reference |
| `type` | enum | 'in_app' \| 'email', NOT NULL | Delivery channel |
| `status` | enum | 'pending' \| 'sent' \| 'delivered' \| 'failed', NOT NULL | Delivery status |
| `content` | jsonb | NOT NULL | Message content (localized) |
| `locale` | enum | 'en' \| 'ar', NOT NULL | Message language |
| `metric_value` | decimal | NOT NULL | Value that triggered the alert |
| `threshold` | decimal | NOT NULL | Threshold that was crossed |
| `error_message` | text | NULLABLE | Error if delivery failed |
| `sent_at` | timestamp | NULLABLE | When notification was sent |
| `delivered_at` | timestamp | NULLABLE | When notification was confirmed delivered |
| `read_at` | timestamp | NULLABLE | When user read the notification |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Record creation time |

**Indexes**:
- `idx_notification_user_id` on `user_id`
- `idx_notification_alert_rule_id` on `alert_rule_id`
- `idx_notification_user_unread` on `(user_id, read_at)` WHERE `read_at IS NULL`

**Content JSON Schema**:
```json
{
  "title": "Low Stock Alert",
  "message": "Pharmacy inventory for Paracetamol has fallen to 45 units (threshold: 50)",
  "action_url": "/dashboard/alerts/123"
}
```

**Validation Rules**:
- `read_at` can only be set if `delivered_at` is set
- `delivered_at` can only be set if `status` is 'delivered'

---

### 5. ScheduledReport

Configuration for a recurring report generation and delivery.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `user_id` | UUID | FK → users.id, NOT NULL | Owner reference |
| `name` | string | NOT NULL, max 150 chars | Report name |
| `description` | text | NULLABLE | Report description |
| `query_id` | UUID | NULLABLE | Reference to saved query |
| `natural_language_query` | text | NOT NULL | Original question text |
| `sql_query` | text | NOT NULL | SQL to generate report data |
| `schedule_type` | enum | 'daily' \| 'weekly' \| 'monthly' \| 'quarterly', NOT NULL | Recurrence pattern |
| `schedule_time` | time | NOT NULL | Time of day to run |
| `schedule_day` | integer | NULLABLE | Day of week (1-7) or month (1-31) |
| `recipients` | jsonb | NOT NULL | List of email recipients |
| `format` | enum | 'pdf' \| 'xlsx' \| 'csv', NOT NULL | Output format |
| `locale` | enum | 'en' \| 'ar', NOT NULL | Report language |
| `include_charts` | boolean | DEFAULT true | Include visualizations |
| `last_run_at` | timestamp | NULLABLE | Last generation time |
| `next_run_at` | timestamp | NULLABLE | Next scheduled run |
| `is_active` | boolean | DEFAULT true | Schedule enabled flag |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Schedule creation time |
| `updated_at` | timestamp | NOT NULL, DEFAULT NOW() | Last modification time |

**Indexes**:
- `idx_scheduled_report_user_id` on `user_id`
- `idx_scheduled_report_next_run` on `(is_active, next_run_at)` for scheduler

**Recipients JSON Schema**:
```json
[
  { "email": "manager@clinic.com", "name": "Clinic Manager" },
  { "email": "finance@clinic.com", "name": "Finance Head" }
]
```

**Schedule Day Logic**:
| Schedule Type | Day Value | Meaning |
|---------------|-----------|---------|
| `daily` | NULL | Every day at schedule_time |
| `weekly` | 1-7 | Day of week (1=Monday, 7=Sunday) |
| `monthly` | 1-31 | Day of month (capped to month's last day) |
| `quarterly` | 1 | First day of quarter (Jan 1, Apr 1, Jul 1, Oct 1) |

**Validation Rules**:
- `recipients` must contain at least one valid email address
- `schedule_day` must be valid for `schedule_type`
- `sql_query` must be SELECT-only
- `next_run_at` is auto-calculated on creation and after each run

---

### 6. ScheduledReportRun

Audit record of a report generation attempt.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `report_id` | UUID | FK → scheduled_reports.id, NOT NULL | Parent report |
| `status` | enum | 'pending' \| 'running' \| 'completed' \| 'failed', NOT NULL | Run status |
| `file_path` | string | NULLABLE | Path to generated file in storage |
| `file_size_bytes` | bigint | NULLABLE | Size of generated file |
| `row_count` | integer | NULLABLE | Number of data rows in report |
| `error_message` | text | NULLABLE | Error if run failed |
| `started_at` | timestamp | NOT NULL | When run started |
| `completed_at` | timestamp | NULLABLE | When run finished |

**Indexes**:
- `idx_scheduled_report_run_report_id` on `report_id`
- `idx_scheduled_report_run_status` on `status`

**State Transitions**:
```
[Pending] → status='running' → [Running]
                                   │
                      ┌────────────┼────────────┐
                      ▼            ▼            ▼
                 [Completed]  [Failed]    [Timeout]
```

---

### 7. ChatMessage

A message in a chat conversation, either from user or AI.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK, auto-generated | Primary key |
| `session_id` | UUID | NOT NULL | Conversation session identifier |
| `user_id` | UUID | FK → users.id, NOT NULL | User who sent/received message |
| `role` | enum | 'user' \| 'assistant', NOT NULL | Message sender |
| `content` | text | NOT NULL | Message text content |
| `chart_spec` | jsonb | NULLABLE | ECharts config if response includes chart |
| `table_data` | jsonb | NULLABLE | Tabular data if response includes table |
| `drilldown_query` | text | NULLABLE | SQL for drill-down (if applicable) |
| `confidence_score` | decimal | NULLABLE, 0.0-1.0 | AI confidence (Module A-06) |
| `locale` | enum | 'en' \| 'ar', NOT NULL | Message language |
| `created_at` | timestamp | NOT NULL, DEFAULT NOW() | Message timestamp |

**Indexes**:
- `idx_chat_message_session_id` on `session_id`
- `idx_chat_message_user_id` on `user_id`
- `idx_chat_message_created_at` on `created_at` DESC

**Table Data JSON Schema**:
```json
{
  "columns": ["Department", "Revenue", "Change"],
  "rows": [
    ["Cardiology", 125000, "+12%"],
    ["Orthopedics", 98000, "-3%"]
  ],
  "total_rows": 15
}
```

**Validation Rules**:
- `confidence_score` must be between 0.0 and 1.0
- `chart_spec` and `table_data` are mutually exclusive for same message
- `session_id` is generated client-side and persisted for conversation continuity

---

## Database Migrations

### Migration Order

```text
migrations/
├── 010_user_preferences.up.sql    # UserPreference table
├── 011_pinned_charts.up.sql       # PinnedChart table
├── 012_alert_rules.up.sql         # AlertRule + Notification tables
├── 013_scheduled_reports.up.sql   # ScheduledReport + ScheduledReportRun tables
└── 014_chat_messages.up.sql       # ChatMessage table
```

### Sample Migration (010_user_preferences.up.sql)

```sql
-- 010_user_preferences.up.sql
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    locale VARCHAR(2) NOT NULL DEFAULT 'en' CHECK (locale IN ('en', 'ar')),
    numeral_system VARCHAR(20) NOT NULL DEFAULT 'western'
        CHECK (numeral_system IN ('western', 'eastern_arabic')),
    calendar_system VARCHAR(20) NOT NULL DEFAULT 'gregorian'
        CHECK (calendar_system IN ('gregorian', 'hijri')),
    report_language VARCHAR(2) NOT NULL DEFAULT 'en'
        CHECK (report_language IN ('en', 'ar')),
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Dubai',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_preference_user_id UNIQUE (user_id)
);

CREATE INDEX idx_user_preference_user_id ON user_preferences(user_id);

-- Trigger for updated_at
CREATE TRIGGER update_user_preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

## Relationships Summary

| From | To | Relationship | Description |
|------|-----|--------------|-------------|
| UserPreference | User | N:1 | Each user has one preference set |
| PinnedChart | User | N:1 | Users can have multiple pinned charts |
| AlertRule | User | N:1 | Users can have multiple alert rules |
| Notification | AlertRule | N:1 | Each alert can generate multiple notifications |
| Notification | User | N:1 | Notifications belong to a user |
| ScheduledReport | User | N:1 | Users can have multiple scheduled reports |
| ScheduledReportRun | ScheduledReport | N:1 | Each report has multiple run records |
| ChatMessage | User | N:1 | Messages belong to a user |

---

## Data Access Patterns

### Dashboard Load
```sql
-- Get all active pinned charts for user with latest data
SELECT id, title, chart_spec, chart_type, position, last_refreshed_at
FROM pinned_charts
WHERE user_id = $1 AND is_active = true
ORDER BY position->>'row', position->>'col';
```

### Alert Evaluation
```sql
-- Get alerts due for evaluation
SELECT ar.*, up.locale, up.timezone
FROM alert_rules ar
JOIN user_preferences up ON ar.user_id = up.user_id
WHERE ar.is_active = true
  AND (ar.last_triggered_at IS NULL
       OR ar.last_triggered_at + (ar.cooldown_period || ' seconds')::interval < NOW())
  AND ar.last_triggered_at + (ar.check_interval || ' seconds')::interval < NOW();
```

### Report Scheduling
```sql
-- Get reports due for generation
SELECT sr.*, up.locale, up.timezone
FROM scheduled_reports sr
JOIN user_preferences up ON sr.user_id = up.user_id
WHERE sr.is_active = true
  AND sr.next_run_at <= NOW();
```

### Chat History
```sql
-- Get recent messages for session
SELECT id, role, content, chart_spec, table_data, confidence_score, created_at
FROM chat_messages
WHERE session_id = $1
ORDER BY created_at ASC
LIMIT 100;
```

---

*Data Model Version: 1.0 | Last Updated: 2026-02-20*
