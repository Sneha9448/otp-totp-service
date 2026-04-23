package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	DB  *redis.Client
	Ctx = context.Background()
)

// InitRedis initializes the Redis connection
func InitRedis() error {
	DB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := DB.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
	return nil
}
