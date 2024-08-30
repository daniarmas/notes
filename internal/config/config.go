package config

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	DatabaseUrl   string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDb       int
}

func LoadConfig() *Configuration {
	config := Configuration{
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
	if config.DatabaseUrl == "" {
		config.DatabaseUrl = "postgresql://root@localhost:26257/notes_database?sslmode=disable"
	}
	if config.RedisHost == "" {
		config.RedisHost = "localhost"
	}
	if config.RedisPort == "" {
		config.RedisPort = "6379"
	}
	if os.Getenv("REDIS_DB") == "" {
		config.RedisDb = 0
	} else {
		number, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			log.Fatalf("REDIS_DB enviroment variable must be a valid integer value")
		}
		config.RedisDb = number
	}
	return &config
}
