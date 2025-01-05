package oss

type ObjectStorageService interface {
	GetPresignedUrl(objectName string) (string, error)
	ObjectExists(objectName string) (bool, error)
	HealthCheck() error
}
