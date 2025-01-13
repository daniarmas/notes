package oss

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"errors"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type oss struct {
	client *minio.Client
	cfg    *config.Configuration
}

func NewDigitalOceanWithMinio(cfg *config.Configuration) ObjectStorageService {
	// Initialize minio client object.
	minioClient, err := minio.New(cfg.ObjectStorageServiceEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.ObjectStorageServiceAccessKey, cfg.ObjectStorageServiceSecretKey, ""),
		Secure: true,
	})
	if err != nil {
		clog.Error(context.Background(), "error creating minio client", err)
	}
	return &oss{
		client: minioClient,
		cfg:    cfg,
	}
}

func (o *oss) GetPresignedUrl(objectName string) (string, error) {
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	presignedURL, err := o.client.PresignedPutObject(context.Background(), o.cfg.ObjectStorageServiceBucket, objectName, expiry)
	if err != nil {
		clog.Error(context.Background(), "error generating presigned URL", err)
		return "", err
	}
	return presignedURL.String(), err
}

func (o *oss) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bucketExists, err := o.client.BucketExists(ctx, o.cfg.ObjectStorageServiceBucket)
	if err != nil {
		clog.Info(context.Background(), "Connection error to Object Storage server", err)
		return err
	} else if bucketExists {
		clog.Info(context.Background(), "Connection sucessfull to Object Storage server", err)
		return nil
	}
	return nil
}

func (o *oss) ObjectExists(objectName string) error {
	_, err := o.client.StatObject(context.Background(), o.cfg.ObjectStorageServiceBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return errors.New("object not found")
		} else {
			return err
		}
	}
	return nil
}

// GetObject download an object from the object storage service and return a file path
func (i *oss) GetObject(ctx context.Context, objectName string) (string, error) {
	// Download the object from the object storage service
	object, err := i.client.GetObject(context.Background(), i.cfg.ObjectStorageServiceBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		clog.Error(ctx, "error getting object", err)
		return "", err
	}
	defer object.Close()

	// Attempt to read from the object to trigger any errors
	if _, err = object.Stat(); err != nil {
		switch minio.ToErrorResponse(err).Code {
		case "NoSuchKey":
			return "", errors.New("object not found")
		default:
			return "", err
		}
	}

	// Create a local file to store the object
	baseName := filepath.Base(objectName)
	path := fmt.Sprintf("/tmp/%s", baseName)

	localFile, err := os.Create(path)
	if err != nil {
		clog.Error(ctx, "error creating local file", err)
		return "", err
	}
	defer localFile.Close()

	// Copy the object to the local file
	if _, err = io.Copy(localFile, object); err != nil {
		clog.Error(ctx, "error copying object to local file", err)
		return "", err
	}
	return path, nil
}
