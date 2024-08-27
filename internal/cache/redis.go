package cache

import (
	"log"
	"net"
	"time"

	"github.com/daniarmas/notes/internal/config"
	"github.com/redis/go-redis/v9"
)

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
