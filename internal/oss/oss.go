package oss

import "context"

type ObjectStorageService interface {
	GetPresignedUrl(ctx context.Context, bucketName, objectName string) (string, error)
	GetObject(ctx context.Context, bucketName, objectName string) (string, error)
	PutObject(ctx context.Context, bucketName, objectName, filePath string) error
	ObjectExists(ctx context.Context, bucketName, objectName string) error
	HealthCheck() error
}
