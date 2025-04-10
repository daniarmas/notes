package cache

import (
	"context"
	"net"

	"github.com/daniarmas/notes/internal/config"
	"github.com/redis/go-redis/v9"
)

// Open the redis connection
func OpenRedis(ctx context.Context, cfg *config.Configuration) (*redis.Client, error) {
	address := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
