package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, pw string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
		Protocol: 2, // Strictly enforce RESP2 for Redis 5.x
	})

	ctx := context.Background()
	// Test connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Connect failed:", err)
	}
	fmt.Println("Connected:", pong)

	return rdb
}
