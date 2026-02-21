// Package storage provides object storage implementations for MediSync.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ObjectStorage provides S3-compatible object storage operations.
type ObjectStorage struct {
	client *minio.Client
	bucket string
}

// ObjectStorageConfig holds configuration for the ObjectStorage.
type ObjectStorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// NewObjectStorage creates a new ObjectStorage client.
func NewObjectStorage(cfg ObjectStorageConfig) (*ObjectStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: failed to create client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("storage: failed to create bucket: %w", err)
		}
	}

	return &ObjectStorage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Upload uploads a file from a reader.
func (s *ObjectStorage) Upload(ctx context.Context, path string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: getContentType(path),
	})
	if err != nil {
		return fmt.Errorf("storage: failed to upload object: %w", err)
	}
	return nil
}

// UploadBytes uploads a file from a byte slice.
func (s *ObjectStorage) UploadBytes(ctx context.Context, path string, data []byte) error {
	reader := bytes.NewReader(data)
	return s.Upload(ctx, path, reader, int64(len(data)))
}

// Download downloads a file to a writer.
func (s *ObjectStorage) Download(ctx context.Context, path string, writer io.Writer) error {
	object, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("storage: failed to get object: %w", err)
	}
	defer object.Close()

	_, err = io.Copy(writer, object)
	if err != nil {
		return fmt.Errorf("storage: failed to download object: %w", err)
	}
	return nil
}

// DownloadBytes downloads a file as a byte slice.
func (s *ObjectStorage) DownloadBytes(ctx context.Context, path string) ([]byte, error) {
	var buf bytes.Buffer
	if err := s.Download(ctx, path, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Delete removes a file.
func (s *ObjectStorage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("storage: failed to delete object: %w", err)
	}
	return nil
}

// GetPresignedURL generates a presigned URL for direct access.
func (s *ObjectStorage) GetPresignedURL(path string, expiry time.Duration) string {
	url, err := s.client.PresignedGetObject(context.Background(), s.bucket, path, expiry, nil)
	if err != nil {
		return ""
	}
	return url.String()
}

// GetDocumentPath generates a storage path for a document.
func (s *ObjectStorage) GetDocumentPath(tenantID, documentID uuid.UUID, filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("documents/%s/%s%s", tenantID, documentID, ext)
}

// GetUploadPath generates a temporary upload path.
func (s *ObjectStorage) GetUploadPath(tenantID, uploadID uuid.UUID, filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("uploads/%s/%s%s", tenantID, uploadID, ext)
}

// Stat returns information about an object.
func (s *ObjectStorage) Stat(ctx context.Context, path string) (ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return ObjectInfo{}, fmt.Errorf("storage: failed to stat object: %w", err)
	}
	return ObjectInfo{
		Key:          info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
	}, nil
}

// ObjectInfo contains information about a stored object.
type ObjectInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
}

// ListObjects lists objects with a given prefix.
func (s *ObjectStorage) ListObjects(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo
	ch := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for info := range ch {
		if info.Err != nil {
			return nil, fmt.Errorf("storage: failed to list objects: %w", info.Err)
		}
		objects = append(objects, ObjectInfo{
			Key:          info.Key,
			Size:         info.Size,
			ContentType:  info.ContentType,
			LastModified: info.LastModified,
		})
	}

	return objects, nil
}

// getContentType returns the content type based on file extension.
func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".tiff", ".tif":
		return "image/tiff"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".csv":
		return "text/csv"
	default:
		return "application/octet-stream"
	}
}
