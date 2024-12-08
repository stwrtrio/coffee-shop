package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"golang.org/x/net/context"
)

// ConnectRedis initializes and returns a Redis client.
func ConnectRedis(config *utils.RedisConfig) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password,
		DB:       0, // default DB
	})

	// Check Redis connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return nil, err
	}

	return client, nil
}
