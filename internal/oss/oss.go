package oss

type ObjectStorageService interface {
	GetPresignedUrl(objectName string) (string, error)
	HealthCheck() error
}
