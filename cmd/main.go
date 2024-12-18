package main

import (
	"context"
	"log"

	"github.com/stwrtrio/coffee-shop/internal/delivery/http/handlers"
	"github.com/stwrtrio/coffee-shop/internal/domain/repositories"
	"github.com/stwrtrio/coffee-shop/internal/domain/services"
	"github.com/stwrtrio/coffee-shop/internal/routes"
	"github.com/stwrtrio/coffee-shop/pkg/database"
	"github.com/stwrtrio/coffee-shop/pkg/email"
	"github.com/stwrtrio/coffee-shop/pkg/kafka"
	"github.com/stwrtrio/coffee-shop/pkg/middlewares"
	"github.com/stwrtrio/coffee-shop/pkg/redis"
	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"github.com/stwrtrio/coffee-shop/scheduler"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

	// Repository
	customerRepo := repositories.NewUserRepository(db)
	emailConfirmationRepo := repositories.NewEmailConfirmationRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	menuRepo := repositories.NewMenuRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Services
	customerService := services.NewUserService(config, customerRepo, kafkaClient)
	emailConfirmationService := services.NewEmailConfirmationRepository(config, kafkaClient, emailConfirmationRepo)
	menuService := services.NewMenuService(config, redisClient, menuRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	orderService := services.NewOrderService(orderRepo, menuRepo)

	// Start Kafka Consumer
	go emailConfirmationService.ConsumeEmailConfirmation()

	// Initialize email service
	emailService := email.NewEmailService(
		config.Email.SMTPHost,
		config.Email.SMTPPort,
		config.Email.SenderEmail,
		config.Email.SenderPasswd,
	)

	// Start the scheduler
	notificationScheduler := scheduler.NewScheduler(notificationRepo, orderRepo, emailService)
	go notificationScheduler.StartScheduler()

	// Handler
	customerHandler := handlers.NewUserHandler(customerService)
	menuHandler := handlers.NewMenuHandler(menuService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Set up Echo
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	// Customer Routes
	routes.RegisterUserRoutes(e, config, customerHandler)

	// Menu Routes
	routes.RegisterMenuRoutes(e, config, menuHandler, categoryHandler)

	// Order Routes
	routes.RegisterOrderRoutes(e, config, orderHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))

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
