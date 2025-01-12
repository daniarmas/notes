package oss

import "context"

type ObjectStorageService interface {
	GetPresignedUrl(objectName string) (string, error)
	GetObject(ctx context.Context, bucketName, objectName string) (string, error)
	ObjectExists(objectName string) error
	HealthCheck() error
}
