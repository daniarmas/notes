package cache

import (
	"log"
	"net"
	"time"

	"github.com/daniarmas/notes/internal/config"
	"github.com/redis/go-redis/v9"
)

// Open the redis connection
func OpenRedis(cfg *config.Configuration) *redis.Client {
	address := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Fatalf("Could not connect to Redis server at %s: %v\n", address, err)
	}
	defer conn.Close()

	log.Printf("Connected to Redis server at %s\n", address)
	return rdb
}

// Close the Redis connection gracefully
func CloseRedis(client *redis.Client) {
	err := client.Close()
	if err != nil {
		log.Fatalf("Error closing Redis connection: %v", err)
	}
	log.Println("Redis connection closed")
}
