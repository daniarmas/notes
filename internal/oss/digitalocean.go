package oss

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
)

type oss struct {
	s3Client *s3.S3
}

// Implement this method
func NewOssDigitalOcean(cfg config.Configuration) *oss {
	// AWS S3 client configuration
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.ObjectStorageServiceAccessKey, cfg.ObjectStorageServiceSecretKey, ""), // Specifies your credentials.
		Endpoint:         aws.String(cfg.ObjectStorageServiceEndpoint),                                                               // Find your endpoint in the control panel, under Settings. Prepend "https://".
		S3ForcePathStyle: aws.Bool(false),                                                                                            // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
		Region:           aws.String(cfg.ObjectStorageServiceRegion),                                                                 // Must be "us-east-1" when creating new Spaces. Otherwise, use the region in your endpoint, such as "nyc3".
	}
	// The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		clog.Error(context.Background(), "error creating aws s3 new session", err)
	}
	return &oss{
		s3Client: s3.New(newSession),
	}
}

// Implement this method
func (o *oss) HealthCheck() error {
	return nil
}
