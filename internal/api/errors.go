// Package api provides comprehensive error handling with localization for MediSync.
//
// This file implements the APIError type for structured error responses with
// full internationalization support (English and Arabic). All API errors follow
// a consistent format with error codes, localized messages, and optional retry hints.
//
// Error codes follow the pattern: MODULE_ERROR_TYPE (e.g., "QUERY_INVALID_REQUEST")
//
// Usage:
//
//	err := api.NewAPIError("INVALID_REQUEST", "en")
//	err.WriteJSON(w, http.StatusBadRequest)
//
// Or with additional context:
//
//	err := api.NewAPIErrorWithDetails("QUERY_TIMEOUT", "en", "Query took longer than 30 seconds")
//	err.WithRetryAfter(60).WriteJSON(w, http.StatusGatewayTimeout)
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error codes organized by module/functionality.
const (
	// General errors (0000-0999)
	ErrInvalidRequest     = "INVALID_REQUEST"
	ErrUnauthorized       = "UNAUTHORIZED"
	ErrForbidden          = "FORBIDDEN"
	ErrNotFound           = "NOT_FOUND"
	ErrMethodNotAllowed   = "METHOD_NOT_ALLOWED"
	ErrConflict           = "CONFLICT"
	ErrRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	ErrInternalError      = "INTERNAL_ERROR"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrGatewayTimeout     = "GATEWAY_TIMEOUT"
	ErrRequestTooLarge    = "REQUEST_TOO_LARGE"
	ErrUnsupportedMedia   = "UNSUPPORTED_MEDIA_TYPE"

	// Query/BI Module errors (1000-1999)
	ErrQueryInvalid       = "QUERY_INVALID"
	ErrQueryTooLong       = "QUERY_TOO_LONG"
	ErrQueryTimeout       = "QUERY_TIMEOUT"
	ErrQueryComplexity    = "QUERY_COMPLEXITY_EXCEEDED"
	ErrSQLGeneration      = "SQL_GENERATION_FAILED"
	ErrSQLExecution       = "SQL_EXECUTION_ERROR"
	ErrSQLPermission      = "SQL_PERMISSION_DENIED"
	ErrNoResults          = "NO_RESULTS_FOUND"
	ErrVisualizationError = "VISUALIZATION_ERROR"

	// OCR/Document Module errors (2000-2999)
	ErrDocumentNotFound  = "DOCUMENT_NOT_FOUND"
	ErrDocumentTooLarge  = "DOCUMENT_TOO_LARGE"
	ErrDocumentFormat    = "DOCUMENT_FORMAT_UNSUPPORTED"
	ErrOCRFailed         = "OCR_PROCESSING_FAILED"
	ErrOCRQualityLow     = "OCR_QUALITY_TOO_LOW"
	ErrHandwritingFailed = "HANDWRITING_RECOGNITION_FAILED"

	// Ledger/Accounting Module errors (3000-3999)
	ErrLedgerNotFound  = "LEDGER_NOT_FOUND"
	ErrLedgerMapping   = "LEDGER_MAPPING_FAILED"
	ErrLedgerAmbiguous = "LEDGER_MAPPING_AMBIGUOUS"
	ErrAccountNotFound = "ACCOUNT_NOT_FOUND"
	ErrJournalEntry    = "JOURNAL_ENTRY_INVALID"

	// Tally Sync Module errors (4000-4999)
	ErrTallyConnection  = "TALLY_CONNECTION_FAILED"
	ErrTallySync        = "TALLY_SYNC_FAILED"
	ErrTallyValidation  = "TALLY_VALIDATION_ERROR"
	ErrApprovalRequired = "APPROVAL_REQUIRED"
	ErrApprovalPending  = "APPROVAL_PENDING"
	ErrApprovalRejected = "APPROVAL_REJECTED"

	// Report Module errors (5000-5999)
	ErrReportNotFound   = "REPORT_NOT_FOUND"
	ErrReportGeneration = "REPORT_GENERATION_FAILED"
	ErrReportExport     = "REPORT_EXPORT_FAILED"
	ErrReportTooLarge   = "REPORT_TOO_LARGE"

	// Analytics Module errors (6000-6999)
	ErrAnalyticsSearch    = "ANALYTICS_SEARCH_FAILED"
	ErrInsightsGeneration = "INSIGHTS_GENERATION_FAILED"
	ErrDataNotFound       = "DATA_NOT_FOUND"

	// i18n Module errors (7000-7999)
	ErrLanguageNotSupported = "LANGUAGE_NOT_SUPPORTED"
	ErrTranslationFailed    = "TRANSLATION_FAILED"
	ErrLocaleInvalid        = "LOCALE_INVALID"

	// Authentication/Authorization errors (8000-8999)
	ErrTokenExpired      = "TOKEN_EXPIRED"
	ErrTokenInvalid      = "TOKEN_INVALID"
	ErrTokenMissing      = "TOKEN_MISSING"
	ErrSessionExpired    = "SESSION_EXPIPED"
	ErrInsufficientScope = "INSUFFICIENT_SCOPE"
	ErrTenantMismatch    = "TENANT_MISMATCH"

	// Validation errors (9000-9999)
	ErrValidationFailed = "VALIDATION_FAILED"
	ErrInvalidUUID      = "INVALID_UUID"
	ErrInvalidDate      = "INVALID_DATE_FORMAT"
	ErrInvalidNumber    = "INVALID_NUMBER_FORMAT"
	ErrInvalidEmail     = "INVALID_EMAIL_FORMAT"
	ErrRequiredField    = "REQUIRED_FIELD_MISSING"
)

// ErrorMessages contains localized error messages for all error codes.
// Structure: map[errorCode]map[locale]message
var ErrorMessages = map[string]map[string]string{
	// General errors
	ErrInvalidRequest: {
		"en": "The request is invalid or malformed. Please check the request format and try again.",
		"ar": "الطلب غير صالح أو مشوه. يرجى التحقق من تنسيق الطلب والمحاولة مرة أخرى.",
	},
	ErrUnauthorized: {
		"en": "Authentication is required to access this resource. Please log in and try again.",
		"ar": "المصادقة مطلوبة للوصول إلى هذا المورد. يرجى تسجيل الدخول والمحاولة مرة أخرى.",
	},
	ErrForbidden: {
		"en": "You do not have permission to perform this action. Contact your administrator for access.",
		"ar": "ليس لديك إذن لتنفيذ هذا الإجراء. اتصل بالمسؤول للحصول على الوصول.",
	},
	ErrNotFound: {
		"en": "The requested resource was not found. It may have been deleted or moved.",
		"ar": "لم يتم العثور على المورد المطلوب. قد تم حذفه أو نقله.",
	},
	ErrMethodNotAllowed: {
		"en": "The HTTP method is not allowed for this endpoint. Check the API documentation.",
		"ar": "طريقة HTTP غير مسموح بها لهذه النقطة النهائية. تحقق من وثائق API.",
	},
	ErrConflict: {
		"en": "The request conflicts with the current state of the resource. Refresh and try again.",
		"ar": "يتعارض الطلب مع الحالة الحالية للمورد. قم بالتحديث والمحاولة مرة أخرى.",
	},
	ErrRateLimitExceeded: {
		"en": "Too many requests. You have exceeded the rate limit. Please wait before trying again.",
		"ar": "طلبات كثيرة جدًا. لقد تجاوزت حد المعدل. يرجى الانتظار قبل المحاولة مرة أخرى.",
	},
	ErrInternalError: {
		"en": "An internal server error occurred. Our team has been notified. Please try again later.",
		"ar": "حدث خطأ داخلي في الخادم. تم إبلاغ فريقنا. يرجى المحاولة مرة أخرى لاحقًا.",
	},
	ErrServiceUnavailable: {
		"en": "The service is temporarily unavailable. Please try again in a few moments.",
		"ar": "الخدمة غير متاحة مؤقتًا. يرجى المحاولة مرة أخرى في لحظات.",
	},
	ErrGatewayTimeout: {
		"en": "The request took too long to process. Please simplify your query or try again later.",
		"ar": "استغرق الطلب وقتًا طويلاً للمعالجة. يرجى تبسيط استعلامك أو المحاولة مرة أخرى لاحقًا.",
	},
	ErrRequestTooLarge: {
		"en": "The request body is too large. Please reduce the size of your request.",
		"ar": "جسم الطلب كبير جدًا. يرجى تقليل حجم طلبك.",
	},
	ErrUnsupportedMedia: {
		"en": "The media type is not supported. Please use application/json.",
		"ar": "نوع الوسائط غير مدعوم. يرجى استخدام application/json.",
	},

	// Query/BI Module errors
	ErrQueryInvalid: {
		"en": "The query could not be understood. Please rephrase your question and try again.",
		"ar": "تعذر فهم الاستعلام. يرجى إعادة صياغة سؤالك والمحاولة مرة أخرى.",
	},
	ErrQueryTooLong: {
		"en": "The query is too long. Please break it down into smaller, more specific questions.",
		"ar": "الاستعلام طويل جدًا. يرجى تقسيمه إلى أسئلة أصغر وأكثر تحديدًا.",
	},
	ErrQueryTimeout: {
		"en": "The query took too long to execute. Please simplify your question or add more filters.",
		"ar": "استغرق الاستعلام وقتًا طويلاً للتنفيذ. يرجى تبسيط سؤالك أو إضافة المزيد من المرشحات.",
	},
	ErrQueryComplexity: {
		"en": "The query is too complex. Please break it down into simpler parts.",
		"ar": "الاستعلام معقد جدًا. يرجى تقسيمه إلى أجزاء أبسط.",
	},
	ErrSQLGeneration: {
		"en": "Could not generate a valid database query from your question. Please rephrase it.",
		"ar": "تعذر إنشاء استعلام قاعدة بيانات صالح من سؤالك. يرجى إعادة صياغته.",
	},
	ErrSQLExecution: {
		"en": "An error occurred while executing the query. Please try again or contact support.",
		"ar": "حدث خطأ أثناء تنفيذ الاستعلام. يرجى المحاولة مرة أخرى أو الاتصال بالدعم.",
	},
	ErrSQLPermission: {
		"en": "This query requires permissions you do not have. Contact your administrator.",
		"ar": "يتطلب هذا الاستعلام أذونات ليس لديك. اتصل بالمسؤول الخاص بك.",
	},
	ErrNoResults: {
		"en": "No results found for your query. Try different search terms or filters.",
		"ar": "لم يتم العثور على نتائج لاستعلامك. جرب مصطلحات بحث أو مرشحات مختلفة.",
	},
	ErrVisualizationError: {
		"en": "Could not create visualization for the results. The data format may be incompatible.",
		"ar": "تعذر إنشاء تصور للنتائج. قد يكون تنسيق البيانات غير متوافق.",
	},

	// OCR/Document Module errors
	ErrDocumentNotFound: {
		"en": "The document was not found. It may have been deleted or the ID is incorrect.",
		"ar": "لم يتم العثور على المستند. قد تم حذفه أو المعرف غير صحيح.",
	},
	ErrDocumentTooLarge: {
		"en": "The document is too large. Maximum file size is 10 MB.",
		"ar": "المستند كبير جدًا. الحد الأقصى لحجم الملف هو 10 ميجابايت.",
	},
	ErrDocumentFormat: {
		"en": "The document format is not supported. Please use PDF, JPG, or PNG.",
		"ar": "تنسيق المستند غير مدعوم. يرجى استخدام PDF أو JPG أو PNG.",
	},
	ErrOCRFailed: {
		"en": "OCR processing failed. The document may be corrupted or unreadable.",
		"ar": "فشلت معالجة OCR. قد يكون المستند تالفًا أو غير قابل للقراءة.",
	},
	ErrOCRQualityLow: {
		"en": "The document quality is too low for accurate OCR. Please upload a clearer image.",
		"ar": "جودة المستند منخفضة جدًا للحصول على OCR دقيق. يرجى تحميل صورة أوضح.",
	},
	ErrHandwritingFailed: {
		"en": "Handwriting recognition failed. Please use a printed or typed document.",
		"ar": "فشل التعرف على الخط اليدوي. يرجى استخدام مستند مطبوع أو مكتوب.",
	},

	// Ledger/Accounting Module errors
	ErrLedgerNotFound: {
		"en": "The ledger account was not found. Check the account code and try again.",
		"ar": "لم يتم العثور على حساب دفتر الأستاذ. تحقق من رمز الحساب وحاول مرة أخرى.",
	},
	ErrLedgerMapping: {
		"en": "Could not determine the appropriate ledger account. Please select manually.",
		"ar": "تعذر تحديد حساب دفتر الأستاذ المناسب. يرجى التحديد يدويًا.",
	},
	ErrLedgerAmbiguous: {
		"en": "Multiple ledger accounts match. Please clarify the transaction type.",
		"ar": "تطابق حسابات دفتر أستاذ متعددة. يرجى توضيح نوع المعاملة.",
	},
	ErrAccountNotFound: {
		"en": "The account was not found in Tally. Verify the account exists.",
		"ar": "لم يتم العثور على الحساب في Tally. تحقق من وجود الحساب.",
	},
	ErrJournalEntry: {
		"en": "The journal entry is invalid. Debits must equal credits.",
		"ar": "إدخال دفتر اليومية غير صالح. يجب أن تساوي المدينون الدائنون.",
	},

	// Tally Sync Module errors
	ErrTallyConnection: {
		"en": "Could not connect to Tally ERP. Ensure Tally is running and accessible.",
		"ar": "تعذر الاتصال بـ Tally ERP. تأكد من تشغيل Tally وإمكانية الوصول إليه.",
	},
	ErrTallySync: {
		"en": "Failed to sync data with Tally. Please try again or check Tally logs.",
		"ar": "فشل مزامنة البيانات مع Tally. يرجى المحاولة مرة أخرى أو التحقق من سجلات Tally.",
	},
	ErrTallyValidation: {
		"en": "Tally rejected the data due to validation errors. Check the data format.",
		"ar": "رفض Tally البيانات بسبب أخطاء التحقق. تحقق من تنسيق البيانات.",
	},
	ErrApprovalRequired: {
		"en": "This action requires approval from an authorized person before it can be completed.",
		"ar": "يتطلب هذا الإجراء موافقة من شخص مخول قبل أن يتم إكماله.",
	},
	ErrApprovalPending: {
		"en": "This request is pending approval. You will be notified once it is reviewed.",
		"ar": "هذا الطلب في انتظار الموافقة. سيتم إعلامك بمجرد مراجعته.",
	},
	ErrApprovalRejected: {
		"en": "This request was rejected. Please review the feedback and resubmit.",
		"ar": "تم رفض هذا الطلب. يرجى مراجعة التعليقات وإعادة الإرسال.",
	},

	// Report Module errors
	ErrReportNotFound: {
		"en": "The report was not found. It may have been deleted or the ID is incorrect.",
		"ar": "لم يتم العثور على التقرير. قد تم حذفه أو المعرف غير صحيح.",
	},
	ErrReportGeneration: {
		"en": "Failed to generate the report. Please check your parameters and try again.",
		"ar": "فشل في إنشاء التقرير. يرجى التحقق من المعلمات والمحاولة مرة أخرى.",
	},
	ErrReportExport: {
		"en": "Failed to export the report. The format may not be supported.",
		"ar": "فشل في تصدير التقرير. قد لا يكون التنسيق مدعومًا.",
	},
	ErrReportTooLarge: {
		"en": "The report is too large to generate. Please add filters to reduce data size.",
		"ar": "التقرير كبير جدًا لإنشائه. يرجى إضافة مرشحات لتقليل حجم البيانات.",
	},

	// Analytics Module errors
	ErrAnalyticsSearch: {
		"en": "The analytics search failed. Please try with different parameters.",
		"ar": "فشل بحث التحليلات. يرجى المحاولة مع معلمات مختلفة.",
	},
	ErrInsightsGeneration: {
		"en": "Could not generate insights from the available data.",
		"ar": "تعذر إنشاء رؤى من البيانات المتاحة.",
	},
	ErrDataNotFound: {
		"en": "No data found for the specified criteria.",
		"ar": "لم يتم العثور على بيانات للمعايير المحددة.",
	},

	// i18n Module errors
	ErrLanguageNotSupported: {
		"en": "The requested language is not supported. Available languages: English, Arabic.",
		"ar": "اللغة المطلوبة غير مدعومة. اللغات المتاحة: الإنجليزية، العربية.",
	},
	ErrTranslationFailed: {
		"en": "Translation service failed. The response may not be in your preferred language.",
		"ar": "فشلت خدمة الترجمة. قد لا تكون الاستجابة بلغتك المفضلة.",
	},
	ErrLocaleInvalid: {
		"en": "The locale code is invalid. Use 'en' for English or 'ar' for Arabic.",
		"ar": "رمز اللغة غير صالح. استخدم 'en' للإنجليزية أو 'ar' للعربية.",
	},

	// Authentication/Authorization errors
	ErrTokenExpired: {
		"en": "Your session has expired. Please log in again to continue.",
		"ar": "انتهت صلاحية جلستك. يرجى تسجيل الدخول مرة أخرى للمتابعة.",
	},
	ErrTokenInvalid: {
		"en": "The authentication token is invalid. Please log in again.",
		"ar": "رمز المصادقة غير صالح. يرجى تسجيل الدخول مرة أخرى.",
	},
	ErrTokenMissing: {
		"en": "Authentication token is missing. Include a valid Bearer token in the header.",
		"ar": "رمز المصادقة مفقود. قم بتضمين رمز Bearer صالح في الرأس.",
	},
	ErrSessionExpired: {
		"en": "Your session has expired due to inactivity. Please log in again.",
		"ar": "انتهت صلاحية جلستك بسبب عدم النشاط. يرجى تسجيل الدخول مرة أخرى.",
	},
	ErrInsufficientScope: {
		"en": "Your account does not have the required permissions for this operation.",
		"ar": "ليس لحسابك الأذونات المطلوبة لهذه العملية.",
	},
	ErrTenantMismatch: {
		"en": "You are not authorized to access resources from this organization.",
		"ar": "أنت غير مخول للوصول إلى موارد من هذه المؤسسة.",
	},

	// Validation errors
	ErrValidationFailed: {
		"en": "Validation failed. Please check the input values and correct any errors.",
		"ar": "فشل التحقق. يرجى التحقق من قيم الإدخال وتصحيح أي أخطاء.",
	},
	ErrInvalidUUID: {
		"en": "The ID format is invalid. It must be a valid UUID.",
		"ar": "تنسيق المعرف غير صالح. يجب أن يكون UUID صالحًا.",
	},
	ErrInvalidDate: {
		"en": "The date format is invalid. Use ISO 8601 format (YYYY-MM-DD).",
		"ar": "تنسيق التاريخ غير صالح. استخدم تنسيق ISO 8601 (YYYY-MM-DD).",
	},
	ErrInvalidNumber: {
		"en": "The number format is invalid. Please provide a valid numeric value.",
		"ar": "تنسيق الرقم غير صالح. يرجى تقديم قيمة رقمية صالحة.",
	},
	ErrInvalidEmail: {
		"en": "The email address format is invalid. Please provide a valid email.",
		"ar": "تنسيق عنوان البريد الإلكتروني غير صالح. يرجى تقديم بريد إلكتروني صالح.",
	},
	ErrRequiredField: {
		"en": "A required field is missing. Please fill in all mandatory fields.",
		"ar": "حقل مطلوب مفقود. يرجى ملء جميع الحقول الإلزامية.",
	},
}

// APIError represents a structured API error with localization support.
type APIError struct {
	// Code is the machine-readable error code.
	Code string `json:"code"`

	// Message is the human-readable error message (localized).
	Message string `json:"message"`

	// Locale is the language code used for the message (not included in JSON).
	Locale string `json:"-"`

	// Details contains additional context about the error.
	Details string `json:"details,omitempty"`

	// RetryAfter indicates seconds to wait before retrying (for rate limits).
	RetryAfter int `json:"retry_after,omitempty"`

	// Field indicates which field caused the error (for validation errors).
	Field string `json:"field,omitempty"`

	// Suggestion provides a hint for resolving the error.
	Suggestion string `json:"suggestion,omitempty"`

	// HTTPStatus is the HTTP status code to use (not included in JSON).
	HTTPStatus int `json:"-"`
}

// NewAPIError creates a new API error with the localized message.
func NewAPIError(code string, locale string) *APIError {
	if locale == "" {
		locale = "en"
	}

	messages, ok := ErrorMessages[code]
	if !ok {
		// Fall back to internal error for unknown codes
		messages = ErrorMessages[ErrInternalError]
	}

	message, ok := messages[locale]
	if !ok {
		// Fall back to English if locale not found
		message = messages["en"]
	}

	return &APIError{
		Code:       code,
		Message:    message,
		Locale:     locale,
		HTTPStatus: getHTTPStatus(code),
	}
}

// NewAPIErrorWithDetails creates a new API error with additional details.
func NewAPIErrorWithDetails(code string, locale string, details string) *APIError {
	err := NewAPIError(code, locale)
	err.Details = details
	return err
}

// WithRetryAfter sets the retry-after hint for rate-limited errors.
func (e *APIError) WithRetryAfter(seconds int) *APIError {
	e.RetryAfter = seconds
	return e
}

// WithField sets the field name for validation errors.
func (e *APIError) WithField(field string) *APIError {
	e.Field = field
	return e
}

// WithSuggestion sets a suggestion for resolving the error.
func (e *APIError) WithSuggestion(suggestion string) *APIError {
	e.Suggestion = suggestion
	return e
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WriteJSON writes the error as a JSON response.
func (e *APIError) WriteJSON(w http.ResponseWriter, status int) {
	if status == 0 {
		status = e.HTTPStatus
		if status == 0 {
			status = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Set Retry-After header for rate limiting
	if e.RetryAfter > 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", e.RetryAfter))
	}

	w.WriteHeader(status)

	// Encode error response
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		},
	}

	if e.Details != "" {
		response["error"].(map[string]interface{})["details"] = e.Details
	}

	if e.RetryAfter > 0 {
		response["error"].(map[string]interface{})["retry_after"] = e.RetryAfter
	}

	if e.Field != "" {
		response["error"].(map[string]interface{})["field"] = e.Field
	}

	if e.Suggestion != "" {
		response["error"].(map[string]interface{})["suggestion"] = e.Suggestion
	}

	json.NewEncoder(w).Encode(response)
}

// WriteJSONWithStatus writes the error with a specific HTTP status.
func (e *APIError) WriteJSONWithStatus(w http.ResponseWriter, status int) {
	e.WriteJSON(w, status)
}

// getHTTPStatus maps error codes to HTTP status codes.
func getHTTPStatus(code string) int {
	switch code {
	case ErrInvalidRequest, ErrQueryInvalid, ErrQueryTooLong, ErrQueryComplexity,
		ErrDocumentTooLarge, ErrDocumentFormat, ErrValidationFailed,
		ErrInvalidUUID, ErrInvalidDate, ErrInvalidNumber, ErrInvalidEmail,
		ErrRequiredField, ErrLocaleInvalid, ErrRequestTooLarge, ErrUnsupportedMedia:
		return http.StatusBadRequest

	case ErrUnauthorized, ErrTokenExpired, ErrTokenInvalid, ErrTokenMissing,
		ErrSessionExpired:
		return http.StatusUnauthorized

	case ErrForbidden, ErrSQLPermission, ErrInsufficientScope, ErrTenantMismatch:
		return http.StatusForbidden

	case ErrNotFound, ErrDocumentNotFound, ErrLedgerNotFound, ErrAccountNotFound,
		ErrReportNotFound, ErrDataNotFound:
		return http.StatusNotFound

	case ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed

	case ErrConflict, ErrApprovalPending:
		return http.StatusConflict

	case ErrRateLimitExceeded:
		return http.StatusTooManyRequests

	case ErrQueryTimeout, ErrGatewayTimeout:
		return http.StatusGatewayTimeout

	case ErrInternalError, ErrSQLGeneration, ErrSQLExecution, ErrOCRFailed,
		ErrLedgerMapping, ErrTallySync, ErrReportGeneration, ErrAnalyticsSearch,
		ErrTranslationFailed:
		return http.StatusInternalServerError

	case ErrServiceUnavailable, ErrTallyConnection:
		return http.StatusServiceUnavailable

	case ErrApprovalRequired:
		return http.StatusAccepted

	default:
		return http.StatusInternalServerError
	}
}

// ErrorResponse writes a generic error response.
func ErrorResponse(w http.ResponseWriter, r *http.Request, code string, status int) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(code, locale)
	err.WriteJSON(w, status)
}

// ErrorResponseWithDetails writes an error response with additional details.
func ErrorResponseWithDetails(w http.ResponseWriter, r *http.Request, code string, status int, details string) {
	locale := getLocaleFromRequest(r)
	err := NewAPIErrorWithDetails(code, locale, details)
	err.WriteJSON(w, status)
}

// ValidationErrorResponse writes a validation error response.
func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, field, message string) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrValidationFailed, locale)
	err.WithField(field).Details = message
	err.WriteJSON(w, http.StatusBadRequest)
}

// NotFoundResponse writes a 404 not found response.
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrNotFound, locale)
	err.WriteJSON(w, http.StatusNotFound)
}

// UnauthorizedResponse writes a 401 unauthorized response.
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrUnauthorized, locale)
	w.Header().Set("WWW-Authenticate", `Bearer realm="medisync"`)
	err.WriteJSON(w, http.StatusUnauthorized)
}

// ForbiddenResponse writes a 403 forbidden response.
func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrForbidden, locale)
	err.WriteJSON(w, http.StatusForbidden)
}

// RateLimitResponse writes a 429 rate limit exceeded response.
func RateLimitResponse(w http.ResponseWriter, r *http.Request, retryAfter int) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrRateLimitExceeded, locale)
	err.WithRetryAfter(retryAfter)
	err.WriteJSON(w, http.StatusTooManyRequests)
}

// TimeoutResponse writes a 504 gateway timeout response.
func TimeoutResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrGatewayTimeout, locale)
	err.WriteJSON(w, http.StatusGatewayTimeout)
}

// InternalErrorResponse writes a 500 internal server error response.
func InternalErrorResponse(w http.ResponseWriter, r *http.Request) {
	locale := getLocaleFromRequest(r)
	err := NewAPIError(ErrInternalError, locale)
	err.WriteJSON(w, http.StatusInternalServerError)
}

// getLocaleFromRequest extracts the locale from the request context or headers.
func getLocaleFromRequest(r *http.Request) string {
	// First check context (set by locale middleware)
	if ctx := r.Context(); ctx != nil {
		if locale, ok := ctx.Value(contextKey("locale")).(string); ok && locale != "" {
			return locale
		}
	}

	// Fall back to Accept-Language header
	acceptLanguage := r.Header.Get("Accept-Language")
	if acceptLanguage != "" {
		// Simple extraction - just get first language
		for i, c := range acceptLanguage {
			if c == ',' || c == ';' {
				acceptLanguage = acceptLanguage[:i]
				break
			}
		}
		if acceptLanguage == "ar" || acceptLanguage == "ar-SA" {
			return "ar"
		}
	}

	// Check URL parameter
	if lang := r.URL.Query().Get("lang"); lang != "" {
		if lang == "ar" {
			return "ar"
		}
	}

	return "en"
}

// contextKey is a type for context keys.
type contextKey string

// MultiError represents multiple validation errors.
type MultiError struct {
	Errors []*APIError `json:"errors"`
}

// NewMultiError creates a new multi-error container.
func NewMultiError() *MultiError {
	return &MultiError{
		Errors: make([]*APIError, 0),
	}
}

// Add adds an error to the multi-error.
func (m *MultiError) Add(err *APIError) {
	m.Errors = append(m.Errors, err)
}

// HasErrors returns true if there are any errors.
func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

// WriteJSON writes all errors as a JSON response.
func (m *MultiError) WriteJSON(w http.ResponseWriter, status int) {
	if status == 0 {
		status = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    ErrValidationFailed,
			"message": "Multiple validation errors occurred",
			"errors":  m.Errors,
		},
	}

	json.NewEncoder(w).Encode(response)
}
