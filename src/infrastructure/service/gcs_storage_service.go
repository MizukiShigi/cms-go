package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type storageService struct {
	client *storage.Client
}

func NewStorageService(client *storage.Client) *storageService {
	return &storageService{client: client}
}

func (s *storageService) Upload(ctx context.Context, bucketName string, path string, data io.Reader) (string, error) {
	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(path)

	writer := obj.NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = s.detectImageContentType(path)

	if _, err := io.Copy(writer, data); err != nil {
		return "", fmt.Errorf("failed to upload image to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to complete image upload: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, path)
	return publicURL, nil
}

func (s *storageService) detectImageContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			return contentType
		}
		return "image/jpeg"
	}
}
