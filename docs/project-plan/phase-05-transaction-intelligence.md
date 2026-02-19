# Phase 05 — Transaction Intelligence (AI Accountant — Part 2)

**Phase Duration:** Weeks 16–19 (4 weeks)  
**Module(s):** Module B (AI Accountant), Module E (E-04)  
**Status:** Planning  
**Milestone:** Ledger Mapping & Approval Workflow Live  
**Depends On:** Phase 04 complete (OCR extraction pipeline operational)  
**Cross-ref:** [PROJECT-PLAN.md](./PROJECT-PLAN.md) | [ARCHITECTURE.md §5.2](../ARCHITECTURE.md)

---

## 1. Objectives

Transform extracted document data into validated, AI-mapped accounting entries ready for Tally. This phase adds vendor matching, intelligent GL ledger mapping, duplicate invoice detection, sub-ledger assignment, and the 4-level human approval workflow that gates all Tally write operations. Also delivers E-04 multilingual PDF report generation.

---

## 2. Scope

### In Scope
- B-04 Vendor Matching Agent
- B-05 Ledger Mapping Agent (with learning feedback loop)
- B-06 Sub-Ledger & Cost Centre Assignment Agent
- B-07 Duplicate Invoice Detection Agent
- B-08 Approval Workflow Agent (full 4-level chain)
- E-04 Multilingual Report Generator (PDF/Excel in EN + AR)
- Approval Workflow UI (accountant → manager → finance → posted)
- Transaction review and edit UI
- Confidence scoring UI (traffic-light badges)
- Learning feedback loop (corrections improve future mappings)

### Out of Scope
- Tally sync (Phase 6)
- Bank reconciliation (Phase 7)

---

## 3. Deliverables

| # | Deliverable | Owner | Acceptance Criteria |
|---|---|---|---|
| D-01 | B-04 Vendor Matching Agent | AI Engineer | Matches extracted vendor name to Tally vendor master; creates new vendor record if no match; ≥ 90% match accuracy |
| D-02 | B-05 Ledger Mapping Agent | AI Engineer | Suggests correct GL ledger per transaction with confidence score; ≥ 85% first-suggestion accuracy on 50-transaction test set |
| D-03 | B-06 Sub-Ledger & Cost Centre Agent | AI Engineer | Assigns sub-ledger and cost centre based on context; tested against 30 transactions |
| D-04 | B-07 Duplicate Invoice Detection | AI Engineer | Detects duplicate: same vendor + amount + date ± 3 days; SHA-256 hash check; 100% detection on 20-duplicate test set |
| D-05 | B-08 Approval Workflow Agent | AI Engineer | Routes transactions through 4-level chain; sends reminders for stale (> 24h) approvals; state machine correct |
| D-06 | Transaction Review UI | Frontend Engineer | Side-by-side: document preview + AI-suggested fields; confidence badges; editable fields |
| D-07 | Approval Workflow UI | Frontend Engineer | Clear approval chain visualisation; approve/reject per level; comments field |
| D-08 | Confidence Badges | Frontend Engineer | Green (≥ 90%), Amber (70–89%), Red (< 70%) traffic-light badges on every AI-suggested field |
| D-09 | Learning Feedback Loop | AI Engineer | User corrections stored in `app.mapping_corrections`; A-05 Ledger Mapping picks up corrections as few-shot examples in next run |
| D-10 | E-04 Multilingual Report Generator | AI Engineer | PDFs correctly render in English OR Arabic with RTL layout; WeasyPrint + Cairo fonts verified |

---

## 4. AI Agents Deployed

### B-04 Vendor Matching Agent

**Trigger:** NATS `document.extracted`  
**Strategy:**
1. Extract vendor name from extraction result
2. Fuzzy string match against `tally_analytics.dim_ledgers` (vendor ledger group)
3. Vector similarity match (vendor name embedding vs. dim_ledgers name embeddings)
4. If match confidence ≥ 85%: auto-assign vendor
5. If 70–84%: suggest + flag for human confirmation
6. If < 70% OR no close match: propose creating new vendor record → HITL

**Duplicate vendor prevention:** Before creating new vendor, checks for existing ones with normalized name similarity > 90%

### B-05 Ledger Mapping Agent

**Trigger:** NATS `document.extracted` (runs after B-04)  
**Strategy:**

```
Transaction context (vendor, amount, description, document type)
    │
    ▼ Load user's historical mapping corrections as few-shot examples
    │
    ▼ Load Tally Chart of Accounts from dim_ledgers
    │
    ▼ LLM prompt:
    │  "Given this transaction: [amount, vendor, description], 
    │   and these Tally ledger options: [list],
    │   and these past mappings: [few-shot examples],
    │   suggest the most appropriate GL ledger. 
    │   Respond with: ledger_id, ledger_name, confidence (0-100), reason"
    │
    ▼ Top-3 suggestions with confidence scores
    │
    ▼ Best suggestion ≥ 90% → auto-map + flag as recommended
    │   Best suggestion < 70% → HITL required
```

**Learning feedback:** `app.mapping_corrections` table stores: original suggestion, corrected ledger, transaction features. Retrieved as few-shot examples in future prompts for same vendor/description pattern.

**Confidence routing:**
- ≥ 90%: auto-map with green badge; accountant can override
- 70–89%: amber badge; accountant reviews before proceeding
- < 70%: red badge; HITL required before advancing in pipeline

### B-06 Sub-Ledger & Cost Centre Assignment

**Trigger:** Runs after B-05 (same pipeline step)  
**Logic:**
- Cost centre assignment based on: invoice department (from vendor category or description keywords)
- Sub-ledger assignment based on: parent ledger type + historical pattern
- Both suggestions have confidence scores
- All assignments are overridable by accountant

### B-07 Duplicate Invoice Detection Agent

**Two-layer detection:**

1. **Hash-based exact dupe:**
   - SHA-256 hash of `{vendor_id}:{invoice_number}:{amount}`
   - Match → immediate detection, hard block

2. **Fuzzy dupe (soft duplicate):**
   - Same vendor + amount within ± 2% + date within ± 7 days
   - Match → warning flag, requires human decision

**Action on detected dupe:**
- Hard dupe: block pipeline; alert accountant; show original + duplicate side-by-side
- Soft dupe: flag with warning badge; accountant can "Confirm Not Duplicate" or "Mark as Duplicate"

### B-08 Approval Workflow Agent

**Type:** Reactive + Proactive (sends reminders)  

**Approval chain:**
```
accountant (Level 1: initial review + submit)
    │
    ▼ STATUS: pending_manager_approval
accountant_lead OR manager (Level 2: business review)
    │
    ▼ STATUS: pending_finance_approval
finance_head (Level 3: financial authorisation)
    │
    ▼ STATUS: approved
    │
    ▼ (Phase 6) B-09 Tally Sync (Level 4: actual Tally posting)
```

**State machine states:** `draft` → `pending_l1` → `pending_l2` → `pending_l3` → `approved` → `synced_to_tally` | `rejected` | `cancelled`

**Stale approval reminders:**
- 24 hours without action → email reminder to approver
- 48 hours without action → escalate to approver's manager + email

**Self-approval prevention:** OPA policy blocks a user from approving their own submission at any level.

**HITL:** This agent ALWAYS requires humans. There is no auto-approval path.

### E-04 Multilingual Report Generator

**Trigger:** Report generation request with `locale` parameter  
**PDF rendering stack:**
```
GoHTML template (locale-aware)
    │ WeasyPrint (HTML → PDF)
    ▼
PDF with correct:
    - Font: Cairo (Arabic) / Roboto (English)
    - Direction: RTL for AR, LTR for EN
    - Number format: ١٬٢٣٤ (AR) / 1,234 (EN)
    - Date format: ١٩ فبراير ٢٠٢٦ (AR) / 19 Feb 2026 (EN)
```

**Excel rendering:**
```
excelize Go library
    - RightToLeft: true for Arabic workbooks
    - All strings: UTF-8 with Arabic Unicode support
    - Number formats: locale-specific
```

---

## 5. Database Schema Additions

```sql
CREATE TABLE transaction_queue (
    txn_id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doc_id              UUID REFERENCES documents(doc_id),
    extraction_id       UUID REFERENCES extracted_documents(extraction_id),
    vendor_id           UUID,             -- resolved Tally vendor ledger ID
    vendor_confidence   NUMERIC(5,2),
    gl_ledger_id        UUID,             -- suggested Tally GL ledger
    ledger_confidence   NUMERIC(5,2),
    sub_ledger_id       UUID,
    cost_centre_id      UUID,
    amount              NUMERIC(15,2),
    tax_amount          NUMERIC(15,2),
    txn_date            DATE,
    narration           TEXT,
    is_duplicate        BOOLEAN DEFAULT FALSE,
    duplicate_of        UUID,             -- FK to another txn_id if dupe
    approval_status     VARCHAR(30) DEFAULT 'draft',
    current_approver_role VARCHAR(30),
    created_by          UUID REFERENCES users(user_id),
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE approval_history (
    approval_id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    txn_id          UUID REFERENCES transaction_queue(txn_id),
    approver_id     UUID REFERENCES users(user_id),
    approver_role   VARCHAR(30),
    action          VARCHAR(20),  -- 'approved' | 'rejected' | 'returned'
    comments        TEXT,
    approved_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mapping_corrections (
    correction_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    txn_id          UUID REFERENCES transaction_queue(txn_id),
    field_name      VARCHAR(50),   -- 'gl_ledger_id' | 'cost_centre_id' | 'vendor_id'
    ai_suggestion   TEXT,
    human_correction TEXT,
    corrected_by    UUID REFERENCES users(user_id),
    corrected_at    TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 6. Approval Workflow UI

**Screens:**

1. **Transaction List** — tabs: Draft | Pending Review | Pending Approval | Approved | Rejected  
2. **Transaction Detail** — document preview + all AI-mapped fields + confidence badges + edit controls  
3. **Approval Panel** — shows current approval level, chain status, approve/reject/return buttons + comment  
4. **Approval History** — timeline of all approval actions for a transaction  
5. **Notification Inbox** — all approval requests and reminders for the logged-in user

**Roles and what they see:**
- `accountant`: Draft + Pending L1 + their own submitted transactions
- `accountant_lead` / `manager`: Pending L2 transactions
- `finance_head`: Pending L3 transactions; also sees full approval history
- `admin`: All transactions

---

## 7. NATS Event Topics (Phase 05)

```
document.extracted      → consumed by: B-04 Vendor Matching
vendor.matched          → consumed by: B-05 Ledger Mapping
ledger.mapped           → consumed by: B-06 Sub-Ledger Assignment
transaction.prepared    → consumed by: B-07 Duplicate Detection
transaction.validated   → consumed by: B-08 Approval Workflow
approval.level.completed → consumed by: B-08 next level trigger
approval.completed      → consumed by: B-09 Tally Sync (Phase 6)
approval.stale          → consumed by: Notification Dispatcher (reminders)
```

---

## 8. Testing Requirements

| Test Type | Scope | Target |
|---|---|---|
| B-04 vendor matching | 50 invoices (known + unknown vendors) | ≥ 90% correct match; correct new-vendor creation |
| B-05 ledger mapping | 50 transactions across all expense categories | ≥ 85% first-suggestion correct |
| B-05 learning feedback | Submit 10 corrections; re-run same vendors | Corrected patterns should appear in next-run suggestions |
| B-07 duplicate detection | 20 exact + 10 fuzzy duplicates | 100% exact detection; ≥ 90% fuzzy detection |
| B-08 approval chain | Walk all 4-level paths (approve, reject, return) | State machine correct; notifications sent at each level |
| B-08 self-approval block | Attempt to approve own submission | OPA rejects request |
| B-08 stale reminders | Simulate 24h delay | Reminder notification sent |
| E-04 PDF Arabic | Generate 5 report types with AR locale | RTL layout correct; Arabic fonts render; validated by Arabic reviewer |
| OPA self-approval | Finance head approves own txn | Blocked with clear error |

---

## 9. Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Ledger mapping < 85% accuracy for new business categories | High | Expand few-shot training set; allow accountants to seed correction history before go-live |
| Vendor deduplication creating unintended merges | Medium | Conservative merge threshold (> 95% similarity); all new-vendor creations require accountant confirmation |
| Approval chain bottleneck (finance_head busy) | Medium | Configurable delegation: finance_head can designate backup approver |
| Arabic PDF rendering font substitution | Medium | Embed fonts (Cairo + Noto Sans Arabic) in WeasyPrint config; test all PDF variants before release |

---

## 10. Phase Exit Criteria

- [ ] Full B-04 → B-05 → B-06 → B-07 → B-08 pipeline processing documents end-to-end
- [ ] B-07 duplicate detection at 100% exact + ≥ 90% fuzzy on test set
- [ ] B-08 approval workflow state machine correct across all paths
- [ ] OPA self-approval block verified
- [ ] Stale approval reminders firing at 24h and 48h
- [ ] Ledger correction feedback stored and confirmed improving B-05 suggestions
- [ ] E-04 Arabic PDFs rendered correctly (signed off by Arabic reviewer)
- [ ] Transaction review and approval UI usable by accountant, manager, and finance head roles
- [ ] Phase gate reviewed and signed off

---

*Phase 05 | Version 1.0 | February 19, 2026*
