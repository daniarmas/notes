package config

import (
	"context"
	"os"
	"strconv"

	"github.com/daniarmas/notes/internal/clog"
)

type Configuration struct {
	DatabaseUrl   string
	JwtSecret     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDb       int
}

func LoadServerConfig() *Configuration {
	ctx := context.Background()
	config := Configuration{
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		JwtSecret:     os.Getenv("JWT_SECRET"),
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
