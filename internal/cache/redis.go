package cache

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/redis/go-redis/v9"
)

// Open the redis connection
func OpenRedis(ctx context.Context, cfg *config.Configuration) *redis.Client {
	address := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		msg := fmt.Sprintf("Could not connect to Redis server at %s: %v\n", address, err)
		clog.Error(ctx, msg, err)
	}
	defer conn.Close()

	msg := fmt.Sprintf("Connected to Redis server at %s", address)
	clog.Info(ctx, msg, nil)

	return rdb
}

// Close the Redis connection gracefully
func CloseRedis(ctx context.Context, client *redis.Client) {
	err := client.Close()
	if err != nil {
		msg := fmt.Sprintf("Error closing Redis connection: %v", err)
		clog.Error(ctx, msg, err)
	}
	clog.Info(ctx, "Redis connection closed", nil)
}
