// Package middleware provides HTTP middleware for the MediSync API.
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/medisync/medisync/internal/warehouse/models"
)

// UploadValidator validates file uploads.
type UploadValidator struct {
	maxFileSize    int64
	allowedFormats map[string]bool
}

// UploadValidatorConfig holds configuration for the upload validator.
type UploadValidatorConfig struct {
	MaxFileSize    int64
	AllowedFormats []string
}

// NewUploadValidator creates a new upload validator.
func NewUploadValidator(cfg UploadValidatorConfig) *UploadValidator {
	if cfg.MaxFileSize == 0 {
		cfg.MaxFileSize = 25 << 20 // 25MB default
	}

	allowedFormats := make(map[string]bool)
	defaultFormats := []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
		"image/tiff",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/csv",
	}

	formats := cfg.AllowedFormats
	if len(formats) == 0 {
		formats = defaultFormats
	}

	for _, format := range formats {
		allowedFormats[format] = true
	}

	return &UploadValidator{
		maxFileSize:    cfg.MaxFileSize,
		allowedFormats: allowedFormats,
	}
}

// ValidateFile validates an uploaded file.
func (v *UploadValidator) ValidateFile(contentType string, size int64) error {
	// Check file size
	if size > v.maxFileSize {
		return errors.New("file size exceeds maximum allowed size")
	}

	// Normalize content type
	contentType = strings.Split(contentType, ";")[0]
	contentType = strings.TrimSpace(contentType)

	// Check format
	if !v.allowedFormats[contentType] {
		return errors.New("file format not supported")
	}

	return nil
}

// ValidateUploadMiddleware returns a middleware that validates file uploads.
func (v *UploadValidator) ValidateUploadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			http.Error(w, `{"error":{"message":"Content-Type must be multipart/form-data","code":"Bad Request"}}`, http.StatusBadRequest)
			return
		}

		// Parse multipart form to check size
		err := r.ParseMultipartForm(v.maxFileSize)
		if err != nil {
			http.Error(w, `{"error":{"message":"Request too large or invalid","code":"Bad Request"}}`, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetFormatFromMIME returns the FileFormat from MIME type.
func GetFormatFromMIME(mime string) (models.FileFormat, error) {
	mime = strings.Split(mime, ";")[0]
	mime = strings.TrimSpace(mime)

	switch mime {
	case "application/pdf":
		return models.FileFormatPDF, nil
	case "image/jpeg", "image/jpg":
		return models.FileFormatJPEG, nil
	case "image/png":
		return models.FileFormatPNG, nil
	case "image/tiff", "image/tif":
		return models.FileFormatTIFF, nil
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return models.FileFormatXLSX, nil
	case "text/csv":
		return models.FileFormatCSV, nil
	default:
		return "", errors.New("unsupported MIME type")
	}
}
