package config

import (
	"os"
)

type Configuration struct {
	DatabaseUrl   string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDb       string
}

func LoadConfig() *Configuration {
	config := Configuration{
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDb:       os.Getenv("REDIS_DB"),
	}
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "postgresql://root@localhost:26257/defaultdb?sslmode=disable"
	}
	if config.RedisHost == "" {
		config.RedisHost = "localhost"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
	if config.RedisDb == "" {
		config.RedisHost = "0"
	}
	return &config
}
