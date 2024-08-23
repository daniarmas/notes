package config

import (
	"os"
)

type Configuration struct {
	AppRestHost string
	AppRestPort string
	DatabaseUrl string
}

func LoadConfig() *Configuration {
	config := Configuration{
		AppRestHost: os.Getenv("APP_REST_HOST"),
		AppRestPort: os.Getenv("APP_REST_PORT"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}
	if config.AppRestHost == "" {
		config.AppRestHost = "localhost"
	}
	if config.AppRestPort == "" {
		config.AppRestPort = "8080"
	}
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	}
	return &config
}
