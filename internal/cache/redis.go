package cache

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

// validateRedisConfig validates the Redis configuration
func validateRedisConfig(host, port string, db int) error {
	if host == "" {
		return fmt.Errorf("host is required")
	} else {
		if host != "localhost" && net.ParseIP(host) == nil {
			return fmt.Errorf("host must be a valid IP address")
		}
	}
	if port == "" {
		return fmt.Errorf("port is required")
	} else {
		if _, err := net.LookupPort("tcp", port); err != nil {
			return fmt.Errorf("port must be a valid port number")
		}
	}
	if db < 0 {
		return fmt.Errorf("db must be a non-negative integer")
	}
	return nil
}

func OpenRedis(ctx context.Context, host, port, password string, db int) (*redis.Client, error) {
	// Validate the Redis configuration
	if err := validateRedisConfig(host, port, db); err != nil {
		return nil, err
	}
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, port),
		Password: password,
		DB:       db,
	})
	// Check if the connection is successful
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
