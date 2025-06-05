package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	domainservice "github.com/MizukiShigi/cms-go/internal/domain/service"
	"github.com/google/uuid"
)

type storageService struct {
	client *storage.Client
}

func NewStorageService(client *storage.Client) *storageService {
	return &storageService{client: client}
}

func (s *storageService) UploadImage(ctx context.Context, bucketName string, fileName valueobject.ImageFilename, data io.Reader) (domainservice.UploadResult, error) {
	bucket := s.client.Bucket(bucketName)
	ext := strings.ToLower(filepath.Ext(fileName.String()))
	path := generateStoredFilename(ext)
	obj := bucket.Object(path)

	writer := obj.NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = detectImageContentType(path)

	if _, err := io.Copy(writer, data); err != nil {
		return domainservice.UploadResult{}, fmt.Errorf("failed to upload image to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return domainservice.UploadResult{}, fmt.Errorf("failed to complete image upload: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, path)
	return domainservice.UploadResult{
		StoredFilename: filepath.Base(path),
		URL:            publicURL,
	}, nil
}

func detectImageContentType(ext string) string {
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

func generateStoredFilename(ext string) string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	uuidStr := uuid.New().String()

	return fmt.Sprintf("images/%s/%s/%s/%s%s", year, month, day, uuidStr, ext)
}
