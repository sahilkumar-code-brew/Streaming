package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const pingTimeout = 2 * time.Second

func NewRedis(ctx context.Context, addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis at %s: %w", addr, err)
	}
	return client, nil
}
