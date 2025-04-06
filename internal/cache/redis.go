package cache

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/daniarmas/clogg"
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
		msg := fmt.Sprintf("could not connect to redis server at %s: %v\n", address, err)
		clogg.Error(ctx, msg)
	}
	defer conn.Close()

	msg := fmt.Sprintf("connected to redis server at %s", address)
	clogg.Info(ctx, msg)

	return rdb
}

// Close the Redis connection gracefully
func CloseRedis(ctx context.Context, client *redis.Client) {
	err := client.Close()
	if err != nil {
		msg := fmt.Sprintf("error closing Redis connection")
		clogg.Error(ctx, msg, clogg.String("error", err.Error()))
	}
	clogg.Info(ctx, "redis connection closed")
}
