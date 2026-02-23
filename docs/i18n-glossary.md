# AnySync — Bilingual Term Glossary (English / Arabic)

**Version:** 1.0 | **Created:** February 19, 2026  
**Status:** Living Document — Update via E-07 Bilingual Glossary Sync Agent  
**Governance:** Changes require review by Medical Advisor + Finance Advisor (HITL gate E-07)  
**Used by:** Agent A-04 (Domain Terminology Normaliser), Agent B-02 (OCR post-processor), AI system prompts, human translator style guide

---

## How to Use This Glossary

- **English** terms are the canonical source of truth for SQL identifiers and data warehouse column names.
- **Arabic** terms are the display-layer equivalents used in UI, chat responses, and reports.
- **Notes** column captures domain nuances, disambiguation guidance, or multiple acceptable Arabic equivalents.
- Terms marked **🔒 Locked** have been reviewed and approved — do not change without E-07 governance workflow.
- Terms marked **⚠️ Review** are draft translations pending expert sign-off.

---

## 1. Healthcare & Clinical Terms

| English | Arabic | Transliteration | Context | Notes |
|---------|--------|----------------|---------|-------|
| Patient | مريض / مرضى (pl.) | Mareed / Marda | HIMS, chat | Singular: مريض; plural: مرضى |
| Patient Footfall | زيارات المرضى | Ziyarat al-Marda | BI Dashboard, KPI | Lit: "patient visits" — do NOT translate as "traffic" (حركة) |
| Appointment | موعد / مواعيد (pl.) | Maw'id / Mawa'eed | HIMS | Medical appointment context |
| Clinic | عيادة / عيادات (pl.) | 'Iyada | Facility reference | |
| Pharmacy | صيدلية / صيدليات (pl.) | Saydalia | Facility reference | |
| Doctor / Physician | طبيب / أطباء (pl.) | Tabeeb | HIMS | |
| Patient Demographics | بيانات المرضى الديموغرافية | — | HIMS | Acceptable short form: بيانات المرضى |
| Diagnosis | تشخيص | Tashkhis | Clinical | |
| Prescription | وصفة طبية | Wasfa Tibbiya | Clinical, Pharmacy | |
| Drug / Medication | دواء / أدوية (pl.) | Dawa' | Pharmacy | Use دواء for individual drug; أدوية for plural/list |
| Drug Dispensation | صرف الأدوية | Sarf al-Adwiya | HIMS, Pharmacy | |
| Expiry Date | تاريخ الانتهاء / تاريخ الصلاحية | — | Pharmacy | Both forms acceptable; تاريخ الانتهاء preferred for system labels |
| Low Stock | مخزون منخفض | Makhzoun Munkhafid | Pharmacy alerts | |
| Stock Out | نفاد المخزون | Nafad al-Makhzoun | Pharmacy alerts | |
| Patient Satisfaction | رضا المرضى | Rida al-Marda | KPI | |
| No-Show | غياب بدون إشعار | — | Appointments | Lit: "absence without notice" — common term in Arabic healthcare |
| Bill / Invoice (Medical) | فاتورة | Fatura | Billing | Same word for bill and invoice |
| Billing | الفوترة | al-Fawtura | HIMS Module | |
| Consultation Fee | رسوم الاستشارة | Rusum al-Istishara | Financial, HIMS | |
| Inpatient | مريض داخلي | — | Clinical | |
| Outpatient | مريض خارجي | — | Clinical | |

---

## 2. Accounting & Finance Terms

| English | Arabic | Transliteration | Context | Notes |
|---------|--------|----------------|---------|-------|
| Revenue | إيرادات | Iyadat | P&L, KPIs | 🔒 Locked — do NOT use مدخولات (informal) |
| Income | دخل | Dukhl | Financial statements | |
| Expense | مصروف / مصاريف (pl.) | Masroof | P&L | |
| Profit | ربح / أرباح (pl.) | Ribh | P&L | |
| Loss | خسارة / خسائر (pl.) | Khasara | P&L | |
| Gross Profit | الربح الإجمالي | al-Ribh al-Ijmali | P&L | |
| Net Profit | صافي الربح | Safi al-Ribh | P&L | |
| Profit Margin | هامش الربح | Hamish al-Ribh | KPI | |
| Gross Profit Margin | هامش الربح الإجمالي | — | KPI, Reports | |
| Cost of Goods Sold (COGS) | تكلفة البضاعة المباعة | — | P&L | |
| Operating Expenses | المصاريف التشغيلية | — | P&L | |
| Receivables / Accounts Receivable | الذمم المدينة / المدينون | al-Dhimam al-Madina | Balance Sheet | Both terms acceptable; الذمم المدينة preferred in formal reports |
| Outstanding Receivables | الذمم المدينة المستحقة | — | Reports, Alerts | |
| Payables / Accounts Payable | الذمم الدائنة / الدائنون | al-Dhimam al-Da'ina | Balance Sheet | |
| Outstanding Invoices | الفواتير المستحقة | al-Fawatir al-Mustahiqqa | AI Accountant | |
| Invoice | فاتورة / فواتير (pl.) | Fatura | AI Accountant | |
| Bill (vendor) | فاتورة مورد | Fatura Murid | AI Accountant | |
| Receipt | إيصال / إيصالات (pl.) | Eesal | Accounting | |
| Ledger | دفتر الأستاذ | Daftar al-Ustad | Tally / Accounting | 🔒 Locked |
| General Ledger (GL) | دفتر الأستاذ العام | — | Tally | |
| Ledger Mapping | تعيين دفتر الأستاذ | — | AI Accountant | |
| Chart of Accounts | دليل الحسابات | — | Tally | |
| Journal Entry | قيد يومية | Qayd Yawmiya | Accounting | |
| Voucher | قسيمة / قسائم (pl.) | Qasima | Tally | |
| Bank Reconciliation | مطابقة الحسابات البنكية | — | AI Accountant | 🔒 Locked |
| Outstanding Payments | المدفوعات المعلقة | — | Reconciliation | |
| Outstanding Receipts | الإيصالات المعلقة | — | Reconciliation | |
| Balance Sheet | الميزانية العمومية | al-Mizaniya al-'Umumiya | Reports | |
| Profit & Loss Statement (P&L) | بيان الأرباح والخسائر | — | Reports | |
| Cash Flow Statement | بيان التدفق النقدي | — | Reports | |
| Trial Balance | ميزان المراجعة | Mizan al-Muraja'a | Reports | |
| Tax | ضريبة / ضرائب (pl.) | Dariba | Compliance | |
| VAT | ضريبة القيمة المضافة | — | Compliance | Common abbreviation: ض.ق.م |
| GST | ضريبة السلع والخدمات | — | Compliance | Context: India/Australia |
| Days Sales Outstanding (DSO) | أيام المبيعات المعلقة | — | KPI | ⚠️ Review — some use متوسط أيام التحصيل |
| Days Payable Outstanding (DPO) | أيام الذمم الدائنة المعلقة | — | KPI | ⚠️ Review |
| Cash Flow | التدفق النقدي | al-Tadfuq al-Naqdi | Analytics, Forecasting | |
| Cash Position | المركز النقدي | al-Markaz al-Naqdi | Dashboard | |
| Budget | الميزانية التقديرية / الموازنة | — | Reports | الموازنة preferred for budget plans |
| Budget vs. Actual | الميزانية مقابل الفعلي | — | Reports | |
| Variance | الانحراف / الفارق | al-Inkhiraf | Reports | |
| Overhead | التكاليف العامة | al-Takalif al-'Amma | Cost Accounting | |
| Cost Centre | مركز التكلفة / مراكز التكلفة (pl.) | Markaz al-Taklifa | Tally, Reports | |
| Profit Centre | مركز الربح | Markaz al-Ribh | Reports | |
| Contribution Margin | هامش المساهمة | — | Analytics | |
| Depreciation | الإهلاك | al-Ihlak | Fixed Assets | |
| Audit Trail | مسار التدقيق | Masar al-Tadqiq | Compliance | 🔒 Locked |
| Compliance | الامتثال | al-Imtithal | General | |
| Reconciliation | المطابقة / التسوية | al-Mutabaqa | Accounting | |

---

## 3. Inventory & Supply Chain Terms

| English | Arabic | Transliteration | Context | Notes |
|---------|--------|----------------|---------|-------|
| Inventory | المخزون | al-Makhzoun | Pharmacy, Reports | 🔒 Locked |
| Inventory Aging | تقادم المخزون | Taqadum al-Makhzoun | Reports | |
| Slow-Moving Stock | مخزون بطيء الحركة | — | Reports, Alerts | |
| Obsolete Stock | مخزون متقادم / راكد | — | Reports | |
| Stock Turnover | معدل دوران المخزون | — | KPI | |
| Reorder Point | نقطة إعادة الطلب | — | Alerts | |
| Reorder Quantity | كمية إعادة الطلب | — | Recommendations | |
| Supplier | مورد / موردون (pl.) | Murid | AI Accountant | |
| Vendor | مورد | Murid | AI Accountant | Same term as Supplier in Arabic |
| Purchase Order | أمر شراء | Amr Shiraa | Procurement | |
| Lead Time | وقت التسليم | Waqt al-Tasleem | Procurement | |

---

## 4. Business Intelligence & Analytics Terms

| English | Arabic | Transliteration | Context | Notes |
|---------|--------|----------------|---------|-------|
| Dashboard | لوحة التحكم | Lawhat al-Tahakum | UI / Navigation | 🔒 Locked |
| Report | تقرير / تقارير (pl.) | Taqrir | Reports Module | |
| Chart | مخطط / مخططات (pl.) | Mukhatat | Visualization | |
| Table | جدول / جداول (pl.) | Jadwal | Visualization | |
| Bar Chart | مخطط شريطي | — | Visualization | |
| Line Chart | مخطط خطي | — | Visualization | |
| Pie Chart | مخطط دائري | — | Visualization | |
| KPI (Key Performance Indicator) | مؤشر الأداء الرئيسي | — | Dashboard | Common abbreviation: م.أ.ر |
| Trend | اتجاه / اتجاهات (pl.) | Itijah | Analytics | |
| Forecast | توقع / توقعات (pl.) | Tawaqo' | Analytics | |
| Anomaly | شذوذ | Shudhudh | Analytics, Alerts | |
| Insight | رؤية تحليلية | Ru'ya Tahliliya | AI responses | Use رؤية for insights; NOT نظرة (which is casual glance) |
| Recommendation | توصية / توصيات (pl.) | Tawsiya | AI responses | |
| Confidence Score | درجة الثقة | Darajat al-Thiqa | AI responses | |
| Drill-Down | التعمق / استعراض التفاصيل | — | BI Dashboard | Use "استعراض التفاصيل" in UI label context |
| Filter | تصفية | Tasfiya | UI | |
| Export | تصدير | Tasdeer | UI | |
| Download | تحميل | Tahmil | UI | |
| Pin to Dashboard | تثبيت في لوحة التحكم | — | UI | |
| Query | استعلام / استفسار | Isti'lam | AI chat | استعلام for technical context; استفسار for conversational |
| Time Period | الفترة الزمنية | al-Fatra al-Zamaniya | Query context | |
| Year-Over-Year (YoY) | مقارنة سنوية | — | Analytics | |
| Month-Over-Month (MoM) | مقارنة شهرية | — | Analytics | |
| Year-to-Date (YTD) | من بداية السنة | — | Date filters | |
| Scheduled Report | تقرير مجدول | — | Reports | |
| Alert / Notification | تنبيه / إشعار | Tanbih / Ish'ar | Alerts | تنبيه for urgent alerts; إشعار for general notifications |

---

## 5. System & UI Terms

| English | Arabic | Transliteration | Context | Notes |
|---------|--------|----------------|---------|-------|
| Settings | الإعدادات | al-I'dadat | Navigation | 🔒 Locked |
| Profile | الملف الشخصي | al-Malaf al-Shakhsi | Navigation | |
| Language | اللغة | al-Lugha | Settings | |
| Display Language | لغة العرض | Lughat al-'Ard | Settings | |
| Save | حفظ | Hifz | Action buttons | 🔒 Locked |
| Cancel | إلغاء | Ilgha' | Action buttons | 🔒 Locked |
| Confirm | تأكيد | Ta'kid | Action buttons | |
| Search | بحث | Baht | UI | |
| Loading | جارٍ التحميل | — | Status | |
| Error | خطأ / أخطاء (pl.) | Khata' | Status | |
| Success | نجاح | Najah | Status | |
| Warning | تحذير / تحذيرات (pl.) | Tahdeer | Status | |
| Sync | مزامنة | Muzawana | AI Accountant | |
| Sync Now | مزامنة الآن | — | Action button | 🔒 Locked |
| Connected | متصل | Muttasil | Status | |
| Disconnected | غير متصل | — | Status | |
| Upload | رفع | Raf' | File actions | |
| Approve | اعتماد / موافقة | I'timad | Approval workflow | اعتماد for formal approval; موافقة for confirmation |
| Reject | رفض | Rafd | Approval workflow | |
| Pending | قيد الانتظار / معلق | — | Status | قيد الانتظار for general; معلق for items on hold |

---

## 6. AI-Specific Phrases (Chat Responses)

| English | Arabic | Context |
|---------|--------|---------|
| "I can only answer questions about your business data." | "يمكنني فقط الإجابة عن أسئلة تتعلق ببيانات عملك." | Off-topic deflection |
| "Your question is ambiguous. Please clarify the time period." | "سؤالك غير واضح. يرجى توضيح الفترة الزمنية المقصودة." | Clarification request |
| "Low confidence — please verify this result." | "ثقة منخفضة — يرجى التحقق من هذه النتيجة." | Confidence warning |
| "Generating your report..." | "جارٍ إنشاء تقريرك..." | Loading state |
| "Here is the data you requested:" | "إليك البيانات المطلوبة:" | Response prefix |
| "No data found for this query." | "لم يتم العثور على بيانات لهذا الاستعلام." | Empty state |
| "Sync completed successfully." | "اكتملت المزامنة بنجاح." | Tally sync status |
| "Sync failed. Please check your connection." | "فشلت المزامنة. يرجى التحقق من الاتصال." | Tally sync error |
| "Based on your data, I recommend:" | "بناءً على بياناتك، أوصي بما يلي:" | Recommendations |
| "This trend shows..." | "يُظهر هذا الاتجاه..." | Insight narrative |

---

## Governance & Change Log

| Version | Date | Changes | Reviewer |
|---------|------|---------|----------|
| 1.0 | Feb 19, 2026 | Initial glossary — seeded from PRD healthcare and accounting terminology | Architecture Team |

> To propose changes: open a PR modifying this file. E-07 Bilingual Glossary Sync Agent will flag the change for Medical Advisor + Finance Advisor review before merging.
