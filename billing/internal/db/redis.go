package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(addr, pass string, db int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Connection to Redis error: %v", err)
	}

	return redisClient
}
