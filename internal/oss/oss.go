package oss

import (
	"context"
	"time"
)

type ObjectStorageService interface {
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
	PresignedPutObject(ctx context.Context, bucketName, objectName string) (string, error)
	GetObject(ctx context.Context, bucketName, objectName string) (string, error)
	PutObject(ctx context.Context, bucketName, objectName, filePath string) error
	ObjectExists(ctx context.Context, bucketName, objectName string) error
	HealthCheck() error
	RemoveObject(ctx context.Context, bucketName string, objectName string) error
}
