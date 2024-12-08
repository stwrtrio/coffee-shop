package main

import (
	"context"
	"log"

	"github.com/stwrtrio/coffee-shop/pkg/database"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/redis"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
)

func main() {
	// Load Configuration
	config, err := utils.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to Database
	db, err := database.Connect(&config.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Connected to MySQL Database")

	// Initialize Kafka
	kafkaClient, err := kafka.InitKafka(&config.Kafka)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	log.Println("Connected to Kafka")
	defer kafkaClient.Close()

	// Connect to Redis
	redisClient, err := redis.ConnectRedis(&config.Redis)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Connected to Redis")
	defer redisClient.Close()

	// Setup graceful shutdown logic
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register cleanup functions for graceful shutdown
	utils.GracefulShutdown(ctx, cancel,
		func() error {
			log.Println("Cleaning up Redis...")
			return redisClient.Close()
		},
		func() error {
			log.Println("Cleaning up database...")
			if sqlDB, err := db.DB(); err == nil && sqlDB != nil {
				return sqlDB.Close()
			}
			return nil
		},
		func() error {
			log.Println("Cleaning up Kafka...")
			kafkaClient.Close()
			return nil
		},
	)

	log.Println("Application stopped.")
}
