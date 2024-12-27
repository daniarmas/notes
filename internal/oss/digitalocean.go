package oss

type oss struct {
}

// Implement this method
func NewOssDigitalOcean() *oss {
	return &oss{}
}

// Implement this method
func (o *oss) HealthCheck() error {
	return nil
}
