package database

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

type RedisConfig struct {
    Host string
    Port string
}

func NewRedisConnection(cfg RedisConfig) (*redis.Client, error) {
    addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return client, nil
}