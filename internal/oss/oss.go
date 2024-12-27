package oss

type ObjectStorageService interface {
	HealthCheck() error
}
