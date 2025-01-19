package config

import (
	"context"
	"os"
	"strconv"

	"github.com/daniarmas/notes/internal/clog"
)

type Configuration struct {
	Environment                   string
	DatabaseUrl                   string
	JwtSecret                     string
	RedisHost                     string
	RedisPort                     string
	RedisPassword                 string
	RedisDb                       int
	ObjectStorageServiceAccessKey string
	ObjectStorageServiceSecretKey string
	ObjectStorageServiceEndpoint  string
	ObjectStorageServiceRegion    string
	ObjectStorageServiceBucket    string
	InK8s                         bool
	DockerImageName               string
}

func LoadServerConfig() *Configuration {
	ctx := context.Background()
	config := Configuration{
		Environment:                   os.Getenv("ENVIRONMENT"),
		DatabaseUrl:                   os.Getenv("DATABASE_URL"),
		RedisHost:                     os.Getenv("REDIS_HOST"),
		RedisPort:                     os.Getenv("REDIS_PORT"),
		RedisPassword:                 os.Getenv("REDIS_PASSWORD"),
		JwtSecret:                     os.Getenv("JWT_SECRET"),
		ObjectStorageServiceAccessKey: os.Getenv("OBJECT_STORAGE_SERVICE_ACCESS_KEY"),
		ObjectStorageServiceSecretKey: os.Getenv("OBJECT_STORAGE_SERVICE_SECRET_KEY"),
		ObjectStorageServiceEndpoint:  os.Getenv("OBJECT_STORAGE_SERVICE_ENDPOINT"),
		ObjectStorageServiceRegion:    os.Getenv("OBJECT_STORAGE_SERVICE_REGION"),
		ObjectStorageServiceBucket:    os.Getenv("OBJECT_STORAGE_SERVICE_BUCKET"),
		InK8s:                         os.Getenv("IN_K8S") == "true",
		DockerImageName:               os.Getenv("DOCKER_IMAGE_NAME"),
	}
	if config.Environment == "" {
		config.Environment = "development"
	}
	if config.DockerImageName == "" {
		config.DockerImageName = "ghcr.io/daniarmas/notes"
	}
	if config.ObjectStorageServiceAccessKey == "" {
		clog.Warn(ctx, "OBJECT_STORAGE_SERVICE_ACCESS_KEY enviroment variable is required", nil)
	}
	if config.ObjectStorageServiceSecretKey == "" {
		clog.Warn(ctx, "OBJECT_STORAGE_SERVICE_SECRET_KEY enviroment variable is required", nil)
	}
	if config.ObjectStorageServiceEndpoint == "" {
		clog.Warn(ctx, "OBJECT_STORAGE_SERVICE_ENDPOINT enviroment variable is required", nil)
	}
	if config.ObjectStorageServiceRegion == "" {
		clog.Warn(ctx, "OBJECT_STORAGE_SERVICE_REGION enviroment variable is required", nil)
	}
	if config.ObjectStorageServiceBucket == "" {
		clog.Warn(ctx, "OBJECT_STORAGE_SERVICE_BUCKET enviroment variable is required", nil)
	}
	if config.JwtSecret == "" {
		clog.Warn(ctx, "JWT_SECRET enviroment variable is required", nil)
	}
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "postgresql://root@localhost:26257/postgres?sslmode=disable"
	}
	if config.RedisHost == "" {
		config.RedisHost = "localhost"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
	if os.Getenv("REDIS_DB") == "" {
		config.RedisDb = 0
	} else if number, err := strconv.Atoi(os.Getenv("REDIS_DB")); err != nil {
		clog.Error(ctx, "REDIS_DB enviroment variable must be a valid integer value", err)
	} else {
		config.RedisDb = number
	}
	return &config
}
