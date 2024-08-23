package config

import (
	"os"
)

type Configuration struct {
	AppRestPort string
	DatabaseUrl string
}

func LoadConfig() *Configuration {
	config := Configuration{
		AppRestPort: os.Getenv("APP_REST_PORT"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}
	if config.AppRestPort == "" {
		config.AppRestPort = "8080"
	}
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	}
	return &config
}
